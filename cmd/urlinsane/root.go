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
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DefualtConfig = []byte(``)

const CONFIG_DIR = ".urlinsane"
const CONFIG_FILE = "urlinsane"

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
	var k = koanf.New(".")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// rootCmd.AddCommand(domainCmd)

	// Plugins
	rootCmd.PersistentFlags().StringP("languages", "l", "all", "IDs of languages to use for linguistic algorithms")
	viper.BindPFlag("languages", rootCmd.Flags().Lookup("languages"))

	rootCmd.PersistentFlags().StringP("keyboards", "k", "all", "IDs of keyboard layouts to use of the given languages")
	viper.BindPFlag("keyboards", rootCmd.Flags().Lookup("keyboards"))

	rootCmd.PersistentFlags().StringP("algorithms", "a", "all", "IDs of typo algorithms to use for generating typos")
	viper.BindPFlag("algorithms", rootCmd.Flags().Lookup("algorithms"))

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

	viper.SetConfigName(CONFIG_FILE) // name of config file (without extension)
	viper.SetConfigType("hcl")       // REQUIRED if the config file does not have the extension in the name
	// viper.AddConfigPath("/etc/urlinsane/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/" + CONFIG_DIR) // call multiple times to add many search paths
	// viper.AddConfigPath(".")                // optionally look for config in the working directory
	// err := viper.ReadInConfig()             // Find and read the config file
	// if err != nil {                         // Handle errors reading the config file
	// 	panic(fmt.Errorf("fatal error config file: %w", err))
	// }

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			if err := CreateConfigPath(CONFIG_DIR); err != nil {
				fmt.Print(err)
			}

			// viper.ReadConfig(bytes.NewBuffer(DefualtConfig))
			// appDir := filepath.Join(userDir, appdir)
			if err := k.Load(file.Provider(c), hcl.Parser(true)); err != nil {
				log.Fatalf("error loading file: %v", err)
			}

			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Print(err)
			}

		} else {
			// Config file was found but another error was produced
		}
	}
}

func CreateConfigPath(appdir string) (err error) {
	var userDir string

	if userDir, err = os.UserHomeDir(); err != nil {
		if userDir, err = os.Getwd(); err != nil {
			userDir = ""
		}
	}

	appDir := filepath.Join(userDir, appdir)
	err = os.MkdirAll(appDir, 0750)
	if err != nil {
		return err
	}

	return
}

// viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

// // any approach to require this configuration into your program.
// var yamlExample = []byte(`
// Hacker: true
// name: steve
// hobbies:
// - skateboarding
// - snowboarding
// - go
// clothing:
//   jacket: leather
//   trousers: denim
// age: 35
// eyes : brown
// beard: true
// `)

// viper.ReadConfig(bytes.NewBuffer(yamlExample))

// func (c *AppDir) Init() (err error) {
// 	if c.homeDir == "" {
// 		if c.homeDir, err = os.UserHomeDir(); err != nil {
// 			if c.homeDir, err = os.Getwd(); err != nil {
// 				c.homeDir = ""
// 			}
// 		}
// 	}
// 	if c.appDir == "" {
// 		c.appDir = filepath.Join(c.homeDir, ".urlinsane")
// 	}
// 	if c.appCfg == "" {
// 		c.appCfg = filepath.Join(c.appDir, "config.yml")
// 	}

// 	return c.getOrCreate()
// }

// func (c *AppDir) getOrCreate() (err error) {
// 	// Create app directory in user's home directory or local directory
// 	err = os.MkdirAll(c.appDir, 0750)
// 	if err != nil {
// 		return err
// 	}
