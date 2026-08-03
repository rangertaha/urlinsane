package config

import (
	"os"
	"path/filepath"
	"testing"
)

// The corruption fingerprint must be exactly four bytes long, no more and no
// less. A prefix check on a short read is the obvious way to get this wrong:
// a one-byte file must not match, and a file that is exactly the fingerprint
// must.
func TestShippedCorruptionBoundaries(t *testing.T) {
	for _, tc := range []struct {
		name string
		b    []byte
		want bool
	}{
		{"nil", nil, false},
		{"empty", []byte{}, false},
		{"one byte", []byte{0x1f}, false},
		{"three bytes", []byte{0x1f, 0xef, 0xbf}, false},
		{"exactly four", []byte{0x1f, 0xef, 0xbf, 0xbd}, true},
		{"four plus", []byte{0x1f, 0xef, 0xbf, 0xbd, 0x08}, true},
		{"real gzip", goodGzip, false},
		{"sqlite", goodSQLite, false},
	} {
		if got := isShippedCorruption(tc.b); got != tc.want {
			t.Errorf("%s: isShippedCorruption(% x) = %v, want %v", tc.name, tc.b, got, tc.want)
		}
	}
}

// A truncated file must not crash, and must not be mistaken for either case.
func TestOptionalReportsTruncatedFiles(t *testing.T) {
	for _, b := range [][]byte{nil, {0x1f}, {0x1f, 0xef}, {0x1f, 0xef, 0xbf}} {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, MaxMindDB), b, 0o640); err != nil {
			t.Fatal(err)
		}
		f := optional(dir, MaxMindDB, isGzip)
		// Present with an error: a file that is there and cannot work.
		if !f.Present {
			t.Errorf("% x: Present is false for a file that exists", b)
		}
		if f.Err == nil {
			t.Errorf("% x: a truncated database was accepted", b)
		}
	}
}
