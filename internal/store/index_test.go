// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func at(min int) time.Time {
	return time.Date(2026, 8, 2, 12, min, 0, 0, time.UTC)
}

func TestIndexIsCreatedOnFirstAdd(t *testing.T) {
	dir := t.TempDir()
	ix, err := OpenIndex(dir)
	if err != nil {
		t.Fatal("a missing index must not be an error:", err)
	}
	if err := ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "bafy1", At: at(0)}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, IndexFile)); err != nil {
		t.Fatal("Add wrote no index:", err)
	}

	reopened, err := OpenIndex(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := reopened.Latest("domain", "acme.com"); !ok || got.Root != "bafy1" {
		t.Errorf("reopened index lost the entry: %+v", got)
	}
}

// A malformed index is an error, not a fresh start: silently starting over
// strands every scan already in the blockstore with no way left to name it.
func TestMalformedIndexIsAnError(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, IndexFile), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenIndex(dir); err == nil {
		t.Fatal("a corrupt index was accepted as empty")
	}
}

func TestLatestIsTheMostRecent(t *testing.T) {
	ix, _ := OpenIndex(t.TempDir())
	// Added out of order: recency is the timestamp, not insertion order.
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "old", At: at(0)})
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "new", At: at(30)})
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "mid", At: at(15)})

	got, ok := ix.Latest("domain", "acme.com")
	if !ok || got.Root != "new" {
		t.Errorf("Latest = %q, want new", got.Root)
	}
	scans := ix.Scans("domain", "acme.com")
	if len(scans) != 3 || scans[0].Root != "new" || scans[2].Root != "old" {
		t.Errorf("Scans is not newest-first: %v", scans)
	}
}

// The CID is the scan's identity, so re-saving one must not create a second
// entry for the same object.
func TestAddReplacesTheSameRoot(t *testing.T) {
	ix, _ := OpenIndex(t.TempDir())
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "same", At: at(0)})
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "same", At: at(10)})

	if got := ix.Scans("domain", "acme.com"); len(got) != 1 {
		t.Fatalf("one root produced %d entries", len(got))
	}
	if got, _ := ix.Latest("domain", "acme.com"); !got.At.Equal(at(10)) {
		t.Errorf("re-adding kept the older timestamp: %v", got.At)
	}
}

// Targets are matched exactly: an email scan is not a scan of the domain
// inside it, and reporting one for the other would be silently wrong.
func TestTargetsAreMatchedExactly(t *testing.T) {
	ix, _ := OpenIndex(t.TempDir())
	_ = ix.Add(Entry{Type: "email", Key: "bob@acme.com", Root: "e", At: at(0)})

	if _, ok := ix.Latest("domain", "acme.com"); ok {
		t.Error("an email scan answered a lookup for its domain")
	}
	if _, ok := ix.Latest("domain", "bob@acme.com"); ok {
		t.Error("the type is not part of the match")
	}
	if _, ok := ix.Latest("email", "bob@acme.com"); !ok {
		t.Error("the exact target did not match")
	}
}

func TestTargetsListsEachOnce(t *testing.T) {
	ix, _ := OpenIndex(t.TempDir())
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "a", At: at(0)})
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "b", At: at(10)})
	_ = ix.Add(Entry{Type: "domain", Key: "other.com", Root: "c", At: at(5)})

	got := ix.Targets()
	if len(got) != 2 {
		t.Fatalf("Targets = %d, want 2 distinct", len(got))
	}
	if got[0].Key != "acme.com" {
		t.Errorf("Targets is not newest-first: %v", got[0])
	}
}

// Partial is a fact about the scan, so it survives the round trip and every
// re-render reports it.
func TestPartialSurvivesReopening(t *testing.T) {
	dir := t.TempDir()
	ix, _ := OpenIndex(dir)
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "p", At: at(0), Partial: true})

	reopened, err := OpenIndex(dir)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := reopened.Latest("domain", "acme.com")
	if !got.Partial {
		t.Error("partial was lost across a save and reopen")
	}
}

// A crash mid-write must not leave an index that parses as neither the old nor
// the new one, so the write goes to a temp file and renames.
func TestSaveLeavesNoTempFile(t *testing.T) {
	dir := t.TempDir()
	ix, _ := OpenIndex(dir)
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "a", At: at(0)})

	if _, err := os.Stat(filepath.Join(dir, IndexFile+".tmp")); !os.IsNotExist(err) {
		t.Error("the temp file was left behind")
	}
}

func TestAtFindsOneRoot(t *testing.T) {
	ix, _ := OpenIndex(t.TempDir())
	_ = ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "wanted", At: at(0)})

	if got, ok := ix.At("wanted"); !ok || got.Key != "acme.com" {
		t.Errorf("At = %+v, %v", got, ok)
	}
	if _, ok := ix.At("absent"); ok {
		t.Error("At matched a root that is not there")
	}
}

// Two scans saving at once must both survive. Without a lock around the
// read-modify-write each writes the index it read at open time, and the loser's
// blocks stay in the store with nothing naming them — `report <target>` can
// never find that scan again.
func TestConcurrentAddsAreNotLost(t *testing.T) {
	dir := t.TempDir()
	const n = 8

	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ix, err := OpenIndex(dir)
			if err != nil {
				errs <- err
				return
			}
			errs <- ix.Add(Entry{
				Type: "domain",
				Key:  fmt.Sprintf("t%d.com", i),
				Root: fmt.Sprintf("root%d", i),
				At:   at(i),
			})
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}

	ix, err := OpenIndex(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ix.Entries) != n {
		t.Fatalf("index holds %d of %d concurrent saves; updates were lost", len(ix.Entries), n)
	}
	for i := 0; i < n; i++ {
		if _, ok := ix.At(fmt.Sprintf("root%d", i)); !ok {
			t.Errorf("root%d is unreachable: its scan can never be reported", i)
		}
	}
}

// A lock left behind by a killed process must not make the index permanently
// unwritable — that is a worse failure than the race the lock prevents.
func TestStaleLockIsBroken(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, lockName)
	if err := os.WriteFile(path, nil, 0o640); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}

	ix, _ := OpenIndex(dir)
	if err := ix.Add(Entry{Type: "domain", Key: "acme.com", Root: "a", At: at(0)}); err != nil {
		t.Fatalf("a stale lock blocked the write: %v", err)
	}
}

// The lock is released on success, or the next save in the same process hangs.
func TestLockIsReleased(t *testing.T) {
	dir := t.TempDir()
	ix, _ := OpenIndex(dir)
	for i := 0; i < 3; i++ {
		if err := ix.Add(Entry{Type: "domain", Key: "acme.com",
			Root: fmt.Sprintf("r%d", i), At: at(i)}); err != nil {
			t.Fatalf("add %d: %v", i, err)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, lockName)); !os.IsNotExist(err) {
		t.Error("the lockfile was left behind")
	}
}
