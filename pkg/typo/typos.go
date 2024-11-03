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

// func SubdomainInsertion(subdomains []string, name string) (names []string) {
// 	return PrefixInsertion(name, subdomains...)
// }

// func TldInsertion(subdomains []string, name string) (names []string) {
// 	return PrefixInsertion(name, subdomains...)
// }

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

// GraphemeInsertion
// Graphemes insertion typosquatting is a type of typosquatting where additional
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
// Grapheme replacement is where additional characters are replaced with
// characters from a list of graphemes to produce slightly modified version.
// For example: "axample", "bxample", "cxample", "dxample", "eaample"
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


// GraphemeRepetition
// Grapheme repetition typos are created by repeating a letter in the name.
// For example: "eexample", "exaample", "exammple", "examplee", "examplle"
func GraphemeRepetition(name string) (names []string) {
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

// DoubleGraphemeAdjacentReplacement
// Double grapheme adjacent replacement typos are created by replacing identical,
// consecutive letters in the name with adjacent keys on the keyboard.
// For example, www.gppgle.com and www.giigle.com.
func DoubleGraphemeAdjacentReplacement(name string, keyboard ...string) (names []string) {
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

// // stripDashesFunc typos are created by omitting a dot from the domain.
// // For example, www.a-b-c.com becomes www.abc.com
// func stripDashesFunc(tc Result) (results []Result) {
// 	for _, str := range replaceCharFunc(tc.Original.Domain, "-", "") {
// 		if tc.Original.Domain != str {
// 			dm := Domain{tc.Original.Subdomain, str, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 		}
// 	}
// 	return
// }

// GraphemeOmission
// Grapheme omission leaves out one character from the name.
// For google: "gogle", "gogle", "googe", "googl", "goole", "oogle",
func GraphemeOmission(name string) (names []string) {
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
	pluralize := NewClient()
	if pluralize.IsPlural(name) {
		names = append(names, pluralize.Singular(name))
	}
	if pluralize.IsSingular(name) {
		names = append(names, pluralize.Plural(name))
	}

	return
}

// // CcommonMisspellingsFunc are created with common misspellings in the given
// // language. For example, www.youtube.com becomes www.youtub.com and
// // www.abseil.com becomes www.absail.com
// func commonMisspellingsFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for _, word := range keyboard.Language.SimilarSpellings(tc.Original.Domain) {
// 			dm := Domain{tc.Original.Subdomain, word, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})

// 		}
// 	}
// 	return
// }

// // vowelSwappingFunc swaps vowels within the domain name except for the first letter.
// // For example, www.google.com becomes www.gaagle.com.
// func vowelSwappingFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for _, vchar := range keyboard.Language.Vowels {
// 			if strings.Contains(tc.Original.Domain, vchar) {
// 				for _, vvchar := range keyboard.Language.Vowels {
// 					new := strings.Replace(tc.Original.Domain, vchar, vvchar, -1)
// 					if new != tc.Original.Domain {
// 						dm := Domain{tc.Original.Subdomain, new, tc.Original.Suffix, Meta{}, false}
// 						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return
// }

// // homophonesFunc are created from sets of words that sound the same when spoken.
// // For example, www.base.com becomes www .bass.com.
// func homophonesFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for _, word := range keyboard.Language.SimilarSounds(tc.Original.Domain) {
// 			dm := Domain{tc.Original.Subdomain, word, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 		}
// 	}
// 	return
// }

// // homoglyphFunc when one or more characters that look similar to another
// // character but are different are called homogylphs. An example is that the
// // lower case l looks similar to the numeral one, e.g. l vs 1. For example,
// // google.com becomes goog1e.com.
// func homoglyphFunc(tc Result) (results []Result) {
// 	for i, char := range tc.Original.Domain {
// 		// Check the alphabet of the language associated with the keyboard for
// 		// homoglyphs
// 		for _, keyboard := range tc.Keyboards {
// 			for _, kchar := range keyboard.Language.SimilarChars(string(char)) {
// 				domain := fmt.Sprint(tc.Original.Domain[:i], kchar, tc.Original.Domain[i+1:])
// 				if tc.Original.Domain != domain {
// 					dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}
// 		// Check languages given with the (-l --language) CLI options for homoglyphs.
// 		for _, language := range tc.Languages {
// 			for _, lchar := range language.SimilarChars(string(char)) {
// 				domain := fmt.Sprint(tc.Original.Domain[:i], lchar, tc.Original.Domain[i+1:])
// 				if tc.Original.Domain != domain {
// 					dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}

// 	}
// 	return results
// }

// // wrongTopLevelDomain for example, www.google.co.nz becomes www.google.co.ns
// // and www.google.com becomes www.google.org. uses the 19 most common top level
// // domains.
// func wrongTopLevelDomainFunc(tc Result) (results []Result) {
// 	labels := strings.Split(tc.Original.Suffix, ".")
// 	length := len(labels)
// 	for _, suffix := range datasets.TLD {
// 		suffixLen := len(strings.Split(suffix, "."))
// 		if length == suffixLen && length == 1 {
// 			if suffix != tc.Original.Suffix {
// 				dm := Domain{tc.Original.Subdomain, tc.Original.Domain, suffix, Meta{}, false}
// 				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 			}
// 		}
// 	}
// 	return
// }

// // wrongSecondLevelDomain uses an alternate, valid second level domain for the
// // top level domain. For example, www.trademe.co.nz becomes www.trademe.ac.nz
// // and www.trademe.iwi.nz
// func wrongSecondLevelDomainFunc(tc Result) (results []Result) {
// 	labels := strings.Split(tc.Original.Suffix, ".")
// 	length := len(labels)
// 	//fmt.Println(length, labels)
// 	for _, suffix := range datasets.TLD {
// 		suffixLbl := strings.Split(suffix, ".")
// 		suffixLen := len(suffixLbl)
// 		if length == suffixLen && length == 2 {
// 			if suffixLbl[1] == labels[1] {
// 				if suffix != tc.Original.Suffix {
// 					dm := Domain{tc.Original.Subdomain, tc.Original.Domain, suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}
// 	}
// 	return
// }

// // wrongThirdLevelDomainFunc uses an alternate, valid third level domain.
// func wrongThirdLevelDomainFunc(tc Result) (results []Result) {
// 	labels := strings.Split(tc.Original.Suffix, ".")
// 	length := len(labels)
// 	for _, suffix := range datasets.TLD {
// 		suffixLbl := strings.Split(suffix, ".")
// 		suffixLen := len(suffixLbl)
// 		if length == suffixLen && length == 3 {
// 			if suffixLbl[1] == labels[1] && suffixLbl[2] == labels[2] {
// 				if suffix != tc.Original.Suffix {
// 					dm := Domain{tc.Original.Subdomain, tc.Original.Domain, suffix, Meta{}, false}
// 					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 				}
// 			}
// 		}
// 	}
// 	return
// }

// // bitsquattingFunc relies on random bit- errors to redirect connections
// // intended for popular domains
// func bitsquattingFunc(tc Result) (results []Result) {
// 	// TOOO: need to improve.
// 	masks := []int{1, 2, 4, 8, 16, 32, 64, 128}
// 	charset := make(map[string][]string)
// 	for _, board := range tc.Keyboards {
// 		for _, alpha := range board.Language.Graphemes {
// 			for _, mask := range masks {
// 				new := int([]rune(alpha)[0]) ^ mask
// 				for _, a := range board.Language.Graphemes {
// 					if string(a) == string(new) {
// 						charset[string(alpha)] = append(charset[string(alpha)], string(new))
// 					}
// 				}
// 			}
// 		}
// 	}

// 	for d, dchar := range tc.Original.Domain {
// 		for _, char := range charset[string(dchar)] {

// 			dnew := tc.Original.Domain[:d] + string(char) + tc.Original.Domain[d+1:]
// 			dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
// 			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 		}
// 	}
// 	return
// }

// // numeralSwapFunc are created by swapping numbers and corresponding words
// func numeralSwapFunc(tc Result) (results []Result) {
// 	for _, keyboard := range tc.Keyboards {
// 		for inum, words := range keyboard.Language.Numerals {
// 			for _, snum := range words {
// 				{
// 					dnew := strings.Replace(tc.Original.Domain, snum, inum, -1)
// 					dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
// 					if dnew != tc.Original.Domain {
// 						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 					}
// 				}
// 				{
// 					dnew := strings.Replace(tc.Original.Domain, inum, snum, -1)
// 					dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
// 					if dnew != tc.Original.Domain {
// 						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return
// }

// // missingCharFunc removes a character one at a time from the string.
// // For example, wwwgoogle.com and www.googlecom
// func missingCharFunc(str, character string) (results []string) {
// 	for i, char := range str {
// 		if character == string(char) {
// 			results = append(results, str[:i]+str[i+1:])
// 		}
// 	}
// 	return
// }

// // replaceCharFunc omits a character from the entire string.
// // For example, www.a-b-c.com becomes www.abc.com
// func replaceCharFunc(str, old, new string) (results []string) {
// 	results = append(results, strings.Replace(str, old, new, -1))
// 	return
// }
