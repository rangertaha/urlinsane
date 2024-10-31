// Copyright (C) 2024 Rangertaha
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
package nlp

import (
	"github.com/arpitgogia/rake"
)

// Rake (Rapid Automatic Keyword Extraction) algorithm, is a domain independent
// keyword extraction algorithm which tries to determine key phrases in a body
// of text by analyzing the frequency of word appearance and its co-occurance
// with other words in the text.
func Rake(s string) map[string]float64 {
	return rake.WithText(s)
}

