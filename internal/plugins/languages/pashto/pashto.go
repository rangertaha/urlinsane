// Copyright 2026 Rangertaha. All Rights Reserved.
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
package pashto

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "ps"

type Pashto struct {
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

func (l *Pashto) Id() string {
	return l.code
}
func (l *Pashto) Name() string {
	return l.name
}
func (l *Pashto) Description() string {
	return l.description
}
func (l *Pashto) Numerals() map[string][]string {
	return l.numerals
}
func (l *Pashto) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *Pashto) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *Pashto) Graphemes() []string {
	return l.graphemes
}

func (l *Pashto) Vowels() []string {
	return l.vowels
}

func (l *Pashto) Misspellings() [][]string {
	return l.misspellings
}

func (l *Pashto) Homophones() [][]string {
	return l.homophones
}

func (l *Pashto) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *Pashto) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *Pashto) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *Pashto) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *Pashto) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *Pashto) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	psMisspellings = [][]string{
		// Common Arabic/Persian character variants that show up in Pashto text / IDNs
		{"کې", "كي"}, // Pashto "kē" with Arabic yeh
		{"یې", "يې"},
		{"ګ", "گ"},
		{"ک", "ك"},
	}

	psHomophones = [][]string{
		{"نقطه", "."},
		{"ټکي", "."},
		{"اټ", "@"},
		{"ډش", "-"},
	}

	psAntonyms = map[string][]string{
		"ښه":  {"بد"},
		"بد":  {"ښه"},
		"هو":  {"نه"},
		"نه":  {"هو"},
		"لوړ": {"ټیټ"},
		"ټیټ": {"لوړ"},
		"سړه": {"ګرمه"},
		"ګرمه": {"سړه"},
	}

	Language = Pashto{
		code:        LANGUAGE,
		name:        "Pashto",
		description: "Pashto is an Eastern Iranian language written in an Arabic-based script (Afghanistan and Pakistan).",

		numerals: map[string][]string{
			// Pashto digits (Persian-style)
			"۰":  {"صفر"},
			"۱":  {"یو", "لومړی"},
			"۲":  {"دوه", "دوهم"},
			"۳":  {"درې", "درېیم"},
			"۴":  {"څلور", "څلورم"},
			"۵":  {"پنځه", "پنځم"},
			"۶":  {"شپږ", "شپږم"},
			"۷":  {"اووه", "اووم"},
			"۸":  {"اته", "اتم"},
			"۹":  {"نهه", "نهم"},
			"۱۰": {"لس", "لسم"},
			"۲۰": {"شل"},
			"۳۰": {"دېرش"},
			"۴۰": {"څلوېښت"},
			"۵۰": {"پنځوس"},
			"۶۰": {"شپېته"},
			"۷۰": {"اویا"},
			"۸۰": {"اتیا"},
			"۹۰": {"نوي"},
			"۱۰۰": {"سل"},
			"۱۰۰۰": {"زر"},

			// Arabic-Indic digits
			"٠":  {"صفر"},
			"١":  {"یو", "لومړی"},
			"٢":  {"دوه", "دوهم"},
			"٣":  {"درې", "درېیم"},
			"٤":  {"څلور", "څلورم"},
			"٥":  {"پنځه", "پنځم"},
			"٦":  {"شپږ", "شپږم"},
			"٧":  {"اووه", "اووم"},
			"٨":  {"اته", "اتم"},
			"٩":  {"نهه", "نهم"},
			"١٠": {"لس", "لسم"},

			// ASCII digits (common in domains)
			"0":  {"صفر"},
			"1":  {"یو", "لومړی"},
			"2":  {"دوه", "دوهم"},
			"3":  {"درې", "درېیم"},
			"4":  {"څلور", "څلورم"},
			"5":  {"پنځه", "پنځم"},
			"6":  {"شپږ", "شپږم"},
			"7":  {"اووه", "اووم"},
			"8":  {"اته", "اتم"},
			"9":  {"نهه", "نهم"},
			"10": {"لس", "لسم"},
		},
		graphemes: []string{
			"ا", "ب", "پ", "ت", "ټ", "ث", "ج", "چ", "ځ", "څ",
			"ح", "خ", "د", "ډ", "ذ", "ر", "ړ", "ز", "ژ", "ږ",
			"س", "ش", "ښ", "ص", "ض", "ط", "ظ", "ع", "غ", "ف",
			"ق", "ک", "ګ", "ل", "م", "ن", "ڼ", "و", "ه", "ی",
			"ې", "ۍ",
		},
		vowels:       []string{"ا", "و", "ی", "ې", "ۍ"},
		misspellings: psMisspellings,
		homophones:   psHomophones,
		antonyms:     psAntonyms,
		homoglyphs: map[string][]string{
			// Arabic/Persian confusables commonly seen in Pashto text
			"ی": {"ي", "ى", "ئ", "ﻱ"},
			"ي": {"ی", "ى", "ئ", "ﻱ"},
			"ک": {"ك", "ﻙ"},
			"ك": {"ک", "ﻙ"},
			"ګ": {"گ"},
			"گ": {"ګ"},
			"ه": {"ة", "ۀ", "ﻩ"},
			"ة": {"ه", "ۀ", "ﻩ"},
			"ا": {"آ", "أ", "إ", "ٱ"},
			"و": {"ؤ"},
			"ؤ": {"و"},

			// Pashto-specific near-confusables (keep conservative)
			"ټ": {"ت"},
			"ت": {"ټ"},
			"ډ": {"د"},
			"د": {"ډ", "ذ"},
			"ړ": {"ر"},
			"ر": {"ړ"},
			"ږ": {"ژ"},
			"ژ": {"ږ"},
			"ښ": {"ش"},
			"ش": {"ښ"},
			"ځ": {"ج"},
			"څ": {"چ"},
			"ن": {"ڼ"},
			"ڼ": {"ن"},
			"ې": {"ی"},
			"ۍ": {"ی"},

			// Digit confusables (Persian/Arabic-Indic/ASCII)
			"۰": {"0", "٠"},
			"۱": {"1", "١"},
			"۲": {"2", "٢"},
			"۳": {"3", "٣"},
			"۴": {"4", "٤"},
			"۵": {"5", "٥"},
			"۶": {"6", "٦"},
			"۷": {"7", "٧"},
			"۸": {"8", "٨"},
			"۹": {"9", "٩"},
			"٠": {"0", "۰"},
			"١": {"1", "۱"},
			"٢": {"2", "۲"},
			"٣": {"3", "۳"},
			"٤": {"4", "۴"},
			"٥": {"5", "۵"},
			"٦": {"6", "۶"},
			"٧": {"7", "۷"},
			"٨": {"8", "۸"},
			"٩": {"9", "۹"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &Language
	})
}
