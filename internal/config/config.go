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
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rangertaha/urlinsane/internal/dataset"
	"gorm.io/gorm"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/entity"
	"github.com/rangertaha/urlinsane/internal/pkg/dns"
	"github.com/rangertaha/urlinsane/internal/pkg/manifest"
	"github.com/rangertaha/urlinsane/internal/store"

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
		targets    []string    // All targets to scan (manifest deps, or [domain])
		entityType entity.Type // Kind of target being analyzed
		directory  string
		dataset    *gorm.DB
		store      *store.Store // content-addressed result store (IPLD)

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

	// Fail fast on an unknown output format (otherwise a nil output plugin
	// would later panic in the banner / output stage). An empty format means no
	// output was requested (e.g. the dataset import tooling), which is fine.
	if s.format != "" && s.output == nil {
		return s, fmt.Errorf("unknown output format %q (supported: list, json)", s.format)
	}

	return s, nil
}

func CliOptions(cli *cli.Context) func(*Config) {
	var (

		// Basic input options
		domain     string   = cli.Args().First()                // Target value
		languages  []string = csSplit(cli.String("languages"))  // Language IDs
		keyboards  []string = csSplit(cli.String("keyboards"))  // Keyboards IDs
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

	// Resolve the entity type. "auto" (the default) and any unrecognized value
	// classify the target from its shape; an explicit known type is honored.
	entityType := entity.Classify(domain)
	if t, ok := entity.Parse(strings.ToLower(strings.TrimSpace(cli.String("type")))); ok {
		entityType = t
	}

	// --manifest scans the dependency names declared by a project manifest
	// (requirements.txt, package.json, go.mod, ...). The targets become the
	// parsed package names and the entity type is forced to package.
	var manifestTargets []string
	if mpath := cli.String("manifest"); mpath != "" {
		names, err := manifest.Parse(mpath)
		if err != nil {
			log.Errorf("manifest: %s", err)
		} else {
			manifestTargets = names
			entityType = entity.Package
			if strings.TrimSpace(domain) == "" {
				domain = filepath.Base(mpath) // display / scan name
			}
		}
	}

	// Point the DNS collectors at custom servers if --nameservers is set.
	dns.SetResolver(csSplit(cli.String("nameservers")))

	// Logs are discarded by default. --debug and --verbose enable them on
	// stderr (not stdout, which carries scan results).
	if debug {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.DebugLevel)
	} else if verbose {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.InfoLevel)
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

	base := ConfigOption(
		domain,     // Target value
		entityType, // Kind of target
		keyboards,  // Keyboards IDs
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

	// Carry the manifest-derived target list (if any) onto the config.
	return func(c *Config) {
		base(c)
		c.targets = manifestTargets
	}
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

		c.analyzers = analyzersForType(analyzers.List(anlyzrs...), entityType)

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

		// Scan results are stored in the content-addressed store (opened below),
		// not a primary GORM database.

		// Create the dataset (read-only reference data) if it does not exist
		c.dataset = createDatasets(c.directory)

		// Open the content-addressed result store (IPLD blockstore + index)
		if c.store, err = store.OpenDir(c.directory); err != nil {
			log.Error("opening result store: ", err)
		}

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

// Targets returns every entity to scan: the manifest dependency names when
// --manifest was given, otherwise the single Target().
func (c *Config) Targets() []string {
	if len(c.targets) > 0 {
		return c.targets
	}
	return []string{c.domain}
}

// Store returns the content-addressed result store.
func (c *Config) Store() *store.Store { return c.store }

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

// analyzersForType keeps only the analyzers that apply to the given entity type.
func analyzersForType(in []internal.Analyzer, t entity.Type) (out []internal.Analyzer) {
	for _, a := range in {
		if entity.Supports(a.Types(), t) {
			out = append(out, a)
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
