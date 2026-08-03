// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package csv renders the node rows as a spreadsheet.
//
// Leading =, +, - and @ are defused: a variant named "=cmd|' /c calc'!A1" is a
// formula-injection payload in every spreadsheet program, and a squatting tool
// that hands one to an analyst has done the attacker's delivery for them.
package csv

import (
	"encoding/csv"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/rangertaha/urlinsane/internal/plugins/report"
)

// renderCSV writes the node table plus the ledger.
//
// CSV is flat, so edges have no place in it — but dropping the ledger to keep
// the shape tidy would let a truncated scan export as complete, which §11
// forbids. report.Declined candidates therefore appear as rows with an empty existence
// and a reason, and the `declined` column is what distinguishes them.
func Render(w io.Writer, r report.Report) error {
	c := csv.NewWriter(w)
	write := func(row []string) error {
		for i := range row {
			row[i] = defuse(row[i])
		}
		return c.Write(row)
	}
	if err := write([]string{
		"type", "key", "depth", "existence", "risk", "in_closure", "props", "findings", "declined",
	}); err != nil {
		return err
	}

	// A run that stopped early says so in every other format. Without a row for
	// it, a truncated export is byte-identical to a complete one — the failure
	// §11 forbids, and the one an analyst is least able to detect from a
	// spreadsheet.
	if r.Partial {
		if err := write([]string{
			"run", r.Target, "", "partial", "", "", "", r.PartialWhy, "",
		}); err != nil {
			return err
		}
	}

	byNode := map[string][]string{}
	for _, f := range r.Findings {
		for _, n := range f.Nodes {
			byNode[n] = append(byNode[n], f.Severity+":"+f.Kind)
		}
	}

	// Findings are never filtered (report.Build, §11): hiding one behind
	// --filter would hide exactly what --fail-on gates on. Nodes are filtered,
	// so a finding reachable only through a node the filter removed has no row
	// to hang its column on and used to vanish — silently, and only from this
	// format. Those get a row of their own, distinguished the way declined
	// candidates already are: by an existence no node row carries.
	shown := make(map[string]bool, len(r.Nodes))
	for _, n := range r.Nodes {
		shown[n.Type+":"+n.Key] = true
	}
	var orphaned []string
	for _, f := range r.Findings {
		for _, n := range f.Nodes {
			if !shown[n] {
				orphaned = append(orphaned, n+"\x00"+f.Severity+":"+f.Kind)
			}
		}
	}
	sort.Strings(orphaned)

	for _, n := range r.Nodes {
		var props []string
		for _, p := range n.Props {
			props = append(props, p.Name+"="+p.Text())
		}
		if err := write([]string{
			n.Type, n.Key, strconv.Itoa(n.Depth), n.Existence, n.Risk,
			strconv.FormatBool(n.InClosure),
			strings.Join(props, " "),
			strings.Join(byNode[n.Type+":"+n.Key], " "),
			"",
		}); err != nil {
			return err
		}
	}
	for _, o := range orphaned {
		ref, finding, _ := strings.Cut(o, "\x00")
		typ, key, _ := strings.Cut(ref, ":")
		if err := write([]string{
			typ, key, "", "filtered", "", "", "", finding, "",
		}); err != nil {
			return err
		}
	}
	for _, d := range r.Ledger {
		if err := write([]string{
			d.Type, d.Key, strconv.Itoa(d.Depth), "", "", "", "", "", d.Reason,
		}); err != nil {
			return err
		}
	}
	c.Flush()
	return c.Error()
}

// defuse neutralises spreadsheet formula injection.
//
// This matters more here than in most tools: urlinsane exists to collect
// attacker-chosen strings, and a username or package name may legitimately
// begin with "=", "+", "-" or "@". Excel and LibreOffice treat such a cell as a
// formula, so `=cmd|' /c calc'!a0` as a variant name becomes code execution on
// the analyst's machine when they open the report — an attacker-controlled
// payload delivered by the security tool examining it.
//
// A leading apostrophe is the portable fix: spreadsheets read it as "this cell
// is text" and strip it, while csv, jq and every other consumer see one extra
// character rather than a mangled value. Quoting alone does not help — the
// formula still evaluates.
func defuse(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + s
	}
	return s
}
