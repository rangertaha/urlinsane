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
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/engine"
	"github.com/rangertaha/urlinsane/internal/pkg"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/internal/plugins/algorithms/all"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/all"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
	_ "github.com/rangertaha/urlinsane/internal/plugins/outputs/all"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "languages",
		Aliases: []string{"l"},
		Value:   "en",
		Usage:   "language IDs to use `[ID]`",
	},
	&cli.StringFlag{
		Name:    "keyboards",
		Aliases: []string{"k"},
		Value:   "en1,en2,en3,en4",
		Usage:   "keyboard layout IDs to use `[ID]`",
	},
	&cli.StringFlag{
		Name:    "algorithms",
		Aliases: []string{"a"},
		Value:   "all",
		Usage:   "algorithm IDs to use `[ID]`",
	},
	&cli.StringFlag{
		Name:    "collectors",
		Aliases: []string{"c"},
		Value:   "",
		Usage:   "collectors IDs to use `[ID]`",
	},
	&cli.StringFlag{
		Name:     "regex",
		Aliases:  []string{"e"},
		Value:    "",
		Category: "CONSTRAINTS",
		Usage:    "regular expressions to match `[PATTERN]`",
	},
	&cli.IntFlag{
		Name:     "workers",
		Aliases:  []string{"w"},
		Value:    50,
		Category: "PERFORMANCE",
		Usage:    "number of concurrent workers `NUM`",
	},
	&cli.IntFlag{
		Name:     "random",
		Value:    1,
		Category: "PERFORMANCE",
		Usage:    "random network delay multiplier `NUM`",
	},
	&cli.IntFlag{
		Name:     "delay",
		Value:    1,
		Category: "PERFORMANCE",
		Usage:    "delay between network calls `NUM`",
	},
	&cli.DurationFlag{
		Name:     "timeout",
		Aliases:  []string{"t"},
		Value:    0 * time.Minute,
		Category: "PERFORMANCE",
		Hidden:   true,
		Usage:    "maximim duration tasks need to complete `DURATION`",
	},
	&cli.DurationFlag{
		Name:     "ttl",
		Value:    0 * time.Hour,
		Category: "PERFORMANCE",
		Usage:    "duration to cache results, 0 clears the cache `DURATION`",
	},
	&cli.IntFlag{
		Name:     "distance",
		Aliases:  []string{"d"},
		Value:    5,
		Category: "CONSTRAINTS",
		Usage:    "minimum Levenshtein distance `NUM`",
	},
	&cli.BoolFlag{
		Name:     "progress",
		Aliases:  []string{"p"},
		Value:    false,
		Category: "OUTPUT",
		Hidden:   true,
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
		Usage:    "filename to save scan output `FILE`",
	},
	&cli.StringFlag{
		Name:     "format",
		Aliases:  []string{"f"},
		Value:    "list",
		Category: "OUTPUT",
		Usage:    "output format: (csv,tsv,table,list,html,md,json) `FORMAT`",
	},
	&cli.StringFlag{
		Name:     "nameservers",
		Aliases:  []string{"n"},
		Value:    "",
		Hidden:   true,
		Category: "PERFORMANCE",
		Usage:    "DNS or DoH servers to query (separated with commas) `[NAMES..]`",
	},
	&cli.BoolFlag{
		Name:     "registered",
		Aliases:  []string{"r"},
		Value:    false,
		Category: "OUTPUT",
		Usage:    "show only registered domain names",
	},
	&cli.BoolFlag{
		Name:     "unregistered",
		Aliases:  []string{"u"},
		Value:    false,
		Category: "OUTPUT",
		Usage:    "show only unregistered domain names",
	},
	&cli.BoolFlag{
		Name:     "summary",
		Aliases:  []string{"s"},
		Value:    true,
		Hidden:   true,
		Category: "OUTPUT",
		Usage:    "show summary of scan results",
	},
	&cli.PathFlag{
		Name:     "dir",
		Value:    "domains",
		Category: "OUTPUT",
		Usage:    "directory name to save files `DIR`",
		Action: func(ctx *cli.Context, v string) error {
			return nil
		},
	},
	&cli.BoolFlag{
		Name:     "rua",
		Value:    false,
		Hidden:   true,
		Category: "PERFORMANCE",
		Usage:    "use random user agent for HTTP requests",
	},
}

var TypoCmd = cli.Command{
	Name:                   "typo",
	Aliases:                []string{"t"},
	Usage:                  "Generate domain variations and collect information on them",
	Description:            "URLInsane is designed to detect domain typosquatting by using advanced algorithms, information-gathering techniques, and data analysis to identify potentially harmful variations of targeted domains that cybercriminals might exploit. This tool is essential for defending against threats like typosquatting, brandjacking, URL hijacking, fraud, phishing, and corporate espionage. By detecting malicious domain variations, it provides an added layer of protection to brand integrity and user trust. Additionally, URLInsane enhances threat intelligence capabilities, strengthening proactive cybersecurity measures.",
	UsageText:              "urlinsane [g opts..] typo [opts..] [domain]",
	UseShortOptionHandling: true,
	Flags:                  Flags,
	Action: func(cCtx *cli.Context) error {
		if cCtx.NArg() == 0 {
			fmt.Println(text.FgRed.Sprint("\n  a domain name is needed!\n"))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)

		}
		if cCtx.NArg() > 1 {
			fmt.Println(text.FgRed.Sprint("\n  only one domain name at at time!\n"))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}

		cfg, err := config.New(config.CliOptions(cCtx))
		if err != nil {
			log.Error(err)
			fmt.Println(text.FgRed.Sprint(err))
			cli.ShowSubcommandHelpAndExit(cCtx, 1)
		}
		return engine.New(cfg).Execute()
	},
	CustomHelpTemplate: ShowSubcommandHelp(cli.SubcommandHelpTemplate),
}

func init() {

}

func ShowSubcommandHelp(template string) string {
	collectors := CollectorTable()
	algorithms := AlgorithmTable()
	languages := LanguageTable()
	keyboards := KeyboardTable()
	outputs := OutputTable()

	return fmt.Sprintf(`%sKEYBOARDS:
%s

			eg: urlinsane typo -k en1,en2,en3,en4 example.com

LANGUAGES:
%s

			eg: urlinsane typo -l ru,hy,en example.com

ALGORITHMS:
%s

			eg: urlinsane typo -a cs,gr,cm example.com

COLLECTORS:
%s

			eg: urlinsane typo -c ip,idn example.com

OUTPUTS:
%s

			eg: urlinsane typo -f table example.com

EXAMPLE:

    urlinsane typo example.com
    urlinsane typo -a co example.com
    urlinsane typo -a co,oi,oy -c ip,idna,ns example.com
    urlinsane typo -l fr,en -k en1,en2 example.com

AUTHOR:
   Rangertaha (rangertaha@gmail.com)
     
     `, template, keyboards, languages, algorithms, collectors, outputs)
}

func CollectorTable() string {
	t := table.NewWriter()
	t.SetStyle(pkg.StyleClear)
	t.AppendHeader(table.Row{"  ", "ID", "Description"})
	for _, p := range collectors.List() {
		t.AppendRow([]interface{}{"  ", p.Id(), p.Description()})
	}
	return t.Render()
}

// func CollectorFields() (fields string) {
// 	headers := []string{}
// 	for _, i := range collectors.List() {
// 		for _, header := range i.Headers() {
// 			headers = append(headers, strings.ToLower(header))
// 		}
// 	}
// 	return strings.Join(headers, ",")
// }

func AlgorithmTable() string {
	t := table.NewWriter()
	t.SetStyle(pkg.StyleClear)
	t.AppendHeader(table.Row{"  ", "ID", "Name"})
	for _, p := range algorithms.List() {
		t.AppendRow([]interface{}{"  ", p.Id(), p.Name()})
	}
	return t.Render()
}

func LanguageTable() string {
	t := table.NewWriter()
	t.SetStyle(pkg.StyleClear)
	t.AppendHeader(table.Row{"  ", "ID", "Name", "Glyphs", "Homophones",
		"Antonyms", "Typos", "Cardinal", "Ordinal", "Stems"})
	for _, p := range languages.Languages() {
		t.AppendRow([]interface{}{"  ", p.Id(), p.Name(), len(p.Homoglyphs()),
			len(p.Homophones()), len(p.Antonyms()), len(p.Misspellings()),
			len(p.Cardinal()), len(p.Ordinal()), 0})
	}
	return t.Render()
}

func KeyboardTable() string {
	t := table.NewWriter()
	t.SetStyle(pkg.StyleClear)
	rows := []table.Row{}
	for _, lang := range languages.Languages() {
		row := table.Row{" "}
		row = append(row, strings.ToUpper(lang.Name()))
		for _, board := range lang.Keyboards() {
			row = append(row, fmt.Sprintf("%s: %s", board.Id(), board.Name()))
		}
		rows = append(rows, row)
	}
	t.AppendHeader(table.Row{" ", "LANGUAGE", "ID:NAME..."})
	for _, row := range rows {
		t.AppendRow(row)
	}
	return t.Render()
}

func OutputTable() string {
	t := table.NewWriter()
	t.SetStyle(pkg.StyleClear)
	t.AppendHeader(table.Row{"  ", "ID", "Name"})
	for _, p := range outputs.List() {
		t.AppendRow([]interface{}{"  ", p.Id(), p.Description()})
	}
	return t.Render()
}
