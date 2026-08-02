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
	"context"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"golang.org/x/net/idna"
)

// --- registry ---------------------------------------------------------------

// testRegistry performs the merge the real registrar has to perform:
// decompose's own field lists, plus this package's Fields() and RelFields()
// appended after them.
//
// It does not call decompose.Register, deliberately. That would test whether
// decompose has caught up with this package rather than whether this package is
// correct, and it would couple every test here to another package's schema
// version. What it does instead is prove the append works — which is the whole
// contract between the two.
func testRegistry(t *testing.T) *graph.Registry {
	t.Helper()
	r := graph.NewRegistry()
	fields := Fields()
	relFields := RelFields()

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
	// FieldCreated lives here, not in Fields(), because it is already declared.
	registrarDomainFields := []graph.FieldDef{
		{Name: "punycode", Kind: graph.KindString},
		{Name: "live", Kind: graph.KindBool},
		{Name: FieldCreated, Kind: graph.KindTime, Merge: graph.Precedence("rdap", "whois")},
		{Name: "rank", Kind: graph.KindInt},
	}

	types := []graph.NodeTypeDef{
		{
			Name: TypeDomain, Cap: graph.Nameable, Version: 1, Canonical: domain,
			Fields: append(registrarDomainFields, fields[TypeDomain]...),
		},
		{Name: TypeUsername, Cap: graph.Nameable, Version: 1, Canonical: lower},
		{Name: TypePackage, Cap: graph.Nameable, Version: 1, Canonical: lower},
		{Name: TypeRepo, Cap: graph.Nameable, Version: 1, Canonical: lower},
		{Name: TypeIP, Cap: graph.Observed, Version: 1, Canonical: address, Fields: fields[TypeIP]},
		{Name: TypeRegistrant, Cap: graph.Observed, Version: 1, Canonical: lower},
		{Name: TypePlatform, Cap: graph.Observed, Version: 1, Canonical: lower, Fields: fields[TypePlatform]},
	}
	for _, d := range types {
		if _, err := r.AddType(d); err != nil {
			t.Fatalf("register type %s: %v", d.Name, err)
		}
	}
	for _, name := range []string{RelResolvesTo, RelNS, RelMX, RelPTRTo, RelRegisteredBy, RelExistsOn} {
		d := graph.RelDef{Name: name, Class: graph.Observation, Version: 1, Fields: relFields[name]}
		if _, err := r.AddRel(d); err != nil {
			t.Fatalf("register relation %s: %v", name, err)
		}
	}
	return r
}

// --- harness ----------------------------------------------------------------

// run seeds a graph and expands it with the given operators. Driving the real
// scheduler rather than calling Exec directly is deliberate: it proves the
// applier accepts every type, relation and field name the operators emit, which
// a delta compared against itself never would.
func run(t *testing.T, seedType, seedKey string, ops ...graph.Operator) *graph.Graph {
	t.Helper()
	g := graph.New(testRegistry(t))
	if _, err := g.Seed(seedType, seedKey); err != nil {
		t.Fatalf("seed: %v", err)
	}
	// Attempts 1 keeps a failing fake from being called twice; the retry path
	// belongs to the scheduler's own tests.
	s := graph.NewScheduler(g, ops, graph.Limits{Workers: 1, Attempts: 1, MaxRounds: 8})
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("run: %v", err)
	}
	return g
}

// nodeID resolves a (type, key) pair the way the applier would.
func nodeID(t *testing.T, g *graph.Graph, typ, key string) graph.NodeID {
	t.Helper()
	for _, n := range g.Nodes() {
		if n.Type.Name() == typ && n.Key == key {
			return n.ID
		}
	}
	t.Fatalf("no %s node with key %q; graph has %s", typ, key, dump(g))
	return graph.NodeID{}
}

func hasNode(g *graph.Graph, typ, key string) bool {
	for _, n := range g.Nodes() {
		if n.Type.Name() == typ && n.Key == key {
			return true
		}
	}
	return false
}

// hasEdge reports whether an edge of rel joins the two keyed nodes.
func hasEdge(g *graph.Graph, fromType, fromKey, rel, toType, toKey string) bool {
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

// prop reads a materialized prop off a node.
func prop(t *testing.T, g *graph.Graph, typ, key, field string) (graph.Value, bool) {
	t.Helper()
	n, ok := g.Node(nodeID(t, g, typ, key))
	if !ok {
		return graph.Value{}, false
	}
	f, ok := n.Type.Field(field)
	if !ok {
		t.Fatalf("type %s has no field %q", typ, field)
	}
	return n.Props.Get(f)
}

// edgeProp reads a materialized prop off the first matching edge.
func edgeProp(t *testing.T, g *graph.Graph, rel, toKey, field string) (graph.Value, bool) {
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

// wantStatus asserts the terminal status the scheduler recorded for a pair.
// This is the assertion that matters most in this package: the ok/empty/failed
// /timeout split is the signal, not a detail of it.
func wantStatus(t *testing.T, g *graph.Graph, typ, key, op string, want graph.Status) {
	t.Helper()
	got, ok := g.Status(nodeID(t, g, typ, key), op)
	if !ok {
		t.Fatalf("%s on %s/%s recorded no status, want %s", op, typ, key, want)
	}
	if got != want {
		t.Fatalf("%s on %s/%s = %s, want %s", op, typ, key, got, want)
	}
}

func dump(g *graph.Graph) string {
	var b strings.Builder
	for _, n := range g.Nodes() {
		fmt.Fprintf(&b, "%s(%s) ", n.Type.Name(), n.Key)
	}
	return b.String()
}

// --- stub view --------------------------------------------------------------

// stubView is the "literal view in, expected delta out" harness of DESIGN §13,
// for the cases where driving the scheduler would say less than calling Exec
// directly. No operator here declares Trigger.Reads, so an empty prop and edge
// surface is a faithful view.
type stubView struct{ typ, key string }

func (stubView) ID() graph.NodeID                { return graph.NodeID{} }
func (v stubView) Type() string                  { return v.typ }
func (v stubView) Key() string                   { return v.key }
func (stubView) Depth() int                      { return 0 }
func (stubView) Prop(string) (graph.Value, bool) { return graph.Value{}, false }
func (stubView) Edges(string) []graph.EdgeView   { return nil }
func (v stubView) Ref() graph.NodeRef            { return graph.NodeRef{Type: v.typ, Key: v.key} }

// --- fakes ------------------------------------------------------------------

// fakeResolver answers from a table. err, when set, is returned by every lookup
// so a test can pin one failure mode at a time.
type fakeResolver struct {
	ips   map[string][]net.IP
	ns    map[string][]*net.NS
	mx    map[string][]*net.MX
	txt   map[string][]string
	cname map[string]string
	addr  map[string][]string
	err   error
	calls int
}

func (f *fakeResolver) LookupIP(_ context.Context, _, host string) ([]net.IP, error) {
	f.calls++
	return f.ips[host], f.err
}

func (f *fakeResolver) LookupNS(_ context.Context, host string) ([]*net.NS, error) {
	f.calls++
	return f.ns[host], f.err
}

func (f *fakeResolver) LookupMX(_ context.Context, host string) ([]*net.MX, error) {
	f.calls++
	return f.mx[host], f.err
}

func (f *fakeResolver) LookupTXT(_ context.Context, host string) ([]string, error) {
	f.calls++
	return f.txt[host], f.err
}

func (f *fakeResolver) LookupCNAME(_ context.Context, host string) (string, error) {
	f.calls++
	return f.cname[host], f.err
}

func (f *fakeResolver) LookupAddr(_ context.Context, addr string) ([]string, error) {
	f.calls++
	return f.addr[addr], f.err
}

// nxdomain is what the Go resolver reports for a name that does not exist. It
// is the single most important error in this package: it must never be Failed.
func nxdomain() error {
	return &net.DNSError{Err: "no such host", Name: "example.com", IsNotFound: true}
}

func servfail() error {
	return &net.DNSError{Err: "server misbehaving", Name: "example.com"}
}

func dnsTimeout() error {
	return &net.DNSError{Err: "i/o timeout", Name: "example.com", IsTimeout: true}
}

// --- package-level wiring ---------------------------------------------------

func TestNewOmitsOperatorsWithoutTheirDependency(t *testing.T) {
	// Geo has no locator here, and the source operators have no source list, so
	// neither may appear: an operator that can only fail would make --explain
	// claim work that will never happen.
	ops := New(Options{Resolver: &fakeResolver{}, Whois: fakeWhois{}})
	for _, op := range ops {
		switch op.Id() {
		case "geo", "pkg", "usr", "repo":
			t.Fatalf("operator %q planned without its dependency", op.Id())
		}
	}

	ops = New(Options{
		Resolver: &fakeResolver{},
		Whois:    fakeWhois{},
		Geo:      fakeGeo{},
		Sources:  fakeSources{},
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
	ops := New(Options{
		Resolver: &fakeResolver{},
		Whois:    fakeWhois{},
		Geo:      fakeGeo{},
		Sources:  fakeSources{},
		Prober:   &fakeProber{},
	})
	want := map[string]string{
		"idn": TypeDomain, "dns-a": TypeDomain, "dns-ns": TypeDomain,
		"dns-mx": TypeDomain, "dns-txt": TypeDomain, "dns-cname": TypeDomain,
		"whois": TypeDomain,
		// The two that would have been producer dependencies: geo and ptr
		// consume addresses, so they bind to ip.
		"geo": TypeIP, "ptr": TypeIP,
		"pkg": TypePackage, "usr": TypeUsername, "repo": TypeRepo,
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
	ops := New(Options{
		Resolver: &fakeResolver{}, Whois: fakeWhois{}, Geo: fakeGeo{},
		Sources: fakeSources{}, Prober: &fakeProber{},
	})
	want := map[string]string{
		"idn": "", "dns-a": ResourceDNS, "dns-ns": ResourceDNS, "dns-mx": ResourceDNS,
		"dns-txt": ResourceDNS, "dns-cname": ResourceDNS, "ptr": ResourceDNS,
		"geo": ResourceGeo, "whois": ResourceWhois,
		"pkg": ResourceHTTP, "usr": ResourceHTTP, "repo": ResourceHTTP,
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
	r := testRegistry(t)

	// Every prop any operator declares in Effects must resolve, or the operator
	// is emitting into a field that does not exist.
	ops := New(Options{
		Resolver: &fakeResolver{}, Whois: fakeWhois{}, Geo: fakeGeo{},
		Sources: fakeSources{}, Prober: &fakeProber{},
	})
	nodeFields := map[string]bool{}
	for _, defs := range Fields() {
		for _, f := range defs {
			nodeFields[f.Name] = true
		}
	}
	nodeFields[FieldCreated] = true // declared by the registrar, asserted here
	relProps := map[string]bool{}
	for _, defs := range RelFields() {
		for _, f := range defs {
			relProps[f.Name] = true
		}
	}
	for _, op := range ops {
		for _, name := range op.Emits().Props {
			if !nodeFields[name] && !relProps[name] {
				t.Errorf("%s declares prop %q that no field list defines", op.Id(), name)
			}
		}
	}

	for name, defs := range RelFields() {
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
