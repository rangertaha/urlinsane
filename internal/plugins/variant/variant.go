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

// Package variant holds the typo/squatting algorithms as graph operators.
//
// A variant operator is any operator that declares VARIANT_OF in Emits().Rels.
// That declaration — not a naming convention and not a separate plugin family —
// is what subjects it to the terminal rule and the seed-closure restriction the
// engine enforces (docs/DESIGN.md §4.1, §8).
//
// Every operator here binds by *capability* rather than by type, so a single
// omission algorithm covers domain, username, package, repo and email instead
// of being registered once per type. The few algorithms that are genuinely
// type-specific — TLD swapping is meaningless for a package name — narrow their
// selector by naming types explicitly.
//
// Almost every algorithm is a pure string -> []string function over the
// entity's name, so they are declared as data (id, version, generator) and
// share one adapter rather than being thirty hand-written operator structs.
package variant

import (
	"context"

	"errors"
	"sort"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
)

// Rel is the relation every operator in this package emits. It is re-exported
// from the graph package so a caller registering the relation and an operator
// declaring it cannot drift apart.
const Rel = graph.VariantRel

// Edge props carried by every VARIANT_OF edge. The per-node analyzers read them
// straight off the edge, which is why the View hands out edges rather than bare
// neighbour nodes.
const (
	// PropAlgorithm is the id of the operator that generated the variant.
	PropAlgorithm = "algorithm"
	// PropDistance is the Levenshtein distance between origin and variant.
	PropDistance = "distance"
)

// Node type names this package knows about. They mirror the type registry in
// docs/DESIGN.md §2; every one of them is Nameable.
const (
	TypeDomain   = "domain"
	TypeEmail    = "email"
	TypePackage  = "package"
	TypeRepo     = "repo"
	TypeUsername = "username"
)

// NameableTypes is every type a capability-bound variant operator may emit.
// Effects.Nodes wants concrete type names for plan compilation, but the trigger
// binds by capability, so the honest answer is "any Nameable type" — listed
// here as the *may* that Emits is defined to be (docs/DESIGN.md §5).
var NameableTypes = []string{TypeDomain, TypeEmail, TypePackage, TypeRepo, TypeUsername}

// RelDef is the VARIANT_OF relation as this package needs it registered.
//
// It exists so a caller with its own registry — a test, or a tool that loads
// only the variant operators — can register the relation without restating the
// field list, since an operator asserting a prop the relation never declared is
// rejected at apply time. Where the shared operator schema is installed instead
// (internal/operators/decompose.Register), that is authoritative and this must
// stay identical to it; a second registration of the same name is an error, not
// a merge.
func RelDef() graph.RelDef {
	return graph.RelDef{
		Name:    Rel,
		Class:   graph.Variant,
		Version: 1,
		Fields: []graph.FieldDef{
			{Name: PropAlgorithm, Kind: graph.KindString},
			{Name: PropDistance, Kind: graph.KindInt},
		},
	}
}

// errNoKey is returned when the matched node has no key at all, which means the
// view was built against a node that is not in the graph.
var errNoKey = errors.New("variant: matched node has no key")

// Generate turns one name into its variations. Implementations must be pure —
// same input, same output — because the scheduler caches an operator's result
// against a digest of its declared reads and will not call it twice for the
// same input.
//
// Order is deliberately *not* part of the contract: several of the underlying
// pkg/typo functions accumulate through a map and so return their results in
// Go's randomized map order. The adapter sorts, which is what makes the emitted
// delta byte-identical across runs regardless.
type Generate func(name string) []string

// Spec declares one algorithm. Everything an operator needs beyond the shared
// adapter lives here, so adding an algorithm is a literal rather than a type.
type Spec struct {
	// ID is the operator id. It is also the value of the edge's algorithm prop,
	// and it is the short code the CLI has always used ("co", "cs", ...).
	ID string
	// Title is the human-readable name, for --explain and reports.
	Title string
	// Version invalidates cached results; bump it when Gen's output changes.
	Version int
	// Types narrows the trigger to specific node types. Empty means bind by
	// capability, which is the default and the point of the design.
	Types []string
	// Whole varies the entire key instead of just the registrable name. TLD
	// swapping and namespace confusion need the whole thing; character omission
	// must not touch the suffix or it would silently generate a different TLD.
	Whole bool
	// Gen produces the variations.
	Gen Generate
}

// generator adapts a Spec to graph.Operator. One implementation covers every
// algorithm; the algorithms themselves are data.
type generator struct {
	spec  Spec
	split Splitter
}

// New builds an operator from a spec. split decomposes a key into the part the
// algorithm may vary and the part it must preserve; pass nil for DefaultSplit.
func New(s Spec, split Splitter) graph.Operator {
	if split == nil {
		split = DefaultSplit
	}
	if s.Version < 1 {
		s.Version = 1
	}
	return &generator{spec: s, split: split}
}

// Id implements graph.Operator.
func (g *generator) Id() string { return g.spec.ID }

// Name is the human-readable title. It is not part of graph.Operator, but
// plan rendering and reports want it.
func (g *generator) Name() string { return g.spec.Title }

// Version implements graph.Operator.
func (g *generator) Version() int { return g.spec.Version }

// Resource implements graph.Operator. Variant generation is pure computation
// and belongs to no rate-limit class; naming one would throttle CPU work behind
// a network budget.
func (g *generator) Resource() string { return "" }

// Trigger implements graph.Operator.
//
// The selector binds by capability so one algorithm covers every Nameable type.
// InClosure is an optimization only: the applier rejects an out-of-closure
// VARIANT_OF edge regardless, and the scheduler already skips variant operators
// outside the closure. Declaring it here keeps the operator honest about what
// it assumes rather than relying on two other layers to be correct.
//
// Reads is empty because the algorithms are functions of the node's key alone,
// which the View always provides. An empty read-set means the digest never
// changes, so an algorithm runs exactly once per node — which is right: nothing
// a collector later learns can change what a string mutation produces.
func (g *generator) Trigger() graph.Trigger {
	t := graph.Trigger{
		Where: []graph.Condition{graph.InClosure()},
	}
	if len(g.spec.Types) > 0 {
		t.On = graph.Selector{Types: append([]string(nil), g.spec.Types...)}
	} else {
		t.On = graph.Selector{Caps: []graph.Capability{graph.Nameable}}
	}
	return t
}

// Emits implements graph.Operator. Declaring VARIANT_OF is what makes this a
// variant operator as far as the engine is concerned.
func (g *generator) Emits() graph.Effects {
	nodes := NameableTypes
	if len(g.spec.Types) > 0 {
		nodes = g.spec.Types
	}
	return graph.Effects{
		Nodes: append([]string(nil), nodes...),
		Rels:  []string{Rel},
		Props: []string{PropAlgorithm, PropDistance},
	}
}

// Exec generates the variations of the matched node's key and returns them as
// VARIANT_OF edges from the origin to each variant.
//
// Nodes are named by NodeRef carrying the *raw* key: canonicalization belongs
// to the registry and the applier, and an operator that minted a NodeID would
// quietly break the convergence the whole design rests on.
//
// The delta lists no bare nodes. Every variant is introduced by the edge that
// justifies it, so a variant edge the applier refuses — an origin outside the
// seed closure — cannot leave an orphan variant node behind.
func (g *generator) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	key := v.Key()
	if key == "" {
		return graph.Delta{}, graph.Failed(errNoKey)
	}

	origin := v.Ref()
	parts := g.split(v.Type(), key)

	source := parts.Core
	if g.spec.Whole {
		source = key
	}
	if source == "" {
		return graph.Delta{}, graph.Empty()
	}

	// Dedupe against the origin as well as against each other: two positions in
	// a name often collapse to the same variation, and a "variant" equal to its
	// own origin would be a self-edge.
	seen := map[string]bool{key: true}
	keys := make([]string, 0, 16)
	for _, raw := range g.spec.Gen(source) {
		if raw == "" || raw == source {
			continue
		}
		variant := raw
		if !g.spec.Whole {
			variant = parts.Join(raw)
		}
		if variant == "" || seen[variant] {
			continue
		}
		seen[variant] = true
		keys = append(keys, variant)
	}
	if len(keys) == 0 {
		// Not a failure: this algorithm authoritatively has nothing to say
		// about this name. A name with no hyphens has no hyphen omissions.
		return graph.Delta{}, graph.Empty()
	}

	// Sorting is what makes two runs byte-identical. Several pkg/typo functions
	// return results in map order, and the delta must not inherit that.
	sort.Strings(keys)

	d := graph.Delta{
		Edges: make([]graph.EdgeRef, 0, len(keys)),
		Props: make([]graph.PropSet, 0, 2*len(keys)),
	}
	for _, k := range keys {
		ref := graph.EdgeRef{
			From: origin,
			Rel:  Rel,
			To:   graph.NodeRef{Type: v.Type(), Key: k},
		}
		edge := ref // addressed per iteration; sharing one address would alias
		d.Edges = append(d.Edges, ref)
		d.Props = append(d.Props,
			graph.PropSet{Edge: &edge, Field: PropAlgorithm, Value: graph.String(g.spec.ID)},
			// Distance is measured between the keys as emitted. The applier
			// canonicalizes afterwards, so a variant whose canonical form
			// differs from its raw form carries the pre-canonical distance.
			graph.PropSet{Edge: &edge, Field: PropDistance, Value: graph.Int(int64(fuzzy.Levenshtein(key, k)))},
		)
	}
	return d, graph.OK()
}
