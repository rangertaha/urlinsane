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
	"context"
	"github.com/likexian/whois"
	parser "github.com/likexian/whois-parser"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/entity"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
)

type Plugin struct {
	collectors.Plugin
}

func (p *Plugin) Exec(ctx context.Context, domain *db.Domain) (variant *db.Domain, err error) {
	raw, err := whois.Whois(domain.Name)
	p.Log.Debug(raw)
	if err != nil {
		p.Log.Error(err)
	}

	r, err := parser.Parse(raw)
	p.Log.Debug(r)
	if err != nil {
		p.Log.Error(err)
	}

	record := db.Whois{}
	if r.Domain != nil {
		record.Created = r.Domain.CreatedDateInTime
		record.Updated = r.Domain.UpdatedDateInTime
		record.Expiration = r.Domain.ExpirationDateInTime
	}
	// Assign each role to its own field (the previous code assigned every role
	// to Administrative). Results are mutated in place; persistence happens at
	// the store boundary.
	record.Administrative = contact(r.Administrative)
	record.Billing = contact(r.Billing)
	record.Registrant = contact(r.Registrant)
	record.Technical = contact(r.Technical)
	record.Registrar = contact(r.Registrar)

	domain.Whois = append(domain.Whois, record)
	return domain, err
}

// contact converts a parsed whois contact into a db.Contact (nil-safe).
func contact(c *parser.Contact) *db.Contact {
	if c == nil {
		return nil
	}
	return &db.Contact{
		Name:         c.Name,
		Organization: c.Organization,
		Street:       c.Street,
		City:         c.City,
		Province:     c.Province,
		PostalCode:   c.PostalCode,
		Country:      c.Country,
		Phone:        c.Phone,
		PhoneExt:     c.PhoneExt,
		Fax:          c.Fax,
		FaxExt:       c.FaxExt,
		Email:        c.Email,
		ReferralURL:  c.ReferralURL,
	}
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
				Entities:  []entity.Type{entity.Domain},
			},
		}
	})
}
