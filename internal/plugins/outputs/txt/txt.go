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
package txt

import (
	"fmt"
	"os"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
)

const (
	CODE        = "txt"
	DESCRIPTION = "Text outputs one record per line and is the default"
)

type Text struct {
	config internal.Config
	output string
	typos  []internal.Typo
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
	if n.config.Progress() {
		n.typos = append(n.typos, in)
	} else {
		n.Stream(in)
	}
}

func (n *Text) Stream(in internal.Typo) {
	var data []interface{}
	data = append(data, in.Ld())
	data = append(data, in.Algorithm().Name())
	data = append(data, in.Variant().Name())

	for _, v := range in.Variant().Meta() {
		data = append(data, v)
	}
	fmt.Println(data...)
	n.output = n.output + fmt.Sprint(data...) + "\n"
}

func (n *Text) Summary(report map[string]int64) {
	fmt.Println("")
	for k, v := range report {
		fmt.Printf("%s %d   ", k, v)
	}
	fmt.Println("")
}

func (n *Text) Save() {
	if n.config.Progress() {
		for _, typo := range n.typos {
			n.Stream(typo)
		}
	}

	if n.config.File() != "" {
		results := []byte(n.output)
		if err := os.WriteFile(n.config.File(), results, 0644); err != nil {
			fmt.Printf("Error: %s", err)
		}
	}
}

// Register the plugin
func init() {
	outputs.Add(CODE, func() internal.Output {
		return &Text{}
	})
}
