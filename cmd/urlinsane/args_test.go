// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"io"
	"reflect"
	"strings"
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

// An unknown command must fail, not print the help and succeed.
//
// urfave/cli falls through to the app-level Action when it cannot match a
// subcommand, so "no command" and "a command this tool does not have" arrive at
// the same place. Treating both as "show the help, exit 0" made every mistyped
// verb a silently successful run: `urlinsane typoo acme.com --fail-on high`
// printed the help to stdout and exited 0, so a CI gate keyed on the 0/1/2
// contract recorded a clean pass for a scan that never happened.
func TestUnknownCommandIsAnError(t *testing.T) {
	app := newApp()
	app.Writer = io.Discard
	app.ErrWriter = io.Discard
	// urfave calls os.Exit for an ExitCoder by default, which would kill the
	// test binary; this hands the error back instead.
	app.ExitErrHandler = func(*cli.Context, error) {}

	err := app.Run([]string{"urlinsane", "typoo", "acme.com", "--fail-on", "high"})
	if err == nil {
		t.Fatal("an unknown command succeeded; a mistyped verb is a green CI build")
	}
	if !strings.Contains(err.Error(), "typoo") {
		t.Errorf("error does not name the unknown command: %v", err)
	}
	// Exit 1, not the 2 that means "a finding was found".
	var coder cli.ExitCoder
	if errors.As(err, &coder) && coder.ExitCode() != exitError {
		t.Errorf("exit code = %d, want %d", coder.ExitCode(), exitError)
	}
}

// No command at all is not an error: the help is the right answer and 0 is the
// right code.
func TestNoCommandPrintsHelpWithoutFailing(t *testing.T) {
	app := newApp()
	app.Writer = io.Discard
	app.ErrWriter = io.Discard
	app.ExitErrHandler = func(*cli.Context, error) {}

	if err := app.Run([]string{"urlinsane", "--help"}); err != nil {
		t.Fatalf("--help failed: %v", err)
	}
}
