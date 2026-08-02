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
	"strings"
	"time"

	wclient "github.com/likexian/whois"
	parser "github.com/likexian/whois-parser"
	"github.com/rangertaha/urlinsane/internal/graph"
)

// WhoisClient fetches a raw whois record. It is an interface so tests answer
// from a fixture; the real client's own method has this shape.
type WhoisClient interface {
	Whois(domain string, servers ...string) (string, error)
}

func defaultWhois(timeout time.Duration) WhoisClient {
	c := wclient.NewClient()
	if timeout <= 0 {
		timeout = defaultTimeout
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
type whoisOp struct {
	base
	client WhoisClient
}

func newWhois(o Options, c WhoisClient) graph.Operator {
	return whoisOp{base: o.base(), client: c}
}

func (whoisOp) Id() string       { return "whois" }
func (whoisOp) Version() int     { return 1 }
func (whoisOp) Resource() string { return ResourceWhois }

func (whoisOp) Trigger() graph.Trigger {
	return graph.Trigger{On: graph.Selector{Types: []string{TypeDomain}}}
}

func (whoisOp) Emits() graph.Effects {
	return graph.Effects{
		Nodes: []string{TypeRegistrant},
		Rels:  []string{RelRegisteredBy},
		Props: []string{FieldCreated, FieldUpdated, FieldExpires, FieldRegistrar},
	}
}

func (o whoisOp) Exec(ctx context.Context, v graph.View) (graph.Delta, graph.Outcome) {
	// The whois client takes no context: its timeout is set on the client at
	// construction. Deriving one here anyway keeps the cancellation path honest
	// for the day the client grows a context-aware call.
	ctx, cancel := o.call(ctx)
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
		when(FieldCreated, record.Domain.CreatedDateInTime)
		when(FieldUpdated, record.Domain.UpdatedDateInTime)
		when(FieldExpires, record.Domain.ExpirationDateInTime)
	}
	if name := contactName(record.Registrar); name != "" {
		d.Props = append(d.Props, graph.PropSet{
			Node: &self, Field: FieldRegistrar, Value: graph.String(name),
		})
	}

	// Registrant details are redacted far more often than not, so an absent
	// registrant is normal and must not be mistaken for a failed lookup.
	if name := contactName(record.Registrant); name != "" {
		ref := graph.NodeRef{Type: TypeRegistrant, Key: name}
		d.Nodes = append(d.Nodes, ref)
		d.Edges = append(d.Edges, graph.EdgeRef{From: self, Rel: RelRegisteredBy, To: ref})
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
