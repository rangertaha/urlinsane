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
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
)

// func Levenshtein(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
// 	out := make(chan internal.Domain)
// 	go func() {
// 		for domain := range in {
// 			// Set Levenshtein distance
// 			//   https://en.wikipedia.org/wiki/Levenshtein_distance
// 			dist := fuzzy.Levenshtein(domain.String(), c.Target())

// 			if dist <= c.Distance() {
// 				domain.Ld(dist)
// 				out <- domain
// 			}
// 		}
// 		close(out)
// 	}()
// 	return out
// }

// func Dedup(in <-chan internal.Domain, c *config.Config, total *int64) <-chan internal.Domain {
// 	out := make(chan internal.Domain)
// 	variants := make(map[string]bool)
// 	go func() {
// 		for domain := range in {
// 			if _, ok := variants[domain.String()]; !ok {
// 				variants[domain.String()] = true
// 				out <- domain
// 				*total++
// 			}
// 		}
// 		close(out)
// 	}()
// 	return out
// }

// func ReadCache(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
// 	out := make(chan internal.Domain)

// 	go func() {
// 		for domain := range in {
// 			if data, _ := c.Database().Read(domain.String()); data != "" {
// 				domain.Json(data)
// 			}
// 			out <- domain
// 		}
// 		close(out)
// 	}()

//		return out
//	}
func ReadCacheFilter() func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
	return func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
		out := make(chan internal.Domain)
		go func() {
			for domain := range in {
				if data, _ := c.Database().Read(domain.String()); data != "" {
					domain.Json(data)
				}
				out <- domain
			}
			close(out)
		}()
		return out
	}
}

func ExampleFilter() func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
	return func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
		out := make(chan internal.Domain)
		go func() {
			for domain := range in {

				out <- domain

			}
			close(out)
		}()
		return out
	}
}

func LevenshteinFilter() func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
	return func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
		out := make(chan internal.Domain)
		go func() {
			for domain := range in {
				// Set Levenshtein distance
				//   https://en.wikipedia.org/wiki/Levenshtein_distance
				dist := fuzzy.Levenshtein(domain.String(), c.Target())

				if dist <= c.Distance() {
					domain.Ld(dist)
					out <- domain
				}
			}
			close(out)
		}()
		return out
	}
}

func DedupFilter() func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
	return func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
		out := make(chan internal.Domain)
		variants := make(map[string]bool)
		go func() {
			for domain := range in {
				if _, ok := variants[domain.String()]; !ok {
					variants[domain.String()] = true
					out <- domain
				}
			}
			close(out)
		}()
		return out
	}
}

func GetTotal() func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
	return func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
		out := make(chan internal.Domain)

		go func() {
			var count int
			for domain := range in {
				count++

				out <- domain

			}
			c.Count(count)
			close(out)
		}()
		return out
	}
}

func RegexFilter() func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
	return func(in <-chan internal.Domain, c *config.Config) <-chan internal.Domain {
		out := make(chan internal.Domain)
		go func() {
			for domain := range in {
				if match, _ := regexp.MatchString(c.Regex(), domain.String()); match {
					out <- domain
				}

			}
			close(out)
		}()
		return out
	}
}
