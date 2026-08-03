// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

import (
	"context"
	"sort"
	"sync"
	"time"
)

// Outcome is what an operator reports alongside its delta. The status is the
// operator's own judgement: NXDOMAIN is StatusEmpty, not StatusFailed, and the
// difference is the signal a squatting scanner exists to collect.
type Outcome struct {
	Status Status
	Err    error
}

// OK, Empty, Failed and Timeout are shorthand for the common outcomes.
func OK() Outcome             { return Outcome{Status: StatusOK} }
func Empty() Outcome          { return Outcome{Status: StatusEmpty} }
func Failed(e error) Outcome  { return Outcome{Status: StatusFailed, Err: e} }
func Timeout(e error) Outcome { return Outcome{Status: StatusTimeout, Err: e} }

// retriable reports whether an outcome should be retried inside its round.
// Only transient conditions qualify: an authoritative "absent" is an answer.
func (o Outcome) retriable() bool {
	return o.Status == StatusFailed || o.Status == StatusTimeout
}

// Limits bound expansion. Zero means unbounded except where noted.
type Limits struct {
	MaxDepth   int            // observation hops from the seed
	MaxRounds  int            // backstop for a type-flow that never converges
	Revisions  int            // per-pair re-runs; default 3
	Attempts   int            // per-pair attempts within a round; default 2
	Workers    int            // concurrent operator calls; default 1
	NodeBudget int            // global admitted-node cap
	TypeBudget map[string]int // per-type admitted-node cap
	Frontier   int            // cap on candidates admitted per round
	// OpTimeout bounds one Exec call. Zero means unbounded. It lives here
	// rather than inside each operator so a slow resource cannot stall a round
	// past the scheduler's knowledge — an operator holding its own private
	// deadline is one the barrier cannot reason about.
	OpTimeout time.Duration
}

// withDefaults fills in the unset bounds.
//
// Negative is treated as unset rather than passed through, at every field, so
// that no caller can reach the scheduler with one. A negative Workers panicked
// outright in make(chan) -- the loud failure of the set. The quiet ones were
// worse: a negative Attempts or MaxRounds made the loops they bound run zero
// times, so the scan expanded nothing, reported no error, and printed an empty
// report as a completed run.
//
// The CLI rejects negatives before they get here, with a message naming the
// flag. This is the backstop for every other caller, and it clamps rather than
// erroring because a Limits value has nowhere to put an error.
func (l Limits) withDefaults() Limits {
	if l.Revisions <= 0 {
		l.Revisions = 3
	}
	if l.Attempts <= 0 {
		l.Attempts = 2
	}
	if l.Workers <= 0 {
		l.Workers = 1
	}
	if l.MaxRounds <= 0 {
		l.MaxRounds = 64
	}
	// Zero is unbounded for these, which is their documented default, so a
	// negative means the same thing rather than a bound of minus four.
	if l.MaxDepth < 0 {
		l.MaxDepth = 0
	}
	if l.NodeBudget < 0 {
		l.NodeBudget = 0
	}
	if l.Frontier < 0 {
		l.Frontier = 0
	}
	if l.OpTimeout < 0 {
		l.OpTimeout = 0
	}
	return l
}

type pairKey struct {
	node NodeID
	op   string
}

type seenKey struct {
	node   NodeID
	op     string
	digest [32]byte
}

// Scheduler drives expansion. It responds to three kinds of event — a delta
// applied, a timer for a retriable failure, and a round barrier — and every
// irreversible decision happens at a barrier.
type Scheduler struct {
	g     *Graph
	ops   []Operator
	lim   Limits
	cache *Cache
	rate  *Limiter

	seen      map[seenKey]bool
	revisions map[pairKey]int
	round     int

	// Stats, for progress and tests.
	Dispatched int
	CacheHits  int
	Retries    int
	Rounds     int
}

// NewScheduler binds a scheduler to a graph and an operator set.
func NewScheduler(g *Graph, ops []Operator, lim Limits) *Scheduler {
	sorted := append([]Operator(nil), ops...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Id() < sorted[j].Id() })

	// Tell the graph which operators actually look something up, so Existence
	// does not read a decomposer's "I parsed this" as "this exists" (§9).
	var observers []string
	for _, o := range sorted {
		if o.Resource() != "" {
			observers = append(observers, o.Id())
		}
	}
	g.SetObservers(observers)

	return &Scheduler{
		g:         g,
		ops:       sorted,
		lim:       lim.withDefaults(),
		cache:     NewCache(),
		rate:      NewLimiter(),
		seen:      map[seenKey]bool{},
		revisions: map[pairKey]int{},
	}
}

// Cache exposes the result cache so a caller can pre-warm or inspect it.
func (s *Scheduler) Cache() *Cache { return s.cache }

// Limiter exposes the rate limiter so resource classes can be configured.
func (s *Scheduler) Limiter() *Limiter { return s.rate }

// Round is the current round number.
func (s *Scheduler) Round() int { return s.round }

// candidate is one eligible (node, operator) pair.
type candidate struct {
	id     NodeID
	op     Operator
	digest [32]byte
}

// Run expands until a round produces no new eligible work, or a limit binds.
// Cancelling ctx stops at the end of the current round, so the barrier still
// runs and parents, belief and the ledger are finalized rather than left half
// computed.
func (s *Scheduler) Run(ctx context.Context) error {
	s.barrier() // barrier 0: seeding
	for {
		if ctx.Err() != nil {
			return nil
		}
		if s.round >= s.lim.MaxRounds {
			s.g.noteStoppedEarly(ReasonRoundCap, Provenance{Operator: "engine", Round: s.round})
			return nil
		}
		s.round++
		s.Rounds = s.round

		work := s.eligible()
		if len(work) == 0 {
			return nil
		}
		results := s.dispatch(ctx, work)
		s.applyAll(work, results)
		s.barrier()
	}
}

// eligible collects every pair whose trigger matches and which has not already
// run against this exact read-set. The order is deterministic — (depth, type,
// key, operator) — which is what makes the whole round reproducible regardless
// of the order operators happen to finish in.
func (s *Scheduler) eligible() []candidate {
	var out []candidate
	for _, id := range s.g.order {
		if s.lim.MaxDepth > 0 && s.g.depth[id] > s.lim.MaxDepth {
			continue
		}
		for _, op := range s.ops {
			t := op.Trigger()
			if !s.g.matches(t, id, true) {
				continue
			}
			// A variant operator may only root on a seed-closure member whose
			// type the run's scope admits. The applier enforces both; skipping
			// here saves generating thousands of candidates only to reject
			// every edge that carries them.
			if declaresVariant(op) && (!s.g.closure[id] || !s.g.inScope(id)) {
				s.g.SetStatus(id, op.Id(), StatusSkipped)
				continue
			}
			d := s.g.readDigest(t, id)
			if s.seen[seenKey{node: id, op: op.Id(), digest: d}] {
				continue
			}
			if s.revisions[pairKey{node: id, op: op.Id()}] >= s.lim.Revisions {
				continue
			}
			out = append(out, candidate{id: id, op: op, digest: d})
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if da, db := s.g.depth[a.id], s.g.depth[b.id]; da != db {
			return da < db
		}
		na, nb := s.g.nodes[a.id], s.g.nodes[b.id]
		if na.Type.name != nb.Type.name {
			return na.Type.name < nb.Type.name
		}
		if na.Key != nb.Key {
			return na.Key < nb.Key
		}
		return a.op.Id() < b.op.Id()
	})
	return out
}

// declaresVariant reports whether an operator emits VARIANT_OF. That
// declaration — not a naming convention — is what makes it a variant operator.
func declaresVariant(op Operator) bool {
	for _, r := range op.Emits().Rels {
		if r == VariantRel {
			return true
		}
	}
	return false
}

// VariantRel is the relation name that marks a generated variation.
const VariantRel = "VARIANT_OF"

type dispatchResult struct {
	delta   Delta
	outcome Outcome
	cached  bool
}

// dispatch runs a round's pairs. Calls run concurrently, but nothing is applied
// here: results are collected and applied afterwards in the deterministic order
// of the work list.
func (s *Scheduler) dispatch(ctx context.Context, work []candidate) []dispatchResult {
	results := make([]dispatchResult, len(work))
	sem := make(chan struct{}, s.lim.Workers)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, c := range work {
		wg.Add(1)
		go func(i int, c candidate) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			key := s.cache.Key(c.op, c.id, c.digest)
			if d, o, ok := s.cache.Get(key); ok {
				mu.Lock()
				s.CacheHits++
				mu.Unlock()
				results[i] = dispatchResult{delta: d, outcome: o, cached: true}
				return
			}

			// Retry lives inside the round: the barrier waits for it, so the
			// outcome set is fixed before belief is computed.
			var d Delta
			var o Outcome
			for attempt := 0; attempt < s.lim.Attempts; attempt++ {
				if ctx.Err() != nil {
					o = Timeout(ctx.Err())
					break
				}
				s.rate.Acquire(ctx, c.op.Resource())
				v := s.g.viewFor(c.op.Trigger(), c.id)
				d, o = s.exec(ctx, c.op, v)
				if !o.retriable() {
					break
				}
				mu.Lock()
				s.Retries++
				mu.Unlock()
			}
			mu.Lock()
			s.Dispatched++
			mu.Unlock()
			if o.Status != StatusFailed && o.Status != StatusTimeout {
				s.cache.Put(key, d, o)
			}
			results[i] = dispatchResult{delta: d, outcome: o}
		}(i, c)
	}
	wg.Wait()
	return results
}

// exec runs one operator under the round's context, bounded by OpTimeout.
//
// The deadline is applied here rather than left to each operator so that every
// operator is bounded by the same rule, including one that forgot to bound
// itself. An operator that ignores ctx still runs to completion — Go cannot
// preempt it — but its result arrives against an expired context and the round
// is not blocked from concluding.
func (s *Scheduler) exec(ctx context.Context, op Operator, v View) (Delta, Outcome) {
	if s.lim.OpTimeout <= 0 {
		return op.Exec(ctx, v)
	}
	ctx, cancel := context.WithTimeout(ctx, s.lim.OpTimeout)
	defer cancel()
	return op.Exec(ctx, v)
}

// applyAll applies a round's results in work order — never completion order —
// and records each pair's terminal status and read-set.
func (s *Scheduler) applyAll(work []candidate, results []dispatchResult) {
	for i, c := range work {
		r := results[i]
		by := Provenance{Operator: c.op.Id(), Round: s.round}
		s.g.Apply(by, c.id, r.delta)
		s.g.SetStatus(c.id, c.op.Id(), r.outcome.Status)

		// This exact read-set is now closed, whatever the outcome. Retry already
		// happened inside the round, bounded by the attempt count — leaving a
		// failed pair eligible as well would multiply attempts by rounds and
		// hammer a service that is already unwell. If the graph later changes
		// something the operator declared it reads, the digest changes and the
		// pair becomes eligible again on its own.
		s.seen[seenKey{node: c.id, op: c.op.Id(), digest: c.digest}] = true
		// Only an execution counts against the revision cap. A pair that was
		// gated off never ran, so there is nothing to revise.
		if r.outcome.Status != StatusSkipped {
			s.revisions[pairKey{node: c.id, op: c.op.Id()}]++
		}
	}
}

// barrier finalizes everything irreversible. Belief is recomputed for every
// node from its props as of now — a node created this round has only the props
// its creating operator set, so computing belief once at creation would leave
// it a bare prior forever.
func (s *Scheduler) barrier() {
	s.g.recomputeBelief()
	s.g.endRound(Provenance{Operator: "engine", Round: s.round})
}
