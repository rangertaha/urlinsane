// The MIT License (MIT)
//
// Copyright © 2019 Rangertaha <rangertaha@gmail.com>
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

package typo

import (
	"fmt"
	"strings"

	"github.com/cybersectech-org/urlinsane/pkg/datasets"
)

// Typos ...
var Typos = NewRegistry()

var missingDot = Module{
	Code:        "MD",
	Name:        "Missing Dot",
	Description: "Missing Dot is created by omitting a dot from the domain.",
	Exe:         missingDotFunc,
}
var subdomainInsertion = Module{
	Code:        "SI",
	Name:        "Subdomain Insertion",
	Description: "Inserting common subdomain at the beginning of the domain.",
	Exe:         subdomainInsertionFunc,
}
var missingDashes = Module{
	Code:        "MDS",
	Name:        "Missing Dashes",
	Description: "Missing Dashes is created by stripping all dashes from the domain.",
	Exe:         missingDashFunc,
}
var stripDashes = Module{
	Code:        "SD",
	Name:        "Strip Dashes",
	Description: "Strip Dashes is created by omitting a dash from the domain",
	Exe:         stripDashesFunc,
}
var characterOmission = Module{
	Code:        "CO",
	Name:        "Character Omission",
	Description: "Character Omission Omitting a character from the domain.",
	Exe:         characterOmissionFunc,
}
var characterSwap = Module{
	Code:        "CS",
	Name:        "Character Swap",
	Description: "Character Swap Swapping two consecutive characters in a domain",
	Exe:         characterSwapFunc,
}
var adjacentCharacterSubstitution = Module{
	Code:        "ACS",
	Name:        "Adjacent Character Substitution",
	Description: "Adjacent Character Substitution replaces adjacent characters",
	Exe:         adjacentCharacterSubstitutionFunc,
}
var adjacentCharacterInsertion = Module{
	Code:        "ACI",
	Name:        "Adjacent Character Insertion",
	Description: "Adjacent Character Insertion inserts adjacent character ",
	Exe:         adjacentCharacterInsertionFunc,
}
var homoglyphs = Module{
	Code:        "HG",
	Name:        "Homoglyphs",
	Description: "Homoglyphs replaces characters with characters that look similar",
	Exe:         homoglyphFunc,
}
var singularPluralise = Module{
	Code:        "SP",
	Name:        "Singular Pluralise",
	Description: "Singular Pluralise creates a singular domain plural and vice versa",
	Exe:         singularPluraliseFunc,
}

var characterRepeat = Module{
	Code:        "CR",
	Name:        "Character Repeat",
	Description: "Character Repeat Repeats a character of the domain name twice",
	Exe:         characterRepeatFunc,
}

var doubleCharacterReplacement = Module{
	Code:        "DCR",
	Name:        "Double Character Replacement",
	Description: "Double Character Replacement repeats a character twice.",
	Exe:         doubleCharacterReplacementFunc,
}
var commonMisspellings = Module{
	Code:        "CM",
	Name:        "Common Misspellings",
	Description: "Common Misspellings are created from a dictionary of commonly misspelled words",
	Exe:         commonMisspellingsFunc,
}
var homophones = Module{
	Code:        "HP",
	Name:        "Homophones",
	Description: "Homophones Modules are created from sets of words that sound the same",
	Exe:         homophonesFunc,
}

var vowelSwapping = Module{
	Code:        "VS",
	Name:        "Vowel Swapping",
	Description: "Vowel Swapping is created by swaps vowels",
	Exe:         vowelSwappingFunc,
}

var bitFlipping = Module{
	Code:        "BF",
	Name:        "Bit Flipping",
	Description: "Bitsquatting relies on random bit-errors to redirect connections",
	Exe:         bitsquattingFunc,
}

var wrongTopLevelDomain = Module{
	Code:        "WTLD",
	Name:        "Wrong TLD",
	Description: "Wrong Top Level Domain",
	Exe:         wrongTopLevelDomainFunc,
}

var wrongSecondLevelDomain = Module{
	Code:        "W2TLD",
	Name:        "Wrong 2nd TLD",
	Description: "Wrong Second Level Domain",
	Exe:         wrongSecondLevelDomainFunc,
}

var wrongThirdLevelDomain = Module{
	Code:        "W3TLD",
	Name:        "Wrong 3rd TLD",
	Description: "Wrong Third Level Domain",
	Exe:         wrongThirdLevelDomainFunc,
}

var numeralSwap = Module{
	Code:        "NS",
	Name:        "Numeral Swap",
	Description: "Numeral Swap numbers, words and vice versa",
	Exe:         numeralSwapFunc,
}

var periodInsertion = Module{
	Code:        "PI",
	Name:        "Period Insertion",
	Description: "Inserting periods in the target domain",
	Exe:         periodInsertionFunc,
}

var hyphenInsertion = Module{
	Code:        "HI",
	Name:        "Dash Insertion",
	Description: "Inserting hyphens in the target domain",
	Exe:         hyphenInsertionFunc,
}

var alphabetInsertion = Module{
	Code:        "AI",
	Name:        "Alphabet Insertion",
	Description: "Inserting the language specific alphabet in the target domain",
	Exe:         alphabetInsertionFunc,
}

var alphabetReplacement = Module{
	Code:        "AR",
	Name:        "Alphabet Replacement",
	Description: "Replacing the language specific alphabet in the target domain",
	Exe:         alphabetReplacementnFunc,
}

func init() {
	Typos.Set("MD", missingDot)
	Typos.Set("SI", subdomainInsertion)

	Typos.Set("MD", missingDot)
	Typos.Set("SI", subdomainInsertion)
	Typos.Set("MDS", missingDashes)
	Typos.Set("CO", characterOmission)
	Typos.Set("CS", characterSwap)
	Typos.Set("ACS", adjacentCharacterSubstitution)
	Typos.Set("ACI", adjacentCharacterInsertion)
	Typos.Set("CR", characterRepeat)
	Typos.Set("DCR", doubleCharacterReplacement)
	Typos.Set("SD", stripDashes)
	Typos.Set("SP", singularPluralise)
	Typos.Set("CM", commonMisspellings)
	Typos.Set("VS", vowelSwapping)
	Typos.Set("HG", homoglyphs)
	Typos.Set("WTLD", wrongTopLevelDomain)
	Typos.Set("W2TLD", wrongSecondLevelDomain)
	Typos.Set("W3TLD", wrongThirdLevelDomain)
	Typos.Set("HP", homophones)
	Typos.Set("BF", bitFlipping)
	Typos.Set("NS", numeralSwap)
	Typos.Set("PI", periodInsertion)
	Typos.Set("HI", hyphenInsertion)
	Typos.Set("AI", alphabetInsertion)
	Typos.Set("AR", alphabetReplacement)

	Typos.Set("ALL",
		missingDot,
		subdomainInsertion,
		missingDashes,
		characterOmission,
		characterSwap,
		adjacentCharacterSubstitution,
		adjacentCharacterInsertion,
		characterRepeat,
		doubleCharacterReplacement,
		stripDashes,
		singularPluralise,
		commonMisspellings,
		vowelSwapping,
		homoglyphs,
		wrongTopLevelDomain,
		wrongSecondLevelDomain,
		wrongThirdLevelDomain,
		homophones,
		bitFlipping,
		numeralSwap,
		periodInsertion,
		hyphenInsertion,
		alphabetInsertion,
		alphabetReplacement,
	)

}

// missingDotFunc typos are created by omitting a dot from the domain. For example, wwwgoogle.com and www.googlecom
func missingDotFunc(tc Result) (results []Result) {
	for _, str := range missingCharFunc(tc.Original.String(), ".") {
		if tc.Original.Domain != str {
			dm := Domain{tc.Original.Subdomain, str, tc.Original.Suffix, Meta{}, false}
			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
		}
	}
	dm := Domain{tc.Original.Subdomain, strings.Replace(tc.Original.Domain, ".", "", -1), tc.Original.Suffix, Meta{}, false}
	results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
	return results
}

// subdomainInsertionFunc typos are created by inserting common subdomains at the begining of the domain. wwwgoogle.com and ftpgoogle.com
func subdomainInsertionFunc(tc Result) (results []Result) {
	for _, str := range datasets.SUBDOMAINS {
		dm := Domain{tc.Original.Subdomain, str + tc.Original.Domain, tc.Original.Suffix, Meta{}, false}
		results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
	}
	return results
}

// missingDashFunc typos are created by omitting a dash from the domain.
// For example, www.a-b-c.com becomes www.ab-c.com, www.a-bc.com, and ww.abc.com
func missingDashFunc(tc Result) (results []Result) {
	for _, str := range missingCharFunc(tc.Original.Domain, "-") {
		if tc.Original.Domain != str {
			dm := Domain{tc.Original.Subdomain, str, tc.Original.Suffix, Meta{}, false}
			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
		}
	}
	dm := Domain{tc.Original.Subdomain, strings.Replace(tc.Original.Domain, "-", "", -1), tc.Original.Suffix, Meta{}, false}
	results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
	return results
}

// characterOmissionFunc typos are when one character in the original domain name is omitted.
// For example: www.exmple.com
func characterOmissionFunc(tc Result) (results []Result) {
	for i := range tc.Original.Domain {
		if i <= len(tc.Original.Domain)-1 {
			domain := fmt.Sprint(
				tc.Original.Domain[:i],
				tc.Original.Domain[i+1:],
			)
			if tc.Original.Domain != domain {
				dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})

			}
		}
	}
	return results
}

// characterSwapFunc typos are when two consecutive characters are swapped in the original domain name.
// Example: www.examlpe.com
func characterSwapFunc(tc Result) (results []Result) {
	for i := range tc.Original.Domain {
		if i <= len(tc.Original.Domain)-2 {
			domain := fmt.Sprint(
				tc.Original.Domain[:i],
				string(tc.Original.Domain[i+1]),
				string(tc.Original.Domain[i]),
				tc.Original.Domain[i+2:],
			)
			if tc.Original.Domain != domain {
				dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
			}
		}
	}
	return results
}

// adjacentCharacterSubstitutionFunc typos are when characters are replaced in the original domain name by their
// adjacent ones on a specific keyboard layout, e.g., www.ezample.com, where “x” was replaced by the adjacent
// character “z” in a the QWERTY keyboard layout.
func adjacentCharacterSubstitutionFunc(tc Result) (results []Result) {
	for _, keyboard := range tc.Keyboards {
		for i, char := range tc.Original.Domain {
			for _, key := range keyboard.Adjacent(string(char)) {
				domain := tc.Original.Domain[:i] + string(key) + tc.Original.Domain[i+1:]
				dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
			}
		}
	}
	return
}

// adjacentCharacterInsertionFunc are created by inserting letters adjacent of each letter. For example, www.googhle.com
// and www.goopgle.com
func adjacentCharacterInsertionFunc(tc Result) (results []Result) {
	for _, keyboard := range tc.Keyboards {
		for i, char := range tc.Original.Domain {
			for _, key := range keyboard.Adjacent(string(char)) {
				d1 := tc.Original.Domain[:i] + string(key) + string(char) + tc.Original.Domain[i+1:]
				dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
				results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})

				d2 := tc.Original.Domain[:i] + string(char) + string(key) + tc.Original.Domain[i+1:]
				dm2 := Domain{tc.Original.Subdomain, d2, tc.Original.Suffix, Meta{}, false}
				results = append(results, Result{Original: tc.Original, Variant: dm2, Typo: tc.Typo, Data: tc.Data})
			}
		}
	}
	return
}

// adjacentCharacterInsertionFunc are created by inserting letters adjacent of each letter. For example, www.googhle.com
// and www.goopgle.com
func hyphenInsertionFunc(tc Result) (results []Result) {

	for i, char := range tc.Original.Domain {
		d1 := tc.Original.Domain[:i] + "-" + string(char) + tc.Original.Domain[i+1:]
		if i == len(tc.Original.Domain)-1 {
			d1 = tc.Original.Domain[:i] + string(char) + "-" + tc.Original.Domain[i+1:]
		}
		dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
		results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
	}
	return
}

func alphabetInsertionFunc(tc Result) (results []Result) {
	alphabet := map[string]bool{}
	for _, keyboard := range tc.Keyboards {
		for _, a := range keyboard.Language.Graphemes {
			alphabet[a] = true
		}
	}
	for i, char := range tc.Original.Domain {
		for alp := range alphabet {
			d1 := tc.Original.Domain[:i] + alp + string(char) + tc.Original.Domain[i+1:]
			if i == len(tc.Original.Domain)-1 {
				d1 = tc.Original.Domain[:i] + string(char) + alp + tc.Original.Domain[i+1:]
			}
			dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
			results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
		}
	}
	return
}

func alphabetReplacementnFunc(tc Result) (results []Result) {
	alphabet := map[string]bool{}
	for _, keyboard := range tc.Keyboards {
		for _, a := range keyboard.Language.Graphemes {
			alphabet[a] = true
		}
	}

	for i := range tc.Original.Domain {
		for alp, _ := range alphabet {

			d1 := tc.Original.Domain[:i] + alp + tc.Original.Domain[i+1:]
			dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
			results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})

			if i == len(tc.Original.Domain)-1 {
				d1 = tc.Original.Domain[:i] + alp + tc.Original.Domain[i+1:]
				dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
				results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
			}
		}
	}
	return
}

// adjacentCharacterInsertionFunc are created by inserting letters adjacent of each letter. For example, www.googhle.com
// and www.goopgle.com
func periodInsertionFunc(tc Result) (results []Result) {

	for i, char := range tc.Original.Domain {

		d1 := tc.Original.Domain[:i] + "." + string(char) + tc.Original.Domain[i+1:]
		dm1 := Domain{tc.Original.Subdomain, d1, tc.Original.Suffix, Meta{}, false}
		results = append(results, Result{Original: tc.Original, Variant: dm1, Typo: tc.Typo, Data: tc.Data})
	}

	return
}

// characterRepeatFunc are created by repeating a letter of the domain name.
// Example, www.ggoogle.com and www.gooogle.com
func characterRepeatFunc(tc Result) (results []Result) {
	for i := range tc.Original.Domain {
		if i <= len(tc.Original.Domain) {
			domain := fmt.Sprint(
				tc.Original.Domain[:i],
				string(tc.Original.Domain[i]),
				string(tc.Original.Domain[i]),
				tc.Original.Domain[i+1:],
			)
			if tc.Original.Domain != domain {
				dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
			}
		}
	}
	return results
}

// doubleCharacterReplacementFunc are created by replacing identical, consecutive
// letters of the domain name with adjacent letters on the keyboard.
// For example, www.gppgle.com and www.giigle.com
func doubleCharacterReplacementFunc(tc Result) (results []Result) {
	for _, keyboard := range tc.Keyboards {
		for i, char := range tc.Original.Domain {
			if i < len(tc.Original.Domain)-1 {
				if tc.Original.Domain[i] == tc.Original.Domain[i+1] {
					for _, key := range keyboard.Adjacent(string(char)) {
						domain := tc.Original.Domain[:i] + string(key) + string(key) + tc.Original.Domain[i+2:]
						dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
					}
				}
			}
		}
	}
	return
}

// stripDashesFunc typos are created by omitting a dot from the domain.
// For example, www.a-b-c.com becomes www.abc.com
func stripDashesFunc(tc Result) (results []Result) {
	for _, str := range replaceCharFunc(tc.Original.Domain, "-", "") {
		if tc.Original.Domain != str {
			dm := Domain{tc.Original.Subdomain, str, tc.Original.Suffix, Meta{}, false}
			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
		}
	}
	return
}

// singularPluraliseFunc are created by making a singular domain plural and
// vice versa. For example, www.google.com becomes www.googles.com and
// www.games.co.nz becomes www.game.co.nz
func singularPluraliseFunc(tc Result) (results []Result) {
	for _, pchar := range []string{"s", "ing"} {
		var domain string
		if strings.HasSuffix(tc.Original.Domain, pchar) {
			domain = strings.TrimSuffix(tc.Original.Domain, pchar)
		} else {
			domain = tc.Original.Domain + pchar
		}
		dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
		results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
	}
	return
}

// CcommonMisspellingsFunc are created with common misspellings in the given
// language. For example, www.youtube.com becomes www.youtub.com and
// www.abseil.com becomes www.absail.com
func commonMisspellingsFunc(tc Result) (results []Result) {
	for _, keyboard := range tc.Keyboards {
		for _, word := range keyboard.Language.SimilarSpellings(tc.Original.Domain) {
			dm := Domain{tc.Original.Subdomain, word, tc.Original.Suffix, Meta{}, false}
			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})

		}
	}
	return
}

// vowelSwappingFunc swaps vowels within the domain name except for the first letter.
// For example, www.google.com becomes www.gaagle.com.
func vowelSwappingFunc(tc Result) (results []Result) {
	for _, keyboard := range tc.Keyboards {
		for _, vchar := range keyboard.Language.Vowels {
			if strings.Contains(tc.Original.Domain, vchar) {
				for _, vvchar := range keyboard.Language.Vowels {
					new := strings.Replace(tc.Original.Domain, vchar, vvchar, -1)
					if new != tc.Original.Domain {
						dm := Domain{tc.Original.Subdomain, new, tc.Original.Suffix, Meta{}, false}
						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
					}
				}
			}
		}
	}
	return
}

// homophonesFunc are created from sets of words that sound the same when spoken.
// For example, www.base.com becomes www .bass.com.
func homophonesFunc(tc Result) (results []Result) {
	for _, keyboard := range tc.Keyboards {
		for _, word := range keyboard.Language.SimilarSounds(tc.Original.Domain) {
			dm := Domain{tc.Original.Subdomain, word, tc.Original.Suffix, Meta{}, false}
			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
		}
	}
	return
}

// homoglyphFunc when one or more characters that look similar to another
// character but are different are called homogylphs. An example is that the
// lower case l looks similar to the numeral one, e.g. l vs 1. For example,
// google.com becomes goog1e.com.
func homoglyphFunc(tc Result) (results []Result) {
	for i, char := range tc.Original.Domain {
		// Check the alphabet of the language associated with the keyboard for
		// homoglyphs
		for _, keyboard := range tc.Keyboards {
			for _, kchar := range keyboard.Language.SimilarChars(string(char)) {
				domain := fmt.Sprint(tc.Original.Domain[:i], kchar, tc.Original.Domain[i+1:])
				if tc.Original.Domain != domain {
					dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
				}
			}
		}
		// Check languages given with the (-l --language) CLI options for homoglyphs.
		for _, language := range tc.Languages {
			for _, lchar := range language.SimilarChars(string(char)) {
				domain := fmt.Sprint(tc.Original.Domain[:i], lchar, tc.Original.Domain[i+1:])
				if tc.Original.Domain != domain {
					dm := Domain{tc.Original.Subdomain, domain, tc.Original.Suffix, Meta{}, false}
					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
				}
			}
		}

	}
	return results
}

// wrongTopLevelDomain for example, www.google.co.nz becomes www.google.co.ns
// and www.google.com becomes www.google.org. uses the 19 most common top level
// domains.
func wrongTopLevelDomainFunc(tc Result) (results []Result) {
	labels := strings.Split(tc.Original.Suffix, ".")
	length := len(labels)
	for _, suffix := range datasets.TLD {
		suffixLen := len(strings.Split(suffix, "."))
		if length == suffixLen && length == 1 {
			if suffix != tc.Original.Suffix {
				dm := Domain{tc.Original.Subdomain, tc.Original.Domain, suffix, Meta{}, false}
				results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
			}
		}
	}
	return
}

// wrongSecondLevelDomain uses an alternate, valid second level domain for the
// top level domain. For example, www.trademe.co.nz becomes www.trademe.ac.nz
// and www.trademe.iwi.nz
func wrongSecondLevelDomainFunc(tc Result) (results []Result) {
	labels := strings.Split(tc.Original.Suffix, ".")
	length := len(labels)
	//fmt.Println(length, labels)
	for _, suffix := range datasets.TLD {
		suffixLbl := strings.Split(suffix, ".")
		suffixLen := len(suffixLbl)
		if length == suffixLen && length == 2 {
			if suffixLbl[1] == labels[1] {
				if suffix != tc.Original.Suffix {
					dm := Domain{tc.Original.Subdomain, tc.Original.Domain, suffix, Meta{}, false}
					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
				}
			}
		}
	}
	return
}

// wrongThirdLevelDomainFunc uses an alternate, valid third level domain.
func wrongThirdLevelDomainFunc(tc Result) (results []Result) {
	labels := strings.Split(tc.Original.Suffix, ".")
	length := len(labels)
	for _, suffix := range datasets.TLD {
		suffixLbl := strings.Split(suffix, ".")
		suffixLen := len(suffixLbl)
		if length == suffixLen && length == 3 {
			if suffixLbl[1] == labels[1] && suffixLbl[2] == labels[2] {
				if suffix != tc.Original.Suffix {
					dm := Domain{tc.Original.Subdomain, tc.Original.Domain, suffix, Meta{}, false}
					results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
				}
			}
		}
	}
	return
}

// bitsquattingFunc relies on random bit- errors to redirect connections
// intended for popular domains
func bitsquattingFunc(tc Result) (results []Result) {
	// TOOO: need to improve.
	masks := []int{1, 2, 4, 8, 16, 32, 64, 128}
	charset := make(map[string][]string)
	for _, board := range tc.Keyboards {
		for _, alpha := range board.Language.Graphemes {
			for _, mask := range masks {
				new := int([]rune(alpha)[0]) ^ mask
				for _, a := range board.Language.Graphemes {
					if string(a) == string(new) {
						charset[string(alpha)] = append(charset[string(alpha)], string(new))
					}
				}
			}
		}
	}

	for d, dchar := range tc.Original.Domain {
		for _, char := range charset[string(dchar)] {

			dnew := tc.Original.Domain[:d] + string(char) + tc.Original.Domain[d+1:]
			dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
			results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
		}
	}
	return
}

// numeralSwapFunc are created by swapping numbers and corresponding words
func numeralSwapFunc(tc Result) (results []Result) {
	for _, keyboard := range tc.Keyboards {
		for inum, words := range keyboard.Language.Numerals {
			for _, snum := range words {
				{
					dnew := strings.Replace(tc.Original.Domain, snum, inum, -1)
					dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
					if dnew != tc.Original.Domain {
						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
					}
				}
				{
					dnew := strings.Replace(tc.Original.Domain, inum, snum, -1)
					dm := Domain{tc.Original.Subdomain, dnew, tc.Original.Suffix, Meta{}, false}
					if dnew != tc.Original.Domain {
						results = append(results, Result{Original: tc.Original, Variant: dm, Typo: tc.Typo, Data: tc.Data})
					}
				}
			}
		}
	}
	return
}

// missingCharFunc removes a character one at a time from the string.
// For example, wwwgoogle.com and www.googlecom
func missingCharFunc(str, character string) (results []string) {
	for i, char := range str {
		if character == string(char) {
			results = append(results, str[:i]+str[i+1:])
		}
	}
	return
}

// replaceCharFunc omits a character from the entire string.
// For example, www.a-b-c.com becomes www.abc.com
func replaceCharFunc(str, old, new string) (results []string) {
	results = append(results, strings.Replace(str, old, new, -1))
	return
}
