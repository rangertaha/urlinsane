# URLInsane

[![Go Report Card](https://goreportcard.com/badge/github.com/rangertaha/urlinsane?style=flat-square)](https://goreportcard.com/report/github.com/rangertaha/urlinsane) [![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/rangertaha/urlinsane) [![PkgGoDev](https://pkg.go.dev/badge/github.com/rangertaha/urlinsane)](https://pkg.go.dev/github.com/github.com/rangertaha/urlinsane) [![Release](https://img.shields.io/github/release/rangertaha/urlinsane.svg?style=flat-square)](https://github.com/rangertaha/urlinsane/releases/latest) [![Build Status](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml/badge.svg)](https://github.com/rangertaha/urlinsane/actions/workflows/go.yml)





Urlinsane is a tool for detecting domain typosquatting and supporting OSINT investigations, designed to operate on multilingual target domains. It helps identify threats such as typosquatting, brandjacking, URL hijacking, phishing, fraud, corporate espionage, supply chain attacks, and more. This command-line tool generates and scans for potential typosquatting variants of a domain, assisting in uncovering and mitigating security risks.

It's inspired by [URLCrazy](https://morningstarsecurity.com/research/urlcrazy), [Dnstwist](https://github.com/elceef/dnstwist), and a few other libraries and tools I was researching at the time.

**Full documentation: [rangertaha.github.io/urlinsane](https://rangertaha.github.io/urlinsane/)** —
the CLI reference, the engine design, and the keyboard model, as a book. The
same pages live in this repo as [docs/CLI.md](docs/CLI.md),
[docs/DESIGN.md](docs/DESIGN.md) and [docs/KB.md](docs/KB.md).



## Installation

This tool is primarily intended for Linux operating systems.

* [urlinsane-0.8.2-darwin-amd64](https://github.com/rangertaha/urlinsane/releases/download/0.8.2/urlinsane-0.8.2-darwin-amd64)
* [urlinsane-0.8.2-linux-amd64](https://github.com/rangertaha/urlinsane/releases/download/0.8.2/urlinsane-0.8.2-linux-amd64)
* [urlinsane-0.8.2-windows-amd64.exe](https://github.com/rangertaha/urlinsane/releases/download/0.8.2/urlinsane-0.8.2-windows-amd64.exe)

### Linux
Download the binary, remove the previous version, and install it in /usr/local/bin:

```bash
wget https://github.com/rangertaha/urlinsane/releases/download/0.8.2/urlinsane-0.8.2-linux-amd64 
chmod +x urlinsane-0.8.2-linux-amd64 
sudo mv urlinsane-0.8.2-linux-amd64  /usr/local/bin/urlinsane
```

## Usage

Generate variations of a target and gather information on them with the `typo` command:

```bash
urlinsane typo example.com
```

**Typosquatting is not limited to domains, and the target's kind is detected
from the string alone** — there is no `--type` flag:

```bash
urlinsane typo acme.com                 # domain
urlinsane typo bob@acme.com             # email: varies bob, acme.com, and the address
urlinsane typo npm:lodash               # package on a named registry
urlinsane typo github.com/acme/tool     # repo
urlinsane typo bobsmith                 # username
```

An optional first positional narrows what gets varied, without changing how the
target is read:

```bash
urlinsane typo username acme.com/bob    # vary only bob
urlinsane typo domain bob@acme.com      # vary only acme.com
```

### `typo` options

`urlinsane typo [<scope>] <target> [flags]`

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--depth` | `-d` | `3` | Observation hops from the seed |
| `--algorithm` | `-a` | all | Restrict variant generation to these algorithm IDs; `^id` excludes |
| `--filter` | `-f` | | Select report rows: `live`, `absent`, `unknown`, `untried`, `risk>SEV`, `type=NAME`, `depth<=N` |
| `--output` | `-o` | `table` | `table`, `json`, `ndjson`, `csv`, `dot` |
| `--save` | | | Write the report to a path; format from the extension |
| `--save-graph` | | | Persist the graph to the store and print its root CID; `urlinsane report <target>` renders it again |
| `--fail-on` | | | Exit `2` if any finding reaches a severity — the CI gate |
| `--verbose` | `-v` | | Include provenance and engine belief |
| `--explain` | | | Compile and print the plan without running it |
| `--list` | | | `types`, `relations`, `operators`, `algorithms`, `languages`, `keyboards`, `formats`, `filters` |

`--filter` selects rows **in the report**, never work in the scan — narrowing
the scan is what `--depth`, `--algorithm` and the scope positional do.

Exit codes: `0` clean, `1` execution error, `2` a finding at or above
`--fail-on`.

**[docs/CLI.md](docs/CLI.md) is the full reference**, including the flags that
are specified but not yet built. [docs/DESIGN.md](docs/DESIGN.md) §12 is the
reasoning behind the interface. Both are also in the
[book](https://rangertaha.github.io/urlinsane/).

List what a build has registered:

```bash
urlinsane typo --list algorithms
urlinsane typo --list keyboards
```

## What a build registers

| Kind | Count | |
|---|---|---|
| Node types | 10 | asn, domain, email, ip, package, platform, registrant, repo, tld, username |
| Algorithms | 30 | generate variants of a name |
| Operators | 42 | 27 variant operators, one per algorithm, plus 15 that decompose and observe |
| Keyboards | 133 | distinct key-adjacency sets, from the 203 layouts `pkg/kb` ships |
| Languages | 113 | codes in `dataset.db`; 30 have curated trees under `datasets/languages/` |
| Formats | 5 | table, json, ndjson, csv, dot |

Counts are what `--list` prints on a build of this tree with `internal/config/dataset.db`
imported. Types, algorithms, operators and formats come from Go registries, so
they are fixed by the binary; keyboards come from `pkg/kb`, also compiled in.
**Languages come from the dataset database in `~/.config/urlinsane/`**, which is
extracted only when absent — an older copy left there from a previous version
will list something else. `geo` and the `pkg`/`usr`/`repo` operators are
conditional (see below), so 42 is the count when their data is present.

Languages and keyboards are **data**, not plugins: a language is a directory
under `datasets/languages/`, a keyboard a layout in `pkg/kb`, and neither needs
Go code. Output formats are a closed set the report projects into. What remains
extensible is `internal/plugins` — operators, analyzers and algorithms — one
directory per plugin, grouped by kind (`decompose`, `variant`, `observe`,
`analyze`, `report`).

### Language IDs

```bash
urlinsane typo --list languages
```

Languages are two-letter directory names under `datasets/languages/`; Pashto is
`ps` and Latin is `la`. Note that this lists every code the *dataset database*
knows — 113 — not the ones with data behind them: 31 carry vocabulary, and 30
have a curated tree in this repo.

### Language Datasets (`datasets/languages/`)

The repo ships with a `datasets/languages/<lang>/` structure (e.g. `numeral.lst`, `homoglyph.lst`, `homophone.lst`, `positive.lst`, `negative.lst`, etc).

These files are the authored source, hand-curated per language, and nothing
generates them. Load them into the dataset database with:

```bash
go run ./cmd/datasets import datasets
```

A `sync-languages` command used to generate this tree *from* the language
plugins. That pointed the wrong way — the curated lists were the artefact and
the plugins were built from them — and it was removed along with the language
plugins themselves. Languages are data now, not code: adding one means adding a
directory here and re-importing.

### Keyboard Layouts

Keyboard layouts are data, compiled in from `pkg/kb`: 203 shipped layouts,
which collapse to the 133 distinct key-adjacency sets the algorithms actually
run over. To list them:

```bash
urlinsane typo --list keyboards
```


## Algorithms

Algorithms generate plausible variants of a name. `--list algorithms` prints
this table for the build you have.

**Applies to** is blank where an algorithm binds by capability rather than by
type — those run on any nameable node, domain or package or handle alike.

| ID | Name | Applies to | Description |
|---|---|---|---|
| `aci` | Adjacent Character Insertion | any | Insert a character adjacent on the keyboard. |
| `acs` | Adjacent Character Substitution | any | Replace a character with a keyboard neighbour. |
| `afx` | Affix Squatting | package, repo, username | Add a plausible prefix or suffix. |
| `bf` | Bit Flipping | any | Flip one bit of a character — bitsquatting. |
| `cb` | Combo Squatting | any | Append or prepend a common keyword. |
| `cm` | Common Misspellings | any | Apply a curated misspelling for the language. |
| `cns` | Cardinal Substitution | any | Swap a number for its cardinal word, and back. |
| `co` | Character Omission | any | Drop a character. |
| `cr` | Character Repetition | any | Double a character. |
| `cs` | Character Swapping | any | Transpose two adjacent characters. |
| `dhs` | Dot Hyphen Substitution | any | Swap dots and hyphens. |
| `di` | Dot Insertion | any | Insert a period. |
| `do` | Dot Omission | any | Remove a period. |
| `gi` | Grapheme Insertion | any | Insert a grapheme from the language's set. |
| `gr` | Grapheme Replacement | any | Replace a grapheme with another. |
| `hi` | Hyphen Insertion | any | Insert a hyphen. |
| `ho` | Hyphen Omission | any | Remove a hyphen. |
| `hr` | Homoglyph Replacement | any | Replace a character with one that looks the same. |
| `hs` | Homophone Substitution | any | Replace a word with one that sounds the same. |
| `nsc` | Namespace Confusion | package, repo | Move a name between namespaces or scopes. |
| `ons` | Ordinal Substitution | any | Swap a number for its ordinal word, and back. |
| `rar` | Repetition Adjacent Replacement | any | Double a character, then replace the double with a neighbour. |
| `sep` | Separator Substitution | package, repo, username | Swap the separator a registry allows. |
| `si` | Subdomain Insertion | domain | Insert a subdomain label. |
| `sld` | Wrong Second-Level Domain | domain | Swap the second level under a ccTLD: `bbc.co.uk` → `bbc.org.uk`. |
| `sp` | Singular Pluralise | any | Make a word singular or plural. |
| `tld` | Wrong TLD | domain | Substitute a different public suffix. |
| `tli` | TLD Insertion | domain | Append a suffix so the whole name becomes a subdomain: `example.com.br`. |
| `vs` | Vowel Swapping | any | Swap one vowel for another. |
| `xhs` | Cross-language Homophone | any | Swap for a spelling that sounds the same in another language: `youtube` → `yutup`. |

## Operators

An operator is what expands or observes the graph. Where the old collectors ran
in a fixed order over a list of domains, an operator declares what data pattern
it **binds to** and what it **emits**, and the scheduler decides what runs when.
`--list operators` prints the plan-eligible set:

| ID | Binds on | Emits |
|---|---|---|
| `decompose.domain` | domain | `TLD_OF` |
| `decompose.email` | email | `LOCAL_PART`, `DOMAIN_OF` |
| `decompose.package` | package | `OWNER` |
| `decompose.repo` | repo | `HOSTED_ON`, `OWNER` |
| `dns-a` | domain | `RESOLVES_TO` |
| `dns-mx` | domain | `MX` |
| `dns-ns` | domain | `NS` |
| `dns-cname` `dns-txt` | domain | props only |
| `ptr` | ip | `PTR_TO` |
| `whois` | domain | `REGISTERED_BY` |
| `idn` | domain | props only |
| `geo` | ip | props only — needs the geolocation database |
| `pkg` `usr` `repo` | package, username, repo | `EXISTS_ON` — needs the source lists |

Plus one operator per algorithm, all emitting `VARIANT_OF`.

**Binding is by data, not by producer.** `ptr` binds to any `ip`, so it runs on
addresses whether `dns-a` found them or something else did — which is what lets
a new operator slot in without anyone rewiring an order.

**Three-state existence.** An operator that cannot reach a registry reports
*unknown*, never *absent*. "We asked, it is not there" and "we could not tell"
are opposite conclusions, and collapsing them turns a broken network into a
clean bill of health.

`geo`, `pkg`, `usr` and `repo` are omitted from the plan when the data they
need is missing, rather than failing at runtime — so `--list operators` shows
them only on a build that has it. `pkg`, `usr` and `repo` need the source lists,
which `dataset.db` now carries, and appear; `geo` needs the MaxMind database
extracted into `~/.config/urlinsane/`, and drops out when that is absent or
unreadable.

## Output formats

| Format | Description |
|---|---|
| `table` | Pretty table with colour styling; the default |
| `json` | One document, written when the scan ends |
| `ndjson` | One object per node |
| `csv` | Comma-separated values |
| `dot` | Graphviz — the graph, not a flattened list |

```bash
urlinsane typo acme.com -o json | jq '.nodes[] | select(.exists=="live")'
urlinsane typo acme.com --save report.csv      # format from the extension
```

`--save` also accepts `.txt`/`.text` for `table` and `.gv` for `dot`. Anything
else is an error rather than a guess, and a saved file is never coloured.

**`--filter` selects rows, not columns**, and it applies to the report rather
than the scan — so re-filtering never costs another lookup.

## Status

The engine is mid-rewrite, from a linear plugin pipeline to a graph engine.
`docs/DESIGN.md` is the design; `docs/CLI.md` §9 tracks what is specified but
not yet wired up.

**Done, and worth saying how it differs from the plan:**

- **The DAG replaced the pipeline.** The original idea was Terraform-style
  declared dependencies between plugins. That is not what shipped: an operator
  declares what data pattern it *binds to*, and the scheduler matches. Declared
  dependencies made plugin order load-bearing and the cache unsound; binding by
  data means a new operator needs no rewiring.
- **Reference data moved into SQLite** (`dataset.db`): vocabulary and weighted
  transitions, replacing a large body of generated Go. Results did **not** —
  they are an IPLD content-addressed graph, so two identical scans address
  identically and "what changed since last week" is a CID comparison.
- **Languages and keyboards stopped being plugins.** 30 curated dataset
  directories and 203 keyboard layouts built from
  [kbdlayout.info](http://kbdlayout.info/), neither needing Go code.
- **Cross-scan diffing** exists in `internal/store`.
- **Saving and replaying scans.** `typo --save-graph` writes the graph to the
  store; `urlinsane report <target>` renders it again from the stored blocks,
  and `report --scans <target>` lists what has been saved. Plugin settings in
  `~/.config/urlinsane/config.yaml` reach the plugins that declare them.

**Open:**

- Flags that are specified but not built — `--quiet`, `--why`, `--ledger`,
  `--tui`, and `--ttl`/`--resume` cross-run caching (`docs/CLI.md` §9).
- An advanced keyboard model with layer-shifting.
- DNS queries against several resolvers.
- Dataset updates downloadable rather than embedded, to cut binary size.
- **[LLM](https://en.wikipedia.org/wiki/Large_language_model)** assistance for
  generating language datasets, and as a judgement operator over variants.
- Reporting confirmed squats back to a shared corpus, so the transition weights
  can be learned from observed cases instead of being uniform placeholders.

## Other Tools

| Name  | Language | Description                    |
|-------|-----------|--------------------------------|
| [Urlcrazy](https://github.com/urbanadventurer/urlcrazy) |  Ruby  |  URLCrazy is an OSINT tool to generate and test domain typos or variations to detect or perform typo squatting, URL hijacking, phishing, and corporate espionage.  |
| [DNSTwist](https://github.com/elceef/dnstwist) | Python   | Domain name permutation engine for detecting homograph phishing attacks, typo squatting, and brand impersonation     |
| [DomainFuzz](https://github.com/monkeym4ster/DomainFuzz) | JavaScript   | Domain name permutation engine for detecting typo squatting, phishing and corporate espionage   |






## Authors

* [Rangertaha (rangertaha@gmail.com)](https://github.com/rangertaha)

## License

This project is licensed under the GPLv3 License - see the [LICENSE](LICENSE) file for details






