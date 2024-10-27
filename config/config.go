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
	"strings"
	"time"

	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/plugins/algorithms/all"
	"github.com/rangertaha/urlinsane/plugins/information"
	_ "github.com/rangertaha/urlinsane/plugins/information/all"
	"github.com/rangertaha/urlinsane/plugins/languages"
	_ "github.com/rangertaha/urlinsane/plugins/languages/all"
	"github.com/rangertaha/urlinsane/plugins/outputs"
	_ "github.com/rangertaha/urlinsane/plugins/outputs/all"
	"github.com/spf13/cobra"
)

type Config struct {
	// Types of targets for typosquattting
	// Domain urlinsane.Domain
	target string

	// Plugins
	keyboards   []urlinsane.Keyboard
	languages   []urlinsane.Language
	algorithms  []urlinsane.Algorithm
	information []urlinsane.Information
	output      urlinsane.Output

	// Performance
	concurrency int
	delay       time.Duration
	random      time.Duration

	// Output
	verbose bool
	format  string
	file    string
	count   int64
}

func (c *Config) Target() string {
	return c.target
}
func (c *Config) Keyboards() []urlinsane.Keyboard {
	return c.keyboards
}
func (c *Config) Languages() []urlinsane.Language {
	return c.languages
}
func (c *Config) Algorithms() []urlinsane.Algorithm {
	return c.algorithms
}
func (c *Config) Information() []urlinsane.Information {
	return c.information
}
func (c *Config) Output() urlinsane.Output {
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
func (c *Config) Format() string {
	return c.format
}
func (c *Config) File() string {
	return c.file
}
func (c *Config) Count(n ...int64) int64 {
	if len(n) > 0 {
		c.count = n[0]
	}
	return c.count
}

// CliDomainConfig creates a configuration from a cobra cli options and arguments
func CobraConfig(cmd *cobra.Command, args []string) (c Config, err error) {

	if len(args) == 0 {
		return c, fmt.Errorf("at least one argument required")
	}
	if len(args) > 1 {
		return c, fmt.Errorf("only one argument is allowed")
	}
	c.target = strings.TrimSpace(args[0])

	if langs, err := commaSplit(cmd.PersistentFlags().GetStringArray("languages")); err == nil {
		c.languages = languages.Languages(langs...)
	}

	if keybs, err := commaSplit(cmd.PersistentFlags().GetStringArray("keyboards")); err == nil {
		c.keyboards = languages.Keyboards(keybs...)
	}

	if typos, err := commaSplit(cmd.PersistentFlags().GetStringArray("typos")); err == nil {
		c.algorithms = algorithms.List(typos...)
	}

	if infos, err := commaSplit(cmd.PersistentFlags().GetStringArray("info")); err == nil {
		c.information = information.List(infos...)
	}

	if c.format, err = cmd.PersistentFlags().GetString("format"); err == nil {
		if c.output, err = outputs.Get(c.format); err != nil {
			return c, err
		}
	}

	if c.concurrency, err = cmd.PersistentFlags().GetInt("concurrency"); err != nil {
		return c, err
	}

	// Output options
	if c.file, err = cmd.PersistentFlags().GetString("file"); err != nil {
		return c, err
	}

	if c.format, err = cmd.PersistentFlags().GetString("format"); err != nil {
		return c, err
	}

	if c.verbose, err = cmd.PersistentFlags().GetBool("verbose"); err != nil {
		return c, err
	}

	if c.random, err = cmd.PersistentFlags().GetDuration("random"); err != nil {
		return c, err
	}

	if c.delay, err = cmd.PersistentFlags().GetDuration("delay"); err != nil {
		return c, err
	}
	c.count = 1

	return c, err
}

// commaSplit splits comma seperated values into an array
func commaSplit(values []string, err error) ([]string, error) {
	return values, err
}
