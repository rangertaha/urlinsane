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
package engine

import (
	"regexp"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
)

// ExampleFilter is a no-op pass-through filter, kept as a template for writing
// new constraint filters.
func ExampleFilter() func(in <-chan *db.Domain, c internal.Config) <-chan *db.Domain {
	return func(in <-chan *db.Domain, c internal.Config) <-chan *db.Domain {
		out := make(chan *db.Domain)
		go func() {
			for domain := range in {

				out <- domain

			}
			close(out)
		}()
		return out
	}
}

// Levenshtein keeps only variants whose edit distance from the target is within
// the configured maximum (--distance).
func Levenshtein() func(in <-chan *db.Domain, c internal.Config) <-chan *db.Domain {
	return func(in <-chan *db.Domain, c internal.Config) <-chan *db.Domain {
		out := make(chan *db.Domain)
		go func() {
			for domain := range in {
				// Keep variants within the maximum Levenshtein distance.
				//   https://en.wikipedia.org/wiki/Levenshtein_distance
				if domain.Levenshtein <= c.Distance() {
					out <- domain
				}
			}
			close(out)
		}()
		return out
	}
}

// Dedup drops variants whose name has already been seen.
func Dedup() func(in <-chan *db.Domain, c internal.Config) <-chan *db.Domain {
	return func(in <-chan *db.Domain, c internal.Config) <-chan *db.Domain {
		out := make(chan *db.Domain)
		variants := make(map[string]bool)
		go func() {
			for domain := range in {
				if _, ok := variants[domain.Name]; !ok {
					variants[domain.Name] = true
					out <- domain
				}
			}
			close(out)
		}()
		return out
	}
}

// Regex keeps only variants whose name matches the configured pattern (--regex).
func Regex() func(in <-chan *db.Domain, c internal.Config) <-chan *db.Domain {
	return func(in <-chan *db.Domain, c internal.Config) <-chan *db.Domain {
		out := make(chan *db.Domain)
		go func() {
			for domain := range in {
				if match, _ := regexp.MatchString(c.Regex(), domain.Name); match {
					out <- domain
				}
			}
			close(out)
		}()
		return out
	}
}
