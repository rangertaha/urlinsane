---
title: 20 · Extending URLInsane
parent: Part III · Inside the engine
nav_order: 8
---

# Extending URLInsane
{: .no_toc }

- TOC
{:toc}

There are three registries and no more, because the graph engine has three
places where behaviour can be added: an operator, an analyzer, and a variant
algorithm. Languages and keyboards are not on that list — they are **data**, and
adding one is adding files, not code. Report formats are a closed set the report
projects into.

| Adding | Where it goes | What registers it |
|---|---|---|
| A variant algorithm | `internal/plugins/variant/<id>/` | a line in `variant/all`, or `plugins.AddSpec` |
| An observation operator | `internal/plugins/observe/<name>/` | a line in `observe/all`, or `plugins.AddOperator` |
| An analyzer | `internal/plugins/analyze/<name>/` | a line in `analyze/all`, or `plugins.AddAnalyzer` |
| A language | `datasets/languages/<code>/*.lst` | `make dataset` |
| A keyboard layout | `pkg/kb/data/` | `go generate ./pkg/kb` |

Every family is a library plus its plugins, and the composition lives in
`<family>/all`. That separation is load-bearing rather than tidy: a family
package that both held the shared library and listed its plugins would make
every plugin import its siblings through it, which is a cycle the compiler
rejects.

## 1. A variant algorithm

Almost every algorithm is a pure `string -> []string` function over the entity's
name, so they are declared as data and share one adapter rather than being
thirty hand-written operator structs:

```go
// Generate turns one name into its variations. Implementations must be pure —
// same input, same output — because the scheduler caches an operator's result
// against a digest of its declared reads and will not call it twice for the
// same input.
type Generate func(name string) []string
```

The declaration is a `variant.Spec`. One directory per algorithm, one exported
`Spec` function, nothing else. This is the whole of
`internal/plugins/variant/co/co.go`:

```go
// Package co is the character omission algorithm.
package co

// Spec declares the algorithm.
func Spec() variant.Spec {
	return variant.Spec{
		ID: "co", Title: "Character Omission", Version: 1,
		Gen: typo.CharacterOmission,
	}
}
```

An algorithm that runs over data takes it as a parameter rather than reading a
global, and folds itself over the set with one of the two combinators. This is
`variant/vs/vs.go`, the vowel swap:

```go
func Spec(langs []variant.Language) variant.Spec {
	return variant.Spec{
		ID: "vs", Title: "Vowel Swapping", Version: 1,
		Gen: variant.OverLanguages(langs, func(l variant.Language, name string) []string {
			return typo.VowelSwapping(name, l.Vowels()...)
		}),
	}
}
```

`OverKeyboards` is the same shape for `[]*kb.Layout`. Both snapshot their input
at construction, because resolving it per call would make an operator's output
depend on registration order at call time — which the scheduler's cache assumes
cannot happen.

### The four fields that change behaviour

```go
type Spec struct {
	ID      string
	Title   string
	Version int
	Types   []string
	Whole   bool
	Gen     Generate
}
```

- **`Types`** narrows the trigger. Leave it empty and the operator binds by
  *capability* — `graph.Selector{Caps: []graph.Capability{graph.Nameable}}` —
  so one algorithm covers domain, email, package, repo and username at once.
  Name types only when the algorithm is genuinely type-specific: `tld` declares
  `Types: []string{variant.TypeDomain}` because swapping a public suffix is
  meaningless for a package name.
- **`Whole`** varies the entire key instead of just the registrable core. The
  default is false and should stay false. Character omission over the whole of
  `example.com` produces `example.cm`, which is a different registry rather than
  a typo of the name — `tld` owns that axis deliberately, and every other
  algorithm would otherwise generate it by accident.
- **`Version`** invalidates cached results. Bump it whenever `Gen`'s output
  changes; leaving it means a stale cached delta stays valid.
- **`ID`** is three things at once: the operator id, the value written to the
  edge's `algorithm` prop, and the short code `--algorithm` selects on.

Order is deliberately not part of the `Generate` contract — several `pkg/typo`
functions accumulate through a map — because the adapter sorts. That sort is
what makes the emitted delta byte-identical across runs.

### Wiring it in

For an algorithm that ships, add the import and one line to
`variant/all/all.go`, which is the only list of them:

```go
specs := []variant.Spec{
	co.Spec(),
	…
	vs.Spec(o.Languages),
	cb.Spec(o.Combos),
}
```

An algorithm added there without a directory, or a directory without a line
there, is exactly the drift the old registry made possible. `All` panics on a
duplicate id rather than silently dropping one, because the scheduler keys its
seen-set and cache on the id.

For an algorithm that lives outside the tree, register it in `init`:

```go
func init() {
	plugins.AddSpec(variant.Spec{ID: "acme", Title: "Acme Squatting", Gen: gen})
}
```

`AddSpec` is sugar over `AddAlgorithm`, which takes a builder and settings
defaults instead — use it when how many algorithms you contribute depends on the
data:

```go
func AddAlgorithm(id string, defaults map[string]any, fn AlgorithmFunc)
type AlgorithmFunc func(Env) ([]variant.Spec, error)
```

The registered id names the *plugin*, not the algorithms: a keyboard plugin
registers once and yields one algorithm per distinct layout. Settings are keyed
on the plugin id; `--algorithm` selects on the `Spec` ids. Registering the same
id twice panics.

### Why declaring `VARIANT_OF` is the whole story

The adapter builds one `graph.Operator` per `Spec`, and its `Emits` says:

```go
return graph.Effects{
	Nodes: append([]string(nil), nodes...),
	Rels:  []string{Rel},          // graph.VariantRel
	Props: []string{PropAlgorithm, PropDistance},
}
```

An operator is a *variant* operator if and only if it declares `VARIANT_OF` in
`Emits().Rels`. Not by naming convention, not by living under `variant/` — by
declaration. That one line is what subjects it to the seed-closure restriction
and the terminal rule from [Limits](../limits/): its edges are refused unless
the origin is in the seed closure, and nothing it produces is ever handed back
to a variant operator.

{: .warning }
> The converse also holds. An observation operator that declares `VARIANT_OF`
> because it "produces variant-looking names" has just made every name it finds
> a variant root, and the closure rule will then refuse most of its edges. If
> the edge is not a generated variation of its origin, it is not `VARIANT_OF`.

### Testing it

Algorithm behaviour is tested in `internal/plugins/variant/all/variant_test.go`
against a registry built in the test file, not against the shipped one. The
tests use tiny substitute datasets — the real lists are 5,000 subdomains and
8,600 suffixes, which would make every test a benchmark:

```go
func testOptions() variant.Options {
	return variant.Options{
		Subdomains: []string{"www", "mail", "login"},
		Suffixes:   []string{"com", "net", "org", "co.uk"},
	}
}
```

The invariants already covered there apply to a new algorithm for free once it
is in `Specs`: every operator emits `VARIANT_OF`, ids are unique and sorted,
every operator has a version, deltas carry raw keys only, variants are
introduced only by their edge, and the same input yields a byte-identical delta.
Add a test for what your algorithm *says*; the structure is checked already.

An algorithm with nothing to say about a name returns `graph.Empty()`, not a
failure — a name with no hyphens has no hyphen omissions. The adapter does this
for you when `Gen` yields nothing.

## 2. An observation operator

An observation operator implements `graph.Operator` in full. It goes in
`internal/plugins/observe/<name>/`, embeds `observe.Base` for the per-call
timeout, and takes its external dependency as an interface parameter so a test
can supply a fake. This is `observe/ptr/ptr.go` entire, minus the shared struct:

```go
func newReverse(o observe.Options, r observe.Resolver) graph.Operator {
	return reverse{dnsOp{
		Base: o.Base(), id: "ptr", ver: 1, on: observe.TypeIP, res: r,
		eff: graph.Effects{Nodes: []string{observe.TypeDomain}, Rels: []string{observe.RelPTRTo}},
	}}
}

func (o reverse) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.Call(ctx)
	defer cancel()

	names, err := o.res.LookupAddr(ctx, v.Key())

	var d graph.Delta
	from := v.Ref()
	for _, n := range names {
		name := observe.Host(n)
		if name == "" {
			continue
		}
		ref := graph.NodeRef{Type: observe.TypeDomain, Key: name}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, graph.EdgeRef{From: from, Rel: observe.RelPTRTo, To: ref})
	}
	return d, observe.DNSOutcome(err, len(d.Nodes) > 0)
}
```

Five things in there are the whole contract.

**Bind to a type, never to a producer.** The trigger is
`graph.Trigger{On: graph.Selector{Types: []string{o.on}}}` and nothing else. A
`Where` clause here would be a producer dependency in disguise — the coupling
the DAG exists to remove — and there is a test in `ptr_test.go` asserting `ptr`
declares no conditions. Bind on `ip` and you run against addresses `dns-a`
found, addresses a future RDAP operator finds, and addresses nobody has written
the producer for yet.

**Derive the context from the one `Exec` was given.** `o.Call(ctx)` applies the
per-operator timeout on top of the caller's deadline. An earlier version rooted
this at `context.Background()`, which silently discarded every cancellation the
engine tried to deliver: Ctrl-C then waited out an in-flight whois instead of
stopping at the round boundary.

**Emit raw keys.** `NodeRef` and `EdgeRef` carry the key as the service returned
it. Canonicalization belongs to the registry and the applier; an operator that
minted a `NodeID` would quietly break convergence.

**`Resource()` does two jobs.** It names the rate-limit class — `dns`, `whois`,
`http`, `geo`, or `""` — and each class gets its own token bucket, so the
interval protecting whois does not throttle DNS to the same crawl. It is *also*
how the scheduler decides which operators count as existence observers:

```go
// Tell the graph which operators actually look something up, so Existence
// does not read a decomposer's "I parsed this" as "this exists" (§9).
var observers []string
for _, o := range sorted {
	if o.Resource() != "" {
		observers = append(observers, o.Id())
	}
}
g.SetObservers(observers)
```

{: .warning }
> A pure operator that returns `""` can never contribute to a node reading as
> live, whatever status it records — and an operator that names a resource it
> does not actually call will make every name it parses read as live. This
> failure has already happened once: without the observer set, `-google.com` and
> `'oogle.com` were reported as live typosquats on the strength of having been
> successfully parsed.

**Return the right kind of nothing.** `DNSOutcome` is the judgement the whole
package exists to make:

| Condition | Outcome | Meaning |
|---|---|---|
| records came back | `graph.OK()` | something positive was learned |
| NXDOMAIN, or NODATA | `graph.Empty()` | authoritatively determined absence |
| SERVFAIL, REFUSED, malformed | `graph.Failed(err)` | the lookup itself broke |
| deadline, cancellation | `graph.Timeout(err)` | nothing was learned |

`StatusEmpty` is an answer and is never retried; `Failed` and `Timeout` are the
absence of one and are retriable. Collapsing "confirmed absent" into "could not
determine" would make a run full of timeouts read as a run full of free names.

### New props need registering, and only ever by appending

An operator may only assert fields the registry declares. A prop that was never
registered is rejected as `RejectUnknownField` — and rejected *quietly*, because
the applier records a rejection rather than failing the run, so the symptom is a
column that is silently always empty.

Shipped observation props are declared in `observe/schema.go` and merged into
the schema by `scan.Registry()`, which passes them to `decompose.Register` as an
`Extension`:

```go
func Fields() map[string][]graph.FieldDef {
	return map[string][]graph.FieldDef{
		TypeDomain: {
			{Name: FieldUnicode, Kind: graph.KindString},
			…
			{Name: FieldRegistrar, Kind: graph.KindString, Merge: graph.Precedence("rdap", "whois")},
		},
	}
}
```

{: .warning }
> Append to the end of the list. Never reorder, never delete. A field's position
> is its stable index and part of the content address of every node already in
> the store, so appending is safe and reordering silently corrupts diffs rather
> than failing loudly. Bump the type's version when you append; see
> [Types and relations](../types/).

The same rule governs the order extensions are passed in, which is why
`scan.Registry()` passes them as a fixed literal and never from a map or from
whichever operators the plan happened to select.

{: .todo }
> There is no hook for a registered plugin to contribute schema fields:
> `scan.Registry()` takes `observe.Fields()` and `observe.RelFields()` and
> nothing else. An out-of-tree operator today can assert only fields that are
> already registered, or emit relations rather than props.

### Wiring it in

In-tree, add it to `observe/all/all.go`, which builds every operator the
supplied options can support:

```go
ops := []graph.Operator{idn.New(o)}
ops = append(ops, dns.New(o, res)...)
ops = append(ops, ptr.New(o, res), whois.New(o, who))
if o.Geo != nil {
	ops = append(ops, geo.New(o, o.Geo))
}
```

That `if` is the pattern for a missing dependency: **omit the operator rather
than register one that can only fail.** Plan compilation reports what may run,
and listing an operator guaranteed to error would make `--explain` lie. Out of
tree, the same discipline applies — an `OperatorFunc` returning no operators is
how a plugin declines:

```go
func init() {
	plugins.AddOperator("acme", map[string]any{"timeout": 5},
		func(e plugins.Env) ([]graph.Operator, error) {
			if e.Observe.Prober == nil {
				return nil, nil
			}
			return []graph.Operator{&acmeOp{
				timeout: e.Settings.Int("timeout", 5),
				prober:  e.Observe.Prober,
			}}, nil
		})
}
```

`Env` carries the same resolver, prober, geolocator and datasets the shipped
operators get — including the fakes a test injects — because an operator must be
a pure function of its inputs to be testable as one. A plugin that reached for a
package-level resolver could not be given a fake, and two concurrent tests would
fight over it.

### Testing it

`internal/plugins/observe/observetest` holds the fakes and graph assertions the
observation tests share — a package rather than a `_test.go` file, so that a
fake resolver copied into five operator packages cannot drift:

```go
func TestReverseBindsToIPAndEmitsDomain(t *testing.T) {
	res := &observetest.FakeResolver{Addr: map[string][]string{
		"93.184.216.34": {"example.com.", "www.example.com."},
	}}
	g := observetest.Run(t, observe.TypeIP, "93.184.216.34", newReverse(observe.Options{}, res))

	for _, want := range []string{"example.com", "www.example.com"} {
		if !observetest.HasNode(g, observe.TypeDomain, want) { … }
		if !observetest.HasEdge(g, observe.TypeIP, "93.184.216.34", observe.RelPTRTo,
			observe.TypeDomain, want) { … }
	}
	observetest.WantStatus(t, g, observe.TypeIP, "93.184.216.34", "ptr", graph.StatusOK)
}
```

`Run` seeds a graph, dispatches your operators through the real scheduler and
returns the graph. `HasNode`, `HasEdge`, `Prop`, `EdgeProp`, `WantStatus` and
`Dump` assert against it. `Nxdomain()`, `Servfail()` and `DNSTimeout()` produce
the three error shapes, so the Empty/Failed/Timeout split is testable without a
network. `TestRegistry` performs the same schema merge the real registrar does,
which proves the append works — the whole contract between an operator package
and `decompose`.

## 3. An analyzer

Analyzers run once, over the finished graph, after expansion stops — for any
reason, including an interrupt. There is exactly one lifetime.

```go
type Analyzer interface {
	Id() string
	Exec(ctx context.Context, a *Analysis) ([]Finding, error)
}
```

`Analysis` is a read-only surface: nodes, edges, `Existence`, `Depth`,
`InClosure`, `Incoming`, `Outgoing`, provenance, the truncation ledger and the
rejections. It deliberately does **not** expose the engine's belief. Withholding
it is what makes "the execution model never contributes to a reported number"
true by construction rather than by convention: an analyzer that could read
belief could launder it into a `Severity` and no reviewer would spot it.

`analyze/scoring` is the shape to copy:

```go
type Scoring struct{}

func (Scoring) Id() string { return "scoring" }

func (Scoring) Exec(_ context.Context, a *graph.Analysis) ([]graph.Finding, error) {
	var out []graph.Finding
	for _, n := range a.Nodes() {
		if !analyze.IsVariant(a, n.ID) {
			continue
		}
		if a.Existence(n.ID) != graph.Live {
			continue
		}
		…
		out = append(out, graph.Finding{
			Kind:     "live-variant",
			Severity: sev,
			Nodes:    []graph.NodeID{n.ID},
			Summary:  fmt.Sprintf("%s is live: %v", n.Key, signals),
		})
	}
	return out, nil
}
```

Two helpers in `internal/plugins/analyze` are the shared queries: `IsVariant`
reports whether a node has an inbound `VARIANT_OF` edge, and `EditDistance`
reads the distance *off that edge* rather than recomputing it — the algorithm
already computed it, and recomputing risks disagreeing with the edge.

Severity is a named, ordered level rather than a bare int, because `--fail-on
high` has to mean the same thing in every release. Findings may also name
declined candidates by `LedgerRef`, which is how an analysis about something
that *does not* exist — dependency confusion — says what it means.

In-tree, add it to `analyze/all/all.go`, which returns the shipped analyzers in
run order. Out of tree:

```go
func init() {
	plugins.AddAnalyzer("acme", nil, func(e plugins.Env) (graph.Analyzer, error) {
		return acme{}, nil
	})
}
```

Returning a nil analyzer is allowed and drops it. Note that an explicit
`scan.Options.Analyzers` **replaces** both the shipped set and the registered
plugins rather than adding to them: a caller naming its analyzers means those
and no others.

## 4. A language

A language is not a plugin, and adding one involves no Go at all. It was a
plugin once — thirty Go files, each restating its vowels, graphemes, homoglyphs
and misspellings as literals — and every one of them had to be edited, compiled
and shipped to correct a single word.

Make a directory and drop `.lst` files in it:

```
datasets/languages/<code>/
  vowel.lst        one token per line
  grapheme.lst     one token per line
  numeral.lst      1 one first
  homoglyph.lst    a à á â ã ä å ɑ а ạ ǎ ă ȧ ӓ ٨
  homophone.lst    groups of words that sound alike
  misspelling.lst  hwile while
```

The six names above are the ones `internal/dataset/lang.go` reads, one per
`Language` method — `vowel.lst` becomes `Vowels()`, `homoglyph.lst` becomes
`Homoglyphs()`, and so on. Any other `.lst` in the directory is imported too,
under its own dataset name, and English ships several: `synonym`, `antonym`,
`word`, `stopword`, `token`, `positive`, `negative`.

Every file is read the same way. A line is its words, runs of whitespace
collapse, a leading `#` is a comment. A line of two or more words associates
them with each other; a single-word line joins the vocabulary and associates
with nothing. That is the difference between the language datasets, whose lines
are groups, and the list-shaped corpora under `domains/`, `entities/` and
`packages/`, which are plain vocabularies. Association counts become weighted
transitions, so a pair seen in several groups outranks one seen once.

Then rebuild the embedded database:

```sh
make dataset       # go run ./cmd/datasets build datasets
```

The output is `internal/config/dataset.db`, which `//go:embed` picks up — so
this must run *before* `make build` for a data change to reach the binary. The
build deletes and recreates the database rather than migrating it, which is the
point of a build: migrating through three schema generations leaves the columns
of all three.

`--list languages` and every language-driven algorithm pick the new language up
with no further change, because `variant.RegisteredLanguages()` queries the
dataset and the algorithms fold over whatever it returns.

One file takes a different route. `synonym.lst` is the combosquatting
vocabulary, and `datasets/synonyms.go` embeds it straight from the repo tree
with `//go:embed languages/*/synonym.lst`, so a change there needs a rebuild
rather than `make dataset`. Despite the name these are not thesaurus entries:
they are the brand-adjacent vocabulary an attacker bolts onto a target name —
login, verify, pay, invoice, delivery — grouped by theme. Adding one does not
widen what `cb` uses, either: `Options.Combos` defaults to English alone and is
deliberately not derived from `Options.Languages`, because combo vocabulary is a
property of the target's audience rather than of the name. Pass
`variant.ComboKeywords(langs)` to opt into the multi-language form.

The `Language` table itself is seeded from `pkg/kb`, not from the dataset
directories: a language the tool can reason about has a keyboard whether or not
anyone has curated a word list for it yet, and seeding from
`datasets/languages/` would make a language blink out of `--list languages` the
moment its corpus emptied. A curated language with no layout still gets its row
from the import that follows.

{: .warning }
> Three curated directories use a code `kb` does not — `iw` (kb has `he`), `no`
> (kb has `nb`) and `la` (kb has no Latin layout at all). Each lands as a second
> row for a language `kb` already listed, so homoglyphs get looked up under one
> code while keyboard adjacency uses the other. Pick a code `kb` already knows
> unless you intend that split.

The contents of these files are algorithm *input*, not prose. `misspelling.lst`
is a list of misspellings on purpose, and "correcting" one removes a variant the
tool exists to generate.

## 5. A keyboard layout

Layouts live in `pkg/kb`, which is a standalone library — 203 layouts across 110
languages, taken from [kbdlayout.info](https://kbdlayout.info) and stored as one
protobuf file per driver plus an index, compiled in with `go:embed`.

```sh
go generate ./pkg/kb
```

That runs the two directives in `pkg/kb/kb.go`:

```go
//go:generate protoc --proto_path=internal/kbpb --go_out=. --go_opt=module=github.com/rangertaha/urlinsane/pkg/kb internal/kbpb/kb.proto
//go:generate go run ./gen -out data
```

`protoc` regenerates the Go bindings from `internal/kbpb/kb.proto`; `gen` walks
the KLID list on kbdlayout.info, follows each KLID to the driver that serves it,
and converts the driver's XML export. Responses are cached under `pkg/kb/.cache`
so a re-run after a parsing change costs no requests. Encoding is deterministic,
so a rebuild that changed nothing leaves the dataset byte for byte identical —
and because the files are opaque to a diff, every layout is decoded and compared
before it is written.

{: .note }
> Regenerating needs `protoc` and `protoc-gen-go` on `PATH`. Building and using
> the package does not, and neither does building URLInsane.

Key positions are not stored. They are recomputed from the scan code and the
layout's form when a layout loads, which keeps the geometry in one table instead
of copying it into every file — and it is why adjacency respects the stagger
between rows rather than treating the board as equal-width columns.

A new layout reaches the algorithms only if it types differently from every
layout already shipped:

```go
// RegisteredKeyboards returns the layouts the keyboard-driven algorithms run
// over: every layout pkg/kb ships, one per distinct Adjacency behaviour.
```

The 203 layouts collapse to about 30 distinct neighbour sets, because most Latin
boards share QWERTY geometry. `adjacencySignature` identifies a layout by how it
types over `a-z0-9`, not by what it is called, and the first layout of each
behaviour wins in id order. Running all 203 would do seven times the work for
the same candidates.

## Getting a plugin into a build

The shipped operators, algorithms and analyzers are composed directly by
`internal/scan` through `decompose/all`, `variant/all`, `observe/all` and
`analyze/all`. They are not registered, because a registry entry for what always
runs would be a second place to look for it.

Everything that does *not* always run is linked in by a blank import of
`internal/plugins/all`, which `cmd/urlinsane/typo.go` already carries:

```go
_ "github.com/rangertaha/urlinsane/internal/plugins/all"
```

That package is empty today, and that is a fact worth stating rather than a file
worth deleting: without it a new plugin has no obvious place to be linked from,
and the first one would be wired into `cmd/urlinsane` by hand, where the second
would not find it.

### Settings

A plugin declares its defaults at registration and never has to be configured:

```go
plugins.AddOperator("acme", map[string]any{"timeout": 5, "retries": 2}, build)
```

The defaults reach `~/.config/urlinsane/config.yaml`, so its `plugins:` section
lists what *can* be configured rather than leaving the user to discover it. What
a plugin receives is its defaults with the file's overrides applied **per key**,
so overriding one setting does not mean restating the rest. The accessors —
`Settings.String`, `Int`, `Bool`, `Strings` — coerce rather than error, and
accept the `float64` and `[]any` that YAML produces for the same literals,
because a plugin that had to validate every field would either duplicate that
logic or ignore it.

Unknown keys in the file are kept rather than dropped: a setting belonging to a
plugin this build does not link in is not necessarily a mistake.

## Gotchas, in one place

**Deltas are additive only.** Nothing removes a node, edge or prop. An operator
that wants to correct an earlier assertion asserts again and lets the field's
merge policy decide — `graph.Precedence("rdap", "whois")` is how two sources of
the same fact are reconciled by declaration rather than by whichever answered
first.

**Fields are append-only.** Appending bumps the version and leaves existing CIDs
untouched. Reordering or deleting changes the meaning of every content address
already stored.

**`Emits()` must be honest.** `--explain` builds the type-flow graph from
declared effects, so an operator that emits something it did not declare is
invisible to the plan — and declaring `VARIANT_OF` is what makes it a variant
operator, with everything that follows.

**`Reads` feeds the cache key.** An operator sees a `View` scoped to what it
declared, and the same declaration is hashed into the read-set digest that
`(node, operator, digest)` keys on. Conditions contribute to that digest
automatically, so a `Where: HasProp("live")` adds `live` even if `Reads` does
not mention it. Under-declaring means the operator never re-runs when the data
changes; the variant adapter declares nothing at all, which is correct, because
nothing a collector later learns changes what a string mutation produces.

**Ids must be unique and are sorted.** Registration order is package
initialisation order, and letting it through would give the same run a different
plan hash between builds. Both the registry and `variant/all` refuse a duplicate
loudly, because the second registration silently winning is the kind of thing
found much later.

**A plugin cannot reorder the run.** There is deliberately no way to say "after
X". The pipeline's `DependsOn` list is what made plugin order load-bearing and
its cache unsound.

{: .todo }
> `plugins.IDs()` exists and is documented as backing `--list plugins`, but
> `listTopic` registers no such topic and nothing calls it. Two related gaps:
> `--list algorithms` renders `variantall.Specs(variant.Options{})`, so
> plugin-contributed algorithms are missing from it — they do appear in
> `--list operators`, which goes through `scan.Operators`.

## What is not a plugin

Report formats are a closed set: `table`, `json`, `ndjson`, `csv`, `dot`, all
projected from one intermediate `report.Report` built once. Adding a sixth means
a renderer package under `report/`, a case in `report/all.Write` and an entry in
`report.Formats()`. Five renderers reading the graph directly would drift, and
"the formats disagree about what the scan found" is not a defect a user can work
around.

Node types and relations are not plugins either. They are the schema, registered
once by `decompose.Register` against a fresh registry — see
[Types and relations](../types/).

---

Reference material — the full CLI, the design document, the `pkg/kb` model, the
glossary and the bibliography — is in the
[Reference section]({{ site.baseurl }}/reference/).
