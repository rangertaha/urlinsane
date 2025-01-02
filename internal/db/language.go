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
package db

type Language struct {
	ID          uint
	Code        string `gorm:"unique" json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	Keyboards []*Keyboard `gorm:"many2many:langboards;" json:"keyboads,omitempty"`
	Stopwords []*Word     `gorm:"many2many:stopwords;"  json:"stopwords,omitempty"`
	Numerals  []*Word     `gorm:"many2many:numerals;"   json:"numerals,omitempty"`
	Words     []*Word     `gorm:"many2many:langwords;"  json:"words,omitempty"`
	Graphemes []*Char     `gorm:"many2many:graphemes;"  json:"graphemes,omitempty"`
	Vowels    []*Char     `gorm:"many2many:langvowels;" json:"vowels,omitempty"`
}

type Char struct {
	ID         uint
	// Code       string      `gorm:"unique" json:"code,omitempty"`
	Text       string      `gorm:"unique" json:"text,omitempty"`
	Languages  []*Language `gorm:"many2many:langchars;"  json:"languages,omitempty"`
	Homoglyphs []*Char     `gorm:"many2many:homoglyphs;" json:"homoglyphs,omitempty"`
}

type Word struct {
	ID           uint
	Text         string      `gorm:"unique" json:"text,omitempty"`
	Languages    []*Language `gorm:"many2many:langwords;"    json:"languages,omitempty"`
	Antonyms     []*Word     `gorm:"many2many:antonyms;"     json:"antonyms,omitempty"`
	Homophones   []*Word     `gorm:"many2many:homophones;"   json:"homophones,omitempty"`
	Misspellings []*Word     `gorm:"many2many:misspellings;" json:"misspellings,omitempty"`
	Translations []*Word     `gorm:"many2many:translations;" json:"translations,omitempty"`
}
