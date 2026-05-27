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
package collectors

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/rangertaha/urlinsane/internal/dataset"
	"github.com/rangertaha/urlinsane/internal/db"
)

// CheckSources looks up the sources of the given type (e.g. "package",
// "username") and reports which ones the name exists at (HTTP 2xx). Each
// source's Template is a URL with %s replaced by name.
func CheckSources(ctx context.Context, sourceType, name string) []*db.Hit {
	if dataset.DB == nil {
		return nil
	}
	var sources []dataset.Source
	if err := dataset.DB.Where("type = ?", sourceType).Find(&sources).Error; err != nil {
		return nil
	}

	var hits []*db.Hit
	for _, s := range sources {
		if ctx.Err() != nil {
			break
		}
		// Test existence against Check (an API that 404s cleanly) when set,
		// but record the human display URL in the hit.
		check := s.CheckURL
		if check == "" {
			check = s.Template
		}
		if sourceExists(ctx, strings.ReplaceAll(check, "%s", name)) {
			hits = append(hits, &db.Hit{Service: s.Code, URL: strings.ReplaceAll(s.Template, "%s", name)})
		}
	}
	return hits
}

// sourceExists reports whether a GET to url returns a 2xx status, bounding each
// request to 8s (or the parent context deadline, whichever is sooner).
func sourceExists(ctx context.Context, url string) bool {
	reqCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "urlinsane")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
