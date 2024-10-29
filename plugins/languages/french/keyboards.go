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
package french

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/languages"
)

type Keyboard struct {
	lang        string
	code        string
	name        string
	description string
	layout      []string
}

func (k *Keyboard) Id() string {
	return k.code
}
func (k *Keyboard) Language() string {
	return k.lang
}
func (k *Keyboard) Name() string {
	return k.name
}
func (k *Keyboard) Description() string {
	return k.description
}
func (k *Keyboard) Layouts() []string {
	return k.layout
}
func (k *Keyboard) Adjacent(char string)  []string {
	return languages.Adjacent(k.layout, char)
}
func (k *Keyboard) Languages() []urlinsane.Language {
	return languages.Languages(k.lang)
}

var Keyboards = []Keyboard{
	// https://en.wikipedia.org/wiki/CSA_keyboard
	{
		lang:        LANGUAGE,
		code:        "fr1",
		name:        "ACNOR",
		description: " French Canadian keyboard layout",
		layout: []string{
			"ù1234567890-  ",
			" qwertyuiop çà",
			" asdfghjkl è  ",
			"  zxcvbnm  é  "},
	},
}

func init() {
	for _, kb := range Keyboards {
		languages.AddKeyboard(kb.code, func() urlinsane.Keyboard {
			return &kb
		})
	}
}
