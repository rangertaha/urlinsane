// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
)

func registry(t *testing.T) *graph.Registry {
	t.Helper()
	r := graph.NewRegistry()
	if err := decompose.Register(r); err != nil {
		t.Fatalf("Register: %v", err)
	}
	return r
}

// expand seeds a graph and runs every decomposer to fixpoint. Using the real
// scheduler rather than calling Exec directly is deliberate: decomposition is
// multi-round — an email's domain only appears after the email operator ran —
// and a test that stepped the operators by hand would not prove that chain
// closes on its own.
func expand(t *testing.T, seedType, seedKey string) (*graph.Graph, graph.NodeID) {
	t.Helper()
	g := graph.New(registry(t))
	id, err := g.Seed(seedType, seedKey)
	if err != nil {
		t.Fatalf("seed %s(%q): %v", seedType, seedKey, err)
	}
	if err := graph.NewScheduler(g, Operators(), graph.Limits{}).Run(context.Background()); err != nil {
		t.Fatalf("run: %v", err)
	}
	return g, id
}

// decomposition is one seed and the complete graph it should expand into.
type decomposition struct {
	name     string
	seedType string
	seedKey  string
	nodes    []string // "type:key", the seed included
	edges    []string // "type:key -REL-> type:key"
}

var decompositions = []decomposition{
	{
		name:     "email splits into local part and domain, and the domain into its tld",
		seedType: decompose.TypeEmail,
		seedKey:  "bob@example.com",
		nodes: []string{
			"email:bob@example.com",
			"username:bob",
			"domain:example.com",
			"tld:com",
		},
		edges: []string{
			"email:bob@example.com -DOMAIN_OF-> domain:example.com",
			"email:bob@example.com -LOCAL_PART-> username:bob",
			"domain:example.com -TLD_OF-> tld:com",
		},
	},
	{
		// The applier canonicalizes, so a shouted seed converges on the same
		// graph. The local part keeps its case; the username derived from it
		// does not.
		name:     "email canonicalization folds the domain but not the local part",
		seedType: decompose.TypeEmail,
		seedKey:  "Bob@Example.COM.",
		nodes: []string{
			"email:Bob@example.com",
			"username:bob",
			"domain:example.com",
			"tld:com",
		},
		edges: []string{
			"email:Bob@example.com -DOMAIN_OF-> domain:example.com",
			"email:Bob@example.com -LOCAL_PART-> username:bob",
			"domain:example.com -TLD_OF-> tld:com",
		},
	},
	{
		name:     "repo splits into platform and owner",
		seedType: decompose.TypeRepo,
		seedKey:  "github.com/acme/tool",
		nodes: []string{
			"repo:github.com/acme/tool",
			"platform:github.com",
			"username:acme",
		},
		edges: []string{
			"repo:github.com/acme/tool -HOSTED_ON-> platform:github.com",
			"repo:github.com/acme/tool -OWNER-> username:acme",
		},
	},
	{
		name:     "a clone URL decomposes identically to a bare repo path",
		seedType: decompose.TypeRepo,
		seedKey:  "git@github.com:Acme/Tool.git",
		nodes: []string{
			"repo:github.com/acme/tool",
			"platform:github.com",
			"username:acme",
		},
		edges: []string{
			"repo:github.com/acme/tool -HOSTED_ON-> platform:github.com",
			"repo:github.com/acme/tool -OWNER-> username:acme",
		},
	},
	{
		name:     "a registrable suffix is kept whole",
		seedType: decompose.TypeDomain,
		seedKey:  "example.co.uk",
		nodes:    []string{"domain:example.co.uk", "tld:co.uk"},
		edges:    []string{"domain:example.co.uk -TLD_OF-> tld:co.uk"},
	},
	{
		// A private suffix is not a registry. Keying on it would rank a
		// squatter's own hosting domain alongside a real TLD.
		name:     "a private suffix falls back to the real tld",
		seedType: decompose.TypeDomain,
		seedKey:  "acme.blogspot.com",
		nodes:    []string{"domain:acme.blogspot.com", "tld:com"},
		edges:    []string{"domain:acme.blogspot.com -TLD_OF-> tld:com"},
	},
	{
		name:     "a scoped package yields its owning namespace",
		seedType: decompose.TypePackage,
		seedKey:  "npm:@Acme/Tool",
		nodes:    []string{"package:npm:@acme/tool", "username:acme"},
		edges:    []string{"package:npm:@acme/tool -OWNER-> username:acme"},
	},
	{
		name:     "an unscoped package has nothing to decompose",
		seedType: decompose.TypePackage,
		seedKey:  "npm:lodash",
		nodes:    []string{"package:npm:lodash"},
	},
	{
		name:     "a single-label host has no suffix to separate out",
		seedType: decompose.TypeDomain,
		seedKey:  "localhost",
		nodes:    []string{"domain:localhost"},
	},
}

func TestDecompose(t *testing.T) {
	for _, tc := range decompositions {
		t.Run(tc.name, func(t *testing.T) {
			g, _ := expand(t, tc.seedType, tc.seedKey)
			assertEqual(t, "nodes", nodeNames(g), tc.nodes)
			assertEqual(t, "edges", edgeNames(g), tc.edges)
		})
	}
}

// TestDecompositionCostsNoDepth guards the reason decomposition edges are
// structural. Under whole-edge counting, bob@example.com would spend its depth
// budget getting to its own domain and never reach an IP — gutting the
// composite-input case the design exists for.
func TestDecompositionCostsNoDepth(t *testing.T) {
	for _, tc := range decompositions {
		t.Run(tc.name, func(t *testing.T) {
			g, _ := expand(t, tc.seedType, tc.seedKey)
			for _, n := range g.Nodes() {
				if d := g.Depth(n.ID); d != 0 {
					t.Errorf("%s:%s is at depth %d, want 0", n.Type.Name(), n.Key, d)
				}
			}
		})
	}
}

// TestDecomposedNodesAreInSeedClosure guards variant eligibility. Only closure
// members may root variant generation, so a decomposed part that fell outside
// it would be silently unvariable — the scan would run, report nothing, and
// read as complete.
func TestDecomposedNodesAreInSeedClosure(t *testing.T) {
	for _, tc := range decompositions {
		t.Run(tc.name, func(t *testing.T) {
			g, _ := expand(t, tc.seedType, tc.seedKey)
			for _, n := range g.Nodes() {
				if !g.InClosure(n.ID) {
					t.Errorf("%s:%s is outside the seed closure", n.Type.Name(), n.Key)
				}
			}
		})
	}
}

// TestDecompositionEdgesAreStructural is the invariant behind the previous two:
// both depth and closure are decided by edge class, not by the operator family.
func TestDecompositionEdgesAreStructural(t *testing.T) {
	for _, tc := range decompositions {
		t.Run(tc.name, func(t *testing.T) {
			g, _ := expand(t, tc.seedType, tc.seedKey)
			for _, e := range g.Edges() {
				if c := e.Rel.Class(); c != graph.Structural {
					t.Errorf("%s is %s, want structural", e.Rel.Name(), c)
				}
			}
		})
	}
}

// TestNothingToDecomposeIsEmptyNotFailed pins the outcome distinction. Empty
// closes the pair; Failed would have the scheduler retry a pure string parse
// inside its round, forever, to no purpose.
func TestNothingToDecomposeIsEmptyNotFailed(t *testing.T) {
	for _, tc := range []struct {
		seedType, seedKey, operator string
		want                        graph.Status
	}{
		{decompose.TypePackage, "npm:lodash", "decompose.package", graph.StatusEmpty},
		{decompose.TypePackage, "npm:@acme/tool", "decompose.package", graph.StatusOK},
		{decompose.TypeDomain, "localhost", "decompose.domain", graph.StatusEmpty},
		{decompose.TypeDomain, "example.com", "decompose.domain", graph.StatusOK},
		{decompose.TypeEmail, "bob@example.com", "decompose.email", graph.StatusOK},
	} {
		t.Run(tc.seedKey, func(t *testing.T) {
			g, seed := expand(t, tc.seedType, tc.seedKey)
			got, ok := g.Status(seed, tc.operator)
			if !ok {
				t.Fatalf("%s recorded no status for the seed", tc.operator)
			}
			if got != tc.want {
				t.Errorf("%s = %s, want %s", tc.operator, got, tc.want)
			}
		})
	}
}

// TestRegisterMatchesDesign holds the schema to docs/DESIGN.md §2 and §1.1. A
// relation registered under the wrong class changes depth accounting and
// closure growth everywhere, without breaking anything visibly.
func TestRegisterMatchesDesign(t *testing.T) {
	r := registry(t)

	types := map[string]graph.Capability{
		decompose.TypeDomain:     graph.Nameable,
		decompose.TypeUsername:   graph.Nameable,
		decompose.TypePackage:    graph.Nameable,
		decompose.TypeRepo:       graph.Nameable,
		decompose.TypeEmail:      graph.Nameable,
		decompose.TypeTLD:        graph.Observed,
		decompose.TypeIP:         graph.Observed,
		decompose.TypeASN:        graph.Observed,
		decompose.TypeRegistrant: graph.Observed,
		decompose.TypePlatform:   graph.Observed,
	}
	for name, want := range types {
		t.Run("type/"+name, func(t *testing.T) {
			got, ok := r.Type(name)
			if !ok {
				t.Fatalf("%s is not registered", name)
			}
			if got.Cap() != want {
				t.Errorf("%s is %s, want %s", name, got.Cap(), want)
			}
		})
	}
	if n := len(r.Types()); n != len(types) {
		t.Errorf("registry holds %d types, want %d; §2 is authoritative", n, len(types))
	}

	rels := map[string]graph.Class{
		decompose.RelLocalPart:    graph.Structural,
		decompose.RelDomainOf:     graph.Structural,
		decompose.RelTLDOf:        graph.Structural,
		decompose.RelOwner:        graph.Structural,
		decompose.RelHostedOn:     graph.Structural,
		decompose.RelVariantOf:    graph.Variant,
		decompose.RelManifest:     graph.Observation,
		decompose.RelResolvesTo:   graph.Observation,
		decompose.RelNS:           graph.Observation,
		decompose.RelMX:           graph.Observation,
		decompose.RelPTRTo:        graph.Observation,
		decompose.RelInASN:        graph.Observation,
		decompose.RelRegisteredBy: graph.Observation,
		decompose.RelExistsOn:     graph.Observation,
	}
	for name, want := range rels {
		t.Run("rel/"+name, func(t *testing.T) {
			got, ok := r.Rel(name)
			if !ok {
				t.Fatalf("%s is not registered", name)
			}
			if got.Class() != want {
				t.Errorf("%s is %s, want %s", name, got.Class(), want)
			}
		})
	}
	if n := len(r.Rels()); n != len(rels) {
		t.Errorf("registry holds %d relations, want %d; §1.1 is authoritative", n, len(rels))
	}

	// VARIANT_OF's props are what the per-node analyzers read off the edge.
	variant, _ := r.Rel(decompose.RelVariantOf)
	for _, f := range []string{decompose.FieldAlgorithm, decompose.FieldDistance} {
		if _, ok := variant.Field(f); !ok {
			t.Errorf("%s has no %q field", decompose.RelVariantOf, f)
		}
	}
}

func TestCanonicalization(t *testing.T) {
	for _, tc := range []struct {
		typ, raw, want string
		wantErr        bool
	}{
		{typ: decompose.TypeDomain, raw: "Example.COM.", want: "example.com"},
		{typ: decompose.TypeDomain, raw: " münchen.de ", want: "xn--mnchen-3ya.de"},
		{typ: decompose.TypeDomain, raw: "xn--mnchen-3ya.de", want: "xn--mnchen-3ya.de"},
		// Variant algorithms produce names a strict resolver would reject.
		// Refusing them here would drop the candidates the scan is for.
		{typ: decompose.TypeDomain, raw: "-exampl-e.com", want: "-exampl-e.com"},
		{typ: decompose.TypeDomain, raw: "   ", wantErr: true},
		{typ: decompose.TypeDomain, raw: "http://example.com/x", wantErr: true},

		{typ: decompose.TypeEmail, raw: "Bob@Example.COM", want: "Bob@example.com"},
		{typ: decompose.TypeEmail, raw: "bob", wantErr: true},
		{typ: decompose.TypeEmail, raw: "@example.com", wantErr: true},

		{typ: decompose.TypeUsername, raw: "Acme", want: "acme"},
		{typ: decompose.TypeUsername, raw: "", wantErr: true},

		{typ: decompose.TypeRepo, raw: "github.com/Acme/Tool", want: "github.com/acme/tool"},
		{typ: decompose.TypeRepo, raw: "https://GitHub.com/acme/tool.git", want: "github.com/acme/tool"},
		{typ: decompose.TypeRepo, raw: "git@github.com:acme/tool.git", want: "github.com/acme/tool"},
		{typ: decompose.TypeRepo, raw: "https://github.com/acme/tool/tree/main", want: "github.com/acme/tool"},
		{typ: decompose.TypeRepo, raw: "github.com/acme", wantErr: true},

		{typ: decompose.TypePackage, raw: "npm:Lodash", want: "npm:lodash"},
		{typ: decompose.TypePackage, raw: "NPM:@Acme/Tool", want: "npm:@acme/tool"},
		{typ: decompose.TypePackage, raw: "pypi:Foo_Bar.baz", want: "pypi:foo-bar-baz"},
		{typ: decompose.TypePackage, raw: "pypi:foo--bar", want: "pypi:foo-bar"},
		// Unqualified: npm's lodash and PyPI's lodash are different entities.
		{typ: decompose.TypePackage, raw: "lodash", wantErr: true},

		{typ: decompose.TypeTLD, raw: ".COM", want: "com"},
		{typ: decompose.TypeIP, raw: " 1.2.3.4 ", want: "1.2.3.4"},
		{typ: decompose.TypeIP, raw: "::ffff:1.2.3.4", want: "1.2.3.4"},
		{typ: decompose.TypeIP, raw: "2001:0DB8:0000::1", want: "2001:db8::1"},
		{typ: decompose.TypeIP, raw: "example.com", wantErr: true},
		{typ: decompose.TypeASN, raw: "15169", want: "AS15169"},
		{typ: decompose.TypeASN, raw: "as015169", want: "AS15169"},
		{typ: decompose.TypeASN, raw: "ASN15169", want: "AS15169"},
		{typ: decompose.TypeASN, raw: "AS-15169", wantErr: true},
		{typ: decompose.TypeRegistrant, raw: "  Google   LLC ", want: "google llc"},
		{typ: decompose.TypeRegistrant, raw: " ", wantErr: true},
		{typ: decompose.TypePlatform, raw: "GitHub.com.", want: "github.com"},
	} {
		t.Run(tc.typ+"/"+tc.raw, func(t *testing.T) {
			g := graph.New(registry(t))
			id, err := g.Seed(tc.typ, tc.raw)
			switch {
			case tc.wantErr && err == nil:
				t.Fatalf("canonicalized %q, want a refusal", tc.raw)
			case tc.wantErr:
				return
			case err != nil:
				t.Fatalf("refused %q: %v", tc.raw, err)
			}
			n, _ := g.Node(id)
			if n.Key != tc.want {
				t.Errorf("key = %q, want %q", n.Key, tc.want)
			}
		})
	}
}

// TestDecomposersDeclareNoResource keeps decomposers out of every rate-limit
// class. They parse a string they were handed; throttling them behind a network
// budget would stall expansion for no reason.
func TestDecomposersDeclareNoResource(t *testing.T) {
	for _, op := range Operators() {
		if r := op.Resource(); r != "" {
			t.Errorf("%s declares resource %q, want none", op.Id(), r)
		}
		if op.Version() < 1 {
			t.Errorf("%s has version %d, want >= 1", op.Id(), op.Version())
		}
	}
}

// TestEmittedRelationsAreDeclaredAndStructural checks each operator's Effects
// against the registry. A decomposer that declared VARIANT_OF would be treated
// as a variant operator by the scheduler, which is exactly the sort of thing
// this catches at test time rather than at plan-compile time.
func TestEmittedRelationsAreDeclaredAndStructural(t *testing.T) {
	r := registry(t)
	for _, op := range Operators() {
		t.Run(op.Id(), func(t *testing.T) {
			e := op.Emits()
			if len(e.Rels) == 0 || len(e.Nodes) == 0 {
				t.Fatalf("%s declares no effects", op.Id())
			}
			for _, name := range e.Rels {
				rel, ok := r.Rel(name)
				if !ok {
					t.Errorf("emits unregistered relation %q", name)
					continue
				}
				if rel.Class() != graph.Structural {
					t.Errorf("emits %s, which is %s, not structural", name, rel.Class())
				}
			}
			for _, name := range e.Nodes {
				if _, ok := r.Type(name); !ok {
					t.Errorf("emits unregistered node type %q", name)
				}
			}
		})
	}
}

// --- helpers ----------------------------------------------------------------

func nodeNames(g *graph.Graph) []string {
	var out []string
	for _, n := range g.Nodes() {
		out = append(out, n.Type.Name()+":"+n.Key)
	}
	return out
}

func edgeNames(g *graph.Graph) []string {
	var out []string
	for _, e := range g.Edges() {
		from, _ := g.Node(e.From)
		to, _ := g.Node(e.To)
		out = append(out, fmt.Sprintf("%s:%s -%s-> %s:%s",
			from.Type.Name(), from.Key, e.Rel.Name(), to.Type.Name(), to.Key))
	}
	return out
}

// assertEqual compares two sets rendered as strings. Order is not the property
// under test here — the graph package already pins its own iteration order —
// so both sides are sorted.
func assertEqual(t *testing.T, what string, got, want []string) {
	t.Helper()
	g := append([]string(nil), got...)
	w := append([]string(nil), want...)
	sort.Strings(g)
	sort.Strings(w)
	if strings.Join(g, "\n") == strings.Join(w, "\n") {
		return
	}
	t.Errorf("%s:\n got: %s\nwant: %s", what, format(g), format(w))
}

func format(s []string) string {
	if len(s) == 0 {
		return "(none)"
	}
	return "\n  " + strings.Join(s, "\n  ")
}
