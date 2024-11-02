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
)

type Typo struct {
	algorithm internal.Algorithm
	original  internal.Target
	variant   internal.Target
	active    bool
}

// func (t *Typo) Id(num ...int64) (id string) {
// 	if len(num) > 0 {
// 		t.id = num[0]
// 	}

// 	if t.id > 0 {
// 		return fmt.Sprintf("%v", t.id)
// 	}

// 	// if t.keyboard != nil {
// 	// 	id = fmt.Sprintf("%s%s", id, t.keyboard.Id())
// 	// }
// 	// if t.language != nil {
// 	// 	id = fmt.Sprintf("%s%s", id, t.language.Id())
// 	// }
// 	if t.algorithm != nil {
// 		id = fmt.Sprintf("%s%s", id, t.algorithm.Id())
// 	}
// 	if t.original != nil {
// 		id = fmt.Sprintf("%s%s", id, t.original.String())
// 	}
// 	if t.variant != nil {
// 		id = fmt.Sprintf("%s%s", id, t.variant.String())
// 	}

// 	return
// }

// func (t *Typo) Keyboards() []internal.Keyboard {
// 	return t.keyboards
// }

// func (t *Typo) Languages() []internal.Language {
// 	return t.languages
// }

func (t *Typo) Algorithm() internal.Algorithm {
	return t.algorithm
}

func (t *Typo) Original() internal.Target {
	return t.original
}

func (t *Typo) Variant() internal.Target {
	return t.variant
}

// func (t *Typo) SetType(val string) {
// 	t.ttype = val
// }

// func (t *Typo) GetType() string {
// 	return t.ttype
// }

func (t *Typo) Active(a ...bool) bool {
	if len(a) > 0 {
		t.active = a[0]
	}
	return t.active
}

func (t *Typo) Clone(name string) internal.Typo {
	return &Typo{
		algorithm: t.algorithm,
		original:  t.original,
		variant:   target.New(name),
	}
}

func (t *Typo) String() (val string) {
	if t.variant != nil {
		return t.variant.Name()
	}

	return ""
}
