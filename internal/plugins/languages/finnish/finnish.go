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
package finnish

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "fi"

type Finnish struct {
	code         string
	name         string
	description  string
	numerals     map[string][]string
	graphemes    []string
	vowels       []string
	misspellings [][]string
	homophones   [][]string
	antonyms     map[string][]string
	homoglyphs   map[string][]string
}

func (l *Finnish) Id() string {
	return l.code
}
func (l *Finnish) Name() string {
	return l.name
}
func (l *Finnish) Description() string {
	return l.description
}
func (l *Finnish) Numerals() map[string][]string {
	return l.numerals
}
func (l *Finnish) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *Finnish) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *Finnish) Graphemes() []string {
	return l.graphemes
}

func (l *Finnish) Vowels() []string {
	return l.vowels
}

func (l *Finnish) Misspellings() [][]string {
	return l.misspellings
}

func (l *Finnish) Homophones() [][]string {
	return l.homophones
}

func (l *Finnish) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *Finnish) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *Finnish) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *Finnish) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *Finnish) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *Finnish) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	// fiMisspellings are common misspellings
	fiMisspellings = [][]string{
		[]string{"", ""},
	}

	// fiHomophones are words that sound alike
	fiHomophones = [][]string{
		[]string{"piste", "."},
	}

	// fiAntonyms are words opposite in meaning to another (e.g. bad and good ).
	fiAntonyms = map[string][]string{
		"hyvä": []string{"huono"},
	}

	Language = Finnish{
		code:        LANGUAGE,
		name:        "Finnish",
		description: "Finnish is one of the two official languages of Finland",

		// http://www.languagesandnumbers.com/how-to-count-in-finnish/en/fin/
		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":  []string{"nolla"},
			"1":  []string{"yksi", "ensimmäinen"},
			"2":  []string{"kaksi"},
			"3":  []string{"kolme"},
			"4":  []string{"neljä"},
			"5":  []string{"viisi"},
			"6":  []string{"kuusi"},
			"7":  []string{"seitsemän"},
			"8":  []string{"kahdeksan"},
			"9":  []string{"yhdeksän"},
			"10": []string{"kymmenen"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z", "å", "ä", "ö"},
		vowels:       []string{"a", "e", "i", "o", "u", "y", "ä", "ö"},
		misspellings: fiMisspellings,
		homophones:   fiHomophones,
		antonyms:     fiAntonyms,
		homoglyphs: map[string][]string{
			"a": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ"},
			"b": []string{"d", "lb", "ib", "ʙ", "Ь", `b̔"`, "ɓ", "Б"},
			"c": []string{"ϲ", "с", "ƈ", "ċ", "ć", "ç"},
			"d": []string{"b", "cl", "dl", "di", "ԁ", "ժ", "ɗ", "đ"},
			"e": []string{"é", "ê", "ë", "ē", "ĕ", "ě", "ė", "е", "ẹ", "ę", "є", "ϵ", "ҽ"},
			"f": []string{"Ϝ", "ƒ", "Ғ"},
			"g": []string{"q", "ɢ", "ɡ", "Ԍ", "Ԍ", "ġ", "ğ", "ց", "ǵ", "ģ"},
			"h": []string{"lh", "ih", "һ", "հ", "Ꮒ", "н"},
			"i": []string{"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			"j": []string{"ј", "ʝ", "ϳ", "ɉ"},
			"k": []string{"lk", "ik", "lc", "κ", "ⲕ", "κ"},
			"l": []string{"1", "i", "ɫ", "ł"},
			"m": []string{"n", "nn", "rn", "rr", "ṃ", "ᴍ", "м", "ɱ"},
			"n": []string{"m", "r", "ń"},
			"o": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
			"p": []string{"ρ", "р", "ƿ", "Ϸ", "Þ"},
			"q": []string{"g", "զ", "ԛ", "գ", "ʠ"},
			"r": []string{"ʀ", "Г", "ᴦ", "ɼ", "ɽ"},
			"s": []string{"Ⴝ", "Ꮪ", "ʂ", "ś", "ѕ"},
			"t": []string{"τ", "т", "ţ"},
			"u": []string{"μ", "υ", "Ս", "ս", "ц", "ᴜ", "ǔ", "ŭ"},
			"v": []string{"ѵ", "ν", "v̇"},
			"w": []string{"vv", "ѡ", "ա", "ԝ"},
			"x": []string{"х", "ҳ", "ẋ"},
			"y": []string{"ʏ", "γ", "у", "Ү", "ý"},
			"z": []string{"ʐ", "ż", "ź", "ʐ", "ᴢ"},
			"å": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ"},
			"ä": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ"},
			"ö": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &Language
	})
}
