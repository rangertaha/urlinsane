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
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Mod is a set of held modifiers, as a bitmask. Shift, AltGr and Caps Lock are
// modelled because between them they cover what almost every layout does.
//
// Two things are left out, and they are not all control codes. Ctrl and Alt on
// their own do produce control codes rather than text, and go. But a couple of
// layouts reach real characters through a modifier this package has no bit
// for: the Japanese board puts its kana on a Kana lock, and the Canadian
// Multilingual Standard board hangs a level off VK_OEM_8. Those characters are
// absent from the dataset — kbdjpn is Latin-only here, and kbdcan is missing
// its ¹ ² ³ ¼ level. The generator reports what it dropped, so a rebuild says
// so rather than letting it pass unnoticed.
type Mod uint8

const (
	// Base is the key pressed on its own.
	Base Mod = 0
	// Shift is either Shift key held down.
	Shift Mod = 1 << 0
	// AltGr is the right Alt key on layouts that treat it as AltGr, or the
	// Ctrl+Alt combination that stands in for it elsewhere.
	AltGr Mod = 1 << 1
	// Caps is Caps Lock engaged.
	Caps Mod = 1 << 2
)

// mods lists every state a key can produce text in, in a stable order.
var mods = []Mod{
	Base,
	Shift,
	Caps,
	Shift | Caps,
	AltGr,
	Shift | AltGr,
	AltGr | Caps,
	Shift | AltGr | Caps,
}

var modNames = map[Mod]string{
	Base:                 "base",
	Shift:                "shift",
	Caps:                 "caps",
	Shift | Caps:         "shiftcaps",
	AltGr:                "altgr",
	Shift | AltGr:        "shiftaltgr",
	AltGr | Caps:         "altgrcaps",
	Shift | AltGr | Caps: "shiftaltgrcaps",
}

// String returns the name this modifier set is stored under.
func (m Mod) String() string {
	if s, ok := modNames[m]; ok {
		return s
	}
	return fmt.Sprintf("mod(%d)", uint8(m))
}

// parseMod is the inverse of Mod.String.
func parseMod(s string) (Mod, bool) {
	for m, name := range modNames {
		if name == s {
			return m, true
		}
	}
	return 0, false
}

// Out is what a key produces in one modifier state.
type Out struct {
	// Text is the character typed. It is occasionally more than one rune:
	// a few layouts have ligature keys, and dead-key accents are given as
	// the spacing form of the accent.
	Text string `json:"t"`

	// Dead reports that the key does not type Text directly but arms an
	// accent that combines with the next keystroke. The circumflex on a
	// French keyboard is a dead key: it types nothing until "o" follows,
	// at which point "ô" appears.
	Dead bool `json:"d,omitempty"`
}

// Key is one physical switch on the board, together with everything it types.
type Key struct {
	// SC is the set 1 scan code, in uppercase hex. It identifies the switch
	// by position and means the same thing on every layout.
	SC string `json:"sc"`

	// VK is the Windows virtual key the switch maps to, such as "VK_Q".
	VK string `json:"vk,omitempty"`

	// Name is the label printed on the key, when the layout gives one.
	Name string `json:"name,omitempty"`

	// Pos is where the key sits, filled in from the layout's form when the
	// layout loads rather than stored in the dataset.
	Pos Pos `json:"-"`

	out map[Mod]Out
}

// NewKey builds a key on a scan code. Use Set to give it outputs. It is what
// the dataset generator calls, and what a caller building a layout of their
// own would use.
func NewKey(sc, vk, name string) Key {
	return Key{SC: strings.ToUpper(sc), VK: vk, Name: name, out: map[Mod]Out{}}
}

// Set records what the key types in a modifier state. Setting a state twice
// keeps the first value, since the dataset lists the plainest spelling first.
//
// A Key holds its outputs in a map, and copying a Key copies the reference to
// it, so writing through one copy would otherwise be visible through all of
// them — including through the layout the copy came from, which Get shares
// between every caller in the program. Set therefore replaces the map rather
// than writing into it, which makes a Key behave like the value it looks like.
func (k *Key) Set(m Mod, o Out) {
	if _, exists := k.out[m]; exists {
		return
	}

	out := make(map[Mod]Out, len(k.out)+1)
	for km, ko := range k.out {
		out[km] = ko
	}
	out[m] = o

	k.out = out
}

// Text returns what the key types in the given modifier state, and whether it
// types anything at all.
//
// A Caps Lock state the key does not list falls back to capsFallback, the
// ordinary behaviour of Caps Lock on a letter. Only the keys that depart from
// it — the digit row of a German keyboard, where Caps Lock does nothing at all
// — are stored.
func (k Key) Text(m Mod) (Out, bool) {
	if o, ok := k.out[m]; ok {
		return o, true
	}
	if m&Caps != 0 {
		if o, ok := k.out[capsFallback(m)]; ok {
			return o, true
		}
	}
	return Out{}, false
}

// capsFallback maps a Caps Lock state onto the unlocked state that behaves the
// same way. Caps Lock inverts Shift, so Caps alone reads as Shift, and Shift
// with Caps reads as neither.
func capsFallback(m Mod) Mod { return (m &^ Caps) ^ Shift }

// Compact drops the modifier states that Text would reproduce on its own,
// which is most of them: a letter key carries four states in the source data
// and needs to store two. It is called when the dataset is generated.
//
// Like Set, it replaces the map rather than writing into it, so compacting one
// copy of a Key leaves every other copy alone.
func (k *Key) Compact() {
	out := make(map[Mod]Out, len(k.out))
	for m, o := range k.out {
		if m&Caps != 0 {
			if fb, ok := k.out[capsFallback(m)]; ok && fb == o {
				continue
			}
		}
		out[m] = o
	}

	k.out = out
}

// Base returns what the key types unmodified, or "" if it types nothing.
func (k Key) Base() string { o, _ := k.Text(Base); return o.Text }

// Shift returns what the key types with Shift held, or "" if it types nothing.
func (k Key) Shift() string { o, _ := k.Text(Shift); return o.Text }

// AltGr returns what the key types with AltGr held, or "" if it types nothing.
func (k Key) AltGr() string { o, _ := k.Text(AltGr); return o.Text }

// Mods returns every modifier state this key produces text in, in a stable
// order.
func (k Key) Mods() []Mod {
	out := make([]Mod, 0, len(k.out))
	for _, m := range mods {
		if _, ok := k.out[m]; ok {
			out = append(out, m)
		}
	}
	return out
}

// Texts returns every distinct string the key can type, in modifier order.
// Dead keys contribute the spacing form of their accent.
func (k Key) Texts() []string {
	seen := make(map[string]bool, len(k.out))
	out := make([]string, 0, len(k.out))
	for _, m := range k.Mods() {
		t := k.out[m].Text
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

// Types reports whether the key produces s in any modifier state.
func (k Key) Types(s string) bool {
	for _, m := range k.Mods() {
		if k.out[m].Text == s {
			return true
		}
	}
	return false
}

// Blank reports whether the key types nothing at all — a modifier, Enter, or
// a switch the layout leaves unassigned.
func (k Key) Blank() bool {
	for _, o := range k.out {
		if o.Text != "" {
			return false
		}
	}
	return true
}

// String renders the key as its scan code and outputs, for diagnostics.
func (k Key) String() string {
	parts := make([]string, 0, len(k.out))
	for _, m := range k.Mods() {
		o := k.out[m]
		s := fmt.Sprintf("%s=%q", m, o.Text)
		if o.Dead {
			s += "(dead)"
		}
		parts = append(parts, s)
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%s[%s]", k.SC, k.Name)
	}
	return fmt.Sprintf("%s %s", k.SC, strings.Join(parts, " "))
}

// keyWire is the JSON shape of a key: no position, since that is recomputed
// from the scan code, and modifier states keyed by name so the output reads
// well. The dataset itself is protocol buffers — see proto.go — and JSON is
// here for callers exporting a layout, and for tests that want a layout
// spelled out in full.
type keyWire struct {
	SC   string         `json:"sc"`
	VK   string         `json:"vk,omitempty"`
	Name string         `json:"name,omitempty"`
	Out  map[string]Out `json:"out,omitempty"`
}

// MarshalJSON writes the key as JSON.
func (k Key) MarshalJSON() ([]byte, error) {
	w := keyWire{SC: k.SC, VK: k.VK, Name: k.Name}
	if len(k.out) > 0 {
		w.Out = make(map[string]Out, len(k.out))
		for m, o := range k.out {
			w.Out[m.String()] = o
		}
	}
	return json.Marshal(w)
}

// UnmarshalJSON reads a key from JSON. Position is left unset; the layout
// fills it in when it indexes its keys.
func (k *Key) UnmarshalJSON(b []byte) error {
	var w keyWire
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}

	*k = Key{SC: strings.ToUpper(w.SC), VK: w.VK, Name: w.Name, out: make(map[Mod]Out, len(w.Out))}
	for name, o := range w.Out {
		m, ok := parseMod(name)
		if !ok {
			return fmt.Errorf("kb: key %s: unknown modifier state %q", w.SC, name)
		}
		k.out[m] = o
	}

	return nil
}

// byPosition orders keys the way they are read off a keyboard: top row first,
// then left to right.
func byPosition(keys []Key) {
	sort.SliceStable(keys, func(i, j int) bool {
		if keys[i].Pos.Y != keys[j].Pos.Y {
			return keys[i].Pos.Y < keys[j].Pos.Y
		}
		return keys[i].Pos.X < keys[j].Pos.X
	})
}
