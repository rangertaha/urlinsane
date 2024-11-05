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

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var nameTransformer = text.Transformer(func(val interface{}) string {
	if val.(string) == "MD" {
		return text.Colors{text.BgBlack, text.FgGreen}.Sprint(val)
	}
	return fmt.Sprint(val)
})

var ColumnConfig = []table.ColumnConfig{
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
		VAlign:       text.VAlignTop,
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
		VAlign:       text.VAlignTop, // VAlignMiddle
		VAlignFooter: text.VAlignTop,
		VAlignHeader: text.VAlignBottom,
		WidthMin:     1,
		WidthMax:     3,
	},
	{
		Name:        "A",
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
		VAlign:       text.VAlignTop, // VAlignMiddle
		VAlignFooter: text.VAlignTop,
		VAlignHeader: text.VAlignBottom,
		WidthMin:     10,
		WidthMax:     20,
	},
}
