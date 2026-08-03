// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package train fits an execution model from recorded scans.
//
// The engine's belief model is a hidden Markov model over the expansion tree:
// a node's belief is its parent's latent state pushed through the relation that
// admitted it and conditioned on the node's own props. internal/model is that
// HMM and can already be trained; internal/graph is where the trees come from;
// nothing joined them, so every scan has run the uniform model, which returns 1
// for everything and reduces expansion to unranked breadth-first.
//
// This package is the join. It featurizes a graph, walks it into the
// root-to-leaf paths Baum-Welch is defined over, and fits a model that
// graph.SetBeliefModel can take.
//
// # Train and serve must featurize identically
//
// The one hazard worth naming. Training reads a finished graph; inference reads
// a graph.View mid-scan, restricted to what an operator declared. Two code
// paths, two shapes, one meaning — and if they ever disagree the model is fitted
// on symbols that never occur at run time and silently degrades to its prior.
// Nothing fails, no test goes red, and the scan simply stops ranking.
//
// So there is one feature function, over the narrowest interface both can
// satisfy: Featurable. It is the only place a symbol is named.
package train

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// Featurable is the least a node has to expose to be featurized.
//
// graph.View satisfies it as it stands, and the training-time adapter below is
// written to match rather than the other way round: the run-time surface is the
// constrained one, so anything the trainer can see but an operator cannot is a
// feature that could not be computed during a scan.
type Featurable interface {
	Type() string
	Key() string
	Depth() int
	Prop(field string) (graph.Value, bool)
	// EdgeProps returns the props of each outgoing edge of a relation.
	//
	// Props rather than graph.EdgeView because an EdgeView cannot be built
	// outside internal/graph, and this package adapts to the engine rather than
	// asking the engine to widen for it. Bag is what both sides can produce.
	EdgeProps(rel string) []Bag
}

// Bag is anything with props: a graph.EdgeView satisfies it as it stands.
type Bag interface {
	Prop(field string) (graph.Value, bool)
}

// Fields and Rels are what the featurizer reads.
//
// They are declared, not discovered, because an operator's trigger has to name
// its read set for the digest to be sound. A model that read a field the
// engine did not declare would be served stale views forever.
// "live" is deliberately absent. It is declared on the domain schema
// (decompose.FieldLive) and no production operator sets it, so the symbol never
// fired — but the reason not to add it is stronger than that: it is the label.
// Existence is what the model is asked to anticipate, so feeding it in as a
// feature would make belief a restatement of the answer rather than an estimate
// of it, and the resulting score would look excellent and mean nothing. The
// observation *relations* below have a weaker version of the same problem, which
// is documented on Evaluate.
//
// "rank" is also set by nothing today, so it never fires either. It is kept
// because it is legitimate signal if it ever arrives — popularity is known
// before any lookup runs — and an unset field costs one map miss.
var (
	Fields = []string{"punycode", "created", "rank"}
	Rels   = []string{graph.VariantRel, "RESOLVES_TO", "NS", "MX", "REGISTERED_BY", "EXISTS_ON"}
)

// Features turns a node into emission symbols.
//
// Symbols are coarse on purpose. The alphabet is the model's identity — it is
// part of the artifact and part of its CID — so a symbol per edit distance or
// per registrar would give an alphabet that grows with the corpus and a model
// that cannot be compared across runs. Everything here is a bucket chosen once.
//
// Sorted, because the emission alphabet is a set and two orderings of the same
// observation must not train as two observations.
func Features(n Featurable) []string {
	var out []string
	add := func(f string, args ...any) {
		if len(args) > 0 {
			f = fmt.Sprintf(f, args...)
		}
		out = append(out, f)
	}

	add("type:%s", n.Type())

	// Depth is bucketed: the difference between the seed and one hop is the
	// whole signal, and between hop nine and ten there is none.
	switch d := n.Depth(); {
	case d == 0:
		add("depth:0")
	case d == 1:
		add("depth:1")
	case d <= 3:
		add("depth:2-3")
	default:
		add("depth:deep")
	}

	// How the name is shaped. A model has no characters to look at — it sees
	// symbols — so anything about the string has to be turned into one here.
	key := n.Key()
	switch l := len([]rune(key)); {
	case l <= 8:
		add("len:short")
	case l <= 16:
		add("len:mid")
	default:
		add("len:long")
	}
	if strings.Contains(key, "-") {
		add("has:hyphen")
	}
	if strings.ContainsAny(key, "0123456789") {
		add("has:digit")
	}
	if !isASCII(key) {
		add("has:nonascii")
	}

	// Props the observation operators set. Not "live" — see Fields.
	if _, ok := n.Prop("punycode"); ok {
		add("has:punycode")
	}
	if _, ok := n.Prop("created"); ok {
		add("has:created")
	}

	// What the node is attached to. Presence, not count: "resolves at all" is
	// the signal, and a symbol per address count would be alphabet growth for
	// no information.
	for _, rel := range Rels {
		if len(n.EdgeProps(rel)) > 0 {
			add("edge:%s", rel)
		}
	}

	// The algorithm that produced a variant is NOT here, and it is the feature
	// this model most wants.
	//
	// It lives on the VARIANT_OF edge, which runs origin -> variant, so on the
	// variant it is an *incoming* edge — and graph.View.Edges returns outgoing
	// edges only. An operator mid-scan therefore cannot see which algorithm
	// produced the node it is looking at, so neither may this featurizer:
	// training on a symbol inference can never observe fits the model to
	// noise and it degrades to its prior in production, silently.
	//
	// Two things soften it. The relation itself is the HMM's transition label,
	// so "arrived by VARIANT_OF" is modelled even though "arrived by co" is
	// not. And the gap is one accessor wide: an incoming-edge view, declared in
	// a trigger like any other read, would make algo: and dist: available to
	// both sides at once. Until then the honest alphabet is the one below.

	sort.Strings(out)
	return dedupe(out)
}

// viewFeaturable adapts a run-time graph.View to Featurable.
type viewFeaturable struct{ v graph.View }

func (f viewFeaturable) Type() string { return f.v.Type() }
func (f viewFeaturable) Key() string  { return f.v.Key() }
func (f viewFeaturable) Depth() int   { return f.v.Depth() }

func (f viewFeaturable) Prop(field string) (graph.Value, bool) { return f.v.Prop(field) }

func (f viewFeaturable) EdgeProps(rel string) []Bag {
	edges := f.v.Edges(rel)
	out := make([]Bag, 0, len(edges))
	for _, e := range edges {
		out = append(out, e)
	}
	return out
}

// Featurizer is Features as internal/model wants it, for inference.
//
// This is the whole of the run-time side: the same function, over the same
// interface, reached through a graph.View. Nothing else may featurize.
func Featurizer(v graph.View) []string { return Features(viewFeaturable{v}) }

// Trigger is what an operator or the engine must declare to be served the
// fields Features reads.
//
// Handing back the declaration alongside the featurizer is what keeps them
// together: a caller that reads Fields without declaring them gets a View that
// reports every one of them unset, and the model quietly sees nothing.
func Trigger() graph.Reads {
	return graph.Reads{
		Fields: append([]string(nil), Fields...),
		Rels:   append([]string(nil), Rels...),
	}
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

func dedupe(sorted []string) []string {
	if len(sorted) < 2 {
		return sorted
	}
	out := sorted[:1]
	for _, s := range sorted[1:] {
		if s != out[len(out)-1] {
			out = append(out, s)
		}
	}
	return out
}
