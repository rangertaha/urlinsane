// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package all dispatches a report to the renderer for its format.
//
// Separate from the report library so each renderer can import that library —
// for Report and its row types — without the library importing them back.
package all

import (
	"fmt"
	"io"
	"strings"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/report"
	"github.com/rangertaha/urlinsane/internal/plugins/report/csv"
	"github.com/rangertaha/urlinsane/internal/plugins/report/dot"
	"github.com/rangertaha/urlinsane/internal/plugins/report/json"
	"github.com/rangertaha/urlinsane/internal/plugins/report/ndjson"
	"github.com/rangertaha/urlinsane/internal/plugins/report/table"
)

// Render writes a report in the requested format.
func Render(w io.Writer, g *graph.Graph, o report.Options) error {
	return Write(w, report.Build(g, o), o)
}

// Write renders an already-built report. Splitting this from Render lets a
// caller render one scan to several sinks — `-o table` on stdout and
// `--save out.json` — without rebuilding, and guarantees the two sinks describe
// the same scan.
func Write(w io.Writer, r report.Report, o report.Options) error {
	switch o.Format {
	case "", "table":
		return table.Render(w, r, o)
	case "json":
		return json.Render(w, r)
	case "ndjson":
		return ndjson.Render(w, r)
	case "csv":
		return csv.Render(w, r)
	case "dot":
		return dot.Render(w, r)
	}
	return fmt.Errorf("report: unknown format %q (want one of %s)",
		o.Format, strings.Join(report.Formats(), ", "))
}
