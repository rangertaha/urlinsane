// Copyright © 2019 rangertaha rangertaha@gmail.com
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

	"strings"

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

		data map[string]map[string]string

		typoWG sync.WaitGroup
		funcWG sync.WaitGroup
	}
	Domain struct {
		Subdomain string `json:"subdomain"`
		Domain    string `json:"domain"`
		Suffix    string `json:"suffix"`
	}
	Extra struct {
		Code        string
		Name        string
		Description string
		Headers     []string
		Exec        ExtraFunc
	}
	Typo struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Exec        TypoFunc
	}
	TypoConfig struct {
		Original  Domain
		Variant   Domain
		Keyboards []languages.Keyboard
		Languages []languages.Language
		Typo      Typo
	}

	TypoResult struct {
		Original Domain            `json:"original"`
		Variant  Domain            `json:"variant"`
		Typo     Typo              `json:"typo"`
		Live     bool              `json:"live"`
		Data     map[string]string `json:"data"`
	}

	// TypoFunc defines a function to register typos.
	TypoFunc func(TypoConfig) []TypoConfig

	// ExtraFunc defines a function to register typos.
	ExtraFunc func(TypoResult) []TypoResult
)

const (
	VERSION = "0.3.0"
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
				out <- TypoConfig{domain, Domain{}, urli.keyboards, urli.languages, typo}
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
			record := TypoResult{Variant: r.Variant, Original: r.Original, Typo: r.Typo}

			// Initialize a place to store extra data for a record
			record.Data = make(map[string]string)

			// Add record placeholder for consistent records
			for _, name := range urli.headers {
				_, ok := record.Data[name]
				if !ok {
					record.Data[name] = ""
				}
			}

			out <- record
		}
		close(out)
	}()
	return out
}

// FuncChain creates a chain of extra functions
func (urli *URLInsane) FuncChain(funcs []Extra, in <-chan TypoResult) <-chan TypoResult {
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
		return urli.FuncChain(funcs, out)
	} else {
		return out
	}
}

// DistChain creates workers of chained functions
func (urli *URLInsane) DistChain(in <-chan TypoResult) <-chan TypoResult {
	out := make(chan TypoResult)
	for w := 1; w <= urli.concurrency; w++ {
		urli.funcWG.Add(1)
		go func(in <-chan TypoResult, out chan<- TypoResult) {
			defer urli.funcWG.Done()
			output := urli.FuncChain(urli.funcs, in)
			for c := range output {
				out <- c
			}
		}(in, out)
	}
	go func() {
		urli.funcWG.Wait()
		close(out)
	}()
	return out
}

// Execute program returning a channel with typos
func (urli *URLInsane) Execute() <-chan TypoResult {

	// Generate typosquatting options
	typoConfigs := urli.GenTypoConfig()

	// Execute typosquatting algorithms
	typos := urli.Typos(typoConfigs)

	// Converting typos to results and remove duplicates
	results := urli.Results(typos)

	// Execute extra functions
	output := urli.DistChain(results)

	return output
}

// Dedup filters the results for unique variations of domains
func (urli *URLInsane) Dedup(in <-chan TypoResult) <-chan TypoResult {
	out := make(chan TypoResult)

	return out
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
	domain = strings.TrimSpace(domain)
	return
}
