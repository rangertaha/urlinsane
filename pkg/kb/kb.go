// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package kb provides the keyboard layouts published by kbdlayout.info as a Go
// library: pick a layout, read its keys, and ask what sits next to a character
// or above it on the Shift level.
//
//	board := kb.MustGet("kbdus")
//	board.Adjacent("e")     // [w r d 3 4 s]
//	board.Shifted("4")      // [$]
//	board.Translate("hello", kb.MustGet("kbdru"))  // "руддщ"
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
