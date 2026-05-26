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

package store

import (
	"sort"
	"time"

	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/entity"
)

// ToEntity converts an in-flight pipeline *db.Domain into the content-addressed
// Entity form. Result slices are normalized (sorted) so that re-scans of
// unchanged data produce identical CIDs. Pipeline-only metadata (Algorithm,
// Levenshtein, Origin) is intentionally excluded — it is derivation context,
// not a result fact, and including it would make CIDs non-deterministic.
func ToEntity(d *db.Domain) *Entity {
	switch d.EntityType() {
	case entity.Name:
		return &Entity{NameEntity: &NameEntity{Name: d.Name, Rank: d.Rank}}
	case entity.User:
		return &Entity{UserEntity: &UserEntity{Name: d.Name, Rank: d.Rank}}
	case entity.Package:
		return &Entity{PackageEntity: &PackageEntity{Name: d.Name, Rank: d.Rank}}
	default: // domain
		return &Entity{DomainEntity: toDomainEntity(d)}
	}
}

func toDomainEntity(d *db.Domain) *DomainEntity {
	de := &DomainEntity{
		Name:     d.Name,
		Punycode: d.Punycode,
		Rank:     d.Rank,
	}
	if d.Redirect != nil {
		de.Redirect = d.Redirect.Name
	}

	for _, r := range d.Dns {
		de.Dns = append(de.Dns, DnsRecord{Type: r.Type, Value: r.Value, Ttl: r.Ttl})
	}
	sort.Slice(de.Dns, func(i, j int) bool {
		if de.Dns[i].Type != de.Dns[j].Type {
			return de.Dns[i].Type < de.Dns[j].Type
		}
		return de.Dns[i].Value < de.Dns[j].Value
	})

	for _, a := range d.IPs {
		de.Ips = append(de.Ips, toAddress(a))
	}
	sort.Slice(de.Ips, func(i, j int) bool { return de.Ips[i].Addr < de.Ips[j].Addr })

	for _, w := range d.Whois {
		de.Whois = append(de.Whois, toWhois(w))
	}
	sort.Slice(de.Whois, func(i, j int) bool { return de.Whois[i].Created < de.Whois[j].Created })

	return de
}

func toAddress(a *db.Address) Address {
	out := Address{Addr: a.Addr, Type: a.Type}
	for _, p := range a.Ports {
		out.Ports = append(out.Ports, Port{Proto: p.Proto, Number: int64(p.Number), State: p.State, Service: p.Service})
	}
	sort.Slice(out.Ports, func(i, j int) bool { return out.Ports[i].Number < out.Ports[j].Number })
	if a.Location != nil {
		out.Location = &Location{
			Code:       a.Location.Code,
			Name:       a.Location.Name,
			Timezone:   a.Location.TimeZone,
			Latitude:   a.Location.Latitude,
			Longitude:  a.Location.Longitude,
			PostalCode: a.Location.PostalCode,
		}
	}
	return out
}

func toWhois(w db.Whois) Whois {
	return Whois{
		Created:        fmtTime(w.Created),
		Updated:        fmtTime(w.Updated),
		Expiration:     fmtTime(w.Expiration),
		Registrar:      toContact(w.Registrar),
		Registrant:     toContact(w.Registrant),
		Administrative: toContact(w.Administrative),
		Technical:      toContact(w.Technical),
		Billing:        toContact(w.Billing),
	}
}

func toContact(c *db.Contact) *Contact {
	if c == nil {
		return nil
	}
	return &Contact{
		Name: c.Name, Organization: c.Organization, Street: c.Street, City: c.City,
		Province: c.Province, PostalCode: c.PostalCode, Country: c.Country, Phone: c.Phone,
		PhoneExt: c.PhoneExt, Fax: c.Fax, FaxExt: c.FaxExt, Email: c.Email, ReferralURL: c.ReferralURL,
	}
}

func fmtTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// ToDomain converts a stored Entity back into a *db.Domain for the pipeline /
// output. Pipeline-only metadata is left zero (the caller carries it over).
func ToDomain(e *Entity) *db.Domain {
	switch {
	case e.NameEntity != nil:
		return &db.Domain{Type: entity.Name, Name: e.NameEntity.Name, Rank: e.NameEntity.Rank}
	case e.UserEntity != nil:
		return &db.Domain{Type: entity.User, Name: e.UserEntity.Name, Rank: e.UserEntity.Rank}
	case e.PackageEntity != nil:
		return &db.Domain{Type: entity.Package, Name: e.PackageEntity.Name, Rank: e.PackageEntity.Rank}
	case e.DomainEntity != nil:
		return fromDomainEntity(e.DomainEntity)
	default:
		return &db.Domain{}
	}
}

func fromDomainEntity(de *DomainEntity) *db.Domain {
	d := &db.Domain{
		Type:     entity.Domain,
		Name:     de.Name,
		Punycode: de.Punycode,
		Rank:     de.Rank,
	}
	if de.Redirect != "" {
		d.Redirect = &db.Domain{Name: de.Redirect}
	}
	for _, r := range de.Dns {
		d.Dns = append(d.Dns, &db.Dns{Type: r.Type, Value: r.Value, Ttl: r.Ttl})
	}
	for _, a := range de.Ips {
		d.IPs = append(d.IPs, fromAddress(a))
	}
	for _, w := range de.Whois {
		d.Whois = append(d.Whois, fromWhois(w))
	}
	return d
}

func fromAddress(a Address) *db.Address {
	out := &db.Address{Addr: a.Addr, Type: a.Type}
	for _, p := range a.Ports {
		out.Ports = append(out.Ports, &db.Port{Proto: p.Proto, Number: int(p.Number), State: p.State, Service: p.Service})
	}
	if a.Location != nil {
		out.Location = &db.Location{
			Code:       a.Location.Code,
			Name:       a.Location.Name,
			TimeZone:   a.Location.Timezone,
			Latitude:   a.Location.Latitude,
			Longitude:  a.Location.Longitude,
			PostalCode: a.Location.PostalCode,
		}
	}
	return out
}

func fromWhois(w Whois) db.Whois {
	return db.Whois{
		Created:        parseTime(w.Created),
		Updated:        parseTime(w.Updated),
		Expiration:     parseTime(w.Expiration),
		Registrar:      fromContact(w.Registrar),
		Registrant:     fromContact(w.Registrant),
		Administrative: fromContact(w.Administrative),
		Technical:      fromContact(w.Technical),
		Billing:        fromContact(w.Billing),
	}
}

func fromContact(c *Contact) *db.Contact {
	if c == nil {
		return nil
	}
	return &db.Contact{
		Name: c.Name, Organization: c.Organization, Street: c.Street, City: c.City,
		Province: c.Province, PostalCode: c.PostalCode, Country: c.Country, Phone: c.Phone,
		PhoneExt: c.PhoneExt, Fax: c.Fax, FaxExt: c.FaxExt, Email: c.Email, ReferralURL: c.ReferralURL,
	}
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}
