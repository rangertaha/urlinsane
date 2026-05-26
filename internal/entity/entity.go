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

// Package entity defines the kinds of named entities urlinsane analyzes.
// Typosquatting is not specific to domains: the same string-mutation engine
// applies to usernames, person names, and package names in public registries.
// Each algorithm and collector declares which entity Types it supports, and the
// engine runs only the plugins that apply to the target's type.
package entity

import (
	"regexp"
	"strings"
)

// Type classifies a named entity.
type Type string

const (
	Domain  Type = "domain"  // e.g. example.com
	Name    Type = "name"    // a person or brand name
	User    Type = "user"    // a username / handle
	Package Type = "package" // a package in a registry (PyPI, npm, ...)

	// Auto is not an entity type but a CLI sentinel: classify the target.
	Auto = "auto"
)

// domainRe matches a hostname with at least one dot and an alphabetic TLD.
var domainRe = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

// Classify infers the entity type of a target from its shape:
//   - contains whitespace        -> name   ("John Smith")
//   - starts with '@' or has '@' -> user   ("@handle", "a@b.com")
//   - looks like a hostname      -> domain ("example.com")
//   - otherwise (bare token)     -> package ("requests", "lodash")
//
// Bare identifiers are ambiguous (a username, a package, or a brand); they
// default to package given urlinsane's supply-chain focus. Pass --type to
// override.
func Classify(s string) Type {
	s = strings.TrimSpace(s)
	switch {
	case s == "":
		return Domain
	case strings.ContainsAny(s, " \t"):
		return Name
	case strings.HasPrefix(s, "@"), strings.Contains(s, "@"):
		return User
	case domainRe.MatchString(s):
		return Domain
	default:
		return Package
	}
}

// All returns every known entity type. A plugin that supports all types (a pure
// string mutation, or a registry-agnostic collector) reports this set.
func All() []Type {
	return []Type{Domain, Name, User, Package}
}

// Parse resolves a string to a Type, reporting whether it is known.
func Parse(s string) (Type, bool) {
	switch Type(s) {
	case Domain, Name, User, Package:
		return Type(s), true
	default:
		return "", false
	}
}

// Supports reports whether t is in the set types. An empty set is treated as
// "all types" so plugins that do not restrict themselves apply everywhere.
func Supports(types []Type, t Type) bool {
	if len(types) == 0 {
		return true
	}
	for _, x := range types {
		if x == t {
			return true
		}
	}
	return false
}
