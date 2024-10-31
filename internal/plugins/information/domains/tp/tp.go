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
package tp

// NLP Topic extraction

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/algorithms"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
)

const CODE = "tp"

type None struct {
	types []string
}

func (n *None) Id() string {
	return CODE
}

func (n *None) Name() string {
	return "Topics"
}

func (n *None) IsType(str string) bool {
	return algorithms.IsType(n.types, str)
}

func (n *None) Description() string {
	return "Topics associated with the name"
}

func (n *None) Headers() []string {
	return []string{"TOPICS"}
}

func (n *None) Exec(in internal.Typo) (out internal.Typo) {
	in.Variant().Add("TOPICS", "News")
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &None{
			types: []string{internal.ENTITY, internal.DOMAIN},
		}
	})
}
