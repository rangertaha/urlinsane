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

package graph

import "sort"

// Status is the terminal outcome of one (node, operator) pair. How a lookup
// failed is itself the finding: NXDOMAIN proves a name is free, a timeout
// proves nothing at all, and collapsing the two would discard the signal a
// squatting scanner exists to collect.
type Status uint8

const (
	// StatusOK means the operator learned something positive.
	StatusOK Status = iota + 1
	// StatusEmpty means it authoritatively determined absence.
	StatusEmpty
	// StatusFailed means the lookup itself broke.
	StatusFailed
	// StatusTimeout means nothing was learned.
	StatusTimeout
	// StatusSkipped means it was never attempted. Unlike the others this is
	// not terminal: a pair gated off at one barrier may run at a later one.
	StatusSkipped
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "ok"
	case StatusEmpty:
		return "empty"
	case StatusFailed:
		return "failed"
	case StatusTimeout:
		return "timeout"
	case StatusSkipped:
		return "skipped"
	}
	return "invalid"
}

// Terminal reports whether recording this status closes the pair. Skipped never
// does; recording it terminally would make the first belief gate permanent.
func (s Status) Terminal() bool { return s != StatusSkipped }

// Provenance records who asserted something and when. It is kept out of the
// content-addressed form so identical scans produce identical CIDs.
type Provenance struct {
	Operator string
	Round    int
}

// Assertion is one operator's claim about a field, retained whether or not it
// won the merge. Disagreement between two sources is signal, not noise.
type Assertion struct {
	Field string
	Value Value
	By    Provenance
	Won   bool
}

// Reason explains why a candidate was declined.
type Reason uint8

const (
	// ReasonBelief — below the execution model's threshold.
	ReasonBelief Reason = iota + 1
	// ReasonBudget — a per-type or global node budget was exhausted.
	ReasonBudget
	// ReasonFrontier — the in-flight candidate bound was reached.
	ReasonFrontier
	// ReasonRoundCap — the round cap was hit.
	ReasonRoundCap
	// ReasonDeadline — the round deadline passed before it could be admitted.
	ReasonDeadline
)

func (r Reason) String() string {
	switch r {
	case ReasonBelief:
		return "belief"
	case ReasonBudget:
		return "budget"
	case ReasonFrontier:
		return "frontier"
	case ReasonRoundCap:
		return "round-cap"
	case ReasonDeadline:
		return "deadline"
	}
	return "invalid"
}

// LedgerRow records a candidate the engine declined to admit. The ledger is
// reported like any other section — a truncated graph that reads as complete is
// a correctness bug — and it doubles as a denylist, so a later operator
// re-emitting the same candidate cannot quietly resurrect it.
type LedgerRow struct {
	Type   string
	Key    string // canonical
	Depth  int
	Belief float64
	Reason Reason
	By     Provenance
}

// RejectKind classifies an invariant violation. Unlike a ledger row, a
// rejection does not deny the candidate forever: a node refused as the source
// of a VARIANT_OF edge may still be admitted legitimately by another edge.
type RejectKind uint8

const (
	// RejectCanonical — the key could not be canonicalized.
	RejectCanonical RejectKind = iota + 1
	// RejectUnknownType — no such node type is registered.
	RejectUnknownType
	// RejectUnknownRel — no such relation is registered.
	RejectUnknownRel
	// RejectUnknownField — the field is not declared by the target's type.
	RejectUnknownField
	// RejectKindMismatch — the value's kind is not the field's kind.
	RejectKindMismatch
	// RejectClosure — a variant edge whose source is outside the seed closure.
	RejectClosure
	// RejectSelfVariant — a variant edge whose source and target canonicalize
	// to the same node.
	RejectSelfVariant
	// RejectScope — a variant edge whose source type the run's scope excludes.
	// Distinct from RejectClosure: the node is a legitimate variant root, the
	// user simply asked for a narrower scan, and conflating the two would make
	// a scoped run look like an invariant violation.
	RejectScope
	// RejectDenied — the candidate is in the truncation ledger.
	RejectDenied
	// RejectMissingNode — an edge or prop referenced a node that was refused.
	RejectMissingNode
)

func (k RejectKind) String() string {
	switch k {
	case RejectCanonical:
		return "canonicalization"
	case RejectUnknownType:
		return "unknown-type"
	case RejectUnknownRel:
		return "unknown-relation"
	case RejectUnknownField:
		return "unknown-field"
	case RejectKindMismatch:
		return "kind-mismatch"
	case RejectClosure:
		return "outside-seed-closure"
	case RejectSelfVariant:
		return "self-variant"
	case RejectScope:
		return "outside-scope"
	case RejectDenied:
		return "declined-earlier"
	case RejectMissingNode:
		return "missing-node"
	}
	return "invalid"
}

// Rejection records one refused item from a delta.
type Rejection struct {
	Kind   RejectKind
	Type   string
	Key    string
	Rel    string
	Field  string
	Detail string
	By     Provenance
}

// ledgerKey identifies a declined candidate.
type ledgerKey struct {
	typ string
	key string
}

// sortLedger orders rows the way the report renders them: (type, key, reason).
// Sorting by (type, key) alone is not total — one candidate can be declined for
// different reasons across rounds.
func sortLedger(rows []LedgerRow) {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Type != rows[j].Type {
			return rows[i].Type < rows[j].Type
		}
		if rows[i].Key != rows[j].Key {
			return rows[i].Key < rows[j].Key
		}
		return rows[i].Reason < rows[j].Reason
	})
}
