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

type arKeyboard struct {
	Lang        string
	Code        string
	Name        string
	Description string
	Layout      []string
}

var arKeyboards = []arKeyboard{
	{
		Lang:        LANGUAGE,
		Code:        "ar1",
		Name:        "غفقثصض",
		Description: "Arabic keyboard layout",
		Layout: []string{
			"١٢٣٤٥٦٧٨٩٠- ",
			"ةجحخهعغفقثصض",
			"  كمنتالبيسش",
			"     ورزدذطظ"},
	},
	{
		Lang:        LANGUAGE,
		Code:        "ar2",
		Name:        "AZERTY",
		Description: "Arabic PC keyboard layout",
		Layout: []string{
			` é   -è çà   `,
			"ذدجحخهعغفقثصض",
			"  طكمنتالبيسش",
			"   ظزوةىلارؤءئ"},
	},
	{
		Lang:        LANGUAGE,
		Code:        "ar3",
		Name:        "غفقثصض",
		Description: "Arabic North african keyboard layout",
		Layout: []string{
			"1234567890  ",
			"ةجحخهعغفقثصض",
			"  كمنتالبيسش",
			"     ورزدذطظ"},
	},
	{
		Lang:        LANGUAGE,
		Code:        "ar4",
		Name:        "QWERTY",
		Description: "Arabic keyboard layout",
		Layout: []string{
			"١٢٣٤٥٦٧٨٩٠  ",
			"ظثةهيوطترعشق",
			"   لكجحغفدسا",
			"     منبذصخز"},
	},
}

func (k *arKeyboard) Id() string {
	return k.Code
}
func (k *arKeyboard) Language() string {
	return k.Lang
}
func (k *arKeyboard) Title() string {
	return k.Name
}
func (k *arKeyboard) Summary() string {
	return k.Description
}
func (k *arKeyboard) Layouts() []string {
	return k.Layout
}

func (k *arKeyboard) Languages() []urlinsane.Language {
	return languages.Languages(k.Lang)
}

func init() {
	for _, kb := range arKeyboards {
		languages.AddKeyboard(kb.Code, func() urlinsane.Keyboard {
			return &kb
		})
	}

}
