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

package graphstore

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/rangertaha/urlinsane/internal/graph"
)

// --- fixture ----------------------------------------------------------------

// testRegistry mirrors the shape of the real one: a composite nameable seed
// (email), the structural chain it decomposes into, one variant relation and
// one observation relation.
func testRegistry(t *testing.T) *graph.Registry {
	t.Helper()
	r := graph.NewRegistry()

	lower := func(s string) (string, error) {
		s = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(s, ".")))
		if s == "" {
			return "", fmt.Errorf("empty key")
		}
		return s, nil
	}

	types := []graph.NodeTypeDef{
		{Name: "email", Cap: graph.Nameable, Version: 1, Canonical: lower},
		{Name: "user", Cap: graph.Nameable, Version: 1, Canonical: lower},
		{Name: "domain", Cap: graph.Nameable, Version: 1, Canonical: lower, Fields: []graph.FieldDef{
			{Name: "punycode", Kind: graph.KindString},
			{Name: "rank", Kind: graph.KindInt},
			{Name: "created", Kind: graph.KindTime, Merge: graph.Precedence("rdap", "whois")},
			{Name: "wildcard", Kind: graph.KindBool},
			{Name: "score", Kind: graph.KindFloat},
			{Name: "digest", Kind: graph.KindBytes},
		}},
		{Name: "tld", Cap: graph.Observed, Version: 1, Canonical: lower},
		{Name: "ip", Cap: graph.Observed, Version: 1, Canonical: lower, Fields: []graph.FieldDef{
			{Name: "asn", Kind: graph.KindString},
		}},
	}
	for _, d := range types {
		if _, err := r.AddType(d); err != nil {
			t.Fatalf("%s: %v", d.Name, err)
		}
	}

	rels := []graph.RelDef{
		{Name: "LOCAL_PART", Class: graph.Structural, Version: 1},
		{Name: "DOMAIN_OF", Class: graph.Structural, Version: 1},
		{Name: "TLD_OF", Class: graph.Structural, Version: 1},
		{Name: "VARIANT_OF", Class: graph.Variant, Version: 1, Fields: []graph.FieldDef{
			{Name: "algorithm", Kind: graph.KindString},
			{Name: "distance", Kind: graph.KindInt},
		}},
		{Name: "RESOLVES_TO", Class: graph.Observation, Version: 1},
	}
	for _, d := range rels {
		if _, err := r.AddRel(d); err != nil {
			t.Fatalf("%s: %v", d.Name, err)
		}
	}
	return r
}

// buildOpts perturbs the fixture so a test can say exactly what differs between
// two scans.
type buildOpts struct {
	rank       int64 // rank prop on example.com; 0 uses the default
	extraTypo  bool  // admit one more variant domain
	dropTypo   bool  // omit the variant domain and everything under it
	shiftRound int   // add to every provenance round, changing nothing addressed
}

type fixture struct {
	g    *graph.Graph
	seed graph.NodeID
	ids  map[string]graph.NodeID
}

func ref(typ, key string) graph.NodeRef { return graph.NodeRef{Type: typ, Key: key} }

// build assembles a small but structurally complete graph: a three-hop
// structural chain from the seed, a variant hanging off the middle of it, and
// an observation below that.
func build(t *testing.T, o buildOpts) *fixture {
	t.Helper()
	g := graph.New(testRegistry(t))
	seed, err := g.Seed("email", "Bob@Example.COM")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	f := &fixture{g: g, seed: seed, ids: map[string]graph.NodeID{}}
	f.ids["email:bob@example.com"] = seed

	by := func(op string, round int) graph.Provenance {
		return graph.Provenance{Operator: op, Round: round + o.shiftRound}
	}

	emailRef := ref("email", "bob@example.com")
	userRef := ref("user", "bob")
	domRef := ref("domain", "example.com")
	tldRef := ref("tld", "com")

	// Structural decomposition. The chain is three hops deep on purpose: a
	// replay that applied edges in the wrong order would leave the tld outside
	// the seed closure.
	f.apply(t, by("decompose", 1), seed, graph.Delta{Edges: []graph.EdgeRef{
		{From: emailRef, Rel: "LOCAL_PART", To: userRef},
		{From: emailRef, Rel: "DOMAIN_OF", To: domRef},
	}})
	f.apply(t, by("decompose", 1), seed, graph.Delta{Edges: []graph.EdgeRef{
		{From: domRef, Rel: "TLD_OF", To: tldRef},
	}})
	f.record(t, userRef, domRef, tldRef)

	dom := f.ids["domain:example.com"]

	rank := o.rank
	if rank == 0 {
		rank = 1000
	}
	created := time.Date(2019, 3, 4, 5, 6, 7, 0, time.UTC)
	f.apply(t, by("rdap", 2), dom, graph.Delta{Props: []graph.PropSet{
		{Node: &domRef, Field: "punycode", Value: graph.String("example.com")},
		{Node: &domRef, Field: "rank", Value: graph.Int(rank)},
		{Node: &domRef, Field: "wildcard", Value: graph.Bool(true)},
		{Node: &domRef, Field: "score", Value: graph.Float(0.5)},
		{Node: &domRef, Field: "digest", Value: graph.Bytes([]byte{0xde, 0xad})},
		{Node: &domRef, Field: "created", Value: graph.Time(created)},
	}})
	// A losing assertion: whois ranks below rdap, so the materialized value is
	// unchanged but the claim is retained in the provenance table.
	f.apply(t, by("whois", 2), dom, graph.Delta{Props: []graph.PropSet{
		{Node: &domRef, Field: "created", Value: graph.Time(created.Add(48 * time.Hour))},
	}})

	if !o.dropTypo {
		typoRef := ref("domain", "exmaple.com")
		f.apply(t, by("omission", 3), dom, graph.Delta{
			Edges: []graph.EdgeRef{{From: domRef, Rel: "VARIANT_OF", To: typoRef}},
			Props: []graph.PropSet{
				{Edge: &graph.EdgeRef{From: domRef, Rel: "VARIANT_OF", To: typoRef},
					Field: "algorithm", Value: graph.String("omission")},
				{Edge: &graph.EdgeRef{From: domRef, Rel: "VARIANT_OF", To: typoRef},
					Field: "distance", Value: graph.Int(1)},
			},
		})
		f.record(t, typoRef)

		ipRef := ref("ip", "93.184.216.34")
		typo := f.ids["domain:exmaple.com"]
		f.apply(t, by("dns", 4), typo, graph.Delta{
			Edges: []graph.EdgeRef{{From: typoRef, Rel: "RESOLVES_TO", To: ipRef}},
			Props: []graph.PropSet{{Node: &ipRef, Field: "asn", Value: graph.String("AS15133")}},
		})
		f.record(t, ipRef)
		g.SetStatus(typo, "dns", graph.StatusOK)
		g.SetScore(typo, "parked", 0.75)
	}

	if o.extraTypo {
		otherRef := ref("domain", "examp1e.com")
		f.apply(t, by("homoglyph", 3), dom, graph.Delta{
			Edges: []graph.EdgeRef{{From: domRef, Rel: "VARIANT_OF", To: otherRef}},
		})
		f.record(t, otherRef)
	}

	g.SetStatus(dom, "dns", graph.StatusOK)
	g.SetStatus(f.ids["tld:com"], "dns", graph.StatusEmpty)
	g.SetScore(dom, "parked", 0.25)

	if err := g.Decline("domain", "Pruned.COM", 1, 0.05, graph.ReasonBelief, by("omission", 3)); err != nil {
		t.Fatalf("decline: %v", err)
	}
	g.NoteTruncation(graph.ReasonRoundCap, 4, "round cap reached")
	g.AddFindings(graph.Finding{
		Kind:     "high-risk-variant",
		Severity: graph.SeverityHigh,
		Nodes:    []graph.NodeID{dom},
		Declined: []graph.LedgerRef{{Type: "domain", Key: "pruned.com"}},
		Summary:  "a variant resolves and was registered recently",
		Evidence: []graph.Provenance{by("dns", 4)},
	})
	return f
}

func (f *fixture) apply(t *testing.T, by graph.Provenance, subject graph.NodeID, d graph.Delta) {
	t.Helper()
	res := f.g.Apply(by, subject, d)
	if len(res.Rejected) > 0 {
		t.Fatalf("delta rejected: %+v", res.Rejected)
	}
}

// record resolves refs to ids after they have been admitted.
func (f *fixture) record(t *testing.T, refs ...graph.NodeRef) {
	t.Helper()
	for _, r := range refs {
		res := f.g.Apply(graph.Provenance{Operator: "test"}, f.seed, graph.Delta{Nodes: []graph.NodeRef{r}})
		if len(res.Nodes) != 1 {
			t.Fatalf("cannot resolve %s/%s: %+v", r.Type, r.Key, res.Rejected)
		}
		f.ids[r.Type+":"+r.Key] = res.Nodes[0]
	}
}

func (f *fixture) save(t *testing.T, s *Store) cid.Cid {
	t.Helper()
	c, err := s.Save(f.g, SaveOptions{Seed: f.seed})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	return c
}

func memStore() *Store { return New(NewMemBlockstore()) }

// --- block round trip -------------------------------------------------------

func TestNodeBlockRoundTrips(t *testing.T) {
	f := build(t, buildOpts{})
	for _, n := range f.g.Nodes() {
		block, c, err := EncodeNode(n)
		if err != nil {
			t.Fatalf("encode %s: %v", n.Key, err)
		}
		nb, err := DecodeNode(block)
		if err != nil {
			t.Fatalf("decode %s: %v", n.Key, err)
		}
		if nb.Type != n.Type.Name() || nb.Key != n.Key {
			t.Fatalf("decoded %s/%s, want %s/%s", nb.Type, nb.Key, n.Type.Name(), n.Key)
		}
		// Every declared slot is written, set or not, so the decoded arity is
		// the type's field count rather than the number of set props.
		set := 0
		n.Props.Each(func(graph.Field, graph.Value) { set++ })
		got := 0
		for _, v := range nb.Values {
			if v.Set {
				got++
			}
		}
		if got != set {
			t.Fatalf("%s: decoded %d set slots, want %d", n.Key, got, set)
		}
		if _, c2, err := EncodeNode(n); err != nil || c2 != c {
			t.Fatalf("%s: re-encoding is not stable", n.Key)
		}
	}
}

func TestEdgeBlockRoundTrips(t *testing.T) {
	f := build(t, buildOpts{})
	for _, e := range f.g.Edges() {
		block, _, err := EncodeEdge(e)
		if err != nil {
			t.Fatalf("encode %s: %v", e.Rel.Name(), err)
		}
		eb, err := DecodeEdge(block)
		if err != nil {
			t.Fatalf("decode %s: %v", e.Rel.Name(), err)
		}
		if eb.From != e.From || eb.To != e.To || eb.Rel != e.Rel.Name() {
			t.Fatalf("edge round trip lost identity: %+v", eb)
		}
	}
}

// TestDecodedSlotsBindToDeclaredKinds covers the one thing the block cannot
// carry: a KindTime and a KindInt are both a CBOR integer, so only the
// registry's declared kind tells them apart.
func TestDecodedSlotsBindToDeclaredKinds(t *testing.T) {
	f := build(t, buildOpts{})
	dom, _ := f.g.Node(f.ids["domain:example.com"])
	block, _, err := EncodeNode(dom)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	nb, err := DecodeNode(block)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Slots follow the registry's declaration order: punycode, rank, created,
	// wildcard, score, digest.
	for _, tc := range []struct {
		slot int
		kind graph.Kind
		want string
	}{
		{0, graph.KindString, "example.com"},
		{1, graph.KindInt, "1000"},
		{2, graph.KindTime, "2019-03-04 05:06:07 +0000 UTC"},
		{3, graph.KindBool, "true"},
		{4, graph.KindFloat, "0.5"},
	} {
		v, ok, err := nb.Slot(tc.slot, tc.kind)
		if err != nil || !ok {
			t.Fatalf("slot %d: ok=%v err=%v", tc.slot, ok, err)
		}
		var got string
		switch tc.kind {
		case graph.KindString:
			got = v.Str()
		case graph.KindInt:
			got = fmt.Sprint(v.Num())
		case graph.KindTime:
			got = v.Time().String()
		case graph.KindBool:
			got = fmt.Sprint(v.Flag())
		case graph.KindFloat:
			got = fmt.Sprint(v.Real())
		}
		if got != tc.want {
			t.Fatalf("slot %d = %q, want %q", tc.slot, got, tc.want)
		}
	}

	// The same integer read as the wrong declared kind must not silently
	// succeed as a different type.
	if _, _, err := nb.Slot(0, graph.KindInt); err == nil {
		t.Fatal("binding a string slot as an int should fail")
	}
	if k := nb.Values[2].Kind; k != datamodel.Kind_Int {
		t.Fatalf("a time slot is %s on the wire, want int", k)
	}
}

// --- determinism ------------------------------------------------------------

// TestIdenticalScansProduceIdenticalRoots is the property everything else rests
// on: with it, "what changed since last week" is a CID comparison; without it,
// every re-scan looks like a total rewrite.
func TestIdenticalScansProduceIdenticalRoots(t *testing.T) {
	a := build(t, buildOpts{}).save(t, memStore())
	b := build(t, buildOpts{}).save(t, memStore())
	if a != b {
		t.Fatalf("identical scans produced different roots:\n  %s\n  %s", a, b)
	}
	// And again, in case an iteration order happened to agree twice.
	if c := build(t, buildOpts{}).save(t, memStore()); c != a {
		t.Fatalf("third identical scan produced %s, want %s", c, a)
	}
}

// TestProvenanceStaysOutOfAddressedForm is §1.2's whole argument. Two scans
// that learned the same facts in different rounds must produce the same node
// and edge addresses, or a diff would report every node as changed every week.
func TestProvenanceStaysOutOfAddressedForm(t *testing.T) {
	s := memStore()
	base := build(t, buildOpts{}).save(t, s)
	shifted := build(t, buildOpts{shiftRound: 7}).save(t, s)

	d, err := s.Diff(base, shifted)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if !d.Empty() {
		t.Fatalf("shifting every round changed the graph: %+v", d)
	}
	// The root itself does differ: it links the side block, and the side block
	// is where provenance is supposed to live.
	if base == shifted {
		t.Fatal("the side tables did not reach the scan root at all")
	}
}

// --- change isolation -------------------------------------------------------

// TestChangedPropTouchesOneNodeAndTheRoot is the second half of the addressing
// argument: a change must propagate up to the root, and nowhere sideways.
func TestChangedPropTouchesOneNodeAndTheRoot(t *testing.T) {
	s := memStore()
	before := build(t, buildOpts{}).save(t, s)
	after := build(t, buildOpts{rank: 42}).save(t, s)

	if before == after {
		t.Fatal("changing a prop did not change the scan root")
	}

	d, err := s.Diff(before, after)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if len(d.Added) != 0 || len(d.Removed) != 0 {
		t.Fatalf("a prop change added or removed nodes: %+v", d)
	}
	if len(d.Changed) != 1 {
		t.Fatalf("%d nodes changed, want exactly 1: %+v", len(d.Changed), d.Changed)
	}
	if got := d.Changed[0]; got.Type != "domain" || got.Key != "example.com" {
		t.Fatalf("the wrong node changed: %s/%s", got.Type, got.Key)
	}
	if slots := d.Changed[0].Slots; len(slots) != 1 || slots[0] != 1 {
		t.Fatalf("changed slots = %v, want just the rank at index 1", slots)
	}
	// Edges address nodes by NodeID, which is identity and does not move when
	// props change.
	if len(d.EdgesAdded)+len(d.EdgesRemoved)+len(d.EdgesChanged) != 0 {
		t.Fatalf("a node prop change disturbed the edges: %+v", d)
	}
	if len(d.Same) == 0 {
		t.Fatal("no nodes reported unchanged")
	}
}

// --- diff -------------------------------------------------------------------

func TestDiffDetectsAddedRemovedAndChanged(t *testing.T) {
	s := memStore()
	older := build(t, buildOpts{}).save(t, s)
	newer := build(t, buildOpts{rank: 7, extraTypo: true, dropTypo: true}).save(t, s)

	d, err := s.Diff(older, newer)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}

	if !hasNode(d.Added, "domain", "examp1e.com") {
		t.Fatalf("added nodes = %+v, want examp1e.com", d.Added)
	}
	if !hasNode(d.Removed, "domain", "exmaple.com") || !hasNode(d.Removed, "ip", "93.184.216.34") {
		t.Fatalf("removed nodes = %+v, want exmaple.com and the ip", d.Removed)
	}
	if !hasNode(d.Changed, "domain", "example.com") {
		t.Fatalf("changed nodes = %+v, want example.com", d.Changed)
	}
	if !hasNode(d.Same, "tld", "com") || !hasNode(d.Same, "user", "bob") {
		t.Fatalf("unchanged nodes = %+v, want the structural chain", d.Same)
	}
	if len(d.EdgesAdded) != 1 || d.EdgesAdded[0].Rel != "VARIANT_OF" {
		t.Fatalf("added edges = %+v, want one VARIANT_OF", d.EdgesAdded)
	}
	if len(d.EdgesRemoved) != 2 {
		t.Fatalf("removed edges = %+v, want the variant and its resolution", d.EdgesRemoved)
	}
}

func TestDiffOfAScanWithItselfIsEmpty(t *testing.T) {
	s := memStore()
	root := build(t, buildOpts{}).save(t, s)
	d, err := s.Diff(root, root)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if !d.Empty() {
		t.Fatalf("a scan differs from itself: %+v", d)
	}
	if len(d.Same) != len(build(t, buildOpts{}).g.Nodes()) {
		t.Fatalf("Same holds %d nodes, want all of them", len(d.Same))
	}
}

func hasNode(s []NodeChange, typ, key string) bool {
	for _, c := range s {
		if c.Type == typ && c.Key == key {
			return true
		}
	}
	return false
}

// --- rehydration ------------------------------------------------------------

// TestRehydrateThenReEncodeIsStable is the resume contract: a graph rebuilt
// from a scan root re-encodes to that same root, so resuming and immediately
// saving is a no-op rather than a rewrite.
func TestRehydrateThenReEncodeIsStable(t *testing.T) {
	s := memStore()
	root := build(t, buildOpts{}).save(t, s)

	r, err := s.Rehydrate(root, testRegistry(t))
	if err != nil {
		t.Fatalf("rehydrate: %v", err)
	}
	again, err := s.Save(r.Graph, SaveOptions{Seed: r.Seed})
	if err != nil {
		t.Fatalf("re-save: %v", err)
	}
	if again != root {
		t.Fatalf("rehydrate/re-encode moved the root:\n  before %s\n  after  %s", root, again)
	}
}

func TestRehydrateRestoresStructureAndSideTables(t *testing.T) {
	s := memStore()
	f := build(t, buildOpts{})
	root := f.save(t, s)

	r, err := s.Rehydrate(root, testRegistry(t))
	if err != nil {
		t.Fatalf("rehydrate: %v", err)
	}
	g := r.Graph

	if got, want := len(g.Nodes()), len(f.g.Nodes()); got != want {
		t.Fatalf("rebuilt %d nodes, want %d", got, want)
	}
	if got, want := len(g.Edges()), len(f.g.Edges()); got != want {
		t.Fatalf("rebuilt %d edges, want %d", got, want)
	}

	// Depth and closure are scheduler state, and a resume that got them wrong
	// would silently re-expand or refuse to expand.
	for _, n := range f.g.Nodes() {
		if g.Depth(n.ID) != f.g.Depth(n.ID) {
			t.Fatalf("%s depth = %d, want %d", n.Key, g.Depth(n.ID), f.g.Depth(n.ID))
		}
		if g.InClosure(n.ID) != f.g.InClosure(n.ID) {
			t.Fatalf("%s closure = %v, want %v", n.Key, g.InClosure(n.ID), f.g.InClosure(n.ID))
		}
	}
	if !g.InClosure(f.ids["tld:com"]) {
		t.Fatal("the far end of the structural chain fell out of the seed closure")
	}
	if g.InClosure(f.ids["domain:exmaple.com"]) {
		t.Fatal("a variant was rehydrated into the seed closure")
	}
	if d := g.Depth(f.ids["ip:93.184.216.34"]); d != 1 {
		t.Fatalf("the observed ip is at depth %d, want 1", d)
	}

	// The losing whois assertion has to survive: §1.4 keeps disagreement as
	// signal rather than resolving it away.
	dom := f.ids["domain:example.com"]
	if got, want := len(g.Assertions(dom)), len(f.g.Assertions(dom)); got != want {
		t.Fatalf("%d assertions restored, want %d", got, want)
	}
	var sawLoser bool
	for _, a := range g.Assertions(dom) {
		if a.By.Operator == "whois" && a.Field == "created" && !a.Won {
			sawLoser = true
		}
	}
	if !sawLoser {
		t.Fatal("the losing whois assertion was not restored")
	}

	if st, ok := g.Status(f.ids["tld:com"], "dns"); !ok || st != graph.StatusEmpty {
		t.Fatalf("status = %v/%v, want empty", st, ok)
	}
	if v, ok := g.Score(dom, "parked"); !ok || v != 0.25 {
		t.Fatalf("score = %v/%v, want 0.25", v, ok)
	}

	ledger := g.Ledger()
	if len(ledger) != 1 || ledger[0].Key != "pruned.com" || ledger[0].Reason != graph.ReasonBelief {
		t.Fatalf("ledger = %+v, want one belief-pruned row", ledger)
	}
	if tr := g.Truncations(); len(tr) != 1 || tr[0].Reason != graph.ReasonRoundCap {
		t.Fatalf("truncations = %+v, want one round-cap row", tr)
	}
	if fs := g.Findings(); len(fs) != 1 || fs[0].Severity != graph.SeverityHigh {
		t.Fatalf("findings = %+v, want one high finding", fs)
	}
}

// TestRehydratedLedgerStillDenies is why the ledger is restored at all: it is a
// denylist as much as a record, and a resume that forgot it would re-admit a
// candidate the first run pruned.
func TestRehydratedLedgerStillDenies(t *testing.T) {
	s := memStore()
	root := build(t, buildOpts{}).save(t, s)
	r, err := s.Rehydrate(root, testRegistry(t))
	if err != nil {
		t.Fatalf("rehydrate: %v", err)
	}
	res := r.Graph.Apply(
		graph.Provenance{Operator: "another", Round: 9}, r.Seed,
		graph.Delta{Nodes: []graph.NodeRef{ref("domain", "PRUNED.com")}},
	)
	if len(res.Nodes) != 0 {
		t.Fatal("a declined candidate was re-admitted after a resume")
	}
	if len(res.Rejected) != 1 || res.Rejected[0].Kind != graph.RejectDenied {
		t.Fatalf("rejections = %+v, want one RejectDenied", res.Rejected)
	}
}

// TestResumeAfterRehydrateKeepsTheRootStable checks the actual resume path:
// expanding a rehydrated graph and saving produces a root that differs from the
// stored one only by the new work.
func TestResumeAfterRehydrateKeepsTheRootStable(t *testing.T) {
	s := memStore()
	partial := build(t, buildOpts{dropTypo: true})
	root := partial.save(t, s)

	r, err := s.Rehydrate(root, testRegistry(t))
	if err != nil {
		t.Fatalf("rehydrate: %v", err)
	}
	domRef := ref("domain", "example.com")
	typoRef := ref("domain", "exmaple.com")
	res := r.Graph.Apply(graph.Provenance{Operator: "omission", Round: 5}, r.Seed, graph.Delta{
		Edges: []graph.EdgeRef{{From: domRef, Rel: "VARIANT_OF", To: typoRef}},
	})
	if len(res.Rejected) > 0 {
		t.Fatalf("resumed expansion was rejected: %+v", res.Rejected)
	}

	resumed, err := s.Save(r.Graph, SaveOptions{Seed: r.Seed})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	d, err := s.Diff(root, resumed)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if len(d.Added) != 1 || !hasNode(d.Added, "domain", "exmaple.com") {
		t.Fatalf("resume added %+v, want exactly the new variant", d.Added)
	}
	if len(d.Changed) != 0 || len(d.Removed) != 0 {
		t.Fatalf("resume disturbed existing nodes: %+v", d)
	}
}

// --- blockstore -------------------------------------------------------------

func TestBlockstorePutGetHas(t *testing.T) {
	for name, bs := range map[string]Blockstore{
		"mem": NewMemBlockstore(),
		"fs":  mustFS(t),
	} {
		block, c, err := encodeList(2, func(e *enc) { e.str("hello"); e.i64(3) })
		if err != nil {
			t.Fatalf("%s: encode: %v", name, err)
		}
		if ok, err := bs.Has(c); err != nil || ok {
			t.Fatalf("%s: an empty store already has the block", name)
		}
		if _, err := bs.Get(c); err == nil {
			t.Fatalf("%s: Get on a missing block should fail", name)
		}
		if err := bs.Put(c, block); err != nil {
			t.Fatalf("%s: put: %v", name, err)
		}
		if ok, err := bs.Has(c); err != nil || !ok {
			t.Fatalf("%s: Has = %v after Put", name, ok)
		}
		got, err := bs.Get(c)
		if err != nil || string(got) != string(block) {
			t.Fatalf("%s: Get returned %q, want %q", name, got, block)
		}
		// Content addressing makes Put idempotent.
		if err := bs.Put(c, block); err != nil {
			t.Fatalf("%s: second put: %v", name, err)
		}
	}
}

func mustFS(t *testing.T) Blockstore {
	t.Helper()
	bs, err := NewFSBlockstore(t.TempDir())
	if err != nil {
		t.Fatalf("fs blockstore: %v", err)
	}
	return bs
}

// TestFilesystemStoreRoundTrips exercises the production path end to end.
func TestFilesystemStoreRoundTrips(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	f := build(t, buildOpts{})
	root := f.save(t, s)

	scan, err := s.Load(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if scan.Root.SeedType != "email" || scan.Root.SeedKey != "bob@example.com" {
		t.Fatalf("seed = %s/%s", scan.Root.SeedType, scan.Root.SeedKey)
	}
	if len(scan.Nodes) != len(f.g.Nodes()) || len(scan.Edges) != len(f.g.Edges()) {
		t.Fatalf("loaded %d nodes and %d edges", len(scan.Nodes), len(scan.Edges))
	}
	if _, ok := scan.Node("tld", "com"); !ok {
		t.Fatal("the tld node is missing from the loaded scan")
	}
	if len(scan.Side.NodeProps) == 0 || len(scan.Side.EdgeProps) == 0 {
		t.Fatal("the side block lost its provenance rows")
	}
	if len(scan.Side.Sched) != len(scan.Nodes) {
		t.Fatalf("%d sched rows for %d nodes", len(scan.Side.Sched), len(scan.Nodes))
	}
}

// --- guards -----------------------------------------------------------------

func TestSaveRequiresASeed(t *testing.T) {
	f := build(t, buildOpts{})
	if _, err := memStore().Save(f.g, SaveOptions{}); err == nil {
		t.Fatal("saving without a seed should fail: the seed cannot be inferred")
	}
}

// TestStatusAndScoresSurviveInSortedOrder guards the two side tables that are
// built from Go maps inside the graph: they reach the block only through the
// analysis surface, which sorts them.
func TestStatusAndScoresSurviveInSortedOrder(t *testing.T) {
	s := memStore()
	f := build(t, buildOpts{})
	scan, err := s.Load(f.save(t, s))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(scan.Side.Status) != 3 {
		t.Fatalf("%d status rows, want 3", len(scan.Side.Status))
	}
	if len(scan.Side.Scores) != 2 {
		t.Fatalf("%d score rows, want 2", len(scan.Side.Scores))
	}
	// Rows follow the node order, then operator order within a node.
	for i := 1; i < len(scan.Side.Status); i++ {
		if scan.Side.Status[i-1].Node == scan.Side.Status[i].Node &&
			scan.Side.Status[i-1].Operator > scan.Side.Status[i].Operator {
			t.Fatal("status rows are not sorted by operator within a node")
		}
	}
}
