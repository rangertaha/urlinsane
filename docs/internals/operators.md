---
title: 15 · Operators
parent: Part III · Inside the engine
nav_order: 3
---

# Operators
{: .no_toc }

- TOC
{:toc}

## The interface

```go
type Operator interface {
    Id() string
    Version() int
    Trigger() Trigger
    Emits() Effects
    // Resource names the rate-limit class this operator's calls belong to.
    Resource() string
    Exec(ctx context.Context, v View) (Delta, Outcome)
}
```

Five declarations and one method. An operator says **when** it runs (`Trigger`),
**what it may produce** (`Emits`), **what shared thing it consumes**
(`Resource`), and then does the work.

Nothing in there says what runs before it, or after it. That is the whole point.

## Binding is by data, never by producer

```go
// Selector chooses which nodes an operator binds to, by type or by capability.
// Binding by capability is what lets one omission algorithm cover every
// Nameable type instead of being registered once per type.
type Selector struct {
    Types []string
    Caps  []Capability
}
```

```go
// Condition is an extra requirement on the matched node. Conditions are data
// conditions, never producer dependencies: "there is an IP", not "the ip
// operator has run".
type Condition interface { … }
```

That comment is the design in one line. The available conditions:

| Condition | Requires |
|---|---|
| `HasProp(field)` | the field is set on the matched node |
| `HasEdge(rel)` | at least one outgoing edge of that relation |
| `InClosure()` | seed-closure membership |
| `BeliefAbove(t)` | the execution model's belief clears a threshold |

`ptr` binds to any `ip`. It does not bind to "whatever `dns-a` produced". So it
runs on addresses discovered by `dns-a`, by a future `rdap` operator, by a
manifest parser, or by an operator nobody has written yet — and none of those
authors has to know `ptr` exists.

This is what makes the system extensible in the way the old pipeline was not.
Adding an operator is adding a file; there is no order to insert it into and no
downstream consumer to notify.

`BeliefAbove` is the odd one out: it is evaluated only at a barrier, never during
delta-driven re-dispatch, so it takes no part in the read-set digest below.

## Reads, and why they are checked

```go
// Reads declares the props and relations an operator consumes. It does double
// duty: it scopes the View, and it is the input to the read-set digest that
// decides re-dispatch and cache validity.
type Reads struct {
    Fields []string
    Rels   []string
}
```

An operator sees a `View` scoped to what it declared it reads — not the whole
graph. The same declaration is hashed into a **read-set digest**, and the pair
`(node, operator, digest)` is what the seen-set and the cache key on.

The consequence is the cache soundness the old design could not have: a cached
result is valid exactly as long as the data the operator actually read has not
changed. When it changes, the digest changes, and the pair becomes eligible
again on its own.

There is a trap here that the code closes explicitly:

```go
// effectiveReads merges the operator's declared reads with whatever its
// conditions inspect. A condition that reads a field the operator forgot to
// declare would otherwise leave that field out of the digest, and the operator
// would never re-run when it changed.
```

So `Where: []Condition{HasProp("live")}` contributes `live` to the digest even
if `Reads` does not mention it.

## Emission

```go
// Effects declares everything an operator may produce. It covers relations and
// props, not just node types, so plan compilation can see a prop-only operator
// and can detect a Where nothing in the plan will ever satisfy.
type Effects struct {
    Nodes []string
    Rels  []string
    Props []string
}
```

Two things depend on this being honest:

- **The plan.** `--explain` builds the type-flow graph from declared effects. An
  operator that emits something it did not declare is invisible to the plan.
- **The variant rule.** An operator is a *variant* operator if and only if it
  declares `VARIANT_OF` in `Emits().Rels`. Not by naming convention, not by
  package — by declaration. That declaration is what subjects it to the
  seed-closure restriction and the terminal rule.

## Outcome: how it failed is the finding

```go
type Outcome struct {
    Status Status
    Err    error
}

func OK() Outcome             // learned something positive
func Empty() Outcome          // authoritatively determined absence
func Failed(e error) Outcome  // the lookup itself broke
func Timeout(e error) Outcome // nothing was learned
```

> How a lookup failed is itself the finding: NXDOMAIN proves a name is free, a
> timeout proves nothing at all, and collapsing the two would discard the signal
> a squatting scanner exists to collect.

`StatusEmpty` is an *answer*. `StatusFailed` and `StatusTimeout` are the absence
of one, and only those two are retriable — an authoritative "absent" is never
retried. There is a fifth status, `StatusSkipped`, which is deliberately **not
terminal**: a pair gated off at one barrier may run at a later one, and
recording it terminally would make the first belief gate permanent.

This is where the three-state existence a user sees comes from. It is not
computed at the end; it is an aggregation of per-pair statuses recorded as they
happen.

## The applier is the only writer

Operators return a `Delta` of `NodeRef`/`EdgeRef`/`PropSet`, all carrying **raw**
keys. The applier canonicalizes, assigns identity, applies merge policy, and
enforces the invariants:

| Rejection | Why |
|---|---|
| `RejectCanonical` | the key could not be canonicalized |
| `RejectUnknownType` | no such node type is registered |
| `RejectUnknownRel` | no such relation is registered |
| `RejectClosure` | a variant edge from outside the seed closure |
| `RejectScope` | a variant edge from a type the run's scope excludes |
| `RejectSelfVariant` | a node claimed as a variant of itself |

Two of those are worth reading closely.

**The closure check runs before the endpoint is admitted, not after.** Rejecting
the edge while admitting its target would leave the variant in the graph,
unrooted and unreachable — and the invariant exists to bound combinatorial
expansion, so a check that stops the edge but not the node does not bound
anything.

**Self-variants are refused at the applier, not at the operator.** Operators emit
raw keys and cannot see canonical form, so an algorithm can produce a string that
differs from its origin and canonicalizes straight back onto it — bit-flipping
`google.com` flips a case bit to `Google.com`, which folds back. The operator's
own dedupe compares raw strings and cannot catch that. Left in, the self-edge
makes the seed a variant of itself, and since analyzers treat any node with an
inbound `VARIANT_OF` as a variant, the target gets scored as a live typosquat of
itself and joins every campaign cluster it hosts.

## Resources and rate limiting

```go
// Limiter throttles operator calls per resource class. A single global delay is
// meaningless once one run talks to DNS, whois, npm, PyPI and GitHub at once:
// the limit protecting the strictest service would throttle everything else to
// the same crawl.
```

`Resource()` returns `dns`, `whois`, `http` or `""`. Each class gets its own
minimum interval, acquired before every call. An operator with no resource — a
pure generator, a decomposer — returns the empty string and is never throttled.

The resource declaration does a second job: the scheduler uses it to decide which
operators count as *observers*, so that a decomposer's "I parsed this" is never
read as "this exists".

```go
// Tell the graph which operators actually look something up, so Existence
// does not read a decomposer's "I parsed this" as "this exists".
```

## Anatomy of a real operator

An algorithm is the simplest case — no resource, no conditions beyond closure
membership, one relation emitted:

```
ID   TRIGGER                                    EMITS
co   on nameable, where in-seed-closure         VARIANT_OF
```

An observation operator declares a resource and binds to one type:

```
ID     TRIGGER      RESOURCE  EMITS
dns-a  on domain    dns       RESOLVES_TO
ptr    on ip        dns       PTR_TO
whois  on domain    whois     REGISTERED_BY
pkg    on package   http      EXISTS_ON
```

And a decomposer binds to a type, has no resource, and emits structural
relations:

```
ID                TRIGGER      EMITS
decompose.email   on email     LOCAL_PART, DOMAIN_OF
decompose.repo    on repo      HOSTED_ON, OWNER
```

Writing one is [chapter 20](../extending/).

---

Next: **[The plan and the scheduler](../plan/)**.
