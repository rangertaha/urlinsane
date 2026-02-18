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
package hebrew

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "iw"

type Hebrew struct {
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

func (l *Hebrew) Id() string {
	return l.code
}
func (l *Hebrew) Name() string {
	return l.name
}
func (l *Hebrew) Description() string {
	return l.description
}
func (l *Hebrew) Numerals() map[string][]string {
	return l.numerals
}
func (l *Hebrew) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *Hebrew) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *Hebrew) Graphemes() []string {
	return l.graphemes
}

func (l *Hebrew) Vowels() []string {
	return l.vowels
}

func (l *Hebrew) Misspellings() [][]string {
	return l.misspellings
}

func (l *Hebrew) Homophones() [][]string {
	return l.homophones
}

func (l *Hebrew) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *Hebrew) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *Hebrew) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *Hebrew) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *Hebrew) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *Hebrew) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	// iwMisspellings are common misspellings
	iwMisspellings = [][]string{
		// Domain-friendly variants (avoid spaces)
		{"כשהגעתי", "שהגעתי"},
		{"נחרט", "נחרת"},
		{"לתבוע", "לטבוע"},
		{"הנידון", "הנדון"},

		// Common Hebrew spelling variants (mater lectionis / common alternations)
		{"אימייל", "מייל", "דואל"},
		{"טכנולוגיה", "טכנולגיה"},
		{"אינטרנט", "אינטרנת"},
	}

	// iwHomophones are words that sound alike
	iwHomophones = [][]string{
		{"נקודה", "."},
		{"מקף", "-"},
	}

	// iwAntonyms are words opposite in meaning to another (e.g. bad and good ).
	iwAntonyms = map[string][]string{
		"טוב":         {"רע"},
		"רע":          {"טוב"},
		"גדול":        {"קטן"},
		"קטן":         {"גדול"},
		"גבוה":        {"נמוך"},
		"נמוך":        {"גבוה"},
		"מהר":         {"לאט"},
		"לאט":         {"מהר"},
		"חזק":         {"חלש"},
		"חלש":         {"חזק"},
		"חדש":         {"ישן"},
		"ישן":         {"חדש"},
		"יום":         {"לילה"},
		"לילה":        {"יום"},
		"כן":          {"לא"},
		"לא":          {"כן"},
		"בתוך":        {"בחוץ"},
		"בחוץ":        {"בתוך"},
		"התחלה":       {"סוף"},
		"סוף":         {"התחלה"},
		"אפשרי":       {"בלתי אפשרי"},
		"בלתי אפשרי": {"אפשרי"},
		"חם":          {"קר"},
		"קר":          {"חם"},
		"קרוב":        {"רחוק"},
		"רחוק":        {"קרוב"},
		"אור":         {"חושך"},
		"חושך":       {"אור"},
		"פתוח":        {"סגור"},
		"סגור":        {"פתוח"},
		"למעלה":       {"למטה"},
		"למטה":        {"למעלה"},
		"לפני":        {"אחרי"},
		"אחרי":        {"לפני"},
		"מוקדם":       {"מאוחר"},
		"מאוחר":       {"מוקדם"},
		"אמת":         {"שקר"},
		"שקר":         {"אמת"},
		"נכון":        {"שגוי"},
		"שגוי":        {"נכון"},
		"עשיר":        {"עני"},
		"עני":         {"עשיר"},
		"שמח":         {"עצוב"},
		"עצוב":        {"שמח"},
		"קל":          {"קשה"},
		"קשה":         {"קל"},
		"מלא":         {"ריק"},
		"ריק":         {"מלא"},
		"נקי":         {"מלוכלך"},
		"מלוכלך":      {"נקי"},
		"לקנות":       {"למכור"},
		"למכור":       {"לקנות"},
		"לתת":         {"לקחת"},
		"לקחת":        {"לתת"},
		"לאהוב":       {"לשנוא"},
		"לשנוא":       {"לאהוב"},
		"לנצח":        {"להפסיד"},
		"להפסיד":      {"לנצח"},
		"חיים":        {"מוות"},
		"מוות":        {"חיים"},
	}

	// Hebrew language
	Language = Hebrew{
		code:        LANGUAGE,
		name:        "Hebrew",
		description: "Hebrew is one of the official languages of the State of Israel",

		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":          {"אפס"},
			"1":          {"אחד", "ראשון"},
			"2":          {"שתיים", "שני"},
			"3":          {"שלוש", "שלישי"},
			"4":          {"ארבע", "רביעי"},
			"5":          {"חמש", "חמישי"},
			"6":          {"שש", "שישי"},
			"7":          {"שבע", "שביעי"},
			"8":          {"שמונה", "שמיני"},
			"9":          {"תשע", "תשיעי"},
			"10":         {"עשר", "עשירי"},
			"11":         {"אחדעשר"},
			"12":         {"שתיםעשרה"},
			"13":         {"שלושעשרה"},
			"14":         {"ארבעעשרה"},
			"15":         {"חמשעשרה"},
			"16":         {"ששעשרה"},
			"17":         {"שבעעשרה"},
			"18":         {"שמונהעשרה"},
			"19":         {"תשעעשרה"},
			"20":         {"עשרים"},
			"30":         {"שלושים"},
			"40":         {"ארבעים"},
			"50":         {"חמישים"},
			"60":         {"שישים"},
			"70":         {"שבעים"},
			"80":         {"שמונים"},
			"90":         {"תשעים"},
			"100":        {"מאה"},
			"1000":       {"אלף"},
			"1000000":    {"מיליון"},
			"1000000000": {"מיליארד"},

			// (Digits are already ASCII keys in this dataset)
		},
		graphemes: []string{
			"א", "ב", "ג", "ד", "ה", "ו",
			"ז", "ח", "ט", "י", "כ", "ל",
			"מ", "נ", "ס", "ע", "פ", "צ",
			"ק", "ר", "ש", "ת"},
		misspellings: iwMisspellings,
		homophones:   iwHomophones,
		antonyms:     iwAntonyms,
		homoglyphs: map[string][]string{
			"א": []string{"x", "X"},
			"ב": []string{"1", "l"},
			"ג": []string{"i"},
			"ד": []string{"T", "t"},
			"ה": []string{"n"},
			"ו": []string{"i"},
			"ז": []string{"t", "T"},
			"ח": []string{"n"},
			"ט": []string{"u", "U"},
			"י": []string{"-"},
			"כ": []string{"J", "j"},
			"ל": []string{"7"},
			"מ": []string{"D"},
			"נ": []string{"l"},
			"ס": []string{"o", "0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ", "ه", "ة"},
			"ע": []string{"v", "y"},
			"פ": []string{"g"},
			"צ": []string{"y"},
			"ק": []string{"p", "P"},
			"ר": []string{"l"},
			"ש": []string{"w"},
			"ת": []string{"n"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &Language
	})
}
