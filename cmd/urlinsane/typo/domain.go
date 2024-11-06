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
	"os"
	"strings"
	"text/template"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/config"

	_ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
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

    urlinsane typo domain google.com
    urlinsane typo domain  -a co google.com
    urlinsane typo domain  -t co,oi,oy -i ip,idna,ns google.com
    urlinsane typo domain -l fr,en -k en1,en2 google.com

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
	Use:   "domain [flags] [name]",
	Short: "Detects domain typosquatting",
	Long:  `Detects domain typosquatting`,
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
	TypoCmd.AddCommand(domainCmd)
	domainHelpOptions := DomainHelpOptions{
		LanguageTable(),
		KeyboardTable(),
		AlgorithmTable(),
		DomainInformationTable(),
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

	// Filtering
	domainCmd.Flags().Bool("all", false, "Scan all generated variants equivalent to: --ld 100")
	domainCmd.Flags().Bool("show", false, "Show all generated variants")
	domainCmd.Flags().Int("ld", 3, "Minimum levenshtein distance to scan")
	domainCmd.Flags().String("filter",  DomainInformationFields(), "Output filter options")

}

func DomainInformationTable() string {
	t := table.NewWriter()
	t.SetStyle(utils.StyleClear)
	t.AppendHeader(table.Row{"  ", "ID", "Description"})
	for _, p := range domains.List() {
		t.AppendRow([]interface{}{"  ", p.Id(), p.Description()})
	}
	return t.Render()
}

func DomainInformationFields() (fields string) {
	headers := []string{}
	for _, i := range domains.List() {
		for _, header := range i.Headers(){
			headers = append(headers, strings.ToLower(header))
		}
	}
	return strings.Join(headers, ",")
}
