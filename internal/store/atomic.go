// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"os"
	"path/filepath"
)

// writeAtomic publishes a file in one step: a temporary beside it, flushed to
// disk, then renamed over the target.
//
// Every durable write in this package goes through here. Two did not, and both
// were the same bug in different clothes:
//
//   - FSBlockstore.Put used os.WriteFile, which creates and truncates before it
//     writes. A crash mid-write left a short file at a block's path, and because
//     the path existed Has reported the block stored and every later Put skipped
//     it — in a content-addressed store, bytes that no longer hash to their own
//     name, permanently, with nothing to notice.
//   - Index.save wrote a temp and renamed but never flushed it, so its own
//     comment ("a crash mid-write must not leave an index that parses as neither
//     the old nor the new one") was not enforced. After a power loss the rename
//     could reach the disk ahead of the data, leaving scans.json short or empty:
//     every saved scan unreachable and no new scan recordable until someone
//     deleted the file by hand.
//
// The rename is what makes a reader see either the old file or the new one and
// never a half of either; the Sync before it is what makes that true across a
// power loss rather than only across a process crash. A new durable write in
// this package calls this rather than os.WriteFile.
func writeAtomic(path string, data []byte, perm os.FileMode) error {
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
