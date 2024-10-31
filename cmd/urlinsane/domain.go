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
	"fmt"
	"os"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/config"
	"github.com/rangertaha/urlinsane/internal/engine"
	_ "github.com/rangertaha/urlinsane/internal/plugins/algorithms/all"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/spf13/cobra"
)

// const templateBase = `USAGE:{{if .Runnable}}
//   {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
//   {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

// ALIASES:
//   {{.NameAndAliases}}{{end}}{{if .HasExample}}

// EXAMPLES:
// {{.Example}}{{end}}{{if .HasAvailableSubCommands}}
// Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
//   {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

// OPTIONS:
// {{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

// GLOBAL OPTIONS:
// {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
// Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
//   {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

// Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`

// const helpTemplate = `

// ALGORITHMS:
//     Typosquatting algorithm plugins that generate typos.

//     ID | Name    | Description
//     -----------------------------------------{{range .Algorithms}}
//     {{.Id}}{{"\t"}}{{.Name}}{{"\t"}}{{.Description}}{{end}}

// INFORMATION:
//     Information-gathering plugins that collect information on each typo

//     ID | Name    | Description
//     -----------------------------------------{{range .Information}}
//     {{.Id}}{{"\t"}}{{.Name}}{{"\t"}}{{.Description}}{{end}}

// LANGUAGES:
//     ID | Name    | Description
//     -----------------------------------------{{range .Languages}}
//     {{.Id}}{{"\t"}}{{.Name}}{{"\t"}}{{.Description}}{{end}}

// KEYBOARDS:
//     ID | Name    | Description
//     -----------------------------------------{{range .Keyboards}}
//     {{.Id}}{{"\t"}}{{.Name}}{{"\t"}}{{.Description}}{{end}}

// EXAMPLE:

//     urlinsane google.com
//     urlinsane -t co google.com
//     urlinsane -t co,oi,oy -i ip,idna,ns google.com
//     urlinsane -l fr,en -k en1,en2 google.com

// AUTHOR:
//    Rangertaha (rangertaha@gmail.com)

// `

// type HelpOptions struct {
// 	Languages   []internal.Language
// 	Keyboards   []internal.Keyboard
// 	Algorithms  []internal.Algorithm
// 	Information []internal.Information
// }

// var cliOptions bytes.Buffer

// rootCmd represents the typo command
var domainCmd = &cobra.Command{
	Use:   "domain [flags] [name]",
	Short: "Generates and detects possible typosquatting domain names and arbitrary names",
	Long: `Urlinsane is used to perform or detect typosquatting, brandjacking, URL hijacking, 
	fraud, phishing attacks, corporate espionage, and threat intelligence.

Urlinsane is built around linguistic modeling, natural language processing, 
information gathering, and analysis. It's easily extensible with plugins for typo algorithms, 
Information gathering and analysis. Its linguistic models also allow us to easily add new 
languages and keyboard layouts. Currently, it supports 9 languages, 19 keyboard layouts, 
24 algorithms, 8 information gathering, and 2 analysis modules.

`,
	// PersistentPreRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
	// },
	//   PreRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
	//   },
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		config, err := config.CobraConfig(cmd, args)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(0)
		}
		fmt.Print(internal.LOGO)
		t := engine.NewDomainTypos(config)
		t.Execute()

	},
	//   PostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
	//   },
	//   PersistentPostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
	//   },
}

func init() {
	// // fmt.Println(languages.Languages())
	// helpOptions := HelpOptions{
	// 	languages.Languages(),
	// 	languages.Keyboards(),
	// 	algorithms.List(),
	// 	information.List(),
	// }

	// // Create a new template and parse into it.
	// tmpl := template.Must(template.New("help").Parse(helpTemplate))

	// // Run the template to verify the output.
	// err := tmpl.Execute(&cliOptions, helpOptions)
	// if err != nil {
	// 	fmt.Printf("Execution: %s", err)
	// }

	// DomainCmd.SetUsageTemplate(templateBase + cliOptions.String())
	// rootCmd.SetVersionTemplate(internal.VERSION)
	// DomainCmd.CompletionOptions.DisableDefaultCmd = true

	// Plugins
	domainCmd.PersistentFlags().StringArrayP("typos", "t", []string{"all"}, "IDs of typo algorithms to use for generating typos")
	domainCmd.PersistentFlags().StringArrayP("info", "i", []string{"all"}, "IDs of info gathering functions to apply")

}
