// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

// The first interrupt cancels, so expansion stops at the round barrier and
// still reports what it found.
func TestFirstInterruptCancels(t *testing.T) {
	sig := make(chan os.Signal, 2)
	ctx, cancel := context.WithCancel(context.Background())
	var out bytes.Buffer

	done := make(chan struct{})
	go func() { watch(sig, &out, cancel, func(int) { t.Error("exited on the first signal") }); close(done) }()

	sig <- os.Interrupt
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("the first interrupt did not cancel the context")
	}
	if !strings.Contains(out.String(), "Ctrl-C again") {
		t.Errorf("the user was not told a second interrupt aborts: %q", out.String())
	}
}

// The second must abort. signal.NotifyContext swallowed it, so during a long
// round drain the process could not be stopped from its own terminal at all.
func TestSecondInterruptAborts(t *testing.T) {
	sig := make(chan os.Signal, 2)
	_, cancel := context.WithCancel(context.Background())
	var out bytes.Buffer

	codes := make(chan int, 1)
	go watch(sig, &out, cancel, func(c int) { codes <- c })

	sig <- os.Interrupt
	sig <- os.Interrupt
	select {
	case got := <-codes:
		if got != exitInterrupted {
			t.Errorf("exit code = %d, want %d (128+SIGINT)", got, exitInterrupted)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the second interrupt was swallowed; the run cannot be aborted")
	}
}

// A clean finish closes the channel; the watcher must return rather than leak.
func TestWatchReturnsWhenTheChannelCloses(t *testing.T) {
	sig := make(chan os.Signal)
	_, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		watch(sig, &bytes.Buffer{}, cancel, func(int) { t.Error("exited on a clean close") })
		close(done)
	}()
	close(sig)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("the signal watcher leaked after a clean finish")
	}
}
