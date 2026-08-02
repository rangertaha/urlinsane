// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/multiformats/go-multihash"
)

// NodeID is a node's stable identity: hash(type, canonical key). It is fixed
// for the node's whole life and is what the scheduler, seen-set, cache, edges
// and side tables key on. It is deliberately *not* the content address — props
// accumulate, so the CID changes while the identity does not.
type NodeID [32]byte

// EdgeID is hash(from, relation, to).
type EdgeID [32]byte

func (id NodeID) String() string { return hex.EncodeToString(id[:8]) }
func (id EdgeID) String() string { return hex.EncodeToString(id[:8]) }

// IsZero reports the zero identity, which no real node has.
func (id NodeID) IsZero() bool { return id == NodeID{} }

// writeField writes a length-prefixed string so that concatenation is
// unambiguous — without the prefix, ("ab","c") and ("a","bc") would hash alike.
func writeField(h interface{ Write([]byte) (int, error) }, s string) {
	var n [8]byte
	binary.BigEndian.PutUint64(n[:], uint64(len(s)))
	_, _ = h.Write(n[:])
	_, _ = h.Write([]byte(s))
}

// newNodeID derives a node identity from its type and *already canonicalized*
// key. Passing a raw key here is the bug that silently breaks convergence, so
// callers go through the applier rather than calling this directly.
func newNodeID(typeName, canonicalKey string) NodeID {
	h := sha256.New()
	writeField(h, typeName)
	writeField(h, canonicalKey)
	var id NodeID
	copy(id[:], h.Sum(nil))
	return id
}

func newEdgeID(from NodeID, rel string, to NodeID) EdgeID {
	h := sha256.New()
	_, _ = h.Write(from[:])
	writeField(h, rel)
	_, _ = h.Write(to[:])
	var id EdgeID
	copy(id[:], h.Sum(nil))
	return id
}

// Node is an admitted graph node. It carries no provenance, status or findings:
// those live in side tables, so that two identical scans produce identical
// content addresses.
type Node struct {
	ID    NodeID
	Type  *NodeType
	Key   string // canonical
	Props Props
}

// Edge is an admitted relation between two nodes.
type Edge struct {
	ID    EdgeID
	From  NodeID
	To    NodeID
	Rel   *Rel
	Props Props
}

// cidPrefix matches the store's: CIDv1, dag-cbor, sha2-256.
var cidPrefix = cid.Prefix{
	Version:  1,
	Codec:    cid.DagCBOR,
	MhType:   multihash.SHA2_256,
	MhLength: -1,
}

// Addressed returns the node's content-addressed encoding: the dag-cbor block
// and its CID. The form is a positional list — [type, key, [values...]] — not a
// map, so field names are not repeated per node and no key-sort step is needed
// to make encoding deterministic.
func (n *Node) Addressed() ([]byte, cid.Cid, error) {
	nb := basicnode.Prototype.List.NewBuilder()
	la, err := nb.BeginList(3)
	if err != nil {
		return nil, cid.Undef, err
	}
	if err := la.AssembleValue().AssignString(n.Type.name); err != nil {
		return nil, cid.Undef, err
	}
	if err := la.AssembleValue().AssignString(n.Key); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleProps(la.AssembleValue(), n.Props); err != nil {
		return nil, cid.Undef, err
	}
	if err := la.Finish(); err != nil {
		return nil, cid.Undef, err
	}
	return encodeBlock(nb.Build())
}

// Addressed returns the edge's content-addressed encoding, as
// [from, relation, to, [values...]].
func (e *Edge) Addressed() ([]byte, cid.Cid, error) {
	nb := basicnode.Prototype.List.NewBuilder()
	la, err := nb.BeginList(4)
	if err != nil {
		return nil, cid.Undef, err
	}
	if err := la.AssembleValue().AssignBytes(e.From[:]); err != nil {
		return nil, cid.Undef, err
	}
	if err := la.AssembleValue().AssignString(e.Rel.name); err != nil {
		return nil, cid.Undef, err
	}
	if err := la.AssembleValue().AssignBytes(e.To[:]); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleProps(la.AssembleValue(), e.Props); err != nil {
		return nil, cid.Undef, err
	}
	if err := la.Finish(); err != nil {
		return nil, cid.Undef, err
	}
	return encodeBlock(nb.Build())
}

// assembleProps writes every declared slot in declaration order, null for
// unset. Writing all slots rather than only the set ones keeps the encoding a
// function of the schema, so a node encodes identically wherever it was built.
func assembleProps(na datamodel.NodeAssembler, p Props) error {
	n := 0
	if p.sch != nil {
		n = len(p.slots)
	}
	la, err := na.BeginList(int64(n))
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		sl := p.slots[i]
		va := la.AssembleValue()
		if !sl.set {
			if err := va.AssignNull(); err != nil {
				return err
			}
			continue
		}
		if err := assembleValue(va, sl.val); err != nil {
			return err
		}
	}
	return la.Finish()
}

func assembleValue(na datamodel.NodeAssembler, v Value) error {
	switch v.kind {
	case KindString:
		return na.AssignString(v.str)
	case KindInt, KindTime:
		return na.AssignInt(v.num)
	case KindBool:
		return na.AssignBool(v.num != 0)
	case KindFloat:
		return na.AssignFloat(v.real)
	case KindBytes:
		return na.AssignBytes(v.raw)
	}
	return fmt.Errorf("graph: cannot encode value of kind %s", v.kind)
}

func encodeBlock(n datamodel.Node) ([]byte, cid.Cid, error) {
	var buf bytes.Buffer
	if err := dagcbor.Encode(n, &buf); err != nil {
		return nil, cid.Undef, err
	}
	block := buf.Bytes()
	c, err := cidPrefix.Sum(block)
	if err != nil {
		return nil, cid.Undef, err
	}
	return block, c, nil
}
