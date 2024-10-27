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
package hebrew

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/languages"
)

type Keyboard struct {
	Lang        string
	Code        string
	Name        string
	Description string
	Layout      []string
}

func (k *Keyboard) Id() string {
	return k.Code
}
func (k *Keyboard) Language() string {
	return k.Lang
}
func (k *Keyboard) Title() string {
	return k.Name
}
func (k *Keyboard) Summary() string {
	return k.Description
}
func (k *Keyboard) Layouts() []string {
	return k.Layout
}

func (k *Keyboard) Languages() []urlinsane.Language {
	return languages.Languages(k.Lang)
}

var Keyboards = []Keyboard{
	{Lang: LANGUAGE,
		Code:        "iw1",
		Name:        "Hebrew",
		Description: "Hebrew standard layout",
		Layout: []string{
			"1234567890 ",
			` פםןוטארק  `,
			` ףךלחיעכגדש `,
			` ץתצמנהבסז  `},
	},
}

func init() {
	for _, kb := range Keyboards {
		languages.AddKeyboard(kb.Code, func() urlinsane.Keyboard {
			return &kb
		})
	}

}
