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

**Start from this language's own alphabet, not from the confusables table.**

The file is keyed by the characters the target audience types, and each line
answers one question: *what, in any other script, looks like this character?*
A Russian brand is attacked with Latin look-alikes of Cyrillic letters; a Latin
brand is attacked with Cyrillic and Greek look-alikes of Latin ones. Same
computation, opposite direction, and the direction is decided by whose alphabet
supplies the keys.

Get the keys from the same CLDR exemplar set step 4 uses, then look each one up.
Working from the confusables file alone produces a Latin-keyed file whatever
language you meant to fill, because Latin is what most of that table targets —
which is how `el/homoglyph.lst` ended up opening with `a`, `b`, `c`.

Check the result before writing it:

```bash
# how many keys are actually in the language's script?
awk '{print $1}' datasets/languages/<code>/homoglyph.lst | \
  python3 -c "import sys,unicodedata as u
from collections import Counter
print(Counter(u.name(l.strip()[0]).split()[0] for l in sys.stdin if l.strip()))"
```

`ru` is the shape to aim for: 43 Cyrillic keys and 14 Latin ones — its own
alphabet, plus the Latin letters a Russian brand is also written in. A file
whose keys are *all* Latin for a non-Latin language has been copied, not
derived.

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

Read it as a **skeleton function**: every line says "this character reduces to
that one". Two characters are confusable when they reduce to the same skeleton,
which is the relation the file actually encodes and the one to build:

```python
# confusables.txt -> skeleton(char), and the inverse: skeleton -> everything
import unicodedata as u
from collections import defaultdict

skeleton, family = {}, defaultdict(set)
for line in open("confusables.txt", encoding="utf-8"):
    line = line.split("#")[0].strip()
    if not line:
        continue
    src, tgt, *_ = [c.strip() for c in line.split(";")]
    if " " in tgt:            # multi-character skeletons are not homoglyphs
        continue
    s, t = chr(int(src, 16)), chr(int(tgt, 16))
    skeleton[s] = t
    family[t].add(s)

def lookalikes(ch):
    """Everything that could be mistaken for ch, in any script."""
    root = skeleton.get(ch, ch)          # ch's own skeleton, or ch itself
    return {c for c in family[root] | {root} if c != ch}
```

4,462 of the file's ~10,000 lines are single-character mappings; the rest are
multi-character skeletons (`rn` → `m`) and are not homoglyphs in this sense.

Then drive it from the alphabet, which is what makes the file this language's:

```python
alphabet = "αβγδεζηθικλμνξοπρστυφχψω"        # CLDR exemplars for the code
for ch in alphabet:
    row = sorted(lookalikes(ch) | set(nfd_variants(ch)))
    print(ch, *row)
```

Measured output, so you know what to expect:

```
ο  →  79 look-alikes: o σ ϭ о օ ס ه ٥ ھ ہ ە ۵ ० ০ ੦ ૦ ୦ ௦ ౦ ഠ ൦ ๐ ໐ ဝ ၀ …
α  →  23:             a ɑ а
```

Greek `ο` pulls Latin `o`, Cyrillic `о`, Armenian `օ`, Hebrew `ס`, Arabic `ه`
and a long tail of zero-shaped digits from a dozen scripts. Run the same
function on Cyrillic `а` and it returns `a ɑ α` — the attack on a Russian brand,
from the same table, because the key changed.

The tail of zero-shapes is why the filtering below matters: `੦` is a real
look-alike and no registrar will accept it next to Greek.

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
entry page names the correct form, which is the second column.

**Use the MediaWiki API, not the category page.** The HTML paginates at 200 and
the API does not:

```bash
curl -s 'https://en.wiktionary.org/w/api.php?action=query&list=categorymembers\
&cmtitle=Category:German%20misspellings&cmlimit=500&format=json' | jq -r '.query.categorymembers[].title'
# ab messen, Abszeße, Abszeßes, Ahmadinedschad, an's, Archillesferse, Armeisen, auf's …
```

`cmcontinue` in the response pages the rest. Fetch each title to get the
correct form it points at.

Most large Wikipedias also keep their own version of the list above; search
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

### When a source needs a browser

**Look for the API first.** Nearly everything in this catalogue is a static file
or a JSON endpoint, and a page that looks like it needs scraping usually does
not:

| Looks like scraping | Actually |
|---|---|
| Wiktionary category pages | `api.php?action=query&list=categorymembers` |
| OpenThesaurus | `openthesaurus.de/synonyme/search?q=<w>&format=application/json` |
| crates.io, NuGet, RubyGems | documented JSON APIs |
| kaikki, CLDR, WikiPron, stopwords-iso | plain files in a repo |

Reaching for a browser where an API exists costs you pagination, rate limits and
reproducibility for nothing.

**A browser is right when the data is rendered client-side** and there is no
endpoint behind it. In practice that is the lure-word research in step 3 — bank
and CERT advisory portals, phishing-report dashboards — rather than any
linguistic corpus.

In this harness, drive the real browser rather than fetching:

```
mcp__claude-in-chrome__tabs_create_mcp   → open a tab
mcp__claude-in-chrome__navigate          → go to the page
mcp__claude-in-chrome__get_page_text     → the rendered text, after JS
mcp__claude-in-chrome__find              → locate an element to click
mcp__claude-in-chrome__javascript_tool   → pull structured values out of the DOM
```

`get_page_text` after `navigate` is the whole job for a JS-rendered list —
`WebFetch` on the same URL returns the empty SPA shell. For something
repeatable, script Playwright against headless Chromium instead and commit the
script beside the data.

Rules that are not optional:

- **Cache the fetch to disk** and work from the cache. Re-running a parser is
  free; re-running a scrape is rude and gets you blocked.
- **One request at a time**, and stop at the first sign of a rate limit.
- **Respect robots.txt and the terms.** A tool for investigating abuse cannot
  be sloppy about how it collects.
- **Never authenticate to scrape.** If the data needs a login, it is not a
  source for a public dataset.
- **Record the URL and the date** at the top of the `.lst`. A scraped page is
  the one source that cannot be re-derived later.

### One source covers four files

Before working through the per-file lists: **Wiktionary, machine-readable,
covers synonyms, antonyms, homophones and misspellings in one download.**

[Wiktextract](https://github.com/tatuylonen/wiktextract) parses Wiktionary into
JSONL, and [kaikki.org](https://kaikki.org/) publishes the output for **500+
languages**. One object per word, one JSON per line:

```
https://kaikki.org/dictionary/downloads/<code>/<code>-extract.jsonl.gz
```

`de` is 2.8GB uncompressed, `es` 1.1GB — stream it, do not load it. The fields
that matter here:

| Field | Fills |
|---|---|
| `synonyms[].word` | `synonym.lst` |
| `antonyms[].word` | `antonym.lst` |
| `sounds[].ipa` | `homophone.lst`, by grouping words with equal IPA |
| `forms[]` with a `misspelling` tag, and `alt_of` | `misspelling.lst` |
| `word` | `word.lst` |

```bash
# every German word tagged as a misspelling, with what it is a misspelling of
zcat de-extract.jsonl.gz | jq -r '
  select(.forms[]?.tags[]? == "misspelling")
  | "\(.word) \(.alt_of[0].word // empty)"'
```

Use the edition in the language you are filling (`downloads/de/` is
de.wiktionary), not the English edition's coverage of it: a language's own
Wiktionary carries far more of its own misspellings and regional forms.

Licence is Wiktionary's — CC BY-SA. The tool is MIT.

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
| kaikki `forms[].tags` (above) | 500+ languages, structured | CC BY-SA |
| [AWB typo lists](https://en.wikipedia.org/wiki/Wikipedia:AutoWikiBrowser/Typos) | English, regex rules used by editing bots | CC BY-SA |
| [de: Liste von Tippfehlern](https://de.wikipedia.org/wiki/Wikipedia:Liste_von_Tippfehlern) | German, large and maintained | CC BY-SA |
| [GitHub Typo Corpus](https://github.com/mhagiwara/github-typo-corpus) | 350k edits, 15+ languages, mined from commits | check the repo |
| [Birkbeck spelling error corpora](https://www.dcs.bbk.ac.uk/~roger/corpora.html) | English, human error data | check |

The AutoWikiBrowser typo lists are what Wikipedia's own cleanup bots run, so
they are the errors that survive into published prose rather than the ones a
dictionary predicts. Most large Wikipedias keep an equivalent — German's is
`Wikipedia:Liste von Tippfehlern`, misspelling first with the correction in
parentheses. They are regex rules, not pairs, so take the literal alternations
and skip the patterned ones.

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
| [Open Multilingual WordNet](https://omwn.org/) | 60 wordnets, 49 languages, synsets | open, but **per-wordnet** — check each |
| kaikki `synonyms[].word` (above) | 500+ languages | CC BY-SA |
| [LibreOffice thesauri](https://github.com/LibreOffice/dictionaries) | ~20 languages, MyThes format | per-language |
| [OpenThesaurus](https://www.openthesaurus.de/) | German, large and current | LGPL / GPL |
| [ConceptNet](https://conceptnet.io/) | ~300 languages, weaker per-language | CC BY-SA |

The LibreOffice repo is one directory per locale (`de/`, `fr_FR/`, `en/`) each
holding `.aff` + `.dic` for spelling and `th_<locale>_v2.dat` + `.idx` for the
thesaurus. MyThes `.dat` is plain text and trivially parsed — an entry line
`word|n` followed by `n` sense lines of `(pos)|syn1|syn2|…`. It ships spelling
and synonyms for the same locale in one clone, which is why it is worth knowing
even though its language coverage is narrower than kaikki's.

### antonym.lst

Same shape as synonyms, same sources — Open Multilingual WordNet carries
antonymy as a relation, so it is one query rather than a second corpus.

### homophone.lst

Nothing ships a homophone list per language. Derive one: **words sharing an IPA
transcription are homophones by definition.**

| Source | Coverage | Licence |
|---|---|---|
| [WikiPron](https://github.com/CUNY-CL/wikipron) | 1.7M pronunciations, **165 languages**, mined from Wiktionary | Apache-2.0 tool, CC BY-SA data |
| kaikki `sounds[].ipa` (above) | 500+ languages, same origin | CC BY-SA |
| [CMU Pronouncing Dictionary](https://github.com/cmusphinx/cmudict) | English, 134k words | permissive — check |

WikiPron ships the scraped data in the repo, no tool run needed:

```
https://raw.githubusercontent.com/CUNY-CL/wikipron/master/data/scrape/tsv/<iso639-3>_<script>_broad.tsv
```

Filenames are `<lang>_<script>_<dialect?>_<broad|narrow>.tsv` — `deu_latn_broad.tsv`,
`ell_grek_broad.tsv`, `ben_beng_dhaka_broad.tsv`. Two tab-separated columns,
word then IPA with phones space-separated. German is 60,277 rows. Take `broad`
rather than `narrow`: narrow transcription encodes allophonic detail that
splits words a listener hears as identical, which is the opposite of what a
homophone list wants.

```bash
curl -sL .../deu_latn_broad.tsv \
  | awk -F'\t' '{k=$2; w[k]=w[k]" "$1} END{for(k in w) if (split(w[k],a," ")>1) print w[k]}'
```

Every bucket with more than one member is a homophone group. That is the whole
derivation, and it is the single highest-yield source in this catalogue for
languages nobody has curated.

German yields 2,332 groups from that one command — but look at them before
writing:

```
Slowakisch slowakisch          case only, worthless: names are case-folded
achtunddreissigmal achtunddreißigmal   orthographic, and a real confusion
Amen amen                      case only
```

**Drop the groups that differ only by case**, which is a large fraction and
contributes nothing to a scan of a case-insensitive namespace. Keep the
orthographic pairs — `ss`/`ß`, `ae`/`ä` — because those are exactly what someone
types when their keyboard cannot reach the other form.

Note the English file also holds symbol names — `dot .`, `dash -` — which no
corpus will give you. Those are typed-aloud spellings and have to be written by
hand.

### grapheme.lst and vowel.lst

| Source | Coverage | Licence |
|---|---|---|
| [CLDR exemplar characters](https://github.com/unicode-org/cldr-json) | per-locale alphabet | Unicode |
| Unicode `Script` property | what belongs to a script at all | Unicode |
| [PHOIBLE](https://phoible.org/) | phoneme inventories, 2,000+ languages | CC BY-SA |

```
https://raw.githubusercontent.com/unicode-org/cldr-json/main/cldr-json/
  cldr-misc-full/main/<code>/characters.json
```

`.main.<code>.characters.exemplarCharacters` is a bracketed string —
`[a b c … x y z]` — with `{}` around multi-character graphemes (Hungarian
`{cs}`, `{sz}`) and `-` denoting ranges. Parse those two, do not split on
whitespace and hope. The same object carries `auxiliary` (borrowed and accented
forms a name may still contain), `index`, and `punctuation`.

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

```
https://raw.githubusercontent.com/hermitdave/FrequencyWords/master/content/2018/<code>/<code>_50k.txt
```

`word count`, space separated, already sorted descending — `ich 5890279`,
`sie 3806767`. The `_50k` file is the top 50,000; the unsuffixed `<code>.txt`
beside it is the full list and is mostly corpus noise. Cut at a few thousand.

Hunspell's `.dic` counts the word list on line 1 and then gives one entry per
line as `word/FLAGS`, where the flags index affix rules in the neighbouring
`.aff`. Strip everything from `/` onward; expanding the affixes needs a
Hunspell binding and is rarely worth it for this file.

### stopword.lst

| Source | Coverage | Licence |
|---|---|---|
| [stopwords-iso](https://github.com/stopwords-iso/stopwords-iso) | **57 languages**, ISO 639-1 keyed, one JSON | MIT |

```
https://raw.githubusercontent.com/stopwords-iso/stopwords-iso/master/stopwords-iso.json
```

One JSON object, ISO 639-1 keys, array of words per language:
`{"de": ["ab","aber",…], "el": [...]}`. Permissively licensed and keyed by the
same codes this tree uses, so it is `jq -r '.<code>[]'` and done. There is no
reason to source this anywhere else.

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

### phonetics/homophone.lst — cross-language sound-squats

Not a language's data, so it does not live under `languages/`: a group like
`boutique boetiek butik` belongs to no single code. It is embedded straight from
the `.lst` by `datasets/phonetics.go` rather than imported into the database,
because every row in that database hangs off a language id.

Derived, not curated. Regenerate from WikiPron:

```bash
# fetch the broad transcriptions you want, e.g. deu_latn_broad.tsv, then
python3 scripts/crosshomophones.py <tsv-dir> datasets/phonetics/homophone.lst
```

The rule is: normalise the IPA (strip stress, length, syllable marks and
combining diacritics), group spellings by the result, and keep a group only
when it draws on **two or more languages** — a single-language group is an
ordinary homophone and `hs` already has it. Then filter to typeable spellings
(`[a-z][a-z0-9-]{2,}`), four characters minimum, two to six members: shorter
words and larger groups are phonetic sinks that collide with everything.

1.2M pronunciations across 15 languages yields about 7,600 groups. The header
of the file records the source and the filters; keep it accurate if you change
them.

### domains/private.lst — hosts that give subdomains away

The PRIVATE section of the public suffix list: the part submitted by domain
owners who delegate subdomains to third parties, as opposed to the ICANN
section, which is registries. `duckdns.org`, `github.io`, `herokuapp.com`,
`vercel.app`, `ddns.net`.

The split is a risk one, which is why it is a separate file rather than a flag
on the suffix list. A name under `.com` costs money and leaves a WHOIS trail; a
name under `duckdns.org` costs nothing, takes a minute and leaves no
registration record. The `fsd` algorithm generates only these, so a scan can ask
the cheap question first. ~3,000 entries against the ICANN section's ~7,000.

```bash
curl -s https://publicsuffix.org/list/public_suffix_list.dat |
  awk '/BEGIN PRIVATE/{p=1;next} /END PRIVATE/{p=0} p && !/^\/\// && NF && !/^[*!]/'
```

Wildcard (`*.`) and exception (`!`) rules are dropped: they are matching
instructions, not names anyone can take. Source is publicsuffix.org, MPL-2.0.

### packages/*.lst — popular library names

A different tree and a different job: these are the names a dependency-confusion
or combosquat scan works against, one per line, one file per registry. The repo
ships four — `npm`, `pypi`, `crates`, `rubygems` — at 30–50 names each, while
`sources/packages.lst` can already *probe* thirteen registries. **`nuget`,
`packagist`, `hex`, `conda`, `homebrew`, `pub`, `cpan`, `hackage` and
`cocoapods` have no corpus**, which is the gap to fill first.

Rank by downloads, not by opinion. A squat targets what people install.

**One API covers almost all of it.** [ecosyste.ms](https://packages.ecosyste.ms/)
indexes **100 registries** with a uniform interface:

```bash
curl -s 'https://packages.ecosyste.ms/api/v1/registries/npmjs.org/packages?sort=downloads&order=desc&per_page=100' \
  | jq -r '.[].name'
```

Swap the host segment for the registry: `pypi.org`, `rubygems.org`,
`crates.io`, `nuget.org`, `packagist.org`, `repo1.maven.org`,
`proxy.golang.org`, `hub.docker.com`. Package counts as of writing — npm 5.7M,
Go 2.2M, PyPI 921k, NuGet 834k, Maven 614k, Packagist 508k, crates 318k,
RubyGems 210k — so the sort matters more than the endpoint.

Per-registry alternatives, all verified:

| Ecosystem | Source | Notes |
|---|---|---|
| Python | [top-pypi-packages](https://hugovk.dev/top-pypi-packages/top-pypi-packages.min.json) | monthly dump, top 15,000 by real download counts, JSON and CSV |
| Rust | `https://crates.io/api/v1/crates?per_page=100&sort=downloads` | official, JSON, `.crates[].id` |
| C# / .NET | `https://azuresearch-usnc.nuget.org/query?q=&take=100&sortBy=totalDownloads-desc` | the endpoint nuget.org's own search uses; `.data[].id` |
| Ruby | ecosyste.ms `rubygems.org` | RubyGems has no public most-downloaded endpoint |
| npm | ecosyste.ms `npmjs.org` | the registry's own `/-/v1/search` needs a text query, so it cannot list a global top |

**C and C++ have no registry**, so there is no download ranking to sort. Use
what the package managers ship, which is the closest thing to a popularity
signal that exists:

```bash
# vcpkg — 2,858 ports
curl -s https://raw.githubusercontent.com/microsoft/vcpkg/master/versions/baseline.json \
  | jq -r '.default | keys[]'

# Conan Center — one directory per recipe
curl -s https://api.github.com/repos/conan-io/conan-center-index/contents/recipes \
  | jq -r '.[].name'
```

The intersection of the two is a better list than either alone: a library both
ecosystems bothered to package is one people actually build against.

**Take hundreds, not thousands.** These files are a corpus of plausible targets,
not a mirror of the registry — the shipped lists are 30–50 and the useful
ceiling is a few hundred. Past that you are scanning the long tail, where a name
collision is more likely to be a coincidence than a squat.

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