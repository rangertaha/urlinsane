// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ipfs/go-cid"
)

// IndexFile is the scan index, beside the blockstore it indexes.
const IndexFile = "scans.json"

// Entry is one saved scan: what was scanned, where it landed, and the two facts
// the root deliberately does not carry.
//
// A scan root records the seed and the content and nothing else — no timestamp,
// no run id — which is what makes two identical scans produce the same CID and
// "what changed since last week" a CID comparison (Root's own doc says so). The
// cost is that a root cannot answer "which was most recent" or "was this one
// interrupted", so those live out here, beside the store rather than inside it.
type Entry struct {
	Type string `json:"type"`
	// Key is the canonical seed key, matched exactly. `report bob@acme.com`
	// finds the email scan, not the scan of acme.com nested inside it.
	Key  string    `json:"key"`
	Root string    `json:"root"`
	At   time.Time `json:"at"`
	// Partial marks a scan that stopped early — an interrupt, a deadline, a
	// budget. It is a fact about the scan, so every re-render reports it.
	Partial bool `json:"partial"`

	// Scope, Plan and Rounds are the rest of the run facts the root does not
	// carry, for the same reason At and Partial are here: they describe the
	// invocation rather than the content, so putting them in the root would
	// make two identical scans produce different CIDs.
	//
	// They were missing, and the report command therefore rendered them as
	// their zero values off bytes that had them: `typo --save-graph -o json`
	// printed rounds 4 with a scope and a plan hash, and `report` off the same
	// scan printed rounds 0 with neither. Same bytes, two renderings — the
	// failure the observer set had before it was persisted (9944703).
	Scope  []string `json:"scope,omitempty"`
	Plan   string   `json:"plan,omitempty"`
	Rounds int      `json:"rounds,omitempty"`
}

// Index is the list of saved scans, newest first.
type Index struct {
	path    string
	Entries []Entry `json:"scans"`
}

// OpenIndex reads the index in dir, or returns an empty one if there is none.
//
// A missing index is not an error: the first save creates it. A malformed one
// is, because silently starting over would strand every scan already in the
// blockstore with no way left to name it.
func OpenIndex(dir string) (*Index, error) {
	ix := &Index{path: filepath.Join(dir, IndexFile)}
	b, err := os.ReadFile(ix.path)
	if os.IsNotExist(err) {
		return ix, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, ix); err != nil {
		return nil, err
	}
	return ix, nil
}

// Add records a scan and writes the index back.
//
// Re-adding the same root replaces its entry rather than appending: the CID is
// the scan's identity, so two entries for one root would be two names for one
// thing, and "most recent" would depend on which was read first.
//
// The whole read-modify-write happens under a lock, and the index is re-read
// inside it. Without that, two scans saving at once each write the index they
// read at open time and the second erases the first: the losing scan's blocks
// are still in the store but nothing names them, so `report <target>` can never
// find it again. Two terminals is all it takes.
func (ix *Index) Add(e Entry) error {
	unlock, err := lock(filepath.Dir(ix.path))
	if err != nil {
		return err
	}
	defer unlock()

	// Re-read under the lock: another process may have written since OpenIndex.
	fresh, err := OpenIndex(filepath.Dir(ix.path))
	if err != nil {
		return err
	}

	out := fresh.Entries[:0]
	for _, x := range fresh.Entries {
		if x.Root != e.Root {
			out = append(out, x)
		}
	}
	fresh.Entries = append(out, e)
	sort.SliceStable(fresh.Entries, func(i, j int) bool {
		return fresh.Entries[i].At.After(fresh.Entries[j].At)
	})
	if err := fresh.save(); err != nil {
		return err
	}
	ix.Entries = fresh.Entries
	return nil
}

// lockName is the exclusive-create lockfile guarding the index.
const lockName = ".scans.lock"

// lock takes an exclusive lock on the index in dir.
//
// O_CREATE|O_EXCL rather than flock: this ships to fourteen platforms including
// Windows, and exclusive create is the one mechanism every one of them has.
//
// A lock older than lockStale is broken rather than waited on. A process killed
// between creating the lockfile and removing it would otherwise make the index
// permanently unwritable, which is a worse failure than the race the lock
// exists to prevent.
func lock(dir string) (func(), error) {
	const (
		lockStale = 30 * time.Second
		poll      = 10 * time.Millisecond
		timeout   = 10 * time.Second
	)
	path := filepath.Join(dir, lockName)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}

	deadline := time.Now().Add(timeout)
	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o640)
		if err == nil {
			f.Close()
			return func() { os.Remove(path) }, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}
		// Breaking a stale lock removes only the file that was observed to be
		// stale. An unconditional Remove let two waiters both break the same
		// lock: the first removed it and created its own, the second then
		// removed *that* fresh lock and created a third, and both proceeded
		// into the critical section — reintroducing exactly the lost update
		// the lock exists to prevent, with the loser's scan silently absent
		// from the index and its blocks left unreferenced.
		//
		// Re-stat and compare identity, so a lock created between the two
		// calls is a different file and is left alone.
		if fi, serr := os.Stat(path); serr == nil && time.Since(fi.ModTime()) > lockStale {
			if again, serr := os.Stat(path); serr == nil && os.SameFile(fi, again) {
				os.Remove(path)
			}
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("store: index locked by another process (%s)", path)
		}
		time.Sleep(poll)
	}
}

func (ix *Index) save() error {
	if err := os.MkdirAll(filepath.Dir(ix.path), 0o750); err != nil {
		return err
	}
	b, err := json.MarshalIndent(ix, "", "  ")
	if err != nil {
		return err
	}
	// A crash mid-write must not leave an index that parses as neither the old
	// nor the new one — which needs the flush writeAtomic does, not just the
	// rename this used to do on its own.
	return writeAtomic(ix.path, append(b, '\n'), 0o640)
}

// Scans returns every saved scan of one target, newest first.
func (ix *Index) Scans(typ, key string) []Entry {
	var out []Entry
	for _, e := range ix.Entries {
		if e.Type == typ && e.Key == key {
			out = append(out, e)
		}
	}
	return out
}

// Latest returns the most recent scan of a target.
func (ix *Index) Latest(typ, key string) (Entry, bool) {
	if s := ix.Scans(typ, key); len(s) > 0 {
		return s[0], true
	}
	return Entry{}, false
}

// At returns the entry for one root CID.
func (ix *Index) At(root string) (Entry, bool) {
	for _, e := range ix.Entries {
		if e.Root == root {
			return e, true
		}
	}
	return Entry{}, false
}

// Targets lists every distinct target in the index, for an empty `report`.
func (ix *Index) Targets() []Entry {
	seen := map[string]bool{}
	var out []Entry
	for _, e := range ix.Entries {
		k := e.Type + "\x00" + e.Key
		if !seen[k] {
			seen[k] = true
			out = append(out, e)
		}
	}
	return out
}

// ParseRoot is cid.Decode, named for what callers are doing with it.
func ParseRoot(s string) (cid.Cid, error) { return cid.Decode(s) }
