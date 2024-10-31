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
package idn

// An internationalized domain name (IDN) is a domain name that includes characters
// outside of the Latin alphabet, such as letters from Arabic, Chinese, Cyrillic, or
// Devanagari scripts. IDNs allow users to use domain names in their local languages
// and scripts.

// // Idna ...
// func (d *Domain) Idna() (punycode string) {
// 	punycode, _ = idna.Punycode.ToASCII(d.String())
// 	return
// }

// // idnaFunc
// func idnaFunc(tr Result) (results []Result) {
// 	tr.Data["IDNA"] = tr.Variant.Idna()
// 	tr.Variant.Meta.IDNA = tr.Variant.Idna()
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
)

const CODE = "idn"

type None struct {
	types []string
}

func (n *None) Id() string {
	return CODE
}

func (n *None) Name() string {
	return "IDN"
}

func (n *None) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *None) Description() string {
	return "Internationalized Domain Name"
}

func (n *None) Fields() []string {
	return []string{}
}

func (n *None) Headers() []string {
	return []string{"IDN"}
}

func (n *None) Exec(in internal.Typo) (out internal.Typo) {
	in.Variant().Add("IDN", "adsf-adfsaf-s-faf")
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &None{
			types: []string{internal.ENTITY, internal.DOMAIN},
		}
	})
}
