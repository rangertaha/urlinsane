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
package config

import (
	"fmt"
	"os"

	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/languages"
	"github.com/spf13/cobra"
)

type Config struct {
	Domain      urlinsane.Domain
	Keyboards   []urlinsane.Keyboard
	Languages   []urlinsane.Language
	Algorithms  []urlinsane.Module
	Inforamtion []urlinsane.Module

	Headers     []string
	Format      string
	File        string
	Verbose     bool
	Concurrency int
	Delay       int
}

// CliConfig creates a configuration from a cobra command and arguments
func CliConfig(cmd *cobra.Command, args []string) (c Config, err error) {

	if len(args) == 0 {
		return c, fmt.Errorf("A domain is expected!")
	}
	domain = args[0]

	if c.domain, err = getDomain(args[0]); err != nil {
		return c, fmt.Errorf("A domain is expected!")
	}

	// if langs, err := CSplit(cmd.PersistentFlags().GetStringArray("languages")); err != nil {
	// 	c.languages = languages.List(langs...)
	// }

	// if keybs, err := CSplit(cmd.PersistentFlags().GetStringArray("keyboards")); err != nil {
	// 	c.keyboards = languages.KeyboardList(keybs...)
	// }

	// if typos, err := CSplit(cmd.PersistentFlags().GetStringArray("typos")); err != nil {
	// 	c.algorithms = algorithms.List(typos...)
	// }

	// if infos, err := CSplit(cmd.PersistentFlags().GetStringArray("infos")); err != nil {
	// 	c.inforamtion = information.List(infos...)

	// 	c.headers = []string{"Live", "Type", "Typo", "Suffix"}
	// 	for _, fnc := range c.inforamtion {
	// 		c.headers = append(c.headers, fnc.Headers()...)
	// 	}
	// }

	// // Basic options
	// c.GetDomains(args)

	// keyboards, err := CSplit(cmd.PersistentFlags().GetStringArray("keyboards"))
	// errHandler(err)
	// c.GetKeyboards(keyboards)

	// // Registered functions
	// var algorithms []string
	// typos, err := cmd.PersistentFlags().GetStringArray("typos")
	// for _, typo := range typos {
	// 	algorithms = append(algorithms, strings.ToUpper(typo))
	// }
	// errHandler(err)
	// c.GetTypos(algorithms)

	// var funcs []string
	// functions, err := cmd.PersistentFlags().GetStringArray("funcs")
	// for _, function := range functions {
	// 	funcs = append(funcs, strings.ToUpper(function))
	// }
	// errHandler(err)
	// c.GetFuncs(funcs)

	// var fltrs []string
	// filters, err := cmd.PersistentFlags().GetStringArray("filters")
	// for _, filter := range filters {
	// 	fltrs = append(fltrs, strings.ToUpper(filter))
	// }
	// errHandler(err)
	// c.GetFilters(fltrs)

	// // Processing option
	// if concurrency, err := cmd.PersistentFlags().GetInt("concurrency"); err != nil {
	// 	c.concurrency = concurrency
	// }
	// // errHandler(err)
	// // c.GetConcurrency(concurrency)

	// // Output options
	// if file, err := cmd.PersistentFlags().GetString("file"); err != nil {
	// 	c.file = file
	// }
	// // errHandler(err)
	// // c.GetFile(file)

	// if format, err := cmd.PersistentFlags().GetString("format"); err != nil {
	// 	c.format = format
	// }
	// // errHandler(err)
	// // c.GetFormat(format)

	// if verbose, err := cmd.PersistentFlags().GetBool("verbose"); err != nil {
	// 	c.verbose = verbose
	// }
	// // errHandler(err)
	// // c.GetVerbose(verbose)

	// if delay, err := cmd.PersistentFlags().GetInt("delay"); err != nil {
	// 	c.delay = delay
	// }
	// // errHandler(err)

	// errHandler(err)
	// c.GetTiming(time.Duration(delay), time.Duration(rdelay))

	// Requires c.funcs
	// c.GetHeaders(c.funcs)

	// c.Domains, _ = CSplit(args, nil)

	if langs, err := CSplit(cmd.PersistentFlags().GetStringArray("languages")); err == nil {
		c.Languages = languages.Languages(langs...)
	}

	// if c.Keyboards, err = CSplit(cmd.PersistentFlags().GetStringArray("keyboards")); err != nil {
	// 	errHandler(err)
	// }

	// if c.Algorithms, err = CSplit(cmd.PersistentFlags().GetStringArray("typos")); err != nil {
	// 	errHandler(err)
	// }

	// if c.Inforamtion, err = CSplit(cmd.PersistentFlags().GetStringArray("info")); err != nil {
	// 	errHandler(err)

	// }

	if c.Concurrency, err = cmd.PersistentFlags().GetInt("concurrency"); err != nil {
		errHandler(err)
	}

	// Output options
	if c.File, err = cmd.PersistentFlags().GetString("file"); err != nil {
		errHandler(err)
	}

	if c.Format, err = cmd.PersistentFlags().GetString("format"); err != nil {
		errHandler(err)
	}

	if c.Verbose, err = cmd.PersistentFlags().GetBool("verbose"); err != nil {
		errHandler(err)
	}

	if c.Delay, err = cmd.PersistentFlags().GetInt("delay"); err != nil {
		errHandler(err)
	}

	// // Create config from cli options/arguments
	// tool.CobraConfig(cmd, args)

	// // // Create a new instance of urlinsane
	// typosquating := tool.New(config)

	// // // Start generating results
	// typosquating.Execute()

	// Print logo
	fmt.Println(urlinsane.LOGO)

	return
}

// // NewConfig ...
// func NewConfig(basic BasicConfig) (config *Config) {
// 	// Basic options
// 	config.GetDomains(basic.Domains)
// 	config.GetKeyboards(basic.Keyboards)

// 	// Registered functions
// 	config.GetTypos(basic.Typos)
// 	config.GetFuncs(basic.Funcs)
// 	config.GetFuncs(basic.Filters)

// 	// Processing option
// 	config.GetConcurrency(basic.Concurrency)

// 	// Output options
// 	config.GetFile(basic.File)
// 	config.GetFormat(basic.Format)
// 	config.GetVerbose(basic.Verbose)

// 	// Requires config.GetFuncs()
// 	config.GetHeaders(config.infos)

// 	config.GetTiming(basic.Timing.Delay, basic.Timing.Random)

// 	return
// }

// // Config creates a Config
// func (b *BasicConfig) Config() (c Config) {
// 	// Basic options
// 	c.GetDomains(b.Domains)
// 	c.GetKeyboards(b.Keyboards)

// 	// Registered functions
// 	c.GetTypos(b.Typos)
// 	c.GetFuncs(b.Funcs)
// 	c.GetFilters(b.Filters)

// 	// Processing option
// 	c.GetConcurrency(b.Concurrency)

// 	// Output options
// 	c.GetFile(b.File)
// 	c.GetFormat(b.Format)
// 	c.GetVerbose(b.Verbose)

// 	// Requires c.GetFuncs()
// 	c.GetHeaders(c.infos)

// 	c.GetTiming(c.timing.Delay, c.timing.Random)

// 	return
// }

// // GetDomains ...
// func GetDomains(args []string) {
// 	dmns := []Domain{}
// 	for _, str := range args {
// 		subdomain := domainutil.Subdomain(str)
// 		domain := domainutil.DomainPrefix(str)
// 		suffix := domainutil.DomainSuffix(str)
// 		if domain == "" {
// 			domain = str
// 		}
// 		dmns = append(dmns, Domain{
// 			Subdomain: subdomain,
// 			Domain:    domain,
// 			Suffix:    suffix})
// 	}
// 	c.domains = dmns
// }

// // GetKeyboards retrieves a list of keyboards
// func (c *Config) GetKeyboards(keyboards []string) {
// 	c.keyboards = languages.Keyboards(keyboards...)
// }

// // GetTypos ...
// func (c *Config) GetTypos(typos []string) {
// 	c.typos = Typos.Get(typos...)
// }

// // GetFuncs ...
// func (c *Config) GetFuncs(funcs []string) {
// 	if funcs := Extras.Get(funcs...); len(funcs) > 0 {
// 		c.funcs = funcs
// 	} else {
// 		c.funcs = Extras.Get("idna", "ld")
// 	}
// }

// // GetFilters ...
// func (c *Config) GetFilters(filters []string) {
// 	if f := Filters.Get(filters...); len(filters) > 0 {
// 		c.filters = f
// 	}
// }

// // GetTiming ...
// func (c *Config) GetTiming(delay, random time.Duration) {
// 	c.timing.Delay = delay
// 	c.timing.Random = random
// }

// // GetHeaders ...
// func (c *Config) GetHeaders(funcs []Module) {
// 	c.headers = []string{"Live", "Type", "Typo", "Suffix"}
// 	for _, fnc := range funcs {
// 		for _, h := range fnc.Headers() {
// 			c.headers = appendIfMissing(c.headers, h)
// 		}
// 	}
// }

// func appendIfMissing(slice []string, i string) []string {
// 	for _, ele := range slice {
// 		if ele == i {
// 			return slice
// 		}
// 	}
// 	return append(slice, i)
// }

// // GetConcurrency ...
// func (c *Config) GetConcurrency(concurrency int) {
// 	c.concurrency = concurrency
// }

// // GetFile ...
// func (c *Config) GetFile(file string) {
// 	c.file = file
// }

// // GetFormat ...
// func (c *Config) GetFormat(format string) {
// 	c.format = format
// }

// // GetVerbose ...
// func (c *Config) GetVerbose(verbose bool) {
// 	c.verbose = verbose
// }

// // errHandler
// func errHandler(err error) {
// 	// TODO
// }

func errHandler(err error) {
	fmt.Println(err)
	os.Exit(0)
}

// CSplit splits comma seperated values into an array
func CSplit(values []string, err error) ([]string, error) {
	return values, err
}
