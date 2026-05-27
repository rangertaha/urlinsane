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

// Package dc implements the dependency-confusion analyzer. After the package
// collector records which registries a name exists on, this analyzer flags the
// gap: a name claimed on some public registries but still available on others is
// squattable there — an attacker can publish a same-named package on the free
// registries and catch installs that resolve to the wrong ecosystem. It only
// fires when the pkg collector has run (Hits populated) and the name exists on
// some but not all registries.
package dc

import (
	"context"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/dataset"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/entity"
	"github.com/rangertaha/urlinsane/internal/plugins/analyzers"
)

const (
	CODE        = "dc"
	ORDER       = 2
	DESCRIPTION = "Dependency-confusion gap (squattable registries)"
)

type Plugin struct {
	// registries are the package-registry codes known to the dataset, loaded
	// once at Init.
	registries []string
}

func (p *Plugin) Id() string          { return CODE }
func (p *Plugin) Order() int          { return ORDER }
func (p *Plugin) Description() string { return DESCRIPTION }
func (p *Plugin) Headers() []string   { return []string{"CONFUSION"} }

// Types restricts the analyzer to package entities.
func (p *Plugin) Types() []entity.Type { return []entity.Type{entity.Package} }

// Init loads the set of package-registry codes from the dataset so the analyzer
// can compute which registries a name is missing from.
func (p *Plugin) Init(conf internal.Config) {
	if dataset.DB == nil {
		return
	}
	var sources []dataset.Source
	if err := dataset.DB.Where("type = ?", "package").Find(&sources).Error; err != nil {
		return
	}
	for _, s := range sources {
		p.registries = append(p.registries, s.Code)
	}
}

// Exec annotates the variant with the registries where it is still available
// when it already exists on some (but not all) of them.
func (p *Plugin) Exec(ctx context.Context, origin *db.Domain, variant *db.Domain) (*db.Domain, error) {
	if len(p.registries) == 0 {
		return variant, nil
	}

	found := make(map[string]bool, len(variant.Hits))
	for _, h := range variant.Hits {
		if h != nil {
			found[h.Service] = true
		}
	}

	var present, available []string
	for _, r := range p.registries {
		if found[r] {
			present = append(present, r)
		} else {
			available = append(available, r)
		}
	}

	// Dependency-confusion gap: claimed on some registries, free on others.
	if len(present) > 0 && len(available) > 0 {
		variant.Notes = append(variant.Notes, "squattable: "+strings.Join(available, ", "))
	}
	return variant, nil
}

// Register the plugin
func init() {
	analyzers.Add(CODE, func() internal.Analyzer {
		return &Plugin{}
	})
}
