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
package pkg

import "unicode/utf8"


const minLengthThreshold = 32

// ComputeDistance computes the levenshtein distance between the two
// strings passed as an argument. The return value is the levenshtein distance
func Levenshtein(a, b string) int {
	if len(a) == 0 {
		return utf8.RuneCountInString(b)
	}

	if len(b) == 0 {
		return utf8.RuneCountInString(a)
	}

	if a == b {
		return 0
	}


	s1 := []rune(a)
	s2 := []rune(b)

	// swap to save some memory O(min(a,b)) instead of O(a)
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}

	// remove trailing identical runes.
	for i := 0; i < len(s1); i++ {
		if s1[len(s1)-1-i] != s2[len(s2)-1-i] {
			s1 = s1[:len(s1)-i]
			s2 = s2[:len(s2)-i]
			break
		}
	}

	// Remove leading identical runes.
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			s1 = s1[i:]
			s2 = s2[i:]
			break
		}
	}

	lenS1 := len(s1)
	lenS2 := len(s2)

	// Init the row.
	var x []uint16
	if lenS1+1 > minLengthThreshold {
		x = make([]uint16, lenS1+1)
	} else {
		// We make a small optimization here for small strings.
		// Because a slice of constant length is effectively an array,
		// it does not allocate. So we can re-slice it to the right length
		// as long as it is below a desired threshold.
		x = make([]uint16, minLengthThreshold)
		x = x[:lenS1+1]
	}

	// we start from 1 because index 0 is already 0.
	for i := 1; i < len(x); i++ {
		x[i] = uint16(i)
	}

	// make a dummy bounds check to prevent the 2 bounds check down below.
	// The one inside the loop is particularly costly.
	_ = x[lenS1]
	// fill in the rest
	for i := 1; i <= lenS2; i++ {
		prev := uint16(i)
		for j := 1; j <= lenS1; j++ {
			current := x[j-1] // match
			if s2[i-1] != s1[j-1] {
				current = min(x[j-1]+1, prev+1, x[j]+1)
			}
			x[j-1] = prev
			prev = current
		}
		x[lenS1] = prev
	}
	return int(x[lenS1])
}
