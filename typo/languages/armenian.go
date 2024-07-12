// Copyright (C) 2024  Rangertaha <rangertaha@gmail.com>
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

// https://en.wikipedia.org/wiki/Armenian_alphabet

var (
	// hyMisspellings are common misspellings
	hyMisspellings = [][]string{
		[]string{"", ""},
	}

	// hyHomophones are words that sound alike
	hyHomophones = [][]string{
		[]string{"կետը", "."},
	}

	// hyAntonyms are words opposite in meaning to another (e.g. bad and good ).
	hyAntonyms = map[string][]string{
		"լավ": []string{"վատը"},
	}

	hyLanguage = Language{
		// https://www.loc.gov/standards/iso639-2/php/code_list.php
		Code: "HY",
		Name: "Armenian",
		// http://mylanguages.org/armenian_numbers.php
		Numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":       []string{"զրո"},
			"1":       []string{"մեկ"},
			"2":       []string{"երկու"},
			"3":       []string{"երեք"},
			"4":       []string{"չորս"},
			"5":       []string{"հինգ"},
			"6":       []string{"վեց"},
			"7":       []string{"յոթ"},
			"8":       []string{"ութ"},
			"9":       []string{"ինը"},
			"10":      []string{"տաս"},
			"11":      []string{"տասնմեկ"},
			"12":      []string{"տասներկու"},
			"13":      []string{"տասներեք"},
			"14":      []string{"տասնչորս"},
			"15":      []string{"տասնհինգ"},
			"16":      []string{"տասնվեց"},
			"17":      []string{"տասնյոթ"},
			"18":      []string{"տասնութ"},
			"19":      []string{"տասնիննը"},
			"20":      []string{"քսան"},
			"100":     []string{"հարյուր"},
			"1000":    []string{"հազար"},
			"1000000": []string{"միլիոն"},
		},
		// http://mylanguages.org/armenian_alphabet.php
		Graphemes: []string{
			"ա", "բ", "գ", "դ", "ե", "զ", "է", "ը",
			"թ", "ժ", "ի", "լ", "խ", "ծ", "կ", "հ",
			"ձ", "ղ", "ճ", "մ", "յ", "ն", "շ", "ո",
			"չ", "պ", "ջ", "ռ", "ս", "վ", "տ", "ր",
			"ց", "փ", "ք", "և", "օ", "ֆ",
		},
		Vowels:       []string{},
		Misspellings: hyMisspellings,
		Homophones:   hyHomophones,
		Antonyms:     hyAntonyms,
		Homoglyphs: map[string][]string{
			//"a": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ", "٨"},
			//"b": []string{"d", "lb", "ib", "ʙ", "Ь", `b̔"`, "ɓ", "Б"},
			//"c": []string{"ϲ", "с", "ƈ", "ċ", "ć", "ç"},
			//"d": []string{"b", "cl", "dl", "di", "ԁ", "ժ", "ɗ", "đ"},
			//"e": []string{"é", "ê", "ë", "ē", "ĕ", "ě", "ė", "е", "ẹ", "ę", "є", "ϵ", "ҽ"},
			//"f": []string{"Ϝ", "ƒ", "Ғ"},
			//"g": []string{"q", "ɢ", "ɡ", "Ԍ", "Ԍ", "ġ", "ğ", "ց", "ǵ", "ģ"},
			//"h": []string{"lh", "ih", "һ", "հ", "Ꮒ", "н"},
			//"i": []string{"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			//"j": []string{"ј", "ʝ", "ϳ", "ɉ"},
			//"k": []string{"lk", "ik", "lc", "κ", "ⲕ", "κ"},
			//"l": []string{"1", "i", "ɫ", "ł", "١", "ا", "", ""},
			//"m": []string{"n", "nn", "rn", "rr", "ṃ", "ᴍ", "м", "ɱ"},
			//"n": []string{"m", "r", "ń"},
			//"o": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ", "ه", "ة"},
			//"p": []string{"ρ", "р", "ƿ", "Ϸ", "Þ"},
			//"q": []string{"g", "զ", "ԛ", "գ", "ʠ"},
			//"r": []string{"ʀ", "Г", "ᴦ", "ɼ", "ɽ"},
			//"s": []string{"Ⴝ", "Ꮪ", "ʂ", "ś", "ѕ"},
			//"t": []string{"τ", "т", "ţ"},
			//"u": []string{"μ", "υ", "Ս", "ս", "ц", "ᴜ", "ǔ", "ŭ"},
			//"v": []string{"ѵ", "ν", "v̇"},
			//"w": []string{"vv", "ѡ", "ա", "ԝ"},
			//"x": []string{"х", "ҳ", "ẋ"},
			//"y": []string{"ʏ", "γ", "у", "Ү", "ý"},
			//"z": []string{"ʐ", "ż", "ź", "ʐ", "ᴢ"},

			"ա": []string{"vv", "ѡ", "ա", "ԝ"},
			"բ": []string{"ρ", "р", "ƿ", "Ϸ", "Þ"},
			"գ": []string{},
			"դ": []string{},
			"ե": []string{"d", "lb", "ib", "ʙ", "Ь", `b̔"`, "ɓ", "Б"},
			"զ": []string{},
			"է": []string{},
			"ը": []string{},
			"թ": []string{},
			"ժ": []string{"b", "cl", "dl", "di", "ԁ", "ժ", "ɗ", "đ"},
			"ի": []string{},
			"լ": []string{"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			"խ": []string{},
			"ծ": []string{},
			"կ": []string{},
			"հ": []string{"lh", "ih", "һ", "հ", "Ꮒ", "н"},
			"ձ": []string{},
			"ղ": []string{},
			"ճ": []string{"6"},
			"մ": []string{},
			"յ": []string{},
			"ն": []string{},
			"շ": []string{"2", "չ", "ջ"},
			"ո": []string{"m", "r", "ń"},
			"չ": []string{"2", "շ", "ջ"},
			"պ": []string{},
			"ջ": []string{"2", "չ", "շ"},
			"ռ": []string{},
			"ս": []string{"μ", "υ", "Ս", "ц", "ᴜ", "ǔ", "ŭ", "u"},
			"վ": []string{},
			"տ": []string{"un"},
			"ր": []string{},
			"ց": []string{"q", "ɢ", "ɡ", "Ԍ", "Ԍ", "ġ", "ğ", "g", "ǵ", "ģ"},
			"փ": []string{},
			"ք": []string{"ρ", "р", "ƿ", "Ϸ", "Þ"},
			"և": []string{},
			"օ": []string{"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ", "ه", "ة"},
			"ֆ": []string{},
		},
	}

	hyKeyboards = []Keyboard{
		{
			Code:        "HY1",
			Name:        "HM QWERTY",
			Description: "Armenian QWERTY keyboard layout",
			Language:    hyLanguage,
			Layout: []string{
				"1234567890-",
				"ճւերտյւիոպ ",
				"ասդֆգհձկլ  ",
				" զխծվբնմ   ",
			},
		},
		{
			Code:        "HY2",
			Name:        "Armenian, Western QWERTY",
			Description: "Armenian, Western QWERTY keyboard layout",
			Language:    hyLanguage,
			Layout: []string{
				" ձյ՛ -   օռժ",
				"խվէրդեըիոբչջ",
				"աստֆկհճքլթփ ",
				" զցգւպնմշղծ ",
			},
		},
		//{
		//	Code:        "HY3",
		//	Name:        "Easter QWERTY",
		//	Description: "Easter QWERTY keyboard layout",
		//	Language:    ENGLISH,
		//	Layout: []string{
		//		"",
		//		"",
		//		"",
		//		"",
		//	},
		//},

	}
)

func init() {
	KEYBOARDS.Add(hyKeyboards)
	KEYBOARDS.Append("HY", hyKeyboards)
	KEYBOARDS.Append("ALL", hyKeyboards)
}
