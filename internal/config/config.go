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
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/internal/plugins/algorithms/all"
	_ "github.com/rangertaha/urlinsane/internal/plugins/analyzers/all"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/all"
	"github.com/rangertaha/urlinsane/internal/plugins/databases"
	_ "github.com/rangertaha/urlinsane/internal/plugins/databases/all"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
	_ "github.com/rangertaha/urlinsane/internal/plugins/outputs/all"
)

const DIR = ".urlinsane"
const FILE = "urlinsane.hcl"

var (
	Conf = koanf.New(".")

	DefualtConfig = []byte(`
database = "badger"
algorithms = "all"
concurrency = 25
delay = 1 
format = "table"
keyboards = "all"
languages = "all"
progress = false
random = 1
verbose = false
banner = true
`)
)

type (
	Dns struct {
		RetryCount       int      `yaml:"retry"`
		QueriesPerSecond int      `yaml:"qps"`
		Concurrency      int      `yaml:"concurrency"`
		Servers          []string `yaml:"servers"`
	}

	Config struct {
		domain string // Target domain

		// Plugins
		keyboards  []internal.Keyboard
		languages  []internal.Language
		algorithms []internal.Algorithm
		collectors []internal.Collector
		database   internal.Database
		analyzers  []internal.Analyzer
		output     internal.Output

		// Performance
		concurrency int           `yaml:"concurrency"`
		delay       time.Duration `yaml:"delay"`
		random      time.Duration `yaml:"random"`
		levenshtein int           `yaml:"levenshtein"`

		// DNS
		Dns Dns `yaml:"dns"`

		// Output
		verbose bool     `yaml:"verbose"`
		banner  bool     `yaml:"verbose"`
		format  string   `yaml:"format"`
		filters []string `yaml:"filters"`
		file    string   `yaml:"file"`
		// showAll  bool     `yaml:"showAll"`
		// scanAll  bool     `yaml:"scanAll"`
		progress bool `yaml:"progress"`
	}

	Infos             []internal.Collector
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
}

func New() Config {
	return Config{}
}

func (c *Config) Target() string {
	return c.domain
}

// func (c *Config) Dir() *AppDir {
// 	return c.appDir
// }

// Plugins
func (c *Config) Keyboards() []internal.Keyboard {
	return c.keyboards
}
func (c *Config) Languages() []internal.Language {
	return c.languages
}
func (c *Config) Algorithms() []internal.Algorithm {
	return c.algorithms
}
func (c *Config) Collectors() []internal.Collector {
	sort.Sort(InfosReverseOrder{c.collectors})
	return c.collectors
}
func (c *Config) Analyzers() []internal.Analyzer {
	return c.analyzers
}
func (c *Config) Output() internal.Output {
	return c.output
}
func (c *Config) Database() internal.Database {
	return c.database
}

// Dist is the Levenshtein_distance
// See: https://en.wikipedia.org/wiki/Levenshtein_distance
// func (c *Config) Dist() int {
// 	return c.levenshtein
// }

// Performance
func (c *Config) Concurrency() int {
	return c.concurrency
}

func (c *Config) Delay() time.Duration {
	return c.delay
}
func (c *Config) Random() time.Duration {
	return c.random
}

// Outputs
func (c *Config) Filters() (fields []string) {
	return c.filters
}

func (c *Config) Verbose() bool {
	return c.verbose
}

func (c *Config) Banner() bool {
	return c.banner
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

func (c *Config) BadgerOptions() badger.Options {
	return badger.DefaultOptions(Conf.String("database.file"))
}

// CobraConfig creates a configuration from a cobra command options and arguments
func CliConfig(target string) (c Config, err error) {
	c.domain = target

	c.languages = languages.Languages(csSplit(Conf.String("languages"))...)
	c.keyboards = languages.Keyboards(csSplit(Conf.String("keyboards"))...)
	c.algorithms = algorithms.List(csSplit(Conf.String("algorithms"))...)
	c.collectors = collectors.List(csSplit(Conf.String("collectors"))...)
	// c.languages = languages.Languages(csSplit(Conf.String("languages"))...)
	if c.database, err = databases.Get(Conf.String("database")); err != nil {
		return c, err
	}

	if c.output, err = outputs.Get(Conf.String("format")); err != nil {
		return c, err
	}

	c.file = Conf.String("file")
	c.verbose = Conf.Bool("verbose")
	c.progress = Conf.Bool("progress")
	c.concurrency = Conf.Int("concurrency")
	c.random = Conf.Duration("random")
	c.delay = Conf.Duration("delay")
	c.banner = Conf.Bool("banner")

	return c, err
}

func csSplit(value string) []string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, " ", ",")
	return strings.Split(value, ",")
}

func LoadOrCreateConfig(conf *koanf.Koanf, dirname, filename string, defaultFile []byte) (err error) {
	var userDir string

	if userDir, err = os.UserHomeDir(); err != nil {
		if userDir, err = os.Getwd(); err != nil {
			userDir = ""
		}
	}
	configDir := filepath.Join(userDir, dirname)
	configFile := filepath.Join(configDir, filename)

	databaseDir := filepath.Join(userDir, DIR, "badger")

	if err := conf.Load(file.Provider(configFile), hcl.Parser(true)); err != nil {
		log.Printf("Unable to load config file: %s, Error: %s", configFile, err)
		if err = os.MkdirAll(databaseDir, 0750); err != nil {
			return err
		}
		_, err := os.Stat(configFile)
		if os.IsNotExist(err) {
			err = os.WriteFile(configFile, defaultFile, 0644)
			if err != nil {
				panic(err)
			}
		}
	}

	// conf.Print()

	return
}
