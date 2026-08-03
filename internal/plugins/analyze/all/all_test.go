// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze/campaign"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze/depconfusion"
	"github.com/rangertaha/urlinsane/internal/plugins/analyze/scoring"
)

func registry(t *testing.T) *graph.Registry {
	t.Helper()
	r := graph.NewRegistry()
	lower := func(s string) (string, error) {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			return "", fmt.Errorf("empty key")
		}
		return s, nil
	}
	for _, d := range []graph.NodeTypeDef{
		{Name: "domain", Cap: graph.Nameable, Version: 1, Canonical: lower},
		{Name: "package", Cap: graph.Nameable, Version: 1, Canonical: lower},
		{Name: "ip", Cap: graph.Observed, Version: 1, Canonical: lower},
	} {
		if _, err := r.AddType(d); err != nil {
			t.Fatalf("%s: %v", d.Name, err)
		}
	}
	for _, d := range []graph.RelDef{
		{Name: graph.VariantRel, Class: graph.Variant, Version: 1, Fields: []graph.FieldDef{
			{Name: "algorithm", Kind: graph.KindString},
			{Name: "distance", Kind: graph.KindInt},
		}},
		{Name: "RESOLVES_TO", Class: graph.Observation, Version: 1},
		{Name: "MX", Class: graph.Observation, Version: 1},
		{Name: "EXISTS_ON", Class: graph.Observation, Version: 1},
	} {
		if _, err := r.AddRel(d); err != nil {
			t.Fatalf("%s: %v", d.Name, err)
		}
	}
	return r
}

func by(op string) graph.Provenance { return graph.Provenance{Operator: op, Round: 1} }

// variantsOnOneIP builds the shape campaign clustering exists to find.
func variantsOnOneIP(t *testing.T, keys ...string) (*graph.Graph, graph.NodeID) {
	t.Helper()
	g := graph.New(registry(t))
	seed, err := g.Seed("domain", "example.com")
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	origin := graph.NodeRef{Type: "domain", Key: "example.com"}
	for _, k := range keys {
		v := graph.NodeRef{Type: "domain", Key: k}
		e := graph.EdgeRef{From: origin, Rel: graph.VariantRel, To: v}
		g.Apply(by("omission"), seed, graph.Delta{
			Edges: []graph.EdgeRef{e},
			Props: []graph.PropSet{{Edge: &e, Field: "distance", Value: graph.Int(1)}},
		})
		id := mustNode(t, g, "domain", k)
		g.Apply(by("dns"), id, graph.Delta{Edges: []graph.EdgeRef{{
			From: v, Rel: "RESOLVES_TO", To: graph.NodeRef{Type: "ip", Key: "1.2.3.4"},
		}}})
		g.SetStatus(id, "dns", graph.StatusOK)
	}
	return g, seed
}

func mustNode(t *testing.T, g *graph.Graph, typ, key string) graph.NodeID {
	t.Helper()
	for _, n := range g.Nodes() {
		if n.Type.Name() == typ && n.Key == key {
			return n.ID
		}
	}
	t.Fatalf("node %s:%s not found", typ, key)
	return graph.NodeID{}
}

func TestCampaignClustersSharedInfrastructure(t *testing.T) {
	g, _ := variantsOnOneIP(t, "examp1e.com", "exarnple.com", "exampie.com")
	if err := g.RunAnalyzers(context.Background(), []graph.Analyzer{campaign.Campaign{}}); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	var found *graph.Finding
	for _, f := range g.Findings() {
		if f.Kind == "campaign" {
			ff := f
			found = &ff
		}
	}
	if found == nil {
		t.Fatal("three variants on one IP produced no campaign finding")
	}
	if found.Severity != graph.SeverityHigh {
		t.Fatalf("severity = %v, want high", found.Severity)
	}
}

func TestCampaignIgnoresSingletons(t *testing.T) {
	g, _ := variantsOnOneIP(t, "examp1e.com")
	_ = g.RunAnalyzers(context.Background(), []graph.Analyzer{campaign.Campaign{}})
	for _, f := range g.Findings() {
		if f.Kind == "campaign" {
			t.Fatal("one variant on an address is not a campaign")
		}
	}
}

func TestScoringRatesLiveVariants(t *testing.T) {
	g, _ := variantsOnOneIP(t, "examp1e.com")
	id := mustNode(t, g, "domain", "examp1e.com")
	v := graph.NodeRef{Type: "domain", Key: "examp1e.com"}
	g.Apply(by("dns"), id, graph.Delta{Edges: []graph.EdgeRef{{
		From: v, Rel: "MX", To: graph.NodeRef{Type: "domain", Key: "mail.examp1e.com"},
	}}})

	if err := g.RunAnalyzers(context.Background(), []graph.Analyzer{scoring.Scoring{}}); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if got := g.Risk(id); got < graph.SeverityHigh {
		t.Fatalf("risk = %v; a live one-edit variant accepting mail should rate high", got)
	}
}

func TestScoringSkipsNonLive(t *testing.T) {
	g := graph.New(registry(t))
	seed, _ := g.Seed("domain", "example.com")
	origin := graph.NodeRef{Type: "domain", Key: "example.com"}
	v := graph.NodeRef{Type: "domain", Key: "gone.com"}
	g.Apply(by("omission"), seed, graph.Delta{Edges: []graph.EdgeRef{{From: origin, Rel: graph.VariantRel, To: v}}})
	id := mustNode(t, g, "domain", "gone.com")
	g.SetStatus(id, "dns", graph.StatusEmpty)

	_ = g.RunAnalyzers(context.Background(), []graph.Analyzer{scoring.Scoring{}})
	if len(g.Findings()) != 0 {
		t.Fatalf("findings = %+v, want none for a confirmed-absent variant", g.Findings())
	}
}

func TestDepConfusionDistinguishesAbsentFromUnknown(t *testing.T) {
	// The whole point of a third status value: an unpublished package is a
	// critical gap, but one we could not check is not a finding about the
	// package — it is a finding about the scan.
	g := graph.New(registry(t))
	seed, _ := g.Seed("package", "npm:internal-lib")
	g.SetStatus(seed, "npm", graph.StatusEmpty)

	unknown := graph.NodeRef{Type: "package", Key: "npm:maybe-lib"}
	g.Apply(by("manifest"), seed, graph.Delta{Nodes: []graph.NodeRef{unknown}})
	uid := mustNode(t, g, "package", "npm:maybe-lib")
	g.SetStatus(uid, "npm", graph.StatusTimeout)

	if err := g.RunAnalyzers(context.Background(), []graph.Analyzer{depconfusion.DepConfusion{}}); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	var critical, info int
	for _, f := range g.Findings() {
		switch f.Kind {
		case "dep-confusion":
			critical++
			if f.Severity != graph.SeverityCritical {
				t.Fatalf("absent package severity = %v, want critical", f.Severity)
			}
		case "dep-confusion-unknown":
			info++
		}
	}
	if critical != 1 || info != 1 {
		t.Fatalf("findings critical=%d unknown=%d, want 1 and 1", critical, info)
	}
}

func TestDepConfusionSkipsPublishedPackages(t *testing.T) {
	g := graph.New(registry(t))
	seed, _ := g.Seed("package", "npm:lodash")
	g.Apply(by("npm"), seed, graph.Delta{Edges: []graph.EdgeRef{{
		From: graph.NodeRef{Type: "package", Key: "npm:lodash"},
		Rel:  "EXISTS_ON",
		To:   graph.NodeRef{Type: "package", Key: "npm:lodash-registry"},
	}}})
	g.SetStatus(seed, "npm", graph.StatusOK)

	_ = g.RunAnalyzers(context.Background(), []graph.Analyzer{depconfusion.DepConfusion{}})
	for _, f := range g.Findings() {
		if strings.HasPrefix(f.Kind, "dep-confusion") && len(f.Nodes) > 0 && f.Nodes[0] == seed {
			t.Fatal("a published package was reported as a confusion gap")
		}
	}
}

func TestFindingsAreOrderedMostSevereFirst(t *testing.T) {
	g, _ := variantsOnOneIP(t, "examp1e.com", "exarnple.com")
	if err := g.RunAnalyzers(context.Background(), All()); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	f := g.Findings()
	for i := 1; i < len(f); i++ {
		if f[i-1].Severity < f[i].Severity {
			t.Fatalf("findings not ordered by severity: %v before %v", f[i-1].Severity, f[i].Severity)
		}
	}
}

func TestAnalysisIsDeterministic(t *testing.T) {
	fingerprint := func() string {
		g, _ := variantsOnOneIP(t, "b.com", "a.com", "c.com")
		if err := g.RunAnalyzers(context.Background(), All()); err != nil {
			t.Fatalf("analyze: %v", err)
		}
		out := ""
		for _, f := range g.Findings() {
			out += fmt.Sprintf("%s|%s|%s\n", f.Severity, f.Kind, f.Summary)
		}
		return out
	}
	first := fingerprint()
	for i := 0; i < 5; i++ {
		if got := fingerprint(); got != first {
			t.Fatalf("analysis varied between runs:\n%s\nvs\n%s", got, first)
		}
	}
}

func TestExistenceRollup(t *testing.T) {
	g := graph.New(registry(t))
	seed, _ := g.Seed("domain", "example.com")
	a := g.Analyze()

	if got := a.Existence(seed); got.String() != "untried" {
		t.Fatalf("untried node = %v", got)
	}
	g.SetStatus(seed, "dns", graph.StatusTimeout)
	if got := a.Existence(seed); got != graph.Unknown {
		t.Fatalf("timeout-only = %v, want unknown", got)
	}
	g.SetStatus(seed, "whois", graph.StatusEmpty)
	if got := a.Existence(seed); got != graph.Absent {
		t.Fatalf("with an empty = %v, want absent", got)
	}
	g.SetStatus(seed, "http", graph.StatusOK)
	if got := a.Existence(seed); got != graph.Live {
		t.Fatalf("with an ok = %v, want live", got)
	}
}

// A generated variant is not a dependency. Being unpublished is the normal
// condition of almost every variant — that is exactly what makes it registrable
// — so flagging it as a supply-chain gap made every variant of a package target
// a CRITICAL finding, hundreds per scan, and `--fail-on` returned 2 for every
// package scan whatever level it was given. The CI gate was unusable.
func TestDepConfusionIgnoresGeneratedVariants(t *testing.T) {
	g := graph.New(registry(t))
	seed, _ := g.Seed("package", "npm:lodash")
	g.SetStatus(seed, "npm", graph.StatusEmpty) // absent: a real gap

	// A typo variant of it, equally absent, reached by a VARIANT_OF edge.
	origin := graph.NodeRef{Type: "package", Key: "npm:lodash"}
	v := graph.NodeRef{Type: "package", Key: "npm:ldash"}
	g.Apply(by("omission"), seed, graph.Delta{Edges: []graph.EdgeRef{
		{From: origin, Rel: graph.VariantRel, To: v},
	}})
	vid := mustNode(t, g, "package", "npm:ldash")
	g.SetStatus(vid, "npm", graph.StatusEmpty)

	if err := g.RunAnalyzers(context.Background(), []graph.Analyzer{depconfusion.DepConfusion{}}); err != nil {
		t.Fatalf("analyze: %v", err)
	}

	for _, f := range g.Findings() {
		for _, id := range f.Nodes {
			if id == vid {
				t.Errorf("%s reported %s on a generated variant; an unpublished variant is not a dependency gap",
					f.Severity, f.Kind)
			}
		}
	}
	// The seed itself is still the finding the analyzer exists for.
	var seedFlagged bool
	for _, f := range g.Findings() {
		if f.Kind == "dep-confusion" && len(f.Nodes) > 0 && f.Nodes[0] == seed {
			seedFlagged = true
		}
	}
	if !seedFlagged {
		t.Error("the absent seed package was not flagged; the fix removed the true positive too")
	}
}
