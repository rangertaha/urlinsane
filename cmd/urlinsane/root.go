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
	"os"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
EXAMPLE:

    urlinsane typo example.com
    urlinsane typo -a co example.com
    urlinsane typo -a co,oi,oy -c ip,idna,ns example.com
    urlinsane typo -l fr,en -k en1,en2 example.com

AUTHOR:
   Rangertaha (rangertaha@gmail.com)

     
     `, cli.AppHelpTemplate)

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
		Usage:       "Urlinsane is an advanced cybersecurity typosquatting tool",
		Description: "",
		UsageText:   "urlinsane [global opts..] [command] [opts..]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "Log debug messags for development",
				Action: func(ctx *cli.Context, v bool) error {
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
		},
	}

	// HandleExitCoder, not log.Fatal: exit codes are part of the interface
	// (§12.4), and log.Fatal collapses every one of them to 1 — so --fail-on
	// could never report 2 and the tool could not be a CI gate. It also stamps
	// a timestamp on messages meant for a human.
	if err := app.Run(os.Args); err != nil {
		if _, ok := err.(cli.ExitCoder); ok {
			cli.HandleExitCoder(err)
			return
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
