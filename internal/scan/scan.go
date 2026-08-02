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

// Package scan assembles a run: schema, operators, plan, expansion, analysis.
//
// It is the one place the operator packages are composed, and it exists because
// composition has an order that matters. Schema registration has to happen
// before any operator is built, extensions have to be appended in a fixed
// sequence because field position is part of every stored content address, and
// the plan has to be compiled before the scheduler is handed anything. Spread
// across the CLI those constraints are invisible and get broken silently — a
// missing extension does not fail, it just makes a collector's output vanish.
package scan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rangertaha/urlinsane/internal/analyze"
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/operators/decompose"
	"github.com/rangertaha/urlinsane/internal/operators/observe"
	"github.com/rangertaha/urlinsane/internal/operators/variant"
	"github.com/rangertaha/urlinsane/internal/report"
)

// Options configures one scan. The zero value is a valid domain-or-whatever
// scan of Target with every shipped operator.
type Options struct {
	// Target is the string the user typed. Its type is detected from the
	// string alone (§12).
	Target string
	// Scope narrows which nameable types get varied. Empty means every
	// nameable node in the seed closure, which includes the seed itself.
	Scope []string

	Limits  graph.Limits
	Variant variant.Options
	Observe observe.Options

	// Algorithms restricts variant operators by id. Empty means all of them.
	Algorithms []string
	// Analyzers overrides the analyzer set. Nil means analyze.All().
	Analyzers []graph.Analyzer
	// Belief overrides the execution model. Nil means the engine's uniform
	// default, under which expansion is unranked (§10.5).
	Belief graph.BeliefModel
}

// Result is everything a run produced. The graph is included because --save,
// diffing and resume all need the graph itself, not a rendering of it.
type Result struct {
	Graph  *graph.Graph
	Plan   *graph.Plan
	Report report.Report

	Seed      graph.NodeID
	SeedType  string
	SeedKey   string
	Rounds    int
	Elapsed   time.Duration
	Interrupt bool
}

// Registry builds the schema for a run.
//
// Extensions are passed in a fixed literal order, never from a map and never
// from the operator set the plan happened to select. A field's position is its
// stable index and part of the content address of every node already stored
// (§1.3), so an order that varies with configuration would make two runs of the
// same target produce different CIDs for the same node.
func Registry() (*graph.Registry, error) {
	r := graph.NewRegistry()
	if err := decompose.Register(r, decompose.Extension{
		Fields:    observe.Fields(),
		RelFields: observe.RelFields(),
	}); err != nil {
		return nil, fmt.Errorf("scan: schema: %w", err)
	}
	return r, nil
}

// Operators returns every operator a run may use, in a deterministic order.
func Operators(o Options) ([]graph.Operator, error) {
	ops := decompose.Operators()

	vs, err := variantOps(o)
	if err != nil {
		return nil, err
	}
	ops = append(ops, vs...)
	return append(ops, observe.New(o.Observe)...), nil
}

func variantOps(o Options) ([]graph.Operator, error) {
	if len(o.Algorithms) == 0 {
		return variant.All(o.Variant), nil
	}
	ops, err := variant.Select(o.Variant, o.Algorithms...)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return ops, nil
}

// Plan compiles a run without executing it. This is what --explain and --plan
// print, and it is the same call Run makes, so what --explain shows is what
// runs rather than a second implementation that can drift from it.
func Plan(o Options) (*graph.Registry, []graph.Operator, *graph.Plan, error) {
	reg, err := Registry()
	if err != nil {
		return nil, nil, nil, err
	}
	typ, key, err := decompose.DetectSeed(o.Target)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scan: %w", err)
	}
	if err := checkScope(reg, o.Scope); err != nil {
		return nil, nil, nil, err
	}
	ops, err := Operators(o)
	if err != nil {
		return nil, nil, nil, err
	}
	p, err := graph.Compile(reg, ops, graph.PlanInput{
		Seed:   graph.SeedSpec{Type: typ, Key: key, Scope: o.Scope},
		Limits: o.Limits,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scan: compile: %w", err)
	}
	return reg, ops, p, nil
}

// checkScope rejects a scope naming a type that cannot be varied, before any
// work is done. §12 promises the positional is *narrowing*; a scope that
// silently matched nothing would instead be a scan that quietly did nothing.
func checkScope(reg *graph.Registry, scope []string) error {
	for _, s := range scope {
		t, ok := reg.Type(s)
		if !ok {
			return fmt.Errorf("scan: unknown scope type %q; see --list types", s)
		}
		if t.Cap() != graph.Nameable {
			return fmt.Errorf(
				"scan: %q is an observed type and cannot be varied; scope must name a nameable type", s)
		}
	}
	return nil
}

// Run expands the graph, analyzes it and builds the report.
//
// Cancelling ctx stops expansion at the end of the current round (§12.4): the
// barrier still runs, so parents, belief and the truncation ledger are
// finalized rather than left half-computed, and analysis then runs over a
// coherent prefix. A ten-minute scan has produced something worth keeping, and
// round-by-round expansion is what makes that prefix coherent rather than an
// arbitrary cross-section.
func Run(ctx context.Context, o Options, ropts report.Options) (*Result, error) {
	start := time.Now()

	reg, ops, plan, err := Plan(o)
	if err != nil {
		return nil, err
	}

	g := graph.New(reg)
	if o.Belief != nil {
		g.SetBeliefModel(o.Belief)
	}
	g.SetBudgets(graph.Budgets{
		Global:  o.Limits.NodeBudget,
		PerType: o.Limits.TypeBudget,
	})
	// From the plan, not from Options: the plan is what --explain printed and
	// what --plan pins, so reading scope from anywhere else would let a pinned
	// plan and the run that used it disagree.
	g.SetScope(plan.Seed.Scope)
	seed, err := g.Seed(plan.Seed.Type, plan.Seed.Key)
	if err != nil {
		return nil, fmt.Errorf("scan: seed: %w", err)
	}

	s := graph.NewScheduler(g, plan.Select(ops), plan.Limits)
	runErr := s.Run(ctx)
	interrupted := errors.Is(runErr, context.Canceled) || ctx.Err() != nil
	if runErr != nil && !interrupted {
		return nil, fmt.Errorf("scan: expand: %w", runErr)
	}

	// Analysis runs over whatever exists, including after an interrupt: there
	// is exactly one analyzer lifetime (§9), and skipping it on the partial
	// path would make Ctrl-C the one way to get a report with no findings in it.
	analyzers := o.Analyzers
	if analyzers == nil {
		analyzers = analyze.All()
	}
	if err := g.RunAnalyzers(context.WithoutCancel(ctx), analyzers); err != nil {
		return nil, fmt.Errorf("scan: analyze: %w", err)
	}

	res := &Result{
		Graph:     g,
		Plan:      plan,
		Seed:      seed,
		SeedType:  plan.Seed.Type,
		SeedKey:   plan.Seed.Key,
		Rounds:    s.Rounds,
		Elapsed:   time.Since(start),
		Interrupt: interrupted,
	}

	ropts.Target = o.Target
	ropts.Scope = o.Scope
	ropts.Plan = plan.Hash
	ropts.Rounds = s.Rounds
	ropts.Elapsed = res.Elapsed
	if interrupted {
		ropts.Partial = true
		ropts.PartialWhy = "interrupted; stopped at a round boundary"
	}
	res.Report = report.Build(g, ropts)
	return res, nil
}
