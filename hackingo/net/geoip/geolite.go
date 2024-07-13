// The MIT License (MIT)
//
// Copyright © 2019 CYBINT
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

package geoip

import (
	"net"

	"github.com/oschwald/geoip2-golang"

	"github.com/cybint/hackingo/datasets"
)

type (

	// The City struct corresponds to the data in the GeoIP2/GeoLite2 City
	// databases.
	City struct {
		City struct {
			GeoNameID uint              `json:"geoname_id,omitempty"`
			Names     map[string]string `json:"names,omitempty"`
		} `json:"city,omitempty"`
		Continent struct {
			Code      string            `json:"code,omitempty"`
			GeoNameID uint              `json:"geoname_id,omitempty"`
			Names     map[string]string `json:"names,omitempty"`
		} `json:"continent,omitempty"`
		Country struct {
			GeoNameID         uint              `json:"geoname_id,omitempty"`
			IsInEuropeanUnion bool              `json:"is_in_european_union,omitempty"`
			IsoCode           string            `json:"iso_code,omitempty"`
			Names             map[string]string `json:"names,omitempty"`
		} `json:"country,omitempty"`
		Location struct {
			AccuracyRadius uint16  `json:"accuracy_radius,omitempty"`
			Latitude       float64 `json:"latitude,omitempty"`
			Longitude      float64 `json:"longitude,omitempty"`
			MetroCode      uint    `json:"metro_code,omitempty"`
			TimeZone       string  `json:"time_zone,omitempty"`
		} `json:"location,omitempty"`
		Postal struct {
			Code string `json:"code,omitempty"`
		} `json:"postal,omitempty"`
		RegisteredCountry struct {
			GeoNameID         uint              `json:"geoname_id,omitempty"`
			IsInEuropeanUnion bool              `json:"is_in_european_union,omitempty"`
			IsoCode           string            `json:"iso_code,omitempty"`
			Names             map[string]string `json:"names,omitempty"`
		} `json:"registered_country,omitempty"`
		RepresentedCountry struct {
			GeoNameID         uint              `json:"geoname_id,omitempty"`
			IsInEuropeanUnion bool              `json:"is_in_european_union,omitempty"`
			IsoCode           string            `json:"iso_code,omitempty"`
			Names             map[string]string `json:"names,omitempty"`
			Type              string            `json:"type,omitempty"`
		} `json:"represented_country,omitempty"`
		Subdivisions []struct {
			GeoNameID uint              `json:"geoname_id,omitempty"`
			IsoCode   string            `json:"iso_code,omitempty"`
			Names     map[string]string `json:"names,omitempty"`
		} `json:"subdivisions,omitempty"`
		Traits struct {
			IsAnonymousProxy    bool `json:"is_anonymous_proxy,omitempty"`
			IsSatelliteProvider bool `json:"is_satellite_provider,omitempty"`
		} `json:"traits,omitempty"`
	}

	// The Country struct corresponds to the data in the GeoIP2/GeoLite2
	// Country databases.
	Country struct {
		Continent struct {
			Code      string            `json:"code,omitempty"`
			GeoNameID uint              `json:"geoname_id,omitempty"`
			Names     map[string]string `json:"names,omitempty"`
		} `json:"continent,omitempty"`
		Country struct {
			GeoNameID         uint              `json:"geoname_id,omitempty"`
			IsInEuropeanUnion bool              `json:"is_in_european_union,omitempty"`
			IsoCode           string            `json:"iso_code,omitempty"`
			Names             map[string]string `json:"names,omitempty"`
		} `json:"country,omitempty"`
		RegisteredCountry struct {
			GeoNameID         uint              `json:"geoname_id,omitempty"`
			IsInEuropeanUnion bool              `json:"is_in_european_union,omitempty"`
			IsoCode           string            `json:"iso_code,omitempty"`
			Names             map[string]string `json:"names,omitempty"`
		} `json:"registered_country,omitempty"`
		RepresentedCountry struct {
			GeoNameID         uint              `json:"geoname_id,omitempty"`
			IsInEuropeanUnion bool              `json:"is_in_european_union,omitempty"`
			IsoCode           string            `json:"iso_code,omitempty"`
			Names             map[string]string `json:"names,omitempty"`
			Type              string            `json:"type,omitempty"`
		} `json:"represented_country,omitempty"`
		Traits struct {
			IsAnonymousProxy    bool `json:"is_anonymous_proxy,omitempty"`
			IsSatelliteProvider bool `json:"is_satellite_provider,omitempty"`
		} `json:"traits,omitempty"`
	}
)

// NewCountry ...
func NewCountry(g *geoip2.Country) (c *Country) {
	c = &Country{}

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

// GeoIP ...
func GeoIP(ipaddr string) (country *Country, err error) {
	geolite2CityMmdb, err := datasets.Asset("GeoLite2-Country.mmdb")
	if err != nil {
		// Asset was not found
		return nil, err
	}

	db, err := geoip2.FromBytes(geolite2CityMmdb)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ip := net.ParseIP(ipaddr)
	cntry, err := db.Country(ip)
	if err == nil {
		country = NewCountry(cntry)
	}
	return
}
