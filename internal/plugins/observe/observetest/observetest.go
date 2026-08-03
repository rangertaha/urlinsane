// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package observetest holds the fakes and graph assertions the observation
// tests share.
//
// A package rather than a _test.go file because the operators now live in one
// package each, and a fake resolver copied into five of them would drift — the
// point of the fakes is that every operator is exercised against identical
// behaviour.
package observetest

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"golang.org/x/net/idna"
)

// TestRegistry performs the merge the real registrar has to perform:
// decompose's own field lists, plus this package's observe.Fields() and observe.RelFields()
// appended after them.
//
// It does not call decompose.Register, deliberately. That would test whether
// decompose has caught up with this package rather than whether this package is
// correct, and it would couple every test here to another package's schema
// version. What it does instead is prove the append works — which is the whole
// contract between the two.
func TestRegistry(t *testing.T) *graph.Registry {
	t.Helper()
	r := graph.NewRegistry()
	fields := observe.Fields()
	relFields := observe.RelFields()

	lower := func(s string) (string, error) {
		s = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(s, ".")))
		if s == "" {
			return "", fmt.Errorf("empty key")
		}
		return s, nil
	}
	domain := func(s string) (string, error) {
		s, err := lower(s)
		if err != nil {
			return "", err
		}
		return idna.Lookup.ToASCII(s)
	}
	address := func(s string) (string, error) {
		ip := net.ParseIP(strings.TrimSpace(s))
		if ip == nil {
			return "", fmt.Errorf("not an ip address: %q", s)
		}
		return ip.String(), nil
	}

	// The registrar's own domain fields, which this package's list appends to.
	// observe.FieldCreated lives here, not in observe.Fields(), because it is already declared.
	registrarDomainFields := []graph.FieldDef{
		{Name: "punycode", Kind: graph.KindString},
		{Name: "live", Kind: graph.KindBool},
		{Name: observe.FieldCreated, Kind: graph.KindTime, Merge: graph.Precedence("rdap", "whois")},
		{Name: "rank", Kind: graph.KindInt},
	}

	types := []graph.NodeTypeDef{
		{
			Name: observe.TypeDomain, Cap: graph.Nameable, Version: 1, Canonical: domain,
			Fields: append(registrarDomainFields, fields[observe.TypeDomain]...),
		},
		// Production's canonicalizers, not `lower`, for the types whose keys
		// carry a qualifier. See decompose.CanonicalFor: `lower` accepted
		// "lodash", which canonPackage rejects, so these tests seeded a key the
		// tool cannot produce and never exercised the operator on the
		// "npm:lodash" it really gets — which it was mangling into a URL for a
		// package that cannot exist.
		{Name: observe.TypeUsername, Cap: graph.Nameable, Version: 1,
			Canonical: decompose.CanonicalFor(observe.TypeUsername)},
		{Name: observe.TypePackage, Cap: graph.Nameable, Version: 1,
			Canonical: decompose.CanonicalFor(observe.TypePackage)},
		{Name: observe.TypeRepo, Cap: graph.Nameable, Version: 1,
			Canonical: decompose.CanonicalFor(observe.TypeRepo)},
		{Name: observe.TypeIP, Cap: graph.Observed, Version: 1, Canonical: address, Fields: fields[observe.TypeIP]},
		{Name: observe.TypeRegistrant, Cap: graph.Observed, Version: 1, Canonical: lower},
		{Name: observe.TypePlatform, Cap: graph.Observed, Version: 1, Canonical: lower, Fields: fields[observe.TypePlatform]},
	}
	for _, d := range types {
		if _, err := r.AddType(d); err != nil {
			t.Fatalf("register type %s: %v", d.Name, err)
		}
	}
	for _, name := range []string{observe.RelResolvesTo, observe.RelNS, observe.RelMX, observe.RelPTRTo, observe.RelRegisteredBy, observe.RelExistsOn} {
		d := graph.RelDef{Name: name, Class: graph.Observation, Version: 1, Fields: relFields[name]}
		if _, err := r.AddRel(d); err != nil {
			t.Fatalf("register relation %s: %v", name, err)
		}
	}
	return r
}

// --- harness ----------------------------------------------------------------

// Run seeds a graph and expands it with the given operators. Driving the real
// scheduler rather than calling Exec directly is deliberate: it proves the
// applier accepts every type, relation and field name the operators emit, which
// a delta compared against itself never would.
func Run(t *testing.T, seedType, seedKey string, ops ...graph.Operator) *graph.Graph {
	t.Helper()
	g := graph.New(TestRegistry(t))
	if _, err := g.Seed(seedType, seedKey); err != nil {
		t.Fatalf("seed: %v", err)
	}
	// Attempts 1 keeps a failing fake from being called twice; the retry path
	// belongs to the scheduler's own tests.
	s := graph.NewScheduler(g, ops, graph.Limits{Workers: 1, Attempts: 1, MaxRounds: 8})
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	return g
}

// NodeID resolves a (type, key) pair the way the applier would.
func NodeID(t *testing.T, g *graph.Graph, typ, key string) graph.NodeID {
	t.Helper()
	for _, n := range g.Nodes() {
		if n.Type.Name() == typ && n.Key == key {
			return n.ID
		}
	}
	t.Fatalf("no %s node with key %q; graph has %s", typ, key, Dump(g))
	return graph.NodeID{}
}

func HasNode(g *graph.Graph, typ, key string) bool {
	for _, n := range g.Nodes() {
		if n.Type.Name() == typ && n.Key == key {
			return true
		}
	}
	return false
}

// HasEdge reports whether an edge of rel joins the two keyed nodes.
func HasEdge(g *graph.Graph, fromType, fromKey, rel, toType, toKey string) bool {
	for _, e := range g.Edges() {
		if e.Rel.Name() != rel {
			continue
		}
		from, okf := g.Node(e.From)
		to, okt := g.Node(e.To)
		if !okf || !okt {
			continue
		}
		if from.Type.Name() == fromType && from.Key == fromKey &&
			to.Type.Name() == toType && to.Key == toKey {
			return true
		}
	}
	return false
}

// Prop reads a materialized Prop off a node.
func Prop(t *testing.T, g *graph.Graph, typ, key, field string) (graph.Value, bool) {
	t.Helper()
	n, ok := g.Node(NodeID(t, g, typ, key))
	if !ok {
		return graph.Value{}, false
	}
	f, ok := n.Type.Field(field)
	if !ok {
		t.Fatalf("type %s has no field %q", typ, field)
	}
	return n.Props.Get(f)
}

// EdgeProp reads a materialized Prop off the first matching edge.
func EdgeProp(t *testing.T, g *graph.Graph, rel, toKey, field string) (graph.Value, bool) {
	t.Helper()
	for _, e := range g.Edges() {
		if e.Rel.Name() != rel {
			continue
		}
		if to, ok := g.Node(e.To); !ok || to.Key != toKey {
			continue
		}
		f, ok := e.Rel.Field(field)
		if !ok {
			t.Fatalf("relation %s has no field %q", rel, field)
		}
		return e.Props.Get(f)
	}
	t.Fatalf("no %s edge to %q", rel, toKey)
	return graph.Value{}, false
}

// WantStatus asserts the terminal status the scheduler recorded for a pair.
// This is the assertion that matters most in this package: the ok/empty/failed
// /timeout split is the signal, not a detail of it.
func WantStatus(t *testing.T, g *graph.Graph, typ, key, op string, want graph.Status) {
	t.Helper()
	got, ok := g.Status(NodeID(t, g, typ, key), op)
	if !ok {
		t.Fatalf("%s on %s/%s recorded no status, want %s", op, typ, key, want)
	}
	if got != want {
		t.Fatalf("%s on %s/%s = %s, want %s", op, typ, key, got, want)
	}
}

func Dump(g *graph.Graph) string {
	var b strings.Builder
	for _, n := range g.Nodes() {
		fmt.Fprintf(&b, "%s(%s) ", n.Type.Name(), n.Key)
	}
	return b.String()
}

// --- stub view --------------------------------------------------------------

// StubView is the "literal view in, expected delta out" harness of DESIGN §13,
// for the cases where driving the scheduler would say less than calling Exec
// directly. No operator here declares Trigger.Reads, so an empty Prop and edge
// surface is a faithful view.
// StubView is a graph.View over nothing but a type and a key, for operators
// whose Exec reads no props.
type StubView struct{ NodeType, NodeKey string }

func (StubView) ID() graph.NodeID                { return graph.NodeID{} }
func (v StubView) Type() string                  { return v.NodeType }
func (v StubView) Key() string                   { return v.NodeKey }
func (StubView) Depth() int                      { return 0 }
func (StubView) Prop(string) (graph.Value, bool) { return graph.Value{}, false }
func (StubView) Edges(string) []graph.EdgeView   { return nil }
func (v StubView) Ref() graph.NodeRef            { return graph.NodeRef{Type: v.NodeType, Key: v.NodeKey} }

// --- fakes ------------------------------------------------------------------

// FakeResolver answers from a table. err, when set, is returned by every lookup
// so a test can pin one failure mode at a time.
type FakeResolver struct {
	IPs   map[string][]net.IP
	NS    map[string][]*net.NS
	MX    map[string][]*net.MX
	TXT   map[string][]string
	CNAME map[string]string
	Addr  map[string][]string
	Err   error
	Calls int
}

func (f *FakeResolver) LookupIP(_ context.Context, _, host string) ([]net.IP, error) {
	f.Calls++
	return f.IPs[host], f.Err
}

func (f *FakeResolver) LookupNS(_ context.Context, host string) ([]*net.NS, error) {
	f.Calls++
	return f.NS[host], f.Err
}

func (f *FakeResolver) LookupMX(_ context.Context, host string) ([]*net.MX, error) {
	f.Calls++
	return f.MX[host], f.Err
}

func (f *FakeResolver) LookupTXT(_ context.Context, host string) ([]string, error) {
	f.Calls++
	return f.TXT[host], f.Err
}

func (f *FakeResolver) LookupCNAME(_ context.Context, host string) (string, error) {
	f.Calls++
	return f.CNAME[host], f.Err
}

func (f *FakeResolver) LookupAddr(_ context.Context, addr string) ([]string, error) {
	f.Calls++
	return f.Addr[addr], f.Err
}

// Nxdomain is what the Go resolver reports for a name that does not exist. It
// is the single most important error in this package: it must never be Failed.
func Nxdomain() error {
	return &net.DNSError{Err: "no such host", Name: "example.com", IsNotFound: true}
}

func Servfail() error {
	return &net.DNSError{Err: "server misbehaving", Name: "example.com"}
}

// DNSTimeout is what the resolver reports when a lookup exceeds its deadline.
func DNSTimeout() error {
	return &net.DNSError{Err: "i/o timeout", Name: "example.com", IsTimeout: true}
}

// FakeWhois answers from a fixed record, so the whois wiring can be exercised
// without a registry.
type FakeWhois struct {
	Raw string
	Err error
}

func (f FakeWhois) Whois(_ string, _ ...string) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	if f.Raw != "" {
		return f.Raw, nil
	}
	return "Domain Name: example.com", nil
}

// FakeGeo answers from a table, so no MaxMind database is needed.
type FakeGeo struct {
	Places map[string]*observe.Location
	Err    error
}

func (f FakeGeo) Locate(addr string) (*observe.Location, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return f.Places[addr], nil
}
