// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package aitypo

import (
	"bytes"
	"strings"
	"testing"
)

func testData() Data {
	return Data{
		Graphemes:  []string{"a", "b"},
		Vowels:     []string{"a", "e", "i", "o", "u"},
		Homoglyphs: map[string][]string{"o": {"0", "о"}},
		Homophones: [][]string{{"to", "too", "two"}},
		Keyboard:   []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"},
	}
}

// The oracle is the label source, so it has to be the same code the scanner
// runs. A task registry that reimplemented a generator would train a model
// against a second definition of correct.
func TestOracleIsTheShippedGenerator(t *testing.T) {
	tasks := Tasks(testData())
	co, ok := tasks["co"]
	if !ok {
		t.Fatal("co is not registered")
	}
	got := co.Expect("google")
	want := []string{"gogle", "googe", "googl", "goole", "oogle"}
	if len(got) != len(want) {
		t.Fatalf("co(google) = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("co(google) = %q, want %q", got, want)
		}
	}
}

// Expect is a set: sorted, deduplicated, and never containing the input.
func TestExpectIsASet(t *testing.T) {
	noisy := Task{ID: "noisy", Oracle: func(s string) []string {
		return []string{"b", "a", "b", s, "", "a"}
	}}
	got := noisy.Expect("seed")
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("Expect = %q, want [a b] — sorted, deduped, input dropped", got)
	}
}

// Needs separates a rule from a table, which is the distinction the whole
// package exists to keep: reporting one accuracy across both averages a rule a
// model learned with a lookup it memorised.
func TestNeedsSeparatesRulesFromTables(t *testing.T) {
	tasks := Tasks(testData())
	for id, want := range map[string]Needs{
		"co":  NeedsNothing,
		"cs":  NeedsNothing,
		"tos": NeedsNothing,
		"hr":  NeedsLanguage,
		"hs":  NeedsLanguage,
		"acs": NeedsKeyboard,
		"aci": NeedsKeyboard,
	} {
		task, ok := tasks[id]
		if !ok {
			t.Errorf("%s is not registered", id)
			continue
		}
		if task.Needs != want {
			t.Errorf("%s needs %v, want %v", id, task.Needs, want)
		}
	}
}

// A missing table silences its tasks rather than registering one that always
// answers nothing — which would otherwise score as perfectly learned.
func TestMissingDataSilencesItsTasks(t *testing.T) {
	tasks := Tasks(Data{})
	for _, id := range []string{"hr", "hs", "acs", "aci", "vs", "gi"} {
		if _, ok := tasks[id]; ok {
			t.Errorf("%s registered with no data behind it", id)
		}
	}
	// The pure rules are still there: a corpus of rules must not require a
	// language database to exist.
	for _, id := range []string{"co", "cs", "hi", "tos"} {
		if _, ok := tasks[id]; !ok {
			t.Errorf("%s should not depend on data", id)
		}
	}
}

// Splits are assigned by input, so no name appears in two of them.
func TestSplitDoesNotLeakAcrossTasks(t *testing.T) {
	tasks, err := Tasks(testData()).Select("co", "cs", "cr", "hi")
	if err != nil {
		t.Fatal(err)
	}
	corpus := []string{"google", "example", "paypal", "amazon", "github", "reddit", "twitch", "signal"}
	ex := Emit(tasks, corpus, "test")
	Assign(ex, DefaultRatio, "")

	if bad := Leakage(ex); len(bad) != 0 {
		t.Errorf("inputs in more than one split: %q", bad)
	}
	for _, e := range ex {
		if e.Split == "" {
			t.Fatalf("%s/%s got no split", e.Task, e.Input)
		}
	}
}

// The same name lands in the same split on every machine and every run. A
// process-seeded hash would move it, and two experiments would then differ by
// their corpora as well as their models.
func TestSplitIsStableAcrossRuns(t *testing.T) {
	corpus := []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot"}
	first := make([]Example, len(corpus))
	for i, c := range corpus {
		first[i] = Example{Task: "co", Input: c}
	}
	Assign(first, DefaultRatio, "salt")

	for round := 0; round < 20; round++ {
		again := make([]Example, len(corpus))
		for i, c := range corpus {
			again[i] = Example{Task: "co", Input: c}
		}
		Assign(again, DefaultRatio, "salt")
		for i := range again {
			if again[i].Split != first[i].Split {
				t.Fatalf("%q moved from %s to %s on round %d",
					corpus[i], first[i].Split, again[i].Split, round)
			}
		}
	}

	// A different salt is allowed to disagree; that is what it is for.
	other := []Example{{Task: "co", Input: "alpha"}}
	Assign(other, DefaultRatio, "different")
	_ = other
}

// Adding names must not move the ones already there. A shuffle-and-slice split
// silently invalidates every earlier measurement when the corpus grows.
func TestSplitIsStableAsTheCorpusGrows(t *testing.T) {
	small := []Example{{Input: "alpha"}, {Input: "bravo"}, {Input: "charlie"}}
	Assign(small, DefaultRatio, "")
	was := map[string]string{}
	for _, e := range small {
		was[e.Input] = e.Split
	}

	big := []Example{{Input: "alpha"}, {Input: "bravo"}, {Input: "charlie"}}
	for _, extra := range []string{"delta", "echo", "foxtrot", "golf", "hotel"} {
		big = append(big, Example{Input: extra})
	}
	Assign(big, DefaultRatio, "")
	for _, e := range big {
		if prev, ok := was[e.Input]; ok && prev != e.Split {
			t.Errorf("%q moved from %s to %s when the corpus grew", e.Input, prev, e.Split)
		}
	}
}

// Exact match is the headline, and it is strict: four right and one invented
// is not the rule.
func TestScoreIsStrictAboutTheWholeSet(t *testing.T) {
	co := Tasks(testData())["co"]

	perfect := Score(co, "google", co.Expect("google"))
	if !perfect.Exact() || perfect.F1() != 1 {
		t.Errorf("a correct answer scored exact=%v f1=%v", perfect.Exact(), perfect.F1())
	}

	full := co.Expect("google")
	invented := append(append([]string{}, full...), "gooooogle")
	r := Score(co, "google", invented)
	if r.Exact() {
		t.Error("a prediction with an invented variant counted as exact")
	}
	if r.Recall() != 1 {
		t.Errorf("recall = %v, want 1: nothing was missed", r.Recall())
	}
	if r.Precision() >= 1 {
		t.Errorf("precision = %v, want < 1: one prediction was wrong", r.Precision())
	}
	if len(r.Spuri) != 1 || r.Spuri[0] != "gooooogle" {
		t.Errorf("spurious = %q, want [gooooogle]", r.Spuri)
	}
}

// A model cannot inflate its score by repeating itself or echoing the input,
// because the prediction is normalised the same way the oracle's output is.
func TestScoreNormalisesThePrediction(t *testing.T) {
	co := Tasks(testData())["co"]
	full := co.Expect("google")

	padded := []string{"google", "google"}
	for _, v := range full {
		padded = append(padded, v, v)
	}
	if r := Score(co, "google", padded); !r.Exact() {
		t.Errorf("duplicates and the input changed the verdict: missed=%q spurious=%q",
			r.Missed, r.Spuri)
	}
}

// Macro-averaged: every input is one vote, so the score does not mostly
// measure the corpus's length distribution.
func TestSummarizeIsMacroAveraged(t *testing.T) {
	results := []Result{
		// One input answered perfectly, with many variants.
		{Task: "co", Input: "aaaaaaaa", Expect: []string{"a", "b", "c", "d"},
			Predict: []string{"a", "b", "c", "d"}, Hit: []string{"a", "b", "c", "d"}},
		// One input answered wrongly, with few.
		{Task: "co", Input: "ab", Expect: []string{"x"}, Predict: []string{"y"},
			Missed: []string{"x"}, Spuri: []string{"y"}},
	}
	s := Summarize(results, map[string]Needs{"co": NeedsNothing})
	if len(s) != 1 {
		t.Fatalf("got %d summaries, want 1", len(s))
	}
	if s[0].ExactMatch != 0.5 {
		t.Errorf("exact = %v, want 0.5 — one of two inputs, not one of five variants",
			s[0].ExactMatch)
	}
	if s[0].N != 2 {
		t.Errorf("n = %d, want 2", s[0].N)
	}
}

// A corpus round-trips, including expectation sets far past the default
// scanner line limit.
func TestCorpusRoundTrips(t *testing.T) {
	tasks, err := Tasks(testData()).Select("co", "cs")
	if err != nil {
		t.Fatal(err)
	}
	ex := Emit(tasks, []string{"google", "example"}, "unit")
	Assign(ex, DefaultRatio, "")

	var buf bytes.Buffer
	if err := WriteJSONL(&buf, ex); err != nil {
		t.Fatalf("write: %v", err)
	}
	back, err := ReadJSONL(&buf)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(back) != len(ex) {
		t.Fatalf("read %d examples, wrote %d", len(back), len(ex))
	}
	for i := range ex {
		if back[i].Task != ex[i].Task || back[i].Input != ex[i].Input ||
			back[i].Split != ex[i].Split || len(back[i].Expect) != len(ex[i].Expect) {
			t.Fatalf("example %d changed: %+v vs %+v", i, back[i], ex[i])
		}
	}
}

// A malformed line is an error, not a silently dropped row: a corpus that
// drops rows trains on less than it reports.
func TestReadRefusesAMalformedLine(t *testing.T) {
	_, err := ReadJSONL(strings.NewReader("{\"task\":\"co\"}\nnot json\n"))
	if err == nil {
		t.Fatal("a malformed line was accepted")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("error = %v, want it to name the line", err)
	}
}

// An unknown task id is an error rather than a quiet omission.
func TestSelectRefusesUnknownTasks(t *testing.T) {
	_, err := Tasks(testData()).Select("co", "nope")
	if err == nil {
		t.Fatal("an unknown task id was accepted")
	}
	if !strings.Contains(err.Error(), "nope") {
		t.Errorf("error = %v, want it to name the unknown id", err)
	}
}

// Emit drops empty expectations by default and keeps them on request, so
// "sometimes the answer is empty" is a deliberate sample rather than an
// accident of which names the corpus held.
func TestEmitEmptyIsOptedInto(t *testing.T) {
	tasks, err := Tasks(testData()).Select("rar")
	if err != nil {
		t.Skip("rar needs a keyboard")
	}
	// "example" has no doubled letters, so rar produces nothing.
	if got := Emit(tasks, []string{"example"}, "u"); len(got) != 0 {
		t.Errorf("Emit kept an empty expectation: %+v", got)
	}
	got := EmitWithEmpty(tasks, []string{"example"}, "u")
	if len(got) != 1 || len(got[0].Expect) != 0 {
		t.Errorf("EmitWithEmpty = %+v, want one example with an empty set", got)
	}
}
