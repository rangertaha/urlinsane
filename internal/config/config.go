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
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const DIR_PRIMARY = ".urlinsane"
const DIR = ".urlinsane"
const FILE = "urlinsane.hcl"

var (
	Conf = koanf.New(".")

	DefualtConfig = []byte(`
database = "badger"
algorithms = "all"
workers = 25
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
	HTTPConfig struct {
	}

	DNSConfig struct {
		RetryCount       int
		QueriesPerSecond int
		Concurrency      int
		Servers          []string
	}

	Config struct {
		domain string // Target domain

		directory  string
		configfile string

		// Plugins
		keyboards  []internal.Keyboard
		languages  []internal.Language
		algorithms []internal.Algorithm
		collectors []internal.Collector
		database   internal.Database
		analyzers  []internal.Analyzer
		output     internal.Output

		// Performance
		workers int
		delay   time.Duration
		random  time.Duration
		timeout time.Duration
		ttl     time.Duration

		// Output
		verbose      bool
		banner       bool
		format       string
		filters      []string
		file         string
		progress     bool
		registered   bool
		unregistered bool

		// DNS
		// dns  DNSConfig
		// http HTTPConfig
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

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(io.Discard)

	// Only log the warning severity or above.
	log.SetLevel(log.ErrorLevel)
}

func New() Config {
	cdir := CreateAppDir(DIR)
	cfile := CreateAppConfig(cdir, FILE, DefualtConfig)
	return Config{directory: cdir, configfile: cfile}
}

func (c *Config) Target() string {
	return c.domain
}

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

// Performance
func (c *Config) Workers() int {
	return c.workers
}

func (c *Config) Delay() time.Duration {
	return c.delay
}

func (c *Config) TTL() time.Duration {
	return c.ttl
}

func (c *Config) Random() time.Duration {
	return c.random
}

func (c *Config) Timeout() time.Duration {
	return c.timeout
}

// Outputs
func (c *Config) Filters() (fields []string) {
	return c.filters
}

func (c *Config) Verbose() bool {
	return c.verbose
}

func (c *Config) Registered() bool {
	return c.registered
}
func (c *Config) Unregistered() bool {
	return c.unregistered
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

func (c *Config) Dir() string {
	return c.directory
}

func (c *Config) Mkdir(name string) (dir string, err error) {
	dir = filepath.Join(c.directory, name)
	if err = os.MkdirAll(dir, 0750); err != nil {
		return
	}
	return
}
func (c *Config) Mkfile(dir, name string, content []byte) (file string, err error) {
	file = filepath.Join(dir, name)
	_, err = os.Stat(file)
	if os.IsNotExist(err) {
		err = os.WriteFile(file, content, 0644)
		if err != nil {
			return
		}
	}
	return
}

// CliConfig creates a configuration from a cobra command options and arguments
func CliConfig(cli *cli.Context) (c Config, err error) {

	c = New()
	c.domain = cli.Args().First()

	c.languages = languages.Languages(csSplit(cli.String("languages"))...)
	c.keyboards = languages.Keyboards(csSplit(cli.String("keyboards"))...)
	c.algorithms = algorithms.List(csSplit(cli.String("algorithms"))...)
	c.collectors = collectors.List(csSplit(cli.String("collectors"))...)
	if c.database, err = databases.Get("badger"); err != nil {
		return c, err
	}

	if c.output, err = outputs.Get(cli.String("format")); err != nil {
		return c, err
	}

	c.file = cli.String("file")
	c.verbose = cli.Bool("verbose")
	c.progress = cli.Bool("progress")
	c.workers = cli.Int("workers")
	c.random = cli.Duration("random")
	c.delay = cli.Duration("delay")
	c.ttl = cli.Duration("ttl")
	c.timeout = cli.Duration("timeout")
	c.registered = cli.Bool("registered")
	c.unregistered = cli.Bool("unregistered")
	c.banner = true
	if cli.Bool("debug") {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		// Conf.Print()
	}

	return c, err
}

func csSplit(value string) []string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, " ", ",")
	return strings.Split(value, ",")
}

// func LoadOrCreateConfig(conf *koanf.Koanf, dirname, filename string, defaultFile []byte) (err error) {
// 	var userDir string

// 	if userDir, err = os.UserHomeDir(); err != nil {
// 		if userDir, err = os.Getwd(); err != nil {
// 			userDir = ""
// 		}
// 	}
// 	configDir := filepath.Join(userDir, dirname)
// 	configFile := filepath.Join(configDir, filename)

// 	if err := conf.Load(file.Provider(configFile), hcl.Parser(true)); err != nil {
// 		log.Printf("Unable to load config file: %s, Error: %s", configFile, err)
// 		if err = os.MkdirAll(configDir, 0750); err != nil {
// 			return err
// 		}
// 		_, err := os.Stat(configFile)
// 		if os.IsNotExist(err) {
// 			err = os.WriteFile(configFile, defaultFile, 0644)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}

// 	// conf.Print()

// 	return
// }

func CreateAppConfig(dirname, filename string, defaultFile []byte) string {

	if err := os.MkdirAll(dirname, 0750); err != nil {
		log.Error(err)
	}

	fpath := filepath.Join(dirname, filename)

	// if err := Conf.Load(file.Provider(fpath), hcl.Parser(true)); err != nil {
	// 	log.Errorf("Unable to load config file: %s, Error: %s", fpath, err)

	// 	if _, err := os.Stat(fpath); os.IsNotExist(err) {
	// 		if err = os.WriteFile(fpath, defaultFile, 0644); err != nil {
	// 			log.Error(err)
	// 		}
	// 	}
	// }

	// conf.Print()

	return fpath
}

func CreateAppDir(dirname string) string {
	var userDir string
	var err error

	if userDir, err = os.UserHomeDir(); err != nil {
		if userDir, err = os.Getwd(); err != nil {
			userDir = ""
		}
	}

	// If .config exits lets put it in there
	configDir := filepath.Join(userDir, ".config")
	if _, err := os.Stat(configDir); !os.IsNotExist(err) {
		configDir = filepath.Join(configDir, dirname)
	} else {
		configDir = filepath.Join(userDir, dirname)
	}

	if err = os.MkdirAll(configDir, 0750); err != nil {
		log.Error(err)
	}

	return configDir
}
