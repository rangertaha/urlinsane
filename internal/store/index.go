// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"encoding/json"
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
func (ix *Index) Add(e Entry) error {
	out := ix.Entries[:0]
	for _, x := range ix.Entries {
		if x.Root != e.Root {
			out = append(out, x)
		}
	}
	ix.Entries = append(out, e)
	sort.SliceStable(ix.Entries, func(i, j int) bool {
		return ix.Entries[i].At.After(ix.Entries[j].At)
	})
	return ix.save()
}

func (ix *Index) save() error {
	if err := os.MkdirAll(filepath.Dir(ix.path), 0o750); err != nil {
		return err
	}
	b, err := json.MarshalIndent(ix, "", "  ")
	if err != nil {
		return err
	}
	// Write and rename: a crash mid-write must not leave an index that parses
	// as neither the old nor the new one.
	tmp := ix.path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o640); err != nil {
		return err
	}
	return os.Rename(tmp, ix.path)
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
