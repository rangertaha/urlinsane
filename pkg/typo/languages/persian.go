// The MIT License (MIT)
//
// Copyright © 2019 Rangertaha <rangertaha@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
	}

	faKeyboards = []Keyboard{
		{
			Code:        "FA1",
			Name:        "Persian",
			Description: "Persian standard layout",
			Language:    faLanguage,
			Layout: []string{
				"۱۲۳۴۵۶۷۸۹۰-  ",
				" چجحخهعغفقثصض",
				"  گکمنتالبیسش",
				"     وپدذرزطظ"},
		},
	}
)

func init() {
	KEYBOARDS.Add(faKeyboards)
	KEYBOARDS.Append("FA", faKeyboards)
	KEYBOARDS.Append("ALL", faKeyboards)
}
