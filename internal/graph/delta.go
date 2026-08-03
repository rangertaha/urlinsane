// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

import (
	"fmt"
	"sort"
)

// NodeRef is how an operator names a node. It carries a *raw* key: operators
// cannot compute a NodeID because canonicalization belongs to the registry and
// the applier, and letting a plugin mint an identity is how convergence quietly
// breaks.
type NodeRef struct {
	Type string
	Key  string
}

// EdgeRef is how an operator names an edge.
type EdgeRef struct {
	From NodeRef
	Rel  string
	To   NodeRef
}

// PropSet is a single field assertion against a node or an edge. Exactly one of
// Node or Edge must be set. Every prop write goes through this so the field's
// merge policy applies uniformly.
type PropSet struct {
	Node  *NodeRef
	Edge  *EdgeRef
	Field string
	Value Value
}

// Delta is everything one operator produced. Operators never mutate the graph;
// they return a Delta and the applier is the single writer. Deltas are additive
// only — nothing removes a node, edge or prop — which is what makes the graph
// monotonic within a run and a delta safely replayable.
type Delta struct {
	Nodes []NodeRef
	Edges []EdgeRef
	Props []PropSet
}

// Result reports what an Apply admitted and what it refused.
type Result struct {
	Nodes    []NodeID
	Edges    []EdgeID
	Changed  []NodeID
	Rejected []Rejection
}

// Graph is the applier and the single writer. It is not safe for concurrent
// use; the scheduler serialises deltas through it.
type Graph struct {
	reg *Registry

	nodes map[NodeID]*Node
	edges map[EdgeID]*Edge

	order []NodeID // admission order, for deterministic iteration
	eord  []EdgeID

	depth   map[NodeID]int
	provis  map[NodeID]bool // depth is the subject's, not yet derived from an in-edge
	closure map[NodeID]bool
	scope   map[string]bool // nameable types that may be varied; nil means all
	seed    NodeID

	assertions map[NodeID][]Assertion
	eassert    map[EdgeID][]Assertion
	status     map[statusKey]Status

	observers   map[string]bool // operator ids whose status attests existence
	ledger      map[ledgerKey]LedgerRow
	rejections  []Rejection
	truncations []RunTruncation

	model      BeliefModel
	belief     map[NodeID]float64
	bstate     map[NodeID]State
	parent     map[NodeID]parentRef
	candidates map[NodeID][]parentRef

	budgets  Budgets
	counts   map[string]int
	total    int
	frontier int // per-round admission cap; 0 is unbounded
	admitted int // admissions so far this round, reset at the barrier

	findings []Finding
	scores   map[NodeID]map[string]float64
}

// Budgets cap how many nodes may be admitted. Zero means unbounded.
type Budgets struct {
	Global  int
	PerType map[string]int
}

// RunTruncation is a run-level limit that bound expansion, as opposed to a
// per-candidate ledger row. Hitting the round cap is reported like any other
// truncation: a truncated graph that reads as complete is a correctness bug.
type RunTruncation struct {
	Reason Reason
	Round  int
	Detail string
}

type statusKey struct {
	node NodeID
	op   string
}

// New returns an empty graph bound to a registry.
func New(reg *Registry) *Graph {
	return &Graph{
		reg:        reg,
		nodes:      map[NodeID]*Node{},
		edges:      map[EdgeID]*Edge{},
		depth:      map[NodeID]int{},
		provis:     map[NodeID]bool{},
		closure:    map[NodeID]bool{},
		assertions: map[NodeID][]Assertion{},
		eassert:    map[EdgeID][]Assertion{},
		status:     map[statusKey]Status{},
		ledger:     map[ledgerKey]LedgerRow{},
		observers:  map[string]bool{},
		model:      uniformModel{},
		belief:     map[NodeID]float64{},
		bstate:     map[NodeID]State{},
		parent:     map[NodeID]parentRef{},
		candidates: map[NodeID][]parentRef{},
		counts:     map[string]int{},
		scores:     map[NodeID]map[string]float64{},
	}
}

// SetBudgets caps admissions. Exceeding a budget declines the candidate to the
// truncation ledger rather than dropping it silently.
func (g *Graph) SetBudgets(b Budgets) { g.budgets = b }

// SetScope restricts which node types may root a variant. An empty or nil list
// means every Nameable type in the seed closure, which is the default.
//
// This is the CLI's scope positional (§12), and it is enforced here rather than
// only at dispatch for the same reason the closure rule is: dispatch-side
// gating is an optimization an operator can be written around, while an
// applier rejection is an invariant it cannot. Scope that was merely validated
// and printed — compiled into the plan, shown by --explain, and then ignored by
// execution — would make `typo username bob@example.com` silently identical to
// the unscoped run, which is the one thing the positional exists to prevent.
func (g *Graph) SetScope(types []string) {
	if len(types) == 0 {
		g.scope = nil
		return
	}
	g.scope = make(map[string]bool, len(types))
	for _, t := range types {
		g.scope[t] = true
	}
}

// inScope reports whether a node's type may root a variant.
func (g *Graph) inScope(id NodeID) bool {
	if g.scope == nil {
		return true
	}
	n, ok := g.nodes[id]
	return ok && g.scope[n.Type.name]
}

// InScope reports whether a node may root a variant under the current scope.
// The scheduler uses it to skip dispatching variant operators it knows will be
// rejected; the rejection in the applier remains the invariant.
func (g *Graph) InScope(id NodeID) bool { return g.inScope(id) }

// NoteTruncation records a run-level limit.
func (g *Graph) NoteTruncation(r Reason, round int, detail string) {
	g.truncations = append(g.truncations, RunTruncation{Reason: r, Round: round, Detail: detail})
}

// Truncations returns run-level limits that bound this expansion.
func (g *Graph) Truncations() []RunTruncation { return g.truncations }

// noteStoppedEarly records that expansion ended with work still eligible.
//
// It was called declineFrontier and took a Reason, which read as though it
// declined candidates for the frontier cap; its one caller passes
// ReasonRoundCap, and the frontier is enforced in admitNode. The name was the
// whole reason --frontier looked implemented from here.
func (g *Graph) noteStoppedEarly(r Reason, by Provenance) {
	g.NoteTruncation(r, by.Round, "expansion stopped before the frontier was exhausted")
}

// endRound is the barrier hook for limits that are counted per round.
//
// Budgets are cumulative and applied at admission, so there is nothing to undo
// for them here. The frontier is not: it caps admissions *within* a round, so
// its counter resets at the barrier and nowhere else — reset it at the start of
// a dispatch instead and a round with two operator generations would get two
// frontiers.
func (g *Graph) endRound(Provenance) { g.admitted = 0 }

// SetFrontier caps how many nodes may be admitted in one round. Zero is
// unbounded.
//
// The cap applies in admission order, which the scheduler has already made
// deterministic — work is sorted by (depth, type, key, operator) and applied in
// that order rather than in completion order — so the same scan truncates at
// the same place on every run.
//
// DESIGN §8 specifies the survivors as the prefix of candidates sorted by
// (-belief, depth, type, key). While the belief model is uniform those two
// orderings are the same list, so this is that rule with the constant factored
// out. A model that returns anything other than 1 makes them differ, and at
// that point the frontier has to become a queue drained at the barrier rather
// than a counter checked at admission.
func (g *Graph) SetFrontier(n int) { g.frontier = n }

// overFrontier reports whether this round has already admitted its allowance.
func (g *Graph) overFrontier() bool {
	return g.frontier > 0 && g.admitted >= g.frontier
}

// overBudget reports whether admitting one more node of this type would exceed
// a budget.
func (g *Graph) overBudget(typeName string) bool {
	if g.budgets.Global > 0 && g.total >= g.budgets.Global {
		return true
	}
	if n, ok := g.budgets.PerType[typeName]; ok && n > 0 && g.counts[typeName] >= n {
		return true
	}
	return false
}

// Seed admits the target node at depth 0 and opens the seed closure. It is the
// only admission that does not descend from an existing node.
func (g *Graph) Seed(typeName, rawKey string) (NodeID, error) {
	if !g.seed.IsZero() {
		return NodeID{}, fmt.Errorf("graph: seed already set")
	}
	t, ok := g.reg.Type(typeName)
	if !ok {
		return NodeID{}, fmt.Errorf("graph: unknown node type %q", typeName)
	}
	key, err := t.canonical(rawKey)
	if err != nil {
		return NodeID{}, fmt.Errorf("graph: cannot canonicalize %q as %s: %w", rawKey, typeName, err)
	}
	id := newNodeID(t.name, key)
	g.nodes[id] = &Node{ID: id, Type: t, Key: key, Props: newProps(t.sch)}
	g.order = append(g.order, id)
	g.depth[id] = 0
	g.closure[id] = true
	g.seed = id
	// The seed counts. It is not budget-*checked* — refusing to admit the
	// target because the budget is zero would be absurd — but leaving it out of
	// the tally makes g.total disagree with len(g.Nodes()) forever after, so
	// `--budget 5` admits six nodes and the off-by-one is invisible.
	g.counts[t.name]++
	g.total++
	return id, nil
}

// Seed returns the target node. The seed cannot be inferred from the graph:
// structural edges cost no depth, so a composite target puts several nodes at
// depth 0 inside the closure and "the one at depth 0" names no single node.
func (g *Graph) SeedID() NodeID { return g.seed }

// Decline records a candidate the engine refused to admit, and denies it
// thereafter. Truncation of every kind lands here, which is what makes
// "pruning is irreversible" true even when a second operator finds the same
// candidate later.
func (g *Graph) Decline(typeName, rawKey string, depth int, belief float64, r Reason, by Provenance) error {
	t, ok := g.reg.Type(typeName)
	if !ok {
		return fmt.Errorf("graph: unknown node type %q", typeName)
	}
	key, err := t.canonical(rawKey)
	if err != nil {
		return fmt.Errorf("graph: cannot canonicalize %q as %s: %w", rawKey, typeName, err)
	}
	k := ledgerKey{typ: t.name, key: key}
	if _, exists := g.ledger[k]; exists {
		return nil // first decline wins; the row is kept
	}
	g.ledger[k] = LedgerRow{Type: t.name, Key: key, Depth: depth, Belief: belief, Reason: r, By: by}
	return nil
}

// Apply applies one operator's delta. subject is the node the operator ran on;
// nodes it emits inherit that depth, and edges adjust it by their class.
func (g *Graph) Apply(by Provenance, subject NodeID, d Delta) Result {
	var res Result
	if _, ok := g.nodes[subject]; !ok && !subject.IsZero() {
		res.Rejected = append(res.Rejected, Rejection{
			Kind: RejectMissingNode, Detail: "subject not in graph", By: by,
		})
		return res
	}
	base := g.depth[subject]

	// Resolved refs for this delta, so an edge or prop can name a node the
	// same delta introduced.
	resolved := map[NodeRef]NodeID{}

	admit := func(ref NodeRef) (NodeID, bool) {
		if id, done := resolved[ref]; done {
			return id, !id.IsZero()
		}
		id, rej, ok := g.admitNode(ref, base, by)
		if !ok {
			resolved[ref] = NodeID{}
			res.Rejected = append(res.Rejected, rej)
			return NodeID{}, false
		}
		resolved[ref] = id
		if _, seen := g.nodes[id]; seen {
			res.Nodes = append(res.Nodes, id)
		}
		return id, true
	}

	for _, ref := range d.Nodes {
		admit(ref)
	}

	for _, eref := range d.Edges {
		id, ok := g.admitEdge(eref, admit, by, &res)
		if ok {
			res.Edges = append(res.Edges, id)
		}
	}

	for _, ps := range d.Props {
		g.applyProp(ps, admit, by, &res)
	}

	// Rejections are recorded run-wide as well as returned. Returning them only
	// in the Result made Graph.Rejections() an accessor that always answered
	// "nothing was refused" — worse than not existing, because a caller asking
	// what the applier turned away gets a confident wrong answer.
	g.rejections = append(g.rejections, res.Rejected...)
	return res
}

// admitNode canonicalizes, checks the denylist and admits. Canonicalization
// runs before the admission decision — not alongside it — so the ledger's
// denylist compares canonical keys and "Example.com" cannot slip past a row
// recorded for "example.com".
func (g *Graph) admitNode(ref NodeRef, base int, by Provenance) (NodeID, Rejection, bool) {
	t, ok := g.reg.Type(ref.Type)
	if !ok {
		return NodeID{}, Rejection{Kind: RejectUnknownType, Type: ref.Type, Key: ref.Key, By: by}, false
	}
	key, err := t.canonical(ref.Key)
	if err != nil {
		return NodeID{}, Rejection{
			Kind: RejectCanonical, Type: t.name, Key: ref.Key, Detail: err.Error(), By: by,
		}, false
	}
	if _, denied := g.ledger[ledgerKey{typ: t.name, key: key}]; denied {
		return NodeID{}, Rejection{Kind: RejectDenied, Type: t.name, Key: key, By: by}, false
	}
	id := newNodeID(t.name, key)
	if _, exists := g.nodes[id]; exists {
		return id, Rejection{}, true
	}
	if g.overBudget(t.name) {
		_ = g.Decline(t.name, key, base, g.Belief(id), ReasonBudget, by)
		return NodeID{}, Rejection{Kind: RejectDenied, Type: t.name, Key: key, Detail: "budget", By: by}, false
	}
	if g.overFrontier() {
		_ = g.Decline(t.name, key, base, g.Belief(id), ReasonFrontier, by)
		return NodeID{}, Rejection{Kind: RejectDenied, Type: t.name, Key: key, Detail: "frontier", By: by}, false
	}
	g.nodes[id] = &Node{ID: id, Type: t, Key: key, Props: newProps(t.sch)}
	g.order = append(g.order, id)
	g.depth[id] = base
	g.provis[id] = true
	g.counts[t.name]++
	g.total++
	g.admitted++
	return id, Rejection{}, true
}

func (g *Graph) admitEdge(ref EdgeRef, admit func(NodeRef) (NodeID, bool), by Provenance, res *Result) (EdgeID, bool) {
	rel, ok := g.reg.Rel(ref.Rel)
	if !ok {
		res.Rejected = append(res.Rejected, Rejection{Kind: RejectUnknownRel, Rel: ref.Rel, By: by})
		return EdgeID{}, false
	}
	from, ok := admit(ref.From)
	if !ok {
		return EdgeID{}, false
	}

	// The seed-closure invariant. Dispatch-side gating is the optimization;
	// this rejection is the invariant, so an operator that emits a variant edge
	// from outside the closure cannot smuggle it in.
	//
	// This runs before the endpoint is admitted, not after. Rejecting the edge
	// while admitting its target would leave the variant in the graph, unrooted
	// and unreachable — the invariant exists to bound combinatorial expansion,
	// and a check that stops the edge but not the node does not bound anything.
	if rel.class == Variant && !g.closure[from] {
		res.Rejected = append(res.Rejected, Rejection{
			Kind: RejectClosure, Rel: rel.name, Type: g.nodes[from].Type.name, Key: g.nodes[from].Key, By: by,
		})
		return EdgeID{}, false
	}

	// The scope restriction, checked in the same place and for the same reason.
	if rel.class == Variant && !g.inScope(from) {
		res.Rejected = append(res.Rejected, Rejection{
			Kind: RejectScope, Rel: rel.name, Type: g.nodes[from].Type.name, Key: g.nodes[from].Key, By: by,
		})
		return EdgeID{}, false
	}

	to, ok := admit(ref.To)
	if !ok {
		return EdgeID{}, false
	}

	// A node is not a variant of itself. Operators emit raw keys and cannot see
	// canonical form, so an algorithm can produce a string that differs from its
	// origin and canonicalizes back onto it — bit-flipping "google.com" flips a
	// case bit to "Google.com", which folds straight back. The operator's own
	// dedupe compares raw strings and cannot catch that; only the applier, which
	// owns canonicalization, can.
	//
	// Left in, the self-edge makes the seed a variant of itself, which is not
	// cosmetic: analyzers treat any node with an inbound VARIANT_OF as a
	// variant, so the target gets scored as a live typosquat of itself and
	// joins every campaign cluster it hosts.
	if rel.class == Variant && from == to {
		res.Rejected = append(res.Rejected, Rejection{
			Kind: RejectSelfVariant, Rel: rel.name,
			Type: g.nodes[from].Type.name, Key: g.nodes[from].Key, By: by,
		})
		return EdgeID{}, false
	}

	id := newEdgeID(from, rel.name, to)
	if _, exists := g.edges[id]; !exists {
		g.edges[id] = &Edge{ID: id, From: from, To: to, Rel: rel, Props: newProps(rel.sch)}
		g.eord = append(g.eord, id)
	}

	// Depth is the shortest observation distance from the seed; structural and
	// variant edges are free. A node admitted bare inherits the subject's depth
	// provisionally — the first in-edge overrides it outright, and later
	// in-edges lower it. Without the override a node listed in the same delta
	// as its edge would keep the subject's depth and never count the hop.
	if d := g.depth[from] + rel.class.DepthCost(); to != g.seed {
		cur, seen := g.depth[to]
		switch {
		case !seen || g.provis[to]:
			g.depth[to] = d
			delete(g.provis, to)
			g.relax(to, improveDepth)
		case d < cur:
			g.depth[to] = d
			g.relax(to, improveDepth)
		}
	}
	if rel.class == Structural && g.closure[from] && !g.closure[to] {
		g.closure[to] = true
		g.relax(to, improveClosure)
	}
	g.considerParent(to, from, rel.name)
	return id, true
}

func (g *Graph) applyProp(ps PropSet, admit func(NodeRef) (NodeID, bool), by Provenance, res *Result) {
	switch {
	case ps.Node != nil:
		id, ok := admit(*ps.Node)
		if !ok {
			return
		}
		n := g.nodes[id]
		f, ok := n.Type.sch.Field(ps.Field)
		if !ok {
			res.Rejected = append(res.Rejected, Rejection{
				Kind: RejectUnknownField, Type: n.Type.name, Key: n.Key, Field: ps.Field, By: by,
			})
			return
		}
		changed, err := n.Props.assert(f, ps.Value, by.Operator)
		if err != nil {
			res.Rejected = append(res.Rejected, Rejection{
				Kind: RejectKindMismatch, Type: n.Type.name, Key: n.Key, Field: ps.Field,
				Detail: err.Error(), By: by,
			})
			return
		}
		won, _ := n.Props.Setter(f)
		g.assertions[id] = append(g.assertions[id], Assertion{
			Field: ps.Field, Value: ps.Value, By: by, Won: won == by.Operator,
		})
		if changed {
			res.Changed = append(res.Changed, id)
		}
	case ps.Edge != nil:
		eid, ok := g.admitEdge(*ps.Edge, admit, by, res)
		if !ok {
			return
		}
		e := g.edges[eid]
		f, ok := e.Rel.sch.Field(ps.Field)
		if !ok {
			res.Rejected = append(res.Rejected, Rejection{
				Kind: RejectUnknownField, Rel: e.Rel.name, Field: ps.Field, By: by,
			})
			return
		}
		if _, err := e.Props.assert(f, ps.Value, by.Operator); err != nil {
			res.Rejected = append(res.Rejected, Rejection{
				Kind: RejectKindMismatch, Rel: e.Rel.name, Field: ps.Field, Detail: err.Error(), By: by,
			})
			return
		}
		won, _ := e.Props.Setter(f)
		g.eassert[eid] = append(g.eassert[eid], Assertion{
			Field: ps.Field, Value: ps.Value, By: by, Won: won == by.Operator,
		})
	default:
		res.Rejected = append(res.Rejected, Rejection{
			Kind: RejectMissingNode, Field: ps.Field, Detail: "PropSet names neither node nor edge", By: by,
		})
	}
}

// SetObservers names the operators whose status attests to a node's existence
// (§9): those that actually looked something up.
//
// The distinction is load-bearing and was missing. Existence rolls up "did any
// operator return ok", and a decomposer returns ok when it successfully *parses*
// a name — which says nothing about whether that name exists. Without this,
// every syntactically valid variant read as live, so "-google.com" and
// "'oogle.com" were reported as live typosquats on the strength of having been
// parsed, with every DNS and whois lookup against them empty.
//
// Membership is "the operator declared a rate-limit resource", i.e. it talks to
// something outside the process. A pure computation cannot attest to existence
// no matter what it returns.
func (g *Graph) SetObservers(ids []string) {
	g.observers = make(map[string]bool, len(ids))
	for _, id := range ids {
		g.observers[id] = true
	}
}

// Observers returns the registered observer set, sorted.
//
// It exists so the set can be persisted with the scan. Existence is computed
// from it, so a graph rebuilt without it answers differently from the run that
// produced it — the same bytes yielding two different verdicts, which is the
// one thing this store exists to rule out.
func (g *Graph) Observers() []string {
	out := make([]string, 0, len(g.observers))
	for id := range g.observers {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

// observes reports whether an operator's status counts toward existence.
//
// With no observer set registered the answer is yes for everything. That is the
// honest fallback rather than a safe-looking "no": a caller that never declared
// which operators observe has given us nothing to discriminate on, and
// answering "no" would make every node read as untried. The scheduler always
// registers the set, so this only affects tests driving the applier directly.
func (g *Graph) observes(op string) bool {
	if len(g.observers) == 0 {
		return true
	}
	return g.observers[op]
}

// SetStatus records the terminal outcome of a (node, operator) pair. Skipped is
// not terminal and does not close the pair.
func (g *Graph) SetStatus(id NodeID, op string, s Status) {
	k := statusKey{node: id, op: op}
	if !s.Terminal() {
		delete(g.status, k)
		return
	}
	g.status[k] = s
}

// Status returns the recorded status of a pair.
func (g *Graph) Status(id NodeID, op string) (Status, bool) {
	s, ok := g.status[statusKey{node: id, op: op}]
	return s, ok
}

// Live reports whether any observation operator returned ok for a node.
func (g *Graph) Live(id NodeID) bool {
	for k, s := range g.status {
		if k.node == id && s == StatusOK {
			return true
		}
	}
	return false
}

// Node returns an admitted node.
func (g *Graph) Node(id NodeID) (*Node, bool) { n, ok := g.nodes[id]; return n, ok }

// Edge returns an admitted edge.
func (g *Graph) Edge(id EdgeID) (*Edge, bool) { e, ok := g.edges[id]; return e, ok }

// Depth returns a node's shortest observation distance from the seed.
func (g *Graph) Depth(id NodeID) int { return g.depth[id] }

// InClosure reports seed-closure membership: the seed plus everything reachable
// from it by structural edges. Only members may root variant generation.
func (g *Graph) InClosure(id NodeID) bool { return g.closure[id] }

// Assertions returns every claim made about a node's fields, winning or not.
func (g *Graph) Assertions(id NodeID) []Assertion { return g.assertions[id] }

// Rejections returns every invariant violation refused so far.
func (g *Graph) Rejections() []Rejection { return g.rejections }

// Ledger returns the truncation ledger in report order.
func (g *Graph) Ledger() []LedgerRow {
	rows := make([]LedgerRow, 0, len(g.ledger))
	for _, r := range g.ledger {
		rows = append(rows, r)
	}
	sortLedger(rows)
	return rows
}

// Nodes returns every admitted node sorted by (type, key) — the canonical
// report order, and what makes two runs byte-comparable.
func (g *Graph) Nodes() []*Node {
	out := make([]*Node, 0, len(g.nodes))
	for _, id := range g.order {
		out = append(out, g.nodes[id])
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Type.name != out[j].Type.name {
			return out[i].Type.name < out[j].Type.name
		}
		return out[i].Key < out[j].Key
	})
	return out
}

// Edges returns every admitted edge sorted by (from, relation, to).
func (g *Graph) Edges() []*Edge {
	out := make([]*Edge, 0, len(g.edges))
	for _, id := range g.eord {
		out = append(out, g.edges[id])
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if c := compareID(a.From[:], b.From[:]); c != 0 {
			return c < 0
		}
		if a.Rel.name != b.Rel.name {
			return a.Rel.name < b.Rel.name
		}
		return compareID(a.To[:], b.To[:]) < 0
	})
	return out
}

func compareID(a, b []byte) int {
	for i := range a {
		switch {
		case a[i] < b[i]:
			return -1
		case a[i] > b[i]:
			return 1
		}
	}
	return 0
}

// relax propagates a monotone node property along out-edges until it stops
// improving, starting from a node whose value has just improved.
//
// It exists because depth and closure are the same shape and were both wrong in
// the same way. Each was computed once, where an edge was admitted, from the
// value its source happened to hold at that moment — and a graph is not built
// in dependency order. A later edge can put a node nearer the seed, or bring a
// node into the seed closure, after its descendants have already been given the
// old answer.
//
// Both were real, and both failed silently, which is why they get a shared
// mechanism rather than two patches:
//
//   - Closure gates variant rooting. A structural edge admitted while its source
//     was outside the closure left the target outside it forever, even once the
//     source joined — so that target and everything under it generated no
//     variants at all. A squat nobody generated is the failure this tool exists
//     to prevent, and nothing reported it.
//   - Depth bounds expansion and is printed. A descendant kept the depth it was
//     first given, so a node could be pruned for exceeding --depth when its real
//     distance was inside it, and the report showed the wrong number.
//
// A new monotone node property belongs here too: give it an improve func and
// relax from wherever it improves. Computing it once at edge-admission time is
// the mistake this retires, and it is not visible in the code that makes it.
//
// improve must return true only on a strict improvement. That is what
// terminates the walk on a cyclic graph — depth is bounded below and closure is
// a latch, so neither can improve forever.
func (g *Graph) relax(start NodeID, improve func(g *Graph, from, to NodeID, rel *Rel) bool) {
	// Adjacency is built per call rather than maintained. A relaxation only runs
	// when a property actually improved, which is rare — a shorter route
	// appearing, or a node joining the closure late — so paying O(edges) then
	// beats keeping an index correct on every admission.
	var adj map[NodeID][]*Edge

	queue := []NodeID{start}
	for i := 0; i < len(queue); i++ {
		if adj == nil {
			adj = make(map[NodeID][]*Edge, len(g.eord))
			for _, eid := range g.eord {
				e := g.edges[eid]
				adj[e.From] = append(adj[e.From], e)
			}
		}
		for _, e := range adj[queue[i]] {
			if improve(g, e.From, e.To, e.Rel) {
				queue = append(queue, e.To)
			}
		}
	}
}

// improveDepth lowers a node's depth when a shorter route to it appears.
func improveDepth(g *Graph, from, to NodeID, rel *Rel) bool {
	if to == g.seed {
		return false
	}
	d := g.depth[from] + rel.class.DepthCost()
	if cur, seen := g.depth[to]; seen && d >= cur {
		return false
	}
	g.depth[to] = d
	delete(g.provis, to)
	return true
}

// improveClosure brings a node into the seed closure when its structural parent
// is in it.
func improveClosure(g *Graph, from, to NodeID, rel *Rel) bool {
	if rel.class != Structural || !g.closure[from] || g.closure[to] {
		return false
	}
	g.closure[to] = true
	return true
}
