// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package observe

import (
	"time"

	wclient "github.com/likexian/whois"
)

// WhoisClient fetches a raw whois record. It is an interface so tests answer
// from a fixture; the real client's own method has this shape.
type WhoisClient interface {
	Whois(domain string, servers ...string) (string, error)
}

func DefaultWhois(timeout time.Duration) WhoisClient {
	c := wclient.NewClient()
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	c.SetTimeout(timeout)
	return c
}

// whoisOp reads a domain's registration.
//
// It emits a registrant node — an entity several domains genuinely converge on,
// which is what makes registrant clustering possible — and keeps the dates and
// the registrar as props, because a registration date has no identity to
// converge on (DESIGN §2).
//
// The old collector declared DependsOn("ip"), which was never true: whois has
// nothing to do with addresses. That is the failure mode producer dependencies
// invite — a list nobody rechecks, quietly serialising work that could have run
// in the first round.
