// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

// width is the column the help wraps at.
//
// 80 because that is what a terminal is until told otherwise, and help that
// wraps in an 80-column window is help nobody finishes reading. urfave/cli's
// default template puts a flag's whole specification and its description on one
// line, which reached 157 characters here — the description of --list wrapped
// three times and lost its alignment, which is exactly the flag a new user
// looks at first.
const width = 80

// helpTemplates replaces urfave/cli's layouts with ones that wrap.
//
// The flag list is rendered by flagLines rather than by the template, because
// the template only has each flag's String(), and that is the long form we are
// trying to get away from.
func helpTemplates() {
	cli.AppHelpTemplate = `{{.Name}} - {{.Usage}}

USAGE:
   {{usage .Name .UsageText}}
{{if .Description}}
{{wrap .Description 3}}
{{end}}
COMMANDS:
{{range .VisibleCommands}}   {{join .Names ", " | printf "%-10s"}} {{.Usage}}
{{end}}
OPTIONS:
{{flags .VisibleFlags}}
EXAMPLES:
   urlinsane typo example.com              scan a domain
   urlinsane typo -a hr,co example.com     only these algorithms
   urlinsane typo --save-graph acme.com    keep the graph
   urlinsane report acme.com               render it again, no network

Run 'urlinsane <command> --help' for a command's own options.
`

	cli.CommandHelpTemplate = `{{.HelpName}} - {{.Usage}}

USAGE:
   {{usage .HelpName .UsageText}}
{{if .Description}}
{{wrap .Description 3}}
{{end}}
OPTIONS:
{{flags .VisibleFlags}}`

	cli.HelpPrinter = func(w io.Writer, tmpl string, data any) {
		cli.HelpPrinterCustom(w, tmpl, data, map[string]any{
			"flags": flagLines,
			"wrap":  func(s string, indent int) string { return wrap(s, indent) },
			"usage": usageLine,
		})
	}
}

// flagLines renders flags as "  -a, --algorithm ID" against a wrapped
// description, in two aligned columns.
//
// Short name first: a user scanning for -a is looking for the thing they type,
// not its long spelling.
func flagLines(flags []cli.Flag) string {
	var b strings.Builder
	t := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)

	for _, f := range flags {
		spec := "   " + names(f)
		if d, ok := f.(cli.DocGenerationFlag); ok && d.TakesValue() {
			if v := placeholder(d); v != "" {
				spec += " " + v
			}
		}
		fmt.Fprintf(t, "%s\t%s\n", spec, usage(f))
	}
	t.Flush()

	// Wrap only after alignment, so a long description folds under its own
	// column instead of resetting to the margin.
	return fold(b.String())
}

func names(f cli.Flag) string {
	var short, long []string
	for _, n := range f.Names() {
		if len(n) == 1 {
			short = append(short, "-"+n)
		} else {
			long = append(long, "--"+n)
		}
	}
	return strings.Join(append(short, long...), ", ")
}

func placeholder(d cli.DocGenerationFlag) string {
	// urfave/cli encodes a placeholder as `backticks` in the usage string.
	u := d.GetUsage()
	if i := strings.Index(u, "`"); i >= 0 {
		if j := strings.Index(u[i+1:], "`"); j >= 0 {
			return u[i+1 : i+1+j]
		}
	}
	return "value"
}

func usage(f cli.Flag) string {
	d, ok := f.(cli.DocGenerationFlag)
	if !ok {
		return ""
	}
	u := strings.ReplaceAll(d.GetUsage(), "`", "")
	if def := d.GetDefaultText(); def != "" && def != "false" {
		u += " (default: " + def + ")"
	}
	return u
}

// fold rewraps aligned two-column text so nothing exceeds width, continuing
// long descriptions under the description column.
func fold(s string) string {
	var out strings.Builder
	for _, line := range strings.Split(strings.TrimRight(s, "\n"), "\n") {
		if len(line) <= width {
			out.WriteString(line + "\n")
			continue
		}
		col := strings.Index(strings.TrimLeft(line, " "), "  ")
		col = len(line) - len(strings.TrimLeft(line, " ")) + col
		for col < len(line) && line[col] == ' ' {
			col++
		}
		out.WriteString(wrapAt(line, col) + "\n")
	}
	return out.String()
}

// wrapAt folds a line at width, indenting continuations to col.
func wrapAt(line string, col int) string {
	head, rest := line[:col], strings.Fields(line[col:])
	var b strings.Builder
	b.WriteString(head)
	n := col
	for i, w := range rest {
		if i > 0 && n+1+len(w) > width {
			b.WriteString("\n" + strings.Repeat(" ", col))
			n = col
		} else if i > 0 {
			b.WriteString(" ")
			n++
		}
		b.WriteString(w)
		n += len(w)
	}
	return b.String()
}

// wrap folds prose to width at the given indent, leaving pre-formatted lines
// alone.
//
// A line that arrives already indented is an example, and reflowing one turns a
// column of commands into a paragraph of words — which is what happened the
// first time this existed. Indentation is the only signal available here, and
// it is the one every help text already uses.
func wrap(s string, indent int) string {
	pad := strings.Repeat(" ", indent)
	var out []string
	var para []string

	flush := func() {
		if len(para) == 0 {
			return
		}
		line, n := pad, indent
		for i, w := range strings.Fields(strings.Join(para, " ")) {
			if i > 0 && n+1+len(w) > width {
				out = append(out, line)
				line, n = pad, indent
			} else if i > 0 {
				line += " "
				n++
			}
			line += w
			n += len(w)
		}
		out = append(out, line)
		para = nil
	}

	for _, line := range strings.Split(strings.TrimRight(s, "\n"), "\n") {
		switch {
		case strings.TrimSpace(line) == "":
			flush()
			out = append(out, "")
		case line[0] == ' ' || line[0] == '\t':
			flush()
			out = append(out, pad+strings.TrimRight(line, " "))
		default:
			para = append(para, line)
		}
	}
	flush()
	return strings.TrimRight(strings.Join(out, "\n"), "\n")
}

// usageLine is what USAGE shows. cli leaves UsageText empty when a command does
// not set one, and an empty USAGE heading is worse than none.
func usageLine(name, text string) string {
	if text != "" {
		return text
	}
	return name + " [options]"
}
