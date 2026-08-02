// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package json renders a report as one document.
//
// Written at the end, because a JSON document has no valid prefix — which is
// exactly why ndjson exists beside it.
package json

import (
	"encoding/json"
	"io"

	"github.com/rangertaha/urlinsane/internal/plugins/report"
)

// renderJSON writes the whole report as one object. Field order follows the
// struct, and every slice is already canonically sorted, so two identical scans
// produce byte-identical bytes.
func Render(w io.Writer, r report.Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(r)
}
