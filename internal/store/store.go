// Copyright 2024 Rangertaha. All Rights Reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package store

import (
	"time"

	"github.com/ipfs/go-cid"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/entity"
	"gorm.io/gorm"
)

// Store is the content-addressed result store: an IPLD blockstore (source of
// truth) plus a SQLite secondary index for name-based lookups and diffing.
type Store struct {
	bs  Blockstore
	idx *Index
}

// Open creates a Store with a filesystem blockstore under dir/blocks and a
// SQLite index on the supplied gorm connection.
func Open(dir string, idxDB *gorm.DB) (*Store, error) {
	bs, err := newFSStore(dir)
	if err != nil {
		return nil, err
	}
	idx, err := newIndex(idxDB)
	if err != nil {
		return nil, err
	}
	return &Store{bs: bs, idx: idx}, nil
}

// Put encodes a domain to an Entity block, stores it, indexes (type,name)→CID,
// and returns the CID. Idempotent for identical content.
func (s *Store) Put(d *db.Domain) (cid.Cid, error) {
	block, err := encodeEntity(ToEntity(d))
	if err != nil {
		return cid.Undef, err
	}
	c, err := cidOf(block)
	if err != nil {
		return cid.Undef, err
	}
	if err := s.bs.Put(c, block); err != nil {
		return cid.Undef, err
	}
	if err := s.idx.PutLatest(d.EntityType(), d.Name, c); err != nil {
		return cid.Undef, err
	}
	return c, nil
}

// Get loads the latest stored entity for (type, name).
func (s *Store) Get(t entity.Type, name string) (*db.Domain, bool, error) {
	c, ok, err := s.idx.LatestCID(t, name)
	if err != nil || !ok {
		return nil, false, err
	}
	d, err := s.GetCID(c)
	if err != nil {
		return nil, false, err
	}
	return d, true, nil
}

// GetCID decodes the block at a specific CID into a *db.Domain.
func (s *Store) GetCID(c cid.Cid) (*db.Domain, error) {
	block, err := s.bs.Get(c)
	if err != nil {
		return nil, err
	}
	e, err := decodeEntity(block)
	if err != nil {
		return nil, err
	}
	return ToDomain(e), nil
}

// PutScan stores each result, then a Scan node linking their CIDs, and records
// the scan in the index. Returns the Scan CID.
func (s *Store) PutScan(query string, results []*db.Domain) (cid.Cid, error) {
	scan := &Scan{Query: query, Created: time.Now().UTC().Format(time.RFC3339)}
	rows := make([]ResultRow, 0, len(results))
	for _, d := range results {
		c, err := s.Put(d)
		if err != nil {
			return cid.Undef, err
		}
		scan.Results = append(scan.Results, cidlink.Link{Cid: c})
		rows = append(rows, ResultRow{Name: d.Name, CID: c})
	}
	block, err := encodeScan(scan)
	if err != nil {
		return cid.Undef, err
	}
	sc, err := cidOf(block)
	if err != nil {
		return cid.Undef, err
	}
	if err := s.bs.Put(sc, block); err != nil {
		return cid.Undef, err
	}
	if err := s.idx.PutScan(query, sc, rows); err != nil {
		return cid.Undef, err
	}
	return sc, nil
}

// EncodeJSON returns the dag-json (nested) form of a domain, for output plugins.
func (s *Store) EncodeJSON(d *db.Domain) ([]byte, error) {
	return encodeEntityJSON(ToEntity(d))
}

// DiffResult reports how a query's results changed between its two most recent
// scans, by comparing content CIDs per entity name.
type DiffResult struct {
	Added   []string // names only in the newer scan
	Removed []string // names only in the older scan
	Changed []string // names present in both with a different CID
	Same    []string // names present in both with the same CID
}

// Diff compares the two most recent scans for a query.
func (s *Store) Diff(query string) (*DiffResult, error) {
	scans, err := s.idx.LatestScanCIDs(query, 2)
	if err != nil {
		return nil, err
	}
	out := &DiffResult{}
	if len(scans) < 2 {
		return out, nil
	}
	newer, err := s.nameToCID(scans[0])
	if err != nil {
		return nil, err
	}
	older, err := s.nameToCID(scans[1])
	if err != nil {
		return nil, err
	}
	for name, nc := range newer {
		oc, ok := older[name]
		switch {
		case !ok:
			out.Added = append(out.Added, name)
		case oc == nc:
			out.Same = append(out.Same, name)
		default:
			out.Changed = append(out.Changed, name)
		}
	}
	for name := range older {
		if _, ok := newer[name]; !ok {
			out.Removed = append(out.Removed, name)
		}
	}
	return out, nil
}

func (s *Store) nameToCID(scanCID string) (map[string]string, error) {
	rows, err := s.idx.ScanRows(scanCID)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(rows))
	for _, r := range rows {
		m[r.EntityName] = r.ResultCID
	}
	return m, nil
}
