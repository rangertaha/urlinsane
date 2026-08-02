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

package variant

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rangertaha/urlinsane/datasets"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

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
	// registered language plugin, in id order.
	Languages []internal.Language
	// Keyboards restricts the keyboard-driven algorithms. Nil means every
	// registered keyboard plugin, in id order.
	Keyboards []internal.Keyboard
	// Subdomains feeds the subdomain-insertion algorithm. Nil means the
	// compiled-in list.
	Subdomains []string
	// Suffixes feeds the TLD-swap algorithm. Nil means the compiled-in public
	// suffix list, wildcard and exception markers stripped.
	Suffixes []string
}

func (o Options) withDefaults() Options {
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
	return o
}

// Specs returns every algorithm declaration, in id order.
func Specs(o Options) []Spec {
	o = o.withDefaults()
	specs := PureSpecs()
	specs = append(specs, DomainSpecs(o.Subdomains, o.Suffixes)...)
	specs = append(specs, KeyboardSpecs(o.Keyboards)...)
	specs = append(specs, LanguageSpecs(o.Languages)...)
	sort.Slice(specs, func(i, j int) bool { return specs[i].ID < specs[j].ID })
	return specs
}

// All returns every algorithm as an operator, in id order. Two ids that collide
// are a programming error rather than a runtime condition — the scheduler keys
// its seen-set and cache on the id — so this panics rather than silently
// dropping one.
func All(o Options) []graph.Operator {
	o = o.withDefaults()
	specs := Specs(o)
	seen := make(map[string]bool, len(specs))
	ops := make([]graph.Operator, 0, len(specs))
	for _, s := range specs {
		if seen[s.ID] {
			panic(fmt.Sprintf("variant: duplicate operator id %q", s.ID))
		}
		seen[s.ID] = true
		ops = append(ops, New(s, o.Split))
	}
	return ops
}

// Select returns the named algorithms as operators, in id order. An unknown id
// is an error rather than a silent omission: a run that quietly dropped half
// the algorithms the user asked for would report an incomplete graph as
// complete.
func Select(o Options, ids ...string) ([]graph.Operator, error) {
	want := make(map[string]bool, len(ids))
	for _, id := range ids {
		want[id] = true
	}
	var ops []graph.Operator
	for _, op := range All(o) {
		if want[op.Id()] {
			delete(want, op.Id())
			ops = append(ops, op)
		}
	}
	if len(want) > 0 {
		missing := make([]string, 0, len(want))
		for id := range want {
			missing = append(missing, id)
		}
		sort.Strings(missing)
		return nil, fmt.Errorf("variant: unknown algorithm(s): %s", strings.Join(missing, ", "))
	}
	return ops, nil
}

// RegisteredLanguages returns every language plugin currently registered, in id
// order. The plugin registry is a map, so iterating it directly would order the
// algorithms' inputs differently on every run.
func RegisteredLanguages() []internal.Language {
	langs := languages.Languages()
	sort.Slice(langs, func(i, j int) bool { return langs[i].Id() < langs[j].Id() })
	return langs
}

// RegisteredKeyboards returns every keyboard plugin currently registered, in id
// order, for the same reason.
func RegisteredKeyboards() []internal.Keyboard {
	boards := languages.Keyboards()
	sort.Slice(boards, func(i, j int) bool { return boards[i].Id() < boards[j].Id() })
	return boards
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
