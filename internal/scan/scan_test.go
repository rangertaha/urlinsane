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

package scan

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/operators/decompose"
	"github.com/rangertaha/urlinsane/internal/operators/observe"
	"github.com/rangertaha/urlinsane/internal/operators/variant"
	"github.com/rangertaha/urlinsane/internal/report"
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

func TestLanguageAndKeyboardPluginsAreRegistered(t *testing.T) {
	// These register via init() in plugins/languages/all. Without that import
	// the registry is empty and every language- or keyboard-driven algorithm
	// iterates an empty list, generating nothing — while still appearing in
	// --list algorithms and still running. The failure is completely silent,
	// which is why it needs a test rather than a comment.
	if n := len(variant.RegisteredLanguages()); n == 0 {
		t.Error("no language plugins registered; language-driven algorithms are no-ops")
	}
	if n := len(variant.RegisteredKeyboards()); n == 0 {
		t.Error("no keyboard plugins registered; keyboard-driven algorithms are no-ops")
	}
}

func TestLanguageDrivenAlgorithmsActuallyGenerate(t *testing.T) {
	// The registry being populated is necessary but not sufficient — the specs
	// snapshot it at construction, so an operator built before registration
	// would still be empty.
	res, err := Run(context.Background(), Options{
		Target: "example.com", Algorithms: []string{"vs"}, // vowel swapping, language-driven
		Limits: graph.Limits{MaxDepth: 1}, Observe: offline(),
	}, report.Options{})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if countVariants(res.Graph) == 0 {
		t.Fatal("vowel swapping produced no variants; its language list is empty")
	}
}
