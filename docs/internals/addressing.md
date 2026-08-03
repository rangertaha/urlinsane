---
title: 19 · Content addressing
parent: Part III · Inside the engine
nav_order: 7
---

# Content addressing
{: .no_toc }

- TOC
{:toc}

## The claim

> Two identical scans produce the same root CID, byte for byte.

Everything in this chapter exists to make that true, because if it is true then
"what changed since last week" is a comparison of two hashes rather than a
re-scan, and a scan result becomes a thing you can store, share and verify
instead of a rendering you have to trust.

## Blocks

Results are stored as IPLD blocks: CIDv1, dag-cbor, sha2-256. A node's addressed
form is a positional list:

```go
// The form is a positional list — [type, key, [values...]] — not a map, so
// field names are not repeated per node and no key-sort step is needed to make
// encoding deterministic.
```

An edge is `[from, relation, to, [values...]]`.

Props are written as **every declared slot in declaration order, null for
unset**:

```go
// Writing all slots rather than only the set ones keeps the encoding a
// function of the schema, so a node encodes identically wherever it was built.
```

That is subtle and worth restating. If only set slots were written, the same node
would encode differently depending on which operators had run, even when the
*values* were identical — and the CID would change for reasons that have nothing
to do with content.

## What is deliberately not in the block

| In the addressed form | In side tables |
|---|---|
| node type, canonical key, props | provenance (who asserted what, in which round) |
| edge from, relation, to, props | per-pair status |
| | belief and latent state |
| | findings |
| | truncation ledger |
| | depth, seed-closure membership |

Provenance is the important exclusion. With run ids and timestamps inside the
block, every node would differ on every run and cross-run diffing — a primary use
case — would be impossible.

Depth and closure membership are excluded for a subtler reason: both are derived
from structure and are *scheduler* state. Keeping them out stops them perturbing
CIDs when, say, a scan reaches the same node by a shorter route.

## The rules that keep it deterministic

From `Store.Save`:

> The root is a pure function of the graph's content: two identical scans produce
> the same CID, byte for byte. Nothing here reads a clock or a random source, and
> every collection comes from an accessor with a defined order — **no Go map is
> ranged over anywhere on the path to a block**.

Four rules, and every one of them has a counterpart elsewhere in the engine:

1. **No clocks, no randomness** on the encoding path.
2. **No map iteration** — every collection comes from an ordered accessor.
3. **Deterministic application order** in the scheduler: work is applied in work
   order, never completion order ([chapter 16](../plan/)).
4. **Deterministic truncation** at barriers, sorted by `(-belief, depth, type,
   key)` ([chapter 17](../limits/)) — so even a truncated scan is byte-identical
   across runs.

Break any one and the property is gone.

## The scan root

A root records the seed and the content and nothing else. That is what makes the
CID stable — and it has a cost, stated in the code:

> The cost is that a root cannot answer "which was most recent" or "was this one
> interrupted", so those live out here, beside the store rather than inside it.

Hence a small index file next to the blockstore:

```go
type Entry struct {
    Type string
    // Key is the canonical seed key, matched exactly. `report bob@acme.com`
    // finds the email scan, not the scan of acme.com nested inside it.
    Key     string
    Root    string
    At      time.Time
    // Partial marks a scan that stopped early — an interrupt, a deadline, a
    // budget. It is a fact about the scan, so every re-render reports it.
    Partial bool
}
```

This is what `report --scans` prints, and it is why you address a saved scan by
target rather than by CID: `acme.com` is what you scanned and what you remember.

## Rehydration, and why it replays

Loading a saved scan does not write into the graph's internals:

> Everything is replayed through the applier rather than written into the graph's
> internals, because the applier is the single writer and identity,
> canonicalization, merge policy and the seed-closure invariant all live there. A
> rebuild that bypassed it would be a second, divergent definition of what the
> graph means.

And then it checks itself:

> The replay finishes by re-encoding every node and edge and checking the CIDs
> against the ones stored. That check is the contract: if it passes, the rebuilt
> graph is byte-identical to the one that was saved, and re-saving it produces
> the same scan root.

Which is also the mechanism that catches a schema drift. A binary whose registry
has had a field reordered will fail the CID check rather than silently rendering
a graph that means something different.

## Diffing

```go
// NodeChange is one node's fate between two scans.
type NodeChange struct {
    Type string
    Key  string
    Old  cid.Cid
    New  cid.Cid
    // Slots are the positional prop indices whose values differ …
    Slots []int
}
```

Positions rather than field names, because the addressed form is positional and
the registry publishes no way to map a position back to a name — a caller holding
the type's declared field list can resolve them.

Edges are paired on `(from, rel, to)`, their identity, so a change means the
edge's own props moved rather than that it points somewhere else.

The comparison itself is cheap: equal root CIDs mean nothing changed anywhere,
and below that, equal node CIDs prune whole subtrees of comparison.

{: .todo }
> `internal/store` implements save, load, rehydrate, CID-verify and diff, and
> `report` uses the first four. A `--resume` flag and a user-facing diff command
> are designed but not wired up; today, comparing two scans means rendering both
> and diffing the output.

## Why this and not a database

The previous design kept results in SQLite through GORM. Three things it could
not do:

- **Verify.** A row can be edited; a block whose hash does not match its content
  is detectably wrong.
- **Compare cheaply.** "Did anything change?" was a full-table comparison rather
  than a hash equality.
- **Be the same twice.** Auto-increment ids, timestamps and insertion order all
  varied between runs of the same scan.

Reference data — vocabulary, misspellings, source lists — still lives in SQLite,
because it is read-only input that nobody needs to address by content. Results do
not.

---

Next: **[Extending URLInsane](../extending/)**.
