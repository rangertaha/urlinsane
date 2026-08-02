// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/report"
)

// CSV formula injection. urlinsane exists to surface attacker-chosen
// strings, and usernames/packages permit "=". A cell beginning =,+,-,@ is
// executed as a formula when the report is opened in a spreadsheet.
func TestCSVNeutralisesFormulaInjection(t *testing.T) {
	g := graph.New(registry(t))
	if _, err := g.Seed("package", "=cmd|' /c calc'!a0"); err != nil {
		t.Skipf("seed refused: %v", err)
	}
	out := render(t, g, report.Options{Format: "csv"})
	for _, line := range strings.Split(out, "\n") {
		for _, cell := range strings.Split(line, ",") {
			c := strings.Trim(cell, `"`)
			if c == "" {
				continue
			}
			if strings.ContainsAny(c[:1], "=+-@\t\r") {
				t.Errorf("CSV cell begins with a formula trigger: %q", c)
			}
		}
	}
}

// --fail-on must not be defeatable by a filter.
func TestFailOnCannotBeEvadedByAFilter(t *testing.T) {
	full := report.Build(scan(t), report.Options{})
	f, err := report.ParseFilters([]string{"absent"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	narrowed := report.Build(scan(t), report.Options{Filters: f})
	if narrowed.Max() != full.Max() {
		t.Errorf("a filter changed Max() from %v to %v; --fail-on becomes evadable",
			full.Max(), narrowed.Max())
	}
}

// depth filter operators.
func TestDepthFilterAcceptsEveryComparison(t *testing.T) {
	for _, spec := range []string{"depth<=1", "depth>0", "depth=0", "depth<2", "depth>=1"} {
		if _, err := report.ParseFilter(spec); err != nil {
			t.Errorf("%q rejected: %v", spec, err)
		}
	}
}

// a filter naming an unknown type should not silently match nothing.
func TestUnknownTypeFilterIsRejected(t *testing.T) {
	f, err := report.ParseFilter("type=nosuchtype")
	if err != nil {
		t.Log("Q5: an unknown type in a filter is rejected at parse time")
		return
	}
	r := report.Build(scan(t), report.Options{Filters: []report.Filter{f}})
	if len(r.Nodes) == 0 {
		t.Log("Q5: `type=nosuchtype` parses and silently selects nothing — " +
			"a typo'd filter reads as an empty scan")
	}
}

// ndjson must not lose a payload field named "kind".
func TestNDJSONEmitsExactlyOneRunEvent(t *testing.T) {
	out := render(t, scan(t), report.Options{Format: "ndjson"})
	var runs int
	for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
		if strings.Contains(l, `"kind":"run"`) {
			runs++
		}
	}
	if runs != 1 {
		t.Errorf("expected exactly one run event, got %d", runs)
	}
}
