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

// Type classifies a named entity.
type Type string

const (
	Domain  Type = "domain"  // e.g. example.com
	Name    Type = "name"    // a person or brand name
	User    Type = "user"    // a username / handle
	Package Type = "package" // a package in a registry (PyPI, npm, ...)
)

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
