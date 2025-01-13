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
package list

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/pkg"
)

func (p *Plugin) Rows(domains ...*db.Domain) (rows []string) {
	for _, domain := range domains {
		rows = append(rows, p.Row(domain))
	}
	return
}

func (p *Plugin) Row(domain *db.Domain) (row string) {
	return p.Print(domain)
}

func (p *Plugin) Print(d *db.Domain) (output string) {
	tb := table.NewWriter()
	tb.SetStyle(pkg.StyleClear)

	// Basic domain info
	if d.Redirect != nil {
		tb.AppendRow(table.Row{d.Levenshtein, d.Algorithm.Code, d.Name, d.Punycode, d.Redirect.Name})
	} else {
		tb.AppendRow(table.Row{d.Levenshtein, d.Algorithm.Code, d.Name, d.Punycode})
	}

	if len(d.Dns) > 0 || len(d.Whois) > 0 {
		output += tb.Render() + "\n"
	} else {
		output += tb.Render()
	}

	// DNS Records
	tb = table.NewWriter()
	tb.SetStyle(pkg.StyleClear)
	for _, record := range d.Dns {
		tb.AppendRow(table.Row{"    ", record.Type, record.Value, record.Ttl})
	}
	if len(d.Dns) > 0 {
		output += tb.Render() + "\n"
	}

	// Whois Records
	tb = table.NewWriter()
	tb.SetStyle(pkg.StyleClear)
	for _, record := range d.Whois {
		if record.Registrant != nil {
			c := record.Registrant
			for _, field := range []string{c.Name, c.Organization, c.Email, c.Phone, c.PhoneExt} {
				if field != "" {
					tb.AppendRow(table.Row{"     ", field})
				}
			}
		}
		if record.Registrar != nil {
			c := record.Registrar
			for _, field := range []string{c.Name, c.Organization, c.Email, c.Phone, c.PhoneExt} {
				if field != "" {
					tb.AppendRow(table.Row{"     ", field})
				}
			}
		}
		if record.Administrative != nil {
			c := record.Administrative
			for _, field := range []string{c.Name, c.Organization, c.Email, c.Phone, c.PhoneExt} {
				if field != "" {
					tb.AppendRow(table.Row{"     ", field})
				}
			}
		}
		if record.Billing != nil {
			c := record.Billing
			for _, field := range []string{c.Name, c.Organization, c.Email, c.Phone, c.PhoneExt} {
				if field != "" {
					tb.AppendRow(table.Row{"     ", field})
				}
			}
		}
		if record.Technical != nil {
			c := record.Technical
			for _, field := range []string{c.Name, c.Organization, c.Email, c.Phone, c.PhoneExt} {
				if field != "" {
					tb.AppendRow(table.Row{"     ", field})
				}
			}
		}
	}
	if len(d.Whois) > 0 {
		tb.AppendRow(table.Row{""})
		output += tb.Render()
	}

	return
}
