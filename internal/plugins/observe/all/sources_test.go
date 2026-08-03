// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package all

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

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
	g := observetest.Run(t, observe.TypePackage, "npm:lodash", ops...)

	if !observetest.HasNode(g, observe.TypePlatform, "www.npmjs.com") {
		t.Fatalf("no platform node; graph has %s", observetest.Dump(g))
	}
	if !observetest.HasEdge(g, observe.TypePackage, "npm:lodash", observe.RelExistsOn, observe.TypePlatform, "www.npmjs.com") {
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
	observetest.WantStatus(t, g, observe.TypePackage, "npm:lodash", "pkg", graph.StatusOK)
}

// TestSourcesAbsentEverywhereIsEmpty: this is the dependency-confusion signal.
// An internal package name that exists on no public registry is only detectable
// because a clean sweep of 404s is recorded as absence rather than as nothing.
func TestSourcesAbsentEverywhereIsEmpty(t *testing.T) {
	prober := &fakeProber{}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "npm:acme-internal-utils", ops...)

	observetest.WantStatus(t, g, observe.TypePackage, "npm:acme-internal-utils", "pkg", graph.StatusEmpty)
	if n := len(g.Nodes()); n != 1 {
		t.Errorf("an absent package admitted %d extra nodes", n-1)
	}
	// One registry, and the right one. A package key names its registry, so
	// asking pypi about an npm package is a request for something that could
	// never exist — and a 404 from it used to count towards the clean sweep
	// that produces this very signal.
	if len(prober.asked) != 1 ||
		!strings.Contains(prober.asked[0], "registry.npmjs.org") {
		t.Errorf("probed %v, want only the npm registry", prober.asked)
	}
}

// A package naming a registry this build has no source for is undetermined, not
// absent.
//
// It is the same rule as an empty source list, and it matters more here: a real
// package on a registry we do not know about would otherwise read as an
// unclaimed name and raise a CRITICAL dependency-confusion finding.
func TestAPackageOnAnUnknownRegistryIsNotAbsent(t *testing.T) {
	prober := &fakeProber{}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "conda:numpy", ops...)

	observetest.WantStatus(t, g, observe.TypePackage, "conda:numpy", "pkg", graph.StatusFailed)
	if len(prober.asked) != 0 {
		t.Errorf("probed %v for a registry with no source", prober.asked)
	}
}

// The whole key is not the name. A package key is registry-qualified, and the
// registry templates take the bare name, so substituting the key produced
// https://registry.npmjs.org/npm%3Alodash — a 404 from every registry, which
// reported lodash as unpublished with a CRITICAL finding attached.
func TestThePackageProbeUsesTheBareName(t *testing.T) {
	prober := &fakeProber{exists: map[string]bool{
		"https://registry.npmjs.org/lodash": true,
	}}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "npm:lodash", ops...)

	for _, asked := range prober.asked {
		if strings.Contains(asked, "npm%3A") || strings.Contains(asked, "npm:") {
			t.Errorf("probed %q: the registry qualifier reached the URL", asked)
		}
	}
	observetest.WantStatus(t, g, observe.TypePackage, "npm:lodash", "pkg", graph.StatusOK)
}

// TestSourcesUnreachableIsFailedNotEmpty is the bug this port fixes. The old
// collector returned false on a transport error, so an unreachable registry was
// reported exactly like a name nobody had taken — the worst possible confusion
// for a dependency-confusion check.
func TestSourcesUnreachableIsFailedNotEmpty(t *testing.T) {
	prober := &fakeProber{fails: map[string]error{
		"https://registry.npmjs.org/acme-internal-utils": errors.New("connection refused"),
	}}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "npm:acme-internal-utils", ops...)
	observetest.WantStatus(t, g, observe.TypePackage, "npm:acme-internal-utils", "pkg", graph.StatusFailed)
}

// TestSourcesPartialFailureIsNotAbsence: one registry answering no while
// another never answered does not prove the name is free anywhere.
func TestSourcesPartialFailureIsNotAbsence(t *testing.T) {
	prober := &fakeProber{fails: map[string]error{
		"https://registry.npmjs.org/acme-internal-utils": errors.New("503 Service Unavailable"),
	}}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "npm:acme-internal-utils", ops...)
	observetest.WantStatus(t, g, observe.TypePackage, "npm:acme-internal-utils", "pkg", graph.StatusFailed)
}

// TestSourcesWithNoListIsFailed: nothing was checked, so nothing was
// determined. Reporting it as absence would turn a missing dataset into a page
// of apparently free names.
func TestSourcesWithNoListIsFailed(t *testing.T) {
	ops := allSourceOps(observe.Options{}, fakeSources{}, &fakeProber{})
	g := observetest.Run(t, observe.TypePackage, "npm:lodash", ops...)
	observetest.WantStatus(t, g, observe.TypePackage, "npm:lodash", "pkg", graph.StatusFailed)
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

// slowProber blocks each call until its context expires, so a probe consumes
// exactly the deadline it was given and no more.
type slowProber struct{ asked []string }

func (p *slowProber) Exists(ctx context.Context, rawURL string) (bool, error) {
	p.asked = append(p.asked, rawURL)
	<-ctx.Done()
	return false, ctx.Err()
}

// The per-source timeout is per source. Taken once around the whole loop, the
// first registry consumed the entire budget and every later one was skipped
// without being contacted -- so a sweep of sixty-six username platforms could
// report on one of them and stay silent about the rest.
func TestEachSourceGetsItsOwnTimeout(t *testing.T) {
	// A username, not a package: a username key carries no qualifier, so every
	// platform legitimately applies — which is the sweep this invariant is
	// about, and the one the original bug silenced sixty-five sixty-sixths of.
	prober := &slowProber{}
	ops := allSourceOps(observe.Options{Timeout: 10 * time.Millisecond}, multiUser(), prober)
	observetest.Run(t, observe.TypeUsername, "lodash", ops...)

	if len(prober.asked) != 2 {
		t.Errorf("probed %v, want both platforms contacted; a shared deadline stops after the first",
			prober.asked)
	}
}

// multiUser is registries() with a second username platform, so a sweep that
// must contact every source has more than one to contact.
func multiUser() fakeSources {
	r := registries()
	r.byKind["username"] = []observe.Source{
		{Code: "github", URL: "https://github.com/%s"},
		{Code: "gitlab", URL: "https://gitlab.com/%s"},
	}
	return r
}

// A source that times out is that source undetermined, not the sweep
// abandoned -- and an undetermined sweep is never absence.
func TestATimedOutSourceDoesNotStopTheSweep(t *testing.T) {
	prober := &slowProber{}
	ops := allSourceOps(observe.Options{Timeout: 10 * time.Millisecond}, registries(), prober)
	g := observetest.Run(t, observe.TypePackage, "npm:acme-internal-utils", ops...)

	// Every source timed out, so nothing is proven absent.
	observetest.WantStatus(t, g, observe.TypePackage, "npm:acme-internal-utils", "pkg", graph.StatusTimeout)
}

// A repo is undetermined by this operator, not absent.
//
// datasets/sources/repos.lst is by its own header a list of *namespace*
// endpoints — "%s is the org/user/project namespace" — while a repo key is
// host/owner/name. api.github.com/users/rangertaha/urlinsane 404s for every real
// repository (the answer lives at /repos/, which returns 200), so this operator
// has never been able to answer for a repo. Before the qualifier was honoured it
// asked all nine forges and reported "live" off the one that 200s on nonsense.
//
// The namespace question these endpoints can answer is already answered
// elsewhere: canonRepo decomposes the key into platform:github.com and
// username:rangertaha, and the username sweep probes that. So the honest state
// is unknown — turning "I could not check" into "this name is free" is the one
// inference this codebase exists to prevent.
func TestARepoIsUndeterminedNotAbsent(t *testing.T) {
	prober := &fakeProber{}
	ops := allSourceOps(observe.Options{}, registries(), prober)
	g := observetest.Run(t, observe.TypeRepo, "github.com/acme/tool", ops...)

	observetest.WantStatus(t, g, observe.TypeRepo, "github.com/acme/tool", "repo", graph.StatusFailed)
	if len(prober.asked) != 0 {
		t.Errorf("probed %v with namespace endpoints that cannot answer for a repo path",
			prober.asked)
	}
}
