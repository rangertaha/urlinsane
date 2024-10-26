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
package urlinsane

// type Module interface {
// 	Code() string
// 	Name() string
// 	Description() string
// 	Fields() []string
// 	Headers() []string
// 	Exec(Typo) []Typo
// }

type Algorithm interface {
	Code() string
	Name() string
	Description() string
	Fields() []string
	Headers() []string
	Exec(Typo) []Typo
}

type Information interface {
	Code() string
	Name() string
	Description() string
	Fields() []string
	Headers() []string
	Exec(Typo) Typo
}

type Typo interface {
	Keyboard() Keyboard
	Language() Language
	Algorithm() Algorithm
	Original() Domain
	Variant() Domain
}

// type Result interface {
// 	Keyboards() []Keyboard
// 	Languages() []Language
// 	Original() Domain
// 	Variant() Domain
// 	Algo() Module
// 	Data() map[string]string
// }

type Domain interface {
	Subdomain() string
	Domain() string
	Suffix() string
	Live() bool
	Meta() map[string]interface{}
}

type Language interface {
	Code() string
	Name() string

	// Numerals in the broadest sense is a word or phrase that
	// describes a numerical quantity.
	Numerals() map[string][]string

	// Graphemes is the smallest functional unit of a writing system.
	Graphemes() []string

	// Vowels are syllabic speech sound pronounced without any stricture in the vocal tract.
	Vowels() []string

	Misspellings() [][]string

	Homophones() [][]string

	Antonyms() map[string][]string

	Homoglyphs() map[string][]string

	Keyboards() []Keyboard
}

type Keyboard interface {
	Id() string
	Title() string
	Summary() string
	Layouts() []string
	Language() string
}
