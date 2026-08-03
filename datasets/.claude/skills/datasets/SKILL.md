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

**Source: Unicode confusables (UTS #39).** This is the authoritative data and it
is a file, not a search.

```
https://www.unicode.org/Public/security/latest/confusables.txt
```

Format is `SOURCE ; TARGET ; TYPE # comment`, hex code points, semicolon
separated, one *source* per line:

```
0430 ;	0061 ;	MA	# ( а → a ) CYRILLIC SMALL LETTER A → LATIN SMALL LETTER A
```

It maps confusable → skeleton, which is the direction this file needs
**inverted**: group every line by its TARGET, and each group becomes one
`homoglyph.lst` line with the target first.

```python
# confusables.txt -> cross-script look-alikes, grouped by target
groups = {}
for line in open("confusables.txt", encoding="utf-8"):
    line = line.split("#")[0].strip()
    if not line:
        continue
    src, tgt, *_ = [c.strip() for c in line.split(";")]
    if " " in tgt:            # multi-character skeletons are not homoglyphs
        continue
    groups.setdefault(chr(int(tgt, 16)), []).append(chr(int(src, 16)))
```

**Confusables alone is not enough, and this is the part that is easy to get
wrong.** It carries cross-script look-alikes but *not* accented forms: the
group for `a` contains `ɑ α а` and no `à á â ã ä å`, because those are handled
by normalisation rather than by confusion. The curated files contain both. Get
the accented half from Unicode decomposition:

```python
# every character whose NFD base is the target -> the accented variants
import unicodedata as u
variants = [chr(cp) for cp in range(0x20, 0x2500)
            if (d := u.normalize("NFD", chr(cp)))[:1] == ("a",)
            and len(d) > 1 and all(u.combining(x) for x in d[1:])]
# a: à á â ã ä å ā ă ą ǎ ǟ ǡ ǻ ȁ ȃ ȧ ḁ ạ ả ấ ...  (29 in the BMP range above)
```

Union the two, then **drop what nobody can type into a name**: the
mathematical alphanumeric blocks (`𝐚 𝑎 𝒂 𝓪 …`, U+1D400–U+1D7FF) and fullwidth
forms dominate the raw confusables output and are neither valid in a domain
label nor plausible in a package name. Keep diacritics, cross-script letters
and digit-shaped look-alikes.

A line holding only diacritic variants of its own script is doing half the job;
so is one holding only the cross-script half.

### Step 2 — misspelling.lst

**English — Wikipedia's machine-readable list**, which is already pairs:

```
https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/For_machines
```

Format is `wrong->right`, one per line, so the conversion is
`line.replace("->", " ")`.

**Other languages** — most large Wikipedias keep the same page under their own
name; search `site:<lang>.wikipedia.org` for the local equivalent of "commonly
misspelled words". Where none exists, generate from two sources instead:

1. **Keyboard-adjacent slips**, derived from the layout `pkg/kb` already ships
   for that code. These are mechanical and language-specific, and no external
   source is needed — `urlinsane typo --list keyboards | grep -i <language>`.
2. **A spellchecker word list** (see below) plus that language's known
   confusion pairs — verb forms, agreement, borrowed spellings.

Scale for calibration: English 4,256 pairs, German 1,349. A few hundred is
already useful; do not stall trying to match English.

### Step 3 — synonym.lst

**No corpus exists for this one**, because it is not a thesaurus. These are the
brand-adjacent lure words an attacker appends to a name, and the way to get them
is to look at what real lures in that language say:

- search for phishing advisories published by that country's bank, postal
  service or CERT — they quote the lure text
- the local words on a login button, a parcel-tracking notice, an invoice
- APWG and national CERT advisories often reproduce the subject lines

Translate the *concepts* (login, verify, secure, account, pay, invoice,
billing, support, delivery, update, confirm), never the English file line by
line: what a German phishing page puts on its button is not the dictionary
translation of what an English one does. One theme per line, most natural word
first. About 65 groups matches the curated set.

### Step 4 — grapheme.lst, vowel.lst and word.lst

**Alphabet — CLDR exemplar characters**, which is exactly this data per locale:

```
https://github.com/unicode-org/cldr-json
  -> cldr-misc-full/main/<code>/characters.json
```

`exemplarCharacters` is a bracketed list, e.g. Greek
`[α ά β γ δ ε έ ζ η ή θ ι ί ϊ ΐ κ λ μ ν ξ ο ό π ρ σ ς τ υ ύ ϋ ΰ φ χ ψ ω ώ]`.
Strip the brackets and split: that is `grapheme.lst`, one per line.
`vowel.lst` is the subset that is syllabic in that script — decide from the
script, not from the Latin five.

**Word lists — Hunspell dictionaries:**

```
https://github.com/wooorm/dictionaries   # 92 languages, BCP-47 directories
```

Each holds `index.dic` (one word per line after the count on line 1) and
`index.aff`.

> **Check the licence before copying anything.** That repo is MIT but each
> dictionary keeps its original licence, and they differ — GPL-2.0, LGPL-2.1,
> MIT, Apache-2.0. This repo is **GPL-3.0-or-later**, so a GPL-2.0-**only**
> dictionary cannot be copied into it. Read `dictionary-<code>/license`, and if
> it is incompatible use the dictionary to *check* words you sourced elsewhere
> rather than as the source.

`numeral.lst` is `digit word ordinal` per line (`0 zero zeroth`) and comes from
CLDR's spelled-out number rules or a grammar reference, not from a corpus.

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
provider domain. `%s` is replaced by the name. Where a third column exists it
is the endpoint that answers existence cleanly — usually an API — and it is
what gets stored and probed; the page URL is for humans.

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
