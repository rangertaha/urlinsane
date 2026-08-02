// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"context"
	"errors"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/observetest"
)

// fakeSources lists sources from a table. Empty means "this kind has none",
// which is itself a case worth testing.
type fakeSources struct {
	byKind map[string][]observe.Source
	err    error
}

func (f fakeSources) Sources(kind string) ([]observe.Source, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.byKind[kind], nil
}

// fakeProber answers existence from a table: present in exists means found,
// present in fails means no answer at all.
type fakeProber struct {
	exists map[string]bool
	fails  map[string]error
	asked  []string
}

func (p *fakeProber) Exists(_ context.Context, rawURL string) (bool, error) {
	p.asked = append(p.asked, rawURL)
	if err, ok := p.fails[rawURL]; ok {
		return false, err
	}
	return p.exists[rawURL], nil
}

func registries() fakeSources {
	return fakeSources{byKind: map[string][]observe.Source{
		"package": {
			{Code: "pypi", URL: "https://pypi.org/project/%s/", CheckURL: "https://pypi.org/pypi/%s/json"},
			{Code: "npm", URL: "https://www.npmjs.com/package/%s", CheckURL: "https://registry.npmjs.org/%s"},
		},
		"username":   {{Code: "github", URL: "https://github.com/%s"}},
		"repository": {{Code: "github", URL: "https://github.com/%s"}},
	}}
}

// TestSourcesEmitExistsOnEdges: a name found on a registry becomes an
// EXISTS_ON edge to a platform node, with the human-facing URL on the edge and
// the registry code on the platform.
func TestSourcesEmitExistsOnEdges(t *testing.T) {
	prober := &fakeProber{exists: map[string]bool{
		"https://registry.npmjs.org/lodash": true,
	}}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "lodash", ops...)

	if !observetest.HasNode(g, observe.TypePlatform, "www.npmjs.com") {
		t.Fatalf("no platform node; graph has %s", observetest.Dump(g))
	}
	if !observetest.HasEdge(g, observe.TypePackage, "lodash", observe.RelExistsOn, observe.TypePlatform, "www.npmjs.com") {
		t.Error("no EXISTS_ON edge to npm")
	}
	if v, ok := observetest.EdgeProp(t, g, observe.RelExistsOn, "www.npmjs.com", observe.FieldURL); !ok ||
		v.Str() != "https://www.npmjs.com/package/lodash" {
		t.Errorf("edge url = %q (set=%v), want the display url, not the check url", v.Str(), ok)
	}
	if v, ok := observetest.Prop(t, g, observe.TypePlatform, "www.npmjs.com", observe.FieldCode); !ok || v.Str() != "npm" {
		t.Errorf("platform code = %q (set=%v), want npm", v.Str(), ok)
	}
	// PyPI answered a clean no, so nothing links to it.
	if observetest.HasNode(g, observe.TypePlatform, "pypi.org") {
		t.Error("a registry that answered 404 produced a platform node")
	}
	observetest.WantStatus(t, g, observe.TypePackage, "lodash", "pkg", graph.StatusOK)
}

// TestSourcesAbsentEverywhereIsEmpty: this is the dependency-confusion signal.
// An internal package name that exists on no public registry is only detectable
// because a clean sweep of 404s is recorded as absence rather than as nothing.
func TestSourcesAbsentEverywhereIsEmpty(t *testing.T) {
	prober := &fakeProber{}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "acme-internal-utils", ops...)

	observetest.WantStatus(t, g, observe.TypePackage, "acme-internal-utils", "pkg", graph.StatusEmpty)
	if n := len(g.Nodes()); n != 1 {
		t.Errorf("an absent package admitted %d extra nodes", n-1)
	}
	if len(prober.asked) != 2 {
		t.Errorf("probed %v, want both registries checked", prober.asked)
	}
}

// TestSourcesUnreachableIsFailedNotEmpty is the bug this port fixes. The old
// collector returned false on a transport error, so an unreachable registry was
// reported exactly like a name nobody had taken — the worst possible confusion
// for a dependency-confusion check.
func TestSourcesUnreachableIsFailedNotEmpty(t *testing.T) {
	prober := &fakeProber{fails: map[string]error{
		"https://registry.npmjs.org/acme-internal-utils": errors.New("connection refused"),
		"https://pypi.org/pypi/acme-internal-utils/json": errors.New("connection refused"),
	}}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "acme-internal-utils", ops...)
	observetest.WantStatus(t, g, observe.TypePackage, "acme-internal-utils", "pkg", graph.StatusFailed)
}

// TestSourcesPartialFailureIsNotAbsence: one registry answering no while
// another never answered does not prove the name is free anywhere.
func TestSourcesPartialFailureIsNotAbsence(t *testing.T) {
	prober := &fakeProber{fails: map[string]error{
		"https://registry.npmjs.org/acme-internal-utils": errors.New("503 Service Unavailable"),
	}}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "acme-internal-utils", ops...)
	observetest.WantStatus(t, g, observe.TypePackage, "acme-internal-utils", "pkg", graph.StatusFailed)
}

// TestSourcesWithNoListIsFailed: nothing was checked, so nothing was
// determined. Reporting it as absence would turn a missing dataset into a page
// of apparently free names.
func TestSourcesWithNoListIsFailed(t *testing.T) {
	ops := allSourceOps(observe.Options{}, fakeSources{}, &fakeProber{})
	g := observetest.Run(t, observe.TypePackage, "lodash", ops...)
	observetest.WantStatus(t, g, observe.TypePackage, "lodash", "pkg", graph.StatusFailed)
}

// TestSourceOperatorsBindToTheirOwnType: one implementation, three patterns.
// The username operator must not observetest.Run against a package, and vice versa.
func TestSourceOperatorsBindToTheirOwnType(t *testing.T) {
	prober := &fakeProber{exists: map[string]bool{"https://github.com/acme": true}}
	ops := allSourceOps(observe.Options{}, registries(), prober)

	g := observetest.Run(t, observe.TypeUsername, "acme", ops...)
	if !observetest.HasEdge(g, observe.TypeUsername, "acme", observe.RelExistsOn, observe.TypePlatform, "github.com") {
		t.Fatalf("usr did not link the username; graph has %s", observetest.Dump(g))
	}
	seed := observetest.NodeID(t, g, observe.TypeUsername, "acme")
	if _, ran := g.Status(seed, "pkg"); ran {
		t.Error("the package operator ran against a username node")
	}
	if _, ran := g.Status(seed, "repo"); ran {
		t.Error("the repository operator ran against a username node")
	}
}

// TestExpandEscapesTheName keeps a scoped npm name from silently becoming a
// different URL path: without escaping, @acme/tool would address the "acme"
// namespace's "tool" resource rather than the package.
func TestExpandEscapesTheName(t *testing.T) {
	if got := observe.Expand("https://registry.npmjs.org/%s", "@acme/tool"); got != "https://registry.npmjs.org/@acme%2Ftool" {
		t.Errorf("expand = %q, want the separator escaped", got)
	}
}

// TestPlatformKeyFallsBackToCode: a source whose template is not a URL has only
// its code for identity.
func TestPlatformKeyFallsBackToCode(t *testing.T) {
	if got := observe.PlatformKey(observe.Source{Code: "gmail", URL: "gmail.com"}); got != "gmail" {
		t.Errorf("platformKey = %q, want the code", got)
	}
	if got := observe.PlatformKey(observe.Source{Code: "npm", URL: "https://www.npmjs.com/package/%s"}); got != "www.npmjs.com" {
		t.Errorf("platformKey = %q, want the host", got)
	}
}

// allSourceOps builds the three source operators for tests that exercise the
// family as one.
func allSourceOps(o observe.Options, list observe.SourceLister, prober observe.Prober) []graph.Operator {
	return []graph.Operator{
		observe.NewSourceOp(o, "pkg", observe.TypePackage, "package", list, prober),
		observe.NewSourceOp(o, "usr", observe.TypeUsername, "username", list, prober),
		observe.NewSourceOp(o, "repo", observe.TypeRepo, "repository", list, prober),
	}
}
