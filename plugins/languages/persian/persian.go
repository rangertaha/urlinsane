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
package persian

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/languages"
)

const LANGUAGE string = "fa"

type Persian struct {
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

func (l *Persian) Id() string {
	return l.code
}
func (l *Persian) Name() string {
	return l.name
}
func (l *Persian) Description() string {
	return l.description
}
func (l *Persian) Numerals() map[string][]string {
	return l.numerals
}
func (l *Persian) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *Persian) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *Persian) Graphemes() []string {
	return l.graphemes
}

func (l *Persian) Vowels() []string {
	return l.vowels
}

func (l *Persian) Misspellings() [][]string {
	return l.misspellings
}

func (l *Persian) Homophones() [][]string {
	return l.homophones
}

func (l *Persian) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *Persian) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *Persian) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *Persian) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *Persian) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *Persian) Keyboards() []urlinsane.Keyboard {
	return languages.Keyboards()
}

var (
	// faMisspellings are common misspellings
	faMisspellings = [][]string{
		[]string{"a", "a"},
	}

	// faHomophones are words that sound alike
	faHomophones = [][]string{
		[]string{"نقطه", "."},
	}

	// faAntonyms are words opposite in meaning to another (e.g. bad and good ).
	faAntonyms = map[string][]string{
		"خوب": []string{"بد"},
	}

	// Persian language
	Language = Persian{
		code:        LANGUAGE,
		name:        "Persian",
		description: "Persian is a member of the Western Iranian group of the Iranian languages",

		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"۰":  []string{"صفر"},
			"۱":  []string{"يك"},
			"۲":  []string{"دو"},
			"۳":  []string{"سه"},
			"۴":  []string{"چهار"},
			"۵":  []string{"پنج"},
			"۶":  []string{"شش"},
			"۷":  []string{"هفت"},
			"۸":  []string{"هشت"},
			"۹":  []string{"نه"},
			"۱۰": []string{"ده"},
		},
		graphemes: []string{
			"ا", "ب", "پ", "ت", "ث", "ج",
			"چ", "ح", "خ", "د", "ذ", "ر",
			"ز", "ژ", "س", "ش", "ص", "ض",
			"ط", "ظ", "ع", "غ", "ف", "ق",
			"ک", "گ", "ل", "م", "ن", "و",
			"ه", "ی"},
		misspellings: faMisspellings,
		homophones:   faHomophones,
		antonyms:     faAntonyms,
		homoglyphs: map[string][]string{
			"ض": []string{""},
			"ص": []string{""},
			"ث": []string{""},
			"ق": []string{""},
			"ف": []string{""},
			"غ": []string{""},
			"ع": []string{""},
			"ه": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
			"خ": []string{"ج", "ح"},
			"ح": []string{"خ", "ج"},
			"ج": []string{"خ", "ح"},
			"ة": []string{""},
			"ش": []string{"ش"},
			"س": []string{"vv", "ѡ", "ա", "ԝ"},
			"ي": []string{""},
			"ب": []string{""},
			"ل": []string{""},
			"ا": []string{"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			"ت": []string{""},
			"ن": []string{""},
			"م": []string{""},
			"ك": []string{""},
			"ظ": []string{""},
			"ط": []string{""},
			"ذ": []string{""},
			"د": []string{""},
			"ز": []string{""},
			"ر": []string{""},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() urlinsane.Language {
		return &Language
	})
}
