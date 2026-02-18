// Copyright 2026 Rangertaha. All Rights Reserved.
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
package latin

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
)

const LANGUAGE string = "la"

type Latin struct {
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

func (l *Latin) Id() string          { return l.code }
func (l *Latin) Name() string        { return l.name }
func (l *Latin) Description() string { return l.description }
func (l *Latin) Numerals() map[string][]string {
	return l.numerals
}
func (l *Latin) Cardinal() map[string]string { return languages.NumeralMap(l.numerals, 0) }
func (l *Latin) Ordinal() map[string]string  { return languages.NumeralMap(l.numerals, 1) }
func (l *Latin) Graphemes() []string         { return l.graphemes }
func (l *Latin) Vowels() []string            { return l.vowels }
func (l *Latin) Misspellings() [][]string    { return l.misspellings }
func (l *Latin) Homophones() [][]string      { return l.homophones }
func (l *Latin) Antonyms() map[string][]string {
	return l.antonyms
}
func (l *Latin) Homoglyphs() map[string][]string { return l.homoglyphs }
func (l *Latin) SimilarChars(char string) []string {
	return languages.SimilarChars(l.homoglyphs, char)
}
func (l *Latin) SimilarSpellings(word string) []string {
	return languages.SimilarSpellings(l.misspellings, word)
}
func (l *Latin) SimilarSounds(word string) []string {
	return languages.SimilarSounds(l.homophones, word)
}

func (l *Latin) Keyboards() (boards []internal.Keyboard) {
	for _, b := range languages.Keyboards() {
		if b.Language() == l.code {
			boards = append(boards, b)
		}
	}
	return
}

var (
	laMisspellings = [][]string{
		// Common classical/medieval orthography variants
		{"v", "u"},
		{"i", "j"},
		{"ae", "e"},
		{"oe", "e"},
	}

	laHomophones = [][]string{
		{"punctum", "."},
		{"linea", "-"},
		{"at", "@"},
		{"et", "&"},
	}

	laAntonyms = map[string][]string{
		"bonus": {"malus"},
		"malus": {"bonus"},
		"verus": {"falsus"},
		"falsus": {"verus"},
		"magnus": {"parvus"},
		"parvus": {"magnus"},
		"novus": {"vetus"},
		"vetus": {"novus"},
		"tutus": {"periculosus"},
		"periculosus": {"tutus"},
		"ita": {"non"},
		"non": {"ita"},
	}

	Language = Latin{
		code:        LANGUAGE,
		name:        "Latin",
		description: "Latin is a classical language of the Roman Empire and the basis of the Romance languages.",
		numerals: map[string][]string{
			"0":  {"nulla", "nullus"},
			"1":  {"unus", "primus"},
			"2":  {"duo", "secundus"},
			"3":  {"tres", "tertius"},
			"4":  {"quattuor", "quartus"},
			"5":  {"quinque", "quintus"},
			"6":  {"sex", "sextus"},
			"7":  {"septem", "septimus"},
			"8":  {"octo", "octavus"},
			"9":  {"novem", "nonus"},
			"10": {"decem", "decimus"},
			"11": {"undecim"},
			"12": {"duodecim"},
			"20": {"viginti"},
			"30": {"triginta"},
			"40": {"quadraginta"},
			"50": {"quinquaginta"},
			"60": {"sexaginta"},
			"70": {"septuaginta"},
			"80": {"octoginta"},
			"90": {"nonaginta"},
			"100": {"centum"},
			"1000": {"mille"},
		},
		graphemes: []string{
			"a", "b", "c", "d", "e", "f", "g",
			"h", "i", "j", "k", "l", "m", "n",
			"o", "p", "q", "r", "s", "t", "u",
			"v", "w", "x", "y", "z",
		},
		vowels:       []string{"a", "e", "i", "o", "u", "y"},
		misspellings: laMisspellings,
		homophones:   laHomophones,
		antonyms:     laAntonyms,
		homoglyphs:   languages.DefaultLatinHomoglyphs(),
	}
)

func init() {
	languages.AddLanguage(LANGUAGE, func() internal.Language {
		return &Language
	})
}

