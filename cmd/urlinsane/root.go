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
	app := newApp()

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

// newApp builds the command tree.
//
// Separate from main so a test can run the app without running the process:
// main's only remaining job is to turn an error into an exit code, and the
// dispatch behaviour worth testing — that an unknown command fails rather than
// printing help and succeeding — lives in here.
func newApp() *cli.App {
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
		// Reached with no command at all, and also with one this app does not
		// have: urfave/cli falls through to the app Action when it cannot match
		// a subcommand, so the two cases arrive here together and have to be
		// told apart.
		//
		// Showing the help and exiting 0 for both meant every mistyped verb was
		// a silently successful run. `urlinsane typoo acme.com --fail-on high`
		// printed the help to stdout and exited 0, so a CI gate keyed on the
		// 0/1/2 contract recorded a clean pass for a scan that never happened —
		// and `-o json > out.json` captured help text into the file while the
		// job went green. Suggest:true does not cover it either; urfave only
		// consults SuggestCommand from ShowCommandHelp, never on this path.
		Action: func(ctx *cli.Context) error {
			if name := ctx.Args().First(); name != "" {
				return exit(fmt.Errorf(
					"unknown command %q; see `urlinsane --help` for the commands this tool has", name),
					exitError)
			}
			cli.ShowAppHelpAndExit(ctx, 0)
			return nil
		},
		Commands: []*cli.Command{
			&TypoCmd,
			&ReportCmd,
		},
	}

	return app
}
