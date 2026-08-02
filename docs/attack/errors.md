---
title: 3 · Why names get mistyped
parent: Part I · The attack
nav_order: 3
---

# Why names get mistyped
{: .no_toc }

- TOC
{:toc}

An algorithm list is only as good as the error model behind it. If you generate
variants by "change one character to any other character" you produce a large
set that is mostly noise, and you still miss the cases that matter — because the
mistakes people actually make are not uniformly distributed over the alphabet,
and several of them are not letter substitutions at all.

This chapter is the error model. Every source of error below corresponds to
algorithms in the tool, named at the end of each section.

## Four sources, not one

It is worth separating these because the defences and the priorities differ.

| Source | Happens in | Example |
|---|---|---|
| **Motor** | the hands | `exmaple` — fingers out of order |
| **Perceptual** | the eyes | `rn` read as `m`; Cyrillic `а` read as Latin `a` |
| **Knowledge** | the memory | `recieve`, `dependancy` — the typist believes it is right |
| **Hardware** | the machine | a flipped bit in a resolver's cache |

Motor errors are self-correcting on average — the typist often notices. Knowledge
errors are not: someone who thinks a name is spelled a certain way will type it
that way every time, which makes those squats reliably profitable. Perceptual
errors are the worst case, because the victim who checks carefully still cannot
see the difference. And hardware errors are not human at all.

## Motor errors: the keyboard as geometry

Most typing errors are physical. The finger goes to the wrong place, arrives
twice, arrives early, or does not arrive. The four classic categories, from
decades of typing research, are **substitution** (wrong key), **insertion**
(extra key), **omission** (missed key) and **transposition** (right keys, wrong
order).

What makes these tractable is that they are not random. A substitution is
overwhelmingly likely to be a key *near* the intended one, so "near" has to be
defined properly — and that means physical geometry, not a grid of characters.

URLInsane measures adjacency from real key positions in key units, with the
row stagger included. On a US layout the neighbours of `e` come out as
`w r d 3 4 s`, in that order — `w` and `r` at exactly 1.0 units, then `d` at
1.031, then the digits at 1.118, then `s` at 1.25. That `d` beats `s` is the
quarter-key stagger showing through, and a naive grid model gets it wrong: it
would call `d` and `s` equally likely, or pick `s` because it is directly below
in the column.

Three further things a grid cannot express:

- **Boards differ in shape.** An ISO keyboard has a key between the left Shift
  and `z`; an ANSI one does not, and that key neighbours `a`. So `a`'s
  neighbours differ between two keyboards both called "UK English".
- **Keys differ in width.** The space bar is six units wide and reachable from
  anywhere along the bottom row.
- **Layouts differ entirely.** On QWERTZ, `z`'s neighbours are `t u h 6 7 g`. On
  AZERTY, `a`'s are `z q & é`. On Dvorak, `e`'s are `o u . q j p`. Assuming
  QWERTY is assuming your target's users are anglophone.

The tool covers **203 layouts across 110 languages**, taken from
[kbdlayout.info](https://kbdlayout.info)'s tables of the keyboard drivers
Windows ships. See [Keyboards]({{ site.baseurl }}/KB/) for the model in full.

There is also a fifth motor case that only exists in a multilingual world:
**typing on the wrong layout**. A name typed with the wrong keyboard selected in
the OS produces a different, registrable string — `google.com` typed on a US
board while the system is set to Russian yields `пщщпдуюсщь`, and on French
AZERTY it yields `google:co,` with not one letter changed, only the punctuation
moved. That last case is the instructive one: an English speaker looking at
`google:co,` sees garbage; a French user who has done it sees a familiar
accident.

*Algorithms:* `acs` (adjacent substitution), `aci` (adjacent insertion), `co`
(omission), `cr` (repetition), `cs` (swap/transposition), `rar` (repetition then
adjacent replacement), `gi`/`gr` (grapheme insertion and replacement, for
substitutions outside the neighbourhood).

## Perceptual errors: things that look the same

Reading is not character-by-character. It is pattern matching over word shapes,
and it is fast because it skips detail — which is exactly the property an
attacker exploits.

**Within one script**, some pairs are simply hard: `l`/`I`/`1`, `0`/`O`,
`rn`/`m`, `vv`/`w`, `cl`/`d`. In many sans-serif UI fonts, lowercase `l` and
capital `I` are rendered by identical glyphs. There is no amount of care that
distinguishes them.

**Across scripts** it is worse. Unicode contains many characters that are
visually indistinguishable from Latin letters in ordinary rendering: Cyrillic
`а е о р с у х`, Greek `ο ν`, and a long tail of mathematical and fullwidth
forms. A domain using them registers as punycode (`xn--...`) and renders as the
original. Browsers apply mixed-script heuristics to decide when to show the
punycode instead, and those heuristics are per-browser, incomplete, and
regularly bypassed by using a *single* script consistently — a wholly Cyrillic
string that happens to spell a Latin brand triggers no mixed-script warning at
all.

This is the case where enumeration is genuinely necessary. You cannot eyeball a
homoglyph attack. You can only generate the confusable set and check it.

*Algorithms:* `hr` (homoglyph replacement), plus the `idn` operator which
surfaces the punycode form of anything generated.

## Knowledge errors: things people believe

Some misspellings are not slips. The typist is typing what they think the word
is, and they will do it again tomorrow. These come in several flavours:

**Common misspellings** of ordinary words — `recieve`, `seperate`, `definately`
— which are language-specific and have to be curated per language rather than
generated. URLInsane ships misspelling lists as data, currently about 60,000
entries across all languages.

**Homophones** — words that sound identical: `to`/`too`/`two`, `site`/`sight`,
`by`/`buy`. These matter more than they used to, because names increasingly
travel by voice: dictation, voice assistants, and someone reading a URL aloud
over a phone. Someone who has only ever *heard* a name has no way to choose
between its homophones.

**Number words** — `4`/`four`, `2nd`/`second`. Brands with digits in them are
squattable in both directions, and the spelled-out form often reads as more
official than the original.

**Singular and plural** — `example.com` and `examples.com`. A one-character
change that a reader's eye actively normalises away, because English readers
routinely ignore trailing `s` when skimming.

**Second-language patterns.** A non-native speaker's errors are not the same as
a native speaker's. They reflect their first language's orthography: which
letters double, which vowels are ambiguous, where accents belong. A brand
targeting a market gets squatted with that market's error patterns, which is why
the language datasets are per-language rather than a single English list.

*Algorithms:* `cm` (common misspellings), `hs` (homophones), `cns`/`ons`
(cardinal and ordinal number words), `sp` (singular/plural), `vs` (vowel
swapping).

## Structural errors: the syntax of names

Some errors are not about characters at all but about the *structure* of the
name — the dots, hyphens and separators that carry meaning.

- **Missing dot**: `wwwexample.com`, or `mailexample.com` for
  `mail.example.com`. The dot is a small target next to the space bar and it is
  the boundary that determines who owns the name.
- **Extra dot**: `exa.mple.com`, which is a different registrable domain.
- **Hyphen added or dropped**: `my-example.com` versus `myexample.com`. Both
  look intentional; neither is obviously wrong.
- **Dot for hyphen**: `my.example.com` versus `my-example.com` — this one
  changes who controls the name entirely, from a subdomain of `example.com` to a
  separate registration.
- **Wrong TLD**: `example.co` for `example.com`, or `.cm`, or `.om`. The `.com`
  neighbourhood in the suffix list is large and cheap.
- **Subdomain insertion**: `login.example.com.attacker.net`, exploiting address
  bars that truncate from the left.
- **Separator swap** in package names: `-` for `_` for `.`, which several
  registries normalise inconsistently.

*Algorithms:* `do`, `di`, `hi`, `ho`, `dhs`, `tld`, `si`, `sep`.

## Composition errors: no error at all

Combosquatting deserves its own heading because it breaks the model this chapter
has been building. There is no error. `example-support.com` is not a mistyping
of anything; it is a plausible name that the organisation might well have
registered and did not.

It works on trust rather than on slips, and the literature (Kintis et al., CCS
2017) finds it both more common and longer-lived than typo variants, because
there is nothing for a careful reader to catch. The same applies to the package
world in the form of affix squatting — `acme-utils`, `node-acme`, `acme-js` —
where ecosystem naming conventions make the fake name look *more* idiomatic than
the real one.

*Algorithms:* `cb` (combo squatting), `afx` (affix squatting), `nsc` (namespace
confusion).

## Machine errors: bitsquatting

The last source is not a person.

A bit in a DRAM cell flips — from cosmic rays, from heat, from a marginal
module. If the flipped bit is in a cached hostname, the resolver looks up a
name one bit away from the one requested. `example.com` becomes `axample.com`,
`egample.com`, `examtle.com` and so on: 32 or so single-bit neighbours per
character, most of them not typeable and only some of them valid hostnames.

Artem Dinaburg demonstrated at Black Hat 2011 that this is not theoretical.
Registering bit-flip variants of popular domains produced real traffic, from
real hosts, at a steady rate — and the affected hosts were disproportionately
mobile and embedded devices with no ECC memory. At internet scale, a rare event
per machine is a constant stream in aggregate.

This is why the `bf` algorithm's output looks like nonsense. It is not modelling
a human at all, and it is the one algorithm whose results you should not sanity
check by asking "would someone type that?"

*Algorithm:* `bf`.

## What this implies for scanning

Three practical consequences, which the rest of the book builds on:

1. **Algorithm choice is a threat model.** Running everything is the default
   and is often right for a one-off audit, but if you are monitoring a
   consumer brand, `cm`/`hs`/`sp` matter more than `bf`; if you are monitoring
   a package, `afx`/`nsc`/`sep` matter more than keyboard adjacency.
2. **Language is not optional.** A tool that generates English keyboard slips
   for a Russian brand is doing a fraction of the job.
3. **Generation is cheap and observation is not.** Producing 10,000 candidates
   costs microseconds; finding out which exist costs 10,000 network round trips
   against rate-limited services. Every design decision in
   [Part III](../../internals/) about depth, budgets and caching comes from that
   asymmetry.

---

Next: **[What defenders do about it](../defending/)**.
