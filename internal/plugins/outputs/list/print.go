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
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/pkg"
)

// Table renders the collected variants as a single, column-aligned table with a
// header. Full per-record detail (DNS/whois) is available via -f json.
func (p *Plugin) Table() string {
	tb := table.NewWriter()
	tb.SetStyle(pkg.StyleClear)
	tb.Style().Options.SeparateHeader = true
	tb.AppendHeader(table.Row{"DIST", "ALGORITHM", "NAME", "FOUND"})

	for _, d := range p.Domains {
		if p.hidden(d) {
			continue
		}
		tb.AppendRow(table.Row{d.Levenshtein, d.Algorithm.Name, d.Name, found(d)})
	}
	return tb.Render()
}

// hidden reports whether a variant is excluded by --registered/--unregistered.
func (p *Plugin) hidden(d *db.Domain) bool {
	if p.Config.Registered() && !d.Live() {
		return true
	}
	if p.Config.Unregistered() && d.Live() {
		return true
	}
	return false
}

// found summarizes where a variant was found: resolved IP addresses for
// domains, or the registries/platforms (hit services) for packages/usernames.
func found(d *db.Domain) string {
	seen := make(map[string]bool, len(d.IPs))
	var out []string
	for _, ip := range d.IPs {
		if ip != nil && ip.Addr != "" && !seen[ip.Addr] {
			seen[ip.Addr] = true
			out = append(out, ip.Addr)
		}
	}
	if len(out) > 0 {
		return strings.Join(out, ", ")
	}
	for _, h := range d.Hits {
		if h != nil && h.Service != "" {
			out = append(out, h.Service)
		}
	}
	return strings.Join(out, ", ")
}
