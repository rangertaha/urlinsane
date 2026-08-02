// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

// Package observe holds the observation operators — the former collectors.
// Each one looks something up against an external service and expands the graph
// with what it found: nodes, observation edges and props.
//
// Two rules separate this package from the collectors it replaces.
//
// Operators bind to a node *type* through a data pattern, never to a producing
// operator. The old collectors carried a Dependencies() list, so geo declared
// that it needed the ip collector to have run; that broke the moment a second
// collector also produced addresses, or ip was deselected. Here geo binds
// On: ip and simply runs wherever an address exists, whoever put it there. No
// operator in this package declares an ordering of any kind.
//
// How a lookup failed is the finding, not a log line. NXDOMAIN is an
// authoritative determination of absence and returns graph.Empty(); SERVFAIL is
// a broken lookup and returns graph.Failed(); a deadline returns
// graph.Timeout(), meaning nothing was learned. Collapsing the three discards
// exactly the signal a squatting scanner exists to collect. See docs/DESIGN.md
// §4.1 and §7.
package observe

import (
	"context"
	"time"

	"github.com/rangertaha/urlinsane/internal/plugins/decompose"
)

// Node types this package reads and emits.
//
// The names are taken from the package that registers them rather than
// restated. Two definitions of "domain" is precisely how convergence breaks,
// and a typo in a duplicated literal would not fail to compile — it would fail
// at admission, one operator at a time, in production.
const (
	TypeDomain     = decompose.TypeDomain
	TypeIP         = decompose.TypeIP
	TypeRegistrant = decompose.TypeRegistrant
	TypePlatform   = decompose.TypePlatform
	TypeUsername   = decompose.TypeUsername
	TypePackage    = decompose.TypePackage
	TypeRepo       = decompose.TypeRepo
)

// Observation relations. Every one of these costs a depth hop, because every
// one of them required a network call (DESIGN §1.1).
//
// A nameserver is a domain reached by an NS edge, not a type of its own — a
// separate type would split ns1.example.com into two nodes and destroy the
// convergence the whole design rests on. Mail hosts are the same, by MX.
const (
	RelResolvesTo   = decompose.RelResolvesTo
	RelNS           = decompose.RelNS
	RelMX           = decompose.RelMX
	RelPTRTo        = decompose.RelPTRTo
	RelRegisteredBy = decompose.RelRegisteredBy
	RelExistsOn     = decompose.RelExistsOn
)

// Props asserted on a domain node. Registration dates are props rather than
// nodes because they have no identity worth converging on (DESIGN §2).
const (
	FieldUnicode   = "unicode"
	FieldCNAME     = "cname"
	FieldTXT       = "txt"
	FieldCreated   = decompose.FieldCreated
	FieldUpdated   = "updated"
	FieldExpires   = "expires"
	FieldRegistrar = "registrar"
)

// Props asserted on an ip node. A location is a prop for the same reason a
// registration date is: nothing else in the graph needs to converge on it.
const (
	FieldIPVersion   = "version"
	FieldCity        = "city"
	FieldCountry     = "country"
	FieldCountryCode = "country_code"
	FieldContinent   = "continent"
	FieldLatitude    = "latitude"
	FieldLongitude   = "longitude"
	FieldTimeZone    = "timezone"
	FieldPostalCode  = "postal_code"
)

// FieldCode is the registry's short name on a platform node — "npm", "pypi".
// The node's key is its host, so without this the dependency-confusion analysis
// would have to map registry.npmjs.org back to a registry by hand.
const FieldCode = "code"

// Props asserted on observation edges.
const (
	// FieldPreference is the MX record's preference value.
	FieldPreference = "preference"
	// FieldURL is the human-facing page an EXISTS_ON hit was found at, which
	// is not always the URL the existence check itself used.
	FieldURL = "url"
)

// Rate-limit classes (DESIGN §6.3). One token bucket per class, so the limit
// protecting whois does not throttle DNS to the same crawl.
const (
	ResourceDNS   = "dns"
	ResourceWhois = "whois"
	ResourceHTTP  = "http"
	// ResourceGeo covers the local MaxMind database. It is a real class so the
	// knob exists if geo is ever backed by a remote service; with no interval
	// configured, acquiring it is free.
	ResourceGeo = "geo"
)

// DefaultTimeout bounds a single external call. It is deliberately short: a
// round waits for its slowest operator (DESIGN §6.2), so an unbounded lookup
// stalls the whole round.
const DefaultTimeout = 10 * time.Second

// Options supplies the operators' external dependencies. Every one of them is
// an interface so tests can run the whole package without touching the network.
type Options struct {
	// Timeout bounds one external call. Zero uses DefaultTimeout.
	//
	// There is deliberately no Context field: Exec receives one, so a parent
	// stored here would be the wrong parent — fixed when the operator was built
	// rather than scoped to the round that is running it.
	Timeout time.Duration

	// Resolver answers DNS. Nil uses the process resolver, which honours
	// --nameservers.
	Resolver Resolver
	// Whois answers registration queries. Nil uses a real whois client.
	Whois WhoisClient
	// Geo locates addresses. Nil omits the geo operator entirely rather than
	// registering one that can only fail — an operator that cannot possibly
	// succeed should not appear in the compiled plan.
	Geo GeoLocator
	// Prober decides whether a URL exists. Nil uses an HTTP client.
	Prober Prober
	// Sources lists the registries and platforms to probe. Nil uses the
	// dataset; if that is unavailable the source operators are omitted.
	Sources SourceLister
	// SourceResource overrides the rate-limit class of the source-checking
	// operators. Zero uses ResourceHTTP.
	SourceResource string
}

func (o Options) Base() Base {
	t := o.Timeout
	if t <= 0 {
		t = DefaultTimeout
	}
	return Base{timeout: t}
}

// Base is the per-operator lookup timeout.
type Base struct {
	timeout time.Duration
}

// call bounds one external lookup, derived from the context Exec was given.
//
// Deriving from the caller is the whole point: the scheduler's round deadline
// and the interrupt reach the in-flight lookup, so Ctrl-C stops at the round
// boundary rather than waiting out a hung whois (§6.2, §12.4). An earlier
// version rooted this at context.Background() because Exec carried no context —
// which silently discarded every cancellation the engine tried to deliver.
//
// The timeout here still applies: it is the per-operator bound, tighter than
// the scheduler's OpTimeout, and whichever expires first wins.
func (b Base) Call(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	if b.timeout <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, b.timeout)
}
