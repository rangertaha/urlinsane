// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package eval measures the algorithm set against names that were actually
// registered, rather than against itself.
//
// The question it answers is recall: of the lookalike domains someone really
// registered for a brand, how many would urlinsane have generated? Without that
// number, "the algorithms produce plausible variants" is unfalsifiable, and a
// change that adds ten thousand candidates looks the same as one that finds a
// real squat.
//
// Nothing here runs during a scan. It is a development harness.
package eval

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/rangertaha/urlinsane/internal/plugins/variant"
)

// Record is one observed lookalike: a name that exists, attributed to the brand
// it imitates.
//
// Truth sets are JSONL so they diff line-by-line in review. These are candidate
// observations from an automated filter, not adjudicated ground truth — see
// Reviewed.
type Record struct {
	// Brand is the seed this record is a lookalike of, as a full key.
	Brand string `json:"brand"`
	// Name is the observed name, normalized.
	Name string `json:"name"`
	// Type is the node type both Brand and Name are ("domain", "package", ...).
	Type string `json:"type"`
	// Source names where the observation came from ("crt.sh", ...) so a set
	// assembled from several feeds stays auditable.
	Source string `json:"source"`
	// FirstSeen is an ISO date if the source supplies one, else empty.
	FirstSeen string `json:"first_seen,omitempty"`
	// Reviewed records whether a human confirmed this is a lookalike rather
	// than a legitimate brand-owned name the filter failed to exclude. An
	// unreviewed set still produces a usable recall number, but the number is
	// a lower bound on precision, not a measurement of it.
	Reviewed bool `json:"reviewed"`
}

// Core is the part of the name the variant algorithms actually vary — the
// registrable label for a domain, the whole key otherwise.
func (r Record) Core() string { return coreOf(r.Type, r.Name) }

func coreOf(nodeType, key string) string {
	if nodeType == variant.TypeDomain {
		_, core, _ := variant.SplitDomain(key)
		return core
	}
	return key
}

// Normalize puts an observed name into the form the graph would canonicalize it
// to: lowercase, no wildcard label, no trailing dot, no surrounding space.
// Returns "" for anything that is not a usable name.
func Normalize(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.TrimPrefix(name, "*.")
	name = strings.Trim(name, ".")
	if name == "" || strings.ContainsAny(name, " \t/\\") {
		return ""
	}
	return name
}

// ReadTruth loads a JSONL truth set. Blank lines are skipped; a malformed line
// is an error rather than a silent drop, because a truth set that quietly lost
// half its records would report an inflated recall.
func ReadTruth(r io.Reader) ([]Record, error) {
	var out []Record
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for line := 1; sc.Scan(); line++ {
		text := strings.TrimSpace(sc.Text())
		if text == "" {
			continue
		}
		var rec Record
		if err := json.Unmarshal([]byte(text), &rec); err != nil {
			return nil, fmt.Errorf("eval: truth line %d: %w", line, err)
		}
		if rec.Name == "" || rec.Brand == "" {
			return nil, fmt.Errorf("eval: truth line %d: record needs both brand and name", line)
		}
		if rec.Type == "" {
			rec.Type = variant.TypeDomain
		}
		out = append(out, rec)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("eval: reading truth: %w", err)
	}
	return out, nil
}

// WriteTruth writes records as JSONL, sorted by brand then name so that
// re-running a fetch produces a diff of what changed rather than a reshuffle.
func WriteTruth(w io.Writer, records []Record) error {
	sorted := append([]Record(nil), records...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Brand != sorted[j].Brand {
			return sorted[i].Brand < sorted[j].Brand
		}
		return sorted[i].Name < sorted[j].Name
	})
	enc := json.NewEncoder(w)
	for _, rec := range sorted {
		if err := enc.Encode(rec); err != nil {
			return fmt.Errorf("eval: writing truth: %w", err)
		}
	}
	return nil
}

// ByBrand groups records by their brand key.
func ByBrand(records []Record) map[string][]Record {
	out := make(map[string][]Record)
	for _, r := range records {
		out[r.Brand] = append(out[r.Brand], r)
	}
	return out
}

// Brands returns the distinct brand keys in a truth set, sorted.
func Brands(records []Record) []string {
	seen := map[string]bool{}
	var out []string
	for _, r := range records {
		if !seen[r.Brand] {
			seen[r.Brand] = true
			out = append(out, r.Brand)
		}
	}
	sort.Strings(out)
	return out
}
