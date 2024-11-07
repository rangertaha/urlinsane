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
package internal

import (
	"time"
)

const (
	PACKAGE = iota
	DOMAIN
	EMAIL
	NAME
)

type Initializer interface {
	Init(Config)
}

type Closer interface {
	Close()
}

type Config interface {
	Target() Target
	Keyboards() []Keyboard
	Languages() []Language
	Algorithms() []Algorithm
	Information() []Information
	Output() Output
	Concurrency() int
	DnsServers() []string
	Delay() time.Duration
	Random() time.Duration
	Verbose() bool
	Progress() bool
<<<<<<< HEAD
=======
	Screenshot() bool
>>>>>>> develop
	Format() string
	File() string
	Dist() int
	Type() int
	ShowAll() bool
	Filters() []string
}

type Algorithm interface {
	Id() string
	Name() string
	Description() string
}

type DomainAlgorithm interface {
	Domain(Typo) []Typo
}
type PackageAlgorithm interface {
	Package(Typo) []Typo
}
type EmailAlgorithm interface {
	Email(Typo) []Typo
}
type UserAlgorithm interface {
	Username(Typo) []Typo
}
type ExecAlgorithm interface {
	Exec(Typo) []Typo
}

type Information interface {
	Id() string
	Description() string
	Headers() []string
	Exec(Typo) Typo
	Order() int
}

type Storage interface {
	Init(Config)
	Read(key string) (error, interface{})
	Write(key string, value interface{}) error
}

type Output interface {
	Id() string
	Description() string
	Write(Typo)
	Summary(map[string]int64)
	Save()
}

type Target interface {
	Meta() map[string]interface{}
	Add(string, interface{})
	Get(string) interface{}
	Ready(...bool) bool
	Live(...bool) bool
	Name() string
	Domain() (string, string, string)
	Email() (string, string)
	Json() ([]byte, error)
}

type Typo interface {
	Algorithm() Algorithm
	Original() Target
	Variant() Target
	Active() bool
	Clone(string) Typo
	String() string
	Ld() int
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
	GetMeta() map[string]interface{}
	AddMeta(string, interface{})

	GetUsername() string
	SetUsername(string)

	GetSubdomain() string
	SetSubdomain(string)

	GetDomain() string
	SetDomain(string)

	GetSuffix() string
	SetSuffix(string)

	GetUrl() string
	SetUrl(string)

	Live() bool
	Name() string
	String() string
}

type Language interface {
	Id() string
	Name() string
	Description() string

	// Numerals in the broadest sense a word or phrase that
	// describes a numerical quantity. Example: one, first
	Numerals() map[string][]string

	// Cardinal numbers are the words of numbers that are used for counting
	// Example: one, two, three, four, five, six, seven, eight, nine, ten
	// See: https://byjus.com/maths/cardinal-numbers/
	Cardinal() map[string]string

	// They are used to denote the rank position or order of something
	// Example: Examples: 1st, 2nd, 5th, 6th, 9th or first, second, third
	// See: https://byjus.com/maths/cardinal-numbers/
	Ordinal() map[string]string

	// Graphemes is the smallest functional unit of a writing system.
	Graphemes() []string

	// Vowels are syllabic speech sounds pronounced without any stricture in the vocal tract.
	Vowels() []string

	Misspellings() [][]string

	Homophones() [][]string

	Antonyms() map[string][]string

	Homoglyphs() map[string][]string

	SimilarChars(char string) []string

	SimilarSpellings(word string) []string

	SimilarSounds(word string) []string

	Keyboards() []Keyboard
}

type Keyboard interface {
	Id() string
	Name() string
	Description() string
	Layouts() []string
	Adjacent(string) []string
	Language() string
}
