// Copyright © 2018 rangertaha rangertaha@gmail.com
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

import "strings"

type (
	Language struct {
		Code         string
		Name         string
		Numerals     map[string][]string
		Graphemes    []string
		Misspellings [][]string
		Homophones   [][]string
		Homoglyphs   map[string][]string
	}

	Keyboard struct {
		Code     string
		Name     string
		Language Language
		Layout   []string
	}
)

// Languages
var LANGUAGES = map[string]Language{
	"en": ENGLISH,
	"ar": ARABIC,
}

var KEYBOARDS = map[string]Keyboard{}

// GetLanguages looks up returns Language types.
func GetLanguages(codes []string) (lgs []Language) {
	for _, name := range codes {
		lang, ok := LANGUAGES[strings.ToLower(name)]
		if ok {
			lgs = append(lgs, lang)
		}

	}
	return
}

// GetKeyboards looks up returns Language types.
func GetKeyboards(names []string) (kbs []Keyboard) {
	for _, name := range names {
		if strings.ToLower(name) == "all" {
			for _, kb := range KEYBOARDS {
				kbs = append(kbs, kb)
			}
		} else {
			keyboard, ok := KEYBOARDS[strings.ToLower(name)]
			if ok {
				kbs = append(kbs, keyboard)
			}
		}
	}
	return
}

// KBRegister adds keyboards to a registry
func KBRegister(keyboards []Keyboard) {
	for _, board := range keyboards {
		KEYBOARDS[strings.ToLower(board.Code)] = board
	}
}

// Adjacent
func (urli *Keyboard) Adjacent(char string) (chars []string) {
	for r, row := range urli.Layout {
		for c, _ := range row {
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

// Similar
func (lang *Language) Similar(key string) (chars []string) {
	char, ok := lang.Homoglyphs[key]
	if ok {
		chars = append(chars, char...)
	}
	return chars
}
