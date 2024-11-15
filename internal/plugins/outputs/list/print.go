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
package list

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal"
)

func (n *Plugin) Rows(domains ...internal.Domain) (rows []string) {
	for _, domain := range domains {
		rows = append(rows, n.Row(domain))
	}
	return
}

func (n *Plugin) Row(domain internal.Domain) (row string) {
	var data []interface{}
	if domain.Cached(){
		data = append(data, fmt.Sprintf("*%d ", domain.Ld()))
	} else {
		data = append(data, fmt.Sprintf("%d  ", domain.Ld()))
	}
	if n.config.Verbose() {
		data = append(data, fmt.Sprintf("%s  ", domain.Algorithm().Name()))
	} else {
		data = append(data, fmt.Sprintf("%s  ", strings.ToUpper(domain.Algorithm().Id())))
	}
	data = append(data, fmt.Sprintf("%s  ", domain.String()))

	for _, v := range domain.Meta() {
		data = append(data, fmt.Sprintf("%s  ", v))
	}
	//  Build content for output file
	row = row + fmt.Sprint(data...)

	if domain.Live() {
		return text.FgGreen.Sprint(row)
	}
	return text.FgRed.Sprint(row)
}
