---
title: 8 · Algorithms
parent: Part II · Using URLInsane
nav_order: 4
---

# Algorithms
{: .no_toc }

- TOC
{:toc}

An algorithm generates plausible variants of a name. There are 29, each one
modelling a specific way a name goes wrong — the error taxonomy behind them is
[chapter 3](../../attack/errors/).

```bash
urlinsane typo --list algorithms      # what this build has
```

## The full set

**Applies to** is blank where an algorithm binds by *capability* rather than by
type: those run on any nameable node, domain or package or handle alike.

| ID | Name | Applies to | What it does |
|---|---|---|---|
| `aci` | Adjacent Character Insertion | any | Insert a character adjacent on the keyboard: `googhle` |
| `acs` | Adjacent Character Substitution | any | Replace a character with a keyboard neighbour: `ezample` |
| `afx` | Affix Squatting | package, repo, username | Add a plausible prefix or suffix: `node-acme`, `acme-js` |
| `bf` | Bit Flipping | any | Flip one bit of a character — bitsquatting |
| `cb` | Combo Squatting | any | Append or prepend a common keyword: `acme-login` |
| `cm` | Common Misspellings | any | Apply a curated misspelling for the language |
| `cns` | Cardinal Substitution | any | Swap a number for its cardinal word and back: `file2` ⇄ `filetwo` |
| `co` | Character Omission | any | Drop a character: `gogle` |
| `cr` | Character Repetition | any | Double a character: `gooogle` |
| `cs` | Character Swapping | any | Transpose two adjacent characters: `examlpe` |
| `dhs` | Dot Hyphen Substitution | any | Swap dots and hyphens: `my-acme.com` ⇄ `my.acme.com` |
| `di` | Dot Insertion | any | Insert a period: `exa.mple` |
| `do` | Dot Omission | any | Remove a period: `wwwacme.com` |
| `gi` | Grapheme Insertion | any | Insert a grapheme from the language's alphabet |
| `gr` | Grapheme Replacement | any | Replace a grapheme with another from the alphabet |
| `hi` | Hyphen Insertion | any | Insert a hyphen: `ex-ample` |
| `ho` | Hyphen Omission | any | Remove a hyphen: `oneforall` |
| `hr` | Homoglyph Replacement | any | Replace a character with one that looks the same |
| `hs` | Homophone Substitution | any | Replace a word with one that sounds the same |
| `nsc` | Namespace Confusion | package, repo | Move a name between namespaces or scopes |
| `ons` | Ordinal Substitution | any | Swap a number for its ordinal word and back |
| `rar` | Repetition Adjacent Replacement | any | Double a character, then replace the double with a neighbour: `gppgle` |
| `sep` | Separator Substitution | package, repo, username | Swap the separator a registry allows: `-` ⇄ `_` ⇄ `.` |
| `si` | Subdomain Insertion | domain | Insert a subdomain label |
| `sld` | Wrong Second-Level Domain | domain | Swap the second level under a ccTLD: `bbc.co.uk` → `bbc.org.uk` |
| `sp` | Singular Pluralise | any | Make a word singular or plural |
| `tld` | Wrong TLD | domain | Substitute a different public suffix |
| `tli` | TLD Insertion | domain | Append a suffix, making the whole name a subdomain: `example.com.br` |
| `vs` | Vowel Swapping | any | Swap one vowel for another: `ixample` |

## Selecting them

```bash
urlinsane typo -a cs acme.com              # only transpositions
urlinsane typo -a cs,co,acs acme.com       # three of them
urlinsane typo -a '^bf' acme.com           # everything except bit flipping
urlinsane typo acme.com                    # everything (the default)
```

`^id` excludes. You can check what a selection actually compiled to without
running anything:

```console
$ urlinsane typo --explain -a cs,co acme.com | sed -n '/operators/,$p'
operators
  co             on nameable           where in-seed-closure
  cs             on nameable           where in-seed-closure
  decompose.domain on domain
  …
```

Note the `where in-seed-closure` condition. Variant operators only fire on nodes
inside the seed closure, which is what stops the scan from generating variants
of every nameserver it happens to discover. See
[Limits](../../internals/limits/).

## Which ones to run

Running all 29 is the default and is right for a one-off audit where you have
time. For anything repeated, pick by threat model.

**A consumer-facing brand.** Human error dominates, and reading errors dominate
over typing errors:

```bash
urlinsane typo -a cs,co,acs,aci,cr,vs,cm,hs,sp,hr,tld acme.com
```

**A package or a repo.** Convention exploitation dominates; keyboard geometry
barely matters, because the name is usually copied rather than typed:

```bash
urlinsane typo -a afx,nsc,sep,cb,co,cs,hr npm:acme-utils
```

**A quick triage pass.** The four highest-yield generators against domains, and
fast:

```bash
urlinsane typo -a co,cs,acs,tld -d 1 acme.com
```

**Everything except the noisy one.** `bf` produces the largest and least
human-plausible set:

```bash
urlinsane typo -a '^bf' acme.com
```

## Notes on individual algorithms

**`bf` (bit flipping)** does not model a human. Its output looks like nonsense
because the error source is memory corruption in a resolver, not a typist —
see [chapter 3](../../attack/errors/#machine-errors-bitsquatting). It generates
a lot. Exclude it unless you specifically want bitsquat coverage.

**`hr` (homoglyphs)** is the algorithm you cannot replace with careful reading.
Its output often looks *identical* to the seed in a terminal. Pair it with the
`idn` operator's punycode output to see what is really there.

**`cm`, `hs`, `vs`, `gi`, `gr`, `cns`, `ons`** are **language-driven**: they read
vocabulary out of the dataset database, and what they generate depends on which
languages that database carries. An empty or stale dataset makes them generate
little or nothing, silently. If these look quiet, check
[Datasets](../datasets/).

**`acs`, `aci`, `rar`** are **keyboard-driven**, and the keyboard is geometry
rather than a grid of rows — `e`'s neighbours on a US layout come out as
`w r d 3 4 s`, in distance order. On QWERTZ, AZERTY or Dvorak the answers differ
entirely. [Keyboards]({{ site.baseurl }}/KB/) has the model.

**`tld`** substitutes from the public suffix list, which is large. It is the
main reason a default scan produces so many candidates against a domain target.

**`si` (subdomain insertion)**, **`nsc` (namespace confusion)** and **`tli`
(TLD insertion)** model *structural* attacks rather than errors: the attacker is
not imitating a typo but exploiting how a name is parsed.

**`tli`** is the one with no misspelling in it at all. `example.com.br` contains
the target name in full, spelled correctly, and is a subdomain of somebody
else's registration — so it defeats the check most people actually perform,
which is "does the address contain the name I expect". It is also what a
truncating mobile address bar shows first.

**`sld`** exists because `tld` cannot produce it. `tld` replaces the whole
public suffix, so `bbc.co.uk` becomes `bbc.de` — a different country. `sld`
keeps the country and changes the category: `bbc.org.uk`, `bbc.ac.uk`. A reader
checking "is this the UK site?" gets the right answer and still lands on the
wrong name.

{: .todo }
> `--language` and `--keyboard` are registered flags but do not yet reach the
> plan, so language and keyboard selection is currently whatever the dataset
> holds rather than something you choose per run.

## Cost

Generation is free; observation is not. One algorithm at depth 1 against
`example.com` took 43 seconds in the example in
[chapter 6](../first-scan/) — essentially all of it network. Doubling the
algorithm count roughly doubles the candidate set, and the candidate set is what
drives the number of DNS and WHOIS calls.

The cheap way to explore is to generate without observing much: `-d 1` keeps the
scan on the variants themselves and their immediate records, rather than walking
out through nameservers and address space.

---

Next: **[Observation and depth](../observing/)**.
