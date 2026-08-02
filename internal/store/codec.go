// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package store persists a graph.Graph as content-addressed dag-cbor
// blocks: one block per node, one per edge, one side-table block and one scan
// root linking them all.
//
// The whole point of the scheme is cross-run diffing, and that only works if
// two identical scans produce identical CIDs. Three rules make it true and are
// load-bearing everywhere in this package:
//
//   - Provenance, per-pair status, findings, scheduler state (depth, closure)
//     and the truncation ledger live in the side block, never inside a node or
//     edge block. Put a round number or an operator id in the addressed form
//     and every node differs on every run (docs/DESIGN.md §1.2).
//   - Everything is encoded as a positional list, never a map. There are no
//     keys to sort and no map to iterate, so the encoding is deterministic by
//     construction rather than by a sort step someone can forget (§1.3).
//   - Every collection written here comes from an accessor with a defined
//     order — graph.Nodes(), graph.Edges(), graph.Ledger() — or is sorted by
//     this package before encoding. Nothing iterates a Go map on the way to a
//     block.
//
// Node and edge blocks are produced by graph.Node.Addressed and
// graph.Edge.Addressed; this package supplies the inverse and everything
// around it.
package store

import (
	"bytes"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/multiformats/go-multihash"
)

// FormatVersion is the first element of every block this package authors
// itself (the side block and the scan root). It is not carried by node and
// edge blocks: those are the addressed form and adding a version byte to them
// would invalidate every content address already in a store.
const FormatVersion = 1

// cidPrefix matches internal/graph: CIDv1, dag-cbor, sha2-256. Both must agree
// or a node's CID would depend on who hashed it.
var cidPrefix = cid.Prefix{
	Version:  1,
	Codec:    cid.DagCBOR,
	MhType:   multihash.SHA2_256,
	MhLength: -1,
}

// encodeBlock serialises an IPLD node to dag-cbor and returns it with its CID.
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

// decodeBlock parses a dag-cbor block into an untyped IPLD node.
func decodeBlock(block []byte) (datamodel.Node, error) {
	nb := basicnode.Prototype.Any.NewBuilder()
	if err := dagcbor.Decode(nb, bytes.NewReader(block)); err != nil {
		return nil, err
	}
	return nb.Build(), nil
}

// enc builds a positional list. It carries the first error rather than
// returning one per write, so a twelve-field row reads as twelve lines instead
// of thirty-six.
type enc struct {
	la  datamodel.ListAssembler
	err error
}

// encodeList builds a dag-cbor list block. n is a size hint; passing the exact
// length keeps the CBOR header definite-length, which is required for the
// encoding to be canonical.
func encodeList(n int, fn func(*enc)) ([]byte, cid.Cid, error) {
	nb := basicnode.Prototype.List.NewBuilder()
	la, err := nb.BeginList(int64(n))
	if err != nil {
		return nil, cid.Undef, err
	}
	e := &enc{la: la}
	fn(e)
	if e.err != nil {
		return nil, cid.Undef, e.err
	}
	if err := la.Finish(); err != nil {
		return nil, cid.Undef, err
	}
	return encodeBlock(nb.Build())
}

func (e *enc) fail(err error) {
	if e.err == nil {
		e.err = err
	}
}

func (e *enc) str(s string) {
	if e.err == nil {
		e.fail(e.la.AssembleValue().AssignString(s))
	}
}

func (e *enc) i64(i int64) {
	if e.err == nil {
		e.fail(e.la.AssembleValue().AssignInt(i))
	}
}

func (e *enc) f64(f float64) {
	if e.err == nil {
		e.fail(e.la.AssembleValue().AssignFloat(f))
	}
}

func (e *enc) flag(b bool) {
	if e.err == nil {
		e.fail(e.la.AssembleValue().AssignBool(b))
	}
}

func (e *enc) raw(b []byte) {
	if e.err == nil {
		e.fail(e.la.AssembleValue().AssignBytes(b))
	}
}

func (e *enc) link(c cid.Cid) {
	if e.err == nil {
		e.fail(e.la.AssembleValue().AssignLink(cidlink.Link{Cid: c}))
	}
}

// sub writes a nested positional list.
func (e *enc) sub(n int, fn func(*enc)) {
	if e.err != nil {
		return
	}
	la, err := e.la.AssembleValue().BeginList(int64(n))
	if err != nil {
		e.fail(err)
		return
	}
	inner := &enc{la: la}
	fn(inner)
	if inner.err != nil {
		e.fail(inner.err)
		return
	}
	e.fail(la.Finish())
}

// dec reads a positional list. Every dec derived from a root dec shares one
// error cell, so a chain of accessors can be written without checking each
// step; the single check at the end catches any failure along the way.
type dec struct {
	n    datamodel.Node
	errp *error
}

func newDec(n datamodel.Node) *dec {
	var err error
	return &dec{n: n, errp: &err}
}

func (d *dec) fail(format string, a ...any) {
	if *d.errp == nil {
		*d.errp = fmt.Errorf(format, a...)
	}
}

func (d *dec) err() error { return *d.errp }

func (d *dec) len() int {
	if *d.errp != nil {
		return 0
	}
	if d.n.Kind() != datamodel.Kind_List {
		d.fail("store: expected a list, got %s", d.n.Kind())
		return 0
	}
	return int(d.n.Length())
}

// expect asserts an exact arity. Blocks are positional, so a wrong length means
// the block was written by an incompatible version and every subsequent read
// would be silently misaligned.
func (d *dec) expect(what string, n int) *dec {
	if got := d.len(); *d.errp == nil && got != n {
		d.fail("store: %s has %d elements, want %d", what, got, n)
	}
	return d
}

func (d *dec) at(i int) *dec {
	child := &dec{n: d.n, errp: d.errp}
	if *d.errp != nil {
		return child
	}
	if d.n.Kind() != datamodel.Kind_List {
		d.fail("store: expected a list, got %s", d.n.Kind())
		return child
	}
	c, err := d.n.LookupByIndex(int64(i))
	if err != nil {
		d.fail("store: element %d: %w", i, err)
		return child
	}
	child.n = c
	return child
}

// each yields every element of a list in order.
func (d *dec) each(fn func(*dec)) {
	n := d.len()
	for i := 0; i < n; i++ {
		fn(d.at(i))
	}
}

func (d *dec) str() string {
	if *d.errp != nil {
		return ""
	}
	s, err := d.n.AsString()
	if err != nil {
		d.fail("store: expected a string: %w", err)
		return ""
	}
	return s
}

func (d *dec) i64() int64 {
	if *d.errp != nil {
		return 0
	}
	i, err := d.n.AsInt()
	if err != nil {
		d.fail("store: expected an int: %w", err)
		return 0
	}
	return i
}

func (d *dec) f64() float64 {
	if *d.errp != nil {
		return 0
	}
	f, err := d.n.AsFloat()
	if err != nil {
		d.fail("store: expected a float: %w", err)
		return 0
	}
	return f
}

func (d *dec) flag() bool {
	if *d.errp != nil {
		return false
	}
	b, err := d.n.AsBool()
	if err != nil {
		d.fail("store: expected a bool: %w", err)
		return false
	}
	return b
}

func (d *dec) raw() []byte {
	if *d.errp != nil {
		return nil
	}
	b, err := d.n.AsBytes()
	if err != nil {
		d.fail("store: expected bytes: %w", err)
		return nil
	}
	return b
}

func (d *dec) link() cid.Cid {
	if *d.errp != nil {
		return cid.Undef
	}
	l, err := d.n.AsLink()
	if err != nil {
		d.fail("store: expected a link: %w", err)
		return cid.Undef
	}
	cl, ok := l.(cidlink.Link)
	if !ok {
		d.fail("store: link is not a CID link")
		return cid.Undef
	}
	return cl.Cid
}

// id32 reads a 32-byte NodeID or EdgeID.
func (d *dec) id32() [32]byte {
	var out [32]byte
	b := d.raw()
	if *d.errp != nil {
		return out
	}
	if len(b) != 32 {
		d.fail("store: identity is %d bytes, want 32", len(b))
		return out
	}
	copy(out[:], b)
	return out
}
