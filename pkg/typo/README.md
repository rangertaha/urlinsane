# Typo Algorithms

## Writing a generator

A generator takes a name and returns names. Two rules, both enforced by
`TestGeneratorsAreRuneSafe` in `unicode_test.go`:

**Index runes, not bytes.** `for i, char := range token` yields a *byte* offset,
so `token[:i]`, `token[i+1:]` and `string(token[i])` are only correct while
every character is one byte wide. Given "яндекс" they cut multi-byte characters
in half: eight of these generators did exactly that, and `CharacterOmission`
returned six variants, all six invalid UTF-8, which the engine then admitted as
node keys. Domains escape this because they are punycoded before admission —
usernames, packages, repositories and email local parts do not. Build on
`runesOf` and `joinRunes` in `runes.go` rather than slicing the string.

**Return variants in a stable order.** Collect through `uniq`, which dedupes
while keeping first-seen order. Several generators deduplicated through a map
and ranged it; that order reaches admission order in the engine, and admission
order decides which candidates survive a frontier or a budget — so an unstable
order here changes the content address of a scan, and two identical scans stop
being identical.

`uniq.add` also drops the input itself, which is not a typo of anything, and
empty results. `BitFlipping` is the one deliberate exception to the rune rule:
it models a bit flipping in a resolver's memory rather than a human
mistyping, so byte-level output is its whole point.

## The algorithms

`pkg/typo` implements the character-level and language-driven generators. The
domain- and registry-shaped ones (`tld`, `sld`, `tli`, `si`, `fsd`, `afx`, `cb`,
`nsc`, `sep`, `tos`, `xhs`) live in `internal/plugins/variant`, because they
need the public suffix list or a namespace model rather than a string edit.

The complete reference — all thirty-two, what each one models, the dataset it
reads, worked examples, and the tool or paper each came from — is
**[docs/reference/algorithms.md](../../docs/reference/algorithms.md)**. It is
kept there rather than here so that one page covers both packages.

| Implemented here | |
|---|---|
| `co` `cr` `cs` `bf` | character edits |
| `hi` `ho` `di` `do` `dhs` | separator edits |
| `sp` | inflection |
| `cm` `hs` `hr` `gi` `gr` `vs` `cns` `ons` | language-driven, reading `datasets/languages/<code>/*.lst` |
| `aci` `acs` `rar` | keyboard-driven, via `pkg/kb` |
