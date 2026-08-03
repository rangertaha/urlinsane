// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"fmt"
	"sort"

	"github.com/ipfs/go-cid"
	"github.com/rangertaha/urlinsane/internal/graph"
)

// replayOp is the provenance used for admissions that carry none of their own.
// Admitting a node or an edge records no provenance in the graph, so this id
// only ever surfaces in a rejection message — which is exactly where a caller
// wants to see that the replay, not an operator, was refused.
const replayOp = "store.replay"

// Rehydrated is a graph rebuilt from a stored scan, with the seed the caller
// needs to keep scanning.
type Rehydrated struct {
	Graph *graph.Graph
	Seed  graph.NodeID
	Scan  *Scan
}

// Rehydrate loads a scan root and rebuilds the graph it describes.
func (s *Store) Rehydrate(root cid.Cid, reg *graph.Registry) (*Rehydrated, error) {
	scan, err := s.Load(root)
	if err != nil {
		return nil, err
	}
	return RehydrateScan(scan, reg)
}

// RehydrateScan rebuilds a graph from a loaded scan, so a partial expansion can
// be resumed.
//
// Everything is replayed through the applier rather than written into the
// graph's internals, because the applier is the single writer and identity,
// canonicalization, merge policy and the seed-closure invariant all live there.
// A rebuild that bypassed it would be a second, divergent definition of what
// the graph means.
//
// The replay finishes by re-encoding every node and edge and checking the CIDs
// against the ones stored. That check is the contract: if it passes, the
// rebuilt graph is byte-identical to the one that was saved, and re-saving it
// produces the same scan root.
func RehydrateScan(scan *Scan, reg *graph.Registry) (*Rehydrated, error) {
	if scan.Side == nil {
		return nil, fmt.Errorf("store: scan %s has no side tables", scan.CID)
	}
	g := graph.New(reg)

	seed, err := g.Seed(scan.Root.SeedType, scan.Root.SeedKey)
	if err != nil {
		return nil, fmt.Errorf("store: reseeding: %w", err)
	}

	byID, err := replayNodes(g, scan, seed)
	if err != nil {
		return nil, err
	}
	if err := replayEdges(g, reg, scan, byID); err != nil {
		return nil, err
	}
	if err := replayProps(g, scan, byID); err != nil {
		return nil, err
	}
	if err := replaySideState(g, scan, byID); err != nil {
		return nil, err
	}
	if err := verify(g, scan); err != nil {
		return nil, err
	}
	return &Rehydrated{Graph: g, Seed: seed, Scan: scan}, nil
}

// replayNodes admits every stored node and maps identities back to refs.
//
// Every node is admitted against the seed, which fixes its provisional depth at
// 0 until an in-edge overrides it in the next phase. A node that never had an
// in-edge therefore comes back at depth 0 whatever depth it was recorded at:
// Apply takes a bare node's depth from the subject the operator ran on, and
// nothing in the addressed form or the side tables says which subject that was.
// The settled check below catches the discrepancy and fails the rehydration
// rather than resuming a scan with wrong depths.
//
// Nodes are admitted one at a time because Apply reports admitted ids in a
// slice that only lines up with the delta when nothing was refused; a
// one-element delta makes the mapping exact instead of positional-and-hopeful.
// The recomputed NodeID must equal the stored one — both are
// hash(type, canonical key) — which is what lets the stored edges, keyed on the
// old ids, resolve at all.
func replayNodes(g *graph.Graph, scan *Scan, seed graph.NodeID) (map[graph.NodeID]graph.NodeRef, error) {
	byID := make(map[graph.NodeID]graph.NodeRef, len(scan.Nodes))
	for _, nb := range scan.Nodes {
		ref := graph.NodeRef{Type: nb.Type, Key: nb.Key}
		res := g.Apply(graph.Provenance{Operator: replayOp}, seed, graph.Delta{Nodes: []graph.NodeRef{ref}})
		if len(res.Nodes) != 1 {
			return nil, fmt.Errorf("store: replaying node %s/%s: %v", nb.Type, nb.Key, res.Rejected)
		}
		byID[res.Nodes[0]] = ref
	}
	return byID, nil
}

// edgePlan is one stored edge with the ordering information the replay needs.
type edgePlan struct {
	ref   graph.EdgeRef
	depth int // the stored depth of From
	class graph.Class
}

// replayEdges re-admits every edge, then iterates until depth and closure
// settle.
//
// Order matters and a single pass is not enough in general. Depth is derived
// from the source's depth, and closure membership propagates along structural
// edges, so an edge applied before its source is settled produces the wrong
// answer — and a variant edge applied before its source has entered the closure
// is refused outright by the applier's invariant. Sorting by the stored depth
// of the source, structural edges first, gets a single pass right in almost
// every graph; the loop is what makes it correct in the rest, including cyclic
// ones, since depth only ever decreases and closure only ever grows.
func replayEdges(g *graph.Graph, reg *graph.Registry, scan *Scan, byID map[graph.NodeID]graph.NodeRef) error {
	depths := make(map[graph.NodeID]int, len(scan.Side.Sched))
	for _, r := range scan.Side.Sched {
		depths[r.Node] = r.Depth
	}

	plans := make([]edgePlan, 0, len(scan.Edges))
	for _, eb := range scan.Edges {
		from, ok := byID[eb.From]
		if !ok {
			return fmt.Errorf("store: edge %s references unknown source %s", eb.Rel, eb.From)
		}
		to, ok := byID[eb.To]
		if !ok {
			return fmt.Errorf("store: edge %s references unknown target %s", eb.Rel, eb.To)
		}
		rel, ok := reg.Rel(eb.Rel)
		if !ok {
			return fmt.Errorf("store: relation %q is not registered", eb.Rel)
		}
		plans = append(plans, edgePlan{
			ref:   graph.EdgeRef{From: from, Rel: eb.Rel, To: to},
			depth: depths[eb.From],
			class: rel.Class(),
		})
	}

	sort.SliceStable(plans, func(i, j int) bool {
		if plans[i].depth != plans[j].depth {
			return plans[i].depth < plans[j].depth
		}
		return plans[i].class < plans[j].class // Structural < Variant < Observation
	})

	by := graph.Provenance{Operator: replayOp}
	// One pass per node is the worst case: each pass settles at least one more
	// node's depth or closure membership, or the state has stopped moving.
	maxPasses := len(scan.Nodes) + 2
	var last error
	for pass := 0; pass < maxPasses; pass++ {
		for _, p := range plans {
			g.Apply(by, graph.NodeID{}, graph.Delta{Edges: []graph.EdgeRef{p.ref}})
		}
		if last = settled(g, scan); last == nil {
			return nil
		}
	}
	return fmt.Errorf("store: edge replay did not settle after %d passes: %w", maxPasses, last)
}

// settled reports whether the replayed graph matches the stored edge set and
// scheduler state. It returns the first discrepancy as an error so a replay
// that never converges says which node it disagreed about.
func settled(g *graph.Graph, scan *Scan) error {
	have := make(map[edgeKey]bool, len(scan.Edges))
	for _, e := range g.Edges() {
		have[edgeKey{from: e.From, rel: e.Rel.Name(), to: e.To}] = true
	}
	for _, eb := range scan.Edges {
		if !have[eb.key()] {
			return fmt.Errorf("edge %s -> %s -> %s was not admitted", eb.From, eb.Rel, eb.To)
		}
	}
	if len(have) != len(scan.Edges) {
		return fmt.Errorf("replay produced %d edges, stored %d", len(have), len(scan.Edges))
	}
	for _, r := range scan.Side.Sched {
		if d := g.Depth(r.Node); d != r.Depth {
			return fmt.Errorf("node %s has depth %d, stored %d", r.Node, d, r.Depth)
		}
		if c := g.InClosure(r.Node); c != r.InClosure {
			return fmt.Errorf("node %s closure membership is %v, stored %v", r.Node, c, r.InClosure)
		}
	}
	return nil
}

// replayProps re-asserts every stored prop with its original provenance.
//
// Node assertions are replayed one at a time and in stored order, losers
// included. That restores the merge outcome — the policy is order-independent
// by construction — and the provenance table with it, so §1.4's "disagreement
// between two sources is preserved as signal" survives a resume.
func replayProps(g *graph.Graph, scan *Scan, byID map[graph.NodeID]graph.NodeRef) error {
	for _, r := range scan.Side.NodeProps {
		id := graph.NodeID(r.Subject)
		ref, ok := byID[id]
		if !ok {
			return fmt.Errorf("store: prop %q references unknown node %s", r.Field, id)
		}
		res := g.Apply(
			graph.Provenance{Operator: r.Operator, Round: r.Round},
			id,
			graph.Delta{Props: []graph.PropSet{{Node: &ref, Field: r.Field, Value: r.Value}}},
		)
		if len(res.Rejected) > 0 {
			return fmt.Errorf("store: replaying %s.%s: %v", ref.Key, r.Field, res.Rejected)
		}
	}

	for _, r := range scan.Side.EdgeProps {
		from, ok := byID[r.From]
		if !ok {
			return fmt.Errorf("store: edge prop %q references unknown source %s", r.Field, r.From)
		}
		to, ok := byID[r.To]
		if !ok {
			return fmt.Errorf("store: edge prop %q references unknown target %s", r.Field, r.To)
		}
		eref := graph.EdgeRef{From: from, Rel: r.Rel, To: to}
		res := g.Apply(
			graph.Provenance{Operator: r.Operator, Round: r.Round},
			graph.NodeID{},
			graph.Delta{Props: []graph.PropSet{{Edge: &eref, Field: r.Field, Value: r.Value}}},
		)
		if len(res.Rejected) > 0 {
			return fmt.Errorf("store: replaying %s.%s: %v", r.Rel, r.Field, res.Rejected)
		}
	}
	return nil
}

// replaySideState restores the tables that are pure records: status, scores,
// findings, run truncations and the truncation ledger.
//
// The ledger goes last on purpose. It is a denylist as well as a record, and
// restoring it before the nodes would have the applier refuse any node that had
// also, at some point, been declined. Restoring it before the scan resumes is
// what keeps "pruning is irreversible" true across a resume.
func replaySideState(g *graph.Graph, scan *Scan, byID map[graph.NodeID]graph.NodeRef) error {
	// Before the statuses, because Existence is computed from both and reading
	// a status without knowing who may attest to existence is what made `typo`
	// and `report` disagree about the same bytes: with no observer set the
	// graph falls back to "everything observes", so a decomposer's successful
	// parse counted as proof the name exists and every syntactically valid
	// variant re-rendered as live.
	g.SetObservers(scan.Side.Observers)

	for _, r := range scan.Side.Status {
		if _, ok := byID[r.Node]; !ok {
			return fmt.Errorf("store: status references unknown node %s", r.Node)
		}
		g.SetStatus(r.Node, r.Operator, r.Status)
	}
	for _, r := range scan.Side.Scores {
		if _, ok := byID[r.Node]; !ok {
			return fmt.Errorf("store: score references unknown node %s", r.Node)
		}
		g.SetScore(r.Node, r.Key, r.Score)
	}
	for _, t := range scan.Side.Truncations {
		g.NoteTruncation(t.Reason, t.Round, t.Detail)
	}
	if len(scan.Side.Findings) > 0 {
		g.AddFindings(scan.Side.Findings...)
	}
	for _, r := range scan.Side.Ledger {
		if err := g.Decline(r.Type, r.Key, r.Depth, r.Belief, r.Reason, r.By); err != nil {
			return fmt.Errorf("store: restoring ledger row %s/%s: %w", r.Type, r.Key, err)
		}
	}
	return nil
}

// verify re-encodes the rebuilt graph and compares every address with the one
// stored. Nodes and edges come back in the same canonical order they were
// written in, so the comparison is positional.
func verify(g *graph.Graph, scan *Scan) error {
	nodes := g.Nodes()
	if len(nodes) != len(scan.Nodes) {
		return fmt.Errorf("store: rebuilt %d nodes, stored %d", len(nodes), len(scan.Nodes))
	}
	for i, n := range nodes {
		_, c, err := EncodeNode(n)
		if err != nil {
			return err
		}
		if c != scan.Nodes[i].CID {
			return fmt.Errorf("store: node %s/%s re-encodes to %s, stored %s",
				n.Type.Name(), n.Key, c, scan.Nodes[i].CID)
		}
	}

	edges := g.Edges()
	if len(edges) != len(scan.Edges) {
		return fmt.Errorf("store: rebuilt %d edges, stored %d", len(edges), len(scan.Edges))
	}
	for i, e := range edges {
		_, c, err := EncodeEdge(e)
		if err != nil {
			return err
		}
		if c != scan.Edges[i].CID {
			return fmt.Errorf("store: edge %s re-encodes to %s, stored %s", e.Rel.Name(), c, scan.Edges[i].CID)
		}
	}
	return nil
}
