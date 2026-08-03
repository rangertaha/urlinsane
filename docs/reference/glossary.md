---
title: Glossary
parent: Reference
nav_order: 4
---

# Glossary
{: .no_toc }

Terms this book uses precisely. Each entry links to the chapter that covers it
in depth.

- TOC
{:toc}

## The attack

**affix squatting** — Bracketing a package or repository name with an ecosystem
affix — `py-acme`, `node-acme`, `acme-utils`, `acme-js`. No error is involved;
the fake name often looks *more* idiomatic than the real one. Algorithm `afx`.
See [Why names get mistyped]({{ site.baseurl }}/attack/errors/).

**bitsquatting** — A variant one *bit* away from the target, not one keystroke:
`axample.com` for `example.com`. The error source is memory corruption in a
resolving host rather than a human, demonstrated at internet scale by Dinaburg
(Black Hat 2011). Algorithm `bf`. See
[Why names get mistyped]({{ site.baseurl }}/attack/errors/).

**brandjacking** — The umbrella term for the commercial-harm framing of every
squatting variety. It describes intent and damage rather than a generation
technique. See [What typosquatting is]({{ site.baseurl }}/attack/typosquatting/).

**combosquatting** — The target name plus a real word: `example-login.com`,
`secure-example.com`. There is no misspelling for a careful reader to catch,
which is why Kintis et al. (CCS 2017) found these both more numerous and
longer-lived than typo variants. Algorithm `cb`. See
[What typosquatting is]({{ site.baseurl }}/attack/typosquatting/).

**cybersquatting** — The legal term for registering a name in bad faith with
respect to a trademark. It overlaps typosquatting without coinciding with it:
many typosquats infringe nothing, and many cybersquats contain no typo. See
[What typosquatting is]({{ site.baseurl }}/attack/typosquatting/).

**dependency confusion** — A build resolves a package name against both a
private and a public registry, the public one has no such package, and an
attacker who registers it publicly wins resolution. Absence on the public
registry is the finding, which is why it can only be expressed against
three-state existence. See [The wider surface]({{ site.baseurl }}/attack/surface/)
and [Analysis]({{ site.baseurl }}/internals/analysis/).

**doppelganger domain** — The missing-dot case: `wwwexample.com`, or
`marketingexample.com` for `marketing.example.com`. A subdomain boundary the
typist failed to type. Algorithms `do` and `di`. See
[What typosquatting is]({{ site.baseurl }}/attack/typosquatting/).

**homoglyph / homograph attack** — A variant that is not a different string to a
human, only to a machine: a Cyrillic `а` for a Latin `a`, or any of the Unicode
characters rendering close to `l`, `1`, `I`, `0`, `O`. Registered as punycode,
rendered as the original. Algorithm `hr`. See
[Why names get mistyped]({{ site.baseurl }}/attack/errors/).

**levelsquatting** — The target name appears as a *subdomain* of an
attacker-controlled domain: `example.com.login-secure.net`. Address bars that
truncate from the left are the usual target. Algorithm `si` in the generative
direction. See [What typosquatting is]({{ site.baseurl }}/attack/typosquatting/).

**namespace confusion** — Moving a name between the namespacing conventions
registries use — npm's `@org/pkg`, a repo's `org/pkg`, the flat `org-pkg` — so a
scoped package is impersonated by an unscoped one or vice versa. Algorithm
`nsc`. See [The wider surface]({{ site.baseurl }}/attack/surface/).

**parking** — A registered name that hosts no real service, typically pointed at
one of a handful of well-known parking nameservers and monetised by ads or
resale. Parked squats cluster tightly on shared nameservers, which is what makes
them visible in a scan. See [Defending a name]({{ site.baseurl }}/attack/defending/).

**public suffix list** — The list of suffixes under which registrations are
made (`com`, `co.uk`, `github.io`). URLInsane uses it to decide that a dotted
string is a domain, to find the registrable label, and as the substitution
source for the `tld` algorithm. See
[Targets and scope]({{ site.baseurl }}/guide/targets/).

**punycode / IDN** — An internationalized domain name holds non-ASCII
characters; punycode (`xn--...`) is its ASCII encoding, and what is actually
registered. Browsers decide per-vendor whether to render the Unicode form or the
punycode, so the same name looks different in different clients. See
[Why names get mistyped]({{ site.baseurl }}/attack/errors/).

**soundsquatting** — A variant that sounds the same read aloud: a homophone of
the target or of a word inside it. Matters for voice interfaces, dictation, and
users who learned a name by ear. Algorithm `hs`. See
[What typosquatting is]({{ site.baseurl }}/attack/typosquatting/).

**typosquatting** — Strictly, a variant one *typing error* away from the target:
`exmaple.com` for `example.com`. Used loosely for the whole family, it is used
narrowly here — motor error, keyboard geometry, single edits. See
[What typosquatting is]({{ site.baseurl }}/attack/typosquatting/).

**UDRP** — The Uniform Domain-Name Dispute-Resolution Policy, the arbitration
route for taking a domain off a bad-faith registrant. One of four responses to a
confirmed squat, and the slowest. See
[Defending a name]({{ site.baseurl }}/attack/defending/).

## The engine

**addressed form** — A node's or edge's content-addressed encoding: a positional
dag-cbor list, `[type, key, [values...]]` for a node and
`[from, relation, to, [values...]]` for an edge, with every declared prop slot
written in declaration order and `null` for unset. Because the form is a
function of the schema, a node encodes identically wherever it was built. See
[Addressing]({{ site.baseurl }}/internals/addressing/).

**analyzer** — A plugin that runs once over the whole graph after expansion
stops, for any reason including an interrupt, and returns findings. It reads a
restricted `Analysis` surface: nodes, edges, statuses, provenance, plugin
scores, the ledger — never engine belief. See
[Analysis]({{ site.baseurl }}/internals/analysis/).

**applier** — The single writer. Operators never mutate the graph; they return a
delta and `Graph.Apply` canonicalizes keys, admits nodes and edges, enforces the
closure, scope and self-variant invariants, and merges props. Every rule that
must not be opt-in for plugin authors lives here. See
[Operators]({{ site.baseurl }}/internals/operators/).

**barrier** — The end-of-round synchronisation point where every irreversible
decision happens: tree parents are finalized, belief is recomputed for every
node, budgets are enforced, and belief conditions are evaluated. Barrier 0 runs
at seeding, before any dispatch. See [The plan]({{ site.baseurl }}/internals/plan/).

**belief** — A scalar the execution model attaches to a node, used only for
frontier ordering, pruning and operator gating. It never contributes to a
reported number, and analyzers cannot read it at all — the `Analysis` surface
simply does not expose it. See [Analysis]({{ site.baseurl }}/internals/analysis/).

**belief model / HMM** — The pluggable `BeliefModel` that produces belief: an
`Initial` prior for the seed and a `Step` that pushes the parent's latent state
through the relation that admitted the node and conditions on its props — the
forward filtering step of a hidden Markov model. The shipped default is uniform,
which reduces expansion to unranked breadth-first. See
[Analysis]({{ site.baseurl }}/internals/analysis/).

**canonicalization** — The per-type function that normalizes a raw key into the
single form the graph converges on; returning an error refuses the candidate.
It runs before every admission decision, so identity, the ledger denylist and
deduplication all compare canonical keys. See
[The graph]({{ site.baseurl }}/internals/graph/).

**capability (Nameable / Observed)** — What may be done to a node type.
`Nameable` types can root variant generation; `Observed` types are only ever
discovered. Nameable is necessary but not sufficient: eligibility also requires
seed-closure membership and, if set, scope. See
[Types]({{ site.baseurl }}/internals/types/).

**CID** — The content identifier of a stored block: CIDv1, dag-cbor,
sha2-256. Identical scans produce identical CIDs, which is what makes "what
changed since last week" a hash comparison rather than a database query. See
[Addressing]({{ site.baseurl }}/internals/addressing/).

**condition** — An extra `Where` requirement on a matched node — `HasProp`,
`HasEdge`, `InClosure`, `BeliefAbove`. Conditions are data conditions, never
producer dependencies: "there is an IP", not "the ip operator has run". Each
declares what it reads, so it lands in the read-set automatically. See
[Operators]({{ site.baseurl }}/internals/operators/).

**delta** — What one operator returns: nodes, edges and prop assertions, named
by raw key. Deltas are additive only — nothing removes a node, edge or prop —
which makes the graph monotonic within a run and a delta safely replayable. See
[The graph]({{ site.baseurl }}/internals/graph/).

**depth** — A node's shortest *observation* distance from the seed. Structural
and variant edges cost nothing, so decomposing a composite target does not spend
the budget. `MaxDepth` of 0 means unbounded, not zero hops. See
[Observation and depth]({{ site.baseurl }}/guide/observing/) and
[Limits]({{ site.baseurl }}/internals/limits/).

**dispatch** — Running a round's eligible pairs. Calls are concurrent, bounded
by `Workers` and by resource-class rate limits, but nothing is applied during
dispatch: results are collected and applied afterwards in the deterministic work
order. See [The plan]({{ site.baseurl }}/internals/plan/).

**edge** — An admitted typed relation between two nodes, identified by
`hash(from, relation, to)` and carrying its own positional props. See
[The graph]({{ site.baseurl }}/internals/graph/).

**edge class (structural / variant / observation)** — A relation's class, which
decides depth accounting and closure propagation. `Structural` comes from
parsing the target string alone and carries closure membership; `Variant`
connects an origin to a generated variation; `Observation` required a lookup
against an external service and is the only class that costs depth. The dividing
line is whether producing the edge needed a network call. See
[Types]({{ site.baseurl }}/internals/types/).

**effects** — An operator's `Emits()` declaration of the node types, relations
and props it *may* produce. It is a may, not a must; plan compilation uses it
for reachability pruning and dead-operator detection, and declaring `VARIANT_OF`
in it is what makes an operator a variant operator. See
[Operators]({{ site.baseurl }}/internals/operators/).

**finding** — An analyzer's conclusion: a kind, a severity, the admitted nodes
and declined ledger rows it concerns, a summary and its evidence. A node's risk
is the maximum severity among findings referencing it, and the only user-facing
score in the system. See [Analysis]({{ site.baseurl }}/internals/analysis/).

**frontier** — The admitted candidates not yet expanded, which belief orders.
`Limits.Frontier` (`--frontier`) is specified as a cap on candidates admitted
per round and is hashed into the plan, but the scheduler does not enforce it
today, so no ledger row currently cites the `frontier` reason. See
[Limits]({{ site.baseurl }}/internals/limits/).

**node** — An admitted graph entity: a type, a canonical key and its props. It
carries no provenance, status or findings — those live in side tables, so that
two identical scans produce identical content addresses. See
[The graph]({{ site.baseurl }}/internals/graph/).

**NodeID** — A node's stable identity, `hash(type, canonical key)`, fixed for
its whole life. It is deliberately not the content address: props accumulate, so
the CID moves while the identity does not. The scheduler, seen-set, cache, edges
and side tables all key on it. See [The graph]({{ site.baseurl }}/internals/graph/).

**operator** — A unit of work bound to a pattern in the graph: an id, a version,
a trigger, its effects, a resource class, and an `Exec` that returns a delta and
an outcome. Operators are the only place external lookups happen and the only
thing the scheduler dispatches. See
[Operators]({{ site.baseurl }}/internals/operators/).

**outcome / status (ok / empty / failed / timeout / skipped)** — An operator's
own judgement of what happened, recorded per (node, operator) pair. `ok` learned
something positive, `empty` authoritatively determined absence, `failed` means
the lookup broke, `timeout` means nothing was learned; only `failed` and
`timeout` are retried. `skipped` means never attempted and is deliberately *not*
terminal — recording it terminally would make the first gate permanent. See
[Operators]({{ site.baseurl }}/internals/operators/).

**plan** — The compiled, inspectable, pinnable answer to "what will this do?":
the seed and scope, the limits, and every operator binding with its trigger,
reads, effects, config and model CIDs. Execution dispatches only what the plan
selects. See [The plan]({{ site.baseurl }}/internals/plan/).

**plan hash** — The SHA-256 over everything in the plan: seed, scope, limits,
engine model, and each binding's id, version, resource, trigger, reads, effects,
config and plugin model CIDs. Reloading a pinned plan re-derives it and refuses
a mismatch. See [The plan]({{ site.baseurl }}/internals/plan/).

**prop** — One typed field on a node or edge, stored positionally in its type's
declared field order rather than by name. Fields are append-only and removal is
by tombstone, because a position is part of every content address already
written. Competing assertions resolve by the field's merge policy, never by
arrival order. See [The graph]({{ site.baseurl }}/internals/graph/).

**reads / read-set digest** — `Reads` declares the props and relations an
operator consumes, merged with whatever its conditions inspect. It scopes the
`View` — an undeclared field reads as unset — and it is hashed into the read-set
digest that decides re-dispatch and cache validity. Undeclared data is invisible
in both senses, which is why the two must come from one declaration. See
[Operators]({{ site.baseurl }}/internals/operators/).

**rejection** — An invariant violation the applier refused: unknown type,
relation or field, a kind mismatch, a failed canonicalization, a variant edge
from outside the closure or scope, or a self-variant. Unlike a ledger row it
does not deny the candidate forever — the same node may be admitted legitimately
by another edge. See [Limits]({{ site.baseurl }}/internals/limits/).

**relation** — A registered edge type: a name, a class, a version and its own
ordered field list, encoded exactly like a node's. Direction follows the design
convention that an edge runs from the composite to the part. See
[Types]({{ site.baseurl }}/internals/types/).

**resource class** — The rate-limit bucket an operator declares with
`Resource()` — `dns`, `whois`, `npm`. A single global delay is meaningless when
one run talks to five services, so the limiter keeps a minimum interval per
class. Declaring one also marks the operator as an *observer*, whose status
counts toward existence. See [Operators]({{ site.baseurl }}/internals/operators/).

**revision cap** — `Limits.Revisions` (default 3), the number of times one
(node, operator) pair may re-run across the whole scan. Only an actual execution
counts against it; a pair that was gated off never ran, so there is nothing to
revise. See [The plan]({{ site.baseurl }}/internals/plan/).

**round** — One expansion cycle: collect eligible pairs, dispatch them, apply
every result in work order, then barrier. Rounds are the unit of synchronisation
— nodes at depth 1 and depth 5 can be worked in the same one — and `MaxRounds`
(default 64) is the backstop against a type flow that never converges. See
[The plan]({{ site.baseurl }}/internals/plan/).

**SCC condensation** — Collapsing each strongly connected component of the type
flow into one vertex, which yields a DAG that can be layered. It is what lets
`--explain` render the plan in readable layers without pretending the cycles are
not there; cyclic components are marked as such. See
[The plan]({{ site.baseurl }}/internals/plan/).

**scope** — The CLI positional restricting which node types may root a variant,
for example varying only the domain part of an email target. An empty scope
means every Nameable type in the seed closure. It is enforced by the applier as
a rejection, not merely at dispatch. See
[Targets and scope]({{ site.baseurl }}/guide/targets/).

**seed closure** — The seed plus everything reachable from it by *structural*
edges. Only members may root variant generation, which is what stops a
reverse-PTR domain or a manifest-derived package from becoming a new variant
root. See [Limits]({{ site.baseurl }}/internals/limits/).

**seen-set** — The record of (node, operator, read-set digest) triples already
run. A pair is not re-dispatched for the same read-set whatever its outcome; if
something it declared it reads later changes, the digest changes and the pair
becomes eligible again on its own. See
[Limits]({{ site.baseurl }}/internals/limits/).

**selector** — The `On` half of a trigger: which nodes an operator binds to, by
type name or by capability. Binding by capability is what lets one omission
algorithm cover every Nameable type instead of being registered once per type.
See [Operators]({{ site.baseurl }}/internals/operators/).

**severity** — An ordered, named level on a finding: `info`, `low`, `medium`,
`high`, `critical`. It is deliberately not a bare integer, so `--fail-on high`
means the same thing in every release. See
[Analysis]({{ site.baseurl }}/internals/analysis/).

**side table** — Everything held about a node or edge that is deliberately kept
out of its content address: provenance, per-pair status, depth and closure
membership, plugin scores, findings, the truncation ledger. It is stored as one
block linked from the scan root. See
[Addressing]({{ site.baseurl }}/internals/addressing/).

**terminal variant rule** — A node reached by `VARIANT_OF` is never handed back
to a variant operator. Without it `example.com` yields `exmaple.com` yields
`exmalpe.com` and the tool generates the string space; with it, a variant edge
can cost zero depth without unbounding anything. See
[Limits]({{ site.baseurl }}/internals/limits/).

**three-state existence (live / absent / unknown / untried)** — The rollup of a
node's *observation* statuses: `live` if any observer returned ok, `absent` if
none did and at least one determined absence, `unknown` if every attempt failed,
timed out or was skipped, and `untried` if no observer ran at all. "We asked, it
is not there" and "we could not tell" are opposite conclusions, and collapsing
them turns a broken resolver into a clean bill of health. See
[Observation and depth]({{ site.baseurl }}/guide/observing/).

**trigger** — When an operator runs: a selector, zero or more conditions, and
its declared reads. It is matched against graph data only; belief conditions are
the one part evaluated at barriers rather than on every re-dispatch. See
[Operators]({{ site.baseurl }}/internals/operators/).

**truncation ledger** — The record of candidates the engine declined to admit —
by belief, budget, frontier, round cap or deadline — with the type, canonical
key, depth, belief and reason. It is reported like any other section, because a
truncated graph that reads as complete is a correctness bug, and it doubles as a
denylist so a later operator cannot resurrect a declined candidate. See
[Limits]({{ site.baseurl }}/internals/limits/).

**type flow** — The graph of "an operator matching type A may emit type B",
built during plan compilation. Reachability from the seed type decides which
operators survive; it is genuinely cyclic (`domain → ip → PTR → domain`), so
pruning is transitive closure rather than a topological walk. See
[The plan]({{ site.baseurl }}/internals/plan/).
