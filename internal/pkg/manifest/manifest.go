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

// Package manifest extracts the declared dependency names from common project
// manifests so urlinsane can typosquat-check the packages a project actually
// depends on — the highest-value supply-chain surface. Parsing is intentionally
// lightweight (no TOML/JSON-schema dependencies beyond encoding/json): it pulls
// out package names, not full version constraints.
package manifest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Parse reads a project manifest and returns the dependency package names it
// declares. The manifest format is detected from the file name.
func Parse(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	base := strings.ToLower(filepath.Base(path))
	var names []string
	switch {
	case base == "package.json":
		names, err = jsonDeps(data, "dependencies", "devDependencies",
			"peerDependencies", "optionalDependencies")
	case base == "composer.json":
		names, err = jsonDeps(data, "require", "require-dev")
	case base == "go.mod":
		names = goMod(data)
	case base == "cargo.toml":
		names = tomlSectionKeys(data, "dependencies", "dev-dependencies", "build-dependencies")
	case base == "pyproject.toml":
		names = pyproject(data)
	case base == "gemfile":
		names = gemfile(data)
	case strings.HasSuffix(base, ".txt"): // requirements.txt, requirements-dev.txt, ...
		names = requirements(data)
	default:
		return nil, fmt.Errorf("unsupported manifest %q (supported: requirements.txt, "+
			"package.json, go.mod, Cargo.toml, pyproject.toml, Gemfile, composer.json)", base)
	}
	if err != nil {
		return nil, err
	}
	return dedupe(names), nil
}

// requirementName extracts the bare package name from a PEP 508 / pip
// requirement spec, stopping at the first version operator, extra, marker, or
// whitespace (e.g. "requests[security]>=2.0 ; python_version<'3'" -> "requests").
func requirementName(spec string) string {
	spec = strings.TrimSpace(spec)
	idx := strings.IndexFunc(spec, func(r rune) bool {
		return strings.ContainsRune(" \t=<>!~;[(@", r)
	})
	if idx >= 0 {
		spec = spec[:idx]
	}
	return strings.TrimSpace(spec)
}

// requirements parses a pip requirements.txt file.
func requirements(data []byte) (names []string) {
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		// Skip blanks, pip options (-r, -e, --hash) and VCS/URL installs.
		if line == "" || strings.HasPrefix(line, "-") || strings.Contains(line, "://") {
			continue
		}
		if name := requirementName(line); name != "" {
			names = append(names, name)
		}
	}
	return
}

// jsonDeps collects the dependency keys from the named objects of a JSON
// manifest (package.json, composer.json).
func jsonDeps(data []byte, sections ...string) (names []string, err error) {
	var doc map[string]json.RawMessage
	if err = json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	for _, section := range sections {
		raw, ok := doc[section]
		if !ok {
			continue
		}
		var deps map[string]json.RawMessage
		if json.Unmarshal(raw, &deps) != nil {
			continue
		}
		for name := range deps {
			// composer pseudo-packages (php, ext-*, lib-*) are not registry names.
			if name == "php" || strings.HasPrefix(name, "ext-") || strings.HasPrefix(name, "lib-") {
				continue
			}
			names = append(names, name)
		}
	}
	return names, nil
}

// goMod extracts module paths from a go.mod require directive(s).
func goMod(data []byte) (names []string) {
	sc := bufio.NewScanner(bytes.NewReader(data))
	inBlock := false
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if i := strings.Index(line, "//"); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		switch {
		case line == "require (":
			inBlock = true
			continue
		case inBlock && line == ")":
			inBlock = false
			continue
		}
		fields := strings.Fields(line)
		if inBlock {
			if len(fields) >= 1 && strings.Contains(fields[0], ".") {
				names = append(names, fields[0])
			}
			continue
		}
		// Single-line: require module/path v1.2.3
		if len(fields) >= 2 && fields[0] == "require" && strings.Contains(fields[1], ".") {
			names = append(names, fields[1])
		}
	}
	return
}

// tomlSectionKeys returns the keys declared under the given TOML tables,
// covering both "name = ..." entries and "[section.name]" sub-tables. Good
// enough for Cargo.toml dependency tables without a full TOML parser.
func tomlSectionKeys(data []byte, sections ...string) (names []string) {
	want := make(map[string]bool, len(sections))
	for _, s := range sections {
		want[s] = true
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	current := ""
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			header := strings.TrimSpace(strings.Trim(line, "[]"))
			switch {
			case want[header]:
				// Whole table is a wanted section (e.g. [dependencies],
				// [tool.poetry.dependencies]); its keys are dependencies.
				current = header
			default:
				// [<section>.<name>] sub-table declares dependency <name>.
				current = ""
				if dot := strings.LastIndexByte(header, '.'); dot >= 0 && want[header[:dot]] {
					names = append(names, header[dot+1:])
				}
			}
			continue
		}
		if want[current] {
			if eq := strings.IndexByte(line, '='); eq > 0 {
				names = append(names, strings.TrimSpace(line[:eq]))
			}
		}
	}
	return
}

// pyproject extracts dependencies from both PEP 621 ([project] dependencies =
// [...]) and Poetry ([tool.poetry.dependencies]) layouts.
func pyproject(data []byte) (names []string) {
	// Poetry-style key tables.
	for _, n := range tomlSectionKeys(data, "tool.poetry.dependencies",
		"tool.poetry.dev-dependencies", "tool.poetry.group.dev.dependencies") {
		if n != "python" {
			names = append(names, n)
		}
	}

	// PEP 621 arrays of quoted PEP 508 strings inside a *dependencies block.
	sc := bufio.NewScanner(bytes.NewReader(data))
	inArray := false
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		if !inArray {
			if (strings.Contains(line, "dependencies") || strings.Contains(line, "requires")) &&
				strings.Contains(line, "[") {
				inArray = true
			}
		}
		if inArray {
			for _, q := range quotedStrings(line) {
				if name := requirementName(q); name != "" {
					names = append(names, name)
				}
			}
			if strings.Contains(line, "]") {
				inArray = false
			}
		}
	}
	return
}

// gemfile extracts gem names from a Ruby Gemfile (gem 'name', ...).
func gemfile(data []byte) (names []string) {
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if !strings.HasPrefix(line, "gem ") {
			continue
		}
		if q := quotedStrings(line); len(q) > 0 {
			names = append(names, q[0])
		}
	}
	return
}

// quotedStrings returns the contents of every single- or double-quoted span in s.
func quotedStrings(s string) (out []string) {
	for i := 0; i < len(s); i++ {
		if c := s[i]; c == '\'' || c == '"' {
			if j := strings.IndexByte(s[i+1:], c); j >= 0 {
				out = append(out, s[i+1:i+1+j])
				i += j + 1
			}
		}
	}
	return
}

// dedupe removes empty and duplicate names, preserving first-seen order.
func dedupe(in []string) (out []string) {
	seen := make(map[string]bool, len(in))
	for _, n := range in {
		n = strings.TrimSpace(n)
		if n == "" || seen[n] {
			continue
		}
		seen[n] = true
		out = append(out, n)
	}
	return
}
