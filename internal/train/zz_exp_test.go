// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

//go:build exp

package train

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	scanpkg "github.com/rangertaha/urlinsane/internal/scan"
	"github.com/rangertaha/urlinsane/internal/store"
)

// Leave-one-out over the saved scans in the local store.
//
//	go test -tags exp ./internal/train/ -run Experiment -v
//
// Build-tagged because it reads whatever scans that machine happens to have, so
// it is not a test — it cannot fail, and its output is a measurement rather than
// an assertion. It lives here anyway because the alternative is what it replaced:
// ad-hoc code written per experiment, unreviewed, reporting a number nobody can
// reproduce. The first version of that reported a held-out AUC of 0.802 off
// graphs whose labels Finalize had corrupted (TestFinalizeKeepsTheObserverSet);
// the corrected figures are below chance.
//
// It rotates the held-out scan rather than splitting once, because with a
// handful of scans a single split is a single sample, and the interesting
// question is whether the answer depends on which one was held out.
func TestExperiment(t *testing.T) {
	dir := filepath.Join(os.Getenv("HOME"), ".config", "urlinsane")
	ix, err := store.OpenIndex(dir)
	if err != nil {
		t.Skipf("no scan index in %s: %v", dir, err)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Skipf("no store in %s: %v", dir, err)
	}
	reg, err := scanpkg.Registry()
	if err != nil {
		t.Fatal(err)
	}

	targets := ix.Targets()
	sort.Slice(targets, func(i, j int) bool { return targets[i].Key < targets[j].Key })

	var scans []Scan
	var keys []string
	for _, e := range targets {
		root, err := store.ParseRoot(e.Root)
		if err != nil {
			continue
		}
		r, err := st.Rehydrate(root, reg)
		if err != nil {
			t.Logf("skip %s: %v", e.Key, err)
			continue
		}
		if err := Finalize(r.Graph); err != nil {
			t.Fatalf("finalize %s: %v", e.Key, err)
		}
		scans = append(scans, Scan{Graph: r.Graph, Seed: r.Seed})
		keys = append(keys, e.Key)
	}
	if len(scans) < 2 {
		t.Skipf("only %d saved scans; need at least two to hold one out", len(scans))
	}

	// Label census over everything, so the corrupted-labels question is
	// answerable from the output rather than by trusting the fix.
	live, absent, unknown, untried := 0, 0, 0, 0
	for _, s := range scans {
		a := s.Graph.Analyze()
		for _, n := range a.Nodes() {
			switch Outcome(a, n.ID) {
			case "live":
				live++
			case "absent":
				absent++
			case "unknown":
				unknown++
			default:
				untried++
			}
		}
	}
	t.Logf("%d scans, labels: live=%d absent=%d unknown=%d untried=%d",
		len(scans), live, absent, unknown, untried)

	// Leave-one-out over every scan. One split on three scans is one number
	// from one sample; rotating at least says whether the answer depends on
	// which scan was held out.
	for h := range scans {
		var fit []Scan
		for i, s := range scans {
			if i != h {
				fit = append(fit, s)
			}
		}
		res, corpus, err := Fit(DefaultConfig(), fit...)
		if err != nil {
			t.Fatalf("fit without %s: %v", keys[h], err)
		}
		held := scans[h]

		base := Evaluate(held.Graph)
		held.Graph.SetBeliefModel(BeliefFrom(res.Model))
		if err := Finalize(held.Graph); err != nil {
			t.Fatalf("rescore: %v", err)
		}
		fitted := Evaluate(held.Graph)

		fmt.Printf("hold out %-14s train=%s\n", keys[h], Describe(corpus, res))
		fmt.Printf("  base   %s\n  fitted %s\n", base, fitted)
	}
}
