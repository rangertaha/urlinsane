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

package observe

import (
	"context"
	"errors"
	"net"

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

// systemResolver returns the process-wide resolver, which --nameservers
// configures.
func systemResolver() Resolver { return idns.Resolver }

// dnsOutcome maps a lookup onto the status taxonomy of DESIGN §7. This is the
// judgement the whole package exists to make, and the ok/empty split is the one
// that matters: NXDOMAIN is a *successful determination of absence*, so it is
// Empty rather than Failed. A report that could not tell "confirmed absent"
// from "could not determine" would be worthless for squatting work, and a run
// full of timeouts would read as a run full of free names.
//
// found reports whether any record came back. Records with an error still count
// as ok — something positive was learned, whatever else broke.
func dnsOutcome(err error, found bool) graph.Outcome {
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
