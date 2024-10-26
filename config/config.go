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

	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/plugins/algorithms/all"
	"github.com/rangertaha/urlinsane/plugins/information"
	_ "github.com/rangertaha/urlinsane/plugins/information/all"
	"github.com/rangertaha/urlinsane/plugins/languages"
	_ "github.com/rangertaha/urlinsane/plugins/languages/all"
	"github.com/spf13/cobra"
)

type Config struct {
	Domain       urlinsane.Domain
	Keyboards    []urlinsane.Keyboard
	Languages    []urlinsane.Language
	Algorithms   []urlinsane.Algorithm
	Information []urlinsane.Information

	Headers     []string
	Format      string
	File        string
	Verbose     bool
	Concurrency int
	Delay       int
}

// CliConfig creates a configuration from a cobra cli options and arguments
func CliConfig(cmd *cobra.Command, args []string) (c Config, err error) {

	if c.Domain, err = getDomain(args); err != nil {
		return c, err
	}

	if langs, err := commaSplit(cmd.PersistentFlags().GetStringArray("languages")); err == nil {
		c.Languages = languages.Languages(langs...)
	}

	if keybs, err := commaSplit(cmd.PersistentFlags().GetStringArray("keyboards")); err == nil {
		c.Keyboards = languages.Keyboards(keybs...)
	}

	if typos, err := commaSplit(cmd.PersistentFlags().GetStringArray("typos")); err == nil {
		c.Algorithms = algorithms.List(typos...)
	}

	if infos, err := commaSplit(cmd.PersistentFlags().GetStringArray("info")); err == nil {
		c.Information = information.List(infos...)
	}

	if c.Concurrency, err = cmd.PersistentFlags().GetInt("concurrency"); err != nil {
		return c, err
	}

	// Output options
	if c.File, err = cmd.PersistentFlags().GetString("file"); err != nil {
		return c, err
	}

	if c.Format, err = cmd.PersistentFlags().GetString("format"); err != nil {
		return c, err
	}

	if c.Verbose, err = cmd.PersistentFlags().GetBool("verbose"); err != nil {
		return c, err
	}

	if c.Delay, err = cmd.PersistentFlags().GetInt("delay"); err != nil {
		return c, err
	}

	// Print logo
	fmt.Print(urlinsane.LOGO)

	return
}

// commaSplit splits comma seperated values into an array
func commaSplit(values []string, err error) ([]string, error) {
	return values, err
}
