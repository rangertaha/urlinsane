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


package main

import (
	"fmt"
	"os"
	"time"

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
				Usage:   "language IDs to use",
			},
			&cli.StringFlag{
				Name:    "keyboards",
				Aliases: []string{"k"},
				Value:   "all",
				Usage:   "keyboard layout IDs to use",
			},
			&cli.StringFlag{
				Name:    "algorithms",
				Aliases: []string{"a"},
				Value:   "all",
				Usage:   "algorithm IDs to use",
			},
			&cli.StringFlag{
				Name:    "collectors",
				Aliases: []string{"c"},
				Value:   "all",
				Usage:   "collectors IDs to use",
			},
			&cli.IntFlag{
				Name:     "workers",
				Aliases:  []string{"w"},
				Value:    50,
				Category: "PERFORMANCE",
				Usage:    "number of concurrent workers",
			},
			&cli.IntFlag{
				Name:     "random",
				Value:    1,
				Category: "PERFORMANCE",
				Usage:    "random network delay multiplier",
			},
			&cli.IntFlag{
				Name:     "delay",
				Value:    1,
				Category: "PERFORMANCE",
				Usage:    "delay between network calls",
			},
			&cli.DurationFlag{
				Name:     "timeout",
				Aliases:  []string{"t"},
				Value:    5 * time.Second,
				Category: "PERFORMANCE",
				Usage:    "Maximim duration tasks need to complete",
			},
			&cli.DurationFlag{
				Name:     "ttl",
				Value:    168 * time.Hour,
				Category: "PERFORMANCE",
				Usage:    "Maximim duration to cache results",
			},
			&cli.BoolFlag{
				Name:     "progress",
				Aliases:  []string{"p"},
				Value:    false,
				Category: "OUTPUT",
				Usage:    "show progress bar",
			},
			&cli.BoolFlag{
				Name:     "verbose",
				Aliases:  []string{"v"},
				Value:    false,
				Category: "OUTPUT",
				Usage:    "more details in the output",
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"o"},
				Value:    "",
				Category: "OUTPUT",
				Usage:    "output filename defaults to stdout",
			},
			&cli.StringFlag{
				Name:     "format",
				Aliases:  []string{"f"},
				Value:    "table",
				Category: "OUTPUT",
				Usage:    "output format: (csv,tsv,table,txt,html,md,json)",
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

			// Target domain
			// config.Conf.Set("target", )

			// // PLUGINS
			// config.Conf.Set("languages", cCtx.String("languages"))
			// config.Conf.Set("keyboards", cCtx.String("keyboards"))
			// config.Conf.Set("algorithms", cCtx.String("algorithms"))
			// config.Conf.Set("collectors", cCtx.String("collectors"))
			// config.Conf.Set("database", "badger")

			// // PERFORMANCE
			// config.Conf.Set("workers", cCtx.String("workers"))
			// config.Conf.Set("random", cCtx.String("random"))
			// config.Conf.Set("delay", cCtx.String("delay"))
			// config.Conf.Set("timeout", cCtx.String("timeout"))
			// config.Conf.Set("ttl", cCtx.String("ttl"))

			// // OUTPUT
			// config.Conf.Set("verbose", cCtx.String("verbose"))
			// config.Conf.Set("file", cCtx.String("file"))
			// config.Conf.Set("format", cCtx.String("format"))
			// config.Conf.Set("progress", cCtx.String("progress"))

			cfg, err := config.CliConfig(cCtx)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(0)
			}


			t := engine.New(cfg)
			return t.Execute()
		},
	},
}
