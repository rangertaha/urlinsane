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
package mx

import (
	"net"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
)

type Plugin struct {
	collectors.Plugin
}

func (p *Plugin) Exec(domain *db.Domain) (vaiant *db.Domain, err error) {

	records, err := net.LookupMX(domain.Name)
	if err != nil {
		p.Log.Error("MX Lookup: ", err)
	}
	for _, record := range records {
		record := strings.TrimSpace(record.Host)
		record = strings.Trim(record, ".")
		if record != "" {
			domain.Dns = append(domain.Dns, &db.DnsRecord{Type: "MX", Value: record})
		}
	}

	return domain, err

}

// Register the plugin
func init() {
	var CODE = "mx"
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{
			Plugin: collectors.Plugin{
				Num:       2,
				Code:      CODE,
				Title:     "MX Records",
				Summary:   "DNS MX Records",
				DependsOn: []string{},
			},
		}
	})
}
