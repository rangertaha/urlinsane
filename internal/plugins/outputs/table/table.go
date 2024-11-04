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
package table

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
	"golang.org/x/term"
)

const (
	CODE        = "table"
	DESCRIPTION = "Pretty table output format with color"
)

type Text struct {
	table  table.Writer
	config internal.Config
}

func (n *Text) Id() string {
	return CODE
}

func (n *Text) Description() string {
	return "Text outputs one record per line and is the default"
}

func (n *Text) Init(conf internal.Config) {
	n.config = conf
	n.table = table.NewWriter()

	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		n.table.SetAllowedRowLength(width)
	}

	n.table.SetOutputMirror(os.Stdout)
	n.table.AppendHeader(n.Header())

	n.Config()
}

func (n *Text) Header() (row table.Row) {
	row = append(row, "TYPE")
	row = append(row, "LD")
	row = append(row, "TYPO")

	for _, info := range n.config.Information() {
		for _, headers := range info.Headers() {
			row = append(row, headers)
		}
	}
	return
}

func (n *Text) Row(typo internal.Typo) (row table.Row) {
	if n.config.Verbose() {
		row = append(row, typo.Algorithm().Name())
	} else {
		row = append(row, strings.ToUpper(typo.Algorithm().Id()))
	}
	row = append(row, typo.Ld())
	row = append(row, typo.String())

	for _, info := range n.config.Information() {
		for _, header := range info.Headers() {
			meta := typo.Variant().Meta()
			row = append(row, meta[header])
		}

	}

	return
}

func (n *Text) Config() (row table.Row) {
	n.table.SetStyle(StyleDefault)
	n.table.AppendFooter(table.Row{})

	nameTransformer := text.Transformer(func(val interface{}) string {
		if val.(string) == "MD" {
			return text.Colors{text.BgBlack, text.FgGreen}.Sprint(val)
		}
		return fmt.Sprint(val)
	})

	// n.table.SetRowPainter()

	n.table.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:        "TYPE",
			Align:       text.AlignLeft,
			AlignFooter: text.AlignLeft,
			AlignHeader: text.AlignLeft,
			// Colors:       text.Colors{text.BgBlack, text.FgRed},
			// ColorsHeader: text.Colors{text.BgRed, text.FgBlack, text.Bold},
			// ColorsFooter: text.Colors{text.BgRed, text.FgBlack},
			Hidden:      false,
			Transformer: nameTransformer,
			// TransformerFooter: nameTransformer,
			// TransformerHeader: nameTransformer,
			VAlign:       text.VAlignMiddle,
			VAlignFooter: text.VAlignTop,
			VAlignHeader: text.VAlignBottom,
			WidthMin:     1,
			WidthMax:     64,
		},
		{
			Name:        "ID",
			Align:       text.AlignLeft,
			AlignFooter: text.AlignLeft,
			AlignHeader: text.AlignLeft,
			// Colors:       text.Colors{text.BgBlack, text.FgRed},
			// ColorsHeader: text.Colors{text.BgRed, text.FgBlack, text.Bold},
			// ColorsFooter: text.Colors{text.BgRed, text.FgBlack},
			Hidden:      false,
			Transformer: nameTransformer,
			// TransformerFooter: nameTransformer,
			// TransformerHeader: nameTransformer,
			VAlign:       text.VAlignMiddle,
			VAlignFooter: text.VAlignTop,
			VAlignHeader: text.VAlignBottom,
			WidthMin:     1,
			WidthMax:     3,
		},
	})
	return
}

func (n *Text) Write(in internal.Typo) {
	n.table.AppendRow(n.Row(in))
}

func (n *Text) Summary(report map[string]int64) {
	for k, v := range report {
		fmt.Printf("%s %d   ", k, v)
	}
	fmt.Println("")
}

func (n *Text) Save() {
	// We need a little space between the progress bar and this output
	fmt.Println("")
	output := n.table.Render()

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
