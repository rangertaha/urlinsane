// Copyright (C) 2024  Tal Hatchi (Rangertaha)
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
package languages

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
	faLanguage = Language{
		Code: "FA",
		Name: "Persian",
		Description: "Persian is a member of the Western Iranian group of the Iranian languages",

		Numerals: map[string][]string{
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
		Graphemes: []string{
			"ا", "ب", "پ", "ت", "ث", "ج",
			"چ", "ح", "خ", "د", "ذ", "ر",
			"ز", "ژ", "س", "ش", "ص", "ض",
			"ط", "ظ", "ع", "غ", "ف", "ق",
			"ک", "گ", "ل", "م", "ن", "و",
			"ه", "ی"},
		Misspellings: faMisspellings,
		Homophones:   faHomophones,
		Antonyms:     faAntonyms,
		Homoglyphs: map[string][]string{
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
		Keyboards: []Keyboard{
			{
				Code:        "FA1",
				Name:        "Persian",
				Description: "Persian standard layout",
				Layout: []string{
					"۱۲۳۴۵۶۷۸۹۰-  ",
					" چجحخهعغفقثصض",
					"  گکمنتالبیسش",
					"     وپدذرزطظ"},
			},
		},
	}
)

func init() {
	Add("fa", faLanguage)
}
