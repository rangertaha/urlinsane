// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// replace must not leave a truncated file behind. A plain WriteFile truncates
// before it writes, so an interrupted write leaves half a suffix list, which
// still loads and still parses -- wrongly, for every domain past the cut.
func TestReplaceLeavesTheOldFileWhenTheWriteFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "suffix.lst")
	if err := os.WriteFile(path, []byte("com\nnet\n"), 0o640); err != nil {
		t.Fatal(err)
	}

	// A directory that cannot be written to is the closest reliable stand-in
	// for a failed write: CreateTemp fails, and the target must be untouched.
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Skip("cannot make the directory read-only:", err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	if err := replace(path, []byte("org\n")); err == nil {
		t.Skip("the directory is still writable; running as root?")
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "com\nnet\n" {
		t.Errorf("suffix.lst = %q after a failed write, want the original", body)
	}
}

func TestReplaceIsAtomicAndOverwrites(t *testing.T) {
	path := filepath.Join(t.TempDir(), "suffix.lst")
	if err := replace(path, []byte("com\n")); err != nil {
		t.Fatal(err)
	}
	if err := replace(path, []byte("org\n")); err != nil {
		t.Fatal(err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "org\n" {
		t.Errorf("suffix.lst = %q, want the second write", body)
	}

	// No temporaries left in the directory beside it.
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "suffix.lst.") {
			t.Errorf("left a temporary behind: %s", e.Name())
		}
	}
}
