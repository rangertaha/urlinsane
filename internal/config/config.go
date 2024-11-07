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
	"path/filepath"
	"sort"
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

type (
	Config struct {
		target internal.Target // Config target
		ctype  int             // Config type
		appDir string

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

		screenShot bool

<<<<<<< HEAD
	// Output
	verbose  bool
	format   string
	filters []string
	file     string
	showAll  bool
	scanAll  bool
	progress bool
=======
		// DNS
		dnsRetryCount       int
		dnsQueriesPerSecond int
		dnsConcurrency      int
		dnsServers          []string

		// Output
		verbose  bool
		format   string
		filters  []string
		file     string
		showAll  bool
		scanAll  bool
		progress bool
	}

	Infos             []internal.Information
	InfosOrder        struct{ Infos }
	InfosReverseOrder struct{ Infos }
)

func (o Infos) Len() int      { return len(o) }
func (o Infos) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (l InfosReverseOrder) Less(i, j int) bool {
	return l.Infos[i].Order() > l.Infos[j].Order()
}
func (l InfosOrder) Less(i, j int) bool {
	return l.Infos[i].Order() < l.Infos[j].Order()
>>>>>>> develop
}

func New() Config {
	return Config{}
}

func (c *Config) Target() internal.Target {
	return c.target
}

func (c *Config) Dir() string {
	return filepath.Join(c.appDir, strings.ToLower(c.target.Name()))
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
	// for _, field := range fields {
	// 	field = strings.TrimSpace(field)
	// 	if field != "" {
	// 		fields = append(fields, field)
	// 	}
	// }
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
func (c *Config) Screenshot() bool {
	return c.screenShot
}
func (c *Config) DnsServers() []string {
	return c.dnsServers
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
	if err, appDir := NewAppDir(); err == nil {
		c = appDir.Config(c)
	}

	c.ctype = ttype

	c.target = target.New(args[0])

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
		sort.Sort(InfosOrder{c.information})
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

	if c.screenShot, err = cmd.Flags().GetBool("image"); err != nil {
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
