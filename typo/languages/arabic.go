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
	arLanguage = Language{
		Code: "AR",
		Name: "Arabic",
		// https://www2.rocketlanguages.com/arabic/lessons/numbers-in-arabic/
		Numerals: map[string][]string{
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
		Graphemes: []string{
			"ض", "ص", "ث", "ق", "ف", "غ", "ع",
			"ه", "خ", "ح", "ج", "ة", "ش", "س", "ي", "ب",
			"ل", "ا", "ت", "ن", "م", "ك", "ظ", "ط", "ذ",
			"د", "ز", "ر", "و"},
		Misspellings: arMisspellings,
		Homophones:   arHomophones,
		Antonyms:     arAntonyms,
		Homoglyphs: map[string][]string{
			"ض": []string{},
			"ص": []string{},
			"ث": []string{},
			"ق": []string{},
			"ف": []string{},
			"غ": []string{},
			"ع": []string{},
			"ه": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
			"خ": []string{"ج", "ح"},
			"ح": []string{"خ", "ج"},
			"ج": []string{"خ", "ح"},
			"ة": []string{},
			"ش": []string{"ش"},
			"س": []string{"vv", "ѡ", "ա", "ԝ"},
			"ي": []string{},
			"ب": []string{},
			"ل": []string{},
			"ا": []string{"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			"ت": []string{},
			"ن": []string{},
			"م": []string{},
			"ك": []string{},
			"ظ": []string{},
			"ط": []string{},
			"ذ": []string{},
			"د": []string{},
			"ز": []string{},
			"ر": []string{},
		},
	}

	arKeyboards = []Keyboard{
		{
			Code:        "AR1",
			Name:        "غفقثصض",
			Description: "Arabic keyboard layout",
			Language:    arLanguage,
			Layout: []string{
				"١٢٣٤٥٦٧٨٩٠- ",
				"ةجحخهعغفقثصض",
				"  كمنتالبيسش",
				"     ورزدذطظ"},
		},
		{
			Code:        "AR2",
			Name:        "AZERTY PC",
			Description: "Arabic PC keyboard layout",
			Language:    arLanguage,
			Layout: []string{
				` é   -è çà   `,
				"ذدجحخهعغفقثصض",
				"  طكمنتالبيسش",
				"   ظزوةىلارؤءئ"},
		},
		{
			Code:        "AR3",
			Name:        "North Africa",
			Description: "Arabic North african keyboard layout",
			Language:    arLanguage,
			Layout: []string{
				"1234567890  ",
				"ةجحخهعغفقثصض",
				"  كمنتالبيسش",
				"     ورزدذطظ"},
		},
		{
			Code:        "AR4",
			Name:        "QWERTY",
			Description: "Arabic keyboard layout",
			Language:    arLanguage,
			Layout: []string{
				"١٢٣٤٥٦٧٨٩٠  ",
				"ظثةهيوطترعشق",
				"   لكجحغفدسا",
				"     منبذصخز"},
		},
	}
)

func init() {
	KEYBOARDS.Add(arKeyboards)
	KEYBOARDS.Append("AR", arKeyboards)
	KEYBOARDS.Append("ALL", arKeyboards)
}
