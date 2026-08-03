---
title: 12 · Datasets and languages
parent: Part II · Using URLInsane
nav_order: 8
---

# Datasets and languages
{: .no_toc }

- TOC
{:toc}

Several algorithms are only as good as the data behind them. `cm` cannot invent
misspellings, `hr` cannot invent homoglyphs, and `hs` cannot invent homophones —
they read vocabulary out of a reference database, and if that database is thin
they generate little, quietly.

## Languages are data, not code

They used to be plugins: thirty Go files, each restating its vowels, graphemes,
homoglyphs and misspellings as literals, every one of which had to be edited,
compiled and shipped to correct a single word. That is gone. A language is now a
directory:

```
datasets/languages/en/
├── antonym.lst
├── grapheme.lst
├── homoglyph.lst
├── homophone.lst
├── misspelling.lst
├── negative.lst
├── numeral.lst
├── positive.lst
├── stopword.lst
├── synonym.lst
├── token.lst
├── vowel.lst
└── word.lst
```

The repository ships 30 such directories: `ar cs da de el en es fa fi fr hi hy
it iw ja ka ko la nl no pl ps pt ru sv th tr uk vi zh`.

The formats are deliberately dumb — one record per line, space-separated:

```
# homoglyph.lst — a character, then everything that looks like it
a à á â ã ä å ɑ а ạ ǎ ă ȧ ӓ ٨
b d lb ib ʙ Ь b̔" ɓ Б
c ϲ с ƈ ċ ć ç
```

```
# misspelling.lst — the misspelling, then the correct form
hwile while
authrorities authorities
emmision emission
```

{: .warning }
> Those "misspellings" are not typos in the repository. They are the algorithm's
> input data. Do not let a spell-checker or a well-meaning contributor "fix"
> them.

Beside `languages/` sit four more trees, which are not language-specific:

| Tree | Holds |
|---|---|
| `datasets/domains/` | common domain words, prefixes and suffixes — feeds `cb` |
| `datasets/entities/` | given names, family names — feeds username and email work |
| `datasets/packages/` | popular package names per registry: npm, PyPI, crates, RubyGems |
| `datasets/sources/` | where to check existence: registries, forges, username platforms, mail providers |

`sources/packages.lst` is the one worth understanding, because it is how the
`pkg` operator knows what to ask:

```
# name   human-facing URL                        existence-check endpoint
pypi     https://pypi.org/project/%s/            https://pypi.org/pypi/%s/json
npm      https://www.npmjs.com/package/%s        https://registry.npmjs.org/%s
crates   https://crates.io/crates/%s             https://crates.io/api/v1/crates/%s
```

Two URLs because the page a human should be shown and the endpoint that answers
existence cleanly are rarely the same thing.

## What the shipped database holds

The `.lst` files are the authored source. The binary reads a compiled SQLite
database, embedded with `go:embed`. As shipped it holds:

| | |
|---|---|
| Vocabulary tokens | 280,176 |
| Weighted transitions | 149,948 |
| Languages | 113 |
| Package registries | 13 |
| Repository hosts | 9 |
| Username platforms | 66 |
| Mail providers | 70 |

The largest tables are `word` (73,990), `misspelling` (60,127), the name lists
(58,257), `domain` (28,645), `synonym` (7,766), `homoglyph` (5,405) and
`homophone` (4,435).

{: .note }
> The language rows come from the **keyboard catalogue** in `pkg/kb`, not from
> the directories: a language the tool can reason about at all is listed whether
> or not its corpus has been curated yet. So `--list languages` answers "what can
> be named", not "what has data behind it" — the keyboard-driven algorithms work
> for all of them, the vocabulary-driven ones only for the curated set.
>
> `datasets/languages/` now has a directory for every one of those languages, but
> most hold only empty `.lst` files waiting to be filled: **30 are curated**, the
> rest are scaffolding. A directory is a place to put data, not data.

Check what *your* installation actually has:

```bash
urlinsane typo --list languages
urlinsane typo --list keyboards
```

{: .warning }
> **An upgrade does not refresh your local copy.** `dataset.db` is extracted from
> the binary only if `~/.config/urlinsane/dataset.db` is absent, so a new release
> with richer data will keep using your old file. If the language list looks
> short, delete it and let it be re-extracted:
> ```bash
> rm ~/.config/urlinsane/dataset.db
> ```

## Rebuilding

The `datasets` tool is the maintainer's side of this. It is built by
`make build` alongside the main binary.

```console
$ datasets --help
COMMANDS:
   import, i        Import datasets into the database
   download, d      Download datasets
   build, b         Build the shipped dataset database from a datasets tree
   languages, l     Create a directory of .lst files for every language with a keyboard
```

```bash
# Rebuild internal/config/dataset.db from the datasets/ tree
make dataset
# equivalently
go run ./cmd/datasets build datasets

# Import a tree into the working database rather than the embedded one
go run ./cmd/datasets import datasets

# Fetch the upstream sources that some trees are derived from
go run ./cmd/datasets download datasets
```

`make dataset` writes `internal/config/dataset.db`, which `go:embed` picks up —
so it must run **before** `make build` for a data change to reach the binary.

## Adding a language

1. Find `datasets/languages/<code>/`. If the language has a keyboard layout the
   directory is already there, scaffolded with an empty `.lst` per relation and a
   comment in each saying what belongs in it. If it has no layout — Latin is
   `la` — create the directory yourself.
2. Fill in the `.lst` files you can. None is mandatory; an algorithm whose file
   is empty simply generates nothing for that language.
3. `make dataset` to rebuild the embedded database.
4. `make build` and check `urlinsane typo --list languages`.

The scaffolding itself comes from `go run ./cmd/datasets languages datasets`,
which `build` also runs. It creates what is missing and never overwrites a file,
so a language kb starts shipping a layout for gets somewhere to put its data
without anyone noticing the catalogue grew. Use `--dry-run` to see what it would
add.

There is no Go code to write, no registration, and no release needed for anyone
who builds from source. That is the entire point of the change.

{: .note }
> A `sync-languages` command used to *generate* this tree from the language
> plugins. It pointed the wrong way — the curated lists were the artefact and the
> plugins were built from them — and it was removed along with the plugins.

## Keyboards

Keyboards are data too, but of a different kind: they are compiled into
[`pkg/kb`](https://pkg.go.dev/github.com/rangertaha/urlinsane/pkg/kb) as protocol
buffers rather than living in the dataset database, so they need no database at
all and `--list keyboards` works on a fresh install.

The dataset covers **203 layouts across 110 languages**, built from
[kbdlayout.info](https://kbdlayout.info), which publishes the key tables of
every keyboard driver Windows ships. Rebuilding it:

```bash
go generate ./pkg/kb      # needs protoc and protoc-gen-go on PATH
```

The full model — geometric adjacency in key units, ISO versus ANSI form factors,
wrong-layout translation — is [Keyboards]({{ site.baseurl }}/KB/).

## Geolocation

Geolocation is **off unless you supply a database**, and nothing ships one.

The database that used to be embedded was corrupt — its gzip magic read
`1f ef bf bd`, the `0x8b` replaced by the UTF-8 replacement character — so
`geo` had never worked from a shipped binary, while 49MB of unusable bytes rode
along in every release and warned on every run. Removing it halves the binary,
from 96MB to 46MB.

To turn geolocation on, put a GeoLite2 database at
`~/.config/urlinsane/maxmind.db.gz`:

```bash
MAXMIND_LICENSE_KEY=... scripts/mmdb.sh
```

The `geo` operator appears in `--list operators` and in `--explain` the next
time you scan, and stays out of the plan until then. A licence key is needed
because MaxMind's terms make redistributing the data a decision rather than a
detail, and a scanner that quietly ships a stale geolocation database is worse
than one that asks.

**Absent is silent; present but broken is not.** A file that is not there is a
feature nobody turned on. A file that is there and unusable is somebody's failed
attempt, and the run says so, names it, and tells you to remove or replace it —
which is what an upgrading user needs, since older releases left a corrupt copy
in that exact path.

What that failure demonstrates is the general rule for missing reference data:
**the operators that needed it are left out of the plan rather than failing at
run time**, because a scan with no geolocation must not look like a target with
no geolocation. Which is why `geo` does not appear in `--list operators` or in
`--explain` output today.

## `config.yaml`

A plugin declares its defaults in code and never needs configuring. The file
exists only to override them:

```yaml
plugins:
  dns:
    timeout: 30      # overrides just this key; other dns settings keep theirs
```

Overrides are per key, not per section. Sections for plugins this build does not
know are preserved, since they may belong to a plugin that is simply not loaded.

The file is read on every scan — `typo` calls `config.Load` and hands the
resolved settings to the plugin registry, where each plugin's declared defaults
are overlaid with whatever the file says. A malformed file is reported and the
scan continues on defaults:

```
plugin settings unavailable: settings: parsing ~/.config/urlinsane/config.yaml: yaml: line 1: did not find expected node content
```

Nothing writes the file for you; a missing one simply means no overrides.

---

That is Part II. **[Part III](../../internals/)** is the engine.
