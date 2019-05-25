// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cybersectech-org/urlinsane"
	"github.com/cybersectech-org/urlinsane/languages"
	"github.com/spf13/cobra"
)

const TEMPLATE_BASE = `USAGE:{{if .Runnable}}
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

const TEMPLATE = `
{{if .Keyboards}}
KEYBOARDS:{{range .Keyboards}}
  {{.Code}}	{{.Description}}{{end}}
  ALL	Use all keyboards
{{end}}{{if .Typos}}
TYPOS: These are the types of typo/error algorithms that generate the domain variants{{range .Typos}}
  {{.Code}}	{{.Description}}{{end}}
  ALL   Apply all typosquatting algorithms
{{end}}{{if .Funcs}}
FUNCTIONS: Post processig functions that retieve aditional information on each domain variant.{{range .Funcs}}
  {{.Code}}	{{.Description}}{{end}}
  ALL  	Apply all post typosquating functions
{{end}}
EXAMPLE:

    urlinsane google.com
    urlinsane google.com -t co
    urlinsane google.com -t co -x ip -x idna -x ns

AUTHOR:
  Written by Rangertaha <rangertaha@gmail.com>

`

type HelpOptions struct {
	Keyboards []languages.Keyboard
	Typos     []urlinsane.Typo
	Funcs     []urlinsane.Extra
}

var cliOptions bytes.Buffer

// typoCmd represents the typo command
var typoCmd = &cobra.Command{
	Use:   "typo [domains]",
	Short: "Generates domain typos and variations",
	Long:  `Multilingual domain typo permutation engine used to perform or detect typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage and threat intelligence.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create config from cli options/arguments
		config := urlinsane.CobraConfig(cmd, args)

		// Create a new instance of urlinsane
		urli := urlinsane.New(config)

		// Start generating results
		urli.Start()
	},
}

func init() {
	helpOptions := HelpOptions{
		languages.KEYBOARDS.Keyboards("all"),
		urlinsane.TRetrieve("all"),
		urlinsane.FRetrieve("all"),
	}

	// Create a new template and parse the letter into it.
	tmpl := template.Must(template.New("help").Parse(TEMPLATE))

	// Run the template to verify the output.
	err := tmpl.Execute(&cliOptions, helpOptions)
	if err != nil {
		fmt.Printf("Execution: %s", err)
	}

	typoCmd.SetUsageTemplate(TEMPLATE_BASE + cliOptions.String())

	// Basic options
	typoCmd.PersistentFlags().StringArrayP("keyboards", "k", []string{"en"},
		"Keyboards/layouts ID to use")
	//rootCmd.PersistentFlags().StringArrayP("languages", "l", []string{"all"},
	//	"Language ID to use for linguistic typos")

	// Processing
	typoCmd.PersistentFlags().IntP("concurrency", "c", 50,
		"Number of concurrent workers")
	typoCmd.PersistentFlags().StringArrayP("typos", "t", []string{"all"},
		"Types of typos to perform")

	// Post Processing options for retrieving additional data
	typoCmd.PersistentFlags().StringArrayP("funcs", "x", []string{"idna"},
		"Extra functions or filters")

	// Output options
	typoCmd.PersistentFlags().StringP("file", "f", "", "Output filename")
	typoCmd.PersistentFlags().StringP("format", "o", "text", "Output format (csv, text)")
	typoCmd.PersistentFlags().BoolP("verbose", "v", false, "Output additional details")
	rootCmd.AddCommand(typoCmd)
}
