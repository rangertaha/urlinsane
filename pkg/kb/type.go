// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package kb

import "strings"

// Stroke is a key being pressed: which switch, and what was held down with it.
// It says nothing about what gets typed — that depends on the layout the
// keystroke lands on, which is the whole point.
type Stroke struct {
	// SC is the scan code of the key pressed.
	SC string

	// Mod is the set of modifiers held.
	Mod Mod
}

// Zero reports whether the stroke stands for a character no key could produce.
func (s Stroke) Zero() bool { return s.SC == "" }

// Strokes returns the keystrokes that type s on this layout, one per rune.
// A rune the layout cannot type yields a zero Stroke, so the result always has
// as many entries as s has runes.
//
// Where a character can be reached more than one way the plainest is chosen:
// unmodified before shifted, shifted before AltGr.
func (l *Layout) Strokes(s string) []Stroke {
	out := make([]Stroke, 0, len(s))
	for _, r := range s {
		out = append(out, l.stroke(string(r)))
	}
	return out
}

// stroke finds the simplest keystroke that types t.
func (l *Layout) stroke(t string) Stroke {
	for _, k := range l.KeysFor(t) {
		// Mods runs from the plainest state to the most encumbered, so the
		// first hit is the one a typist would actually use.
		for _, m := range k.Mods() {
			if o, _ := k.Text(m); o.Text == t {
				return Stroke{SC: k.SC, Mod: m}
			}
		}
	}
	return Stroke{}
}

// Type returns what the given keystrokes produce on this layout. Strokes that
// this layout has no key for, or whose key types nothing in that modifier
// state, contribute nothing.
//
//	ru.Type(us.Strokes("hello"))  // "руддщ"
//
// That pairing is the interesting one: the keys someone actually pressed,
// read through a different layout.
func (l *Layout) Type(strokes []Stroke) string {
	var b strings.Builder

	for _, s := range strokes {
		if s.Zero() {
			continue
		}
		k, ok := l.Key(s.SC)
		if !ok {
			continue
		}
		if o, ok := k.Text(s.Mod); ok {
			b.WriteString(o.Text)
		}
	}

	return b.String()
}

// Translate returns what s becomes when the keystrokes that type it here land
// on another layout instead — someone typing a familiar word with the wrong
// layout selected:
//
//	us.Translate("hello", ru)   // "руддщ"
//	us.Translate("google", ru)  // "пщщпду"
//
// Characters this layout cannot type, and keys the other layout leaves bare,
// are passed through unchanged rather than dropped, so a domain keeps its
// dots. Use Strokes and Type together when you would rather see exactly what
// fell through.
//
// Two things follow from that which are easy to assume away:
//
// The result is not always the same length as the input. A handful of layouts
// have ligature keys that type two characters at once — the Arabic boards
// reach "لا" on one key — so a rune can come back as two.
//
// And it does not round-trip unless this layout can type every character
// given. A character that falls through untouched here is still an ordinary
// character on the way back, and will be translated then: Russian cannot type
// "abc", so it passes through, and translating the result back through
// Russian yields "фис" rather than "abc".
//
// A key that is dead on the far layout contributes its accent, since that is
// what the dataset records; it would really arm the accent instead.
func (l *Layout) Translate(s string, as *Layout) string {
	var b strings.Builder

	for _, r := range s {
		t := string(r)

		stroke := l.stroke(t)
		if stroke.Zero() {
			b.WriteString(t)
			continue
		}

		k, ok := as.Key(stroke.SC)
		if !ok {
			b.WriteString(t)
			continue
		}

		o, ok := k.Text(stroke.Mod)
		if !ok || o.Text == "" {
			b.WriteString(t)
			continue
		}

		b.WriteString(o.Text)
	}

	return b.String()
}
