// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package kb

import (
	"math"
	"strings"
)

// Form is the physical shape of the alphanumeric block. The three forms differ
// only in the bottom-left and right-hand edges: ISO and JIS keyboards carry
// extra keys that ANSI does not, which shifts where the Enter and Shift keys
// sit. Everything between Q and M is identical across all three.
type Form string

const (
	// ANSI is the 101/104-key form used in the US and most of Asia. It has a
	// wide left Shift and a backslash key above Enter.
	ANSI Form = "ansi"

	// ISO is the 102/105-key form used across Europe. It splits the left
	// Shift to make room for an extra key, and moves backslash next to Enter.
	ISO Form = "iso"

	// JIS is the 106/109-key Japanese form: ISO's tall Enter plus two extra
	// keys, one on the number row and one on the bottom letter row.
	JIS Form = "jis"
)

// Pos is a key's physical position on the board, measured in key units: one
// unit is the width of a plain letter key. X grows rightwards from the left
// edge of the alphanumeric block, Y grows downwards from the number row. The
// function row is not modelled, since it produces no text.
type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
}

// Center returns the coordinates of the middle of the key. It is where a
// finger aims, and for every ordinary key it is also where Distance measures
// from — see strike for the one case where it is not.
func (p Pos) Center() (x, y float64) {
	return p.X + p.W/2, p.Y
}

// Distance returns the distance between two keys, in key units. Two
// side-by-side letter keys are exactly 1.0 apart.
//
// It is measured between the points a finger would plausibly strike rather
// than between the centers, which for ordinary one-unit keys is the same
// thing. It stops being the same thing for the space bar: that key is six
// units wide, so its center sits under "b" and a center-to-center reading
// would call "b" its only neighbour, when in truth the bar runs the length of
// the bottom row. A wide key is struck wherever the hand already is, and
// Distance treats it that way.
//
// The measurement is symmetric — each key's strike point is taken against the
// other's center, so swapping the arguments only swaps the two terms.
func (p Pos) Distance(q Pos) float64 {
	px, py := p.Center()
	qx, qy := q.Center()
	return math.Hypot(p.strike(qx)-q.strike(px), py-qy)
}

// strike returns the point on a key that a finger reaching from x would land
// on: the center of an ordinary key, and for a wider one the nearest point
// that is still half a key-width inside its edges.
func (p Pos) strike(x float64) float64 {
	cx, _ := p.Center()

	lo, hi := p.X+0.5, p.X+p.W-0.5
	if hi < lo {
		// Narrower than a normal key; there is only one place to hit it.
		return cx
	}

	return math.Min(math.Max(x, lo), hi)
}

// DefaultRadius is how far apart two keys may be and still count as adjacent.
// At 1.3 units a letter key reaches its left and right neighbours (1.0), the
// two keys above and the two below (1.03 to 1.25), the space bar if it runs
// beneath it (1.0), and nothing else. Raising it to 1.9 pulls in the
// diagonal-but-one keys as well.
//
// The stagger is what these numbers are for: on QWERTY the "e" key sits a
// quarter unit left of "d" and three quarters left of "s", so "d" comes back
// as the likelier slip. A layout stored as rows of equal-width columns cannot
// tell the two apart.
const DefaultRadius = 1.3

// Row numbers within the alphanumeric block.
const (
	rowDigit  = 0 // ` 1 2 3 4 5 6 7 8 9 0 - =
	rowUpper  = 1 // Tab Q W E R T Y U I O P [ ]
	rowHome   = 2 // Caps A S D F G H J K L ; '
	rowLower  = 3 // Shift Z X C V B N M , . /
	rowBottom = 4 // Space
)

// common holds the scan codes whose position is the same on every form. Scan
// codes are set 1, as reported by kbdlayout.info, and are a property of the
// physical switch rather than of the layout printed on it: "10" is the key one
// row up and one column right of Caps Lock whether it types Q, A, or ф.
var common = map[string]Pos{
	// Number row.
	"29": {0, rowDigit, 1},
	"02": {1, rowDigit, 1},
	"03": {2, rowDigit, 1},
	"04": {3, rowDigit, 1},
	"05": {4, rowDigit, 1},
	"06": {5, rowDigit, 1},
	"07": {6, rowDigit, 1},
	"08": {7, rowDigit, 1},
	"09": {8, rowDigit, 1},
	"0A": {9, rowDigit, 1},
	"0B": {10, rowDigit, 1},
	"0C": {11, rowDigit, 1},
	"0D": {12, rowDigit, 1},

	// Upper letter row.
	"0F": {0, rowUpper, 1.5}, // Tab
	"10": {1.5, rowUpper, 1},
	"11": {2.5, rowUpper, 1},
	"12": {3.5, rowUpper, 1},
	"13": {4.5, rowUpper, 1},
	"14": {5.5, rowUpper, 1},
	"15": {6.5, rowUpper, 1},
	"16": {7.5, rowUpper, 1},
	"17": {8.5, rowUpper, 1},
	"18": {9.5, rowUpper, 1},
	"19": {10.5, rowUpper, 1},
	"1A": {11.5, rowUpper, 1},
	"1B": {12.5, rowUpper, 1},

	// Home row.
	"3A": {0, rowHome, 1.75}, // Caps Lock
	"1E": {1.75, rowHome, 1},
	"1F": {2.75, rowHome, 1},
	"20": {3.75, rowHome, 1},
	"21": {4.75, rowHome, 1},
	"22": {5.75, rowHome, 1},
	"23": {6.75, rowHome, 1},
	"24": {7.75, rowHome, 1},
	"25": {8.75, rowHome, 1},
	"26": {9.75, rowHome, 1},
	"27": {10.75, rowHome, 1},
	"28": {11.75, rowHome, 1},

	// Lower letter row.
	"2C": {2.25, rowLower, 1},
	"2D": {3.25, rowLower, 1},
	"2E": {4.25, rowLower, 1},
	"2F": {5.25, rowLower, 1},
	"30": {6.25, rowLower, 1},
	"31": {7.25, rowLower, 1},
	"32": {8.25, rowLower, 1},
	"33": {9.25, rowLower, 1},
	"34": {10.25, rowLower, 1},
	"35": {11.25, rowLower, 1},

	"39": {3.75, rowBottom, 6.25}, // Space
}

// geometries holds the finished position table for each form. They are built
// once, at startup, and only ever read: Positioned is called for every key of
// every layout that loads, and rebuilding a fifty-odd entry map each time was
// pure waste.
var geometries = map[Form]map[string]Pos{
	ANSI: geometry(ANSI),
	ISO:  geometry(ISO),
	JIS:  geometry(JIS),
}

// geometry builds the scan-code-to-position table for a form. The three forms
// share everything in common and differ only in the handful of keys below.
func geometry(f Form) map[string]Pos {
	g := make(map[string]Pos, len(common)+4)
	for sc, p := range common {
		g[sc] = p
	}

	switch f {
	case ISO, JIS:
		// A short left Shift, with SC 56 filling the gap it leaves, and the
		// backslash key moved up beside a tall Enter.
		g["2A"] = Pos{0, rowLower, 1.25}
		g["56"] = Pos{1.25, rowLower, 1}
		g["2B"] = Pos{12.75, rowHome, 1}
	default: // ANSI
		g["2A"] = Pos{0, rowLower, 2.25}
		g["2B"] = Pos{13.5, rowUpper, 1.5}
	}

	if f == JIS {
		// Two keys ANSI and ISO do not have: yen on the number row and the
		// ro key squeezed in before the right Shift.
		g["7D"] = Pos{13, rowDigit, 1}
		g["73"] = Pos{12.25, rowLower, 1}
	}

	return g
}

// Positioned reports whether a scan code belongs to the alphanumeric block,
// and returns where it sits on the given form. Scan codes outside the block —
// function keys, the navigation cluster, the numeric keypad — are not
// modelled: they either produce no text or duplicate text the block already
// covers, and including them would put two physical keys behind a digit.
//
// A form this package does not know is read as ANSI, which is the shape the
// others are variations on. Scan codes are matched however they are spelled,
// since hex has two spellings and the source data does not agree with itself
// about which to use.
func Positioned(f Form, sc string) (Pos, bool) {
	g, ok := geometries[f]
	if !ok {
		g = geometries[ANSI]
	}
	p, ok := g[strings.ToUpper(strings.TrimSpace(sc))]
	return p, ok
}
