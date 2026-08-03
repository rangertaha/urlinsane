// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package gen

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/rangertaha/urlinsane/internal/dataset"
)

// fresh opens an empty database in a temp directory and returns its path.
//
// A real SQLite file rather than :memory: because that is what Build produces
// and what ships; an in-memory double would not have caught the schema drift
// this package exists to prevent.
func fresh(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	dataset.Config(path)
	if dataset.DB == nil {
		t.Fatal("no database")
	}
	return path
}

// write lays out a datasets tree from a map of relative path to contents.
func write(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, body := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func count(t *testing.T, model any) int64 {
	t.Helper()
	var n int64
	if err := dataset.DB.Model(model).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	return n
}

func TestExtractSplitsLinesIntoWords(t *testing.T) {
	root := write(t, map[string]string{"a.lst": "one   two\tthree\n\n  four  \n"})
	got, err := Extract(filepath.Join(root, "a.lst"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d lines, want 2 (the blank one is dropped): %v", len(got), got)
	}
	if len(got[0]) != 3 || got[0][2] != "three" {
		t.Errorf("runs of whitespace should collapse to one separator: %v", got[0])
	}
	if len(got[1]) != 1 || got[1][0] != "four" {
		t.Errorf("surrounding whitespace should be trimmed: %v", got[1])
	}
}

// An unreadable file must be an error, not an empty corpus. It used to be
// printed to stdout and treated as zero words, so a missing file imported
// successfully as nothing.
func TestExtractReportsAMissingFile(t *testing.T) {
	if _, err := Extract(filepath.Join(t.TempDir(), "nope.lst")); err == nil {
		t.Fatal("a missing file was accepted")
	}
}

func TestOneStoresVocabularyAndTransitions(t *testing.T) {
	fresh(t)
	root := write(t, map[string]string{"syn.lst": "login signin\npay\n"})

	if err := One("en", "synonym", filepath.Join(root, "syn.lst")); err != nil {
		t.Fatal(err)
	}
	if got := count(t, &dataset.Vocabulary{}); got != 3 {
		t.Errorf("vocabulary = %d, want 3 (login, signin, pay)", got)
	}
	// A one-word line has nothing to associate with, so it contributes no edge.
	// Both directions of the two-word line do.
	if got := count(t, &dataset.Transition{}); got != 2 {
		t.Errorf("transitions = %d, want 2 (login<->signin only)", got)
	}
}

// Probabilities out of a word must sum to 1, or the weights mean nothing.
func TestTransitionProbabilitiesAreNormalised(t *testing.T) {
	fresh(t)
	root := write(t, map[string]string{"g.lst": "a b c\n"})
	if err := One("en", "group", filepath.Join(root, "g.lst")); err != nil {
		t.Fatal(err)
	}

	var a dataset.Vocabulary
	if err := dataset.DB.Where("token = ?", "a").First(&a).Error; err != nil {
		t.Fatal(err)
	}
	var edges []dataset.Transition
	dataset.DB.Where("src = ?", a.ID).Find(&edges)
	if len(edges) != 2 {
		t.Fatalf("a has %d edges, want 2", len(edges))
	}
	var sum float64
	for _, e := range edges {
		sum += e.Probability
	}
	if math.Abs(sum-1) > 1e-9 {
		t.Errorf("probabilities out of a sum to %v, want 1", sum)
	}
}

// Re-importing replaces rather than duplicating. The importer is run repeatedly
// against the same tree, and a merge would double every row each time.
func TestOneIsIdempotent(t *testing.T) {
	fresh(t)
	root := write(t, map[string]string{"s.lst": "login signin\n"})
	p := filepath.Join(root, "s.lst")

	for i := 0; i < 3; i++ {
		if err := One("en", "synonym", p); err != nil {
			t.Fatal(err)
		}
	}
	if got := count(t, &dataset.Vocabulary{}); got != 2 {
		t.Errorf("vocabulary = %d after three imports, want 2", got)
	}
	if got := count(t, &dataset.Transition{}); got != 2 {
		t.Errorf("transitions = %d after three imports, want 2", got)
	}
}

// A word dropped from the corpus must leave, taking its edges with it —
// otherwise a transition points at a vocabulary id that no longer exists.
func TestReimportLeavesNoDanglingTransitions(t *testing.T) {
	fresh(t)
	root := t.TempDir()
	p := filepath.Join(root, "s.lst")

	if err := os.WriteFile(p, []byte("login signin verify\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := One("en", "synonym", p); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("login signin\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := One("en", "synonym", p); err != nil {
		t.Fatal(err)
	}

	var edges []dataset.Transition
	dataset.DB.Find(&edges)
	ids := map[uint]bool{}
	var vocab []dataset.Vocabulary
	dataset.DB.Find(&vocab)
	for _, v := range vocab {
		ids[v.ID] = true
	}
	for _, e := range edges {
		if !ids[e.Src] || !ids[e.Dest] {
			t.Fatalf("transition %d->%d survives with no vocabulary row", e.Src, e.Dest)
		}
	}
}

// languages/<code>/<rel>.lst is keyed by language; <group>/<name>.lst is not.
func TestAllKeysLanguageFilesByCode(t *testing.T) {
	fresh(t)
	root := write(t, map[string]string{
		"languages/en/synonym.lst": "login signin\n",
		"languages/fr/synonym.lst": "connexion identifiant\n",
		"packages/npm.lst":         "react\nvue\n",
	})
	if err := All(root); err != nil {
		t.Fatal(err)
	}

	var en dataset.Language
	if err := dataset.DB.Where("code = ?", "en").First(&en).Error; err != nil {
		t.Fatal("en was not created:", err)
	}
	var ds dataset.Dataset
	if err := dataset.DB.Where("name = ?", "packages/npm").First(&ds).Error; err != nil {
		t.Fatal("a group file should be named by its path:", err)
	}

	// The list-shaped corpus belongs to no language.
	var n int64
	dataset.DB.Model(&dataset.Vocabulary{}).
		Where("dataset = ? AND language = 0", ds.ID).Count(&n)
	if n != 2 {
		t.Errorf("packages/npm has %d language-less rows, want 2", n)
	}
}

func TestLanguagesComeFromTheKeyboardCatalogue(t *testing.T) {
	fresh(t)
	if err := Languages(); err != nil {
		t.Fatal(err)
	}
	// More than the curated directories: a language with a layout is listed
	// whether or not anyone has written a word list for it.
	if got := count(t, &dataset.Language{}); got < 100 {
		t.Fatalf("languages = %d, want the full keyboard catalogue", got)
	}
	var fr dataset.Language
	if err := dataset.DB.Where("code = ?", "fr").First(&fr).Error; err != nil {
		t.Fatal(err)
	}
	if fr.Name != "French" {
		t.Errorf("fr resolved to %q, want a display name", fr.Name)
	}
}

func TestSourcesLoadEachShape(t *testing.T) {
	fresh(t)
	dir := write(t, map[string]string{
		"packages.lst":  "# comment\npypi https://pypi.org/project/%s/ https://pypi.org/pypi/%s/json\n",
		"usernames.lst": "github https://github.com/%s\n",
		"email.lst":     "gmail.com\n",
	})
	if err := Sources(dir); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct{ kind, code, url string }{
		// Three columns: the check URL is what the prober calls, so that is the
		// one stored.
		{"package", "pypi", "https://pypi.org/pypi/%s/json"},
		{"username", "github", "https://github.com/%s"},
		{"email", "gmail.com", "gmail.com"},
	} {
		var got dataset.Source
		if err := dataset.DB.Where("type = ? AND code = ?", tc.kind, tc.code).
			First(&got).Error; err != nil {
			t.Fatalf("%s/%s missing: %v", tc.kind, tc.code, err)
		}
		if got.URL != tc.url {
			t.Errorf("%s/%s url = %q, want %q", tc.kind, tc.code, got.URL, tc.url)
		}
	}
	if got := count(t, &dataset.Source{}); got != 3 {
		t.Errorf("sources = %d, want 3; the comment line must not become a row", got)
	}
}

// A source removed from the list must disappear, or a retired registry keeps
// being probed forever.
func TestSourcesReplaceRatherThanMerge(t *testing.T) {
	fresh(t)
	dir := t.TempDir()
	p := filepath.Join(dir, "usernames.lst")

	if err := os.WriteFile(p, []byte("github https://github.com/%s\ngitlab https://gitlab.com/%s\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Sources(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("github https://github.com/%s\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Sources(dir); err != nil {
		t.Fatal(err)
	}

	if got := count(t, &dataset.Source{}); got != 1 {
		t.Errorf("sources = %d after removing one, want 1", got)
	}
}

func TestBuildProducesAWholeDatabase(t *testing.T) {
	root := write(t, map[string]string{
		"languages/en/synonym.lst": "login signin\n",
		"sources/usernames.lst":    "github https://github.com/%s\n",
	})
	path := filepath.Join(t.TempDir(), "built.db")

	if err := Build(path, root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("Build wrote no file:", err)
	}
	if count(t, &dataset.Language{}) < 100 {
		t.Error("Build did not seed languages")
	}
	if count(t, &dataset.Vocabulary{}) == 0 {
		t.Error("Build imported no vocabulary")
	}
	if count(t, &dataset.Source{}) != 1 {
		t.Error("Build loaded no sources")
	}
}

// Build starts from nothing. Migrating an existing file is what left the
// shipped database carrying columns from three schema generations, with the
// current ones empty because AutoMigrate can add a column but cannot fill it.
func TestBuildDiscardsWhateverWasThere(t *testing.T) {
	path := filepath.Join(t.TempDir(), "built.db")
	root := write(t, map[string]string{"languages/en/synonym.lst": "login signin\n"})

	if err := Build(path, root); err != nil {
		t.Fatal(err)
	}
	dataset.DB.Create(&dataset.Vocabulary{Language: 999, Dataset: 999, Token: "stale"})
	before := count(t, &dataset.Vocabulary{})

	if err := Build(path, root); err != nil {
		t.Fatal(err)
	}
	after := count(t, &dataset.Vocabulary{})
	if after >= before {
		t.Errorf("rows = %d after rebuild, was %d with the injected row; the file was reused", after, before)
	}
	var stale int64
	dataset.DB.Model(&dataset.Vocabulary{}).Where("token = ?", "stale").Count(&stale)
	if stale != 0 {
		t.Error("a row from the previous database survived the rebuild")
	}
}

// Dataset directories were named by hand and some use a retired BCP 47 code.
// Left alone, "iw" produced a second Language row for Hebrew alongside kb's
// "he" -- one with words and no keyboard, one with a keyboard and no words.
func TestCanonicalCodeMergesRetiredCodes(t *testing.T) {
	if got := CanonicalCode("iw"); got != "he" {
		t.Errorf("CanonicalCode(iw) = %q, want he", got)
	}
}

// A code with no keyboard layout keeps its own identity: "la" and "no" are real
// languages with curated data, and renaming them to tidy the table would lose
// that data.
func TestCanonicalCodeKeepsLanguagesKbDoesNotShip(t *testing.T) {
	for _, code := range []string{"la", "no"} {
		if got := CanonicalCode(code); got != code {
			t.Errorf("CanonicalCode(%s) = %q, want it unchanged", code, got)
		}
	}
}

func TestCanonicalCodeLeavesNonsenseAlone(t *testing.T) {
	for _, code := range []string{"", "zzzz", "not-a-tag"} {
		if got := CanonicalCode(code); got != code {
			t.Errorf("CanonicalCode(%q) = %q, want it unchanged", code, got)
		}
	}
}

// The whole point: importing the iw directory must not create a second row.
func TestImportUsesTheCanonicalLanguageRow(t *testing.T) {
	fresh(t)
	if err := Languages(); err != nil {
		t.Fatal(err)
	}
	before := count(t, &dataset.Language{})

	root := write(t, map[string]string{"languages/iw/synonym.lst": "שלום היי\n"})
	if err := All(root); err != nil {
		t.Fatal(err)
	}
	if got := count(t, &dataset.Language{}); got != before {
		t.Errorf("importing iw added %d language rows; it should reuse he", got-before)
	}

	var he dataset.Language
	if err := dataset.DB.Where("code = ?", "he").First(&he).Error; err != nil {
		t.Fatal(err)
	}
	var n int64
	dataset.DB.Model(&dataset.Vocabulary{}).Where("language = ?", he.ID).Count(&n)
	if n == 0 {
		t.Error("the iw corpus did not land on the he row")
	}
	var iw int64
	dataset.DB.Model(&dataset.Language{}).Where("code = ?", "iw").Count(&iw)
	if iw != 0 {
		t.Error("a separate iw row was created")
	}
}
