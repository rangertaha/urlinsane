// Copyright 2024 Rangertaha. All Rights Reserved.
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
	_ "embed"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rangertaha/urlinsane/internal/dataset"
	"github.com/rangertaha/urlinsane/internal/db"
	"gorm.io/gorm"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/entity"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	_ "github.com/rangertaha/urlinsane/internal/plugins/algorithms/all"
	"github.com/rangertaha/urlinsane/internal/plugins/analyzers"
	_ "github.com/rangertaha/urlinsane/internal/plugins/analyzers/all"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
	_ "github.com/rangertaha/urlinsane/internal/plugins/collectors/all"
	"github.com/rangertaha/urlinsane/internal/plugins/languages"
	_ "github.com/rangertaha/urlinsane/internal/plugins/languages/all"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
	_ "github.com/rangertaha/urlinsane/internal/plugins/outputs/all"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	DIR_PRIMARY = ".config/urlinsane"
	PRIMARY_DB  = "urlinsane.db"
	DATASET_DB  = "dataset.db"
	MAXMIND_DB  = "maxmind.db.gz"
)

//go:embed maxmind.db.gz
var MaxMindDB []byte

//go:embed dataset.db
var datasetDB []byte

type (
	Config struct {
		domain     string      // Target value (domain, username, package, ...)
		entityType entity.Type // Kind of target being analyzed
		directory  string
		database  *gorm.DB
		dataset   *gorm.DB

		// Plugins
		keyboards  []internal.Keyboard
		languages  []internal.Language
		algorithms []internal.Algorithm
		collectors []internal.Collector
		analyzers  []internal.Analyzer
		output     internal.Output

		// Constraints
		regex       string
		levenshtein int

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
		summary      bool
		assets       string

		// Metrics
		total int
	}
)

func init() {
	// Outputs to nowhere
	log.SetOutput(io.Discard)

	// Only log the errors
	log.SetLevel(log.ErrorLevel)
}

// New creates a new configuration
func New(options ...func(*Config)) (*Config, error) {
	s := &Config{
		format: "list",
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
		domain     string   = cli.Args().First()                // Target value
		languages  []string = csSplit(cli.String("languages"))  // Language IDs
		keyboards  []string = csSplit(cli.String("keyboards"))  // Keybards IDs
		algorithms []string = csSplit(cli.String("algorithms")) // algorithms IDs
		collectors []string = csSplit(cli.String("collectors")) // Collectors IDs
		analyzers  []string = csSplit(cli.String("analyzers"))  // Analyzers IDs

		// Constraints
		regex    string = cli.String("regex") //
		distance int    = cli.Int("distance") //

		// Outputs options
		format       string = cli.String("format")     // Output format ID/Name
		file         string = cli.String("file")       //
		summary      bool   = cli.Bool("summary")      //
		registered   bool   = cli.Bool("registered")   //
		unregistered bool   = cli.Bool("unregistered") //
		verbose      bool   = cli.Bool("verbose")      //
		progress     bool   = cli.Bool("progress")     //
		debug        bool   = cli.Bool("debug")        //
		assets       string = cli.Path("dir")          //
		banner       bool   = true                     //

		// Performance options
		workers int           = cli.Int("workers")      //
		random  time.Duration = cli.Duration("random")  //
		delay   time.Duration = cli.Duration("delay")   //
		ttl     time.Duration = cli.Duration("ttl")     //
		timeout time.Duration = cli.Duration("timeout") //
	)

	// Resolve the entity type, defaulting to domain for unknown/empty values.
	entityType, ok := entity.Parse(strings.ToLower(strings.TrimSpace(cli.String("type"))))
	if !ok {
		entityType = entity.Domain
	}

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

	// if ttl == 0 {
	// 	deleteCacheDir(DIR_PRIMARY)
	// }

	return ConfigOption(
		domain,     // Target value
		entityType, // Kind of target
		keyboards,  // Keybards IDs
		languages,  // Language IDs
		algorithms, // algorithms IDs
		collectors, // Collectors IDs
		analyzers,  // Analyzers IDs

		// Constraints
		regex,
		distance,

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
		assets,

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
	entityType entity.Type,
	boards []string,
	langs []string,
	algos []string,
	cols []string,
	anlyzrs []string,

	// Constraints
	regex string,
	distance int,

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
	assets string,

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
		c.entityType = entityType
		c.languages = languages.Languages(langs...)
		c.keyboards = languages.Keyboards(boards...)

		// Keep only the algorithms/collectors that apply to this entity type
		// (e.g. TLD swaps and DNS collectors apply to domains only).
		c.algorithms = algosForType(algorithms.List(algos...), entityType)

		// Collector execution order is determined by the dependency DAG at
		// run time (internal/engine/dag), not by a static sort here.
		c.collectors = collectorsForType(collectors.List(cols...), entityType)

		c.analyzers = analyzers.List(anlyzrs...)

		if c.output, err = outputs.Get(format); err != nil {
			log.Error(err)
		}

		// Constraints
		c.regex = regex
		c.levenshtein = distance

		// Outputs options
		c.format = format
		c.file = file
		c.summary = summary
		c.registered = registered
		c.unregistered = unregistered
		c.verbose = verbose
		c.progress = progress
		c.banner = banner
		c.assets = assets

		// Performance options
		c.workers = workers
		c.random = random
		c.delay = delay
		c.ttl = ttl
		c.timeout = timeout

		// Create app directory if it does not exits
		c.directory = createAppDir(DIR_PRIMARY)

		// Create app database if it does not exits
		c.database = createDatabase(c.directory)

		// Create app database if it does not exits
		c.dataset = createDatasets(c.directory)

		// Create app database if it does not exits
		createMaxMindDB(c.directory)

		logger := log.WithFields(log.Fields{
			"domain":     domain,
			"languages":  langs,
			"keyboards":  boards,
			"algorithms": algos,
			"collectors": cols,
			"analyzers":  anlyzrs,
			"regex":      regex,
			"format":     format,
			"file":       file,
			"summary":    summary,
			"save":       assets,
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

// EntityType returns the kind of target being analyzed (defaults to domain).
func (c *Config) EntityType() entity.Type {
	if c.entityType == "" {
		return entity.Domain
	}
	return c.entityType
}

// algosForType keeps only the algorithms that apply to the given entity type.
func algosForType(in []internal.Algorithm, t entity.Type) (out []internal.Algorithm) {
	for _, a := range in {
		if entity.Supports(a.Types(), t) {
			out = append(out, a)
		}
	}
	return
}

// collectorsForType keeps only the collectors that apply to the given entity type.
func collectorsForType(in []internal.Collector, t entity.Type) (out []internal.Collector) {
	for _, c := range in {
		if entity.Supports(c.Types(), t) {
			out = append(out, c)
		}
	}
	return
}

// Plugins options
func (c *Config) Keyboards() []internal.Keyboard   { return c.keyboards }
func (c *Config) Languages() []internal.Language   { return c.languages }
func (c *Config) Algorithms() []internal.Algorithm { return c.algorithms }
func (c *Config) Collectors() []internal.Collector { return c.collectors }
func (c *Config) Analyzers() []internal.Analyzer   { return c.analyzers }
func (c *Config) Output() internal.Output          { return c.output }
func (c *Config) Database() *gorm.DB               { return c.database }
func (c *Config) Dataset() *gorm.DB                { return c.dataset }

// Constraint options
func (c *Config) Regex() string { return c.regex }
func (c *Config) Distance() int { return c.levenshtein }

// Performance options
func (c *Config) Workers() int           { return c.workers }
func (c *Config) Delay() time.Duration   { return c.delay }
func (c *Config) TTL() time.Duration     { return c.ttl }
func (c *Config) Random() time.Duration  { return c.random }
func (c *Config) Timeout() time.Duration { return c.timeout }

// Outputs options
func (c *Config) Filters() (fields []string) { return c.filters }
func (c *Config) Verbose() bool              { return c.verbose }
func (c *Config) Registered() bool           { return c.registered }
func (c *Config) Unregistered() bool         { return c.unregistered }
func (c *Config) Banner() bool               { return c.banner }
func (c *Config) Progress() bool             { return c.progress }
func (c *Config) Summary() bool              { return c.summary }
func (c *Config) Format() string             { return c.format }
func (c *Config) File() string               { return c.file }
func (c *Config) AssetDir() string           { return c.assets }

func (c *Config) Dir() string { return c.directory }
func (c *Config) Count(v ...int) int {
	if len(v) > 0 {
		c.total = v[0]
	}

	return c.total
}
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

func createDatabase(dirname string) *gorm.DB {
	db.Config(filepath.Join(dirname, PRIMARY_DB))
	return db.DB
}

func createDatasets(dirname string) *gorm.DB {
	datasetPath := filepath.Join(dirname, DATASET_DB)

	if _, err := os.Stat(datasetPath); os.IsNotExist(err) {
		if err := os.WriteFile(datasetPath, datasetDB, 0666); err != nil {
			log.Fatal(err)
		}
	}

	dataset.Config(datasetPath)
	return dataset.DB
}

func createMaxMindDB(dirname string) {
	datasetPath := filepath.Join(dirname, MAXMIND_DB)

	if _, err := os.Stat(datasetPath); os.IsNotExist(err) {
		if err := os.WriteFile(datasetPath, MaxMindDB, 0666); err != nil {
			log.Fatal(err)
		}
	}
}
