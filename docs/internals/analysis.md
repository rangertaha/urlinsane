---
title: 18 · Analysis
parent: Part III · Inside the engine
nav_order: 6
---

# Analysis
{: .no_toc }

- TOC
{:toc}

Expansion produces a graph. Analysis turns it into conclusions. The two are
separated because they answer different questions and must not be allowed to
contaminate each other.

## Belief is not risk

There are two numbers in the system and confusing them would be a serious bug,
so the code makes it structurally impossible:

```go
// BeliefModel scores a node for execution control only: frontier ordering,
// pruning and operator gating. It never contributes to a reported number, which
// is why analyzers cannot see belief at all.
```

**Belief** is the engine's own estimate of whether a node is worth spending
network calls on. It decides what gets scanned next when a budget binds. It is
engine-internal and analyzers have no access to it.

**Risk** is an analyzer's conclusion about the finished graph, computed from
observable facts — does it resolve, does it accept mail, how many edits from the
target — and it is what the report shows.

If belief leaked into risk, the tool would report a domain as dangerous partly
*because it decided to look at it*. That is circular, and it is the kind of
circularity that survives review because both numbers look plausible.

## The belief model

```go
type BeliefModel interface {
    Initial() (float64, State)
    Step(parent State, rel string, v View) (float64, State)
}
```

Belief is a pure function of the parent's belief-state and the node's props as
of the current barrier. Nothing reaches sideways into the graph, so the same run
recomputes the same values in the same order.

`State` is a model's latent state, carried from parent to child and never
inspected by the graph. It exists because the design specifies a hidden Markov
model, and forward filtering in an HMM propagates a *distribution over latent
states* — a vector, not a scalar. An earlier version of the interface passed the
parent's scalar belief, which forced a model to reconstruct a plausible
distribution from one number and collapse it again on the way out:

> That round trip is exact for a two-state model and lossy for three or more, so
> the interface silently answered the open question about state cardinality as
> "two" — not by anyone's decision, but as a consequence of a type signature.
> Numbers from a larger model would have looked entirely reasonable and been
> wrong.

The shipped default is deliberately trivial:

```go
type uniformModel struct{}

func (uniformModel) Initial() (float64, State)                 { return 1, nil }
func (uniformModel) Step(State, string, View) (float64, State) { return 1, nil }
```

It reduces expansion to breadth-first and unranked, which means the engine ships
and runs correctly before any model exists, and a model that turns out poor can
be dropped without invalidating a single result. This is why `--verbose` shows
`belief=1.000` on every row today.

Belief is recomputed for **every** node at every barrier, not once at creation —
a node created this round has only the props its creating operator set, so
computing belief once would leave it a bare prior forever.

## Severity

```go
const (
    SeverityInfo Severity = iota + 1
    SeverityLow
    SeverityMedium
    SeverityHigh
    SeverityCritical
)
```

> It is deliberately not a bare int: `--fail-on high` has to mean the same thing
> in every release, and an integer scale drifts silently the first time someone
> inserts a level.

## The analyzers

An analyzer runs once, over the finished graph, and returns findings.

### `scoring` — is this variant dangerous?

Ranks a live variant by how plausible the confusion is:

```go
if len(a.Outgoing(n.ID, "RESOLVES_TO")) > 0 {
    score += 2
    signals = append(signals, "resolves")
}
if len(a.Outgoing(n.ID, "MX")) > 0 {
    // Mail is the signal that separates a parked name from one built
    // to receive credentials.
    score += 3
    signals = append(signals, "accepts mail")
}
if d := analyze.EditDistance(a, n.ID); d > 0 && d <= 1 {
    score += 2
    signals = append(signals, "one edit from the target")
}
```

Which is why the report reads:

```
HIGH    live-variant  eaxmple.com is live: [resolves accepts mail]
MEDIUM  live-variant  examlpe.com is live: [resolves]
```

MX weighs more than an A record because a name that accepts mail is built to
receive something, and the something is usually credentials or invoices.

Note the two guards at the top: it only looks at nodes with an inbound
`VARIANT_OF`, and only at ones whose existence is `live`. That first guard is why
the applier's self-variant rejection matters — a seed that became a variant of
itself would be scored as a live typosquat of itself.

### `campaign` — is this one actor?

```go
// Package campaign finds variants that cluster: shared addresses, shared
// nameservers, shared registrant. One squatted lookalike is an incident; forty
// on one nameserver is a campaign, and only the graph can see the difference.
```

It clusters on `RESOLVES_TO`, `NS` and `REGISTERED_BY` by default, with a floor
of two — "two is the honest floor: one variant on an address is not a campaign".

This is the analyzer that justifies the graph model. In a flat list of domains,
"these six variants share a nameserver" is not expressible; here it is an
incoming-edge count on a hub node.

### `dep-confusion` — is absence the finding?

```go
// Package depconfusion flags the supply-chain case: a dependency that exists
// on one registry and nowhere else. Absence is the finding, which is why it can
// only be written against three-state existence — under a two-state model an
// unreachable registry and a free name are the same answer.
```

A package with no `EXISTS_ON` edge and existence `absent` is **critical**: it is
not published on its registry, so a public package of that name would win
resolution.

The same package with existence `unknown` produces a different finding at
`info` severity:

```go
// Reporting this as a gap would be a guess. Saying so is the point
// of having a third value.
```

Two analyzers, same input shape, opposite conclusions — decided entirely by
whether the registry answered. This is the clearest payoff of the three-state
model anywhere in the codebase.

## Findings

```go
type Finding struct {
    Kind     string
    Severity Severity
    Nodes    []NodeID
    …
}
```

Findings attach to nodes, which is what lets the report show a `RISK` column per
row *and* an explanatory `FINDINGS` block, both from one pass. `--filter
risk>SEV` selects on the maximum severity of a node's findings, and `--fail-on`
compares against the maximum across the whole run.

A finding may also reference a **declined** candidate via `LedgerRef`, so an
analyzer can say something about a name the engine chose not to admit.

---

Next: **[Content addressing](../addressing/)**.
