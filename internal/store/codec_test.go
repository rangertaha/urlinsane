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
	"testing"
)

func sampleDomainEntity() *Entity {
	return &Entity{DomainEntity: &DomainEntity{
		Name:     "exmaple.com",
		Punycode: "exmaple.com",
		Rank:     42,
		Dns: []DnsRecord{
			{Type: "A", Value: "93.184.216.34", Ttl: "300"},
			{Type: "AAAA", Value: "2606:2800:220:1:248:1893:25c8:1946"},
		},
		Ips: []Address{
			{Addr: "93.184.216.34", Type: "IPv4", Location: &Location{
				Code: "US", Name: "United States", Latitude: 37.75, Longitude: -97.82,
			}},
		},
		Whois: []Whois{
			{Created: "2020-01-01T00:00:00Z", Registrant: &Contact{Name: "ACME", Country: "US"}},
		},
	}}
}

func TestEntityRoundTrip(t *testing.T) {
	in := sampleDomainEntity()
	block, err := encodeEntity(in)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	out, err := decodeEntity(block)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.DomainEntity == nil {
		t.Fatalf("expected DomainEntity union member, got %+v", out)
	}
	d := out.DomainEntity
	if d.Name != "exmaple.com" || d.Rank != 42 {
		t.Fatalf("scalar mismatch: %+v", d)
	}
	if len(d.Dns) != 2 || d.Dns[0].Type != "A" || d.Dns[0].Value != "93.184.216.34" {
		t.Fatalf("dns mismatch: %+v", d.Dns)
	}
	if len(d.Ips) != 1 || d.Ips[0].Location == nil || d.Ips[0].Location.Code != "US" {
		t.Fatalf("ips/location mismatch: %+v", d.Ips)
	}
	if len(d.Whois) != 1 || d.Whois[0].Registrant == nil || d.Whois[0].Registrant.Name != "ACME" {
		t.Fatalf("whois mismatch: %+v", d.Whois)
	}
}

func TestEntityCIDDeterminism(t *testing.T) {
	a, err := encodeEntity(sampleDomainEntity())
	if err != nil {
		t.Fatal(err)
	}
	b, err := encodeEntity(sampleDomainEntity())
	if err != nil {
		t.Fatal(err)
	}
	ca, err := cidOf(a)
	if err != nil {
		t.Fatal(err)
	}
	cb, err := cidOf(b)
	if err != nil {
		t.Fatal(err)
	}
	if !ca.Equals(cb) {
		t.Fatalf("identical entities produced different CIDs: %s vs %s", ca, cb)
	}

	// A changed record must change the CID.
	changed := sampleDomainEntity()
	changed.DomainEntity.Dns[0].Value = "10.0.0.1"
	cc, err := encodeEntity(changed)
	if err != nil {
		t.Fatal(err)
	}
	cChanged, err := cidOf(cc)
	if err != nil {
		t.Fatal(err)
	}
	if ca.Equals(cChanged) {
		t.Fatal("changed entity produced the same CID")
	}
}

func TestUnionKeyedSelectsMember(t *testing.T) {
	pkg := &Entity{PackageEntity: &PackageEntity{Name: "reqeusts", Rank: 1, Registry: "pypi"}}
	block, err := encodeEntity(pkg)
	if err != nil {
		t.Fatal(err)
	}
	out, err := decodeEntity(block)
	if err != nil {
		t.Fatal(err)
	}
	if out.PackageEntity == nil || out.DomainEntity != nil {
		t.Fatalf("expected PackageEntity member only, got %+v", out)
	}
	if out.PackageEntity.Registry != "pypi" {
		t.Fatalf("registry mismatch: %+v", out.PackageEntity)
	}
}

func TestEntityJSON(t *testing.T) {
	js, err := encodeEntityJSON(sampleDomainEntity())
	if err != nil {
		t.Fatalf("json encode: %v", err)
	}
	if len(js) == 0 {
		t.Fatal("empty json")
	}
}
