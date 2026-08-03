// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

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
// Every type whose canonical key carries a qualifier keeps that qualifier out
// of the core. A package key is "npm:lodash" and a repo key is
// "github.com/acme/tool", so varying the whole key does not vary the name — it
// varies the registry and the forge. Character omission over "npm:lodash"
// produced "pm:lodash": a package in a registry that does not exist, generated
// by every Whole:false algorithm, for every package and repo the tool was
// pointed at.
//
// What is left in the core is still varied whole. "acme/tool" is one name and
// splitting it would exclude exactly the org-name confusions that make repo
// squatting work, and a bare username has no qualifier to strip.
//
// The qualifiers dropped here are not lost. A repo decomposes into
// platform:github.com and username:acme, and a scoped package into its owner,
// each reached structurally as its own node and varied there by these same
// operators — the same division of labour that lets the email case vary only
// the local part.
func DefaultSplit(nodeType, key string) Parts {
	switch nodeType {
	case TypeDomain:
		prefix, core, suffix := SplitDomain(key)
		if core == "" {
			return wholeKey(key)
		}
		return Parts{Prefix: prefix, Core: core, Suffix: suffix, join: JoinDomain}
	case TypeEmail:
		// Vary the local part only. The domain half of an email is reached
		// structurally as its own domain node and is varied there, by the same
		// operators, without this one having to know how.
		at := strings.LastIndex(key, "@")
		if at <= 0 || at == len(key)-1 {
			return wholeKey(key)
		}
		return Parts{Prefix: "", Core: key[:at], Suffix: key[at+1:], join: joinEmail}
	case TypePackage:
		// "registry:name" — canonPackage requires the qualifier, so a key
		// without one did not come from the decomposer and is left whole
		// rather than guessed at.
		registry, name, ok := strings.Cut(key, ":")
		if !ok || registry == "" || name == "" {
			return wholeKey(key)
		}
		return Parts{Prefix: registry + ":", Core: name, join: joinQualified}
	case TypeRepo:
		// "host/owner/name" — the host is the forge, and a variant of it is a
		// domain squat that the platform node covers.
		host, path, ok := strings.Cut(key, "/")
		if !ok || host == "" || path == "" {
			return wholeKey(key)
		}
		return Parts{Prefix: host + "/", Core: path, join: joinQualified}
	}
	return wholeKey(key)
}

func wholeKey(key string) Parts {
	return Parts{Core: key, join: func(_, core, _ string) string { return core }}
}

// joinDomain reassembles a domain, dropping empty labels. Algorithms legally
// produce cores containing dots (dot insertion) and empty strings (omission of
// a one-character name), and a naive join would leave "..com".
func JoinDomain(prefix, core, suffix string) string {
	labels := make([]string, 0, 3)
	for _, p := range []string{prefix, core, suffix} {
		p = strings.Trim(strings.TrimSpace(p), ".")
		if p != "" {
			labels = append(labels, p)
		}
	}
	return strings.Join(labels, ".")
}

// joinQualified restores a qualifier the algorithm was not allowed to see.
//
// An empty core yields an empty key rather than a bare "npm:", matching what
// joinEmail does with an emptied local part: omission of a one-character name
// has no result, and "npm:" is not a package.
func joinQualified(prefix, core, _ string) string {
	if core == "" {
		return ""
	}
	return prefix + core
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
