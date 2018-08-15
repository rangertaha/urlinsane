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
	"fmt"
	"os"

	"github.com/rangertaha/urlinsane"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "urlinsane [domains]",
	Short: "Generates domain typos and variations",
	Long: `Generates domain typos and variations to detect and perform typo squatting, URL hijacking, phishing, and
	corporate espionage.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, _ := cmd.PersistentFlags().GetBool("list-keyboards")
		types, _ := cmd.PersistentFlags().GetBool("list-typos")
		langs, _ := cmd.PersistentFlags().GetBool("list-languages")

		if db {
			urlinsane.ListKeyboards(cmd, args)
		} else if langs {
			urlinsane.ListLanguages(cmd, args)
		} else if types {
			urlinsane.ListTypos(cmd, args)
		} else {
			urlinsane.Run(cmd, args)
		}

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
	// Basic options
	rootCmd.PersistentFlags().StringArrayP("keyboards", "k",
		[]string{"en1"}, "Keyboards/layouts ID to use for typing errors. "+
			"\n\tUse 'urlinsane --list-keyboards' for options.")
	rootCmd.PersistentFlags().StringArrayP("languages", "l",
		[]string{"all"}, "Language ID to use for linguistic typos."+
			"\n\tUse 'urlinsane --list-languages' for options.")

	// List keyboard and language options
	rootCmd.PersistentFlags().Bool("list-keyboards", false, "List keyboard options.")
	rootCmd.PersistentFlags().Bool("list-languages", false, "List language options")
	rootCmd.PersistentFlags().Bool("list-typos", false, "List typosquatting techniques")

	// Processing
	rootCmd.PersistentFlags().IntP("concurrency", "c", 50, "Number of workers generating results")
	rootCmd.PersistentFlags().StringArrayP("typos", "t", []string{"all"}, "Types of typos to perform."+
		"\n\tUse 'urlinsane --list-typos' for options.")

	// Post Processing options for retrieving additional data
	//rootCmd.PersistentFlags().BoolP("popularity", "p", false, "Check domain popularity with Google")
	//rootCmd.PersistentFlags().Bool("geoip", false, "Perform GeoIp location lookup")
	//rootCmd.PersistentFlags().Bool("ssdeep", false, "Get webpage and compare fuzzy hashes")
	//rootCmd.PersistentFlags().Bool("whois", false, "Query WHOIS for records")
	//rootCmd.PersistentFlags().Bool("dns", false, "Query DNS for records")

	// Output options
	//rootCmd.PersistentFlags().StringP("format", "f", "stdout", "Output format." +
	//	"Options: json, csv, and stdout")
	//rootCmd.PersistentFlags().StringP("out", "o", "", "Output file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Output additional details")
}
