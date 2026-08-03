// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/observetest"
)

// --- registry ---------------------------------------------------------------

// --- package-level wiring ---------------------------------------------------

func TestNewOmitsOperatorsWithoutTheirDependency(t *testing.T) {
	// Geo has no locator here, and the source operators have no source list, so
	// neither may appear: an operator that can only fail would make --explain
	// claim work that will never happen.
	ops := New(observe.Options{Resolver: &observetest.FakeResolver{}, Whois: observetest.FakeWhois{}})
	for _, op := range ops {
		switch op.Id() {
		case "geo", "pkg", "usr", "repo":
			t.Fatalf("operator %q planned without its dependency", op.Id())
		}
	}

	// A lister that lists nothing is not a dependency met. It was accepted as
	// one, so a dataset whose source tables were empty planned all three
	// operators and every call returned "no sources configured" -- the same
	// --explain lie, one layer down.
	ops = New(observe.Options{
		Resolver: &observetest.FakeResolver{},
		Whois:    observetest.FakeWhois{},
		Sources:  fakeSources{},
		Prober:   &fakeProber{},
	})
	for _, op := range ops {
		switch op.Id() {
		case "pkg", "usr", "repo":
			t.Errorf("operator %q planned against a lister with no sources", op.Id())
		}
	}

	ops = New(observe.Options{
		Resolver: &observetest.FakeResolver{},
		Whois:    observetest.FakeWhois{},
		Geo:      observetest.FakeGeo{},
		Sources:  registries(),
		Prober:   &fakeProber{},
	})
	ids := map[string]bool{}
	for _, op := range ops {
		ids[op.Id()] = true
	}
	for _, want := range []string{
		"idn", "dns-a", "dns-ns", "dns-mx", "dns-txt", "dns-cname",
		"ptr", "geo", "whois", "pkg", "usr", "repo",
	} {
		if !ids[want] {
			t.Errorf("operator %q missing from New", want)
		}
	}
}

// TestTriggersBindToTypesNotProducers is the central invariant of DESIGN §4.1.
// The old collectors carried a Dependencies() list; nothing here may reintroduce
// one, and every operator must bind to the type of the data it consumes rather
// than to whatever happened to produce it.
func TestTriggersBindToTypesNotProducers(t *testing.T) {
	ops := New(observe.Options{
		Resolver: &observetest.FakeResolver{},
		Whois:    observetest.FakeWhois{},
		Geo:      observetest.FakeGeo{},
		Sources:  fakeSources{},
		Prober:   &fakeProber{},
	})
	want := map[string]string{
		"idn": observe.TypeDomain, "dns-a": observe.TypeDomain, "dns-ns": observe.TypeDomain,
		"dns-mx": observe.TypeDomain, "dns-txt": observe.TypeDomain, "dns-cname": observe.TypeDomain,
		"whois": observe.TypeDomain,
		// The two that would have been producer dependencies: geo and ptr
		// consume addresses, so they bind to ip.
		"geo": observe.TypeIP, "ptr": observe.TypeIP,
		"pkg": observe.TypePackage, "usr": observe.TypeUsername, "repo": observe.TypeRepo,
	}
	for _, op := range ops {
		trigger := op.Trigger()
		if len(trigger.On.Types) != 1 || trigger.On.Types[0] != want[op.Id()] {
			t.Errorf("%s binds to %v, want [%s]", op.Id(), trigger.On.Types, want[op.Id()])
		}
		if len(op.Emits().Nodes)+len(op.Emits().Rels)+len(op.Emits().Props) == 0 {
			t.Errorf("%s declares no effects, so plan compilation cannot see it", op.Id())
		}
	}
}

// TestResourceClassesAreDeclared guards DESIGN §6.3: one token bucket per class,
// so the limit protecting whois does not throttle DNS to the same crawl.
func TestResourceClassesAreDeclared(t *testing.T) {
	ops := New(observe.Options{
		Resolver: &observetest.FakeResolver{}, Whois: observetest.FakeWhois{}, Geo: observetest.FakeGeo{},
		Sources: fakeSources{}, Prober: &fakeProber{},
	})
	want := map[string]string{
		"idn": "", "dns-a": observe.ResourceDNS, "dns-ns": observe.ResourceDNS, "dns-mx": observe.ResourceDNS,
		"dns-txt": observe.ResourceDNS, "dns-cname": observe.ResourceDNS, "ptr": observe.ResourceDNS,
		"geo": observe.ResourceGeo, "whois": observe.ResourceWhois,
		"pkg": observe.ResourceHTTP, "usr": observe.ResourceHTTP, "repo": observe.ResourceHTTP,
	}
	for _, op := range ops {
		if got := op.Resource(); got != want[op.Id()] {
			t.Errorf("%s resource = %q, want %q", op.Id(), got, want[op.Id()])
		}
	}
}

// TestSchemaAppendsCleanly proves the field lists this package publishes
// register when appended to the registrar's own — the contract decompose has to
// honour — and that no name collides with one already declared there.
func TestSchemaAppendsCleanly(t *testing.T) {
	r := observetest.TestRegistry(t)

	// Every observetest.Prop any operator declares in Effects must resolve, or the operator
	// is emitting into a field that does not exist.
	ops := New(observe.Options{
		Resolver: &observetest.FakeResolver{}, Whois: observetest.FakeWhois{}, Geo: observetest.FakeGeo{},
		Sources: fakeSources{}, Prober: &fakeProber{},
	})
	nodeFields := map[string]bool{}
	for _, defs := range observe.Fields() {
		for _, f := range defs {
			nodeFields[f.Name] = true
		}
	}
	nodeFields[observe.FieldCreated] = true // declared by the registrar, asserted here
	relProps := map[string]bool{}
	for _, defs := range observe.RelFields() {
		for _, f := range defs {
			relProps[f.Name] = true
		}
	}
	for _, op := range ops {
		for _, name := range op.Emits().Props {
			if !nodeFields[name] && !relProps[name] {
				t.Errorf("%s declares observetest.Prop %q that no field list defines", op.Id(), name)
			}
		}
	}

	for name, defs := range observe.RelFields() {
		rel, ok := r.Rel(name)
		if !ok {
			t.Fatalf("relation %s did not register", name)
		}
		if rel.Class() != graph.Observation {
			t.Errorf("%s is class %s, want observation: every relation here cost a network call", name, rel.Class())
		}
		for _, f := range defs {
			if _, ok := rel.Field(f.Name); !ok {
				t.Errorf("relation %s is missing field %q", name, f.Name)
			}
		}
	}
}
