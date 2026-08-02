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

package observe

import (
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// TestIDNRecordsTheUnicodeForm: the key is punycode by canonicalization, so the
// useful assertion is the readable form. A homograph attack is only legible in
// Unicode, and xn--80ak6aa92e.com tells a reader nothing.
func TestIDNRecordsTheUnicodeForm(t *testing.T) {
	g := run(t, TypeDomain, "xn--80ak6aa92e.com", newIDN(Options{}))

	v, ok := prop(t, g, TypeDomain, "xn--80ak6aa92e.com", FieldUnicode)
	if !ok {
		t.Fatal("unicode prop not set for an internationalized domain")
	}
	if v.Str() == "xn--80ak6aa92e.com" || v.Str() == "" {
		t.Errorf("unicode = %q, want the decoded form", v.Str())
	}
	wantStatus(t, g, TypeDomain, "xn--80ak6aa92e.com", "idn", graph.StatusOK)
}

// TestIDNPlainDomainIsEmpty: an ASCII domain has no Unicode form. Nothing was
// learned, and that absence is recorded rather than reported as a failure.
func TestIDNPlainDomainIsEmpty(t *testing.T) {
	g := run(t, TypeDomain, "example.com", newIDN(Options{}))
	if _, ok := prop(t, g, TypeDomain, "example.com", FieldUnicode); ok {
		t.Error("an ASCII domain was given a unicode prop")
	}
	wantStatus(t, g, TypeDomain, "example.com", "idn", graph.StatusEmpty)
}
