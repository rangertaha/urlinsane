// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
EXAMPLE:

    datasets download datasets
    go run ./cmd/datasets download datasets

    datasets import datasets
    go run ./cmd/datasets import datasets

AUTHOR:
   Rangertaha (rangertaha@gmail.com)

     
     `, cli.AppHelpTemplate)

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print the version",
	}

	app := &cli.App{
		Name:        "data",
		Version:     internal.VERSION,
		Compiled:    time.Now(),
		Suggest:     true,
		HelpName:    "data",
		Usage:       "data is used to import and process data models",
		Description: "",
		UsageText:   "data [command] [opts..] [directory]",
		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelpAndExit(ctx, 0)
			return nil
		},
		Commands: []*cli.Command{
			&ImportCmd,
			&DownloadCmd,
			&BuildCmd,
			&LanguagesCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
