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

func characterDeletion(token string, character string) (tokens []string) {
	var nmap = map[string]bool{}

	for i, char := range token {
		if character == string(char) {
			nmap[token[:i]+token[i+1:]] = true
			// tokens = append(tokens, token[:i]+token[i+1:])
		}
	}
	nmap[strings.Replace(token, character, "", -1)] = true

	for n := range nmap {
		tokens = append(tokens, n)
	}

	return
}

// PrefixInsertion creates tokens by prepending each prefix from the given
// list to the specified token. Example:
// Inputs:
//
//	prefixes = ["www", "ftp", "shop"]
//	token = "example"
//
// Outputs: ["wwwexample", "ftpexample", "shopexample"]
func PrefixInsertion(token string, prefixes ...string) (tokens []string) {
	for _, prefix := range prefixes {
		tokens = append(tokens, prefix+token)
	}
	return
}

// SuffixInsertion creates tokens by appending each suffix from the provided
// list to the end (right side) of the given token. Example:
// Inputs:
//
//	suffixes = ["com", "net", "io"]
//	token = "example"
//
// Outputs: ["examplecom", "examplenet", "exampleio"]
func SuffixInsertion(token string, suffixes ...string) (tokens []string) {
	for _, suffix := range suffixes {
		tokens = append(tokens, token+suffix)
	}
	return
}

// CharacterSwapping refers to a type of typo where two adjacent characters in 
// the original token are exchanged or swapped. This often occurs when characters 
// are unintentionally reversed in order, resulting in a misspelling.For example, 
// the word "example" could become "examlpe" by swapping the position of the 
// letters "l" and "p".
func CharacterSwapping(token string) (tokens []string) {
	for i := range token {
		if i <= len(token)-2 {
			variant := fmt.Sprint(token[:i], string(token[i+1]), string(token[i]), token[i+2:])
			if token != variant {
				tokens = append(tokens, variant)
			}
		}
	}
	return
}

// AdjacentCharacterSubstitution typos happen when a character in the original 
// token is mistakenly replaced by a neighboring character from the same keyboard 
// layout. This type of error often occurs due to hitting an adjacent key by accident. 
// For example, the token "ezample" contains a typo where the letter "x" is 
// substituted with "z," which is the neighboring key on an English QWERTY 
// keyboard layout.
func AdjacentCharacterSubstitution(token string, keyboard ...string) (tokens []string) {
	for i, char := range token {
		for _, key := range AdjacentCharacters(string(char), keyboard...) {
			variant := token[:i] + string(key) + token[i+1:]
			tokens = append(tokens, variant)
		}
	}
	return
}

// AdjacentCharacterInsertion typos occur when characters adjacent of each
// letter are inserted. For example, googhle inserts "h" next to it's
// adjacent character "g" on an English QWERTY keyboard layout.
func AdjacentCharacterInsertion(token string, keyboard ...string) (tokens []string) {
	for i, char := range token {
		for _, key := range AdjacentCharacters(string(char), keyboard...) {
			d1 := token[:i] + string(key) + string(char) + token[i+1:]
			tokens = append(tokens, d1)

			d2 := token[:i] + string(char) + string(key) + token[i+1:]
			tokens = append(tokens, d2)
		}
	}
	return
}

// HyphenInsertion typos happen when hyphens are mistakenly placed between 
// characters in a token, often occurring in various positions around the 
// letters. This type of error can lead to unnecessary fragmentation of the 
// word, with hyphens inserted at different points throughout the token. 
// For example, the word "example" might be incorrectly written as "-example",
//  "e-xample", "ex-ample", "exa-mple", "exam-ple", "examp-le", or even 
// "example-", with hyphens appearing before, between, or after the letters.
func HyphenInsertion(token string) (tokens []string) {
	for i, char := range token {
		variant := token[:i] + "-" + string(char) + token[i+1:]
		if i == len(token)-1 {
			variant = token[:i] + string(char) + "-" + token[i+1:]
			// variant = strings.Trim(variant, "-")
		}
		tokens = append(tokens, variant)
	}
	return
}

// HyphenOmission typos occur when hyphens are unintentionally left out of a 
// token, resulting in a version of the token that misses the expected hyphenation. 
// For example, the token "one-for-all" might be mistakenly written as "onefor-all", 
// "one-forall", or even "oneforall", where the hyphens are omitted.
func HyphenOmission(token string) (tokens []string) {
	return characterDeletion(token, "-")
}

// DotInsertion typos take place when periods (.) are mistakenly added at 
// various points within a token, leading to an incorrect placement of dots in 
// the target token. This type of error typically happens due to inadvertent 
// key presses or misplacement while typing. For instance, the word "example" 
// may be mistakenly written as "e.xample", "ex.ample", "exa.mple", "exam.ple", 
// or "examp.le", where the dot appears at different locations 
// within the token, disrupting the original structure.
func DotInsertion(token string) (tokens []string) {
	var nmap = map[string]bool{}
	for i, char := range token {
		variant := token[:i] + "." + string(char) + token[i+1:]
		if i == len(token)-1 {
			variant = token[:i] + string(char) + "." + token[i+1:]
		}
		variant = strings.Trim(variant, ".")
		nmap[variant] = true
		// tokens = append(tokens, variant)
	}

	for n := range nmap {
		tokens = append(tokens, n)
	}

	return
}

// DotOmission typos happen when periods (.) that should be present in the target 
// token are unintentionally omitted or left out. This type of error typically 
// occurs when the user fails to input the expected dots, often resulting in a 
// word or sequence that appears as a single string without proper separation. 
// For example, the sequence "one.two.three" might be mistakenly written 
// as "one.twothree", "onetwo.three", or even "onetwothree", where the dots 
// are missing between certain parts of the token, causing it to lose the 
// intended structure or meaning.
func DotOmission(token string) (tokens []string) {
	return characterDeletion(token, ".")
}

// GraphemeInsertion, also known as alphabet insertion, occurs when one or more 
// unintended letters are added to a valid token, leading to a modified or 
// misspelled version of the original token. These extra characters are typically 
// inserted either at the beginning or within the token, causing it to deviate 
// from its intended form. This type of error is often the result of a slip 
// of the finger or an accidental keystroke. For example, the token "example" 
// might be mistakenly typed as "aexample", "eaxample", "exaample", "examaple",
//  or "eaxampale", where additional letter like "a" are inserted throughout 
// the token, distorting its original structure.
func GraphemeInsertion(token string, graphemes ...string) (tokens []string) {
	alphabet := map[string]bool{}
	for _, a := range graphemes {
		alphabet[a] = true
	}
	for i, char := range token {
		for alp := range alphabet {
			variant := token[:i] + alp + string(char) + token[i+1:]
			if i == len(token)-1 {
				variant = token[:i] + string(char) + alp + token[i+1:]
			}
			tokens = append(tokens, variant)
		}
	}
	return
}

// GraphemeReplacement, also known as alphabet replacement, occurs when characters 
// from the original token are replaced by other letters from the alphabet, 
// resulting in a modified version of the token. This type of error typically leads 
// to small changes in the original token, where one or more letters are swapped 
// for different characters. For example, the token "example" could be mistakenly 
// written as "axample", "bxample", "cxample", "dxample", or "eaample", where 
// letters like "a", "b", "c", "d", or "e" are substituted, altering the 
// word slightly but keeping its general structure.
func GraphemeReplacement(token string, graphemes ...string) (tokens []string) {
	alphabet := map[string]bool{}

	for _, a := range graphemes {
		alphabet[a] = true
	}

	for i := range token {
		for alp := range alphabet {
			variant := token[:i] + alp + token[i+1:]

			if i == len(token)-1 {
				variant = token[:i] + alp + token[i+1:]
			}
			tokens = append(tokens, variant)
		}
	}
	return
}


// CharacterRepetition typos occur when a letter is unintentionally repeated 
// within a token, leading to a misspelled version. This type of error typically
// happens when a key is pressed twice or a letter is accidentally duplicated. 
// For example, the token "example" might be mistakenly written as "eexample", 
// "exaample", "exammple", "examplee", or "examplle", where one or more 
// characters are repeated, causing the token to diverge from its original form.
func CharacterRepetition(token string) (tokens []string) {
	for i := range token {
		if i <= len(token) {
			variant := fmt.Sprint(token[:i], string(token[i]), string(token[i]), token[i+1:])
			if token != variant {
				tokens = append(tokens, variant)
			}
		}
	}
	return
}


// RepetitionAdjacentReplacement typos occur when consecutive, identical letters 
// in a token are replaced with adjacent keys on the keyboard, resulting in a 
// slight alteration of the original word. This type of error often happens due 
// to accidental key presses of nearby characters. For example, the token 
// "google" might be mistakenly typed as "gppgle" or "giigle", where the repeated 
// letters are swapped with neighboring keys on the keyboard, causing the word 
// to be misspelled.
func RepetitionAdjacentReplacement(token string, keyboard ...string) (tokens []string) {
	// for _, keyboard := range tc.Keyboards {
	for i, char := range token {
		if i < len(token)-1 {
			if token[i] == token[i+1] {
				for _, key := range AdjacentCharacters(string(char), keyboard...) {
					variant := token[:i] + string(key) + string(key) + token[i+2:]

					tokens = append(tokens, variant)
				}
			}
		}
	}
	// }
	return
}


// CharacterOmission occurs when one character is unintentionally omitted from 
// the token, leading to an incomplete version of the original word. This type 
// of typo can happen when a key is accidentally skipped or overlooked while 
// typing. For example, the word "google" might be mistakenly written as "gogle", 
// "gogle", "googe", "googl", "goole", or "oogle", where a single character is 
// missing from different positions in the word, causing it to deviate from 
// the correct spelling.
func CharacterOmission(token string) (tokens []string) {
	for i := range token {
		if i <= len(token)-1 {
			variant := fmt.Sprint(token[:i], token[i+1:])
			if token != variant {
				tokens = append(tokens, variant)
			}
		}
	}
	return
}

// SingularPluralise typos are where a word is altered by switching between its 
// singular and plural forms. This subtle change can create a word that looks 
// similar to the original, but with a small variation that is easy to overlook. 
// For example, if the original word is 'example', a Singular-Plural might result
// in 'examples', or vice versa.  
func SingularPluralise(token string) (tokens []string) {
	pluralize := nlp.NewClient()
	if pluralize.IsPlural(token) {
		tokens = append(tokens, pluralize.Singular(token))
	}
	if pluralize.IsSingular(token) {
		tokens = append(tokens, pluralize.Plural(token))
	}

	return
}





// CommonMisspellings refers to typos created by frequent spelling errors or 
// missteps that occur in the target language. These errors often involve slight
// changes to the spelling of a word, making them appear similar to the original 
// but incorrect. For instance, the word "youtube" could be mistyped as 
// "youtub", and "abseil" could become "absail", where common mistakes in 
// spelling lead to slightly altered but recognizable versions of the original.
func CommonMisspellings(token string, dataset ...[]string) (words []string) {
	words = []string{}
	for _, wordset := range dataset {
		for _, word := range wordset {
			if strings.Contains(token, word) {
				for _, w := range wordset {
					if w != word {
						words = append(words, strings.Replace(token, word, w, -1))
					}
				}

			}
		}
	}
	return
}


// VowelSwapping occurs when the vowels in the target token are swapped with each other, leading to a slightly altered version of the original word. 
// This type of error typically involves exchanging one vowel for another, which can still make the altered token look similar to the original, 
// but with a subtle change. For example, the word "example" could become "ixample", "exomple", or "exaple", where vowels like "a", "e", and "o" 
// are swapped, causing the token to differ from its correct form.
func VowelSwapping(token string, vowels ...string) (words []string) {
	for _, vchar := range vowels {
		if strings.Contains(token, vchar) {
			for _, vvchar := range vowels {
				new := strings.Replace(token, vchar, vvchar, -1)
				if new != token {
					words = append(words, new)
				}
			}
		}
	}
	return
}


// HomophoneSwapping occurs when words that sound the same but have different 
// meanings or spellings are substituted for one another. This type of error 
// arises from words that are homophones—words that are pronounced the same but 
// may differ in spelling or meaning. For example, the word "base" could be 
// swapped with "bass", where "base" and "bass" are homophones, making the 
// altered word sound the same when spoken, yet look different in writing.
func HomophoneSwapping(token string, homophones ...[]string) (words []string) {
	words = []string{}
	for _, wordset := range homophones {
		for _, word := range wordset {
			if strings.Contains(token, word) {
				for _, w := range wordset {
					if w != word {
						words = append(words, strings.Replace(token, word, w, -1))
					}
				}

			}
		}
	}
	return
}

// HomoglyphSwapping is a technique where visually similar characters, called
// homoglyphs, are swapped for one another in text. These characters look alike
// but are actually different in code, often coming from different alphabets
// or character sets. For example, an attacker might replace the letter "o" with
// the Cyrillic letter "о" (which looks nearly identical) in a URL or word. This
// can trick people into clicking a fraudulent link or misreading text.
func HomoglyphSwapping(token string, homoglyphs map[string][]string) (tokens []string) {
	for i, char := range token {
		for _, kchar := range SimilarChars(string(char), homoglyphs) {
			variant := fmt.Sprint(token[:i], kchar, token[i+1:])
			if token != variant {
				tokens = append(tokens, variant)
			}
		}
	}
	return
}


// BitFlipping
func BitFlipping(token string, graphemes ...string) (variations []string) {
	// Flip a single bit in a byte
	flipBit := func(b byte, pos uint) byte {
		mask := byte(1 << pos)
		return b ^ mask
	}

	// Flip each bit in each byte of the token
	for i := 0; i < len(token); i++ {
		for bit := 0; bit < 8; bit++ {
			flippedChar := flipBit(token[i], uint(bit))
			// Construct new variation
			variant := token[:i] + string(flippedChar) + token[i+1:]
			variations = append(variations, variant)
		}
	}
	return
}

func CardinalSwap(token string, numerals map[string][]string) (variations []string) {
	var fn func(map[string]string, string, bool) map[string]bool

	cardinals := NumeralMap(numerals, 0)

	fn = func(data map[string]string, str string, reverse bool) (tokens map[string]bool) {
		tokens = make(map[string]bool)

		for num, word := range data {
			{
				var variant string
				if !reverse {
					variant = strings.Replace(str, word, num, -1)
				} else {
					variant = strings.Replace(str, num, word, -1)
				}

				if str != variant {
					if _, ok := tokens[variant]; !ok {
						tokens[variant] = true
						for k, v := range fn(cardinals, variant, reverse) {
							tokens[k] = v
						}

						fn(cardinals, variant, reverse)
					}
				}
			}
		}
		return tokens
	}

	for token := range fn(cardinals, token, false) {
		variations = append(variations, token)
	}
	for token := range fn(cardinals, token, true) {
		variations = append(variations, token)
	}
	return
}

func OrdinalSwap(token string, numerals map[string][]string) (variations []string) {
	var fn func(map[string]string, string, bool) map[string]bool
	ordinals := NumeralMap(numerals, 1)

	fn = func(data map[string]string, str string, reverse bool) (tokens map[string]bool) {
		tokens = make(map[string]bool)

		for num, word := range data {
			{
				var variant string
				if !reverse {
					variant = strings.Replace(str, word, num, -1)
				} else {
					variant = strings.Replace(str, num, word, -1)
				}

				if str != variant {
					if _, ok := tokens[variant]; !ok {
						tokens[variant] = true
						for k, v := range fn(ordinals, variant, reverse) {
							tokens[k] = v
						}

						fn(ordinals, variant, reverse)
					}
				}
			}
		}
		return tokens
	}

	for token := range fn(ordinals, token, false) {
		variations = append(variations, token)
	}
	for token := range fn(ordinals, token, true) {
		variations = append(variations, token)
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

// DHUSubstitution substitutes (.-_) in a given token to produce variants that
// look similar. Primarily used in package tokens
func DHUSubstitution(token string) (variations []string) {

	return
}

// DotHyphenSubstitution substitutes (.-) in a given token to produce variants that
// look similar.
func DotHyphenSubstitution(token string) (variations []string) {

	return
}

func NumeralMap(data map[string][]string, pos int) (words map[string]string) {
	words = make(map[string]string)

	for num, tokens := range data {
		for i, token := range tokens {
			if i == pos {
				words[num] = token
				// words[token] = num
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
