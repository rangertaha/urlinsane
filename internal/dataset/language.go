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
package dataset

type Language struct {
	ID          uint
	Code        string `gorm:"unique"`
	Name        string
	Description string

	Keyboards []*Keyboard `gorm:"many2many:langboards;"`
	Stopwords []*Word     `gorm:"many2many:stopwords;"`
	Numerals  []*Word     `gorm:"many2many:numerals;"`
	Words     []*Word     `gorm:"many2many:langwords;"`
	Graphemes []*Char     `gorm:"many2many:graphemes;"`
	Vowels    []*Char     `gorm:"many2many:vowels;"`
}

type Sym struct {
	ID    uint
	Name  string  `gorm:"unique"`
	Value string  `gorm:"unique"`
	Words []*Char `gorm:"many2many:symbols;"`
}

type Char struct {
	ID         uint
	Text       string      `gorm:"unique"`
	Homoglyphs []*Char     `gorm:"many2many:homoglyphs;"`
	Languages  []*Language `gorm:"many2many:graphemes;"`
}

type Word struct {
	ID           uint
	Text         string      `gorm:"unique"`
	Symbols      []*Sym      `gorm:"many2many:symbols;"`
	Languages    []*Language `gorm:"many2many:langwords;"`
	Antonyms     []*Word     `gorm:"many2many:antonyms;"`
	Homophones   []*Word     `gorm:"many2many:homophones;"`
	Misspellings []*Word     `gorm:"many2many:misspellings;"`
	Translations []*Word     `gorm:"many2many:translations;"`
}
