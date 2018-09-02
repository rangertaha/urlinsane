// Copyright © 2018 rangertaha rangertaha@gmail.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"os"
	"fmt"
	"bytes"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/languages"
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
	Typos []urlinsane.Typo
	Funcs []urlinsane.Extra
}

var cliOptions bytes.Buffer

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "urlinsane [domains]",
	Short: "Generates domain typos and variations",
	Long: `Multilingual domain typo permutation engine used to perform or detect typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, corporate espionage and threat intelligence.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Create config from cli options/arguments
		config := urlinsane.CobraConfig(cmd, args)

		// Create a new instance of urlinsane
		urli := urlinsane.New(config)

		// Start generating results
		urli.Start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	helpOptions := HelpOptions{
		languages.GetKeyboards([]string{"all"}),
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

	rootCmd.SetUsageTemplate(TEMPLATE_BASE + cliOptions.String())

	// Basic options
	rootCmd.PersistentFlags().StringArrayP("keyboards", "k", []string{"en1"},
		"Keyboards/layouts ID to use")
	//rootCmd.PersistentFlags().StringArrayP("languages", "l", []string{"all"},
	//	"Language ID to use for linguistic typos")

	// Processing
	rootCmd.PersistentFlags().IntP("concurrency", "c", 50,
		"Number of concurrent workers")
	rootCmd.PersistentFlags().StringArrayP("typos", "t", []string{"all"},
		"Types of typos to perform")

	// Post Processing options for retrieving additional data
	rootCmd.PersistentFlags().StringArrayP("funcs", "x", []string{"idna"},
		"Extra functions for data or filtering")

	// Output options
	rootCmd.PersistentFlags().StringP("file", "f", "", "Output filename")
	rootCmd.PersistentFlags().StringP("format", "o", "text", "Output format (csv, text)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Output additional details")
}

