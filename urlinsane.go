// Copyright © 2019 CyberSecTech Inc
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
	"strings"

	"golang.org/x/net/idna"

	"github.com/cybersectech-org/urlinsane/languages"
)

type (
	URLInsane struct {
		domains   []Domain
		keyboards []languages.Keyboard
		languages []languages.Language

		types   []Typo
		funcs   []Extra
		filters []Extra

		file        string
		count       int
		format      string
		verbose     bool
		headers     []string
		concurrency int

		data map[string]map[string]string

		typoWG sync.WaitGroup
		funcWG sync.WaitGroup
		fltrWG sync.WaitGroup
	}
	Domain struct {
		Subdomain string `json:"subdomain,omitempty"`
		Domain    string `json:"domain,omitempty"`
		Suffix    string `json:"suffix,omitempty"`
	}
	Extra struct {
		Code        string    `json:"code,omitempty"`
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		Headers     []string  `json:"code,omitempty"`
		Exec        ExtraFunc `json:"-"`
	}
	Typo struct {
		Code        string   `json:"code,omitempty"`
		Name        string   `json:"name,omitempty"`
		Description string   `json:"description,omitempty"`
		Exec        TypoFunc `json:"-"`
	}
	TypoConfig struct {
		Original  Domain               `json:"original,omitempty"`
		Variant   Domain               `json:"variant,omitempty"`
		Keyboards []languages.Keyboard `json:"keyboards,omitempty"`
		Languages []languages.Language `json:"languages,omitempty"`
		Typo      Typo                 `json:"typo,omitempty"`
	}

	TypoResult struct {
		Original Domain            `json:"original,omitempty"`
		Variant  Domain            `json:"variant,omitempty"`
		Typo     Typo              `json:"typo,omitempty"`
		Live     bool              `json:"live,omitempty"`
		Data     map[string]string `json:"data,omitempty"`
	}

	OutputResult map[string]interface{}

	// TypoFunc defines a function to register typos.
	TypoFunc func(TypoConfig) []TypoConfig

	// ExtraFunc defines a function to register typos.
	ExtraFunc func(TypoResult) []TypoResult
)

const (
	VERSION = "0.5.2"
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
		filters:     c.filters,
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
				// Execute typo function returning typo results
				for _, t := range c.Typo.Exec(c) {
					if t.Variant.String() != t.Original.String() {
						out <- t
					}
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

// Results ...
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

// FilterChain creates workers of chained functions
// func (urli *URLInsane) FilterChain(in <-chan TypoResult) <-chan TypoResult {
// 	out := make(chan TypoResult)
// 	for w := 1; w <= urli.concurrency; w++ {
// 		urli.fltrWG.Add(1)
// 		go func(in <-chan TypoResult, out chan<- TypoResult) {
// 			defer urli.fltrWG.Done()
// 			output := urli.FuncChain(urli.filters, in)
// 			for c := range output {

// 				out <- c

// 			}
// 		}(in, out)
// 	}
// 	go func() {
// 		urli.fltrWG.Wait()
// 		close(out)
// 	}()
// 	return out
// }
func (urli *URLInsane) FilterChain(in <-chan TypoResult) <-chan TypoResult {
	//var xfunc Extra
	out := make(chan TypoResult)
	// xfunc, funcs = funcs[len(funcs)-1], funcs[:len(funcs)-1]
	go func() {
		for i := range in {
			if len(urli.filters) > 0 {
				for _, filter := range urli.filters {
					for _, result := range filter.Exec(i) {
						out <- result
					}
				}
			} else {
				out <- i
			}
		}
		close(out)
	}()
	return out
}

// Execute program returning results
func (urli *URLInsane) Execute() (res []TypoResult) {

	for r := range urli.Stream() {
		res = append(res, r)
	}
	return res
}

// Stream results via channels
func (urli *URLInsane) Stream() <-chan TypoResult {

	// Generate typosquatting options
	typoConfigs := urli.GenTypoConfig()

	// Execute typosquatting algorithms
	typos := urli.Typos(typoConfigs)

	// Converting typos to results and remove duplicates
	results := urli.Results(typos)

	// Execute extra functions
	output := urli.DistChain(results)

	// Execute filter functions
	filtred := urli.FilterChain(output)

	return urli.Dedup(filtred)
}

// Dedup filters the results for unique variations of domains
func (urli *URLInsane) Dedup(in <-chan TypoResult) <-chan TypoResult {
	duplicates := make(map[string]int)
	out := make(chan TypoResult)
	go func(in <-chan TypoResult, out chan<- TypoResult) {
		for c := range in {

			// Count and remove deplicates
			dup, ok := duplicates[c.Variant.String()]
			if ok {
				duplicates[c.Variant.String()] = dup + 1

			} else {
				duplicates[c.Variant.String()] = 1
				out <- c
			}

		}
		close(out)
	}(in, out)
	return out
}

// Start executes the program and outputs results. Primarily used for CLI tools
func (urli *URLInsane) Start() {

	// Execute program returning a channel with results
	output := urli.Stream()

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
