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
	// Common misspellings from https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/For_machines
	// updated (25/April/2018)
	langMisspellings = [][]string{
		[]string{"", ""},
	}

	// Homophones are words that sound alike. See https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/Homophones (25/June/2018)
	langHomophones = [][]string{
		[]string{"dot", "."},
	}

	// enAntonyms are words opposite in meaning to another (e.g. bad and good ).
	langAntonyms = map[string][]string{
		"about": []string{"exactly"},
	}

	langLanguage = Language{
		Code:        "EN",
		Name:        "English",
		Description: "English language",

		Numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0": []string{"zero"},
			"1": []string{"one", "first"},
			"2": []string{"two", "second"},
			"3": []string{"three", "third"},
		},
		Graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z"},
		Vowels:       []string{"a", "e", "i", "o", "u"},
		Misspellings: langMisspellings,
		Homophones:   langHomophones,
		Antonyms:     langAntonyms,
		Homoglyphs: map[string][]string{
			"a": []string{"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ", "٨"},
			"b": []string{"d", "lb", "ib", "ʙ", "Ь", `b̔"`, "ɓ", "Б"},
			"c": []string{"ϲ", "с", "ƈ", "ċ", "ć", "ç"},
			"d": []string{"b", "cl", "dl", "di", "ԁ", "ժ", "ɗ", "đ"},
			"e": []string{"é", "ê", "ë", "ē", "ĕ", "ě", "ė", "е", "ẹ", "ę", "є", "ϵ", "ҽ"},
			"f": []string{"Ϝ", "ƒ", "Ғ"},
			"g": []string{"q", "ɢ", "ɡ", "Ԍ", "Ԍ", "ġ", "ğ", "ց", "ǵ", "ģ"},
			"h": []string{"lh", "ih", "һ", "հ", "Ꮒ", "н"},
			"i": []string{"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
			"j": []string{"ј", "ʝ", "ϳ", "ɉ"},
			"k": []string{"lk", "ik", "lc", "κ", "ⲕ", "κ"},
			"l": []string{"1", "i", "ɫ", "ł", "١", "ا", "", ""},
			"m": []string{"n", "nn", "rn", "rr", "ṃ", "ᴍ", "м", "ɱ"},
			"n": []string{"m", "r", "ń"},
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
		},
		Keyboards: []Keyboard{
			{
				Code:        "EN1",
				Name:        "QWERTY",
				Description: "English QWERTY keyboard layout",
				Layout: []string{
					"1234567890-",
					"qwertyuiop ",
					"asdfghjkl  ",
					"zxcvbnm    ",
				},
			},
			{
				Code:        "EN2",
				Name:        "AZERTY",
				Description: "English AZERTY keyboard layout",
				Layout: []string{
					"1234567890",
					"azertyuiop",
					"qsdfghjklm",
					"wxcvbn    ",
				},
			},
			{
				Code:        "EN3",
				Name:        "QWERTZ",
				Description: "English QWERTZ keyboard layout",
				Layout: []string{
					"1234567890",
					"qwertzuiop",
					"asdfghjkl ",
					"yxcvbnm   ",
				},
			},
			{
				Code:        "EN4",
				Name:        "DVORAK",
				Description: "English DVORAK keyboard layout",
				Layout: []string{
					"1234567890",
					"   pyfgcrl",
					"aoeuidhtns",
					" qjkxbmwvz",
				},
			},
		},
	}
)

func init() {
	// Add("ID", langLanguage)
}
