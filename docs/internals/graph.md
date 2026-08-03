---
title: 13 · The graph
parent: Part III · Inside the engine
nav_order: 1
---

# The graph
{: .no_toc }

- TOC
{:toc}

## The inversion

The old engine was a pipeline. A target produced a list of domains; the list was
pushed through a fixed sequence of collectors; the result was printed. A DAG did
exist, but it ordered *plugins*, with declared dependencies between them.

That design failed in three specific ways, and they are worth naming because
they are the reason for everything in this part:

1. **Plugin order was load-bearing.** Adding a collector meant working out where
   in the sequence it belonged and what it could assume had already run.
2. **The cache was unsound.** If results are keyed by plugin and the plugin's
   inputs depend on an implicit ordering, the key does not describe the inputs.
3. **It could only produce one shape of result.** A list of domains cannot
   express "this variant and that variant resolve to the same address, which
   reverse-resolves to a name registered by the same registrant".

What replaced it:

> **The graph is the data, not the execution plan.**

A target expands into a typed property graph. Operators match patterns in that
graph, do work, and expand it further. Expansion stops when a round produces no
new eligible work or a limit binds. Analyzers then run over the whole graph, and
the report is a projection of it.

Typosquatting is the first workload expressed this way, not the shape of the
engine.

## Nodes and edges

```go
type Node struct {
    ID    NodeID   // hash(Type, Canonical(Key)) — stable for the node's life
    Type  *NodeType
    Key   string   // canonicalized at admission
    Props Props
}

type Edge struct {
    ID       EdgeID // hash(From, Rel, To)
    From, To NodeID
    Rel      *Rel
    Props    Props
}
```

Operators never construct these. They emit `NodeRef`/`EdgeRef`/`PropSet` — which
carry *raw* keys — and the applier canonicalizes, assigns identity and applies
merge policy. Letting a plugin mint an identity is how convergence quietly
breaks.

Nodes and edges carry **no provenance, status or findings**. Those live in side
tables keyed by `NodeID`, and the reason is the next section.

## Three layers of identity

| Layer | Contents | Purpose |
|---|---|---|
| `NodeID` | `hash(Type, Canonical(Key))` | in-memory identity — scheduler, seen-set, cache, edges and side tables all key on it |
| Addressed form | node `(Type, Key, Props)`, edge `(From, Rel, To, Props)` — **only** | gets the CID; this is what diffs between runs |
| Side tables | provenance, per-pair status, belief, findings, truncation ledger | keyed by `NodeID`, persisted beside the block, never inside it |

Keeping provenance out of the addressed form is what makes two identical scans
produce identical CIDs. With run ids and timestamps inside, every node would
differ on every run and cross-run diffing — a primary use case — would be
impossible.

`NodeID` is stable for the node's whole life; the CID changes as props
accumulate. Both are needed and they are not the same thing.

## Canonicalization comes first

Canonicalization is a required registry field, and it runs **before the
admission decision** — every candidate is canonicalized as it leaves the
applier's front door, whether or not it is ultimately admitted.

Two consequences:

- Without it, convergence fails silently the first time one operator emits
  `Example.com` and another `example.com`. They would be two nodes, each with
  half the evidence.
- If it ran *after* admission, the truncation ledger's denylist would compare
  raw keys, and a candidate recorded as declined under `example.com` would slip
  back in as `Example.com`.

A key that fails to canonicalize is refused with a recorded status and never
becomes a node.

## Props are an ordered field list

`map[string]any` cannot be deterministically encoded — map iteration order is
unspecified and `any` has no stable IPLD representation. Since a content address
over the props is load-bearing, that is disqualifying.

Each type registers an **ordered** field list; a field's position is its stable
index, and values are stored positionally. Values come from a closed kind set:
`String, Int, Float, Bool, Bytes, Time`. Field handles are resolved from the
registry when an operator registers, so an unknown field name fails at
registration rather than returning a boolean at every access.

What ordering buys:

- **Deterministic encoding with no sort step** — order is a property of the
  type, so identical values encode identically by construction.
- **Meaningful presentation order** — report columns follow declaration order,
  not `created, live, punycode, rank` alphabetical noise.
- **Compact encoding** — values are positional, so field-name strings are not
  repeated per node. At the 10⁵-node scale a wide scan reaches, that decides
  whether the graph fits in memory.
- **Cheap access** — an index, not a map hash.

The cost is a rule that must be enforced: **fields are append-only.** Reordering
or deleting one changes the meaning of every content address already in the
store — corrupting diffs rather than failing loudly. Removal is by tombstone;
the slot stays, marked deprecated. Each type carries a schema version, so a
decoder meeting a block from a newer binary *detects* trailing fields it does not
know.

It does not preserve them. `Props` is sized from the registry schema, so there
is nowhere to hold a slot the running binary does not declare. The implemented
behaviour is to **refuse**: the store fails the rehydrate CID check rather than
writing back a truncated node. Silently dropping unknown values and changing the
CID on re-save would be data loss disguised as a successful write.

## Conflicting assertions

Two operators may assert the same field on the same node — `whois` and `rdap`
both produce registration dates. Under concurrent dispatch, last-write-wins would
be decided by network timing, which is nondeterministic.

Every field therefore declares a merge policy:

```go
{Name: "created", Kind: Time, Merge: Precedence("rdap", "whois")}
```

All assertions are retained in the provenance side table; the **materialized**
value is selected by declared precedence, falling back to lowest operator id.
Deterministic regardless of arrival order — and disagreement between two sources
is preserved as signal for analyzers rather than silently resolved by whoever
answered first:

```go
// Assertion is one operator's claim about a field, retained whether or not it
// won the merge. Disagreement between two sources is signal, not noise.
type Assertion struct {
    Field string
    Value Value
    By    Provenance
    Won   bool
}
```

## Deltas are additive

```go
type Delta struct {
    Nodes []NodeRef
    Edges []EdgeRef
    Props []PropSet
}
```

Operators never mutate the graph; they return a `Delta` and the applier is the
single writer. Deltas are **additive only** — nothing removes a node, edge or
prop — which is what makes the graph monotonic within a run and a delta safely
replayable.

Monotonicity is not an aesthetic choice. It is what lets the scheduler dispatch
concurrently and apply in a deterministic order afterwards, and what lets a saved
graph be rehydrated by replaying blocks through the same applier.

## Cycles are correct

The entity graph is **cyclic by nature**. `domain → ip → PTR → domain` is correct
data, not a bug, and the type-flow graph is cyclic for the same reason. Neither
can be topologically sorted, and termination never relies on acyclicity — see
[Limits](../limits/).

So nothing is topologically sorted to decide execution order. The `dag` package
survives for **plan presentation** only: condensing the type-flow graph into
strongly connected components and ordering that condensation, so `--explain`
renders in a readable layered form. That is Tarjan SCC plus a topological sort of
the condensation, and it decides nothing about what runs.

---

Next: **[Types and relations](../types/)**.
