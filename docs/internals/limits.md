---
title: 17 · Limits and termination
parent: Part III · Inside the engine
nav_order: 5
---

# Limits and termination
{: .no_toc }

- TOC
{:toc}

The graph is cyclic and the candidate space is combinatorial. Nothing about the
structure guarantees a scan ends, so termination is enforced explicitly — by six
mechanisms, none of which relies on acyclicity.

## 1. The seen-set

A pair `(NodeID, operator id, read-set digest)` runs once. Convergent and cyclic
paths therefore never redo work: coming back to `example.com` by a different
route finds every operator already closed against the read-set it has.

This is the mechanism that makes cycles safe. Not detection, not a visited flag
on the node — a hash of what each operator actually read.

## 2. Depth counts observation hops only

Default 3. Structural and variant edges cost nothing.

The alternative — counting every edge — was tried and is wrong: under whole-edge
counting `bob@example.com` spent its budget on decomposition and never reached an
IP, silently gutting the composite-target case the design exists for.

Depth is a budget and termination concept only. Synchronisation is by dispatch
round, not by depth; nodes at depth 5 and depth 1 can be worked in the same
round.

{: .note }
> `MaxDepth` of 0 means **unbounded** — the check is `if MaxDepth > 0 && depth >
> MaxDepth`. So `-d 0` is not "no hops".

## 3. Variants are terminal

A node reached by `VARIANT_OF` is never handed back to a variant operator.

Without this, `example.com` produces `exmaple.com`, which produces
`exmalpe.com`, and the tool generates the entire string space at a polynomial
rate. The rule is absolute, and it is why a variant edge can cost zero depth
without unbounding anything.

## 4. Variant roots are seed-closure members

The **seed closure** is the seed plus everything reachable from it by
*structural* edges. Only members may root variant generation.

This closes a hole the terminal rule does not cover. `ptr` emits `domain` nodes,
and `domain` is Nameable — so without a closure rule, every reverse-PTR domain on
every variant's IP address becomes a new variant root, which is exactly the
explosion the capability flag was meant to prevent, arriving by an edge the
terminal rule says nothing about.

It also has a consequence that surprises people until they see the reason:
because `MANIFEST` is an *observation* edge, manifest-derived packages fall
outside the closure. So `urlinsane typo github.com/acme/tool` does **not**
silently generate typo variants of several hundred dependencies. Varying those is
an explicit `urlinsane typo npm:<name>`, or the dependency-confusion analysis,
which needs no variants at all.

**The engine enforces this, not the plugin**, in two places:

- The scheduler will not dispatch a variant operator against a node outside the
  closure, so the work is never done.
- The applier rejects any `VARIANT_OF` edge whose source is outside it, so an
  operator that emits one anyway cannot smuggle it in.

Dispatch-side gating is the optimisation; applier-side rejection is the
invariant. An invariant whose whole purpose is preventing combinatorial explosion
must not be opt-in for plugin authors.

## 5. Budgets

Global and per-type caps on admitted nodes:

```go
type Budgets struct {
    Global  int
    PerType map[string]int
}
```

`--budget` sets the global one; `0` is unbounded. Enforcement is at **admission
time**, not at the barrier: `admit` checks `overBudget` before creating the node
and, when it binds, writes a `budget` row to the ledger and refuses the
candidate.

{: .todo }
> The design puts every limit check at the barrier, and `enforceBudgets` exists
> as the hook for that — but it is currently an empty function, and admission-time
> checking is what actually runs.
>
> `--frontier` is worse off: the flag is accepted and hashed into the plan, and
> `ReasonFrontier` exists in the ledger's reason set, but **nothing enforces the
> cap** — `declineFrontier` is only ever called with `ReasonRoundCap`. Setting
> `--frontier` today changes the plan hash and nothing else.

## 6. The round cap

Depth bounds observation hops, but a cyclic type flow (`ptr → domain → dns → ip →
ptr`) keeps generating rounds until depth bites, and an oscillating trigger
pattern might not converge at all. `--rounds` (default 64) is the backstop the
per-pair revision cap cannot provide.

Hitting it is reported like any other truncation, not swallowed.

## The truncation ledger

Belief pruning, budgets, the frontier cap, the round cap and the round deadline
all decline candidates. A declined candidate was never admitted, so it has no
`NodeID` and cannot carry a per-pair status — which means without somewhere to
record it, a truncated graph would read as a complete one.

So there is a ledger:

```go
// LedgerRow records a candidate the engine declined to admit. The ledger is
// reported like any other section — a truncated graph that reads as complete is
// a correctness bug — and it doubles as a denylist, so a later operator
// re-emitting the same candidate cannot quietly resurrect it.
type LedgerRow struct {
    Type   string
    Key    string // canonical
    Depth  int
    Belief float64
    Reason Reason
    By     Provenance
}
```

The reasons: `belief`, `budget`, `frontier`, `round-cap`, `deadline`.

The **denylist** half matters as much as the record. Without it, "pruning is
irreversible" would be false in the one case that matters — a second operator
reaching the same candidate would re-admit it, and whether a node exists would
depend on how many operators happened to find it.

This is also why canonicalization must run before the admission decision. A
ledger keyed on raw keys would let `Example.com` slip past a row recorded for
`example.com`.

## Truncation is not rejection

Two different things, deliberately kept apart:

| | Truncation (ledger) | Rejection |
|---|---|---|
| Cause | a limit bound: belief, budget, frontier, rounds, deadline | an invariant was violated |
| Effect | denies the candidate **for the rest of the run** | refuses this one assertion only |
| Examples | pruned by budget | key would not canonicalize; variant rooted outside the closure |

Putting rejections in the ledger would deny the candidate forever — and a node
refused as the source of one bad `VARIANT_OF` edge may still be perfectly
legitimate when another operator reaches it by an observation edge.

## Truncation is deterministic

It happens at a barrier, where belief is final. Candidates for the next round are
sorted by `(-belief, depth, type, key)` — belief first, so a bound budget spends
itself on the most promising frontier rather than alphabetically — and the prefix
is admitted.

The ordering is total and fixed under a pinned plan, so output stays
byte-identical across runs **even when truncated**.

## Failure is data

The thread running through all of this: an engine that cannot distinguish
"nothing there" from "could not tell" produces confident nonsense.

- Per-pair statuses are recorded terminally, including `failed` and `timeout`.
- `StatusSkipped` is explicitly **not** terminal — recording it terminally would
  make the first belief gate permanent.
- Existence is derived from those statuses, and only operators that declare a
  resource count as observers, so a decomposer's "I parsed this" is never read as
  "this exists".
- Every truncation is in the ledger and rendered in the report.
- An interrupted scan is marked `PARTIAL` in every format, including JSON.

None of that is defensive programming. Each one is a place where the tool could
have produced a clean-looking report that was wrong.

---

Next: **[Analysis](../analysis/)**.
