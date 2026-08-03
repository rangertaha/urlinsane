// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// fakeOp is a synthetic operator. The scheduler is exercised entirely through
// these: no plugins, no network, and every failure mode reachable on demand.
type fakeOp struct {
	id    string
	ver   int
	trig  Trigger
	emits Effects
	res   string
	fn    func(v View) (Delta, Outcome)
	// ctxfn is for operators that must observe cancellation. Most tests do not
	// care and use fn.
	ctxfn func(ctx context.Context, v View) (Delta, Outcome)

	// mu guards the counters below. The scheduler dispatches concurrently, so
	// an operator's own bookkeeping races with itself whenever Workers > 1 —
	// which is the configuration the concurrency tests exist to exercise.
	// Reads in assertions are unguarded on purpose: they run after Run returns,
	// and the scheduler's wait group already establishes happens-before.
	mu     sync.Mutex
	calls  int
	rounds []int // the round each call landed in, filled by the harness
}

func (o *fakeOp) Id() string       { return o.id }
func (o *fakeOp) Version() int     { return o.ver }
func (o *fakeOp) Trigger() Trigger { return o.trig }
func (o *fakeOp) Emits() Effects   { return o.emits }
func (o *fakeOp) Resource() string { return o.res }
func (o *fakeOp) Exec(ctx context.Context, v View) (Delta, Outcome) {
	o.mu.Lock()
	o.calls++
	o.mu.Unlock()
	if o.ctxfn != nil {
		return o.ctxfn(ctx, v)
	}
	return o.fn(v)
}

// onDomain is the common trigger: bind to every domain node.
func onDomain(reads Reads, where ...Condition) Trigger {
	return Trigger{On: Selector{Types: []string{"domain"}}, Where: where, Reads: reads}
}

func run(t *testing.T, g *Graph, ops []Operator, lim Limits) *Scheduler {
	t.Helper()
	s := NewScheduler(g, ops, lim)
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("run: %v", err)
	}
	return s
}

// --- rounds ------------------------------------------------------------------

func TestRoundIsOneDispatchGeneration(t *testing.T) {
	// dns emits an ip; geo binds to ip. If a mid-round delta dispatched
	// immediately, geo would run in the same round as dns and its belief would
	// not yet be computed. Rounds must separate them.
	g, _ := seeded(t)

	var geoRound int
	dns := &fakeOp{
		id: "dns", ver: 1, trig: onDomain(Reads{}),
		emits: Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}},
		fn: func(v View) (Delta, Outcome) {
			return Delta{Edges: []EdgeRef{{
				From: v.Ref(), Rel: "RESOLVES_TO", To: NodeRef{Type: "ip", Key: "1.2.3.4"},
			}}}, OK()
		},
	}
	var s *Scheduler
	geo := &fakeOp{
		id: "geo", ver: 1,
		trig:  Trigger{On: Selector{Types: []string{"ip"}}},
		emits: Effects{Props: []string{"asn"}},
		fn: func(v View) (Delta, Outcome) {
			geoRound = s.Round()
			ref := v.Ref()
			return Delta{Props: []PropSet{{Node: &ref, Field: "asn", Value: String("AS13335")}}}, OK()
		},
	}
	s = NewScheduler(g, []Operator{dns, geo}, Limits{})
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("run: %v", err)
	}

	if dns.calls != 1 || geo.calls != 1 {
		t.Fatalf("calls dns=%d geo=%d, want 1 and 1", dns.calls, geo.calls)
	}
	if geoRound < 2 {
		t.Fatalf("geo ran in round %d; it must not share a round with the dns call that created its node", geoRound)
	}
}

func TestExpansionStopsWhenNoWorkRemains(t *testing.T) {
	g, _ := seeded(t)
	noop := &fakeOp{
		id: "noop", ver: 1, trig: onDomain(Reads{}),
		fn: func(View) (Delta, Outcome) { return Delta{}, Empty() },
	}
	s := run(t, g, []Operator{noop}, Limits{})
	if noop.calls != 1 {
		t.Fatalf("calls = %d, want 1 — nothing changed, so nothing should re-dispatch", noop.calls)
	}
	if s.Rounds > 2 {
		t.Fatalf("ran %d rounds for a single no-op operator", s.Rounds)
	}
}

// --- seen-set and re-dispatch -----------------------------------------------

func TestUndeclaredChangeDoesNotRedispatch(t *testing.T) {
	// The operator declares nothing, so a prop another operator sets must not
	// make it eligible again.
	g, seed := seeded(t)
	watcher := &fakeOp{
		id: "watch", ver: 1, trig: onDomain(Reads{}),
		fn: func(View) (Delta, Outcome) { return Delta{}, OK() },
	}
	setter := &fakeOp{
		id: "set", ver: 1, trig: onDomain(Reads{}),
		emits: Effects{Props: []string{"rank"}},
		fn: func(v View) (Delta, Outcome) {
			ref := v.Ref()
			return Delta{Props: []PropSet{{Node: &ref, Field: "rank", Value: Int(1)}}}, OK()
		},
	}
	run(t, g, []Operator{watcher, setter}, Limits{})
	_ = seed
	if watcher.calls != 1 {
		t.Fatalf("watcher ran %d times; an undeclared field changed and must not re-trigger it", watcher.calls)
	}
}

func TestDeclaredChangeDoesRedispatch(t *testing.T) {
	g, _ := seeded(t)
	watcher := &fakeOp{
		id: "watch", ver: 1, trig: onDomain(Reads{Fields: []string{"rank"}}),
		fn: func(View) (Delta, Outcome) { return Delta{}, OK() },
	}
	setter := &fakeOp{
		id: "set", ver: 1, trig: onDomain(Reads{}),
		emits: Effects{Props: []string{"rank"}},
		fn: func(v View) (Delta, Outcome) {
			ref := v.Ref()
			return Delta{Props: []PropSet{{Node: &ref, Field: "rank", Value: Int(1)}}}, OK()
		},
	}
	run(t, g, []Operator{watcher, setter}, Limits{})
	if watcher.calls < 2 {
		t.Fatalf("watcher ran %d times; it declared rank, which changed, so it must re-run", watcher.calls)
	}
}

func TestRevisionCapBoundsRework(t *testing.T) {
	// An operator that keeps changing what it reads would loop forever without
	// the cap.
	g, _ := seeded(t)
	n := 0
	churn := &fakeOp{
		id: "churn", ver: 1, trig: onDomain(Reads{Fields: []string{"rank"}}),
		emits: Effects{Props: []string{"rank"}},
		fn: func(v View) (Delta, Outcome) {
			n++
			ref := v.Ref()
			return Delta{Props: []PropSet{{Node: &ref, Field: "rank", Value: Int(int64(n))}}}, OK()
		},
	}
	run(t, g, []Operator{churn}, Limits{Revisions: 3})
	if churn.calls > 3 {
		t.Fatalf("operator ran %d times, cap is 3", churn.calls)
	}
}

// --- failure and retry -------------------------------------------------------

func TestRetryHappensInsideItsRound(t *testing.T) {
	// A failure produces no delta, so nothing would re-trigger it. Retry has to
	// be scheduled, and it has to complete before the barrier.
	g, seed := seeded(t)
	flaky := &fakeOp{
		id: "flaky", ver: 1, trig: onDomain(Reads{}),
		fn: func(View) (Delta, Outcome) { return Delta{}, Failed(errors.New("SERVFAIL")) },
	}
	s := run(t, g, []Operator{flaky}, Limits{Attempts: 3})
	if flaky.calls != 3 {
		t.Fatalf("attempts = %d, want 3", flaky.calls)
	}
	if s.Retries == 0 {
		t.Fatal("no retries recorded")
	}
	if st, ok := g.Status(seed, "flaky"); !ok || st != StatusFailed {
		t.Fatalf("status = %v/%v, want failed recorded terminally", st, ok)
	}
}

func TestSuccessAfterRetryStops(t *testing.T) {
	g, _ := seeded(t)
	n := 0
	flaky := &fakeOp{
		id: "flaky", ver: 1, trig: onDomain(Reads{}),
		fn: func(View) (Delta, Outcome) {
			n++
			if n < 2 {
				return Delta{}, Timeout(errors.New("deadline"))
			}
			return Delta{}, OK()
		},
	}
	run(t, g, []Operator{flaky}, Limits{Attempts: 5})
	if flaky.calls != 2 {
		t.Fatalf("attempts = %d, want 2 — retry must stop on success", flaky.calls)
	}
}

func TestAbsenceIsNotRetried(t *testing.T) {
	// NXDOMAIN is an answer. Retrying it would waste calls and, worse, imply
	// the tool could not tell absence from failure.
	g, _ := seeded(t)
	absent := &fakeOp{
		id: "absent", ver: 1, trig: onDomain(Reads{}),
		fn: func(View) (Delta, Outcome) { return Delta{}, Empty() },
	}
	run(t, g, []Operator{absent}, Limits{Attempts: 5})
	if absent.calls != 1 {
		t.Fatalf("attempts = %d, want 1 — an authoritative absence is not transient", absent.calls)
	}
}

// --- gating ------------------------------------------------------------------

func TestVariantOperatorSkippedOutsideClosure(t *testing.T) {
	g, _ := seeded(t)

	// Put a domain outside the closure, reached by observation.
	ptr := &fakeOp{
		id: "ptr", ver: 1, trig: onDomain(Reads{}),
		emits: Effects{Nodes: []string{"domain"}, Rels: []string{"RESOLVES_TO"}},
		fn: func(v View) (Delta, Outcome) {
			if v.Key() != "example.com" {
				return Delta{}, Empty()
			}
			return Delta{Edges: []EdgeRef{{
				From: v.Ref(), Rel: "RESOLVES_TO", To: NodeRef{Type: "domain", Key: "parked.net"},
			}}}, OK()
		},
	}
	var rooted []string
	variant := &fakeOp{
		id: "omission", ver: 1, trig: onDomain(Reads{}),
		emits: Effects{Nodes: []string{"domain"}, Rels: []string{VariantRel}},
		fn: func(v View) (Delta, Outcome) {
			rooted = append(rooted, v.Key())
			return Delta{}, Empty()
		},
	}
	run(t, g, []Operator{ptr, variant}, Limits{})

	for _, k := range rooted {
		if k == "parked.net" {
			t.Fatal("a variant operator was dispatched against a node outside the seed closure")
		}
	}
	if len(rooted) == 0 || rooted[0] != "example.com" {
		t.Fatalf("variant roots = %v, want the seed", rooted)
	}
	id := newNodeID("domain", "parked.net")
	if st, ok := g.Status(id, "omission"); ok && st != StatusSkipped {
		t.Fatalf("status = %v, want skipped or unrecorded", st)
	}
}

func TestBeliefGateBlocksDispatch(t *testing.T) {
	g, _ := seeded(t)
	g.SetBeliefModel(zeroModel{})
	gated := &fakeOp{
		id: "whois", ver: 1,
		trig: onDomain(Reads{}, BeliefAbove(0.5)),
		fn:   func(View) (Delta, Outcome) { return Delta{}, OK() },
	}
	run(t, g, []Operator{gated}, Limits{})
	if gated.calls != 0 {
		t.Fatalf("gated operator ran %d times despite belief below threshold", gated.calls)
	}
}

type zeroModel struct{}

func (zeroModel) Initial() (float64, State)                 { return 0, nil }
func (zeroModel) Step(State, string, View) (float64, State) { return 0, nil }

func TestUniformModelLeavesEverythingEligible(t *testing.T) {
	// The engine must work correctly before any model exists.
	g, _ := seeded(t)
	gated := &fakeOp{
		id: "whois", ver: 1,
		trig: onDomain(Reads{}, BeliefAbove(0.5)),
		fn:   func(View) (Delta, Outcome) { return Delta{}, OK() },
	}
	run(t, g, []Operator{gated}, Limits{})
	if gated.calls != 1 {
		t.Fatalf("calls = %d, want 1 under the uniform default", gated.calls)
	}
}

// --- budgets and caps --------------------------------------------------------

func TestBudgetDeclinesToLedger(t *testing.T) {
	g, _ := seeded(t)
	g.SetBudgets(Budgets{PerType: map[string]int{"domain": 2}})
	gen := &fakeOp{
		id: "gen", ver: 1, trig: onDomain(Reads{}),
		emits: Effects{Nodes: []string{"domain"}, Rels: []string{VariantRel}},
		fn: func(v View) (Delta, Outcome) {
			if v.Key() != "example.com" {
				return Delta{}, Empty()
			}
			var d Delta
			for _, k := range []string{"a.com", "b.com", "c.com"} {
				d.Edges = append(d.Edges, EdgeRef{
					From: v.Ref(), Rel: VariantRel, To: NodeRef{Type: "domain", Key: k},
				})
			}
			return d, OK()
		},
	}
	run(t, g, []Operator{gen}, Limits{})

	rows := g.Ledger()
	if len(rows) == 0 {
		t.Fatal("budget bound but nothing was written to the truncation ledger")
	}
	for _, r := range rows {
		if r.Reason != ReasonBudget {
			t.Fatalf("ledger row %+v, want ReasonBudget", r)
		}
	}
}

func TestRoundCapIsReported(t *testing.T) {
	g, _ := seeded(t)
	n := 0
	churn := &fakeOp{
		id: "churn", ver: 1, trig: onDomain(Reads{Fields: []string{"rank"}}),
		emits: Effects{Props: []string{"rank"}},
		fn: func(v View) (Delta, Outcome) {
			n++
			ref := v.Ref()
			return Delta{Props: []PropSet{{Node: &ref, Field: "rank", Value: Int(int64(n))}}}, OK()
		},
	}
	s := run(t, g, []Operator{churn}, Limits{MaxRounds: 2, Revisions: 100})
	if s.Rounds > 2 {
		t.Fatalf("ran %d rounds, cap was 2", s.Rounds)
	}
	tr := g.Truncations()
	if len(tr) != 1 || tr[0].Reason != ReasonRoundCap {
		t.Fatalf("truncations = %+v, want one round-cap row", tr)
	}
}

func TestMaxDepthStopsExpansion(t *testing.T) {
	g, _ := seeded(t)
	chain := &fakeOp{
		id: "chain", ver: 1,
		trig:  Trigger{On: Selector{Caps: []Capability{Nameable, Observed}}},
		emits: Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}},
		fn: func(v View) (Delta, Outcome) {
			next := fmt.Sprintf("10.0.0.%d", v.Depth()+1)
			return Delta{Edges: []EdgeRef{{
				From: v.Ref(), Rel: "RESOLVES_TO", To: NodeRef{Type: "ip", Key: next},
			}}}, OK()
		},
	}
	run(t, g, []Operator{chain}, Limits{MaxDepth: 2, MaxRounds: 10})
	for _, n := range g.Nodes() {
		if d := g.Depth(n.ID); d > 3 {
			t.Fatalf("node %s at depth %d, past the depth limit", n.Key, d)
		}
	}
}

// --- cache -------------------------------------------------------------------

func TestCacheHitStillOccupiesItsRound(t *testing.T) {
	// A warm run and a cold run must produce the same rounds, the same barriers
	// and the same graph; otherwise plan pinning would depend on cache state.
	build := func(cache *Cache) (*Scheduler, *Graph) {
		g, _ := seeded(t)
		dns := &fakeOp{
			id: "dns", ver: 1, trig: onDomain(Reads{}),
			emits: Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}},
			fn: func(v View) (Delta, Outcome) {
				return Delta{Edges: []EdgeRef{{
					From: v.Ref(), Rel: "RESOLVES_TO", To: NodeRef{Type: "ip", Key: "1.2.3.4"},
				}}}, OK()
			},
		}
		s := NewScheduler(g, []Operator{dns}, Limits{})
		if cache != nil {
			s.cache = cache
		}
		if err := s.Run(context.Background()); err != nil {
			t.Fatalf("run: %v", err)
		}
		return s, g
	}

	cold, gCold := build(nil)
	warm, gWarm := build(cold.cache)

	if warm.CacheHits == 0 {
		t.Fatal("warm run had no cache hits")
	}
	if cold.Rounds != warm.Rounds {
		t.Fatalf("rounds cold=%d warm=%d; a cache hit must not short-circuit the round structure",
			cold.Rounds, warm.Rounds)
	}
	if a, b := len(gCold.Nodes()), len(gWarm.Nodes()); a != b {
		t.Fatalf("nodes cold=%d warm=%d", a, b)
	}
	for i, n := range gCold.Nodes() {
		if gWarm.Nodes()[i].Key != n.Key {
			t.Fatal("warm run produced a different graph")
		}
	}
}

func TestCacheKeyCoversVersionAndReadSet(t *testing.T) {
	g, seed := seeded(t)
	c := NewCache()
	op1 := &fakeOp{id: "x", ver: 1, trig: onDomain(Reads{})}
	op2 := &fakeOp{id: "x", ver: 2, trig: onDomain(Reads{})}
	d := g.readDigest(op1.Trigger(), seed)

	if c.Key(op1, seed, d) == c.Key(op2, seed, d) {
		t.Fatal("bumping Version() did not change the cache key")
	}
	var other [32]byte
	other[0] = 9
	if c.Key(op1, seed, d) == c.Key(op1, seed, other) {
		t.Fatal("a different read-set did not change the cache key")
	}
	c.SetModels("x", []string{"bafyOLD"})
	k1 := c.Key(op1, seed, d)
	c.SetModels("x", []string{"bafyNEW"})
	if k1 == c.Key(op1, seed, d) {
		t.Fatal("a retrained plugin model did not invalidate the cache")
	}
}

// --- determinism -------------------------------------------------------------

func TestRunIsReproducible(t *testing.T) {
	// The property nearly every constraint in the design exists to protect.
	fingerprint := func() string {
		g, _ := seeded(t)
		gen := &fakeOp{
			id: "gen", ver: 1, trig: onDomain(Reads{}),
			emits: Effects{Nodes: []string{"domain"}, Rels: []string{VariantRel}},
			fn: func(v View) (Delta, Outcome) {
				if v.Key() != "example.com" {
					return Delta{}, Empty()
				}
				var d Delta
				for _, k := range []string{"zeta.com", "alpha.com", "mid.com"} {
					d.Edges = append(d.Edges, EdgeRef{
						From: v.Ref(), Rel: VariantRel, To: NodeRef{Type: "domain", Key: k},
					})
				}
				return d, OK()
			},
		}
		dns := &fakeOp{
			id: "dns", ver: 1, trig: onDomain(Reads{}),
			emits: Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}},
			fn: func(v View) (Delta, Outcome) {
				return Delta{Edges: []EdgeRef{{
					From: v.Ref(), Rel: "RESOLVES_TO", To: NodeRef{Type: "ip", Key: "1.2.3.4"},
				}}}, OK()
			},
		}
		run(t, g, []Operator{gen, dns}, Limits{Workers: 4})

		out := ""
		for _, n := range g.Nodes() {
			_, c, err := n.Addressed()
			if err != nil {
				t.Fatalf("addressed: %v", err)
			}
			out += fmt.Sprintf("%s:%s=%s d=%d\n", n.Type.Name(), n.Key, c, g.Depth(n.ID))
		}
		for _, e := range g.Edges() {
			out += fmt.Sprintf("%s->%s\n", e.Rel.Name(), e.To)
		}
		return out
	}

	first := fingerprint()
	for i := 0; i < 4; i++ {
		if got := fingerprint(); got != first {
			t.Fatalf("run %d differed under concurrency:\n%s\nvs\n%s", i, got, first)
		}
	}
}

func TestParentIsDeterministicUnderConvergence(t *testing.T) {
	// Two variants resolve to the same IP. Whichever answers first must not
	// decide the IP's tree parent.
	pick := func() string {
		g, _ := seeded(t)
		gen := &fakeOp{
			id: "gen", ver: 1, trig: onDomain(Reads{}),
			emits: Effects{Nodes: []string{"domain"}, Rels: []string{VariantRel}},
			fn: func(v View) (Delta, Outcome) {
				if v.Key() != "example.com" {
					return Delta{}, Empty()
				}
				return Delta{Edges: []EdgeRef{
					{From: v.Ref(), Rel: VariantRel, To: NodeRef{Type: "domain", Key: "aaa.com"}},
					{From: v.Ref(), Rel: VariantRel, To: NodeRef{Type: "domain", Key: "zzz.com"}},
				}}, OK()
			},
		}
		dns := &fakeOp{
			id: "dns", ver: 1, trig: onDomain(Reads{}),
			emits: Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}},
			fn: func(v View) (Delta, Outcome) {
				if v.Key() == "example.com" {
					return Delta{}, Empty()
				}
				return Delta{Edges: []EdgeRef{{
					From: v.Ref(), Rel: "RESOLVES_TO", To: NodeRef{Type: "ip", Key: "1.2.3.4"},
				}}}, OK()
			},
		}
		run(t, g, []Operator{gen, dns}, Limits{Workers: 8})
		p, rel, ok := g.Parent(newNodeID("ip", "1.2.3.4"))
		if !ok {
			t.Fatal("ip node has no tree parent")
		}
		n, _ := g.Node(p)
		return n.Key + " via " + rel
	}
	first := pick()
	for i := 0; i < 5; i++ {
		if got := pick(); got != first {
			t.Fatalf("tree parent varied between runs: %q vs %q", got, first)
		}
	}
}

// --- view scoping ------------------------------------------------------------

func TestViewHidesUndeclaredPropsAndEdges(t *testing.T) {
	g, seed := seeded(t)
	ref := NodeRef{Type: "domain", Key: "example.com"}
	g.Apply(op("seedprops"), seed, Delta{
		Props: []PropSet{{Node: &ref, Field: "rank", Value: Int(5)}},
		Edges: []EdgeRef{{From: ref, Rel: "TLD_OF", To: NodeRef{Type: "tld", Key: "com"}}},
	})

	var sawProp, sawEdge bool
	probe := &fakeOp{
		id: "probe", ver: 1, trig: onDomain(Reads{}),
		fn: func(v View) (Delta, Outcome) {
			_, sawProp = v.Prop("rank")
			sawEdge = len(v.Edges("TLD_OF")) > 0
			return Delta{}, OK()
		},
	}
	run(t, g, []Operator{probe}, Limits{})
	if sawProp {
		t.Fatal("view exposed a prop the trigger never declared")
	}
	if sawEdge {
		t.Fatal("view exposed a relation the trigger never declared")
	}
}

func TestViewExposesDeclaredEdgeProps(t *testing.T) {
	// VARIANT_OF carries algorithm and distance; ported per-node analyzers read
	// them off the edge, so bare neighbour nodes would not be enough.
	g, seed := seeded(t)
	origin := NodeRef{Type: "domain", Key: "example.com"}
	variant := NodeRef{Type: "domain", Key: "exmple.com"}
	e := EdgeRef{From: origin, Rel: VariantRel, To: variant}
	g.Apply(op("gen"), seed, Delta{
		Edges: []EdgeRef{e},
		Props: []PropSet{
			{Edge: &e, Field: "algorithm", Value: String("omission")},
			{Edge: &e, Field: "distance", Value: Int(1)},
		},
	})

	var algo string
	var dist int64
	probe := &fakeOp{
		id: "probe", ver: 1,
		trig: onDomain(Reads{Rels: []string{VariantRel}}),
		fn: func(v View) (Delta, Outcome) {
			for _, ev := range v.Edges(VariantRel) {
				if a, ok := ev.Prop("algorithm"); ok {
					algo = a.Str()
				}
				if d, ok := ev.Prop("distance"); ok {
					dist = d.Num()
				}
			}
			return Delta{}, OK()
		},
	}
	run(t, g, []Operator{probe}, Limits{})
	if algo != "omission" || dist != 1 {
		t.Fatalf("edge props = %q/%d, want omission/1", algo, dist)
	}
}

// --- cancellation reaches the operator ---------------------------------------

func TestCancellationReachesAnInFlightOperator(t *testing.T) {
	// §12.4: Ctrl-C stops expansion at the end of the current round. That is
	// only achievable if the operator can see the cancellation — otherwise the
	// round cannot end until a hung lookup returns on its own, and "stops at the
	// round boundary" becomes "stops whenever the network feels like it".
	g, _ := seeded(t)
	ctx, cancel := context.WithCancel(context.Background())

	entered := make(chan struct{})
	var observed bool
	op := &fakeOp{
		id: "slow", ver: 1, trig: onDomain(Reads{}),
		ctxfn: func(ctx context.Context, _ View) (Delta, Outcome) {
			close(entered)
			<-ctx.Done() // hangs forever if the context never arrives
			observed = true
			return Delta{}, Timeout(ctx.Err())
		},
	}

	go func() {
		<-entered
		cancel()
	}()

	s := NewScheduler(g, []Operator{op}, Limits{})
	if err := s.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		t.Fatalf("run: %v", err)
	}
	if !observed {
		t.Fatal("the operator never saw the cancellation")
	}
}

func TestOpTimeoutBoundsOneCall(t *testing.T) {
	// An operator that forgot to bound itself is still bounded, and by a rule
	// the scheduler knows about rather than one hidden inside the plugin.
	g, _ := seeded(t)
	var deadline bool
	op := &fakeOp{
		id: "hang", ver: 1, trig: onDomain(Reads{}),
		ctxfn: func(ctx context.Context, _ View) (Delta, Outcome) {
			<-ctx.Done()
			deadline = errors.Is(ctx.Err(), context.DeadlineExceeded)
			return Delta{}, Timeout(ctx.Err())
		},
	}
	s := NewScheduler(g, []Operator{op}, Limits{OpTimeout: 20 * time.Millisecond})
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("run: %v", err)
	}
	if !deadline {
		t.Fatal("OpTimeout did not bound the call")
	}
}

func TestOperatorDeadlineDoesNotCancelTheRun(t *testing.T) {
	// One operator timing out must not take the round with it: the per-call
	// deadline is derived, so cancelling it leaves the parent context intact and
	// the other operators in the same round still run.
	g, _ := seeded(t)
	slow := &fakeOp{
		id: "slow", ver: 1, trig: onDomain(Reads{}),
		ctxfn: func(ctx context.Context, _ View) (Delta, Outcome) {
			<-ctx.Done()
			return Delta{}, Timeout(ctx.Err())
		},
	}
	fast := &fakeOp{
		id: "fast", ver: 1, trig: onDomain(Reads{}),
		fn: func(View) (Delta, Outcome) { return Delta{}, OK() },
	}
	s := NewScheduler(g, []Operator{slow, fast}, Limits{
		OpTimeout: 20 * time.Millisecond, Workers: 2,
	})
	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("run: %v", err)
	}
	if fast.calls == 0 {
		t.Fatal("a peer operator's timeout cancelled the whole round")
	}
}

// --- frontier ----------------------------------------------------------------

// The frontier caps admissions within a round, and the candidates it turns away
// are recorded rather than dropped.
//
// This was accepted as a flag, hashed into the plan, and enforced nowhere:
// --frontier changed the plan hash and nothing else, so a run asked to bound
// its per-round expansion silently did not.
func TestFrontierCapsAdmissionsPerRound(t *testing.T) {
	g, _ := seeded(t)
	g.SetFrontier(2)

	// One operator, one round, five candidates. Two may land.
	fan := &fakeOp{
		id: "fan", ver: 1, trig: onDomain(Reads{}),
		emits: Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}},
		fn: func(v View) (Delta, Outcome) {
			var d Delta
			for i := 1; i <= 5; i++ {
				d.Edges = append(d.Edges, EdgeRef{
					From: v.Ref(), Rel: "RESOLVES_TO",
					To: NodeRef{Type: "ip", Key: fmt.Sprintf("10.0.0.%d", i)},
				})
			}
			return d, OK()
		},
	}
	run(t, g, []Operator{fan}, Limits{MaxRounds: 1, Revisions: 1})

	var ips int
	for _, n := range g.Nodes() {
		if n.Type.Name() == "ip" {
			ips++
		}
	}
	if ips != 2 {
		t.Errorf("admitted %d ip nodes, want 2 — the frontier did not bind", ips)
	}

	var declined int
	for _, row := range g.Ledger() {
		if row.Reason == ReasonFrontier {
			declined++
		}
	}
	if declined != 3 {
		t.Errorf("ledger holds %d frontier rows, want 3 — declined candidates must be reported", declined)
	}
}

// The allowance is per round, not per run: a later round gets a fresh one.
func TestFrontierResetsEachRound(t *testing.T) {
	g, _ := seeded(t)
	g.SetFrontier(1)

	// Each domain emits one ip, and each ip emits one domain, so every round
	// has exactly one candidate and none of them should ever be declined.
	chain := &fakeOp{
		id: "chain", ver: 1,
		trig:  Trigger{On: Selector{Caps: []Capability{Nameable, Observed}}},
		emits: Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}},
		fn: func(v View) (Delta, Outcome) {
			return Delta{Edges: []EdgeRef{{
				From: v.Ref(), Rel: "RESOLVES_TO",
				To: NodeRef{Type: "ip", Key: fmt.Sprintf("10.0.0.%d", v.Depth()+1)},
			}}}, OK()
		},
	}
	run(t, g, []Operator{chain}, Limits{MaxDepth: 3, MaxRounds: 6})

	for _, row := range g.Ledger() {
		if row.Reason == ReasonFrontier {
			t.Fatalf("declined %s at depth %d for the frontier, but each round had one candidate",
				row.Key, row.Depth)
		}
	}
	var ips int
	for _, n := range g.Nodes() {
		if n.Type.Name() == "ip" {
			ips++
		}
	}
	if ips < 2 {
		t.Errorf("admitted %d ip nodes across the run, want at least 2 — the cap did not reset", ips)
	}
}

// Zero means unbounded, so the default run is unaffected.
func TestFrontierZeroIsUnbounded(t *testing.T) {
	g, _ := seeded(t)
	g.SetFrontier(0)

	fan := &fakeOp{
		id: "fan", ver: 1, trig: onDomain(Reads{}),
		emits: Effects{Nodes: []string{"ip"}, Rels: []string{"RESOLVES_TO"}},
		fn: func(v View) (Delta, Outcome) {
			var d Delta
			for i := 1; i <= 5; i++ {
				d.Edges = append(d.Edges, EdgeRef{
					From: v.Ref(), Rel: "RESOLVES_TO",
					To: NodeRef{Type: "ip", Key: fmt.Sprintf("10.0.0.%d", i)},
				})
			}
			return d, OK()
		},
	}
	run(t, g, []Operator{fan}, Limits{MaxRounds: 1, Revisions: 1})

	var ips int
	for _, n := range g.Nodes() {
		if n.Type.Name() == "ip" {
			ips++
		}
	}
	if ips != 5 {
		t.Errorf("admitted %d ip nodes with no frontier, want all 5", ips)
	}
}

// A negative bound reached the scheduler unchanged, because withDefaults only
// looked for zero. Workers panicked in make(chan); the rest were quieter and
// worse -- a negative Attempts or MaxRounds made the loop it bounds run zero
// times, so the scan expanded nothing and reported a successful, empty run.
func TestNegativeLimitsAreTreatedAsUnset(t *testing.T) {
	l := Limits{
		MaxDepth:   -1,
		MaxRounds:  -1,
		Revisions:  -1,
		Attempts:   -1,
		Workers:    -1,
		NodeBudget: -1,
		Frontier:   -1,
		OpTimeout:  -time.Second,
	}.withDefaults()

	for _, tc := range []struct {
		name string
		got  int
		want int
	}{
		{"Revisions", l.Revisions, 3},
		{"Attempts", l.Attempts, 2},
		{"Workers", l.Workers, 1},
		{"MaxRounds", l.MaxRounds, 64},
		// Zero is this field's documented "unbounded", so a negative means
		// that rather than a bound of minus one.
		{"MaxDepth", l.MaxDepth, 0},
		{"NodeBudget", l.NodeBudget, 0},
		{"Frontier", l.Frontier, 0},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %d, want %d", tc.name, tc.got, tc.want)
		}
	}
	if l.OpTimeout != 0 {
		t.Errorf("OpTimeout = %s, want 0", l.OpTimeout)
	}
}

// The panic the clamp exists to stop.
func TestANegativeWorkerCountDoesNotPanic(t *testing.T) {
	l := Limits{Workers: -4}.withDefaults()
	// This is the expression that panicked with "makechan: size out of range".
	sem := make(chan struct{}, l.Workers)
	if cap(sem) < 1 {
		t.Fatalf("semaphore capacity = %d, want at least 1", cap(sem))
	}
}

// The read digest decides two things: whether a pair is dispatched again, and
// what its cached result is keyed on. It has to cover everything the operator
// can see, and View.Edges hands over the props of every edge on a declared
// relation -- not just the edges' identities.
//
// It hashed only edge ids. An id is a content address of the endpoints and the
// relation, not of the props hung off the edge, so an operator reading an edge
// prop saw the digest stay identical when that prop changed: the pair stayed
// marked seen, it was never dispatched again, and the cache served a result
// computed from the old value.
func TestReadDigestCoversEdgeProps(t *testing.T) {
	g, seed := seeded(t)
	trig := onDomain(Reads{Rels: []string{"VARIANT_OF"}})

	edge := EdgeRef{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "exarnple.com"},
	}
	g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{edge}})
	before := g.readDigest(trig, seed)

	// A second operator fills in a field the first left unset. The edge is the
	// same edge -- same endpoints, same relation, same id -- so nothing about
	// its identity moves.
	g.Apply(op("distance"), seed, Delta{
		Props: []PropSet{{Edge: &edge, Field: "distance", Value: Int(3)}},
	})

	if after := g.readDigest(trig, seed); before == after {
		t.Fatal("setting an edge prop left the read digest identical; the operator would never re-run")
	}
}

// An operator that declares no relation reads must not be disturbed by one.
func TestReadDigestIgnoresEdgePropsItDidNotDeclare(t *testing.T) {
	g, seed := seeded(t)
	trig := onDomain(Reads{}) // declares nothing

	edge := EdgeRef{
		From: NodeRef{Type: "domain", Key: "example.com"},
		Rel:  "VARIANT_OF",
		To:   NodeRef{Type: "domain", Key: "exarnple.com"},
	}
	g.Apply(op("omission"), seed, Delta{Edges: []EdgeRef{edge}})
	before := g.readDigest(trig, seed)

	g.Apply(op("distance"), seed, Delta{
		Props: []PropSet{{Edge: &edge, Field: "distance", Value: Int(3)}},
	})
	if after := g.readDigest(trig, seed); before != after {
		t.Fatal("an undeclared edge prop changed the digest; it is invisible to this operator")
	}
}
