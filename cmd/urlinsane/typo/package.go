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
	"text/template"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/config"

	_ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
	"github.com/rangertaha/urlinsane/internal/plugins/information/packages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/rangertaha/urlinsane/internal/urlinsane"
	"github.com/rangertaha/urlinsane/internal/utils"
	"github.com/spf13/cobra"
)


const pkgHelpTemplate = `

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

    urlinsane typo pkg example


AUTHOR:
   Rangertaha (rangertaha@gmail.com)

`

var pkgCliOptions bytes.Buffer

// rootCmd represents the typo command
var pkgCmd = &cobra.Command{
	Use:     "pkg [flags] [name]",
	Aliases: []string{"package", "module", "p"},
	Short:   "Detects software package typosquatting",
	Long:    `Detects software package typosquatting`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) == 0 {
		// 	cmd.Help()
		// }

		config, err := config.CobraConfig(cmd, args, internal.PACKAGE)
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
	TypoCmd.AddCommand(pkgCmd)
	pkgHelpOptions := HelpOptions{
		LanguageTable(),
		KeyboardTable(),
		AlgorithmTable(),
		PackageInformationTable(),
	}

	// Create a new template and parse into it.
	tmpl := template.Must(template.New("help").Parse(pkgHelpTemplate))

	// Run the template to verify the output.
	err := tmpl.Execute(&pkgCliOptions, pkgHelpOptions)
	if err != nil {
		fmt.Printf("Execution: %s", err)
	}

	pkgCmd.SetUsageTemplate(templateBase + pkgCliOptions.String())
	pkgCmd.CompletionOptions.DisableDefaultCmd = true

	// Plugins
	pkgCmd.Flags().StringP("info", "i", "all", "Information plugin IDs to apply")

	// Filtering
	pkgCmd.Flags().Bool("all", false, "Scan all generated variants equivalent to: --ld 100")
	pkgCmd.Flags().Bool("show", false, "Show all generated variants")
	pkgCmd.Flags().Int("ld", 3, "Minimum levenshtein distance to scan")

}

func PackageInformationTable() string {
	t := table.NewWriter()
	t.SetStyle(utils.StyleClear)
	t.AppendHeader(table.Row{"  ", "ID", "Description"})
	for _, p := range packages.List() {
		t.AppendRow([]interface{}{"  ", p.Id(), p.Description()})
	}
	return t.Render()
}
