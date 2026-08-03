// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package gen

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// brokenRoot is a datasets tree whose import fails partway through.
//
// The failure is a single line past the scanner's 4MB token limit rather than
// an unreadable file: a chmod-based fixture passes trivially when the suite
// runs as root, and a directory in place of a corpus is skipped outright by
// All. An over-long line reaches Extract and errors there, which is the case
// that matters -- a build that dies after it has begun writing.
func brokenRoot(t *testing.T) string {
	t.Helper()
	return write(t, map[string]string{
		"languages/en/synonym.lst":     "login signin\n",
		"languages/en/homoglyph.lst":   strings.Repeat("a", 5<<20) + "\n",
		"languages/en/misspelling.lst": "hwile while\n",
	})
}

// digest is the file's content, so "unchanged" means byte-identical rather than
// merely still present.
func digest(t *testing.T, path string) [32]byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return sha256.Sum256(b)
}

// A build that fails must leave the previous database exactly where it was.
//
// Build used to remove the target before opening it, which made any later
// failure unrecoverable rather than merely unsuccessful: internal/config
// carries `//go:embed dataset.db`, so once the file is gone the module stops
// compiling -- including cmd/datasets, the binary needed to rebuild it. The
// only way back was `git checkout`. This is the regression test for that.
func TestFailedBuildLeavesThePreviousDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dataset.db")
	good := write(t, map[string]string{"languages/en/synonym.lst": "login signin\n"})

	if err := Build(path, good); err != nil {
		t.Fatal(err)
	}
	before := digest(t, path)

	if err := Build(path, brokenRoot(t)); err == nil {
		t.Fatal("Build reported success over an unreadable corpus")
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("a failed build destroyed the database: %v", err)
	}
	if after := digest(t, path); after != before {
		t.Error("a failed build overwrote the database; it should be untouched")
	}
}

// The temporary the build writes through must not survive it, in either
// direction -- a stray dataset.db.building beside the real file would be
// embedded by nothing and confusing to everyone.
func TestBuildLeavesNoTemporaryBehind(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dataset.db")
	root := write(t, map[string]string{"languages/en/synonym.lst": "login signin\n"})

	if err := Build(path, root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path + ".building"); !os.IsNotExist(err) {
		t.Error("a successful build left its temporary behind")
	}

	if err := Build(path, brokenRoot(t)); err == nil {
		t.Fatal("Build reported success over an unreadable corpus")
	}
	if _, err := os.Stat(path + ".building"); !os.IsNotExist(err) {
		t.Error("a failed build left its temporary behind")
	}
}
