// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"github.com/rangertaha/urlinsane/internal/dataset"
	"github.com/rangertaha/urlinsane/pkg/kb"
)

// Build creates the default dataset database from scratch: schema, languages,
// vocabulary, transitions and sources.
//
// It deletes any existing file first. That is the point of a build rather than
// an import: migrating a database through three schema generations leaves the
// columns of all three, and the shipped database had exactly that — `template`,
// `check` and `check_url` beside the `url`, `success` and `failed` that
// replaced them, with the new columns empty because AutoMigrate can add a
// column but cannot fill it. A build has one schema by construction.
func Build(dbPath, root string) error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o750); err != nil {
		return err
	}
	if err := os.Remove(dbPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	dataset.Config(dbPath)
	if dataset.DB == nil {
		return fmt.Errorf("gen: could not open %s", dbPath)
	}

	if err := Languages(); err != nil {
		return err
	}
	if err := All(root); err != nil {
		return err
	}
	return Sources(filepath.Join(root, "sources"))
}

// Languages seeds the Language table from the keyboard catalogue.
//
// pkg/kb is the authority rather than the dataset directories, because a
// language is a thing the tool can reason about — it has a keyboard, so it has
// adjacency — whether or not anyone has curated a word list for it yet. Seeding
// from datasets/languages/ would have made the table a mirror of whichever
// directories exist, so a language would blink out of --list languages the
// moment its corpus was empty.
//
// The two sets overlap but neither contains the other: kb ships 110 languages,
// the repo curates 30, and a curated language with no layout still gets a row
// from the import that follows.
//
// Three curated directories use a code kb does not: iw (kb has he), no (kb has
// nb) and la (kb has no Latin layout at all). Each therefore lands as a second
// row for a language kb already listed, so homoglyphs would be looked up under
// one code while keyboard adjacency uses the other. Reconciling them is a
// dataset change rather than a code one, so this records the split instead of
// silently folding the codes together.
func Languages() error {
	names := display.English.Languages()
	for _, code := range kb.Languages() {
		name := code
		if tag, err := language.Parse(code); err == nil {
			if n := names.Name(tag); n != "" {
				name = n
			}
		}
		row := dataset.Language{Code: code}
		if err := dataset.DB.Where(dataset.Language{Code: code}).
			Assign(dataset.Language{Name: name}).
			FirstOrCreate(&row).Error; err != nil {
			return fmt.Errorf("gen: language %s: %w", code, err)
		}
	}
	return nil
}

// Sources loads the platform lists into the Source table: the registries,
// forges and providers a name is checked against.
//
// Three shapes, one per file, because the kinds genuinely differ. A package
// registry has a page for humans and an API that answers existence cleanly;
// a mail provider is a bare domain with neither.
//
//	packages.lst   code  page-url  check-url
//	repos.lst      code  page-url  check-url
//	usernames.lst  code  url
//	email.lst      domain
func Sources(dir string) error {
	kinds := []struct{ file, kind string }{
		{"packages.lst", "package"},
		{"repos.lst", "repository"},
		{"usernames.lst", "username"},
		{"email.lst", "email"},
	}
	for _, k := range kinds {
		path := filepath.Join(dir, k.file)
		lines, err := Extract(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		// Replace rather than merge: a source removed from the list must
		// disappear from the database, or a retired registry keeps being
		// probed forever.
		dataset.DB.Where("type = ?", k.kind).Delete(&dataset.Source{})

		rows := make([]*dataset.Source, 0, len(lines))
		for _, f := range lines {
			if len(f) == 0 || strings.HasPrefix(f[0], "#") {
				continue
			}
			row := &dataset.Source{Type: k.kind, Code: f[0]}
			switch {
			case len(f) >= 3:
				// The check URL is what the prober calls, so it is the one
				// that goes in URL. The page URL has no column in this schema.
				row.URL = f[2]
			case len(f) == 2:
				row.URL = f[1]
			default:
				// A bare domain: the code is the whole record.
				row.URL = f[0]
			}
			rows = append(rows, row)
		}
		if len(rows) == 0 {
			continue
		}
		if err := dataset.DB.Create(&rows).Error; err != nil {
			return fmt.Errorf("gen: sources %s: %w", k.file, err)
		}
		if Progress != nil {
			Progress(k.kind+" sources", path, len(rows), 0)
		}
	}
	return nil
}
