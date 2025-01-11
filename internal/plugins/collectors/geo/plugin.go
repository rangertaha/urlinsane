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
package geo

import (
	"embed"
	"io"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"

	"github.com/rainycape/geoip"
)

//go:embed GeoLite2-Country.mmdb
var dataFile embed.FS

type Plugin struct {
	collectors.Plugin
}

func (i *Plugin) Exec(domain *db.Domain) (vaiant *db.Domain, err error) {

	for _, ip := range domain.IPs {
		if ip.Type == "IPv4" {
			ip.Location, err = GeoLookup(ip.Addr)
		}
	}

	return domain, err
}

// Register the plugin
func init() {
	var CODE = "geo"
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{
			Plugin: collectors.Plugin{
				Num:       10,
				Code:      CODE,
				Title:     "GeoIP Lookup",
				Summary:   "Retrieves location of IP addresses",
				DependsOn: []string{"ip"},
			},
		}
	})
}

func GeoLookup(ip string) (*db.Location, error) {
	var err error
	var r *geoip.Record
	var loc db.Location
	r, err = GetGeo(ip)
	if err != nil {
		return nil, err
	}
	loc.Latitude = r.Latitude
	loc.Longitude = r.Longitude
	loc.IsAnonymousProxy = r.IsAnonymousProxy
	loc.IsSatelliteProvider = r.IsSatelliteProvider
	loc.MetroCode = r.MetroCode
	loc.PostalCode = r.PostalCode
	loc.TimeZone = r.TimeZone
	if r.Continent != nil {

	}

	if r.Country != nil {
		place := &db.Place{
			Code:      r.Country.Code,
			GeonameID: r.Country.GeonameID,
			Name:      r.Country.Name.String(),
		}
		loc.Country = place

	}

	if r.City != nil {
		place := &db.Place{
			Code:      r.City.Code,
			GeonameID: r.City.GeonameID,
			Name:      r.City.Name.String(),
		}
		loc.City = place
	}

	if r.RegisteredCountry != nil {
		place := &db.Place{
			Code:      r.RegisteredCountry.Code,
			GeonameID: r.RegisteredCountry.GeonameID,
			Name:      r.RegisteredCountry.Name.String(),
		}
		loc.RegisteredCountry = place
	}

	if r.RepresentedCountry != nil {
		place := &db.Place{
			Code:      r.RepresentedCountry.Code,
			GeonameID: r.RepresentedCountry.GeonameID,
			Name:      r.RepresentedCountry.Name.String(),
		}
		loc.RepresentedCountry = place
	}

	for _, subd := range r.Subdivisions {
		loc.Subdivisions = append(loc.Subdivisions, &db.Place{
			Code:      subd.Code,
			GeonameID: subd.GeonameID,
			Name:      subd.Name.String(),
		})

	}

	db.DB.FirstOrCreate(&loc, db.Location{Latitude: loc.Latitude, Longitude: loc.Longitude})

	return &loc, err
}

func GetGeo(ip string) (r *geoip.Record, err error) {
	file, err := dataFile.Open("GeoLite2-Country.mmdb")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	db, err := geoip.New(file.(io.ReadSeeker))
	if err != nil {
		return nil, err
	}

	return db.Lookup(ip)
}
