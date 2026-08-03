---
title: Algorithms
parent: Reference
nav_order: 4
---

# Algorithms
{: .no_toc }

- TOC
{:toc}

A **variant algorithm** takes one name and returns the names a person or a
machine could reach instead of it. Each one models a *specific* failure: a
finger landing on the next key, a vowel heard rather than read, a bit flipping
in a resolver's memory, a package manager resolving an unscoped name. They are
deliberately narrow, because the value of the set is that every member is
independently justifiable — each exists because a measurement study found the
technique in the wild, and the
[bibliography](bibliography/) maps each one back to the paper that found it.

Run `urlinsane typo --list algorithms` for the shipped list.

## Using one algorithm

`-a` restricts generation to the ids you name; the target is positional and its
type is inferred from the string, so the same flag works across every target
kind:

```bash
urlinsane typo -a cs example.com              # domain
urlinsane typo -a nsc npm:left-pad            # package
urlinsane typo -a afx github.com/acme/tool    # repository
urlinsane typo -a hs bob@example.com          # email
urlinsane typo -a hr,hs,cm example.com        # several at once
```

An optional scope positional narrows what gets *varied* without changing what
the target *is* — `urlinsane typo username bob@example.com` varies only the
local part. To see which operators a run would compile without executing it:

```bash
urlinsane typo --explain -a hr example.com
```

## The generator contract

Every generator upholds four rules. The first two are enforced by
`TestGeneratorsAreRuneSafe` in `pkg/typo/unicode_test.go`; the last two by
`TestNeverReturnsItsInput` and `TestReturnsNoDuplicates`, which run in all
thirty-two plugin packages.

1. **Index runes, not bytes.** `for i, char := range token` yields a *byte*
   offset, so `token[:i]` cuts multi-byte characters in half. Domains escape
   this because they are punycoded before admission — usernames, packages,
   repositories and email local parts do not. Build on `runesOf` and
   `joinRunes` in `pkg/typo/runes.go`.
2. **Return variants in a stable order.** Deduplicating through a map and
   ranging it yields Go's randomised order, and that order reaches admission
   order in the engine, which decides which candidates survive a budget. An
   unstable order changes the content address of a scan.
3. **Never return the input.** A name is not a typo of itself, and every
   consumer downstream would have to filter it out again.
4. **Never invent a name from nothing.** An empty input yields nothing.

Character-level generators get 2–4 free by collecting through `uniq` in
`pkg/typo/runes.go`. The domain- and registry-shaped generators assemble their
output by appending to a slice, and use `variant.Clean` for the same guarantee.

`bf` is the one deliberate exception to rule 1: it models a bit flipping in
memory rather than a human mistyping, so byte-level output is the whole point.

## The set

| ID | Title | Applies to | Reads |
|---|---|---|---|
| `aci` | Adjacent Character Insertion | any | keyboard layout |
| `acs` | Adjacent Character Substitution | any | keyboard layout |
| `rar` | Repetition Adjacent Replacement | any | keyboard layout |
| `co` | Character Omission | any | — |
| `cr` | Character Repetition | any | — |
| `cs` | Character Swapping | any | — |
| `bf` | Bit Flipping | any | — |
| `hi` | Hyphen Insertion | any | — |
| `ho` | Hyphen Omission | any | — |
| `di` | Dot Insertion | any | — |
| `do` | Dot Omission | any | — |
| `dhs` | Dot Hyphen Substitution | any | — |
| `sep` | Separator Substitution | package, repo, username | — |
| `tos` | Token Order Swap | any | — |
| `sp` | Singular Pluralise | any | — |
| `cm` | Common Misspellings | any | `misspelling.lst` |
| `hs` | Homophone Substitution | any | `homophone.lst` |
| `xhs` | Cross-language Homophone | any | `phonetics/homophone.lst` |
| `hr` | Homoglyph Replacement | any | `homoglyph.lst` |
| `gi` | Grapheme Insertion | any | `grapheme.lst` |
| `gr` | Grapheme Replacement | any | `grapheme.lst` |
| `vs` | Vowel Swapping | any | `vowel.lst` |
| `cns` | Cardinal Substitution | any | `numeral.lst` |
| `ons` | Ordinal Substitution | any | `numeral.lst` |
| `tld` | Wrong TLD | domain | public suffix list |
| `sld` | Wrong Second-Level Domain | domain | public suffix list |
| `tli` | TLD Insertion | domain | public suffix list |
| `si` | Subdomain Insertion | domain | subdomain list |
| `fsd` | Delegated Subdomain | domain | PSL private section |
| `afx` | Affix Squatting | package, repo, username | — |
| `cb` | Combo Squatting | any | keyword list |
| `nsc` | Namespace Confusion | package, repo | — |

---

## Keyboard-driven

These read a `pkg/kb` layout and ask what is physically next to a key. Adjacency
is *geometric* — computed from key positions in the layout, not from a
hand-written neighbour table — so it is correct for AZERTY, QWERTZ and Dvorak
without anyone enumerating them. Because most Latin boards share QWERTY
geometry, the 203 shipped layouts collapse to about 30 distinct neighbour sets.

### `aci` — Adjacent Character Insertion

A neighbouring key struck *in addition to* the intended one, on either side.

```
abc  ->  sabc  asbc  qabc  aqbc  zabc  azbc  ...
```

Inserting a neighbour before a character and after the character behind it are
the same edit, so every interior position is reachable twice; the generator
deduplicates.

### `acs` — Adjacent Character Substitution

A neighbouring key struck *instead of* the intended one. The classic fat-finger
error, and the single most productive keyboard model.

```
abc  ->  sbc  qbc  zbc  wbc  avc  anc  agc  ahc  abx  abv  ...
```

### `rar` — Repetition Adjacent Replacement

A doubled character where the repeat landed on a neighbour instead — the hand
was already moving. Only fires on names that already contain a doubled
character.

```
aab  ->  ssb  qqb  zzb  wwb
```

---

## Character-level

No data, no language, no layout — pure string edits. These are the models in
Szurdi et al.'s *The Long "Taile" of Typosquatting*, and they account for most
registered typo domains.

### `co` — Character Omission

Each character dropped in turn. A key that did not register.

```
abcd  ->  bcd  acd  abd  abc
```

### `cr` — Character Repetition

Each character typed twice. A key that bounced.

```
abc  ->  aabc  abbc  abcc
```

### `cs` — Character Swapping

Each adjacent pair transposed — the two-finger race, where the second finger
lands first.

```
abcd  ->  bacd  acbd  abdc
```

### `bf` — Bit Flipping

One bit flipped in one octet of the name. This models **hardware**, not humans:
a cosmic ray or a failing DRAM cell corrupts a cached name in a resolver or
client, and the query goes somewhere else. Dinaburg showed at Black Hat 2011
that this reaches the internet at measurable scale, and it is why the algorithm
is byte-level rather than rune-level.

```
ab  ->  `b  cb  eb  ib  qb  Ab  !b  áb  ac  a`  af  aj  ar  aB
```

### `hi` / `ho` — Hyphen Insertion and Omission

A hyphen added at every gap including the ends, or every hyphen removed
individually and all at once.

```
abc    ->  -abc  a-bc  ab-c  abc-
a-b-c  ->  ab-c  a-bc  abc
```

### `di` / `do` — Dot Insertion and Omission

The same edit for dots, which in a domain moves a label boundary and therefore
changes what the registrable name is.

```
abc    ->  a.bc  ab.c
a.b.c  ->  ab.c  a.bc  abc
```

### `dhs` — Dot Hyphen Substitution

Dot and hyphen exchanged for one another. They are visually similar at small
sizes and adjacent in many layouts.

```
a.b-c  ->  a-b-c  a.b.c
```

### `sep` — Separator Substitution

Every separator normalised to each of the others, including to nothing. This
matters most in package registries, where `my-lib`, `my_lib`, `my.lib` and
`mylib` are four different packages that read as one.

```
a-b_c  ->  a-b-c  a_b_c  a.b.c  abc
```

### `tos` — Token Order Swap

The hyphen-separated tokens reordered. `foo-bar` and `bar-foo` are different
names that mean the same thing to a reader.

```
foo-bar-baz  ->  foo-baz-bar  bar-foo-baz  bar-baz-foo  baz-bar-foo  baz-foo-bar
```

### `sp` — Singular Pluralise

The name pluralised, or a plural made singular.

```
example  ->  examples
```

---

## Language-driven

These read the curated `.lst` data described in the
[datasets guide]({{ site.baseurl }}/guide/datasets/). Each runs once per
selected language, so the same algorithm produces different candidates for
`-l en` and `-l ru`. The examples below use a small fixture rather than the
shipped data.

### `cm` — Common Misspellings

A habitual misspelling exchanged for its correct form, and the reverse. Reads
`misspelling.lst`, which is a corpus of errors people actually make — not a
generated edit distance. English ships 4,256 pairs from Wikipedia's
machine-readable list.

```
seperate  ->  separate
```

### `hs` — Homophone Substitution

A word replaced by one that sounds the same. This is the **soundsquatting**
vector: the attack surface is anyone who hears a name rather than reads it,
including screen-reader users and anyone typing from dictation.

```
to  ->  two  too
```

### `xhs` — Cross-language Homophone

One sound spelled the way a *different* language writes it — `youtube` reaching
`yutup`, which is how a Turkish or Indonesian speaker would render those
syllables. Reads `phonetics/homophone.lst`, which belongs to no single language.
Valentim et al. found roughly 15% of such candidates already carry TLS
certificates.

```
foo  ->  phoo
```

### `hr` — Homoglyph Replacement

A character replaced by one that *looks* like it — the IDN homograph attack.
Reads `homoglyph.lst`, built from Unicode confusables (UTS #39) unioned with NFD
decomposition, so it covers both cross-script look-alikes (Cyrillic `а` for
Latin `a`) and accented forms.

```
look  ->  1ook  l0ok  lo0k
```

### `gi` / `gr` — Grapheme Insertion and Replacement

A character of the language's own alphabet inserted at each gap, or substituted
for each existing character. Reads `grapheme.lst`. Unlike `aci`/`acs` these are
alphabet-driven rather than layout-driven, which is what catches errors in
scripts whose keyboard the tool has no layout for.

```
ab   ->  aab  aba  bab  abb  cab  acb  abc
abc  ->  bbc  cbc  aac  acc  aba  abb
```

### `vs` — Vowel Swapping

Each vowel replaced by every other vowel of that language. Vowels carry less
information than consonants, so a swapped vowel is disproportionately likely to
go unnoticed.

```
ab  ->  eb  ib  ob  ub
```

### `cns` / `ons` — Cardinal and Ordinal Substitution

Digits and their spelled-out forms exchanged in both directions — `cns` for
cardinals, `ons` for ordinals. Reads `numeral.lst`, whose line shape is
`digit word ordinal`.

```
one2three  ->  12three  onetwothree
myfirstsite  ->  my1site
```

---

## Domain-shaped

These vary the *structure* of a domain rather than its characters, and are the
only algorithms allowed to touch the public suffix. Every character-level
algorithm deliberately preserves the suffix, so that changing it is one
algorithm's job and shows up as one algorithm's finding.

### `tld` — Wrong TLD

The public suffix replaced by every other known suffix, keeping any subdomain.

```
example.com  ->  example.net  example.co.uk  ...
```

### `sld` — Wrong Second-Level Domain

The second-level label swapped for a sibling under the same TLD. `tld` cannot
produce these: it would turn `example.co.uk` into `example.de`, a different
*country*. This keeps the country and changes the category, which is the
variation a reader checking "is this the UK site?" will not catch.

```
example.co.uk  ->  example.org.uk  example.ac.uk
```

### `tli` — TLD Insertion

A second suffix appended after the real one, so the true registrable domain is
the attacker's while the familiar name still reads left-to-right.

```
example.com  ->  example.com.net  example.com.co.uk
```

### `si` — Subdomain Insertion

A common label prepended. Harmless-looking, and the reason `www.example.com.
attacker.tld` works on a reader who stops scanning at the first familiar token.

```
example.com  ->  www.example.com  mail.example.com
```

### `fsd` — Delegated Subdomain

The name re-hosted under a provider that delegates subdomains to its users —
`github.io`, `blogspot.com`, the *private* section of the public suffix list.
These are registrable by anyone, for free, in seconds.

```
example.com  ->  example.github.io  example.com.github.io
```

---

## Registry- and namespace-shaped

Built for package registries, repositories and usernames, where the namespace
rather than the DNS hierarchy is what gets confused. This is the supply-chain
half of the tool.

### `afx` — Affix Squatting

Ecosystem and role brackets glued to the name — `py`, `python-`, `node-`,
`-js`, `-cli`, `-sdk`, `-official`. A reader takes them as metadata about the
same project rather than as a different name.

```
example  ->  pyexample  python-example  libexample  example2  example-js  example-official
```

### `cb` — Combo Squatting

The brand plus a lure keyword, joined both ways and in both orders. Kintis et
al. showed combosquatting domains outnumber typo variants and live far longer,
because nothing about them is a *mistake* — the victim reads the name and
believes it.

```
example  ->  example-login  examplelogin  login-example  loginexample
```

### `nsc` — Namespace Confusion

A scoped name confused for an unscoped one and the reverse — `@scope/pkg` for
`scope-pkg` for `pkg`. This is the shape behind dependency confusion, where a
resolver offered both an internal and a public name picks the wrong one.

```
scope/pkg     ->  pkg  scope-pkg  @scope/pkg
npm:left-pad  ->  @npm:left/pad  npm:left/pad
```

---

## Prior art

### Tools

| Tool | Language | What it contributes here |
|---|---|---|
| [dnstwist](https://github.com/elceef/dnstwist) | Python | The reference implementation of the character-level set, and the closest comparison for `co` `cs` `cr` `hr` `bf` |
| [URLCrazy](https://github.com/urbanadventurer/urlcrazy) | Ruby | `wrong_sld`, which `sld` implements; also the cardinal/ordinal idea |
| [ail-typo-squatting](https://github.com/typosquatter/ail-typo-squatting) | Python | `WrongSld`, and the broadest published algorithm list to check coverage against |
| [DomainFuzz](https://github.com/monkeym4ster/DomainFuzz) | Node | Domain-shaped variations |
| [pypi-squatting](https://github.com/typosquatter/pypi-squatting) | Python | Registry-side prior art for `afx` and `nsc` |

`sld` was the one domain-shaped generator present in both URLCrazy and
ail-typo-squatting that urlinsane had no answer for; it exists because of that
gap.

### Data and libraries

| Source | Used for |
|---|---|
| [Unicode UTS #39](https://www.unicode.org/reports/tr39/) confusables | `homoglyph.lst`, and so `hr` |
| [Unicode CLDR](https://github.com/unicode-org/cldr-json) exemplar characters | `grapheme.lst` and `vowel.lst` |
| [Public Suffix List](https://publicsuffix.org/) | `tld` `sld` `tli`, and its *private* section for `fsd` |
| [kbdlayout.info](http://kbdlayout.info/) | The 203 layouts behind `pkg/kb`, and so `aci` `acs` `rar` |
| [Hunspell dictionaries](https://github.com/wooorm/dictionaries) | `word.lst` where the licence permits — several are GPL-2.0-only and cannot be used here |
| [Wikipedia lists of common misspellings](https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/For_machines) | `misspelling.lst`, and so `cm` |
| [Birkbeck spelling error corpora](https://www.dcs.bbk.ac.uk/~roger/corpora.html) | Human error models behind `cm` |

### Papers

The [bibliography](bibliography/) has the full list with local PDFs and the
algorithm-to-paper mapping. The load-bearing ones:

- **Dinaburg**, *Bitsquatting: DNS Hijacking Without Exploitation* (Black Hat
  2011) — `bf`.
- **Szurdi et al.**, *The Long "Taile" of Typosquatting Domain Names* (USENIX
  Security 2014) — the character-level set.
- **Nikiforakis et al.**, *Soundsquatting: Uncovering the Use of Homophones in
  Domain Squatting* — `hs`.
- **Valentim et al.**, *X-squatter* (ACM TOPS 2024) — `xhs`, cross-language
  sound squatting.
- **Kintis et al.**, *Hiding in Plain Sight: A Longitudinal Study of
  Combosquatting Abuse* (CCS 2017) — `cb`.
- **Agten et al.**, *Seven Months' Worth of Mistakes* (NDSS 2015) — the
  domain-shaped set and the case for measuring registration rather than
  guessing.
- **Duan et al.**, *Towards Measuring Supply Chain Attacks on Package Managers*
  (NDSS 2021) — `afx`, `nsc`, `sep`.
