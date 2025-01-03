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

	"github.com/rangertaha/urlinsane/internal/dataset"
)

func Words(language, file string) (err error) {
	lng := &dataset.Language{Code: language}
	dataset.Data.FirstOrCreate(lng)

	fmt.Printf("Importing words from %s\n", file)
	for _, wordslist := range Extract(file) {
		for _, w := range wordslist {
			if strings.TrimSpace(w) != "" {
				var word dataset.Word
				dataset.Data.FirstOrInit(&word, dataset.Word{Text: w})
				lng.Words = append(lng.Words, &word)
			}
		}
	}
	dataset.Data.Save(&lng)
	return
}

func Vowels(language, file string) (err error) {
	lng := &dataset.Language{Code: language}
	dataset.Data.FirstOrCreate(lng)

	// fmt.Println("Importing vowels...")
	fmt.Printf("Importing vowels from %s\n", file)
	for _, lines := range Extract(file) {
		for _, c := range lines {
			var char dataset.Char
			dataset.Data.FirstOrInit(&char, dataset.Char{Text: c})
			lng.Vowels = append(lng.Vowels, &char)
		}
	}
	dataset.Data.Save(&lng)
	return
}

func Graphemes(language, file string) (err error) {
	lng := &dataset.Language{Code: language}
	dataset.Data.FirstOrCreate(lng)

	// fmt.Println("Importing graphemes...")
	fmt.Printf("Importing graphemes from %s\n", file)
	for _, lines := range Extract(file) {
		for _, c := range lines {
			var char dataset.Char
			dataset.Data.FirstOrInit(&char, dataset.Char{Text: c})
			lng.Graphemes = append(lng.Graphemes, &char)
		}
	}
	dataset.Data.Save(&lng)
	return
}

func Antonyms(language, file string) (err error) {
	lng := &dataset.Language{Code: language}
	dataset.Data.FirstOrCreate(lng)

	// fmt.Println("Importing antonyms...")
	fmt.Printf("Importing antonyms from %s\n", file)
	var words []*dataset.Word
	for _, wordslist := range Extract(file) {
		var word dataset.Word
		dataset.Data.FirstOrInit(&word, dataset.Word{Text: wordslist[0]})
		for _, w := range wordslist[1:] {
			var related dataset.Word
			dataset.Data.FirstOrInit(&related, dataset.Word{Text: w})
			word.Antonyms = append(word.Antonyms, &related)
		}
		words = append(words, &word)
	}
	dataset.Data.Save(&words)
	return
}

func Homophones(language, file string) (err error) {
	lng := &dataset.Language{Code: language}
	dataset.Data.FirstOrCreate(lng)

	// fmt.Println("Importing homophones...")
	fmt.Printf("Importing homophones from %s\n", file)
	var words []*dataset.Word
	for _, wordslist := range Extract(file) {
		var word dataset.Word
		dataset.Data.FirstOrInit(&word, dataset.Word{Text: wordslist[0]})
		for _, w := range wordslist[1:] {
			var related dataset.Word
			dataset.Data.FirstOrInit(&related, dataset.Word{Text: w})
			word.Homophones = append(word.Homophones, &related)
		}
		words = append(words, &word)
	}
	dataset.Data.Save(&words)
	return
}

func Homoglyphs(language, file string) (err error) {
	lng := &dataset.Language{Code: language}
	dataset.Data.FirstOrCreate(lng)

	// fmt.Println("Importing homoglyphs...")
	fmt.Printf("Importing homoglyphs from %s\n", file)
	var chars []*dataset.Char

	for _, lines := range Extract(file) {
		var char dataset.Char
		dataset.Data.FirstOrInit(&char, dataset.Char{Text: lines[0]})
		for _, c := range lines[1:] {
			var related dataset.Char
			dataset.Data.FirstOrInit(&related, dataset.Char{Text: c})
			char.Homoglyphs = append(char.Homoglyphs, &related)
		}
		chars = append(chars, &char)

	}
	dataset.Data.Save(&chars)
	return
}

func StopWords(language, file string) (err error) {
	lng := &dataset.Language{Code: language}
	dataset.Data.FirstOrCreate(lng)

	// fmt.Println("Importing stopwords...")
	fmt.Printf("Importing stopwords from %s\n", file)
	for _, wordslist := range Extract(file) {
		for _, w := range wordslist {
			var word dataset.Word
			dataset.Data.FirstOrInit(&word, dataset.Word{Text: w})
			lng.Stopwords = append(lng.Stopwords, &word)
		}
	}
	dataset.Data.Save(&lng)
	return
}

func Numerals(language, file string) (err error) {
	lng := &dataset.Language{Code: language}
	dataset.Data.FirstOrCreate(lng)

	// fmt.Println("Importing numerals...")
	fmt.Printf("Importing numerals from %s\n", file)
	for _, wordslist := range Extract(file) {
		for _, w := range wordslist {
			var word dataset.Word
			dataset.Data.FirstOrInit(&word, dataset.Word{Text: w})
			lng.Numerals = append(lng.Numerals, &word)
		}
	}
	dataset.Data.Save(&lng)
	return
}

func Misspellings(language, file string) (err error) {
	lng := &dataset.Language{Code: language}
	dataset.Data.FirstOrCreate(lng)

	// fmt.Println("Importing misspellings...")
	fmt.Printf("Importing misspellings from %s\n", file)
	var words []*dataset.Word
	for _, wordslist := range Extract(file) {
		var word dataset.Word
		dataset.Data.FirstOrInit(&word, dataset.Word{Text: wordslist[0]})
		for _, w := range wordslist[1:] {
			var related dataset.Word
			dataset.Data.FirstOrInit(&related, dataset.Word{Text: w})
			word.Misspellings = append(word.Misspellings, &related)
		}
		words = append(words, &word)
	}
	dataset.Data.Save(&words)
	return
}
