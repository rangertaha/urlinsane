// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"fmt"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/rangertaha/urlinsane/internal/graph"
)

// EncodeNode returns a node's addressed block and CID. It is a thin wrapper over
// graph.Node.Addressed: the addressed form belongs to the graph package, and a
// second encoder here would be a second definition of identity.
func EncodeNode(n *graph.Node) ([]byte, cid.Cid, error) { return n.Addressed() }

// EncodeEdge returns an edge's addressed block and CID.
func EncodeEdge(e *graph.Edge) ([]byte, cid.Cid, error) { return e.Addressed() }

// NodeBlock is a decoded node block: [type, key, [values...]].
//
// It deliberately stops short of a graph.Node. Rebuilding one needs a registry
// (for the type handle and the field kinds) and a Graph (which owns identity
// and admission), so decoding yields the block's literal content and
// rehydration does the rest.
type NodeBlock struct {
	CID    cid.Cid // set by the store on load; zero when decoded standalone
	Type   string
	Key    string // canonical
	Values []RawValue
}

// EdgeBlock is a decoded edge block: [from, relation, to, [values...]].
type EdgeBlock struct {
	CID    cid.Cid
	From   graph.NodeID
	To     graph.NodeID
	Rel    string
	Values []RawValue
}

// key is the edge's cross-run identity. EdgeID is hash(from, rel, to) and all
// three are in the block, so an edge can be paired across runs without the
// identity ever being stored.
func (e *EdgeBlock) key() edgeKey { return edgeKey{from: e.From, rel: e.Rel, to: e.To} }

type edgeKey struct {
	from graph.NodeID
	rel  string
	to   graph.NodeID
}

// RawValue is one decoded prop slot. The addressed form carries positions and
// values but no field kinds — those live in the registry — so a decoded slot
// keeps whatever the codec produced and is bound to a graph.Value only once a
// declared Kind says what it should be. A KindTime and a KindInt are both a
// CBOR integer on the wire; nothing in the block distinguishes them.
type RawValue struct {
	Set   bool           // false for a declared-but-unset slot (encoded as null)
	Kind  datamodel.Kind // the wire kind, valid when Set
	Str   string
	Int   int64
	Float float64
	Bool  bool
	Bytes []byte
}

// Bind converts a decoded slot to a graph.Value under a declared field kind,
// failing when the wire kind cannot supply it.
func (r RawValue) Bind(k graph.Kind) (graph.Value, error) {
	if !r.Set {
		return graph.Value{}, fmt.Errorf("store: slot is unset")
	}
	mismatch := func(want datamodel.Kind) error {
		return fmt.Errorf("store: field of kind %s wants %s on the wire, got %s", k, want, r.Kind)
	}
	switch k {
	case graph.KindString:
		if r.Kind != datamodel.Kind_String {
			return graph.Value{}, mismatch(datamodel.Kind_String)
		}
		return graph.String(r.Str), nil
	case graph.KindInt:
		if r.Kind != datamodel.Kind_Int {
			return graph.Value{}, mismatch(datamodel.Kind_Int)
		}
		return graph.Int(r.Int), nil
	case graph.KindTime:
		// Times are Unix nanoseconds on the wire: an integer with no location
		// and no formatting, so the same instant always encodes identically.
		if r.Kind != datamodel.Kind_Int {
			return graph.Value{}, mismatch(datamodel.Kind_Int)
		}
		return graph.Time(time.Unix(0, r.Int).UTC()), nil
	case graph.KindFloat:
		if r.Kind != datamodel.Kind_Float {
			return graph.Value{}, mismatch(datamodel.Kind_Float)
		}
		return graph.Float(r.Float), nil
	case graph.KindBool:
		if r.Kind != datamodel.Kind_Bool {
			return graph.Value{}, mismatch(datamodel.Kind_Bool)
		}
		return graph.Bool(r.Bool), nil
	case graph.KindBytes:
		if r.Kind != datamodel.Kind_Bytes {
			return graph.Value{}, mismatch(datamodel.Kind_Bytes)
		}
		return graph.Bytes(r.Bytes), nil
	}
	return graph.Value{}, fmt.Errorf("store: unknown field kind %d", k)
}

// equal reports slot equality, which is how the diff names the fields that
// changed without needing a registry.
func (r RawValue) equal(o RawValue) bool {
	if r.Set != o.Set {
		return false
	}
	if !r.Set {
		return true
	}
	if r.Kind != o.Kind {
		return false
	}
	switch r.Kind {
	case datamodel.Kind_String:
		return r.Str == o.Str
	case datamodel.Kind_Int:
		return r.Int == o.Int
	case datamodel.Kind_Float:
		return r.Float == o.Float
	case datamodel.Kind_Bool:
		return r.Bool == o.Bool
	case datamodel.Kind_Bytes:
		return string(r.Bytes) == string(o.Bytes)
	}
	return false
}

// Slot resolves one positional prop slot of a stored node under a declared
// kind, so report code can read a stored scan without rehydrating a whole
// graph. The caller supplies the index because the registry exposes fields by
// name only: neither graph.NodeType nor graph.Field publishes the declared
// field list or a field's position, and the position is exactly what the
// addressed form is keyed on.
func (n *NodeBlock) Slot(i int, k graph.Kind) (graph.Value, bool, error) {
	if i < 0 || i >= len(n.Values) || !n.Values[i].Set {
		return graph.Value{}, false, nil
	}
	v, err := n.Values[i].Bind(k)
	if err != nil {
		return graph.Value{}, false, err
	}
	return v, true, nil
}

// DecodeNode parses a node block. It is the inverse of graph.Node.Addressed.
func DecodeNode(block []byte) (*NodeBlock, error) {
	n, err := decodeBlock(block)
	if err != nil {
		return nil, err
	}
	d := newDec(n).expect("node block", 3)
	nb := &NodeBlock{Type: d.at(0).str(), Key: d.at(1).str()}
	nb.Values = decodeValues(d.at(2))
	if err := d.err(); err != nil {
		return nil, err
	}
	return nb, nil
}

// DecodeEdge parses an edge block. It is the inverse of graph.Edge.Addressed.
func DecodeEdge(block []byte) (*EdgeBlock, error) {
	n, err := decodeBlock(block)
	if err != nil {
		return nil, err
	}
	d := newDec(n).expect("edge block", 4)
	eb := &EdgeBlock{
		From: graph.NodeID(d.at(0).id32()),
		Rel:  d.at(1).str(),
		To:   graph.NodeID(d.at(2).id32()),
	}
	eb.Values = decodeValues(d.at(3))
	if err := d.err(); err != nil {
		return nil, err
	}
	return eb, nil
}

// decodeValues reads the positional prop list. Null is a declared-but-unset
// slot: the encoder writes every declared slot so a node encodes identically
// wherever it was built, which means the decoder sees holes rather than gaps.
func decodeValues(d *dec) []RawValue {
	out := make([]RawValue, 0, d.len())
	d.each(func(v *dec) {
		out = append(out, decodeValue(v))
	})
	return out
}

func decodeValue(d *dec) RawValue {
	if d.err() != nil {
		return RawValue{}
	}
	switch k := d.n.Kind(); k {
	case datamodel.Kind_Null:
		return RawValue{}
	case datamodel.Kind_String:
		return RawValue{Set: true, Kind: k, Str: d.str()}
	case datamodel.Kind_Int:
		return RawValue{Set: true, Kind: k, Int: d.i64()}
	case datamodel.Kind_Float:
		return RawValue{Set: true, Kind: k, Float: d.f64()}
	case datamodel.Kind_Bool:
		return RawValue{Set: true, Kind: k, Bool: d.flag()}
	case datamodel.Kind_Bytes:
		return RawValue{Set: true, Kind: k, Bytes: d.raw()}
	default:
		d.fail("store: prop slot has unencodable kind %s", k)
		return RawValue{}
	}
}
