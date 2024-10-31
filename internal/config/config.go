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

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/internal/plugins/algorithms/all"

	// "github.com/rangertaha/urlinsane/internal/plugins/information"
	// _ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
	_ "github.com/rangertaha/urlinsane/internal/plugins/information/all"
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
	"github.com/rangertaha/urlinsane/internal/plugins/information/emails"
	"github.com/rangertaha/urlinsane/internal/plugins/information/packages"
	"github.com/rangertaha/urlinsane/internal/plugins/information/usernames"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
	_ "github.com/rangertaha/urlinsane/internal/plugins/outputs/all"
	"github.com/spf13/cobra"
)

type Config struct {
	// Types of targets for typosquatting
	// Domain internal.Domain
	target string

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
	progress bool
}

func (c *Config) Target() string {
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

// Count sets and gets the count of variants for processing
func (c *Config) Count(n ...int64) int64 {
	if len(n) > 0 {
		c.count = n[0]
	}
	return c.count
}

// CobraConfig creates a configuration from a cobra command options and arguments
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

	// if infos, err := commaSplit(cmd.Flags().GetStringArray("info")); err == nil {
	// 	c.information = information.List(infos...)
	// }

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

	if c.progress, err = cmd.PersistentFlags().GetBool("progress"); err != nil {
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

// CobraConfig creates a configuration from a cobra command options and arguments
func DomainConfig(cmd *cobra.Command, args []string) (c Config, err error) {
	c, err = CobraConfig(cmd, args)

	if infos, err := commaSplit(cmd.Flags().GetStringArray("info")); err == nil {
		c.information = domains.List(infos...)
	}
	return
}

func EmailConfig(cmd *cobra.Command, args []string) (c Config, err error) {
	c, err = CobraConfig(cmd, args)

	if infos, err := commaSplit(cmd.Flags().GetStringArray("info")); err == nil {
		c.information = emails.List(infos...)
	}
	return
}

func UsernameConfig(cmd *cobra.Command, args []string) (c Config, err error) {
	c, err = CobraConfig(cmd, args)

	if infos, err := commaSplit(cmd.Flags().GetStringArray("info")); err == nil {
		c.information = usernames.List(infos...)
	}
	return
}

func PackageConfig(cmd *cobra.Command, args []string) (c Config, err error) {
	c, err = CobraConfig(cmd, args)

	if infos, err := commaSplit(cmd.Flags().GetStringArray("info")); err == nil {
		c.information = packages.List(infos...)
	}
	return
}

// commaSplit splits comma separated values into an array
func commaSplit(values []string, err error) ([]string, error) {
	return values, err
}
