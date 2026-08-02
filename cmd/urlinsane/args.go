// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"strings"

	"github.com/urfave/cli/v2"
)

// reorder moves flags ahead of positional arguments.
//
// urfave/cli v2 stops parsing flags at the first non-flag argument, so
//
//	urlinsane typo example.com -o json
//
// treats "-o" and "json" as two more positionals and fails with "expected at
// most a scope and a target, got 4 arguments" — an error that describes the
// parser's difficulty rather than the user's mistake, for a command line that
// is not a mistake at all. Every other tool a user has typed today accepts it.
//
// Rewriting the arguments before the parser sees them is the smallest fix that
// works for both commands and stays true whatever flags they later grow: the
// flag set is read from the command, so a new flag is handled by having been
// declared, not by being added to a list here.
//
// Everything after a bare "--" is positional, as always.
func reorder(cmds []*cli.Command, args []string) []string {
	if len(args) < 2 {
		return args
	}
	cmd := find(cmds, args[1])
	if cmd == nil {
		return args
	}

	// Only boolean flags stand alone; every other flag consumes the token
	// after it, which must not be mistaken for a positional.
	standalone := map[string]bool{}
	for _, f := range cmd.Flags {
		if _, ok := f.(*cli.BoolFlag); !ok {
			continue
		}
		for _, n := range f.Names() {
			standalone[n] = true
		}
	}

	head := append([]string{}, args[:2]...) // program, command
	rest := args[2:]

	var flags, positional []string
	for i := 0; i < len(rest); i++ {
		a := rest[i]
		switch {
		case a == "--":
			positional = append(positional, rest[i+1:]...)
			i = len(rest)
		case strings.HasPrefix(a, "-") && a != "-":
			flags = append(flags, a)
			// --flag=value carries its own value; --flag value does not.
			name := strings.TrimLeft(a, "-")
			if !strings.Contains(a, "=") && !standalone[name] && i+1 < len(rest) {
				i++
				flags = append(flags, rest[i])
			}
		default:
			positional = append(positional, a)
		}
	}
	return append(append(head, flags...), positional...)
}

func find(cmds []*cli.Command, name string) *cli.Command {
	for _, c := range cmds {
		if c.Name == name {
			return c
		}
		for _, a := range c.Aliases {
			if a == name {
				return c
			}
		}
	}
	return nil
}
