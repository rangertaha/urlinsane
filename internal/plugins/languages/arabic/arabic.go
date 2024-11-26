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
package arabic

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "ar"

type Arabic struct {
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

func (l *Arabic) Id() string {
	return l.code
}
func (l *Arabic) Name() string {
	return l.name
}
func (l *Arabic) Description() string {
	return l.description
}
func (l *Arabic) Numerals() map[string][]string {
	return l.numerals
}
func (l *Arabic) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *Arabic) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *Arabic) Graphemes() []string {
	return l.graphemes
}

func (l *Arabic) Vowels() []string {
	return l.vowels
}

func (l *Arabic) Misspellings() [][]string {
	return l.misspellings
}

func (l *Arabic) Homophones() [][]string {
	return l.homophones
}

func (l *Arabic) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *Arabic) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *Arabic) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *Arabic) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *Arabic) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *Arabic) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	// arMisspellings are common misspellings
	arMisspellings = [][]string{
		// []string{"", ""},
	}

	// arHomophones are words that sound alike
	arHomophones = [][]string{
		[]string{"نقطة", "."},
	}

	// arAntonyms are words opposite in meaning to another (e.g. bad and good ).
	arAntonyms = map[string][]string{
		"حسن": []string{"سيئة"},
	}

	// Arabic language
	arLanguage = Arabic{
		code:        LANGUAGE,
		name:        "Arabic",
		description: "Arabic is spoken primarily in the Arab world",

		// https://www2.rocketlanguages.com/arabic/lessons/numbers-in-arabic/
		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"٠":  []string{"صفر", "sifr"},
			"١":  []string{"واحد", "أول", "wa7ed"},
			"٢":  []string{"اثنان", "اتنين", "ثانيا", "etneyn", "athnan"},
			"٣":  []string{"تلاتة", "الثالث", "talata"},
			"٤":  []string{"اربعة", "رابع", "arba3a"},
			"٥":  []string{"خمسة", "خامس", "7amsa"},
			"٦":  []string{"ستة", "السادس", "setta"},
			"٧":  []string{"سابعة", "سابع", "sab3a"},
			"٨":  []string{"تمانية", "ثامن", "tamanya"},
			"٩":  []string{"تسعة", "تاسع", "tes3a"},
			"١٠": []string{"عشرة", "العاشر", "3ashara"},
		},
		graphemes: []string{
			"ض", "ص", "ث", "ق", "ف", "غ", "ع",
			"ه", "خ", "ح", "ج", "ة", "ش", "س", "ي", "ب",
			"ل", "ا", "ت", "ن", "م", "ك", "ظ", "ط", "ذ",
			"د", "ز", "ر", "و"},
		misspellings: arMisspellings,
		homophones:   arHomophones,
		antonyms:     arAntonyms,
		homoglyphs: map[string][]string{
			"ض": []string{"ص", "ظ", "ط", "ড", "b", "в"},
			"ص": []string{"ض", "ظ", "ط"},
			"ث": []string{"ت", "ن"},
			"ق": []string{"م"},
			"ف": []string{"م"},
			"غ": []string{"ع", "خ"},
			"ع": []string{"غ", "خ"},
			"ه": []string{"ة", "0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
			"خ": []string{"ج", "ح", "ع"},
			"ح": []string{"خ", "ج", "ع"},
			"ج": []string{"خ", "ح", "ع"},
			"ة": []string{"ن", "ق"},
			"ش": []string{"س", "ث"},
			"س": []string{"vv", "ѡ", "ա", "ԝ"},
			"ي": []string{"ف"},
			"ب": []string{"ث", "ت", "ن"},
			"ل": []string{"j", "J"},
			"ا": []string{"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			"ت": []string{"ن", "ث"},
			"ن": []string{"ت", "ث"},
			"م": []string{"ق", "ف", "غ", "ع"},
			"ك": []string{"ل", "ا"},
			"ظ": []string{"ط", "ص", "ض"},
			"ط": []string{"ظ", "ص", "ض"},
			"ذ": []string{"ز", "د", "ر"},
			"د": []string{"ز", "ذ", "ر"},
			"ز": []string{"ر", "د", "ذ"},
			"ر": []string{"ز", "د", "ذ"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &arLanguage
	})
}
