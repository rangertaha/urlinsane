// Copyright (C) 2024 Rangertaha
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

/*
Urlinsane OSINT tool
It uses tabs for indentation and blanks for alignment.
Alignment assumes that an editor is using a fixed-width font.

Without an explicit path, it processes the standard input. Given a file,
it operates on that file; given a directory, it operates on all .go files in
that directory, recursively. (Files starting with a period are ignored.)
By default, gofmt prints the reformatted sources to standard output.

Usage:

	gofmt [flags] [path ...]

The flags are:

	-d
	    Do not print reformatted sources to standard output.
	    If a file's formatting is different than gofmt's, print diffs
	    to standard output.
	-w
	    Do not print reformatted sources to standard output.
	    If a file's formatting is different from gofmt's, overwrite it
	    with gofmt's version. If an error occurred during overwriting,
	    the original file is restored from an automatic backup.

When gofmt reads from standard input, it accepts either a full Go program
or a program fragment. A program fragment must be a syntactically
valid declaration list, statement list, or expression. When formatting
such a fragment, gofmt preserves leading indentation as well as leading
and trailing spaces, so that individual sections of a Go program can be
formatted by piping them through gofmt.
*/
package main

// import "github.com/rangertaha/urlinsane/cmd/urlinsane"

// func main() {
// 	urlinsane.Execute()
// }

// package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/urfave/cli/v2"
)

func main() {
	config.LoadOrCreateConfig(config.Conf, config.DIR, config.FILE, config.DefualtConfig)

	// EXAMPLE: Append to an existing template

	// // EXAMPLE: Override a template
	// cli.AppHelpTemplate = `NAME:
	//     {{.Name}} - {{.Usage}}
	//  USAGE:
	//     {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
	//     {{if len .Authors}}
	//  AUTHOR:
	//     {{range .Authors}}{{ . }}{{end}}
	//     {{end}}{{if .Commands}}
	//  COMMANDS:
	//  {{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
	//  GLOBAL OPTIONS:
	//     {{range .VisibleFlags}}{{.}}
	//     {{end}}{{end}}{{if .Copyright }}
	//  COPYRIGHT:
	//     {{.Copyright}}
	//     {{end}}{{if .Version}}
	//  VERSION:
	//     {{.Version}}
	//     {{end}}
	//  `
	cli.AppHelpTemplate = fmt.Sprintf(`%s
EXAMPLE:

    urlinsane typo example.com
    urlinsane typo -a co example.com
    urlinsane typo -t co,oi,oy -i ip,idna,ns example.com
    urlinsane typo -l fr,en -k en1,en2 example.com

AUTHOR:
   Rangertaha (rangertaha@gmail.com)

     
     `, cli.AppHelpTemplate)

	// // EXAMPLE: Replace the `HelpPrinter` func
	// cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
	// 	fmt.Println("Ha HA.  I pwnd the help!!1")
	// }

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print the version",
	}

	app := &cli.App{
		Name:     "urlinsane",
		Version:  internal.VERSION,
		Compiled: time.Now(),
		Suggest:  true,
		// Authors: []*cli.Author{
		// 	{
		// 		Name:  "Rangertaha",
		// 		Email: "rangertaha@gmail.com",
		// 	},
		// },
		// Copyright: "(c) 2024 Rangertaha <rangertaha@gmail.com>",
		HelpName:    "urlinsane",
		Usage:       "Urlinsane is an advanced cybersecurity typosquatting tool",
		Description: "",
		UsageText:   "urlinsane [command] [options..]",
		Flags:       []cli.Flag{
			// &cli.BoolFlag{
			// 	Name:     "ll",
			// 	Aliases:  []string{"p"},
			// 	Value:    false,
			// 	Category: "OUTPUT",
			// 	Usage:    "Show progress bar",
			// },
		},
		Action: func(ctx *cli.Context) error {
			cli.ShowAppHelpAndExit(ctx, 0)
			return nil
		},
		Commands: TypoCmd,
	}

	// man, _ := app.ToMan()
	// fmt.Println(man)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
