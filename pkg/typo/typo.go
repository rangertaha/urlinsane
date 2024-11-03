// Copyright (C) 2024 Rangertaha
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

import (
	"fmt"
	"strings"

	"github.com/rangertaha/urlinsane/pkg/nlp"
)

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/rangertaha/urlinsane/internal/utils/datasets"
// )

// // Adjacent Characters

// // Inserting adjacent character from the keyboard

// // // Adjacent Character Insertion
// // func AdjacentCharacterInsertion(keyboard []string, word string) (typos []string) {
// // 	// for i, char := range word {
// // 	// 	for _, row := range keyboard {
// // 	// 		for _, kchar := range keyboard.Adjacent(string(char)) {
// // 	// 			variant := fmt.Sprint(word[:i], kchar, word[i+1:])
// // 	// 			typos = append(typos, variant)
// // 	// 		}
// // 	// 	}
// // 	// }
// // 	return
// // }

// // // Adjacent Character Substitution
// // func AdjacentCharacterSubstitution(keyboard []string, word string) (typos []string) {
// // 	// for i, char := range word {
// // 	// 	for _, row := range keyboard {
// // 	// 		for _, kchar := range keyboard.Adjacent(string(char)) {
// // 	// 			variant := fmt.Sprint(word[:i], kchar, word[i+1:])
// // 	// 			typos = append(typos, variant)
// // 	// 		}
// // 	// 	}
// // 	// }
// // 	return
// // }

// // // Adjacent Character Swapping
// // func AdjacentCharacterSwapping(keyboard []string, word string) (typos []string) {
// // 	// for i, char := range word {
// // 	// 	for _, row := range keyboard {
// // 	// 		for _, kchar := range keyboard.Adjacent(string(char)) {
// // 	// 			variant := fmt.Sprint(word[:i], kchar, word[i+1:])
// // 	// 			typos = append(typos, variant)
// // 	// 		}
// // 	// 	}
// // 	// }
// // 	return
// // }

// // Character substitution is the process of replacing a character or set of characters with another character or set of characters
// func CharSubstitution(name string, substitutes ...string) (names []string) {

// 	return
// }

// func CharOmission(name string, substitutes ...string) (names []string) {

// 	return
// }

// // Dot deletion typos are created by omitting a dot from the domain. For example, wwwgoogle.com and www.googlecom
// func DotDeletion(name string) (names []string) {
// 	return CharOmission(name, ".")
// }

// // www.google.com  ->  www-google.com
// //
// // Dot substitution typos are created by substituting a dot with a dash.
// func DotSubstitution(name string) (names []string) {
// 	return CharSubstitution(name, ".")
// }

// // // missingDotFunc typos are created by omitting a dot from the domain. For example, wwwgoogle.com and www.googlecom
// // func missingDotFunc(tc Result) (results []Result) {
// // 	for _, str := range missingCharFunc(tc.Original.String(), ".") {
// // 		if tc.Original.Domain != str {
// // 			dm := Domain{tc.Original.Subdomain, str, tc.Original.Suffix, Meta{}, false}
// // 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// // 		}
// // 	}
// // 	dm := Domain{tc.Original.Subdomain, strings.Replace(tc.Original.Domain, ".", "", -1), tc.Original.Suffix, Meta{}, false}
// // 	results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// // 	return results
// // }

func PrefixInsertion(name string, prefixes ...string) (names []string) {
	for _, prefix := range prefixes {
		names = append(names, prefix+name)
	}
	return
}

func SuffixInsertion(name string, suffixes ...string) (names []string) {
	for _, suffix := range suffixes {
		names = append(names, name+suffix)
	}
	return
}

func SubdomainInsertion(subdomains []string, name string) (names []string) {
	return PrefixInsertion(name, subdomains...)
}

func TldInsertion(subdomains []string, name string) (names []string) {
	return PrefixInsertion(name, subdomains...)
}

// // subdomainInsertionFunc typos are created by inserting common subdomains at the begining of the domain. wwwgoogle.com and ftpgoogle.com
// // func subdomainInsertionFunc(tc Result) (results []Result) {
// // 	for _, str := range datasets.SUBDOMAINS {
// // 		dm := Domain{tc.Original.Subdomain, str + tc.Original.Domain, tc.Original.Suffix, Meta{}, false}
// // 		results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// // 	}
// // 	return results
// // }

// // missingDashFunc typos are created by omitting a dash from the domain.
// // For example, www.a-b-c.com becomes www.ab-c.com, www.a-bc.com, and ww.abc.com
// // func missingDashFunc(tc Result) (results []Result) {
// // 	for _, str := range missingCharFunc(tc.Original.Domain, "-") {
// // 		if tc.Original.Domain != str {
// // 			dm := Domain{tc.Original.Subdomain, str, tc.Original.Suffix, Meta{}, false}
// // 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// // 		}
// // 	}
// // 	dm := Domain{tc.Original.Subdomain, strings.Replace(tc.Original.Domain, "-", "", -1), tc.Original.Suffix, Meta{}, false}
// // 	results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// // 	return results
// // }

// // characterOmissionFunc typos are when one character in the original domain name is omitted.
// // For example: www.exmple.com
// // func characterOmissionFunc(tc Result) (results []Result) {
// // 	for i := range tc.Original.Domain {
// // 		if i <= len(tc.Original.Domain)-1 {
// // 			domain := fmt.Sprint(
// // 				tc.Original.Domain[:i],
// // 				tc.Original.Domain[i+1:],
// // 			)
// // 			if tc.Original.Domain != domain {
// // 				dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// // 				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})

// //				}
// //			}
// //		}
// //		return results
// //	}

// // characterSwapFunc typos are when two consecutive characters are swapped in the original domain name.
// // Example: www.examlpe.com
// // func characterSwapFunc(tc Result) (results []Result) {
// // 	for i := range tc.Original.Domain {
// // 		if i <= len(tc.Original.Domain)-2 {
// // 			domain := fmt.Sprint(
// // 				tc.Original.Domain[:i],
// // 				string(tc.Original.Domain[i+1]),
// // 				string(tc.Original.Domain[i]),
// // 				tc.Original.Domain[i+2:],
// // 			)
// // 			if tc.Original.Domain != domain {
// // 				dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// // 				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// // 			}
// // 		}
// // 	}
// // 	return results
// // }

// characterSwapFunc typos are when two consecutive characters are swapped in the original domain name.
// Example: www.examlpe.com
func CharacterSwap(name string) (names []string) {
	for i := range name {
		if i <= len(name)-2 {
			variant := fmt.Sprint(name[:i], string(name[i+1]), string(name[i]), name[i+2:])
			if name != variant {
				names = append(names, variant)
			}
		}
	}
	return
}

// AdjacentCharacterSubstitution
// Adjacent character substitution typos occur when characters in the original
// domain name are replaced by neighboring characters on a specific keyboard
// layout. For example, www.ezample.com uses "z" instead of "x," substituting it
// with the adjacent character on a QWERTY keyboard.
func AdjacentCharacterSubstitution(name string, keyboard ...string) (names []string) {
	for i, char := range name {
		for _, key := range AdjacentCharacters(string(char), keyboard...) {
			variant := name[:i] + string(key) + name[i+1:]
			names = append(names, variant)
		}
	}
	return
}

// AdjacentCharacterInsertion
// Adjacent character insertion typos occur when characters adjacent of each
// letter are inserted. For example, www.googhle.com inserts "h" next to it's
// adjacent character "g" on a QWERTY keyboard.
func AdjacentCharacterInsertion(name string, keyboard ...string) (names []string) {
	for i, char := range name {
		for _, key := range AdjacentCharacters(string(char), keyboard...) {
			d1 := name[:i] + string(key) + string(char) + name[i+1:]
			names = append(names, d1)

			d2 := name[:i] + string(char) + string(key) + name[i+1:]
			names = append(names, d2)
		}
	}
	return
}

// HyphenInsertion
// Hyphen insertion typos occur when hyphens are inserted adjacent to each
// letter in a name. For example: "-example", "e-xample", "ex-ample", "exa-mple",
// "exam-ple", "examp-le", "example-"
func HyphenInsertion(name string) (names []string) {
	for i, char := range name {
		variant := name[:i] + "-" + string(char) + name[i+1:]
		if i == len(name)-1 {
			variant = name[:i] + string(char) + "-" + name[i+1:]
		}
		names = append(names, variant)
	}
	return
}

func HyphenOmission(name string) (names []string) {
	return CharacterDeletion(name, "-")
}

// DotInsertion
// Dot insertion typos occur when dots(.) are inserted the target name
// For example: "e.xample", "ex.ample", "exa.mple", "exam.ple", "examp.le"
func DotInsertion(name string) (names []string) {
	var nmap = map[string]bool{}
	for i, char := range name {
		variant := name[:i] + "." + string(char) + name[i+1:]
		if i == len(name)-1 {
			variant = name[:i] + string(char) + "." + name[i+1:]
		}
		variant = strings.Trim(variant, ".")
		nmap[variant] = true
		// names = append(names, variant)
	}

	for n := range nmap {
		names = append(names, n)
	}

	return
}

// DotOmission
// Dot ommission typos occur when dots(.) are left out of the target name
// For one.two.three: "one.twothree", "onetwo.three", "onetwothree",
func DotOmission(name string) (names []string) {
	return CharacterDeletion(name, ".")
}

// GraphemeInsertion
// Grapheme insertion also known as alphabet insertion where additional
// letters are inserted into a legitimate name to create a slightly modified
// version. For example: "aexample", "bexample", "cexample", "dexample", "eaxample"
func GraphemeInsertion(name string, graphemes ...string) (names []string) {
	alphabet := map[string]bool{}
	for _, a := range graphemes {
		alphabet[a] = true
	}
	for i, char := range name {
		for alp := range alphabet {
			variant := name[:i] + alp + string(char) + name[i+1:]
			if i == len(name)-1 {
				variant = name[:i] + string(char) + alp + name[i+1:]
			}
			names = append(names, variant)
		}
	}
	return
}

// GraphemeReplacement
// Grapheme replacement also known as alphabet replacement is where additional
// characters from the alphabet are replaced with characters from the target name
// to produce slightly modified version. For example: "axample", "bxample",
// "cxample", "dxample", "eaample"
func GraphemeReplacement(name string, graphemes ...string) (names []string) {
	alphabet := map[string]bool{}

	for _, a := range graphemes {
		alphabet[a] = true
	}

	for i := range name {
		for alp, _ := range alphabet {
			variant := name[:i] + alp + name[i+1:]

			if i == len(name)-1 {
				variant = name[:i] + alp + name[i+1:]
			}
			names = append(names, variant)
		}
	}
	return
}

// CharacterRepetition
// Character repetition typos are created by repeating a letter in the name.
// For example: "eexample", "exaample", "exammple", "examplee", "examplle"
func CharacterRepetition(name string) (names []string) {
	for i := range name {
		if i <= len(name) {
			variant := fmt.Sprint(name[:i], string(name[i]), string(name[i]), name[i+1:])
			if name != variant {
				names = append(names, variant)
			}
		}
	}
	return
}

// Example keyboard layout
//
//	var keyboard = []string{
//		"1234567890-",
//		"qwertyuiop ",
//		"asdfghjkl  ",
//		"zxcvbnm    ",
//	}
//
// DoubleCharacterAdjacentReplacement
// Double character adjacent replacement typos are created by replacing identical,
// consecutive letters in the name with adjacent keys on the keyboard.
// For example, www.gppgle.com and www.giigle.com.
func DoubleCharacterAdjacentReplacement(name string, keyboard ...string) (names []string) {
	// for _, keyboard := range tc.Keyboards {
	for i, char := range name {
		if i < len(name)-1 {
			if name[i] == name[i+1] {
				for _, key := range AdjacentCharacters(string(char), keyboard...) {
					variant := name[:i] + string(key) + string(key) + name[i+2:]

					names = append(names, variant)
				}
			}
		}
	}
	// }
	return
}

// CharacterOmission
// Grapheme omission leaves out one character from the name.
// For google: "gogle", "gogle", "googe", "googl", "goole", "oogle",
func CharacterOmission(name string) (names []string) {
	for i := range name {
		if i <= len(name)-1 {
			variant := fmt.Sprint(name[:i], name[i+1:])
			if name != variant {
				names = append(names, variant)
			}
		}
	}
	return
}

// The technique of creating typosquatting domains by switching between singular
// and plural forms of words is often referred to as Singular-Plural Substitution
// or Singular-Plural Manipulation.

// SingularPluralise
// For instance, if the original domain is 'example', a Singular-Plural
// Substitution typo might be 'examples', or vice versa. This subtle variation
// can make the fake domain look credible, especially when users are quickly
// scanning URLs.
func SingularPluraliseSubstitution(name string) (names []string) {
	pluralize := nlp.NewClient()
	if pluralize.IsPlural(name) {
		names = append(names, pluralize.Singular(name))
	}
	if pluralize.IsSingular(name) {
		names = append(names, pluralize.Plural(name))
	}

	return
}

func CharacterDeletion(name string, character string) (names []string) {
	var nmap = map[string]bool{}

	for i, char := range name {
		if character == string(char) {
			nmap[name[:i]+name[i+1:]] = true
			// names = append(names, name[:i]+name[i+1:])
		}
	}
	nmap[strings.Replace(name, character, "", -1)] = true

	for n := range nmap {
		names = append(names, n)
	}

	return
}

// CommonMisspellings
// Created from  common misspellings in the given language.
// For example, www.youtube.com becomes www.youtub.com and www.abseil.com
// becomes www.absail.com
func CommonMisspellings(name string, dataset ...[]string) (words []string) {
	words = []string{}
	for _, wordset := range dataset {
		for _, word := range wordset {
			if strings.Contains(name, word) {
				for _, w := range wordset {
					if w != word {
						words = append(words, strings.Replace(name, word, w, -1))
					}
				}

			}
		}
	}
	return
}

// VowelSwapping
// Created from vowels of the target name
// For example,
func VowelSwapping(name string, vowels ...string) (words []string) {
	for _, vchar := range vowels {
		if strings.Contains(name, vchar) {
			for _, vvchar := range vowels {
				new := strings.Replace(name, vchar, vvchar, -1)
				if new != name {
					words = append(words, new)
				}
			}
		}
	}
	return
}

// HomophoneSwapping
// homophonesFunc are created from sets of words that sound the same when spoken.
// For example, www.base.com becomes www .bass.com.
func HomophoneSwapping(name string, homophones ...[]string) (words []string) {
	words = []string{}
	for _, wordset := range homophones {
		for _, word := range wordset {
			if strings.Contains(name, word) {
				for _, w := range wordset {
					if w != word {
						words = append(words, strings.Replace(name, word, w, -1))
					}
				}

			}
		}
	}
	return
}

// HomoglyphSwapping
// Homoglyph swapping is a technique where visually similar characters, called
// homoglyphs, are swapped for one another in text. These characters look alike
// but are actually different in code, often coming from different alphabets
// or character sets. For example, an attacker might replace the letter "o" with
// the Cyrillic letter "о" (which looks nearly identical) in a URL or word. This
// can trick people into clicking a fraudulent link or misreading text.
func HomoglyphSwapping(name string, homoglyphs map[string][]string) (names []string) {
	for i, char := range name {
		for _, kchar := range SimilarChars(string(char), homoglyphs) {
			variant := fmt.Sprint(name[:i], kchar, name[i+1:])
			if name != variant {
				names = append(names, variant)
			}
		}
	}
	return
}

// WrongTopLevelDomain
func TopLevelDomain(tld string, tlds ...string) (records []string) {
	for _, suffix := range tlds {
		if len(strings.Split(suffix, ".")) == 1 {
			records = append(records, suffix)
		}
	}
	return
}

// SecondLevelDomain
func SecondLevelDomain(tld string, tlds ...string) (records []string) {
	for _, suffix := range tlds {
		if len(strings.Split(suffix, ".")) == 2 {
			records = append(records, suffix)
		}
	}
	return

}

// ThirdLevelDomain
func ThirdLevelDomain(tld string, tlds ...string) (records []string) {
	for _, suffix := range tlds {
		if len(strings.Split(suffix, ".")) == 3 {
			records = append(records, suffix)
		}
	}
	return
}

// BitFlipping
func BitFlipping(name string, graphemes ...string) (variations []string) {
	// Flip a single bit in a byte
	flipBit := func(b byte, pos uint) byte {
		mask := byte(1 << pos)
		return b ^ mask
	}

	// Flip each bit in each byte of the name
	for i := 0; i < len(name); i++ {
		for bit := 0; bit < 8; bit++ {
			flippedChar := flipBit(name[i], uint(bit))
			// Construct new variation
			variant := name[:i] + string(flippedChar) + name[i+1:]
			variations = append(variations, variant)
		}
	}
	return
}

func CardinalSwap(name string, numerals map[string][]string) (variations []string) {
	var fn func(map[string]string, string, bool) map[string]bool

	cardinals := NumeralMap(numerals, 0)

	fn = func(data map[string]string, str string, reverse bool) (names map[string]bool) {
		names = make(map[string]bool)

		for num, word := range data {
			{
				var variant string
				if !reverse {
					variant = strings.Replace(str, word, num, -1)
				} else {
					variant = strings.Replace(str, num, word, -1)
				}

				if str != variant {
					if _, ok := names[variant]; !ok {
						names[variant] = true
						for k, v := range fn(cardinals, variant, reverse) {
							names[k] = v
						}

						fn(cardinals, variant, reverse)
					}
				}
			}
		}
		return names
	}

	for name := range fn(cardinals, name, false) {
		variations = append(variations, name)
	}
	for name := range fn(cardinals, name, true) {
		variations = append(variations, name)
	}
	return
}

func OrdinalSwap(name string, numerals map[string][]string) (variations []string) {
	var fn func(map[string]string, string, bool) map[string]bool
	ordinals := NumeralMap(numerals, 1)

	fn = func(data map[string]string, str string, reverse bool) (names map[string]bool) {
		names = make(map[string]bool)

		for num, word := range data {
			{
				var variant string
				if !reverse {
					variant = strings.Replace(str, word, num, -1)
				} else {
					variant = strings.Replace(str, num, word, -1)
				}

				if str != variant {
					if _, ok := names[variant]; !ok {
						names[variant] = true
						for k, v := range fn(ordinals, variant, reverse) {
							names[k] = v
						}

						fn(ordinals, variant, reverse)
					}
				}
			}
		}
		return names
	}

	for name := range fn(ordinals, name, false) {
		variations = append(variations, name)
	}
	for name := range fn(ordinals, name, true) {
		variations = append(variations, name)
	}

	return
}

// SimilarChars returns homoglyphs, characters that look alike from other languages
func SimilarChars(key string, data map[string][]string) (chars []string) {
	chars = []string{}
	char, ok := data[key]
	if ok {
		chars = append(chars, char...)
	}
	return chars
}

// SimilarSounds returns common homophones, words that sound alike
func SimilarSounds(str string, data ...[]string) (words []string) {
	words = []string{}
	for _, wordset := range data {
		for _, word := range wordset {
			if strings.Contains(str, word) {
				for _, w := range wordset {
					if w != word {
						words = append(words, strings.Replace(str, word, w, -1))
					}
				}

			}
		}
	}
	return
}

func NumeralMap(data map[string][]string, pos int) (words map[string]string) {
	words = make(map[string]string)

	for num, names := range data {
		for i, name := range names {
			if i == pos {
				words[num] = name
				// words[name] = num
			}
		}
	}

	return
}

// Adjacent returns adjacent characters on a given keyboard
func AdjacentCharacters(char string, layout ...string) (chars []string) {
	chars = []string{}
	for r, row := range layout {
		for c := range row {
			var top, bottom, left, right string
			if char == string(layout[r][c]) {
				if r > 0 {
					top = string(layout[r-1][c])
					if top != " " {
						chars = append(chars, top)
					}
				}
				if r < len(layout)-1 {
					bottom = string(layout[r+1][c])
					if bottom != " " {
						chars = append(chars, bottom)
					}
				}
				if c > 0 {
					left = string(layout[r][c-1])
					if left != " " {
						chars = append(chars, left)
					}
				}
				if c < len(row)-1 {
					right = string(layout[r][c+1])
					if right != " " {
						chars = append(chars, right)
					}
				}
			}
		}
	}
	return chars
}
