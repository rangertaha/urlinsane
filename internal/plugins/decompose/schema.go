// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package decompose

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// The node types of docs/DESIGN.md §2. That table is authoritative: a thing is
// a node only if it has identity worth converging on across the graph, which is
// why a nameserver is a domain reached by an NS edge rather than a type, and
// why a registration date is a prop rather than a node.
const (
	TypeDomain     = "domain"
	TypeUsername   = "username"
	TypePackage    = "package"
	TypeRepo       = "repo"
	TypeEmail      = "email"
	TypeTLD        = "tld"
	TypeIP         = "ip"
	TypeASN        = "asn"
	TypeRegistrant = "registrant"
	TypePlatform   = "platform"
)

// Structural relations (§1.1). Each is derived from parsing the target string
// alone, so it costs no depth and carries seed-closure membership across.
// Direction follows §3: an edge runs from the composite to the part.
const (
	RelLocalPart = "LOCAL_PART"
	RelDomainOf  = "DOMAIN_OF"
	RelTLDOf     = "TLD_OF"
	RelOwner     = "OWNER"
	RelHostedOn  = "HOSTED_ON"
)

// RelVariantOf marks a generated variation. The name is taken from the graph
// package rather than restated: declaring it in Effects is what subjects an
// operator to the terminal rule and the seed-closure restriction, and the
// scheduler recognises it by exactly this string.
const RelVariantOf = graph.VariantRel

// Observation relations (§1.1). The line between these and the structural set
// is whether producing the edge required a network call — which is why MANIFEST
// lives here despite looking like decomposition, and why the packages it yields
// fall outside the seed closure and are never varied.
const (
	RelManifest     = "MANIFEST"
	RelResolvesTo   = "RESOLVES_TO"
	RelNS           = "NS"
	RelMX           = "MX"
	RelPTRTo        = "PTR_TO"
	RelInASN        = "IN_ASN"
	RelRegisteredBy = "REGISTERED_BY"
	RelExistsOn     = "EXISTS_ON"
)

// Field names. Fields are append-only: a position is part of every content
// address already in the store, so reordering or deleting one corrupts diffs
// rather than failing loudly (§1.3).
const (
	FieldPunycode  = "punycode"
	FieldLive      = "live"
	FieldCreated   = "created"
	FieldRank      = "rank"
	FieldAlgorithm = "algorithm"
	FieldDistance  = "distance"
)

// Extension is one package's addition to the schema: the props its operators
// assert that the base declaration does not already name, keyed by node type
// and by relation.
//
// It exists because a package cannot register its own fields. A field lives on
// the *type*, and a type has exactly one registration — so an operator package
// that declared its own would either duplicate the type (two "domain"s with
// different canonicalization, which is how convergence breaks) or import the
// registering package and be imported by it in turn.
//
// The alternative — leaving each package to register what it asserts — fails
// quietly rather than loudly: the applier records an unknown-field rejection
// and carries on, so a collector's entire output vanishes with the run still
// reporting success.
type Extension struct {
	Fields    map[string][]graph.FieldDef // by node type
	RelFields map[string][]graph.FieldDef // by relation
}

// Register installs every node type and relation of docs/DESIGN.md §2 and §1.1
// into r, extended by each Extension in argument order. It is the one
// definition of the schema; operators refer to it by the constants above rather
// than declaring types of their own.
//
// Extension fields are appended after the base fields and after each other, in
// the order given. That order is load-bearing: a field's position is its stable
// index and part of the content address of every node already in the store
// (§1.3), so appending is safe and reordering silently corrupts diffs. Callers
// must therefore pass extensions in a fixed order — not from a map, and not
// from a set the plan happens to select.
//
// Registration is not concurrency-safe and a duplicate name is an error, so
// call this once against a fresh registry.
func Register(r *graph.Registry, extra ...Extension) error {
	types := nodeTypes()
	rels := relations()

	// Validated before anything is installed, so a bad extension leaves the
	// registry untouched rather than half-built.
	if err := unknownTargets(types, rels, extra); err != nil {
		return err
	}

	for _, e := range extra {
		for i := range types {
			types[i].Fields = append(types[i].Fields, e.Fields[types[i].Name]...)
		}
		for i := range rels {
			rels[i].Fields = append(rels[i].Fields, e.RelFields[rels[i].Name]...)
		}
	}

	for _, d := range types {
		if _, err := r.AddType(d); err != nil {
			return err
		}
	}
	for _, d := range rels {
		if _, err := r.AddRel(d); err != nil {
			return err
		}
	}
	return nil
}

// unknownTargets reports an extension field aimed at a type or relation that
// does not exist. Silently dropping it would reproduce exactly the failure
// Extension exists to prevent, one level up: the field would never be
// registered, and every assertion against it would be refused at runtime.
func unknownTargets(types []graph.NodeTypeDef, rels []graph.RelDef, extra []Extension) error {
	known := map[string]bool{}
	for _, t := range types {
		known[t.Name] = true
	}
	knownRel := map[string]bool{}
	for _, r := range rels {
		knownRel[r.Name] = true
	}
	for _, e := range extra {
		for name := range e.Fields {
			if !known[name] {
				return fmt.Errorf("decompose: extension declares fields for unknown node type %q", name)
			}
		}
		for name := range e.RelFields {
			if !knownRel[name] {
				return fmt.Errorf("decompose: extension declares fields for unknown relation %q", name)
			}
		}
	}
	return nil
}

// nodeTypes is §2's table. Nameable means variant operators *may* apply —
// necessary but not sufficient, since eligibility also requires seed-closure
// membership, which the applier enforces independently.
//
// Field lists start at the minimum the design names. They grow by appending as
// collectors land; the append-only rule makes that safe and makes inventing
// speculative fields now the more expensive option.
func nodeTypes() []graph.NodeTypeDef {
	return []graph.NodeTypeDef{
		{
			Name: TypeDomain, Cap: graph.Nameable, Version: 1, Canonical: canonDomain,
			// Declaration order is presentation order (§1.3): identity first,
			// then reachability, then registration facts.
			Fields: []graph.FieldDef{
				{Name: FieldPunycode, Kind: graph.KindString},
				{Name: FieldLive, Kind: graph.KindBool},
				// Two sources assert this, so the winner is declared rather
				// than decided by whichever answered first (§1.4).
				{Name: FieldCreated, Kind: graph.KindTime, Merge: graph.Precedence("rdap", "whois")},
				{Name: FieldRank, Kind: graph.KindInt},
			},
		},
		{Name: TypeUsername, Cap: graph.Nameable, Version: 1, Canonical: canonUsername},
		{Name: TypePackage, Cap: graph.Nameable, Version: 1, Canonical: canonPackage},
		{Name: TypeRepo, Cap: graph.Nameable, Version: 1, Canonical: canonRepo},
		{Name: TypeEmail, Cap: graph.Nameable, Version: 1, Canonical: canonEmail},
		{Name: TypeTLD, Cap: graph.Observed, Version: 1, Canonical: canonTLD},
		{Name: TypeIP, Cap: graph.Observed, Version: 1, Canonical: canonIP},
		{Name: TypeASN, Cap: graph.Observed, Version: 1, Canonical: canonASN},
		{Name: TypeRegistrant, Cap: graph.Observed, Version: 1, Canonical: canonRegistrant},
		{Name: TypePlatform, Cap: graph.Observed, Version: 1, Canonical: canonPlatform},
	}
}

// relations is §1.1's table. A relation's class is the single place depth
// accounting and closure growth are decided, so getting one wrong here silently
// changes termination everywhere.
func relations() []graph.RelDef {
	return []graph.RelDef{
		{Name: RelLocalPart, Class: graph.Structural, Version: 1},
		{Name: RelDomainOf, Class: graph.Structural, Version: 1},
		{Name: RelTLDOf, Class: graph.Structural, Version: 1},
		{Name: RelOwner, Class: graph.Structural, Version: 1},
		{Name: RelHostedOn, Class: graph.Structural, Version: 1},

		// Relation props carry data operators need: the per-node analyzers read
		// the algorithm and edit distance straight off the edge, which bare
		// neighbour nodes would hide (§4.2).
		{Name: RelVariantOf, Class: graph.Variant, Version: 1, Fields: []graph.FieldDef{
			{Name: FieldAlgorithm, Kind: graph.KindString},
			{Name: FieldDistance, Kind: graph.KindInt},
		}},

		{Name: RelManifest, Class: graph.Observation, Version: 1},
		{Name: RelResolvesTo, Class: graph.Observation, Version: 1},
		{Name: RelNS, Class: graph.Observation, Version: 1},
		{Name: RelMX, Class: graph.Observation, Version: 1},
		{Name: RelPTRTo, Class: graph.Observation, Version: 1},
		{Name: RelInASN, Class: graph.Observation, Version: 1},
		{Name: RelRegisteredBy, Class: graph.Observation, Version: 1},
		{Name: RelExistsOn, Class: graph.Observation, Version: 1},
	}
}
