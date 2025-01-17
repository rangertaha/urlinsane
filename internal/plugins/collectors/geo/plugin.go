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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"

	"github.com/rainycape/geoip"
)

type Plugin struct {
	collectors.Plugin
	geoip *geoip.GeoIP
}

func (p *Plugin) Init(c internal.Config) {
	p.Plugin.Init(c)
	var err error

	p.geoip, err = geoip.Open(filepath.Join(c.Dir(), "maxmind.db.gz"))
	if err != nil {
		p.Log.Error(err)
	}
}

func (p *Plugin) Exec(domain *db.Domain) (vaiant *db.Domain, err error) {
	for _, ip := range domain.IPs {
		p.GeoLookup(ip)
	}
	return domain, err
}

func (p *Plugin) GeoLookup(ip *db.Address) {
	record, err := p.geoip.Lookup(ip.Addr)
	if err != nil {
		p.Log.Error(err)
	}
	if record == nil {
		return
	}

	if record.City != nil {
		city := db.Location{}
		db.DB.FirstOrInit(&city, db.Location{
			Code:       SetCode(record, record.City.Name.String()),
			Name:       record.City.Name.String(),
			Latitude:   record.Latitude,
			Longitude:  record.Longitude,
			TimeZone:   record.TimeZone,
			PostalCode: record.PostalCode,
		})
		ip.Location = &city

	} else if record.Country != nil {
		country := db.Location{}
		db.DB.FirstOrInit(&country, db.Location{
			Code:       SetCode(record, record.Country.Name.String()),
			Name:       record.Country.Name.String(),
			Latitude:   record.Latitude,
			Longitude:  record.Longitude,
			TimeZone:   record.TimeZone,
			PostalCode: record.PostalCode,
		})
		ip.Location = &country

	} else if record.Continent != nil {
		continent := db.Location{}
		db.DB.FirstOrInit(&continent, db.Location{
			Code:       SetCode(record, record.Continent.Name.String()),
			Name:       record.Continent.Name.String(),
			Latitude:   record.Latitude,
			Longitude:  record.Longitude,
			TimeZone:   record.TimeZone,
			PostalCode: record.PostalCode,
		})
		ip.Location = &continent
	}
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

func SetCode(l *geoip.Record, name string) string {
	hasher := md5.New()
	_, err := io.WriteString(hasher, fmt.Sprintf("%f-%f-%s-%s", l.Latitude, l.Longitude, l.PostalCode, name))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
