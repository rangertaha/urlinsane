// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package dataset

import (
	"fmt"
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

type Vocabulary struct {
	ID       uint
	Language uint // language ID
	Dataset  uint // dataset ID
	Token    string
}

type Transition struct {
	ID          uint
	Src         uint // vocab ID
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
