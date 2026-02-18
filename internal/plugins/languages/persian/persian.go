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
package persian

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
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

func (l *Persian) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	// faMisspellings are common misspellings
	faMisspellings = [][]string{
		// Common Arabic/Persian character variant spellings (word-level)
		{"مي", "می"}, // Arabic yeh -> Persian yeh (common in older content)

		// Common Persian orthographic misspellings
		{"مسول", "مسئول"},
		{"مساله", "مسئله"},
		{"لطفا", "لطفاً"},
		{"اصلا", "اصلاً"},
		{"اتفاقا", "اتفاقاً"},
	}

	// faHomophones are words that sound alike
	faHomophones = [][]string{
		{"نقطه", "."},

		// Commonly confused in speech/writing (near-homophones)
		{"گذارش", "گزارش"},
		{"ذره", "زره"},
	}

	// faAntonyms are words opposite in meaning to another (e.g. bad and good ).
	faAntonyms = map[string][]string{
		"خوب":     {"بد"},
		"بد":      {"خوب"},
		"بزرگ":    {"کوچک"},
		"کوچک":    {"بزرگ"},
		"بالا":    {"پایین"},
		"پایین":   {"بالا"},
		"روشن":    {"تاریک"},
		"تاریک":   {"روشن"},
		"سریع":    {"کند"},
		"کند":     {"سریع"},
		"گرم":     {"سرد"},
		"سرد":     {"گرم"},
		"قدیم":    {"جدید"},
		"جدید":    {"قدیم"},
		"قوی":     {"ضعیف"},
		"ضعیف":    {"قوی"},
		"نزدیک":   {"دور"},
		"دور":     {"نزدیک"},
		"شروع":    {"پایان"},
		"پایان":   {"شروع"},
		"دوست":    {"دشمن"},
		"دشمن":    {"دوست"},
		"همیشه":   {"هرگز"},
		"هرگز":    {"همیشه"},
		"بله":     {"نه"},
		"نه":      {"بله"},
		"درست":    {"غلط"},
		"غلط":     {"درست"},
		"زیاد":    {"کم"},
		"کم":      {"زیاد"},
		"بلند":    {"کوتاه"},
		"کوتاه":   {"بلند"},
		"زود":     {"دیر"},
		"دیر":     {"زود"},
		"آسان":    {"سخت"},
		"سخت":     {"آسان"},
		"باز":     {"بسته"},
		"بسته":    {"باز"},
		"داخل":    {"خارج"},
		"خارج":    {"داخل"},
		"قبل":     {"بعد"},
		"بعد":     {"قبل"},
		"راست":    {"چپ"},
		"چپ":      {"راست"},
		"روز":     {"شب"},
		"شب":      {"روز"},
		"خریدن":   {"فروختن"},
		"فروختن":  {"خریدن"},
		"دادن":    {"گرفتن"},
		"گرفتن":   {"دادن"},
		"بردن":    {"باختن"},
		"باختن":   {"بردن"},
		"زندگی":   {"مرگ"},
		"مرگ":     {"زندگی"},
	}

	// Persian language
	Language = Persian{
		code:        LANGUAGE,
		name:        "Persian",
		description: "Persian is a member of the Western Iranian group of the Iranian languages",

		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			// Persian digits
			"۰":          {"صفر"},
			"۱":          {"یک", "اول"},
			"۲":          {"دو", "دوم"},
			"۳":          {"سه", "سوم"},
			"۴":          {"چهار", "چهارم"},
			"۵":          {"پنج", "پنجم"},
			"۶":          {"شش", "ششم"},
			"۷":          {"هفت", "هفتم"},
			"۸":          {"هشت", "هشتم"},
			"۹":          {"نه", "نهم"},
			"۱۰":         {"ده", "دهم"},
			"۱۱":         {"یازده"},
			"۱۲":         {"دوازده"},
			"۱۳":         {"سیزده"},
			"۱۴":         {"چهارده"},
			"۱۵":         {"پانزده"},
			"۱۶":         {"شانزده"},
			"۱۷":         {"هفده"},
			"۱۸":         {"هجده"},
			"۱۹":         {"نوزده"},
			"۲۰":         {"بیست"},
			"۳۰":         {"سی"},
			"۴۰":         {"چهل"},
			"۵۰":         {"پنجاه"},
			"۶۰":         {"شصت"},
			"۷۰":         {"هفتاد"},
			"۸۰":         {"هشتاد"},
			"۹۰":         {"نود"},
			"۱۰۰":        {"صد"},
			"۱۰۰۰":       {"هزار"},
			"۱۰۰۰۰۰۰":    {"میلیون"},
			"۱۰۰۰۰۰۰۰۰۰": {"میلیارد"},

			// ASCII digits (common in domains)
			"0":  {"صفر"},
			"1":  {"یک", "اول"},
			"2":  {"دو", "دوم"},
			"3":  {"سه", "سوم"},
			"4":  {"چهار", "چهارم"},
			"5":  {"پنج", "پنجم"},
			"6":  {"شش", "ششم"},
			"7":  {"هفت", "هفتم"},
			"8":  {"هشت", "هشتم"},
			"9":  {"نه", "نهم"},
			"10": {"ده", "دهم"},
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
			// Persian/Arabic confusables
			"ی": {"ي", "ى", "ئ", "ﻱ"},
			"ي": {"ی", "ى", "ئ", "ﻱ"},
			"ک": {"ك", "ﻙ"},
			"ك": {"ک", "ﻙ"},
			"ه": {"ة", "ۀ", "ﻩ"},
			"ة": {"ه", "ۀ", "ﻩ"},
			"ا": {"آ", "أ", "إ", "ٱ"},
			"و": {"ؤ"},
			"ؤ": {"و"},
			"ب": {"پ", "ت", "ث"},
			"پ": {"ب"},
			"ت": {"ث", "ب"},
			"ث": {"ت", "ب"},
			"ج": {"چ", "ح", "خ"},
			"چ": {"ج"},
			"ح": {"ج", "خ"},
			"خ": {"ج", "ح"},
			"س": {"ش"},
			"ش": {"س"},
			"ص": {"س"},
			"ض": {"ظ", "ص"},
			"ز": {"ژ"},
			"ژ": {"ز"},
			"د": {"ذ"},
			"ذ": {"د"},
			"ر": {"ز"},
			"ل": {"1", "I", "l"},
			// Latin/Greek digit-like confusables
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
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &Language
	})
}
