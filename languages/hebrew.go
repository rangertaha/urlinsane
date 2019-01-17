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

// Common Hebrew misspellings
var IW_SPELLINGS = [][]string{
	// []string{"misspelling", "correct1", "correct2"},
}

// Common Hebrew homophones
var IW_HOMOPHONES = [][]string{}

// Hebrew language
var IW_LANGUAGE = Language{
	Code: "IW",
	Name: "Hebrew",
	Numerals: map[string][]string{
		"0":  []string{""},
		"1":  []string{""},
		"2":  []string{""},
		"3":  []string{""},
		"4":  []string{""},
		"5":  []string{""},
		"6":  []string{""},
		"7":  []string{""},
		"8":  []string{""},
		"9":  []string{""},
		"10": []string{""},
	},
	Graphemes:    []string{""},
	Misspellings: IW_SPELLINGS,
	Homophones:   IW_HOMOPHONES,
	Homoglyphs: map[string][]string{
		"": []string{},
	},
}

var IW_KEYBOARDS = []Keyboard{
	{
		Code:        "IW1",
		Name:        "Hebrew",
		Description: "Hebrew standard layout",
		Language:    IW_LANGUAGE,
		Layout: []string{
			"1234567890 ",
			" פםןוטארק  ",
			" ףךלחיעכגדש",
			" ץתצמנהבסז "},
	},
}

func init() {
	KEYBOARDS.Add(IW_KEYBOARDS)
	KEYBOARDS.Append("IW", IW_KEYBOARDS)
	KEYBOARDS.Append("ALL", IW_KEYBOARDS)
}
