// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package aitypo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteJSONL writes examples one per line.
//
// JSONL because a corpus is appended to, streamed, and reviewed in a diff, and
// because every training stack reads it. One example per line means a corpus
// can be shuffled with sort -R, split with head, and inspected with grep
// without a parser.
//
// Examples are written in the order given. Emit produces them in corpus order
// and task order, so a corpus built twice from the same inputs is byte
// identical — which is what lets a training run be identified by the hash of
// its data.
func WriteJSONL(w io.Writer, examples []Example) error {
	bw := bufio.NewWriter(w)
	enc := json.NewEncoder(bw)
	for _, e := range examples {
		if err := enc.Encode(e); err != nil {
			return fmt.Errorf("aitypo: writing example %s/%s: %w", e.Task, e.Input, err)
		}
	}
	return bw.Flush()
}

// ReadJSONL reads a corpus back.
//
// A malformed line is an error naming the line number rather than a skipped
// record. A corpus that silently drops rows trains on less data than it
// reports, and the report is what an experiment is judged on.
func ReadJSONL(r io.Reader) ([]Example, error) {
	var out []Example
	sc := bufio.NewScanner(r)
	// Expectation sets can be long — grapheme insertion over a 26-letter
	// alphabet is hundreds of variants — so the default 64KB line limit is not
	// enough.
	sc.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)

	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e Example
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("aitypo: line %d: %w", n, err)
		}
		out = append(out, e)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("aitypo: reading corpus: %w", err)
	}
	return out, nil
}

// Filter returns the examples in one split.
func Filter(examples []Example, split string) []Example {
	out := make([]Example, 0, len(examples))
	for _, e := range examples {
		if e.Split == split {
			out = append(out, e)
		}
	}
	return out
}

// Stats is what a corpus is made of, for the line a training run should print
// before it starts.
type Stats struct {
	Examples int
	Inputs   int
	Tasks    int
	Splits   map[string]int
	ByTask   map[string]int
	// Variants is the total size of all expectation sets — the number of
	// output strings a model has to produce across the corpus, which is the
	// figure that predicts training cost, not the example count.
	Variants int
}

func (s Stats) String() string {
	splits := make([]string, 0, len(s.Splits))
	for _, k := range []string{SplitTrain, SplitVal, SplitTest} {
		if n, ok := s.Splits[k]; ok {
			splits = append(splits, fmt.Sprintf("%s=%d", k, n))
		}
	}
	return fmt.Sprintf("%d examples, %d inputs, %d tasks, %d variants [%s]",
		s.Examples, s.Inputs, s.Tasks, s.Variants, strings.Join(splits, " "))
}

// Describe summarises a corpus.
func Describe(examples []Example) Stats {
	s := Stats{Splits: map[string]int{}, ByTask: map[string]int{}}
	inputs := map[string]bool{}
	for _, e := range examples {
		s.Examples++
		s.Variants += len(e.Expect)
		inputs[e.Input] = true
		s.ByTask[e.Task]++
		if e.Split != "" {
			s.Splits[e.Split]++
		}
	}
	s.Inputs = len(inputs)
	s.Tasks = len(s.ByTask)
	return s
}

// Leakage reports inputs that appear in more than one split.
//
// It should always return nothing, because Assign groups by input. It exists
// because that is a property worth asserting rather than trusting: a corpus
// built by concatenating two files, or split by two different salts, breaks it
// silently, and the symptom is a test score that looks unusually good.
func Leakage(examples []Example) []string {
	splits := map[string]map[string]bool{}
	for _, e := range examples {
		if e.Split == "" {
			continue
		}
		if splits[e.Input] == nil {
			splits[e.Input] = map[string]bool{}
		}
		splits[e.Input][e.Split] = true
	}
	var bad []string
	for in, s := range splits {
		if len(s) > 1 {
			bad = append(bad, in)
		}
	}
	sort.Strings(bad)
	return bad
}
