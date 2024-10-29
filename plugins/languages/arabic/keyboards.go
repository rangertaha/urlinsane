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
package arabic

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
func (k *Keyboard) Adjacent(char string) []string {
	return languages.Adjacent(k.layout, char)
}
func (k *Keyboard) Languages() []urlinsane.Language {
	return languages.Languages(k.lang)
}

var Keyboards = []Keyboard{
	{
		lang:        LANGUAGE,
		code:        "ar1",
		name:        "غفقثصض",
		description: "Arabic keyboard layout",
		layout: []string{
			"١٢٣٤٥٦٧٨٩٠- ",
			"ةجحخهعغفقثصض",
			"  كمنتالبيسش",
			"     ورزدذطظ"},
	},
	{
		lang:        LANGUAGE,
		code:        "ar2",
		name:        "AZERTY",
		description: "Arabic PC keyboard layout",
		layout: []string{
			` é   -è çà   `,
			"ذدجحخهعغفقثصض",
			"  طكمنتالبيسش",
			"   ظزوةىلارؤءئ"},
	},
	{
		lang:        LANGUAGE,
		code:        "ar3",
		name:        "غفقثصض",
		description: "Arabic North african keyboard layout",
		layout: []string{
			"1234567890  ",
			"ةجحخهعغفقثصض",
			"  كمنتالبيسش",
			"     ورزدذطظ"},
	},
	{
		lang:        LANGUAGE,
		code:        "ar4",
		name:        "QWERTY",
		description: "Arabic keyboard layout",
		layout: []string{
			"١٢٣٤٥٦٧٨٩٠  ",
			"ظثةهيوطترعشق",
			"   لكجحغفدسا",
			"     منبذصخز"},
	},
}

func init() {
	for _, kb := range Keyboards {
		languages.AddKeyboard(kb.code, func() urlinsane.Keyboard {
			return &kb
		})
	}

}
