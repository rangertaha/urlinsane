// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/dataset"
	"github.com/rangertaha/urlinsane/pkg/kb"
)

func TestScaffoldCoversEveryLanguageWithAKeyboard(t *testing.T) {
	root := write(t, map[string]string{"languages/en/synonym.lst": "login signin\n"})

	if _, err := Scaffold(root); err != nil {
		t.Fatal(err)
	}

	for _, code := range kb.Languages() {
		dir := filepath.Join(root, "languages", code)
		if _, err := os.Stat(dir); err != nil {
			t.Errorf("%s has a keyboard but no directory: %v", code, err)
		}
	}
}

func TestScaffoldWritesEveryRelation(t *testing.T) {
	root := t.TempDir()
	if _, err := Scaffold(root); err != nil {
		t.Fatal(err)
	}

	for _, r := range Relations {
		p := filepath.Join(root, "languages", "fr", r.Name+".lst")
		body, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("%s: %v", r.Name, err)
		}
		// The corpus is empty, but the file says what belongs in it: a
		// scaffolded tree that is eighty directories of nameless empty files
		// tells a curator nothing.
		if !strings.HasPrefix(string(body), "# ") {
			t.Errorf("%s has no header comment: %q", r.Name, body)
		}
	}
}

// The scaffold runs over a tree that people have hand-curated for years. It may
// add files; it may not touch one.
func TestScaffoldNeverOverwritesCuratedData(t *testing.T) {
	const curated = "login signin logon\n"
	root := write(t, map[string]string{"languages/en/synonym.lst": curated})

	for i := 0; i < 2; i++ {
		if _, err := Scaffold(root); err != nil {
			t.Fatal(err)
		}
	}

	body, err := os.ReadFile(filepath.Join(root, "languages", "en", "synonym.lst"))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != curated {
		t.Errorf("synonym.lst = %q, want the curated %q", body, curated)
	}
}

// A curated directory missing one of the relations gains that file, so adding a
// relation reaches the languages someone already wrote data for.
func TestScaffoldFillsGapsInACuratedDirectory(t *testing.T) {
	root := write(t, map[string]string{"languages/en/synonym.lst": "login signin\n"})

	if _, err := Scaffold(root); err != nil {
		t.Fatal(err)
	}

	for _, r := range Relations {
		p := filepath.Join(root, "languages", "en", r.Name+".lst")
		if _, err := os.Stat(p); err != nil {
			t.Errorf("en/%s.lst was not filled in: %v", r.Name, err)
		}
	}
}

// A directory named with a code kb does not use still covers the language kb
// knows by another name. Creating the other name would split one language
// across two directories -- words under one code, keyboard adjacency under the
// other -- which is the thing CanonicalCode exists to stop.
func TestScaffoldDoesNotSplitALanguageAcrossTwoCodes(t *testing.T) {
	root := write(t, map[string]string{
		"languages/iw/word.lst": "shalom\n", // retired code for Hebrew, kb has "he"
		"languages/no/word.lst": "hei\n",    // Norwegian, kb ships Bokmal as "nb"
	})

	if _, err := Scaffold(root); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct{ made, curated string }{
		{"he", "iw"},
		{"nb", "no"},
	} {
		if _, err := os.Stat(filepath.Join(root, "languages", tc.made)); err == nil {
			t.Errorf("scaffolded %s/ beside the curated %s/; one language, two directories",
				tc.made, tc.curated)
		}
	}
}

// Latin has curated data and no keyboard anywhere in kb. A scaffold that
// reconciled the tree to the catalogue in both directions would delete it.
func TestScaffoldLeavesLanguagesKbDoesNotShip(t *testing.T) {
	root := write(t, map[string]string{"languages/la/word.lst": "aqua\n"})

	if _, err := Scaffold(root); err != nil {
		t.Fatal(err)
	}

	body, err := os.ReadFile(filepath.Join(root, "languages", "la", "word.lst"))
	if err != nil {
		t.Fatal("the scaffold removed a language with no keyboard:", err)
	}
	if string(body) != "aqua\n" {
		t.Errorf("la/word.lst = %q, want it untouched", body)
	}
}

func TestMissingReportsWhatScaffoldWouldCreateWithoutWriting(t *testing.T) {
	root := t.TempDir()

	missing, err := Missing(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != len(kb.Languages()) {
		t.Errorf("Missing = %d languages, want all %d", len(missing), len(kb.Languages()))
	}
	if _, err := os.Stat(filepath.Join(root, "languages")); !os.IsNotExist(err) {
		t.Error("Missing wrote to the tree; it is meant to be a dry run")
	}

	if _, err := Scaffold(root); err != nil {
		t.Fatal(err)
	}
	after, err := Missing(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != 0 {
		t.Errorf("Missing = %v after scaffolding, want nothing left", after)
	}
}

// The scaffolded headers are comments, so a scaffolded language imports as an
// empty corpus rather than as thirteen files of English prose.
func TestScaffoldedFilesImportAsEmpty(t *testing.T) {
	fresh(t)
	root := t.TempDir()
	if _, err := Scaffold(root); err != nil {
		t.Fatal(err)
	}

	if err := All(root); err != nil {
		t.Fatal(err)
	}

	if n := count(t, &dataset.Vocabulary{}); n != 0 {
		t.Errorf("a scaffolded tree imported %d words, want none", n)
	}
}
