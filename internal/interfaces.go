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
	"encoding/json"
	"time"
)

type Initializer interface {
	Init(Config)
}

type Closer interface {
	Close()
}

type Config interface {
	Target() string

	// Plugins
	Keyboards() []Keyboard
	Languages() []Language
	Algorithms() []Algorithm
	Collectors() []Collector
	Analyzers() []Analyzer
	Database() Database
	Output() Output

	// Performance
	Workers() int
	Delay() time.Duration
	Random() time.Duration
	Timeout() time.Duration
	TTL() time.Duration

	// Processing
	Distance() int

	// Output
	Verbose() bool
	Progress() bool
	Banner() bool
	Format() string
	Filters() []string
	Registered() bool
	Unregistered() bool
	Summary() bool

	Dir() string
	File() string
}

type Algorithm interface {
	Id() string
	Name() string
	Description() string
	Exec(origin Domain, acc Accumulator) error
}

type Analyzer interface {
	Id() string
	Order() int
	Description() string
	Headers() []string
	Exec(origin Domain, variant Domain, acc Accumulator) error
}

type Collector interface {
	Id() string
	Order() int
	Description() string
	Headers() []string
	Exec(Accumulator) error
}

type Database interface {
	Id() string
	Init(Config)
	Read(key ...string) (string, error)
	Write(value string, key ...string) error
	Close()
}

type Output interface {
	Id() string
	Description() string

	// Batching and streaming
	Read(Domain)
	Write()
	Save(filename string)
	Summary(report map[string]string)
}

type Domain interface {
	Prefix(...string) string
	Name(...string) string
	Suffix(...string) string
	String(...string) string

	// Metadata
	Meta() map[string]string
	SetMeta(key string, value string)
	GetMeta(key string) (value string)
	SetData(key string, value json.RawMessage)
	GetData(key string) (value json.RawMessage)

	Algorithm() Algorithm
	Valid() bool
	Live(...bool) bool
	Active(...bool) bool
	Ld(...int) int
	Json() string
	Idn(...string) string
}

// type Table interface {
// 	Metatable() map[string]interface{}
// 	Set(key string, value interface{})
// 	Get(key string) interface{}
// }

// type Typo interface {
// 	Algo() Algorithm
// 	Set(models.Domain, models.Domain)
// 	Get() (models.Domain, models.Domain)
// 	New(algo Algorithm, original models.Domain, variant models.Domain) Typo

// 	String() string
// 	Dist() int
// 	// Threat() int
// 	Live() bool
// 	Valid() bool

// 	Origin(...string) models.Domain
// 	Derived(...string) models.Domain

// 	// Metatable
// 	Metatable() map[string]interface{}
// 	SetMeta(key string, value interface{})
// 	GetMeta(key string) interface{}
// }

type Accumulator interface {
	Add(Domain)
	// Mkdir(root, dirname string) (string, error)
	// Mkfile(dirname, filename string, content []byte) (string, error)
	Save(filename string, content []byte) error

	Domain() Domain
	GetJson(key string) json.RawMessage
	SetJson(key string, data json.RawMessage)
	Unmarshal(key string, v interface{}) error
	GetMeta(key string) (data string)
	SetMeta(key string, data string)
	Next() (err error)
	Live(...bool) bool
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
