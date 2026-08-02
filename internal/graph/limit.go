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
