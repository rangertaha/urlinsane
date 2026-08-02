// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package idn

import (
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/observetest"
)

// TestIDNRecordsTheUnicodeForm: the key is punycode by canonicalization, so the
// useful assertion is the readable form. A homograph attack is only legible in
// Unicode, and xn--80ak6aa92e.com tells a reader nothing.
func TestIDNRecordsTheUnicodeForm(t *testing.T) {
	g := observetest.Run(t, observe.TypeDomain, "xn--80ak6aa92e.com", newIDN(observe.Options{}))

	v, ok := observetest.Prop(t, g, observe.TypeDomain, "xn--80ak6aa92e.com", observe.FieldUnicode)
	if !ok {
		t.Fatal("unicode observetest.Prop not set for an internationalized domain")
	}
	if v.Str() == "xn--80ak6aa92e.com" || v.Str() == "" {
		t.Errorf("unicode = %q, want the decoded form", v.Str())
	}
	observetest.WantStatus(t, g, observe.TypeDomain, "xn--80ak6aa92e.com", "idn", graph.StatusOK)
}

// TestIDNPlainDomainIsEmpty: an ASCII domain has no Unicode form. Nothing was
// learned, and that absence is recorded rather than reported as a failure.
func TestIDNPlainDomainIsEmpty(t *testing.T) {
	g := observetest.Run(t, observe.TypeDomain, "example.com", newIDN(observe.Options{}))
	if _, ok := observetest.Prop(t, g, observe.TypeDomain, "example.com", observe.FieldUnicode); ok {
		t.Error("an ASCII domain was given a unicode observetest.Prop")
	}
	observetest.WantStatus(t, g, observe.TypeDomain, "example.com", "idn", graph.StatusEmpty)
}
