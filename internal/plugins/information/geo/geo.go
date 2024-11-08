// Copyright (C) 2024 Rangertaha
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
	"strings"

	"github.com/rainycape/geoip"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/models"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
)

const (
	ORDER       = 10
	CODE        = "geo"
	NAME        = "GeoIP Lookup"
	DESCRIPTION = "Retrieves location of IP addresses"
)

//go:embed GeoLite2-City.mmdb
var dataFile embed.FS

type Plugin struct {
}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Order() int {
	return ORDER
}

func (n *Plugin) Name() string {
	return NAME
}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (n *Plugin) Headers() []string {
	return []string{"GEO"}
}

func (n *Plugin) Exec(in internal.Typo) (out internal.Typo) {
	orig, vari := in.Get()

	countries := []string{}
	for i, ip := range vari.IPv4 {
		loc := n.getIp(ip.Address)
		// fmt.Println(orig, vari.FQDN)

		if loc.Country != nil {
			countries = append(countries, loc.Country.String())
			vari.IPv4[i].Location = models.Location{
				Continent: &models.Place{
					Code:      loc.Continent.Code,
					Name:      loc.Continent.Name.String(),
					GeonameID: loc.Continent.GeonameID,
				},
				Country: &models.Place{
					Code:      loc.Country.Code,
					Name:      loc.Country.Name.String(),
					GeonameID: loc.Country.GeonameID,
				},
				Latitude:            loc.Latitude,
				Longitude:           loc.Longitude,
				MetroCode:           loc.MetroCode,
				PostalCode:          loc.PostalCode,
				TimeZone:            loc.TimeZone,
				IsAnonymousProxy:    loc.IsAnonymousProxy,
				IsSatelliteProvider: loc.IsSatelliteProvider,
			}

			if loc.City != nil {
				vari.IPv4[i].Location.City = &models.Place{
					Code:      loc.City.Code,
					Name:      loc.City.Name.String(),
					GeonameID: loc.City.GeonameID,
				}
			}

			if loc.RegisteredCountry != nil {
				vari.IPv4[i].Location.RegisteredCountry = &models.Place{
					Code:      loc.RegisteredCountry.Code,
					Name:      loc.RegisteredCountry.Name.String(),
					GeonameID: loc.RegisteredCountry.GeonameID,
				}
			}

			if loc.RepresentedCountry != nil {
				vari.IPv4[i].Location.RepresentedCountry = &models.Place{
					Code:      loc.RepresentedCountry.Code,
					Name:      loc.RepresentedCountry.Name.String(),
					GeonameID: loc.RepresentedCountry.GeonameID,
				}
			}
		}
	}

	in.Set(orig, vari)
	in.SetMeta("GEO", strings.Join(countries, " "))
	return in
}

func (n *Plugin) getIp(ip string) *geoip.Record {
	file, err := dataFile.Open("GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	db, err := geoip.New(file.(io.ReadSeeker))
	if err != nil {
		panic(err)
	}

	res, err := db.Lookup(ip)
	if err != nil {
		panic(err)
	}
	return res
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &Plugin{}
	})
}
