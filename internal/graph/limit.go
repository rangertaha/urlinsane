// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package graph

import (
	"context"
	"sync"
	"time"
)

// Limiter throttles operator calls per resource class. A single global delay is
// meaningless once one run talks to DNS, whois, npm, PyPI and GitHub at once:
// the limit protecting the strictest service would throttle everything else to
// the same crawl.
//
// This is the operator-declared layer. Per-host limiting belongs beneath it in
// the transport.
type Limiter struct {
	mu       sync.Mutex
	interval map[string]time.Duration
	next     map[string]time.Time
	sleep    func(context.Context, time.Duration)
}

func NewLimiter() *Limiter {
	return &Limiter{
		interval: map[string]time.Duration{},
		next:     map[string]time.Time{},
		sleep:    realSleep,
	}
}

func realSleep(ctx context.Context, d time.Duration) {
	if d <= 0 {
		return
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
	case <-ctx.Done():
	}
}

// Set gives a resource class a minimum interval between calls.
func (l *Limiter) Set(resource string, minInterval time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.interval[resource] = minInterval
}

// SetSleep replaces the sleep function. Tests use it to keep rate limiting
// observable without spending wall-clock time.
func (l *Limiter) SetSleep(fn func(context.Context, time.Duration)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sleep = fn
}

// Acquire blocks until this resource class may be called again.
func (l *Limiter) Acquire(ctx context.Context, resource string) {
	if resource == "" {
		return
	}
	l.mu.Lock()
	iv := l.interval[resource]
	if iv <= 0 {
		l.mu.Unlock()
		return
	}
	now := time.Now()
	at := l.next[resource]
	if at.Before(now) {
		at = now
	}
	l.next[resource] = at.Add(iv)
	wait := at.Sub(now)
	sleep := l.sleep
	l.mu.Unlock()
	sleep(ctx, wait)
}
