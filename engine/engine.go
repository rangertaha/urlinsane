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
package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/config"
)

type (

	// Typosquatting ...
	Typosquatting struct {
		Config config.Config
		Typos  map[string]urlinsane.Typo

		typoWG sync.WaitGroup
		funcWG sync.WaitGroup

		// stats <-chan Statser
		errs <-chan interface{}
	}
)

// New ...
func New(conf config.Config) Typosquatting {
	return Typosquatting{
		Config: conf,
	}
}

// Generate typo config options
func (t *Typosquatting) GenOptions() <-chan urlinsane.Typo {
	// selects between names and domains algos

	out := make(chan urlinsane.Typo)
	go func() {
		for _, lang := range t.Config.Languages {
			for _, board := range t.Config.Keyboards {
				for _, algo := range t.Config.Algorithms {
					if algo.IsType(t.Config.Type) {
						out <- Typo{
							language:  lang,
							keyboard:  board,
							algorithm: algo,
							name:      t.Config.Name,
							original:  t.Config.Domain,
						}
					}
				}
			}
		}
		close(out)
	}()
	return t.Algorithms(out)
}

// Algorithms generates typos using the algorithm plugins
func (t *Typosquatting) Algorithms(in <-chan urlinsane.Typo) <-chan urlinsane.Typo {
	out := make(chan urlinsane.Typo)

	for w := 1; w <= t.Config.Concurrency; w++ {
		t.typoWG.Add(1)
		go func(id int, in <-chan urlinsane.Typo, out chan<- urlinsane.Typo) {
			defer t.typoWG.Done()
			for c := range in {
				// Execute typo algorithm returning typos
				for _, t := range c.Algorithm().Exec(c) {
					out <- t
				}
			}
		}(w, in, out)
	}
	go func() {
		t.typoWG.Wait()
		close(out)
	}()

	return t.Dedup(out)
}

func (t *Typosquatting) Dedup(in <-chan urlinsane.Typo) <-chan urlinsane.Typo {
	out := make(chan urlinsane.Typo)
	go func() {
		for typo := range in {
			t.Typos[typo.Repr()] = typo
		}
		for _, typ := range t.Typos {
			out <- typ
		}
		close(out)
	}()

	return t.Information(out)
}

func (t *Typosquatting) Information(in <-chan urlinsane.Typo) <-chan urlinsane.Typo {
	out := make(chan urlinsane.Typo)
	for w := 1; w <= t.Config.Concurrency; w++ {
		t.funcWG.Add(1)
		go func(in <-chan urlinsane.Typo, out chan<- urlinsane.Typo) {
			defer t.funcWG.Done()
			output := t.InfoChain(t.Config.Information, in)
			for c := range output {
				out <- c
			}
		}(in, out)
	}
	go func() {
		t.funcWG.Wait()
		close(out)
	}()
	return t.Cache(out)
}

// // DistChain creates workers of chained functions
// func (t *Typosquatting) DistChain(in <-chan urlinsane.Typo) <-chan urlinsane.Typo {
// 	out := make(chan urlinsane.Typo)
// 	for w := 1; w <= t.Config.Concurrency; w++ {
// 		t.funcWG.Add(1)
// 		go func(in <-chan urlinsane.Typo, out chan<- urlinsane.Typo) {
// 			defer t.funcWG.Done()
// 			output := t.InfoChain(t.Config.Information, in)
// 			for c := range output {
// 				out <- c
// 			}
// 		}(in, out)
// 	}
// 	go func() {
// 		t.funcWG.Wait()
// 		close(out)
// 	}()
// 	return t.Information(out)
// }


// FuncChain creates a chain of information gathering functions
func (t *Typosquatting) InfoChain(funcs []urlinsane.Information, in <-chan urlinsane.Typo) <-chan urlinsane.Typo {
	var xfunc urlinsane.Information
	out := make(chan urlinsane.Typo)
	xfunc, funcs = funcs[len(funcs)-1], funcs[:len(funcs)-1]
	go func() {
		for i := range in {
			time.Sleep(t.Config.Random * t.Config.Delay)
			// fmt.Println(typ.config.timing.Random * typ.config.timing.Delay * time.Millisecond)

			out <- xfunc.Exec(i)

		}
		close(out)
	}()

	if len(funcs) > 0 {
		return t.InfoChain(funcs, out)
	}
	return out
}


func (t *Typosquatting) Cache(in <-chan urlinsane.Typo) <-chan urlinsane.Typo {
	// out := make(chan urlinsane.Typo)
	// go func() {
	// 	for typo := range in {
	// 		t.Typos[typo.Repr()] = typo
	// 	}
	// 	for _, typ := range t.Typos {
	// 		out <- typ
	// 	}
	// 	close(out)
	// }()

	return in
}


func (t *Typosquatting) Output(in <-chan urlinsane.Typo) {
	fmt.Println(t)
}

func (t *Typosquatting) Execute() {
	t.Output(t.GenOptions())
}

// // FilterChain ...
// func (typ *Typosquatting) FilterChain(in <-chan Result) <-chan Result {
// 	out := make(chan Result)
// 	go func() {
// 		for i := range in {
// 			if len(typ.config.filters) > 0 {
// 				for _, filter := range typ.config.filters {
// 					for _, result := range filter.Exec(i) {
// 						out <- result
// 					}
// 				}
// 			} else {
// 				out <- i
// 			}
// 		}
// 		close(out)
// 	}()
// 	return typ.Dedup(out)
// }

// // Start ...
// func (typ *Typosquatting) Start() <-chan Result {
// 	for i, dmname := range typ.config.domains {
// 		records, _ := net.LookupIP(dmname.String())
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// }
// 		for _, record := range uniqIP(records) {
// 			dotlen := strings.Count(record, ".")
// 			if dotlen == 3 {
// 				if !stringInSlice(record, dmname.Meta.DNS.IPv4) {
// 					typ.config.domains[i].Meta.DNS.IPv4 = append(dmname.Meta.DNS.IPv4, record)
// 				}
// 				typ.config.domains[i].Live = true
// 			}
// 			clen := strings.Count(record, ":")
// 			if clen == 5 {
// 				if !stringInSlice(record, dmname.Meta.DNS.IPv6) {
// 					typ.config.domains[i].Meta.DNS.IPv6 = append(dmname.Meta.DNS.IPv6, record)
// 				}
// 				typ.config.domains[i].Live = true
// 			}
// 		}
// 		if len(typ.config.domains[i].Meta.DNS.IPv4) > 0 {
// 			httpRes, gerr := http.Get("http://" + typ.config.domains[i].Meta.DNS.IPv4[0])
// 			if gerr == nil {
// 				res := httpLib.NewResponse(httpRes)
// 				// spew.Dump(res)
// 				typ.config.domains[i].Meta.HTTP = res

// 				// spew.Dump(original)

// 				str := httpRes.Request.URL.String()
// 				subdomain := domainutil.Subdomain(str)
// 				domain := domainutil.DomainPrefix(str)
// 				suffix := domainutil.DomainSuffix(str)
// 				if domain == "" {
// 					domain = str
// 				}
// 				dm := Domain{subdomain, domain, suffix, Meta{}, true}
// 				if dmname.String() != dm.String() {
// 					typ.config.domains[i].Meta.Redirect = dm.String()
// 				}
// 			}
// 		}
// 	}

// 	return typ.GenTypoConfig()
// }

// // GenTypoConfig ...
// func (typ *Typosquatting) GenTypoConfig() <-chan Result {
// 	out := make(chan Result)
// 	go func() {
// 		for _, domain := range typ.config.domains {
// 			for _, typo := range typ.config.typos {
// 				out <- Result{
// 					Original:  domain,
// 					Variant:   Domain{},
// 					Typo:      typo,
// 					Keyboards: typ.config.keyboards,
// 					Languages: typ.config.languages,
// 				}
// 			}
// 		}
// 		close(out)
// 	}()
// 	return typ.Typos(out)
// }

// // Typos gives typo options to a pool of workers
// func (typ *Typosquatting) Typos(in <-chan Result) <-chan Result {
// 	out := make(chan Result)
// 	for w := 1; w <= typ.config.concurrency; w++ {
// 		typ.typoWG.Add(1)
// 		go func(id int, in <-chan Result, out chan<- Result) {
// 			defer typ.typoWG.Done()
// 			for c := range in {
// 				// Execute typo function returning typo results
// 				for _, t := range c.Typo.Exec(c) {
// 					if t.Variant.String() != t.Original.String() {
// 						out <- t
// 					}
// 				}
// 			}
// 		}(w, in, out)
// 	}
// 	go func() {
// 		typ.typoWG.Wait()
// 		close(out)
// 	}()
// 	return typ.Results(out)
// }

// // Results ...
// func (typ *Typosquatting) Results(in <-chan Result) <-chan Result {
// 	out := make(chan Result)
// 	go func() {
// 		for r := range in {
// 			record := Result{Variant: r.Variant, Original: r.Original, Typo: r.Typo}
// 			// Initialize a place to store extra data for a record
// 			record.Data = make(map[string]string)

// 			// Initialize a place to store meta data
// 			record.Variant.Meta = Meta{}

// 			// Add record placeholder for consistent records
// 			for _, name := range typ.config.headers {
// 				_, ok := record.Data[name]
// 				if !ok {
// 					record.Data[name] = ""
// 				}
// 			}
// 			out <- record
// 		}
// 		close(out)
// 	}()
// 	return typ.DistChain(out)
// }

// // FuncChain creates a chain of information gathering functions
// func (typ *Typosquatting) FuncChain(funcs []Module, in <-chan Result) <-chan Result {
// 	var xfunc Module
// 	out := make(chan Result)
// 	xfunc, funcs = funcs[len(funcs)-1], funcs[:len(funcs)-1]
// 	go func() {
// 		for i := range in {
// 			time.Sleep(typ.config.timing.Random * typ.config.timing.Delay * time.Millisecond)
// 			// fmt.Println(typ.config.timing.Random * typ.config.timing.Delay * time.Millisecond)
// 			for _, result := range xfunc.Exec(i) {
// 				out <- result
// 			}
// 		}
// 		close(out)
// 	}()

// 	if len(funcs) > 0 {
// 		return typ.FuncChain(funcs, out)
// 	}
// 	return out
// }

// // DistChain creates workers of chained functions
// func (typ *Typosquatting) DistChain(in <-chan Result) <-chan Result {
// 	out := make(chan Result)
// 	for w := 1; w <= typ.config.concurrency; w++ {
// 		typ.funcWG.Add(1)
// 		go func(in <-chan Result, out chan<- Result) {
// 			defer typ.funcWG.Done()
// 			output := typ.FuncChain(typ.config.funcs, in)
// 			for c := range output {
// 				out <- c
// 			}
// 		}(in, out)
// 	}
// 	go func() {
// 		typ.funcWG.Wait()
// 		close(out)
// 	}()
// 	return typ.FilterChain(out)
// }

// // FilterChain ...
// func (typ *Typosquatting) FilterChain(in <-chan Result) <-chan Result {
// 	out := make(chan Result)
// 	go func() {
// 		for i := range in {
// 			if len(typ.config.filters) > 0 {
// 				for _, filter := range typ.config.filters {
// 					for _, result := range filter.Exec(i) {
// 						out <- result
// 					}
// 				}
// 			} else {
// 				out <- i
// 			}
// 		}
// 		close(out)
// 	}()
// 	return typ.Dedup(out)
// }

// // Dedup filters the results for unique variations of domains
// func (typ *Typosquatting) Dedup(in <-chan Result) <-chan Result {
// 	duplicates := make(map[string]int)
// 	out := make(chan Result)
// 	go func(in <-chan Result, out chan<- Result) {
// 		for c := range in {
// 			// Count and remove deplicates
// 			dup, ok := duplicates[c.Variant.String()]
// 			if ok {
// 				duplicates[c.Variant.String()] = dup + 1

// 			} else {
// 				duplicates[c.Variant.String()] = 1
// 				out <- c
// 			}
// 		}
// 		close(out)
// 	}(in, out)
// 	return out
// }

// // Stream returns the results one at a time
// func (typ *Typosquatting) Stream() <-chan Result {
// 	return typ.Start()
// }

// // Batch returns all the results at once
// func (typ *Typosquatting) Batch() (res []Result) {
// 	for r := range typ.Stream() {
// 		res = append(res, r)
// 	}
// 	return
// }

// // Execute starts the program and outputs results. Primarily used for CLI tools
// func (typ *Typosquatting) Execute() {
// 	// Output results based on config
// 	typ.Output(typ.Stream())
// }

// // Output ...
// func (typ *Typosquatting) Output(in <-chan Result) {
// 	if typ.config.format == "json" {
// 		typ.jsonOutput(in)
// 	}
// 	if typ.config.format == "csv" {
// 		typ.csvOutput(in)
// 	}
// 	if typ.config.format == "text" {
// 		typ.stdOutput(in)
// 	}
// }

// func stringInSlice(a string, list []string) bool {
// 	for _, b := range list {
// 		if b == a {
// 			return true
// 		}
// 	}
// 	return false
// }
