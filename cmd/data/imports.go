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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/dataset"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type importFunc func(string, string) error

var DATASETS = map[string]map[string]importFunc{
	"languages": {
		"word.lst":        Words,
		"antonym.lst":     Antonyms,
		"stopword.lst":    StopWords,
		"numeral.lst":     Numerals,
		"misspelling.lst": Misspellings,
		"homophone.lst":   Homophones,
		"vowel.lst":       Vowels,
		"grapheme.lst":    Graphemes,
		"homoglyph.lst":   Homoglyphs,
	},
	"domains": {
		"domain.lst": Domains,
		"prefix.lst": Prefixes,
		"suffix.lst": Suffixes,
	},
}

var importFlags = []cli.Flag{}

var ImportCmd = cli.Command{
	Name:                   "import",
	Aliases:                []string{"i"},
	Usage:                  "Import datasets into the database",
	UsageText:              "import [opt..] [directory]",
	UseShortOptionHandling: true,
	Flags:                  importFlags,
	Action: func(cCtx *cli.Context) error {
		if cCtx.NArg() == 0 {
			fmt.Println(text.FgRed.Sprint("\n  a directory is needed!\n"))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}

		_, err := config.New(config.CliOptions(cCtx))
		if err != nil {
			log.Error(err)
			fmt.Println(text.FgRed.Sprint(err))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}
		return Import(cCtx)
	},
}

func Import(cli *cli.Context) error {
	folder := cli.Args().First()
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			var namespace string
			for dir, file := range DATASETS {
				if strings.Contains(path, dir) {
					for filename, processor := range file {
						paths := strings.Split(path, "/")
						if paths[len(paths)-1] == filename {
							if paths[len(paths)-2] != dir {
								namespace = paths[len(paths)-2]
							} else {
								namespace = ""
							}
							processor(namespace, path)
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func Extract(file string) (lines [][]string) {
	readFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		m := regexp.MustCompile(`\s+`)
		line = m.ReplaceAllString(line, " ")
		line = strings.Trim(line, " ")
		line = strings.TrimSpace(line)

		lines = append(lines, strings.Split(line, " "))
	}
	readFile.Close()
	return
}

//
// Domains
//

func Domains(language, file string) (err error) {
	fmt.Printf("Importing domains from %s\n", file)
	var domains []*dataset.Domain
	for _, line := range Extract(file) {
		for _, domain := range line {
			var dm dataset.Domain
			dataset.Data.FirstOrInit(&dm, dataset.Domain{Name: domain})
			domains = append(domains, &dm)
		}
	}

	for chunk := range slices.Chunk(domains, 1000) {
		fmt.Printf("Saving %v domains \n", len(chunk))
		dataset.Data.Save(&chunk)
		time.Sleep(time.Second)
	}

	return
}

func Suffixes(language, file string) (err error) {
	fmt.Printf("Importing suffixes from %s\n", file)
	var suffixes []*dataset.Suffix
	for _, line := range Extract(file) {
		for _, suffix := range line {
			var s dataset.Suffix
			dataset.Data.FirstOrInit(&s, dataset.Suffix{Name: suffix})
			suffixes = append(suffixes, &s)
		}
	}

	for chunk := range slices.Chunk(suffixes, 1000) {
		fmt.Printf("Saving %v suffixes \n", len(chunk))
		dataset.Data.Save(&chunk)
		time.Sleep(time.Second)
	}

	return
}

func Prefixes(language, file string) (err error) {
	fmt.Printf("Importing prefixes from %s\n", file)
	var prefixes []*dataset.Prefix
	for _, line := range Extract(file) {
		for _, domain := range line {
			var p dataset.Prefix
			dataset.Data.FirstOrInit(&p, dataset.Prefix{Name: domain})
			prefixes = append(prefixes, &p)
		}
	}

	for chunk := range slices.Chunk(prefixes, 1000) {
		fmt.Printf("Saving %v prefixes \n", len(chunk))
		dataset.Data.Save(&chunk)
		time.Sleep(time.Second)
	}

	return
}

//
// Language
//

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
