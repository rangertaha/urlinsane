// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"path/filepath"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/dataset/gen"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// BuildCmd creates the shipped dataset database from the datasets tree.
//
// Separate from import because it means something different: import adds to
// whatever database is already there, build produces the artifact that gets
// embedded. Building writes to internal/config/dataset.db by default, which is
// the file //go:embed picks up, so the next compile carries it.
var BuildCmd = cli.Command{
	Name:                   "build",
	Aliases:                []string{"b"},
	Usage:                  "Build the shipped dataset database from a datasets tree",
	UsageText:              "build [opt..] [directory]",
	UseShortOptionHandling: true,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "out",
			Value: filepath.Join("internal", "config", "dataset.db"),
			Usage: "where to write the database; the default is the embedded copy"},
	},
	Action: func(cCtx *cli.Context) error {
		root := cCtx.Args().First()
		if root == "" {
			root = "datasets"
		}
		out := cCtx.String("out")

		gen.Progress = func(name, file string, words, transitions int) {
			fmt.Printf("  %s: %d words, %d transitions\n", name, words, transitions)
		}
		fmt.Printf("Building %s from %s\n\n", out, root)
		if err := gen.Build(out, root); err != nil {
			log.Error(err)
			fmt.Println(text.FgRed.Sprint(err))
			return cli.Exit("", 1)
		}
		fmt.Printf("\n%s written\n", out)
		return nil
	},
}
