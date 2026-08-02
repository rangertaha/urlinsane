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

package graph

import (
	"bytes"
	"strings"
	"testing"
)

func planOps() []Operator {
	dns := &fakeOp{
		id: "dns", ver: 1, res: "dns",
		trig:  Trigger{On: Selector{Types: []string{"domain"}}, Where: []Condition{HasProp("punycode")}},
		emits: Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}},
	}
	geo := &fakeOp{
		id: "geo", ver: 1, res: "http",
		trig:  Trigger{On: Selector{Types: []string{"ip"}}},
		emits: Effects{Props: []string{"asn"}},
	}
	ptr := &fakeOp{
		id: "ptr", ver: 1, res: "dns",
		trig:  Trigger{On: Selector{Types: []string{"ip"}}},
		emits: Effects{Nodes: []string{"domain"}, Rels: []string{"PTR_TO"}},
	}
	idn := &fakeOp{
		id: "idn", ver: 1,
		trig:  Trigger{On: Selector{Types: []string{"domain"}}},
		emits: Effects{Props: []string{"punycode"}},
	}
	return []Operator{dns, geo, ptr, idn}
}

func planInput() PlanInput {
	return PlanInput{Seed: SeedSpec{Type: "domain", Key: "example.com"}}
}

func TestPlanCompiles(t *testing.T) {
	p, err := Compile(testRegistry(t), planOps(), planInput())
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	if len(p.Operators) != 4 {
		t.Fatalf("operators = %d, want 4", len(p.Operators))
	}
	if p.Hash == "" {
		t.Fatal("plan has no hash")
	}
}

func TestPlanHashIsStableAndSensitive(t *testing.T) {
	reg := testRegistry(t)
	a, _ := Compile(reg, planOps(), planInput())
	b, _ := Compile(reg, planOps(), planInput())
	if a.Hash != b.Hash {
		t.Fatal("identical inputs produced different plan hashes")
	}

	// A retrained plugin model must change the hash: a variant operator whose
	// model changed emits different variants, so the graph itself would differ.
	in := planInput()
	in.Models = map[string][]string{"dns": {"bafyNEW"}}
	c, _ := Compile(reg, planOps(), in)
	if c.Hash == a.Hash {
		t.Fatal("a plugin model CID did not change the plan hash")
	}

	// So must a changed flag value: a pinned plan reproduces configuration,
	// not merely operator selection.
	in = planInput()
	in.Config = map[string]map[string]string{"dns": {"nameservers": "9.9.9.9"}}
	d, _ := Compile(reg, planOps(), in)
	if d.Hash == a.Hash {
		t.Fatal("an operator config value did not change the plan hash")
	}

	// And the engine model.
	in = planInput()
	in.Model = "bafyENGINE"
	e, _ := Compile(reg, planOps(), in)
	if e.Hash == a.Hash {
		t.Fatal("the engine model CID did not change the plan hash")
	}
}

func TestPlanPrunesUnreachableOperators(t *testing.T) {
	reg := testRegistry(t)
	// Nothing in this plan produces a tld, so an operator bound to tld can
	// never fire and should not be listed as active.
	ops := append(planOps(), &fakeOp{
		id: "tldstat", ver: 1,
		trig:  Trigger{On: Selector{Types: []string{"tld"}}},
		emits: Effects{Props: []string{"rank"}},
	})
	p, err := Compile(reg, ops, planInput())
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	var pruned bool
	for _, id := range p.Pruned {
		if id == "tldstat" {
			pruned = true
		}
	}
	if !pruned {
		t.Fatalf("pruned = %v, want tldstat — nothing emits a tld from a domain seed", p.Pruned)
	}
}

func TestPlanDetectsDeadOperator(t *testing.T) {
	reg := testRegistry(t)
	// Drop idn, the only producer of punycode. dns requires it, so dns can
	// never fire — and the plan must say so rather than listing it as active.
	var ops []Operator
	for _, o := range planOps() {
		if o.Id() != "idn" {
			ops = append(ops, o)
		}
	}
	p, err := Compile(reg, ops, planInput())
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	var dead *OpBinding
	for i := range p.Operators {
		if p.Operators[i].Id == "dns" {
			dead = &p.Operators[i]
		}
	}
	if dead == nil || !dead.Dead {
		t.Fatal("dns requires punycode which nothing emits, but was not marked dead")
	}
	if !strings.Contains(dead.DeadWhy, "punycode") {
		t.Fatalf("DeadWhy = %q, want it to name the unsatisfiable prop", dead.DeadWhy)
	}
}

func TestDeadOperatorsAreNotSelected(t *testing.T) {
	reg := testRegistry(t)
	var ops []Operator
	for _, o := range planOps() {
		if o.Id() != "idn" {
			ops = append(ops, o)
		}
	}
	p, _ := Compile(reg, ops, planInput())
	for _, o := range p.Select(ops) {
		if o.Id() == "dns" {
			t.Fatal("a dead operator was selected for execution")
		}
	}
}

func TestPlanRoundTripsAndVerifies(t *testing.T) {
	p, _ := Compile(testRegistry(t), planOps(), planInput())
	var buf bytes.Buffer
	if err := p.Write(&buf); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := ReadPlan(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.Hash != p.Hash {
		t.Fatalf("round-tripped hash %s != %s", got.Hash, p.Hash)
	}
}

func TestTamperedPlanIsRejected(t *testing.T) {
	// A pinned plan is a reproducibility claim. If the contents no longer hash
	// to the recorded value, the claim is false and loading must fail loudly.
	p, _ := Compile(testRegistry(t), planOps(), planInput())
	var buf bytes.Buffer
	_ = p.Write(&buf)
	tampered := strings.Replace(buf.String(), `"version": 1`, `"version": 7`, 1)
	if _, err := ReadPlan(strings.NewReader(tampered)); err == nil {
		t.Fatal("a tampered plan was accepted")
	}
}

func TestExplainRendersCyclesAsLayers(t *testing.T) {
	// domain -> ip -> domain is a genuine cycle. It must render, not hang or
	// pretend to be a DAG.
	p, _ := Compile(testRegistry(t), planOps(), planInput())
	var buf bytes.Buffer
	if err := p.Explain(&buf); err != nil {
		t.Fatalf("explain: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"plan ", "seed domain example.com", "type flow", "operators", "(cycle)"} {
		if !strings.Contains(out, want) {
			t.Fatalf("explain output missing %q:\n%s", want, out)
		}
	}
}

func TestExplainListsDeadAndPruned(t *testing.T) {
	reg := testRegistry(t)
	ops := append(planOps(), &fakeOp{
		id: "tldstat", ver: 1,
		trig: Trigger{On: Selector{Types: []string{"tld"}}},
	})
	var kept []Operator
	for _, o := range ops {
		if o.Id() != "idn" {
			kept = append(kept, o)
		}
	}
	p, _ := Compile(reg, kept, planInput())
	var buf bytes.Buffer
	_ = p.Explain(&buf)
	out := buf.String()
	if !strings.Contains(out, "dead") || !strings.Contains(out, "pruned") {
		t.Fatalf("explain must surface dead and pruned operators:\n%s", out)
	}
}

func TestCompileRejectsBadScope(t *testing.T) {
	reg := testRegistry(t)
	in := planInput()
	in.Seed.Scope = []string{"ip"} // observed, not nameable
	if _, err := Compile(reg, planOps(), in); err == nil {
		t.Fatal("accepted an observed type as a variant scope")
	}
	in.Seed.Scope = []string{"nosuch"}
	if _, err := Compile(reg, planOps(), in); err == nil {
		t.Fatal("accepted an unregistered scope type")
	}
}

func TestPlanSelectDrivesTheScheduler(t *testing.T) {
	// Execution dispatches only operators present in the plan; runtime never
	// reaches past it.
	reg := testRegistry(t)
	g := New(reg)
	if _, err := g.Seed("domain", "example.com"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	live := &fakeOp{
		id: "live", ver: 1, trig: onDomain(Reads{}),
		fn: func(View) (Delta, Outcome) { return Delta{}, OK() },
	}
	orphan := &fakeOp{
		id: "orphan", ver: 1,
		trig: Trigger{On: Selector{Types: []string{"tld"}}},
		fn:   func(View) (Delta, Outcome) { return Delta{}, OK() },
	}
	ops := []Operator{live, orphan}
	p, err := Compile(reg, ops, planInput())
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	run(t, g, p.Select(ops), Limits{})
	if live.calls == 0 {
		t.Fatal("planned operator never ran")
	}
	if orphan.calls != 0 {
		t.Fatal("a pruned operator was dispatched")
	}
}
