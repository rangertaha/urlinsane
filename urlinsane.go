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
	"sync"
	"golang.org/x/net/idna"
	"github.com/rangertaha/urlinsane/languages"
)

type (
	URLInsane struct {
		domains   []Domain
		keyboards []languages.Keyboard
		languages []languages.Language

		types []Typo
		funcs []Extra

		file        string
		count       int
		format      string
		verbose     bool
		headers     []string
		concurrency int

		typoWG sync.WaitGroup
		// extraWG     sync.WaitGroup
	}
	Domain struct {
		Subdomain string
		Domain    string
		Suffix    string
	}
	Extra struct {
		Code        string
		Name        string
		Description string
		Headers     []string
		Exec        ExtraFunc
	}
	Typo struct {
		Code        string
		Name        string
		Description string
		Exec        TypoFunc
	}
	TypoConfig struct {
		Domain    Domain
		Keyboards []languages.Keyboard
		Languages []languages.Language
		Typo      Typo
	}

	TypoResult struct {
		Domain Domain
		Typo   Typo
		Live   bool
		Data   map[string]string
	}

	// TypoFunc defines a function to register typos.
	TypoFunc func(TypoConfig) []TypoConfig

	// ExtraFunc defines a function to register typos.
	ExtraFunc func(TypoResult) []TypoResult
)

const (
	VERSION = "0.1.0"
	DEBUG   = false
	LOGO    = `
 _   _  ____   _      ___
| | | ||  _ \ | |    |_ _| _ __   ___   __ _  _ __    ___
| | | || |_) || |     | | | '_ \ / __| / _' || '_ \  / _ \
| |_| ||  _ < | |___  | | | | | |\__ \| (_| || | | ||  __/
 \___/ |_| \_\|_____||___||_| |_||___/ \__,_||_| |_| \___|

 Version: ` + VERSION + "\n"
)

// New
func New(c Config) (i URLInsane) {
	i = URLInsane{
		domains:     c.domains,
		keyboards:   c.keyboards,
		languages:   c.languages,
		types:       c.typos,
		funcs:       c.funcs,
		concurrency: c.concurrency,
		headers:     c.headers,
		format:      c.format,
		file:        c.file,
		verbose:     c.verbose,
	}
	return
}

// GenTypoConfig
func (urli *URLInsane) GenTypoConfig() <-chan TypoConfig {
	out := make(chan TypoConfig)
	go func() {
		for _, domain := range urli.domains {
			for _, typo := range urli.types {
				urli.count++
				out <- TypoConfig{domain, urli.keyboards, urli.languages, typo}
			}
		}
		close(out)
	}()
	return out
}

// Typos gives typo options to a pool of workers
func (urli *URLInsane) Typos(in <-chan TypoConfig) <-chan TypoConfig {
	out := make(chan TypoConfig)
	for w := 1; w <= urli.concurrency; w++ {
		urli.typoWG.Add(1)
		go func(id int, in <-chan TypoConfig, out chan<- TypoConfig) {
			defer urli.typoWG.Done()
			for c := range in {
				for _, t := range c.Typo.Exec(c) {
					out <- t
				}
			}
		}(w, in, out)
	}
	go func() {
		urli.typoWG.Wait()
		close(out)
	}()
	return out
}

// Results
func (urli *URLInsane) Results(in <-chan TypoConfig) <-chan TypoResult {
	out := make(chan TypoResult)
	go func() {
		for r := range in {
			record := TypoResult{Domain: r.Domain, Typo: r.Typo}

			// Initialize a place to store extra data for a record
			record.Data = make(map[string]string)

			// Add record placeholder for consistent records
			for _, name := range urli.headers {
				_, ok := record.Data[name]
				if !ok {
					record.Data[name] = " "
				}
			}

			out <- record
		}
		close(out)
	}()
	return out
}

// Funcs
func (urli *URLInsane) Funcs(funcs []Extra, in <-chan TypoResult) <-chan TypoResult {
	var xfunc Extra
	out := make(chan TypoResult)
	xfunc, funcs = funcs[len(funcs)-1], funcs[:len(funcs)-1]
	go func() {
		for i := range in {
			for _, result := range xfunc.Exec(i) {
				out <- result
			}
		}
		close(out)
	}()

	if len(funcs) > 0 {
		return urli.Funcs(funcs, out)
	} else {
		return out
	}
}

// Execute program returning a channel with typos
func (urli *URLInsane) Execute() <-chan TypoResult {

	// Generate typosquatting options
	typoConfigs := urli.GenTypoConfig()

	// Execute typosquatting algorithms
	typos := urli.Typos(typoConfigs)

	// Converting typos to results and remove duplicates
	results := urli.Results(typos)

	// Execute functions to retrieve additional info
	output := urli.Funcs(urli.funcs, results)

	return output
}

// Start executes the program and outputs results. Primarily used for CLI tools
func (urli *URLInsane) Start() {

	// Execute program returning a channel with results
	output := urli.Execute()

	// Output results based on config
	urli.Output(output)
}

// Idna
func (d *Domain) Idna() (punycode string) {
	punycode, _ = idna.Punycode.ToASCII(d.String())
	return
}

// String
func (d *Domain) String() (domain string) {
	if d.Subdomain != "" {
		domain = d.Subdomain + "."
	}
	if d.Domain != "" {
		domain = domain + d.Domain
	}
	if d.Suffix != "" {
		domain = domain + "." + d.Suffix
	}
	return
}