// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package dataset

import (
	"fmt"
	"io"
	"os"

	"github.com/glebarez/sqlite"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Language struct {
	ID   uint
	Code string `gorm:"unique"`
	Name string
}

type Dataset struct {
	ID          uint
	Name        string `gorm:"unique"`
	Description string
}

// Vocabulary is one token of one relation of one language.
//
// (language, dataset) is indexed because it is how every read starts: lang.go's
// tokens, edges and groups all select a whole relation by that pair. Without
// the index SQLite full-scans the table, and the table is 436k rows in the
// shipped database -- so a scan over the language-driven algorithms paid a full
// scan per language per relation, and adding a corpus for one language slowed
// the reads for all the others.
type Vocabulary struct {
	ID       uint
	Language uint `gorm:"index:idx_vocabularies_language_dataset,priority:1"` // language ID
	Dataset  uint `gorm:"index:idx_vocabularies_language_dataset,priority:2"` // dataset ID
	Token    string
}

// Transition is one weighted edge between two vocabulary tokens.
//
// Src is indexed for the same reason: transitions are always read as "the edges
// leaving this relation's vocabulary", which is a lookup by Src, over 386k rows.
type Transition struct {
	ID          uint
	Src         uint `gorm:"index"` // vocab ID
	Dest        uint // vocab ID
	Probability float64
}

type Source struct {
	ID     uint
	Type   string `gorm:"index"` // package | repository | username | email
	Code   string
	Params string // {}
	// URL is the page a human would open: https://www.npmjs.com/package/%s.
	// It names the platform and is what a report links to.
	URL string
	// CheckURL is the endpoint the prober calls to decide existence:
	// https://registry.npmjs.org/%s. Empty means URL answers both questions,
	// which is the shape the username and email lists have.
	//
	// The two are separate columns because they are separate facts, and folding
	// them into one lost whichever was not stored. Keeping only the check URL
	// named the platform after the API host, so GitHub arrived as
	// api.github.com from the repo list and github.com from the username list —
	// two platform nodes for one forge, and no analyzer able to see that a
	// username and a repository sit on the same one.
	CheckURL string
	Success  string // regex for a successful response
	Failed   string // regex for a failed response
}

// models is the dataset schema (read-only reference data).
func models() []interface{} {
	return []interface{}{

		// Weighted word tables, built from every dataset
		&Dataset{},
		&Language{},
		&Vocabulary{},
		&Transition{},
		&Source{},
	}
}

// Config opens the dataset database, falling back to an in-memory one so a bad
// dataset can never crash the tool (previously a malformed file panicked inside
// the gorm/sqlite migrator). It removes the file only when it is provably not a
// database; see the comment in the body for why that distinction matters.
func Config(path string) {
	err := open(path)
	if err == nil {
		return
	}

	// Delete only what is provably not a database, and never leave an empty one
	// in its place.
	//
	// This used to remove the file on *any* open failure and immediately create
	// a fresh empty database at the same path. open() fails for transient
	// reasons too -- AutoMigrate runs DDL whenever the on-disk schema predates a
	// change, and the driver is opened without a busy timeout, so a second
	// process running at the same moment gets "database is locked" -- and the
	// remedy destroyed the reference data with nothing to restore it:
	// config.extract only re-extracts when the file is *absent*, and the empty
	// replacement is present and opens cleanly. Every later scan then generated
	// no language-driven variants at all and still reported success.
	//
	// A file carrying the SQLite header is a database that would not open right
	// now, which is not the same as a corrupt one. Leave it where it is and run
	// from memory; the next run, without the contention, opens it.
	if isSQLiteFile(path) {
		log.Errorf("dataset db %q could not be opened (%v); leaving it in place "+
			"and running from memory for this run", path, err)
	} else {
		// Not a database at all. Remove it and do NOT recreate: an empty file
		// here would look valid forever, where an absent one is re-extracted.
		_ = os.Remove(path)
		log.Errorf("dataset db %q is not a database (%v); it has been removed and "+
			"will be re-extracted on the next run", path, err)
	}

	// Last resort: an in-memory database is always valid, keeping the tool
	// runnable (with empty reference data) rather than crashing.
	if err := open(":memory:"); err != nil {
		log.Errorf("in-memory dataset db failed: %v", err)
		DB = nil
	}
}

// isSQLiteFile reports whether path begins with the SQLite file header.
//
// The point is to tell "a database I cannot open right now" from "not a
// database", because only the second is safe to delete. Reading the first
// sixteen bytes is enough: the header is a fixed string and every SQLite file
// starts with it.
func isSQLiteFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	var hdr [16]byte
	if _, err := io.ReadFull(f, hdr[:]); err != nil {
		return false
	}
	return string(hdr[:]) == "SQLite format 3\x00"
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

func Tokens(name string) []string {
	if DB == nil {
		return nil
	}
	var ds Dataset
	if err := DB.Where("name = ?", name).First(&ds).Error; err != nil {
		return nil
	}
	var rows []Vocabulary
	if err := DB.Where("dataset = ?", ds.ID).Order("id").Find(&rows).Error; err != nil {
		return nil
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.Token)
	}
	return out
}
