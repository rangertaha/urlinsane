// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package model

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/multiformats/go-multihash"
)

// Algorithm names recorded in a model's provenance.
const (
	// AlgorithmUniform is the untrained model — no corpus, no fitting.
	AlgorithmUniform = "uniform"
	// AlgorithmBaumWelch is unsupervised EM over recorded expansion traces.
	AlgorithmBaumWelch = "baum-welch"
	// AlgorithmManual is a hand-written model, for tests and bootstrapping.
	AlgorithmManual = "manual"
)

// blockVersion is the layout version of the encoded model. It is the first
// element of the block so that a decoder can reject a layout it does not
// understand before misreading positions.
const blockVersion = 1

// cidPrefix matches the store's and the graph's: CIDv1, dag-cbor, sha2-256.
var cidPrefix = cid.Prefix{
	Version:  1,
	Codec:    cid.DagCBOR,
	MhType:   multihash.SHA2_256,
	MhLength: -1,
}

// negInfSentinel stands in for -Inf inside an encoded block.
//
// dag-cbor is asymmetric about infinities: the encoder emits them happily and
// the decoder rejects them ("infinite float value rejected"), so a block
// containing a literal -Inf would encode and then fail to read back. An
// impossible cell therefore travels as the most negative finite float64, which
// is never a legitimate log-probability and round-trips bit-exactly.
const negInfSentinel = -math.MaxFloat64

// Provenance records where a model came from (§10.4). It is part of the
// addressed block, not a side file, because a model's identity includes what it
// was fitted on: two models with identical tables but different corpora are not
// interchangeable, and a plan pinning one must not silently accept the other.
type Provenance struct {
	// Algorithm is how the tables were produced.
	Algorithm string
	// Seed is the RNG seed used to initialize training, recorded so that
	// training reproduces exactly.
	Seed int64
	// Date is when training ran, in UTC. It is deliberately inside the block:
	// unlike a scan, where timestamps would break cross-run CID equality, a
	// model is a build artifact and when it was built is part of what it is.
	Date time.Time
	// Corpus are the CIDs of the trace blocks trained on, encoded as IPLD
	// links so the model block references its corpus in the DAG.
	Corpus []cid.Cid
	// Iterations is the number of EM iterations actually run.
	Iterations int
	// LogLikelihood is the corpus log-likelihood of the final tables. It is a
	// training diagnostic and never a user-facing number.
	LogLikelihood float64
}

// Addressed returns the model's dag-cbor block and its CID.
//
// The encoding is a positional list throughout — no maps anywhere — so there is
// no key-ordering question to get wrong and no per-cell key repeated across a
// table with thousands of entries. Field names live in this file, not in the
// block; the layout version is what protects readers from drift.
//
// The same model encodes to the same bytes and therefore the same CID on every
// run and every machine. That matters beyond tidiness: the model CID enters the
// plan hash (§10.4), so a plan pins a traversal only if this is stable.
func (h *HMM) Addressed() ([]byte, cid.Cid, error) {
	nb := basicnode.Prototype.List.NewBuilder()
	la, err := nb.BeginList(10)
	if err != nil {
		return nil, cid.Undef, err
	}
	if err := la.AssembleValue().AssignInt(blockVersion); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleStrings(la.AssembleValue(), h.states); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleInts(la.AssembleValue(), h.focus); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleStrings(la.AssembleValue(), h.rels); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleStrings(la.AssembleValue(), h.symbols); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleFloats(la.AssembleValue(), h.logInit); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleCube(la.AssembleValue(), h.logTrans); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleMatrix(la.AssembleValue(), h.logEmit); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleFloats(la.AssembleValue(), []float64{h.smooth.Init, h.smooth.Trans, h.smooth.Emit}); err != nil {
		return nil, cid.Undef, err
	}
	if err := assembleProvenance(la.AssembleValue(), h.prov); err != nil {
		return nil, cid.Undef, err
	}
	if err := la.Finish(); err != nil {
		return nil, cid.Undef, err
	}

	var buf bytes.Buffer
	if err := dagcbor.Encode(nb.Build(), &buf); err != nil {
		return nil, cid.Undef, err
	}
	block := buf.Bytes()
	c, err := cidPrefix.Sum(block)
	if err != nil {
		return nil, cid.Undef, err
	}
	return block, c, nil
}

// CID returns the model's content address. It is what `--model NAME@cid` pins
// and what enters the plan hash.
func (h *HMM) CID() (cid.Cid, error) {
	_, c, err := h.Addressed()
	return c, err
}

// Decode rebuilds a model from its dag-cbor block.
//
// It installs the log tables verbatim rather than re-normalizing them, so an
// encode/decode round trip is bit-exact and the CID is unchanged. Re-deriving
// the tables would be a slow way to introduce a rounding difference that
// changes a model's identity every time it is loaded.
func Decode(block []byte) (*HMM, error) {
	nb := basicnode.Prototype.List.NewBuilder()
	if err := dagcbor.Decode(nb, bytes.NewReader(block)); err != nil {
		return nil, fmt.Errorf("model: decode block: %w", err)
	}
	n := nb.Build()
	if n.Length() != 10 {
		return nil, fmt.Errorf("model: block has %d fields, want 10", n.Length())
	}
	v, err := readInt(n, 0)
	if err != nil {
		return nil, err
	}
	if v != blockVersion {
		return nil, fmt.Errorf("model: block layout version %d, want %d", v, blockVersion)
	}

	h := &HMM{}
	if h.states, err = readStrings(n, 1); err != nil {
		return nil, err
	}
	focus, err := readInts(n, 2)
	if err != nil {
		return nil, err
	}
	if h.rels, err = readStrings(n, 3); err != nil {
		return nil, err
	}
	if h.symbols, err = readStrings(n, 4); err != nil {
		return nil, err
	}
	if h.logInit, err = readFloats(n, 5); err != nil {
		return nil, err
	}
	if h.logTrans, err = readCube(n, 6); err != nil {
		return nil, err
	}
	if h.logEmit, err = readMatrix(n, 7); err != nil {
		return nil, err
	}
	sm, err := readFloats(n, 8)
	if err != nil {
		return nil, err
	}
	if len(sm) != 3 {
		return nil, fmt.Errorf("model: smoothing has %d entries, want 3", len(sm))
	}
	h.smooth = Smoothing{Init: sm[0], Trans: sm[1], Emit: sm[2]}
	if h.prov, err = readProvenance(n, 9); err != nil {
		return nil, err
	}

	if h.stateIdx, err = index("state", h.states); err != nil {
		return nil, err
	}
	if h.relIdx, err = index("relation", h.rels); err != nil {
		return nil, err
	}
	if h.symIdx, err = index("symbol", h.symbols); err != nil {
		return nil, err
	}
	names := make([]string, 0, len(focus))
	for _, i := range focus {
		if i < 0 || i >= len(h.states) {
			return nil, fmt.Errorf("model: focus index %d out of range", i)
		}
		names = append(names, h.states[i])
	}
	if err := h.setFocus(names); err != nil {
		return nil, err
	}
	if err := h.checkShape(); err != nil {
		return nil, err
	}
	oovRel, ok := h.relIdx[OOVRelation]
	if !ok {
		return nil, fmt.Errorf("model: block has no %q relation", OOVRelation)
	}
	oovSym, ok := h.symIdx[OOVSymbol]
	if !ok {
		return nil, fmt.Errorf("model: block has no %q symbol", OOVSymbol)
	}
	h.oovRel, h.oovSym = oovRel, oovSym
	return h, nil
}

// checkShape rejects tables whose dimensions disagree with the alphabets. A
// decoder that trusted the block would index out of range deep inside Forward,
// mid-scan, rather than at load.
func (h *HMM) checkShape() error {
	ns, nr, nk := len(h.states), len(h.rels), len(h.symbols)
	if ns == 0 {
		return fmt.Errorf("model: block has no states")
	}
	if len(h.logInit) != ns {
		return fmt.Errorf("model: initial distribution has %d entries, want %d", len(h.logInit), ns)
	}
	if len(h.logTrans) != nr {
		return fmt.Errorf("model: %d transition tables, want %d", len(h.logTrans), nr)
	}
	for ri, t := range h.logTrans {
		if len(t) != ns {
			return fmt.Errorf("model: transition table %d has %d rows, want %d", ri, len(t), ns)
		}
		for i, row := range t {
			if len(row) != ns {
				return fmt.Errorf("model: transition row %d of table %d has %d entries, want %d", i, ri, len(row), ns)
			}
		}
	}
	if len(h.logEmit) != ns {
		return fmt.Errorf("model: emission table has %d rows, want %d", len(h.logEmit), ns)
	}
	for i, row := range h.logEmit {
		if len(row) != nk {
			return fmt.Errorf("model: emission row %d has %d entries, want %d", i, len(row), nk)
		}
	}
	return nil
}

// --- assembly -------------------------------------------------------------

func assembleStrings(na datamodel.NodeAssembler, ss []string) error {
	la, err := na.BeginList(int64(len(ss)))
	if err != nil {
		return err
	}
	for _, s := range ss {
		if err := la.AssembleValue().AssignString(s); err != nil {
			return err
		}
	}
	return la.Finish()
}

func assembleInts(na datamodel.NodeAssembler, is []int) error {
	la, err := na.BeginList(int64(len(is)))
	if err != nil {
		return err
	}
	for _, i := range is {
		if err := la.AssembleValue().AssignInt(int64(i)); err != nil {
			return err
		}
	}
	return la.Finish()
}

func assembleFloats(na datamodel.NodeAssembler, fs []float64) error {
	la, err := na.BeginList(int64(len(fs)))
	if err != nil {
		return err
	}
	for _, f := range fs {
		if err := la.AssembleValue().AssignFloat(encodeFloat(f)); err != nil {
			return err
		}
	}
	return la.Finish()
}

func assembleMatrix(na datamodel.NodeAssembler, m [][]float64) error {
	la, err := na.BeginList(int64(len(m)))
	if err != nil {
		return err
	}
	for _, row := range m {
		if err := assembleFloats(la.AssembleValue(), row); err != nil {
			return err
		}
	}
	return la.Finish()
}

func assembleCube(na datamodel.NodeAssembler, c [][][]float64) error {
	la, err := na.BeginList(int64(len(c)))
	if err != nil {
		return err
	}
	for _, m := range c {
		if err := assembleMatrix(la.AssembleValue(), m); err != nil {
			return err
		}
	}
	return la.Finish()
}

func assembleProvenance(na datamodel.NodeAssembler, p Provenance) error {
	la, err := na.BeginList(6)
	if err != nil {
		return err
	}
	if err := la.AssembleValue().AssignString(p.Algorithm); err != nil {
		return err
	}
	if err := la.AssembleValue().AssignInt(p.Seed); err != nil {
		return err
	}
	// A zero time encodes as 0 rather than as its nanosecond representation,
	// which is a large negative number that reads as a real date.
	var when int64
	if !p.Date.IsZero() {
		when = p.Date.UTC().UnixNano()
	}
	if err := la.AssembleValue().AssignInt(when); err != nil {
		return err
	}
	cl, err := la.AssembleValue().BeginList(int64(len(p.Corpus)))
	if err != nil {
		return err
	}
	for _, c := range p.Corpus {
		if err := cl.AssembleValue().AssignLink(cidlink.Link{Cid: c}); err != nil {
			return err
		}
	}
	if err := cl.Finish(); err != nil {
		return err
	}
	if err := la.AssembleValue().AssignInt(int64(p.Iterations)); err != nil {
		return err
	}
	if err := la.AssembleValue().AssignFloat(encodeFloat(p.LogLikelihood)); err != nil {
		return err
	}
	return la.Finish()
}

// encodeFloat maps -Inf onto the sentinel and refuses the values that have no
// meaning in a probability table at all.
func encodeFloat(f float64) float64 {
	switch {
	case math.IsInf(f, -1):
		return negInfSentinel
	case math.IsNaN(f), math.IsInf(f, 1):
		// A NaN in a table is a bug upstream; encoding it would hide the bug
		// behind a decode failure much later. Clamp to the sentinel so the
		// block stays readable and the value stays obviously impossible.
		return negInfSentinel
	}
	return f
}

func decodeFloat(f float64) float64 {
	if f == negInfSentinel {
		return LogZero
	}
	return f
}

// --- reading --------------------------------------------------------------

func at(n datamodel.Node, i int) (datamodel.Node, error) {
	v, err := n.LookupByIndex(int64(i))
	if err != nil {
		return nil, fmt.Errorf("model: field %d: %w", i, err)
	}
	return v, nil
}

func readInt(n datamodel.Node, i int) (int64, error) {
	v, err := at(n, i)
	if err != nil {
		return 0, err
	}
	return v.AsInt()
}

func readStrings(n datamodel.Node, i int) ([]string, error) {
	v, err := at(n, i)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, v.Length())
	it := v.ListIterator()
	for it != nil && !it.Done() {
		_, e, err := it.Next()
		if err != nil {
			return nil, err
		}
		s, err := e.AsString()
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func readInts(n datamodel.Node, i int) ([]int, error) {
	v, err := at(n, i)
	if err != nil {
		return nil, err
	}
	out := make([]int, 0, v.Length())
	it := v.ListIterator()
	for it != nil && !it.Done() {
		_, e, err := it.Next()
		if err != nil {
			return nil, err
		}
		x, err := e.AsInt()
		if err != nil {
			return nil, err
		}
		out = append(out, int(x))
	}
	return out, nil
}

func readFloats(n datamodel.Node, i int) ([]float64, error) {
	v, err := at(n, i)
	if err != nil {
		return nil, err
	}
	return floatsOf(v)
}

func floatsOf(v datamodel.Node) ([]float64, error) {
	out := make([]float64, 0, v.Length())
	it := v.ListIterator()
	for it != nil && !it.Done() {
		_, e, err := it.Next()
		if err != nil {
			return nil, err
		}
		f, err := e.AsFloat()
		if err != nil {
			return nil, err
		}
		out = append(out, decodeFloat(f))
	}
	return out, nil
}

func readMatrix(n datamodel.Node, i int) ([][]float64, error) {
	v, err := at(n, i)
	if err != nil {
		return nil, err
	}
	return matrixOf(v)
}

func matrixOf(v datamodel.Node) ([][]float64, error) {
	out := make([][]float64, 0, v.Length())
	it := v.ListIterator()
	for it != nil && !it.Done() {
		_, e, err := it.Next()
		if err != nil {
			return nil, err
		}
		row, err := floatsOf(e)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, nil
}

func readCube(n datamodel.Node, i int) ([][][]float64, error) {
	v, err := at(n, i)
	if err != nil {
		return nil, err
	}
	out := make([][][]float64, 0, v.Length())
	it := v.ListIterator()
	for it != nil && !it.Done() {
		_, e, err := it.Next()
		if err != nil {
			return nil, err
		}
		m, err := matrixOf(e)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
}

func readProvenance(n datamodel.Node, i int) (Provenance, error) {
	var p Provenance
	v, err := at(n, i)
	if err != nil {
		return p, err
	}
	if v.Length() != 6 {
		return p, fmt.Errorf("model: provenance has %d fields, want 6", v.Length())
	}
	if p.Algorithm, err = mustString(v, 0); err != nil {
		return p, err
	}
	if p.Seed, err = readInt(v, 1); err != nil {
		return p, err
	}
	when, err := readInt(v, 2)
	if err != nil {
		return p, err
	}
	if when != 0 {
		p.Date = time.Unix(0, when).UTC()
	}
	corpus, err := at(v, 3)
	if err != nil {
		return p, err
	}
	it := corpus.ListIterator()
	for it != nil && !it.Done() {
		_, e, err := it.Next()
		if err != nil {
			return p, err
		}
		lnk, err := e.AsLink()
		if err != nil {
			return p, err
		}
		cl, ok := lnk.(cidlink.Link)
		if !ok {
			return p, fmt.Errorf("model: corpus entry is not a CID link")
		}
		p.Corpus = append(p.Corpus, cl.Cid)
	}
	iter, err := readInt(v, 4)
	if err != nil {
		return p, err
	}
	p.Iterations = int(iter)
	ll, err := at(v, 5)
	if err != nil {
		return p, err
	}
	f, err := ll.AsFloat()
	if err != nil {
		return p, err
	}
	p.LogLikelihood = decodeFloat(f)
	return p, nil
}

func mustString(n datamodel.Node, i int) (string, error) {
	v, err := at(n, i)
	if err != nil {
		return "", err
	}
	return v.AsString()
}
