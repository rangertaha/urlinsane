---
name: datasets
description: Work on urlinsane's reference data — the .lst tree under datasets/, how it becomes dataset.db, where to source each relation (homoglyphs, misspellings, synonyms, homophones, graphemes, vowels, dictionaries, stopwords) and the rules for editing linguistic data without destroying it. Use when adding or editing a language, a source list, or any .lst file, or when rebuilding the shipped database.
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

## The twelve files

Every `languages/<code>/` holds the same twelve files. `go run ./cmd/datasets
languages` scaffolds all of them, each as a stub holding a one-line comment, so
an empty language is twelve stubs rather than a missing directory — absence of
data never shows up as an absent file, and `ls` will not tell you what is done.
Counts below are `en`, the fullest set, for calibration.

**`word.lst`** (8783) — the language's vocabulary, one word per line, and the
largest file by an order of magnitude. It is list-shaped, so it contributes no
edges; its job is to be the thing a candidate name is checked against. It is not
a dictionary: `en` holds `agregate`, `agreing` and `ahev` alongside their correct
forms, because a squatted name is a real string that has to be recognised, not a
well-formed one. It also holds digits, `-`, `.` and `/`. Source it from a
Hunspell dictionary whose licence is GPL-3.0 compatible, and fold in the wrong
halves of `misspelling.lst`.

**`misspelling.lst`** (4256) — the errors people actually make, group-shaped,
`wrong right` per line. This is the highest-value file after `word.lst` and the
one with a real corpus behind it: Wikipedia's machine-readable list for English,
its local equivalent on other large Wikipedias, and keyboard-adjacency slips
derived from the `pkg/kb` layout where no list exists. A misspelling that is
itself a real word still belongs here.

**`homophone.lst`** (487) — words that sound alike, group-shaped: `acts ax axe`,
`hoard horde`. The first few lines are the special case that matters most for
names — the spoken form of a character, `dot .`, `dash -`, `at @`, `slash /`,
`underscore _` — which is what turns a dictated or voice-phished address into a
typo. Author those five before the rhyming pairs.

**`antonym.lst`** (93) — a word and its opposite, one pair per line: `above
below`, `accept refuse`. Group-shaped, so the pair becomes an edge in both
directions, which is the point: an attacker swaps `secure` for `insecure` as
readily as for a synonym.

**`synonym.lst`** (65) — not a thesaurus. These are the brand-adjacent lure words
appended to a name — `login signin logon access`, `verify confirm validate
check` — one theme per line, most natural word first. No corpus exists; get them
from what real lures in that language say, and translate the concept rather than
the English line. About 65 groups is the curated size, and every populated
language has exactly that, which makes it the quickest way to spot a language
somebody started and abandoned.

**`stopword.lst`** (61) — words carrying no distinguishing sense, list-shaped.
They mark the parts of a candidate name that should not be weighted, so `the` in
`thebank` does not read as signal. Any standard stopword list for the language
does.

**`numeral.lst`** (36) — `digit word ordinal` per line: `0 zero zeroth`, `1 one
first`. Group-shaped, so it wires the digit to its names in both directions,
which is what catches `1` for `one` and `4` for `four` in a label. It comes from
CLDR's spelled-out number rules or a grammar reference, never from a corpus;
`en`'s 36 lines cover 0–24, the tens to 90, then 100, 1000, and the millions and
billions.

**`grapheme.lst`** (26) — the smallest functional units of the writing system,
one per line. Take it from CLDR `exemplarCharacters` for the locale, which is
exactly this data: strip the brackets and split. For Latin languages it is the
alphabet including the accented letters the language actually uses; for an
abugida it is the consonants and the dependent vowel signs both.

**`homoglyph.lst`** (26) — a character followed by everything that looks like it:
`a à á â ã ä å ɑ а ạ ǎ ă ȧ ӓ ٨`. Group-shaped and the heart of the whole tool.
Build it from Unicode confusables (UTS #39), which needs **inverting** — it maps
confusable → skeleton and this file is keyed the other way — then unioned with
the NFD-derived accented forms, which confusables does not carry. A line holding
only the accents or only the cross-script look-alikes is doing half the job. `en`
has one row per ASCII letter; the other Latin languages extend that same block
with rows headed by their own accented letters, to 42, and scripts that must map
onto Latin run longer still — `el` to 65.

**`negative.lst`** (20) and **`positive.lst`** (20) — connotation, one word per
line, list-shaped. Deliberately small and deliberately domain-specific rather
than general sentiment: `positive` is `safe trusted secure official`, `negative`
is `scam malicious fraud phishing spam`. They are the words a squatted name
leans on to look legitimate or the words that flag one, not a sentiment lexicon.
Keep them ASCII-folded for Latin-script languages, as the existing files do.

**`vowel.lst`** (5) — the syllabic characters, list-shaped, feeding the
vowel-swap algorithms. Decide the subset from the script rather than from the
Latin five: for an abugida the dependent vowel signs belong here, and a language
whose `y` is syllabic should list it, as `cs` does.

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

def variants_of(base):
    out = []
    for cp in range(0x20, 0x2500):
        d = u.normalize("NFD", chr(cp))
        if len(d) > 1 and d[0] == base and all(u.combining(x) for x in d[1:]):
            out.append(chr(cp))
    return out
# variants_of("a") -> à á â ã ä å ā ă ą ǎ ǟ ǡ ǻ ȁ ȃ ȧ ḁ ạ ả ấ ...  (29)
```

Union the two, then **drop what nobody can type into a name**: the
mathematical alphanumeric blocks (`𝐚 𝑎 𝒂 𝓪 …`, U+1D400–U+1D7FF) and fullwidth
forms dominate the raw confusables output and are neither valid in a domain
label nor plausible in a package name. Keep diacritics, cross-script letters
and digit-shaped look-alikes.

A line holding only diacritic variants of its own script is doing half the job;
so is one holding only the cross-script half.

Two files beside it in the same directory are worth knowing about:
`intentional.txt` is 6KB against confusables' 728KB and holds the pairs Unicode
considers *intended* to look alike, which is closer to what a squatter picks;
`IdentifierType.txt` says which characters are permitted in identifiers at all,
and filtering by it drops the mathematical-alphanumeric noise more precisely
than excluding blocks by hand. See the source catalogue below.

### Step 2 — misspelling.lst

**English — Wikipedia's machine-readable list**, which is already pairs:

```
https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/For_machines
```

Format is `wrong->right`, one per line, so the conversion is
`line.replace("->", " ")`.

**Other languages — Wiktionary's misspelling categories cover 144 of them**,
under one naming pattern, so the same scrape works everywhere:

```
https://en.wiktionary.org/wiki/Category:<Language>_misspellings
```

Spanish 274 entries, Portuguese 312, Italian 208, German 71, French 77. Each
entry page names the correct form, which is the second column. Most large
Wikipedias also keep their own version of the list above; search
`site:<lang>.wikipedia.org` for the local equivalent of "commonly misspelled
words".

Where neither exists, generate from two sources instead:

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

## Source catalogue

Where the data comes from, per file. The procedure above names the one source to
reach for first; this is the rest, for when that one has no coverage for your
language.

**Read the licence before copying anything.** This repo is
**GPL-3.0-or-later**, and these sources are not uniformly compatible with it.
Three groups, and the difference matters:

| | |
|---|---|
| **Unicode / CLDR** | Unicode licence, permissive. Copy freely. |
| **Wikipedia / Wiktionary / OpenSubtitles-derived** | CC BY-SA 4.0. Usable, but attribution is a condition — record the source URL and the licence in a comment at the top of any file you fill from it. |
| **Research-only corpora** | Cannot be copied in at all. Use them to *check* words you sourced elsewhere. |

A comment line at the top of a `.lst` costs nothing and is the only record of
where a line came from once it is one word among four thousand.

### homoglyph.lst

| Source | Coverage | Licence |
|---|---|---|
| [confusables.txt](https://www.unicode.org/Public/security/latest/confusables.txt) | ~6,500 characters, cross-script | Unicode |
| [intentional.txt](https://www.unicode.org/Public/security/latest/intentional.txt) | the deliberately-confusable subset — small, high signal | Unicode |
| Unicode NFD decomposition | the accented half confusables omits | n/a, computed |
| [codebox/homoglyph](https://github.com/codebox/homoglyph) | pre-grouped char sets, several languages | MIT |

`intentional.txt` is the one worth adding to the procedure's step 1: it is 6KB
rather than 728KB and holds the pairs Unicode considers *intended* to look
alike in a harmonised typeface, which is a closer match to what a squatter
picks than the full confusables table.

The `IdentifierStatus.txt` and `IdentifierType.txt` files in the same directory
say which characters are *allowed* in identifiers at all. Filtering by those is
a cheaper way to drop the mathematical-alphanumeric noise than the block
exclusion in step 1.

### misspelling.lst

| Source | Coverage | Licence |
|---|---|---|
| [Wikipedia machine-readable list](https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/For_machines) | English, ~4,000 pairs, already `wrong->right` | CC BY-SA |
| [Wiktionary misspelling categories](https://en.wiktionary.org/wiki/Category:Misspellings_by_language) | **144 languages** | CC BY-SA |
| [GitHub Typo Corpus](https://github.com/mhagiwara/github-typo-corpus) | 350k edits, 15+ languages, mined from commits | check the repo |
| [Birkbeck spelling error corpora](https://www.dcs.bbk.ac.uk/~roger/corpora.html) | English, human error data | check |

**The Wiktionary categories are the answer to the gap step 2 admits.** They
follow one naming pattern — `Category:<Language> misspellings` — so the same
scrape works for all 144: English 2,490 entries, Portuguese 312, Spanish 274,
Italian 208, French 77, German 71, Dutch 73. Each entry page names the correct
form, which is the second column.

The GitHub Typo Corpus is real typing errors rather than curated ones, which
makes it the better source for keyboard-shaped slips and the worse one for the
knowledge errors an attacker banks on. Its paper is already in
`docs/papers/2020.lrec-1.835.pdf`.

### synonym.lst

The procedure is right that no corpus fits: this file is brand-adjacent **lure
words**, not a thesaurus, and translating an English thesaurus entry gives you
the dictionary word rather than the one on the phishing button.

Thesauri are still useful for *widening* a group you have already established
from a real lure:

| Source | Coverage | Licence |
|---|---|---|
| [Open Multilingual WordNet](https://omwn.org/) | 60 wordnets, 49 languages | open, but **per-wordnet** — check each |
| [OpenThesaurus](https://www.openthesaurus.de/) | German, large and current | LGPL / GPL |
| [ConceptNet](https://conceptnet.io/) | ~300 languages, weaker per-language | CC BY-SA |

### antonym.lst

Same shape as synonyms, same sources — Open Multilingual WordNet carries
antonymy as a relation, so it is one query rather than a second corpus.

### homophone.lst

Nothing ships a homophone list per language. Derive one: **words sharing an IPA
transcription are homophones by definition.**

| Source | Coverage | Licence |
|---|---|---|
| [WikiPron](https://github.com/CUNY-CL/wikipron) | 1.7M pronunciations, **165 languages**, mined from Wiktionary | Apache-2.0 tool, CC BY-SA data |
| [CMU Pronouncing Dictionary](https://github.com/cmusphinx/cmudict) | English, 134k words | permissive — check |

Group a WikiPron TSV by its transcription column and every bucket with more
than one member is a homophone group. That is the whole derivation, and it is
the single highest-yield source in this catalogue for languages nobody has
curated.

Note the English file also holds symbol names — `dot .`, `dash -` — which no
corpus will give you. Those are typed-aloud spellings and have to be written by
hand.

### grapheme.lst and vowel.lst

| Source | Coverage | Licence |
|---|---|---|
| [CLDR exemplar characters](https://github.com/unicode-org/cldr-json) | per-locale alphabet, `cldr-misc-full/main/<code>/characters.json` | Unicode |
| Unicode `Script` property | what belongs to a script at all | Unicode |
| [PHOIBLE](https://phoible.org/) | phoneme inventories, 2,000+ languages | CC BY-SA |

CLDR's `exemplarCharacters` is the alphabet; `auxiliaryExemplarCharacters` in
the same file is the borrowed-and-accented set a name may still contain, which
is worth unioning in for `grapheme.lst` and worth *excluding* from a strict
alphabet.

PHOIBLE is phonemic, not orthographic, so it does not give graphemes. It
answers the one question `vowel.lst` actually poses — which sounds are syllabic
in this language — for scripts where the Latin five are no guide.

### word.lst

| Source | Coverage | Licence |
|---|---|---|
| [wooorm/dictionaries](https://github.com/wooorm/dictionaries) | Hunspell, 92 languages | MIT wrapper, **per-dictionary licence** |
| [FrequencyWords](https://github.com/hermitdave/FrequencyWords) | 61 languages, ranked by frequency | CC BY-SA (OpenSubtitles) |
| [SCOWL](http://wordlist.aspell.net/) | English, size-tiered | permissive — check |

**Frequency ranking is what Hunspell does not give you**, and it is what this
file wants: a combosquat is built from words people actually use, so the top
few thousand by frequency beat a complete dictionary sorted alphabetically.
FrequencyWords is one file per language, `word count` per line.

### stopword.lst

| Source | Coverage | Licence |
|---|---|---|
| [stopwords-iso](https://github.com/stopwords-iso/stopwords-iso) | **57 languages**, ISO 639-1 keyed, one JSON | MIT |

One file, permissively licensed, keyed by the same codes this tree uses. There
is no reason to source this anywhere else.

### positive.lst and negative.lst

| Source | Coverage | Licence |
|---|---|---|
| [NRC Emotion Lexicon](https://saifmohammad.com/WebPages/NRC-Emotion-Lexicon.htm) | 27k terms, 100+ languages | **research only — do not copy in** |
| Wiktionary / OMW | per-language, hand-picked | CC BY-SA / per-wordnet |

The NRC lexicon is the obvious source and the one you cannot use: it is
research-licensed, so copying it into a GPL-3.0 repo is not open to us. It is
still legitimate to *consult* it and write your own twenty words. Both curated
files are twenty lines — this is an afternoon of judgement, not a scrape.

### numeral.lst

CLDR's spelled-out number rules (RBNF) or a grammar reference. `0 zero zeroth`
per line — digit, cardinal, ordinal — and it stops being worth automating at
about thirty lines.

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






Researching a language from datasets/languages/* and reasearch each of the .lst files is a multi-step process. The steps are: