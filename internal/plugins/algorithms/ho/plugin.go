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
package ho

// Missing Dashes
//
// Func omits a dash from the name. For example:
//
// Original: www.one-two-three.com
//
// Variants: www.onetwo-three.com
//           www.one-twothree.com

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	algo "github.com/rangertaha/urlinsane/pkg/typo"
)

const (
	CODE        = "ho"
	NAME        = "Hyphen Omission"
	DESCRIPTION = "Created by removing hyphens from the name"
)

type Plugin struct {
	config internal.Config
}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Init(conf internal.Config) {
	n.config = conf
}

func (n *Plugin) Name() string {
	return NAME
}
func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (n *Plugin) Exec(original *db.Domain) (domains []*db.Domain, err error) {
	for _, variant := range algo.HyphenOmission(original.Name) {
		if original.Name != variant {
			domains = append(domains, &db.Domain{Name: variant})
			// acc.Add(domain.Variant(n, original.Prefix(), variant, original.Suffix()))
		}
	}

	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Plugin{}
	})
}
