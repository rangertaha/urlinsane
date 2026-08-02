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

	"github.com/ipfs/go-cid"
	"github.com/rangertaha/urlinsane/internal/graph"
)

// Store persists graphs as content-addressed blocks.
type Store struct {
	bs Blockstore
}

// Open returns a Store backed by a filesystem blockstore under dir.
func Open(dir string) (*Store, error) {
	bs, err := NewFSBlockstore(dir)
	if err != nil {
		return nil, err
	}
	return &Store{bs: bs}, nil
}

// New returns a Store over any blockstore.
func New(bs Blockstore) *Store { return &Store{bs: bs} }

// SaveOptions supplies the one thing a graph holds but does not publish.
//
// The seed has no accessor and cannot be inferred: structural edges cost no
// depth, so a composite seed puts several nodes at depth 0 inside the closure
// and only the caller that called Seed knows which one it started from.
type SaveOptions struct {
	Seed graph.NodeID
}

// Save writes every node, edge and side table and returns the scan root's CID.
//
// The root is a pure function of the graph's content: two identical scans
// produce the same CID, byte for byte. Nothing here reads a clock or a random
// source, and every collection comes from an accessor with a defined order —
// no Go map is ranged over anywhere on the path to a block.
func (s *Store) Save(g *graph.Graph, opts SaveOptions) (cid.Cid, error) {
	if opts.Seed.IsZero() {
		return cid.Undef, fmt.Errorf("graphstore: SaveOptions.Seed is required")
	}
	seed, ok := g.Node(opts.Seed)
	if !ok {
		return cid.Undef, fmt.Errorf("graphstore: seed %s is not in the graph", opts.Seed)
	}

	root := &Root{Version: FormatVersion, SeedType: seed.Type.Name(), SeedKey: seed.Key}

	nodes := g.Nodes()
	for _, n := range nodes {
		block, c, err := EncodeNode(n)
		if err != nil {
			return cid.Undef, fmt.Errorf("graphstore: encoding node %s/%s: %w", n.Type.Name(), n.Key, err)
		}
		if err := s.bs.Put(c, block); err != nil {
			return cid.Undef, err
		}
		root.Nodes = append(root.Nodes, c)
	}

	edges := g.Edges()
	for _, e := range edges {
		block, c, err := EncodeEdge(e)
		if err != nil {
			return cid.Undef, fmt.Errorf("graphstore: encoding edge %s: %w", e.Rel.Name(), err)
		}
		if err := s.bs.Put(c, block); err != nil {
			return cid.Undef, err
		}
		root.Edges = append(root.Edges, c)
	}

	side := collectSide(g, nodes, edges)
	sideBlock, sideCID, err := encodeSide(side)
	if err != nil {
		return cid.Undef, fmt.Errorf("graphstore: encoding side tables: %w", err)
	}
	if err := s.bs.Put(sideCID, sideBlock); err != nil {
		return cid.Undef, err
	}
	root.Side = sideCID

	rootBlock, rootCID, err := encodeRoot(root)
	if err != nil {
		return cid.Undef, fmt.Errorf("graphstore: encoding scan root: %w", err)
	}
	if err := s.bs.Put(rootCID, rootBlock); err != nil {
		return cid.Undef, err
	}
	return rootCID, nil
}

// collectSide gathers everything kept out of the addressed form. Every loop
// here walks a slice in a defined order; introducing a map range in this
// function would silently change the scan root's CID between identical runs.
func collectSide(g *graph.Graph, nodes []*graph.Node, edges []*graph.Edge) *Side {
	side := &Side{Version: FormatVersion}

	// Status and scores live in graph-internal maps. They are read through the
	// analysis surface, which sorts them by operator and by key; ranging the
	// maps would put Go's iteration order into a content address.
	view := g.Analyze()

	for _, n := range nodes {
		for _, as := range g.Assertions(n.ID) {
			side.NodeProps = append(side.NodeProps, PropRow{
				Subject:  n.ID,
				Field:    as.Field,
				Kind:     as.Value.Kind(),
				Value:    as.Value,
				Operator: as.By.Operator,
				Round:    as.By.Round,
				Won:      as.Won,
			})
		}
	}

	// Edge props are read back from the materialized values because the edge
	// assertion table has no accessor. The winning operator is recoverable via
	// Props.Setter, the round is not, so it is recorded as 0.
	for _, e := range edges {
		e.Props.Each(func(f graph.Field, v graph.Value) {
			op, _ := e.Props.Setter(f)
			side.EdgeProps = append(side.EdgeProps, EdgePropRow{
				From:     e.From,
				Rel:      e.Rel.Name(),
				To:       e.To,
				Field:    f.Name(),
				Kind:     v.Kind(),
				Value:    v,
				Operator: op,
				Won:      true,
			})
		})
	}

	for _, n := range nodes {
		for _, r := range view.Statuses(n.ID) {
			side.Status = append(side.Status, StatusRow{Node: n.ID, Operator: r.Operator, Status: r.Status})
		}
	}

	for _, n := range nodes {
		side.Sched = append(side.Sched, SchedRow{
			Node:      n.ID,
			Depth:     g.Depth(n.ID),
			InClosure: g.InClosure(n.ID),
		})
	}

	for _, n := range nodes {
		for _, r := range view.Scores(n.ID) {
			side.Scores = append(side.Scores, ScoreRow{Node: n.ID, Key: r.Key, Score: r.Value})
		}
	}

	side.Ledger = g.Ledger()           // already in (type, key, reason) order
	side.Truncations = g.Truncations() // recorded in round order
	side.Findings = g.Findings()       // sorted by severity, kind, summary
	return side
}

// Load reads a scan root and every block it links.
func (s *Store) Load(root cid.Cid) (*Scan, error) {
	block, err := s.bs.Get(root)
	if err != nil {
		return nil, err
	}
	r, err := decodeRoot(block)
	if err != nil {
		return nil, err
	}
	scan := &Scan{CID: root, Root: r}

	for _, c := range r.Nodes {
		b, err := s.bs.Get(c)
		if err != nil {
			return nil, err
		}
		nb, err := DecodeNode(b)
		if err != nil {
			return nil, fmt.Errorf("graphstore: node %s: %w", c, err)
		}
		nb.CID = c
		scan.Nodes = append(scan.Nodes, nb)
	}

	for _, c := range r.Edges {
		b, err := s.bs.Get(c)
		if err != nil {
			return nil, err
		}
		eb, err := DecodeEdge(b)
		if err != nil {
			return nil, fmt.Errorf("graphstore: edge %s: %w", c, err)
		}
		eb.CID = c
		scan.Edges = append(scan.Edges, eb)
	}

	sb, err := s.bs.Get(r.Side)
	if err != nil {
		return nil, err
	}
	scan.Side, err = decodeSide(sb)
	if err != nil {
		return nil, fmt.Errorf("graphstore: side tables %s: %w", r.Side, err)
	}
	return scan, nil
}
