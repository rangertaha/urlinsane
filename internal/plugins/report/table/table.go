// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package table renders a report for a terminal: nodes, findings, the declined
// ledger and a summary, with colour when the sink is a TTY.
//
// The default format, and the only one that decides what to leave out — a
// terminal has a width, so detail is folded rather than truncated silently.
package table

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/rangertaha/urlinsane/internal/plugins/report"
)

// ANSI codes, used only when report.Options.Color is set. The caller decides; this
// package never inspects the terminal or the environment, so a test can render
// both forms without setting globals.
const (
	reset  = "\033[0m"
	dim    = "\033[2m"
	bold   = "\033[1m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	blue   = "\033[34m"
)

func Render(w io.Writer, r report.Report, o report.Options) error {
	c := painter(o.Color)

	if r.Target != "" {
		head := r.Target
		if len(r.Scope) > 0 {
			head += " · scope " + strings.Join(r.Scope, ",")
		}
		fmt.Fprintf(w, "%s\n", c(bold, head))
	}
	if r.Partial {
		// Loud, and above the data: a reader who scrolls to the table and stops
		// must not mistake a prefix for the whole scan.
		why := r.PartialWhy
		if why == "" {
			why = "expansion stopped early"
		}
		fmt.Fprintf(w, "%s\n", c(yellow, "PARTIAL — "+why+"; results are a prefix, not the whole scan"))
	}
	fmt.Fprintln(w)

	if err := nodeTable(w, r, o, c); err != nil {
		return err
	}
	if err := findingTable(w, r, c); err != nil {
		return err
	}
	if err := ledgerTable(w, r, c); err != nil {
		return err
	}
	return summary(w, r, c)
}

func nodeTable(w io.Writer, r report.Report, o report.Options, c colorer) error {
	if len(r.Nodes) == 0 {
		fmt.Fprintf(w, "%s\n\n", c(dim, "no nodes matched"))
		return nil
	}
	t := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintf(t, "%s\t%s\t%s\t%s\t%s\t%s\n",
		c(dim, "TYPE"), c(dim, "KEY"), c(dim, "DEPTH"),
		c(dim, "EXISTENCE"), c(dim, "RISK"), c(dim, "DETAIL"))

	for _, n := range r.Nodes {
		fmt.Fprintf(t, "%s\t%s\t%d\t%s\t%s\t%s\n",
			n.Type, n.Key, n.Depth,
			c(existenceColor(n.Existence), n.Existence),
			c(severityColor(n.Risk), n.Risk),
			detail(n, o.Verbose))
	}
	if err := t.Flush(); err != nil {
		return err
	}
	fmt.Fprintln(w)
	return nil
}

// detail is the one-line summary that makes a row worth reading: the props an
// operator learned, and in verbose mode the operators that spoke.
func detail(n report.NodeRow, verbose bool) string {
	var parts []string
	for _, p := range n.Props {
		parts = append(parts, p.Name+"="+oneLine(p.Text()))
	}
	for _, s := range n.Scores {
		parts = append(parts, fmt.Sprintf("%s=%.2f", s.Key, s.Value))
	}
	if verbose {
		for _, s := range n.Statuses {
			parts = append(parts, s.Operator+":"+s.Status)
		}
		if n.Belief != nil {
			parts = append(parts, fmt.Sprintf("belief=%.3f", *n.Belief))
		}
	}
	return strings.Join(parts, " ")
}

// oneLine flattens a value for a column-aligned table.
//
// Some props genuinely hold several records — a domain's TXT set is joined into
// one string, because props are single-valued and the records have no entity to
// hang off. tabwriter measures column widths per line, so an embedded newline
// silently destroys the alignment of every row after it. The other formats keep
// the value intact; only the table has to flatten it.
func oneLine(s string) string {
	if !strings.ContainsAny(s, "\n\r\t") {
		return s
	}
	return strings.Join(strings.FieldsFunc(s, func(r rune) bool {
		return r == '\n' || r == '\r' || r == '\t'
	}), " ⏎ ")
}

func findingTable(w io.Writer, r report.Report, c colorer) error {
	if len(r.Findings) == 0 {
		return nil
	}
	fmt.Fprintf(w, "%s\n", c(bold, "FINDINGS"))
	t := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	for _, f := range r.Findings {
		fmt.Fprintf(t, "  %s\t%s\t%s\n",
			c(severityColor(f.Severity), strings.ToUpper(f.Severity)), f.Kind, f.Summary)
	}
	if err := t.Flush(); err != nil {
		return err
	}
	fmt.Fprintln(w)
	return nil
}

func ledgerTable(w io.Writer, r report.Report, c colorer) error {
	if len(r.Ledger) == 0 && len(r.Truncations) == 0 {
		return nil
	}
	fmt.Fprintf(w, "%s\n", c(bold, "DECLINED"))
	for _, t := range r.Truncations {
		fmt.Fprintf(w, "  %s\n", c(yellow, fmt.Sprintf("round %d: %s — %s", t.Round, t.Reason, t.Detail)))
	}

	// Candidates are summarised by reason rather than listed: a belief gate can
	// decline thousands, and a report that prints every one buries the findings
	// it exists to surface. --verbose is not the escape hatch — `-o json` is,
	// where the full ledger always appears.
	byReason := map[string]int{}
	var order []string
	for _, d := range r.Ledger {
		if _, seen := byReason[d.Reason]; !seen {
			order = append(order, d.Reason)
		}
		byReason[d.Reason]++
	}
	sort.Strings(order)
	for _, reason := range order {
		fmt.Fprintf(w, "  %s\t%s\n",
			c(dim, strconv.Itoa(byReason[reason])+" candidates"), reason)
	}
	fmt.Fprintln(w)
	return nil
}

func summary(w io.Writer, r report.Report, c colorer) error {
	counts := make([]string, 0, len(r.Totals.TypeOrder()))
	types := append([]string(nil), r.Totals.TypeOrder()...)
	sort.Strings(types)
	for _, t := range types {
		counts = append(counts, fmt.Sprintf("%s %d", t, r.Totals.ByType[t]))
	}

	fmt.Fprintf(w, "%s\n", c(dim, strings.Join(counts, "  ")))
	line := fmt.Sprintf("%d nodes  %d edges  %s  %s  %s",
		r.Totals.Nodes, r.Totals.Edges,
		c(green, fmt.Sprintf("%d live", r.Totals.Live)),
		fmt.Sprintf("%d absent", r.Totals.Absent),
		fmt.Sprintf("%d unknown", r.Totals.Unknown))
	if r.Totals.Declined > 0 {
		line += c(yellow, fmt.Sprintf("  %d declined", r.Totals.Declined))
	}
	if r.Totals.Shown != r.Totals.Nodes {
		line += c(dim, fmt.Sprintf("  (%d shown)", r.Totals.Shown))
	}
	fmt.Fprintln(w, line)

	if r.Elapsed > 0 {
		fmt.Fprintf(w, "%s\n", c(dim, fmt.Sprintf("%d rounds in %s", r.Rounds, r.Elapsed.Round(1e6))))
	}
	return nil
}

type colorer func(code, s string) string

func painter(on bool) colorer {
	if !on {
		return func(_, s string) string { return s }
	}
	return func(code, s string) string {
		if s == "" {
			return s
		}
		return code + s + reset
	}
}

func existenceColor(s string) string {
	switch s {
	case "live":
		return green
	case "absent":
		return dim
	case "unknown":
		return yellow
	}
	return dim
}

func severityColor(s string) string {
	switch s {
	case "critical", "high":
		return red
	case "medium":
		return yellow
	case "low":
		return blue
	}
	return dim
}
