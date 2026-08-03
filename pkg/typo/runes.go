// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package typo

// The primitives in this package operate on characters, and a character is a
// rune, not a byte.
//
// Every generator here used to be written as `for i, char := range token` with
// bodies that then sliced `token[:i]`, `token[i+1:]` and `string(token[i])`.
// The index a range over a string yields is a *byte* offset, so those slices
// are only correct while every character is one byte wide. Given "яндекс" they
// cut multi-byte runes in half: CharacterOmission returned "\x8fндекс" — not a
// typo of anything, and not valid UTF-8 — and the graph admitted it as a node
// key. Domains escaped this because they are punycoded before admission, but
// usernames, packages, repositories and email local parts are not, and
// supporting non-Latin names is the point of the language datasets.
//
// Three of the algorithms (acs, aci, rar) were fixed in place, in their own
// plugin directories, leaving those copies correct and these ones wrong — which
// is why this file exists rather than a fourth private fix. New generators
// build on `runesOf` and `joinRunes` and inherit the correct behaviour instead
// of having to remember it.
//
// The other half of the same problem is order. Several generators deduplicated
// through a map and then ranged it, so the *order* of the variants they
// returned changed run to run. That order reaches admission order, and
// admission order decides which candidates survive a frontier or a budget — so
// it reaches the content address of the scan. `uniq` dedupes while keeping
// first-seen order, which is deterministic and costs nothing.

// runesOf is the one conversion every generator starts from.
func runesOf(token string) []rune { return []rune(token) }

// joinRunes builds a string from rune slices and strings in one pass, which is
// what a generator does at every position: prefix, something new, suffix.
func joinRunes(parts ...any) string {
	var b []rune
	for _, p := range parts {
		switch v := p.(type) {
		case []rune:
			b = append(b, v...)
		case rune:
			b = append(b, v)
		case string:
			b = append(b, []rune(v)...)
		}
	}
	return string(b)
}

// uniq collects variants, dropping duplicates and keeping the order they were
// first produced in.
//
// Deterministic by construction: a map-and-range dedupe returns Go's
// randomized iteration order, and that order propagates through admission into
// which candidates a bounded scan keeps.
type uniq struct {
	seen map[string]bool
	out  []string
}

func newUniq() *uniq { return &uniq{seen: map[string]bool{}} }

// add records a variant unless it is empty, already present, or identical to
// the token it came from. Returning the input as its own variant is not a typo,
// and every caller downstream would have to filter it out again.
func (u *uniq) add(origin, variant string) {
	if variant == "" || variant == origin || u.seen[variant] {
		return
	}
	u.seen[variant] = true
	u.out = append(u.out, variant)
}

func (u *uniq) tokens() []string { return u.out }

// insertEverywhere puts sep at every position in token, including before the
// first character and after the last.
//
// An n-character token has n+1 gaps. The generators here used to produce n of
// them: the loop ran once per character and the final iteration was spent on a
// special case that moved the separator to the *end* rather than adding the
// last interior gap, so "example" never yielded "exampl-e" and "abc" never
// yielded "ab-c". The gap between the last two characters is as reachable as
// any other, and for a two-character name it was the only interior one there
// is.
func insertEverywhere(token, sep string) []string {
	rs := runesOf(token)
	u := newUniq()
	for i := 0; i <= len(rs); i++ {
		u.add(token, joinRunes(rs[:i], sep, rs[i:]))
	}
	return u.tokens()
}
