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

var importFlags = []cli.Flag{
	// &cli.StringFlag{
	// 	Name:    "language",
	// 	Usage:   "language code of the content to import",
	// 	Aliases: []string{"l"},
	// 	Value:   "",
	// },
	// &cli.BoolFlag{
	// 	Name:  "clear",
	// 	Usage: "Remove previously imported datasets",
	// 	Value: false,
	// },
	// &cli.BoolFlag{
	// 	Name:     "words",
	// 	Value:    false,
	// 	Category: "TYPES",
	// 	Usage:    "Import dataset as words",
	// },
	// &cli.BoolFlag{
	// 	Name:     "stopwords",
	// 	Value:    false,
	// 	Category: "TYPES",
	// 	Usage:    "Import dataset as stopwords",
	// },
	// &cli.BoolFlag{
	// 	Name:     "graphemes",
	// 	Value:    false,
	// 	Category: "TYPES",
	// 	Usage:    "Import dataset as graphemes",
	// },
	// &cli.BoolFlag{
	// 	Name:     "vowels",
	// 	Value:    false,
	// 	Category: "TYPES",
	// 	Usage:    "Import dataset as vowels",
	// },
}

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

func Directory(dir, name string, processor func(paths []string) error) (err error) {
	var files []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if strings.Contains(path, name) {
				files = append(files, path)
			}

		}
		return nil
	})
	if err != nil {
		return err
	}
	return processor(files)
}

func Import(cli *cli.Context) error {
	folder := cli.Args().First()
	if err := Directory(folder, "languages", Languages); err != nil {
		return err
	}

	return nil
}

func Languages(files []string) (err error) {
	type importer func(string, string) error
	datasets := map[string]importer{
		"words.lst": Words,
	}
	for _, file := range files {
		segs := strings.Split(file, "/")
		language := segs[len(segs)-2]

		for name, ifunc := range datasets {
			if strings.HasSuffix(file, name) {
				ifunc(language, file)
			}
		}
	}
	return
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

		lines = append(lines, strings.Split(line, " "))
	}
	readFile.Close()
	return
}

func Words(lang, file string) (err error) {
	for _, words := range Extract(file) {
		fmt.Println(lang, words)
	}

	return
}
