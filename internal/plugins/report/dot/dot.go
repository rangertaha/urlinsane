// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package dot renders the graph as Graphviz.
//
// The only format that keeps the shape. A table flattens the graph to rows and
// loses precisely what the engine exists to find — that forty variants share
// one nameserver.
package dot

import (
	"fmt"
	"io"
	"strings"

	"github.com/rangertaha/urlinsane/internal/plugins/report"
)

// renderDOT writes a graphviz digraph.
//
// This is the format the whole data model was chosen for: shared infrastructure
// is visible as convergence — several variants pointing at one IP — which is
// exactly the campaign shape §9 clusters on, and which no row-per-domain
// rendering can show at all.
func Render(w io.Writer, r report.Report) error {
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

	// Findings whose nodes the filter removed. Findings are never filtered —
	// hiding one behind --filter would hide exactly what --fail-on gates on —
	// but they are drawn as an attribute of a node, so a stranded finding had
	// nothing to attach to and vanished from this format alone.
	shown := make(map[string]bool, len(r.Nodes))
	for _, n := range r.Nodes {
		shown[n.ID] = true
	}
	stranded := map[string]int{}
	for _, f := range r.Findings {
		for _, n := range f.Nodes {
			if !shown[n] {
				stranded[f.Severity+":"+f.Kind]++
			}
		}
	}
	if len(stranded) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "  subgraph cluster_filtered_findings {")
		fmt.Fprintf(w, "    label=%s; style=dashed; color=\"#e0a0a0\";\n",
			quote(fmt.Sprintf("findings on filtered nodes (%d)", len(stranded))))
		for _, k := range sortedKeys(stranded) {
			fmt.Fprintf(w, "    %s [label=%s shape=note fontcolor=\"#aa4444\" color=\"#e0a0a0\"];\n",
				quote("finding:"+k),
				quote(fmt.Sprintf("%d %s", stranded[k], k)))
		}
		fmt.Fprintln(w, "  }")
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

func edgeLabel(e report.EdgeRow) string {
	label := e.Rel
	for _, p := range e.Props {
		label += "\n" + p.Name + "=" + p.Text()
	}
	return label
}

func fillFor(n report.NodeRow) (string, bool) {
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
