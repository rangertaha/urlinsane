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

package main

import (
	"fmt"
	"strings"

	"github.com/rangertaha/urlinsane/internal/db"
)

func Words(language, file string) (err error) {
	lng := &db.Language{Code: language}
	db.DB.FirstOrCreate(lng)

	fmt.Println("Importing words...")
	for _, wordslist := range Extract(file) {
		for _, w := range wordslist {
			if strings.TrimSpace(w) != "" {
				var word db.Word
				db.DB.FirstOrInit(&word, db.Word{Text: w})
				lng.Words = append(lng.Words, &word)
			}
		}
	}
	db.DB.Save(&lng)
	return
}

func Vowels(language, file string) (err error) {
	lng := &db.Language{Code: language}
	db.DB.FirstOrCreate(lng)

	fmt.Println("Importing vowels...")
	for _, lines := range Extract(file) {
		for _, c := range lines {
			var char db.Char
			db.DB.FirstOrInit(&char, db.Char{Text: c})
			lng.Vowels = append(lng.Vowels, &char)
		}
	}
	db.DB.Save(&lng)
	return
}

func Graphemes(language, file string) (err error) {
	lng := &db.Language{Code: language}
	db.DB.FirstOrCreate(lng)

	fmt.Println("Importing graphemes...")
	for _, lines := range Extract(file) {
		for _, c := range lines {
			var char db.Char
			db.DB.FirstOrInit(&char, db.Char{Text: c})
			lng.Graphemes = append(lng.Graphemes, &char)
		}
	}
	db.DB.Save(&lng)
	return
}

func Antonyms(language, file string) (err error) {
	lng := &db.Language{Code: language}
	db.DB.FirstOrCreate(lng)

	fmt.Println("Importing antonyms...")
	var words []*db.Word
	for _, wordslist := range Extract(file) {
		var word db.Word
		db.DB.FirstOrInit(&word, db.Word{Text: wordslist[0]})
		for _, w := range wordslist[1:] {
			var related db.Word
			db.DB.FirstOrInit(&related, db.Word{Text: w})
			word.Antonyms = append(word.Antonyms, &related)
		}
		words = append(words, &word)
	}
	db.DB.Save(&words)
	return
}

func Homophones(language, file string) (err error) {
	lng := &db.Language{Code: language}
	db.DB.FirstOrCreate(lng)

	fmt.Println("Importing homophones...")
	var words []*db.Word
	for _, wordslist := range Extract(file) {
		var word db.Word
		db.DB.FirstOrInit(&word, db.Word{Text: wordslist[0]})
		for _, w := range wordslist[1:] {
			var related db.Word
			db.DB.FirstOrInit(&related, db.Word{Text: w})
			word.Homophones = append(word.Homophones, &related)
		}
		words = append(words, &word)
	}
	db.DB.Save(&words)
	return
}

func Homoglyphs(language, file string) (err error) {
	lng := &db.Language{Code: language}
	db.DB.FirstOrCreate(lng)

	fmt.Println("Importing homoglyphs...")
	var chars []*db.Char

	for _, lines := range Extract(file) {
		var char db.Char
		db.DB.FirstOrInit(&char, db.Char{Text: lines[0]})
		for _, c := range lines[1:] {
			var related db.Char
			db.DB.FirstOrInit(&related, db.Char{Text: c})
			char.Homoglyphs = append(char.Homoglyphs, &related)
		}
		chars = append(chars, &char)

	}
	db.DB.Save(&chars)
	return
}

func StopWords(language, file string) (err error) {
	lng := &db.Language{Code: language}
	db.DB.FirstOrCreate(lng)

	fmt.Println("Importing stopwords...")
	for _, wordslist := range Extract(file) {
		for _, w := range wordslist {
			var word db.Word
			db.DB.FirstOrInit(&word, db.Word{Text: w})
			lng.Stopwords = append(lng.Stopwords, &word)
		}
	}
	db.DB.Save(&lng)
	return
}

func Numerals(language, file string) (err error) {
	lng := &db.Language{Code: language}
	db.DB.FirstOrCreate(lng)

	fmt.Println("Importing numerals...")
	for _, wordslist := range Extract(file) {
		for _, w := range wordslist {
			var word db.Word
			db.DB.FirstOrInit(&word, db.Word{Text: w})
			lng.Numerals = append(lng.Numerals, &word)
		}
	}
	db.DB.Save(&lng)
	return
}

func Misspellings(language, file string) (err error) {
	lng := &db.Language{Code: language}
	db.DB.FirstOrCreate(lng)

	fmt.Println("Importing misspellings...")
	var words []*db.Word
	for _, wordslist := range Extract(file) {
		var word db.Word
		db.DB.FirstOrInit(&word, db.Word{Text: wordslist[0]})
		for _, w := range wordslist[1:] {
			var related db.Word
			db.DB.FirstOrInit(&related, db.Word{Text: w})
			word.Misspellings = append(word.Misspellings, &related)
		}
		words = append(words, &word)
	}
	db.DB.Save(&words)
	return
}
