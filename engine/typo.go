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
	"github.com/rangertaha/urlinsane"
)

type Typo struct {
	keyboard  urlinsane.Keyboard
	language  urlinsane.Language
	algorithm urlinsane.Algorithm
	original  urlinsane.Domain
	variant   urlinsane.Domain
	name      string
}

func (t Typo) Keyboard() urlinsane.Keyboard {
	return t.keyboard
}

func (t Typo) Language() urlinsane.Language {
	return t.language
}

func (t Typo) Algorithm() urlinsane.Algorithm {
	return t.algorithm
}

func (t Typo) Original() urlinsane.Domain {
	return t.original
}

func (t Typo) Variant() urlinsane.Domain {
	return t.variant
}

func (t Typo) Name() string {
	return t.name
}

func (t Typo) Repr() string {
	if t.name != "" {
		return t.name
	}

	return t.variant.Repr()
}
