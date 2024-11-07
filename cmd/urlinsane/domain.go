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
package urlinsane

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/config"

	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/rangertaha/urlinsane/internal/urlinsane"
	"github.com/rangertaha/urlinsane/internal/utils"
	"github.com/spf13/cobra"
)

const domainHelpTemplate = `

ALGORITHMS:
    Typosquatting algorithm plugins that generate typos.

{{.Algorithms}}


INFORMATION:
    Information-gathering plugins that collect information on each domain

{{.Information}}


LANGUAGES:

{{.Languages}}


KEYBOARDS:

{{.Keyboards}}


EXAMPLE:

    urlinsane typo example.com
    urlinsane typo -a co example.com
    urlinsane typo -t co,oi,oy -i ip,idna,ns example.com
    urlinsane typo -l fr,en -k en1,en2 example.com

AUTHOR:
   Rangertaha (rangertaha@gmail.com)

`

type DomainHelpOptions struct {
	Languages   string
	Keyboards   string
	Algorithms  string
	Information string
}

var domainCliOptions bytes.Buffer

// rootCmd represents the typo command
var domainCmd = &cobra.Command{
	Use:   "typo [flags] [name]",
	Short: "Detects potential typosquatting domains by generating and checking misspelled variations of a given domain name.",
	Long:  `URLInsane is designed to detect domain typosquatting by using advanced algorithms, information-gathering 
  techniques, and data analysis to identify potentially harmful variations of targeted domains that cybercriminals 
  might exploit. This tool is essential for defending against threats like typosquatting, brandjacking, URL hijacking, 
  fraud, phishing, and corporate espionage. By detecting malicious domain variations, it provides an added layer of 
  protection to brand integrity and user trust. Additionally, URLInsane enhances threat intelligence capabilities, 
  strengthening proactive cybersecurity measures.
	
`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) == 0 {
		// 	cmd.Help()
		// }

		config, err := config.CobraConfig(cmd, args, internal.DOMAIN)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(0)
		}
		config.Type()

		t := urlinsane.New(config)
		t.Execute()

	},
}

func init() {
	domainHelpOptions := DomainHelpOptions{
		LanguageTable(),
		KeyboardTable(),
		AlgorithmTable(),
		InformationTable(),
	}

	// Create a new template and parse into it.
	tmpl := template.Must(template.New("help").Parse(domainHelpTemplate))

	// Run the template to verify the output.
	err := tmpl.Execute(&domainCliOptions, domainHelpOptions)
	if err != nil {
		fmt.Printf("Execution: %s", err)
	}

	domainCmd.SetUsageTemplate(templateBase + domainCliOptions.String())
	domainCmd.CompletionOptions.DisableDefaultCmd = true

	// Plugins
	domainCmd.Flags().StringP("info", "i", "all", "Information plugin IDs to apply")
	domainCmd.PersistentFlags().Bool("image", false, "Take screenshot of domain saved to .urlinsane/domains/")

	// Filtering
	domainCmd.Flags().Bool("all", false, "Scan all generated variants equivalent to: --ld 100")
	domainCmd.Flags().Bool("show", false, "Show all generated variants")
	domainCmd.Flags().Int("ld", 3, "Minimum levenshtein distance to scan")
	domainCmd.Flags().String("filter", InformationFields(), "Output filter options")

}

func InformationTable() string {
	t := table.NewWriter()
	t.SetStyle(utils.StyleClear)
	t.AppendHeader(table.Row{"  ", "ID", "Description"})
	for _, p := range domains.List() {
		t.AppendRow([]interface{}{"  ", p.Id(), p.Description()})
	}
	return t.Render()
}

func InformationFields() (fields string) {
	headers := []string{}
	for _, i := range domains.List() {
		for _, header := range i.Headers() {
			headers = append(headers, strings.ToLower(header))
		}
	}
	return strings.Join(headers, ",")
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
