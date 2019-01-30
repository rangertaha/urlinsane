// Copyright © 2019 rangertaha rangertaha@gmail.com
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
	// iwMisspellings are common misspellings
	iwMisspellings = [][]string{
		// []string{"", ""},
	}

	// iwHomophones are words that sound alike
	iwHomophones = [][]string{
		[]string{"נקודה", "."},
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
			"0":  []string{"zeo"},
			"1":  []string{"אחד"},
			"2":  []string{"שתיים"},
			"3":  []string{"שלוש"},
			"4":  []string{"ארבעה"},
			"5":  []string{""},
			"6":  []string{""},
			"7":  []string{""},
			"8":  []string{""},
			"9":  []string{""},
			"10": []string{""},
		},
		Graphemes:    []string{""},
		Misspellings: iwMisspellings,
		Homophones:   iwHomophones,
		Antonyms:     iwAntonyms,
		Homoglyphs:   map[string][]string{
			//"": []string{},
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
