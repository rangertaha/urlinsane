// Copyright (C) 2024  Rangertaha <rangertaha@gmail.com>
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
package typo

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bobesa/go-domain-util/domainutil"
	// dnsLib "github.com/rangertaha/urlinsane/hackingo/net/dns"
	// geoLib "github.com/rangertaha/urlinsane/hackingo/net/geoip"
	// httpLib "github.com/rangertaha/urlinsane/hackingo/net/http"
	"golang.org/x/net/idna"

	httpLib "github.com/rangertaha/urlinsane/hackingo/net/http"
	"github.com/rangertaha/urlinsane/typo/languages"
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

	// Register ...
	Register interface {
		Set(string, ...Module)
		Get(...string) []Module
	}

	// Typosquatting ...
	Typosquatting struct {
		config Config
		// Used to store collected data of the trget domains
		meta Meta

		typoWG sync.WaitGroup
		funcWG sync.WaitGroup
		fltrWG sync.WaitGroup

		stats <-chan Statser
		errs  <-chan interface{}
	}

	// Domain ...
	Domain struct {
		Subdomain string `json:"subdomain,omitempty"`
		Domain    string `json:"domain,omitempty"`
		Suffix    string `json:"suffix,omitempty"`
		Meta      Meta   `json:"meta,omitempty"`
		Live      bool   `json:"live,omitempty"`
	}

	// Module ...
	Module struct {
		Code        string     `json:"code,omitempty"`
		Name        string     `json:"name,omitempty"`
		Description string     `json:"description,omitempty"`
		Fields      []string   `json:"-"`
		Exe         ModuleFunc `json:"-"`
	}

	// Result ...
	Result struct {
		Keyboards []languages.Keyboard `json:"-"`
		Languages []languages.Language `json:"-"`
		Original  Domain               `json:"original,omitempty"`
		Variant   Domain               `json:"variant,omitempty"`
		Typo      Module               `json:"typo,omitempty"`
		Data      map[string]string    `json:"data,omitempty"`
	}
	// DNS ...
	DNS struct {
		IPv4 []string `json:"ipv4,omitempty"`
		IPv6 []string `json:"ip46,omitempty"`
		// NS      []dnsLib.NS `json:"ns,omitempty"`
		// MX      []dnsLib.MX `json:"mx,omitempty"`
		CName   []string `json:"cname,omitempty"`
		TXT     []string `json:"txt,omitempty"`
		ipCheck bool
	}
	// Meta ...
	Meta struct {
		Levenshtein int      `json:"Levenshtein,omitempty"`
		IDNA        string   `json:"idna,omitempty"`
		IP          []string `json:"ip,omitempty"`
		Redirect    string   `json:"redirect,omitempty"`
		// HTTP        httpLib.Response `json:"http,omitempty"`
		// Geo         geoLib.Country   `json:"geo,omitempty"`
		DNS        DNS    `json:"dns,omitempty"`
		SSDeep     string `json:"ssdeep,omitempty"`
		Similarity int    `json:"similarity,omitempty"`
		// Whois    Whois      `json:"whois,omitempty"`
	}

	// OutputResult ...
	OutputResult map[string]interface{}

	// ModuleFunc defines a function to register typos.
	ModuleFunc func(Result) []Result

	// Registry ...
	Registry map[string][]Module
)

// NewRegistry ...
func NewRegistry() Registry {
	return make(Registry)
}

// Get ...
func (reg Registry) Get(names ...string) (mods []Module) {
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
func (reg Registry) Set(name string, mod ...Module) {
	_, registered := reg[strings.ToUpper(name)]
	if !registered {
		reg[strings.ToUpper(name)] = mod
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

// New ...
func New(conf Config) Typosquatting {
	return Typosquatting{config: conf}
}

// Exec ...
func (m *Module) Exec(res Result) []Result {
	return m.Exe(res)
}

// Headers ...
func (m *Module) Headers() []string {
	return m.Fields
}

// // SetMeta ...
// func (m *Result) SetMeta(key string, obj interface{}) {
// 	m.Meta[key] = obj
// }

// // GetMeta ...
// func (m *Result) GetMeta(key string) interface{} {
// 	return m.Meta[key]
// }

// SetData ...
func (m *Result) SetData(key string, obj string) {
	m.Data[key] = obj
}

// GetData ...
func (m *Result) GetData(key string) string {
	return m.Data[key]
}

// Start ...
func (typ *Typosquatting) Start() <-chan Result {
	for i, dmname := range typ.config.domains {
		records, _ := net.LookupIP(dmname.String())
		// if err != nil {
		// 	fmt.Println(err)
		// }
		for _, record := range uniqIP(records) {
			dotlen := strings.Count(record, ".")
			if dotlen == 3 {
				if !stringInSlice(record, dmname.Meta.DNS.IPv4) {
					typ.config.domains[i].Meta.DNS.IPv4 = append(dmname.Meta.DNS.IPv4, record)
				}
				typ.config.domains[i].Live = true
			}
			clen := strings.Count(record, ":")
			if clen == 5 {
				if !stringInSlice(record, dmname.Meta.DNS.IPv6) {
					typ.config.domains[i].Meta.DNS.IPv6 = append(dmname.Meta.DNS.IPv6, record)
				}
				typ.config.domains[i].Live = true
			}
		}
		if len(typ.config.domains[i].Meta.DNS.IPv4) > 0 {
			httpRes, gerr := http.Get("http://" + typ.config.domains[i].Meta.DNS.IPv4[0])
			if gerr == nil {
				res := httpLib.NewResponse(httpRes)
				// spew.Dump(res)
				typ.config.domains[i].Meta.HTTP = res

				// spew.Dump(original)

				str := httpRes.Request.URL.String()
				subdomain := domainutil.Subdomain(str)
				domain := domainutil.DomainPrefix(str)
				suffix := domainutil.DomainSuffix(str)
				if domain == "" {
					domain = str
				}
				dm := Domain{subdomain, domain, suffix, Meta{}, true}
				if dmname.String() != dm.String() {
					typ.config.domains[i].Meta.Redirect = dm.String()
				}
			}
		}
	}

	return typ.GenTypoConfig()
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
			record.Variant.Meta = Meta{}

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
			time.Sleep(typ.config.timing.Random * typ.config.timing.Delay * time.Millisecond)
			// fmt.Println(typ.config.timing.Random * typ.config.timing.Delay * time.Millisecond)
			for _, result := range xfunc.Exec(i) {
				out <- result
			}
		}
		close(out)
	}()

	if len(funcs) > 0 {
		return typ.FuncChain(funcs, out)
	}
	return out
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
	out := make(chan Result)
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
	return typ.Start()
}

// Batch returns all the results at once
func (typ *Typosquatting) Batch() (res []Result) {
	for r := range typ.Stream() {
		res = append(res, r)
	}
	return
}

// Execute starts the program and outputs results. Primarily used for CLI tools
func (typ *Typosquatting) Execute() {
	// Output results based on config
	typ.Output(typ.Stream())
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
