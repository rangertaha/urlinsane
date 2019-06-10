// The MIT License (MIT)
//
// Copyright © 2019 Rangertaha <rangertaha@gmail.com>
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

package typo

import (
	"strings"
	"sync"

	"golang.org/x/net/idna"

	"github.com/cybersectech-org/urlinsane/pkg/typo/languages"
)

type (
	// Moduler  ...
	Moduler interface {
		Exec(Result) []Result
		Headers() []string
	}

	// Metadater ...
	Metadater interface {
		String() string
		Object() interface{}
	}
	// Statser ...
	Statser interface{}

	// Registry ...
	Register interface {
		Set(string, ...Module)
		Get(...string) []Module
	}

	// Typosquatting ...
	Typosquatting struct {
		config Config
		// Used to store collected data of the trget domains
		meta map[string]interface{}

		typoWG sync.WaitGroup
		funcWG sync.WaitGroup
		fltrWG sync.WaitGroup

		stats <-chan Statser
	}

	// Domain ...
	Domain struct {
		Subdomain string `json:"subdomain,omitempty"`
		Domain    string `json:"domain,omitempty"`
		Suffix    string `json:"suffix,omitempty"`
	}

	// Module ...
	Module struct {
		Code        string `json:"code,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		headers     []string
		exec        ModuleFunc
	}

	// Result ...
	Result struct {
		Keyboards []languages.Keyboard
		Languages []languages.Language
		Original  Domain                 `json:"original,omitempty"`
		Variant   Domain                 `json:"variant,omitempty"`
		Typo      Module                 `json:"typo,omitempty"`
		Data      map[string]string      `json:"data,omitempty"`
		Meta      map[string]interface{} `json:"meta,omitempty"`
		Live      bool                   `json:"live,omitempty"`
	}

	// OutputResult ...
	OutputResult map[string]interface{}

	// ModuleFunc defines a function to register typos.
	ModuleFunc func(Result) []Result

	registry map[string][]Module
)

// NewRegistry ...
func NewRegistry() registry {
	return make(registry)
}

// Get ...
func (reg registry) Get(names ...string) (mods []Module) {
	for _, f := range names {
		value, ok := reg[strings.ToUpper(f)]
		if ok {
			mods = append(mods, value...)
		}
	}
	if len(names) == 0 {
		return reg.Get("all")
	}
	return
}

// Set ...
func (reg registry) Set(name string, mod ...Module) {
	_, registered := reg[strings.ToUpper(name)]
	if !registered {
		reg[strings.ToUpper(name)] = mod
	}
}

// New ...
func New(conf Config) Typosquatting {
	return Typosquatting{config: conf}
}

// Exec ...
func (m *Module) Exec(res Result) []Result {
	return m.exec(res)
}

// Headers ...
func (m *Module) Headers() []string {
	return m.headers
}

// SetMeta ...
func (m *Result) SetMeta(key string, obj interface{}) {
	m.Meta[key] = obj
}

// GetMeta ...
func (m *Result) GetMeta(key string) interface{} {
	return m.Meta[key]
}

// SetData ...
func (m *Result) SetData(key string, obj string) {
	m.Data[key] = obj
}

// GetData ...
func (m *Result) GetData(key string) string {
	return m.Data[key]
}

// GenTypoConfig ...
func (typ *Typosquatting) GenTypoConfig() <-chan Result {
	out := make(chan Result)
	go func() {
		for _, domain := range typ.config.domains {
			for _, typo := range typ.config.typos {
				out <- Result{
					Original:  domain,
					Variant:   Domain{},
					Typo:      typo,
					Keyboards: typ.config.keyboards,
					Languages: typ.config.languages,
				}
			}
		}
		close(out)
	}()
	return typ.Typos(out)
}

// Typos gives typo options to a pool of workers
func (typ *Typosquatting) Typos(in <-chan Result) <-chan Result {
	out := make(chan Result)
	for w := 1; w <= typ.config.concurrency; w++ {
		typ.typoWG.Add(1)
		go func(id int, in <-chan Result, out chan<- Result) {
			defer typ.typoWG.Done()
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
		typ.typoWG.Wait()
		close(out)
	}()
	return typ.Results(out)
}

// Results ...
func (typ *Typosquatting) Results(in <-chan Result) <-chan Result {
	out := make(chan Result)
	go func() {
		for r := range in {
			record := Result{Variant: r.Variant, Original: r.Original, Typo: r.Typo}

			// Initialize a place to store extra data for a record
			record.Data = make(map[string]string)

			// Initialize a place to store meta data
			record.Meta = make(map[string]interface{})

			// Add record placeholder for consistent records
			for _, name := range typ.config.headers {
				_, ok := record.Data[name]
				if !ok {
					record.Data[name] = ""
				}
			}
			out <- record
		}
		close(out)
	}()
	return typ.DistChain(out)
}

// FuncChain creates a chain of information gathering functions
func (typ *Typosquatting) FuncChain(funcs []Module, in <-chan Result) <-chan Result {
	var xfunc Module
	out := make(chan Result)
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
		return typ.FuncChain(funcs, out)
	} else {
		return out
	}
}

// DistChain creates workers of chained functions
func (typ *Typosquatting) DistChain(in <-chan Result) <-chan Result {
	out := make(chan Result)
	for w := 1; w <= typ.config.concurrency; w++ {
		typ.funcWG.Add(1)
		go func(in <-chan Result, out chan<- Result) {
			defer typ.funcWG.Done()
			output := typ.FuncChain(typ.config.funcs, in)
			for c := range output {
				out <- c
			}
		}(in, out)
	}
	go func() {
		typ.funcWG.Wait()
		close(out)
	}()
	return typ.FilterChain(out)
}

// FilterChain ...
func (typ *Typosquatting) FilterChain(in <-chan Result) <-chan Result {
	//var xfunc Extra
	out := make(chan Result)
	// xfunc, funcs = funcs[len(funcs)-1], funcs[:len(funcs)-1]
	go func() {
		for i := range in {
			if len(typ.config.filters) > 0 {
				for _, filter := range typ.config.filters {
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
	return typ.Dedup(out)
}

// Dedup filters the results for unique variations of domains
func (typ *Typosquatting) Dedup(in <-chan Result) <-chan Result {
	duplicates := make(map[string]int)
	out := make(chan Result)
	go func(in <-chan Result, out chan<- Result) {
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

// Stream returns the results one at a time
func (typ *Typosquatting) Stream() <-chan Result {
	return typ.GenTypoConfig()
}

// Batch returns all the results at once
func (typ *Typosquatting) Batch() (res []Result) {

	for r := range typ.Stream() {
		res = append(res, r)
	}
	return res
}

// Execute starts the program and outputs results. Primarily used for CLI tools
func (typ *Typosquatting) Execute() {

	// Execute program returning a channel with results
	output := typ.Stream()

	// Output results based on config
	typ.Output(output)
}

// Output ...
func (typ *Typosquatting) Output(in <-chan Result) {
	if typ.config.format == "json" {
		typ.jsonOutput(in)
	}
	if typ.config.format == "csv" {
		typ.csvOutput(in)
	}
	if typ.config.format == "text" {
		typ.stdOutput(in)
	}
}

// Idna ...
func (d *Domain) Idna() (punycode string) {
	punycode, _ = idna.Punycode.ToASCII(d.String())
	return
}

// String ...
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
