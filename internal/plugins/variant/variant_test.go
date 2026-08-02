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

package variant

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"

	// The language-driven algorithms need at least one language plugin
	// registered to have anything to say. English is the smallest realistic
	// choice; its datasets are algorithm input and are used as they ship.
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/english"
)

// --- fixtures ---------------------------------------------------------------

// lower is a stand-in canonicalization: enough to converge two spellings of the
// same key, which is all these tests need from the registry.
func lower(s string) (string, error) {
	s = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(s, ".")))
	if s == "" {
		return "", fmt.Errorf("empty key")
	}
	return s, nil
}

// testRegistry registers every Nameable type plus the relations a variant
// operator interacts with. HOSTED_ON is structural so the seed closure can be
// widened in the closure test; RESOLVES_TO is an observation edge, which is how
// a Nameable node gets reached from outside the closure.
func testRegistry(t *testing.T) *graph.Registry {
	t.Helper()
	r := graph.NewRegistry()
	for _, name := range NameableTypes {
		if _, err := r.AddType(graph.NodeTypeDef{
			Name: name, Cap: graph.Nameable, Version: 1, Canonical: lower,
		}); err != nil {
			t.Fatalf("type %s: %v", name, err)
		}
	}
	if _, err := r.AddType(graph.NodeTypeDef{
		Name: "ip", Cap: graph.Observed, Version: 1, Canonical: lower,
	}); err != nil {
		t.Fatalf("type ip: %v", err)
	}
	if _, err := r.AddRel(RelDef()); err != nil {
		t.Fatalf("rel %s: %v", Rel, err)
	}
	for _, d := range []graph.RelDef{
		{Name: "HOSTED_ON", Class: graph.Structural, Version: 1},
		{Name: "RESOLVES_TO", Class: graph.Observation, Version: 1},
		{Name: "PTR_TO", Class: graph.Observation, Version: 1},
	} {
		if _, err := r.AddRel(d); err != nil {
			t.Fatalf("rel %s: %v", d.Name, err)
		}
	}
	return r
}

// testOptions keeps the domain datasets tiny. The shipped lists are 5000
// subdomains and 8600 suffixes, which would make every test a benchmark without
// exercising anything the short lists do not.
func testOptions() Options {
	return Options{
		Subdomains: []string{"www", "mail", "login"},
		Suffixes:   []string{"com", "net", "org", "co.uk"},
	}
}

// seeded returns a graph with one seed node of the given type.
func seeded(t *testing.T, typeName, key string) (*graph.Graph, graph.NodeID) {
	t.Helper()
	g := graph.New(testRegistry(t))
	id, err := g.Seed(typeName, key)
	if err != nil {
		t.Fatalf("seed %s %q: %v", typeName, key, err)
	}
	return g, id
}

// only returns the single operator with this id, or fails.
func only(t *testing.T, id string) graph.Operator {
	t.Helper()
	ops, err := Select(testOptions(), id)
	if err != nil {
		t.Fatalf("select %q: %v", id, err)
	}
	if len(ops) != 1 {
		t.Fatalf("select %q returned %d operators", id, len(ops))
	}
	return ops[0]
}

// run executes one operator against a freshly seeded node and returns what it
// produced together with the graph it was applied to.
func run(t *testing.T, op graph.Operator, typeName, key string) (*graph.Graph, graph.NodeID, graph.Delta, graph.Outcome) {
	t.Helper()
	g, id := seeded(t, typeName, key)
	d, o := execAgainst(g, op, id)
	g.Apply(graph.Provenance{Operator: op.Id(), Round: 1}, id, d)
	return g, id, d, o
}

// execAgainst calls an operator directly against a node. The graph's own view
// constructor is unexported, so this supplies the equivalent: these operators
// declare no reads, and a view that hands out nothing but identity, type and
// key is precisely what that declaration entitles them to. The end-to-end
// tests below drive the real thing through the scheduler.
func execAgainst(g *graph.Graph, op graph.Operator, id graph.NodeID) (graph.Delta, graph.Outcome) {
	return op.Exec(context.Background(), &stubView{g: g, id: id})
}

// stubView is the smallest View that satisfies what these operators declare
// they read: nothing but the node's identity, type and key. Prop and Edges
// return nothing, which is exactly what an undeclared read must return.
type stubView struct {
	g  *graph.Graph
	id graph.NodeID
}

func (v *stubView) ID() graph.NodeID { return v.id }
func (v *stubView) Type() string {
	n, _ := v.g.Node(v.id)
	return n.Type.Name()
}
func (v *stubView) Key() string {
	n, _ := v.g.Node(v.id)
	return n.Key
}
func (v *stubView) Depth() int { return v.g.Depth(v.id) }
func (v *stubView) Prop(string) (graph.Value, bool) {
	return graph.Value{}, false
}
func (v *stubView) Edges(string) []graph.EdgeView { return nil }
func (v *stubView) Ref() graph.NodeRef {
	n, _ := v.g.Node(v.id)
	return graph.NodeRef{Type: n.Type.Name(), Key: n.Key}
}

// variantKeys returns the keys a delta's VARIANT_OF edges point at.
func variantKeys(d graph.Delta) []string {
	out := make([]string, 0, len(d.Edges))
	for _, e := range d.Edges {
		out = append(out, e.To.Key)
	}
	return out
}

func contains(keys []string, want string) bool {
	for _, k := range keys {
		if k == want {
			return true
		}
	}
	return false
}

// --- generation -------------------------------------------------------------

func TestGeneratesVariants(t *testing.T) {
	_, _, d, o := run(t, only(t, "co"), TypeDomain, "example.com")
	if o.Status != graph.StatusOK {
		t.Fatalf("outcome = %v, want ok", o.Status)
	}
	keys := variantKeys(d)
	if len(keys) == 0 {
		t.Fatal("character omission produced no variants")
	}
	for _, want := range []string{"xample.com", "exampl.com", "exmple.com"} {
		if !contains(keys, want) {
			t.Errorf("missing variant %q; got %v", want, keys)
		}
	}
}

// The registrable name is the only thing a general algorithm may vary. If
// omission could reach the suffix it would generate "example.cm" — a different
// registry, attributed to the wrong algorithm.
func TestSuffixIsPreserved(t *testing.T) {
	_, _, d, _ := run(t, only(t, "co"), TypeDomain, "www.example.co.uk")
	for _, k := range variantKeys(d) {
		if !strings.HasSuffix(k, ".co.uk") {
			t.Fatalf("variant %q lost the public suffix", k)
		}
		if !strings.HasPrefix(k, "www.") {
			t.Fatalf("variant %q lost the subdomain prefix", k)
		}
	}
}

func TestEveryOperatorEmitsVariantRel(t *testing.T) {
	ops := All(testOptions())
	if len(ops) == 0 {
		t.Fatal("no operators registered")
	}
	for _, op := range ops {
		e := op.Emits()
		var found bool
		for _, r := range e.Rels {
			if r == Rel {
				found = true
			}
		}
		if !found {
			t.Errorf("%s does not declare %s, so the engine would not treat it as a variant operator", op.Id(), Rel)
		}
		if len(e.Props) != 2 {
			t.Errorf("%s declares props %v, want algorithm and distance", op.Id(), e.Props)
		}
	}
}

// Binding by capability is what lets one algorithm cover every Nameable type.
func TestBindsByCapability(t *testing.T) {
	op := only(t, "co")
	tr := op.Trigger()
	if len(tr.On.Types) != 0 {
		t.Fatalf("co binds by type %v, want capability", tr.On.Types)
	}
	if len(tr.On.Caps) != 1 || tr.On.Caps[0] != graph.Nameable {
		t.Fatalf("co caps = %v, want [nameable]", tr.On.Caps)
	}
	for _, tc := range []struct{ typ, key string }{
		{TypeDomain, "example.com"},
		{TypeUsername, "example"},
		{TypePackage, "example"},
		{TypeRepo, "acme/example"},
		{TypeEmail, "bob@example.com"},
	} {
		_, _, d, o := run(t, op, tc.typ, tc.key)
		if o.Status != graph.StatusOK || len(d.Edges) == 0 {
			t.Errorf("%s %q produced nothing (%v)", tc.typ, tc.key, o.Status)
		}
		for _, e := range d.Edges {
			if e.To.Type != tc.typ {
				t.Errorf("%s variant has type %q, want %q", tc.typ, e.To.Type, tc.typ)
			}
		}
	}
}

// An email's domain half is reached structurally and varied there; varying it
// here as well would produce the same node by two different provenances.
func TestEmailVariesLocalPartOnly(t *testing.T) {
	_, _, d, _ := run(t, only(t, "co"), TypeEmail, "bob@example.com")
	for _, k := range variantKeys(d) {
		if !strings.HasSuffix(k, "@example.com") {
			t.Fatalf("email variant %q changed the domain half", k)
		}
	}
}

// A type-narrowed algorithm must not fire on types it cannot mean anything for.
func TestDomainOnlyAlgorithmsNarrowTheirSelector(t *testing.T) {
	for _, id := range []string{"si", "tld"} {
		tr := only(t, id).Trigger()
		if len(tr.On.Caps) != 0 {
			t.Errorf("%s binds by capability %v; a package has no TLD", id, tr.On.Caps)
		}
		if len(tr.On.Types) != 1 || tr.On.Types[0] != TypeDomain {
			t.Errorf("%s types = %v, want [domain]", id, tr.On.Types)
		}
	}
}

func TestTLDSwapChangesOnlyTheSuffix(t *testing.T) {
	_, _, d, _ := run(t, only(t, "tld"), TypeDomain, "example.com")
	keys := variantKeys(d)
	want := []string{"example.co.uk", "example.net", "example.org"}
	for _, w := range want {
		if !contains(keys, w) {
			t.Errorf("missing %q; got %v", w, keys)
		}
	}
	if contains(keys, "example.com") {
		t.Error("tld swap emitted the origin as its own variant")
	}
}

func TestSubdomainInsertionPrependsLabels(t *testing.T) {
	_, _, d, _ := run(t, only(t, "si"), TypeDomain, "example.com")
	keys := variantKeys(d)
	for _, w := range []string{"www.example.com", "mail.example.com", "login.example.com"} {
		if !contains(keys, w) {
			t.Errorf("missing %q; got %v", w, keys)
		}
	}
}

// An algorithm with nothing to say reports Empty, not Failed. The distinction
// is the whole point of Outcome: absence is an answer.
func TestNothingToSayIsEmptyNotFailed(t *testing.T) {
	_, _, d, o := run(t, only(t, "ho"), TypeDomain, "example.com")
	if o.Status != graph.StatusEmpty {
		t.Fatalf("hyphen omission on a hyphenless name = %v, want empty", o.Status)
	}
	if len(d.Edges) != 0 {
		t.Fatalf("empty outcome carried %d edges", len(d.Edges))
	}
	if o.Err != nil {
		t.Fatalf("empty outcome carried an error: %v", o.Err)
	}
}

func TestNamespaceConfusion(t *testing.T) {
	_, _, d, _ := run(t, only(t, "nsc"), TypePackage, "@acme/tool")
	keys := variantKeys(d)
	for _, w := range []string{"tool", "acme-tool", "acme/tool", "acmetool"} {
		if !contains(keys, w) {
			t.Errorf("missing %q; got %v", w, keys)
		}
	}
}

func TestKeyboardAndLanguageAlgorithmsProduceVariants(t *testing.T) {
	// These need a keyboard or a language plugin loaded; english is imported
	// by this test file for exactly that reason. "google" is the input because
	// repetition-adjacent replacement needs a doubled letter to work on.
	for _, id := range []string{"aci", "acs", "rar", "vs", "gi", "gr"} {
		_, _, d, o := run(t, only(t, id), TypeDomain, "google.com")
		if o.Status != graph.StatusOK || len(d.Edges) == 0 {
			t.Errorf("%s produced nothing (%v); language/keyboard data missing?", id, o.Status)
		}
	}
}

// --- edge props -------------------------------------------------------------

func TestEdgePropsCarryAlgorithmAndDistance(t *testing.T) {
	op := only(t, "co")
	g, seed, d, _ := run(t, op, TypeDomain, "example.com")

	// The delta itself must carry both props for every edge it emits.
	if got, want := len(d.Props), 2*len(d.Edges); got != want {
		t.Fatalf("%d props for %d edges, want %d", got, len(d.Edges), want)
	}

	// And they must survive the applier onto the real edge.
	var checked int
	for _, e := range g.Edges() {
		if e.Rel.Name() != Rel {
			continue
		}
		if e.From != seed {
			t.Fatalf("edge runs from %s, want the origin %s", e.From, seed)
		}
		algoField, ok := e.Rel.Field(PropAlgorithm)
		if !ok {
			t.Fatal("relation has no algorithm field")
		}
		v, set := e.Props.Get(algoField)
		if !set || v.Str() != "co" {
			t.Fatalf("algorithm = %q (set=%v), want co", v.Str(), set)
		}
		distField, ok := e.Rel.Field(PropDistance)
		if !ok {
			t.Fatal("relation has no distance field")
		}
		dv, set := e.Props.Get(distField)
		if !set {
			t.Fatal("distance not set")
		}
		to, _ := g.Node(e.To)
		if dv.Num() < 1 {
			t.Fatalf("distance to %q = %d, want >= 1", to.Key, dv.Num())
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no VARIANT_OF edges were applied")
	}
}

func TestDistanceIsLevenshtein(t *testing.T) {
	// One omitted character from a five-character name is distance 1.
	_, _, d, _ := run(t, only(t, "co"), TypeUsername, "alice")
	for _, ps := range d.Props {
		if ps.Field != PropDistance {
			continue
		}
		if ps.Value.Num() != 1 {
			t.Fatalf("distance = %d for a single omission, want 1", ps.Value.Num())
		}
	}
	// A doubled character is also one edit.
	_, _, d, _ = run(t, only(t, "cr"), TypeUsername, "alice")
	for _, ps := range d.Props {
		if ps.Field == PropDistance && ps.Value.Num() != 1 {
			t.Fatalf("repetition distance = %d, want 1", ps.Value.Num())
		}
	}
}

// A prop must name the edge it belongs to, not alias a shared address.
func TestEachPropNamesItsOwnEdge(t *testing.T) {
	_, _, d, _ := run(t, only(t, "co"), TypeDomain, "example.com")
	seen := map[string]int{}
	for _, ps := range d.Props {
		if ps.Edge == nil {
			t.Fatal("prop names no edge")
		}
		if ps.Node != nil {
			t.Fatal("prop names both a node and an edge")
		}
		if ps.Edge.Rel != Rel {
			t.Fatalf("prop on relation %q, want %s", ps.Edge.Rel, Rel)
		}
		seen[ps.Edge.To.Key]++
	}
	for k, n := range seen {
		if n != 2 {
			t.Fatalf("edge to %q carries %d props, want 2", k, n)
		}
	}
}

// --- delta shape ------------------------------------------------------------

// Operators emit raw keys and let the applier canonicalize. Nothing here may
// look like a NodeID.
func TestDeltaCarriesRawKeysOnly(t *testing.T) {
	_, _, d, _ := run(t, only(t, "co"), TypeDomain, "Example.COM")
	if len(d.Edges) == 0 {
		t.Fatal("no edges")
	}
	for _, e := range d.Edges {
		if e.From.Key == "" || e.To.Key == "" {
			t.Fatal("edge names a node with an empty key")
		}
		if e.From.Type == "" || e.To.Type == "" {
			t.Fatal("edge names a node with no type")
		}
	}
}

// Bare nodes in the delta would survive an applier rejection of the edge that
// justified them, leaving orphan variants behind.
func TestDeltaIntroducesVariantsOnlyByTheirEdge(t *testing.T) {
	for _, op := range All(testOptions()) {
		_, _, d, _ := run(t, op, TypeDomain, "example.com")
		if len(d.Nodes) != 0 {
			t.Fatalf("%s emitted %d bare nodes", op.Id(), len(d.Nodes))
		}
	}
}

// --- determinism ------------------------------------------------------------

// dump renders a delta in a stable textual form. Two runs of the same operator
// on the same input must produce identical bytes.
func dump(d graph.Delta) string {
	var b strings.Builder
	for _, e := range d.Edges {
		fmt.Fprintf(&b, "E %s %s %s -> %s %s\n", e.From.Type, e.From.Key, e.Rel, e.To.Type, e.To.Key)
	}
	for _, p := range d.Props {
		fmt.Fprintf(&b, "P %s %s %d %q\n", p.Edge.To.Key, p.Field, p.Value.Num(), p.Value.Str())
	}
	return b.String()
}

// Several pkg/typo functions accumulate through a map, so their output order is
// randomized per call. The operator sorts, which is what makes this hold.
func TestSameInputYieldsByteIdenticalDelta(t *testing.T) {
	const runs = 8
	for _, op := range All(testOptions()) {
		var first string
		for i := 0; i < runs; i++ {
			_, _, d, _ := run(t, op, TypeDomain, "one.two-three.example.com")
			got := dump(d)
			if i == 0 {
				first = got
				continue
			}
			if got != first {
				t.Fatalf("%s: run %d differs from run 0\n--- run 0 ---\n%s\n--- run %d ---\n%s",
					op.Id(), i, first, i, got)
			}
		}
	}
}

// The same property end to end: two independent expansions must produce the
// same content addresses, node for node and edge for edge.
func TestExpansionIsReproducible(t *testing.T) {
	expand := func() []string {
		g := graph.New(testRegistry(t))
		if _, err := g.Seed(TypeDomain, "example.com"); err != nil {
			t.Fatalf("seed: %v", err)
		}
		s := graph.NewScheduler(g, All(testOptions()), graph.Limits{MaxRounds: 4, Workers: 4})
		if err := s.Run(context.Background()); err != nil {
			t.Fatalf("run: %v", err)
		}
		var out []string
		for _, n := range g.Nodes() {
			_, c, err := n.Addressed()
			if err != nil {
				t.Fatalf("address node %s: %v", n.Key, err)
			}
			out = append(out, "n "+c.String())
		}
		for _, e := range g.Edges() {
			_, c, err := e.Addressed()
			if err != nil {
				t.Fatalf("address edge: %v", err)
			}
			out = append(out, "e "+c.String())
		}
		return out
	}
	a, b := expand(), expand()
	if len(a) == 0 {
		t.Fatal("expansion produced nothing")
	}
	if len(a) != len(b) {
		t.Fatalf("run sizes differ: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("block %d differs: %s vs %s", i, a[i], b[i])
		}
	}
}

// --- engine invariants ------------------------------------------------------

// Variants are terminal: a node reached by VARIANT_OF is never handed back to a
// variant operator. The engine enforces it, and it holds because a variant edge
// does not extend the seed closure — this test is what proves the operators as
// declared actually get that treatment.
func TestVariantsAreTerminal(t *testing.T) {
	g := graph.New(testRegistry(t))
	seed, err := g.Seed(TypeDomain, "example.com")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	ops, err := Select(testOptions(), "co")
	if err != nil {
		t.Fatalf("select: %v", err)
	}
	s := graph.NewScheduler(g, ops, graph.Limits{MaxRounds: 8, Workers: 1})
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("run: %v", err)
	}

	origins := map[graph.NodeID]bool{}
	for _, e := range g.Edges() {
		if e.Rel.Name() == Rel {
			origins[e.From] = true
		}
	}
	if len(origins) != 1 || !origins[seed] {
		t.Fatalf("%d distinct variant origins, want only the seed", len(origins))
	}
	for _, n := range g.Nodes() {
		if n.ID == seed {
			continue
		}
		if g.InClosure(n.ID) {
			t.Fatalf("variant %q entered the seed closure and would be varied again", n.Key)
		}
	}
}

// A Nameable node reached by an observation edge — a PTR result, say — is
// outside the seed closure, and the applier must refuse to root a variant on it
// even if an operator emits one anyway.
func TestVariantOutsideTheClosureIsRejected(t *testing.T) {
	g := graph.New(testRegistry(t))
	seed, err := g.Seed(TypeDomain, "example.com")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	by := graph.Provenance{Operator: "ptr", Round: 1}
	g.Apply(by, seed, graph.Delta{Edges: []graph.EdgeRef{{
		From: graph.NodeRef{Type: TypeDomain, Key: "example.com"},
		Rel:  "RESOLVES_TO",
		To:   graph.NodeRef{Type: "ip", Key: "203.0.113.7"},
	}}})
	var ipID graph.NodeID
	for _, n := range g.Nodes() {
		if n.Type.Name() == "ip" {
			ipID = n.ID
		}
	}
	g.Apply(by, ipID, graph.Delta{Edges: []graph.EdgeRef{{
		From: graph.NodeRef{Type: "ip", Key: "203.0.113.7"},
		Rel:  "PTR_TO",
		To:   graph.NodeRef{Type: TypeDomain, Key: "neighbour.example.net"},
	}}})

	var far graph.NodeID
	for _, n := range g.Nodes() {
		if n.Key == "neighbour.example.net" {
			far = n.ID
		}
	}
	if far.IsZero() {
		t.Fatal("PTR domain was not admitted")
	}
	if g.InClosure(far) {
		t.Fatal("an observation edge widened the seed closure")
	}

	op := only(t, "co")
	d, _ := execAgainst(g, op, far)
	res := g.Apply(graph.Provenance{Operator: op.Id(), Round: 2}, far, d)
	if len(res.Edges) != 0 {
		t.Fatalf("applier admitted %d variant edges rooted outside the closure", len(res.Edges))
	}
	if len(res.Rejected) == 0 {
		t.Fatal("out-of-closure variant edges were not rejected")
	}
	for _, r := range res.Rejected {
		if r.Kind != graph.RejectClosure {
			t.Fatalf("rejection kind = %v, want closure", r.Kind)
		}
	}
}

// The trigger says so as well, so the scheduler can skip the work rather than
// let the applier throw it away.
func TestTriggerDeclaresClosureMembership(t *testing.T) {
	for _, op := range All(testOptions()) {
		var found bool
		for _, c := range op.Trigger().Where {
			if c == graph.InClosure() {
				found = true
			}
		}
		if !found {
			t.Errorf("%s does not gate on seed-closure membership", op.Id())
		}
	}
}

// --- registry ---------------------------------------------------------------

func TestOperatorIdsAreUniqueAndSorted(t *testing.T) {
	ops := All(testOptions())
	ids := make([]string, 0, len(ops))
	for _, op := range ops {
		ids = append(ids, op.Id())
	}
	if !sort.StringsAreSorted(ids) {
		t.Fatalf("ids are not sorted: %v", ids)
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate id %q", id)
		}
		seen[id] = true
	}
	if len(ids) < 20 {
		t.Fatalf("only %d algorithms registered: %v", len(ids), ids)
	}
}

func TestSelectRejectsUnknownIds(t *testing.T) {
	if _, err := Select(testOptions(), "co", "nope"); err == nil {
		t.Fatal("unknown algorithm id was accepted")
	}
}

func TestEveryOperatorHasAVersion(t *testing.T) {
	for _, op := range All(testOptions()) {
		if op.Version() < 1 {
			t.Errorf("%s has version %d", op.Id(), op.Version())
		}
		if op.Resource() != "" {
			t.Errorf("%s claims resource class %q; variant generation is pure computation",
				op.Id(), op.Resource())
		}
	}
}

// --- splitting --------------------------------------------------------------

func TestSplitDomain(t *testing.T) {
	for _, tc := range []struct{ in, prefix, name, suffix string }{
		{"example.com", "", "example", "com"},
		{"www.example.com", "www", "example", "com"},
		{"example.co.uk", "", "example", "co.uk"},
		{"a.b.example.co.uk", "a.b", "example", "co.uk"},
		{"acme.internal", "", "acme", "internal"},
		{"com", "", "", ""}, // a bare public suffix has no registrable label
	} {
		p, n, s := SplitDomain(tc.in)
		if p != tc.prefix || n != tc.name || s != tc.suffix {
			t.Errorf("SplitDomain(%q) = (%q,%q,%q), want (%q,%q,%q)",
				tc.in, p, n, s, tc.prefix, tc.name, tc.suffix)
		}
	}
}

// A key with no registrable label falls back to varying the whole thing rather
// than producing nothing at all.
func TestBarePublicSuffixFallsBackToWholeKey(t *testing.T) {
	_, _, d, o := run(t, only(t, "co"), TypeDomain, "com")
	if o.Status != graph.StatusOK || len(d.Edges) == 0 {
		t.Fatalf("bare suffix produced nothing (%v)", o.Status)
	}
}

func TestNonDomainTypesAreVariedWhole(t *testing.T) {
	p := DefaultSplit(TypePackage, "acme.tool")
	if p.Core != "acme.tool" {
		t.Fatalf("package core = %q, want the whole key", p.Core)
	}
	if p.Join("other.tool") != "other.tool" {
		t.Fatalf("package join = %q", p.Join("other.tool"))
	}
}

func TestPublicSuffixesDropMatchingMarkers(t *testing.T) {
	for _, s := range PublicSuffixes() {
		if strings.HasPrefix(s, "*") || strings.HasPrefix(s, "!") {
			t.Fatalf("suffix list contains a matching rule rather than a name: %q", s)
		}
	}
}
