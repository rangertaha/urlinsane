// Copyright 2024 Rangertaha. All Rights Reserved.
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

// Package kb provides the keyboard layouts published by kbdlayout.info as a Go
// library: pick a layout, read its keys, and ask what sits next to a character
// or above it on the Shift level.
//
//	board, err := kb.Get("kbdus")
//	board.Adjacent("e")  // [d w r 3 4 s]
//	board.Shifted("4")   // [$]
//
// A layout is identified by the Windows driver it ships in — "kbdus", "kbdfr",
// "kbdgr" — or by any of the KLIDs installed against that driver, so
// Get("00000409") and Get("kbdus") return the same board. ByLanguage selects
// layouts by BCP-47 tag.
//
// Adjacency is geometric rather than tabular. Every key carries the scan code
// of the physical switch it sits on, and scan codes are a property of the
// keyboard rather than of the layout printed on it, so one position table
// serves every layout. Neighbours are then found by distance between key
// strike points, which respects the stagger between rows: on QWERTY the "e" key
// overhangs "d" far more than it overhangs "s", and Adjacent reports them in
// that order. A layout laid out as rows of equal-width columns cannot make
// that distinction, and neither can it tell an ANSI board from an ISO one,
// where the short left Shift puts an extra key below "a".
//
// The dataset covers the alphanumeric block only. The function row and the
// navigation cluster type nothing, and the numeric keypad only duplicates
// digits the number row already has.
package kb

//go:generate protoc --proto_path=internal/kbpb --go_out=. --go_opt=module=github.com/rangertaha/urlinsane/pkg/kb internal/kbpb/kb.proto
//go:generate go run ./gen -out data
