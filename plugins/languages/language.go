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
package languages

type (
// Language type
// Language struct {
// 	code         string
// 	name         string
// 	description  string
// 	numerals     map[string][]string
// 	graphemes    []string
// 	vowels       []string
// 	misspellings [][]string
// 	homophones   [][]string
// 	antonyms     map[string][]string
// 	homoglyphs   map[string][]string
// 	keyboards    []Keyboard
// }

// Keyboard type
// Keyboard struct {
// 	code        string
// 	name        string
// 	description string
// 	layout      []string
// }
// // KeyboardGroup type
// KeyboardGroup struct {
// 	code        string   `json:"code,omitempty"`
// 	Keyboards   []string `json:"keyboards,omitempty"`
// 	description string   `json:"description,omitempty"`
// }

// // KeyboardRegistry stores registered keyboards and groups
//
//	KeyboardRegistry struct {
//		registry map[string][]Keyboard
//	}
)

// func (k *Keyboard) Id() string {
// 	return k.code
// }
// func (k *Keyboard) Name() string {
// 	return k.name
// }
// func (k *Keyboard) Description() string {
// 	return k.description
// }
// func (k *Keyboard) Layouts() []string {
// 	return k.layout
// }

// func (l *Language) Id() string {
// 	return l.code
// }
// func (l *Language) Name() string {
// 	return l.name
// }

// // Numerals in the broadest sense is a word or phrase that
// // describes a numerical quantity.
// func (l *Language) Numerals() map[string][]string {
// 	return l.numerals
// }

// // Graphemes is the smallest functional unit of a writing system.
// func (l *Language) Graphemes() []string {
// 	return l.graphemes
// }

// // Vowels are syllabic speech sound pronounced without any stricture in the vocal tract.
// func (l *Language) Vowels() []string {
// 	return l.vowels
// }

// func (l *Language) Misspellings() [][]string {
// 	return l.misspellings
// }

// func (l *Language) Homophones() [][]string {
// 	return l.homophones
// }

// func (l *Language) Antonyms() map[string][]string {
// 	return l.antonyms
// }

// func (l *Language) Homoglyphs() map[string][]string {
// 	return l.homoglyphs
// }

// func (l *Language) Keyboards() []Keyboard {
// 	return l.keyboards
// }

// func (k *Keyboard) Id() string {
// 	return k.code
// }
// func (k *Keyboard) Name() string {
// 	return k.name
// }
// func (k *Keyboard) Description() string {
// 	return k.description
// }
// func (k *Keyboard) Layout() []string {
// 	return k.layout
// }

// // KEYBOARDS stores all the registered keyboards
// var KEYBOARDS = NewKeyboardRegistry()

// // NewKeyboardRegistry returns a new KeyboardRegistry
// func NewKeyboardRegistry() KeyboardRegistry {
// 	return KeyboardRegistry{
// 		registry: make(map[string][]Keyboard),
// 	}
// }

// // Add allows you to add keyboards to the registry
// func (kb *KeyboardRegistry) Add(keyboards []Keyboard) {
// 	for _, board := range keyboards {
// 		kb.registry[strings.ToUpper(board.code)] = []Keyboard{board}
// 	}
// }

// // Append allows you to append keyboards to a group name
// func (kb *KeyboardRegistry) Append(name string, keyboards []Keyboard) {
// 	key := strings.ToUpper(name)
// 	kbs, ok := kb.registry[key]
// 	if ok {
// 		for _, board := range keyboards {
// 			kbs = append(kbs, board)
// 		}
// 		kb.registry[key] = kbs
// 	} else {
// 		kb.registry[key] = keyboards
// 	}
// }

// // Keyboards looks up and returns Keyboards.
// func (kb *KeyboardRegistry) Keyboards(names ...string) (kbs []Keyboard) {
// 	for _, name := range names {
// 		keyboards, ok := kb.registry[strings.ToUpper(name)]
// 		if ok {
// 			for _, keyboard := range keyboards {
// 				kbs = append(kbs, keyboard)
// 			}
// 		}
// 	}
// 	return
// }

// // Adjacent returns adjacent characters on the given keyboard
// func (k *Keyboard) Adjacent(char string) (chars []string) {
// 	for r, row := range k.layout {
// 		for c := range row {
// 			var top, bottom, left, right string
// 			if char == string(k.layout[r][c]) {
// 				if r > 0 {
// 					top = string(k.layout[r-1][c])
// 					if top != " " {
// 						chars = append(chars, top)
// 					}
// 				}
// 				if r < len(k.layout)-1 {
// 					bottom = string(k.layout[r+1][c])
// 					if bottom != " " {
// 						chars = append(chars, bottom)
// 					}
// 				}
// 				if c > 0 {
// 					left = string(k.layout[r][c-1])
// 					if left != " " {
// 						chars = append(chars, left)
// 					}
// 				}
// 				if c < len(row)-1 {
// 					right = string(k.layout[r][c+1])
// 					if right != " " {
// 						chars = append(chars, right)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return chars
// }

// // SimilarChars ...
// func (l *Language) SimilarChars(key string) (chars []string) {
// 	char, ok := l.homoglyphs[key]
// 	if ok {
// 		chars = append(chars, char...)
// 	}
// 	return
// }

// // SimilarSpellings ...
// func (l *Language) SimilarSpellings(str string) (words []string) {
// 	for _, wordset := range l.misspellings {
// 		for _, word := range wordset {
// 			if strings.Contains(str, word) {
// 				for _, w := range wordset {
// 					if w != word {
// 						words = append(words, strings.Replace(str, word, w, -1))
// 					}
// 				}

// 			}
// 		}
// 	}
// 	return
// }

// // SimilarSounds ...
// func (l *Language) SimilarSounds(str string) (words []string) {
// 	for _, wordset := range l.homophones {
// 		for _, word := range wordset {
// 			if strings.Contains(str, word) {
// 				for _, w := range wordset {
// 					if w != word {
// 						words = append(words, strings.Replace(str, word, w, -1))
// 					}
// 				}

// 			}
// 		}
// 	}
// 	return
// }

// var LANGUAGES = map[string]Language{}

// func Add(name string, lang Language) {
// 	LANGUAGES[name] = lang
// }

// func Get(name string) (lang Language, err error) {
// 	if lang, ok := LANGUAGES[name]; ok {
// 		return lang, err
// 	}

// 	return lang, fmt.Errorf("unable to locate outputs/%s plugin", name)
// }

// func All() (langs []Language) {
// 	for _, lang := range LANGUAGES {
// 		langs = append(langs, lang)
// 	}

// 	return
// }

// func List(IDs ...string) (langs []Language) {
// 	for id, lang := range LANGUAGES {
// 		for _, lid := range IDs {
// 			if id == lid {
// 				langs = append(langs, lang)
// 			}
// 		}
// 	}
// 	for _, lid := range IDs {
// 		if lid == "all" {
// 			IDs = []string{}
// 		}
// 	}

// 	if len(IDs) == 0 {
// 		for _, lang := range LANGUAGES {
// 			langs = append(langs, lang)
// 		}
// 	}

// 	return
// }

// func Keyboards(IDs ...string) (keyboards []Keyboard) {
// 	for id, lang := range LANGUAGES {
// 		for _, lid := range IDs {
// 			if id == lid {
// 				keyboards = append(keyboards, lang.keyboards...)
// 			}
// 		}
// 	}
// 	for _, kid := range IDs {
// 		if kid == "all" {
// 			IDs = []string{}
// 		}
// 	}
// 	if len(IDs) == 0 {
// 		for _, lang := range LANGUAGES {
// 			keyboards = append(keyboards, lang.keyboards...)
// 		}
// 	}

// 	return
// }
