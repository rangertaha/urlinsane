// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"os"

	"github.com/rangertaha/urlinsane/internal/atomicfile"
)

// writeAtomic is atomicfile.Write, kept as a package-local name because every
// durable write in this package already calls it.
//
// The rationale lives on atomicfile: the same partial-write bug was written
// three times across this module, and each time the truncated file it left was
// then trusted by whatever check came next.
func writeAtomic(path string, data []byte, perm os.FileMode) error {
	return atomicfile.Write(path, data, perm)
}
