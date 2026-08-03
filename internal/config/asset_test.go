// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mangle is what happened to the shipped assets: every byte >= 0x80 replaced by
// the UTF-8 replacement character, the signature of a binary file round tripped
// through a text decoder.
func mangle(b []byte) []byte {
	var out []byte
	for _, c := range b {
		if c >= 0x80 {
			out = append(out, 0xef, 0xbf, 0xbd)
			continue
		}
		out = append(out, c)
	}
	return out
}

var (
	goodGzip   = []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03}
	goodSQLite = append([]byte("SQLite format 3\x00"), make([]byte, 16)...)
)

// A mangled asset must not be written. Writing it and reporting "extracted" is
// what shipped: tens of megabytes of unusable bytes landed in the user's config
// directory, the operator that needed them failed later and somewhere else, and
// deleting the file re-extracted the same corruption.
func TestExtractRefusesMangledBytes(t *testing.T) {
	for _, tc := range []struct {
		name  string
		data  []byte
		valid validator
		want  string
	}{
		{"maxmind.db.gz", mangle(goodGzip), isGzip, "not gzip data"},
		// SQLite's magic is ASCII and survives the mangling (see
		// TestManglingIsInvisibleInAnASCIIHeader), so the case that reaches
		// this validator is a wrong or damaged header rather than that one.
		{"dataset.db", []byte("SQLite format 2\x00................"), isSQLite, "not a SQLite database"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			f := extract(dir, tc.name, tc.data, tc.valid)

			if f.Err == nil {
				t.Fatalf("extract accepted mangled %s", tc.name)
			}
			if !strings.Contains(f.Err.Error(), tc.want) {
				t.Errorf("Err = %v, want it to mention %q", f.Err, tc.want)
			}
			if !strings.Contains(f.Err.Error(), "internal/config/"+tc.name) {
				t.Errorf("Err = %v, want it to say which file to replace", f.Err)
			}
			if f.Written {
				t.Error("Written is true for a file that was refused")
			}
			if _, err := os.Stat(filepath.Join(dir, tc.name)); err == nil {
				t.Error("the refused bytes were written to disk anyway")
			}
		})
	}
}

// The mangling is not hypothetical: it destroys a gzip magic specifically,
// because 0x8b is exactly the kind of byte a text decoder cannot represent.
func TestMangleDestroysTheGzipMagic(t *testing.T) {
	if err := isGzip(goodGzip); err != nil {
		t.Fatalf("isGzip rejected valid gzip bytes: %v", err)
	}
	bad := mangle(goodGzip)
	if bad[0] != 0x1f || bad[1] != 0xef {
		t.Fatalf("mangle did not reproduce the observed corruption: % x", bad[:4])
	}
	if err := isGzip(bad); err == nil {
		t.Error("isGzip accepted the corruption that shipped")
	}
}

// What a header check can and cannot catch, recorded so nobody assumes it
// catches more.
//
// gzip's magic contains 0x8b, so a text round trip always destroys it and the
// corruption is detectable from three bytes. SQLite's magic is the ASCII string
// "SQLite format 3", which the same round trip leaves untouched — a mangled
// database would keep a perfect header and be ruined from the first page
// onward. dataset.db happens to be intact; if it ever is not, this check will
// not be what tells us.
func TestManglingIsInvisibleInAnASCIIHeader(t *testing.T) {
	if err := isSQLite(mangle(goodSQLite)); err != nil {
		t.Fatalf("the premise changed: mangling now alters the SQLite header (%v)", err)
	}
	if err := isGzip(mangle(goodGzip)); err == nil {
		t.Fatal("the premise changed: mangling no longer destroys the gzip magic")
	}
}

// Well-formed bytes still extract, and only once.
func TestExtractWritesValidBytesOnce(t *testing.T) {
	dir := t.TempDir()

	f := extract(dir, "dataset.db", goodSQLite, isSQLite)
	if f.Err != nil {
		t.Fatalf("extract refused valid bytes: %v", f.Err)
	}
	if !f.Written {
		t.Error("Written is false for a file that was written")
	}

	// A second call must leave the existing file alone: the user may have
	// replaced it with a better one.
	if err := os.WriteFile(f.Path, goodSQLite, 0o640); err != nil {
		t.Fatal(err)
	}
	again := extract(dir, "dataset.db", goodSQLite, isSQLite)
	if again.Written {
		t.Error("an existing file was overwritten")
	}
	if again.Err != nil {
		t.Errorf("Err = %v, want nil for a file already in place", again.Err)
	}
}

// A truncated or empty asset is refused with a reason rather than written.
func TestExtractRefusesTruncatedAndAbsentAssets(t *testing.T) {
	dir := t.TempDir()

	if f := extract(dir, "a.gz", []byte{0x1f}, isGzip); f.Err == nil ||
		!strings.Contains(f.Err.Error(), "truncated") {
		t.Errorf("truncated gzip: Err = %v, want a truncation error", f.Err)
	}
	if f := extract(dir, "b.gz", nil, isGzip); f.Err == nil ||
		!strings.Contains(f.Err.Error(), "not compiled into this binary") {
		t.Errorf("absent asset: Err = %v, want the not-compiled-in error", f.Err)
	}
}

// The class test: every shipped asset either validates or is reported.
//
// It runs the real validators over the real embedded bytes, so a corrupt asset
// cannot reach a release quietly. It deliberately does not require the assets to
// be valid — maxmind.db.gz is not, and has not been for as long as it has been
// in the repository — because the invariant that matters to a user is weaker and
// more useful: whatever ships, the tool says so rather than pretending.
//
// If this fails, an asset went bad AND Init stayed silent about it. Both halves
// have to break for the tool to mislead anyone.
func TestShippedAssetsValidateOrAreReported(t *testing.T) {
	for _, a := range []struct {
		name  string
		data  []byte
		valid validator
		file  func(Setup) File
	}{
		{DatasetDB, datasetDB, isSQLite, func(s Setup) File { return s.Dataset }},
		{MaxMindDB, maxmindDB, isGzip, func(s Setup) File { return s.GeoIP }},
	} {
		t.Run(a.name, func(t *testing.T) {
			dir := t.TempDir()
			f := extract(dir, a.name, a.data, a.valid)

			switch err := a.valid(a.data); {
			case err == nil:
				if f.Err != nil {
					t.Errorf("%s is well formed but was refused: %v", a.name, f.Err)
				}
			default:
				t.Logf("%s does not validate: %v", a.name, err)
				if f.Err == nil {
					t.Errorf("%s is corrupt and extract reported success", a.name)
				}
				// And the run has to be able to tell the user.
				s := Setup{}
				switch a.name {
				case DatasetDB:
					s.Dataset = f
				case MaxMindDB:
					s.GeoIP = f
				}
				if !s.FirstRun() {
					t.Errorf("%s is corrupt but the setup reports nothing worth showing", a.name)
				}
			}
		})
	}
}
