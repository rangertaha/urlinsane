// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

// View is what an operator sees: the matched node, and exactly the props and
// relations its trigger declared it reads. Both are filtered, not just
// relations — the read-set digest is built from the same declaration, so a
// field an operator could read without declaring could change without changing
// the digest, and the operator would be served a stale cached result forever.
type View interface {
	// ID is the matched node's identity.
	ID() NodeID
	// Type is the matched node's type name.
	Type() string
	// Key is the matched node's canonical key.
	Key() string
	// Depth is the node's observation distance from the seed.
	Depth() int
	// Prop returns a declared field's value, and whether it is set. Reading an
	// undeclared field always reports unset.
	Prop(field string) (Value, bool)
	// Edges returns outgoing edges of a declared relation, with their props.
	// Reading an undeclared relation always returns nothing.
	Edges(rel string) []EdgeView
	// Ref names this node for use in a Delta.
	Ref() NodeRef
}

// EdgeView is one edge as an operator sees it: the far node and the edge's own
// props. Relation props carry data operators need — VARIANT_OF holds the
// algorithm and edit distance — which bare neighbour nodes would hide.
type EdgeView struct {
	Rel  string
	To   NodeRef
	ToID NodeID
	prop Props
}

// Prop returns an edge prop by name.
func (e EdgeView) Prop(field string) (Value, bool) {
	if e.prop.sch == nil {
		return Value{}, false
	}
	f, ok := e.prop.sch.Field(field)
	if !ok {
		return Value{}, false
	}
	return e.prop.Get(f)
}

// view is the concrete, pattern-scoped View handed to one operator call.
type view struct {
	g      *Graph
	id     NodeID
	fields map[string]bool
	rels   map[string]bool
}

func (g *Graph) viewFor(t Trigger, id NodeID) View {
	r := t.effectiveReads()
	v := &view{g: g, id: id, fields: map[string]bool{}, rels: map[string]bool{}}
	for _, f := range r.Fields {
		v.fields[f] = true
	}
	for _, rel := range r.Rels {
		v.rels[rel] = true
	}
	return v
}

func (v *view) ID() NodeID { return v.id }
func (v *view) Type() string {
	n, ok := v.g.nodes[v.id]
	if !ok {
		return ""
	}
	return n.Type.name
}
func (v *view) Key() string {
	n, ok := v.g.nodes[v.id]
	if !ok {
		return ""
	}
	return n.Key
}
func (v *view) Depth() int { return v.g.depth[v.id] }

func (v *view) Prop(field string) (Value, bool) {
	if !v.fields[field] {
		return Value{}, false
	}
	n, ok := v.g.nodes[v.id]
	if !ok {
		return Value{}, false
	}
	f, ok := n.Type.sch.Field(field)
	if !ok {
		return Value{}, false
	}
	return n.Props.Get(f)
}

func (v *view) Edges(rel string) []EdgeView {
	if !v.rels[rel] {
		return nil
	}
	var out []EdgeView
	for _, e := range v.g.outEdges(v.id, rel) {
		to := v.g.nodes[e.To]
		out = append(out, EdgeView{
			Rel:  rel,
			To:   NodeRef{Type: to.Type.name, Key: to.Key},
			ToID: e.To,
			prop: e.Props,
		})
	}
	return out
}

func (v *view) Ref() NodeRef {
	n, ok := v.g.nodes[v.id]
	if !ok {
		return NodeRef{}
	}
	return NodeRef{Type: n.Type.name, Key: n.Key}
}
