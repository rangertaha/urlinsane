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

// import (
// 	"bytes"
// 	"fmt"
// 	"os"
// 	"text/template"

// 	"github.com/jedib0t/go-pretty/v6/table"
// 	"github.com/rangertaha/urlinsane/internal"
// 	"github.com/rangertaha/urlinsane/internal/config"

// 	_ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
// 	"github.com/rangertaha/urlinsane/internal/plugins/information/emails"
// 	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
// 	"github.com/rangertaha/urlinsane/internal/urlinsane"
// 	"github.com/rangertaha/urlinsane/internal/utils"
// 	"github.com/spf13/cobra"
// )


// const emailHelpTemplate = `

// ALGORITHMS:
//     Typosquatting algorithm plugins that generate typos.

// {{.Algorithms}}


// INFORMATION:
//     Information-gathering plugins that collect information on each domain

// {{.Information}}


// LANGUAGES:

// {{.Languages}}


// KEYBOARDS:

// {{.Keyboards}}


// EXAMPLE:

//     urlinsane typo email username@example.com


// AUTHOR:
//    Rangertaha (rangertaha@gmail.com)

// `

// var emailCliOptions bytes.Buffer

// // rootCmd represents the typo command
// var emailCmd = &cobra.Command{
// 	Use:   "email [flags] [name]",
// 	Aliases: []string{"mail", "e"},
// 	Short: "Detects email typosquatting",
// 	Long:  `Detects email typosquatting`,
// 	Args:  cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// if len(args) == 0 {
// 		// 	cmd.Help()
// 		// }

// 		config, err := config.CobraConfig(cmd, args, internal.EMAIL)
// 		if err != nil {
// 			fmt.Printf("%s", err)
// 			os.Exit(0)
// 		}
// 		config.Type()

// 		t := urlinsane.New(config)
// 		t.Execute()

// 	},
// }

// func init() {
// 	// TypoCmd.AddCommand(emailCmd)
// 	emailHelpOptions := HelpOptions{
// 		LanguageTable(),
// 		KeyboardTable(),
// 		AlgorithmTable(),
// 		EmailInformationTable(),
// 	}

// 	// Create a new template and parse into it.
// 	tmpl := template.Must(template.New("help").Parse(emailHelpTemplate))

// 	// Run the template to verify the output.
// 	err := tmpl.Execute(&emailCliOptions, emailHelpOptions)
// 	if err != nil {
// 		fmt.Printf("Execution: %s", err)
// 	}

// 	emailCmd.SetUsageTemplate(templateBase + emailCliOptions.String())
// 	emailCmd.CompletionOptions.DisableDefaultCmd = true

// 	// Plugins
// 	emailCmd.Flags().StringP("info", "i", "all", "Information plugin IDs to apply")

// 	// Filtering
// 	emailCmd.Flags().Bool("all", false, "Scan all generated variants equivalent to: --ld 100")
// 	emailCmd.Flags().Bool("show", false, "Show all generated variants")
// 	emailCmd.Flags().Int("ld", 3, "Minimum levenshtein distance to scan")

// }

// func EmailInformationTable() string {
// 	t := table.NewWriter()
// 	t.SetStyle(utils.StyleClear)
// 	t.AppendHeader(table.Row{"  ", "ID", "Description"})
// 	for _, p := range emails.List() {
// 		t.AppendRow([]interface{}{"  ", p.Id(), p.Description()})
// 	}
// 	return t.Render()
// }
