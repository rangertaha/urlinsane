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
package french

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/languages"
)

const LANGUAGE string = "fr"

type French struct {
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

func (l *French) Id() string {
	return l.code
}
func (l *French) Name() string {
	return l.name
}
func (l *French) Description() string {
	return l.description
}
func (l *French) Numerals() map[string][]string {
	return l.numerals
}
func (l *French) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *French) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *French) Graphemes() []string {
	return l.graphemes
}

func (l *French) Vowels() []string {
	return l.vowels
}

func (l *French) Misspellings() [][]string {
	return l.misspellings
}

func (l *French) Homophones() [][]string {
	return l.homophones
}

func (l *French) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *French) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *French) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *French) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *French) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *French) Keyboards() []urlinsane.Keyboard {
	return languages.Keyboards()
}

var (
	// frMisspellings are common misspellings
	frMisspellings = [][]string{
		[]string{"", ""},
	}

	// frHomophones are words that sound alike
	frHomophones = [][]string{
		[]string{"point", "."},
	}

	// frAntonyms are words opposite in meaning to another (e.g. bad and good ).
	frAntonyms = map[string][]string{
		"bien": []string{"mal"},
	}

	// French language
	Language = French{
		code:        LANGUAGE,
		name:        "French",
		description: "French is an official language in 27 countries",

		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":  []string{"zéro"},
			"1":  []string{"un", "premier"},
			"2":  []string{"deux", "seconde"},
			"3":  []string{"trois", "troisième"},
			"4":  []string{"quatre", "quatrième"},
			"5":  []string{"cinq", "cinquième"},
			"6":  []string{"six", "sixième"},
			"7":  []string{"sept", "septième"},
			"8":  []string{"huit", "huitième"},
			"9":  []string{"neuf", "neuvième"},
			"10": []string{"dix", "dixième"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z", "ê", "û", "î", "ô", "â",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "y"},
		misspellings: frMisspellings,
		homophones:   frHomophones,
		antonyms:     frAntonyms,
		homoglyphs: map[string][]string{
			"a": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ", "٨"},
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
			"l": []string{"1", "i", "ɫ", "ł", "١", "ا", "", ""},
			"m": []string{"n", "nn", "rn", "rr", "ṃ", "ᴍ", "м", "ɱ"},
			"n": []string{"m", "r", "ń"},
			"o": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ", "ه", "ة"},
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
			"â": []string{"à", "á", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ", "٨"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() urlinsane.Language {
		return &Language
	})
}
