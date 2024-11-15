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
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/internal/plugins/algorithms/all"
	"github.com/rangertaha/urlinsane/internal/plugins/analyzers"
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

const DIR_PRIMARY = ".config/urlinsane"

type (
	Config struct {
		domain string // Target domain

		directory string

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
		levenshtein  int
		summary      bool
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

	// Outputs to nowhere
	log.SetOutput(io.Discard)

	// Only log the errors
	log.SetLevel(log.ErrorLevel)
}

// New creates a new configuration
func New(options ...func(*Config)) (*Config, error) {
	s := &Config{
		format: "json",
	} // Default values

	// Apply config options
	for _, opt := range options {
		opt(s)
	}

	// Validate the domain name input
	if err := validateDomain(s); err != nil {
		return s, err
	}

	return s, nil
}

func CliOptions(cli *cli.Context) func(*Config) {
	var (

		// Basic input optoins
		domain     string   = cli.Args().First()                // Target domain
		languages  []string = csSplit(cli.String("languages"))  // Language IDs
		keyboards  []string = csSplit(cli.String("keyboards"))  // Keybards IDs
		algorithms []string = csSplit(cli.String("algorithms")) // algorithms IDs
		collectors []string = csSplit(cli.String("collectors")) // Collectors IDs
		analyzers  []string = csSplit(cli.String("analyzers"))  // Analyzers IDs
		database   string   = "badger"

		// Outputs options
		format       string = cli.String("format")     // Output format ID/Name
		file         string = cli.String("file")       //
		summary      bool   = cli.Bool("summary")      //
		registered   bool   = cli.Bool("registered")   //
		unregistered bool   = cli.Bool("unregistered") //
		verbose      bool   = cli.Bool("verbose")      //
		progress     bool   = cli.Bool("progress")     //
		debug        bool   = cli.Bool("debug")        //
		distance     int    = cli.Int("distance")      //
		banner       bool   = true                     //

		// Performance options
		workers int           = cli.Int("workers")      //
		random  time.Duration = cli.Duration("random")  //
		delay   time.Duration = cli.Duration("delay")   //
		ttl     time.Duration = cli.Duration("ttl")     //
		timeout time.Duration = cli.Duration("timeout") //
	)

	// Logs are disabled by default so we need to setup it up to log to stdout
	if debug {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
	}

	// We need to remove anything that is not json for output to work with
	// the 'jq' processor
	if strings.EqualFold(strings.TrimSpace(format), "json") {
		banner = false
		summary = false
		progress = false
	}

	return ConfigOption(
		domain,     // Target domain
		keyboards,  // Keybards IDs
		languages,  // Language IDs
		algorithms, // algorithms IDs
		collectors, // Collectors IDs
		analyzers,  // Analyzers IDs
		database,   // Database ID

		// Outputs options
		format, //
		file,
		summary,
		registered,
		unregistered,
		verbose,
		progress,
		banner,
		debug,
		distance,

		// Performance options
		workers,
		random,
		delay,
		ttl,
		timeout,
	)
}

func ConfigOption(
	domain string,
	boards []string,
	langs []string,
	algos []string,
	cols []string,
	anlyzrs []string,
	database string,

	// Outputs options
	format string,
	file string,
	summary bool,
	registered bool,
	unregistered bool,
	verbose bool,
	progress bool,
	banner bool,
	debug bool,
	distance int,

	// Performance options
	workers int,
	random time.Duration,
	delay time.Duration,
	ttl time.Duration,
	timeout time.Duration,

) func(*Config) {

	return func(c *Config) {
		var err error
		c.domain = domain
		c.languages = languages.Languages(langs...)
		c.keyboards = languages.Keyboards(boards...)
		c.algorithms = algorithms.List(algos...)

		c.collectors = collectors.List(cols...)
		sort.Sort(InfosReverseOrder{c.collectors})

		c.analyzers = analyzers.List(anlyzrs...)
		if c.database, err = databases.Get(database); err != nil {
			log.Error(err)
		}

		if c.output, err = outputs.Get(format); err != nil {
			log.Error(err)
		}

		// Outputs options
		c.format = format
		c.file = file
		c.summary = summary
		c.registered = registered
		c.unregistered = unregistered
		c.verbose = verbose
		c.progress = progress
		c.banner = banner
		c.levenshtein = distance

		// Performance options
		c.workers = workers
		c.random = random
		c.delay = delay
		c.ttl = ttl
		c.timeout = timeout

		// Create app directory if it does not exits
		c.directory = createAppDir(DIR_PRIMARY)

		logger := log.WithFields(log.Fields{
			"domain":     domain,
			"languages":  langs,
			"keyboards":  boards,
			"algorithms": algos,
			"collectors": cols,
			"analyzers":  anlyzrs,
			"database":   database,
			"format":     format,
			"file":       file,
			"summary":    summary,
			"registered": registered,
			"verbose":    verbose,
			"debug":      debug,
			"progress":   progress,
			"banner":     banner,
			"distance":   distance,
			"workers":    workers,
			"random":     random,
			"delay":      delay,
			"ttl":        ttl,
			"timeout":    timeout,
		})

		logger.Debug("Config options created")

	}
}

func validateDomain(cfg *Config) (err error) {
	name := strings.TrimSpace(cfg.domain)
	name = strings.Trim(name, ".")

	if strings.HasPrefix(name, "http") {
		u, err := url.Parse(name)
		if err != nil {
			return err
		}
		cfg.domain = u.Hostname()
	}

	return
}

func (c *Config) Target() string { return c.domain }

// Plugins options
func (c *Config) Keyboards() []internal.Keyboard   { return c.keyboards }
func (c *Config) Languages() []internal.Language   { return c.languages }
func (c *Config) Algorithms() []internal.Algorithm { return c.algorithms }
func (c *Config) Collectors() []internal.Collector { return c.collectors }
func (c *Config) Analyzers() []internal.Analyzer   { return c.analyzers }
func (c *Config) Output() internal.Output          { return c.output }
func (c *Config) Database() internal.Database      { return c.database }

// Performance options
func (c *Config) Workers() int           { return c.workers }
func (c *Config) Delay() time.Duration   { return c.delay }
func (c *Config) TTL() time.Duration     { return c.ttl }
func (c *Config) Random() time.Duration  { return c.random }
func (c *Config) Timeout() time.Duration { return c.timeout }

// Outputs options
func (c *Config) Filters() (fields []string) { return c.filters }
func (c *Config) Distance() int              { return c.levenshtein }
func (c *Config) Verbose() bool              { return c.verbose }
func (c *Config) Registered() bool           { return c.registered }
func (c *Config) Unregistered() bool         { return c.unregistered }
func (c *Config) Banner() bool               { return c.banner }
func (c *Config) Progress() bool             { return c.progress }
func (c *Config) Summary() bool              { return c.summary }
func (c *Config) Format() string             { return c.format }
func (c *Config) File() string               { return c.file }

func (c *Config) Dir() string { return c.directory }

// func (c *Config) Mkdir(name string) (dir string, err error) {
// 	dir = filepath.Join(c.directory, name)
// 	if err = os.MkdirAll(dir, 0750); err != nil {
// 		return
// 	}
// 	return
// }
// func (c *Config) Mkfile(dir, name string, content []byte) (file string, err error) {
// 	file = filepath.Join(dir, name)
// 	_, err = os.Stat(file)
// 	if os.IsNotExist(err) {
// 		err = os.WriteFile(file, content, 0644)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }

func csSplit(value string) []string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, " ", ",")
	value = strings.ReplaceAll(value, ".", ",")
	value = strings.ReplaceAll(value, "|", ",")
	value = strings.ReplaceAll(value, ":", ",")
	return strings.Split(value, ",")
}

func createAppDir(dirname string) string {
	var userDir string
	var err error

	if userDir, err = os.UserHomeDir(); err != nil {
		if userDir, err = os.Getwd(); err != nil {
			userDir = ""
		}
	}

	// If .config exits lets put it in there
	configDir := filepath.Join(userDir, dirname)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err = os.MkdirAll(configDir, 0750); err != nil {
			log.Error(err)
		}
	}

	return configDir
}

// // CliConfig creates a configuration from a cobra command options and arguments
// func CliConfig(cli *cli.Context) (c Config) {

// 	c = New()
// 	c.domain = cli.Args().First()

// 	// c.languages = languages.Languages(csSplit(cli.String("languages"))...)
// 	// c.keyboards = languages.Keyboards(csSplit(cli.String("keyboards"))...)
// 	// c.algorithms = algorithms.List(csSplit(cli.String("algorithms"))...)
// 	// c.collectors = collectors.List(csSplit(cli.String("collectors"))...)
// 	// if c.database, err = databases.Get("badger"); err != nil {
// 	// 	return c, err
// 	// }

// 	// if c.output, err = outputs.Get(cli.String("format")); err != nil {
// 	// 	return c, err
// 	// }

// 	c.file = cli.String("file")
// 	c.verbose = cli.Bool("verbose")
// 	c.progress = cli.Bool("progress")
// 	c.workers = cli.Int("workers")
// 	c.random = cli.Duration("random")
// 	c.delay = cli.Duration("delay")
// 	c.ttl = cli.Duration("ttl")
// 	c.timeout = cli.Duration("timeout")
// 	c.registered = cli.Bool("registered")
// 	c.unregistered = cli.Bool("unregistered")
// 	c.summary = cli.Bool("summary")
// 	c.banner = true
// 	if cli.Bool("debug") {
// 		log.SetOutput(os.Stdout)
// 		log.SetLevel(log.DebugLevel)
// 	}

// 	if cli.String("format") == "json" {
// 		c.banner = false
// 		c.summary = false
// 		c.progress = false
// 	}

// 	return c
// }

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

// func CreateAppConfig(dirname, filename string, defaultFile []byte) string {

// 	if err := os.MkdirAll(dirname, 0750); err != nil {
// 		log.Error(err)
// 	}

// 	fpath := filepath.Join(dirname, filename)

// 	// if err := Conf.Load(file.Provider(fpath), hcl.Parser(true)); err != nil {
// 	// 	log.Errorf("Unable to load config file: %s, Error: %s", fpath, err)

// 	// 	if _, err := os.Stat(fpath); os.IsNotExist(err) {
// 	// 		if err = os.WriteFile(fpath, defaultFile, 0644); err != nil {
// 	// 			log.Error(err)
// 	// 		}
// 	// 	}
// 	// }

// 	// conf.Print()

// 	return fpath
// }
