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

	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/config"
	"github.com/rangertaha/urlinsane/engine"
	"github.com/rangertaha/urlinsane/plugins/languages"
	_ "github.com/rangertaha/urlinsane/plugins/languages/all"

	"github.com/rangertaha/urlinsane/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/plugins/algorithms/all"

	"github.com/rangertaha/urlinsane/plugins/information"
	_ "github.com/rangertaha/urlinsane/plugins/information/all"
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

ALGORITHMS:
    Typosquating algorithms that generate domain names.

    ID | Description
    -----------------------------------------{{range .Algorithms}}
    {{.Code}}	{{.Name}} {{.Description}}{{end}}

INFORMATION:
    Information gathering functions that collects information on each domain name

    ID | Description
    -----------------------------------------{{range .Information}}
    {{.Code}}	{{.Name}} {{.Description}}{{end}}


LANGUAGES:
    ID | Name    | Description
    -----------------------------------------{{range .Languages}}
    {{.Code}}  {{.Name}}   {{end}}

KEYBOARDS:
    ID | Name     | Description
    -----------------------------------------{{range .Languages}}{{range .Keyboards}}
    {{.Code}}  {{.Name}}    {{.Description}}{{end}}{{end}}



EXAMPLE:

    urlinsane google.com
    urlinsane -t co google.com 
    urlinsane -t co,oi,oy -i ip,idna,ns google.com
    urlinsane -l fr,en -k en1,en2 google.com

AUTHOR:
    Tal Hatchi (Rangertaha)

`

const hTemplate = `
{{if .Algorithms}}
TYPOS: 
  These are the types of typo/error algorithms that generate the domain variants{{range .Algorithms}}
    {{.Code}}	 {{.Name}}	{{.Description}}{{end}}
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
	Languages   []urlinsane.Language
	Algorithms  []urlinsane.Module
	Information []urlinsane.Module
}

var cliOptions bytes.Buffer

// rootCmd represents the typo command
var rootCmd = &cobra.Command{
	Use:   "urlinsane [flags] [domains]",
	Short: "Generates and detects possible typosquatting domains",
	Long: `Urlinsane is used to perform or detect typosquatting, brandjacking,
URL hijacking, fraud, phishing attacks, corporate espionage and threat intelligence.

Urlinsane is built around linguistic modeling, natural language processing, 
information gathering and analysis. It's easily extensible with plugins for typo algorithms, 
inforamtion gathering and analysis. Its linguistic models also allow it us to easily add new 
languages and keyboard layouts. Currently it supports 9 languages, 19 keyboard layouts, 
24 algorithms, 8 information gathering, and 2 analysis modules.

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		c, err := config.CobraConfig(cmd, args)
		if err != nil {
			fmt.Println(err)
			cmd.Help()
			os.Exit(0)
		}

		engine.New(c).Execute()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func init() {
	// fmt.Println(languages.Languages())
	helpOptions := HelpOptions{
		languages.Languages(),
		algorithms.List(),
		information.List(),
	}

	// Create a new template and parse the letter into it.
	tmpl := template.Must(template.New("help").Parse(helpTemplate))

	// Run the template to verify the output.
	err := tmpl.Execute(&cliOptions, helpOptions)
	if err != nil {
		fmt.Printf("Execution: %s", err)
	}

	rootCmd.SetUsageTemplate(templateBase + cliOptions.String())

	// Options
	rootCmd.PersistentFlags().StringArrayP("languages", "l", []string{"en"}, "IDs of languages to use for lingustic algorithms")
	rootCmd.PersistentFlags().StringArrayP("keyboards", "k", []string{"all"}, "IDs of keyboard layouts to use of the given languages")

	// Plugins
	rootCmd.PersistentFlags().StringArrayP("typos", "t", []string{"all"}, "IDs of typo algorithms to use for generating domains")
	rootCmd.PersistentFlags().StringArrayP("info", "i", []string{"all"}, "IDs of info gathering functions to apply to each domain")

	// Processing
	rootCmd.PersistentFlags().BoolP("progress", "p", false, "Show progress bar")
	rootCmd.PersistentFlags().Bool("no-cache", true, "Prevents caching of results")
	rootCmd.PersistentFlags().Bool("online", false, "Only show domains that are online")

	// Timing
	rootCmd.PersistentFlags().IntP("concurrency", "c", 50, "Number of concurrent workers")
	rootCmd.PersistentFlags().Int("delay", 10, "A delay between network calls")

	// Outputs
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show more details and remove truncated columns")
	rootCmd.PersistentFlags().StringP("file", "f", "", "Output filename defaults to stdout")
	rootCmd.PersistentFlags().StringP("format", "o", "text", "Output format (csv, text, json)")
}
