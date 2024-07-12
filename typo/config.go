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
	"fmt"
	"strings"
	"time"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/rangertaha/urlinsane/typo/languages"
	"github.com/spf13/cobra"
)

// BasicConfig ...
type BasicConfig struct {
	Domains     []string `json:"domains,omitempty"`
	Keyboards   []string `json:"keyboards,omitempty"`
	Filters     []string `json:"filters,omitempty"`
	Typos       []string `json:"typos,omitempty"`
	Funcs       []string `json:"funcs,omitempty"`
	Storage     []string `json:"storage,omitempty"`
	Concurrency int      `json:"concurrency,omitempty"`
	Format      string   `json:"format,omitempty"`
	File        string   `json:"file,omitempty"`
	Verbose     bool     `json:"verbose,omitempty"`
	Timing      Timing   `json:"timing,omitempty"`
}

// Timing ...
type Timing struct {
	Delay  time.Duration `json:"delay,omitempty"`
	Random time.Duration `json:"random,omitempty"`
}

// Config ...
type Config struct {
	domains   []Domain
	keyboards []languages.Keyboard
	languages []languages.Language

	typos   []Module
	funcs   []Module
	filters []Module
	// analyzers []Module

	headers     []string
	format      string
	file        string
	verbose     bool
	concurrency int
	timing      Timing
}

// NewConfig ...
func NewConfig(basic BasicConfig) (config *Config) {
	// Basic options
	config.GetDomains(basic.Domains)
	config.GetKeyboards(basic.Keyboards)

	// Registered functions
	config.GetTypos(basic.Typos)
	config.GetFuncs(basic.Funcs)
	config.GetFuncs(basic.Filters)

	// Processing option
	config.GetConcurrency(basic.Concurrency)

	// Output options
	config.GetFile(basic.File)
	config.GetFormat(basic.Format)
	config.GetVerbose(basic.Verbose)

	// Requires config.GetFuncs()
	config.GetHeaders(config.funcs)

	config.GetTiming(basic.Timing.Delay, basic.Timing.Random)

	return
}

// Config creates a Config
func (b *BasicConfig) Config() (c Config) {
	// Basic options
	c.GetDomains(b.Domains)
	c.GetKeyboards(b.Keyboards)

	// Registered functions
	c.GetTypos(b.Typos)
	c.GetFuncs(b.Funcs)
	c.GetFilters(b.Filters)

	// Processing option
	c.GetConcurrency(b.Concurrency)

	// Output options
	c.GetFile(b.File)
	c.GetFormat(b.Format)
	c.GetVerbose(b.Verbose)

	// Requires c.GetFuncs()
	c.GetHeaders(c.funcs)

	c.GetTiming(c.timing.Delay, c.timing.Random)

	return
}

// GetDomains ...
func (c *Config) GetDomains(args []string) {
	dmns := []Domain{}
	for _, str := range args {
		subdomain := domainutil.Subdomain(str)
		domain := domainutil.DomainPrefix(str)
		suffix := domainutil.DomainSuffix(str)
		if domain == "" {
			domain = str
		}
		dmns = append(dmns, Domain{
			Subdomain: subdomain,
			Domain:    domain,
			Suffix:    suffix})
	}
	c.domains = dmns
}

// GetKeyboards retrieves a list of keyboards
func (c *Config) GetKeyboards(keyboards []string) {
	c.keyboards = languages.KEYBOARDS.Keyboards(keyboards...)
}

// GetTypos ...
func (c *Config) GetTypos(typos []string) {
	c.typos = Typos.Get(typos...)
}

// GetFuncs ...
func (c *Config) GetFuncs(funcs []string) {
	if funcs := Extras.Get(funcs...); len(funcs) > 0 {
		c.funcs = funcs
	} else {
		c.funcs = Extras.Get("idna", "ld")
	}
}

// GetFilters ...
func (c *Config) GetFilters(filters []string) {
	if f := Filters.Get(filters...); len(filters) > 0 {
		c.filters = f
	}
}

// GetTiming ...
func (c *Config) GetTiming(delay, random time.Duration) {
	c.timing.Delay = delay
	c.timing.Random = random
}

// GetHeaders ...
func (c *Config) GetHeaders(funcs []Module) {
	c.headers = []string{"Live", "Type", "Typo", "Suffix"}
	for _, fnc := range funcs {
		for _, h := range fnc.Headers() {
			c.headers = appendIfMissing(c.headers, h)
		}
	}
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// GetConcurrency ...
func (c *Config) GetConcurrency(concurrency int) {
	c.concurrency = concurrency
}

// GetFile ...
func (c *Config) GetFile(file string) {
	c.file = file
}

// GetFormat ...
func (c *Config) GetFormat(format string) {
	c.format = format
}

// GetVerbose ...
func (c *Config) GetVerbose(verbose bool) {
	c.verbose = verbose
}

// errHandler
func errHandler(err error) {
	// TODO
}

// CobraConfig creates a configuration from a cobra command and arguments
func CobraConfig(cmd *cobra.Command, args []string) (c Config) {

	// Print logo
	fmt.Println(LOGO)

	// Basic options
	c.GetDomains(args)

	keyboards, err := cmd.PersistentFlags().GetStringArray("keyboards")
	errHandler(err)
	c.GetKeyboards(keyboards)

	// Registered functions
	var algorithms []string
	typos, err := cmd.PersistentFlags().GetStringArray("typos")
	for _, typo := range typos {
		algorithms = append(algorithms, strings.ToUpper(typo))
	}
	errHandler(err)
	c.GetTypos(algorithms)

	var funcs []string
	functions, err := cmd.PersistentFlags().GetStringArray("funcs")
	for _, function := range functions {
		funcs = append(funcs, strings.ToUpper(function))
	}
	errHandler(err)
	c.GetFuncs(funcs)

	var fltrs []string
	filters, err := cmd.PersistentFlags().GetStringArray("filters")
	for _, filter := range filters {
		fltrs = append(fltrs, strings.ToUpper(filter))
	}
	errHandler(err)
	c.GetFilters(fltrs)

	// Processing option
	concurrency, err := cmd.PersistentFlags().GetInt("concurrency")
	errHandler(err)
	c.GetConcurrency(concurrency)

	// Output options
	file, err := cmd.PersistentFlags().GetString("file")
	errHandler(err)
	c.GetFile(file)

	format, err := cmd.PersistentFlags().GetString("format")
	errHandler(err)
	c.GetFormat(format)

	verbose, err := cmd.PersistentFlags().GetBool("verbose")
	errHandler(err)
	c.GetVerbose(verbose)

	delay, err := cmd.PersistentFlags().GetInt64("delay")
	errHandler(err)
	rdelay, err := cmd.PersistentFlags().GetInt64("random-delay")
	errHandler(err)
	c.GetTiming(time.Duration(delay), time.Duration(rdelay))

	// Requires c.funcs
	c.GetHeaders(c.funcs)

	return
}
