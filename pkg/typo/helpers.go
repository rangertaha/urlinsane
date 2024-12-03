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
package typo

import (
	"strings"
)

func characterDeletion(token string, character string) (tokens []string) {
	var nmap = map[string]bool{}

	for i, char := range token {
		if character == string(char) {
			nmap[token[:i]+token[i+1:]] = true
			// tokens = append(tokens, token[:i]+token[i+1:])
		}
	}
	nmap[strings.Replace(token, character, "", -1)] = true

	for n := range nmap {
		tokens = append(tokens, n)
	}

	return
}

func characterReplace(token string, character, replacement string) (tokens []string) {
	var nmap = map[string]bool{}

	for i, char := range token {
		if character == string(char) {
			nmap[token[:i]+token[i+1:]] = true
		}
	}
	nmap[strings.Replace(token, character, replacement, -1)] = true

	for n := range nmap {
		tokens = append(tokens, n)
	}

	return
}

// PrefixInsertion creates tokens by prepending each prefix from the given
// list to the specified token. Example:
// Inputs:
//
//	prefixes = ["www", "ftp", "shop"]
//	token = "example"
//
// Outputs: ["wwwexample", "ftpexample", "shopexample"]
func PrefixInsertion(token string, prefixes ...string) (tokens []string) {
	for _, prefix := range prefixes {
		tokens = append(tokens, prefix+token)
	}
	return
}

// SuffixInsertion creates tokens by appending each suffix from the provided
// list to the end (right side) of the given token. Example:
// Inputs:
//
//	suffixes = ["com", "net", "io"]
//	token = "example"
//
// Outputs: ["examplecom", "examplenet", "exampleio"]
func SuffixInsertion(token string, suffixes ...string) (tokens []string) {
	for _, suffix := range suffixes {
		tokens = append(tokens, token+suffix)
	}
	return
}

func numeralMap(data map[string][]string, pos int) (words map[string]string) {
	words = make(map[string]string)

	for num, tokens := range data {
		for i, token := range tokens {
			if i == pos {
				words[num] = token
				// words[token] = num
			}
		}
	}

	return
}

// Adjacent returns adjacent characters on a given keyboard
func adjacentCharacters(char string, layout ...string) (chars []string) {
	chars = []string{}
	for r, row := range layout {
		for c := range row {
			var top, bottom, left, right string
			if char == string(layout[r][c]) {
				if r > 0 {
					top = string(layout[r-1][c])
					if top != " " {
						chars = append(chars, top)
					}
				}
				if r < len(layout)-1 {
					bottom = string(layout[r+1][c])
					if bottom != " " {
						chars = append(chars, bottom)
					}
				}
				if c > 0 {
					left = string(layout[r][c-1])
					if left != " " {
						chars = append(chars, left)
					}
				}
				if c < len(row)-1 {
					right = string(layout[r][c+1])
					if right != " " {
						chars = append(chars, right)
					}
				}
			}
		}
	}
	return chars
}

// similarChars returns homoglyphs, characters that look alike from other languages
func similarChars(key string, data map[string][]string) (chars []string) {
	chars = []string{}
	char, ok := data[key]
	if ok {
		chars = append(chars, char...)
	}
	return chars
}

// similarSounds returns common homophones, words that sound alike
func similarSounds(str string, data ...[]string) (words []string) {
	words = []string{}
	for _, wordset := range data {
		for _, word := range wordset {
			if strings.Contains(str, word) {
				for _, w := range wordset {
					if w != word {
						words = append(words, strings.Replace(str, word, w, -1))
					}
				}

			}
		}
	}
	return
}
