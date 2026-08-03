// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package report

import (
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
)

func mustFilters(t *testing.T, specs ...string) []Filter {
	t.Helper()
	out, err := ParseFilters(specs)
	if err != nil {
		t.Fatalf("ParseFilters(%v): %v", specs, err)
	}
	return out
}

// Two existence filters are alternatives. keep worked this out by comparing the
// filter's spec to lowercase literals, but spec holds the user's casing, so
// "Live" matched none of them and fell through to the conjunction branch --
// asking for nodes that were both live and absent. The report came back empty
// and the exit code said success.
func TestExistenceFiltersAreAlternativesWhateverTheCasing(t *testing.T) {
	live := NodeRow{Key: "a", existence: graph.Live}
	absent := NodeRow{Key: "b", existence: graph.Absent}
	unknown := NodeRow{Key: "c", existence: graph.Unknown}

	for _, specs := range [][]string{
		{"live", "absent"},
		{"Live", "Absent"},
		{"LIVE", "absent"},
	} {
		filters := mustFilters(t, specs...)
		for _, n := range []NodeRow{live, absent} {
			if !keep(filters, n) {
				t.Errorf("%v dropped %q, want it kept: the two filters are alternatives",
					specs, n.Key)
			}
		}
		if keep(filters, unknown) {
			t.Errorf("%v kept %q, which is neither", specs, unknown.Key)
		}
	}
}

// Across families filters still narrow, which is the behaviour the existence
// special case exists to carve out of.
func TestFiltersFromDifferentFamiliesNarrow(t *testing.T) {
	filters := mustFilters(t, "Live", "risk>=high")

	hot := NodeRow{Key: "hot", existence: graph.Live, severity: graph.SeverityHigh}
	cold := NodeRow{Key: "cold", existence: graph.Live, severity: graph.SeverityLow}

	if !keep(filters, hot) {
		t.Error("a live, high-risk node was dropped")
	}
	if keep(filters, cold) {
		t.Error("a live, low-risk node was kept; the risk filter must still narrow")
	}
}

// A single existence filter is still a filter, not a no-op.
func TestASingleExistenceFilterStillExcludes(t *testing.T) {
	filters := mustFilters(t, "Absent")

	if keep(filters, NodeRow{Key: "a", existence: graph.Live}) {
		t.Error("--filter Absent kept a live node")
	}
	if !keep(filters, NodeRow{Key: "b", existence: graph.Absent}) {
		t.Error("--filter Absent dropped an absent node")
	}
}

// String is the user's text, so it can be quoted back to them unchanged.
func TestFilterStringKeepsTheUsersSpelling(t *testing.T) {
	f, err := ParseFilter("Live")
	if err != nil {
		t.Fatal(err)
	}
	if f.String() != "Live" {
		t.Errorf("String() = %q, want the spec as written", f.String())
	}
}

// Two filters on an attribute a node has exactly one of are alternatives, in
// EVERY family — not just existence.
//
// This is the class test for that rule. The behaviour was special-cased for
// existence, which fixed `--filter Live,Absent` and left type, risk= and depth=
// broken in exactly the same way: `--filter type=domain,type=ip` matched no node
// at all, rendering an empty report with exit 0 — a scan that reads as "found
// nothing" rather than a query that cannot match. A new single-valued filter
// declares its family where it is parsed and is covered here.
func TestSameFamilyFiltersAreAlternatives(t *testing.T) {
	domain := NodeRow{Key: "d", Type: "domain", Depth: 0,
		existence: graph.Live, severity: graph.SeverityHigh}
	ip := NodeRow{Key: "i", Type: "ip", Depth: 1,
		existence: graph.Absent, severity: graph.SeverityCritical}

	for _, tc := range []struct {
		specs []string
		why   string
	}{
		{[]string{"live", "absent"}, "a node has one existence"},
		{[]string{"type=domain", "type=ip"}, "a node has one type"},
		{[]string{"risk=high", "risk=critical"}, "a node has one severity"},
		{[]string{"depth=0", "depth=1"}, "a node has one depth"},
	} {
		filters := mustFilters(t, tc.specs...)
		for _, n := range []NodeRow{domain, ip} {
			if !keep(filters, n) {
				t.Errorf("%v dropped %q: %s, so the two filters must be alternatives",
					tc.specs, n.Key, tc.why)
			}
		}
	}
}

// A comparison is not mutually exclusive, so two of them still narrow — that is
// what the user asked for and the family rule must not swallow it.
func TestComparisonFiltersStillNarrow(t *testing.T) {
	filters := mustFilters(t, "risk>=high", "depth<=0")

	hot := NodeRow{Key: "hot", Depth: 0, severity: graph.SeverityCritical}
	deep := NodeRow{Key: "deep", Depth: 5, severity: graph.SeverityCritical}

	if !keep(filters, hot) {
		t.Error("a shallow, critical node was dropped")
	}
	if keep(filters, deep) {
		t.Error("a deep node was kept; two comparisons must still both apply")
	}
}

// Across families they narrow, which is the behaviour the family rule carves
// the alternatives out of.
func TestDifferentFamiliesStillNarrowAcross(t *testing.T) {
	filters := mustFilters(t, "type=domain", "type=ip", "live")

	if !keep(filters, NodeRow{Key: "a", Type: "domain", existence: graph.Live}) {
		t.Error("a live domain was dropped")
	}
	if keep(filters, NodeRow{Key: "b", Type: "ip", existence: graph.Absent}) {
		t.Error("an absent ip was kept; the existence family must still narrow across")
	}
}
