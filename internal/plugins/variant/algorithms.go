// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package variant

import ()

// affixes are the common ecosystem/role brackets seen in real registry
// squatting. Ported verbatim from the afx plugin.
var (
	affixPrefixes = []string{"py", "py-", "python-", "node-", "js-", "go-", "lib", "lib-", "the-", "get-"}
	affixSuffixes = []string{"2", "3", "js", "-js", "py", "-py", "-python", "-cli", "-dev",
		"-core", "-utils", "-api", "-sdk", "-lib", ".js", "-ng", "-master", "-official"}
)

// separators are the word joiners a registry name may legally use. The empty
// string is one of them: "my-lib" and "mylib" are different packages.
var separators = []string{"-", "_", ".", ""}

// Clean enforces the generator contract on a list a generator built by hand.
//
// The contract is the one pkg/typo's uniq already applies to the character-level
// algorithms: a variant is never empty, never a repeat, and never the name it
// came from. The domain- and registry-shaped generators here assemble their
// output by appending to a slice instead, and each one that did so had drifted
// off the contract in its own way -- tld emitted the name's own suffix back,
// sep rejoined a name that already used that separator, and aci appended the
// same neighbour twice when two positions produced it.
//
// Applying it in one place means a new generator gets the contract by calling
// Clean rather than by remembering three separate rules.
func Clean(origin string, out []string) []string {
	if origin == "" {
		return nil
	}
	seen := make(map[string]bool, len(out))
	kept := out[:0]
	for _, v := range out {
		if v == "" || v == origin || seen[v] {
			continue
		}
		seen[v] = true
		kept = append(kept, v)
	}
	if len(kept) == 0 {
		return nil
	}
	return kept
}
