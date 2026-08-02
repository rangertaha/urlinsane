// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package observe

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/rangertaha/urlinsane/internal/graph"
	idns "github.com/rangertaha/urlinsane/internal/pkg/dns"
)

// Resolver is the DNS surface these operators use. It is an interface, and
// exactly the subset of *net.Resolver they need, so tests answer from a table
// instead of the network.
type Resolver interface {
	LookupIP(ctx context.Context, network, host string) ([]net.IP, error)
	LookupNS(ctx context.Context, host string) ([]*net.NS, error)
	LookupMX(ctx context.Context, host string) ([]*net.MX, error)
	LookupTXT(ctx context.Context, host string) ([]string, error)
	LookupCNAME(ctx context.Context, host string) (string, error)
	LookupAddr(ctx context.Context, addr string) ([]string, error)
}

// SystemResolver returns the process-wide resolver, which --nameservers
// configures.
func SystemResolver() Resolver { return idns.Resolver }

// DNSOutcome maps a lookup onto the status taxonomy of DESIGN §7. This is the
// judgement the whole package exists to make, and the ok/empty split is the one
// that matters: NXDOMAIN is a *successful determination of absence*, so it is
// Empty rather than Failed. A report that could not tell "confirmed absent"
// from "could not determine" would be worthless for squatting work, and a run
// full of timeouts would read as a run full of free names.
//
// found reports whether any record came back. Records with an error still count
// as ok — something positive was learned, whatever else broke.
func DNSOutcome(err error, found bool) graph.Outcome {
	if found {
		return graph.OK()
	}
	if err == nil {
		// NODATA: the name resolves but holds no record of this type. That is
		// an authoritative absence, not a failure.
		return graph.Empty()
	}

	// *net.DNSError is checked before the generic net.Error because it
	// implements it, and only it distinguishes NXDOMAIN from SERVFAIL.
	var derr *net.DNSError
	if errors.As(err, &derr) {
		switch {
		case derr.IsNotFound:
			return graph.Empty()
		case derr.IsTimeout:
			return graph.Timeout(err)
		}
		// SERVFAIL, REFUSED, a malformed response: the lookup itself broke, so
		// nothing is proved about the name.
		return graph.Failed(err)
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return graph.Timeout(err)
	}
	var nerr net.Error
	if errors.As(err, &nerr) && nerr.Timeout() {
		return graph.Timeout(err)
	}
	return graph.Failed(err)
}

// Host normalizes a hostname out of a DNS answer. Answers are fully qualified
// and the graph's keys are not; canonicalization would strip the dot anyway,
// but trimming here means an empty answer is recognisably empty.
func Host(s string) string {
	return strings.Trim(strings.TrimSpace(s), ".")
}
