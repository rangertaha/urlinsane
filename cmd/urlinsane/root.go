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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(domainCmd)

	// Plugins
	rootCmd.PersistentFlags().StringP("languages", "l", "all", "IDs of languages to use for linguistic algorithms")
	// rootCmd.PersistentFlags().Bool("list-languages", false, "List languages and thier IDs")

	rootCmd.PersistentFlags().StringP("keyboards", "k", "all", "IDs of keyboard layouts to use of the given languages")
	// rootCmd.PersistentFlags().Bool("list-keyboards", false, "List keyboards and their IDs")

	rootCmd.PersistentFlags().StringP("algorithms", "a", "all", "IDs of typo algorithms to use for generating typos")
	// rootCmd.PersistentFlags().Bool("plugins", false, "List plugins/IDs for algorithms, langauages, information, and keyboards")

	// Cache
	// rootCmd.PersistentFlags().Duration("ttl", time.Hour*24, "Cache duration for expiration")

	// Timing
	rootCmd.PersistentFlags().IntP("concurrency", "c", 50, "Number of concurrent workers")
	rootCmd.PersistentFlags().Duration("random", 1, "Random delay multiplier for network calls")
	rootCmd.PersistentFlags().Duration("delay", 1, "Duration between network calls")

	// Outputs
	rootCmd.PersistentFlags().BoolP("progress", "p", false, "Show progress bar")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show more details and remove truncated columns")
	rootCmd.PersistentFlags().StringP("file", "f", "", "Output filename defaults to stdout")
	rootCmd.PersistentFlags().StringP("format", "o", "table", "Output format: (csv,tsv,table,txt,html,md)")

}
