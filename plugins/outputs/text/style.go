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
	"github.com/jedib0t/go-pretty/v6/table"
	txt "github.com/jedib0t/go-pretty/v6/text"
)

// // Style declares how to render the Table and provides very fine-grained control
// // on how the Table gets rendered on the Console.
// type Style struct {
// 	Name    string        // name of the Style
// 	Box     BoxStyle      // characters to use for the boxes
// 	Color   ColorOptions  // colors to use for the rows and columns
// 	Format  FormatOptions // formatting options for the rows and columns
// 	HTML    HTMLOptions   // rendering options for HTML mode
// 	Options Options       // misc. options for the table
// 	Size    SizeOptions   // size (width) options for the table
// 	Title   TitleOptions  // formation options for the title text
// }

var (
	// StyleDefault renders a Table like below:
	//  +-----+------------+-----------+--------+-----------------------------+
	//  |   # | FIRST NAME | LAST NAME | SALARY |                             |
	//  +-----+------------+-----------+--------+-----------------------------+
	//  |   1 | Arya       | Stark     |   3000 |                             |
	//  |  20 | Jon        | Snow      |   2000 | You know nothing, Jon Snow! |
	//  | 300 | Tyrion     | Lannister |   5000 |                             |
	//  +-----+------------+-----------+--------+-----------------------------+
	//  |     |            | TOTAL     |  10000 |                             |
	//  +-----+------------+-----------+--------+-----------------------------+
	StyleDefault = table.Style{
		Name:    "StyleDefault",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsDefault,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsDefault,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsDefault,
	}

	// StyleBold renders a Table like below:
	//  ┏━━━━━┳━━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
	//  ┃   # ┃ FIRST NAME ┃ LAST NAME ┃ SALARY ┃                             ┃
	//  ┣━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
	//  ┃   1 ┃ Arya       ┃ Stark     ┃   3000 ┃                             ┃
	//  ┃  20 ┃ Jon        ┃ Snow      ┃   2000 ┃ You know nothing, Jon Snow! ┃
	//  ┃ 300 ┃ Tyrion     ┃ Lannister ┃   5000 ┃                             ┃
	//  ┣━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
	//  ┃     ┃            ┃ TOTAL     ┃  10000 ┃                             ┃
	//  ┗━━━━━┻━━━━━━━━━━━━┻━━━━━━━━━━━┻━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
	StyleBold = table.Style{
		Name:    "StyleBold",
		Box:     StyleBoxBold,
		Color:   ColorOptionsDefault,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsDefault,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsDefault,
	}

	// StyleColoredBright renders a Table without any borders or separators,
	// and with Black text on Cyan background for Header/Footer and
	// White background for other rows.
	StyleColoredBright = table.Style{
		Name:    "StyleColoredBright",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsBright,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsDark,
	}

	// StyleColoredDark renders a Table without any borders or separators, and
	// with Header/Footer in Cyan text and other rows with White text, all on
	// Black background.
	StyleColoredDark = table.Style{
		Name:    "StyleColoredDark",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsDark,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsBright,
	}

	// StyleColoredBlackOnBlueWhite renders a Table without any borders or
	// separators, and with Black text on a Blue background for Header/Footer and
	// White background for other rows.
	StyleColoredBlackOnBlueWhite = table.Style{
		Name:    "StyleColoredBlackOnBlueWhite",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsBlackOnBlueWhite,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsBlueOnBlack,
	}

	// StyleColoredBlackOnCyanWhite renders a Table without any borders or
	// separators, and with Black text on a Cyan background for Header/Footer and
	// White background for other rows.
	StyleColoredBlackOnCyanWhite = table.Style{
		Name:    "StyleColoredBlackOnCyanWhite",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsBlackOnCyanWhite,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsCyanOnBlack,
	}

	// StyleColoredBlackOnGreenWhite renders a Table without any borders or
	// separators, and with Black text on a Green background for Header/Footer and
	// White background for other rows.
	StyleColoredBlackOnGreenWhite = table.Style{
		Name:    "StyleColoredBlackOnGreenWhite",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsBlackOnGreenWhite,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsGreenOnBlack,
	}

	// StyleColoredBlackOnMagentaWhite renders a Table without any borders or
	// separators, and with Black text on a Magenta background for the Header/Footer and
	// White background for other rows.
	StyleColoredBlackOnMagentaWhite = table.Style{
		Name:    "StyleColoredBlackOnMagentaWhite",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsBlackOnMagentaWhite,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsMagentaOnBlack,
	}

	// StyleColoredBlackOnYellowWhite renders a Table without any borders or
	// separators, and with Black text on a Yellow background for Header/Footer and
	// White background for other rows.
	StyleColoredBlackOnYellowWhite = table.Style{
		Name:    "StyleColoredBlackOnYellowWhite",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsBlackOnYellowWhite,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsYellowOnBlack,
	}

	// StyleColoredBlackOnRedWhite renders a Table without any borders or
	// separators, and with Black text on a Red background for Header/Footer and
	// White background for other rows.
	StyleColoredBlackOnRedWhite = table.Style{
		Name:    "StyleColoredBlackOnRedWhite",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsBlackOnRedWhite,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsRedOnBlack,
	}

	// StyleColoredBlueWhiteOnBlack renders a Table without any borders or
	// separators, and with Header/Footer in Blue text and other rows with
	// White text, all on a Black background.
	StyleColoredBlueWhiteOnBlack = table.Style{
		Name:    "StyleColoredBlueWhiteOnBlack",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsBlueWhiteOnBlack,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsBlackOnBlue,
	}

	// StyleColoredCyanWhiteOnBlack renders a Table without any borders or
	// separators, and with Header/Footer in Cyan text and other rows with
	// White text, all on a Black background.
	StyleColoredCyanWhiteOnBlack = table.Style{
		Name:    "StyleColoredCyanWhiteOnBlack",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsCyanWhiteOnBlack,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsBlackOnCyan,
	}

	// StyleColoredGreenWhiteOnBlack renders a Table without any borders or
	// separators, and with Header/Footer in Green text and other rows with
	// White text, all on Black background.
	StyleColoredGreenWhiteOnBlack = table.Style{
		Name:    "StyleColoredGreenWhiteOnBlack",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsGreenWhiteOnBlack,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsBlackOnGreen,
	}

	// StyleColoredMagentaWhiteOnBlack renders a Table without any borders or
	// separators, and with Header/Footer in Magenta text and other rows with
	// White text, all on a Black background.
	StyleColoredMagentaWhiteOnBlack = table.Style{
		Name:    "StyleColoredMagentaWhiteOnBlack",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsMagentaWhiteOnBlack,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsBlackOnMagenta,
	}

	// StyleColoredRedWhiteOnBlack renders a Table without any borders or
	// separators, and with Header/Footer in Red text and other rows with
	// White text, all on a Black background.
	StyleColoredRedWhiteOnBlack = table.Style{
		Name:    "StyleColoredRedWhiteOnBlack",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsRedWhiteOnBlack,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsBlackOnRed,
	}

	// StyleColoredYellowWhiteOnBlack renders a Table without any borders or
	// separators, and with Header/Footer in Yellow text and other rows with
	// White text, all on a Black background.
	StyleColoredYellowWhiteOnBlack = table.Style{
		Name:    "StyleColoredYellowWhiteOnBlack",
		Box:     StyleBoxDefault,
		Color:   ColorOptionsYellowWhiteOnBlack,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsNoBordersAndSeparators,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsBlackOnYellow,
	}

	// StyleDouble renders a Table like below:
	//  ╔═════╦════════════╦═══════════╦════════╦═════════════════════════════╗
	//  ║   # ║ FIRST NAME ║ LAST NAME ║ SALARY ║                             ║
	//  ╠═════╬════════════╬═══════════╬════════╬═════════════════════════════╣
	//  ║   1 ║ Arya       ║ Stark     ║   3000 ║                             ║
	//  ║  20 ║ Jon        ║ Snow      ║   2000 ║ You know nothing, Jon Snow! ║
	//  ║ 300 ║ Tyrion     ║ Lannister ║   5000 ║                             ║
	//  ╠═════╬════════════╬═══════════╬════════╬═════════════════════════════╣
	//  ║     ║            ║ TOTAL     ║  10000 ║                             ║
	//  ╚═════╩════════════╩═══════════╩════════╩═════════════════════════════╝
	StyleDouble = table.Style{
		Name:    "StyleDouble",
		Box:     StyleBoxDouble,
		Color:   ColorOptionsDefault,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsDefault,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsDefault,
	}

	// StyleLight renders a Table like below:
	//  ┌─────┬────────────┬───────────┬────────┬─────────────────────────────┐
	//  │   # │ FIRST NAME │ LAST NAME │ SALARY │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │   1 │ Arya       │ Stark     │   3000 │                             │
	//  │  20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow! │
	//  │ 300 │ Tyrion     │ Lannister │   5000 │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │     │            │ TOTAL     │  10000 │                             │
	//  └─────┴────────────┴───────────┴────────┴─────────────────────────────┘
	StyleLight = table.Style{
		Name:    "StyleLight",
		Box:     StyleBoxLight,
		Color:   ColorOptionsDefault,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsDefault,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsDefault,
	}

	// StyleRounded renders a Table like below:
	//  ╭─────┬────────────┬───────────┬────────┬─────────────────────────────╮
	//  │   # │ FIRST NAME │ LAST NAME │ SALARY │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │   1 │ Arya       │ Stark     │   3000 │                             │
	//  │  20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow! │
	//  │ 300 │ Tyrion     │ Lannister │   5000 │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │     │            │ TOTAL     │  10000 │                             │
	//  ╰─────┴────────────┴───────────┴────────┴─────────────────────────────╯
	StyleRounded = table.Style{
		Name:    "StyleRounded",
		Box:     StyleBoxRounded,
		Color:   ColorOptionsDefault,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsDefault,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsDefault,
	}

	// styleTest renders a Table like below:
	//  (-----^------------^-----------^--------^-----------------------------)
	//  [<  #>|<FIRST NAME>|<LAST NAME>|<SALARY>|<                           >]
	//  {-----+------------+-----------+--------+-----------------------------}
	//  [<  1>|<Arya      >|<Stark    >|<  3000>|<                           >]
	//  [< 20>|<Jon       >|<Snow     >|<  2000>|<You know nothing, Jon Snow!>]
	//  [<300>|<Tyrion    >|<Lannister>|<  5000>|<                           >]
	//  {-----+------------+-----------+--------+-----------------------------}
	//  [<   >|<          >|<TOTAL    >|< 10000>|<                           >]
	//  \-----v------------v-----------v--------v-----------------------------/
	styleTest = table.Style{
		Name:    "styleTest",
		Box:     styleBoxTest,
		Color:   ColorOptionsDefault,
		Format:  FormatOptionsDefault,
		HTML:    DefaultHTMLOptions,
		Options: OptionsDefault,
		Size:    SizeOptionsDefault,
		Title:   TitleOptionsDefault,
	}
)

// // BoxStyle defines the characters/strings to use to render the borders and
// // separators for the Table.
// type BoxStyle struct {
// 	BottomLeft       string
// 	BottomRight      string
// 	BottomSeparator  string
// 	EmptySeparator   string
// 	Left             string
// 	LeftSeparator    string
// 	MiddleHorizontal string
// 	MiddleSeparator  string
// 	MiddleVertical   string
// 	PaddingLeft      string
// 	PaddingRight     string
// 	PageSeparator    string
// 	Right            string
// 	RightSeparator   string
// 	TopLeft          string
// 	TopRight         string
// 	TopSeparator     string
// 	UnfinishedRow    string
// }

var (
	// StyleBoxDefault defines a Boxed-Table like below:
	//  +-----+------------+-----------+--------+-----------------------------+
	//  |   # | FIRST NAME | LAST NAME | SALARY |                             |
	//  +-----+------------+-----------+--------+-----------------------------+
	//  |   1 | Arya       | Stark     |   3000 |                             |
	//  |  20 | Jon        | Snow      |   2000 | You know nothing, Jon Snow! |
	//  | 300 | Tyrion     | Lannister |   5000 |                             |
	//  +-----+------------+-----------+--------+-----------------------------+
	//  |     |            | TOTAL     |  10000 |                             |
	//  +-----+------------+-----------+--------+-----------------------------+
	StyleBoxDefault = table.BoxStyle{
		BottomLeft:       "+",
		BottomRight:      "+",
		BottomSeparator:  "+",
		EmptySeparator:   " ",
		Left:             "|",
		LeftSeparator:    "+",
		MiddleHorizontal: "-",
		MiddleSeparator:  "+",
		MiddleVertical:   "|",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            "|",
		RightSeparator:   "+",
		TopLeft:          "+",
		TopRight:         "+",
		TopSeparator:     "+",
		UnfinishedRow:    " ~",
	}

	// StyleBoxBold defines a Boxed-Table like below:
	//  ┏━━━━━┳━━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
	//  ┃   # ┃ FIRST NAME ┃ LAST NAME ┃ SALARY ┃                             ┃
	//  ┣━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
	//  ┃   1 ┃ Arya       ┃ Stark     ┃   3000 ┃                             ┃
	//  ┃  20 ┃ Jon        ┃ Snow      ┃   2000 ┃ You know nothing, Jon Snow! ┃
	//  ┃ 300 ┃ Tyrion     ┃ Lannister ┃   5000 ┃                             ┃
	//  ┣━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
	//  ┃     ┃            ┃ TOTAL     ┃  10000 ┃                             ┃
	//  ┗━━━━━┻━━━━━━━━━━━━┻━━━━━━━━━━━┻━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
	StyleBoxBold = table.BoxStyle{
		BottomLeft:       "┗",
		BottomRight:      "┛",
		BottomSeparator:  "┻",
		EmptySeparator:   " ",
		Left:             "┃",
		LeftSeparator:    "┣",
		MiddleHorizontal: "━",
		MiddleSeparator:  "╋",
		MiddleVertical:   "┃",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            "┃",
		RightSeparator:   "┫",
		TopLeft:          "┏",
		TopRight:         "┓",
		TopSeparator:     "┳",
		UnfinishedRow:    " ≈",
	}

	// StyleBoxDouble defines a Boxed-Table like below:
	//  ╔═════╦════════════╦═══════════╦════════╦═════════════════════════════╗
	//  ║   # ║ FIRST NAME ║ LAST NAME ║ SALARY ║                             ║
	//  ╠═════╬════════════╬═══════════╬════════╬═════════════════════════════╣
	//  ║   1 ║ Arya       ║ Stark     ║   3000 ║                             ║
	//  ║  20 ║ Jon        ║ Snow      ║   2000 ║ You know nothing, Jon Snow! ║
	//  ║ 300 ║ Tyrion     ║ Lannister ║   5000 ║                             ║
	//  ╠═════╬════════════╬═══════════╬════════╬═════════════════════════════╣
	//  ║     ║            ║ TOTAL     ║  10000 ║                             ║
	//  ╚═════╩════════════╩═══════════╩════════╩═════════════════════════════╝
	StyleBoxDouble = table.BoxStyle{
		BottomLeft:       "╚",
		BottomRight:      "╝",
		BottomSeparator:  "╩",
		EmptySeparator:   " ",
		Left:             "║",
		LeftSeparator:    "╠",
		MiddleHorizontal: "═",
		MiddleSeparator:  "╬",
		MiddleVertical:   "║",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            "║",
		RightSeparator:   "╣",
		TopLeft:          "╔",
		TopRight:         "╗",
		TopSeparator:     "╦",
		UnfinishedRow:    " ≈",
	}

	// StyleBoxLight defines a Boxed-Table like below:
	//  ┌─────┬────────────┬───────────┬────────┬─────────────────────────────┐
	//  │   # │ FIRST NAME │ LAST NAME │ SALARY │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │   1 │ Arya       │ Stark     │   3000 │                             │
	//  │  20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow! │
	//  │ 300 │ Tyrion     │ Lannister │   5000 │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │     │            │ TOTAL     │  10000 │                             │
	//  └─────┴────────────┴───────────┴────────┴─────────────────────────────┘
	StyleBoxLight = table.BoxStyle{
		BottomLeft:       "└",
		BottomRight:      "┘",
		BottomSeparator:  "┴",
		EmptySeparator:   " ",
		Left:             "│",
		LeftSeparator:    "├",
		MiddleHorizontal: "─",
		MiddleSeparator:  "┼",
		MiddleVertical:   "│",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            "│",
		RightSeparator:   "┤",
		TopLeft:          "┌",
		TopRight:         "┐",
		TopSeparator:     "┬",
		UnfinishedRow:    " ≈",
	}

	// StyleBoxRounded defines a Boxed-Table like below:
	//  ╭─────┬────────────┬───────────┬────────┬─────────────────────────────╮
	//  │   # │ FIRST NAME │ LAST NAME │ SALARY │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │   1 │ Arya       │ Stark     │   3000 │                             │
	//  │  20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow! │
	//  │ 300 │ Tyrion     │ Lannister │   5000 │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │     │            │ TOTAL     │  10000 │                             │
	//  ╰─────┴────────────┴───────────┴────────┴─────────────────────────────╯
	StyleBoxRounded = table.BoxStyle{
		BottomLeft:       "╰",
		BottomRight:      "╯",
		BottomSeparator:  "┴",
		EmptySeparator:   " ",
		Left:             "│",
		LeftSeparator:    "├",
		MiddleHorizontal: "─",
		MiddleSeparator:  "┼",
		MiddleVertical:   "│",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            "│",
		RightSeparator:   "┤",
		TopLeft:          "╭",
		TopRight:         "╮",
		TopSeparator:     "┬",
		UnfinishedRow:    " ≈",
	}

	// styleBoxTest defines a Boxed-Table like below:
	//  (-----^------------^-----------^--------^-----------------------------)
	//  [<  #>|<FIRST NAME>|<LAST NAME>|<SALARY>|<                           >]
	//  {-----+------------+-----------+--------+-----------------------------}
	//  [<  1>|<Arya      >|<Stark    >|<  3000>|<                           >]
	//  [< 20>|<Jon       >|<Snow     >|<  2000>|<You know nothing, Jon Snow!>]
	//  [<300>|<Tyrion    >|<Lannister>|<  5000>|<                           >]
	//  {-----+------------+-----------+--------+-----------------------------}
	//  [<   >|<          >|<TOTAL    >|< 10000>|<                           >]
	//  \-----v------------v-----------v--------v-----------------------------/
	styleBoxTest = table.BoxStyle{
		BottomLeft:       "\\",
		BottomRight:      "/",
		BottomSeparator:  "v",
		EmptySeparator:   " ",
		Left:             "[",
		LeftSeparator:    "{",
		MiddleHorizontal: "--",
		MiddleSeparator:  "+",
		MiddleVertical:   "|",
		PaddingLeft:      "<",
		PaddingRight:     ">",
		PageSeparator:    "\n",
		Right:            "]",
		RightSeparator:   "}",
		TopLeft:          "(",
		TopRight:         ")",
		TopSeparator:     "^",
		UnfinishedRow:    " ~~~",
	}
)

var (
	// ColorOptionsDefault defines sensible ANSI color options - basically NONE.
	ColorOptionsDefault = table.ColorOptions{}

	// ColorOptionsBright renders dark text on bright background.
	ColorOptionsBright = ColorOptionsBlackOnCyanWhite

	// ColorOptionsDark renders bright text on dark background.
	ColorOptionsDark = ColorOptionsCyanWhiteOnBlack

	// ColorOptionsBlackOnBlueWhite renders Black text on Blue/White background.
	ColorOptionsBlackOnBlueWhite = table.ColorOptions{
		Footer:       txt.Colors{txt.BgBlue, txt.FgBlack},
		Header:       txt.Colors{txt.BgHiBlue, txt.FgBlack},
		IndexColumn:  txt.Colors{txt.BgHiBlue, txt.FgBlack},
		Row:          txt.Colors{txt.BgHiWhite, txt.FgBlack},
		RowAlternate: txt.Colors{txt.BgWhite, txt.FgBlack},
	}

	// ColorOptionsBlackOnCyanWhite renders Black text on Cyan/White background.
	ColorOptionsBlackOnCyanWhite = table.ColorOptions{
		Footer:       txt.Colors{txt.BgCyan, txt.FgBlack},
		Header:       txt.Colors{txt.BgHiCyan, txt.FgBlack},
		IndexColumn:  txt.Colors{txt.BgHiCyan, txt.FgBlack},
		Row:          txt.Colors{txt.BgHiWhite, txt.FgBlack},
		RowAlternate: txt.Colors{txt.BgWhite, txt.FgBlack},
	}

	// ColorOptionsBlackOnGreenWhite renders Black text on Green/White
	// background.
	ColorOptionsBlackOnGreenWhite = table.ColorOptions{
		Footer:       txt.Colors{txt.BgGreen, txt.FgBlack},
		Header:       txt.Colors{txt.BgHiGreen, txt.FgBlack},
		IndexColumn:  txt.Colors{txt.BgHiGreen, txt.FgBlack},
		Row:          txt.Colors{txt.BgHiWhite, txt.FgBlack},
		RowAlternate: txt.Colors{txt.BgWhite, txt.FgBlack},
	}

	// ColorOptionsBlackOnMagentaWhite renders Black text on Magenta/White
	// background.
	ColorOptionsBlackOnMagentaWhite = table.ColorOptions{
		Footer:       txt.Colors{txt.BgMagenta, txt.FgBlack},
		Header:       txt.Colors{txt.BgHiMagenta, txt.FgBlack},
		IndexColumn:  txt.Colors{txt.BgHiMagenta, txt.FgBlack},
		Row:          txt.Colors{txt.BgHiWhite, txt.FgBlack},
		RowAlternate: txt.Colors{txt.BgWhite, txt.FgBlack},
	}

	// ColorOptionsBlackOnRedWhite renders Black text on Red/White background.
	ColorOptionsBlackOnRedWhite = table.ColorOptions{
		Footer:       txt.Colors{txt.BgRed, txt.FgBlack},
		Header:       txt.Colors{txt.BgHiRed, txt.FgBlack},
		IndexColumn:  txt.Colors{txt.BgHiRed, txt.FgBlack},
		Row:          txt.Colors{txt.BgHiWhite, txt.FgBlack},
		RowAlternate: txt.Colors{txt.BgWhite, txt.FgBlack},
	}

	// ColorOptionsBlackOnYellowWhite renders Black text on Yellow/White
	// background.
	ColorOptionsBlackOnYellowWhite = table.ColorOptions{
		Footer:       txt.Colors{txt.BgYellow, txt.FgBlack},
		Header:       txt.Colors{txt.BgHiYellow, txt.FgBlack},
		IndexColumn:  txt.Colors{txt.BgHiYellow, txt.FgBlack},
		Row:          txt.Colors{txt.BgHiWhite, txt.FgBlack},
		RowAlternate: txt.Colors{txt.BgWhite, txt.FgBlack},
	}

	// ColorOptionsBlueWhiteOnBlack renders Blue/White text on Black background.
	ColorOptionsBlueWhiteOnBlack = table.ColorOptions{
		Footer:       txt.Colors{txt.FgBlue, txt.BgHiBlack},
		Header:       txt.Colors{txt.FgHiBlue, txt.BgHiBlack},
		IndexColumn:  txt.Colors{txt.FgHiBlue, txt.BgHiBlack},
		Row:          txt.Colors{txt.FgHiWhite, txt.BgBlack},
		RowAlternate: txt.Colors{txt.FgWhite, txt.BgBlack},
	}

	// ColorOptionsCyanWhiteOnBlack renders Cyan/White text on Black background.
	ColorOptionsCyanWhiteOnBlack = table.ColorOptions{
		Footer:       txt.Colors{txt.FgCyan, txt.BgHiBlack},
		Header:       txt.Colors{txt.FgHiCyan, txt.BgHiBlack},
		IndexColumn:  txt.Colors{txt.FgHiCyan, txt.BgHiBlack},
		Row:          txt.Colors{txt.FgHiWhite, txt.BgBlack},
		RowAlternate: txt.Colors{txt.FgWhite, txt.BgBlack},
	}

	// ColorOptionsGreenWhiteOnBlack renders Green/White text on Black
	// background.
	ColorOptionsGreenWhiteOnBlack = table.ColorOptions{
		Footer:       txt.Colors{txt.FgGreen, txt.BgHiBlack},
		Header:       txt.Colors{txt.FgHiGreen, txt.BgHiBlack},
		IndexColumn:  txt.Colors{txt.FgHiGreen, txt.BgHiBlack},
		Row:          txt.Colors{txt.FgHiWhite, txt.BgBlack},
		RowAlternate: txt.Colors{txt.FgWhite, txt.BgBlack},
	}

	// ColorOptionsMagentaWhiteOnBlack renders Magenta/White text on Black
	// background.
	ColorOptionsMagentaWhiteOnBlack = table.ColorOptions{
		Footer:       txt.Colors{txt.FgMagenta, txt.BgHiBlack},
		Header:       txt.Colors{txt.FgHiMagenta, txt.BgHiBlack},
		IndexColumn:  txt.Colors{txt.FgHiMagenta, txt.BgHiBlack},
		Row:          txt.Colors{txt.FgHiWhite, txt.BgBlack},
		RowAlternate: txt.Colors{txt.FgWhite, txt.BgBlack},
	}

	// ColorOptionsRedWhiteOnBlack renders Red/White text on a Black background.
	ColorOptionsRedWhiteOnBlack = table.ColorOptions{
		Footer:       txt.Colors{txt.FgRed, txt.BgHiBlack},
		Header:       txt.Colors{txt.FgHiRed, txt.BgHiBlack},
		IndexColumn:  txt.Colors{txt.FgHiRed, txt.BgHiBlack},
		Row:          txt.Colors{txt.FgHiWhite, txt.BgBlack},
		RowAlternate: txt.Colors{txt.FgWhite, txt.BgBlack},
	}

	// ColorOptionsYellowWhiteOnBlack renders Yellow/White text on Black
	// background.
	ColorOptionsYellowWhiteOnBlack = table.ColorOptions{
		Footer:       txt.Colors{txt.FgYellow, txt.BgHiBlack},
		Header:       txt.Colors{txt.FgHiYellow, txt.BgHiBlack},
		IndexColumn:  txt.Colors{txt.FgHiYellow, txt.BgHiBlack},
		Row:          txt.Colors{txt.FgHiWhite, txt.BgBlack},
		RowAlternate: txt.Colors{txt.FgWhite, txt.BgBlack},
	}
)

// FormatOptionsDefault defines sensible formatting options.
var FormatOptionsDefault = table.FormatOptions{
	Footer: txt.FormatUpper,
	Header: txt.FormatUpper,
	Row:    txt.FormatDefault,
}

// DefaultHTMLOptions defines sensible HTML rendering defaults.
var DefaultHTMLOptions = table.HTMLOptions{
	CSSClass:    table.DefaultHTMLCSSClass,
	EmptyColumn: "&nbsp;",
	EscapeText:  true,
	Newline:     "<br/>",
}

// Options defines the global options that determine how the Table is
// rendered.
type Options struct {
	// DoNotColorBordersAndSeparators disables coloring all the borders and row
	// or column separators.
	DoNotColorBordersAndSeparators bool

	// DrawBorder enables or disables drawing the border around the Table.
	// Example of a table where it is disabled:
	//     # │ FIRST NAME │ LAST NAME │ SALARY │
	//  ─────┼────────────┼───────────┼────────┼─────────────────────────────
	//     1 │ Arya       │ Stark     │   3000 │
	//    20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow!
	//   300 │ Tyrion     │ Lannister │   5000 │
	//  ─────┼────────────┼───────────┼────────┼─────────────────────────────
	//       │            │ TOTAL     │  10000 │
	DrawBorder bool

	// SeparateColumns enables or disables drawing borders between columns.
	// Example of a table where it is disabled:
	//  ┌─────────────────────────────────────────────────────────────────┐
	//  │   #  FIRST NAME  LAST NAME  SALARY                              │
	//  ├─────────────────────────────────────────────────────────────────┤
	//  │   1  Arya        Stark        3000                              │
	//  │  20  Jon         Snow         2000  You know nothing, Jon Snow! │
	//  │ 300  Tyrion      Lannister    5000                              │
	//  │                  TOTAL       10000                              │
	//  └─────────────────────────────────────────────────────────────────┘
	SeparateColumns bool

	// SeparateFooter enables or disables drawing the border between the footer and
	// the rows. Example of a table where it is disabled:
	//  ┌─────┬────────────┬───────────┬────────┬─────────────────────────────┐
	//  │   # │ FIRST NAME │ LAST NAME │ SALARY │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │   1 │ Arya       │ Stark     │   3000 │                             │
	//  │  20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow! │
	//  │ 300 │ Tyrion     │ Lannister │   5000 │                             │
	//  │     │            │ TOTAL     │  10000 │                             │
	//  └─────┴────────────┴───────────┴────────┴─────────────────────────────┘
	SeparateFooter bool

	// SeparateHeader enables or disables drawing the border between the header and
	// the rows. Example of a table where it is disabled:
	//  ┌─────┬────────────┬───────────┬────────┬─────────────────────────────┐
	//  │   # │ FIRST NAME │ LAST NAME │ SALARY │                             │
	//  │   1 │ Arya       │ Stark     │   3000 │                             │
	//  │  20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow! │
	//  │ 300 │ Tyrion     │ Lannister │   5000 │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │     │            │ TOTAL     │  10000 │                             │
	//  └─────┴────────────┴───────────┴────────┴─────────────────────────────┘
	SeparateHeader bool

	// SeparateRows enables or disables drawing separators between each row.
	// Example of a table where it is enabled:
	//  ┌─────┬────────────┬───────────┬────────┬─────────────────────────────┐
	//  │   # │ FIRST NAME │ LAST NAME │ SALARY │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │   1 │ Arya       │ Stark     │   3000 │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │  20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow! │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │ 300 │ Tyrion     │ Lannister │   5000 │                             │
	//  ├─────┼────────────┼───────────┼────────┼─────────────────────────────┤
	//  │     │            │ TOTAL     │  10000 │                             │
	//  └─────┴────────────┴───────────┴────────┴─────────────────────────────┘
	SeparateRows bool
}

var (
	// OptionsDefault defines sensible global options.
	OptionsDefault = table.Options{
		// DoNotColorBordersAndSeparators disables coloring all the borders and row
		// or column separators.
		DoNotColorBordersAndSeparators: false,

		// DrawBorder enables or disables drawing the border around the Table.
		// Example of a table where it is disabled:
		//     # │ FIRST NAME │ LAST NAME │ SALARY │
		//  ─────┼────────────┼───────────┼────────┼─────────────────────────────
		//     1 │ Arya       │ Stark     │   3000 │
		//    20 │ Jon        │ Snow      │   2000 │ You know nothing, Jon Snow!
		//   300 │ Tyrion     │ Lannister │   5000 │
		//  ─────┼────────────┼───────────┼────────┼─────────────────────────────
		//       │            │ TOTAL     │  10000 │
		DrawBorder:      false,
		SeparateColumns: false,
		SeparateFooter:  true,
		SeparateHeader:  true,
		SeparateRows:    false,
	}

	// OptionsNoBorders sets up a table without any borders.
	OptionsNoBorders = table.Options{
		DrawBorder:      false,
		SeparateColumns: true,
		SeparateFooter:  true,
		SeparateHeader:  true,
		SeparateRows:    false,
	}

	// OptionsNoBordersAndSeparators sets up a table without any borders or
	// separators.
	OptionsNoBordersAndSeparators = table.Options{
		DrawBorder:      false,
		SeparateColumns: false,
		SeparateFooter:  false,
		SeparateHeader:  false,
		SeparateRows:    false,
	}
)

var (
	// SizeOptionsDefault defines sensible size options - basically NONE.
	SizeOptionsDefault = table.SizeOptions{
		WidthMax: 0,
		WidthMin: 0,
	}
)

var (
	// TitleOptionsDefault defines sensible title options - basically NONE.
	TitleOptionsDefault = table.TitleOptions{}

	// TitleOptionsBright renders Bright Bold text on a Dark background.
	TitleOptionsBright = TitleOptionsBlackOnCyan

	// TitleOptionsDark renders Dark Bold text on the Bright background.
	TitleOptionsDark = TitleOptionsCyanOnBlack

	// TitleOptionsBlackOnBlue renders Black text on a Blue background.
	TitleOptionsBlackOnBlue = table.TitleOptions{
		Colors: append(ColorOptionsBlackOnBlueWhite.Header, txt.Bold),
	}

	// TitleOptionsBlackOnCyan renders Black Bold text on a Cyan background.
	TitleOptionsBlackOnCyan = table.TitleOptions{
		Colors: append(ColorOptionsBlackOnCyanWhite.Header, txt.Bold),
	}

	// TitleOptionsBlackOnGreen renders Black Bold text on Green background.
	TitleOptionsBlackOnGreen = table.TitleOptions{
		Colors: append(ColorOptionsBlackOnGreenWhite.Header, txt.Bold),
	}

	// TitleOptionsBlackOnMagenta renders Black Bold text on a Magenta background.
	TitleOptionsBlackOnMagenta = table.TitleOptions{
		Colors: append(ColorOptionsBlackOnMagentaWhite.Header, txt.Bold),
	}

	// TitleOptionsBlackOnRed renders Black Bold text on a Red background.
	TitleOptionsBlackOnRed = table.TitleOptions{
		Colors: append(ColorOptionsBlackOnRedWhite.Header, txt.Bold),
	}

	// TitleOptionsBlackOnYellow renders Black Bold text on a Yellow background.
	TitleOptionsBlackOnYellow = table.TitleOptions{
		Colors: append(ColorOptionsBlackOnYellowWhite.Header, txt.Bold),
	}

	// TitleOptionsBlueOnBlack renders Blue Bold text on a Black background.
	TitleOptionsBlueOnBlack = table.TitleOptions{
		Colors: append(ColorOptionsBlueWhiteOnBlack.Header, txt.Bold),
	}

	// TitleOptionsCyanOnBlack renders Cyan Bold text on a Black background.
	TitleOptionsCyanOnBlack = table.TitleOptions{
		Colors: append(ColorOptionsCyanWhiteOnBlack.Header, txt.Bold),
	}

	// TitleOptionsGreenOnBlack renders Green Bold text on a Black background.
	TitleOptionsGreenOnBlack = table.TitleOptions{
		Colors: append(ColorOptionsGreenWhiteOnBlack.Header, txt.Bold),
	}

	// TitleOptionsMagentaOnBlack renders Magenta Bold text on a Black background.
	TitleOptionsMagentaOnBlack = table.TitleOptions{
		Colors: append(ColorOptionsMagentaWhiteOnBlack.Header, txt.Bold),
	}

	// TitleOptionsRedOnBlack renders Red Bold text on a Black background.
	TitleOptionsRedOnBlack = table.TitleOptions{
		Colors: append(ColorOptionsRedWhiteOnBlack.Header, txt.Bold),
	}

	// TitleOptionsYellowOnBlack renders Yellow Bold text on a Black background.
	TitleOptionsYellowOnBlack = table.TitleOptions{
		Colors: append(ColorOptionsYellowWhiteOnBlack.Header, txt.Bold),
	}
)
