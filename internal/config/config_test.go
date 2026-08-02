// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// A first run has no file: one is created and the defaults stand.
func TestLoadCreatesTheFileWhenMissing(t *testing.T) {
	dir := t.TempDir()
	Register("dns", map[string]any{"timeout": 5, "retries": 2})

	f, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ConfigFile)); err != nil {
		t.Errorf("config file was not created: %v", err)
	}
	got := f.Apply("dns")
	if got["timeout"] != 5 || got["retries"] != 2 {
		t.Errorf("defaults not applied: %v", got)
	}
}

// The file overrides a default, per key -- overriding one setting must not
// discard the others.
func TestFileOverridesDefaultsPerKey(t *testing.T) {
	dir := t.TempDir()
	Register("dns", map[string]any{"timeout": 5, "retries": 2})

	if err := os.WriteFile(filepath.Join(dir, ConfigFile),
		[]byte("plugins:\n  dns:\n    timeout: 30\n"), 0o640); err != nil {
		t.Fatal(err)
	}

	f, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	got := f.Apply("dns")
	if got["timeout"] != 30 {
		t.Errorf("timeout = %v, want the file's 30", got["timeout"])
	}
	if got["retries"] != 2 {
		t.Errorf("retries = %v, want the default 2 to survive", got["retries"])
	}
	// Defaults themselves must not be mutated by an override.
	if d := Defaults("dns"); d["timeout"] != 5 {
		t.Errorf("Apply mutated the registered defaults: %v", d)
	}
}

// A plugin with no section gets its defaults; an unregistered one gets nothing
// rather than an error.
func TestApplyWithoutSectionOrRegistration(t *testing.T) {
	dir := t.TempDir()
	Register("whois", map[string]any{"timeout": 10})

	f, _ := Load(dir)
	if got := f.Apply("whois"); got["timeout"] != 10 {
		t.Errorf("unmentioned plugin lost its defaults: %v", got)
	}
	if got := f.Apply("nosuchplugin"); len(got) != 0 {
		t.Errorf("unregistered plugin returned %v, want empty", got)
	}
}

// A section for a plugin this build does not know must survive a round trip:
// it may belong to a plugin that is simply not loaded.
func TestUnknownSectionsSurviveSave(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ConfigFile),
		[]byte("plugins:\n  someotherplugin:\n    key: value\n"), 0o640); err != nil {
		t.Fatal(err)
	}

	f, err := Load(dir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	f.Set("dns", "timeout", 30)
	if err := f.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	again, err := Load(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if again.Plugins["someotherplugin"]["key"] != "value" {
		t.Error("an unknown plugin's settings were dropped on save")
	}
	if again.Plugins["dns"]["timeout"] != 30 {
		t.Error("the override was not persisted")
	}
}

// Seed lists what can be configured without overwriting what is already set.
func TestSeedAddsMissingSectionsOnly(t *testing.T) {
	dir := t.TempDir()
	Register("dns", map[string]any{"timeout": 5})
	Register("http", map[string]any{"agent": "urlinsane"})

	f, _ := Load(dir)
	f.Set("dns", "timeout", 99)
	f.Seed()

	if f.Plugins["dns"]["timeout"] != 99 {
		t.Errorf("Seed overwrote an existing value: %v", f.Plugins["dns"])
	}
	if f.Plugins["http"]["agent"] != "urlinsane" {
		t.Errorf("Seed did not add the missing section: %v", f.Plugins)
	}
}

// A malformed file is an error, not a reason to overwrite what the user wrote.
func TestMalformedFileIsNotOverwritten(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ConfigFile)
	bad := []byte("plugins: [this is not a mapping\n")
	if err := os.WriteFile(path, bad, 0o640); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(dir); err == nil {
		t.Fatal("Load accepted a malformed file")
	}
	got, _ := os.ReadFile(path)
	if string(got) != string(bad) {
		t.Error("a malformed config was overwritten instead of reported")
	}
}
