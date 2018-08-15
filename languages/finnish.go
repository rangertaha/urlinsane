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

// Common misspellings and homophones from https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/For_machines
// The first word is a misspelling of a word. All words that follows are correct spellings or homophones
var FI_SPELLINGS = [][]string{
// []string{"misspelling", "correct1", "correct2"},
}
var FI_HOMOPHONES = [][]string{
// []string{"vary", "very"},
}
var FI_LANGUAGE = Language{
	Code: "fi",
	Name: "Finnish",
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
		"a", "b", "c", "d", "e", "f", "g",
		"h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u",
		"v", "w", "x", "y", "z", "å", "ä", "ö"},
	Misspellings: FI_SPELLINGS,
	Homophones:   FI_HOMOPHONES,
	Homoglyphs: map[string][]string{
		"a": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ"},
		"b": []string{"d", "lb", "ib", "ʙ", "Ь", `b̔"`, "ɓ", "Б"},
		"c": []string{"ϲ", "с", "ƈ", "ċ", "ć", "ç"},
		"d": []string{"b", "cl", "dl", "di", "ԁ", "ժ", "ɗ", "đ"},
		"e": []string{"é", "ê", "ë", "ē", "ĕ", "ě", "ė", "е", "ẹ", "ę", "є", "ϵ", "ҽ"},
		"f": []string{"Ϝ", "ƒ", "Ғ"},
		"g": []string{"q", "ɢ", "ɡ", "Ԍ", "Ԍ", "ġ", "ğ", "ց", "ǵ", "ģ"},
		"h": []string{"lh", "ih", "һ", "հ", "Ꮒ", "н"},
		"i": []string{"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ"},
		"j": []string{"ј", "ʝ", "ϳ", "ɉ"},
		"k": []string{"lk", "ik", "lc", "κ", "ⲕ", "κ"},
		"l": []string{"1", "i", "ɫ", "ł"},
		"m": []string{"n", "nn", "rn", "rr", "ṃ", "ᴍ", "м", "ɱ"},
		"n": []string{"m", "r", "ń"},
		"o": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
		"p": []string{"ρ", "р", "ƿ", "Ϸ", "Þ"},
		"q": []string{"g", "զ", "ԛ", "գ", "ʠ"},
		"r": []string{"ʀ", "Г", "ᴦ", "ɼ", "ɽ"},
		"s": []string{"Ⴝ", "Ꮪ", "ʂ", "ś", "ѕ"},
		"t": []string{"τ", "т", "ţ"},
		"u": []string{"μ", "υ", "Ս", "ս", "ц", "ᴜ", "ǔ", "ŭ"},
		"v": []string{"ѵ", "ν", "v̇"},
		"w": []string{"vv", "ѡ", "ա", "ԝ"},
		"x": []string{"х", "ҳ", "ẋ"},
		"y": []string{"ʏ", "γ", "у", "Ү", "ý"},
		"z": []string{"ʐ", "ż", "ź", "ʐ", "ᴢ"},
		"å": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ"},
		"ä": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ"},
		"ö": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ"},
	},
}

var FI_KEYBOARDS = []Keyboard{
	{
		Code:     "fi1",
		Name:     "QWERTY",
		Language: FI_LANGUAGE,
		Layout: []string{
			"1234567890 ",
			"qwertyuiopå",
			"asdfghjklöä",
			"zxcvbnm  - ",
		},
	},
}

func init() {
	KBRegister(FI_KEYBOARDS)
}
