// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package kb

import (
	"io"
	"os"
	"strings"
)

// Scale of the drawing. Every key position in the geometry is a multiple of a
// quarter unit, so four columns to the unit puts every edge on a whole column
// and no key has to be nudged to fit.
const (
	asciiCols = 4 // columns per key unit
	asciiRows = 3 // lines per keyboard row, the last shared with the next
)

// Print writes the layout to standard output. It is the quickest way to see
// what a layout actually looks like:
//
//	kb.MustGet("kbdus").Print()
//
// See String for what the drawing shows.
func (l *Layout) Print() {
	_ = l.Fprint(os.Stdout)
}

// Fprint writes the layout to w.
func (l *Layout) Fprint(w io.Writer) error {
	_, err := io.WriteString(w, l.String())
	return err
}

// String renders the layout as a keyboard diagram: one box per key, placed
// where the key physically sits, with the shifted character above the base
// one, as on a keycap.
//
// The offset between rows is the real stagger, not decoration — it is what
// makes "e" nearer to "d" than to "s", and seeing it is usually the quickest
// way to check that a layout loaded the way you expected.
//
// Keys that type nothing are drawn empty, which is how the modifiers and the
// space bar show up. A dead key is bracketed, since pressing it produces no
// character on its own. The drawing assumes every character takes one column,
// so a layout in a full-width script will not line up.
func (l *Layout) String() string {
	if len(l.Keys) == 0 {
		return l.title() + "\n"
	}

	width, height := 0, 0
	for _, k := range l.Keys {
		if r := int((k.Pos.X + k.Pos.W) * asciiCols); r+1 > width {
			width = r + 1
		}
		if b := int(k.Pos.Y)*asciiRows + asciiRows; b+1 > height {
			height = b + 1
		}
	}

	canvas := make([][]rune, height)
	for i := range canvas {
		canvas[i] = []rune(strings.Repeat(" ", width))
	}

	for _, k := range l.Keys {
		left := int(k.Pos.X * asciiCols)
		right := int((k.Pos.X + k.Pos.W) * asciiCols)
		top := int(k.Pos.Y) * asciiRows
		bottom := top + asciiRows

		// Borders. Corners are drawn as "+" so that keys sharing an edge
		// join up without any junction bookkeeping.
		for c := left; c <= right; c++ {
			canvas[top][c] = '-'
			canvas[bottom][c] = '-'
		}
		for r := top + 1; r < bottom; r++ {
			canvas[r][left] = '|'
			canvas[r][right] = '|'
		}
		canvas[top][left], canvas[top][right] = '+', '+'
		canvas[bottom][left], canvas[bottom][right] = '+', '+'

		// The shifted character sits above the base one, as it does on a
		// keycap. Either may be missing.
		put(canvas[top+1], left, right, shiftLabel(k))
		put(canvas[top+2], left, right, baseLabel(k))
	}

	var b strings.Builder
	b.WriteString(l.title())
	b.WriteByte('\n')
	for _, line := range canvas {
		b.WriteString(strings.TrimRight(string(line), " "))
		b.WriteByte('\n')
	}

	return b.String()
}

// title names the layout above the drawing.
func (l *Layout) title() string {
	var b strings.Builder
	b.WriteString(l.ID)
	if l.Name != "" {
		b.WriteString(" — ")
		b.WriteString(l.Name)
	}
	b.WriteString(" (")
	b.WriteString(string(l.Form))
	b.WriteString(")")
	return b.String()
}

// baseLabel is what to print on the lower half of a keycap.
func baseLabel(k Key) string { return label(k, Base) }

// shiftLabel is what to print on the upper half. A key whose shifted form is
// the same as its base one — a digit on the numeric row of some layouts, or
// the space bar — is left blank rather than printed twice.
func shiftLabel(k Key) string {
	if label(k, Shift) == label(k, Base) {
		return ""
	}
	return label(k, Shift)
}

// label renders one modifier state of a key, or "" when it types nothing
// worth drawing. A dead key is shown in brackets, since pressing it produces
// no character on its own.
func label(k Key, m Mod) string {
	o, ok := k.Text(m)
	if !ok || o.Text == "" || o.Text == " " {
		return ""
	}
	if o.Dead {
		return "[" + o.Text + "]"
	}
	return o.Text
}

// put centres s between two borders, clipping it if the key is too narrow.
func put(line []rune, left, right int, s string) {
	if s == "" {
		return
	}

	runes := []rune(s)
	space := right - left - 1
	if len(runes) > space {
		runes = runes[:space]
	}

	start := left + 1 + (space-len(runes))/2
	for i, r := range runes {
		line[start+i] = r
	}
}
