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
var ES_SPELLINGS = [][]string{
// []string{"critising", "criticising", "criticizing"},
}

// Homophones
var ES_HOMOPHONES = [][]string{
// []string{"dot", "."},

}

var SPANISH = Language{
	Code: "ES",
	Name: "Spanish",
	Numerals: map[string][]string{
		"0":  []string{"zero"},
		"1":  []string{"one"},
		"2":  []string{"two"},
		"3":  []string{"three"},
		"4":  []string{"four", "for"},
		"5":  []string{"five"},
		"6":  []string{"six"},
		"7":  []string{"seven"},
		"8":  []string{"eight"},
		"9":  []string{"nine"},
		"10": []string{"ten"},
	},
	Graphemes: []string{
		"a", "b", "c", "d", "e", "f", "g",
		"h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u",
		"v", "w", "x", "y", "z"},
	Vowels:       []string{"a", "e", "i", "o", "u"},
	Misspellings: ES_SPELLINGS,
	Homophones:   ES_HOMOPHONES,
	Homoglyphs: map[string][]string{
		"a": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ", "٨"},
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
		"l": []string{"1", "i", "ɫ", "ł", "١", "ا", "", ""},
		"m": []string{"n", "nn", "rn", "rr", "ṃ", "ᴍ", "м", "ɱ"},
		"n": []string{"m", "r", "ń", "ñ"},
		"o": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ", "ه", "ة"},
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
		"ñ": []string{"n", "ń", "r"},
	},
}

var ES_KEYBOARDS = []Keyboard{
	{
		Code:        "es1",
		Name:        "QWERTY",
		Description: "Spanish keyboard layout",
		Language:    SPANISH,
		Layout: []string{
			"1234567890-",
			"qwertyuiop ",
			"asdfghjklñ ",
			"zxcvbnm  ç ",
		},
	},
	{
		Code:        "es2",
		Name:        "QWERTY",
		Description: "Spanish ISO keyboard layout",
		Language:    SPANISH,
		Layout: []string{
			"1234567890 ¡",
			"qwertyuiop  ",
			"asdfghjklñ  ",
			"zxcvbnm  -  ",
		},
	},
}

func init() {
	KRegister(ES_KEYBOARDS)
}
