// The MIT License (MIT)
//
// Copyright © 2018 Rangertaha <rangertaha@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package infos

import (
	"fmt"
	"net"
	"strings"

	"github.com/cybersectech-org/urlinsane/pkg/datasets"
	"github.com/cybersectech-org/urlinsane/pkg/typo"
	"github.com/oschwald/geoip2-golang"
)

type (
	Continent struct {
		Code      string            `json:"code,omitempty"`
		GeoNameID uint              `json:"geo_name,omitempty"`
		Names     map[string]string `json:"names,omitempty"`
	}
	Country struct {
		GeoNameID         uint              `json:"code,omitempty"`
		IsInEuropeanUnion bool              `json:"european,omitempty"`
		IsoCode           string            `json:"iso_code,omitempty"`
		Names             map[string]string `json:"names,omitempty"`
	}
	RegisteredCountry struct {
		GeoNameID         uint              `json:"geo_name,omitempty"`
		IsInEuropeanUnion bool              `json:"european,omitempty"`
		IsoCode           string            `json:"iso_code,omitempty"`
		Names             map[string]string `json:"names,omitempty"`
	}
	RepresentedCountry struct {
		GeoNameID         uint              `json:"geo_name,omitempty"`
		IsInEuropeanUnion bool              `json:"european,omitempty"`
		IsoCode           string            `json:"iso_code,omitempty"`
		Names             map[string]string `json:"names,omitempty"`
		Type              string            `json:"type,omitempty"`
	}
	Traits struct {
		IsAnonymousProxy    bool `json:"is_anonymous_proxy,omitempty"`
		IsSatelliteProvider bool `json:"is_satellite_provider,omitempty"`
	}
	GeoCountry struct {
		Continent          Continent          `json:"continent,omitempty"`
		Country            Country            `json:"country,omitempty"`
		RegisteredCountry  RegisteredCountry  `json:"registered_country,omitempty"`
		RepresentedCountry RepresentedCountry `json:"represented_country,omitempty"`
		Traits             Traits             `json:"traits,omitempty"`
	}
)

var geoIPLookup = Module{
	Code:        "GEO",
	Name:        "GeoIP Lookup",
	Description: "Show country location of ip address",
	exe:         geoIPLookupFunc,
	headers:     []string{"IPv4", "IPv6", "GEO"},
}

// NewCountry ...
func toCountry(g *geoip2.Country) (c GeoCountry) {
	c.Continent.Code = g.Continent.Code
	c.Continent.GeoNameID = g.Continent.GeoNameID
	c.Continent.Names = g.Continent.Names

	c.Country.GeoNameID = g.Country.GeoNameID
	c.Country.IsInEuropeanUnion = g.Country.IsInEuropeanUnion
	c.Country.IsoCode = g.Country.IsoCode
	c.Country.Names = g.Country.Names

	c.RegisteredCountry.GeoNameID = g.RegisteredCountry.GeoNameID
	c.RegisteredCountry.IsInEuropeanUnion = g.RegisteredCountry.IsInEuropeanUnion
	c.RegisteredCountry.IsoCode = g.RegisteredCountry.IsoCode
	c.RegisteredCountry.Names = g.RegisteredCountry.Names

	c.RepresentedCountry.GeoNameID = g.RegisteredCountry.GeoNameID
	c.RepresentedCountry.IsInEuropeanUnion = g.RepresentedCountry.IsInEuropeanUnion
	c.RepresentedCountry.IsoCode = g.RepresentedCountry.IsoCode
	c.RepresentedCountry.Names = g.RepresentedCountry.Names
	c.RepresentedCountry.Type = g.RepresentedCountry.Type

	c.Traits.IsAnonymousProxy = g.Traits.IsAnonymousProxy
	c.Traits.IsSatelliteProvider = g.Traits.IsSatelliteProvider

	return
}

// geoIPLookupFunc
func geoIPLookupFunc(tr typo.Result) (results []typo.Result) {
	tr = checkIP(tr)
	if tr.Live {
		geolite2CityMmdb, err := datasets.Asset("GeoLite2-Country.mmdb")
		if err != nil {
			// Asset was not found.
		}

		db, err := geoip2.FromBytes(geolite2CityMmdb)
		if err != nil {
			fmt.Print(err)
		}
		defer db.Close()

		ipv4s, ok := tr.Data["IPv4"]
		if ok {
			ips := strings.Split(ipv4s, "\n")
			for _, ip4 := range ips {
				ip := net.ParseIP(ip4)
				if ip != nil {
					record, err := db.Country(ip)
					if err != nil {
						fmt.Print(err)
					}
					tr.Data["GEO"] = fmt.Sprint(record.Country.Names["en"])
					tr.Meta["geo"] = toCountry(record)
				}
			}
		}
	}

	// If you are using strings that may be invalid, check that ip is not nil
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Live: tr.Live, Data: tr.Data, Meta: tr.Meta})
	return
}
