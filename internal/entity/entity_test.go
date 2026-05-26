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

package entity

import "testing"

func TestParse(t *testing.T) {
	for _, in := range []string{"domain", "name", "user", "package"} {
		if got, ok := Parse(in); !ok || string(got) != in {
			t.Fatalf("Parse(%q) = %q,%v", in, got, ok)
		}
	}
	if _, ok := Parse("bogus"); ok {
		t.Fatal("Parse(bogus) should be unknown")
	}
	if _, ok := Parse(""); ok {
		t.Fatal("Parse(\"\") should be unknown")
	}
}

func TestSupports(t *testing.T) {
	// Empty set means "all types".
	if !Supports(nil, Domain) || !Supports(nil, Package) {
		t.Fatal("empty set should support every type")
	}
	// Restricted set.
	domainOnly := []Type{Domain}
	if !Supports(domainOnly, Domain) {
		t.Fatal("domainOnly should support Domain")
	}
	if Supports(domainOnly, User) {
		t.Fatal("domainOnly should not support User")
	}
}

func TestAll(t *testing.T) {
	if len(All()) != 4 {
		t.Fatalf("expected 4 entity types, got %d", len(All()))
	}
}
