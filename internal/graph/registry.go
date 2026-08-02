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

// Package graph implements the typed property graph the scan engine expands:
// content-addressed nodes, typed relations, positional props and the single
// writer that applies operator deltas. See docs/DESIGN.md.
package graph

import (
	"fmt"
	"sort"
)

// Capability classifies what may be done to a node type. Nameable is necessary
// but not sufficient for variant expansion — eligibility also requires
// seed-closure membership, which the applier enforces.
type Capability uint8

const (
	// Nameable types can be the root of variant generation.
	Nameable Capability = iota + 1
	// Observed types are only ever discovered, never varied.
	Observed
)

func (c Capability) String() string {
	switch c {
	case Nameable:
		return "nameable"
	case Observed:
		return "observed"
	}
	return "invalid"
}

// Class is a relation's edge class. It decides depth accounting and whether an
// edge extends the seed closure. The dividing line is whether producing the
// edge required a network call.
type Class uint8

const (
	// Structural edges come from parsing the target string alone.
	Structural Class = iota + 1
	// Variant edges connect an origin to a generated variation.
	Variant
	// Observation edges required a lookup against some external service.
	Observation
)

func (c Class) String() string {
	switch c {
	case Structural:
		return "structural"
	case Variant:
		return "variant"
	case Observation:
		return "observation"
	}
	return "invalid"
}

// DepthCost is what traversing this class adds to a node's depth. Only
// observation hops count; structural and variant edges are free, so a composite
// seed does not spend its depth budget on decomposition.
func (c Class) DepthCost() int {
	if c == Observation {
		return 1
	}
	return 0
}

// Canonical normalizes a raw key into the single form the graph converges on.
// Returning an error refuses the candidate outright.
type Canonical func(string) (string, error)

// NodeTypeDef declares a node type. Version rises only when Fields is appended
// to or a field is tombstoned.
type NodeTypeDef struct {
	Name      string
	Cap       Capability
	Version   int
	Canonical Canonical
	Fields    []FieldDef
}

// RelDef declares a relation type. Relations carry their own ordered field
// list, encoded exactly like a node's.
type RelDef struct {
	Name    string
	Class   Class
	Version int
	Fields  []FieldDef
}

// NodeType is a registered node type handle.
type NodeType struct {
	name      string
	cap       Capability
	canonical Canonical
	sch       *schema
}

func (t *NodeType) Name() string                 { return t.name }
func (t *NodeType) Cap() Capability              { return t.cap }
func (t *NodeType) Version() int                 { return t.sch.version }
func (t *NodeType) Field(n string) (Field, bool) { return t.sch.Field(n) }

// Rel is a registered relation handle.
type Rel struct {
	name  string
	class Class
	sch   *schema
}

func (r *Rel) Name() string                 { return r.name }
func (r *Rel) Class() Class                 { return r.class }
func (r *Rel) Version() int                 { return r.sch.version }
func (r *Rel) Field(n string) (Field, bool) { return r.sch.Field(n) }

// Registry holds every registered node and relation type. It is not safe for
// concurrent registration; register during init, then read.
type Registry struct {
	types map[string]*NodeType
	rels  map[string]*Rel
}

func NewRegistry() *Registry {
	return &Registry{types: map[string]*NodeType{}, rels: map[string]*Rel{}}
}

// AddType registers a node type.
func (r *Registry) AddType(d NodeTypeDef) (*NodeType, error) {
	if d.Name == "" {
		return nil, fmt.Errorf("graph: node type has no name")
	}
	if _, dup := r.types[d.Name]; dup {
		return nil, fmt.Errorf("graph: node type %q already registered", d.Name)
	}
	if d.Cap != Nameable && d.Cap != Observed {
		return nil, fmt.Errorf("graph: node type %q has invalid capability", d.Name)
	}
	if d.Canonical == nil {
		return nil, fmt.Errorf("graph: node type %q needs a canonicalization function", d.Name)
	}
	if d.Version < 1 {
		return nil, fmt.Errorf("graph: node type %q needs version >= 1", d.Name)
	}
	sch, err := newSchema(d.Name, d.Version, d.Fields)
	if err != nil {
		return nil, err
	}
	t := &NodeType{name: d.Name, cap: d.Cap, canonical: d.Canonical, sch: sch}
	r.types[d.Name] = t
	return t, nil
}

// AddRel registers a relation type.
func (r *Registry) AddRel(d RelDef) (*Rel, error) {
	if d.Name == "" {
		return nil, fmt.Errorf("graph: relation has no name")
	}
	if _, dup := r.rels[d.Name]; dup {
		return nil, fmt.Errorf("graph: relation %q already registered", d.Name)
	}
	if d.Class < Structural || d.Class > Observation {
		return nil, fmt.Errorf("graph: relation %q has invalid class", d.Name)
	}
	if d.Version < 1 {
		return nil, fmt.Errorf("graph: relation %q needs version >= 1", d.Name)
	}
	sch, err := newSchema(d.Name, d.Version, d.Fields)
	if err != nil {
		return nil, err
	}
	rel := &Rel{name: d.Name, class: d.Class, sch: sch}
	r.rels[d.Name] = rel
	return rel, nil
}

// Type looks up a registered node type.
func (r *Registry) Type(name string) (*NodeType, bool) {
	t, ok := r.types[name]
	return t, ok
}

// Rel looks up a registered relation.
func (r *Registry) Rel(name string) (*Rel, bool) {
	rel, ok := r.rels[name]
	return rel, ok
}

// Types returns every registered node type, ordered by name.
func (r *Registry) Types() []*NodeType {
	out := make([]*NodeType, 0, len(r.types))
	for _, t := range r.types {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

// Rels returns every registered relation, ordered by name.
func (r *Registry) Rels() []*Rel {
	out := make([]*Rel, 0, len(r.rels))
	for _, rel := range r.rels {
		out = append(out, rel)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}
