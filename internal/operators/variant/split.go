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
	"sync"

	"github.com/rangertaha/urlinsane/datasets"
)

// Parts is a key decomposed into the substring a variant algorithm may vary and
// the surrounding text it must preserve.
//
// Varying only the core is not cosmetic. Character omission over the whole of
// "example.com" produces "example.cm", which is a *different registry* rather
// than a typo of the name — the kind of variant the TLD algorithm exists to
// generate deliberately, and which every other algorithm would otherwise
// generate by accident.
type Parts struct {
	// Prefix is everything left of the core (domain subdomains).
	Prefix string
	// Core is the part an algorithm may vary.
	Core string
	// Suffix is everything right of the core (a domain's public suffix).
	Suffix string

	join func(prefix, core, suffix string) string
}

// Join puts a varied core back together with the preserved parts.
func (p Parts) Join(core string) string {
	if p.join == nil {
		return core
	}
	return p.join(p.Prefix, core, p.Suffix)
}

// Splitter decomposes a node's key for the node's type. It is injectable so
// that a test — or a future type with its own grammar — can supply its own
// decomposition without every algorithm having to know about it.
type Splitter func(nodeType, key string) Parts

// DefaultSplit is the decomposition the shipped operators use.
//
// Only domains and emails have structure worth preserving. A username, package
// or repo name is varied whole: "acme/tool" is one name, and splitting it would
// exclude exactly the org-name confusions that make repo squatting work.
func DefaultSplit(nodeType, key string) Parts {
	switch nodeType {
	case TypeDomain:
		prefix, core, suffix := SplitDomain(key)
		if core == "" {
			return wholeKey(key)
		}
		return Parts{Prefix: prefix, Core: core, Suffix: suffix, join: joinDomain}
	case TypeEmail:
		// Vary the local part only. The domain half of an email is reached
		// structurally as its own domain node and is varied there, by the same
		// operators, without this one having to know how.
		at := strings.LastIndex(key, "@")
		if at <= 0 || at == len(key)-1 {
			return wholeKey(key)
		}
		return Parts{Prefix: "", Core: key[:at], Suffix: key[at+1:], join: joinEmail}
	}
	return wholeKey(key)
}

func wholeKey(key string) Parts {
	return Parts{Core: key, join: func(_, core, _ string) string { return core }}
}

// joinDomain reassembles a domain, dropping empty labels. Algorithms legally
// produce cores containing dots (dot insertion) and empty strings (omission of
// a one-character name), and a naive join would leave "..com".
func joinDomain(prefix, core, suffix string) string {
	labels := make([]string, 0, 3)
	for _, p := range []string{prefix, core, suffix} {
		p = strings.Trim(strings.TrimSpace(p), ".")
		if p != "" {
			labels = append(labels, p)
		}
	}
	return strings.Join(labels, ".")
}

func joinEmail(_, core, suffix string) string {
	if core == "" || suffix == "" {
		return ""
	}
	return core + "@" + suffix
}

// SplitDomain splits a domain into its subdomain prefix, registrable name and
// public suffix, using the compiled-in public suffix list.
//
// It deliberately does not go through internal/pkg/dns, whose parser reads the
// suffix list out of the dataset database: an operator whose output depends on
// whether a database happens to be open is not the pure function the scheduler
// caches it as. datasets.TLD is compiled in and always present.
func SplitDomain(domain string) (prefix, name, suffix string) {
	domain = strings.Trim(strings.TrimSpace(domain), ".")
	if domain == "" {
		return "", "", ""
	}
	labels := strings.Split(domain, ".")
	n := len(labels)

	sfx := publicSuffixLabels(labels)
	// The whole name is a public suffix ("com", "co.uk"): there is no
	// registrable label to vary, so the caller falls back to the whole key.
	if sfx >= n {
		return "", "", ""
	}
	return strings.Join(labels[:n-sfx-1], "."), labels[n-sfx-1], strings.Join(labels[n-sfx:], ".")
}

// publicSuffixLabels returns how many trailing labels form the public suffix,
// following the PSL matching rules: longest match wins, an exception rule beats
// everything, and an unlisted TLD falls back to the implicit "*" rule.
func publicSuffixLabels(labels []string) int {
	normal, wildcard, exception := psl()
	n := len(labels)

	// Exceptions take priority over every other rule, so they are checked
	// first and return outright.
	for i := 0; i < n; i++ {
		if exception[strings.Join(labels[i:], ".")] {
			return n - i - 1
		}
	}

	best := 1 // the implicit "*" rule: an unknown TLD is one label
	for i := 0; i < n; i++ {
		s := strings.Join(labels[i:], ".")
		if normal[s] && n-i > best {
			best = n - i
		}
		// A wildcard rule "*.foo" is stored under "foo" and consumes one more
		// label than its parent, so it needs a label to its left to match.
		if i > 0 && wildcard[s] && n-i+1 > best {
			best = n - i + 1
		}
	}
	if best > n {
		best = n
	}
	return best
}

var (
	pslOnce      sync.Once
	pslNormal    map[string]bool
	pslWildcard  map[string]bool
	pslException map[string]bool
)

// psl builds the three rule sets once. Wildcards and exceptions are kept apart
// from normal rules rather than stripped of their markers: "*.ck" reduced to
// "ck" would be indistinguishable from a plain "ck" rule and would match one
// label short.
func psl() (normal, wildcard, exception map[string]bool) {
	pslOnce.Do(func() {
		pslNormal = make(map[string]bool, len(datasets.TLD))
		pslWildcard = map[string]bool{}
		pslException = map[string]bool{}
		for _, rule := range datasets.TLD {
			rule = strings.TrimSpace(rule)
			switch {
			case rule == "":
			case strings.HasPrefix(rule, "!"):
				pslException[rule[1:]] = true
			case strings.HasPrefix(rule, "*."):
				pslWildcard[rule[2:]] = true
			default:
				pslNormal[rule] = true
			}
		}
	})
	return pslNormal, pslWildcard, pslException
}
