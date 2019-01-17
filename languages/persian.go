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

// Common Persian misspellings
var FA_SPELLINGS = [][]string{
	// []string{"misspelling", "correct1", "correct2"},
}

// Common Persian homophones
var FA_HOMOPHONES = [][]string{}

// Persian language
var FA_LANGUAGE = Language{
	Code: "FA",
	Name: "Persian",
	Numerals: map[string][]string{
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
	Graphemes:    []string{""},
	Misspellings: FA_SPELLINGS,
	Homophones:   FA_HOMOPHONES,
	Homoglyphs: map[string][]string{
		"": []string{},
	},
}

var FA_KEYBOARDS = []Keyboard{
	{
		Code:        "FA1",
		Name:        "Persian",
		Description: "Persian standard layout",
		Language:    FA_LANGUAGE,
		Layout: []string{
			"۱۲۳۴۵۶۷۸۹۰-  ",
			" چجحخهعغفقثصض",
			"  گکمنتالبیسش",
			"     وپدذرزطظ"},
	},
}

func init() {
	KEYBOARDS.Add(FA_KEYBOARDS)
	KEYBOARDS.Append("FA", FA_KEYBOARDS)
	KEYBOARDS.Append("ALL", FA_KEYBOARDS)
}
