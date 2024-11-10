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
	"os"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/engine"
	"github.com/urfave/cli/v2"
)

var TypoCmd = []*cli.Command{
	{
		Name:        "typo",
		Aliases:     []string{"t"},
		Usage:       "Generate domain variations and collect information on them",
		Description: "URLInsane is designed to detect domain typosquatting by using advanced algorithms, information-gathering techniques, and data analysis to identify potentially harmful variations of targeted domains that cybercriminals might exploit. This tool is essential for defending against threats like typosquatting, brandjacking, URL hijacking, fraud, phishing, and corporate espionage. By detecting malicious domain variations, it provides an added layer of protection to brand integrity and user trust. Additionally, URLInsane enhances threat intelligence capabilities, strengthening proactive cybersecurity measures.",
		UsageText:   "urlinsane typo [options..] [domain]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "languages",
				Aliases: []string{"l"},
				Value:   "all",
				Usage:   "IDs of languages to use",
			},
			&cli.StringFlag{
				Name:    "keyboards",
				Aliases: []string{"k"},
				Value:   "all",
				Usage:   "IDs of keyboard layouts to use",
			},
			&cli.StringFlag{
				Name:    "algorithms",
				Aliases: []string{"a"},
				Value:   "all",
				Usage:   "IDs of typo algorithms to use",
			},
			&cli.IntFlag{
				Name:     "concurrency",
				Aliases:  []string{"c"},
				Value:    50,
				Category: "PERFORMANCE",
				Usage:    "Number of concurrent workers",
			},
			&cli.IntFlag{
				Name:     "random",
				Value:    1,
				Category: "PERFORMANCE",
				Usage:    "Random network delay multiplier",
			},
			&cli.IntFlag{
				Name:     "delay",
				Value:    1,
				Category: "PERFORMANCE",
				Usage:    "Delay between network calls",
			},
			&cli.BoolFlag{
				Name:     "progress",
				Aliases:  []string{"p"},
				Value:    false,
				Category: "OUTPUT",
				Usage:    "Show progress bar",
			},
			&cli.BoolFlag{
				Name:     "verbose",
				Aliases:  []string{"v"},
				Value:    false,
				Category: "OUTPUT",
				Usage:    "More details in the output",
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"o"},
				Value:    "",
				Category: "OUTPUT",
				Usage:    "Output filename defaults to stdout",
			},
			&cli.StringFlag{
				Name:     "format",
				Aliases:  []string{"f"},
				Value:    "table",
				Category: "OUTPUT",
				Usage:    "Output format: (csv,tsv,table,txt,html,md,json)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() == 0 {
				fmt.Println(text.FgRed.Sprint("\n  a domain name is needed!\n"))
				cli.ShowAppHelpAndExit(cCtx, 0)

			}
			if cCtx.NArg() > 1 {
				fmt.Println(text.FgRed.Sprint("\n  only one domain name at at time!\n"))
			}

			config.Conf.Set("languages", cCtx.String("languages"))
			config.Conf.Set("keyboards", cCtx.String("keyboards"))
			config.Conf.Set("algorithms", cCtx.String("algorithms"))
			config.Conf.Set("concurrency", cCtx.String("concurrency"))
			config.Conf.Set("random", cCtx.String("random"))
			config.Conf.Set("delay", cCtx.String("delay"))
			config.Conf.Set("verbose", cCtx.String("verbose"))
			config.Conf.Set("file", cCtx.String("file"))
			config.Conf.Set("format", cCtx.String("format"))
			config.Conf.Set("progress", cCtx.String("progress"))
			// fmt.Println("banner", config.Conf.Get("banner"))
			// fmt.Println("database", config.Conf.Get("database"))

			// Target domain
			config.Conf.Set("target", cCtx.Args().First())

			cfg, err := config.CliConfig(config.Conf.String("target"))
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(0)
			}

			t := engine.New(cfg)
			t.Execute()

			// config.Conf.Print()

			// fmt.Println("DOMAIN: ", config.Conf.String("target"))
			return nil
		},
	},
}
