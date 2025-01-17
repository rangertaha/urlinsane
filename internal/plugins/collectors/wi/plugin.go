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
		db.DB.FirstOrCreate(&record, db.Whois{
			Created:    r.Domain.CreatedDateInTime,
			Updated:    r.Domain.UpdatedDateInTime,
			Expiration: r.Domain.ExpirationDateInTime,
		})
	}

	if r.Administrative != nil {
		administrative := db.Contact{
			Name:         r.Administrative.Name,
			Organization: r.Administrative.Organization,
			Street:       r.Administrative.Street,
			City:         r.Administrative.City,
			Province:     r.Administrative.Province,
			PostalCode:   r.Administrative.PostalCode,
			Country:      r.Administrative.Country,
			Phone:        r.Administrative.Phone,
			PhoneExt:     r.Administrative.PhoneExt,
			Fax:          r.Administrative.Fax,
			FaxExt:       r.Administrative.FaxExt,
			Email:        r.Administrative.Email,
			ReferralURL:  r.Administrative.ReferralURL,
		}
		db.DB.FirstOrInit(&administrative, administrative)
		record.Administrative = &administrative
	}

	if r.Billing != nil {
		billing := db.Contact{
			Name:         r.Billing.Name,
			Organization: r.Billing.Organization,
			Street:       r.Billing.Street,
			City:         r.Billing.City,
			Province:     r.Billing.Province,
			PostalCode:   r.Billing.PostalCode,
			Country:      r.Billing.Country,
			Phone:        r.Billing.Phone,
			PhoneExt:     r.Billing.PhoneExt,
			Fax:          r.Billing.Fax,
			FaxExt:       r.Billing.FaxExt,
			Email:        r.Billing.Email,
			ReferralURL:  r.Billing.ReferralURL,
		}
		db.DB.FirstOrInit(&billing, billing)
		record.Administrative = &billing
	}

	if r.Registrant != nil {
		registrant := db.Contact{
			Name:         r.Registrant.Name,
			Organization: r.Registrant.Organization,
			Street:       r.Registrant.Street,
			City:         r.Registrant.City,
			Province:     r.Registrant.Province,
			PostalCode:   r.Registrant.PostalCode,
			Country:      r.Registrant.Country,
			Phone:        r.Registrant.Phone,
			PhoneExt:     r.Registrant.PhoneExt,
			Fax:          r.Registrant.Fax,
			FaxExt:       r.Registrant.FaxExt,
			Email:        r.Registrant.Email,
			ReferralURL:  r.Registrant.ReferralURL,
		}
		db.DB.FirstOrInit(&registrant, registrant)
		record.Administrative = &registrant
	}

	if r.Technical != nil {
		technical := db.Contact{
			Name:         r.Technical.Name,
			Organization: r.Technical.Organization,
			Street:       r.Technical.Street,
			City:         r.Technical.City,
			Province:     r.Technical.Province,
			PostalCode:   r.Technical.PostalCode,
			Country:      r.Technical.Country,
			Phone:        r.Technical.Phone,
			PhoneExt:     r.Technical.PhoneExt,
			Fax:          r.Technical.Fax,
			FaxExt:       r.Technical.FaxExt,
			Email:        r.Technical.Email,
			ReferralURL:  r.Technical.ReferralURL,
		}
		db.DB.FirstOrInit(&technical, technical)
		record.Administrative = &technical
	}

	if r.Registrar != nil {
		registrar := &db.Contact{
			Name:         r.Registrar.Name,
			Organization: r.Registrar.Organization,
			Street:       r.Registrar.Street,
			City:         r.Registrar.City,
			Province:     r.Registrar.Province,
			PostalCode:   r.Registrar.PostalCode,
			Country:      r.Registrar.Country,
			Phone:        r.Registrar.Phone,
			PhoneExt:     r.Registrar.PhoneExt,
			Fax:          r.Registrar.Fax,
			FaxExt:       r.Registrar.FaxExt,
			Email:        r.Registrar.Email,
			ReferralURL:  r.Registrar.ReferralURL,
		}
		db.DB.FirstOrInit(&registrar, registrar)
		record.Administrative = registrar
	}

	domain.Whois = append(domain.Whois, record)
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
