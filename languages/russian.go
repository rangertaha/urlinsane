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

// Common misspellings from https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/For_machines
var RU_SPELLINGS = [][]string{
// []string{"misspelling", "correct1", "correct2"},
}
var RU_HOMOPHONES = [][]string{
// []string{"vary", "very"},
}

var RUSSIAN = Language{
	Code: "ru",
	Name: "Russian",
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
	Graphemes: []string{
		"а", "б", "в", "г", "д", "е", "ё",
		"ж", "з", "и", "й", "к", "л", "м",
		"н", "о", "п", "р", "с", "т", "у",
		"ф", "х", "ц", "ч", "ш", "щ", "ъ",
		"ы", "ь", "э", "ю", "я", "ѕ", "ѯ",
		"ѱ", "ѡ", "ѫ", "ѧ", "ѭ", "ѩ"},
	Misspellings: RU_SPELLINGS,
	Homophones:   RU_HOMOPHONES,
	Homoglyphs: map[string][]string{
		"а": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ"},
		"б": []string{"6", "b", "Ь", `b̔"`, "ɓ", "Б"},
		"в": []string{"B"},
		"г": []string{"ʀ", "Г", "ᴦ", "ɼ", "ɽ"},
		"д": []string{""},
		"е": []string{""},
		"ё": []string{""},
		"ж": []string{""},
		"з": []string{""},
		"и": []string{""},
		"й": []string{""},
		"к": []string{""},
		"л": []string{""},
		"м": []string{""},
		"н": []string{""},
		"о": []string{""},
		"п": []string{""},
		"р": []string{""},
		"с": []string{""},
		"т": []string{""},
		"у": []string{""},
		"ф": []string{""},
		"х": []string{""},
		"ц": []string{""},
		"ч": []string{""},
		"ш": []string{""},
		"щ": []string{""},
		"ъ": []string{""},
		"ы": []string{""},
		"ь": []string{""},
		"э": []string{""},
		"ю": []string{""},
		"я": []string{""},
		"ѕ": []string{""},
		"ѯ": []string{""},
		"ѱ": []string{""},
		"ѡ": []string{""},
		"ѫ": []string{""},
		"ѧ": []string{""},
		"ѭ": []string{""},
		"ѩ": []string{""},
	},
}

var RU_KEYBOARDS = []Keyboard{
	{
		Code:     "ru1",
		Name:     "Russian",
		Language: RUSSIAN,
		Layout: []string{
			"1234567890   ",
			"йцукенгшщзхъё",
			"фывапролджэ  ",
			"ячсмитьбю    "},
	},
	{
		Code:     "ru2",
		Name:     "Russian - Phonetic",
		Language: RUSSIAN,
		Layout: []string{
			"1234567890ьъ ",
			"яшертыуиопюжэ",
			"асдфгчйкл    ",
			"зхцвбнм      "},
	},
	{
		Code:     "ru3",
		Name:     "Russian - PC",
		Language: RUSSIAN,
		Layout: []string{
			"1234567890- ",
			"йцукенгшщзхъ",
			"фывапролджэ ",
			"ячсмитьбю   "},
	},
}

func init() {
	RegisterKeyboard(RU_KEYBOARDS)
}
