// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/dataset/gen"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// LanguagesCmd creates the language directories the keyboard catalogue implies.
//
// Build runs this too, so the tree is scaffolded whenever the database is
// rebuilt. It is a command of its own because the two have different reasons to
// be run: adding a language is a change to the tree that wants reviewing on its
// own, without a rebuilt database in the same commit.
var LanguagesCmd = cli.Command{
	Name:                   "languages",
	Aliases:                []string{"l", "lang"},
	Usage:                  "Create a directory of .lst files for every language with a keyboard",
	UsageText:              "languages [opt..] [directory]",
	UseShortOptionHandling: true,
	Description: "Every language pkg/kb ships a keyboard layout for gets\n" +
		"languages/<code>/, holding one empty .lst per relation.\n" +
		"Existing files are never touched, so this is safe to run\n" +
		"over a curated tree.",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"n"},
			Usage:   "report what would be created without writing anything"},
	},
	Action: func(cCtx *cli.Context) error {
		root := cCtx.Args().First()
		if root == "" {
			root = "datasets"
		}

		if cCtx.Bool("dry-run") {
			missing, err := gen.Missing(root)
			if err != nil {
				log.Error(err)
				fmt.Println(text.FgRed.Sprint(err))
				return cli.Exit("", 1)
			}
			report(root, missing, "would be created")
			return nil
		}

		created, err := gen.Scaffold(root)
		if err != nil {
			log.Error(err)
			fmt.Println(text.FgRed.Sprint(err))
			return cli.Exit("", 1)
		}
		report(root, created, "created")
		return nil
	},
}

// report prints the codes one per line rather than as a paragraph, because the
// list is usually either empty or eighty long and the eighty-long case is worth
// reading.
func report(root string, codes []string, verb string) {
	if len(codes) == 0 {
		fmt.Printf("%s/languages already covers every language with a keyboard\n", root)
		return
	}
	for _, c := range codes {
		fmt.Printf("  %s\n", c)
	}
	fmt.Printf("\n%d language(s) %s under %s/languages, %d files each\n",
		len(codes), verb, root, len(gen.Relations))
}
