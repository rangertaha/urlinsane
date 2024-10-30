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
package md

// Missing Dot
//
// Created by omitting one dot at a time from the domain, For example
//
// Original: facebook.com.io.uk
//
// Veriants: facebookcom.io.uk
//           facebook.comio.uk
//           facebook.com.iouk

import (
	"github.com/rangertaha/urlinsane"
	"github.com/rangertaha/urlinsane/plugins/algorithms"
)

const (
	CODE        = "md"
	NAME        = "Missing Dot"
	DESCRIPTION = "Created by omitting a dot from the name"
)

type Algo struct {
	types []string
}

func (n *Algo) Id() string {
	return CODE
}
func (n *Algo) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *Algo) Name() string {
	return "Missing Dot"
}

func (n *Algo) Description() string {
	return "Created by omitting a dot from the name"
}

func (n *Algo) Exec(typo urlinsane.Typo) (typos []urlinsane.Typo) {
	for _, variant := range n.Func(typo.Original().Repr(), ".") {
		if typo.Original().Repr() != variant {
			typos = append(typos, typo.New(variant))
		}
	}
	return
}

// Func removes a character one at a time from the string.
// For example, wwwgoogle.com and www.googlecom
func (n *Algo) Func(str, character string) (results []string) {
	for i, char := range str {
		if character == string(char) {
			results = append(results, str[:i]+str[i+1:])
		}
	}
	return
}

// Register the plugin
func init() {
	algorithms.Add(CODE, func() urlinsane.Algorithm {
		return &Algo{
			types: []string{algorithms.ENTITY, algorithms.DOMAIN},
		}
	})
}
