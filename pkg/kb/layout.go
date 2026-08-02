// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package kb

import (
	"encoding/json"
	"sort"
	"strings"
)

// Locale is one of the language identities a layout is installed under.
// Windows ships a single driver for several locales — KBDUS backs the US
// English layout as well as the US variants of Chinese and Bulgarian — so a
// layout carries a list of these rather than a single language.
type Locale struct {
	// KLID is the eight hex digit Windows keyboard layout identifier, such
	// as "00000409".
	KLID string `json:"klid"`

	// Tag is the BCP-47 language tag for the locale, such as "en-US".
	Tag string `json:"tag,omitempty"`

	// Name is how Windows names the layout for this locale.
	Name string `json:"name,omitempty"`
}

// Layout is a keyboard layout: a set of physical keys and the text each one
// types. Layouts are immutable once loaded and safe for concurrent use.
//
// Get one from Get, New, or by unmarshalling; those position the keys and
// build the lookup tables. A Layout assembled as a struct literal has neither,
// and its keys will all sit on top of one another at the origin, so Key,
// KeysFor and Adjacent will not answer usefully.
type Layout struct {
	// ID is the driver name the layout is published under, lowercased, such
	// as "kbdus" or "kbdfr". It is the layout's canonical identifier.
	ID string `json:"id"`

	// Name is a human-readable name, such as "US" or "French".
	Name string `json:"name"`

	// File is the Windows driver the layout comes from, such as "KBDUS.DLL".
	File string `json:"file,omitempty"`

	// Form is the physical board the layout was drawn for, which decides
	// where the keys around the edges of the alphanumeric block sit. A
	// layout can be typed on a board of another shape — use With to say so.
	Form Form `json:"form"`

	// Locales are the language identities the layout is installed under.
	Locales []Locale `json:"locales,omitempty"`

	// Keys are the keys of the alphanumeric block, ordered top-left to
	// bottom-right.
	Keys []Key `json:"keys"`

	// bySC and byText index Keys; both are built once when the layout loads.
	bySC   map[string]int
	byText map[string][]int
}

// index positions each key, drops any that fall outside the alphanumeric
// block, and builds the lookup tables. It is called once, when the layout is
// loaded. Keys with no position are dropped rather than defaulted, so that a
// stray scan code cannot land on top of the number row and skew adjacency.
func (l *Layout) index() {
	// Filter into a new slice rather than in place. New takes the caller's
	// slice, and reusing its backing array would shuffle keys around under
	// them and leave a stale tail behind.
	kept := make([]Key, 0, len(l.Keys))
	for _, k := range l.Keys {
		p, ok := Positioned(l.Form, k.SC)
		if !ok {
			continue
		}
		k.Pos = p
		kept = append(kept, k)
	}
	l.Keys = kept

	byPosition(l.Keys)

	l.bySC = make(map[string]int, len(l.Keys))
	l.byText = make(map[string][]int)
	for i := range l.Keys {
		l.bySC[l.Keys[i].SC] = i
		for _, t := range l.Keys[i].Texts() {
			l.byText[t] = append(l.byText[t], i)
		}
	}
}

// With returns the layout as it would sit on a board of a different shape.
// The characters do not move between keys, but the keys move: an ISO layout
// typed on an ANSI board loses the extra key below the left Shift, and the
// backslash beside Enter moves up a row. Callers who know what hardware they
// are modelling can use this; the rest should take the layout as it comes.
//
// Asking for the shape a layout already has returns the layout itself, not a
// copy, since neither may be modified in any case.
func (l *Layout) With(form Form) *Layout {
	if form == l.Form {
		return l
	}

	keys := make([]Key, len(l.Keys))
	copy(keys, l.Keys)

	out := &Layout{
		ID:      l.ID,
		Name:    l.Name,
		File:    l.File,
		Form:    form,
		Locales: cloneLocales(l.Locales),
		Keys:    keys,
	}
	out.index()

	return out
}

// cloneLocales copies a locale list, so that two layouts sharing an origin do
// not share the slice they were built from.
func cloneLocales(locales []Locale) []Locale {
	if locales == nil {
		return nil
	}
	out := make([]Locale, len(locales))
	copy(out, locales)
	return out
}

// Languages returns the distinct primary language subtags the layout is used
// for, such as ["en"] for the US layout or ["de"] for German. Layouts shared
// across locales return one entry per language.
//
// Subtags come back lowercased. BCP-47 tags are case-insensitive and the
// dataset carries them as Windows spells them, so "ar-SA" and "ar" have to
// fold together here — otherwise they would count as two languages, and the
// one that kept its capitals would be unreachable through ByLanguage, which
// lowercases what it is asked for.
func (l *Layout) Languages() []string {
	seen := make(map[string]bool, len(l.Locales))
	out := make([]string, 0, len(l.Locales))
	for _, loc := range l.Locales {
		tag := strings.ToLower(loc.Tag)
		if i := strings.IndexByte(tag, '-'); i > 0 {
			tag = tag[:i]
		}
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		out = append(out, tag)
	}
	sort.Strings(out)
	return out
}

// Key returns the key at a scan code.
func (l *Layout) Key(sc string) (Key, bool) {
	i, ok := l.bySC[strings.ToUpper(sc)]
	if !ok {
		return Key{}, false
	}
	return l.Keys[i], true
}

// KeysFor returns every key that types s in some modifier state. It usually
// returns one key, but a layout may reach the same character two ways — on
// many European layouts a digit is on the number row unshifted and again as
// the shifted form of a letter key.
func (l *Layout) KeysFor(s string) []Key {
	idx := l.byText[s]
	out := make([]Key, 0, len(idx))
	for _, i := range idx {
		out = append(out, l.Keys[i])
	}
	return out
}

// Types reports whether the layout can type s at all.
func (l *Layout) Types(s string) bool { return len(l.byText[s]) > 0 }

// AdjacentKeys returns the keys physically surrounding k, nearest first,
// within the given radius in key units. Pass DefaultRadius for the eight or so
// keys a finger could slip onto.
func (l *Layout) AdjacentKeys(k Key, radius float64) []Key {
	type hit struct {
		key Key
		d   float64
	}

	var hits []hit
	for _, o := range l.Keys {
		if o.SC == k.SC || o.Blank() {
			continue
		}
		d := k.Pos.Distance(o.Pos)
		if d > radius {
			continue
		}
		hits = append(hits, hit{o, d})
	}

	sort.SliceStable(hits, func(i, j int) bool { return hits[i].d < hits[j].d })

	out := make([]Key, len(hits))
	for i, h := range hits {
		out[i] = h.key
	}
	return out
}

// Adjacent returns the characters typed by the keys physically next to the one
// that types s, nearest first and without duplicates.
//
// The characters come from the same modifier state s was found in, so
// Adjacent("e") gives lowercase neighbours and Adjacent("E") gives uppercase
// ones. This is what separates the result from a naive row-and-column grid:
// the rows of a real keyboard are staggered, so "e" sits a quarter of a key
// left of "d" but three quarters left of "s", and "d" comes back first.
func (l *Layout) Adjacent(s string) []string {
	return l.AdjacentWithin(s, DefaultRadius)
}

// AdjacentWithin is Adjacent with an explicit radius in key units. Widening it
// past about 1.9 starts to include keys two columns away.
func (l *Layout) AdjacentWithin(s string, radius float64) []string {
	seen := map[string]bool{s: true}
	out := []string{}

	for _, k := range l.KeysFor(s) {
		// Match the neighbours to the state s was typed in, so that a
		// shifted character yields shifted neighbours.
		state := Base
		for _, m := range k.Mods() {
			if o, _ := k.Text(m); o.Text == s {
				state = m
				break
			}
		}

		for _, n := range l.AdjacentKeys(k, radius) {
			o, ok := n.Text(state)
			if !ok || o.Text == "" || seen[o.Text] {
				continue
			}
			seen[o.Text] = true
			out = append(out, o.Text)
		}
	}

	return out
}

// Shifted returns what the keys that type s produce with Shift held — the
// uppercase of a letter, or the symbol printed above a digit. It returns
// nothing when s is only reachable shifted already.
func (l *Layout) Shifted(s string) []string {
	return l.mapState(s, func(m Mod) Mod { return m | Shift })
}

// Unshifted returns what the keys that type s produce without Shift, which
// undoes Shifted.
func (l *Layout) Unshifted(s string) []string {
	return l.mapState(s, func(m Mod) Mod { return m &^ Shift })
}

// AltGraphed returns what the keys that type s produce with AltGr held. Most
// keys on most layouts produce nothing, and the result is then empty.
func (l *Layout) AltGraphed(s string) []string {
	return l.mapState(s, func(m Mod) Mod { return m | AltGr })
}

// mapState finds the keys that type s, moves each to a different modifier
// state, and collects what they type there.
func (l *Layout) mapState(s string, to func(Mod) Mod) []string {
	seen := map[string]bool{s: true}
	out := []string{}

	for _, k := range l.KeysFor(s) {
		for _, m := range k.Mods() {
			o, _ := k.Text(m)
			if o.Text != s {
				continue
			}
			target := to(m)
			if target == m {
				continue
			}
			if t, ok := k.Text(target); ok && t.Text != "" && !seen[t.Text] {
				seen[t.Text] = true
				out = append(out, t.Text)
			}
		}
	}

	return out
}

// Rows groups the keys by physical row, top row first and each row ordered
// left to right.
func (l *Layout) Rows() [][]Key {
	var rows [][]Key
	var y float64 = -1

	for _, k := range l.Keys {
		if k.Pos.Y != y || rows == nil {
			rows = append(rows, nil)
			y = k.Pos.Y
		}
		rows[len(rows)-1] = append(rows[len(rows)-1], k)
	}

	return rows
}

// layoutWire is the JSON shape of a layout. It exists so that Layout can have
// JSON methods without recursing into itself. The dataset is stored as
// protocol buffers; see proto.go.
type layoutWire struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	File    string   `json:"file,omitempty"`
	Form    Form     `json:"form"`
	Locales []Locale `json:"locales,omitempty"`
	Keys    []Key    `json:"keys"`
}

// MarshalJSON writes the layout as JSON, which spells out everything it
// carries. It is for exporting and inspecting a layout, not for the dataset.
func (l Layout) MarshalJSON() ([]byte, error) {
	return json.Marshal(layoutWire{
		ID:      l.ID,
		Name:    l.Name,
		File:    l.File,
		Form:    l.Form,
		Locales: l.Locales,
		Keys:    l.Keys,
	})
}

// UnmarshalJSON reads a layout from JSON and indexes it, which positions the
// keys and discards any that sit outside the alphanumeric block.
func (l *Layout) UnmarshalJSON(b []byte) error {
	var w layoutWire
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}

	*l = Layout{
		ID:      w.ID,
		Name:    w.Name,
		File:    w.File,
		Form:    w.Form,
		Locales: w.Locales,
		Keys:    w.Keys,
	}
	l.index()

	return nil
}

// New assembles a layout from a set of keys. Keys outside the alphanumeric
// block are discarded, and the rest are positioned according to the form.
// The dataset generator uses it, as would a caller defining a layout of their
// own.
func New(id, name string, form Form, keys []Key) *Layout {
	l := &Layout{ID: id, Name: name, Form: form, Keys: keys}
	l.index()
	return l
}
