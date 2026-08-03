// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package atomicfile_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rangertaha/urlinsane/internal/atomicfile"
)

func TestWritePublishesTheNewContentAndLeavesNoTemporary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "asset.db")

	if err := atomicfile.Write(path, []byte("first"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := atomicfile.Write(path, []byte("second"), 0o640); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "second" {
		t.Errorf("file = %q, want the second write", got)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("directory holds %v, want only the published file; a leftover "+
			"temporary is a partial write that outlived its write", names)
	}
}

// The class test.
//
// The same partial-write bug was written three times — the blockstore, the scan
// index, and first-run asset extraction — and each time the truncated file that
// was left behind was then trusted by whatever check came next, because a
// partial file is indistinguishable from a complete one to an os.Stat.
//
// Extracting a helper only retires that if new call sites actually use it, and
// nothing makes them. So this asserts the absence of the mistake directly: no
// package under internal/ may call os.WriteFile. A durable write goes through
// atomicfile.Write; a genuinely throwaway write belongs in a test, which this
// does not scan.
func TestNoInternalPackageCallsOsWriteFile(t *testing.T) {
	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command("grep", "-rn", "--include=*.go", "os.WriteFile",
		filepath.Join(root, "internal")).Output()
	if err != nil && len(out) == 0 {
		return // grep exits 1 with no matches, which is the passing case
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}
		file, code := parts[0], parts[2]
		// Tests write scratch files; that is not a durability claim.
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		// A mention in a comment is the rationale, not a call.
		if before, _, found := strings.Cut(code, "os.WriteFile"); found && strings.Contains(before, "//") {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(code), "//") {
			continue
		}
		t.Errorf("%s calls os.WriteFile. It creates and truncates before it "+
			"writes, so an interrupted write leaves a partial file that the next "+
			"os.Stat cannot tell from a complete one. Use atomicfile.Write, or "+
			"move the write into a test if it is throwaway.", strings.TrimPrefix(line, root+"/"))
	}
}
