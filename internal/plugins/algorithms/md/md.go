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

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
)

const (
	CODE        = "md"
	NAME        = "Missing Dot"
	DESCRIPTION = "Created by omitting a dot from the name"
)

type Algo struct {
	ctype     int
	config    internal.Config
	languages []internal.Language
	keyboards []internal.Keyboard
}

func (n *Algo) Id() string {
	return CODE
}

func (n *Algo) Init(conf internal.Config) {
	n.keyboards = conf.Keyboards()
	n.languages = conf.Languages()
	n.ctype = conf.Type()
	n.config = conf
}

func (n *Algo) Name() string {
	return NAME
}

func (n *Algo) Description() string {
	return DESCRIPTION
}

func (n *Algo) Exec(typo internal.Typo) (typos []internal.Typo) {
	if n.config.Type() == internal.DOMAIN {
		return n.domain(typo)
	}

	if n.config.Type() == internal.PACKAGE {
		return n.code(typo)
	}

	if n.config.Type() == internal.NAME {
		return n.name(typo)
	}
	return
}


func (n *Algo) domain(typo internal.Typo) (typos []internal.Typo) {
	original := n.config.Target().Name()
	for _, variant := range n.Func(original, ".") {
		if original != variant {
			typos = append(typos, typo.Clone(variant))
		}
	}
	return
}

func (n *Algo) code(typo internal.Typo) (typos []internal.Typo) {
	original := n.config.Target().Name()
	for _, variant := range n.Func(original, ".") {
		if original != variant {
			typos = append(typos, typo.Clone(variant))
		}
	}
	return
}

func (n *Algo) name(typo internal.Typo) (typos []internal.Typo) {
	original := n.config.Target().Name()
	for _, variant := range n.Func(original, ".") {
		if original != variant {
			typos = append(typos, typo.Clone(variant))
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
	algorithms.Add(CODE, func() internal.Algorithm {
		return &Algo{}
	})
}
