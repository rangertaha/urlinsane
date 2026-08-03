---
title: Design document
parent: Reference
nav_order: 2
---

# Graph Engine Design

Status: implemented. It superseded the linear stage pipeline in
`internal/engine/processor.go`, which no longer exists; where this document and
the code disagree, the code is right and the disagreement is a bug in this file.

## Summary

urlinsane becomes a generic graph engine. A target string expands into a typed
property graph; operator plugins match patterns in that graph, do work, and
expand it further. Expansion stops when a round produces no new eligible work,
or when a limit binds — depth, budget, round cap or interrupt. Analyzers then
run over the whole graph and a report is generated from it.

Typosquatting is the first workload expressed this way, not the shape of the
engine. The same engine covers usernames, packages, repositories and emails
because those are node types with different operators attached.

The inversion: **the graph is the data, not the execution plan.** Today a DAG
orders collectors inside a hardcoded pipeline. Here the graph *is* the result.

Two layers, deliberately:

- A **compiled plan** (§5) — a static, inspectable, pinnable artifact naming
  every operator that may run, its trigger, and the limits. This is what
  `--explain` prints and what `--plan` pins for reproducibility.
- **Pattern dispatch in rounds** (§6) — the runtime mechanism. An operator runs
  when its pattern matches a node, which may be long after that node was
  created; matches are dispatched one generation per round, with a barrier
  between.

The plan answers "what will this do?". Dispatch answers "when?". Conflating
them is what made the earlier revisions of this document inconsistent.

### Where the DAG is, honestly

The entity graph is **cyclic by nature** — `domain → ip → (PTR) → domain` is
correct data, not a bug. The type-flow graph is cyclic for the same reason.
Neither can be topologically sorted, and termination never relies on
acyclicity (§8).

So nothing is topologically sorted to decide execution order. `internal/graph/dag`
survives for **plan presentation**: condensing the type-flow graph into strongly
connected components and ordering that condensation, so `--explain` renders in
a readable layered form. That requires Tarjan SCC plus a topo sort of the
condensation, not the current plain Kahn — a change to that package, not a
deletion.

## 1. Data model

```go
type Node struct {
    ID    NodeID     // hash(Type, Canonical(Key)) — stable for the node's life
    Type  NodeType
    Key   string     // canonicalized at admission
    Props Props      // ordered field list, §1.3
}

type Edge struct {
    ID       EdgeID  // hash(From, Rel, To)
    From, To NodeID
    Rel      Rel
    Props    Props
}
```

These are the applier's types. Operators never construct them — they emit
`NodeRef`/`EdgeRef`/`PropSet` (§4.3) and the applier canonicalizes, assigns
identity and applies merge policy.

Nodes and edges carry **no provenance, status, or findings**. Those live in
side tables keyed by `NodeID` / `EdgeID`. The reason is §1.2.

### 1.1 Edge classes

Every relation declares a class, which governs depth accounting (§8) and
variant eligibility (§8):

| Class | Relations | Depth cost |
|---|---|---|
| Structural | `LOCAL_PART`, `HAS_DOMAIN`, `HAS_TLD`, `OWNER`, `HOSTED_ON` | 0 |
| Variant | `VARIANT_OF` | 0 (bounded by the terminal rule instead) |
| Observation | `MANIFEST`, `RESOLVES_TO`, `NS`, `MX`, `PTR_TO`, `IN_ASN`, `REGISTERED_BY`, `EXISTS_ON` | 1 |

Relation names read in the direction of the arrow: `email --HAS_DOMAIN-->
domain`. Direction is load-bearing — seed closure and depth both propagate from
the `From` node to the `To` node — so a name that reads backwards against its
arrow is a standing invitation to wire an edge the wrong way round.

The line between the first and last class is **whether the edge required a
network call**. Everything structural is derived from parsing the target string
alone. `MANIFEST` looks like decomposition but requires fetching the
repository, so it is an observation and costs depth — which is also what keeps
a repo scan from silently varying every declared dependency (§8).

### 1.2 Three layers of identity

| Layer | Contents | Purpose |
|---|---|---|
| `NodeID` | `hash(Type, Canonical(Key))` | in-memory identity — scheduler, seen-set, cache, edges and side tables all key on this |
| Addressed form | node `(Type, Key, Props)`, edge `(From, Rel, To, Props)` — **only** | gets the CID; this is what diffs between runs |
| Side tables | provenance, per-pair status, belief, findings, truncation ledger | keyed by `NodeID`, except the ledger which is keyed by canonical candidate key; persisted beside the block, never inside it. Belief is engine-internal and not exposed to analyzers (§9) |

Keeping provenance out of the addressed form is what makes two identical scans
produce identical CIDs. With run ids and timestamps inside, every node would
differ on every run and cross-run diffing — a primary use case — would be
impossible.

`NodeID` is stable for the node's whole life; the CID changes as props
accumulate. Both are needed and they are not the same thing.

Canonicalization is a required registry field (§2) and runs **before the
admission decision**, not alongside it. Every candidate is canonicalized as it
leaves the applier's front door, whether or not it is ultimately admitted —
otherwise the truncation ledger's denylist (§8) would compare raw keys and
`Example.com` would slip past a row recorded for `example.com`.

Without canonicalization the convergence the whole design rests on fails
silently the first time one operator emits `Example.com` and another
`example.com`. A key that fails to canonicalize is refused with a recorded
status and never becomes a node.

### 1.3 Props are an ordered field list

`map[string]any` cannot be deterministically encoded — map iteration order is
unspecified and `any` has no stable IPLD representation. Since a content
address over the props is load-bearing, that is disqualifying.

Each type registers an **ordered** field list; a field's position is its stable
index, and values are stored positionally.

```go
type Props struct {
    typ    NodeType
    values []Value          // parallel to typ.Fields, addressed by index
}

func (p Props) Get(f Field) Value            // Field resolved once from the registry
func (p *Props) Set(f Field, v Value) error
func (p Props) Each(fn func(Field, Value))   // declaration order
```

Field handles are resolved from the registry when an operator registers, so an
unknown field name fails at registration rather than returning a per-access
boolean at runtime.
Values come from a closed kind set: `String, Int, Float, Bool, Bytes, Time`.

Ordering is structural rather than imposed, which buys:

- **Deterministic encoding with no sort step** — order is a property of the
  type, so identical values encode identically by construction.
- **Meaningful presentation order** — report columns follow declaration order,
  not `created, live, punycode, rank` alphabetical noise.
- **Compact encoding** — values are positional, so field-name strings are not
  repeated per node. At the 10⁵-node scale a wide scan reaches, that decides
  whether the graph fits in memory.
- **Cheap access** — an index, not a map hash.

The cost is an evolution rule that must be enforced: **fields are append-only.**
Reordering or deleting one changes the meaning of every content address already
in the store — corrupting diffs rather than failing loudly. Removal is by
tombstone; the slot stays, marked deprecated. Each type carries a schema
version so a decoder meeting a block from a newer binary **detects** trailing
fields it does not know.

It does not preserve them. `Props` is sized from the registry schema, so there
is nowhere to hold a slot the running binary does not declare; an earlier draft
of this section promised round-tripping, which would silently drop the unknown
values and change the CID on re-save — data loss disguised as a successful
write. The implemented behaviour is to **refuse**: `store` fails the
rehydrate CID check rather than writing back a truncated node. Preserving them
properly needs an unknown-slots escape hatch in `Props`, which is not built.

### 1.4 Conflicting assertions

Two operators may assert the same field on the same node — `whois` and `rdap`
both produce registration dates. Under concurrent dispatch, last-write-wins
would be decided by network timing, which is nondeterministic.

Every field therefore declares a merge policy:

```go
{Name: "created", Kind: Time, Merge: Precedence("rdap", "whois")}
```

All assertions are retained in the provenance side table; the **materialized**
value is selected by declared precedence, falling back to lowest operator id.
Deterministic regardless of arrival order — and disagreement between two
sources is preserved as signal for analyzers rather than silently resolved by
whoever answered first.

## 2. Type registry

This table is authoritative. Every type named anywhere in this document appears
here; nothing else is a type.

Every entry carries a capability, a canonicalization function and a schema
version (§1.3); versions start at 1 and rise only when a field is appended or
tombstoned.

| Type | Cap | Canonicalization |
|---|---|---|
| `domain` | Nameable | lowercase, IDNA→punycode, strip trailing dot |
| `username` | Nameable | platform-qualified `platform:handle`, lowercased |
| `package` | Nameable | `registry:name`; PEP 503 for PyPI, lowercase for npm |
| `repo` | Nameable | host + owner + name, lowercased |
| `email` | Nameable | lowercase domain, preserve local-part case |
| `tld` | Observed | ICANN public suffix, lowercased |
| `ip` | Observed | normalized v4/v6 text form |
| `asn` | Observed | `AS` + decimal, no leading zeros |
| `registrant` | Observed | normalized org/handle string |
| `platform` | Observed | lowercase host |

Relations are registered in the same registry with their class (§1.1) and their
own ordered field list.

Three canonicalization rules are easy to get wrong and are therefore spelled
out rather than left to the implementer:

- **`username` and `package` keys must be qualified** — `github:acme`,
  `npm:lodash`. An unqualified key is a canonicalization *refusal* (§1.2), not
  a best-effort guess. Without qualification `github:acme` and `pypi:acme`
  collapse into one node, which is exactly the convergence hazard qualification
  exists to prevent; the argument was previously made for packages only and
  applies identically to usernames.
- **`tld` is the ICANN public suffix**, falling back to the final label when the
  suffix is private or unrecognised — so `example.co.uk` yields `co.uk` but
  `acme.blogspot.com` yields `com`. A `tld` node should always name something
  registry-operated, or clustering by TLD means nothing.
- **Host canonicalization maps for lookup without validating labels.** Strict
  IDNA validation refuses precisely the malformed names variant algorithms are
  built to generate, so the scanner would drop its own most interesting output.

**Field declarations live with the type**, in the same registration call as the
capability and canonicalization. This document does not enumerate them: they
are per-type, append-only (§1.3), and belong in code where the tombstone rule
is enforceable.

`Nameable` means variant operators *may* apply — necessary but not sufficient;
eligibility also requires seed-closure membership (§8).

**A thing is a node only if it has identity worth converging on across the
graph.** IPs, domains, registrants and ASNs qualify. Registration dates and
geolocations do not — they are props. A nameserver is a `domain` reached by an
`NS` edge, not a type of its own; making it separate would split
`ns1.example.com` into two nodes and destroy the convergence that justifies
this design. Mail hosts likewise: an `MX` edge to a `domain`.

## 3. Seed expansion

Not special-cased. The target becomes a seed node and *decomposer* operators
match it like any other operator:

```
bob@example.com
  email ──LOCAL_PART──→ username(bob)
        └─HAS_DOMAIN──→ domain(example.com) ──HAS_TLD──→ tld(com)

github.com/acme/tool
  repo(github.com/acme/tool) ──HOSTED_ON──→ platform(github.com)
                             ├─OWNER──────→ username(acme)
                             └─MANIFEST───→ package(...)     ← observation, not decomposition
```

Decomposition edges are structural and cost no depth (§8). `MANIFEST` is the
exception in that diagram: reading a repository's manifest is a network call,
so it is an observation edge, and the packages it yields are outside the seed
closure (§8). Supporting a new input form — a Docker image ref, a bare IP — is
one plugin.

## 4. Operators

```go
type Operator interface {
    Id() string
    Version() int              // bump to invalidate cached results
    Trigger() Trigger          // when this runs
    Emits() Effects            // everything this may produce
    Resource() string          // rate-limit class: "dns", "whois", "npm", "http"
    Exec(ctx context.Context, v View) (Delta, Outcome)
}

type Trigger struct {
    On    Selector             // Types{...} or Caps{Nameable}
    Where []Condition          // HasProp(f), HasEdge(rel), BeliefAbove(t)
    Reads Reads                // []Field and []Rel, resolved at registration
}

type Effects struct {
    Nodes []NodeType
    Rels  []Rel
    Props []Field
}
```

`Emits` declares relations and props, not just node types, for two reasons. An
operator like `geo` emits no node type at all — only props on an `ip` — and
would otherwise be invisible to plan compilation. And without declared props,
the plan cannot detect a `Where: HasProp(punycode)` that no operator in the
plan can ever satisfy: a permanently dead operator the plan would list as
active, which defeats the point of `--explain`.

`Reads` holds resolved `Field` and `Rel` handles rather than strings, so an
unknown name fails at registration — the same rule §1.3 applies to props.

### 4.1 Triggers are patterns, not producer dependencies

An earlier revision gave operators `Requires() []string` — ids of operators
that must run first. That repeats the mistake the current collector DAG makes:
**depending on producers rather than on data.** `geo` does not depend on the
`ip` operator having run; it depends on there being an IP. `Requires("ip")`
breaks the moment a second operator also produces IPs, or `ip` is deselected.

Binding to a pattern dissolves it:

```
dns   : On domain,  Where HasProp(punycode)  → ip, domain (via NS, MX)
geo   : On ip                                 → props on ip
ptr   : On ip                                 → domain
whois : On domain                             → registrant, props
rdap  : On domain                             → registrant, props
```

`dns` needing normalized punycode is expressed as a *data* condition, satisfied
by whichever operator produced it. There is no ordering declaration anywhere in
the system, and no cycle risk from the ordering mechanism because there is no
ordering mechanism.

A `Selector` binds by type or by capability, so one omission algorithm covers
every `Nameable` type instead of being registered four times.

A **variant operator** is any operator that declares `VARIANT_OF` in its
`Effects.Rels`. That declaration — not a naming convention or a separate plugin
family — is what subjects it to the terminal rule and the seed-closure
restriction (§8), and it is checkable at plan-compile time.

### 4.2 Views are patterns, not snapshots of everything

```go
type View interface {
    Type() NodeType
    Key() string
    Prop(f Field) Value          // only fields declared in Trigger.Reads
    Edges(rel Rel) []*Edge       // only relations declared in Trigger.Reads
}
```

**Props are filtered too, not just relations.** An earlier revision handed the
operator the whole `*Node`, which let it read fields it never declared. That
silently breaks the cache: the read-set digest (§6.3) is built from
`Trigger.Reads`, so an undeclared field could change without changing the key,
and the operator would be served a stale result forever. Filtering both props
and relations is what makes "pure function of its declared inputs" true rather
than aspirational, and it is enforced by the view rather than trusted to plugin
authors.

`Edges`, not neighbours: relation props carry data operators need. `VARIANT_OF`
holds `algorithm` and `distance`, which the ported per-node analyzers (§4.4)
read directly. Returning bare nodes would hide them.

Operators cannot query arbitrary graph state — that would make results
order-dependent, runs irreproducible and caching unsound.

Because dispatch waits for the pattern, a view is never the empty
neighbourhood an at-creation snapshot would give. An operator needing the
origin of a variant declares `Where: HasEdge(VARIANT_OF)` and reads it; that is
how today's origin/variant analyzers port over.

Anything needing whole-graph context is an analyzer (§9), which runs on a
stable graph.

### 4.3 Operators never mutate

They return a `Delta`; the scheduler is the single writer.

```go
type Delta struct {
    Nodes []NodeRef   // type + raw key; the applier canonicalizes and assigns NodeID
    Edges []EdgeRef   // From/To as NodeRef, Rel, props
    Props []PropSet   // (NodeRef|EdgeRef, Field, Value)
}
```

Operators emit `NodeRef`, not `Node`: they cannot compute a `NodeID` because
canonicalization (§1.2) belongs to the registry and the applier, and letting a
plugin mint an identity is how convergence quietly breaks. Every prop write goes
through `PropSet` so the field's merge policy (§1.4) applies uniformly — a
`Node` carrying a whole `Props` block would bypass it.

Deltas are additive only. Nothing removes a node, edge or prop, which is what
makes the graph monotonic within a run and a delta safely replayable.

Concurrent operators are race-free by construction, with no locking and nothing
for plugin authors to get wrong, and each step is a diffable, replayable value.
The applier is single-threaded, which is not a bottleneck: operator work is
network-bound and dominates by orders of magnitude.

### 4.4 Family mapping

| Today | Becomes |
|---|---|
| `Algorithm` | operator with `Selector{Caps: Nameable}`, emitting `VARIANT_OF` edges |
| `Collector` | operator emitting observed nodes, edges and props |
| `Analyzer` (per-node) | operator with `Where: HasEdge(VARIANT_OF)`, emitting props |
| `Analyzer` (global) | graph analyzer, §9 |

## 5. The compiled plan

Before anything runs, the seed, the registry and the CLI flags compile into a
plan:

```go
type Plan struct {
    Hash      string       // identifies this plan exactly
    Seed      SeedSpec     // type, canonical key, scope
    Operators []OpBinding  // defined once in §12.2
    Model     ModelRef     // id + CID of the execution model, §10.4
    Limits    Limits       // depth, budgets, ttl, rate classes
}
```

`Hash` covers the engine model and every plugin model (§10.6) as well as the
operators. A plan pinning operator versions but letting models float would
reproduce neither the traversal nor, where a plugin model decides what an
operator emits, the graph.

Compilation prunes operators unreachable from the seed by intersecting `Emits`
against triggers. Because the type-flow graph is cyclic, reachability is
transitive closure, not a topological walk — pruning is an over-approximation
and is presented as such, since `Emits` is a *may* not a *must*.

Compilation also reports **dead operators**: one whose `Where` names a prop or
relation that nothing in the plan declares in its `Effects` can never fire.
`--explain` lists these separately from active operators rather than implying
they will run — which is the whole reason `Effects` covers props and relations
and not just node types (§4).

The plan is the answer to "what will this do?":

- `--explain` renders it, SCC-condensed and layered, and exits 0 — dead
  operators are listed, not fatal, since a plan that cannot fire an operator is
  usually the user narrowing scope rather than an error.
- `--plan FILE` writes it; supplying an existing file pins it. A pinned plan
  plus a warm cache reproduces a run exactly, which is what makes results
  citable.
- Execution dispatches only operators present in the plan. Runtime never
  reaches past it.

## 6. Execution

The scheduler responds to **three kinds of event**, not one. Every earlier
revision of this document assumed deltas were the only trigger, and three
separate bugs followed from that single assumption: a retry that could never
fire, a belief condition that could never re-evaluate, and a tree parent chosen
by network timing.

| Event | Raised by | Effect |
|---|---|---|
| **Delta** | an operator returning | apply it; re-evaluate triggers for touched nodes |
| **Timer** | a retriable failure | re-enqueue the pair after backoff, inside its round |
| **Barrier** | a dispatch round draining | finalize parents, belief, gating and truncation |

### 6.1 Rounds and the barrier

The scheduler runs in **rounds**. Round 0 is seeding: the engine admits the
canonicalized seed node and runs barrier 0, which gives the seed its prior
belief. Every round after that evaluates triggers against the graph as of the
previous barrier, dispatches every eligible pair, and drains — including any
outstanding timers. Then its barrier runs, and only at a barrier are
irreversible decisions taken:

- tree parents are finalized, once and for all (§10.3),
- belief is **recomputed** for every node from its props as of this barrier,
- `BeliefAbove` gates are evaluated for the next round,
- truncation admits its prefix (§8).

Parents are fixed once; belief is not. A node created in round g has only the
props its creating operator set — observation operators on it dispatch in
round g+1 and land at barrier g+1. Computing belief once at creation would
therefore leave it a bare transition prior forever, which is the opposite of
what §10.2 needs. Recomputation is safe because belief is a pure function of
`(parent belief, current props)` and the parent never moves: the same run
recomputes the same values in the same order.

What recomputation cannot do is reverse a decision already taken on an earlier
value. A node pruned at barrier g stays pruned even if barrier g+1 would have
raised it — pruning is irreversible by construction (§6.1), and the truncation
ledger (§8) records it so the loss is visible rather than silent.

All four are irreversible. A pruned node is gone; a whois already sent cannot be
recalled. Taking them on provisional belief — belief derived from whichever
parent happened to arrive first — makes the final graph depend on network
timing, and no amount of later re-parenting repairs it, because re-parenting
cannot un-prune a node or un-send a request.

**Rounds, not depth waves.** An earlier revision keyed barriers on depth. That
was wrong, and wrong in the way that mattered most: depth counts observation
hops only (§8), so decomposition, every variant, *and* the `dns`/`whois` calls
on those variants all sit at depth 0. Belief would first be computed at the
barrier ending depth 0 — after the hundred thousand whois calls it exists to
prevent. Gating was structurally a wave too late.

Rounds fix it because they follow dispatch, not distance. Decomposers run in
round 1, variant operators in round 2, and observation operators on those
variants in round 3 — with a barrier in between, so belief is available exactly
where `BeliefAbove` needs it.

**A round is exactly one dispatch generation.** Deltas arriving mid-round mark
pairs *eligible*; they are dispatched by the **next** round, never the current
one. This is the load-bearing rule. If a delta dispatched immediately, `dns`
emitting an `ip` in round 3 would let `geo` run on that `ip` in round 3 too —
before barrier 3 computed its belief — reproducing the depth-wave bug one level
down. One generation per round is what keeps "belief is final before anything
acts on it" true at every level, not just the first.

**Parents are finalized once.** A node's parent is fixed at the barrier of the
round that created it, chosen by §10.3's rule over the in-edges present at that
barrier. Every operator dispatched in a round drains before its barrier, so
every in-edge created in that round is present. In-edges arriving in later
rounds never re-parent — consistent with §10.3's existing rule that later
in-edges are not used as evidence, and what keeps parent choice a function of
round membership rather than of arrival order.

Depth survives as a budget and termination concept (§8) only. It no longer
defines synchronization.

The cost is real and accepted: **no cross-round pipelining.** Round g+1 cannot
begin before round g drains. Reproducibility is a stated goal of this design;
pipelining is not.

### 6.2 Retry lives inside its round

A retriable failure schedules a timer and stays in its own round; the barrier
waits for it. Each pair has a bounded attempt count and each resource class a
deadline. When that deadline passes, the pair records `timeout` permanently for
the round and the barrier proceeds.

This is what makes retry compatible with the barrier. A retry landing *after*
the barrier would add props to a node whose belief was already final and whose
descendants were already gated on it. Bounding retries inside the round fixes
the outcome set before belief is computed, so a slow service delays its round
rather than corrupting it — and deadline truncation is reported exactly like
budget truncation (§8).

### 6.3 Dispatch

When a delta touches node N, triggers for N are re-evaluated and newly-matching
`(node, operator)` pairs are marked eligible for the **next** round (§6.1) —
never dispatched into the current one. Workers execute the current round's pairs
concurrently; the applier serialises deltas and is the single writer.

**Re-dispatch is bounded by the read-set.** A pair runs once per
`(NodeID, operator id, read-set digest)`. New data the operator declared it
reads makes it eligible again; new data it does not read does not. A hard cap
of 3 revisions per pair is the loop backstop.

Re-evaluating a `BeliefAbove` gate does not consume a revision — a pair that
stayed `skipped` never ran, so there is nothing to revise. Only an actual
execution counts, which keeps a gate that flickers across barriers from
exhausting the cap before the operator has run even once.

Belief is deliberately **not** part of the read-set digest. `BeliefAbove` is a
barrier-time condition, evaluated once per round, so it never participates in
delta-driven re-dispatch and never needs the digest to notice it changed.

**Admission control, not queue dropping.** The applier never blocks and never
discards queued work. When the frontier exceeds its bound — a limit on
in-flight candidates, distinct from the per-type and global node budgets of §8
— the scheduler stops admitting new nodes and writes `frontier` rows to the
truncation ledger. That removes both the deadlock and the silent loss.

**Rate limiting is two-layer.** Each operator declares one `Resource()` class
with its own token bucket; beneath it the HTTP/DNS transport limits per host.
Today's single global `--delay` is meaningless once one run talks to DNS,
whois, npm, PyPI and GitHub — the limit protecting the strictest service
throttles everything else to the same crawl.

**Cache key covers everything the operator reads:**

```
H(operator id, version, modelCIDs, NodeID, readSetDigest, resourceConfigDigest)
```

`Trigger.Reads` supplies the read-set digest — the same declaration doing
double duty. `resourceConfigDigest` covers config that changes what a service
returns: nameservers for `dns`, registry URL and token for `npm`, proxy and UA
policy for `http`. `Version()` is what makes a fixed operator stop serving its old wrong answers,
and `modelCIDs` does the same for a retrained plugin model (§10.6).

A cache hit still occupies its round: the delta is applied at the same point a
live call's would be, so a warm run and a cold run produce the same rounds, the
same barriers and the same graph. A cache that short-circuited the round
structure would make `--plan` reproducibility depend on cache state, which is
exactly what pinning is meant to eliminate.

## 7. Failure is data

An operator error is not merely a log line. For a scanner, *how* a lookup
failed is the finding: `NXDOMAIN` proves the name is free, `SERVFAIL` means
misconfigured, a timeout means nothing was learned. Collapsing all three into
"error" discards the signal the tool exists to collect.

Every `(node, operator)` pair records a terminal status in the side table, with
its provenance:

| Status | Meaning | Example |
|---|---|---|
| `ok` | learned something positive | `A` records returned; package exists on registry |
| `empty` | authoritatively determined absence | `NXDOMAIN`; registry 404 |
| `failed` | the lookup itself broke | `SERVFAIL`; malformed response |
| `timeout` | nothing was learned | deadline passed (§6.2) |
| `skipped` | never attempted | gated off by `BeliefAbove` (§10.2), or a variant operator refused outside the seed closure (§8) |

`skipped` is the one status that is **not** terminal. Belief is recomputed at
every barrier (§6.1), so a pair gated off at barrier g becomes eligible again at
barrier g+1 if belief has since risen above the threshold; recording `skipped`
terminally would make the first gate permanent and quietly starve nodes the
model later favours. Closure refusals are the exception — the seed closure never
grows, so that skip is final.

The `ok`/`empty` split is the one that matters and is easy to get wrong:
`NXDOMAIN` is `empty`, not `failed` and not `ok`. It is a successful
determination of absence, which is exactly the signal a squatting scanner
exists to collect — and it is what makes `live` (§9) mean anything. Reports
distinguish "confirmed absent" from "could not determine", and a run with
widespread timeouts is visibly degraded rather than quietly thin.

Transient failures schedule a timer and retry inside their own round (§6.2),
bounded by an attempt count and the resource's deadline; they are not left to
delta-driven re-dispatch, which would never fire because a failure produces no
delta. Permanent failures record and move on. An operator failing never fails
the run.

## 8. Termination and truncation

- **Seen-set** over `(NodeID, operator id, read-set digest)` — convergent and
  cyclic paths never redo work.
- **Depth counts observation hops only**, default 3; structural and variant
  edges cost nothing (§1.1). It is a budget and termination concept only —
  synchronization is by dispatch round (§6.1), not by depth. Under whole-edge
  counting `bob@example.com` spent its budget on decomposition and never
  reached an IP, silently gutting the composite-input case this design exists
  for.
- **Variants are terminal.** A node reached by `VARIANT_OF` is never handed
  back to a variant operator.
- **Variant roots are seed-closure members.** The seed closure is the seed plus
  everything reachable from it by *structural* edges. `ptr` emits `domain`
  nodes, which are `Nameable`; without this rule every reverse-PTR domain on
  every variant's IP becomes a new variant root — the exact explosion the
  capability flag was meant to prevent, arriving by an edge the terminal rule
  does not cover.

  Because `MANIFEST` is an observation edge (§1.1), manifest-derived packages
  fall outside the closure. `urlinsane typo github.com/acme/tool` therefore
  does not silently generate typo variants of several hundred dependencies;
  varying those is an explicit `urlinsane typo package <name>` or the
  dependency-confusion analysis, which needs no variants at all.
- **The engine enforces this, not the plugin**, in two places. The scheduler
  will not dispatch a variant operator against a node outside the closure, so
  the work is never done; and the applier rejects any `VARIANT_OF` edge whose
  source is outside it, so an operator that emits one anyway cannot smuggle it
  in. Dispatch-side gating is the optimization, applier-side rejection is the
  invariant. An invariant whose whole purpose is preventing combinatorial
  explosion must not be opt-in for plugin authors.
- **Budgets**, per type and global.
- **A round cap.** Depth bounds observation hops, but a cyclic type-flow
  (`ptr → domain → dns → ip → ptr`) keeps generating rounds until depth bites,
  and an oscillating trigger pattern might not converge at all. The cap is the
  backstop the per-pair revision cap cannot provide, and hitting it is reported
  like any other truncation.
- **Depth and seed-closure membership are scheduler state**, held in the side
  table alongside provenance and status. Both are derived from structure, and
  keeping them out of the addressed form stops them perturbing CIDs.

**Belief pruning is truncation and is reported as such.** A pruned candidate
was never admitted, so it has no `NodeID` and cannot carry a per-pair status.
Truncation of every kind — belief threshold, budget, round cap, round deadline
— is therefore recorded in a **truncation ledger**: a side table of
`(type, key, depth, belief, reason)` rows for candidates the engine declined to
admit. The report renders it exactly like any other section.

**Truncation is not the same as an invariant rejection**, and only the former
denies. A candidate refused because its key would not canonicalize, or because
an operator tried to root a variant on it from outside the closure, is recorded
as a *rejection* and nothing more. Putting those in the ledger would deny the
candidate forever — and a node refused as the source of one bad `VARIANT_OF`
edge may still be perfectly legitimate when another operator reaches it by an
observation edge. (Found while implementing phase 0: the two paths look alike
until you ask what a second operator finding the same candidate should see.)

The ledger is also a **denylist**, not merely a record. A declined candidate
that a later operator re-emits is refused again by the applier and the existing
row is kept. Without that, "pruning is irreversible" (§6.1) would be false in
the one case that matters — a second operator reaching the same candidate would
re-admit it, and whether a node exists would depend on how many operators
happened to find it. Without this it is
the quieter version of the same bug: the node never gets a `dns` result, so
`--filter absent` under-reports while the output still reads as complete.

Truncation is **deterministic and reported**. It happens at a barrier (§6.1),
where belief is final: candidates for the next round are sorted by
`(-belief, depth, type, key)` — the execution model's belief (§10.2) first, so
a bound budget spends itself on the most promising frontier rather than
alphabetically — and the prefix is admitted. The ordering is total and fixed
under a pinned plan, so output stays byte-identical across runs even when
truncated. A truncated graph that reads as complete is a correctness bug, not a
cosmetic one.

## 9. Analysis

Analyzers run **once, after expansion stops** — for any reason, including
Ctrl-C. There is exactly one lifetime.

```go
type Analyzer interface {
    Id() string
    Exec(ctx context.Context, g *Graph) ([]Finding, error)
}

type Finding struct {
    Kind     string        // "campaign", "dep-confusion", "high-risk-variant"
    Severity Severity      // info | low | medium | high | critical
    Nodes    []NodeID      // admitted nodes this concerns
    Declined []LedgerRef   // ledger rows this concerns, §8
    Summary  string
    Evidence []Provenance
}
```

A finding can concern a candidate that was never admitted — a pruned package is
still a dependency-confusion gap — which is why `Declined` references ledger
rows and not just `NodeID`s.

`Graph` exposes nodes and edges plus read-only side-table access: status,
provenance, plugin-model scores (§10.6) and the truncation ledger — **but not
the engine's belief**.

The asymmetry is deliberate and worth stating, because both are model output
sitting in the same side table. A plugin score is an *observation about the
entity* — how parked a page looks, how plausible a variant reads — and is
legitimate evidence for a finding. Engine belief is a *scheduling artifact*
about traversal order, and means nothing about the entity itself. Withholding belief is what makes §10.2's "the execution model never contributes
to a reported number" true by construction rather than by convention; an
analyzer that could read it could launder it into a `Severity`, and no reviewer
would spot it. Belief remains visible under `--verbose` for diagnostics,
outside the analysis path. An analyzer that could not see
status could not tell "confirmed absent" from "never tried", which is the
distinction §7 exists to preserve — and one that could not see the ledger would
report a truncated graph as complete.

The valuable analyses exist only because infrastructure is nodes:

- **Campaign clustering** — variants sharing an IP, nameserver or registrant.
- **Scoring** — resolves + MX + recent registration + low edit distance. This
  is rule-based and entirely independent of the execution model (§10.2), which
  never contributes to a reported number.
- **Dependency confusion** — an internal `package` with no `EXISTS_ON` edge to
  a public registry.

Two derived terms the CLI exposes, defined here so nothing else has to invent
them:

- **Existence** is a rollup of §7's statuses, in three values rather than two:
  a node is **live** if any observation operator returned `ok`, **absent** if
  none did and at least one returned `empty`, and **unknown** if every attempt
  ended in `failed`, `timeout` or `skipped`.

  This replaces `registered`/`unregistered`, which were domain-only and could
  not express the third case. Collapsing "confirmed absent" into "unregistered"
  discards exactly the distinction §7 exists to preserve — a variant nobody
  could resolve is not a variant proven free — and the trichotomy generalizes:
  a package absent from a registry is `absent` in the same sense a domain
  returning `NXDOMAIN` is.
- A node's **risk** is the maximum `Severity` among findings referencing it.
  It is what `--filter risk>N` and `--fail-on SEV` select on, and the only
  user-facing score in the system. `Severity` is a named ordered enum —
  `info | low | medium | high | critical` — rather than a bare int, so
  `--fail-on high` means the same thing in every release; an integer scale
  would drift silently the first time a new level was inserted.

Findings land in the side table, so reporting is a graph query rather than a
parallel data path. `Evidence` carrying provenance is what lets a report answer
"why is this here?" instead of asserting a score.

## 10. The execution model

One HMM, trained offline and loaded as an artifact, whose sole job is to steer
expansion. It is not a scoring model for the report.

- **States** — a small set of latent node statuses.
- **Initial distribution** — `P(state)` for the seed, which has no parent.
  Barrier 0 (§6.1) assigns it; every other node's belief derives from a parent.
- **Transitions** — `P(state | parent state, Rel)`, conditioned on the
  relation, so sharing an IP transmits far more belief than sharing a TLD.
- **Emissions** — `P(observed props | state)` over props that operators in the
  plan actually produce: resolves, MX present, registrar, registration age.
  An emission over a prop nothing emits is dead weight in the model, the same
  way a `Where` nothing satisfies is a dead operator (§5).

### 10.1 Inference is the forward algorithm, and nothing else

The expansion tree gives every node exactly one parent — chosen
deterministically at a barrier, §10.3 — so the path seed → node is a
*sequence*, which is what an HMM is defined over. That makes this a literal HMM
rather than a graph model borrowing the name.

Forward filtering is inherently incremental: a node's belief is a pure function
of its parent's belief and its own props, recomputed at every barrier as those
props arrive (§6.1). No backward pass, no Viterbi, no belief propagation, no
whole-graph inference.

This shape is forced by the purpose. A model that runs after expansion can only
annotate a graph that has already been built — by which point every network
call it might have saved has already been made. To affect execution it has to
be available *during* execution, which means forward-only.

### 10.2 What it is used for

Execution control, exclusively. All three are evaluated at a round barrier
(§6.1) against belief recomputed from props *as of that barrier*, so an
admitted node's belief sharpens each round as observations arrive:

- **Frontier ordering** — `(-belief, depth, type, key)` decides which
  candidates enter the next round when a budget binds.
- **Pruning** — candidates below threshold never enter the next round.
- **Operator gating** — `Where: BeliefAbove(t)`, so nothing whois-es what the
  model considers uninteresting.

The first two rank **candidates**, which have no props yet, so their belief is
the transition prior alone — parent belief pushed through the edge's relation.
Only gating reads the belief of an **admitted** node, which does carry
observations. Both are "belief" and both are deterministic, but a candidate's
is necessarily less informed, and the ledger records it so a pruning decision
can be second-guessed after the fact.

Explicitly *not* used for: report scores, findings, severities, or anything a
user reads as a probability. This is what removes the calibration requirement
(§10.5): a miscalibrated model that mis-ranks the frontier wastes network
calls, whereas a miscalibrated model in a report makes a false accusation.

### 10.3 Determinism and convergent paths

A node reachable by several edges has several candidate parents. The tree
parent is the minimum by `(depth, Rel, parent NodeID)` over the **complete**
candidate set for the round that created the node — which is precisely why
parent selection happens at a barrier (§6.1) rather than at admission.
Min-so-far is not min-final: a parent revised after the fact revises belief
after the gating and pruning that belief already drove, and pruning cannot be
undone.

Belief is then recomputed at every barrier from `(parent belief, the node's
props as of that barrier)` rather than accumulated observation by observation,
so it depends neither on which operator returned first nor on which edge
arrived first — only on which round each edge belongs to, which is
deterministic. The parent is fixed; the belief sharpens. Under
a pinned plan the expansion tree, the beliefs, the gating decisions and the
admission order are all reproducible.

Later in-edges are not used as evidence: a domain independently reached from
two suspicious variants gains nothing from the second. That is a real loss of
signal, and recovering it would require belief propagation over the cyclic
graph — not worth the cost for what is an execution heuristic.

### 10.4 Artifact — trainable, loadable, pinnable

A model is a dag-cbor block with a schema and a CID, like everything else here:
states, the initial distribution, transition and emission tables in log space,
smoothing parameters, and training provenance (corpus CIDs, algorithm, RNG
seed, date).

- `--model NAME[@cid]`; the default ships in `datasets/`.
- **The model CID enters the plan hash (§5).** A pinned plan pins the model, so
  a reproduced run reproduces the same traversal. Without this, `--plan`
  reproducibility fails the first time the model is retrained — the graph
  itself would differ, not merely its annotations.
- `datasets train` — a maintainer command, running Baum-Welch over recorded
  expansion traces. The RNG seed is recorded so training reproduces. It sits in
  `cmd/datasets` rather than beside `typo` because it builds reference data, as
  `import` and `download` do, and the scanner only ever consumes what it makes
  (§12).
- Traces are not a by-product: `typo --trace FILE` writes the
  `(parent belief, relation, props, outcome)` tuples training needs. Recording
  is opt-in because it persists observation data a normal scan discards, and a
  tool that silently accumulated a corpus of everything it had ever looked up
  would be making that choice on the user's behalf.

Numerics: log space throughout, Dirichlet priors on transitions, an explicit
OOV symbol for unseen emissions.

### 10.5 Limits, and why they are tolerable

Training needs recorded traces, and supervised training needs labels this
project does not yet have. Unsupervised EM produces clusters that still require
interpretation.

What makes that tolerable is the failure mode. Because the model only steers
execution, a bad model degrades *efficiency*, not correctness — it explores in
a worse order and wastes calls. A uniform model reduces exactly to today's
behaviour: breadth-first, unranked. The engine therefore ships and runs
correctly before any model exists, and a model that turns out to be poor can be
dropped without invalidating a single result.

### 10.6 Models inside plugins

The model above is the engine's, and it is the only one that touches
scheduling. Operators may carry models of their own — a variant operator
scoring which typos a human would plausibly make, a collector classifying
parked pages from response features. Those are separate in every sense that
matters:

| | Execution model (§10.1–§10.5) | Plugin models |
|---|---|---|
| Owner | the engine | the operator |
| Affects | admission, pruning, gating | what the operator emits |
| Output | side-table belief, never data | graph structure, and side-table scores |
| Count | exactly one | any number |
| Corpus | recorded expansion traces | whatever that plugin needs |

`internal/model` is therefore a **library**, not the engine's model: HMM
primitives (forward, Baum-Welch, log-space arithmetic, smoothing), the dag-cbor
artifact format, and loading. The engine instantiates one; plugins instantiate
their own; neither knows about the other's states or alphabet.

Three integration rules, all forced by reproducibility:

1. **Plugin model CIDs enter the plan hash.** `OpBinding` carries
   `Models []ModelRef` beside `Flags` and `Config` (§12.2). This is not
   optional: a variant operator whose name model changed emits *different
   variants*, so a plan pinning operator versions but letting their models
   float would fail to reproduce the graph itself — a much worse failure than
   the engine model's, which only reorders traversal.
2. **Plugin model CIDs enter that operator's cache key** (§6). `Version()`
   covers the operator's code; a model is data the code reads, so a retrained
   model with unchanged code must still invalidate.
3. **Flags follow §12.2 namespacing.** `--model` selects the engine's;
   `--<operator>.model` selects a plugin's. `--list models` groups by owner.

`datasets train` takes a target — `train exec` for the engine model,
`train <operator>` for a plugin's — because their corpora, state spaces and
alphabets have nothing to do with each other.

One placement rule for what a plugin model produces. Structural decisions —
which variants exist — go into the graph and are pinned by the plan. Derived
judgment scores go to the side table, for the same reason execution belief does:
as props they would make every node's CID depend on a model version and break
cross-run diffing the moment anything is retrained.

Unlike engine belief, plugin scores **are** visible to analyzers (§9). They
describe the entity rather than the traversal, so a finding may cite one as
evidence — which is also why they, not belief, are what a plugin model is for.


## 11. Output and determinism

Two modes, which are different things and are documented as such:

- **Report** — every format except ndjson. Emitted after analysis, with nodes
  sorted by `(type, key)`, edges by `(From, Rel, To)`, ledger rows by
  `(type, key, reason)`, and props already in structural order (§1.3). It always
  carries the truncation ledger (§8) and, when expansion stopped early, a
  partial marker (§12.4). Byte-identical across runs, so "what changed since last week"
  is a diff rather than noise.
- **Event stream** — `-o ndjson` only. Emitted as facts are found, explicitly
  unordered, for piping into other tools. It cannot be canonically sorted and
  does not claim to be.

`--save` is a separate sink from stdout, so it never conflicts with `-o`: each
sink gets the form implied by its own format. `--save out.json` writes a sorted
report to the file whatever stdout is doing, and `--save out.ndjson` writes the
stream form — the format determines the mode, not the sink.

### 11.1 What forces the partial marker

Partial is *derived*, not passed in. Any run-level truncation, and any non-empty
truncation ledger, marks the report partial whatever the caller believed —
including belief-gated declines, which are working exactly as configured and
still leave the report incomplete. The caller never observes a budget bounding
expansion mid-run, so trusting a flag here would let the most common truncation
of all render as a complete scan.

## 12. CLI

```
urlinsane typo   [<scope>] <target> [flags]   scan
urlinsane report <target> [flags]             render a saved scan of it
```

Two verbs, and the split is between **scanning and not scanning**. An earlier
draft of this section specified five — `typo`, `explain`, `why`, `list`,
`train` — on the grounds that only the first scans. That grouping is true and
useless: it puts four unlike things together because of what they are not.

`explain`, `why` and `list` are flags on `typo`. Each takes the same target, the
same scope positional and the same selection flags, because what each describes
*is* the scan: which operators the plan compiled, what produced a result, which
algorithms are registered for the target's type. A verb re-accepting that entire
surface is `typo` under another name with a second parser to keep in step, and
the drift shows up as a plan `explain` prints that `typo` would not have run.

The objection that a suppressing flag gives `typo` several exit paths before it
has parsed a target is answered by parsing first: target and scope are resolved
for every invocation, and only then does a suppressing flag decide whether to
run the plan or print it. That ordering is what makes `--explain` describe the
run it suppressed.

`report` is the one that genuinely differs, because it does not scan *and* does
not describe a scan it is about to run — it renders one that already happened.
It takes a target, but rejects everything that shapes a scan, since `--depth`
and `--algorithm` cannot change what a finished graph contains. A flag on `typo`
whose presence invalidated the rest of `typo`'s surface would be worse than a
verb. It reads from the content-addressed store (§12.7): the target names the
scan, `--at CID` names an exact one, and the replay re-encodes and CID-checks
every node and edge, so what is rendered is byte-identical to what was saved.
That is also why filters live there — `--filter` was always a display filter
over built rows, so re-filtering a ten-minute scan costs nothing.

`train` is not a scanner concern at all. It runs Baum-Welch over stored graphs
(§10.4), which is maintainer work on reference data, and it belongs beside the
other reference-data commands in `cmd/datasets`.

The scope positional is a node type, **optional and narrowing**, comma-separated
for several:

```
urlinsane typo bob@example.com             # broad: every nameable in the seed closure
urlinsane typo username bob@example.com    # narrow: vary bob only
urlinsane typo domain bob@example.com      # narrow: vary example.com only
urlinsane typo username,domain bob@example.com
urlinsane typo package npm:lodash
urlinsane typo repo github.com/acme/tool
```

Omitted, scope is every `Nameable` node in the seed closure — **including the
composite node itself**, so a bare email seed also yields whole-email variants.
Supplied, it filters that set. The target parses identically either way; scope
never changes how the string is read, only what gets varied.

Scope is an **applier invariant**, not a dispatch hint. The scheduler skips
variant operators on out-of-scope nodes as an optimization, and the applier
rejects any `VARIANT_OF` edge whose source type the scope excludes — the same
two-layer arrangement the seed-closure rule uses, and for the same reason: a
dispatch gate is something an operator can be written around. A scope that was
only validated, hashed into the plan and printed by `explain` would leave
`typo username bob@example.com` silently identical to the unscoped run, which is
the one thing the positional exists to prevent. The rejection is a distinct kind
(`outside-scope`, not `outside-seed-closure`) and writes no ledger row: an
out-of-scope root is a legitimate variant root the user chose not to expand, and
denying the candidate would poison a later run with a wider scope.

There is no `--scope` flag — the positional is the only spelling. This retires
`--type` (`cmd/urlinsane/typo.go:105`) and `entity.Classify` as the seeding
path (`internal/config/config.go:178`), which today reads `bob@example.com` as
`user` purely because it contains an `@`. Parsing rule: with two positionals
the first must be registered `Nameable` types; with one, it is the target.

### 12.1 Progressive disclosure

`typo` currently carries ~25 flags in one undifferentiated block, with
`Category` values nothing uses. `--help` shows what most runs need;
`--help-all` shows everything.

```
  -d, --depth int        observation hops from the seed (default 3)
  -f, --filter strings   live, absent, unknown, risk>N (all defined in §9)
  -a, --algorithm ids    restrict variant generation (§12.10)
  -l, --language ids     languages the language-driven algorithms run over
  -k, --keyboard ids     keyboards the keyboard-driven algorithms run over
  -c, --collect ids      restrict observation; ^id excludes
  -o, --output string    table | json | ndjson | csv | dot (default table)
      --save PATH        write the report to PATH; format from extension
      --save-graph       persist the graph to the store; prints its root CID
      --fail-on SEV      exit 2 if any finding reaches SEV
      --quiet            silence stderr; the exit code is the whole answer
      --tui              interactive view of the same scan (§12.8)
      --plan FILE        write the plan, or pin an existing one
      --model NAME[@cid] execution model; --<op>.model for plugin models (§10.6)
  -v, --verbose
```

Timeouts, delays, TTL, budgets, workers and user-agent policy move behind
`--help-all`, unchanged in behaviour. Operator-owned flags such as
`--nameservers` and `--registry` are governed by §12.2 instead: they appear in
`--help` when their operator is in the plan.

Four cleanups: `--options`/`--ids`/`--opts` — three aliases for one hidden
flag — collapse into the single `--list TOPIC`; `--registered`/`--unregistered`, two booleans
settable to contradict each other, become `--filter live` and `--filter absent`
— which also gains the `unknown` case they could not express; `--format`,
`--file` and `--dir` collapse into `-o` plus `--save`; and `--explain` becomes
`--explain`, which suppresses the run it describes.

**`--filter` selects rows in the report, never work in the scan.** It is applied
after expansion, over built rows. Narrowing the *scan* is what `--depth`,
`--algorithm` and the scope positional are for. Two selection languages
competing over the same run would be worse than the confusion of one, so
`--filter` gains no scan-side twin — instead §12.6 requires a filtered-empty
report to say what it hid.

### 12.2 Flags are scoped to the plan

`--nameservers` is meaningless for a package scan; `--registry` is meaningless
for a domain scan. Today all ~25 flags apply to every run, so the help text is
a union in which most entries are inert for whatever the user is actually
doing.

Flags are therefore **declared by the operator that consumes them**, and the
compiled plan (§5) decides which are relevant:

```go
type OpBinding struct {
    Id       string
    Version  int
    Trigger  Trigger
    Resource string
    Flags    []Flag      // declared by the operator, namespaced by its id
    Config   Values      // parsed values
    Models   []ModelRef  // this operator's models, §10.6
}
```

This is the single definition; §5 and §10.6 refer to it rather than restating
it.

Three rules:

1. **Every flag is always registered; only help is filtered.** The plan is
   compiled *from* flags, so flags cannot depend on the plan — a two-pass parse
   would be circular. Registering everything and filtering the help text
   sidesteps that entirely, because the real problem is discoverability, not
   parsing. `urlinsane typo package --help` shows registry flags;
   `typo domain --help` shows resolver flags. With no scope the plan is still
   the seed *closure*, which for a bare domain target is narrow anyway; only a
   composite seed like an email or repo yields several scopes, and there help
   groups by scope instead of filtering.
2. **Namespaced, with a bare form when unambiguous.** `--npm.registry` always
   works; bare `--registry` works when exactly one operator in the plan
   declares it. A collision makes only the bare form an error, and at
   registration time.
3. **An inapplicable flag is an error, not a no-op.** `typo domain x.com
   --registry npm` fails with "`--registry` is declared by operator `npm`,
   which is not in this plan (scope: domain)". A silently inert flag is worse
   than a rejected one — the user believes they configured something.

The mechanism already exists: urfave/cli carries `Category` on every flag, and
§12.1 notes nothing currently uses it. `Category` becomes the owning operator,
and the help template filters on the active plan. No new flag library, no
custom parser.

This also tightens two things elsewhere in the design:

- `OpBinding.Config` is part of the plan, so `--plan` pins **flag values**, not
  just operator selection. A pinned plan reproduces configuration.
- `resourceConfigDigest` in the cache key (§6) is computed from `Config`, so
  changing `--nameservers` invalidates the `dns` operator's cached results and
  nothing else's. Scoped ownership is what makes that precise rather than a
  blunt global cache flush.

Help now has three tiers: common flags, plan-relevant flags (`--help`), and
everything (`--help-all`).

### 12.3 Progress, and where output goes

One rule decides every question about where output goes:

> **stderr is for the human, stdout is for the machine.**

Progress and diagnostics go to stderr and are suppressed when stdout is not a
TTY. The report goes to stdout. So
`urlinsane typo acme.com -o json > out.json` shows a live scan on the terminal
and writes clean JSON to the file, and **no flag is needed to say so** — the
common case needs no configuration at all.

The current bar is `progressbar.DefaultSilent(1000)` — hardcoded total, never
displayed. The default view is a one-shot run: a bar while it works, the report
when it finishes.

```
  acme.com · domain · 27 algorithms · 9 observers

  ████████████████████░░░░░░░░░░  68%   1,284 nodes   round 4/64   0:08
  dns 240/s   whois 2/s   geo 11/s
```

Round and depth answer "how much is left" in a way a spinner cannot, because
expansion is round-based. The per-resource rates are the part worth keeping from
a richer display: `whois 2/s` beside `dns 240/s` explains a stalled scan without
the user having to know that whois is rate-limited. The bar is replaced by the
report when the scan ends — it is progress, not output.

`--tui` opens an interactive view of the same scan (§12.8). It is **opt-in**:
the default path is a one-shot command, so nothing about existing scripts, CI or
muscle memory changes, and the interactive view can land later without a
migration. `--tui` with a non-TTY stdout is an error, not a silent downgrade,
under §12.2's rule that an inapplicable flag is an error rather than a no-op.

Findings are not shown during the scan: analysis is a distinct final phase,
rendered as `analyzing…` when expansion stops.

**`-o ndjson` emits a record as soon as that record is done** — when every
operator bound to a node's type has returned, failed or been skipped. `-o json`
is one document and can only be written at the end; `ndjson` is one object per
node and there is no reason to hold them. A node is emitted once, complete: a
record that appeared early and was later amended would force every consumer to
keep state and de-duplicate. `ndjson` is already the one format that declines to
promise an ordering (§11), so completion order costs nothing — and it is more
useful than sorted order, being the order the answers actually arrive in.

### 12.4 Interruption and exit codes

Ctrl-C stops expansion at the **end of the current round**, not mid-round: the
barrier still runs, so parents, belief and the truncation ledger are finalized
rather than left half-computed. Analyzers then run over what exists and the
report is marked partial. A second Ctrl-C aborts immediately without a report,
for when a round is stuck behind a slow resource. The first Ctrl-C says which
of the two just happened, because otherwise an impatient second press looks
like the first did nothing:

```
^C
  stopping at the end of round 4 — press ^C again to abort now
```

A ten-minute scan has produced something worth keeping, and round-by-round
expansion guarantees it is a coherent prefix rather than an arbitrary
cross-section. `PARTIAL` appears in the header of **every** format including
json, so a truncated scan is never mistaken downstream for a complete one.

Exit `0` clean, `1` execution error, `2` a finding at or above `--fail-on`.
That is what makes the tool a CI gate for dependency confusion; today the only
way to react to results is to parse stdout.

`--quiet` silences stderr entirely, making the exit code the whole answer:

```
urlinsane typo --quiet --fail-on high npm:acme-internal   # $? is the result
```

It does not suppress first-run setup (§12.6): a silent thirty-second first run
is indistinguishable from a hang.

`NO_COLOR` and non-TTY stdout are respected.

### 12.5 Existence is three-state, and always shown

`live`, `absent` and `unknown` are the most valuable distinction this tool
draws and the easiest to lose. *Absent* is "we asked, it is not there";
*unknown* is "we could not tell". Collapsing them turns a broken network into a
clean bill of health.

So the split is **unconditional, not a flag**: a glyph in every row, and the
counts in the footer of every run.

```
  ● live 42   ○ absent 1,203   ? unknown 39
```

`●` green, `○` dim, `?` amber. The legend costs one line and prevents the
single most likely misreading of the output.

This is also why `--filter` must explain an empty result (§12.6) rather than
printing nothing: `-f live` over a scan where every lookup timed out yields the
same empty table as a scan that found nothing, and those are opposite outcomes.

### 12.6 Say what was not done

Three kinds of gap are currently silent, and each makes a degraded run
indistinguishable from a complete one.

**Omitted operators.** An operator whose dependency is missing is left out of
the plan rather than included and failed (§4), which is right — but it must be
reported, or a scan with no geolocation looks like a scan of a target with no
geolocation.

**Declined candidates.** The truncation ledger (§8) already records every
candidate refused by a budget, a depth bound or a frontier cap. It is written
and never surfaced.

**First-run setup.** The embedded datasets are extracted on first run. Silently,
today — and when extraction fails, the operators that needed them simply vanish.

```
  ⚠ 2 observers were not in this plan
      geo        no maxmind database        (first run did not complete)
      npm        no registry dataset        (datasets import not run)

  ⚠ 412 candidates were declined
      budget     380   node cap reached at round 6
      depth      32    beyond --depth 3
```

An empty report explains its own emptiness in the same spirit:

```
  no rows matched --filter live

  1,284 hidden:  1,203 absent · 39 unknown · 42 live but below --depth 2

  ⚠ 39 unknown — every dns attempt timed out. the network, not the result.
```

**The last line is the point.** It separates a clean negative from a broken run,
which is the difference between "this brand is unsquatted" and "our resolver is
down".

### 12.7 Provenance: `--why`

The graph records how every node was reached — which operator produced it, under
which relation, with which props (§1.3). An investigator's first question about
a suspicious hit is "how did you get here", and there is no way to ask it.

It is a flag rather than a verb (§12) because the answer is *about* a particular
scan: the same target, scope and selection that produced the hit are what locate
it in the stored graph.

```
urlinsane typo acme.com --save-graph      # earlier: scanned and kept
urlinsane typo acme.com --why acmе.com    # now: ask about one result

  acmе.com · domain · live · critical

  reached by
    acme.com ──VARIANT_OF── acmе.com
      algorithm  hr  homoglyph replacement
      distance   1
      change     e → е  (U+0065 → U+0435, cyrillic small letter ie)

  observed
    dns-a      ✓  91.195.240.1                          0.31s
    dns-ns     ✓  ns1.above.com, ns2.above.com          0.28s
    whois      ✓  created 2026-07-30, registrar Namecheap
    ptr        ✗  timeout after 10s

  findings
    critical   homoglyph in a registrable label
    high       shares nameservers with 4 other variants of this seed
```

`--why` reads a **stored** graph rather than re-scanning: re-observing one name
repeats network work and can disagree with the run being explained, which is
the one thing an explanation must not do.

That requires the scan to have been saved, and saving is **opt-in**
(`--save-graph`) for the reason §10.4 gives about traces: a tool that
silently accumulated a corpus of everything it had ever looked up would be
making that choice on the user's behalf. `--why` without a saved graph is an
error that says which flag would have produced one — not a silent re-scan.

`--save-graph` takes no path. The store is content-addressed and shared, so a
scan's root CID is its name and the store is where `report <target>`
and `--resume` both look; writing one run to a loose `run.dag` would put it
somewhere neither could find it. `--store DIR` selects a different store, and
`--save-graph` prints the root CID it wrote.

The store already exists (`internal/store`), carries no timestamp or run
id in its root, and is therefore content-addressed in a way that makes two
identical scans produce the same CID — so `--why`, cross-run diffing and `--resume`
are all the same mechanism seen from different angles.

### 12.8 The interactive view

`--tui` opens the same scan as a navigable view, borrowing from two programs
that already solved halves of this problem: **htop** for live density, **nvim**
for modal navigation.

```
┌──────────────────────────────────────────────────────────────────────────────┐
│  dns    [||||||||||||||||||||||    ] 240/s  1284      round    4/64          │
│  whois  [|||                       ]   2/s    37      depth    2/3           │
│  geo    [|||||||||                 ]  11/s   118      nodes    1,284         │
│  ──────────────────────────────────────────      live 42  absent 1203  ? 39  │
├──────────────────────────────────────────────────────────────────────────────┤
│ SEV NAME               EX ALG  D  DNS NS  MX  WHO GEO PTR  PROGRESS          │
│ ⣿   acmе.com           ●  hr   1  ✓   ✓   ✓   ✓   ✓   ⠋   [||||||||  ] 5/6  │
│ ⣿   acme-support.com   ●  cb   2  ✓   ✓   ⠋   ·   ·   ·   [|||       ] 2/6  │
│ ⣶   acnne.com          ●  cr   1  ✓   ✓   ✓   ✓   ✓   ✓   [||||||||||] 6/6  │
│     acme.net           ○  tld  1  ✓   −   −   −   −   −   [||||||||||] 1/1  │
│     acmee.com          ?  cr   1  ✗   ⠋   ·   ·   ·   ·   [          ] 0/6  │
├──────────────────────────────────────────────────────────────────────────────┤
│ NORMAL   results   sort:sev   filter:—              acme.com · domain  0:08  │
│ :                                                                            │
└──────────────────────────────────────────────────────────────────────────────┘
```

**Progress is per target, and per observer within it.** The default bar (§12.3)
reports one number for the whole run, which averages away the only thing worth
knowing: which target is stuck, and on what. Every node is observed by several
operators independently, so the row is where progress belongs and the column is
where a stalled resource shows up — `WHO` reading `⠋` down several rows while
`DNS` is `✓` everywhere identifies the bottleneck without the user knowing which
resources are rate-limited. The htop-style meters above say the same thing in
aggregate.

| | |
|---|---|
| `✓` | returned ok |
| `✗` | failed — SERVFAIL, refused, error |
| `⠋` | in flight |
| `·` | not yet attempted |
| `−` | not applicable: the node is absent, so nothing downstream runs |

`−` is load-bearing. `acme.net` is absent, so ns/mx/whois/geo/ptr never run — it
is `1/1` with a full gauge, not `1/6` and stalled. Without it most rows in a
typical scan would look permanently unfinished, because most variants are
absent. It also makes the three-state existence (§12.5) self-evident per row:
`acmee.com` shows `✗` then a retry in flight, which is *why* it is `?` and not
`○`.

**Rows are selected with the standard keys, and `h`/`l` work the pane.** `l`
opens the detail pane on the right for the selected row and moves focus into it;
`h` closes it and returns. The direction is the meaning — `l` moves right, into
the detail; `h` moves left, back to the list — so the pane needs no dedicated
binding and no toggle to remember.

| Key | |
|---|---|
| `j` `k` · `↓` `↑` | select row |
| `gg` `G` | first, last |
| `C-d` `C-u` | half page |
| `C-f` `C-b` | page |
| `l` `→` | open the detail pane, focus it |
| `h` `←` | close the detail pane |
| `/` `n` `N` | search, next, previous |
| `gd` | provenance (§12.7) |
| `C-o` `C-i` | jump back, forward |

```
┌─ results ──────────────────────────────┬─ acmе.com ───────────────────┐
│ SEV NAME             EX ALG D PROGRESS │ domain · live · critical     │
│ ⣿   acmе.com         ●  hr  1 [||||| ] │                              │
│ ⣿   acme-support.com ●  cb  2 [||    ] │ dns    91.195.240.1          │
│ ⣶   acnne.com        ●  cr  1 [||||||] │ ns     ns1.above.com  +4 ▸   │
│     acme.net         ○  tld 1 [||||||] │ whois  2026-07-30 Namecheap  │
│     acmee.com        ?  cr  1 [      ] │ geo    DE · Frankfurt        │
│                                        │ ptr    ✗ timeout after 10s   │
│                                        │                              │
│                                        │ VARIANT_OF acme.com          │
│                                        │   hr · d1 · e→е U+0435       │
└────────────────────────────────────────┴──────────────────────────────┘
```

The pane is closed by default, because at eighty columns it costs half the width
of the table it is describing. `l` is cheap enough that opening it per row is
not a burden, and `h` gets the width back.

**The pane is itself a row list, navigated and highlighted the same way.** Its
fields are rows: `j`/`k` moves through them with the selection highlighted, and
`l` descends again — from the `ns` row into the four other variants sharing that
nameserver, from a registrant into everything they registered. `h` pops back.
Panes stack, so navigation is a path through the graph rather than a toggle
between two fixed views.

That is what makes §12.7's provenance and the traversal below the same gesture
instead of two features: descending from a row *is* following an edge.

### 12.9 One component, composed

The view is deliberately not a screen layout. **Everything visible is a list of
rows**, and the differences between the results table, the detail pane, the
ledger, the plan and the findings are the rows they carry — not the code that
draws them.

| Component | Responsibility |
|---|---|
| `list` | rows, columns, selection, highlight, scroll, sort, filter, fold |
| `meters` | the per-resource gauges of the top panel |
| `status` | mode, active sort and filter, target, elapsed |
| `cmdline` | `:` and `/`, with completion over the same registries `list` uses |
| `overlay` | help (`g?`), confirmations, errors |
| `stack` | the pane stack: `l` pushes, `h` pops |

Three properties follow, and they are the reason for the arrangement rather than
consequences of it:

1. **A new pane is data, not code.** The ledger, the plan and the findings are
   `list`s over different row sets. Adding one — say, a buffer of every node
   sharing a nameserver — is a query and a column spec.
2. **Behaviour is written once.** Sorting, filtering, search, highlight and
   fold live in `list`, so they work identically in the detail pane and the
   results table. A binding that works in one works in all of them, which is
   what makes the keymap learnable.
3. **The engine is not consulted about presentation.** A `list` is fed rows and
   knows nothing about `graph`, so the view can be tested without a scan and
   the engine carries no rendering concerns. This is the §13 layering applied
   to the interface: the same reason operators are pure functions of a view.

The row source is an interface, not a type: anything that can produce rows and
say what a row descends into can be a pane. That is the seam the TUI is built
on, and it is why `--tui` can be added without the engine knowing it exists.

**The graph is navigable.** The engine produces variants, the addresses they
resolve to, the nameservers those share and the registrants behind those. A
report flattens that to a list and discards the edges — which is the part an
investigation needs: *this variant is suspicious → what else shares its
nameserver → who registered those*. That is a traversal, and the modal idioms
already exist for it:

| Idiom | Here |
|---|---|
| buffers | projections of one graph: results, findings, ledger, plan |
| splits | list ▏detail |
| jumplist `C-o` / `C-i` | move along graph edges, and back |
| `gd` | go to provenance — the edge that produced this node (§12.7) |
| `:` | the flags, available mid-run: `:filter live`, `:depth 4` |
| `/` | fuzzy find over names |
| folds | collapse variants by algorithm, or unfold one row to named columns |

htop's F-key hints remain as an affordance for discovery, but every one is also
a `:` command, so the two idioms do not compete for the same action. `g?`
overlays the keymap for the current buffer, so nothing has to be memorised.

The observer columns collapse as the terminal narrows — named, then glyphs, then
a bare ratio — and `zo` unfolds the selected row to the named form without
resizing anything else:

```
  wide (>140)   DNS ✓  NS ✓  MX ✓  WHOIS ✓  GEO ✓  PTR ⠋   [||||||||  ] 5/6
  default       ✓   ✓   ✓   ✓   ✓   ⠋                      [||||||||  ] 5/6
  narrow (<100) ✓✓✓✓✓⠋                                     5/6
```

`:ls` lists the buffers — results, findings, ledger, plan, log — each a
projection of the same graph rather than a separate fetch.

The scan keeps running while the graph is browsed; rounds land in the background
and the list grows under the cursor.

**`:report` runs analysis without waiting for the scan.** Analysis is a distinct
final phase (§9); `:report` runs it over whatever exists at that moment and opens
the result as a buffer, marked `PARTIAL` under the same rule as an interrupted
run (§12.4). A scan that has found something worth acting on at round 6 should
not have to reach round 64 before it can be read, and the round barrier (§6.1)
guarantees the graph it reads is a coherent prefix rather than a half-applied
delta. `:export` then writes that report in any of the formats of §11.

**`:depth 4` mid-run is the argument for the whole view.** In a one-shot command,
discovering the bound was too low means re-running from scratch and paying for
every lookup again. Here it widens the frontier and the existing graph stands.

The truncation ledger (§8) and the omitted operators (§12.6) are buffers rather
than warnings that scroll past, each carrying the command that lifts the limit:
`depth 32 — beyond depth 3 — :depth 4 to admit these`.

### 12.10 Selecting what runs, and what it runs over

Four flags, and the distinction between them is the point:

| Flag | Selects | Default |
|---|---|---|
| `-a, --algorithm` | which variant operators generate | all |
| `-l, --language` | which languages those algorithms run *over* | all registered |
| `-k, --keyboard` | which keyboard layouts they run *over* | all registered |
| `-c, --collect` | which observation operators run | all in the plan |

**`-a` picks the algorithm; `-l` and `-k` pick the data it runs over.** `-a hr`
selects homoglyph replacement; `-l ru` decides whose homoglyphs. They compose
rather than compete, which is why the old `--algorithms` alone could not express
what `--languages` did, and why collapsing them into one flag would lose a
distinction the engine already makes: `variant.Options` carries `Languages` and
`Keyboards` as *parameters to the operators*, while `Algorithms` selects the
operators themselves.

The defaults are "everything registered", and deliberately so. A tool whose
failure mode is missing a squatted name should not narrow its own recall by
default; a scan that generates too many candidates is visibly expensive, while a
scan that generates too few is silently wrong. Locale-derived defaults were
considered and rejected for the same reason — a homoglyph attack on a `.com` is
cross-script by definition, so inferring `ru` from a `.ru` TLD would suppress
exactly the case the algorithm exists to catch.

**Exclusion is first-class, because it is the common case.** A `^` prefix removes
an id from the set that would otherwise run:

```
urlinsane typo acme.com -c ^whois          # everything except whois
urlinsane typo acme.com -c dns,ptr         # only these
urlinsane typo acme.com -a ^cb,^afx        # drop the noisy generators
```

Dropping one slow or rate-limited observer is far more common than enumerating
the ones to keep, and `^` spells that without a second `--no-collect` flag whose
interaction with the positive form would have to be defined. Mixing forms in one
flag is an error rather than a precedence rule to memorise.

`^` rather than a bare `-` because `-whois` is indistinguishable from a flag to
any argument parser, and rather than `!` or `~` because both are shell
metacharacters that would need quoting. In the interactive view the same syntax
works on `:algo` and `:collect`, where no parser is in the way.

An unknown id is an **error**, not an empty selection: `variant.Select` already
reports unknown algorithms rather than silently matching nothing, and `-c` and
`-l` follow it. A typo that quietly scanned with no observers would be the worst
possible outcome — a clean report produced by doing nothing.

This requires one engine change: `observe.Options` currently carries only
dependencies (`Resolver`, `Whois`, `Geo`, `Prober`) and `observe.New` builds the
full set, so there is nowhere to express a selection. It gains an `Ids []string`
alongside them, applied where `variant.Select` applies its own.

## 13. Testing and observability

The purity constraints buy four clean layers:

- **Operators** — pure functions of a view; a literal view in, an expected
  delta out. No engine, no network. An operator carrying a model (§10.6) pins
  it by CID in the test, since the model is an input like any other.
- **Scheduler** — synthetic operators exercise round barriers, timer-driven
  retry and its deadline, parent finalization under convergent in-edges,
  pattern re-dispatch, the revision cap, admission control, budgets, cyclic
  data and failure policy, without a single real plugin. Two invariants deserve
  dedicated tests because violating them is silent: a declined candidate
  re-emitted by a second operator stays declined (§8), and a pair gated off at
  one barrier still runs at a later one once belief rises (§7). This is the subtlest
  machinery in the design and the layer most worth over-testing.
- **End-to-end** — golden graphs from recorded fixtures, compared after
  canonical sorting.
- **Reproducibility** — the same pinned plan replayed against the same fixtures
  must produce byte-identical output, including the truncation ledger, and must
  do so with a cold cache and a warm one. Nearly every constraint in this
  document exists to make that test pass; without it they are unenforced
  conventions.

Observability: `--explain` for the plan, per-operator counts and timings under
`--verbose`, truncation always surfaced in the report.

## 14. Package layout

```
internal/
  graph/
    model.go      Node, Edge, NodeID, EdgeID, edge classes
    registry.go   type + relation registration, capabilities, canonicalization
    props.go      ordered field list, kinds, merge policy, schema evolution
    delta.go      Delta, apply, admission, canonicalization rejection
    view.go       pattern-scoped view
    trigger.go    Selector, Condition, match index
    schedule.go   delta/timer/barrier events, rounds, read-set digests, admission
    limit.go      resource classes and transport-level host limiting
    cache.go      (operator, version, models, node, read-set, resource-config) keys
    plan.go       compilation, pruning, pinning, --explain rendering
    side.go       provenance, status, belief, findings, truncation ledger
    dag/          SCC condensation + topo sort of the condensation, for plan layout
  model/          library — used by the engine and by plugins alike (§10.6)
    hmm.go        forward step, Baum-Welch, log-space arithmetic
    codec.go      dag-cbor model blocks, CID identity, training provenance
    smooth.go     Dirichlet priors, OOV emissions
    trace.go      expansion-trace recording for `train exec` (§10.4)
  store/     content-addressed scan storage, replay, diffing
  plugins/        everything that acts on the graph, one directory per plugin
    plugins.go    the three registries: operators, analyzers, algorithms
    decompose/    domain, email, pkg, repo          + all
    variant/      the 27 algorithms, one each       + all
    observe/      dns, ptr, whois, idn, geo,        + all
                  pkg, usr, repo                    + observetest
    analyze/      campaign, scoring, depconfusion   + all
    report/       table, json, ndjson, csv, dot     + all
```

**Every plugin is one directory, and every family is a library plus its
plugins.** The family package holds what its plugins share — `observe` owns
Options, the per-call timeout and the schema vocabulary; `variant` owns Spec,
the Generate signature and the keyboard and language combinators — and each
plugin imports it. Composition lives in `<family>/all`, which imports the
plugins; without that separation a plugin importing its family would import its
siblings through it, and the cycle is not hypothetical.

An earlier arrangement grouped by *target* instead — network, service, social,
repo — on the reasoning that an npm plugin contributes both an operator and an
algorithm and belongs in one place. That holds for a third-party plugin and not
for the shipped ones: splitting dns from whois from geo meant either
duplicating the shared vocabulary or exporting a package's internals to itself.
Cohesion is real, and the directory tree should not fight it.

## 15. Migration

Hard cutover. `db.Domain` is retired as the runtime payload and every plugin
moves to `Operator` in one pass; there is no adapter shim. The branch is not
shippable until the cutover completes — the accepted cost of not maintaining
two engines.

Retired: `internal/engine/processor.go` with its `Stage`/`stageFunc`;
`db.Domain` as runtime payload with its embedded record slices;
`internal.Algorithm`, `internal.Collector`, `internal.Analyzer`; `entity.Type`
as a closed enum and `entity.Classify` as the seeding path; `--type`,
`--registered` and `--unregistered` (replaced by `--filter live|absent|unknown`,
§9), `--format`, `--file`, `--dir`, `--options`, and the global
`--delay`/`--random` throttle (replaced by per-resource buckets, §6.3).

Kept: `internal/engine/dag`, moved to `internal/graph/dag` and reworked for SCC condensation; the
store's content-addressing approach.

Demoted rather than dropped: the language and keyboard plugin families. Both
were plugin interfaces wrapping what is only data — a language is a directory
of curated `.lst` files, a keyboard is a layout — and neither had behaviour a
plugin needed to supply. Languages are now rows the dataset database answers
(`internal/dataset`), keyboards are layouts compiled in from `pkg/kb`, and a new
language is added by adding a dataset directory rather than by writing Go. What
they feed, the variant operators, is unchanged.

### Phases

| # | Scope | Done when |
|---|---|---|
| 0 | `graph` core: model, registry, props, identity, delta, applier invariants | round-trip and canonicalization tests pass; CID stable across runs; applier rejects out-of-closure `VARIANT_OF` |
| 1 | Scheduler: triggers, match index, rounds and barriers, timers, re-dispatch, admission, budgets, status, cache | synthetic operators exercise cycles, budgets, timer retry, barrier parent finalization, revision cap |
| 2 | Plan: compilation, pruning, pinning, `--explain` | pinned plan reproduces a run byte-identically with a cold cache and a warm one (§13) |
| 3 | Decomposers and variant operators | every current algorithm ported; closure rule (enforced since phase 0) exercised by real variants |
| 4 | Observation operators | every current collector ported; rate classes assigned |
| 5 | Analyzers, report, CLI | feature parity with today's output, plus findings and the truncation ledger |
| 6 | Execution model + `train`: trace recording, forward filtering, Baum-Welch, model blocks, belief gating | beliefs reproduce under a pinned plan; gating measurably cuts network calls |
| 7 | Store rework to nodes and edges; `--resume` | cross-run diff works; resume continues a partial graph |

`--resume` depends on phase 7 and is not available before it — the store today
models an entity union with embedded records and cannot rehydrate a partial
graph.

Phase 6 is separable from the engine: everything before it runs with a uniform
model, which reduces to breadth-first unranked expansion (§10.5), so the engine
ships without waiting on model work.

Plugin models (§10.6) are not a phase. They arrive with whichever operator
wants one, in phase 3 or 4, and pull `internal/model` forward to that point —
the library lands when the first model does, engine or plugin, whichever comes
first.

## 16. Open

Everything previously listed here has been decided and moved into the body.
What remains genuinely undecided:

- Per-type budget defaults (global default and depth are set: unbounded and 3).
- Per-resource retry attempt counts, round deadlines (§6.2) and the round cap
  (§8) — the mechanisms are fixed, the numbers are not.
- Whether `--filter` grows into an expression language or stays a keyword set.
- In the interactive view (§12.8): whether a
  rescan-one-node binding should exist at all, given it would put two
  observations of the same node from different moments into one content-addressed
  graph, which §1.2 says must not blur; whether `:q` mid-scan cancels at the
  round barrier and still writes a partial report or discards; and which TUI
  library, since either `bubbletea` or `tview` is a new dependency.
- Whether a `cert` type is worth adding — it was removed from the registry
  because no operator emits one, so plan compilation would report it dead
  (§5). It returns when a TLS operator does.
- Training corpus for phase 6: whether to record expansion traces from real
  runs, bootstrap from public blocklists, or hand-label a seed set. §10.5 makes
  this non-blocking — a uniform model ships — but it decides when gating
  actually starts paying for itself.
- State cardinality — how many latent statuses earn their keep over a binary
  split, given the model only has to rank a frontier. Still open, but no longer
  constrained by the plumbing: `BeliefModel.Step` carries an opaque
  `graph.State` between parent and child, so a model of any width propagates
  its own posterior. An earlier signature passed the parent's *scalar*, which
  forced a model to reconstruct a distribution from one number — exact at two
  states, and at three or more the maximum-entropy guess rather than the
  posterior the parent had. That silently answered this question as "two" as a
  side effect of a type signature, and a wider model would have produced
  plausible, wrong numbers.
