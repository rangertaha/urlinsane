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
	"strings"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/target"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/internal/plugins/algorithms/all"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
	_ "github.com/rangertaha/urlinsane/internal/plugins/outputs/all"
	"github.com/spf13/cobra"
)

type Config struct {
	target internal.Target // Config target
	ctype  int             // Config type

	// Plugins
	keyboards   []internal.Keyboard
	languages   []internal.Language
	algorithms  []internal.Algorithm
	information []internal.Information
	output      internal.Output

	// Performance
	concurrency int
	delay       time.Duration
	random      time.Duration
	levenshtein int

	// DNS
	dnsRetryCount       int
	dnsQueriesPerSecond int
	dnsConcurrency      int
	dnsServers          []string

	// Output
	verbose  bool
	format   string
	filters []string
	file     string
	showAll  bool
	scanAll  bool
	progress bool
}

func New() Config {
	return Config{}
}

func (c *Config) Target() internal.Target {
	return c.target
}
func (c *Config) Keyboards() []internal.Keyboard {
	return c.keyboards
}
func (c *Config) Languages() []internal.Language {
	return c.languages
}
func (c *Config) Algorithms() []internal.Algorithm {
	return c.algorithms
}
func (c *Config) Information() []internal.Information {
	return c.information
}
func (c *Config) Output() internal.Output {
	return c.output
}
func (c *Config) Concurrency() int {
	return c.concurrency
}
func (c *Config) Filters() (fields []string) {
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field != "" {
			fields = append(fields, field)
		}
	}
	return c.filters
}

// Dist is the Levenshtein_distance
// See: https://en.wikipedia.org/wiki/Levenshtein_distance
func (c *Config) Dist() int {
	return c.levenshtein
}
func (c *Config) Delay() time.Duration {
	return c.delay
}
func (c *Config) Random() time.Duration {
	return c.random
}
func (c *Config) Verbose() bool {
	return c.verbose
}
func (c *Config) Progress() bool {
	return c.progress
}
func (c *Config) Format() string {
	return c.format
}
func (c *Config) File() string {
	return c.file
}
func (c *Config) Type() int {
	return c.ctype
}
func (c *Config) ScanAll() bool {
	return c.scanAll
}
func (c *Config) ShowAll() bool {
	return c.showAll
}

func (c *Config) DnsServers() []string {
	if len(c.dnsServers) != 0 {
		return c.dnsServers
	}
	return []string{
		"127.0.0.53",
		"100.64.100.1",
		"8.8.8.8",         // Google
		"8.8.4.4",         // Google
		"1.1.1.1",         // Cloudflare
		"1.0.0.1",         // Cloudflare
		"1.1.1.2",         // Cloudflare  malware blocking
		"1.0.0.2",         // Cloudflare  malware blocking
		"1.1.1.3",         // Cloudflare  adult blocking
		"1.0.0.3",         // Cloudflare  adult blocking
		"9.9.9.9",         // Quad9
		"149.112.112.112", // Quad9
		"185.228.168.9",   // Cleanbrowsing
		"185.228.169.9",   // Cleanbrowsing
		"8.26.56.26",      // Comodo Secure DNS
		"8.20.247.20",     // Comodo Secure DNS
	}
}

// DnsQps is the queries per second
func (c *Config) DnsConcurrency() int {
	return c.dnsConcurrency
}

func (c *Config) DnsRetry() int {
	return c.dnsQueriesPerSecond
}

// DnsQps is the queries per second
func (c *Config) DnsQps() int {
	return c.dnsRetryCount
}

// CobraConfig creates a configuration from a cobra command options and arguments
func CobraConfig(cmd *cobra.Command, args []string, ttype int) (c Config, err error) {
	c.ctype = ttype

	c.target = target.New(args[0])

	// Target options
	// if name, err := cmd.Flags().GetString("name"); err == nil && name != "" {
	// 	c.ctype = internal.NAME
	// 	c.target = target.New(strings.TrimSpace(name))
	// }
	// if pkg, err := cmd.Flags().GetString("pkg"); err == nil && pkg != "" {
	// 	c.ctype = internal.PACKAGE
	// 	c.target = target.New(strings.TrimSpace(pkg))
	// }
	// if email, err := cmd.Flags().GetString("email"); err == nil && email != "" {
	// 	c.ctype = internal.EMAIL
	// 	c.target = target.New(strings.TrimSpace(email))
	// }
	// if domain, err := cmd.Flags().GetString("domain"); err == nil && domain != "" {
	// 	c.ctype = internal.DOMAIN
	// 	c.target = target.New(strings.TrimSpace(domain))
	// }
	if url, err := cmd.Flags().GetString("url"); err == nil && url != "" {
		if c.target != nil {
			c.target.Add("url", url)
		}
	}

	// Plugin options
	if langs, err := commaSplit(cmd.Flags().GetString("languages")); err == nil {
		c.languages = languages.Languages(langs...)
	}

	if keybs, err := commaSplit(cmd.Flags().GetString("keyboards")); err == nil {
		c.keyboards = languages.Keyboards(keybs...)
	}

	if typos, err := commaSplit(cmd.Flags().GetString("algorithms")); err == nil {
		c.algorithms = algorithms.List(typos...)
	}

	if infos, err := commaSplit(cmd.Flags().GetString("info")); err == nil {
		c.information = information.ListType(c.ctype, infos...)
	}

	// Output options
	if c.filters, err = commaSplit(cmd.Flags().GetString("filter")); err != nil {
		return c, err
	}

	if c.file, err = cmd.Flags().GetString("file"); err != nil {
		return c, err
	}

	if c.format, err = cmd.Flags().GetString("format"); err == nil {
		if c.output, err = outputs.Get(c.format); err != nil {
			return c, err
		}
	}

	if c.verbose, err = cmd.Flags().GetBool("verbose"); err != nil {
		return c, err
	}

	if c.progress, err = cmd.Flags().GetBool("progress"); err != nil {
		return c, err
	}

	// Performance and processing options
	if c.concurrency, err = cmd.Flags().GetInt("concurrency"); err != nil {
		return c, err
	}

	if c.random, err = cmd.Flags().GetDuration("random"); err != nil {
		return c, err
	}

	if c.delay, err = cmd.Flags().GetDuration("delay"); err != nil {
		return c, err
	}

	if c.levenshtein, err = cmd.Flags().GetInt("ld"); err != nil {
		return c, err
	}

	if c.showAll, err = cmd.Flags().GetBool("show"); err != nil {
		return c, err
	}

	if c.scanAll, err = cmd.Flags().GetBool("all"); err != nil {
		return c, err
	}

	if c.scanAll {
		c.levenshtein = 100
	}

	return c, err
}

// commaSplit splits comma separated values into an array
func commaSplit(value string, err error) ([]string, error) {
	return strings.Split(strings.TrimSpace(value), ","), err
}
