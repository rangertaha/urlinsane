// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/dataset/gen"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var importFlags = []cli.Flag{}

// ImportCmd loads a datasets tree into the database.
//
// One pass over every .lst, into Vocabulary and Transition. There used to be an
// importer per relation -- Synonyms, Antonyms, Homoglyphs and fifteen more --
// each filling its own table through a many2many association. Those tables
// recorded only that two words were related; the two that replaced them record
// how strongly, and one generic walk fills them for every dataset, including
// the ones no per-relation importer covered.
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

		if _, err := config.Init(); err != nil {
			log.Error(err)
			fmt.Println(text.FgRed.Sprint(err))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}
		gen.Progress = func(name, file string, words, transitions int) {
			fmt.Printf("  %s (%s): %d words, %d transitions\n", name, file, words, transitions)
		}
		return gen.All(cCtx.Args().First())
	},
}
