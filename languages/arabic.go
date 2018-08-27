// Copyright © 2018 rangertaha rangertaha@gmail.com
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

// Common misspellings
var AR_SPELLINGS = [][]string{
// []string{"misspelling", "correct1", "correct2"},
}

// Common homophones
var AR_HOMOPHONES = [][]string{}

var ARABIC = Language{
	Code: "AR",
	Name: "Arabic",
	Numerals: map[string][]string{
		"٠":  []string{""},
		"١":  []string{""},
		"٢":  []string{""},
		"٣":  []string{""},
		"٤":  []string{""},
		"٥":  []string{""},
		"٦":  []string{""},
		"٧":  []string{""},
		"٨":  []string{""},
		"٩":  []string{""},
		"١٠": []string{""},
	},
	Graphemes: []string{
		"ض", "ص", "ث", "ق", "ف", "غ", "ع",
		"ه", "خ", "ح", "ج", "ة", "ش", "س", "ي", "ب",
		"ل", "ا", "ت", "ن", "م", "ك", "ظ", "ط", "ذ",
		"د", "ز", "ر", "و"},
	Misspellings: AR_SPELLINGS,
	Homophones:   AR_HOMOPHONES,
	Homoglyphs: map[string][]string{
		"ض": []string{},
		"ص": []string{},
		"ث": []string{},
		"ق": []string{},
		"ف": []string{},
		"غ": []string{},
		"ع": []string{},
		"ه": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
		"خ": []string{},
		"ح": []string{},
		"ج": []string{},
		"ة": []string{},
		"ش": []string{},
		"س": []string{"vv", "ѡ", "ա", "ԝ"},
		"ي": []string{},
		"ب": []string{},
		"ل": []string{},
		"ا": []string{"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ"},
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

var AR_KEYBOARDS = []Keyboard{
	{
		Code:        "AR1",
		Name:        "غفقثصض",
		Description: "Arabic keyboard layout",
		Language:    ARABIC,
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
		Language:    ARABIC,
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
		Language:    ARABIC,
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
		Language:    ARABIC,
		Layout: []string{
			"١٢٣٤٥٦٧٨٩٠  ",
			"ظثةهيوطترعشق",
			"   لكجحغفدسا",
			"     منبذصخز"},
	},
}

func init() {
	KRegister(AR_KEYBOARDS)
}
