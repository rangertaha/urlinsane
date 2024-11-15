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
package tld2

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/domain"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	"github.com/rangertaha/urlinsane/internal/utils/datasets"
	algo "github.com/rangertaha/urlinsane/pkg/typo"
)

const (
	CODE        = "tld2"
	NAME        = "Wrong TLD2"
	DESCRIPTION = "Wrong second level domain (TLD2)"
)

type Algo struct{}

func (n *Algo) Id() string {
	return CODE
}

func (n *Algo) Name() string {
	return NAME
}
func (n *Algo) Description() string {
	return DESCRIPTION
}

func (n *Algo) Exec(original internal.Domain, acc internal.Accumulator) (err error) {
	for _, variant := range algo.SecondLevelDomain(original.Suffix(), datasets.TLD...) {
		if original.Suffix() != variant {
			acc.Add(domain.Variant(n, original.Prefix(), original.Name(), variant))
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
