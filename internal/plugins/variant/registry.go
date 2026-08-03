// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package variant

import (
	"sort"
	"strings"

	"github.com/rangertaha/urlinsane/datasets"
	"github.com/rangertaha/urlinsane/internal/dataset"
	"github.com/rangertaha/urlinsane/pkg/kb"
)

// Language is the data the language-driven algorithms read.
//
// Declared here, by the consumer, rather than imported from wherever the data
// comes from. internal/lang satisfies it by querying the dataset, and a test
// satisfies it with a literal — which is the point: an operator is a pure
// function of its inputs (§4.3), and it cannot be tested as one if constructing
// its input needs a database.
type Language interface {
	Code() string
	Name() string
	Vowels() []string
	Graphemes() []string
	Numerals() map[string][]string
	Homoglyphs() map[string][]string
	Homophones() [][]string
	Misspellings() [][]string
}

// Options selects the data the algorithms run over. Every field is optional;
// the zero Options gives the full shipped set.
//
// The datasets are parameters rather than globals read at Exec time because an
// operator is cached against a digest of its declared reads, and a plugin
// registered after the plan compiled would change results without changing the
// digest.
type Options struct {
	// Split decomposes a key by node type. Nil means DefaultSplit.
	Split Splitter
	// Languages restricts the language-driven algorithms. Nil means every
	// registered language, in id order.
	Languages []Language
	// Keyboards restricts the keyboard-driven algorithms. Nil means every
	// layout pkg/kb ships, deduplicated by Adjacency behaviour: 203 layouts
	// collapse to about 30 distinct neighbour sets, because most Latin boards
	// share QWERTY geometry. Running all 203 would do seven times the work for
	// the same candidates.
	Keyboards []*kb.Layout
	// Subdomains feeds the subdomain-insertion algorithm. Nil means the
	// compiled-in list.
	Subdomains []string
	// Suffixes feeds the TLD-swap algorithm. Nil means the compiled-in public
	// suffix list, wildcard and exception markers stripped.
	Suffixes []string
	// Providers feeds the delegated-subdomain algorithm: hosts that hand out
	// subdomains to anyone. Nil reads the public suffix list's private section
	// from the dataset.
	Providers []string
	// CrossHomophones feeds the cross-language homophone algorithm: groups of
	// spellings that sound alike to speakers of different languages. Nil reads
	// them from the dataset.
	//
	// Not derived from Languages, and it cannot be: a group is in the file
	// precisely because it crosses a language boundary, so narrowing it to the
	// run's languages would empty it.
	CrossHomophones [][]string
	// Extra are algorithms contributed from outside this package, which is how
	// a plugin adds one (internal/plugins). Passed in rather than read from a
	// registry here, because the registry imports this package to declare a
	// Spec and importing it back would be a cycle.
	Extra []Spec
	// Combos feeds the combosquatting algorithm. Nil means the default
	// language's vocabulary; an empty non-nil slice means no keywords, which
	// silences the algorithm rather than defaulting it back on.
	//
	// It is deliberately not derived from Languages. Combo vocabulary is a
	// property of the target's audience rather than of the name, so widening
	// the language set to catch more homoglyphs must not also bolt another
	// language's keywords onto the name. Pass ComboKeywords(langs) to opt into
	// the multi-language vocabulary; see DefaultComboLanguage for why that is
	// not the default.
	Combos []string
}

func (o Options) WithDefaults() Options {
	if o.Split == nil {
		o.Split = DefaultSplit
	}
	if o.Languages == nil {
		o.Languages = RegisteredLanguages()
	}
	if o.Keyboards == nil {
		o.Keyboards = RegisteredKeyboards()
	}
	if o.Subdomains == nil {
		o.Subdomains = datasets.SUBDOMAINS
	}
	if o.Suffixes == nil {
		o.Suffixes = PublicSuffixes()
	}
	if o.CrossHomophones == nil {
		o.CrossHomophones = dataset.Groups("phonetics/homophone")
	}
	if o.Providers == nil {
		o.Providers = dataset.Tokens("domains/private")
	}
	if o.Combos == nil {
		o.Combos = DefaultComboKeywords()
	}
	return o
}

// RegisteredLanguages returns every registered language, in id order.
//
// Sorted because an operator built from this set is cached against a digest of
// what it read: registration order is package-initialisation order, and letting
// it through would invalidate the cache on every run.
func RegisteredLanguages() []Language {
	all := dataset.Languages()
	out := make([]Language, 0, len(all))
	for _, l := range all {
		out = append(out, l)
	}
	return out
}

// RegisteredKeyboards returns the layouts the keyboard-driven algorithms run
// over: every layout pkg/kb ships, one per distinct Adjacency behaviour.
//
// Deduplicated rather than enumerated. The 203 shipped layouts produce about 30
// distinct neighbour sets, because most Latin boards share QWERTY geometry —
// so running all of them multiplies the work sevenfold and the candidates not
// at all. The first layout of each behaviour wins, in id order, so the choice
// is stable across runs.
func RegisteredKeyboards() []*kb.Layout {
	seen := map[string]bool{}
	var out []*kb.Layout
	for _, id := range kb.IDs() {
		l, err := kb.Get(id)
		if err != nil {
			continue
		}
		if sig := adjacencySignature(l); !seen[sig] {
			seen[sig] = true
			out = append(out, l)
		}
	}
	return out
}

// adjacencySignature identifies a layout by how it types, not by what it is
// called. Two layouts that differ only in their non-alphanumeric keys produce
// the same variants, so they are the same layout for this purpose.
func adjacencySignature(l *kb.Layout) string {
	var b strings.Builder
	for _, c := range "abcdefghijklmnopqrstuvwxyz0123456789" {
		adj := l.Adjacent(string(c))
		sort.Strings(adj)
		b.WriteString(string(c))
		b.WriteByte(':')
		b.WriteString(strings.Join(adj, ""))
		b.WriteByte('|')
	}
	return b.String()
}

// PublicSuffixes returns the public suffix list as plain suffixes, sorted, with
// the PSL's wildcard and exception markers removed. Those markers are matching
// instructions, not names: "*.ck" is not a domain anyone can register, and
// emitting it as a variant would produce a key nothing could ever resolve.
func PublicSuffixes() []string {
	out := make([]string, 0, len(datasets.TLD))
	seen := make(map[string]bool, len(datasets.TLD))
	for _, rule := range datasets.TLD {
		rule = strings.TrimSpace(rule)
		switch {
		case rule == "", strings.HasPrefix(rule, "*."):
			continue
		case strings.HasPrefix(rule, "!"):
			rule = rule[1:]
		}
		if rule == "" || seen[rule] {
			continue
		}
		seen[rule] = true
		out = append(out, rule)
	}
	sort.Strings(out)
	return out
}
