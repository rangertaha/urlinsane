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
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
)

const (
	CODE        = "html"
	DESCRIPTION = "HTML formatted output"
)

type Text struct {
	table  table.Writer
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
	internal.Banner()
	n.table = table.NewWriter()
	n.table.SetOutputMirror(os.Stdout)
	n.table.AppendHeader(n.Header())
}


func (n *Text) Header() (row table.Row) {
	row = append(row, "LD")
	row = append(row, "TYPE")
	row = append(row, "TYPO")

	for _, info := range n.config.Information() {
		for _, header := range info.Headers() {
			if n.Filter(header) {
				row = append(row, header)
			}
		}
	}
	return
}

func (n *Text) Row(typo internal.Typo) (row table.Row) {
	row = append(row, typo.Ld())
	if n.config.Verbose() {
		row = append(row, typo.Algorithm().Name())
	} else {
		row = append(row, strings.ToUpper(typo.Algorithm().Id()))
	}
	row = append(row, typo.String())

	for _, info := range n.config.Information() {
		for _, header := range info.Headers() {
			if n.Filter(header) {
				meta := typo.Variant().Meta()
				if col, ok := meta[header]; ok {
					row = append(row, col)
				} else {
					row = append(row, "")
				}
			}
		}
	}
	return
}

func (n *Text) Filter(header string) bool {
	header = strings.TrimSpace(header)
	header = strings.ToLower(header)
	for _, filter := range n.config.Filters() {
		filter = strings.TrimSpace(filter)
		filter = strings.ToLower(filter)
		if filter == header {
			return true
		}
	}
	return false
}

func (n *Text) Progress(typo <-chan internal.Typo) <-chan internal.Typo {
	return typo
}

func (n *Text) Write(in internal.Typo) {
	n.table.AppendRow(n.Row(in))
}

func (n *Text) Summary(report []internal.Typo) {
	fmt.Println("")
	for k, v := range report {
		fmt.Printf("%s %d   ", k, v)
	}
	fmt.Println("")
}

func (n *Text) Save() {
	output := n.table.RenderHTML()

	if n.config.File() != "" {
		results := []byte(output)
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
