// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"reflect"
	"testing"

	"github.com/urfave/cli/v2"
)

func testCmds() []*cli.Command {
	return []*cli.Command{{
		Name:    "typo",
		Aliases: []string{"t"},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "output", Aliases: []string{"o"}},
			&cli.StringSliceFlag{Name: "algorithm", Aliases: []string{"a"}},
			&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}},
		},
	}}
}

func TestReorderMovesFlagsAheadOfPositionals(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   []string
		want []string
	}{
		{
			"value flag after the target",
			[]string{"urlinsane", "typo", "example.com", "-o", "json"},
			[]string{"urlinsane", "typo", "-o", "json", "example.com"},
		},
		{
			// The value must travel with its flag, or "co" becomes a target.
			"value is not mistaken for a positional",
			[]string{"urlinsane", "typo", "example.com", "-a", "co"},
			[]string{"urlinsane", "typo", "-a", "co", "example.com"},
		},
		{
			"bool flag consumes nothing",
			[]string{"urlinsane", "typo", "example.com", "-v"},
			[]string{"urlinsane", "typo", "-v", "example.com"},
		},
		{
			"a bool flag must not swallow the target",
			[]string{"urlinsane", "typo", "-v", "example.com"},
			[]string{"urlinsane", "typo", "-v", "example.com"},
		},
		{
			"--flag=value carries its own value",
			[]string{"urlinsane", "typo", "example.com", "--output=json"},
			[]string{"urlinsane", "typo", "--output=json", "example.com"},
		},
		{
			"scope and target keep their order",
			[]string{"urlinsane", "typo", "username", "acme.com/bob", "-v"},
			[]string{"urlinsane", "typo", "-v", "username", "acme.com/bob"},
		},
		{
			"already in order is unchanged",
			[]string{"urlinsane", "typo", "-o", "json", "example.com"},
			[]string{"urlinsane", "typo", "-o", "json", "example.com"},
		},
		{
			"aliases resolve to the same command",
			[]string{"urlinsane", "t", "example.com", "-v"},
			[]string{"urlinsane", "t", "-v", "example.com"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := reorder(testCmds(), tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("reorder(%v)\n = %v\nwant %v", tc.in, got, tc.want)
			}
		})
	}
}

// Everything after "--" is positional, and the separator must SURVIVE: without
// it the escaped arguments arrive after the flags with nothing marking them as
// positional, and cli parses "-weird.com" as an undefined flag. This test
// previously asserted the separator-less output and so encoded the bug.
func TestReorderRespectsTheDoubleDash(t *testing.T) {
	got := reorder(testCmds(), []string{"urlinsane", "typo", "-v", "--", "-weird.com"})
	want := []string{"urlinsane", "typo", "-v", "--", "-weird.com"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// A global flag before the command must not disable reordering. Looking only at
// args[1] made this a no-op and reintroduced the error reorder exists to
// prevent.
func TestReorderFindsTheCommandAfterGlobalFlags(t *testing.T) {
	got := reorder(testCmds(), []string{"urlinsane", "--debug", "typo", "example.com", "-o", "json"})
	want := []string{"urlinsane", "--debug", "typo", "-o", "json", "example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
}

// An unknown command is left alone: cli must produce its own "no such command"
// rather than a confusing error about arguments it never saw.
func TestReorderLeavesUnknownCommandsAlone(t *testing.T) {
	in := []string{"urlinsane", "nosuch", "a", "-b", "c"}
	if got := reorder(testCmds(), in); !reflect.DeepEqual(got, in) {
		t.Errorf("got %v, want it untouched", got)
	}
}

func TestReorderHandlesShortInput(t *testing.T) {
	for _, in := range [][]string{nil, {"urlinsane"}} {
		if got := reorder(testCmds(), in); !reflect.DeepEqual(got, in) {
			t.Errorf("reorder(%v) = %v", in, got)
		}
	}
}
