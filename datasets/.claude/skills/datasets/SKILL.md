---
name: datasets
description: Work on urlinsane's reference data — the .lst tree under datasets/, how it becomes dataset.db, and the rules for editing linguistic data without destroying it. Use when adding or editing a language, a source list, or any .lst file, or when rebuilding the shipped database.
---

# Datasets

`datasets/` is the authored source for everything the scanner knows about
language: homoglyphs, misspellings, synonyms, keyboards' worth of vowels and
graphemes, plus the registries and platforms a name can exist on.

Nothing generates these files. They are hand-curated, they are the artifact,
and the database is built *from* them — never the reverse. A `sync-languages`
command once generated this tree from Go plugins; it pointed the wrong way and
was deleted along with those plugins.

## Layout

```
datasets/
  languages/<code>/*.lst   per-language linguistic data, one dir per language
  domains/*.lst            public suffix list, common prefixes, domain words
  entities/*.lst           person and family names
  packages/*.lst           package names per registry
  sources/*.lst            where a name can be checked for existence
```

`<code>` is a BCP 47 language subtag. It is canonicalized against `pkg/kb` at
build time, so a directory named with a retired code merges into the current
one: `iw` becomes `he`. A code `kb` does not ship — `la`, `no` — keeps its own
identity, because it is a real language with curated data and no keyboard.

## The two line shapes

This is the only thing you need to understand to add data correctly.

**Group-shaped** — every word on a line means the same thing, or is confusable
with the others. Order does not matter.

```
start begin launch commence          # synonym.lst
a à á â ã ä å ɑ а ạ ǎ ă ȧ ӓ           # homoglyph.lst
dot .                                # homophone.lst
hwile while                          # misspelling.lst  (wrong, right)
```

**List-shaped** — one item per line, no relationships.

```
a                                    # vowel.lst
the                                  # stopword.lst
react                                # packages/npm.lst
```

The importer does not need to be told which is which. Every word on a line joins
the vocabulary; a line with **two or more** words additionally becomes weighted
transitions between each pair, in both directions. A one-word line has nothing
to associate with, so it contributes vocabulary and no edges. That single rule
covers both shapes and every file in the tree.

Weights are counted, then normalised per source word, so a word that appears in
many groups is weighted by how often it was associated rather than by whichever
group happened to be imported last.

## Rules for editing linguistic data

**Never "clean" these files.** They look full of errors because they are lists
of errors. This has already gone wrong once: a cleaning pass filtered entries by
length and destroyed every single-character rule — Latin `ae→e`, Armenian
`օ→ո`, Pashto `ک→ك` — and a harakat stripper collapsed Persian pairs like
`لطفا→لطفاً` into self-pairs. In particular:

- A misspelling that is a real word is still a misspelling.
- A homoglyph line will contain characters that render identically. That is the
  point of the file.
- Single-character entries are the most valuable rules in the set.
- Diacritics are data. Stripping them merges rows that must stay distinct.

Conventions the existing files follow: diacritic-free ASCII for Latin-script
languages, native script otherwise, unaccented Greek, Vietnamese syllables
joined. Match the file you are editing rather than imposing these.

`#` at the start of a word makes it a comment. Blank lines are ignored. Runs of
whitespace collapse, so alignment is free.

## Populating a language

A procedure, not advice. Run it start to finish for one language code; it is
resumable, because each step writes its own file.

`go run ./cmd/datasets languages` scaffolds a directory per language `pkg/kb`
ships. Most are stubs holding a single comment line. Find the work:

```bash
for d in datasets/languages/*/; do
  printf "%s %s\n" "$(wc -l < $d/synonym.lst)" "$(basename $d)"
done | sort -n | head -20
```

### Step 0 — establish the script

Do this before any other research; it decides what every later file contains.

```bash
urlinsane typo --list keyboards | grep -i '<language name>'
```

The keyboard layouts `kb` ships for the code tell you which script the target
audience actually types. Azerbaijani, Serbian, Uzbek and Kazakh are each written
in more than one; pick the script of the layouts, not the one an encyclopaedia
lists first. Record the choice in a comment at the top of each file you fill.

### Step 1 — homoglyph.lst

Highest value, and the file general search is worst at. Do not search for
"letters that look like X". Use, in order:

1. **Unicode confusables** (TR39 `confusables.txt`) — the authoritative source.
   Fetch it, then for each character of the target script collect the entries
   that map to the same skeleton.
2. **The script's Unicode block** plus the Latin/Cyrillic/Greek blocks, for
   cross-script neighbours.
3. Only then, a search, to catch script-specific conventions the tables miss.

Format: one line per base character, that character first, then its
look-alikes separated by spaces.

```
a à á â ã ä å ɑ а ạ ǎ ă ȧ ӓ
```

Cross-script entries are the point: a Cyrillic `а` on the Latin `a` line is the
attack this tool exists to find. A line with only diacritic variants of its own
script is doing half the job.

### Step 2 — misspelling.lst

Pairs, **wrong first**, one pair per line: `hwile while`.

Sources that reliably yield results, in order of yield:
- the target language's Wikipedia "commonly misspelled words" page
- open spellchecker correction lists (Hunspell/aspell `.dic` companions)
- keyboard-adjacent slips derived from the layout `kb` ships for that code
- verb/noun forms that native writers habitually confuse

Scale for calibration: English 4,256 pairs, German 1,349. A few hundred is
already useful; do not stall trying to match English.

### Step 3 — synonym.lst

**Not a thesaurus.** These are the brand-adjacent words an attacker appends to a
name: login, verify, secure, account, pay, invoice, billing, support, delivery,
update, confirm. One theme per line, in the target language, the most natural
word first.

```
start beginn anfang starten
```

Do **not** translate the English file line by line. Ask instead what a phishing
page in that language says on its button, and what a parcel-delivery SMS in that
country says. The curated set runs to about 65 groups.

### Step 4 — the mechanical files

`vowel.lst` (one syllabic character per line), `grapheme.lst` (the script's
letters), `numeral.lst` (`0 zero zeroth`), `homophone.lst`, `antonym.lst`.
These come from the script and a dictionary, not from research.

### Step 5 — write

**Keep every existing line verbatim.** Read the file, append what is new,
deduplicate exact repeats only. Never normalise, case-fold, strip diacritics or
filter by length — see the editing rules above for what that destroyed last
time.

### Step 6 — verify

```bash
wc -l datasets/languages/<code>/*.lst          # nothing shrank
go run ./cmd/datasets build datasets
urlinsane typo --list languages | grep '<code>'
```

Then confirm the rows actually landed, per relation:

```bash
sqlite3 internal/config/dataset.db "
  select d.name, count(*) from vocabularies v
    join languages l on v.language = l.id
    join datasets  d on v.dataset  = d.id
   where l.code = '<code>' group by d.name order by 2 desc;"
```

Every file you filled must appear with a plausible count. A relation missing
from that output did not import — check the directory name against
`--list languages`, since a retired code is merged (`iw` becomes `he`) and an
unrecognised one is kept as-is, so data can land under a code you did not
expect.

There is **no `--language` flag** to scan with one language in isolation; `-l`
is specified but not built. The query above is the check.

**A file that shrank is the signal to stop and look, not to continue.**

## Source lists

`sources/*.lst` say where a name can be checked for existence. Three shapes,
because the kinds genuinely differ:

```
pypi   https://pypi.org/project/%s/  https://pypi.org/pypi/%s/json
github https://github.com/%s
gmail.com
```

Three columns is `code page-url check-url`, two is `code url`, one is a bare
provider domain. `%s` is replaced by the name. Where a third column exists it is the endpoint
that answers existence cleanly — usually an API — and it is what gets stored and
probed; the page URL is for humans.

## Commands

```bash
go run ./cmd/datasets languages          # scaffold a dir per language kb ships
go run ./cmd/datasets download datasets  # refresh the public suffix list
go run ./cmd/datasets build datasets     # rebuild internal/config/dataset.db
go run ./cmd/datasets import datasets    # load into an existing database
```

**`build` is what ships.** It writes `internal/config/dataset.db`, which
`//go:embed` picks up, so a data change reaches the binary only after a build
*and* a recompile. It deletes the database first rather than migrating it: the
shipped file once carried columns from three schema generations, with the
current ones empty, because AutoMigrate can add a column but cannot fill it.

`import` adds to whatever database is already there. Use it for a scratch
database; use `build` for the artifact.

`make dataset` is the same as `build` with the default paths.

## After editing

```bash
go run ./cmd/datasets build datasets
go test ./internal/dataset/... ./internal/plugins/variant/...
urlinsane typo --list languages          # the codes that made it in
urlinsane typo -a hr -d 0 example.com    # exercise one algorithm
```

The importer is tested in `internal/dataset/gen`. If you change a file's shape
rather than its contents, add a case there — a malformed line is silently
imported as vocabulary with no edges, which looks like success.
