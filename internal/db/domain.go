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
	"time"

	"gorm.io/gorm"
)

type Prefix struct {
	gorm.Model
	Name     string
	Rank     int64
	Imported bool
}

type Suffix struct {
	gorm.Model
	Name     string
	Rank     int64
	Imported bool
}

type Domain struct {
	gorm.Model
	Name     string `json:"name,omitempty"`
	Punycode string `json:"punycode,omitempty"`
	Rank     int64  `json:"rank,omitempty"`

	// Related Records
	RedirectID *uint
	Redirect   *Domain   `json:"redirect,omitempty"`
	Servers    []*Server `gorm:"many2many:servers;"  json:"servers,omitempty"`
	Pages      []*Page   `gorm:"many2many:pages;"    json:"pages,omitempty"`
	Whois      []*Whois  `gorm:"many2many:whois;"    json:"whois,omitempty"`
	Dns        []*Dns    `gorm:"many2many:dns;"      json:"dns,omitempty"`

	// Language Analysis
	// Languages
	// Keywords
	// Topics
	// Vector
}

type Dns struct {
	gorm.Model
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
	Ttl   string `json:"ttl,omitempty"`

	Domains []*Domain `gorm:"many2many:dns;" json:"domains,omitempty"`
}

type Whois struct {
	gorm.Model
	DomainID         uint
	RegistrarID      uint
	RegistrantID     uint
	AdministrativeID uint
	TechnicalID      uint
	BillingID        uint
	Domain           *Domain    `json:"domain,omitempty"`
	Registrar        *Contact   `json:"registrar,omitempty"`
	Registrant       *Contact   `json:"registrant,omitempty"`
	Administrative   *Contact   `json:"administrative,omitempty"`
	Technical        *Contact   `json:"technical,omitempty"`
	Billing          *Contact   `json:"billing,omitempty"`
	Created          *time.Time `json:"created,omitempty"`
	Updated          *time.Time `json:"updated,omitempty"`
	Expiration       *time.Time `json:"expiration,omitempty"`
}

// Contact storing domain contact info
type Contact struct {
	ID           uint   `json:"id,omitempty"`
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
