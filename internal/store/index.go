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
	"github.com/rangertaha/urlinsane/internal/entity"
	"gorm.io/gorm"
)

// EntityIndex maps (entity type, name) to the latest root CID. IPLD has no
// query layer, so this index serves the name-based lookups the engine relies on.
type EntityIndex struct {
	EntityType string `gorm:"primaryKey"`
	Name       string `gorm:"primaryKey"`
	RootCID    string `gorm:"column:root_cid;index"`
	UpdatedAt  time.Time
}

// ScanIndex records one row per (scan, result), giving the scan→results
// aggregate and the history needed for cross-scan diff.
type ScanIndex struct {
	ID         uint   `gorm:"primaryKey"`
	Query      string `gorm:"index"`
	ScanCID    string `gorm:"column:scan_cid;index"`
	ResultCID  string `gorm:"column:result_cid"`
	EntityName string
	CreatedAt  time.Time
}

// ResultRow is one (name, CID) pair within a scan.
type ResultRow struct {
	Name string
	CID  cid.Cid
}

// Index is the SQLite-backed secondary index over the IPLD blockstore.
type Index struct {
	db *gorm.DB
}

func newIndex(gdb *gorm.DB) (*Index, error) {
	if err := gdb.AutoMigrate(&EntityIndex{}, &ScanIndex{}); err != nil {
		return nil, err
	}
	return &Index{db: gdb}, nil
}

// LatestCID returns the latest root CID for (type, name), if indexed. Uses
// Find (not First) to avoid gorm logging ErrRecordNotFound on the common
// "new entity" path.
func (i *Index) LatestCID(t entity.Type, name string) (cid.Cid, bool, error) {
	var rows []EntityIndex
	if err := i.db.Limit(1).Find(&rows, "entity_type = ? AND name = ?", string(t), name).Error; err != nil {
		return cid.Undef, false, err
	}
	if len(rows) == 0 {
		return cid.Undef, false, nil
	}
	c, err := cid.Decode(rows[0].RootCID)
	if err != nil {
		return cid.Undef, false, err
	}
	return c, true, nil
}

// PutLatest upserts the latest root CID for (type, name).
func (i *Index) PutLatest(t entity.Type, name string, c cid.Cid) error {
	return i.db.Save(&EntityIndex{
		EntityType: string(t),
		Name:       name,
		RootCID:    c.String(),
		UpdatedAt:  time.Now(),
	}).Error
}

// PutScan records a scan and its result rows.
func (i *Index) PutScan(query string, scanCID cid.Cid, results []ResultRow) error {
	rows := make([]ScanIndex, 0, len(results))
	now := time.Now()
	for _, r := range results {
		rows = append(rows, ScanIndex{
			Query:      query,
			ScanCID:    scanCID.String(),
			ResultCID:  r.CID.String(),
			EntityName: r.Name,
			CreatedAt:  now,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	return i.db.Create(&rows).Error
}

// LatestScanCIDs returns up to n distinct scan CIDs for a query, newest first.
func (i *Index) LatestScanCIDs(query string, n int) ([]string, error) {
	var cids []string
	err := i.db.Model(&ScanIndex{}).
		Where("query = ?", query).
		Group("scan_cid").
		Order("MAX(created_at) desc").
		Limit(n).
		Pluck("scan_cid", &cids).Error
	return cids, err
}

// ScanRows returns the result rows for a given scan CID.
func (i *Index) ScanRows(scanCID string) ([]ScanIndex, error) {
	var rows []ScanIndex
	err := i.db.Where("scan_cid = ?", scanCID).Find(&rows).Error
	return rows, err
}
