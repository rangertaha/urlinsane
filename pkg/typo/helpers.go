// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package typo

import (
	"strings"
)

// characterDeletion removes one occurrence of character at a time, then all of
// them at once.
func characterDeletion(token string, character string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i, r := range rs {
		if character == string(r) {
			u.add(token, joinRunes(rs[:i], rs[i+1:]))
		}
	}
	u.add(token, strings.ReplaceAll(token, character, ""))
	return u.tokens()
}

// characterReplace substitutes one occurrence of character at a time, then all
// of them at once.
//
// The per-position branch used to be a copy of characterDeletion's and dropped
// the character instead of writing the replacement, so the only substitution
// this ever produced was the replace-all on the last line. Reached through
// DotHyphenSubstitution, that meant `dhs` on "one.two.three" returned two
// *dot omissions* — variants `do` already generates, mis-attributed to another
// algorithm — and never the single-separator swaps it documents.
func characterReplace(token string, character, replacement string) (tokens []string) {
	rs := runesOf(token)
	u := newUniq()
	for i, r := range rs {
		if character == string(r) {
			u.add(token, joinRunes(rs[:i], replacement, rs[i+1:]))
		}
	}
	u.add(token, strings.ReplaceAll(token, character, replacement))
	return u.tokens()
}

// PrefixInsertion creates tokens by prepending each prefix from the given
// list to the specified token. Example:
// Inputs:
//
//	prefixes = ["www", "ftp", "shop"]
//	token = "example"
//
// Outputs: ["wwwexample", "ftpexample", "shopexample"]
func PrefixInsertion(token string, prefixes ...string) (tokens []string) {
	// Through uniq like every other generator in this package. Appending
	// straight to the result meant an empty prefix -- or a list carrying the
	// same prefix twice -- produced the token itself, which is not a typo of
	// anything, and every caller downstream had to filter it out again.
	u := newUniq()
	for _, prefix := range prefixes {
		u.add(token, prefix+token)
	}
	return u.tokens()
}

// SuffixInsertion creates tokens by appending each suffix from the provided
// list to the end (right side) of the given token. Example:
// Inputs:
//
//	suffixes = ["com", "net", "io"]
//	token = "example"
//
// Outputs: ["examplecom", "examplenet", "exampleio"]
func SuffixInsertion(token string, suffixes ...string) (tokens []string) {
	// Through uniq, for the same reason PrefixInsertion is.
	u := newUniq()
	for _, suffix := range suffixes {
		u.add(token, token+suffix)
	}
	return u.tokens()
}

func numeralMap(data map[string][]string, pos int) (words map[string]string) {
	words = make(map[string]string)

	for num, tokens := range data {
		for i, token := range tokens {
			if i == pos {
				words[num] = token
				// words[token] = num
			}
		}
	}

	return
}

// adjacentCharacters returns the characters neighbouring char on a keyboard
// given as rows of characters.
//
// Rows are ragged — "qwertyuiop" is ten keys and "zxcvbnm" is seven — so the
// key above or below a given column may not exist. The previous version indexed
// layout[r-1][c] and layout[r+1][c] unguarded and panicked with index out of
// range on any character whose column exceeded the neighbouring row's length:
// 'p' at column 9 of the top row reached column 9 of a nine-character row.
// AdjacentCharacterSubstitution("example") with an ordinary QWERTY layout was
// enough to crash the process.
//
// Rows are compared as runes, so a Cyrillic or Greek layout works the same way
// a Latin one does.
//
// This models a keyboard as a grid, which is a poor model — real keys are
// staggered and vary in width. pkg/kb measures adjacency from key geometry
// instead, and the acs, aci and rar plugins use it. This remains for callers of
// pkg/typo that supply their own layout.
func adjacentCharacters(char string, layout ...string) (chars []string) {
	chars = []string{}
	rows := make([][]rune, len(layout))
	for i, row := range layout {
		rows[i] = []rune(row)
	}

	// at reports the character at (r, c), and whether that key exists.
	at := func(r, c int) (string, bool) {
		if r < 0 || r >= len(rows) || c < 0 || c >= len(rows[r]) {
			return "", false
		}
		return string(rows[r][c]), true
	}

	for r := range rows {
		for c := range rows[r] {
			if char != string(rows[r][c]) {
				continue
			}
			for _, n := range [][2]int{{r - 1, c}, {r + 1, c}, {r, c - 1}, {r, c + 1}} {
				if s, ok := at(n[0], n[1]); ok && s != " " {
					chars = append(chars, s)
				}
			}
		}
	}
	return chars
}

// similarChars returns homoglyphs, characters that look alike from other languages
func similarChars(key string, data map[string][]string) (chars []string) {
	chars = []string{}
	char, ok := data[key]
	if ok {
		chars = append(chars, char...)
	}
	return chars
}
