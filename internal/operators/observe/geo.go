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

package observe

import (
	"context"

	"fmt"
	"path/filepath"

	"github.com/rainycape/geoip"
	"github.com/rangertaha/urlinsane/internal/graph"
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
type geo struct {
	base
	db GeoLocator
}

func newGeo(o Options, db GeoLocator) graph.Operator {
	return geo{base: o.base(), db: db}
}

func (geo) Id() string       { return "geo" }
func (geo) Version() int     { return 1 }
func (geo) Resource() string { return ResourceGeo }

func (geo) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{TypeIP}}}
}

func (geo) Emits() graph.Effects {
	// No Nodes and no Rels: a location has no identity worth converging on, so
	// it is props on the address rather than a node of its own (DESIGN §2).
	return graph.Effects{Props: []string{
		FieldCity, FieldCountry, FieldCountryCode, FieldContinent,
		FieldLatitude, FieldLongitude, FieldTimeZone, FieldPostalCode,
	}}
}

func (o geo) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	loc, err := o.db.Locate(v.Key())
	if err != nil {
		return graph.Delta{}, graph.Failed(err)
	}
	if loc == nil {
		// The database answered and has nothing for this address. Absence, not
		// failure — the same distinction NXDOMAIN gets.
		return graph.Delta{}, graph.Empty()
	}

	self := v.Ref()
	var d graph.Delta
	str := func(field, value string) {
		if value == "" {
			return
		}
		d.Props = append(d.Props, graph.PropSet{Node: &self, Field: field, Value: graph.String(value)})
	}
	str(FieldCity, loc.City)
	str(FieldCountry, loc.Country)
	str(FieldCountryCode, loc.CountryCode)
	str(FieldContinent, loc.Continent)
	str(FieldTimeZone, loc.TimeZone)
	str(FieldPostalCode, loc.PostalCode)

	// MaxMind uses 0,0 for "unknown", so recording it would place half the
	// internet in the Gulf of Guinea.
	if loc.Latitude != 0 || loc.Longitude != 0 {
		d.Props = append(d.Props,
			graph.PropSet{Node: &self, Field: FieldLatitude, Value: graph.Float(loc.Latitude)},
			graph.PropSet{Node: &self, Field: FieldLongitude, Value: graph.Float(loc.Longitude)},
		)
	}
	if len(d.Props) == 0 {
		return graph.Delta{}, graph.Empty()
	}
	return d, graph.OK()
}
