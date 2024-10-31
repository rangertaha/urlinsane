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
package spanish

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "es"

type Spanish struct {
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

func (l *Spanish) Id() string {
	return l.code
}
func (l *Spanish) Name() string {
	return l.name
}
func (l *Spanish) Description() string {
	return l.description
}
func (l *Spanish) Numerals() map[string][]string {
	return l.numerals
}
func (l *Spanish) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *Spanish) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *Spanish) Graphemes() []string {
	return l.graphemes
}

func (l *Spanish) Vowels() []string {
	return l.vowels
}

func (l *Spanish) Misspellings() [][]string {
	return l.misspellings
}

func (l *Spanish) Homophones() [][]string {
	return l.homophones
}

func (l *Spanish) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *Spanish) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *Spanish) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *Spanish) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *Spanish) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *Spanish) Keyboards() []internal.Keyboard {
	return languages.Keyboards()
}

var (
	// esMisspellings are common misspellings
	esMisspellings = [][]string{
		[]string{"", ""},
	}

	// esHomophones are words that sound alike
	esHomophones = [][]string{
		[]string{"punto", "."},
	}

	// esAntonyms are words opposite in meaning to another (e.g. bad and good ).
	esAntonyms = map[string][]string{
		"bueno": []string{"malo"},
	}

	// SPANISH Language
	Language = Spanish{
		code:        LANGUAGE,
		name:        "Spanish",
		description: "Spanish is an official language in 20 countries",

		// https://www.donquijote.org/spanish-language/numbers/
		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":  []string{"zero"},
			"1":  []string{"uno"},
			"2":  []string{"dos"},
			"3":  []string{"tres"},
			"4":  []string{"cuatro"},
			"5":  []string{"cinco"},
			"6":  []string{"seis"},
			"7":  []string{"siete"},
			"8":  []string{"ocho"},
			"9":  []string{"nueve"},
			"10": []string{"diez"},
			"11": []string{"once"},
			"12": []string{"doce"},
			"13": []string{"trece"},
			"14": []string{"catorce"},
			"15": []string{"quince"},
			"16": []string{"dieciséis", "dieciseis"},
			"17": []string{"diecisiete"},
			"18": []string{"dieciocho"},
			"19": []string{"diecinueve"},
			"20": []string{"veinte"},
			"21": []string{"veintiuno"},
			"22": []string{"veintidós", "veintidos"},
			"23": []string{"veintitrés", "veintitres"},
			"24": []string{"veinticuatro"},
			"25": []string{"veinticinco"},
			"26": []string{"veintiséis", "veintiseis"},
			"27": []string{"veintisiete"},
			"28": []string{"veintiocho"},
			"29": []string{"veintinueve"},
			"30": []string{"treinta"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"ñ", "o", "p", "q", "r", "s", "t",
			"u", "v", "w", "x", "y", "z"},
		vowels:       []string{"a", "e", "i", "o", "u"},
		misspellings: esMisspellings,
		homophones:   esHomophones,
		antonyms:     esAntonyms,
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
			"n": []string{"m", "r", "ń", "ñ"},
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
			"ñ": []string{"n", "ń", "r"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &Language
	})
}
