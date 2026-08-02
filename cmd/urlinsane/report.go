// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
	"github.com/rangertaha/urlinsane/internal/plugins/report"
	reportall "github.com/rangertaha/urlinsane/internal/plugins/report/all"
	"github.com/rangertaha/urlinsane/internal/scan"
	"github.com/rangertaha/urlinsane/internal/store"
	"github.com/urfave/cli/v2"
)

// ReportCmd renders a scan that already happened.
//
// It is a verb rather than a flag on typo because it does not scan and does not
// describe a scan about to run — it renders one that already did. It takes a
// target, as typo does, but rejects everything that shapes a scan: --depth and
// --algorithm cannot change what a finished graph contains, and a flag on typo
// whose presence invalidated the rest of typo's surface would be worse than a
// verb (DESIGN §12).
var ReportCmd = cli.Command{
	Name:      "report",
	Usage:     "Render a saved scan",
	UsageText: "report <target> [flags]",
	Description: `Renders a scan saved earlier with 'typo --save-graph'.

You name the target, not the file: acme.com is what you scanned and what you
remember, a root CID is not. Rendering never re-observes — the stored blocks are
replayed through the same applier the scan used and CID-checked against what was
stored, so what you see is byte-identical to what was saved.`,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "filter", Aliases: []string{"f"},
			Usage: "select rows: live, absent, unknown, untried, risk>SEV, type=NAME, depth<=N"},
		&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Value: "table",
			Usage: "table | json | ndjson | csv | dot"},
		&cli.StringFlag{Name: "save", Usage: "write the report to `PATH`"},
		&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"},
			Usage: "include provenance and engine belief"},
		&cli.StringFlag{Name: "at", Usage: "render this root `CID` rather than the newest scan"},
		&cli.BoolFlag{Name: "scans", Usage: "list every saved scan of this target"},
		&cli.BoolFlag{Name: "no-color", Usage: "disable ANSI styling"},
	},
	Action: runReport,
}

func runReport(c *cli.Context) error {
	setup, err := config.Init()
	if err != nil {
		return exit(err, exitError)
	}

	ix, err := store.OpenIndex(setup.Dir)
	if err != nil {
		return exit(err, exitError)
	}

	// No target: say what there is rather than an empty usage error.
	if c.NArg() == 0 && c.String("at") == "" {
		return listScans(c, ix.Targets(), "saved scans")
	}

	entry, err := resolve(c, ix)
	if err != nil {
		return exit(err, exitError)
	}
	if c.Bool("scans") {
		return listScans(c, ix.Scans(entry.Type, entry.Key), "scans of "+entry.Key)
	}

	root, err := store.ParseRoot(entry.Root)
	if err != nil {
		return exit(err, exitError)
	}
	st, err := store.Open(setup.Dir)
	if err != nil {
		return exit(err, exitError)
	}
	reg, err := scan.Registry()
	if err != nil {
		return exit(err, exitError)
	}
	re, err := st.Rehydrate(root, reg)
	if err != nil {
		return exit(fmt.Errorf("report: %w", err), exitError)
	}

	filters, err := report.ParseFilters(c.StringSlice("filter"))
	if err != nil {
		return exit(err, exitError)
	}
	o := report.Options{
		Target:  entry.Key,
		Filters: filters,
		Verbose: c.Bool("verbose"),
		Color:   useColor(c),
		// Partial comes from the index, not from a flag: a scan interrupted at
		// a round barrier is partial however often it is re-rendered.
		Partial:    entry.Partial,
		PartialWhy: partialWhy(entry),
		Format:     c.String("output"),
	}
	if err := reportall.Write(c.App.Writer, report.Build(re.Graph, o), o); err != nil {
		return exit(err, exitError)
	}
	if path := c.String("save"); path != "" {
		f, err := os.Create(path)
		if err != nil {
			return exit(err, exitError)
		}
		defer f.Close()
		format, ok := report.FormatFor(path)
		if !ok {
			return exit(fmt.Errorf("report: cannot infer a format from %q", path), exitError)
		}
		so := o
		so.Format, so.Color = format, false
		if err := reportall.Write(f, report.Build(re.Graph, so), so); err != nil {
			return exit(err, exitError)
		}
	}
	return nil
}

func partialWhy(e store.Entry) string {
	if e.Partial {
		return "the scan was interrupted"
	}
	return ""
}

// resolve turns the arguments into one index entry: --at names a root exactly,
// a target names its most recent scan.
func resolve(c *cli.Context, ix *store.Index) (store.Entry, error) {
	if at := c.String("at"); at != "" {
		e, ok := ix.At(at)
		if !ok {
			return e, fmt.Errorf("report: no saved scan with root %s", at)
		}
		return e, nil
	}
	target := c.Args().First()
	if c.NArg() > 1 {
		return store.Entry{}, fmt.Errorf("report takes one target, got %d arguments", c.NArg())
	}
	typ, key, err := decompose.DetectSeed(target)
	if err != nil {
		return store.Entry{}, err
	}
	e, ok := ix.Latest(typ, key)
	if !ok {
		return e, fmt.Errorf(
			"report: no saved scan of %s; run: urlinsane typo --save-graph %s", key, target)
	}
	return e, nil
}

func listScans(c *cli.Context, entries []store.Entry, what string) error {
	if len(entries) == 0 {
		fmt.Fprintf(c.App.Writer, "no %s; run: urlinsane typo --save-graph <target>\n", what)
		return nil
	}
	t := tabwriter.NewWriter(c.App.Writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(t, "WHEN\tTYPE\tTARGET\tROOT\t")
	for _, e := range entries {
		mark := ""
		if e.Partial {
			mark = "partial"
		}
		fmt.Fprintf(t, "%s\t%s\t%s\t%s\t%s\n",
			e.At.Format("2006-01-02 15:04"), e.Type, e.Key, e.Root, mark)
	}
	return t.Flush()
}
