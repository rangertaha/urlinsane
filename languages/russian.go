// Copyright © 2019 Tal Hachi
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
	// ruMisspellings are common misspellings
	ruMisspellings = [][]string{
		[]string{"", ""},
	}

	// ruHomophones are words that sound alike
	ruHomophones = [][]string{
		[]string{"точка", "."},
	}

	// ruAntonyms are words opposite in meaning to another (e.g. bad and good ).
	ruAntonyms = map[string][]string{
		"хорошо": []string{"плохой"},
	}

	ruLanguage = Language{
		Code: "RU",
		Name: "Russian",
		// http://www.russianlessons.net/lessons/lesson2_main.php
		Numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":          []string{"ноль"},
			"1":          []string{"один", "первый"},
			"2":          []string{"два", "второй"},
			"3":          []string{"три", "в третьих"},
			"4":          []string{"четыре", "четвертый"},
			"5":          []string{"пять", "пятый"},
			"6":          []string{"шесть", "шестой"},
			"7":          []string{"семь", "Седьмой"},
			"8":          []string{"восемь", "Восьмой"},
			"9":          []string{"девять", "девятый"},
			"10":         []string{"десять", "десятый"},
			"11":         []string{"одиннадцать"},
			"12":         []string{"двенадцать"},
			"13":         []string{"тринадцать"},
			"14":         []string{"четырнадцать"},
			"15":         []string{"пятнадцать"},
			"16":         []string{"шестнадцать"},
			"17":         []string{"семнадцать"},
			"18":         []string{"восемнадцать"},
			"19":         []string{"девятнадцать"},
			"20":         []string{"двадцать"},
			"21":         []string{"двадцатьодин"},
			"22":         []string{"двадцатьдва"},
			"23":         []string{"двадцатьтри"},
			"24":         []string{"двадцатьчетыре"},
			"30":         []string{"тридцать"},
			"40":         []string{"сорок"},
			"50":         []string{"пятьдесят"},
			"60":         []string{"шестьдесят"},
			"70":         []string{"семьдесят"},
			"80":         []string{"восемьдесят"},
			"90":         []string{"девяносто"},
			"100":        []string{"сто"},
			"200":        []string{"двести"},
			"300":        []string{"триста"},
			"400":        []string{"четыреста"},
			"500":        []string{"пятьсот"},
			"600":        []string{"шестьсот"},
			"700":        []string{"семьсот"},
			"800":        []string{"восемьсот"},
			"900":        []string{"девятьсот"},
			"1000":       []string{"тысяча"},
			"1000000":    []string{"миллион"},
			"1000000000": []string{"миллиард"},
		},
		Graphemes: []string{
			"а", "б", "в", "г", "д", "е", "ё",
			"ж", "з", "и", "й", "к", "л", "м",
			"н", "о", "п", "р", "с", "т", "у",
			"ф", "х", "ц", "ч", "ш", "щ", "ъ",
			"ы", "ь", "э", "ю", "я", "ѕ", "ѯ",
			"ѱ", "ѡ", "ѫ", "ѧ", "ѭ", "ѩ"},
		Vowels:       []string{"a", "о", "у", "э", "ы", "я", "ё", "ю", "е", "и"},
		Misspellings: ruMisspellings,
		Homophones:   ruHomophones,
		Antonyms:     ruAntonyms,
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

	ruKeyboards = []Keyboard{
		{
			Code:        "RU1",
			Name:        "ЙЦУКЕН",
			Description: "Russian keyboard layout",
			Language:    ruLanguage,
			Layout: []string{
				"1234567890   ",
				"йцукенгшщзхъё",
				"фывапролджэ  ",
				"ячсмитьбю    "},
		},
		{
			Code:        "RU2",
			Name:        "ЯШЕРТЫ",
			Description: "Phonetic Russian keybaord layout",
			Language:    ruLanguage,
			Layout: []string{
				"1234567890ьъ ",
				"яшертыуиопюжэ",
				"асдфгчйкл    ",
				"зхцвбнм      "},
		},
		{
			Code:        "RU3",
			Name:        "ЙЦУКЕН",
			Description: "PC Russian keyboard layout",
			Language:    ruLanguage,
			Layout: []string{
				"1234567890- ",
				"йцукенгшщзхъ",
				"фывапролджэ ",
				"ячсмитьбю   "},
		},
	}
)

func init() {
	KEYBOARDS.Add(ruKeyboards)
	KEYBOARDS.Append("RU", ruKeyboards)
	KEYBOARDS.Append("ALL", ruKeyboards)
}
