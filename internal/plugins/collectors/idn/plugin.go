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

// An internationalized domain name (IDN) is a domain name that includes characters
// outside of the Latin alphabet, such as letters from Arabic, Chinese, Cyrillic, or
// Devanagari scripts. IDNs allow users to use domain names in their local languages
// and scripts.
package idn

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
	"golang.org/x/net/idna"
)

type Plugin struct {
	collectors.Plugin
}

func (i *Plugin) Exec(domain *db.Domain) (vaiant *db.Domain, err error) {
	var idn string
	idn, err = idna.Punycode.ToASCII(domain.Name)

	if domain.Name != idn {
		domain.Punycode = idn
	}

	if err != nil {
		i.Log.Error("IDN Lookup: ", err)
	}

	return domain, err
}

// Register the plugin
func init() {
	var CODE = "idn"
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{
			Plugin: collectors.Plugin{
				Num:       0,
				Code:      CODE,
				Title:     "Internationalize",
				Summary:   "Internationalized domain name",
				DependsOn: []string{},
			},
		}
	})
}
