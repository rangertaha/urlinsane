// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package dns resolves domains: addresses, nameservers, mail hosts, text
// records and canonical names.
//
// Five operators rather than one, because they fail independently. A domain
// whose NS lookup times out while its A record answers has said something
// specific, and a single merged outcome would flatten that to "dns failed".
package dns

import (
	"context"
	"strings"

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
// --- A and AAAA -------------------------------------------------------------

// addresses resolves a domain to ip nodes. Addresses are nodes, not props,
// because they have identity worth converging on: two variants sharing an
// address is exactly what campaign clustering looks for.
type addresses struct{ dnsOp }

func newAddresses(o observe.Options, r observe.Resolver) graph.Operator {
	return addresses{dnsOp{
		Base: o.Base(), id: "dns-a", ver: 1, on: observe.TypeDomain, res: r,
		eff: graph.Effects{
			Nodes: []string{observe.TypeIP},
			Rels:  []string{observe.RelResolvesTo},
			Props: []string{observe.FieldIPVersion},
		},
	}}
}

func (o addresses) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.Call(ctx)
	defer cancel()

	ips, err := o.res.LookupIP(ctx, "ip", v.Key())

	var d graph.Delta
	from := v.Ref()
	for _, ip := range ips {
		key := ip.String()
		if key == "" {
			continue
		}
		ref := graph.NodeRef{Type: observe.TypeIP, Key: key}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, graph.EdgeRef{From: from, Rel: observe.RelResolvesTo, To: ref})

		version := int64(6)
		if ip.To4() != nil {
			version = 4
		}
		d.Props = append(d.Props, graph.PropSet{
			Node: &ref, Field: observe.FieldIPVersion, Value: graph.Int(version),
		})
	}
	return d, observe.DNSOutcome(err, len(d.Nodes) > 0)
}

// --- NS ---------------------------------------------------------------------

// nameservers records a domain's authoritative servers. A nameserver is a
// domain node reached by an NS edge, never a type of its own — a separate type
// would split ns1.example.com into two nodes and lose the convergence that
// makes shared-nameserver clustering possible.
type nameservers struct{ dnsOp }

func newNameservers(o observe.Options, r observe.Resolver) graph.Operator {
	return nameservers{dnsOp{
		Base: o.Base(), id: "dns-ns", ver: 1, on: observe.TypeDomain, res: r,
		eff: graph.Effects{Nodes: []string{observe.TypeDomain}, Rels: []string{observe.RelNS}},
	}}
}

func (o nameservers) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.Call(ctx)
	defer cancel()

	records, err := o.res.LookupNS(ctx, v.Key())

	var d graph.Delta
	from := v.Ref()
	for _, rec := range records {
		name := observe.Host(rec.Host)
		if name == "" {
			continue
		}
		ref := graph.NodeRef{Type: observe.TypeDomain, Key: name}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, graph.EdgeRef{From: from, Rel: observe.RelNS, To: ref})
	}
	return d, observe.DNSOutcome(err, len(d.Nodes) > 0)
}

// --- MX ---------------------------------------------------------------------

// mailHosts records a domain's mail exchangers, as domain nodes behind an MX
// edge. The preference is an edge prop: it describes the relation, not the host.
type mailHosts struct{ dnsOp }

func newMailHosts(o observe.Options, r observe.Resolver) graph.Operator {
	return mailHosts{dnsOp{
		Base: o.Base(), id: "dns-mx", ver: 1, on: observe.TypeDomain, res: r,
		eff: graph.Effects{
			Nodes: []string{observe.TypeDomain},
			Rels:  []string{observe.RelMX},
			Props: []string{observe.FieldPreference},
		},
	}}
}

func (o mailHosts) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.Call(ctx)
	defer cancel()

	records, err := o.res.LookupMX(ctx, v.Key())

	var d graph.Delta
	from := v.Ref()
	for _, rec := range records {
		name := observe.Host(rec.Host)
		if name == "" {
			continue
		}
		ref := graph.NodeRef{Type: observe.TypeDomain, Key: name}
		edge := graph.EdgeRef{From: from, Rel: observe.RelMX, To: ref}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, edge)
		d.Props = append(d.Props, graph.PropSet{
			Edge: &edge, Field: observe.FieldPreference, Value: graph.Int(int64(rec.Pref)),
		})
	}
	return d, observe.DNSOutcome(err, len(d.Nodes) > 0)
}

// --- TXT --------------------------------------------------------------------

// text records a domain's TXT records as a prop. TXT has no identity to
// converge on, so it is not a node; and since props hold a single value, the
// record set is joined with newlines rather than dropped. A repeating
// observation with no entity behind it is the one shape this data model does
// not have a good home for.
type text struct{ dnsOp }

func newText(o observe.Options, r observe.Resolver) graph.Operator {
	return text{dnsOp{
		Base: o.Base(), id: "dns-txt", ver: 1, on: observe.TypeDomain, res: r,
		eff: graph.Effects{Props: []string{observe.FieldTXT}},
	}}
}

func (o text) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.Call(ctx)
	defer cancel()

	records, err := o.res.LookupTXT(ctx, v.Key())

	var kept []string
	for _, rec := range records {
		if rec = strings.TrimSpace(rec); rec != "" {
			kept = append(kept, rec)
		}
	}
	if len(kept) == 0 {
		return graph.Delta{}, observe.DNSOutcome(err, false)
	}
	self := v.Ref()
	return graph.Delta{Props: []graph.PropSet{{
		Node: &self, Field: observe.FieldTXT, Value: graph.String(strings.Join(kept, "\n")),
	}}}, observe.DNSOutcome(err, true)
}

// --- CNAME ------------------------------------------------------------------

// canonicalName records a domain's CNAME target as a prop. There is no CNAME
// relation in the registry, and inventing one would give every alias a second
// path into the graph; the target is recorded as data and the A lookup follows
// the chain anyway.
type canonicalName struct{ dnsOp }

func newCanonicalName(o observe.Options, r observe.Resolver) graph.Operator {
	return canonicalName{dnsOp{
		Base: o.Base(), id: "dns-cname", ver: 1, on: observe.TypeDomain, res: r,
		eff: graph.Effects{Props: []string{observe.FieldCNAME}},
	}}
}

func (o canonicalName) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.Call(ctx)
	defer cancel()

	record, err := o.res.LookupCNAME(ctx, v.Key())
	name := observe.Host(record)

	// A resolver answers its own name when there is no alias. That is not a
	// CNAME, and recording it would make every domain look aliased to itself.
	if name == "" || strings.EqualFold(name, v.Key()) {
		return graph.Delta{}, observe.DNSOutcome(err, false)
	}
	self := v.Ref()
	return graph.Delta{Props: []graph.PropSet{{
		Node: &self, Field: observe.FieldCNAME, Value: graph.String(name),
	}}}, observe.DNSOutcome(err, true)
}

// New builds the five domain-facing DNS operators.
//
// Five rather than one because they fail independently: an NS timeout beside an
// A answer is a specific fact, and one merged outcome would erase it.
func New(o observe.Options, r observe.Resolver) []graph.Operator {
	return []graph.Operator{
		newAddresses(o, r), newNameservers(o, r), newMailHosts(o, r),
		newText(o, r), newCanonicalName(o, r),
	}
}
