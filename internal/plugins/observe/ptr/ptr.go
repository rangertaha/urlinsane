// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package ptr resolves an address back to the names that claim it.
//
// It binds On: ip, so it runs wherever an address exists — whoever produced it.
// ptr never names dns-a, which is why a new operator that emits addresses gets
// reverse lookups for free (§4.1).
package ptr

import (
	"context"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
)

// DNS operator wait on whichever operator sets that prop, which is the coupling
// §4.1 exists to remove.
type dnsOp struct {
	observe.Base
	id  string
	ver int
	on  string
	eff graph.Effects
	res observe.Resolver
}

func (o dnsOp) Id() string           { return o.id }
func (o dnsOp) Version() int         { return o.ver }
func (o dnsOp) Resource() string     { return observe.ResourceDNS }
func (o dnsOp) Emits() graph.Effects { return o.eff }
func (o dnsOp) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{o.on}}}
}

// host normalizes a hostname out of a DNS answer. Answers are fully qualified
// and the graph's keys are not; canonicalization would strip the dot anyway,
// --- PTR --------------------------------------------------------------------

// reverse resolves an address back to the names that claim it. It binds On: ip
// and emits domain — the type-flow cycle domain → ip → domain that makes the
// entity graph cyclic by nature, and correct data rather than a bug. The old
// ptr collector instead walked domain.IPs and declared DependsOn("ip"); binding
// to the address itself means it runs for addresses no DNS operator produced.
//
// Nothing here roots a variant: the domains it emits sit outside the seed
// closure, and the applier refuses a VARIANT_OF edge from them, which is what
// stops every reverse-PTR name on every variant's address becoming a new
// variant root (DESIGN §8).
type reverse struct{ dnsOp }

func newReverse(o observe.Options, r observe.Resolver) graph.Operator {
	return reverse{dnsOp{
		Base: o.Base(), id: "ptr", ver: 1, on: observe.TypeIP, res: r,
		eff: graph.Effects{Nodes: []string{observe.TypeDomain}, Rels: []string{observe.RelPTRTo}},
	}}
}

func (o reverse) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.Call(ctx)
	defer cancel()

	names, err := o.res.LookupAddr(ctx, v.Key())

	var d graph.Delta
	from := v.Ref()
	for _, n := range names {
		name := observe.Host(n)
		if name == "" {
			continue
		}
		ref := graph.NodeRef{Type: observe.TypeDomain, Key: name}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, graph.EdgeRef{From: from, Rel: observe.RelPTRTo, To: ref})
	}
	return d, observe.DNSOutcome(err, len(d.Nodes) > 0)
}

// New builds the reverse-lookup operator.
func New(o observe.Options, r observe.Resolver) graph.Operator { return newReverse(o, r) }
