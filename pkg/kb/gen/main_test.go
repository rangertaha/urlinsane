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
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rangertaha/urlinsane/pkg/kb"
)

func TestCodepoints(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want string
		ok   bool
	}{
		{"0041", "A", true},
		{"41", "A", true},   // odd length, padded
		{"0020", " ", true}, // space is graphic
		{"1F600", "\U0001F600", true},
		{"0041 0301", "Á", true}, // a base letter and a combining mark
		{"001B", "", false},       // escape: a control code, not text
		{"0000", "", false},
		{"00410301", "", false}, // one number past the end of Unicode
		{"zz", "", false},
		{"", "", false},
	} {
		got, ok := codepoints(tc.in)
		if got != tc.want || ok != tc.ok {
			t.Errorf("codepoints(%q) = %q, %v; want %q, %v", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestModifiers(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want kb.Mod
		ok   bool
	}{
		{"", kb.Base, true},
		{"VK_SHIFT", kb.Shift, true},
		{"VK_CAPITAL", kb.Caps, true},
		{"VK_SHIFT VK_CAPITAL", kb.Shift | kb.Caps, true},
		{"VK_CONTROL VK_MENU", kb.AltGr, true},
		{"VK_RMENU", kb.AltGr, true},
		{"VK_SHIFT VK_CONTROL VK_MENU", kb.Shift | kb.AltGr, true},
		{"VK_CONTROL VK_MENU VK_SHIFT VK_CAPITAL", kb.Shift | kb.AltGr | kb.Caps, true},

		// Ctrl and Alt alone give control codes, not text.
		{"VK_CONTROL", 0, false},
		{"VK_MENU", 0, false},
		{"VK_SHIFT VK_CONTROL", 0, false},

		// Levels this package has no bit for. Dropping them is deliberate,
		// and the run reports what went with them.
		{"VK_KANA", 0, false},
		{"VK_OEM_8", 0, false},
		{"VK_NUMLOCK", 0, false},
	} {
		got, ok := modifiers(tc.in)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("modifiers(%q) = %v, %v; want %v, %v", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestForm(t *testing.T) {
	key := func(sc, base string) physicalKey {
		return physicalKey{SC: sc, Results: []result{{Text: base}}}
	}

	for _, tc := range []struct {
		name string
		keys []physicalKey
		want kb.Form
	}{
		// Windows defines SC 56 on every layout. One drawn for ANSI just
		// repeats the backslash already on SC 2B there; one drawn for ISO
		// gives the extra switch a character of its own.
		{"repeats 2B", []physicalKey{key("56", `\`), key("2B", `\`)}, kb.ANSI},
		{"differs from 2B", []physicalKey{key("56", "<"), key("2B", "*")}, kb.ISO},
		{"no 56 at all", []physicalKey{key("2B", `\`)}, kb.ANSI},
		{"56 types nothing", []physicalKey{key("56", ""), key("2B", `\`)}, kb.ANSI},
		{"no keys", nil, kb.ANSI},
	} {
		if got := form(tc.keys); got != tc.want {
			t.Errorf("form(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestName(t *testing.T) {
	// Where several locales share a driver, the shortest name is the one
	// naming the layout rather than the locale.
	got := name("kbdus", meta{Locales: []kb.Locale{
		{Name: "Chinese (Traditional) - US"},
		{Name: "US"},
		{Name: "Bulgarian (Latin)"},
	}})
	if got != "US" {
		t.Errorf("name = %q, want %q", got, "US")
	}

	// Ties break alphabetically, so the result does not depend on the order
	// the page happened to list them in.
	if got := name("kbdxx", meta{Locales: []kb.Locale{{Name: "BBB"}, {Name: "AAA"}}}); got != "AAA" {
		t.Errorf("name = %q, want AAA", got)
	}

	// With nothing to go on, fall back to the driver.
	if got := name("kbdxx", meta{Locales: []kb.Locale{{Name: ""}}}); got != "KBDXX" {
		t.Errorf("name = %q, want KBDXX", got)
	}
}

func TestResultText(t *testing.T) {
	dead := result{}
	dead.Dead = &struct {
		Accent string `xml:"Accent,attr"`
	}{Accent: "^"}

	for _, tc := range []struct {
		name string
		in   result
		want string
	}{
		{"plain", result{Text: "a"}, "a"},
		{"codepoints", result{Codepoints: "0041"}, "A"},
		{"control code", result{Codepoints: "001B"}, ""},
		{"dead key", dead, "^"},
		{"nothing", result{}, ""},
	} {
		if got := resultText(tc.in); got != tc.want {
			t.Errorf("resultText(%s) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestWriteAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "entry.cache")

	if err := writeAtomic(path, []byte("first")); err != nil {
		t.Fatal(err)
	}
	if err := writeAtomic(path, []byte("second")); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "second" {
		t.Errorf("file holds %q, want the second write", got)
	}

	// The temporary file must not be left behind.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("%d files in the cache directory, want just the entry", len(entries))
	}
}
