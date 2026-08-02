// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
// Package whois asks a registry who registered a name.
//
// Its central job is telling apart three things the old collector merged. A
// "no match" is an authoritative determination of *absence* — the most
// valuable answer a squatting scan can get, because an unregistered lookalike
// is one an attacker can still take. A refused or malformed response is
// unknown; so is a timeout. Returning false for all three would report a
// rate-limited registry as a page of free names.
package whois

import (
	"context"
	"errors"
	"strings"
	"time"

	parser "github.com/likexian/whois-parser"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
)

type whoisOp struct {
	observe.Base
	client observe.WhoisClient
}

func newWhois(o observe.Options, c observe.WhoisClient) graph.Operator {
	return whoisOp{Base: o.Base(), client: c}
}

func (whoisOp) Id() string       { return "whois" }
func (whoisOp) Version() int     { return 1 }
func (whoisOp) Resource() string { return observe.ResourceWhois }

func (whoisOp) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{observe.TypeDomain}}}
}

func (whoisOp) Emits() graph.Effects {
	return graph.Effects{
		Nodes: []string{observe.TypeRegistrant},
		Rels:  []string{observe.RelRegisteredBy},
		Props: []string{observe.FieldCreated, observe.FieldUpdated, observe.FieldExpires, observe.FieldRegistrar},
	}
}

func (o whoisOp) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	// The whois client takes no context: its timeout is set on the client at
	// construction. Deriving one here anyway keeps the cancellation path honest
	// for the day the client grows a context-aware call.
	ctx, cancel := o.Call(ctx)
	defer cancel()

	raw, err := o.client.Whois(v.Key())
	if err != nil {
		if ctx.Err() != nil {
			return graph.Delta{}, graph.Timeout(err)
		}
		return graph.Delta{}, graph.Failed(err)
	}

	record, err := parser.Parse(raw)
	if err != nil {
		return graph.Delta{}, whoisParseOutcome(err)
	}

	self := v.Ref()
	var d graph.Delta
	when := func(field string, t *time.Time) {
		if t == nil || t.IsZero() {
			return
		}
		d.Props = append(d.Props, graph.PropSet{Node: &self, Field: field, Value: graph.Time(*t)})
	}
	if record.Domain != nil {
		when(observe.FieldCreated, record.Domain.CreatedDateInTime)
		when(observe.FieldUpdated, record.Domain.UpdatedDateInTime)
		when(observe.FieldExpires, record.Domain.ExpirationDateInTime)
	}
	if name := contactName(record.Registrar); name != "" {
		d.Props = append(d.Props, graph.PropSet{
			Node: &self, Field: observe.FieldRegistrar, Value: graph.String(name),
		})
	}

	// Registrant details are redacted far more often than not, so an absent
	// registrant is normal and must not be mistaken for a failed lookup.
	if name := contactName(record.Registrant); name != "" {
		ref := graph.NodeRef{Type: observe.TypeRegistrant, Key: name}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, graph.EdgeRef{From: self, Rel: observe.RelRegisteredBy, To: ref})
	}

	if len(d.Nodes) == 0 && len(d.Props) == 0 {
		// Nothing usable came out of the parse, so the raw text decides. This
		// is not belt and braces: whois-parser only applies its own not-found
		// detection when it cannot find a domain name anywhere in the response,
		// and a registry's "No match for EXAMPLE.COM" contains one — so Parse
		// returns a nil error and an empty record for the single answer this
		// package most needs to get right.
		if looksUnregistered(raw) {
			return graph.Delta{}, graph.Empty()
		}
		return graph.Delta{}, graph.Failed(errors.New("whois: record carried no registration data"))
	}
	return d, graph.OK()
}

// notFoundMarkers are the phrases registries use to say a name is free. The
// list mirrors whois-parser's own, which it does not export.
var notFoundMarkers = []string{
	"no match", "not found", "no found", "not match", "not available",
	"no data found", "nothing found", "no entries found", "no matching record",
	"not registered", "not been registered", "object does not exist",
	"query returned 0 objects", "domain name not known", "status: free",
	"status: available", "no object found",
}

// looksUnregistered reports whether a raw whois response is a registry saying
// the name is free. Only consulted when the parse yielded nothing, so a
// registered domain with an unlucky field value cannot be mistaken for absent.
func looksUnregistered(raw string) bool {
	lower := strings.ToLower(raw)
	for _, marker := range notFoundMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// whoisParseOutcome maps a parse error onto the status taxonomy. Every variant
// of "you cannot have this name" is an authoritative determination about the
// registry's state, so it is Empty; only a rate limit or unreadable data is a
// failure, and only those two are worth retrying.
func whoisParseOutcome(err error) graph.Outcome {
	switch {
	case errors.Is(err, parser.ErrNotFoundDomain),
		errors.Is(err, parser.ErrReservedDomain),
		errors.Is(err, parser.ErrPremiumDomain),
		errors.Is(err, parser.ErrBlockedDomain):
		return graph.Empty()
	}
	return graph.Failed(err)
}

// contactName picks the stable identity out of a whois contact: the
// organization if there is one, the personal name otherwise. Both are
// whitespace-normalized so "ACME  Inc" and "ACME Inc" converge on one node.
func contactName(c *parser.Contact) string {
	if c == nil {
		return ""
	}
	for _, candidate := range []string{c.Organization, c.Name} {
		if name := strings.Join(strings.Fields(candidate), " "); name != "" {
			return name
		}
	}
	return ""
}

// New builds the registration-lookup operator.
func New(o observe.Options, c observe.WhoisClient) graph.Operator { return newWhois(o, c) }
