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
