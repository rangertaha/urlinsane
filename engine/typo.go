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
package engine

import (
	"fmt"

	"github.com/rangertaha/urlinsane"
)

type Typo struct {
	keyboard  []urlinsane.Keyboard
	language  []urlinsane.Language
	algorithm urlinsane.Algorithm
	original  urlinsane.Domain
	variant   urlinsane.Domain
	active    bool
	id        int64
}

func NewTypo(tpy *Typo) *Typo {
	typo := tpy
	return typo
}

func (t *Typo) Id(num ...int64) (id string) {
	if len(num) > 0 {
		t.id = num[0]
	}

	if t.id > 0 {
		return fmt.Sprintf("%v", t.id)
	}

	// if t.keyboard != nil {
	// 	id = fmt.Sprintf("%s%s", id, t.keyboard.Id())
	// }
	// if t.language != nil {
	// 	id = fmt.Sprintf("%s%s", id, t.language.Id())
	// }
	if t.algorithm != nil {
		id = fmt.Sprintf("%s%s", id, t.algorithm.Id())
	}
	if t.original != nil {
		id = fmt.Sprintf("%s%s", id, t.original.Repr())
	}
	if t.variant != nil {
		id = fmt.Sprintf("%s%s", id, t.variant.Repr())
	}

	return
}

func (t *Typo) Keyboards() []urlinsane.Keyboard {
	return t.keyboard
}

func (t *Typo) Languages() []urlinsane.Language {
	return t.language
}

func (t *Typo) Algorithm() urlinsane.Algorithm {
	return t.algorithm
}

func (t *Typo) Original() urlinsane.Domain {
	return t.original
}

func (t *Typo) Variant() urlinsane.Domain {
	return t.variant
}

func (t *Typo) Active(a ...bool) bool {
	if len(a) > 0 {
		t.active = a[0]
	}
	return t.active
}

func (t *Typo) New(str string) urlinsane.Typo {
	return &Typo{
		language:  t.language,
		keyboard:  t.keyboard,
		algorithm: t.algorithm,
		original:  t.original,
		variant:   NewDomain(str),
	}
}

func (t *Typo) Repr() string {
	if t.variant != nil {
		return t.variant.Repr()
	}

	return t.original.Repr()
}
