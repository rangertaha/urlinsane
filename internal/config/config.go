// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package config owns the application directory: where it is, what lives in it,
// and putting the shipped data there on first run.
//
// It used to own the run configuration too — a thirty-field struct, a
// twenty-five-parameter builder, and a CliOptions that read twenty-five flags
// no surviving command registered, so every lookup returned the zero value. The
// graph engine takes its configuration as scan.Options built directly from the
// command line, so none of that had a caller. What remains is the part that
// never had an alternative.
//
// Setup is reported rather than silent. An extraction that fails leaves the
// operators that needed the data out of the compiled plan, and a scan missing
// geolocation must not look like a target with no geolocation (§12.6).
//
// # Plugin settings
//
// The same directory holds config.yaml, and this package owns that too. A
// plugin declares its defaults in code and never has to be configured; the file
// exists only to override them, so a setting appears in two places — the
// default that ships, and the value the user chose — and the file is the one
// that wins.
//
//	Register("dns", map[string]any{"timeout": 5, "retries": 2})
//
//	f, _ := Load(dir)                 // ~/.config/urlinsane/config.yaml
//	conf := f.Apply("dns")            // defaults, with file entries winning
//
// Unknown keys in the file are kept rather than dropped. A setting belonging to
// a plugin not registered in this build is not necessarily a mistake — it may
// belong to a plugin that is simply not loaded — and silently discarding it on
// the next Save would lose the user's configuration.
package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/rangertaha/urlinsane/internal/dataset"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

const (
	// DirName is the application directory, under the user's home.
	DirName = ".config/urlinsane"
	// DatasetDB is the reference data: vocabulary and weighted transitions.
	DatasetDB = "dataset.db"
	// MaxMindDB is the geolocation database.
	MaxMindDB = "maxmind.db.gz"
	// ConfigFile is the plugin settings file. Load and Settings, below,
	// are what read and write it.
	ConfigFile = "config.yaml"
)

//go:embed maxmind.db.gz
var maxmindDB []byte

//go:embed dataset.db
var datasetDB []byte

// File is what became of one shipped file on this run.
type File struct {
	Path string
	// Written is true when the file was absent and has just been extracted,
	// which is what makes a run the *first* run.
	Written bool
	// Err is why the file is unusable. Non-nil does not stop a scan: the
	// operators that needed it are left out of the compiled plan instead (§4),
	// and this is what lets the CLI say which, and why.
	Err error
}

// Setup is what Init did, for the caller to report.
type Setup struct {
	Dir      string
	Created  bool
	Dataset  File
	GeoIP    File
	Settings File
}

// FirstRun reports whether the setup is worth showing the user: something had
// to be created, or something could not be.
//
// The failure half is not decoration. A shipped file that never validates is
// wrong on every run, not only the first, and reporting it once would leave
// every subsequent run silently missing an operator — which is the failure this
// package exists to make loud (§12.6).
func (s Setup) FirstRun() bool {
	return s.Created || s.Dataset.Written || s.GeoIP.Written ||
		s.Dataset.Err != nil || s.GeoIP.Err != nil
}

// Dir returns the application directory, creating it if absent, and whether it
// had to be created.
//
// It resolves under the user's home, falling back to the working directory — a
// tool that cannot find a home directory should still run, with its data beside
// it, rather than refuse to start.
func Dir() (dir string, created bool, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		if home, err = os.Getwd(); err != nil {
			home = ""
		}
	}
	dir = filepath.Join(home, DirName)

	if _, err := os.Stat(dir); err == nil {
		return dir, false, nil
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return dir, false, fmt.Errorf("config: creating %s: %w", dir, err)
	}
	return dir, true, nil
}

// Init prepares the application directory and opens the dataset database.
//
// A failure to extract one file is recorded in the Setup rather than returned:
// losing geolocation should cost the geo operator, not the scan. Only a
// directory that cannot be created is fatal, because nothing else can proceed
// without it.
func Init() (Setup, error) {
	dir, created, err := Dir()
	if err != nil {
		return Setup{Dir: dir}, err
	}
	s := Setup{Dir: dir, Created: created}

	s.Dataset = extract(dir, DatasetDB, datasetDB, isSQLite)
	s.GeoIP = extract(dir, MaxMindDB, maxmindDB, isGzip)
	s.Settings = File{Path: filepath.Join(dir, ConfigFile)}

	// Opening is separate from extracting: a file that is present but corrupt
	// is dataset.Config's problem, and it self-heals by recreating it empty.
	if s.Dataset.Err == nil {
		dataset.Config(s.Dataset.Path)
	}
	return s, nil
}

// validator reports whether shipped bytes are the format their name claims.
//
// It reads the header and nothing else. The point is not to prove a database is
// intact — that costs a full parse of tens of megabytes on every run — but to
// refuse the one failure that has actually happened here: a binary asset round
// tripped through a text decoder, which destroys every byte >= 0x80 and so
// always destroys a magic number.
type validator func([]byte) error

// isGzip checks the two-byte gzip magic and the deflate method that follows it.
func isGzip(b []byte) error {
	if len(b) < 3 {
		return fmt.Errorf("truncated: %d bytes", len(b))
	}
	if b[0] != 0x1f || b[1] != 0x8b {
		return fmt.Errorf("not gzip data: magic is % x, want 1f 8b", b[:min(4, len(b))])
	}
	if b[2] != 0x08 {
		return fmt.Errorf("gzip method is %#02x, want 08 (deflate)", b[2])
	}
	return nil
}

// sqliteMagic is the header every SQLite 3 database begins with.
var sqliteMagic = []byte("SQLite format 3\x00")

func isSQLite(b []byte) error {
	if len(b) < len(sqliteMagic) {
		return fmt.Errorf("truncated: %d bytes", len(b))
	}
	if !bytes.Equal(b[:len(sqliteMagic)], sqliteMagic) {
		return fmt.Errorf("not a SQLite database: header is %q", b[:len(sqliteMagic)])
	}
	return nil
}

// extract writes a shipped file if it is not already there, refusing bytes that
// are not the format they claim to be.
//
// Validating before the write rather than trusting the embed is what stops the
// tool reporting "extracted" for a file it has just made unusable. The shipped
// maxmind.db.gz is exactly that case: its gzip magic reads 1f ef bf bd, the 8b
// replaced by the UTF-8 replacement character, so 49 MB was written out on every
// fresh install and then failed to open — and deleting it, the obvious remedy,
// re-extracted the same corruption.
func extract(dir, name string, data []byte, valid validator) File {
	f := File{Path: filepath.Join(dir, name)}
	if _, err := os.Stat(f.Path); err == nil {
		return f
	}
	if len(data) == 0 {
		f.Err = fmt.Errorf("config: %s was not compiled into this binary", name)
		return f
	}
	if valid != nil {
		if err := valid(data); err != nil {
			f.Err = fmt.Errorf("config: the embedded %s is corrupt (%w); "+
				"replace internal/config/%s and rebuild", name, err, name)
			return f
		}
	}
	if err := os.WriteFile(f.Path, data, 0o640); err != nil {
		f.Err = fmt.Errorf("config: writing %s: %w", f.Path, err)
		return f
	}
	f.Written = true
	return f
}

var (
	mu       sync.RWMutex
	defaults = map[string]map[string]any{}
)

// Register records a plugin's default settings. Called from a plugin's init,
// alongside the registration of the plugin itself.
//
// Registering the same id twice replaces the earlier defaults: a plugin owns
// its own settings, and two plugins sharing an id is a bug that would show up
// long before this.
func Register(id string, values map[string]any) {
	mu.Lock()
	defer mu.Unlock()
	copied := make(map[string]any, len(values))
	for k, v := range values {
		copied[k] = v
	}
	defaults[id] = copied
}

// Defaults returns the registered defaults for a plugin, or nil.
func Defaults(id string) map[string]any {
	mu.RLock()
	defer mu.RUnlock()
	return clone(defaults[id])
}

// Registered lists the plugin ids that have declared defaults, sorted.
func Registered() []string {
	mu.RLock()
	defer mu.RUnlock()
	ids := make([]string, 0, len(defaults))
	for id := range defaults {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Settings is the parsed config file: one section per plugin id.
type Settings struct {
	// Plugins is keyed by plugin id, each holding that plugin's overrides.
	Plugins map[string]map[string]any `yaml:"plugins"`

	path string
}

// Load reads the config file from an application config directory, creating an
// empty one if it does not exist.
//
// A missing file is not an error -- it is the state of every first run -- so it
// is created and an empty configuration returned. A malformed one *is* an
// error: overwriting a file the user has edited, because this build could not
// parse it, would throw away their settings.
func Load(dir string) (*Settings, error) {
	path := filepath.Join(dir, ConfigFile)
	f := &Settings{Plugins: map[string]map[string]any{}, path: path}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return f, f.Save()
	}
	if err != nil {
		return nil, fmt.Errorf("settings: reading %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, f); err != nil {
		return nil, fmt.Errorf("settings: parsing %s: %w", path, err)
	}
	if f.Plugins == nil {
		f.Plugins = map[string]map[string]any{}
	}
	return f, nil
}

// Save writes the file back to where it was loaded from.
func (f *Settings) Save() error {
	if f.path == "" {
		return fmt.Errorf("settings: no path to save to")
	}
	if err := os.MkdirAll(filepath.Dir(f.path), 0o750); err != nil {
		return err
	}
	data, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, data, 0o640)
}

// Path is where the file was loaded from.
func (f *Settings) Path() string { return f.path }

// Apply returns a plugin's effective settings: its registered defaults, with
// anything the file sets for that plugin overriding them.
//
// Per key rather than per section, so overriding one setting does not mean
// restating the rest.
func (f *Settings) Apply(id string) map[string]any {
	out := Defaults(id)
	if out == nil {
		out = map[string]any{}
	}
	if f == nil {
		return out
	}
	for k, v := range f.Plugins[id] {
		out[k] = v
	}
	return out
}

// Set records an override for one setting and leaves the rest alone.
func (f *Settings) Set(id, key string, value any) {
	if f.Plugins == nil {
		f.Plugins = map[string]map[string]any{}
	}
	if f.Plugins[id] == nil {
		f.Plugins[id] = map[string]any{}
	}
	f.Plugins[id][key] = value
}

// Seed fills in a section for every registered plugin that the file does not
// mention, so a freshly written config lists what can be configured instead of
// leaving the user to discover it.
//
// Existing entries are never touched, and sections belonging to unregistered
// plugins are left in place.
func (f *Settings) Seed() {
	if f.Plugins == nil {
		f.Plugins = map[string]map[string]any{}
	}
	for _, id := range Registered() {
		if _, ok := f.Plugins[id]; !ok {
			f.Plugins[id] = Defaults(id)
		}
	}
}

func clone(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
