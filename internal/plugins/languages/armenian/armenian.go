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
package armenian

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

// https://en.wikipedia.org/wiki/Armenian_alphabet

const LANGUAGE string = "hy"

type Armenian struct {
	code         string
	name         string
	description  string
	numerals     map[string][]string
	graphemes    []string
	vowels       []string
	misspellings [][]string
	homophones   [][]string
	antonyms     map[string][]string
	homoglyphs   map[string][]string
}

func (l *Armenian) Id() string {
	return l.code
}
func (l *Armenian) Name() string {
	return l.name
}
func (l *Armenian) Description() string {
	return l.description
}
func (l *Armenian) Numerals() map[string][]string {
	return l.numerals
}
func (l *Armenian) Cardinal() map[string]string {
	return languages.NumeralMap(l.numerals, 0)
}

func (l *Armenian) Ordinal() map[string]string {
	return languages.NumeralMap(l.numerals, 1)
}

func (l *Armenian) Graphemes() []string {
	return l.graphemes
}

func (l *Armenian) Vowels() []string {
	return l.vowels
}

func (l *Armenian) Misspellings() [][]string {
	return l.misspellings
}

func (l *Armenian) Homophones() [][]string {
	return l.homophones
}

func (l *Armenian) Antonyms() map[string][]string {
	return l.antonyms
}

func (l *Armenian) Homoglyphs() map[string][]string {
	return l.homoglyphs
}

func (l *Armenian) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}

func (l *Armenian) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}

func (l *Armenian) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *Armenian) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	// hyMisspellings are common misspellings
	hyMisspellings = [][]string{
		// Domain-friendly Armenian orthography variants
		{"եւ", "և"}, // different Unicode forms seen in the wild
		{"օ", "ո"},  // common confusion in some contexts
		{"տասնիննը", "տասնինը"},
	}

	// hyHomophones are words that sound alike
	hyHomophones = [][]string{
		{"կետ", "."},
		{"կետը", "."},
		{"շնիկ", "@"},
		{"գծիկ", "-"},
	}

	// hyAntonyms are words opposite in meaning to another (e.g. bad and good ).
	hyAntonyms = map[string][]string{
		"լավ":   {"վատ"},
		"մեծ":   {"փոքր"},
		"բարձր": {"ցածր"},
		"արագ":  {"դանդաղ"},
		"ուժեղ": {"թույլ"},
		"նոր":   {"հին"},
		"օր":    {"գիշեր"},
		"սկիզբ": {"վերջ"},
		"այո":   {"ոչ"},
		"լույս": {"մութ"},
	}

	hyLanguage = Armenian{
		// https://www.loc.gov/standards/iso639-2/php/code_list.php
		code:        LANGUAGE,
		name:        "Armenian",
		description: "Armenian is the native language of the Armenian people",

		// http://mylanguages.org/armenian_numbers.php
		numerals: map[string][]string{
			// Number: cardinal..,  ordinal.., other...
			"0":          {"զրո"},
			"1":          {"մեկ", "առաջին"},
			"2":          {"երկու", "երկրորդ"},
			"3":          {"երեք", "երրորդ"},
			"4":          {"չորս", "չորրորդ"},
			"5":          {"հինգ", "հինգերորդ"},
			"6":          {"վեց", "վեցերորդ"},
			"7":          {"յոթ", "յոթերորդ"},
			"8":          {"ութ", "ութերորդ"},
			"9":          {"ինը", "իններորդ"},
			"10":         {"տաս", "տասերորդ"},
			"11":         {"տասնմեկ"},
			"12":         {"տասներկու"},
			"13":         {"տասներեք"},
			"14":         {"տասնչորս"},
			"15":         {"տասնհինգ"},
			"16":         {"տասնվեց"},
			"17":         {"տասնյոթ"},
			"18":         {"տասնութ"},
			"19":         {"տասնինը"},
			"20":         {"քսան"},
			"30":         {"երեսուն"},
			"40":         {"քառասուն"},
			"50":         {"հիսուն"},
			"60":         {"վաթսուն"},
			"70":         {"յոթանասուն"},
			"80":         {"ութանասուն"},
			"90":         {"իննսուն"},
			"100":        {"հարյուր"},
			"1000":       {"հազար"},
			"1000000":    {"միլիոն"},
			"1000000000": {"միլիարդ"},
		},
		// http://mylanguages.org/armenian_alphabet.php
		graphemes: []string{
			"ա", "բ", "գ", "դ", "ե", "զ", "է", "ը",
			"թ", "ժ", "ի", "լ", "խ", "ծ", "կ", "հ",
			"ձ", "ղ", "ճ", "մ", "յ", "ն", "շ", "ո",
			"չ", "պ", "ջ", "ռ", "ս", "վ", "տ", "ր",
			"ց", "փ", "ք", "և", "օ", "ֆ",
		},
		vowels:       []string{},
		misspellings: hyMisspellings,
		homophones:   hyHomophones,
		antonyms:     hyAntonyms,
		homoglyphs: map[string][]string{
			// Focused Armenian confusables for IDN / mixed-script domains
			"օ": {"o", "0", "Ο", "ο", "О", "о", "Օ", "ӧ", "ö"},
			"ո": {"o", "0", "Ο", "ο", "О", "о", "Օ"},
			"ս": {"u", "υ", "ц"},
			"հ": {"h", "һ", "Ꮒ", "н"},
			"ք": {"q"},
			"ր": {"r"},
			"լ": {"l", "1", "I"},
			"ա": {"a", "ɑ", "а"},
			"ե": {"e", "е", "є"},
			"բ": {"b", "Ь"},
			"ճ": {"6"},
			"և": {"ev", "&"},
		},
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &hyLanguage
	})
}
