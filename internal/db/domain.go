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
package db

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rangertaha/urlinsane/internal/entity"
)

// Algorithm identifies the typo algorithm that produced a variant.
type Algorithm struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

// Domain is the in-flight pipeline value for a scanned entity. It is no longer a
// GORM model: results are persisted via the content-addressed store
// (internal/store). The former relational/ID fields were removed with the GORM
// result path.
type Domain struct {
	// Type classifies the named entity (domain, name, user, package). Empty is
	// treated as domain for backward compatibility — see EntityType().
	Type     entity.Type `json:"type,omitempty"`
	Name     string      `json:"name,omitempty"`
	Punycode string      `json:"punycode,omitempty"`
	Rank     int64       `json:"rank,omitempty"`

	// Related records
	Redirect *Domain    `json:"redirect,omitempty"`
	IPs      []*Address `json:"ips,omitempty"`
	Dns      []*Dns     `json:"dns,omitempty"`
	Whois    []Whois    `json:"whois,omitempty"`

	// Pipeline-only metadata (not part of the content-addressed result).
	Algorithm   Algorithm `json:"algorithm"`
	Levenshtein int       `json:"distance"`

	// Origin is the source entity this variant was generated from, used to pair
	// a variant with its origin in the Analyzers stage.
	Origin *Domain `json:"-"`
}

// EntityType returns the entity's type, defaulting to entity.Domain when unset.
func (d *Domain) EntityType() entity.Type {
	if d.Type == "" {
		return entity.Domain
	}
	return d.Type
}

// Dns is a single DNS record (e.g. A, AAAA, NS, MX, TXT, CNAME, PTR).
type Dns struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
	Ttl   string `json:"ttl,omitempty"`
}

// Whois holds the registration record and contact roles for a domain.
type Whois struct {
	Registrar      *Contact   `json:"registrar,omitempty"`
	Registrant     *Contact   `json:"registrant,omitempty"`
	Administrative *Contact   `json:"administrative,omitempty"`
	Technical      *Contact   `json:"technical,omitempty"`
	Billing        *Contact   `json:"billing,omitempty"`
	Created        *time.Time `json:"created,omitempty"`
	Updated        *time.Time `json:"updated,omitempty"`
	Expiration     *time.Time `json:"expiration,omitempty"`
}

// Contact stores entity contact info (e.g. domain whois).
type Contact struct {
	Name         string `json:"name,omitempty"`
	Organization string `json:"organization,omitempty"`
	Street       string `json:"street,omitempty"`
	City         string `json:"city,omitempty"`
	Province     string `json:"province,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	Country      string `json:"country,omitempty"`
	Phone        string `json:"phone,omitempty"`
	PhoneExt     string `json:"phone_ext,omitempty"`
	Fax          string `json:"fax,omitempty"`
	FaxExt       string `json:"fax_ext,omitempty"`
	Email        string `json:"email,omitempty"`
	ReferralURL  string `json:"referral_url,omitempty"`
}

// Live reports whether the entity resolved — i.e. it has any DNS records.
func (d *Domain) Live() bool {
	return len(d.Dns) > 0
}

// Json returns the entity marshaled as JSON.
func (d *Domain) Json() (j string) {
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Error("Marshal:", err)
	}
	return string(jsonData)
}
