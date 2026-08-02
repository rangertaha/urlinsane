// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package observe

import (
	"fmt"
	"path/filepath"

	"github.com/rainycape/geoip"
)

// Location is where an address is, as this package needs it. It is deliberately
// not the MaxMind record type: keeping the interface free of that dependency is
// what lets a test locate an address without a database file.
type Location struct {
	City        string
	Country     string
	CountryCode string
	Continent   string
	TimeZone    string
	PostalCode  string
	Latitude    float64
	Longitude   float64
}

// GeoLocator resolves an address to a location. A nil Location with a nil error
// means the database simply has no entry — an absence, not a failure.
type GeoLocator interface {
	Locate(addr string) (*Location, error)
}

// OpenGeoIP opens the MaxMind database shipped in the data directory and adapts
// it to GeoLocator. It returns an error rather than a silent no-op locator so
// the caller decides whether the geo operator belongs in the plan at all.
func OpenGeoIP(dir string) (GeoLocator, error) {
	db, err := geoip.Open(filepath.Join(dir, "maxmind.db.gz"))
	if err != nil {
		return nil, fmt.Errorf("observe: open geoip database: %w", err)
	}
	if db == nil {
		return nil, fmt.Errorf("observe: geoip database opened as nil")
	}
	return maxmind{db}, nil
}

// maxmind adapts *geoip.GeoIP to GeoLocator, flattening the record's nested
// places into the flat props the graph stores.
type maxmind struct{ db *geoip.GeoIP }

func (m maxmind) Locate(addr string) (*Location, error) {
	rec, err := m.db.Lookup(addr)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, nil
	}
	loc := &Location{
		TimeZone:   rec.TimeZone,
		PostalCode: rec.PostalCode,
		Latitude:   rec.Latitude,
		Longitude:  rec.Longitude,
	}
	if rec.City != nil {
		loc.City = rec.City.Name.String()
	}
	if rec.Country != nil {
		loc.Country = rec.Country.Name.String()
		loc.CountryCode = rec.Country.Code
	}
	if rec.Continent != nil {
		loc.Continent = rec.Continent.Name.String()
	}
	return loc, nil
}

// geo locates an address.
//
// It binds On: ip, which is the whole point of pattern dispatch. The old
// collector declared DependsOn("ip") and then walked domain.IPs, so it could
// only ever see addresses that one collector had produced, on domains it had
// produced them for. Binding to the type means geo runs wherever an address
// exists — from A records, from a manifest, from a bare-IP seed — and nothing
// has to be re-declared when a second producer appears.
//
// It emits no node type at all, only props, which is why Effects declares props
// as well as nodes: a prop-only operator would otherwise be invisible to plan
// compilation (DESIGN §4).
