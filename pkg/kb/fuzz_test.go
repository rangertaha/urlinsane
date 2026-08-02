// Copyright 2024 Rangertaha. All Rights Reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package kb

import (
	"strings"
	"testing"
)

// hostile is the set of inputs a caller might arrive with: empty, blank,
// unprintable, outside the Basic Multilingual Plane, decomposed, joined,
// absurdly long, and a two-rune ligature that really is on a key.
var hostile = []string{
	"", " ", "\x00", "\n", "\t",
	"中", "🙂", "�",
	"é", "‍",
	"..", "-", "a b", "A", "ß", "لا",
	strings.Repeat("a", 5000),
}

// TestNothingPanics drives every exported entry point with those inputs, on
// layouts of each shape plus an empty one, and on forms the package does not
// know. None of it should do anything worse than return nothing.
func TestNothingPanics(t *testing.T) {
	boards := []*Layout{
		MustGet("kbdus"),
		MustGet("kbdgr"),
		MustGet("kbda1"),
		MustGet("kbdus").With(JIS),
		New("empty", "Empty", ANSI, nil),
	}

	for _, l := range boards {
		for _, s := range hostile {
			l.Adjacent(s)
			l.AdjacentWithin(s, -1)
			l.AdjacentWithin(s, 0)
			l.AdjacentWithin(s, 1e9)
			l.Shifted(s)
			l.Unshifted(s)
			l.AltGraphed(s)
			l.KeysFor(s)
			l.Types(s)
			l.Key(s)
			l.Type(l.Strokes(s))

			for _, other := range boards {
				l.Translate(s, other)
			}
		}

		_ = l.String()
		l.Rows()
		l.Languages()

		for _, f := range []Form{ANSI, ISO, JIS, Form(""), Form("bogus")} {
			_ = l.With(f).String()
		}
	}

	for _, s := range hostile {
		Get(s)
		ByLanguage(s)
		Find(s)
		ByKeys(s)
		ByString(s)
		Positioned(ANSI, s)
	}

	// Strokes a caller assembled rather than got from Strokes: a zero one, a
	// scan code no board has, and a modifier set outside the modelled range.
	MustGet("kbdus").Type([]Stroke{
		{},
		{SC: "ZZ"},
		{SC: "12", Mod: Mod(200)},
		{SC: "12", Mod: Base},
	})
}

// TestTypeSkipsStrokesItCannotHonour checks what the fuzzing above only proved
// harmless: unusable strokes contribute nothing, and the usable ones around
// them still come through.
func TestTypeSkipsStrokesItCannotHonour(t *testing.T) {
	us := MustGet("kbdus")

	got := us.Type([]Stroke{
		{SC: "12"},                // e
		{},                        // no key at all
		{SC: "ZZ"},                // a scan code the board does not have
		{SC: "12", Mod: Mod(200)}, // a modifier set it cannot be in
		{SC: "1E", Mod: Shift},    // A
	})
	if got != "eA" {
		t.Errorf("Type = %q, want %q", got, "eA")
	}
}

// TestUnboundedRadiusReachesEverything pins the other end of the range: a
// radius large enough takes in every key that types something, and nothing
// else — not the key asked about, and not the blank modifiers.
func TestUnboundedRadiusReachesEverything(t *testing.T) {
	us := MustGet("kbdus")

	blank := 0
	for _, k := range us.Keys {
		if k.Blank() {
			blank++
		}
	}

	got := len(us.AdjacentWithin("e", 1e9))
	want := len(us.Keys) - blank - 1 // everything but the blanks and "e" itself
	if got != want {
		t.Errorf("an unbounded radius reached %d keys, want %d", got, want)
	}
}
