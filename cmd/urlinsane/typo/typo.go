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
package typo

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/internal/plugins/algorithms/all"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/rangertaha/urlinsane/internal/utils"
	"github.com/spf13/cobra"
)

const templateBase = `USAGE:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

ALIASES:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

EXAMPLES:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

OPTIONS:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

GLOBAL OPTIONS:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`

type HelpOptions struct {
	Languages   string
	Keyboards   string
	Algorithms  string
	Information string
}

var cliOptions bytes.Buffer

// rootCmd represents the typo command
var TypoCmd = &cobra.Command{
	Use:   "typo [flags] [name]",
	Short: "Detects typosquatting across domains, emails, usernames, and software packages",
	Long: `Designed to detect typosquatting across domains, arbitrary names, usernames, and software packages. 
By leveraging advanced algorithms, information-gathering techniques, and data analysis, it identifies 
potentially harmful variations of targeted entities that cybercriminals might exploit. Essential for 
defending against typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, and corporate 
espionage, URLInsane also enhances threat intelligence capabilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func init() {
	// rootCmd.AddCommand(typoCmd)
}

func AlgorithmTable() string {
	t := table.NewWriter()
	t.SetStyle(utils.StyleClear)
	t.AppendHeader(table.Row{"  ", "ID", "Name"})
	for _, p := range algorithms.List() {
		t.AppendRow([]interface{}{"  ", p.Id(), p.Name()})
	}
	return t.Render()
}

func LanguageTable() string {
	t := table.NewWriter()
	t.SetStyle(utils.StyleClear)
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
	t.SetStyle(utils.StyleClear)
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
