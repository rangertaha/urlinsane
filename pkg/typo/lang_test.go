// Copyright 2024 Rangertaha. All Rights Reserved.
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
package typo

type Keyboard struct {
	lang   string
	code   string
	name   string
	layout []string
}

var (
	enMisspellings = [][]string{
		{"hwile", "while"},
		{"authrorities", "authorities"},
		{"emmision", "emission"},
		{"absail", "abseil"},
		{"inadquate", "inadequate"},
		{"adviced", "advised"},
		{"vyre", "very"},
		{"cassowarry", "cassowary"},
		{"abondon", "abandon"},
		{"proclaimation", "proclamation"},
		{"dominaton", "domination"},
	}

	enHomophones = [][]string{
		{"dot", "."},
		{"dash", "-"},
		{"accept", "except"},
		{"acclamation", "acclimation"},
		{"acts", "ax", "axe"},
		{"adds", "adz", "adze"},
		{"affect", "effect"},
		{"aid", "aide"},
		{"ail", "ale"},
		{"air", "ere", "heir"},
		{"aisle", "ill", "isle"},
		{"all", "awl"},
		{"allowed", "aloud"},
		{"allude", "elude"},
		{"bade", "bayed"},
		{"bail", "bale"},
		{"bait", "bate"},
		{"bald", "balled", "bawled"},
	}

	enAntonyms = map[string][]string{
		"asleep":     {"awake"},
		"attack":     {"defence", "protection", "defend"},
		"attic":      {"cellar"},
		"autumn":     {"spring"},
		"awake":      {"asleep"},
		"awful":      {"delicious", "nice", "pleasant"},
		"back":       {"front"},
		"background": {"foreground"},
		"backward":   {"forward"},
		"bad":        {"good"},
		"badluck":    {"fortune", "goodluck"},
		"beauty":     {"ugliness"},
		"before":     {"after"},
		"begin":      {"end", "finish"},
		"beginning":  {"end", "ending"},
		"behind":     {"infront"},
		"below":      {"above"},
		"best":       {"worst"},
	}

	enNumerals = map[string][]string{
		// Number: cardinal..,  ordinal.., other...
		"0": {"zero"},
		"1": {"one", "first"},
		"2": {"two", "second"},
		"3": {"three", "third"},
		"4": {"four", "fourth", "for"},
		"5": {"five", "fifth"},
		"6": {"six", "sixth"},
		"7": {"seven", "seventh"},
		"8": {"eight", "eighth"},
		"9": {"nine", "ninth"},
	}
	enGraphemes = []string{
		"a", "b", "c", "d", "e", "f", "g",
		"h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u",
		"v", "w", "x", "y", "z"}
	enVowels = []string{"a", "e", "i", "o", "u"}

	enHomoglyphs = map[string][]string{
		"a": {"à", "á", "â", "ã", "ä", "å", "ɑ", "а", "ạ", "ǎ", "ă", "ȧ", "ӓ", "٨"},
		"b": {"d", "lb", "ib", "ʙ", "Ь", `b̔"`, "ɓ", "Б"},
		"c": {"ϲ", "с", "ƈ", "ċ", "ć", "ç"},
		"d": {"b", "cl", "dl", "di", "ԁ", "ժ", "ɗ", "đ"},
		"e": {"é", "ê", "ë", "ē", "ĕ", "ě", "ė", "е", "ẹ", "ę", "є", "ϵ", "ҽ"},
		"f": {"Ϝ", "ƒ", "Ғ"},
		"g": {"q", "ɢ", "ɡ", "Ԍ", "Ԍ", "ġ", "ğ", "ց", "ǵ", "ģ"},
		"h": {"lh", "ih", "һ", "հ", "Ꮒ", "н"},
		"i": {"1", "l", "Ꭵ", "í", "ï", "ı", "ɩ", "ι", "ꙇ", "ǐ", "ĭ", "¡"},
		"j": {"ј", "ʝ", "ϳ", "ɉ"},
		"k": {"lk", "ik", "lc", "κ", "ⲕ", "κ"},
		"l": {"1", "i", "ɫ", "ł", "١", "ا", "", ""},
		"m": {"n", "nn", "rn", "rr", "ṃ", "ᴍ", "м", "ɱ"},
		"n": {"m", "r", "ń"},
		"o": {"0", "Ο", "ο", "О", "о", "Օ", "ȯ", "ọ", "ỏ", "ơ", "ó", "ö", "ӧ", "ه", "ة"},
		"p": {"ρ", "р", "ƿ", "Ϸ", "Þ"},
		"q": {"g", "զ", "ԛ", "գ", "ʠ"},
		"r": {"ʀ", "Г", "ᴦ", "ɼ", "ɽ"},
		"s": {"Ⴝ", "Ꮪ", "ʂ", "ś", "ѕ"},
		"t": {"τ", "т", "ţ"},
		"u": {"μ", "υ", "Ս", "ս", "ц", "ᴜ", "ǔ", "ŭ"},
		"v": {"ѵ", "ν", "v̇"},
		"w": {"vv", "ѡ", "ա", "ԝ"},
		"x": {"х", "ҳ", "ẋ"},
		"y": {"ʏ", "γ", "у", "Ү", "ý"},
		"z": {"ʐ", "ż", "ź", "ʐ", "ᴢ"},
	}

	enKeyboards = []Keyboard{
		{
			lang: "en",
			code: "en1",
			name: "QWERTY",
			layout: []string{
				"1234567890-",
				"qwertyuiop ",
				"asdfghjkl  ",
				"zxcvbnm    ",
			},
		},
		{
			lang: "en",
			code: "en2",
			name: "AZERTY",
			layout: []string{
				"1234567890",
				"azertyuiop",
				"qsdfghjklm",
				"wxcvbn    ",
			},
		},
		{
			lang: "en",
			code: "en3",
			name: "QWERTZ",
			layout: []string{
				"1234567890",
				"qwertzuiop",
				"asdfghjkl ",
				"yxcvbnm   ",
			},
		},
		{
			lang: "en",
			code: "en4",
			name: "DVORAK",
			layout: []string{
				"1234567890",
				"   pyfgcrl",
				"aoeuidhtns",
				" qjkxbmwvz",
			},
		},
	}

	tstTLDs = []string{
		"io",
		"uk",
		"co",
		"uk.com",
		"uk.io",
		"uk.eu.org",
	}
)
