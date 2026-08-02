// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

import (
	"fmt"
	"sort"
	"strings"
)

// ExcludePrefix marks an id as one to remove rather than keep (§12.10).
//
// '^' rather than '-' because "-whois" is indistinguishable from a flag to an
// argument parser, and rather than '!' or '~' because both are shell
// metacharacters that would need quoting.
const ExcludePrefix = "^"

// SelectOperators filters a set by id, for --algorithm and --collect.
//
//	nil or empty   every operator
//	"dns", "ptr"   only these
//	"^whois"       everything except this
//
// The two forms cannot be mixed in one call. A precedence rule between "keep
// only these" and "drop these" would have to be memorised to be used safely, and
// the combination expresses nothing the two forms cannot express separately.
//
// An unknown id is an error rather than a silent omission. A run that quietly
// dropped half of what the user asked for — or, worse, selected nothing and
// reported a clean result — would present doing nothing as a finding of
// nothing.
func SelectOperators(all []Operator, ids []string) ([]Operator, error) {
	want, exclude, err := parseSelection(ids)
	if err != nil {
		return nil, err
	}
	if len(want) == 0 {
		return all, nil
	}

	have := make(map[string]bool, len(all))
	for _, op := range all {
		have[op.Id()] = true
	}
	var unknown []string
	for id := range want {
		if !have[id] {
			unknown = append(unknown, id)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		known := make([]string, 0, len(have))
		for id := range have {
			known = append(known, id)
		}
		sort.Strings(known)
		return nil, fmt.Errorf("unknown id(s): %s; known: %s",
			strings.Join(unknown, ", "), strings.Join(known, ", "))
	}

	out := make([]Operator, 0, len(all))
	for _, op := range all {
		if want[op.Id()] != exclude {
			out = append(out, op)
		}
	}
	return out, nil
}

// parseSelection splits a selection into the id set and whether it excludes.
// Ids are comma-split as well as repeatable, so -c dns,ptr and -c dns -c ptr
// are the same.
func parseSelection(ids []string) (set map[string]bool, exclude bool, err error) {
	set = map[string]bool{}
	var sawKeep, sawDrop bool
	for _, raw := range ids {
		for _, id := range strings.Split(raw, ",") {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			if cut, ok := strings.CutPrefix(id, ExcludePrefix); ok {
				if cut == "" {
					return nil, false, fmt.Errorf("%q names no id", id)
				}
				sawDrop = true
				set[cut] = true
				continue
			}
			sawKeep = true
			set[id] = true
		}
	}
	if sawKeep && sawDrop {
		return nil, false, fmt.Errorf(
			"cannot mix kept and excluded ids in one selection; use either a list or %sid exclusions",
			ExcludePrefix)
	}
	return set, sawDrop, nil
}
