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
package dataset

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// models is the dataset schema (read-only reference data).
func models() []interface{} {
	return []interface{}{
		&Keyboard{},

		// Language
		&Language{},
		&Word{},
		&Char{},
		&Sym{},

		// Domain
		&Prefix{},
		&Suffix{},
		&Domain{},
	}
}

// Config opens the dataset database. It self-heals a corrupt or unreadable file
// by recreating it empty, and falls back to an in-memory database, so a bad
// dataset can never crash the tool (previously a malformed file panicked inside
// the gorm/sqlite migrator).
func Config(path string) {
	if err := open(path); err == nil {
		return
	} else {
		log.Warnf("dataset db %q is unusable (%v); recreating it empty", path, err)
	}

	// The on-disk file is corrupt — remove it and try a fresh, empty database.
	_ = os.Remove(path)
	if err := open(path); err == nil {
		return
	} else {
		log.Errorf("could not recreate dataset db %q (%v); falling back to in-memory", path, err)
	}

	// Last resort: an in-memory database is always valid, keeping the tool
	// runnable (with empty reference data) rather than crashing.
	if err := open(":memory:"); err != nil {
		log.Errorf("in-memory dataset db failed: %v", err)
		DB = nil
	}
}

// open opens, probes and migrates the database at path, assigning DB on success.
// A recover guards against driver-level panics on a corrupt file.
func open(path string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic opening dataset db: %v", r)
		}
	}()

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		CreateBatchSize: 10000,
		Logger:          logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	// Probe the schema. A malformed file errors here rather than panicking
	// inside AutoMigrate's HasTable.
	var n int
	if err := db.Raw("SELECT count(*) FROM sqlite_master").Scan(&n).Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(models()...); err != nil {
		return err
	}

	DB = db
	return nil
}
