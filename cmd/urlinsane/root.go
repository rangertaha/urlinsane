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

	// _ "github.com/rangertaha/urlinsane/internal/plugins/algorithms/all"
	// _ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
	// _ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/rangertaha/urlinsane/internal/config"
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
var Config config.Config

// rootCmd represents the typo command
var rootCmd = &cobra.Command{
	Use:   "urlinsane [flags] [name]",
	Short: "Urlinsane is an advanced cybersecurity typosquatting tool",
	Long: `URLInsane is a powerful command-line tool crafted to identify typo-squatting across domains, 
  usernames, and software packages. Utilizing advanced algorithms, information gathering, and data analysis, 
  it uncovers potentially harmful variations of a target's named entity that cybercriminals could exploit. 
  This tool is essential for combating typosquatting, brandjacking, URL hijacking, fraud, phishing attacks, 
  corporate espionage, and enhancing threat intelligence.

`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	//   PreRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
	//   },
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		// if len(args) == 0 {
		// 	cmd.Help()
		// 	os.Exit(0)
		// }

		// config, err := config.CobraConfig(cmd, args)
		// if err != nil {
		// 	fmt.Printf("%s", err)
		// 	os.Exit(0)
		// }
		// fmt.Print(internal.LOGO)
		// t := engine.NewDomainTypos(config)
		// t.Execute()

	},
	//   PostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
	//   },
	//   PersistentPostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
	//   },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func init() {

	// // fmt.Println(languages.Languages())
	// helpOptions := HelpOptions{
	// 	languages.Languages(),
	// 	languages.Keyboards(),
	// 	algorithms.List(),
	// 	information.List(),
	// }

	// // Create a new template and parse the letter into it.
	// tmpl := template.Must(template.New("help").Parse(helpTemplate))

	// // Run the template to verify the output.
	// err := tmpl.Execute(&cliOptions, helpOptions)
	// if err != nil {
	// 	fmt.Printf("Execution: %s", err)
	// }

	// rootCmd.SetUsageTemplate(templateBase + cliOptions.String())
	// rootCmd.SetVersionTemplate(internal.VERSION)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Targets
	// rootCmd.PersistentFlags().Bool("lib", false, "IDs of languages to use for linguistic algorithms")
	// rootCmd.PersistentFlags().Bool("email", false, "IDs of languages to use for linguistic algorithms")
	// rootCmd.PersistentFlags().Bool("username", false, "IDs of keyboard layouts to use of the given languages")
	// rootCmd.PersistentFlags().Bool("domain", false, "IDs of keyboard layouts to use of the given languages")
	// rootCmd.MarkFlagsMutuallyExclusive("email", "lib", "username", "domain")

	// Plugins
	rootCmd.PersistentFlags().StringArrayP("languages", "l", []string{"all"}, "IDs of languages to use for linguistic algorithms \n Use --list-languages to view full list of supported languages and thier IDs\n\n")
	// rootCmd.PersistentFlags().Bool("list-languages", false, "List languages and thier IDs")

	rootCmd.PersistentFlags().StringArrayP("keyboards", "k", []string{"all"}, "IDs of keyboard layouts to use of the given languages \n Use --list-keyboards to view full list of supported keyboard layouts and thier IDs \n\n")
	// rootCmd.PersistentFlags().Bool("list-keyboards", false, "List keyboards and their IDs")

	rootCmd.PersistentFlags().StringArrayP("algorithms", "a", []string{"all"}, "IDs of typo algorithms to use for generating typos \n Use --list-algorithms to view full list of supported algorithms and thier IDs \n\n")
	rootCmd.PersistentFlags().Bool("list", false, "List plugins/IDs for algorithms, langauages, information, and keyboards")

	// rootCmd.PersistentFlags().StringArrayP("info", "i", []string{"all"}, "IDs of info gathering functions to apply")

	// Processing
	rootCmd.PersistentFlags().Bool("no-cache", true, "Prevents caching of results")
	rootCmd.PersistentFlags().Bool("online", false, "Only show domains that are online")

	// Timing
	rootCmd.PersistentFlags().IntP("concurrency", "c", 50, "Number of concurrent workers")
	rootCmd.PersistentFlags().Duration("random", 1, "Random delay multiplier for network calls")
	rootCmd.PersistentFlags().Duration("delay", 1, "Duration between network calls")

	// Outputs
	rootCmd.PersistentFlags().BoolP("progress", "p", false, "Show progress bar")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show more details and remove truncated columns")
	rootCmd.PersistentFlags().StringP("file", "f", "", "Output filename defaults to stdout")
	rootCmd.PersistentFlags().StringP("format", "o", "table", "Output format: (csv,tsv,table,text,html,md)")

}
