// Copyright (C) 2024  Tal Hatchi (Rangertaha)
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
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/rangertaha/urlinsane/languages"
	"github.com/rangertaha/urlinsane/typo"
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

const helpTemplate = `
{{if .Typos}}
TYPOS: 
  These are the types of typo/error algorithms that generate the domain variants{{range .Typos}}
    {{.Code}}	{{.Description}}{{end}}
    ALL	Apply all typosquatting algorithms
{{end}}{{if .Funcs}}
INFORMATION: 
  Retrieve aditional information on each domain variant.{{range .Funcs}}
    {{.Code}}    {{.Description}}{{end}}
    ALL    Apply all post typosquating functions
{{end}}{{if .Filters}}
FILTERS: 
  Filters to reduce the number domain variants returned.{{range .Filters}}
    {{.Code}}   {{.Description}}{{end}}
    ALL    Apply all filters
{{end}}{{if .Keyboards}}
KEYBOARDS:{{range .Keyboards}}
    {{.Code}}	{{.Description}}{{end}}
    ALL	Use all keyboards
{{end}}
EXAMPLE:

    urlinsane google.com
    urlinsane -a co google.com 
    urlinsane -a co -x ip,idna,ns google.com 

AUTHOR:
    Tal Hatchi (Rangertaha)

`

type HelpOptions struct {
	Keyboards []languages.Keyboard
	Typos     []typo.Module
	Funcs     []typo.Module
	Filters   []typo.Module
}

var cliOptions bytes.Buffer

// rootCmd represents the typo command
var rootCmd = &cobra.Command{
	Use:   "urlinsane [flags] [domains]",
	Short: "Generates and detects possible squatting domains",
	Long: `
Multilingual domain typo permutation engine used to perform or detect typosquatting, brandjacking,
URL hijacking, fraud, phishing attacks, corporate espionage and threat intelligence.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		// // Create config from cli options/arguments
		// config := typo.CobraConfig(cmd, args)

		// // Create a new instance of urlinsane
		// typosquating := typo.New(config)

		// // Start generating results
		// typosquating.Execute()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println(err)
		// os.Exit(1)
	}
}
func init() {
	helpOptions := HelpOptions{
		languages.Keyboards(),
		typo.Typos.Get("all"),
		typo.Extras.Get("all"),
		typo.Filters.Get("all"),
	}

	// Create a new template and parse the letter into it.
	tmpl := template.Must(template.New("help").Parse(helpTemplate))

	// Run the template to verify the output.
	err := tmpl.Execute(&cliOptions, helpOptions)
	if err != nil {
		fmt.Printf("Execution: %s", err)
	}

	rootCmd.SetUsageTemplate(templateBase + cliOptions.String())

	// Basic options
	rootCmd.PersistentFlags().StringArrayP("keyboards", "k", []string{"en"},
		"Keyboards/layouts ID to use")
	// viper.BindPFlag("keyboards", rootCmd.PersistentFlags().Lookup("keyboards"))

	// Processing
	rootCmd.PersistentFlags().IntP("concurrency", "c", 50,
		"Number of concurrent workers")
	// viper.BindPFlag("concurrency", rootCmd.PersistentFlags().Lookup("concurrency"))

	rootCmd.PersistentFlags().StringArrayP("typos", "t", []string{"all"},
		"Types of typos to perform")
	// viper.BindPFlag("typos", rootCmd.PersistentFlags().Lookup("typos"))

	// Post Processing options for retrieving additional data
	rootCmd.PersistentFlags().StringArrayP("funcs", "x", []string{"ld", "idna"},
		"Extra functions or filters")
	// viper.BindPFlag("funcs", rootCmd.PersistentFlags().Lookup("funcs"))

	rootCmd.PersistentFlags().StringArrayP("filters", "r", []string{""},
		"Filter results to reduce the number of results")
	// viper.BindPFlag("filters", rootCmd.PersistentFlags().Lookup("filters"))

	rootCmd.PersistentFlags().Int64("delay", 10,
		"A delay between network calls")

	rootCmd.PersistentFlags().Int64("random-delay", 5,
		"Used to randomize the delay between network calls.")

	// Output options
	rootCmd.PersistentFlags().StringP("file", "f", "", "Output filename")
	rootCmd.PersistentFlags().StringP("format", "o", "text", "Output format (csv, text)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Output additional details")
	// viper.BindPFlag("filverboseters", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.AddCommand(rootCmd)
}
