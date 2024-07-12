// Copyright (C) 2024  Rangertaha <rangertaha@gmail.com>
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
	// iwMisspellings are common misspellings
	iwMisspellings = [][]string{
		[]string{"כשהגעתי", "שהגעתי"},
		[]string{"אני יבוא", "אני אבוא"},
		[]string{"נחרט", "נחרת"},
		[]string{"לתבוע", "לטבוע"},
		[]string{"הנידון", "הנדון"},
	}

	// iwHomophones are words that sound alike
	iwHomophones = [][]string{
		[]string{"נקודה", "."},
		[]string{"לזנק", "-"},
	}

	// iwAntonyms are words opposite in meaning to another (e.g. bad and good ).
	iwAntonyms = map[string][]string{
		"טוב": []string{"רע"},
	}

	// Hebrew language
	iwLanguage = Language{
		Code: "IW",
		Name: "Hebrew",
		Numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":  []string{"אפס"},
			"1":  []string{"אחד"},
			"2":  []string{"שתיים"},
			"3":  []string{"שלוש"},
			"4":  []string{"ארבעה"},
			"5":  []string{"חמישה"},
			"6":  []string{"שישה"},
			"7":  []string{"שבע"},
			"8":  []string{"שמונה"},
			"9":  []string{"תשע"},
			"10": []string{"עשר"},
		},
		Graphemes: []string{
			"א", "ב", "ג", "ד", "ה", "ו",
			"ז", "ח", "ט", " י", "כ", "ל",
			"מ", "נ", "ס", "ע", "פ", "צ",
			"ק", "ר", "ש", "ת"},
		Misspellings: iwMisspellings,
		Homophones:   iwHomophones,
		Antonyms:     iwAntonyms,
		Homoglyphs: map[string][]string{
			"א":  []string{"x", "X"},
			"ב":  []string{"1", "l"},
			"ג":  []string{"i"},
			"ד":  []string{"T", "t"},
			"ה":  []string{"n"},
			"ו":  []string{"i"},
			"ז":  []string{"t", "T"},
			"ח":  []string{"n"},
			"ט":  []string{"u", "U"},
			" י": []string{"-"},
			"כ":  []string{"J", "j"},
			"ל":  []string{"7"},
			"מ":  []string{"D"},
			"נ":  []string{"l"},
			"ס":  []string{"o", "0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ", "ه", "ة"},
			"ע":  []string{"v", "y"},
			"פ":  []string{"g"},
			"צ":  []string{"y"},
			"ק":  []string{"p", "P"},
			"ר":  []string{"l"},
			"ש":  []string{"w"},
			"ת":  []string{"n"},
		},
	}

	iwKeyboards = []Keyboard{
		{
			Code:        "IW1",
			Name:        "Hebrew",
			Description: "Hebrew standard layout",
			Language:    iwLanguage,
			Layout: []string{
				"1234567890 ",
				` פםןוטארק  `,
				` ףךלחיעכגדש `,
				` ץתצמנהבסז  `},
		},
	}
)

func init() {
	KEYBOARDS.Add(iwKeyboards)
	KEYBOARDS.Append("IW", iwKeyboards)
	KEYBOARDS.Append("ALL", iwKeyboards)
}
