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

// Package store persists urlinsane scan results as a content-addressed IPLD
// Merkle DAG in a local blockstore, with a lightweight SQLite index for
// name-based lookups (IPLD itself has no query layer). It is the source of
// truth for results; the pipeline keeps mutating a plain Go struct and only
// encodes/decodes at this boundary.
package store

import (
	_ "embed"
	"fmt"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
)

//go:embed schema.ipldsch
var schemaBytes []byte

var (
	entityProto schema.TypedPrototype
	scanProto   schema.TypedPrototype
)

func init() {
	ts, err := ipld.LoadSchemaBytes(schemaBytes)
	if err != nil {
		panic(fmt.Sprintf("store: loading IPLD schema: %v", err))
	}
	entityProto = bindnode.Prototype((*Entity)(nil), ts.TypeByName("Entity"))
	scanProto = bindnode.Prototype((*Scan)(nil), ts.TypeByName("Scan"))
}

// Entity is a keyed union over the entity types. Exactly one field is non-nil;
// the non-nil field selects the union member (and its keyed discriminant).
type Entity struct {
	DomainEntity  *DomainEntity
	NameEntity    *NameEntity
	UserEntity    *UserEntity
	PackageEntity *PackageEntity
}

type DomainEntity struct {
	Name     string
	Punycode string
	Rank     int64
	Redirect string
	Dns      []DnsRecord
	Ips      []Address
	Whois    []Whois
}

type NameEntity struct {
	Name string
	Rank int64
}

type UserEntity struct {
	Name string
	Rank int64
}

type PackageEntity struct {
	Name     string
	Rank     int64
	Registry string
}

type DnsRecord struct {
	Type  string
	Value string
	Ttl   string
}

type Address struct {
	Addr     string
	Type     string
	Ports    []Port
	Location *Location
}

type Port struct {
	Proto   string
	Number  int64
	State   string
	Service string
}

type Location struct {
	Code       string
	Name       string
	Timezone   string
	Latitude   float64
	Longitude  float64
	PostalCode string
}

type Whois struct {
	Created        string
	Updated        string
	Expiration     string
	Registrar      *Contact
	Registrant     *Contact
	Administrative *Contact
	Technical      *Contact
	Billing        *Contact
}

type Contact struct {
	Name         string
	Organization string
	Street       string
	City         string
	Province     string
	PostalCode   string
	Country      string
	Phone        string
	PhoneExt     string
	Fax          string
	FaxExt       string
	Email        string
	ReferralURL  string
}

type Scan struct {
	Query   string
	Created string
	Results []datamodel.Link
}
