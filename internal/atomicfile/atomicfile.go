// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package atomicfile publishes a file in one step, so a reader sees either the
// old contents or the new ones and never half of either.
//
// It exists because the same bug was written three times, in three packages,
// and each time the file that was left behind was then trusted:
//
//   - internal/store's blockstore used os.WriteFile for a content-addressed
//     block. A crash mid-write left a short file at the block's path, and
//     because the path existed Has reported the block stored and every later
//     Put skipped it — bytes that no longer hash to their own name, for good.
//   - internal/store's index wrote a temporary and renamed it but never flushed
//     it, so after a power loss the rename could reach disk ahead of the data.
//     scans.json came back short or empty: every saved scan unreachable and no
//     new scan recordable until someone deleted it by hand.
//   - internal/config's first-run extraction wrote the embedded dataset and
//     GeoIP databases straight to their final paths, and used mere existence as
//     the "already extracted" test. A write that ran out of disk left a
//     truncated file that was never re-extracted and never re-checked.
//
// The shape is the same each time and so is the reason it goes unnoticed: the
// partial file is indistinguishable from a complete one to whatever check comes
// next. A durable write anywhere in this module calls Write rather than
// os.WriteFile.
package atomicfile

import (
	"os"
	"path/filepath"
)

// Write creates path with the given contents and permissions, atomically.
//
// The temporary is made in the target's own directory so the rename is a
// metadata operation on one filesystem rather than a copy across two. The Sync
// before it is what extends the guarantee from "survives a process crash" to
// "survives a power loss": without it the rename can be durable while the bytes
// it published are not.
func Write(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	// Cleans up on every failure path below; a no-op once the rename has moved
	// the file away.
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	// Close is checked: a write buffered by the kernel can still fail here, and
	// renaming a file that failed to close would publish the failure.
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), perm); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
