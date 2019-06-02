// The MIT License (MIT)
//
// Copyright © 2019 Tal Hachi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package languages

import (
	"strings"
)

type (
	// Language type
	Language struct {
		Code         string              `json:"code,omitempty"`
		Name         string              `json:"name,omitempty"`
		Numerals     map[string][]string `json:"numerals,omitempty"`
		Graphemes    []string            `json:"graphemes,omitempty"`
		Vowels       []string            `json:"vowels,omitempty"`
		Misspellings [][]string          `json:"misspellings,omitempty"`
		Homophones   [][]string          `json:"homophones,omitempty"`
		Antonyms     map[string][]string `json:"antonyms,omitempty"`
		Homoglyphs   map[string][]string `json:"homoglyphs,omitempty"`
	}

	// Keyboard type
	Keyboard struct {
		Code        string   `json:"code,omitempty"`
		Name        string   `json:"name,omitempty"`
		Description string   `json:"description,omitempty"`
		Language    Language `json:"language,omitempty"`
		Layout      []string `json:"layout,omitempty"`
	}
	// KeyboardGroup type
	KeyboardGroup struct {
		Code        string   `json:"code,omitempty"`
		Keyboards   []string `json:"keyboards,omitempty"`
		Description string   `json:"description,omitempty"`
	}

	// KeyboardRegistry stores registered keyboards and groups
	KeyboardRegistry struct {
		registry map[string][]Keyboard
	}
)

// KEYBOARDS stores all the registered keyboards
var KEYBOARDS = NewKeyboardRegistry()

// NewKeyboardRegistry returns a new KeyboardRegistry
func NewKeyboardRegistry() KeyboardRegistry {
	return KeyboardRegistry{
		registry: make(map[string][]Keyboard),
	}
}

// Add allows you to add keyboards to the registry
func (kb *KeyboardRegistry) Add(keyboards []Keyboard) {
	for _, board := range keyboards {
		kb.registry[strings.ToUpper(board.Code)] = []Keyboard{board}
	}
}

// Append allows you to append keyboards to a group name
func (kb *KeyboardRegistry) Append(name string, keyboards []Keyboard) {
	key := strings.ToUpper(name)
	kbs, ok := kb.registry[key]
	if ok {
		for _, board := range keyboards {
			kbs = append(kbs, board)
		}
		kb.registry[key] = kbs
	} else {
		kb.registry[key] = keyboards
	}
}

// Keyboards looks up and returns Keyboards.
func (kb *KeyboardRegistry) Keyboards(names ...string) (kbs []Keyboard) {
	for _, name := range names {
		keyboards, ok := kb.registry[strings.ToUpper(name)]
		if ok {
			for _, keyboard := range keyboards {
				kbs = append(kbs, keyboard)
			}
		}
	}
	return
}

// Adjacent returns adjacent characters on the given keyboard
func (urli *Keyboard) Adjacent(char string) (chars []string) {
	for r, row := range urli.Layout {
		for c := range row {
			var top, bottom, left, right string
			if char == string(urli.Layout[r][c]) {
				if r > 0 {
					top = string(urli.Layout[r-1][c])
					if top != " " {
						chars = append(chars, top)
					}
				}
				if r < len(urli.Layout)-1 {
					bottom = string(urli.Layout[r+1][c])
					if bottom != " " {
						chars = append(chars, bottom)
					}
				}
				if c > 0 {
					left = string(urli.Layout[r][c-1])
					if left != " " {
						chars = append(chars, left)
					}
				}
				if c < len(row)-1 {
					right = string(urli.Layout[r][c+1])
					if right != " " {
						chars = append(chars, right)
					}
				}
			}
		}
	}
	return chars
}

// SimilarChars ...
func (lang *Language) SimilarChars(key string) (chars []string) {
	char, ok := lang.Homoglyphs[key]
	if ok {
		chars = append(chars, char...)
	}
	return chars
}

// SimilarSpellings ...
func (lang *Language) SimilarSpellings(str string) (words []string) {
	for _, wordset := range lang.Misspellings {
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

// SimilarSounds ...
func (lang *Language) SimilarSounds(str string) (words []string) {
	for _, wordset := range lang.Homophones {
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
