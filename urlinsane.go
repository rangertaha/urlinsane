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

package urlinsane

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/olekukonko/tablewriter"
	"github.com/rangertaha/urlinsane/languages"
	"strings"
)

type (
	URLInsane struct {
		domains   []string
		types     []Typo
		keyboards []languages.Keyboard
		languages []languages.Language

		concurrency int
		format      string
		out         string
		verbose     bool

		wg sync.WaitGroup
	}
	Domain struct {
		Sub    string
		Domain string
		TLD    string
	}
	Typo struct {
		Code        string
		Name        string
		Description string
		Exec        TypoFunc
	}
	TypoConfig struct {
		Domain    string
		Keyboards []languages.Keyboard
		Languages []languages.Language
		Typo      Typo
	}

	TypoResult struct {
		Domain string
		Typo   Typo
	}

	// TypoFunc defines a function to register typos.
	TypoFunc func(TypoConfig) []TypoConfig
)

// typoRegistry = make(map[string][]Typo)
var Typos = make(map[string][]Typo)

// NewCLI
func NewCLI(cmd *cobra.Command, args []string) (i URLInsane) {
	// Basic options
	keyboards, _ := cmd.PersistentFlags().GetStringArray("keyboards")
	langs, _ := cmd.PersistentFlags().GetStringArray("languages")
	types, _ := cmd.PersistentFlags().GetStringArray("typos")

	// Processing option
	concurrency, _ := cmd.PersistentFlags().GetInt("concurrency")

	// Output options
	output, _ := cmd.PersistentFlags().GetString("out")
	format, _ := cmd.PersistentFlags().GetString("format")
	verbose, _ := cmd.PersistentFlags().GetBool("verbose")

	i = URLInsane{
		domains:     args,
		keyboards:   languages.GetKeyboards(keyboards),
		languages:   languages.GetLanguages(langs),
		types:       GetTypos(types),
		concurrency: concurrency,

		format:  format,
		out:     output,
		verbose: verbose,
	}
	return
}

// Register
func Register(name string, typo ...Typo) {
	_, registered := Typos[strings.ToLower(name)]
	if !registered {
		Typos[strings.ToLower(name)] = typo
	}
}

// GetTypos
func GetTypos(codes []string) (typos []Typo) {
	for _, tcode := range codes {
		typo, ok := Typos[strings.ToLower(tcode)]
		if ok {
			typos = append(typos, typo...)
		}
	}
	return
}

// GenTypoConfigs
func (urli *URLInsane) GenTypoConfigs() <-chan TypoConfig {
	out := make(chan TypoConfig)
	go func() {
		for _, domain := range urli.domains {
			for _, typo := range urli.types {
				out <- TypoConfig{domain, urli.keyboards, urli.languages, typo}
			}
		}

		close(out)
	}()
	return out
}

// worker executes the typosquatting algorithms
func (urli *URLInsane) worker(id int, in <-chan TypoConfig, out chan<- TypoConfig) {
	defer urli.wg.Done()
	for c := range in {
		for _, t := range c.Typo.Exec(c) {
			out <- t
		}
	}
}

// Process gives typo options to a pool of workers
func (urli *URLInsane) Process(in <-chan TypoConfig) <-chan TypoConfig {
	out := make(chan TypoConfig)
	for w := 1; w <= urli.concurrency; w++ {
		urli.wg.Add(1)
		go urli.worker(w, in, out)
	}

	go func() {
		urli.wg.Wait()
		close(out)
	}()

	return out
}

// Results
func (urli *URLInsane) Results(in <-chan TypoConfig) <-chan TypoResult {
	out := make(chan TypoResult)
	go func() {
		for r := range in {
			out <- TypoResult{r.Domain, r.Typo}
		}
		close(out)
	}()
	return out
}

// PostProcess
func (urli *URLInsane) PostProcess(in <-chan TypoResult) <-chan TypoResult {
	out := make(chan TypoResult)
	go func() {
		for r := range in {
			out <- r
		}
		close(out)
	}()
	return out
}

func (urli *URLInsane) jsonOutput(in <-chan TypoResult) {
}

func (urli *URLInsane) csvOutput(in <-chan TypoResult) {
}

func (urli *URLInsane) stdOutput(in <-chan TypoResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Type", "Domain"})

	for v := range in {
		table.Append([]string{v.Typo.Name, v.Domain})
	}
	table.Render()

	//w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	//fmt.Fprintln(w, "Type\tDomain")
	//fmt.Fprintln(w, "-----------------------------------------------------")
	//for c := range in {
	//	//fmt.Println(c.Domain, c.Keyboard.Name, c.Typo)
	//	if urli.verbose {
	//		fmt.Fprintln(w, c.Typo.Name + "\t" + c.Domain + "\t")
	//	} else {
	//		fmt.Fprintln(w, strings.ToUpper(c.Typo.Code) + "\t" + c.Domain + "\t")
	//	}
	//
	//}
	//w.Flush()
}

func (urli *URLInsane) Output(in <-chan TypoResult) {
	urli.stdOutput(in)
}

// Run is the main program
func Run(cmd *cobra.Command, args []string) {

	urli := NewCLI(cmd, args)

	// Generate typosquatting options
	typoConfigs := urli.GenTypoConfigs()

	// Execute typosquatting generating results
	typos := urli.Process(typoConfigs)

	// Converting typos to results
	results := urli.Results(typos)

	// TODO: Add additional info to the results via post process functions
	output := urli.PostProcess(results)

	// Output results
	urli.Output(output)

}

// ListKeyboards lists all supported keyboards
func ListKeyboards(cmd *cobra.Command, args []string) {
	urli := NewCLI(cmd, args)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tLetters")
	for _, keyboard := range urli.keyboards {
		fmt.Fprintln(w, keyboard.Code, keyboard.Name, keyboard.Language.Graphemes)
	}
	w.Flush()
}

// ListLanguages lists all supported languages
func ListLanguages(cmd *cobra.Command, args []string) {
	urli := NewCLI(cmd, args)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tName")
	for _, lang := range urli.languages {
		fmt.Fprintln(w, lang.Code+"\t"+lang.Name)
	}
	w.Flush()
}

// ListTypos lists typosquatting techniques
func ListTypos(cmd *cobra.Command, args []string) {
	urli := NewCLI(cmd, args)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tDescription")
	for _, typo := range urli.types {
		fmt.Fprintln(w, typo.Code+"\t"+typo.Name+"\t"+typo.Description)
	}
	w.Flush()
}
