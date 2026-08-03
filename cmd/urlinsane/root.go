// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	helpTemplates()

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print the version",
	}

	app := &cli.App{
		Name:        "urlinsane",
		Version:     internal.VERSION,
		Compiled:    time.Now(),
		Suggest:     true,
		HelpName:    "urlinsane",
		Usage:       "find the names your target could be mistaken for",
		Description: "",
		UsageText:   "urlinsane [options] <command> [options]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "log debug messages to stderr",
				// It used to be inert: a no-op Action and nothing reading the
				// value, so --debug was a flag the help advertised and the
				// program ignored.
				Action: func(_ *cli.Context, v bool) error {
					if v {
						log.SetLevel(log.DebugLevel)
					}
					return nil
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelpAndExit(ctx, 0)
			return nil
		},
		Commands: []*cli.Command{
			&TypoCmd,
			&ReportCmd,
		},
	}

	// HandleExitCoder, not log.Fatal: exit codes are part of the interface
	// (§12.4), and log.Fatal collapses every one of them to 1 — so --fail-on
	// could never report 2 and the tool could not be a CI gate. It also stamps
	// a timestamp on messages meant for a human.
	if err := app.Run(reorder(app.Commands, os.Args)); err != nil {
		if _, ok := err.(cli.ExitCoder); ok {
			cli.HandleExitCoder(err)
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
