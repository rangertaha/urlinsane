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
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type importFunc func(string, string) error

var DATASETS = map[string]map[string]importFunc{
	// Language datasets
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

// func Directory(dir, name string, processor func(paths []string) error) (err error) {
// 	var files []string
// 	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !info.IsDir() {
// 			if strings.Contains(path, name) {
// 				files = append(files, path)
// 			}

// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return processor(files)
// }

// func Datasets(files []string) (err error) {

// 	for _, file := range files {
// 		segs := strings.Split(file, "/")
// 		language := segs[len(segs)-2]

// 		for name, ifunc := range DATASETS {
// 			if strings.HasSuffix(file, name) {
// 				ifunc(language, file)
// 			}
// 		}
// 	}
// 	return
// }

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
