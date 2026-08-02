// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// exitInterrupted is the conventional shell code for a process killed by
// SIGINT: 128 + 2. It sits outside the documented 0/1/2 contract on purpose —
// an aborted run produced no verdict, so it must not be mistaken for one.
const exitInterrupted = 130

// interrupts wires Ctrl-C to a two-stage stop: the first cancels the context so
// expansion finishes its round and reports what it has, the second aborts.
//
// signal.NotifyContext cannot express this. It cancels on the first signal but
// leaves the signal registered until its stop function runs, so every later
// Ctrl-C is swallowed — and the window where that matters is exactly the one
// where the user wants out, because draining a round with one worker over a
// thousand nodes takes minutes. The only escape was kill -9 from another
// terminal.
//
// The message matters as much as the mechanism: a scan that goes silent after
// Ctrl-C is indistinguishable from one that ignored it.
func interrupts(w io.Writer) (context.Context, context.CancelFunc) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go watch(sig, w, cancel, os.Exit)
	return ctx, func() { signal.Stop(sig); cancel() }
}

// watch implements the two-stage stop. Split out so it can be tested with a
// channel and a fake exit rather than by signalling the test binary.
func watch(sig <-chan os.Signal, w io.Writer, cancel context.CancelFunc, exit func(int)) {
	if _, ok := <-sig; !ok {
		return
	}
	fmt.Fprintln(w, "\n  interrupt: finishing the current round, then reporting what was found")
	fmt.Fprintln(w, "  press Ctrl-C again to abort now")
	cancel()

	if _, ok := <-sig; !ok {
		return
	}
	fmt.Fprintln(w, "  aborted")
	exit(exitInterrupted)
}
