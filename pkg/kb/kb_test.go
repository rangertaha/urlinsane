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
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/rangertaha/urlinsane/pkg/kb/internal/kbpb"
)

// TestAdjacentOnQWERTY pins the neighbours of a few keys on the US layout,
// which anyone with a keyboard in front of them can check.
func TestAdjacentOnQWERTY(t *testing.T) {
	us := MustGet("kbdus")

	for _, tc := range []struct {
		in   string
		want []string
	}{
		// The order is by distance: the keys straight left and right first,
		// then the staggered rows above and below.
		{"e", []string{"w", "r", "d", "3", "4", "s"}},
		{"s", []string{"a", "d", "w", "z", "x", "e"}},
		{"a", []string{"s", "q", "z", "w"}},
		{"1", []string{"`", "2", "q"}},
		// The space bar counts: it runs the length of the bottom row, so
		// it sits under "m" as much as under "b".
		{"m", []string{"n", ",", " ", "j", "k"}},
		// Shift travels with the character.
		{"E", []string{"W", "R", "D", "#", "$", "S"}},
	} {
		if got := us.Adjacent(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Adjacent(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

// TestStaggerOutranksColumns is the reason this package positions keys rather
// than tabulating them. The rows of a keyboard are offset from one another, so
// the key below "e" is not the one in the same column: "d" overlaps it far
// more than "s" does, and a typist's finger is correspondingly likelier to
// land there. A layout stored as equal-width rows cannot tell the two apart.
func TestStaggerOutranksColumns(t *testing.T) {
	us := MustGet("kbdus")

	e, ok := us.Key("12") // the E key
	if !ok {
		t.Fatal("no key at SC 12")
	}
	d, _ := us.Key("20")
	s, _ := us.Key("1F")

	toD := e.Pos.Distance(d.Pos)
	toS := e.Pos.Distance(s.Pos)
	if toD >= toS {
		t.Errorf("e->d is %.3f and e->s is %.3f; d should be the nearer", toD, toS)
	}

	got := us.Adjacent("e")
	if indexOf(got, "d") > indexOf(got, "s") {
		t.Errorf("Adjacent(e) = %v; d should rank before s", got)
	}
}

// TestAdjacentFollowsTheLayout checks that adjacency is read off the layout
// rather than assumed from QWERTY. Each of these is the same physical key in a
// different arrangement.
func TestAdjacentFollowsTheLayout(t *testing.T) {
	for _, tc := range []struct {
		layout string
		in     string
		want   []string
	}{
		// QWERTZ swaps Y and Z, so German "z" sits where English "y" does.
		{"kbdgr", "z", []string{"t", "u", "h", "6", "7", "g"}},
		// AZERTY puts "a" where QWERTY has "q".
		{"kbdfr", "a", []string{"z", "q", "&", "é"}},
		// Dvorak moves "e" to the home row.
		{"kbddv", "e", []string{"o", "u", ".", "q", "j", "p"}},
	} {
		l := MustGet(tc.layout)
		if got := l.Adjacent(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s.Adjacent(%q) = %v, want %v", tc.layout, tc.in, got, tc.want)
		}
	}
}

// TestFormChangesNeighbours covers the other thing a flat table cannot express:
// an ISO keyboard has a key between the left Shift and "z" that an ANSI
// keyboard does not, and it is a neighbour of "a".
func TestFormChangesNeighbours(t *testing.T) {
	uk := MustGet("kbduk")
	if uk.Form != ISO {
		t.Fatalf("kbduk form = %q, want iso", uk.Form)
	}

	iso := uk.Adjacent("a")
	if indexOf(iso, `\`) < 0 {
		t.Errorf("on ISO, Adjacent(a) = %v; want the extra key below left Shift", iso)
	}

	ansi := uk.With(ANSI).Adjacent("a")
	if indexOf(ansi, `\`) >= 0 {
		t.Errorf("on ANSI, Adjacent(a) = %v; the extra key does not exist", ansi)
	}

	if len(uk.With(ANSI).Keys) >= len(uk.Keys) {
		t.Errorf("ANSI has %d keys and ISO %d; ANSI should have fewer",
			len(uk.With(ANSI).Keys), len(uk.Keys))
	}

	// The original is untouched.
	if uk.Form != ISO {
		t.Error("With mutated the layout it was called on")
	}
}

// TestAdjacencyIsSymmetric checks the relation both ways across every layout:
// if b is next to a then a is next to b. Distance is symmetric, so anything
// else would mean a key is missing from an index.
func TestAdjacencyIsSymmetric(t *testing.T) {
	for _, e := range List() {
		l := MustGet(e.ID)

		for _, k := range l.Keys {
			if k.Blank() {
				continue
			}
			for _, n := range l.AdjacentKeys(k, DefaultRadius) {
				back := l.AdjacentKeys(n, DefaultRadius)
				found := false
				for _, b := range back {
					if b.SC == k.SC {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("%s: %s neighbours %s but not the reverse", e.ID, k.SC, n.SC)
				}
			}
		}
	}
}

func TestShiftLevels(t *testing.T) {
	us := MustGet("kbdus")

	if got := us.Shifted("4"); !reflect.DeepEqual(got, []string{"$"}) {
		t.Errorf("Shifted(4) = %v, want [$]", got)
	}
	if got := us.Unshifted("$"); !reflect.DeepEqual(got, []string{"4"}) {
		t.Errorf("Unshifted($) = %v, want [4]", got)
	}
	if got := us.Shifted("e"); !reflect.DeepEqual(got, []string{"E"}) {
		t.Errorf("Shifted(e) = %v, want [E]", got)
	}
	// A character already at the shift level has nothing above it.
	if got := us.Shifted("$"); len(got) != 0 {
		t.Errorf("Shifted($) = %v, want nothing", got)
	}

	// AltGr is where the extra characters of a European layout live: on the
	// French keyboard the key that types "à" gives "@".
	fr := MustGet("kbdfr")
	if got := fr.AltGraphed("à"); !reflect.DeepEqual(got, []string{"@"}) {
		t.Errorf("fr.AltGraphed(à) = %v, want [@]", got)
	}
	if got := us.AltGraphed("e"); len(got) != 0 {
		t.Errorf("us.AltGraphed(e) = %v; the US layout has no AltGr level", got)
	}
}

// TestCapsFallback checks the rule the dataset relies on to avoid storing a
// Caps Lock state for every letter on every layout.
func TestCapsFallback(t *testing.T) {
	us := MustGet("kbdus")

	e, ok := us.Key("12")
	if !ok {
		t.Fatal("no key at SC 12")
	}

	// Only the two states that carry information are stored.
	if got := e.Mods(); !reflect.DeepEqual(got, []Mod{Base, Shift}) {
		t.Errorf("stored states = %v, want [base shift]", got)
	}

	// Caps Lock inverts Shift, and the rest is derived.
	for _, tc := range []struct {
		mod  Mod
		want string
	}{
		{Base, "e"},
		{Shift, "E"},
		{Caps, "E"},
		{Shift | Caps, "e"},
	} {
		o, ok := e.Text(tc.mod)
		if !ok || o.Text != tc.want {
			t.Errorf("Text(%s) = %q (%v), want %q", tc.mod, o.Text, ok, tc.want)
		}
	}

	// Layouts that depart from the rule keep their own value: Caps Lock on
	// the German number row types the digit, not the symbol above it.
	de := MustGet("kbdgr")
	two, _ := de.Key("03")
	if o, _ := two.Text(Caps); o.Text != `"` {
		t.Errorf("de SC03 with caps = %q, want a stored value rather than the fallback", o.Text)
	}
}

func TestDeadKeys(t *testing.T) {
	de := MustGet("kbdgr")

	k, ok := de.Key("29")
	if !ok {
		t.Fatal("no key at SC 29")
	}

	o, ok := k.Text(Base)
	if !ok {
		t.Fatal("SC 29 types nothing")
	}
	if o.Text != "^" || !o.Dead {
		t.Errorf("de SC29 = %q dead=%v, want the circumflex marked dead", o.Text, o.Dead)
	}
}

func TestLookupAcceptsEveryName(t *testing.T) {
	want := "kbdus"

	for _, id := range []string{
		"kbdus", "KBDUS", " kbdus ", "KBDUS.DLL",
		"00000409", "0x00000409", "409", // the US KLID, spelled as Windows might
		"00000804", // a Chinese locale that shares the driver
	} {
		l, err := Get(id)
		if err != nil {
			t.Errorf("Get(%q): %v", id, err)
			continue
		}
		if l.ID != want {
			t.Errorf("Get(%q) = %q, want %q", id, l.ID, want)
		}
	}

	if _, err := Get("nosuchlayout"); err == nil {
		t.Error("Get on an unknown layout should fail")
	}
}

func TestGetIsCached(t *testing.T) {
	a := MustGet("kbdus")
	b := MustGet("00000409")
	if a != b {
		t.Error("Get returned two different pointers for the same layout")
	}
}

func TestByLanguage(t *testing.T) {
	de := ByLanguage("de")
	if len(de) < 2 {
		t.Fatalf("ByLanguage(de) returned %d layouts, want several", len(de))
	}

	var ids []string
	for _, e := range de {
		ids = append(ids, e.ID)
	}
	if indexOf(ids, "kbdgr") < 0 || indexOf(ids, "kbdsg") < 0 {
		t.Errorf("ByLanguage(de) = %v, want the German and Swiss German layouts", ids)
	}

	// A full tag is narrower than the bare subtag.
	swiss := ByLanguage("de-CH")
	if len(swiss) == 0 || len(swiss) >= len(de) {
		t.Errorf("ByLanguage(de-CH) returned %d of %d German layouts", len(swiss), len(de))
	}

	if got := ByLanguage("zz"); len(got) != 0 {
		t.Errorf("ByLanguage(zz) = %v, want nothing", got)
	}
}

func TestFind(t *testing.T) {
	if got := Find("dvorak"); len(got) == 0 {
		t.Error("Find(dvorak) found nothing")
	}
	if got := Find("french"); len(got) == 0 {
		t.Error("Find(french) found nothing")
	}
	if got := Find(""); got != nil {
		t.Errorf("Find(\"\") = %v, want nothing", got)
	}
}

// TestDatasetIsWellFormed walks every layout in the catalogue. It is the check
// that the generated data and the code that reads it still agree.
func TestDatasetIsWellFormed(t *testing.T) {
	entries := List()
	if len(entries) < 100 {
		t.Fatalf("catalogue has %d layouts, want the full set", len(entries))
	}

	onDisk, err := files()
	if err != nil {
		t.Fatal(err)
	}
	if len(onDisk) != len(entries) {
		t.Errorf("%d layout files but %d catalogue entries", len(onDisk), len(entries))
	}

	seen := map[string]bool{}

	for _, e := range entries {
		if seen[e.ID] {
			t.Errorf("duplicate layout %q", e.ID)
		}
		seen[e.ID] = true

		l, err := Get(e.ID)
		if err != nil {
			t.Errorf("Get(%q): %v", e.ID, err)
			continue
		}

		switch {
		case l.ID != e.ID:
			t.Errorf("%s: file declares id %q", e.ID, l.ID)
		case l.Name == "":
			t.Errorf("%s: no name", e.ID)
		case l.Form != ANSI && l.Form != ISO && l.Form != JIS:
			t.Errorf("%s: form %q", e.ID, l.Form)
		case len(l.Keys) == 0:
			t.Errorf("%s: no keys", e.ID)
		case len(l.Locales) == 0:
			t.Errorf("%s: no locales", e.ID)
		}

		// Every key must be positioned and reachable both ways.
		for _, k := range l.Keys {
			if k.Pos.W == 0 {
				t.Errorf("%s: key %s has no position", e.ID, k.SC)
			}
			if _, ok := l.Key(k.SC); !ok {
				t.Errorf("%s: key %s is not indexed", e.ID, k.SC)
			}
			for _, text := range k.Texts() {
				if !l.Types(text) {
					t.Errorf("%s: key %s types %q but the layout says otherwise", e.ID, k.SC, text)
				}
			}
		}

		// A layout that cannot type its own alphabet is a parsing failure.
		if l.Types("a") && len(l.Adjacent("a")) == 0 {
			t.Errorf("%s: 'a' has no neighbours", e.ID)
		}
	}
}

// TestEveryLayoutTypesSomething guards against a layout that parsed into an
// empty shell, which is how a change in the source data would first show up.
//
// The bar is low because some layouts genuinely are: the Ogham and Old Italic
// keyboards leave most of the board unassigned, having only a couple of dozen
// letters to place.
func TestEveryLayoutTypesSomething(t *testing.T) {
	for _, e := range List() {
		l := MustGet(e.ID)

		typed := 0
		for _, k := range l.Keys {
			if !k.Blank() {
				typed++
			}
		}
		if typed < 20 {
			t.Errorf("%s: only %d keys type anything", e.ID, typed)
		}
	}
}

// TestProtoRoundTrip exercises the format the dataset is actually stored in.
// The files are opaque to a diff, so this is what stands in for reading them.
func TestProtoRoundTrip(t *testing.T) {
	for _, e := range List() {
		original := MustGet(e.ID)

		raw, err := original.Marshal()
		if err != nil {
			t.Fatalf("%s: %v", e.ID, err)
		}

		reloaded := new(Layout)
		if err := reloaded.Unmarshal(raw); err != nil {
			t.Fatalf("%s: %v", e.ID, err)
		}

		// Compare through JSON, which spells out everything a layout
		// carries, rather than comparing the bytes to themselves.
		before, err := json.Marshal(original)
		if err != nil {
			t.Fatal(err)
		}
		after, err := json.Marshal(reloaded)
		if err != nil {
			t.Fatal(err)
		}
		if string(before) != string(after) {
			t.Errorf("%s did not survive encoding", e.ID)
		}

		// Encoding is deterministic, so a rebuild that changed nothing
		// leaves the dataset byte for byte identical.
		again, err := reloaded.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(raw, again) {
			t.Errorf("%s does not encode deterministically", e.ID)
		}
	}
}

func TestCatalogueRoundTrip(t *testing.T) {
	entries := List()

	raw, err := MarshalCatalogue(entries)
	if err != nil {
		t.Fatal(err)
	}

	reloaded, err := UnmarshalCatalogue(raw)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(entries, reloaded) {
		t.Errorf("the catalogue did not survive encoding")
	}
}

// TestRejectsUnknownModifier checks both ends of the encoding. A state this
// package cannot name must not be quietly dropped on the way out, and a
// dataset written by a newer version carrying one must fail loudly on the way
// in rather than losing keys.
func TestRejectsUnknownModifier(t *testing.T) {
	k := NewKey("10", "VK_Q", "")
	k.Set(Base, Out{Text: "q"})
	k.Set(Mod(0xFF), Out{Text: "?"})

	if _, err := New("test", "Test", ANSI, []Key{k}).Marshal(); err == nil {
		t.Error("Marshal dropped an unknown modifier state instead of failing")
	}

	// Coming the other way, the state has to be planted in the encoded form
	// directly, since Marshal now refuses to produce one.
	raw, err := proto.Marshal(&kbpb.Layout{
		Id:   "test",
		Name: "Test",
		Form: kbpb.Form_FORM_ANSI,
		Keys: []*kbpb.Key{{
			Sc:      0x10,
			Vk:      kbpb.VirtualKey_VK_Q,
			Outputs: []*kbpb.Output{{Mod: 0xFF, Text: "?"}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := new(Layout).Unmarshal(raw); err == nil {
		t.Error("Unmarshal accepted an unknown modifier state")
	}
}

// TestVirtualKeyEnum checks the mapping the dataset relies on to store virtual
// keys as numbers while the API keeps their names.
func TestVirtualKeyEnum(t *testing.T) {
	// The enum values are the codes Windows uses, not positions in a list.
	for name, want := range map[string]int32{
		"VK_A":       0x41,
		"VK_Z":       0x5A,
		"VK_0":       0x30,
		"VK_SPACE":   0x20,
		"VK_OEM_1":   0xBA,
		"VK_OEM_102": 0xE2,
	} {
		if got := kbpb.VirtualKey_value[name]; got != want {
			t.Errorf("%s = %#x, want %#x", name, got, want)
		}
	}

	// Names survive the round trip, and an unnamed key stays unnamed.
	for _, name := range []string{"", "VK_Q", "VK_OEM_MINUS", "VK_LSHIFT"} {
		enc, err := encodeVK(name)
		if err != nil {
			t.Fatalf("encodeVK(%q): %v", name, err)
		}
		got, err := decodeVK(enc)
		if err != nil {
			t.Fatalf("decodeVK(%q): %v", name, err)
		}
		if got != name {
			t.Errorf("%q came back as %q", name, got)
		}
	}

	// A virtual key the schema has never seen must stop the build.
	if _, err := encodeVK("VK_MADE_UP"); err == nil {
		t.Error("encodeVK accepted an unknown virtual key")
	}
	if _, err := decodeVK(kbpb.VirtualKey(0x7777)); err == nil {
		t.Error("decodeVK accepted an unknown virtual key")
	}
}

// TestScanCodeEncoding checks that the hex spelling the API uses survives
// being stored as a number.
func TestScanCodeEncoding(t *testing.T) {
	for _, sc := range []string{"02", "0A", "10", "1E", "2B", "39", "3A", "56", "7D"} {
		enc, err := encodeSC(sc)
		if err != nil {
			t.Fatalf("encodeSC(%q): %v", sc, err)
		}
		if got := decodeSC(enc); got != sc {
			t.Errorf("%q came back as %q", sc, got)
		}
	}

	if _, err := encodeSC("zz"); err == nil {
		t.Error("encodeSC accepted a non-hex scan code")
	}
}

// TestRejectsUnknownForm checks the other enum the dataset carries.
func TestRejectsUnknownForm(t *testing.T) {
	raw, err := proto.Marshal(&kbpb.Layout{Id: "test", Form: kbpb.Form_FORM_UNSPECIFIED})
	if err != nil {
		t.Fatal(err)
	}
	if err := new(Layout).Unmarshal(raw); err == nil {
		t.Error("Unmarshal accepted a layout with no form")
	}

	if _, err := (&Layout{ID: "test", Form: Form("bogus")}).Marshal(); err == nil {
		t.Error("Marshal accepted a layout with an unknown form")
	}
}

func TestRoundTrip(t *testing.T) {
	for _, id := range []string{"kbdus", "kbdfr", "kbdgr", "kbdjpn"} {
		original := MustGet(id)

		raw, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}

		reloaded := new(Layout)
		if err := json.Unmarshal(raw, reloaded); err != nil {
			t.Fatalf("%s: %v", id, err)
		}

		again, err := json.Marshal(reloaded)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if string(raw) != string(again) {
			t.Errorf("%s did not survive a round trip", id)
		}
		if got, want := reloaded.Adjacent("e"), original.Adjacent("e"); !reflect.DeepEqual(got, want) {
			t.Errorf("%s: Adjacent(e) = %v after reloading, was %v", id, got, want)
		}
	}
}

func TestRows(t *testing.T) {
	us := MustGet("kbdus")

	rows := us.Rows()
	if len(rows) != 5 {
		t.Fatalf("got %d rows, want 5", len(rows))
	}

	// Rows run top to bottom, and each runs left to right.
	for i, row := range rows {
		for j := 1; j < len(row); j++ {
			if row[j].Pos.X < row[j-1].Pos.X {
				t.Errorf("row %d is out of order at %d", i, j)
			}
		}
		if i > 0 && row[0].Pos.Y <= rows[i-1][0].Pos.Y {
			t.Errorf("row %d is not below row %d", i, i-1)
		}
	}

	var home []string
	for _, k := range rows[2] {
		if b := k.Base(); b != "" {
			home = append(home, b)
		}
	}
	if got := strings.Join(home, ""); got != "asdfghjkl;'" {
		t.Errorf("home row = %q, want asdfghjkl;'", got)
	}
}

func TestLanguages(t *testing.T) {
	langs := Languages()
	if len(langs) < 50 {
		t.Fatalf("%d languages, want the full set", len(langs))
	}
	if !sort.StringsAreSorted(langs) {
		t.Error("Languages is not sorted")
	}

	us := MustGet("kbdus")
	if got := us.Languages(); indexOf(got, "en") < 0 {
		t.Errorf("kbdus languages = %v, want en among them", got)
	}
}

func TestRadiusWidensTheNeighbourhood(t *testing.T) {
	us := MustGet("kbdus")

	near := us.Adjacent("e")
	far := us.AdjacentWithin("e", 2.0)

	if len(far) <= len(near) {
		t.Errorf("a wider radius gave %d neighbours against %d", len(far), len(near))
	}
	// Widening only adds; the near neighbours stay, and stay in front.
	for i, c := range near {
		if far[i] != c {
			t.Errorf("far[%d] = %q, want %q", i, far[i], c)
		}
	}
}

// A Layout from Get is shared by every caller in the program, so nothing
// handed out of this package may be a window onto it. These are the ways that
// went wrong once.

// TestKeyIsAValue checks that a Key copied out of a layout is independent of
// it. Keys hold their outputs in a map, so a copy shares that map unless the
// mutators take care not to write through it.
func TestKeyIsAValue(t *testing.T) {
	board := MustGet("kbdus")

	k, ok := board.Key("11") // the W key
	if !ok {
		t.Fatal("no key at SC 11")
	}

	k.Set(AltGr, Out{Text: "ZZZ"})
	if o, ok := k.Text(AltGr); !ok || o.Text != "ZZZ" {
		t.Error("Set did not take effect on the copy")
	}

	fresh, _ := board.Key("11")
	if _, leaked := fresh.Text(AltGr); leaked {
		t.Error("writing to a copied Key leaked into the shared layout")
	}

	// Compact is the other mutator, and must not write through either.
	e, _ := board.Key("12")
	e.Compact()
	if again, _ := board.Key("12"); len(again.Mods()) != len(e.Mods()) {
		t.Error("Compact on a copied Key leaked into the shared layout")
	}
}

// TestNewLeavesCallerSliceAlone checks that building a layout does not reuse
// the caller's backing array. Layouts drop the keys that fall outside the
// alphanumeric block, and filtering in place would shuffle the caller's slice
// and leave a duplicated tail.
func TestNewLeavesCallerSliceAlone(t *testing.T) {
	keys := []Key{
		NewKey("10", "VK_Q", ""),
		NewKey("E0", "VK_X", ""), // outside the block, so the layout drops it
		NewKey("11", "VK_W", ""),
	}
	want := []string{"10", "E0", "11"}

	l := New("test", "Test", ANSI, keys)

	for i, sc := range want {
		if keys[i].SC != sc {
			t.Errorf("caller's keys[%d] = %q after New, want %q", i, keys[i].SC, sc)
		}
	}
	if len(l.Keys) != 2 {
		t.Errorf("layout kept %d keys, want the 2 inside the block", len(l.Keys))
	}
}

// TestCatalogueIsNotShared checks that the entries List and friends hand out
// cannot be written back into the catalogue.
func TestCatalogueIsNotShared(t *testing.T) {
	for name, get := range map[string]func() []Entry{
		"List":       List,
		"ByLanguage": func() []Entry { return ByLanguage("de") },
		"Find":       func() []Entry { return Find("german") },
		"ByKeys":     func() []Entry { return ByKeys("a", "b") },
	} {
		got := get()
		if len(got) == 0 || len(got[0].Locales) == 0 {
			t.Fatalf("%s returned nothing to test with", name)
		}

		was := got[0].Locales[0].Tag
		got[0].Locales[0].Tag = "corrupted"

		if again := get(); again[0].Locales[0].Tag != was {
			t.Errorf("%s hands out the catalogue's own locales", name)
		}

		got[0].Locales[0].Tag = was
	}
}

// TestWithDoesNotShareLocales checks the same for a rearranged layout.
func TestWithDoesNotShareLocales(t *testing.T) {
	uk := MustGet("kbduk")
	if len(uk.Locales) == 0 {
		t.Fatal("kbduk has no locales")
	}

	ansi := uk.With(ANSI)
	if &uk.Locales[0] == &ansi.Locales[0] {
		t.Error("With shares the locale slice with the layout it copied")
	}

	// Asking for the shape it already has is documented to return the
	// receiver rather than a copy.
	if uk.With(ISO) != uk {
		t.Error("With(same form) should return the receiver")
	}
}

// TestPositionedDoesNotAllocate pins the geometry tables being built once.
// Positioned runs for every key of every layout that loads, and it used to
// rebuild a fifty-entry map on each call.
func TestPositionedDoesNotAllocate(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { Positioned(ANSI, "10") }); n != 0 {
		t.Errorf("Positioned allocates %.0f times per call, want 0", n)
	}

	// An unrecognised form falls back to ANSI rather than to no keys at all.
	if _, ok := Positioned(Form("bogus"), "10"); !ok {
		t.Error("an unknown form should be read as ANSI")
	}
}

// TestLanguagesFoldCase checks that Languages and ByLanguage agree on
// spelling. BCP-47 tags are case-insensitive and the dataset carries them as
// Windows writes them, so a tag that arrived capitalised must still fold into
// the language its neighbours are filed under.
func TestLanguagesFoldCase(t *testing.T) {
	l := Layout{Locales: []Locale{
		{KLID: "1", Tag: "EN"},
		{KLID: "2", Tag: "en-GB"},
		{KLID: "3", Tag: "De-DE"},
	}}

	if got := l.Languages(); !reflect.DeepEqual(got, []string{"de", "en"}) {
		t.Errorf("Languages() = %v, want [de en]", got)
	}

	// The live catalogue has to hold to the same rule, or ByLanguage would
	// miss whole layouts.
	for _, e := range List() {
		for _, lang := range e.Languages() {
			if lang != strings.ToLower(lang) {
				t.Errorf("%s reports language %q, which ByLanguage could never match", e.ID, lang)
			}
			if len(ByLanguage(lang)) == 0 {
				t.Errorf("%s reports language %q, which ByLanguage finds nothing for", e.ID, lang)
			}
		}
	}
}

// TestPositionedAcceptsEitherCase covers the exported lookup: hex has two
// spellings and a caller has no reason to guess which one is wanted.
func TestPositionedAcceptsEitherCase(t *testing.T) {
	want, ok := Positioned(ANSI, "1E")
	if !ok {
		t.Fatal("SC 1E is not positioned")
	}
	for _, sc := range []string{"1e", " 1E ", "1E"} {
		got, ok := Positioned(ANSI, sc)
		if !ok || got != want {
			t.Errorf("Positioned(%q) = %v (%v), want %v", sc, got, ok, want)
		}
	}
}

// TestGeometryHasNoOverlaps checks that no two keys are given the same spot on
// any form. Two keys at one position would be zero apart, and each would count
// as the other's nearest neighbour.
func TestGeometryHasNoOverlaps(t *testing.T) {
	for _, f := range []Form{ANSI, ISO, JIS} {
		seen := map[Pos]string{}
		for sc, p := range geometries[f] {
			if prev, ok := seen[p]; ok {
				t.Errorf("%s: SC %s and SC %s are both at %v", f, prev, sc, p)
			}
			seen[p] = sc
		}
	}
}

func TestByKeys(t *testing.T) {
	// Every layout returned really does type everything asked for.
	got := ByKeys("a", "b", "c")
	if len(got) == 0 {
		t.Fatal("ByKeys(a,b,c) found nothing")
	}
	for _, e := range got {
		l := MustGet(e.ID)
		for _, s := range []string{"a", "b", "c"} {
			if !l.Types(s) {
				t.Errorf("%s was returned but cannot type %q", e.ID, s)
			}
		}
	}

	// And no layout that types them all was left out.
	want := 0
	for _, e := range List() {
		l := MustGet(e.ID)
		if l.Types("a") && l.Types("b") && l.Types("c") {
			want++
		}
	}
	if len(got) != want {
		t.Errorf("ByKeys returned %d layouts, want %d", len(got), want)
	}

	// Asking for more can only narrow the field.
	if a, ab := len(ByKeys("a")), len(ByKeys("a", "b")); ab > a {
		t.Errorf("ByKeys(a,b) matched %d layouts but ByKeys(a) only %d", ab, a)
	}

	// The result is ordered like the catalogue.
	if !sort.SliceIsSorted(got, func(i, j int) bool { return got[i].ID < got[j].ID }) {
		t.Error("ByKeys did not return catalogue order")
	}

	// Case matters, because "a" and "A" are different keystrokes.
	if reflect.DeepEqual(ByKeys("a"), ByKeys("A")) {
		t.Error("ByKeys should treat a and A as different texts")
	}

	// Nothing asked for, nothing found; likewise for an untypeable text.
	if n := ByKeys(); n != nil {
		t.Errorf("ByKeys() = %v, want nothing", n)
	}
	if n := ByKeys("a", "\u4e2d"); len(n) != 0 {
		t.Errorf("ByKeys with an untypeable text returned %d layouts", len(n))
	}
}

func TestByString(t *testing.T) {
	// It is ByKeys over the runes, so the two must agree.
	if got, want := ByString("abc"), ByKeys("a", "b", "c"); !reflect.DeepEqual(got, want) {
		t.Errorf("ByString(abc) = %d layouts, ByKeys(a,b,c) = %d", len(got), len(want))
	}

	// Repeated characters are asked about once and change nothing.
	if got, want := len(ByString("hello")), len(ByKeys("h", "e", "l", "o")); got != want {
		t.Errorf("ByString(hello) = %d, want %d", got, want)
	}

	for _, e := range ByString("münchen") {
		l := MustGet(e.ID)
		for _, r := range "münchen" {
			if !l.Types(string(r)) {
				t.Errorf("%s was returned but cannot type %q", e.ID, r)
			}
		}
	}

	if n := ByString(""); n != nil {
		t.Errorf("ByString(\"\") = %v, want nothing", n)
	}

	// A dead key is not a key that types the composed character. German has
	// an acute accent that makes "é" when pressed before "e", but no key of
	// its own for it, and ByString reports what is on the keys.
	de := MustGet("kbdgr")
	if de.Types("é") {
		t.Fatal("kbdgr appears to type é directly; the caveat below is stale")
	}
	for _, e := range ByString("café") {
		if e.ID == "kbdgr" {
			t.Error("ByString counted a dead-key composition as a key")
		}
	}
}

func TestPrint(t *testing.T) {
	out := MustGet("kbdus").String()

	if !strings.Contains(out, "kbdus") || !strings.Contains(out, "ansi") {
		t.Error("the drawing is not labelled with the layout it shows")
	}

	// Every row of the alphanumeric block should be there, base and shift.
	for _, want := range []string{
		"| ` | 1 | 2 |", "| ~ | ! | @ |",
		"| q | w | e |", "| Q | W | E |",
		"| a | s | d |",
		"| z | x | c |",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("drawing is missing %q:\n%s", want, out)
		}
	}

	// Nothing may run past the right-hand edge of the widest key, and no
	// line may carry trailing blanks.
	for i, line := range strings.Split(out, "\n") {
		if line != strings.TrimRight(line, " ") {
			t.Errorf("line %d has trailing spaces", i)
		}
	}

	// The stagger has to survive into the drawing: "a" sits to the right of
	// where "q" starts, which is why the rows are indented differently.
	lines := strings.Split(out, "\n")
	var qLine, aLine string
	for _, line := range lines {
		if strings.Contains(line, "| q |") {
			qLine = line
		}
		if strings.Contains(line, "| a |") {
			aLine = line
		}
	}
	if strings.Index(qLine, "| q |") >= strings.Index(aLine, "| a |") {
		t.Error("the row stagger is missing from the drawing")
	}
}

func TestPrintMarksDeadKeys(t *testing.T) {
	// A dead key types nothing on its own, so it is bracketed rather than
	// shown as though it were an ordinary character.
	out := MustGet("kbdgr").String()
	if !strings.Contains(out, "[^]") {
		t.Errorf("German circumflex is a dead key and should be bracketed:\n%s", out)
	}
}

func TestPrintWritesToAWriter(t *testing.T) {
	// Print goes to stdout; Fprint is how it is captured, and both must
	// produce exactly what String renders.
	var b bytes.Buffer
	l := MustGet("kbdus")
	if err := l.Fprint(&b); err != nil {
		t.Fatal(err)
	}
	if b.String() != l.String() {
		t.Error("Fprint and String disagree")
	}

	// fmt sees the same drawing, since String is the Stringer.
	if fmt.Sprint(l) != l.String() {
		t.Error("fmt.Sprint and String disagree")
	}
}

func TestPrintHandlesEmpty(t *testing.T) {
	// A layout with no keys still names itself rather than panicking.
	out := New("empty", "Empty", ANSI, nil).String()
	if !strings.Contains(out, "empty") {
		t.Errorf("empty layout rendered as %q", out)
	}
}

func TestPrintEveryLayout(t *testing.T) {
	// The drawing is built on a fixed canvas, so a key sitting outside it
	// would panic. Run the whole catalogue through it.
	for _, e := range List() {
		l := MustGet(e.ID)
		if out := l.String(); out == "" {
			t.Errorf("%s rendered nothing", e.ID)
		}
		if out := l.With(ANSI).String(); out == "" {
			t.Errorf("%s rendered nothing on ANSI", e.ID)
		}
		if out := l.With(JIS).String(); out == "" {
			t.Errorf("%s rendered nothing on JIS", e.ID)
		}
	}
}

func TestStrokesAndType(t *testing.T) {
	us, ru := MustGet("kbdus"), MustGet("kbdru")

	// The classic wrong-layout result: the keys that spell "hello" on a US
	// keyboard spell this on a Russian one.
	if got := ru.Type(us.Strokes("hello")); got != "руддщ" {
		t.Errorf("ru.Type(us.Strokes(hello)) = %q, want %q", got, "руддщ")
	}

	// One stroke per rune, always.
	for _, in := range []string{"", "a", "abc", "a中b", "hello.world"} {
		if got, want := len(us.Strokes(in)), len([]rune(in)); got != want {
			t.Errorf("Strokes(%q) gave %d strokes for %d runes", in, got, want)
		}
	}

	// The plainest way to reach a character is the one chosen.
	if got := us.Strokes("a"); len(got) != 1 || got[0].Mod != Base || got[0].SC != "1E" {
		t.Errorf("Strokes(a) = %v, want the unmodified A key", got)
	}
	if got := us.Strokes("A"); len(got) != 1 || got[0].Mod != Shift || got[0].SC != "1E" {
		t.Errorf("Strokes(A) = %v, want shift on the A key", got)
	}

	// A character the layout cannot type has no stroke, and Type skips it.
	zero := us.Strokes("中")
	if len(zero) != 1 || !zero[0].Zero() {
		t.Errorf("Strokes(中) = %v, want one zero stroke", zero)
	}
	if got := ru.Type(zero); got != "" {
		t.Errorf("Type of a zero stroke produced %q", got)
	}

	// Typing strokes back on the layout they came from is the identity.
	const s = "hello.world"
	if got := us.Type(us.Strokes(s)); got != s {
		t.Errorf("us.Type(us.Strokes(%q)) = %q", s, got)
	}
}

func TestTranslate(t *testing.T) {
	us, ru, fr := MustGet("kbdus"), MustGet("kbdru"), MustGet("kbdfr")

	if got := us.Translate("hello", ru); got != "руддщ" {
		t.Errorf("us.Translate(hello, ru) = %q, want %q", got, "руддщ")
	}

	// AZERTY moves the punctuation, which is what makes this a squatting
	// vector rather than a curiosity.
	if got := us.Translate("google.com", fr); got != "google:co," {
		t.Errorf("us.Translate(google.com, fr) = %q, want %q", got, "google:co,")
	}

	// Translating to the same layout changes nothing.
	if got := us.Translate("google.com", us); got != "google.com" {
		t.Errorf("translating to the same layout gave %q", got)
	}

	// And it round-trips, because the strokes are the same either way.
	for _, in := range []string{"google.com", "paypal", "a1b2c3"} {
		there := us.Translate(in, ru)
		if back := ru.Translate(there, us); back != in {
			t.Errorf("%q -> %q -> %q", in, there, back)
		}
	}

	// Characters no key can produce are passed through rather than dropped.
	if got := us.Translate("a中b", ru); got != "ф中и" {
		t.Errorf("us.Translate(a中b, ru) = %q, want %q", got, "ф中и")
	}
}

// TestTranslateIsNotLengthPreserving pins a property it would be natural to
// assume and wrong to rely on: a ligature key types two characters at once, so
// one rune in can be two out.
func TestTranslateCanGrow(t *testing.T) {
	us, ar := MustGet("kbdus"), MustGet("kbda1")

	// The Arabic 101 board reaches "لا" on the key US uses for "b".
	k, ok := ar.Key("30")
	if !ok || len([]rune(k.Base())) != 2 {
		t.Skip("kbda1 SC 30 is no longer a ligature key")
	}

	got := us.Translate("b", ar)
	if len([]rune(got)) != 2 {
		t.Errorf("us.Translate(b, kbda1) = %q, %d runes; want the 2-rune ligature",
			got, len([]rune(got)))
	}
}

// TestTranslateRoundTripNeedsATypeableSource is the other assumption worth
// nailing down. Pass-through is not invertible: a character the source layout
// cannot type survives untouched, but it is an ordinary character on the way
// back and gets translated then.
func TestTranslateRoundTripNeedsATypeableSource(t *testing.T) {
	ru, gr := MustGet("kbdru"), MustGet("kbdgr")

	if ru.Types("a") {
		t.Skip("kbdru now types Latin; this case no longer applies")
	}

	there := ru.Translate("abc", gr)
	if there != "abc" {
		t.Errorf("Russian cannot type abc, so it should pass through; got %q", there)
	}
	if back := gr.Translate(there, ru); back == "abc" {
		t.Error("expected the return trip to translate what the outbound one passed through")
	}
}

// TestLookupsAgreeOnEmpty checks that the catalogue lookups answer "nothing"
// the same way. ByLanguage used to return nil for a full tag and an empty
// non-nil slice for a bare one, from the same call.
func TestLookupsAgreeOnEmpty(t *testing.T) {
	for name, got := range map[string][]Entry{
		"ByLanguage(bare)": ByLanguage("zz"),
		"ByLanguage(full)": ByLanguage("zz-ZZ"),
		"ByKeys":           ByKeys("中"),
		"ByString":         ByString("中"),
		"Find":             Find("zzzznotalayout"),
	} {
		if got != nil {
			t.Errorf("%s returned a non-nil empty slice; the others return nil", name)
		}
		if len(got) != 0 {
			t.Errorf("%s returned %d entries, want none", name, len(got))
		}
	}
}

// TestWideKeysReachTheirWholeLength is the space bar case. Measuring between
// key centers works while every key is one unit wide and falls apart when one
// is six: the bar's center sits under "b", so a center-to-center reading calls
// "b" its only neighbour and misses the rest of the row it physically covers.
func TestWideKeysReachTheirWholeLength(t *testing.T) {
	us := MustGet("kbdus")

	got := us.Adjacent(" ")
	for _, want := range []string{"x", "c", "v", "b", "n", "m"} {
		if indexOf(got, want) < 0 {
			t.Errorf("Adjacent(\" \") = %q, missing %q", got, want)
		}
	}

	// Not everything on the row, though: the bar starts to the right of "z",
	// so there is a real gap there.
	if indexOf(got, "z") >= 0 {
		t.Errorf("Adjacent(\" \") = %q; z sits left of where the bar starts", got)
	}

	// And the letters above it see the bar in turn.
	for _, c := range []string{"c", "v", "b", "n", "m"} {
		if indexOf(us.Adjacent(c), " ") < 0 {
			t.Errorf("Adjacent(%q) does not reach the space bar", c)
		}
	}
}

// TestStrikePointLeavesNormalKeysAlone pins the other half of the bargain: the
// strike-point rule must be a no-op for one-unit keys, or it would quietly
// re-rank every letter on every layout.
func TestStrikePointLeavesNormalKeysAlone(t *testing.T) {
	us := MustGet("kbdus")

	for _, sc := range []string{"12", "1E", "24", "02"} {
		k, ok := us.Key(sc)
		if !ok {
			t.Fatalf("no key at SC %s", sc)
		}
		cx, _ := k.Pos.Center()
		for _, from := range []float64{-5, 0, cx, cx + 3, 100} {
			if got := k.Pos.strike(from); got != cx {
				t.Errorf("SC %s strike(%v) = %v, want its centre %v", sc, from, got, cx)
			}
		}
	}

	// Distance between one-unit keys is still plain centre-to-centre.
	e, _ := us.Key("12")
	d, _ := us.Key("20")
	sKey, _ := us.Key("1F")
	if got := e.Pos.Distance(d.Pos); math.Abs(got-1.031) > 0.001 {
		t.Errorf("e->d = %.3f, want 1.031", got)
	}
	if got := e.Pos.Distance(sKey.Pos); math.Abs(got-1.25) > 0.001 {
		t.Errorf("e->s = %.3f, want 1.250", got)
	}
}

// TestDistanceIsSymmetric checks the property directly, including the mixed
// widths that made the naive fix asymmetric.
func TestDistanceIsSymmetric(t *testing.T) {
	for _, f := range []Form{ANSI, ISO, JIS} {
		var all []Pos
		for _, p := range geometries[f] {
			all = append(all, p)
		}
		for _, a := range all {
			for _, b := range all {
				if ab, ba := a.Distance(b), b.Distance(a); math.Abs(ab-ba) > 1e-9 {
					t.Fatalf("%s: %v->%v is %.4f but %v->%v is %.4f", f, a, b, ab, b, a, ba)
				}
			}
		}
	}
}

func indexOf(haystack []string, needle string) int {
	for i, s := range haystack {
		if s == needle {
			return i
		}
	}
	return -1
}

// The error paths below are the ones the package promises to take rather than
// guess. They are cheap to leave untested and expensive to have wrong.

func TestKeyAccessors(t *testing.T) {
	fr, _ := MustGet("kbdfr").Key("12") // e, with € on AltGr
	if got := fr.AltGr(); got != "€" {
		t.Errorf("AltGr = %q, want €", got)
	}
	if !fr.Types("e") || !fr.Types("E") || fr.Types("nope") {
		t.Error("Types disagrees with what the key produces")
	}

	// String is for diagnostics, so it just has to name the key and say
	// what it does.
	s := fr.String()
	if !strings.Contains(s, "12") || !strings.Contains(s, `"e"`) {
		t.Errorf("Key.String = %q", s)
	}

	// A key that types nothing still renders, with its label.
	blank := NewKey("3A", "VK_CAPITAL", "Caps Lock")
	if got := blank.String(); !strings.Contains(got, "Caps Lock") {
		t.Errorf("blank Key.String = %q", got)
	}

	// A dead key says so.
	de, _ := MustGet("kbdgr").Key("29")
	if got := de.String(); !strings.Contains(got, "dead") {
		t.Errorf("dead Key.String = %q", got)
	}
}

func TestModStringAndParse(t *testing.T) {
	for _, m := range mods {
		name := m.String()
		back, ok := parseMod(name)
		if !ok || back != m {
			t.Errorf("%v -> %q -> %v (%v)", m, name, back, ok)
		}
	}

	// A modifier set outside the modelled range names itself numerically
	// rather than pretending to be one of the known ones.
	if got := Mod(200).String(); got != "mod(200)" {
		t.Errorf("Mod(200).String() = %q", got)
	}
	if _, ok := parseMod("notamodifier"); ok {
		t.Error("parseMod accepted a name it does not know")
	}
}

func TestSetKeepsTheFirstValue(t *testing.T) {
	// The source data lists the plainest spelling first, so a second write
	// to the same state is ignored.
	k := NewKey("10", "VK_Q", "")
	k.Set(Base, Out{Text: "q"})
	k.Set(Base, Out{Text: "OVERWRITTEN"})

	if k.Base() != "q" {
		t.Errorf("Base = %q, want the first value", k.Base())
	}
}

func TestCompactLeavesInformativeStates(t *testing.T) {
	// Caps that matches the fallback goes; caps that carries its own value
	// stays.
	ordinary := NewKey("10", "VK_Q", "")
	ordinary.Set(Base, Out{Text: "q"})
	ordinary.Set(Shift, Out{Text: "Q"})
	ordinary.Set(Caps, Out{Text: "Q"})
	ordinary.Compact()
	if got := ordinary.Mods(); len(got) != 2 {
		t.Errorf("ordinary key kept %v after compacting", got)
	}

	special := NewKey("03", "VK_2", "")
	special.Set(Base, Out{Text: "2"})
	special.Set(Shift, Out{Text: `"`})
	special.Set(Caps, Out{Text: "2"}) // Caps Lock does nothing here
	special.Compact()
	if got := special.Mods(); len(got) != 3 {
		t.Errorf("informative caps state was compacted away: %v", got)
	}
}

func TestStrikeOnAKeyNarrowerThanOne(t *testing.T) {
	// Nothing in the dataset is narrower than a key, but the rule has to
	// mean something if one ever is: there is only one place to hit it.
	narrow := Pos{X: 2, Y: 0, W: 0.5}
	cx, _ := narrow.Center()
	for _, from := range []float64{-10, 0, 2.25, 100} {
		if got := narrow.strike(from); got != cx {
			t.Errorf("strike(%v) = %v, want the centre %v", from, got, cx)
		}
	}
}

func TestIDs(t *testing.T) {
	ids := IDs()
	if len(ids) != len(List()) {
		t.Errorf("IDs returned %d, List %d", len(ids), len(List()))
	}
	if !sort.StringsAreSorted(ids) {
		t.Error("IDs is not sorted")
	}
	if indexOf(ids, "kbdus") < 0 {
		t.Error("IDs does not include kbdus")
	}
}

func TestMustGetPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustGet should panic on a layout that does not exist")
		}
	}()
	MustGet("nosuchlayout")
}

func TestEncodingRejectsBadInput(t *testing.T) {
	// A KLID that is not hex cannot be stored as the number it is meant to
	// be, and saying so beats writing a zero.
	if _, err := encodeKLID("not-hex"); err == nil {
		t.Error("encodeKLID accepted a non-hex KLID")
	}
	if _, err := encodeLocales([]Locale{{KLID: "zzzz"}}); err == nil {
		t.Error("encodeLocales accepted a non-hex KLID")
	}

	// The same on the way out of a layout and out of the catalogue.
	bad := New("t", "T", ANSI, []Key{NewKey("10", "VK_Q", "")})
	bad.Locales = []Locale{{KLID: "zzzz"}}
	if _, err := bad.Marshal(); err == nil {
		t.Error("Marshal accepted a layout with a non-hex KLID")
	}
	if _, err := MarshalCatalogue([]Entry{{ID: "t", Form: ANSI, Locales: []Locale{{KLID: "zzzz"}}}}); err == nil {
		t.Error("MarshalCatalogue accepted a non-hex KLID")
	}
	if _, err := MarshalCatalogue([]Entry{{ID: "t", Form: Form("bogus")}}); err == nil {
		t.Error("MarshalCatalogue accepted an unknown form")
	}

	// A scan code that is not hex either.
	if _, err := New("t", "T", ANSI, []Key{{SC: "zz"}}).Marshal(); err != nil {
		// index() drops it for having no position, so this must not fail
		t.Errorf("a key outside the block should be dropped, not an error: %v", err)
	}
}

func TestDecodingRejectsBadInput(t *testing.T) {
	if err := new(Layout).Unmarshal([]byte("\xff\xff\xff\xff")); err == nil {
		t.Error("Unmarshal accepted rubbish")
	}
	if _, err := UnmarshalCatalogue([]byte("\xff\xff\xff\xff")); err == nil {
		t.Error("UnmarshalCatalogue accepted rubbish")
	}

	// A catalogue entry with no form is not one this package can place.
	raw, err := proto.Marshal(&kbpb.Catalogue{Entries: []*kbpb.Entry{{Id: "t"}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := UnmarshalCatalogue(raw); err == nil {
		t.Error("UnmarshalCatalogue accepted an entry with no form")
	}
}

func TestJSONRejectsBadInput(t *testing.T) {
	if err := new(Key).UnmarshalJSON([]byte(`{"sc":"10","out":{"nosuchstate":{"t":"q"}}}`)); err == nil {
		t.Error("Key.UnmarshalJSON accepted an unknown modifier state")
	}
	if err := new(Key).UnmarshalJSON([]byte(`not json`)); err == nil {
		t.Error("Key.UnmarshalJSON accepted rubbish")
	}
	if err := new(Layout).UnmarshalJSON([]byte(`not json`)); err == nil {
		t.Error("Layout.UnmarshalJSON accepted rubbish")
	}
	if err := new(Layout).UnmarshalJSON([]byte(`{"keys":[{"sc":"10","out":{"bad":{"t":"q"}}}]}`)); err == nil {
		t.Error("Layout.UnmarshalJSON accepted an unknown modifier state")
	}
}
