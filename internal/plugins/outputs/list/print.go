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
package list

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal/db"
)

func (p *Plugin) Rows(domains ...*db.Domain) (rows []string) {
	for _, domain := range domains {
		rows = append(rows, p.Row(domain))
	}
	return
}

func (p *Plugin) Row(domain *db.Domain) (row string) {
	var data []interface{}

	data = append(data, fmt.Sprintf("%d  ", domain.Levenshtein))

	if p.config.Verbose() {
		data = append(data, fmt.Sprintf("%s  ", domain.Algorithm.Name))
	} else {
		data = append(data, fmt.Sprintf("%s  ", strings.ToUpper(domain.Algorithm.Code)))
	}

	data = append(data, fmt.Sprintf("%s  ", domain.Name))

	data = append(data, fmt.Sprintf("%s  ", domain.Punycode))

	if domain.Redirect != nil {
		data = append(data, fmt.Sprintf("%s  ", domain.Redirect.Name))
	}

	for _, record := range domain.Dns {
		data = append(data, fmt.Sprintf("%s  ", record.Value))
	}

	for _, ip := range domain.IPs {
		if ip.Location != nil {
			data = append(data, fmt.Sprintf("%s  ", ip.Location.Name))
			if ip.Location.Latitude > 0 {
				data = append(data, fmt.Sprintf("%f %f  ", ip.Location.Latitude, ip.Location.Longitude))
			}
			if ip.Location.TimeZone != "" {
				data = append(data, fmt.Sprintf("%s  ", ip.Location.TimeZone))
			}
		}
	}

	for _, record := range domain.Whois {
		data = append(data, fmt.Sprintf("%s  ", record.Registrant.Name))
	}


	// for _, record := range domain.Dns {
	// 	data = append(data, fmt.Sprintf("%s  ", record.Value))
	// }

	// for _, record := range domain.Dns {
	// 	data = append(data, fmt.Sprintf("%s  ", record.Value))
	// }

	//  Build content for output file
	row = row + fmt.Sprint(data...)

	if domain.Live() {
		return text.FgGreen.Sprint(row)
	}
	return text.FgRed.Sprint(row)
}
