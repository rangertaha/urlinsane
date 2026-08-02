// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package ndjson renders a report as one object per node.
//
// It promises no ordering, which is what lets a record be emitted the moment
// its node is done rather than when the scan is.
package ndjson

import (
	"encoding/json"
	"io"

	"github.com/rangertaha/urlinsane/internal/plugins/report"
)

// renderNDJSON writes one object per line.
//
// This is the event stream (§11), not a report: it exists to be piped into
// something else while a scan runs. It carries a `kind` discriminator on every
// line because a consumer reading line by line has no enclosing object to tell
// it what it just got. The ordering here follows the graph, but the format
// makes no ordering claim and a consumer must not rely on one.
func Render(w io.Writer, r report.Report) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	emit := func(kind string, v any) error {
		// Merge the discriminator into the payload rather than nesting, so
		// `jq 'select(.kind=="node") | .key'` works without a wrapper path.
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		var m map[string]json.RawMessage
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}
		m["kind"], _ = json.Marshal(kind)
		return enc.Encode(m)
	}

	if err := emit("run", struct {
		Target  string   `json:"target"`
		Scope   []string `json:"scope,omitempty"`
		Plan    string   `json:"plan,omitempty"`
		Rounds  int      `json:"rounds"`
		Partial bool     `json:"partial"`
		Why     string   `json:"partial_reason,omitempty"`
	}{r.Target, r.Scope, r.Plan, r.Rounds, r.Partial, r.PartialWhy}); err != nil {
		return err
	}
	for _, n := range r.Nodes {
		if err := emit("node", n); err != nil {
			return err
		}
	}
	for _, e := range r.Edges {
		if err := emit("edge", e); err != nil {
			return err
		}
	}
	for _, f := range r.Findings {
		if err := emit("finding", f); err != nil {
			return err
		}
	}
	for _, d := range r.Ledger {
		if err := emit("declined", d); err != nil {
			return err
		}
	}
	for _, t := range r.Truncations {
		if err := emit("truncation", t); err != nil {
			return err
		}
	}
	return emit("totals", r.Totals)
}
