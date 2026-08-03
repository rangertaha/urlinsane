// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package plugins holds everything that acts on the graph, and the registries
// that let more be added.
//
// There are three registries and no more, because the graph engine has three
// places where behaviour can be added. The pipeline this replaced had five —
// algorithms, collectors, analyzers, languages, outputs — of which two were not
// behaviour at all. Languages and keyboards are *data*: a language is added by
// adding a dataset directory, a keyboard by a layout in pkg/kb, and neither
// needs code. Output formats are a closed set the report projects into.
//
// # Layout
//
// One package per extension point, each holding the shipped implementations of
// that point:
//
//	internal/plugins/
//	  plugins.go   the three registries
//	  decompose/   domain, email, pkg, repo
//	  variant/     the 27 algorithms, one directory each
//	  observe/     dns, ptr, whois, idn, geo, pkg, usr, repo
//	  analyze/     campaign, scoring, depconfusion
//	  report/      table, json, ndjson, csv, dot
//
// One directory per plugin, and each family is a library plus its plugins. The
// family package holds what its plugins share — observe owns Options, the
// per-call timeout and the schema vocabulary; variant owns Spec, the Generate
// signature and the keyboard and language combinators — and each plugin
// imports it.
//
// Composition lives in <family>/all, which imports the plugins. That
// separation is load-bearing rather than tidy: a family package that both held
// the shared library and listed its plugins would make every plugin import its
// siblings through it, which is a cycle the compiler rejects.
//
// Grouped by what a thing *is* rather than by what it targets. An earlier
// arrangement grouped by target — network, service, social, repo — on the
// grounds that an npm plugin contributes both an operator and an algorithm and
// should live in one directory. That reasoning holds for a third-party plugin
// and not for these: the shipped observers share Options, the per-call timeout
// and the whole schema vocabulary, so splitting dns from whois from geo meant
// either duplicating that or exporting the package's internals to itself. The
// cohesion is real, and the directory tree should not fight it.
//
// A plugin that does target one service is still free to be one directory —
// it just registers through the same three functions, from wherever it lives.
//
// # What is registered and what is composed
//
// The packages above are composed directly by internal/scan, not registered. A
// registry entry for what always runs would be a second place to look for it,
// and scan already has to name them to pass each its own configuration.
//
// The registries exist for what does *not* always run.
//
// # Registration
//
// A plugin registers in init and declares its own settings defaults:
//
//	func init() {
//	    plugins.AddOperator("acme", map[string]any{"timeout": 5},
//	        func(e plugins.Env) ([]graph.Operator, error) {
//	            return []graph.Operator{&acmeOp{
//	                timeout: e.Settings.Int("timeout", 5),
//	                prober:  e.Observe.Prober,
//	            }}, nil
//	        })
//	}
//
// The defaults reach the config file, so `plugins:` lists what can be
// configured rather than leaving the user to discover it, and the values a
// plugin receives are its defaults with the file's overrides applied per key.
//
// Env carries the engine's shared services alongside those settings, so a
// plugin gets the same resolver, prober and datasets the shipped operators do —
// including the fakes a test injects. Returning no operators is how a plugin
// declines when a service it needs is absent, which keeps --explain honest.
//
// # What a plugin cannot do
//
// It cannot reorder the run. Operators bind to data through a trigger and the
// scheduler decides what runs when (§4.1); there is deliberately no way for a
// plugin to say it comes after another. The pipeline's DependsOn list is what
// made its plugin order load-bearing and its cache unsound, and the DAG exists
// so that a plugin declares only what it *reads* and what it *emits*.
package plugins

import (
	"fmt"
	"sort"
	"sync"

	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// Settings are one plugin's resolved settings: its declared defaults with the
// config file's overrides applied per key.
type Settings map[string]any

// String returns a string setting, or def when absent or the wrong type.
//
// The accessors coerce rather than error because a setting is user input from
// a YAML file: a plugin that had to validate every field would either duplicate
// that logic or ignore it, and ignoring it is what a silent zero value causes.
func (s Settings) String(key, def string) string {
	if v, ok := s[key].(string); ok {
		return v
	}
	return def
}

// Int returns an integer setting, accepting the int and float64 that YAML
// produces for the same literal.
func (s Settings) Int(key string, def int) int {
	switch v := s[key].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return def
}

// Bool returns a boolean setting, or def.
func (s Settings) Bool(key string, def bool) bool {
	if v, ok := s[key].(bool); ok {
		return v
	}
	return def
}

// Strings returns a list setting, accepting the []any that YAML produces.
func (s Settings) Strings(key string) []string {
	switch v := s[key].(type) {
	case []string:
		return v
	case []any:
		out := make([]string, 0, len(v))
		for _, e := range v {
			if str, ok := e.(string); ok {
				out = append(out, str)
			}
		}
		return out
	}
	return nil
}

// Env is everything a plugin is built from: its own resolved settings, plus the
// shared services the engine supplies.
//
// The services are here rather than reachable through a global because an
// operator must be a pure function of its inputs to be testable as one (§4.3).
// A plugin that reached for a package-level resolver could not be given a fake,
// and two tests running concurrently would fight over it.
//
// Settings and services are separate fields because they come from different
// places and fail differently: settings are user input from a YAML file that a
// plugin coerces and defaults, services are injected by the caller and are
// either present or deliberately absent.
type Env struct {
	// Settings are this plugin's declared defaults with the config file's
	// per-key overrides applied.
	Settings Settings
	// Observe carries the external services observation plugins need — a
	// resolver, a whois client, a geolocator, an HTTP prober, the source
	// lists. A nil field means the plugin should omit the operators that
	// needed it rather than register ones that can only fail.
	Observe observe.Options
	// Variant carries the data the generating plugins run over: languages,
	// keyboard layouts, suffixes, combo vocabulary.
	Variant variant.Options
}

// OperatorFunc builds the operators a plugin contributes. Returning several is
// allowed: one plugin may cover a family, as the dns operators do.
//
// Returning none is also allowed, and is how a plugin declines: geo with no
// geolocator returns nil, so it is absent from the compiled plan and from
// --list operators, rather than present and guaranteed to fail.
type OperatorFunc func(Env) ([]graph.Operator, error)

// AnalyzerFunc builds the analyzer a plugin contributes.
type AnalyzerFunc func(Env) (graph.Analyzer, error)

// AlgorithmFunc builds the variant algorithms a plugin contributes.
//
// A function rather than a plain Spec because most algorithms are data-driven:
// the keyboard and language families expand to one algorithm per layout or per
// language, and how many there are is not known until Env supplies the data.
type AlgorithmFunc func(Env) ([]variant.Spec, error)

// Source supplies a plugin's resolved settings. *settings.File satisfies it;
// nil means every plugin gets its declared defaults.
//
// An interface rather than the concrete type so the registry does not depend on
// how settings are stored, and so a test can supply values without a file.
type Source interface {
	Apply(id string) map[string]any
}

type opEntry struct {
	id string
	fn OperatorFunc
}

type anEntry struct {
	id string
	fn AnalyzerFunc
}

type algEntry struct {
	id string
	fn AlgorithmFunc
}

var (
	mu         sync.RWMutex
	operators  []opEntry
	analyzers  []anEntry
	algorithms []algEntry
)

// AddOperator registers an operator plugin and its settings defaults.
//
// Registering the same id twice is a programming error rather than an override:
// two operators sharing an id would collide in the plan, in the cache key and
// in --list operators, and the second registration silently winning is the kind
// of thing that is found much later.
func AddOperator(id string, defaults map[string]any, fn OperatorFunc) {
	mu.Lock()
	defer mu.Unlock()
	for _, e := range operators {
		if e.id == id {
			panic(fmt.Sprintf("plugins: operator %q registered twice", id))
		}
	}
	operators = append(operators, opEntry{id: id, fn: fn})
	if len(defaults) > 0 {
		config.Register(id, defaults)
	}
}

// AddAnalyzer registers an analyzer plugin and its settings defaults.
func AddAnalyzer(id string, defaults map[string]any, fn AnalyzerFunc) {
	mu.Lock()
	defer mu.Unlock()
	for _, e := range analyzers {
		if e.id == id {
			panic(fmt.Sprintf("plugins: analyzer %q registered twice", id))
		}
	}
	analyzers = append(analyzers, anEntry{id: id, fn: fn})
	if len(defaults) > 0 {
		config.Register(id, defaults)
	}
}

// AddAlgorithm registers a family of variant algorithms under one plugin id.
//
// The engine wraps each returned Spec in an operator declaring VARIANT_OF,
// which is what subjects it to the terminal rule and the seed-closure
// restriction (§4).
//
// The registered id names the *plugin*, not the algorithms: the keyboard plugin
// registers once and yields one algorithm per distinct layout. Settings and
// --list are keyed on the plugin id; --algorithm selects on the Spec ids.
func AddAlgorithm(id string, defaults map[string]any, fn AlgorithmFunc) {
	mu.Lock()
	defer mu.Unlock()
	for _, e := range algorithms {
		if e.id == id {
			panic(fmt.Sprintf("plugins: algorithm %q registered twice", id))
		}
	}
	algorithms = append(algorithms, algEntry{id: id, fn: fn})
	if len(defaults) > 0 {
		config.Register(id, defaults)
	}
}

// AddSpec registers a single static algorithm, for a plugin whose output does
// not depend on Env. Sugar over AddAlgorithm.
func AddSpec(spec variant.Spec) {
	AddAlgorithm(spec.ID, nil, func(Env) ([]variant.Spec, error) {
		return []variant.Spec{spec}, nil
	})
}

// Operators builds every registered operator plugin, in id order.
//
// Sorted because the plan is compiled from this set and hashed (§5):
// registration order is package-initialisation order, and letting it through
// would give the same run a different plan hash between builds.
func Operators(src Source, env Env) ([]graph.Operator, error) {
	mu.RLock()
	entries := append([]opEntry(nil), operators...)
	mu.RUnlock()
	sort.Slice(entries, func(i, j int) bool { return entries[i].id < entries[j].id })

	var out []graph.Operator
	for _, e := range entries {
		ops, err := e.fn(env.for_(src, e.id))
		if err != nil {
			return nil, fmt.Errorf("plugins: operator %q: %w", e.id, err)
		}
		out = append(out, ops...)
	}
	return out, nil
}

// Analyzers builds every registered analyzer plugin, in id order.
func Analyzers(src Source, env Env) ([]graph.Analyzer, error) {
	mu.RLock()
	entries := append([]anEntry(nil), analyzers...)
	mu.RUnlock()
	sort.Slice(entries, func(i, j int) bool { return entries[i].id < entries[j].id })

	var out []graph.Analyzer
	for _, e := range entries {
		a, err := e.fn(env.for_(src, e.id))
		if err != nil {
			return nil, fmt.Errorf("plugins: analyzer %q: %w", e.id, err)
		}
		if a != nil {
			out = append(out, a)
		}
	}
	return out, nil
}

// Algorithms returns every registered algorithm, in Spec id order.
func Algorithms(src Source, env Env) ([]variant.Spec, error) {
	mu.RLock()
	entries := append([]algEntry(nil), algorithms...)
	mu.RUnlock()
	sort.Slice(entries, func(i, j int) bool { return entries[i].id < entries[j].id })

	var out []variant.Spec
	for _, e := range entries {
		specs, err := e.fn(env.for_(src, e.id))
		if err != nil {
			return nil, fmt.Errorf("plugins: algorithm %q: %w", e.id, err)
		}
		out = append(out, specs...)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// for_ returns the Env one plugin is built from: the shared services unchanged,
// with that plugin's settings resolved — its declared defaults when there is no
// source, and the file's per-key overrides on top of them when there is.
func (e Env) for_(src Source, id string) Env {
	if src == nil {
		e.Settings = Settings(config.Defaults(id))
	} else {
		e.Settings = Settings(src.Apply(id))
	}
	return e
}

// reset clears every registry. Tests only.
func reset() {
	mu.Lock()
	defer mu.Unlock()
	operators, analyzers, algorithms = nil, nil, nil
}
