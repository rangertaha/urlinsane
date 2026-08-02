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
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
	"sync"
)

// The dataset, compiled into any program that imports this package. It is one
// encoded Layout per file plus an index, so that selecting a layout costs the
// index alone and only the layouts actually used are decoded.
//
//go:embed data
var data embed.FS

const (
	dataDir       = "data"
	layoutExt     = ".pb"
	catalogueFile = dataDir + "/index" + layoutExt
)

// Entry is the summary of a layout in the catalogue, enough to choose one
// without paying to parse its keys.
type Entry struct {
	// ID is the driver name, lowercased: "kbdus", "kbdfr".
	ID string `json:"id"`

	// Name is a human-readable name: "US", "French".
	Name string `json:"name"`

	// Form is the physical shape of the board.
	Form Form `json:"form"`

	// Locales are the language identities the layout is installed under.
	Locales []Locale `json:"locales,omitempty"`
}

// Languages returns the distinct primary language subtags for the entry.
func (e Entry) Languages() []string {
	l := Layout{Locales: e.Locales}
	return l.Languages()
}

// clone detaches an entry from the catalogue. An Entry copies by value, but
// its locales are a slice, so handing one out unchanged would let a caller
// write through to the catalogue every other caller reads.
func (e Entry) clone() Entry {
	e.Locales = cloneLocales(e.Locales)
	return e
}

var catalogue struct {
	once sync.Once
	err  error

	entries []Entry
	byID    map[string]int // driver name and every KLID, lowercased
	byLang  map[string][]int
}

// load reads the index once. The index is small; the layouts themselves are
// decoded lazily, one file at a time, only when asked for.
func load() error {
	catalogue.once.Do(func() {
		raw, err := data.ReadFile(catalogueFile)
		if err != nil {
			catalogue.err = fmt.Errorf("kb: reading catalogue: %w", err)
			return
		}

		entries, err := UnmarshalCatalogue(raw)
		if err != nil {
			catalogue.err = fmt.Errorf("kb: parsing catalogue: %w", err)
			return
		}

		sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })

		catalogue.entries = entries
		catalogue.byID = make(map[string]int, len(entries)*2)
		catalogue.byLang = make(map[string][]int)

		for i, e := range entries {
			catalogue.byID[e.ID] = i
			for _, loc := range e.Locales {
				// A KLID belongs to exactly one driver, so this cannot
				// collide the way the driver names could.
				catalogue.byID[strings.ToLower(loc.KLID)] = i
			}
			for _, lang := range e.Languages() {
				catalogue.byLang[lang] = append(catalogue.byLang[lang], i)
			}
		}
	})

	return catalogue.err
}

// candidates lists the catalogue keys an identifier could mean, most literal
// first. It accepts what a caller might reasonably have to hand: a driver name
// as "KBDUS" or "kbdus.dll", and a KLID as "00000409", "0x409", or the bare
// "409" that Windows tools often print.
func candidates(id string) []string {
	id = strings.ToLower(strings.TrimSpace(id))
	id = strings.TrimSuffix(id, ".dll")

	out := []string{id}

	hex := strings.TrimPrefix(id, "0x")
	if hex != id {
		out = append(out, hex)
	}

	// Zero-pad a short KLID, but only after its literal spelling has been
	// tried, so that a name which happens to be all hex digits still wins.
	if n := len(hex); n > 0 && n < 8 && isHex(hex) {
		out = append(out, strings.Repeat("0", 8-n)+hex)
	}

	return out
}

// lookup resolves an identifier to a catalogue index.
func lookup(id string) (int, bool) {
	for _, key := range candidates(id) {
		if i, ok := catalogue.byID[key]; ok {
			return i, true
		}
	}
	return 0, false
}

func isHex(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune("0123456789abcdef", r) {
			return false
		}
	}
	return true
}

var loaded struct {
	sync.Mutex
	layouts map[string]*Layout
}

// Get returns the layout for a driver name or KLID. Layouts are parsed on
// first use and cached, so repeated calls are cheap and return the same
// pointer. The returned layout must not be modified.
func Get(id string) (*Layout, error) {
	if err := load(); err != nil {
		return nil, err
	}

	i, ok := lookup(id)
	if !ok {
		return nil, fmt.Errorf("kb: no layout %q", id)
	}
	name := catalogue.entries[i].ID

	loaded.Lock()
	defer loaded.Unlock()

	if l, ok := loaded.layouts[name]; ok {
		return l, nil
	}

	raw, err := data.ReadFile(path.Join(dataDir, name+layoutExt))
	if err != nil {
		return nil, fmt.Errorf("kb: reading layout %q: %w", name, err)
	}

	l := new(Layout)
	if err := l.Unmarshal(raw); err != nil {
		return nil, fmt.Errorf("kb: parsing layout %q: %w", name, err)
	}

	if loaded.layouts == nil {
		loaded.layouts = make(map[string]*Layout)
	}
	loaded.layouts[name] = l

	return l, nil
}

// MustGet is Get for layouts known to exist, such as constants in a program's
// own source. It panics if the layout is missing.
func MustGet(id string) *Layout {
	l, err := Get(id)
	if err != nil {
		panic(err)
	}
	return l
}

// List returns the catalogue of available layouts, ordered by ID. It does not
// parse any of them.
func List() []Entry {
	if err := load(); err != nil {
		return nil
	}
	out := make([]Entry, 0, len(catalogue.entries))
	for _, e := range catalogue.entries {
		out = append(out, e.clone())
	}
	return out
}

// IDs returns the driver name of every available layout, in order.
func IDs() []string {
	entries := List()
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.ID
	}
	return out
}

// ByLanguage returns the layouts installed for a language, given either a
// primary subtag ("de") or a full BCP-47 tag ("de-CH"). A full tag matches
// only the layouts carrying that exact locale; a bare subtag matches every
// layout for the language.
func ByLanguage(tag string) []Entry {
	if err := load(); err != nil {
		return nil
	}

	tag = strings.ToLower(strings.TrimSpace(tag))
	if tag == "" {
		return nil
	}

	if primary, _, ok := strings.Cut(tag, "-"); ok {
		var out []Entry
		for _, i := range catalogue.byLang[primary] {
			for _, loc := range catalogue.entries[i].Locales {
				if strings.EqualFold(loc.Tag, tag) {
					out = append(out, catalogue.entries[i].clone())
					break
				}
			}
		}
		return out
	}

	idx := catalogue.byLang[tag]
	out := make([]Entry, 0, len(idx))
	for _, i := range idx {
		out = append(out, catalogue.entries[i].clone())
	}
	return out
}

// ByKeys returns the layouts that can type every one of the given texts, in
// catalogue order. It answers "which keyboards is this reachable on" — the
// characters of a domain, a name, a package identifier:
//
//	kb.ByKeys("g", "o", "l", "e", ".", "c", "m")
//	kb.ByKeys(strings.Split("münchen", "")...)
//
// A text matches if some key on the layout produces exactly it, in any
// modifier state, so "a" and "A" are different questions and a layout that
// reaches a character only through AltGr still counts. Passing nothing returns
// nothing, as does passing a text no layout can type.
//
// Unlike ByLanguage this cannot be answered from the index: it has to look at
// the keys, so the first call decodes the whole catalogue. They are cached
// afterwards, and later calls are a map lookup per layout.
func ByKeys(texts ...string) []Entry {
	if err := load(); err != nil {
		return nil
	}
	if len(texts) == 0 {
		return nil
	}

	var out []Entry
	for _, e := range catalogue.entries {
		l, err := Get(e.ID)
		if err != nil {
			continue
		}

		typesAll := true
		for _, t := range texts {
			if !l.Types(t) {
				typesAll = false
				break
			}
		}
		if typesAll {
			out = append(out, e.clone())
		}
	}

	return out
}

// ByString returns the layouts that can type every character of s, in
// catalogue order. It is ByKeys over the runes of s, deduplicated:
//
//	kb.ByString("google.com")
//	kb.ByString("münchen")
//
// The question it answers is whether every character sits on a key of its own.
// A layout that reaches a character only by dead key — an acute accent on a
// German keyboard, pressed before "e" to make "é" — does not count, because
// the dataset records that the accent key is dead but not what it composes
// into. Characters written as a base letter followed by a combining mark are
// likewise treated as the two runes they are.
func ByString(s string) []Entry {
	if s == "" {
		return nil
	}

	seen := make(map[string]bool, len(s))
	texts := make([]string, 0, len(s))
	for _, r := range s {
		t := string(r)
		if seen[t] {
			continue
		}
		seen[t] = true
		texts = append(texts, t)
	}

	return ByKeys(texts...)
}

// Languages returns every language subtag with at least one layout, sorted.
func Languages() []string {
	if err := load(); err != nil {
		return nil
	}
	out := make([]string, 0, len(catalogue.byLang))
	for lang := range catalogue.byLang {
		out = append(out, lang)
	}
	sort.Strings(out)
	return out
}

// Find returns the layouts whose ID, name, or locale names contain q, matched
// case-insensitively. It is for interactive selection — a "--keyboard french"
// flag — rather than for programmatic lookup, which should use Get.
func Find(q string) []Entry {
	if err := load(); err != nil {
		return nil
	}

	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return nil
	}

	var out []Entry
	for _, e := range catalogue.entries {
		if strings.Contains(strings.ToLower(e.ID), q) || strings.Contains(strings.ToLower(e.Name), q) {
			out = append(out, e.clone())
			continue
		}
		for _, loc := range e.Locales {
			if strings.Contains(strings.ToLower(loc.Name), q) || strings.EqualFold(loc.Tag, q) {
				out = append(out, e.clone())
				break
			}
		}
	}
	return out
}

// files lists the layout files present, for the tests that check the dataset
// and the catalogue agree.
func files() ([]string, error) {
	entries, err := fs.ReadDir(data, dataDir)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, layoutExt) || dataDir+"/"+name == catalogueFile {
			continue
		}
		out = append(out, strings.TrimSuffix(name, layoutExt))
	}

	return out, nil
}
