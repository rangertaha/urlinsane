---
title: 16 · The plan and the scheduler
parent: Part III · Inside the engine
nav_order: 4
---

# The plan and the scheduler
{: .no_toc }

- TOC
{:toc}

There are two layers here, deliberately, and conflating them is what made earlier
revisions of the design inconsistent:

- **The compiled plan** — a static, inspectable artifact naming every operator
  that *may* run, its trigger, and the limits. This is what `--explain` prints.
- **Pattern dispatch in rounds** — the runtime mechanism. An operator runs when
  its pattern matches a node, which may be long after that node was created.

The plan answers "what will this do?". Dispatch answers "when?".

## The compiled plan

```console
$ urlinsane typo --explain example.com
plan 7437e47542eddfee
seed domain example.com
scope every nameable node in the seed closure
limits depth=3 rounds=64

type flow
  0  domain, email, ip, package, repo, username  (cycle)
  1  platform
  1  registrant
  1  tld

operators
  aci            on nameable           where in-seed-closure
  …
  dns-a          on domain             [dns]
  ptr            on ip                 [dns]
  whois          on domain             [whois]
```

**The plan hash** is a function of the compiled plan, not of the results. The
same flags against the same build produce the same hash, which is what makes it
usable as a "did anything about how I scan change?" check.

**The type flow** is built from operators' declared `Effects`: if some operator
binds to `domain` and emits `RESOLVES_TO` to an `ip`, then `domain → ip` is an
edge in the type-flow graph. The numbers on the left are layers of the SCC
condensation, and `(cycle)` marks a strongly connected component with more than
one member — here, the whole nameable cluster, because `domain → ip → domain` is
a real cycle in the data.

This is the only place a topological sort appears in the whole engine, and it
decides nothing: it is presentation, so `--explain` renders in a readable layered
form instead of an arbitrary order.

**Operators missing from the plan are the interesting part.** An operator whose
reference data is unavailable is omitted at compile time rather than failing at
run time. So the plan is also the answer to "will this run actually check package
registries?" — if `pkg` is not listed, it will not, and you know before you start.

## Rounds

```go
// Run expands until a round produces no new eligible work, or a limit binds.
func (s *Scheduler) Run(ctx context.Context) error {
    s.barrier() // barrier 0: seeding
    for {
        if ctx.Err() != nil { return nil }
        if s.round >= s.lim.MaxRounds { … return nil }
        s.round++

        work := s.eligible()
        if len(work) == 0 { return nil }
        results := s.dispatch(ctx, work)
        s.applyAll(work, results)
        s.barrier()
    }
}
```

Four phases per round, and the separation between them is the whole concurrency
story:

1. **`eligible()`** — collect every `(node, operator)` pair whose trigger matches
   and which has not already run against this exact read-set.
2. **`dispatch()`** — run them **concurrently**, up to `--workers`. Nothing is
   applied here; results are collected into a slice indexed by work order.
3. **`applyAll()`** — apply the results **in work order, never completion
   order**, through the single-writer applier.
4. **`barrier()`** — recompute belief, enforce budgets. Everything irreversible
   happens here.

Concurrent execution with deterministic application. The network decides when
answers arrive; it does not decide what the graph looks like afterwards.

## Determinism

`eligible()` sorts its work list by `(depth, type, key, operator)`:

```go
// The order is deterministic — (depth, type, key, operator) — which is what
// makes the whole round reproducible regardless of the order operators happen
// to finish in.
```

Combined with apply-in-work-order, this is what makes two runs of the same scan
produce the same graph, and therefore the same CID. Note what is *not* relied on:
no map is ranged over anywhere on the path to a block, and nothing reads a clock
or a random source.

## The seen-set and re-dispatch

A pair is keyed by `(node, operator, read-set digest)`:

```go
type seenKey struct {
    node   NodeID
    op     string
    digest [32]byte
}
```

After a pair runs, that exact read-set is closed **whatever the outcome**:

```go
// Retry already happened inside the round, bounded by the attempt count —
// leaving a failed pair eligible as well would multiply attempts by rounds and
// hammer a service that is already unwell. If the graph later changes something
// the operator declared it reads, the digest changes and the pair becomes
// eligible again on its own.
```

That is the re-dispatch mechanism in one paragraph. An operator does not
subscribe to anything and nothing notifies it; its eligibility is recomputed from
the data, and the data has a hash.

A second bound, `Revisions` (default 3), caps how many times one pair may run
across the whole scan, so a pair whose read-set keeps oscillating cannot spin.
Only an actual execution counts against it — a pair that was gated off never ran,
so there is nothing to revise.

## Retry lives inside the round

```go
for attempt := 0; attempt < s.lim.Attempts; attempt++ {
    if ctx.Err() != nil { o = Timeout(ctx.Err()); break }
    s.rate.Acquire(ctx, c.op.Resource())
    v := s.g.viewFor(c.op.Trigger(), c.id)
    d, o = s.exec(ctx, c.op, v)
    if !o.retriable() { break }
}
```

Retry is *inside* the dispatch call, not across rounds. The barrier therefore
waits for it, and the outcome set is fixed before belief is computed. Retrying
across rounds instead would multiply attempts by rounds against a service that
is, by hypothesis, already struggling.

Only `Failed` and `Timeout` are retriable. `Empty` is an answer.

## Timeouts belong to the scheduler

```go
// The deadline is applied here rather than left to each operator so that every
// operator is bounded by the same rule, including one that forgot to bound
// itself. An operator that ignores ctx still runs to completion — Go cannot
// preempt it — but its result arrives against an expired context and the round
// is not blocked from concluding.
```

`--timeout` bounds one `Exec` call. An operator holding its own private deadline
is one the barrier cannot reason about.

The same applies to interruption: `ctx` carries both the round deadline and the
Ctrl-C signal, and an operator making network calls must pass it down. Without
that, the scheduler can only cancel *between* attempts, so Ctrl-C waits for an
in-flight WHOIS rather than stopping at the round boundary.

## Caching

The cache key is `(operator, node, read-set digest)`. Two rules:

- A `Failed` or `Timeout` outcome is **never cached**. Caching the absence of an
  answer would turn one bad moment into a permanent one.
- A cache hit skips dispatch entirely and is counted, so `CacheHits` against
  `Dispatched` tells you how much of a scan was recomputation.

## The barrier

```go
// barrier finalizes everything irreversible. Belief is recomputed for every
// node from its props as of now — a node created this round has only the props
// its creating operator set, so computing belief once at creation would leave
// it a bare prior forever.
func (s *Scheduler) barrier() {
    s.g.recomputeBelief()
    s.g.enforceBudgets(Provenance{Operator: "engine", Round: s.round})
}
```

**Belief recomputation** is the part that must be here: it needs the round's
props to be in place, and it walks nodes in depth order so a parent's belief is
current before its children read it.

{: .todo }
> `enforceBudgets` is the hook for making the barrier own every limit check, and
> it is currently an empty function. Budgets are enforced at admission time
> instead — `admit` checks `overBudget` and writes a ledger row — which works,
> but means the decision is taken per-candidate rather than over the round's
> whole admission set. See [Limits](../limits/).

Interruption stops expansion at the **end** of the current round for exactly this
reason: the barrier still runs, so parents, belief and the ledger are finalised
rather than left half-computed.

---

Next: **[Limits and termination](../limits/)**.
