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
package models

// type Domain struct {
// 	Prefix   string `json:"prefix,omitempty"`
// 	Name     string `json:"name,omitempty"`
// 	Suffix   string `json:"suffix,omitempty"`
// 	FQDN     string `json:"fqdn"`
// 	Punycode string

// 	Live       bool   `json:"live,omitempty"`
// 	IPv4       []IP   `json:"ipv4,omitempty"`
// 	IPv6       []IP   `json:"ipv6,omitempty"`
// 	Banner     Banner `json:"response,omitempty"`
// 	Screenshot string `json:"screenshot,omitempty"`
// 	Html       string `json:"html,omitempty"`
// 	Ssdeep     string `json:"ssdeep,omitempty"`

// 	Dns   []DnsRecord   `json:"dns,omitempty"`
// 	Whois []WhoisRecord `json:"whois,omitempty"`
// }

// type DnsRecord struct {
// 	Type  string `json:"type,omitempty"`
// 	Value string `json:"value,omitempty"`
// 	Ttl   string `json:"ttl,omitempty"`
// }

// type Banner struct {
// 	Status        string `json:"status,omitempty"`
// 	Protocol      string `json:"protocol,omitempty"`
// 	Headers       string `json:"headers,omitempty"`
// 	ContentLength int64  `json:"length,omitempty"`
// 	Proto         string `json:"proto,omitempty"`
// 	StatusCode    int    `json:"code,omitempty"`
// 	TLS           string `json:"tls,omitempty"`
// 	Cookies       string `json:"cookies,omitempty"`
// 	Trailer       string `json:"trailer,omitempty"`
// }

// type WhoisRecord struct {
// 	Domain         *WhoisDomain `json:"domain,omitempty"`
// 	Registrar      *Contact     `json:"registrar,omitempty"`
// 	Registrant     *Contact     `json:"registrant,omitempty"`
// 	Administrative *Contact     `json:"administrative,omitempty"`
// 	Technical      *Contact     `json:"technical,omitempty"`
// 	Billing        *Contact     `json:"billing,omitempty"`
// }

// // Domain storing domain name info
// type WhoisDomain struct {
// 	ID                   string     `json:"id,omitempty"`
// 	Domain               string     `json:"domain,omitempty"`
// 	Punycode             string     `json:"punycode,omitempty"`
// 	Name                 string     `json:"name,omitempty"`
// 	Extension            string     `json:"extension,omitempty"`
// 	WhoisServer          string     `json:"whois_server,omitempty"`
// 	Status               []string   `json:"status,omitempty"`
// 	NameServers          []string   `json:"name_servers,omitempty"`
// 	DNSSec               bool       `json:"dnssec,omitempty"`
// 	CreatedDate          string     `json:"created_date,omitempty"`
// 	CreatedDateInTime    *time.Time `json:"created_date_in_time,omitempty"`
// 	UpdatedDate          string     `json:"updated_date,omitempty"`
// 	UpdatedDateInTime    *time.Time `json:"updated_date_in_time,omitempty"`
// 	ExpirationDate       string     `json:"expiration_date,omitempty"`
// 	ExpirationDateInTime *time.Time `json:"expiration_date_in_time,omitempty"`
// }

// // Contact storing domain contact info
// type Contact struct {
// 	ID           string `json:"id,omitempty"`
// 	Name         string `json:"name,omitempty"`
// 	Organization string `json:"organization,omitempty"`
// 	Street       string `json:"street,omitempty"`
// 	City         string `json:"city,omitempty"`
// 	Province     string `json:"province,omitempty"`
// 	PostalCode   string `json:"postal_code,omitempty"`
// 	Country      string `json:"country,omitempty"`
// 	Phone        string `json:"phone,omitempty"`
// 	PhoneExt     string `json:"phone_ext,omitempty"`
// 	Fax          string `json:"fax,omitempty"`
// 	FaxExt       string `json:"fax_ext,omitempty"`
// 	Email        string `json:"email,omitempty"`
// 	ReferralURL  string `json:"referral_url,omitempty"`
// }

// func (d *Domain) Fqdn() (name string) {
// 	name = fmt.Sprintf("%s.%s.%s", d.Prefix, d.Name, d.Suffix)
// 	name = strings.ReplaceAll(name, "..", ".")
// 	name = strings.Trim(name, ".")
// 	d.FQDN = name
// 	return
// }

// func (d *Domain) IP() (ip string) {
// 	for _, ip := range d.IPv4 {
// 		return ip.Address
// 	}
// 	for _, ip := range d.IPv6 {
// 		return ip.Address
// 	}
// 	return
// }
