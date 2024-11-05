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

// An internationalized domain name (IDN) is a domain name that includes characters
// outside of the Latin alphabet, such as letters from Arabic, Chinese, Cyrillic, or
// Devanagari scripts. IDNs allow users to use domain names in their local languages
// and scripts.
package idn




import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
	"golang.org/x/net/idna"
)

const (
	CODE        = "idn"
	NAME        = "Internationalize"
	DESCRIPTION = "Internationalized Domain Name"
)

type None struct {
	types []string
}

func (n *None) Id() string {
	return CODE
}

func (n *None) Description() string {
	return DESCRIPTION
}

func (n *None) Headers() []string {
	return []string{"IDN"}
}

func (n *None) Exec(in internal.Typo) (out internal.Typo) {

	if punycode, err := idna.Punycode.ToASCII(in.Variant().Name()); err == nil {
		in.Variant().Add("IDN", punycode)
	}
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &None{}
	})
}
