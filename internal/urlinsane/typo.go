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
package urlinsane

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg/target"
	"github.com/rangertaha/urlinsane/pkg/fuzzy"
)

type Typo struct {
	algorithm internal.Algorithm
	original  internal.Target
	variant   internal.Target
	dist      int
}

func (t *Typo) Algorithm() internal.Algorithm {
	return t.algorithm
}

func (t *Typo) Original() internal.Target {
	return t.original
}

func (t *Typo) Variant() internal.Target {
	return t.variant
}

func (t *Typo) Active() bool {
	return t.variant.Live()
}

func (t *Typo) Ld() int {
	return t.dist
}

func (t *Typo) Clone(name string) internal.Typo {
	dist := fuzzy.Levenshtein(t.original.Name(), name)

	return &Typo{
		algorithm: t.algorithm,
		original:  t.original,
		variant:   target.New(name),
		dist:      dist,
	}
}

func (t *Typo) String() (val string) {
	if t.variant != nil {
		return t.variant.Name()
	}

	return ""
}
