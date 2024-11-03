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
package text

import (
	"fmt"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
)

const (
	CODE        = "text"
	DESCRIPTION = "Text outputs one record per line and is the default"
)

type Text struct {
	config internal.Config
}

func (n *Text) Id() string {
	return CODE
}

func (n *Text) Description() string {
	return DESCRIPTION
}

func (n *Text) Init(conf internal.Config) {
	n.config = conf
}

func (n *Text) Write(in internal.Typo) {
	var data []interface{}
	data = append(data, in.Algorithm().Name())
	data = append(data, in.Variant().Name())

	for _, v := range in.Variant().Meta() {
		data = append(data, v)
	}
	fmt.Println(data...)
}

func (n *Text) Summary(report map[string]int64) {
	// footer := table.Row{}
	// for k, v := range report {
	// 	footer = append(footer, k, v)
	// }

	// n.table.AppendFooter(footer)
	// n.table.SetStyle(StyleDefault)
}

func (n *Text) Save() {}

// Register the plugin
func init() {
	outputs.Add(CODE, func() internal.Output {
		return &Text{}
	})
}
