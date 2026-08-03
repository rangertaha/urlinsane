// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

// Selector chooses which nodes an operator binds to, by type or by capability.
// Binding by capability is what lets one omission algorithm cover every
// Nameable type instead of being registered once per type.
type Selector struct {
	Types []string
	Caps  []Capability
}

func (s Selector) matches(t *NodeType) bool {
	for _, n := range s.Types {
		if n == t.name {
			return true
		}
	}
	for _, c := range s.Caps {
		if c == t.cap {
			return true
		}
	}
	return false
}

// Condition is an extra requirement on the matched node. Conditions are data
// conditions, never producer dependencies: "there is an IP", not "the ip
// operator has run".
type Condition interface {
	match(g *Graph, id NodeID) bool
	// declares reports what the condition itself reads, so a trigger's read-set
	// is complete without the operator having to restate it.
	declares() (fields []string, rels []string)
	describe() string
}

// HasProp requires a field to be set on the matched node.
func HasProp(field string) Condition { return hasProp{field} }

type hasProp struct{ field string }

func (c hasProp) match(g *Graph, id NodeID) bool {
	n, ok := g.nodes[id]
	if !ok {
		return false
	}
	f, ok := n.Type.sch.Field(c.field)
	if !ok {
		return false
	}
	_, set := n.Props.Get(f)
	return set
}
func (c hasProp) declares() ([]string, []string) { return []string{c.field}, nil }
func (c hasProp) describe() string               { return "has-prop:" + c.field }

// HasEdge requires at least one outgoing edge of a relation.
func HasEdge(rel string) Condition { return hasEdge{rel} }

type hasEdge struct{ rel string }

func (c hasEdge) match(g *Graph, id NodeID) bool { return len(g.outEdges(id, c.rel)) > 0 }
func (c hasEdge) declares() ([]string, []string) { return nil, []string{c.rel} }
func (c hasEdge) describe() string               { return "has-edge:" + c.rel }

// InClosure requires seed-closure membership. The engine already refuses
// out-of-closure variant edges at the applier; this lets an operator avoid the
// wasted dispatch as well.
func InClosure() Condition { return inClosure{} }

type inClosure struct{}

func (inClosure) match(g *Graph, id NodeID) bool { return g.closure[id] }
func (inClosure) declares() ([]string, []string) { return nil, nil }
func (inClosure) describe() string               { return "in-seed-closure" }

// BeliefAbove requires the execution model's belief to clear a threshold. It is
// evaluated only at a barrier, never during delta-driven re-dispatch, so it
// takes no part in the read-set digest.
func BeliefAbove(t float64) Condition { return beliefAbove{t} }

type beliefAbove struct{ threshold float64 }

func (c beliefAbove) match(g *Graph, id NodeID) bool { return g.Belief(id) > c.threshold }
func (c beliefAbove) declares() ([]string, []string) { return nil, nil }
func (c beliefAbove) describe() string               { return "belief-above" }

// isBelief reports whether a condition is barrier-time only.
func isBelief(c Condition) bool { _, ok := c.(beliefAbove); return ok }

// Reads declares the props and relations an operator consumes. It does double
// duty: it scopes the View, and it is the input to the read-set digest that
// decides re-dispatch and cache validity.
type Reads struct {
	Fields []string
	Rels   []string
}

// Trigger is when an operator runs.
type Trigger struct {
	On    Selector
	Where []Condition
	Reads Reads
}

// effectiveReads merges the operator's declared reads with whatever its
// conditions inspect. A condition that reads a field the operator forgot to
// declare would otherwise leave that field out of the digest, and the operator
// would never re-run when it changed.
func (t Trigger) effectiveReads() Reads {
	fset := map[string]bool{}
	rset := map[string]bool{}
	for _, f := range t.Reads.Fields {
		fset[f] = true
	}
	for _, r := range t.Reads.Rels {
		rset[r] = true
	}
	for _, c := range t.Where {
		fs, rs := c.declares()
		for _, f := range fs {
			fset[f] = true
		}
		for _, r := range rs {
			rset[r] = true
		}
	}
	return Reads{Fields: sortedKeys(fset), Rels: sortedKeys(rset)}
}

func sortedKeys(m map[string]bool) []string {
	if len(m) == 0 {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Effects declares everything an operator may produce. It covers relations and
// props, not just node types, so plan compilation can see a prop-only operator
// and can detect a Where nothing in the plan will ever satisfy.
type Effects struct {
	Nodes []string
	Rels  []string
	Props []string
}

// Operator is a unit of work bound to a pattern in the graph.
type Operator interface {
	Id() string
	Version() int
	Trigger() Trigger
	Emits() Effects
	// Resource names the rate-limit class this operator's calls belong to.
	Resource() string
	// Exec does the work. The Outcome is the operator's own judgement of what
	// happened — an authoritative absence is Empty, not Failed — because how a
	// lookup failed is itself the finding.
	//
	// ctx carries the round deadline and the interrupt. An operator that makes
	// network calls must pass it down: without it the scheduler can only cancel
	// *between* attempts, so Ctrl-C waits for an in-flight whois rather than
	// stopping at the round boundary (§6.2, §12.4). Pure operators may ignore it.
	Exec(ctx context.Context, v View) (Delta, Outcome)
}

// matches reports whether an operator's trigger matches a node right now.
// barrier selects whether belief conditions participate: they are evaluated at
// barriers only.
func (g *Graph) matches(t Trigger, id NodeID, barrier bool) bool {
	n, ok := g.nodes[id]
	if !ok {
		return false
	}
	if !t.On.matches(n.Type) {
		return false
	}
	for _, c := range t.Where {
		if isBelief(c) && !barrier {
			continue
		}
		if !c.match(g, id) {
			return false
		}
	}
	return true
}

// readDigest hashes exactly what an operator declared it reads. Anything not
// declared is invisible to it (the View enforces that), so anything not
// declared must not affect re-dispatch or the cache key either.
func (g *Graph) readDigest(t Trigger, id NodeID) [32]byte {
	r := t.effectiveReads()
	h := sha256.New()
	n := g.nodes[id]
	for _, name := range r.Fields {
		writeField(h, name)
		f, ok := n.Type.sch.Field(name)
		if !ok {
			writeField(h, "\x00unknown")
			continue
		}
		v, set := n.Props.Get(f)
		if !set {
			writeField(h, "\x00unset")
			continue
		}
		writeValue(h, v)
	}
	for _, rel := range r.Rels {
		writeField(h, rel)

		// Edges sorted by id, so the digest does not inherit adjacency order.
		edges := append([]*Edge(nil), g.outEdges(id, rel)...)
		sort.Slice(edges, func(i, j int) bool {
			return compareID(edges[i].ID[:], edges[j].ID[:]) < 0
		})

		var count [8]byte
		binary.BigEndian.PutUint64(count[:], uint64(len(edges)))
		_, _ = h.Write(count[:])

		for _, e := range edges {
			_, _ = h.Write(e.ID[:])

			// The edge's props as well as its identity. View.Edges hands an
			// operator the props of every edge on a declared relation, so they
			// are part of what it reads — and this function's whole rule is
			// that what is read decides re-dispatch and the cache key. Hashing
			// the id alone meant an edge whose props changed left the digest
			// identical: the pair stayed marked seen, the operator was never
			// dispatched again, and a cached result computed from the old props
			// was served in its place. An id is a content address of the
			// endpoints and relation, not of the props hung off it, so nothing
			// else was covering this.
			if e.Rel == nil || e.Rel.sch == nil {
				continue
			}
			for _, f := range e.Rel.sch.fields {
				writeField(h, f.name)
				v, set := e.Props.Get(f)
				if !set {
					writeField(h, "\x00unset")
					continue
				}
				writeValue(h, v)
			}
		}
	}
	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}

func writeValue(h interface{ Write([]byte) (int, error) }, v Value) {
	var b [9]byte
	b[0] = byte(v.kind)
	switch v.kind {
	case KindString:
		_, _ = h.Write(b[:1])
		writeField(h, v.str)
	case KindBytes:
		_, _ = h.Write(b[:1])
		writeField(h, string(v.raw))
	case KindFloat:
		binary.BigEndian.PutUint64(b[1:], uint64(int64(v.real*1e9)))
		_, _ = h.Write(b[:])
	default:
		binary.BigEndian.PutUint64(b[1:], uint64(v.num))
		_, _ = h.Write(b[:])
	}
}

// outEdges returns this node's outgoing edges of a relation, in EdgeID order.
func (g *Graph) outEdges(id NodeID, rel string) []*Edge {
	var out []*Edge
	for _, eid := range g.eord {
		e := g.edges[eid]
		if e.From == id && e.Rel.name == rel {
			out = append(out, e)
		}
	}
	return out
}
