// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package scan

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/report"
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// --- the schema contract -----------------------------------------------------

func TestRegistryDeclaresEveryFieldTheOperatorsAssert(t *testing.T) {
	// The bug this exists for is silent: a field an operator asserts but nobody
	// registered is refused by the applier as unknown, recorded as a rejection,
	// and the run reports success with that collector's entire output missing.
	// Nothing fails, nothing logs, the report is just quietly wrong.
	r, err := Registry()
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	for typeName, fields := range observe.Fields() {
		nt, ok := r.Type(typeName)
		if !ok {
			t.Fatalf("observe asserts fields on unregistered type %q", typeName)
		}
		for _, f := range fields {
			if _, ok := nt.Field(f.Name); !ok {
				t.Errorf("%s.%s is asserted by observe but not registered", typeName, f.Name)
			}
		}
	}
	for relName, fields := range observe.RelFields() {
		rel, ok := r.Rel(relName)
		if !ok {
			t.Fatalf("observe asserts props on unregistered relation %q", relName)
		}
		for _, f := range fields {
			if _, ok := rel.Field(f.Name); !ok {
				t.Errorf("%s.%s is asserted by observe but not registered", relName, f.Name)
			}
		}
	}
}

func TestRegistryIsDeterministic(t *testing.T) {
	// Field position is part of every stored content address, so a registration
	// order that varies between runs silently invalidates every stored CID.
	first, err := Registry()
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	fingerprint := func(r *graph.Registry) string {
		var b strings.Builder
		for _, tp := range r.Types() {
			b.WriteString(tp.Name() + ":")
			for _, f := range observe.Fields()[tp.Name()] {
				b.WriteString(f.Name + ",")
			}
			b.WriteString(";")
		}
		return b.String()
	}
	want := fingerprint(first)
	for i := 0; i < 5; i++ {
		got, err := Registry()
		if err != nil {
			t.Fatalf("registry: %v", err)
		}
		if fingerprint(got) != want {
			t.Fatal("schema registration varied between runs")
		}
	}
}

func TestExtensionForAnUnknownTypeIsRejected(t *testing.T) {
	// Dropping it silently would reproduce the very failure Extension exists to
	// prevent, one level up.
	r := graph.NewRegistry()
	err := decompose.Register(r, decompose.Extension{
		Fields: map[string][]graph.FieldDef{
			"nosuchtype": {{Name: "x", Kind: graph.KindString}},
		},
	})
	if err == nil {
		t.Fatal("an extension aimed at an unregistered type was accepted")
	}
	if !strings.Contains(err.Error(), "nosuchtype") {
		t.Fatalf("error does not name the offending type: %v", err)
	}
}

// --- planning ----------------------------------------------------------------

func TestPlanDetectsSeedTypeFromTheTargetAlone(t *testing.T) {
	for target, want := range map[string]string{
		"example.com":            "domain",
		"bob@example.com":        "email",
		"npm:lodash":             "package",
		"github.com/acme/tool":   "repo",
		"https://github.com/a/b": "repo",
	} {
		_, _, p, err := Plan(Options{Target: target, Observe: offline()})
		if err != nil {
			t.Fatalf("plan %q: %v", target, err)
		}
		if p.Seed.Type != want {
			t.Errorf("Plan(%q) seeded %s, want %s", target, p.Seed.Type, want)
		}
	}
}

func TestScopeDoesNotChangeSeedDetection(t *testing.T) {
	// §12: the target parses identically either way; scope narrows what gets
	// varied, never what the string is.
	bare, _, a, err := Plan(Options{Target: "bob@example.com", Observe: offline()})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}
	_ = bare
	_, _, b, err := Plan(Options{Target: "bob@example.com", Scope: []string{"username"}, Observe: offline()})
	if err != nil {
		t.Fatalf("plan: %v", err)
	}
	if a.Seed.Type != b.Seed.Type || a.Seed.Key != b.Seed.Key {
		t.Fatalf("scope changed the seed: %s:%s vs %s:%s",
			a.Seed.Type, a.Seed.Key, b.Seed.Type, b.Seed.Key)
	}
}

func TestBadScopeIsRejectedBeforeAnyWork(t *testing.T) {
	// A scope that silently matched nothing would be a scan that quietly did
	// nothing — worse than one that refuses to start.
	_, _, _, err := Plan(Options{Target: "example.com", Scope: []string{"ip"}, Observe: offline()})
	if err == nil {
		t.Fatal("an observed type was accepted as a variant scope")
	}
	if !strings.Contains(err.Error(), "nameable") {
		t.Fatalf("error does not explain the rule: %v", err)
	}
	if _, _, _, err = Plan(Options{Target: "example.com", Scope: []string{"nope"}, Observe: offline()}); err == nil {
		t.Fatal("an unregistered scope type was accepted")
	}
}

func TestPlanIsWhatRuns(t *testing.T) {
	// --explain must not be a second implementation that can drift from the
	// real one.
	opts := Options{
		Target: "example.com", Algorithms: []string{"co"},
		Limits: graph.Limits{MaxDepth: 1}, Observe: offline(),
	}
	_, _, p, err := Plan(opts)
	if err != nil {
		t.Fatalf("plan: %v", err)
	}
	res, err := Run(context.Background(), opts, report.Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Plan.Hash != p.Hash {
		t.Fatalf("Run compiled a different plan than Plan: %s vs %s", res.Plan.Hash, p.Hash)
	}
}

// --- running -----------------------------------------------------------------

func TestRunProducesVariantsAndAReport(t *testing.T) {
	res, err := Run(context.Background(), Options{
		Target:     "example.com",
		Algorithms: []string{"co"}, // character omission: small and deterministic
		Limits:     graph.Limits{MaxDepth: 1},
		Observe:    offline(),
	}, report.Options{Format: "json"})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.SeedType != "domain" || res.SeedKey != "example.com" {
		t.Fatalf("seed = %s:%s", res.SeedType, res.SeedKey)
	}
	var variants int
	for _, e := range res.Graph.Edges() {
		if e.Rel.Name() == graph.VariantRel {
			variants++
		}
	}
	if variants == 0 {
		t.Fatal("a full run of an omission algorithm produced no variants")
	}
	if len(res.Report.Nodes) < variants {
		t.Fatalf("report has %d nodes for %d variants", len(res.Report.Nodes), variants)
	}
}

func TestRunDecomposesACompositeSeed(t *testing.T) {
	// The reason decomposition is an operator rather than a special case: an
	// email seed must yield its local part and its domain, both in the closure
	// and both varyable.
	res, err := Run(context.Background(), Options{
		Target: "bob@example.com", Algorithms: []string{"co"},
		Limits: graph.Limits{MaxDepth: 1}, Observe: offline(),
	}, report.Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	found := map[string]bool{}
	for _, n := range res.Graph.Nodes() {
		found[n.Type.Name()+":"+n.Key] = true
	}
	for _, want := range []string{"email:bob@example.com", "domain:example.com", "username:bob"} {
		if !found[want] {
			t.Errorf("composite seed did not yield %s", want)
		}
	}
}

func TestScopeNarrowsWhatIsVaried(t *testing.T) {
	broad, err := Run(context.Background(), Options{
		Target: "bob@example.com", Algorithms: []string{"co"},
		Limits: graph.Limits{MaxDepth: 1}, Observe: offline(),
	}, report.Options{})
	if err != nil {
		t.Fatalf("broad: %v", err)
	}
	narrow, err := Run(context.Background(), Options{
		Target: "bob@example.com", Scope: []string{"username"}, Algorithms: []string{"co"},
		Limits: graph.Limits{MaxDepth: 1}, Observe: offline(),
	}, report.Options{})
	if err != nil {
		t.Fatalf("narrow: %v", err)
	}
	if countVariants(narrow.Graph) >= countVariants(broad.Graph) {
		t.Fatalf("scope did not narrow: %d variants scoped vs %d unscoped",
			countVariants(narrow.Graph), countVariants(broad.Graph))
	}
	// And what it did vary must be the scoped type.
	for _, e := range narrow.Graph.Edges() {
		if e.Rel.Name() != graph.VariantRel {
			continue
		}
		from, _ := narrow.Graph.Node(e.From)
		if from.Type.Name() != "username" {
			t.Fatalf("scope username still varied a %s", from.Type.Name())
		}
	}
}

func TestRunIsReproducible(t *testing.T) {
	// Two runs of one target must produce identical reports, or "what changed
	// since last week" is noise.
	fingerprint := func() string {
		res, err := Run(context.Background(), Options{
			Target: "example.com", Algorithms: []string{"co", "cr"},
			Limits: graph.Limits{MaxDepth: 1}, Observe: offline(),
		}, report.Options{})
		if err != nil {
			t.Fatalf("run: %v", err)
		}
		var b strings.Builder
		for _, n := range res.Report.Nodes {
			b.WriteString(n.Type + ":" + n.Key + "\n")
		}
		return b.String()
	}
	first := fingerprint()
	for i := 0; i < 3; i++ {
		if got := fingerprint(); got != first {
			t.Fatal("two runs of the same target produced different graphs")
		}
	}
}

func TestInterruptYieldsAPartialReportNotAnError(t *testing.T) {
	// §12.4: a ten-minute scan has produced something worth keeping. Cancelling
	// must produce a report marked partial, not a failure with nothing in it.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := Run(ctx, Options{
		Target: "example.com", Algorithms: []string{"co"}, Observe: offline(),
	}, report.Options{})
	if err != nil {
		t.Fatalf("an interrupted run returned an error instead of a partial report: %v", err)
	}
	if !res.Interrupt {
		t.Fatal("interrupted run not flagged")
	}
	if !res.Report.Partial {
		t.Fatal("report from an interrupted run is not marked partial")
	}
	if res.Report.PartialWhy == "" {
		t.Fatal("partial report does not say why")
	}
}

func TestAnalysisRunsEvenAfterAnInterrupt(t *testing.T) {
	// There is exactly one analyzer lifetime (§9). Skipping analysis on the
	// partial path would make Ctrl-C the one way to get an empty findings list.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	res, err := Run(ctx, Options{Target: "example.com", Observe: offline()}, report.Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	// The seed alone yields no findings, but the analyzers must have run: the
	// report is built and its totals are populated.
	if res.Report.Totals.Nodes == 0 {
		t.Fatal("no report was built for an interrupted run")
	}
}

func TestUnknownAlgorithmIsAnError(t *testing.T) {
	_, err := Run(context.Background(), Options{
		Target: "example.com", Algorithms: []string{"nosuchalgo"}, Observe: offline(),
	}, report.Options{})
	if err == nil {
		t.Fatal("an unknown algorithm id was accepted")
	}
}

func TestBudgetsReachTheGraph(t *testing.T) {
	res, err := Run(context.Background(), Options{
		Target: "example.com", Limits: graph.Limits{NodeBudget: 5, MaxDepth: 1},
		Observe: offline(),
	}, report.Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(res.Graph.Nodes()) > 5 {
		t.Fatalf("node budget of 5 admitted %d nodes", len(res.Graph.Nodes()))
	}
	// And truncation must be visible, not silent.
	if len(res.Graph.Ledger()) == 0 {
		t.Fatal("a budget bound expansion but wrote no ledger rows")
	}
	if !res.Report.Partial {
		t.Fatal("a truncated scan was not marked partial")
	}
}

// --- helpers -----------------------------------------------------------------

func countVariants(g *graph.Graph) int {
	var n int
	for _, e := range g.Edges() {
		if e.Rel.Name() == graph.VariantRel {
			n++
		}
	}
	return n
}

// offline supplies stubs for every external dependency, so the whole scan
// pipeline is exercised without a network.
func offline() observe.Options {
	return observe.Options{
		Timeout:  time.Millisecond,
		Resolver: deadResolver{},
		Whois:    deadWhois{},
		Prober:   deadProber{},
		Sources:  emptySources{},
	}
}

type deadResolver struct{}

func (deadResolver) LookupIP(context.Context, string, string) ([]net.IP, error) {
	return nil, &net.DNSError{Err: "no such host", IsNotFound: true}
}
func (deadResolver) LookupNS(context.Context, string) ([]*net.NS, error) {
	return nil, &net.DNSError{Err: "no such host", IsNotFound: true}
}
func (deadResolver) LookupMX(context.Context, string) ([]*net.MX, error) {
	return nil, &net.DNSError{Err: "no such host", IsNotFound: true}
}
func (deadResolver) LookupTXT(context.Context, string) ([]string, error) {
	return nil, &net.DNSError{Err: "no such host", IsNotFound: true}
}
func (deadResolver) LookupCNAME(context.Context, string) (string, error) {
	return "", &net.DNSError{Err: "no such host", IsNotFound: true}
}
func (deadResolver) LookupAddr(context.Context, string) ([]string, error) {
	return nil, &net.DNSError{Err: "no such host", IsNotFound: true}
}

type deadWhois struct{}

func (deadWhois) Whois(string, ...string) (string, error) { return "", nil }

type deadProber struct{}

func (deadProber) Exists(context.Context, string) (bool, error) { return false, nil }

type emptySources struct{}

func (emptySources) Sources(string) ([]observe.Source, error) { return nil, nil }

// Keyboard layouts ship with the binary (pkg/kb), so they are available with
// no dataset and no registration. A keyboard-driven algorithm that generated
// nothing would be a silent no-op — it still appears in --list algorithms and
// still runs — which is why this is a test and not a comment.
func TestKeyboardLayoutsAreAvailable(t *testing.T) {
	if n := len(variant.RegisteredKeyboards()); n == 0 {
		t.Fatal("no keyboard layouts; keyboard-driven algorithms are no-ops")
	}

	res, err := Run(context.Background(), Options{
		Target: "google.com", Algorithms: []string{"acs"}, // adjacent character substitution
		Limits: graph.Limits{MaxDepth: 1}, Observe: offline(),
	}, report.Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if countVariants(res.Graph) == 0 {
		t.Fatal("adjacent character substitution produced no variants")
	}
}

// Languages come from the dataset, not from registered plugins, so an empty
// language list is the *expected* state of a run with no dataset — an offline
// unit test, or a fresh install before `datasets import`.
//
// That makes it exactly the failure §12.6 exists to surface: the algorithms
// still appear in the plan and still run, and produce nothing. The engine must
// not pretend otherwise, so this pins the two halves — no languages means no
// language-driven variants, and it is observable rather than silent.
func TestLanguageDrivenAlgorithmsAreEmptyWithoutADataset(t *testing.T) {
	if n := len(variant.RegisteredLanguages()); n != 0 {
		t.Skipf("a dataset is loaded (%d languages); this test covers the empty case", n)
	}

	res, err := Run(context.Background(), Options{
		Target: "example.com", Algorithms: []string{"vs"}, // vowel swapping, language-driven
		Limits: graph.Limits{MaxDepth: 1}, Observe: offline(),
	}, report.Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if n := countVariants(res.Graph); n != 0 {
		t.Fatalf("vowel swapping produced %d variants with no language data", n)
	}
}

// Given a language, the same algorithm generates. This is the other half of the
// pair above: the emptiness is the dataset's absence, not a broken operator.
func TestLanguageDrivenAlgorithmsGenerateWithALanguage(t *testing.T) {
	res, err := Run(context.Background(), Options{
		Target: "example.com", Algorithms: []string{"vs"},
		Variant: variant.Options{Languages: []variant.Language{testLanguage{}}},
		Limits:  graph.Limits{MaxDepth: 1}, Observe: offline(),
	}, report.Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if countVariants(res.Graph) == 0 {
		t.Fatal("vowel swapping produced no variants from a language that has vowels")
	}
}

// testLanguage is the minimum a language-driven algorithm needs.
type testLanguage struct{}

func (testLanguage) Code() string                    { return "xx" }
func (testLanguage) Name() string                    { return "Test" }
func (testLanguage) Vowels() []string                { return []string{"a", "e", "i", "o", "u"} }
func (testLanguage) Graphemes() []string             { return nil }
func (testLanguage) Numerals() map[string][]string   { return nil }
func (testLanguage) Homoglyphs() map[string][]string { return nil }
func (testLanguage) Homophones() [][]string          { return nil }
func (testLanguage) Misspellings() [][]string        { return nil }
