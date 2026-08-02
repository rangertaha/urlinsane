// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package decompose splits a composite target into the entities it is made of:
// an email into a local part and a domain, a domain into its TLD, a repository
// URL into a platform and an owner, a scoped package into its namespace.
//
// Seed expansion is not special-cased. A decomposer is an ordinary
// graph.Operator that happens to match the seed's type, so supporting a new
// input form is one more operator rather than a branch in the engine.
//
// Every edge a decomposer emits is structural (docs/DESIGN.md §1.1): it was
// derived by parsing a key that was already in hand, with no network call. Two
// consequences the rest of the engine depends on — decomposition costs no
// depth, so a composite seed does not spend its observation budget before it
// reaches an IP; and the parts land inside the seed closure, so they are
// legitimate roots for variant generation. MANIFEST looks like decomposition
// but requires fetching the repository, which makes it an observation edge and
// somebody else's job.
package decompose

import (
	"context"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// Operators returns every decomposer, ready to hand to a scheduler.

// New builds a decomposer. It is the constructor the per-type packages under
// this one call, and the reason the struct itself stays unexported: the four
// fields are a closed contract, and a literal in another package could grow a
// fifth without this package noticing.
func New(id string, version int, on string, emits graph.Effects, split func(string) graph.Delta) graph.Operator {
	return &decomposer{id: id, version: version, on: on, emits: emits, split: split}
}

// decomposer is the shared shape of all four. A decomposition is a pure
// function from a node's canonical key to the parts it names, so the only thing
// that actually differs between them is the parsing rule; giving each its own
// type would duplicate the Operator boilerplate four times to say nothing.
type decomposer struct {
	id      string
	version int
	on      string
	emits   graph.Effects
	// split parses an already-canonical key. It returns a zero Delta when the
	// key holds no parts, which Exec turns into an Empty outcome.
	split func(key string) graph.Delta
}

// Id names the operator. Ids are dotted so a plan groups its decomposers.
func (d *decomposer) Id() string { return d.id }

// Version invalidates cached results when a parsing rule changes.
func (d *decomposer) Version() int { return d.version }

// Trigger binds to a node type. There are no Where conditions and no Reads: a
// decomposer consumes only the canonical key, which every View exposes
// unconditionally, so its read-set is empty and it runs exactly once per node.
func (d *decomposer) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{d.on}}}
}

// Emits declares every node type and relation this may produce.
func (d *decomposer) Emits() graph.Effects { return d.emits }

// Resource is empty: decomposers parse a string and make no call, so they
// belong to no rate-limit class and must never be throttled behind one.
func (d *decomposer) Resource() string { return "" }

// Exec parses the key and returns the parts. An input with nothing to decompose
// — "npm:lodash", a single-label host — is an authoritative absence, not a
// failure: Empty closes the pair, whereas Failed would have it retried inside
// the round to no purpose.
func (d *decomposer) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	out := d.split(v.Key())
	if len(out.Nodes) == 0 && len(out.Edges) == 0 && len(out.Props) == 0 {
		return graph.Delta{}, graph.Empty()
	}
	return out, graph.OK()
}
