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

	"strings"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// dnsOp is the shared skeleton of the DNS operators: same trigger shape, same
// rate-limit class, same error taxonomy. They are separate operators rather
// than one "dns" because status is per (node, operator): folding five record
// types into one would collapse five judgements — a name with an A record but a
// SERVFAIL on MX — into a single misleading one.
//
// The trigger carries no Where. DESIGN §4.1 sketches
// `Where: HasProp(punycode)`, but §2 already canonicalizes a domain key to
// punycode at admission, so the condition is either redundant or, worse, a
// producer dependency wearing a data condition's clothes: it would make every
// DNS operator wait on whichever operator sets that prop, which is the coupling
// §4.1 exists to remove.
type dnsOp struct {
	base
	id  string
	ver int
	on  string
	eff graph.Effects
	res Resolver
}

func (o dnsOp) Id() string           { return o.id }
func (o dnsOp) Version() int         { return o.ver }
func (o dnsOp) Resource() string     { return ResourceDNS }
func (o dnsOp) Emits() graph.Effects { return o.eff }
func (o dnsOp) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{o.on}}}
}

// host normalizes a hostname out of a DNS answer. Answers are fully qualified
// and the graph's keys are not; canonicalization would strip the dot anyway,
// but trimming here means an empty answer is recognisably empty.
func host(s string) string {
	return strings.Trim(strings.TrimSpace(s), ".")
}

// --- A and AAAA -------------------------------------------------------------

// addresses resolves a domain to ip nodes. Addresses are nodes, not props,
// because they have identity worth converging on: two variants sharing an
// address is exactly what campaign clustering looks for.
type addresses struct{ dnsOp }

func newAddresses(o Options, r Resolver) graph.Operator {
	return addresses{dnsOp{
		base: o.base(), id: "dns-a", ver: 1, on: TypeDomain, res: r,
		eff: graph.Effects{
			Nodes: []string{TypeIP},
			Rels:  []string{RelResolvesTo},
			Props: []string{FieldIPVersion},
		},
	}}
}

func (o addresses) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.call(ctx)
	defer cancel()

	ips, err := o.res.LookupIP(ctx, "ip", v.Key())

	var d graph.Delta
	from := v.Ref()
	for _, ip := range ips {
		key := ip.String()
		if key == "" {
			continue
		}
		ref := graph.NodeRef{Type: TypeIP, Key: key}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, graph.EdgeRef{From: from, Rel: RelResolvesTo, To: ref})

		version := int64(6)
		if ip.To4() != nil {
			version = 4
		}
		d.Props = append(d.Props, graph.PropSet{
			Node: &ref, Field: FieldIPVersion, Value: graph.Int(version),
		})
	}
	return d, dnsOutcome(err, len(d.Nodes) > 0)
}

// --- NS ---------------------------------------------------------------------

// nameservers records a domain's authoritative servers. A nameserver is a
// domain node reached by an NS edge, never a type of its own — a separate type
// would split ns1.example.com into two nodes and lose the convergence that
// makes shared-nameserver clustering possible.
type nameservers struct{ dnsOp }

func newNameservers(o Options, r Resolver) graph.Operator {
	return nameservers{dnsOp{
		base: o.base(), id: "dns-ns", ver: 1, on: TypeDomain, res: r,
		eff: graph.Effects{Nodes: []string{TypeDomain}, Rels: []string{RelNS}},
	}}
}

func (o nameservers) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.call(ctx)
	defer cancel()

	records, err := o.res.LookupNS(ctx, v.Key())

	var d graph.Delta
	from := v.Ref()
	for _, rec := range records {
		name := host(rec.Host)
		if name == "" {
			continue
		}
		ref := graph.NodeRef{Type: TypeDomain, Key: name}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, graph.EdgeRef{From: from, Rel: RelNS, To: ref})
	}
	return d, dnsOutcome(err, len(d.Nodes) > 0)
}

// --- MX ---------------------------------------------------------------------

// mailHosts records a domain's mail exchangers, as domain nodes behind an MX
// edge. The preference is an edge prop: it describes the relation, not the host.
type mailHosts struct{ dnsOp }

func newMailHosts(o Options, r Resolver) graph.Operator {
	return mailHosts{dnsOp{
		base: o.base(), id: "dns-mx", ver: 1, on: TypeDomain, res: r,
		eff: graph.Effects{
			Nodes: []string{TypeDomain},
			Rels:  []string{RelMX},
			Props: []string{FieldPreference},
		},
	}}
}

func (o mailHosts) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.call(ctx)
	defer cancel()

	records, err := o.res.LookupMX(ctx, v.Key())

	var d graph.Delta
	from := v.Ref()
	for _, rec := range records {
		name := host(rec.Host)
		if name == "" {
			continue
		}
		ref := graph.NodeRef{Type: TypeDomain, Key: name}
		edge := graph.EdgeRef{From: from, Rel: RelMX, To: ref}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, edge)
		d.Props = append(d.Props, graph.PropSet{
			Edge: &edge, Field: FieldPreference, Value: graph.Int(int64(rec.Pref)),
		})
	}
	return d, dnsOutcome(err, len(d.Nodes) > 0)
}

// --- TXT --------------------------------------------------------------------

// text records a domain's TXT records as a prop. TXT has no identity to
// converge on, so it is not a node; and since props hold a single value, the
// record set is joined with newlines rather than dropped. A repeating
// observation with no entity behind it is the one shape this data model does
// not have a good home for.
type text struct{ dnsOp }

func newText(o Options, r Resolver) graph.Operator {
	return text{dnsOp{
		base: o.base(), id: "dns-txt", ver: 1, on: TypeDomain, res: r,
		eff: graph.Effects{Props: []string{FieldTXT}},
	}}
}

func (o text) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.call(ctx)
	defer cancel()

	records, err := o.res.LookupTXT(ctx, v.Key())

	var kept []string
	for _, rec := range records {
		if rec = strings.TrimSpace(rec); rec != "" {
			kept = append(kept, rec)
		}
	}
	if len(kept) == 0 {
		return graph.Delta{}, dnsOutcome(err, false)
	}
	self := v.Ref()
	return graph.Delta{Props: []graph.PropSet{{
		Node: &self, Field: FieldTXT, Value: graph.String(strings.Join(kept, "\n")),
	}}}, dnsOutcome(err, true)
}

// --- CNAME ------------------------------------------------------------------

// canonicalName records a domain's CNAME target as a prop. There is no CNAME
// relation in the registry, and inventing one would give every alias a second
// path into the graph; the target is recorded as data and the A lookup follows
// the chain anyway.
type canonicalName struct{ dnsOp }

func newCanonicalName(o Options, r Resolver) graph.Operator {
	return canonicalName{dnsOp{
		base: o.base(), id: "dns-cname", ver: 1, on: TypeDomain, res: r,
		eff: graph.Effects{Props: []string{FieldCNAME}},
	}}
}

func (o canonicalName) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.call(ctx)
	defer cancel()

	record, err := o.res.LookupCNAME(ctx, v.Key())
	name := host(record)

	// A resolver answers its own name when there is no alias. That is not a
	// CNAME, and recording it would make every domain look aliased to itself.
	if name == "" || strings.EqualFold(name, v.Key()) {
		return graph.Delta{}, dnsOutcome(err, false)
	}
	self := v.Ref()
	return graph.Delta{Props: []graph.PropSet{{
		Node: &self, Field: FieldCNAME, Value: graph.String(name),
	}}}, dnsOutcome(err, true)
}

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

func newReverse(o Options, r Resolver) graph.Operator {
	return reverse{dnsOp{
		base: o.base(), id: "ptr", ver: 1, on: TypeIP, res: r,
		eff: graph.Effects{Nodes: []string{TypeDomain}, Rels: []string{RelPTRTo}},
	}}
}

func (o reverse) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	ctx, cancel := o.call(ctx)
	defer cancel()

	names, err := o.res.LookupAddr(ctx, v.Key())

	var d graph.Delta
	from := v.Ref()
	for _, n := range names {
		name := host(n)
		if name == "" {
			continue
		}
		ref := graph.NodeRef{Type: TypeDomain, Key: name}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, graph.EdgeRef{From: from, Rel: RelPTRTo, To: ref})
	}
	return d, dnsOutcome(err, len(d.Nodes) > 0)
}
