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
package wi

import (
	"github.com/likexian/whois"
	parser "github.com/likexian/whois-parser"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
)

type Plugin struct {
	collectors.Plugin
}

func (p *Plugin) Exec(domain *db.Domain) (vaiant *db.Domain, err error) {
	raw, err := whois.Whois(domain.Name)
	if err != nil {
		p.Log.Error(err)
	}
	r, err := parser.Parse(raw)
	if err != nil {
		p.Log.Error(err)
	}
	p.Log.Debug(r)

	// type Whois struct {
	// 	Domain           *Domain    `json:"domain,omitempty"`
	// 	Registrar        *Contact   `json:"registrar,omitempty"`
	// 	Registrant       *Contact   `json:"registrant,omitempty"`
	// 	Administrative   *Contact   `json:"administrative,omitempty"`
	// 	Technical        *Contact   `json:"technical,omitempty"`
	// 	Billing          *Contact   `json:"billing,omitempty"`
	// 	Created          *time.Time `json:"created,omitempty"`
	// 	Updated          *time.Time `json:"updated,omitempty"`
	// 	Expiration       *time.Time `json:"expiration,omitempty"`
	// }

	whois := db.Whois{}
	db.DB.FirstOrInit(&whois, db.Whois{
		// Domain:     domain,
		Created:    r.Domain.CreatedDateInTime,
		Updated:    r.Domain.UpdatedDateInTime,
		Expiration: r.Domain.ExpirationDateInTime,
	})

	return domain, err
}

// Register the plugin
func init() {
	var CODE = "wi"
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{
			Plugin: collectors.Plugin{
				Num:       10,
				Code:      CODE,
				Title:     "Whois Lookup",
				Summary:   "Domain registration lookup",
				DependsOn: []string{"ip"},
			},
		}
	})
}
