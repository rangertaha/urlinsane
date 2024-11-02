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

	// "github.com/rangertaha/urlinsane/internal/plugins/information"
	// _ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
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

	// Output
	verbose  bool
	format   string
	file     string
	count    int64
	live     int64
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

// func (c *Config) IsMode(mode string) bool {
// 	return c.mode == mode
// }

// func (c *Config) IsDomain() string {
// 	return c.file
// }

// Count sets and gets the count of variants for processing
// func (c *Config) Count(n ...int64) int64 {
// 	if len(n) > 0 {
// 		c.count = n[0]
// 	}
// 	return c.count
// }
// func (c *Config) Live(n ...int64) int64 {
// 	if len(n) > 0 {
// 		c.live = n[0]
// 	}
// 	return c.count
// }

// CobraConfig creates a configuration from a cobra command options and arguments
func CobraConfig(cmd *cobra.Command) (c Config, err error) {

	if name, err := cmd.Flags().GetString("name"); err == nil && name != "" {
		c.ctype = internal.NAME
		c.target = target.New(strings.TrimSpace(name))
	}
	if pkg, err := cmd.Flags().GetString("pkg"); err == nil && pkg != "" {
		c.ctype = internal.PACKAGE
		c.target = target.New(strings.TrimSpace(pkg))
	}
	if domain, err := cmd.Flags().GetString("domain"); err == nil && domain != "" {
		c.ctype = internal.DOMAIN
		c.target = target.New(strings.TrimSpace(domain))
	}
	if url, err := cmd.Flags().GetString("url"); err == nil && url != "" {
		if c.target != nil {
			c.target.Add("url", url)
		}
	}

	// if c.name, err = cmd.Flags().GetString("pkg"); err == nil {
	// 	c.mode = internal.PACKAGE
	// }

	// if name, err := cmd.Flags().GetString("domain"); err == nil {
	// 	name = strings.TrimSpace(name)
	// 	c.domain = domain.New(name)
	// 	c.mode = internal.DOMAIN
	// }

	if langs, err := commaSplit(cmd.Flags().GetStringArray("languages")); err == nil {
		c.languages = languages.Languages(langs...)
	}

	if keybs, err := commaSplit(cmd.Flags().GetStringArray("keyboards")); err == nil {
		c.keyboards = languages.Keyboards(keybs...)
	}

	if typos, err := commaSplit(cmd.Flags().GetStringArray("algorithms")); err == nil {
		// c.algorithms = algorithms.List(typos...)
		c.algorithms = algorithms.List(typos...)
	}

	// if typos, err := commaSplit(cmd.Flags().GetStringArray("algorithms")); err == nil {
	// 	algos := algorithms.List(typos...)

	// 	for _, algo := range algos {
	// 		if c.IsMode(internal.DOMAIN) {
	// 			if al, ok := algo.(internal.DomainAlgo); ok {

	// 				c.algorithms = append(c.algorithms, al)
	// 			}
	// 		}

	// 	}
	// }

	// if infos, err := commaSplit(cmd.Flags().GetStringArray("info")); err == nil {
	// 	c.information = information.List(infos...)
	// }

	if c.format, err = cmd.Flags().GetString("format"); err == nil {
		if c.output, err = outputs.Get(c.format); err != nil {
			return c, err
		}
	}

	if c.concurrency, err = cmd.Flags().GetInt("concurrency"); err != nil {
		return c, err
	}

	// Output options
	if c.file, err = cmd.Flags().GetString("file"); err != nil {
		return c, err
	}

	if c.format, err = cmd.Flags().GetString("format"); err != nil {
		return c, err
	}

	if c.verbose, err = cmd.Flags().GetBool("verbose"); err != nil {
		return c, err
	}

	if c.progress, err = cmd.Flags().GetBool("progress"); err != nil {
		return c, err
	}

	if c.random, err = cmd.Flags().GetDuration("random"); err != nil {
		return c, err
	}

	if c.delay, err = cmd.Flags().GetDuration("delay"); err != nil {
		return c, err
	}
	c.count = 1

	return c, err
}

// // CobraConfig creates a configuration from a cobra command options and arguments
// func DomainConfig(cmd *cobra.Command, args []string) (c Config, err error) {
// 	c, err = CobraConfig(cmd, args)

// 	if infos, err := commaSplit(cmd.Flags().GetStringArray("info")); err == nil {
// 		c.information = domains.List(infos...)
// 	}lgos :=
// 	return
// }

// func UsernameConfig(cmd *cobra.Command, args []string) (c Config, err error) {
// 	c, err = CobraConfig(cmd, args)

// 	if infos, err := commaSplit(cmd.Flags().GetStringArray("info")); err == nil {
// 		c.information = usernames.List(infos...)
// 	}
// 	return
// }

// commaSplit splits comma separated values into an array
func commaSplit(values []string, err error) ([]string, error) {
	return values, err
}
