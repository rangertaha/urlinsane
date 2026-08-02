// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"fmt"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/rangertaha/urlinsane/internal/graph"
)

// Side holds everything that is deliberately not part of a node's or an edge's
// content address: provenance, per-pair status, scheduler state, findings and
// the truncation ledger (docs/DESIGN.md §1.2, §8). It is one block, linked from
// the scan root, so a cross-run diff can compare structure without loading it.
//
// Rows reference nodes and edges by their 32-byte identity rather than by index
// into the root's link lists. Indices would be smaller, but they would also
// make every row's meaning depend on a list stored in a different block —
// a class of corruption that no checksum here would catch.
type Side struct {
	Version int

	// NodeProps is every retained assertion about a node's fields, winning or
	// not: §1.4's "disagreement between two sources is signal, not noise".
	NodeProps []PropRow
	// EdgeProps is the winning assertion per set edge field, and only that.
	// graph.Graph keeps a full assertion list for edges internally but exposes
	// no accessor for it, so the losing claims about an edge prop cannot be
	// persisted (see the package's known gaps).
	EdgeProps []EdgePropRow

	// Status is the terminal outcome of each (node, operator) pair.
	Status []StatusRow
	// Sched is depth and seed-closure membership, which §8 keeps out of the
	// addressed form precisely so they cannot perturb a CID.
	Sched []SchedRow

	// Scores are plugin-model judgements about entities (§10.6). They are side
	// table state and not props on purpose: as props they would make every
	// node's CID depend on a model version, so retraining anything would
	// invalidate every stored address and break the diff this store exists for.
	Scores []ScoreRow

	// Ledger is the truncation ledger: candidates the engine declined. It
	// doubles as a denylist, so it must be restored before any node is
	// re-admitted or "pruning is irreversible" stops being true across a
	// resume.
	Ledger []graph.LedgerRow
	// Truncations are run-level limits that bound the expansion.
	Truncations []graph.RunTruncation

	// Findings are the analyzers' output (§9).
	Findings []graph.Finding
}

// PropRow is one assertion, retained with the provenance that produced it.
// Subject is a NodeID for NodeProps and an EdgeID for EdgeProps.
type PropRow struct {
	Subject  [32]byte
	Field    string
	Kind     graph.Kind
	Value    graph.Value
	Operator string
	Round    int
	Won      bool
}

// EdgePropRow is one assertion about an edge's props. The edge is named by
// (from, rel, to) rather than by EdgeID because EdgeID is a hash the graph
// package computes internally and never exposes — and that hash is taken over
// exactly these three fields, so naming them directly loses nothing and keeps
// the row readable without a graph in hand.
type EdgePropRow struct {
	From     graph.NodeID
	Rel      string
	To       graph.NodeID
	Field    string
	Kind     graph.Kind
	Value    graph.Value
	Operator string
	Round    int
	Won      bool
}

// StatusRow is the terminal outcome of one (node, operator) pair. How a lookup
// failed is itself the finding, so the status is persisted rather than
// collapsed into a boolean.
type StatusRow struct {
	Node     graph.NodeID
	Operator string
	Status   graph.Status
}

// SchedRow is the scheduler state for one node: its shortest observation
// distance from the seed and whether it is in the seed closure. Both are
// derived from structure, and both are needed to resume a partial scan.
type SchedRow struct {
	Node      graph.NodeID
	Depth     int
	InClosure bool
}

// ScoreRow is one plugin model's judgement about one node.
type ScoreRow struct {
	Node  graph.NodeID
	Key   string
	Score float64
}

// encodeSide writes the side block:
//
//	[version, nodeProps, edgeProps, status, sched, scores, ledger,
//	 truncations, findings]
func encodeSide(s *Side) ([]byte, cid.Cid, error) {
	return encodeList(9, func(e *enc) {
		e.i64(int64(s.Version))
		e.sub(len(s.NodeProps), func(e *enc) {
			for _, r := range s.NodeProps {
				encodePropRow(e, r)
			}
		})
		e.sub(len(s.EdgeProps), func(e *enc) {
			for _, r := range s.EdgeProps {
				e.sub(9, func(e *enc) {
					e.raw(r.From[:])
					e.str(r.Rel)
					e.raw(r.To[:])
					e.str(r.Field)
					e.i64(int64(r.Kind))
					encodeValue(e, r.Kind, r.Value)
					e.str(r.Operator)
					e.i64(int64(r.Round))
					e.flag(r.Won)
				})
			}
		})
		e.sub(len(s.Status), func(e *enc) {
			for _, r := range s.Status {
				e.sub(3, func(e *enc) {
					e.raw(r.Node[:])
					e.str(r.Operator)
					e.i64(int64(r.Status))
				})
			}
		})
		e.sub(len(s.Sched), func(e *enc) {
			for _, r := range s.Sched {
				e.sub(3, func(e *enc) {
					e.raw(r.Node[:])
					e.i64(int64(r.Depth))
					e.flag(r.InClosure)
				})
			}
		})
		e.sub(len(s.Scores), func(e *enc) {
			for _, r := range s.Scores {
				e.sub(3, func(e *enc) {
					e.raw(r.Node[:])
					e.str(r.Key)
					e.f64(r.Score)
				})
			}
		})
		e.sub(len(s.Ledger), func(e *enc) {
			for _, r := range s.Ledger {
				e.sub(7, func(e *enc) {
					e.str(r.Type)
					e.str(r.Key)
					e.i64(int64(r.Depth))
					e.f64(r.Belief)
					e.i64(int64(r.Reason))
					e.str(r.By.Operator)
					e.i64(int64(r.By.Round))
				})
			}
		})
		e.sub(len(s.Truncations), func(e *enc) {
			for _, r := range s.Truncations {
				e.sub(3, func(e *enc) {
					e.i64(int64(r.Reason))
					e.i64(int64(r.Round))
					e.str(r.Detail)
				})
			}
		})
		e.sub(len(s.Findings), func(e *enc) {
			for _, f := range s.Findings {
				encodeFinding(e, f)
			}
		})
	})
}

func encodePropRow(e *enc, r PropRow) {
	e.sub(7, func(e *enc) {
		e.raw(r.Subject[:])
		e.str(r.Field)
		e.i64(int64(r.Kind))
		encodeValue(e, r.Kind, r.Value)
		e.str(r.Operator)
		e.i64(int64(r.Round))
		e.flag(r.Won)
	})
}

// encodeValue writes a graph.Value under its declared kind. The kind is stored
// alongside so a decoder can detect a type whose field list was reordered or
// truncated — the one schema-evolution mistake §1.3 says corrupts addresses
// instead of failing loudly.
func encodeValue(e *enc, k graph.Kind, v graph.Value) {
	switch k {
	case graph.KindString:
		e.str(v.Str())
	case graph.KindInt, graph.KindTime:
		e.i64(v.Num())
	case graph.KindFloat:
		e.f64(v.Real())
	case graph.KindBool:
		e.flag(v.Flag())
	case graph.KindBytes:
		e.raw(v.Raw())
	default:
		e.fail(fmt.Errorf("store: cannot encode a value of kind %s", k))
	}
}

func decodeTypedValue(d *dec, k graph.Kind) graph.Value {
	switch k {
	case graph.KindString:
		return graph.String(d.str())
	case graph.KindInt:
		return graph.Int(d.i64())
	case graph.KindTime:
		return graph.Time(time.Unix(0, d.i64()).UTC())
	case graph.KindFloat:
		return graph.Float(d.f64())
	case graph.KindBool:
		return graph.Bool(d.flag())
	case graph.KindBytes:
		return graph.Bytes(d.raw())
	}
	d.fail("store: cannot decode a value of kind %d", k)
	return graph.Value{}
}

func encodeFinding(e *enc, f graph.Finding) {
	e.sub(6, func(e *enc) {
		e.str(f.Kind)
		e.i64(int64(f.Severity))
		e.sub(len(f.Nodes), func(e *enc) {
			for _, id := range f.Nodes {
				e.raw(id[:])
			}
		})
		e.sub(len(f.Declined), func(e *enc) {
			for _, r := range f.Declined {
				e.sub(2, func(e *enc) {
					e.str(r.Type)
					e.str(r.Key)
				})
			}
		})
		e.str(f.Summary)
		e.sub(len(f.Evidence), func(e *enc) {
			for _, p := range f.Evidence {
				e.sub(2, func(e *enc) {
					e.str(p.Operator)
					e.i64(int64(p.Round))
				})
			}
		})
	})
}

// decodeSide parses the side block.
func decodeSide(block []byte) (*Side, error) {
	n, err := decodeBlock(block)
	if err != nil {
		return nil, err
	}
	d := newDec(n).expect("side block", 9)
	s := &Side{Version: int(d.at(0).i64())}
	if v := s.Version; v != FormatVersion && d.err() == nil {
		return nil, fmt.Errorf("store: side block format version %d, this build writes %d", v, FormatVersion)
	}

	d.at(1).each(func(r *dec) { s.NodeProps = append(s.NodeProps, decodePropRow(r)) })
	d.at(2).each(func(r *dec) { s.EdgeProps = append(s.EdgeProps, decodeEdgePropRow(r)) })

	d.at(3).each(func(r *dec) {
		r.expect("status row", 3)
		s.Status = append(s.Status, StatusRow{
			Node:     graph.NodeID(r.at(0).id32()),
			Operator: r.at(1).str(),
			Status:   graph.Status(r.at(2).i64()),
		})
	})
	d.at(4).each(func(r *dec) {
		r.expect("sched row", 3)
		s.Sched = append(s.Sched, SchedRow{
			Node:      graph.NodeID(r.at(0).id32()),
			Depth:     int(r.at(1).i64()),
			InClosure: r.at(2).flag(),
		})
	})
	d.at(5).each(func(r *dec) {
		r.expect("score row", 3)
		s.Scores = append(s.Scores, ScoreRow{
			Node:  graph.NodeID(r.at(0).id32()),
			Key:   r.at(1).str(),
			Score: r.at(2).f64(),
		})
	})
	d.at(6).each(func(r *dec) {
		r.expect("ledger row", 7)
		s.Ledger = append(s.Ledger, graph.LedgerRow{
			Type:   r.at(0).str(),
			Key:    r.at(1).str(),
			Depth:  int(r.at(2).i64()),
			Belief: r.at(3).f64(),
			Reason: graph.Reason(r.at(4).i64()),
			By:     graph.Provenance{Operator: r.at(5).str(), Round: int(r.at(6).i64())},
		})
	})
	d.at(7).each(func(r *dec) {
		r.expect("truncation row", 3)
		s.Truncations = append(s.Truncations, graph.RunTruncation{
			Reason: graph.Reason(r.at(0).i64()),
			Round:  int(r.at(1).i64()),
			Detail: r.at(2).str(),
		})
	})
	d.at(8).each(func(r *dec) { s.Findings = append(s.Findings, decodeFinding(r)) })

	if err := d.err(); err != nil {
		return nil, err
	}
	return s, nil
}

func decodePropRow(d *dec) PropRow {
	d.expect("prop row", 7)
	r := PropRow{
		Subject: d.at(0).id32(),
		Field:   d.at(1).str(),
		Kind:    graph.Kind(d.at(2).i64()),
	}
	r.Value = decodeTypedValue(d.at(3), r.Kind)
	r.Operator = d.at(4).str()
	r.Round = int(d.at(5).i64())
	r.Won = d.at(6).flag()
	return r
}

func decodeEdgePropRow(d *dec) EdgePropRow {
	d.expect("edge prop row", 9)
	r := EdgePropRow{
		From:  graph.NodeID(d.at(0).id32()),
		Rel:   d.at(1).str(),
		To:    graph.NodeID(d.at(2).id32()),
		Field: d.at(3).str(),
		Kind:  graph.Kind(d.at(4).i64()),
	}
	r.Value = decodeTypedValue(d.at(5), r.Kind)
	r.Operator = d.at(6).str()
	r.Round = int(d.at(7).i64())
	r.Won = d.at(8).flag()
	return r
}

func decodeFinding(d *dec) graph.Finding {
	d.expect("finding row", 6)
	f := graph.Finding{Kind: d.at(0).str(), Severity: graph.Severity(d.at(1).i64())}
	d.at(2).each(func(v *dec) { f.Nodes = append(f.Nodes, graph.NodeID(v.id32())) })
	d.at(3).each(func(v *dec) {
		v.expect("ledger ref", 2)
		f.Declined = append(f.Declined, graph.LedgerRef{Type: v.at(0).str(), Key: v.at(1).str()})
	})
	f.Summary = d.at(4).str()
	d.at(5).each(func(v *dec) {
		v.expect("provenance", 2)
		f.Evidence = append(f.Evidence, graph.Provenance{Operator: v.at(0).str(), Round: int(v.at(1).i64())})
	})
	return f
}
