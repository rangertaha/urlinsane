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

package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// renderJSON writes the whole report as one object. Field order follows the
// struct, and every slice is already canonically sorted, so two identical scans
// produce byte-identical bytes.
func renderJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(r)
}

// renderNDJSON writes one object per line.
//
// This is the event stream (§11), not a report: it exists to be piped into
// something else while a scan runs. It carries a `kind` discriminator on every
// line because a consumer reading line by line has no enclosing object to tell
// it what it just got. The ordering here follows the graph, but the format
// makes no ordering claim and a consumer must not rely on one.
func renderNDJSON(w io.Writer, r Report) error {
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

// renderCSV writes the node table plus the ledger.
//
// CSV is flat, so edges have no place in it — but dropping the ledger to keep
// the shape tidy would let a truncated scan export as complete, which §11
// forbids. Declined candidates therefore appear as rows with an empty existence
// and a reason, and the `declined` column is what distinguishes them.
func renderCSV(w io.Writer, r Report) error {
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

	byNode := map[string][]string{}
	for _, f := range r.Findings {
		for _, n := range f.Nodes {
			byNode[n] = append(byNode[n], f.Severity+":"+f.Kind)
		}
	}

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

// renderDOT writes a graphviz digraph.
//
// This is the format the whole data model was chosen for: shared infrastructure
// is visible as convergence — several variants pointing at one IP — which is
// exactly the campaign shape §9 clusters on, and which no row-per-domain
// rendering can show at all.
func renderDOT(w io.Writer, r Report) error {
	fmt.Fprintln(w, "digraph urlinsane {")
	fmt.Fprintln(w, `  rankdir=LR;`)
	fmt.Fprintln(w, `  node [shape=box style=rounded fontname="sans-serif"];`)
	fmt.Fprintln(w, `  edge [fontname="sans-serif" fontsize=9];`)
	if r.Partial {
		fmt.Fprintf(w, "  labelloc=t; label=%s;\n", quote("PARTIAL: "+r.PartialWhy))
	}
	fmt.Fprintln(w)

	for _, n := range r.Nodes {
		id := n.Type + ":" + n.Key
		// A real newline: quote escapes it to graphviz's \n. Writing the
		// escape here instead would emit a literal backslash-n on the label.
		label := n.Key
		if n.Type != "" {
			label = n.Key + "\n" + n.Type
		}
		attrs := []string{"label=" + quote(label)}
		if fill, ok := fillFor(n); ok {
			attrs = append(attrs, "style=\"rounded,filled\"", "fillcolor="+quote(fill))
		}
		if n.Seed {
			attrs = append(attrs, "penwidth=2")
		}
		fmt.Fprintf(w, "  %s [%s];\n", quote(id), strings.Join(attrs, " "))
	}
	fmt.Fprintln(w)

	for _, e := range r.Edges {
		attrs := []string{"label=" + quote(edgeLabel(e))}
		switch e.Class {
		case "variant":
			attrs = append(attrs, "color=\"#888888\"")
		case "structural":
			attrs = append(attrs, "style=dashed", "color=\"#bbbbbb\"")
		}
		fmt.Fprintf(w, "  %s -> %s [%s];\n", quote(e.From), quote(e.To), strings.Join(attrs, " "))
	}

	// The ledger is drawn, not merely counted: a diagram that silently omits
	// what was pruned reads as a complete map of the territory.
	if len(r.Ledger) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  subgraph cluster_declined {")
		fmt.Fprintf(w, "    label=%s; style=dashed; color=\"#cccccc\";\n",
			quote(fmt.Sprintf("declined (%d)", len(r.Ledger))))
		byReason := map[string]int{}
		for _, d := range r.Ledger {
			byReason[d.Reason]++
		}
		for _, reason := range sortedKeys(byReason) {
			fmt.Fprintf(w, "    %s [label=%s shape=note fontcolor=\"#888888\" color=\"#cccccc\"];\n",
				quote("declined:"+reason),
				quote(fmt.Sprintf("%d %s", byReason[reason], reason)))
		}
		fmt.Fprintln(w, "  }")
	}

	fmt.Fprintln(w, "}")
	return nil
}

func edgeLabel(e EdgeRow) string {
	label := e.Rel
	for _, p := range e.Props {
		label += "\n" + p.Name + "=" + p.Text()
	}
	return label
}

func fillFor(n NodeRow) (string, bool) {
	switch n.Risk {
	case "critical":
		return "#ffcdd2", true
	case "high":
		return "#ffe0b2", true
	case "medium":
		return "#fff9c4", true
	}
	switch n.Existence {
	case "live":
		return "#e8f5e9", true
	case "absent":
		return "#f5f5f5", true
	}
	return "", false
}

func sortedKeys(m map[string]int) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// quote renders a graphviz string literal. Keys can contain quotes and
// backslashes — a variant algorithm will happily generate them — and an
// unescaped one produces a file graphviz refuses to parse.
func quote(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"', '\\':
			b.WriteByte('\\')
			b.WriteRune(r)
		case '\n':
			b.WriteString("\\n")
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}
