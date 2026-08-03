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

## Researching and populating a language

`go run ./cmd/datasets languages` scaffolds a directory for every language
`pkg/kb` ships. Most are stubs: a comment line and nothing else. As of this
writing 30 of 111 are curated and 81 hold one placeholder line, so "add a
language" almost always means filling in a scaffold rather than making one.

Fill it in this order. The list is ordered by how much each file changes what
the scanner can find, so a language stopped halfway is still useful:

1. **`homoglyph.lst`** — the highest-value file by a wide margin, and the one a
   general search engine is worst at. Work from the script's Unicode blocks and
   from the confusables data (Unicode TR39), not from "letters that look
   similar" in prose. One line per base character, that character first:
   `a à á â ã ä å ɑ а ạ ǎ ă ȧ ӓ`. Include cross-script neighbours — Cyrillic `а`
   on a Latin `a` line is the entire point.
2. **`misspelling.lst`** — pairs, wrong first: `hwile while`. Sources that work:
   a Wikipedia "commonly misspelled words" page in that language, spellchecker
   correction lists, and keyboard-adjacent errors for that language's layout.
   German has 1,349 pairs and English 4,256; a few hundred is already useful.
3. **`synonym.lst`** — despite the name, not a thesaurus. These are the
   brand-adjacent words an attacker bolts onto a name: login, verify, secure,
   pay, invoice, support, delivery, account. One theme per line, in that
   language. The curated set runs to ~65 groups; match that shape.
4. **`vowel.lst`**, **`grapheme.lst`** — small, mechanical, from the script.
5. **`homophone.lst`**, **`numeral.lst`**, **`antonym.lst`** — lower value; fill
   if the language makes them easy.

While researching:

- **Check the script, not the language.** Azerbaijani is written in Latin,
  Cyrillic and Arabic depending on where; the file should carry the script the
  target audience actually types.
- **Verify a claim before adding it.** A wrong homoglyph produces variants that
  nobody would ever mistype, which costs scan time on every future run and is
  invisible once merged.
- **Do not translate the English file.** The English synonym groups reflect
  English phishing lures. Another language's lures are its own.

**Never rewrite entries that are already there.** Add to a file; pass existing
lines through untouched. The one time a bulk pass normalised existing entries it
destroyed data — see the editing rules above — so the safe shape for a
populating pass is: read the file, keep every existing line verbatim, append
what is new, deduplicate exact repeats only.

Then:

```bash
wc -l datasets/languages/<code>/*.lst   # sanity: did anything shrink?
go run ./cmd/datasets build datasets
urlinsane typo --list languages | grep <code>
urlinsane typo -a hr -l <code> -d 0 example.com   # homoglyphs actually fire
```

A file that shrank is the signal to stop and look, not to continue.

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
