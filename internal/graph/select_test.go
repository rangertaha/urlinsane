// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

import (
	"context"
	"strings"
	"testing"
)

type selOp struct{ id string }

func (o selOp) Id() string                                  { return o.id }
func (o selOp) Version() int                                { return 1 }
func (o selOp) Trigger() Trigger                            { return Trigger{} }
func (o selOp) Emits() Effects                              { return Effects{} }
func (o selOp) Resource() string                            { return "" }
func (o selOp) Exec(context.Context, View) (Delta, Outcome) { return Delta{}, Outcome{} }

func ids(ops []Operator) string {
	out := make([]string, 0, len(ops))
	for _, o := range ops {
		out = append(out, o.Id())
	}
	return strings.Join(out, ",")
}

func all() []Operator {
	return []Operator{selOp{"dns"}, selOp{"whois"}, selOp{"geo"}, selOp{"ptr"}}
}

func TestSelectOperatorsEmptySelectsEverything(t *testing.T) {
	for _, in := range [][]string{nil, {}, {""}, {" , "}} {
		got, err := SelectOperators(all(), in)
		if err != nil {
			t.Fatalf("%v: %v", in, err)
		}
		if ids(got) != "dns,whois,geo,ptr" {
			t.Errorf("%v selected %q, want everything", in, ids(got))
		}
	}
}

func TestSelectOperatorsKeepsListed(t *testing.T) {
	got, err := SelectOperators(all(), []string{"dns,ptr"})
	if err != nil {
		t.Fatal(err)
	}
	if ids(got) != "dns,ptr" {
		t.Errorf("got %q, want dns,ptr", ids(got))
	}
	// Repeatable and comma-split are the same thing.
	again, _ := SelectOperators(all(), []string{"dns", "ptr"})
	if ids(again) != ids(got) {
		t.Errorf("repeated %q differs from comma-split %q", ids(again), ids(got))
	}
}

func TestSelectOperatorsExcludes(t *testing.T) {
	got, err := SelectOperators(all(), []string{"^whois"})
	if err != nil {
		t.Fatal(err)
	}
	if ids(got) != "dns,geo,ptr" {
		t.Errorf("got %q, want everything but whois", ids(got))
	}
}

// Order is the input order, not the selection order: the plan is deterministic.
func TestSelectOperatorsPreservesOrder(t *testing.T) {
	got, _ := SelectOperators(all(), []string{"ptr,dns"})
	if ids(got) != "dns,ptr" {
		t.Errorf("got %q, want input order dns,ptr", ids(got))
	}
}

// An unknown id must not silently select nothing.
func TestSelectOperatorsRejectsUnknown(t *testing.T) {
	_, err := SelectOperators(all(), []string{"dns", "nosuch"})
	if err == nil {
		t.Fatal("accepted an unknown id")
	}
	if !strings.Contains(err.Error(), "nosuch") || !strings.Contains(err.Error(), "known:") {
		t.Errorf("error should name the bad id and the known set: %v", err)
	}
	if _, err := SelectOperators(all(), []string{"^nosuch"}); err == nil {
		t.Error("accepted an unknown exclusion")
	}
}

func TestSelectOperatorsRejectsMixedForms(t *testing.T) {
	if _, err := SelectOperators(all(), []string{"dns", "^whois"}); err == nil {
		t.Fatal("accepted a mix of kept and excluded ids")
	}
	if _, err := SelectOperators(all(), []string{"^"}); err == nil {
		t.Error("accepted a bare exclusion prefix")
	}
}
