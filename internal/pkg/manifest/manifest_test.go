// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package manifest

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// write creates name with content in a temp dir and returns its path.
func write(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func wantNames(t *testing.T, got []string, want ...string) {
	t.Helper()
	g := append([]string(nil), got...)
	sort.Strings(g)
	w := append([]string(nil), want...)
	sort.Strings(w)
	if len(g) != len(w) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range g {
		if g[i] != w[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestRequirements(t *testing.T) {
	path := write(t, "requirements.txt", `
# comment
requests==2.31.0
Flask>=2.0  # inline comment
numpy
django[argon2]>=4
-r other.txt
-e .
git+https://github.com/x/y.git
urllib3 ; python_version < "3.10"
`)
	names, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	wantNames(t, names, "requests", "Flask", "numpy", "django", "urllib3")
}

func TestPackageJSON(t *testing.T) {
	path := write(t, "package.json", `{
	"name": "demo",
	"dependencies": {"react": "^18.0.0", "lodash": "4.17.21"},
	"devDependencies": {"jest": "^29"}
}`)
	names, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	wantNames(t, names, "react", "lodash", "jest")
}

func TestComposerJSON(t *testing.T) {
	path := write(t, "composer.json", `{
	"require": {"php": ">=8.1", "ext-json": "*", "monolog/monolog": "^3.0"},
	"require-dev": {"phpunit/phpunit": "^10"}
}`)
	names, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	wantNames(t, names, "monolog/monolog", "phpunit/phpunit")
}

func TestGoMod(t *testing.T) {
	path := write(t, "go.mod", `module example.com/me

go 1.22

require (
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.9.0 // indirect
)

require golang.org/x/sync v0.7.0
`)
	names, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	wantNames(t, names,
		"github.com/spf13/cobra",
		"github.com/stretchr/testify",
		"golang.org/x/sync")
}

func TestCargoToml(t *testing.T) {
	path := write(t, "Cargo.toml", `[package]
name = "demo"

[dependencies]
serde = "1.0"
tokio = { version = "1", features = ["full"] }

[dependencies.reqwest]
version = "0.12"

[dev-dependencies]
criterion = "0.5"
`)
	names, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	wantNames(t, names, "serde", "tokio", "reqwest", "criterion")
}

func TestPyprojectPEP621(t *testing.T) {
	path := write(t, "pyproject.toml", `[project]
name = "demo"
dependencies = [
	"requests>=2.0",
	"rich",
]
`)
	names, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	wantNames(t, names, "requests", "rich")
}

func TestPyprojectPoetry(t *testing.T) {
	path := write(t, "pyproject.toml", `[tool.poetry.dependencies]
python = "^3.11"
httpx = "^0.27"

[tool.poetry.group.dev.dependencies]
pytest = "^8.0"
`)
	names, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	wantNames(t, names, "httpx", "pytest")
}

func TestGemfile(t *testing.T) {
	path := write(t, "Gemfile", `source "https://rubygems.org"

gem "rails", "~> 7.1"
gem 'puma'
# gem "commented"
`)
	names, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	wantNames(t, names, "rails", "puma")
}

func TestUnsupported(t *testing.T) {
	path := write(t, "weird.lock", "stuff")
	if _, err := Parse(path); err == nil {
		t.Fatal("expected error for unsupported manifest")
	}
}
