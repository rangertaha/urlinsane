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
		// Word-level misspellings / common variants (avoid single-letter sets to reduce variant explosion)
		{"الى", "إلى"},
		{"مسول", "مسؤول"},
		{"مساله", "مسألة"},
		{"اسلام", "إسلام"},
		{"امير", "أمير"},
		{"مدرسه", "مدرسة"},
		{"جامعه", "جامعة"},
		{"شركه", "شركة"},
	}

	// arHomophones are words that sound alike
	arHomophones = [][]string{
		{"نقطة", "."},
		{"شرطة", "-"},
		{"آت", "@"},
		{"سلاش", "/"},
	}

	// arAntonyms are words opposite in meaning to another (e.g. bad and good ).
	arAntonyms = map[string][]string{
		"جيد":   {"سيئ"},
		"سيئ":   {"جيد"},
		"كبير":  {"صغير"},
		"صغير":  {"كبير"},
		"سريع":  {"بطيء"},
		"بطيء":  {"سريع"},
		"حار":   {"بارد"},
		"بارد":  {"حار"},
		"قريب":  {"بعيد"},
		"بعيد":  {"قريب"},
		"قديم":  {"جديد"},
		"جديد":  {"قديم"},
		"قوي":   {"ضعيف"},
		"ضعيف":  {"قوي"},
		"سعيد":  {"حزين"},
		"حزين":  {"سعيد"},
		"نعم":   {"لا"},
		"لا":    {"نعم"},
		"داخل":  {"خارج"},
		"خارج":  {"داخل"},
		"بداية": {"نهاية"},
		"نهاية": {"بداية"},
		"طويل":  {"قصير"},
		"قصير":  {"طويل"},
		"كثير":  {"قليل"},
		"قليل":  {"كثير"},
		"فوق":   {"تحت"},
		"تحت":   {"فوق"},
		"قبل":   {"بعد"},
		"بعد":   {"قبل"},
		"صحيح":  {"خطأ"},
		"خطأ":   {"صحيح"},
		"حق":    {"باطل"},
		"باطل":  {"حق"},
		"مفتوح": {"مغلق"},
		"مغلق":  {"مفتوح"},
		"نظيف":  {"وسخ"},
		"وسخ":   {"نظيف"},
		"غني":   {"فقير"},
		"فقير":  {"غني"},
		"سهل":   {"صعب"},
		"صعب":   {"سهل"},
		"نور":   {"ظلام"},
		"ظلام":  {"نور"},
		"شراء":  {"بيع"},
		"بيع":   {"شراء"},
		"يأتي":  {"يذهب"},
		"يذهب":  {"يأتي"},
		"حياة":  {"موت"},
		"موت":   {"حياة"},
		"نهار":  {"ليل"},
		"ليل":   {"نهار"},
	}

	// Arabic language
	arLanguage = Arabic{
		code:        LANGUAGE,
		name:        "Arabic",
		description: "Arabic is spoken primarily in the Arab world",

		// https://www2.rocketlanguages.com/arabic/lessons/numbers-in-arabic/
		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			// Arabic-Indic digits
			"٠":          {"صفر", "sifr"},
			"١":          {"واحد", "أول", "wa7ed"},
			"٢":          {"اثنان", "ثاني", "etneyn", "athnan"},
			"٣":          {"ثلاثة", "ثالث", "talata"},
			"٤":          {"أربعة", "رابع", "arba3a"},
			"٥":          {"خمسة", "خامس", "7amsa"},
			"٦":          {"ستة", "سادس", "setta"},
			"٧":          {"سبعة", "سابع", "sab3a"},
			"٨":          {"ثمانية", "ثامن", "tamanya"},
			"٩":          {"تسعة", "تاسع", "tes3a"},
			"١٠":         {"عشرة", "عاشر", "3ashara"},
			"١١":         {"أحدعشر"},
			"١٢":         {"اثناعشر"},
			"١٣":         {"ثلاثةعشر"},
			"١٤":         {"أربعةعشر"},
			"١٥":         {"خمسةعشر"},
			"٢٠":         {"عشرون"},
			"٣٠":         {"ثلاثون"},
			"٤٠":         {"أربعون"},
			"٥٠":         {"خمسون"},
			"٦٠":         {"ستون"},
			"٧٠":         {"سبعون"},
			"٨٠":         {"ثمانون"},
			"٩٠":         {"تسعون"},
			"١٠٠":        {"مئة"},
			"١٠٠٠":       {"ألف"},
			"١٠٠٠٠٠٠":    {"مليون"},
			"١٠٠٠٠٠٠٠٠٠": {"مليار"},

			// ASCII digits (common in domains)
			"0":  {"صفر"},
			"1":  {"واحد", "أول"},
			"2":  {"اثنان", "ثاني"},
			"3":  {"ثلاثة", "ثالث"},
			"4":  {"أربعة", "رابع"},
			"5":  {"خمسة", "خامس"},
			"6":  {"ستة", "سادس"},
			"7":  {"سبعة", "سابع"},
			"8":  {"ثمانية", "ثامن"},
			"9":  {"تسعة", "تاسع"},
			"10": {"عشرة", "عاشر"},
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
			"ا": []string{"أ", "إ", "آ", "ٱ", "1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			"أ": []string{"ا", "إ", "آ"},
			"إ": []string{"ا", "أ", "آ"},
			"آ": []string{"ا", "أ", "إ"},
			"ت": []string{"ن", "ث"},
			"ن": []string{"ت", "ث"},
			"م": []string{"ق", "ف", "غ", "ع"},
			"ك": []string{"ل", "ا", "ک", "ك"},
			"ى": []string{"ي"},
			"ئ": []string{"ي"},
			"ؤ": []string{"و"},
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
