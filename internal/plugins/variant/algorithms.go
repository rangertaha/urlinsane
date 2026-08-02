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
	"strings"

	"github.com/rangertaha/urlinsane/pkg/typo"
)

// PureSpecs are the algorithms that need no dataset at all: they are functions
// of the name and nothing else. They are listed first because they are the ones
// that can run in any deployment, with no language plugin loaded and no dataset
// database open.
//
// Ids are the short codes the CLI has always used, so `--algorithms co,cs`
// keeps meaning what it meant.
func PureSpecs() []Spec {
	return []Spec{
		{
			ID: "co", Title: "Character Omission", Version: 1,
			Gen: typo.CharacterOmission,
		},
		{
			ID: "cr", Title: "Character Repetition", Version: 1,
			Gen: typo.CharacterRepetition,
		},
		{
			ID: "cs", Title: "Character Swapping", Version: 1,
			Gen: typo.CharacterSwapping,
		},
		{
			ID: "hi", Title: "Hyphen Insertion", Version: 1,
			Gen: typo.HyphenInsertion,
		},
		{
			ID: "ho", Title: "Hyphen Omission", Version: 1,
			Gen: typo.HyphenOmission,
		},
		{
			ID: "di", Title: "Dot Insertion", Version: 1,
			Gen: typo.DotInsertion,
		},
		{
			ID: "do", Title: "Dot Omission", Version: 1,
			Gen: typo.DotOmission,
		},
		{
			ID: "dhs", Title: "Dot Hyphen Substitution", Version: 1,
			Gen: typo.DotHyphenSubstitution,
		},
		{
			// Bitsquatting: a single flipped bit in a transmitted name. It is
			// the one algorithm here that models hardware error rather than
			// human error, which is why it ignores keyboards and languages.
			ID: "bf", Title: "Bit Flipping", Version: 1,
			Gen: func(name string) []string { return typo.BitFlipping(name) },
		},
		{
			ID: "sp", Title: "Singular Pluralise", Version: 1,
			Gen: typo.SingularPluralise,
		},
		{
			// Affix squatting is a whole-key algorithm: "py-" belongs in front
			// of the package name, not in front of a domain's registrable
			// label, and the ecosystem affixes only make sense unsplit.
			ID: "afx", Title: "Affix Squatting", Version: 1,
			Types: []string{TypePackage, TypeRepo, TypeUsername},
			Whole: true,
			Gen:   affixSquatting,
		},
		{
			ID: "sep", Title: "Separator Substitution", Version: 1,
			Types: []string{TypePackage, TypeRepo, TypeUsername},
			Whole: true,
			Gen:   separatorSubstitution,
		},
		{
			ID: "nsc", Title: "Namespace Confusion", Version: 1,
			Types: []string{TypePackage, TypeRepo},
			Whole: true,
			Gen:   namespaceConfusion,
		},
	}
}

// affixes are the common ecosystem/role brackets seen in real registry
// squatting. Ported verbatim from the afx plugin.
var (
	affixPrefixes = []string{"py", "py-", "python-", "node-", "js-", "go-", "lib", "lib-", "the-", "get-"}
	affixSuffixes = []string{"2", "3", "js", "-js", "py", "-py", "-python", "-cli", "-dev",
		"-core", "-utils", "-api", "-sdk", "-lib", ".js", "-ng", "-master", "-official"}
)

// affixSquatting brackets the name with each common ecosystem affix.
func affixSquatting(name string) []string {
	out := make([]string, 0, len(affixPrefixes)+len(affixSuffixes))
	for _, p := range affixPrefixes {
		out = append(out, p+name)
	}
	for _, s := range affixSuffixes {
		out = append(out, name+s)
	}
	return out
}

// separators are the word joiners a registry name may legally use. The empty
// string is one of them: "my-lib" and "mylib" are different packages.
var separators = []string{"-", "_", ".", ""}

// separatorSubstitution re-joins the name's word tokens with each alternate
// separator. A name with a single token has nothing to re-separate.
func separatorSubstitution(name string) []string {
	tokens := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_' || r == '.'
	})
	if len(tokens) < 2 {
		return nil
	}
	out := make([]string, 0, len(separators))
	for _, s := range separators {
		out = append(out, strings.Join(tokens, s))
	}
	return out
}

// namespaceConfusion moves a name between the namespacing conventions the
// registries use — npm's "@org/pkg", a repo's "org/pkg", and the flat
// "org-pkg" — which is how a scoped package gets impersonated by an unscoped
// one and vice versa.
func namespaceConfusion(name string) []string {
	switch {
	case strings.HasPrefix(name, "@") && strings.Contains(name, "/"):
		org, pkg, _ := strings.Cut(strings.TrimPrefix(name, "@"), "/")
		return []string{pkg, org + "-" + pkg, org + "/" + pkg, org + pkg}
	case strings.Contains(name, "/"):
		org, pkg, _ := strings.Cut(name, "/")
		return []string{pkg, org + "-" + pkg, "@" + org + "/" + pkg}
	case strings.Contains(name, "-"):
		org, pkg, _ := strings.Cut(name, "-")
		return []string{"@" + org + "/" + pkg, org + "/" + pkg}
	}
	return nil
}

// DomainSpecs are the algorithms that only make sense for a domain, driven by
// the compiled-in subdomain and public-suffix datasets. They are separate from
// PureSpecs because they are the two that bind by type rather than capability:
// a package name has no TLD to get wrong.
func DomainSpecs(subdomains, suffixes []string) []Spec {
	return []Spec{
		{
			ID: "si", Title: "Subdomain Insertion", Version: 1,
			Types: []string{TypeDomain},
			Whole: true,
			Gen:   subdomainInsertion(subdomains),
		},
		{
			ID: "tld", Title: "Wrong TLD", Version: 1,
			Types: []string{TypeDomain},
			Whole: true,
			Gen:   tldSwap(suffixes),
		},
	}
}

// subdomainInsertion prepends each known subdomain label to the registrable
// domain. It replaces any existing prefix rather than stacking on top of it,
// which is why it varies the whole key and re-splits itself.
func subdomainInsertion(subdomains []string) Generate {
	return func(name string) []string {
		_, core, suffix := SplitDomain(name)
		if core == "" {
			return nil
		}
		out := make([]string, 0, len(subdomains))
		for _, sub := range subdomains {
			out = append(out, joinDomain(sub, core, suffix))
		}
		return out
	}
}

// tldSwap replaces the public suffix with every other known suffix, keeping any
// subdomain prefix. This is the one algorithm allowed to change the registry a
// name lives in; every other algorithm preserves the suffix precisely so this
// one owns that axis.
func tldSwap(suffixes []string) Generate {
	return func(name string) []string {
		prefix, core, _ := SplitDomain(name)
		if core == "" {
			return nil
		}
		out := make([]string, 0, len(suffixes))
		for _, s := range suffixes {
			out = append(out, joinDomain(prefix, core, s))
		}
		return out
	}
}
