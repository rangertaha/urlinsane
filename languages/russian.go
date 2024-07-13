// Copyright (C) 2024  Tal Hatchi (Rangertaha)
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
			// Language:    ruLanguage,
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
			// Language:    ruLanguage,
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
			// Language:    ruLanguage,
			Layout: []string{
				"1234567890- ",
				"йцукенгшщзхъ",
				"фывапролджэ ",
				"ячсмитьбю   "},
		},
	}
)

func init() {
	Add("RU", ruLanguage)
}
