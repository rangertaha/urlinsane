// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/graph"
	_ "github.com/rangertaha/urlinsane/internal/plugins/all"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/report"
	reportall "github.com/rangertaha/urlinsane/internal/plugins/report/all"
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	variantall "github.com/rangertaha/urlinsane/internal/plugins/variant/all"
	"github.com/rangertaha/urlinsane/internal/scan"
	"github.com/rangertaha/urlinsane/internal/store"
	"github.com/urfave/cli/v2"
)

// Exit codes (§12.4). A CI gate has to react to results without parsing
// stdout, which is the only thing the old command left available.
const (
	exitError   = 1
	exitFinding = 2
)

// TypoCmd is the scanning verb.
//
// The scope positional is optional and narrowing (§12): absent, every Nameable
// node in the seed closure is varied — including the composite seed itself, so
// a bare email target also yields whole-address variants. Supplied, it filters
// that set. The target parses identically either way; scope never changes how
// the string is read, only what gets varied.
var TypoCmd = cli.Command{
	Name:                   "typo",
	UsageText:              "urlinsane typo [<scope>] <target> [options]",
	Aliases:                []string{"t"},
	Usage:                  "Scan a target for typosquatting and confusable names",
	ArgsUsage:              "[<scope>] <target>",
	UseShortOptionHandling: true,
	Description: `Expands a target into a graph of the entities it is made of, generates
variants of each, and observes what exists.

The target's kind is detected from the string alone:

    urlinsane typo example.com                 a domain
    urlinsane typo bob@example.com             an email address
    urlinsane typo npm:lodash                  a package on a named registry
    urlinsane typo github.com/acme/tool        a repository

The optional scope positional narrows what gets varied, and never changes how
the target is read:

    urlinsane typo username bob@example.com    vary only bob
    urlinsane typo domain bob@example.com      vary only example.com
    urlinsane typo username,domain bob@example.com`,
	Flags:  typoFlags(),
	Action: runTypo,
}

func typoFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{Name: "depth", Aliases: []string{"d"}, Value: 3,
			Usage: "observation hops from the seed"},
		&cli.StringSliceFlag{Name: "filter", Aliases: []string{"f"},
			Usage: "`live`, absent, unknown, risk>SEV, type=NAME, depth<=N"},
		&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Value: "table",
			Usage: "`table` | json | ndjson | csv | dot"},
		&cli.StringFlag{Name: "save",
			Usage: "write the report to `PATH`; format from the extension"},
		&cli.StringFlag{Name: "fail-on",
			Usage: "exit 2 if any finding reaches `SEV` (info|low|medium|high|critical)"},
		&cli.BoolFlag{Name: "save-graph",
			Usage: "persist the graph to the store; prints its root CID"},
		&cli.BoolFlag{Name: "explain",
			Usage: "print the compiled plan and exit"},
		&cli.StringFlag{Name: "list",
			Usage: "print a `TOPIC` and exit: operators, types, relations, algorithms, languages, keyboards, formats, filters"},
		&cli.StringSliceFlag{Name: "algorithm", Aliases: []string{"a"},
			Usage: "restrict variant generation to these algorithm `ID`s"},
		&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"},
			Usage: "include provenance and engine belief"},

		// §12.1's second tier. Registered always, hidden from the common help:
		// most runs never touch these, and a help screen that is a union of
		// every knob is one nobody reads.
		&cli.IntFlag{Name: "rounds", Hidden: true,
			Usage: "backstop for a type flow that never converges"},
		&cli.IntFlag{Name: "workers", Hidden: true,
			Usage: "concurrent operator calls"},
		&cli.IntFlag{Name: "budget", Hidden: true,
			Usage: "global admitted-node cap; 0 means unbounded"},
		&cli.IntFlag{Name: "frontier", Hidden: true,
			Usage: "cap on candidates admitted per round"},
		&cli.IntFlag{Name: "attempts", Hidden: true,
			Usage: "per-pair attempts within a round"},
		&cli.DurationFlag{Name: "timeout", Hidden: true,
			Usage: "bound on a single operator call"},
		&cli.BoolFlag{Name: "no-color", Hidden: true,
			Usage: "disable ANSI styling"},
	}
}

func runTypo(c *cli.Context) error {
	// Setup runs before --list as well as before a scan: languages come from
	// the dataset database, so `--list languages` on a fresh install must
	// extract and open it rather than report an empty set (§9.1).
	setup, err := config.Init()
	if err != nil {
		return exit(err, exitError)
	}
	reportSetup(c.App.ErrWriter, setup)

	if topic := c.String("list"); topic != "" {
		return listTopic(c.App.Writer, topic)
	}

	scope, target, err := parseArgs(c.Args().Slice())
	if err != nil {
		return exit(err, exitError)
	}

	if err := checkLimitFlags(c); err != nil {
		return exit(err, exitError)
	}

	opts := scan.Options{
		Target:     target,
		Scope:      scope,
		Algorithms: c.StringSlice("algorithm"),
		Limits: graph.Limits{
			MaxDepth:   c.Int("depth"),
			MaxRounds:  c.Int("rounds"),
			Workers:    c.Int("workers"),
			NodeBudget: c.Int("budget"),
			Frontier:   c.Int("frontier"),
			Attempts:   c.Int("attempts"),
			OpTimeout:  c.Duration("timeout"),
		},
	}
	opts.Observe.Timeout = c.Duration("timeout")

	// A missing geolocation database costs the geo operator, not the scan: with
	// a nil locator it is left out of the plan rather than planned and failed
	// (§4). Same for settings — an unreadable config file must not decide
	// whether a scan happens.
	if setup.GeoIP.Err == nil {
		if db, err := observe.OpenGeoIP(setup.Dir); err == nil {
			opts.Observe.Geo = db
		} else {
			fmt.Fprintf(c.App.ErrWriter, "  geolocation unavailable: %v\n", err)
		}
	}
	if f, err := config.Load(setup.Dir); err == nil {
		opts.Settings = f
	} else {
		fmt.Fprintf(c.App.ErrWriter, "  plugin settings unavailable: %v\n", err)
	}

	if c.Bool("explain") {
		_, _, p, err := scan.Plan(opts)
		if err != nil {
			return exit(err, exitError)
		}
		return p.Explain(c.App.Writer)
	}

	filters, err := report.ParseFilters(c.StringSlice("filter"))
	if err != nil {
		return exit(err, exitError)
	}
	if reg, err := scan.Registry(); err == nil {
		var names []string
		for _, t := range reg.Types() {
			names = append(names, t.Name())
		}
		if err := report.ValidateTypes(filters, names); err != nil {
			return exit(err, exitError)
		}
	}
	format := c.String("output")
	if !validFormat(format) {
		return exit(fmt.Errorf("unknown output format %q; want one of %s",
			format, strings.Join(report.Formats(), ", ")), exitError)
	}

	var failOn graph.Severity
	if s := c.String("fail-on"); s != "" {
		v, ok := graph.ParseSeverity(s)
		if !ok {
			return exit(fmt.Errorf(
				"unknown severity %q; want info, low, medium, high or critical", s), exitError)
		}
		failOn = v
	}

	// Ctrl-C cancels the context, which stops expansion at the end of the
	// current round: the barrier still runs, so parents, belief and the
	// truncation ledger are finalized rather than left half-computed, and the
	// report comes out marked partial. A second Ctrl-C falls through to the
	// default handler and aborts immediately — for when a round is stuck behind
	// a slow resource and waiting for the boundary is not worth it (§12.4).
	ctx, stop := interrupts(c.App.ErrWriter)
	defer stop()

	res, err := scan.Run(ctx, opts, report.Options{
		Filters: filters,
		Verbose: c.Bool("verbose"),
		Color:   useColor(c),
	})
	if err != nil {
		return exit(err, exitError)
	}

	// Saving before rendering: a graph worth keeping should survive a renderer
	// that errors on a broken pipe.
	if c.Bool("save-graph") {
		root, err := saveGraph(setup.Dir, res)
		if err != nil {
			return exit(err, exitError)
		}
		fmt.Fprintf(c.App.ErrWriter, "  saved %s\n\n", root)
	}

	if err := write(c.App.Writer, res.Report, format, c); err != nil {
		return exit(err, exitError)
	}
	if path := c.String("save"); path != "" {
		if err := save(path, res.Report, c); err != nil {
			return exit(err, exitError)
		}
	}

	if failOn != 0 && res.Report.Max() >= failOn {
		return cli.Exit(fmt.Sprintf("findings at or above %s", failOn), exitFinding)
	}
	return nil
}

// parseArgs implements §12's positional rule: with two positionals the first is
// the scope, with one it is the target. There is no --scope flag; the
// positional is the only spelling.
//
// Scope names are validated here rather than passed through, because the
// two-positional form is otherwise indistinguishable from a mistyped single
// one: `urlinsane typo exmaple.com example.com` would silently scan the second
// while treating the first as a scope that matches nothing.
func parseArgs(args []string) (scope []string, target string, err error) {
	switch len(args) {
	case 0:
		return nil, "", fmt.Errorf("no target given\n\nusage: urlinsane typo [<scope>] <target>")
	case 1:
		return nil, args[0], nil
	case 2:
		scope = splitScope(args[0])
		if err := knownScope(scope); err != nil {
			return nil, "", err
		}
		return scope, args[1], nil
	default:
		return nil, "", fmt.Errorf(
			"expected at most a scope and a target, got %d arguments: %s",
			len(args), strings.Join(args, " "))
	}
}

func splitScope(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func knownScope(scope []string) error {
	reg, err := scan.Registry()
	if err != nil {
		return err
	}
	for _, s := range scope {
		t, ok := reg.Type(s)
		if !ok {
			return fmt.Errorf("%q is not a node type; see --list types", s)
		}
		if t.Cap() != graph.Nameable {
			return fmt.Errorf(
				"%q is an observed type and cannot be varied; scope must name a nameable type", s)
		}
	}
	return nil
}

func write(w io.Writer, r report.Report, format string, c *cli.Context) error {
	return reportall.Write(w, r, report.Options{
		Format:  format,
		Verbose: c.Bool("verbose"),
		Color:   useColor(c),
	})
}

// save writes to a second sink. The format follows the extension, not -o, so
// `--save out.json` writes a report whatever stdout is doing (§11).
func save(path string, r report.Report, c *cli.Context) error {
	format, ok := report.FormatFor(path)
	if !ok {
		return fmt.Errorf(
			"cannot tell what format %q should be; use an extension: %s",
			path, strings.Join(report.Formats(), ", "))
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	// Never colored: a file is not a terminal.
	return reportall.Write(f, r, report.Options{Format: format, Verbose: c.Bool("verbose")})
}

// useColor honours --no-color, NO_COLOR and a non-TTY stdout, so `-o json` stays
// pipeable.
func useColor(c *cli.Context) bool {
	if c.Bool("no-color") {
		return false
	}
	if _, set := os.LookupEnv("NO_COLOR"); set {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func validFormat(f string) bool {
	for _, v := range report.Formats() {
		if v == f {
			return true
		}
	}
	return false
}

// exit wraps an error with a code, so a shell can tell an execution failure
// from a finding.
func exit(err error, code int) error { return cli.Exit(err.Error(), code) }

// listTopic replaces --options/--ids/--opts, which were three aliases for one
// hidden flag.
func listTopic(w io.Writer, topic string) error {
	reg, err := scan.Registry()
	if err != nil {
		return err
	}
	t := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	defer t.Flush()

	switch strings.ToLower(topic) {
	case "types":
		fmt.Fprintln(t, "NAME\tCAPABILITY\tVERSION")
		for _, d := range reg.Types() {
			fmt.Fprintf(t, "%s\t%s\t%d\n", d.Name(), d.Cap(), d.Version())
		}
	case "relations", "rels":
		fmt.Fprintln(t, "NAME\tCLASS\tDEPTH COST")
		for _, d := range reg.Rels() {
			fmt.Fprintf(t, "%s\t%s\t%d\n", d.Name(), d.Class(), d.Class().DepthCost())
		}
	case "operators", "ops":
		ops, err := scan.Operators(scan.Options{})
		if err != nil {
			return err
		}
		sort.Slice(ops, func(i, j int) bool { return ops[i].Id() < ops[j].Id() })
		fmt.Fprintln(t, "ID\tVERSION\tRESOURCE\tBINDS ON\tEMITS")
		for _, o := range ops {
			fmt.Fprintf(t, "%s\t%d\t%s\t%s\t%s\n",
				o.Id(), o.Version(), o.Resource(),
				strings.Join(o.Trigger().On.Types, ","),
				strings.Join(o.Emits().Rels, ","))
		}
	case "algorithms", "algos":
		fmt.Fprintln(t, "ID\tTITLE\tAPPLIES TO")
		for _, s := range variantall.Specs(variant.Options{}) {
			fmt.Fprintf(t, "%s\t%s\t%s\n", s.ID, s.Title, strings.Join(s.Types, ","))
		}
	case "languages", "langs":
		fmt.Fprintln(t, "ID\tNAME")
		for _, l := range variant.RegisteredLanguages() {
			fmt.Fprintf(t, "%s\t%s\n", l.Code(), l.Name())
		}
	case "keyboards":
		fmt.Fprintln(t, "ID\tNAME")
		for _, k := range variant.RegisteredKeyboards() {
			fmt.Fprintf(t, "%s\t%s\n", k.ID, k.Name)
		}
	case "formats":
		for _, f := range report.Formats() {
			fmt.Fprintln(t, f)
		}
	case "filters":
		// Documented here because the vocabulary is small and the old
		// --registered/--unregistered pair could not express the third case.
		fmt.Fprintln(t, "FILTER\tSELECTS")
		for _, row := range [][2]string{
			{"live", "an observation operator returned ok"},
			{"absent", "none did, and at least one determined absence"},
			{"unknown", "every attempt failed, timed out or was skipped"},
			{"untried", "no operator ran on it"},
			{"risk>SEV", "nodes with findings above a severity"},
			{"type=NAME", "a node type"},
			{"depth<=N", "observation hops from the seed"},
		} {
			fmt.Fprintf(t, "%s\t%s\n", row[0], row[1])
		}
	default:
		return fmt.Errorf(
			"unknown topic %q; want operators, types, relations, algorithms, languages, keyboards, formats or filters",
			topic)
	}
	return nil
}

// reportSetup says what first-run setup did. Silence would be worse than noise:
// an extraction that failed leaves operators out of the plan, and a scan with
// no geolocation must not look like a target with no geolocation (§12.6).
func reportSetup(w io.Writer, s config.Setup) {
	if !s.FirstRun() {
		return
	}
	if s.Created {
		fmt.Fprintf(w, "  created %s\n", s.Dir)
	}
	for _, f := range []struct {
		what string
		file config.File
	}{
		{"reference data", s.Dataset},
		{"geolocation database", s.GeoIP},
	} {
		switch {
		case f.file.Err != nil:
			fmt.Fprintf(w, "  %s unavailable: %v\n", f.what, f.file.Err)
		case f.file.Written:
			fmt.Fprintf(w, "  extracted %s\n", f.file.Path)
		}
	}
	fmt.Fprintln(w)
}

// saveGraph writes the scan to the content-addressed store and records it in
// the index, returning the root CID.
//
// Two writes rather than one because they answer different questions. The store
// answers "what was this scan", keyed by content so two identical scans are one
// object. The index answers "which scans of acme.com are there, and when" —
// facts a content address cannot carry without ceasing to be one.
func saveGraph(dir string, res *scan.Result) (string, error) {
	st, err := store.Open(dir)
	if err != nil {
		return "", err
	}
	root, err := st.Save(res.Graph, store.SaveOptions{Seed: res.Seed})
	if err != nil {
		return "", err
	}
	ix, err := store.OpenIndex(dir)
	if err != nil {
		return "", err
	}
	return root.String(), ix.Add(store.Entry{
		Type:    res.SeedType,
		Key:     res.SeedKey,
		Root:    root.String(),
		At:      time.Now(),
		Partial: res.Interrupt,
		// Every run fact the renderer reads, or `report` prints a different
		// page from the one the scan printed.
		Scope:  res.Report.Scope,
		Plan:   res.Report.Plan,
		Rounds: res.Rounds,
	})
}

// limitFlags are the numeric bounds a scan is run under. Every one of them is
// a count or a duration, so a negative is never meaningful.
var limitFlags = []string{"depth", "rounds", "workers", "budget", "frontier", "attempts"}

// checkLimitFlags rejects a negative bound before the scan starts.
//
// The scheduler clamps these too, but clamping is not the whole answer at this
// level: someone who typed `--workers -4` meant something, and silently running
// with one worker tells them nothing. Before either check existed, `--workers -4`
// panicked in make(chan) and `--attempts -1` was worse than a panic -- the
// attempt loop ran zero times, so nothing was probed and the empty report came
// back as a successful scan.
func checkLimitFlags(c *cli.Context) error {
	for _, name := range limitFlags {
		if v := c.Int(name); v < 0 {
			return fmt.Errorf("--%s must not be negative, got %d", name, v)
		}
	}
	if d := c.Duration("timeout"); d < 0 {
		return fmt.Errorf("--timeout must not be negative, got %s", d)
	}
	return nil
}
