// Copyright (C) 2024  Tal Hatchi (Rangertaha)
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
package main

import (
	"fmt"

	"github.com/rangertaha/urlinsane/typo"
)

func main() {

	conf := typo.BasicConfig{
		Domains:     []string{"google.com"},
		Keyboards:   []string{"en1"},
		Typos:       []string{"co"},
		Funcs:       []string{"ip"},
		Concurrency: 50,
		Format:      "text",
		Verbose:     false,
	}

	urli := typo.New(conf.Config())

	out := urli.Stream()

	for r := range out {
		fmt.Println(r.Variant.Live, r.Variant.Domain, r.Typo.Name, r.Data)
	}
}
