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
package gr

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	algo "github.com/rangertaha/urlinsane/pkg/typo"
)

const (
	CODE        = "gr"
	NAME        = "Grapheme Replacement"
	DESCRIPTION = "Replaces an alphabet in the target domain"
)

type Algo struct {
	config    internal.Config
	languages []internal.Language
}

func (n *Algo) Id() string {
	return CODE
}

func (n *Algo) Init(conf internal.Config) {
	n.config = conf
}

func (n *Algo) Name() string {
	return NAME
}
func (n *Algo) Description() string {
	return DESCRIPTION
}

func (n *Algo) Exec(typo internal.Typo) (typos []internal.Typo) {
	orig, _ := typo.Get()

	for _, language := range n.languages {
		for _, variant := range algo.GraphemeReplacement(orig.Name, language.Graphemes()...) {
			if orig.Name != variant {

				new := typo.New(n, orig, domain.New(orig.Prefix, variant, orig.Suffix))
				typos = append(typos, new)
			}
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
