---
title: 14 · Types and relations
parent: Part III · Inside the engine
nav_order: 2
---

# Types and relations
{: .no_toc }

- TOC
{:toc}

The registry is the schema. Nothing is a type unless it is registered, and every
entry carries a capability, a canonicalization function, an ordered field list
and a schema version.

## Node types

```console
$ urlinsane typo --list types
NAME        CAPABILITY  VERSION
asn         observed    1
domain      nameable    1
email       nameable    1
ip          observed    1
package     nameable    1
platform    observed    1
registrant  observed    1
repo        nameable    1
tld         observed    1
username    nameable    1
```

## Capability: what may be done to a type

```go
const (
    // Nameable types can be the root of variant generation.
    Nameable Capability = iota + 1
    // Observed types are only ever discovered, never varied.
    Observed
)
```

This is a two-value enum doing a lot of work. It is why `urlinsane typo
192.0.2.1` is refused: `ip` is Observed, so it cannot root a variant, and
falling through to `username` would have accepted the scan and emitted nonsense.
It is also why one omission algorithm covers domains, packages, repos, usernames
and emails without being registered five times — it binds to the *capability*,
not to a list of types.

Nameable is necessary but not sufficient. Variant expansion also requires
seed-closure membership and scope admission, both enforced by the applier; see
[Limits](../limits/).

## Relations

```console
$ urlinsane typo --list relations
NAME           CLASS        DEPTH COST
DOMAIN_OF      structural   0
EXISTS_ON      observation  1
HOSTED_ON      structural   0
IN_ASN         observation  1
LOCAL_PART     structural   0
MANIFEST       observation  1
MX             observation  1
NS             observation  1
OWNER          structural   0
PTR_TO         observation  1
REGISTERED_BY  observation  1
RESOLVES_TO    observation  1
TLD_OF         structural   0
VARIANT_OF     variant      0
```

Relation names read **in the direction of the arrow**: `email --DOMAIN_OF-->
domain`. Direction is load-bearing — seed closure and depth both propagate from
the `From` node to the `To` node — so a name that reads backwards against its
arrow is a standing invitation to wire an edge the wrong way round.

## Edge classes

```go
const (
    // Structural edges come from parsing the target string alone.
    Structural Class = iota + 1
    // Variant edges connect an origin to a generated variation.
    Variant
    // Observation edges required a lookup against some external service.
    Observation
)
```

| Class | Depth cost | Extends the seed closure |
|---|---|---|
| Structural | 0 | yes |
| Variant | 0 | no |
| Observation | 1 | no |

The dividing line between the first and last class is **whether the edge
required a network call**. Everything structural is derived from parsing the
target string alone. Everything that had to ask something else costs a hop.

Two consequences that are easy to miss:

**`MANIFEST` looks structural and is not.** Reading a repository's dependency
manifest feels like decomposition — it is just parsing a file. But it requires
fetching the repository, so it is an observation and costs depth. That is also
what keeps a repo scan from silently varying every declared dependency.

**Variant edges cost nothing but do not extend the closure.** A variant is
generated locally, so charging depth for it would make `--depth` mean two
different things. But the *variant* is not part of the seed closure, which is
what stops variants of variants — see the terminal rule in
[Limits](../limits/).

## Canonicalization

```go
// Canonical normalizes a raw key into the single form the graph converges on.
// Returning an error refuses the candidate outright.
type Canonical func(string) (string, error)
```

Each type supplies one. The domain canonicalizer lowercases, strips a trailing
dot, and refuses anything containing a delimiter; the email one splits on the
last `@`; the package one splits `registry:name`. What matters is that there is
exactly one, that it is total (every input either canonicalizes or is refused),
and that it is applied by the applier rather than by whoever emitted the
candidate.

## Schema versions and evolution

Every type carries a version. It starts at 1 and rises **only** when a field is
appended or tombstoned.

The rules, in full:

- Appending a field: version bumps, old blocks still decode, CIDs of existing
  nodes are unchanged.
- Tombstoning a field: version bumps, the slot stays occupied and is marked
  deprecated.
- Reordering or deleting a field: **never**. It changes the meaning of every
  content address already in the store, which corrupts diffs rather than
  failing loudly.

A decoder meeting a block from a newer binary detects the trailing fields it
does not understand and refuses to re-save, rather than silently truncating.

## Registering a type

Types and relations are declared as plain values and added to a `Registry`.
This is the real `domain` declaration, with the comments the code carries:

```go
{
    Name: TypeDomain, Cap: graph.Nameable, Version: 1, Canonical: canonDomain,
    // Declaration order is presentation order: identity first,
    // then reachability, then registration facts.
    Fields: []graph.FieldDef{
        {Name: FieldPunycode, Kind: graph.KindString},
        {Name: FieldLive, Kind: graph.KindBool},
        // Two sources assert this, so the winner is declared rather
        // than decided by whichever answered first.
        {Name: FieldCreated, Kind: graph.KindTime, Merge: graph.Precedence("rdap", "whois")},
        {Name: FieldRank, Kind: graph.KindInt},
    },
},
```

```go
{Name: RelExistsOn, Class: graph.Observation, Version: 1},

// Relation props carry data operators need: the per-node analyzers read the
// algorithm and edit distance straight off the edge.
{Name: RelVariantOf, Class: graph.Variant, Version: 1, Fields: []graph.FieldDef{
    {Name: FieldAlgorithm, Kind: graph.KindString},
    …
}},
```

`Registry.AddType` and `AddRel` validate on the way in: a type with no
canonicalization function, an unknown capability, a version below 1 or a
duplicate name is an error at registration. Registration happens during init,
before any scan — the registry is explicitly not safe for concurrent
registration, and is read-only afterwards.

Note that relations carry their own ordered field list, encoded exactly like a
node's. That is not decoration: `VARIANT_OF` stores which algorithm produced the
variant and its edit distance, and the analyzers read both straight off the
edge. A model with bare edges would have had to hide that somewhere worse.

---

Next: **[Operators](../operators/)**.
