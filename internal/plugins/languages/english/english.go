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
package english

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "en"

type English struct {
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
	stopwords    []string
}

func (l *English) Id() string {
	return l.code
}
func (l *English) Name() string {
	return l.name
}
func (l *English) Description() string {
	return l.description
}
func (l *English) Numerals() map[string][]string {
	return l.numerals
}

func (l *English) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *English) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *English) Graphemes() []string {
	return l.graphemes
}

func (l *English) Vowels() []string {
	return l.vowels
}

func (l *English) Misspellings() [][]string {
	return l.misspellings
}

func (l *English) Homophones() [][]string {
	return l.homophones
}

func (l *English) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *English) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *English) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *English) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *English) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *English) Keyboards() []internal.Keyboard {
	return languages.Keyboards()
}

func (l *English) StopWords() []string {
	return l.stopwords
}

var (
	Language = English{
		code:        LANGUAGE,
		name:        "English",
		description: "English the most spoken language in the world",

		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0": {"zero"},
			"1": {"one", "first"},
			"2": {"two", "second"},
			"3": {"three", "third"},
			"4": {"four", "fourth", "for"},
			"5": {"five", "fifth"},
			"6": {"six", "sixth"},
			"7": {"seven", "seventh"},
			"8": {"eight", "eighth"},
			"9": {"nine", "ninth"},
			// "10": {"ten", "tenth"},
			// "11": {"eleven", "eleventh"},
			// "12": {"twelve", "twelfth"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z"},
		vowels:       []string{"a", "e", "i", "o", "u"},
		misspellings: enMisspellings,
		homophones:   enHomophones,
		antonyms:     enAntonyms,
		stopwords:    enStopWords,
		homoglyphs: map[string][]string{
			"a": {"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ", "٨"},
			"b": {"d", "lb", "ib", "ʙ", "Ь", `b̔"`, "ɓ", "Б"},
			"c": {"ϲ", "с", "ƈ", "ċ", "ć", "ç"},
			"d": {"b", "cl", "dl", "di", "ԁ", "ժ", "ɗ", "đ"},
			"e": {"é", "ê", "ë", "ē", "ĕ", "ě", "ė", "е", "ẹ", "ę", "є", "ϵ", "ҽ"},
			"f": {"Ϝ", "ƒ", "Ғ"},
			"g": {"q", "ɢ", "ɡ", "Ԍ", "Ԍ", "ġ", "ğ", "ց", "ǵ", "ģ"},
			"h": {"lh", "ih", "һ", "հ", "Ꮒ", "н"},
			"i": {"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			"j": {"ј", "ʝ", "ϳ", "ɉ"},
			"k": {"lk", "ik", "lc", "κ", "ⲕ", "κ"},
			"l": {"1", "i", "ɫ", "ł", "١", "ا", "", ""},
			"m": {"n", "nn", "rn", "rr", "ṃ", "ᴍ", "м", "ɱ"},
			"n": {"m", "r", "ń"},
			"o": {"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ", "ه", "ة"},
			"p": {"ρ", "р", "ƿ", "Ϸ", "Þ"},
			"q": {"g", "զ", "ԛ", "գ", "ʠ"},
			"r": {"ʀ", "Г", "ᴦ", "ɼ", "ɽ"},
			"s": {"Ⴝ", "Ꮪ", "ʂ", "ś", "ѕ"},
			"t": {"τ", "т", "ţ"},
			"u": {"μ", "υ", "Ս", "ս", "ц", "ᴜ", "ǔ", "ŭ"},
			"v": {"ѵ", "ν", "v̇"},
			"w": {"vv", "ѡ", "ա", "ԝ"},
			"x": {"х", "ҳ", "ẋ"},
			"y": {"ʏ", "γ", "у", "Ү", "ý"},
			"z": {"ʐ", "ż", "ź", "ʐ", "ᴢ"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &Language
	})
}
