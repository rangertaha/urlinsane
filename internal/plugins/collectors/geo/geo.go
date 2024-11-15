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

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg"
	log "github.com/sirupsen/logrus"
)

//go:embed GeoLite2-City.mmdb
var dataFile embed.FS

type Plugin struct {
	log *log.Entry
}

func (i *Plugin) Init(c internal.Config) {
	i.log = log.WithFields(log.Fields{"plugin": CODE, "method": "Exec"})
}

func (i *Plugin) Exec(acc internal.Accumulator) (err error) {
	l := i.log.WithFields(log.Fields{"domain": acc.Domain().String()})
	// if acc.Domain().Cached() {
	// 	return acc.Next()
	// }

	dns := make(pkg.DnsRecords, 0)
	if err = acc.Unmarshal("DNS", dns); err != nil {
		l.Error("Unmarshal DNS: ", err)
	}

	if gip, err := NewGeoIp(dns.Array("A")...); err == nil {
		acc.SetJson("GEO", gip.Json())
		acc.SetMeta("GEO", gip.String())
	}

	return acc.Next()
}
