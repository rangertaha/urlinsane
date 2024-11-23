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
package domain

// type Place struct {
// 	// Code is the given code for the place. For continents, this
// 	// value is one of AF (Africa), AS (Asia), EU (Europe), OC (Oceania),
// 	// NA (North America) and SA (South America). For countries, its
// 	// their ISO 3166-1 2 letter code (see http://en.wikipedia.org/wiki/ISO_3166-1).
// 	Code string `json:"code,omitempty"`
// 	// GeonameID is the place's ID in the geonames database. See
// 	// http://www.geonames.org for more information.
// 	GeonameID int `json:"geoname,omitempty"`
// 	// Name is the place name, usually with several translations.
// 	Name string `json:"name,omitempty"`
// }

// // Record hold the information returned for a given
// // IP address. See the comments on each field for more
// // information.
// type Location struct {
// 	// Continent contains information about the continent
// 	// where the record is located.
// 	Continent *Place `json:"continent,omitempty"`
// 	// Country contains information about the country
// 	// where the record is located.
// 	Country *Place `json:"country,omitempty"`
// 	// RegisteredCountry contains information about the
// 	// country where the ISP has registered the IP address
// 	// for this record. Note that this field might be
// 	// different from Country.
// 	RegisteredCountry *Place `json:"registered,omitempty"`
// 	// RepresentedCountry is non nil only when the record
// 	// belongs an entity representing a country, like an
// 	// embassy or a military base. Note that it might be
// 	// diferrent from Country.
// 	RepresentedCountry *Place `json:"represented,omitempty"`
// 	// City contains information about the city where the
// 	// record is located.
// 	City *Place `json:"city,omitempty"`
// 	// Subdivisions contains details about the subdivisions
// 	// of the country where the record is located. Subdivisions
// 	// are arranged from largest to smallest and the number of
// 	// them will vary depending on the country.
// 	Subdivisions []*Place `json:"subdivisions,omitempty"`
// 	// Latitude of the location associated with the record.
// 	// Note that a 0 Latitude and a 0 Longitude means the
// 	// coordinates are not known.
// 	Latitude float64 `json:"lat,omitempty"`
// 	// Longitude of the location associated with the record.
// 	// Note that a 0 Latitude and a 0 Longitude means the
// 	// coordinates are not known.
// 	Longitude float64 `json:"lon,omitempty"`
// 	// MetroCode contains the metro code associated with the
// 	// record. These are only available in the US
// 	MetroCode int `json:"metrocode,omitempty"`
// 	// PostalCode associated with the record. These are available in
// 	// AU, CA, FR, DE, IT, ES, CH, UK and US.
// 	PostalCode string `json:"postcode,omitempty"`
// 	// TimeZone associated with the record, in IANA format (e.g.
// 	// America/New_York). See http://www.iana.org/time-zones.
// 	TimeZone string `json:"timezone,omitempty"`
// 	// IsAnonymousProxy is true iff the record belongs
// 	// to an anonymous proxy.
// 	IsAnonymousProxy bool `json:"proxy,omitempty"`
// 	// IsSatelliteProvider is true iff the record is
// 	// in a block managed by a satellite ISP that provides
// 	// service to multiple countries. These IPs might be
// 	// in high risk countries.
// 	IsSatelliteProvider bool `json:"satellite,omitempty"`
// }
