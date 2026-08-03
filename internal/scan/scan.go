// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

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
	"sort"
	"strings"
	"time"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze/all"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
	decomposeall "github.com/rangertaha/urlinsane/internal/plugins/decompose/all"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	observeall "github.com/rangertaha/urlinsane/internal/plugins/observe/all"
	"github.com/rangertaha/urlinsane/internal/plugins/report"
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	variantall "github.com/rangertaha/urlinsane/internal/plugins/variant/all"
	// Language and keyboard plugins register themselves in init(). Without this
	// import the registry is empty, and the algorithms built over it — vowel
	// swapping, homoglyphs, keyboard adjacency, misspellings — iterate an empty
	// list and silently generate nothing. They still appear in --list
	// algorithms and still run, producing no variants, so the loss is invisible.
	//
	// It belongs here because this is the package that composes a run. The old
	// engine got it from internal/config; nothing in the new path imported it.
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

	// Algorithms restricts variant operators by id. Empty means all of them;
	// a "^id" entry excludes instead of selecting (§12.10).
	Algorithms []string
	// Collectors restricts observation operators the same way. Empty means
	// every operator the Observe options support -- which already excludes any
	// whose dependency is missing, so a named-but-absent operator is an error
	// rather than a silent omission.
	Collectors []string
	// Analyzers overrides the analyzer set. Nil means all.All() plus any
	// registered analyzer plugins.
	Analyzers []graph.Analyzer
	// Settings resolves each plugin's configuration. Nil gives every plugin
	// its declared defaults, which is what an unconfigured run wants.
	Settings plugins.Source
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
	ops := decomposeall.Operators()

	vs, err := variantOps(o)
	if err != nil {
		return nil, err
	}
	ops = append(ops, vs...)

	obs, err := observeall.Select(o.Observe, o.Collectors...)
	if err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	ops = append(ops, obs...)

	// Plugin operators come last, but the order here is not the run order:
	// the scheduler decides that from triggers (§4.1). It is only the order
	// the plan lists them in, and Operators sorts by id so it is stable.
	ext, err := plugins.Operators(o.Settings, o.env())
	if err != nil {
		return nil, err
	}
	return append(ops, ext...), nil
}

// env is what registered plugins are built from: the same services and data the
// shipped operators get. Passing them through rather than letting a plugin
// reach for a package-level resolver is what keeps a plugin testable — and what
// lets an offline test run the whole plan without touching the network.
func (o Options) env() plugins.Env {
	return plugins.Env{Observe: o.Observe, Variant: o.Variant}
}

func variantOps(o Options) ([]graph.Operator, error) {
	// Plugin algorithms join the shipped ones before selection, so --algorithm
	// names them the same way and an unknown id is still an error.
	v := o.Variant
	extra, err := plugins.Algorithms(o.Settings, o.env())
	if err != nil {
		return nil, err
	}
	v.Extra = append(append([]variant.Spec(nil), v.Extra...), extra...)

	ops, err := variantall.Select(v, o.Algorithms...)
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
	if err := checkNamedAlgorithms(o.Algorithms, plan); err != nil {
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
	g.SetFrontier(o.Limits.Frontier)
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
		// The standard set, plus whatever registered. An explicit
		// Options.Analyzers replaces both rather than adding to them: a caller
		// naming its analyzers means those and no others.
		analyzers = all.All()
		ext, err := plugins.Analyzers(o.Settings, o.env())
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		analyzers = append(analyzers, ext...)
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

// checkNamedAlgorithms refuses a run whose named algorithms were all pruned.
//
// Selection and compilation are two different gates and only the first was
// honest about it. `--algorithm tld` on a package seed passes selection, because
// tld exists, and is then pruned by Compile as unreachable — nothing produces a
// domain from a package. The run went ahead with no variant operators at all
// and reported a clean, empty result: doing nothing presented as a finding of
// nothing, which is the exact failure SelectOperators documents itself as
// preventing.
//
// An exclusion is not checked. `^tld` asks for tld to be absent, and pruning it
// is that wish granted rather than denied.
//
// Only a total loss is fatal. Naming five algorithms on a seed that can use
// three is a reasonable thing to do — the run still varies something — so the
// survivors run and the plan records the rest, which --explain prints.
func checkNamedAlgorithms(ids []string, plan *graph.Plan) error {
	named, exclude, err := graph.Selection(ids)
	if err != nil || exclude || len(named) == 0 {
		// A malformed selection has already failed in Operators; returning it
		// again here would report it twice.
		return nil
	}

	pruned := make(map[string]bool, len(plan.Pruned))
	for _, id := range plan.Pruned {
		pruned[id] = true
	}

	var lost []string
	for id := range named {
		if pruned[id] {
			lost = append(lost, id)
		}
	}
	if len(lost) < len(named) {
		return nil // something the user named survived
	}
	sort.Strings(lost)

	return fmt.Errorf(
		"scan: every named algorithm (%s) is unreachable from a %s seed, so the scan would generate nothing; see --list algorithms",
		strings.Join(lost, ", "), plan.Seed.Type)
}
