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
	"fmt"
	"testing"
)

// TestLevenshtein ...
func TestLevenshtein(t *testing.T) {
	var str1 = "Asheville"
	var str2 = "Arizona"
	fmt.Println("Distance between Asheville and Arizona:", Levenshtein(str1, str2))

	str1 = "Python"
	str2 = "Peithen"
	fmt.Println("Distance between Python and Peithen:", Levenshtein(str1, str2))

	str1 = "Orange"
	str2 = "Apple"
	fmt.Println("Distance between Orange and Apple:", Levenshtein(str1, str2))
}
