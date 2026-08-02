// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"sort"

	"github.com/ipfs/go-cid"
	"github.com/rangertaha/urlinsane/internal/graph"
)

// NodeChange is one node's fate between two scans.
type NodeChange struct {
	Type string
	Key  string // canonical
	Old  cid.Cid
	New  cid.Cid
	// Slots are the positional prop indices whose values differ, valid only
	// for a changed node. Positions rather than field names because the
	// addressed form is positional and the registry publishes no way to map a
	// position back to a name; a caller holding the type's declared field list
	// can resolve them.
	Slots []int
}

// EdgeChange is one edge's fate between two scans. Edges are paired on
// (from, rel, to) — their identity — so a change means the edge's own props
// moved, not that it points somewhere else.
type EdgeChange struct {
	From  graph.NodeID
	Rel   string
	To    graph.NodeID
	Old   cid.Cid
	New   cid.Cid
	Slots []int
}

// Diff is what changed between two scans. This is the use case the whole
// content-addressing scheme exists for: pairing is by identity — (type,
// canonical key) for nodes — and the verdict is a CID comparison, so "changed"
// means the addressed content genuinely differs rather than that some
// timestamp moved.
//
// Only nodes and edges are compared. Side tables are excluded on purpose: a
// different round number or a re-run analyzer must not read as a changed graph.
type Diff struct {
	Added   []NodeChange
	Removed []NodeChange
	Changed []NodeChange
	Same    []NodeChange

	EdgesAdded   []EdgeChange
	EdgesRemoved []EdgeChange
	EdgesChanged []EdgeChange
}

// Empty reports whether the two scans have identical nodes and edges.
func (d *Diff) Empty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 &&
		len(d.EdgesAdded) == 0 && len(d.EdgesRemoved) == 0 && len(d.EdgesChanged) == 0
}

// Diff loads two scan roots and reports how the newer differs from the older.
func (s *Store) Diff(oldRoot, newRoot cid.Cid) (*Diff, error) {
	older, err := s.Load(oldRoot)
	if err != nil {
		return nil, err
	}
	newer, err := s.Load(newRoot)
	if err != nil {
		return nil, err
	}
	return DiffScans(older, newer), nil
}

type nodeKey struct{ typ, key string }

// DiffScans compares two loaded scans.
func DiffScans(older, newer *Scan) *Diff {
	d := &Diff{}
	diffNodes(d, older, newer)
	diffEdges(d, older, newer)
	return d
}

func diffNodes(d *Diff, older, newer *Scan) {
	old := make(map[nodeKey]*NodeBlock, len(older.Nodes))
	for _, n := range older.Nodes {
		old[nodeKey{n.Type, n.Key}] = n
	}

	for _, n := range newer.Nodes {
		k := nodeKey{n.Type, n.Key}
		o, ok := old[k]
		switch {
		case !ok:
			d.Added = append(d.Added, NodeChange{Type: n.Type, Key: n.Key, New: n.CID})
		case o.CID == n.CID:
			d.Same = append(d.Same, NodeChange{Type: n.Type, Key: n.Key, Old: o.CID, New: n.CID})
		default:
			d.Changed = append(d.Changed, NodeChange{
				Type: n.Type, Key: n.Key, Old: o.CID, New: n.CID,
				Slots: changedSlots(o.Values, n.Values),
			})
		}
		delete(old, k)
	}

	// Whatever is left was in the older scan only. Iterating the leftover map
	// would be nondeterministic, so the result is sorted below.
	for k, o := range old {
		d.Removed = append(d.Removed, NodeChange{Type: k.typ, Key: k.key, Old: o.CID})
	}

	for _, s := range [][]NodeChange{d.Added, d.Removed, d.Changed, d.Same} {
		sortNodeChanges(s)
	}
}

func diffEdges(d *Diff, older, newer *Scan) {
	old := make(map[edgeKey]*EdgeBlock, len(older.Edges))
	for _, e := range older.Edges {
		old[e.key()] = e
	}

	for _, e := range newer.Edges {
		k := e.key()
		o, ok := old[k]
		switch {
		case !ok:
			d.EdgesAdded = append(d.EdgesAdded, edgeChange(e, cid.Undef, e.CID, nil))
		case o.CID != e.CID:
			d.EdgesChanged = append(d.EdgesChanged,
				edgeChange(e, o.CID, e.CID, changedSlots(o.Values, e.Values)))
		}
		delete(old, k)
	}
	for _, o := range old {
		d.EdgesRemoved = append(d.EdgesRemoved, edgeChange(o, o.CID, cid.Undef, nil))
	}

	for _, s := range [][]EdgeChange{d.EdgesAdded, d.EdgesRemoved, d.EdgesChanged} {
		sortEdgeChanges(s)
	}
}

func edgeChange(e *EdgeBlock, old, new cid.Cid, slots []int) EdgeChange {
	return EdgeChange{From: e.From, Rel: e.Rel, To: e.To, Old: old, New: new, Slots: slots}
}

// changedSlots names the positions whose values differ. A slot present in one
// block and absent from the other counts as changed: that is what a field
// appended to the type looks like from here.
func changedSlots(old, new []RawValue) []int {
	var out []int
	n := len(old)
	if len(new) > n {
		n = len(new)
	}
	for i := 0; i < n; i++ {
		var a, b RawValue
		if i < len(old) {
			a = old[i]
		}
		if i < len(new) {
			b = new[i]
		}
		if !a.equal(b) {
			out = append(out, i)
		}
	}
	return out
}

// sortNodeChanges puts rows in the report's (type, key) order, so a diff of two
// diffs is itself meaningful.
func sortNodeChanges(s []NodeChange) {
	sort.Slice(s, func(i, j int) bool {
		if s[i].Type != s[j].Type {
			return s[i].Type < s[j].Type
		}
		return s[i].Key < s[j].Key
	})
}

func sortEdgeChanges(s []EdgeChange) {
	sort.Slice(s, func(i, j int) bool {
		a, b := s[i], s[j]
		if c := compareBytes(a.From[:], b.From[:]); c != 0 {
			return c < 0
		}
		if a.Rel != b.Rel {
			return a.Rel < b.Rel
		}
		return compareBytes(a.To[:], b.To[:]) < 0
	})
}

func compareBytes(a, b []byte) int {
	for i := range a {
		if i >= len(b) {
			return 1
		}
		switch {
		case a[i] < b[i]:
			return -1
		case a[i] > b[i]:
			return 1
		}
	}
	if len(a) < len(b) {
		return -1
	}
	return 0
}
