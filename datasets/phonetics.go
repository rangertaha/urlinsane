// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package datasets

import (
	"embed"
	"strings"
	"sync"
)

// crossFS carries the cross-language homophone groups.
//
// Embedded rather than imported into the dataset database, unlike the language
// relations, because this data is not *of* a language. Every row in that
// database hangs off a language id, and these groups deliberately span
// languages — "youtube" and "yutup" belong to no single one. Attaching them to
// a language would either duplicate every group under each contributing code or
// pick one arbitrarily, and both make the data lie about what it is.
//
//go:embed phonetics/homophone.lst
var crossFS embed.FS

var (
	crossOnce   sync.Once
	crossGroups [][]string
)

func loadCrossHomophones() {
	b, err := crossFS.ReadFile("phonetics/homophone.lst")
	if err != nil {
		// The embed directive guarantees the file exists; a failure here is a
		// build problem, not a runtime condition worth reporting to callers.
		return
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if words := strings.Fields(line); len(words) > 1 {
			crossGroups = append(crossGroups, words)
		}
	}
}

// CrossHomophones returns groups of spellings that sound alike to speakers of
// different languages.
//
// Each group is one pronunciation written several ways: boutique, boetiek,
// butik. Substituting one member for another inside a name is sound-squatting
// across a language boundary — the technique X-squatter measured (ACM TOPS
// 2024) and found to carry TLS certificates about twice as often as other
// squatting types.
//
// The slice is shared and must not be modified. Reading it is cheap after the
// first call; parsing happens once.
func CrossHomophones() [][]string {
	crossOnce.Do(loadCrossHomophones)
	return crossGroups
}
