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

package observe

import (
	"context"
	"errors"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// fakeSources lists sources from a table. Empty means "this kind has none",
// which is itself a case worth testing.
type fakeSources struct {
	byKind map[string][]Source
	err    error
}

func (f fakeSources) Sources(kind string) ([]Source, error) {
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
	return fakeSources{byKind: map[string][]Source{
		"package": {
			{Code: "pypi", Template: "https://pypi.org/project/%s/", CheckURL: "https://pypi.org/pypi/%s/json"},
			{Code: "npm", Template: "https://www.npmjs.com/package/%s", CheckURL: "https://registry.npmjs.org/%s"},
		},
		"username":   {{Code: "github", Template: "https://github.com/%s"}},
		"repository": {{Code: "github", Template: "https://github.com/%s"}},
	}}
}

// TestSourcesEmitExistsOnEdges: a name found on a registry becomes an
// EXISTS_ON edge to a platform node, with the human-facing URL on the edge and
// the registry code on the platform.
func TestSourcesEmitExistsOnEdges(t *testing.T) {
	prober := &fakeProber{exists: map[string]bool{
		"https://registry.npmjs.org/lodash": true,
	}}
	ops := newSourceOps(Options{}, registries(), prober)
	g := run(t, TypePackage, "lodash", ops...)

	if !hasNode(g, TypePlatform, "www.npmjs.com") {
		t.Fatalf("no platform node; graph has %s", dump(g))
	}
	if !hasEdge(g, TypePackage, "lodash", RelExistsOn, TypePlatform, "www.npmjs.com") {
		t.Error("no EXISTS_ON edge to npm")
	}
	if v, ok := edgeProp(t, g, RelExistsOn, "www.npmjs.com", FieldURL); !ok ||
		v.Str() != "https://www.npmjs.com/package/lodash" {
		t.Errorf("edge url = %q (set=%v), want the display url, not the check url", v.Str(), ok)
	}
	if v, ok := prop(t, g, TypePlatform, "www.npmjs.com", FieldCode); !ok || v.Str() != "npm" {
		t.Errorf("platform code = %q (set=%v), want npm", v.Str(), ok)
	}
	// PyPI answered a clean no, so nothing links to it.
	if hasNode(g, TypePlatform, "pypi.org") {
		t.Error("a registry that answered 404 produced a platform node")
	}
	wantStatus(t, g, TypePackage, "lodash", "pkg", graph.StatusOK)
}

// TestSourcesAbsentEverywhereIsEmpty: this is the dependency-confusion signal.
// An internal package name that exists on no public registry is only detectable
// because a clean sweep of 404s is recorded as absence rather than as nothing.
func TestSourcesAbsentEverywhereIsEmpty(t *testing.T) {
	prober := &fakeProber{}
	ops := newSourceOps(Options{}, registries(), prober)
	g := run(t, TypePackage, "acme-internal-utils", ops...)

	wantStatus(t, g, TypePackage, "acme-internal-utils", "pkg", graph.StatusEmpty)
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
	ops := newSourceOps(Options{}, registries(), prober)
	g := run(t, TypePackage, "acme-internal-utils", ops...)
	wantStatus(t, g, TypePackage, "acme-internal-utils", "pkg", graph.StatusFailed)
}

// TestSourcesPartialFailureIsNotAbsence: one registry answering no while
// another never answered does not prove the name is free anywhere.
func TestSourcesPartialFailureIsNotAbsence(t *testing.T) {
	prober := &fakeProber{fails: map[string]error{
		"https://registry.npmjs.org/acme-internal-utils": errors.New("503 Service Unavailable"),
	}}
	ops := newSourceOps(Options{}, registries(), prober)
	g := run(t, TypePackage, "acme-internal-utils", ops...)
	wantStatus(t, g, TypePackage, "acme-internal-utils", "pkg", graph.StatusFailed)
}

// TestSourcesWithNoListIsFailed: nothing was checked, so nothing was
// determined. Reporting it as absence would turn a missing dataset into a page
// of apparently free names.
func TestSourcesWithNoListIsFailed(t *testing.T) {
	ops := newSourceOps(Options{}, fakeSources{}, &fakeProber{})
	g := run(t, TypePackage, "lodash", ops...)
	wantStatus(t, g, TypePackage, "lodash", "pkg", graph.StatusFailed)
}

// TestSourceOperatorsBindToTheirOwnType: one implementation, three patterns.
// The username operator must not run against a package, and vice versa.
func TestSourceOperatorsBindToTheirOwnType(t *testing.T) {
	prober := &fakeProber{exists: map[string]bool{"https://github.com/acme": true}}
	ops := newSourceOps(Options{}, registries(), prober)

	g := run(t, TypeUsername, "acme", ops...)
	if !hasEdge(g, TypeUsername, "acme", RelExistsOn, TypePlatform, "github.com") {
		t.Fatalf("usr did not link the username; graph has %s", dump(g))
	}
	seed := nodeID(t, g, TypeUsername, "acme")
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
	if got := expand("https://registry.npmjs.org/%s", "@acme/tool"); got != "https://registry.npmjs.org/@acme%2Ftool" {
		t.Errorf("expand = %q, want the separator escaped", got)
	}
}

// TestPlatformKeyFallsBackToCode: a source whose template is not a URL has only
// its code for identity.
func TestPlatformKeyFallsBackToCode(t *testing.T) {
	if got := platformKey(Source{Code: "gmail", Template: "gmail.com"}); got != "gmail" {
		t.Errorf("platformKey = %q, want the code", got)
	}
	if got := platformKey(Source{Code: "npm", Template: "https://www.npmjs.com/package/%s"}); got != "www.npmjs.com" {
		t.Errorf("platformKey = %q, want the host", got)
	}
}
