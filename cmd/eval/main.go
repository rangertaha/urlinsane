// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Command urlinsane-eval measures the variant algorithms against names that
// were actually registered.
//
// It is a development harness, not part of a scan. Two steps:
//
//	urlinsane-eval fetch --brand paypal.com --out testdata/truth/paypal.jsonl
//	urlinsane-eval score --truth testdata/truth
//
// fetch needs the network and writes an *unreviewed* candidate set; score is
// offline and is the part that answers whether a change helped.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rangertaha/urlinsane/internal/eval"
	"github.com/rangertaha/urlinsane/internal/plugins/variant"
	variantall "github.com/rangertaha/urlinsane/internal/plugins/variant/all"
	"github.com/urfave/cli/v2"
	// The language-driven algorithms are silent without a language plugin, so
	// a score run that omitted these would under-report recall for a reason
	// that has nothing to do with the algorithms.
)

func main() {
	app := &cli.App{
		Name:  "urlinsane-eval",
		Usage: "measure variant algorithms against registered lookalikes",
		Commands: []*cli.Command{
			{
				Name:      "fetch",
				Usage:     "pull candidate lookalikes from Certificate Transparency",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "brand", Required: true,
						Usage: "the `DOMAIN` to find lookalikes of"},
					&cli.StringSliceFlag{Name: "exclude",
						Usage: "registrable `DOMAIN` known to belong to the brand; repeatable"},
					&cli.IntFlag{Name: "distance", Value: 2,
						Usage: "max edit distance from the brand's registrable name"},
					&cli.StringFlag{Name: "out", Aliases: []string{"o"},
						Usage: "write JSONL to `PATH` instead of stdout"},
					&cli.DurationFlag{Name: "timeout", Value: 120 * time.Second,
						Usage: "bound the crt.sh request"},
				},
				Action: fetch,
			},
			{
				Name:      "score",
				Usage:     "measure recall of the algorithm set against a truth set",
				ArgsUsage: " ",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "truth", Aliases: []string{"t"}, Required: true,
						Usage: "truth set `PATH`: a .jsonl file or a directory of them"},
					&cli.BoolFlag{Name: "missed", Aliases: []string{"m"},
						Usage: "list the names no algorithm reached"},
				},
				Action: score,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func fetch(c *cli.Context) error {
	ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
	defer cancel()

	records, err := eval.FetchCrtSh(ctx, nil, eval.Filter{
		Brand:       c.String("brand"),
		MaxDistance: c.Int("distance"),
		Exclude:     c.StringSlice("exclude"),
	})
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("no lookalikes found for %s; an empty truth set scores as perfect recall, so this is being reported as an error rather than written out", c.String("brand"))
	}

	out := os.Stdout
	if path := c.String("out"); path != "" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}
	if err := eval.WriteTruth(out, records); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr,
		"%d candidate lookalikes for %s -- unreviewed. Read them before trusting the recall number.\n",
		len(records), c.String("brand"))
	return nil
}

func score(c *cli.Context) error {
	records, err := loadTruth(c.String("truth"))
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("truth set is empty; there is nothing to measure")
	}

	specs := variantall.Specs(variant.Options{})
	reports := eval.ScoreAll(records, specs)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "BRAND\tTRUTH\tGENERATED\tCORE HITS\tEXACT\tRECALL\tCANDS/HIT")
	var truth, hits, generated int
	for _, r := range reports {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%.1f%%\t%.0f\n",
			r.Brand, r.Truth, r.Generated, r.CoreHits, r.ExactHits,
			r.CoreRecall()*100, r.CandidatesPerHit())
		truth += r.Truth
		hits += r.CoreHits
		generated += r.Generated
	}
	w.Flush()

	overall := 0.0
	if truth > 0 {
		overall = float64(hits) / float64(truth)
	}
	fmt.Printf("\noverall: %d/%d core recall = %.1f%% across %d brands, %d candidates\n",
		hits, truth, overall*100, len(reports), generated)

	printAlgorithms(reports)

	if c.Bool("missed") {
		fmt.Println("\nmissed:")
		for _, r := range reports {
			for _, m := range r.Missed {
				fmt.Printf("  %-28s (%s)\n", m, r.Brand)
			}
		}
	}
	return nil
}

// printAlgorithms aggregates per-algorithm contribution across brands. This is
// the actionable half of the report: an algorithm with many candidates and no
// unique hits is buying coverage another algorithm already provides.
func printAlgorithms(reports []eval.Report) {
	type agg struct{ cands, hits, unique int }
	totals := map[string]*agg{}
	for _, r := range reports {
		for _, s := range r.ByAlgorithm {
			if totals[s.ID] == nil {
				totals[s.ID] = &agg{}
			}
			t := totals[s.ID]
			t.cands += s.Candidates
			t.hits += s.Hits
			t.unique += s.Unique
		}
	}
	ids := make([]string, 0, len(totals))
	for id := range totals {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		a, b := totals[ids[i]], totals[ids[j]]
		if a.hits != b.hits {
			return a.hits > b.hits
		}
		if a.unique != b.unique {
			return a.unique > b.unique
		}
		return ids[i] < ids[j]
	})

	fmt.Printf("\nper algorithm, summed across all %d brands:\n", len(reports))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ALGORITHM\tCANDIDATES\tHITS\tUNIQUE\tCANDS/HIT")
	for _, id := range ids {
		t := totals[id]
		perHit := "-"
		if t.hits > 0 {
			perHit = fmt.Sprintf("%.0f", float64(t.cands)/float64(t.hits))
		}
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%s\n", id, t.cands, t.hits, t.unique, perHit)
	}
	w.Flush()
}

// loadTruth reads a single .jsonl file or every .jsonl in a directory.
func loadTruth(path string) ([]eval.Record, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var paths []string
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
				paths = append(paths, filepath.Join(path, e.Name()))
			}
		}
		sort.Strings(paths)
	} else {
		paths = []string{path}
	}

	var out []eval.Record
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		recs, err := eval.ReadTruth(f)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		out = append(out, recs...)
	}
	return out, nil
}
