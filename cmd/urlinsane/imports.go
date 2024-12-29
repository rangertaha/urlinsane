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

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/db/imports"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var importFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "language",
		Usage:   "language code of the content to import",
		Aliases: []string{"l"},
		Value:   "",
	},
	&cli.BoolFlag{
		Name:  "clear",
		Usage: "Remove previously imported datasets",
		Value: false,
	},
	&cli.BoolFlag{
		Name:     "words",
		Value:    false,
		Category: "TYPES",
		Usage:    "Import dataset as words",
	},
	&cli.BoolFlag{
		Name:     "stopwords",
		Value:    false,
		Category: "TYPES",
		Usage:    "Import dataset as stopwords",
	},
	&cli.BoolFlag{
		Name:     "graphemes",
		Value:    false,
		Category: "TYPES",
		Usage:    "Import dataset as graphemes",
	},
	&cli.BoolFlag{
		Name:     "vowels",
		Value:    false,
		Category: "TYPES",
		Usage:    "Import dataset as vowels",
	},
}

var ImportCmd = cli.Command{
	Name:                   "import",
	Aliases:                []string{"i"},
	Usage:                  "Import datasets into the database",
	UsageText:              "import [opt..] [file]",
	UseShortOptionHandling: true,
	Flags:                  importFlags,
	Action: func(cCtx *cli.Context) error {
		if cCtx.NArg() == 0 {
			fmt.Println(text.FgRed.Sprint("\n  a file is needed!\n"))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}

		_, err := config.New(config.CliOptions(cCtx))
		if err != nil {
			log.Error(err)
			fmt.Println(text.FgRed.Sprint(err))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}
		return imports.Import(cCtx)
	},
}
