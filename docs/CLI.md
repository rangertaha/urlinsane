---
title: CLI reference
parent: Reference
nav_order: 1
---

# CLI

The command-line reference for `urlinsane`. This is the *what*; `DESIGN.md` §12
is the *why* — invocation grammar, flag scoping, progress, interruption and the
reasoning behind each choice. Where the two overlap, this page cites §12 rather
than restating it.

**Status marks what exists.** The engine is mid-rewrite. `typo` is largely built;
much of the interaction design in §12.3–§12.7 is specified and not yet
implemented. Anything marked ◦ will not work today.

| Mark | Meaning |
|---|---|
| ✓ | implemented |
| ◦ | specified in `DESIGN.md`, not yet built |

Two binaries ship: `urlinsane`, the scanner, and `datasets`, a maintainer tool
for building the reference data the scanner reads.

## 1. Commands

```
urlinsane typo   [<scope>] <target> [flags]   scan                          ✓
urlinsane report <target> [flags]             render a saved scan of it     ✓
```

A bare `urlinsane` prints help and exits `0`; the root does not scan.

**A verb has to do something a flag on `typo` could not.** That is the test, and
`report` passes it: it runs without scanning at all. It takes a target, as
`typo` does, but nothing else — no scope, no `--depth`, no algorithm or
collector selection — because every one of those was decided when the scan ran.
As a flag it would sit on a command whose entire remaining surface it rejects.

What fails the test is the work done *around* a scan: compiling the plan without
running it, listing what is registered, asking how a result was reached. Those
take the same target, the same scope and the same selection flags as the scan
itself, so they are flags on `typo` — `--explain` ✓, `--list TOPIC` ✓, and
`--why NAME` ◦. A verb would have to re-accept that whole surface to say
anything useful, and the two spellings would drift apart.

Training is not in this binary at all. It is a maintainer operation over stored
graphs, which is what `datasets` is for (§8).

**The target's kind is detected from the string alone:**

```bash
urlinsane typo acme.com                 # domain
urlinsane typo bob@acme.com             # email: varies bob, acme.com, and the address
urlinsane typo npm:lodash               # package on a named registry
urlinsane typo github.com/acme/tool     # repo
urlinsane typo bobsmith                 # username
```

Detection is ordered: `@` with a host after it is an email; an owner/name triple
or a git remote is a repo; `registry:name` is a package; a dotted name whose
suffix is in the public suffix list is a domain; anything else that is a legal
handle is a username. That last rule separates `lodash` from `lodash.com`.

### 1.1 The scope positional

The optional first positional narrows **what gets varied**, and never changes
how the target is read:

```bash
urlinsane typo username acme.com/bob    # vary only bob
urlinsane typo domain bob@acme.com      # vary only acme.com
urlinsane typo username,domain bob@acme.com
```

Omitted, scope is every `Nameable` node in the seed closure — *including the
composite node itself*, so a bare email also yields whole-address variants.

**There is no `--scope` flag; the positional is the only spelling.** With two
positionals the first must name registered `Nameable` types; with one, it is the
target. Three or more is an error.

### 1.2 `report` — render a saved scan ✓

```bash
urlinsane typo acme.com --save-graph        # scan once, keep the graph
urlinsane report acme.com                   # the table again, no network
urlinsane report acme.com -o json --save r.json
urlinsane report acme.com -f risk\>high     # re-filter without re-scanning
urlinsane report acme.com --at bafyrei…     # an older scan of the same target
```

**You name the target, not the file.** `acme.com` is what you scanned and what
you remember; a root CID is not. Targets are matched exactly as `typo` detected
them (§1), so `report bob@acme.com` finds the email scan and not the scan of
`acme.com` inside it.

Underneath, the store is content-addressed (`internal/store`) and a scan's
root CID *is* its name — two runs that produce the same CID are the same scan by
construction. `--at CID` names one exactly, for when a target has been scanned
more than once; bare, `report` takes the most recent. Blocks live under
`~/.config/urlinsane/blocks`; `--store DIR` reads another.

**Rendering never re-observes.** `report` replays the stored blocks through the
same applier the scan used, then re-encodes every node and edge and checks the
CIDs against what was stored — so the graph it renders is byte-identical to the
one that was saved, or it errors. Re-observing instead would repeat network work
and could disagree with the run being reported, which is the one thing a report
of a past scan must not do.

This is why filters belong here. `--filter` was always a display filter over
built rows (§4), so re-filtering is a re-render, not a re-scan: a scan that took
ten minutes can be sliced repeatedly for nothing.

| Flag | Alias | Default | Description | |
|---|---|---|---|---|
| `--output` | `-o` | `table` | as `typo` | ✓ |
| `--save` | | | as `typo` | ✓ |
| `--filter` | `-f` | | as `typo` | ✓ |
| `--verbose` | `-v` | | include provenance and engine belief | ✓ |
| `--ledger` | | | candidates declined by a budget or bound | ◦ |
| `--at` | | latest | render this root `CID` rather than the newest scan | ✓ |
| `--scans` | | | every saved scan of this target, newest first | ✓ |
| `--store` | | `~/.config/urlinsane/blocks` | blockstore to read | ◦ |

The scan-shaping flags — `--depth`, `--algorithm`, `--language`, `--keyboard`,
`--collect` — are **not** accepted. They would read as though they could change
what a finished scan contains, and under §12.2 an inapplicable flag is an error
rather than a no-op.

**The index sits beside the store.** Every stored root already records its seed
type and key, so each scan knows its own target — but nothing in the store maps
a target back to its roots, and the root deliberately carries **no timestamp, no
run id and no partial flag**, because a clock reading inside it would make every
re-scan produce a different CID. "Newest first" and "this scan was interrupted"
are therefore not answerable from roots alone.

That is by design, not an oversight: `store.Root` says what a caller
needing the time of a scan should do, which is keep it *beside* the root. So the
index carries, per entry: target, root CID, when it ran, and whether it was
partial. It is `~/.config/urlinsane/scans.json`, written by `--save-graph` and
read by `report`.

Partial-ness comes from that index rather than a flag: a scan interrupted at a
round barrier (§6) is reported partial however often it is re-rendered, because
that is a fact about the scan and not about the rendering.

## 2. `typo` flags

### 2.1 Common

| Flag | Alias | Default | Description | |
|---|---|---|---|---|
| `--depth` | `-d` | `3` | observation hops from the seed | ✓ |
| `--filter` | `-f` | | `live`, absent, unknown, untried, risk>SEV, type=NAME, depth<=N | ✓ |
| `--output` | `-o` | `table` | `table` \| json \| ndjson \| csv \| dot | ✓ |
| `--save` | | | write the report to `PATH`; format from the extension | ✓ |
| `--algorithm` | `-a` | all | restrict variant generation to these algorithm `ID`s | ✓ |
| `--language` | `-l` | all registered | languages the language-driven algorithms run over | ◦ |
| `--keyboard` | `-k` | all registered | keyboard layouts the keyboard-driven algorithms run over | ◦ |
| `--collect` | `-c` | all in plan | restrict observation operators; `^id` excludes | ◦ |
| `--fail-on` | | | exit 2 if any finding reaches `SEV` | ✓ |
| `--verbose` | `-v` | | include provenance and engine belief | ✓ |
| `--quiet` | | | silence stderr; the exit code is the whole answer | ◦ |
| `--tui` | | | interactive view of the same scan (§5.3) | ◦ |
| `--save-graph` | | | persist the graph to the store; prints its root `CID` | ✓ |
| `--why` | | | how `NAME` was reached, from a saved graph (§12.7) | ◦ |
| `--store` | | `~/.config/urlinsane/blocks` | graph store to save into | ◦ |
| `--ledger` | | | list candidates declined by a budget or bound | ◦ |
| `--plan` | | | write the compiled plan, or pin an existing one | ◦ |
| `--model` | | | execution model `NAME[@cid]`; `--<op>.model` per plugin | ◦ |
| `--trace` | | | record tuples for training (§10.4) | ◦ |

`--explain`, `--list` and `--why` suppress the scan rather than adding to it:
each prints its own thing and exits `0`.

Today they are checked in order and `--list` wins over `--explain` silently ✓ —
a precedence rule nobody asked for. Supplying more than one should be an error,
since there is no sensible reading of "explain the plan *and* list the formats"
◦.

`--filter` and `--algorithm` are repeatable *and* comma-split, so
`-f live -f risk>high` and `-f live,risk>high` are the same.

**`--filter` selects rows in the report, never work in the scan.** It runs after
expansion, over built rows — it never makes a scan faster or cheaper. Narrowing
the scan is what `--depth`, `--algorithm`, `--collect` and the scope positional
do. See §4.

### 2.1.1 Selecting what runs, and what it runs over

**`-a` picks the algorithm; `-l` and `-k` pick the data it runs over.** `-a hr`
selects homoglyph replacement; `-l ru` decides whose homoglyphs. They compose
rather than compete (§12.10).

```bash
urlinsane typo acme.com -a hr -l ru,el     # homoglyphs, cyrillic and greek only
urlinsane typo acme.com -c ^whois          # everything except whois
urlinsane typo acme.com -c dns,ptr         # only these observers
urlinsane typo acme.com -a ^cb,^afx        # drop the noisy generators
```

A `^` prefix excludes; a bare list restricts. Mixing the two forms in one flag is
an error rather than a precedence rule to memorise. `^` rather than `-` because
`-whois` is indistinguishable from a flag, and rather than `!` or `~` because
both need shell quoting.

**An unknown id is an error, not an empty selection.** A typo that quietly
scanned with no observers would produce a clean report by doing nothing.

Defaults are everything registered: a tool whose failure mode is missing a
squatted name should not narrow its own recall by default.

Languages and keyboards are no longer plugins, so "all registered" now means
all the *data* present. Keyboards come from `pkg/kb` and are always available:
203 shipped layouts collapse to 133 distinct adjacency sets, and it is the
distinct sets the algorithms run over. Languages come from the dataset
database, which `typo` opens on startup: the shipped one carries 113, and the
language-driven algorithms run over all of them. What is missing is `-l` and
`-k` themselves — neither flag is registered, so the sets cannot be narrowed
(§9).

### 2.2 Advanced

Registered but hidden from `--help`: most runs never touch them, and a help
screen that is a union of every knob is one nobody reads (§12.1). The tier that
would reveal them, `--help-all`, is not built ◦.

| Flag | Default | Description | |
|---|---|---|---|
| `--rounds` | `64` | backstop for a type flow that never converges | ✓ |
| `--workers` | `1` | concurrent operator calls | ✓ |
| `--budget` | `0` | global admitted-node cap; 0 means unbounded | ✓ |
| `--frontier` | | cap on candidates admitted per round | ✓ |
| `--attempts` | `2` | per-pair attempts within a round | ✓ |
| `--timeout` | | bound on a single operator call | ✓ |
| `--no-color` | | disable ANSI styling | ✓ |

Defaults come from `graph.Limits.withDefaults`; an unset flag takes the engine
default rather than zero.

Under §12.2, operator-owned flags (`--nameservers`, `--npm.registry`) are
declared by the operator that consumes them and appear in `--help` only when
that operator is in the plan; an inapplicable flag is an **error**, not a no-op.
Not built ◦.

### 2.3 Global

| Flag | Alias | Description | |
|---|---|---|---|
| `--version` | `-V` | print the version | ✓ |
| `--debug` | | log debug messages for development | ◦ |

`--debug` is registered but inert — its action is a no-op and nothing reads the
value.

## 3. `--list` topics

| Topic | Aliases | Columns | |
|---|---|---|---|
| `types` | | `NAME  CAPABILITY  VERSION` | ✓ |
| `relations` | `rels` | `NAME  CLASS  DEPTH COST` | ✓ |
| `operators` | `ops` | `ID  VERSION  RESOURCE  BINDS ON  EMITS` | ✓ |
| `algorithms` | `algos` | `ID  TITLE  APPLIES TO` | ✓ |
| `languages` | `langs` | `ID  NAME` — from the dataset database; 113 in the shipped one | ✓ |
| `keyboards` | | `ID  NAME` — the 133 distinct layouts in `pkg/kb` | ✓ |
| `formats` | | one per line | ✓ |
| `filters` | | `FILTER  SELECTS` | ✓ |
| `models` | | available execution models (§10.6) | ◦ |

Unknown topics error and name the known set, rather than printing nothing.

## 4. Filters

`--filter` selects which nodes appear in the report:

| Filter | Selects |
|---|---|
| `live` | an observation operator returned ok |
| `absent` | none did, and at least one determined absence |
| `unknown` | every attempt failed, timed out or was skipped |
| `untried` | no operator ran on it |
| `risk>SEV` | nodes with findings above a severity |
| `type=NAME` | a node type |
| `depth<=N` | observation hops from the seed |

Comparisons take `>`, `>=`, `<`, `<=`, `=`, `==` and `!=`, on any comparable
field — `<=` is not special to `depth`. Severities are `info`, `low`, `medium`,
`high`, `critical`.

**Filters in the same existence family are alternatives; across families they
narrow.** `-f live -f absent` means *either*, because a node cannot be both and
requiring both would always render empty. `-f live -f risk>medium` means *both*.
Treating every filter as a conjunction would make the commonest two-value query
silently return nothing.

This replaced `--registered`/`--unregistered`, two booleans that could be set to
contradict each other and could not express the third case at all.

**A filtered-empty report must explain itself** (§12.6) — an empty table after
`-f live` looks identical whether nothing was squatted or every lookup timed
out, and those are opposite conclusions ◦:

```
  no rows matched --filter live

  1,284 hidden:  1,203 absent · 39 unknown · 42 live but below --depth 2

  ⚠ 39 unknown — every dns attempt timed out. the network, not the result.
```

## 5. Output

**stderr is for the human, stdout is for the machine** (§12.3). Progress,
streamed findings and diagnostics go to stderr and are suppressed when stdout is
not a TTY; the report goes to stdout. So this needs no flag to do the right
thing ◦:

```bash
urlinsane typo acme.com -o json > out.json   # live progress on the terminal,
                                             # clean json in the file
```

`-o/--output` chooses the rendering; `--save PATH` writes a second copy, format
from the extension:

| Extension | Format |
|---|---|
| `.json` `.ndjson` `.csv` `.dot` | itself |
| `.txt` `.text` | `table` |
| `.gv` | `dot` |

Anything else is an error rather than a guess. A saved file is never coloured.

Every format is byte-identical across runs of the same scan, except `ndjson`,
which makes no ordering claim. Colour is suppressed by `--no-color`, by
`NO_COLOR`, and when stdout is not a TTY.

### 5.1 Existence is three-state, and always shown

```
  ● live 42   ○ absent 1,203   ? unknown 39
```

*Absent* is "we asked, it is not there"; *unknown* is "we could not tell".
Collapsing them turns a broken network into a clean bill of health, so **the
split is unconditional, not a flag** — a glyph in every row, counts in the
footer of every run ◦ (§12.5).

### 5.2 Progress — the default view

A one-shot run: a bar on stderr while it works, the report on stdout when it
finishes ◦ (§12.3).

```
  acme.com · domain · 27 algorithms · 9 observers

  ████████████████████░░░░░░░░░░  68%   1,284 nodes   round 4/64   0:08
  dns 240/s   whois 2/s   geo 11/s
```

Round and depth answer "how much is left" in a way a spinner cannot, because
expansion is round-based. `whois 2/s` beside `dns 240/s` explains a stalled scan
without the user having to know whois is rate-limited. The bar is replaced by
the report when the scan ends — it is progress, not output.

### 5.3 `--tui` — the interactive view

Opt-in ◦ (§12.8). The default path is unchanged for anyone who has one, and
`--tui` on a non-TTY stdout is an **error**, not a silent downgrade.

**Progress is per target rather than global**, because every node is observed by
several operators independently and a single bar averages away which one is
stuck:

```
   SEV  NAME               EX  ALG  D  OBSERVE
 ▸ ⣿    acmе.com           ●   hr   1  ✓✓✓✓⠋   4/5
   ⣿    acme-support.com   ●   cb   2  ✓✓⠋··   2/5
   ⣶    acnne.com          ●   cr   1  ✓✓✓✓✓   5/5
        acme.net           ○   tld  1  ✓−−−−   1/1
        acmee.com          ?   cr   1  ✗⠋···   0/5
```

One glyph per operator bound to that node's type, in fixed order, so a stalled
operator reads as a **vertical** pattern:

| | |
|---|---|
| `✓` | returned ok |
| `✗` | failed |
| `⠋` | in flight |
| `·` | not yet attempted |
| `−` | not applicable — the node is absent, so nothing downstream runs |

`−` matters: an absent domain is `1/1` complete, not `1/5` stalled.

The top panel is htop's — one meter per resource rather than one bar for the
run, because the resources are rate-limited differently and `whois 2/s` beside
`dns 240/s` names the bottleneck at a glance. The bottom two lines are nvim's: a
mode indicator, the active sort and filter, and a `:` command line.

Rows are selected with the standard keys; `l` opens the detail pane on the right
and focuses it, `h` closes it. The direction is the meaning, so the pane needs no
toggle to remember — and it stays closed by default, because at eighty columns it
costs half the width of the table it describes.

**The pane is itself a row list**, navigated and highlighted identically: `j`/`k`
moves through its fields, and `l` descends again — from the `ns` row into the
four other variants sharing that nameserver. `h` pops back. Panes stack, so
navigation is a path through the graph rather than a toggle between two views
(§12.8).

| Key | |
|---|---|
| `j` `k` · `↓` `↑` | select row |
| `gg` `G` | first, last |
| `C-d` `C-u` · `C-f` `C-b` | half page, page |
| `l` `→` | open the detail pane, focus it |
| `h` `←` | close the detail pane |
| `/` `n` `N` | search, next, previous |
| `gd` | provenance — the edge that produced this node |
| `C-o` `C-i` | jump back, forward along the graph |
| `zo` | unfold the row to named observer columns |
| `g?` | keymap for the current buffer |

```
┌─ results ──────────────────────────────┬─ acmе.com ───────────────────┐
│ SEV NAME             EX ALG D PROGRESS │ domain · live · critical     │
│ ⣿   acmе.com         ●  hr  1 [||||| ] │                              │
│ ⣶   acnne.com        ●  cr  1 [||||||] │ dns    91.195.240.1          │
│     acme.net         ○  tld 1 [||||||] │ ns     ns1.above.com  +4 ▸   │
│     acmee.com        ?  cr  1 [      ] │ whois  2026-07-30 Namecheap  │
│                                        │ ptr    ✗ timeout after 10s   │
│                                        │ VARIANT_OF acme.com          │
│                                        │   hr · d1 · e→е U+0435       │
└────────────────────────────────────────┴──────────────────────────────┘
```

Commands available mid-scan:

| Command | |
|---|---|
| `:filter live risk>high` | narrow the view (display only) |
| `:depth 4` | raise the bound; the scan keeps going |
| `:sort risk \| name \| depth` | re-sort live |
| `:report` | run analysis over what exists now, marked `PARTIAL` |
| `:export json out.json` | write the current report |
| `:ledger` · `:plan` | open the declined candidates, or the compiled plan |
| `:q` | quit; the scan cancels at the round barrier |

**`:report` does not wait for the scan.** Analysis normally runs when expansion
stops; this runs it over whatever exists at that moment, so a scan that found
something worth acting on at round 6 need not reach round 64 to be read.

Every view — results, detail, findings, ledger, plan — is the same `list`
*component* over different rows (§12.9) — unrelated to the `--list` flag — so sorting, filtering, search, highlight
and folding behave identically in all of them and a binding learned in one works
in the rest.

### 5.4 `ndjson` emits per completed record

`-o json` is one document and can only be written at the end. `-o ndjson` is one
object per node, and **a record is emitted as soon as that node is done** — every
operator bound to its type has returned, failed or been skipped ◦.

```bash
urlinsane typo acme.com -o ndjson | jq -c 'select(.exists=="live")'
```

A node is emitted once, complete; there are no partial or amended records to
de-duplicate. `ndjson` already promises no ordering, so completion order costs
nothing and is the order the answers actually arrive in.

## 6. Exit codes and interruption

| Code | Meaning |
|---|---|
| `0` | clean |
| `1` | execution error |
| `2` | a finding at or above `--fail-on` |

**That is what makes the tool a CI gate** — without it the only way to react to
results is to parse stdout.

```bash
urlinsane typo --quiet --fail-on high npm:acme-internal   # $? is the answer
```

Ctrl-C stops expansion at the **end of the current round**: the barrier still
runs, so parents, belief and the truncation ledger are finalized rather than
left half-computed. Analyzers run over what exists and the report is marked
`PARTIAL` — in every format including json, so a truncated scan is never
mistaken downstream for a complete one. A second Ctrl-C aborts without a
report (§12.4).

## 7. Configuration

Application data lives in `~/.config/urlinsane`, created on first run.

| File | Contents | |
|---|---|---|
| `dataset.db` | reference data: vocabulary and weighted transitions | ✓ |
| `maxmind.db.gz` | geolocation database — **not shipped**, supply your own | ✓ |
| `blocks/` | the content-addressed graph store | ✓ |
| `scans.json` | the index from a target to its saved scans (§1.2) | ✓ |
| `config.yaml` | plugin settings | ✓ |

`dataset.db` is embedded in the binary and written out if absent, by every
command that calls `config.Init()` — `typo`, `report`, and `datasets import` and
`download`. `blocks/` and `scans.json` are created by the first
`typo --save-graph`.

`maxmind.db.gz` is **not** embedded and never written. Geolocation is opt-in:
fetch the database with `scripts/mmdb.sh` (it needs `MAXMIND_LICENSE_KEY`) and
drop it at that path, and the `geo` operator appears in the plan the next time
you scan. Absent, it stays out of the plan and nothing is printed — reporting the
absence of a feature nobody enabled is how people learn to ignore warnings.

It was embedded until 189a5ef, and the copy that shipped had never worked (§9.1),
so 49 MB of unusable bytes rode along in every release and warned on every run.
If you ran one of those releases, the broken file is still in your config
directory; it is recognised by its signature and ignored, so you can delete it or
leave it. **First-run setup is reported, not silent** ✓ — when
extraction fails, the operators that needed it are omitted from the plan, and a
scan with no geolocation must not look like a target with no geolocation
(§12.6).

### 7.1 `config.yaml`

A plugin declares its defaults in code and never has to be configured; the file
exists only to override them.

```yaml
plugins:
  dns:
    timeout: 30      # overrides the default; other dns settings keep theirs
```

**Overrides are per key, not per section** — setting one value does not mean
restating the rest. A missing file is created silently; a malformed one is an
error rather than a reason to overwrite what was hand-edited. Sections for
plugins this build does not know are preserved, since they may belong to a
plugin that is simply not loaded.

## 8. `datasets` — the maintainer tool

Builds the reference data the scanner reads. Not needed to run a scan.

```
datasets [command] [opts..] [directory]
```

The binary is `datasets`; its `App.Name` is still `data`, so its own `--help`
and error messages say `data`. That is drift in the source, not a second name.

| Command | Alias | Usage | |
|---|---|---|---|
| `import` | `i` | `import [opt..] [directory]` — import datasets into the database | ✓ |
| `download` | `d` | `download [opt..] [directory]` — download datasets | ✓ |

```bash
datasets download datasets      # fetch the public suffix list
datasets import datasets        # load every .lst into the database
```

`import` walks the tree in one pass, loading every `.lst` into the `Vocabulary`
and `Transition` tables. A language file whose line holds several words becomes
weighted transitions between them; a list-shaped file (`packages/npm.lst`)
becomes vocabulary with no transitions, which is still worth storing because it
answers membership.

`download` currently fetches only the public suffix list, despite the plural.

**`datasets/` is the source of truth and `import` runs one way.** There used to
be a `sync-languages` command that generated `datasets/languages/<lang>/` *from*
the language plugins, which pointed the wrong way: the curated lists were the
authored artefact and the plugins were built from them. It went with the
language plugins themselves.

## 9. Not yet implemented

Collected so nothing above is mistaken for a promise.

**Commands**: none. `typo` and `report` are both built. A stale
`~/.config/urlinsane/index.db` on an older install is a leftover of the deleted
SQLite model and is not the scan index — that is `scans.json` (§1.2).

**Flags**: `--quiet`, `--why`, `--ledger`, `--store`, `--plan`, `--model`,
`--<op>.model`, `--trace`, `--help-all`, and §12.2's operator-scoped flags.
`--language`, `--keyboard` and `--collect` are not registered at all, so they
are rejected as unknown flags rather than accepted and ignored.

**Views**: the progress bar, `--tui` and its per-target progress, `ndjson`
emitting per completed record rather than all of them once the scan ends, and
the three-state legend — the counts are in the footer, the glyphs are not.

**Interaction** (§12.3, §12.5–§12.7): declined-candidate warnings and the
self-explaining empty report. First-run setup reporting is built; the only
omitted-operator warnings are the two `typo` prints itself, for geolocation and
for plugin settings.

**Elsewhere**: `--ttl` cross-run caching, and `--resume` and graph diffing.

### 9.1 Gaps between what builds and what runs

Everything below compiles and is tested; what fails is the data it is handed.

- **Geolocation is off unless you supply a database.** The `maxmind.db.gz` that
  used to be embedded was not gzip: its header read `1f ef bf bd 08 08`, a gzip
  magic whose `0x8b` had been replaced by the UTF-8 replacement character, so
  the file had at some point been round-tripped through a text conversion. It
  had therefore never worked from a shipped binary.

  It is no longer embedded (189a5ef, 7a2c6ea): 49 MB left every release, and geo
  became opt-in via `scripts/mmdb.sh` (§7). The header is now validated before
  the bytes are written, so this class of corruption is refused rather than
  extracted and blamed on the operator later. What is still missing is the
  database itself, which needs a MaxMind licence key.
