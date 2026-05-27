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
package usr

import (
	"context"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/entity"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
)

type Plugin struct {
	collectors.Plugin
}

// Exec checks each username platform for the variant handle and records the
// platforms where it exists.
func (p *Plugin) Exec(ctx context.Context, domain *db.Domain) (variant *db.Domain, err error) {
	domain.Hits = append(domain.Hits, collectors.CheckSources(ctx, "username", domain.Name)...)
	return domain, nil
}

// Register the plugin
func init() {
	var CODE = "usr"
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{
			Plugin: collectors.Plugin{
				Num:      6,
				Code:     CODE,
				Title:    "Username Platforms",
				Summary:  "Check social/dev platforms for the username",
				Entities: []entity.Type{entity.User},
			},
		}
	})
}
