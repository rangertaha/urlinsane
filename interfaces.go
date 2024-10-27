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

import "time"

const (
	ENTITY = "ENTITY"
	DOMAIN = "DOMAIN"
)

type Config interface {
	Target() string
	Keyboards() []Keyboard
	Languages() []Language
	Algorithms() []Algorithm
	Information() []Information
	Output() Output
	Concurrency() int
	Delay() time.Duration
	Random() time.Duration
	Verbose() bool
	Format() string
	File() string
	Count(...int64) int64
}

type Algorithm interface {
	Id() string
	Name() string
	IsType(string) bool
	Description() string
	// Fields() []string
	// Headers() []string
	Exec(Typo) []Typo
}

type Information interface {
	Id() string
	Name() string
	IsType(string) bool
	Description() string
	// Fields() []string
	Headers() []string
	Exec(Typo) Typo
}

type Output interface {
	Id() string
	Init(Config)
	Description() string
	Write(Typo) // Write(interface{})
	Save()
}

type Target interface {
	Repr() string
	Live(...bool) bool
	Meta() map[string]interface{}
	Add(string, interface{})
}

type Typo interface {
	Id(...int64) string
	Keyboard() Keyboard
	Language() Language
	Algorithm() Algorithm
	Original() Domain
	Variant() Domain
	Active(...bool) bool
	New(string) Typo
	Repr() string
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
	Repr() string
	Live() bool
	Meta() map[string]interface{}
	Add(string, interface{})
}

type Language interface {
	Id() string
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
