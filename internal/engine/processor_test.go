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

package engine

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/entity"
	"github.com/rangertaha/urlinsane/internal/store"
)

// fakeConfig is a minimal internal.Config for exercising the collector stage
// without the heavy real config (filesystem, embedded databases).
type fakeConfig struct {
	algorithms []internal.Algorithm
	collectors []internal.Collector
	analyzers  []internal.Analyzer
	output     internal.Output
	st         *store.Store
	workers    int
	timeout    time.Duration
	delay      time.Duration
	random     time.Duration
}

func (c *fakeConfig) Target() string                   { return "example.com" }
func (c *fakeConfig) EntityType() entity.Type          { return entity.Domain }
func (c *fakeConfig) Keyboards() []internal.Keyboard   { return nil }
func (c *fakeConfig) Languages() []internal.Language   { return nil }
func (c *fakeConfig) Algorithms() []internal.Algorithm { return c.algorithms }
func (c *fakeConfig) Collectors() []internal.Collector { return c.collectors }
func (c *fakeConfig) Analyzers() []internal.Analyzer   { return c.analyzers }
func (c *fakeConfig) Store() *store.Store              { return c.st }
func (c *fakeConfig) Output() internal.Output          { return c.output }
func (c *fakeConfig) Workers() int {
	if c.workers == 0 {
		return 1
	}
	return c.workers
}
func (c *fakeConfig) Delay() time.Duration   { return c.delay }
func (c *fakeConfig) Random() time.Duration  { return c.random }
func (c *fakeConfig) Timeout() time.Duration { return c.timeout }
func (c *fakeConfig) TTL() time.Duration     { return 0 }
func (c *fakeConfig) Distance() int          { return 0 }
func (c *fakeConfig) Regex() string          { return "" }
func (c *fakeConfig) Verbose() bool          { return false }
func (c *fakeConfig) Progress() bool         { return false }
func (c *fakeConfig) Banner() bool           { return false }
func (c *fakeConfig) Format() string         { return "" }
func (c *fakeConfig) Filters() []string      { return nil }
func (c *fakeConfig) Registered() bool       { return false }
func (c *fakeConfig) Unregistered() bool     { return false }
func (c *fakeConfig) Summary() bool          { return false }
func (c *fakeConfig) Dir() string            { return "" }
func (c *fakeConfig) File() string           { return "" }
func (c *fakeConfig) AssetDir() string       { return "" }

// fakeCollector records its execution and optionally blocks on ctx.
type fakeCollector struct {
	id     string
	order  int
	deps   []string
	onExec func(id string) // invoked at the start of Exec
	block  bool            // if true, block until ctx is cancelled
	calls  *int64          // optional atomic call counter
}

func (f *fakeCollector) Id() string             { return f.id }
func (f *fakeCollector) Order() int             { return f.order }
func (f *fakeCollector) Dependencies() []string { return f.deps }
func (f *fakeCollector) Types() []entity.Type   { return nil }
func (f *fakeCollector) Description() string    { return f.id }
func (f *fakeCollector) Exec(ctx context.Context, d *db.Domain) (*db.Domain, error) {
	if f.calls != nil {
		atomic.AddInt64(f.calls, 1)
	}
	if f.onExec != nil {
		f.onExec(f.id)
	}
	if f.block {
		<-ctx.Done()
		return d, ctx.Err()
	}
	// Realistic mutation of the shared variant.
	d.Dns = append(d.Dns, &db.Dns{Type: "X", Value: f.id})
	return d, nil
}

// fakeAlgorithm produces two fixed variants from an origin.
type fakeAlgorithm struct{ id string }

func (a *fakeAlgorithm) Id() string           { return a.id }
func (a *fakeAlgorithm) Name() string         { return a.id }
func (a *fakeAlgorithm) Description() string  { return a.id }
func (a *fakeAlgorithm) Types() []entity.Type { return nil }
func (a *fakeAlgorithm) Exec(origin *db.Domain) ([]*db.Domain, error) {
	return []*db.Domain{
		{Name: origin.Name + ".v1"},
		{Name: origin.Name + ".v2"},
	}, nil
}

// fakeAnalyzer records the (origin, variant) pairs it is given.
type fakeAnalyzer struct {
	id    string
	mu    *sync.Mutex
	pairs *[][2]string
}

func (a *fakeAnalyzer) Id() string          { return a.id }
func (a *fakeAnalyzer) Order() int          { return 0 }
func (a *fakeAnalyzer) Description() string { return a.id }
func (a *fakeAnalyzer) Headers() []string   { return nil }
func (a *fakeAnalyzer) Exec(ctx context.Context, origin, variant *db.Domain) (*db.Domain, error) {
	a.mu.Lock()
	*a.pairs = append(*a.pairs, [2]string{origin.Name, variant.Name})
	a.mu.Unlock()
	return variant, nil
}

// fakeOutput is a no-op internal.Output for exercising the Output stage.
type fakeOutput struct{}

func (fakeOutput) Id() string          { return "fake" }
func (fakeOutput) Description() string { return "fake" }
func (fakeOutput) Read(*db.Domain)     {}
func (fakeOutput) Write()              {}
func (fakeOutput) Save(string)         {}
func (fakeOutput) Report()             {}

// TestOutput_PersistsToStore verifies the Output stage persists live results to
// the IPLD store as the sole source of truth (no GORM).
func TestOutput_PersistsToStore(t *testing.T) {
	s, err := store.OpenDir(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	cfg := &fakeConfig{output: fakeOutput{}, st: s}
	u := New(cfg)

	live := &db.Domain{Type: entity.Domain, Name: "exmaple.com", Dns: []*db.Dns{{Type: "A", Value: "1.1.1.1"}}}
	u.Output(context.Background(), feed(live))

	got, ok, err := s.Get(entity.Domain, "exmaple.com")
	if err != nil || !ok {
		t.Fatalf("store missing dual-written entity: ok=%v err=%v", ok, err)
	}
	if len(got.Dns) != 1 || got.Dns[0].Value != "1.1.1.1" {
		t.Fatalf("dual-written entity wrong: %+v", got)
	}
}

// TestLoad_PassThrough verifies Load forwards variants unchanged and does NOT
// pre-populate collected data from the store (which caused records to duplicate
// across re-scans, since collectors append).
func TestLoad_PassThrough(t *testing.T) {
	s, err := store.OpenDir(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	// Seed the store with cached result data for the same name.
	if _, err := s.Put(&db.Domain{Type: entity.Domain, Name: "exmaple.com", Dns: []*db.Dns{{Type: "A", Value: "9.9.9.9"}}}); err != nil {
		t.Fatalf("seed store: %v", err)
	}

	u := New(&fakeConfig{st: s})
	in := &db.Domain{Type: entity.Domain, Name: "exmaple.com", Levenshtein: 1}

	out := u.Load(context.Background(), feed(in))
	var got *db.Domain
	for d := range out {
		got = d
	}
	if got != in {
		t.Fatal("Load should pass the in-flight variant through unchanged")
	}
	if len(got.Dns) != 0 {
		t.Fatalf("Load must not pre-populate collected data; got %d dns records", len(got.Dns))
	}
}

// drain reads the channel to completion, failing if it does not close in time.
func drain(t *testing.T, out <-chan *db.Domain, within time.Duration) int {
	t.Helper()
	n := 0
	timeout := time.After(within)
	for {
		select {
		case _, ok := <-out:
			if !ok {
				return n
			}
			n++
		case <-timeout:
			t.Fatalf("collector stage did not finish within %s", within)
			return n
		}
	}
}

func feed(domains ...*db.Domain) <-chan *db.Domain {
	in := make(chan *db.Domain, len(domains))
	for _, d := range domains {
		in <- d
	}
	close(in)
	return in
}

// TestCollectors_DependencyOrder verifies a collector never runs before a
// collector it declares a dependency on (c <- b <- a).
func TestCollectors_DependencyOrder(t *testing.T) {
	var order []string
	rec := func(id string) { order = append(order, id) }
	cfg := &fakeConfig{
		workers: 1,
		collectors: []internal.Collector{
			&fakeCollector{id: "c", deps: []string{"b"}, onExec: rec},
			&fakeCollector{id: "b", deps: []string{"a"}, onExec: rec},
			&fakeCollector{id: "a", onExec: rec},
		},
	}
	u := New(cfg)
	out := u.Collectors(context.Background(), feed(&db.Domain{Name: "x"}))
	drain(t, out, 5*time.Second)

	want := []string{"a", "b", "c"}
	if len(order) != 3 || order[0] != want[0] || order[1] != want[1] || order[2] != want[2] {
		t.Fatalf("dependency order wrong: got %v, want %v", order, want)
	}
}

// TestCollectors_ContextCancel verifies the timeout/cancel context reaches
// Exec and the stage drains promptly instead of hanging.
func TestCollectors_ContextCancel(t *testing.T) {
	cfg := &fakeConfig{
		workers:    1,
		collectors: []internal.Collector{&fakeCollector{id: "blocker", block: true}},
	}
	u := New(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	out := u.Collectors(ctx, feed(&db.Domain{Name: "x"}))
	cancel() // collector is blocked on ctx.Done()

	if n := drain(t, out, 5*time.Second); n > 1 {
		t.Fatalf("expected at most 1 emitted variant, got %d", n)
	}
}

// TestCollectors_NoRace runs many variants through a multi-level DAG across a
// worker pool. Intended to be run under `go test -race`.
func TestCollectors_NoRace(t *testing.T) {
	var calls int64
	cfg := &fakeConfig{
		workers: 8,
		collectors: []internal.Collector{
			&fakeCollector{id: "a", calls: &calls},
			&fakeCollector{id: "b", deps: []string{"a"}, calls: &calls},
			&fakeCollector{id: "geo", deps: []string{"a"}, calls: &calls},
		},
	}
	u := New(cfg)

	const variants = 200
	domains := make([]*db.Domain, variants)
	for i := range domains {
		domains[i] = &db.Domain{Name: "v"}
	}
	out := u.Collectors(context.Background(), feed(domains...))
	got := drain(t, out, 10*time.Second)

	if got != variants {
		t.Fatalf("emitted %d variants, want %d", got, variants)
	}
	if want := int64(variants * 3); atomic.LoadInt64(&calls) != want {
		t.Fatalf("collector calls = %d, want %d", calls, want)
	}
}

// TestAlgorithms_StampsOrigin verifies the algorithms stage tags every variant
// with the origin it was generated from, so analyzers can pair them later.
func TestAlgorithms_StampsOrigin(t *testing.T) {
	cfg := &fakeConfig{algorithms: []internal.Algorithm{&fakeAlgorithm{id: "f"}}}
	u := New(cfg)

	origin := &db.Domain{Name: "example.com"}
	out := u.Algorithms(context.Background(), feed(origin))

	var got []*db.Domain
	timeout := time.After(5 * time.Second)
	for {
		select {
		case d, ok := <-out:
			if !ok {
				goto done
			}
			got = append(got, d)
		case <-timeout:
			t.Fatal("algorithms stage did not finish in time")
		}
	}
done:
	if len(got) != 2 {
		t.Fatalf("want 2 variants, got %d", len(got))
	}
	for _, v := range got {
		if v.Origin != origin {
			t.Fatalf("variant %s missing origin pointer", v.Name)
		}
	}
}

// TestAnalyzers_OriginVariantPairing verifies the analyzer stage pairs each
// variant with the origin it was derived from (variant.Origin).
func TestAnalyzers_OriginVariantPairing(t *testing.T) {
	var pairs [][2]string
	var mu sync.Mutex
	cfg := &fakeConfig{
		analyzers: []internal.Analyzer{&fakeAnalyzer{id: "fake", mu: &mu, pairs: &pairs}},
	}
	u := New(cfg)

	origin := &db.Domain{Name: "example.com"}
	v1 := &db.Domain{Name: "exmaple.com", Origin: origin}
	v2 := &db.Domain{Name: "exampel.com", Origin: origin}
	out := u.Analyzers(context.Background(), feed(v1, v2))
	drain(t, out, 5*time.Second)

	if len(pairs) != 2 {
		t.Fatalf("want 2 analyzed pairs, got %d: %v", len(pairs), pairs)
	}
	for _, p := range pairs {
		if p[0] != "example.com" {
			t.Fatalf("origin should be example.com, got %v", p)
		}
		if p[1] != "exmaple.com" && p[1] != "exampel.com" {
			t.Fatalf("unexpected variant in pair %v", p)
		}
	}
}
