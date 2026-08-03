// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package whois

import (
	"errors"
	"testing"

	parser "github.com/likexian/whois-parser"
	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/observetest"
)

// fakeWhois returns a canned raw record, so the parser is exercised for real
// while the network is not.
type fakeWhois struct {
	raw map[string]string
	err error
}

func (f fakeWhois) Whois(domain string, _ ...string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.raw[domain], nil
}

// A minimal but genuine registrar-format record: enough for the parser to find
// dates, a registrar and a registrant organization.
const registeredRaw = `Domain Name: example.com
Registry Domain ID: 2336799_DOMAIN_COM-VRSN
Creation Date: 1995-08-14T04:00:00Z
Updated Date: 2023-08-14T07:01:44Z
Registry Expiry Date: 2024-08-13T04:00:00Z
Registrar: RESERVED-Internet Assigned Numbers Authority
Registrant Organization: Example Holdings Inc
Registrant Country: US
Name Server: A.IANA-SERVERS.NET
`

const unregisteredRaw = `No match for "NOSUCHDOMAIN-XYZ.COM".
`

// TestWhoisEmitsRegistrantAndDateProps: the registrant is an entity several
// domains converge on, so it is a node; the dates are not, so they are props.
func TestWhoisEmitsRegistrantAndDateProps(t *testing.T) {
	client := fakeWhois{raw: map[string]string{"example.com": registeredRaw}}
	g := observetest.Run(t, observe.TypeDomain, "example.com", newWhois(observe.Options{}, client))

	if !observetest.HasNode(g, observe.TypeRegistrant, "example holdings inc") {
		t.Fatalf("no registrant node; graph has %s", observetest.Dump(g))
	}
	if !observetest.HasEdge(g, observe.TypeDomain, "example.com", observe.RelRegisteredBy, observe.TypeRegistrant, "example holdings inc") {
		t.Error("no REGISTERED_BY edge")
	}
	if v, ok := observetest.Prop(t, g, observe.TypeDomain, "example.com", observe.FieldCreated); !ok || v.Time().Year() != 1995 {
		t.Errorf("created = %v (set=%v), want 1995", v.Time(), ok)
	}
	if v, ok := observetest.Prop(t, g, observe.TypeDomain, "example.com", observe.FieldRegistrar); !ok || v.Str() == "" {
		t.Errorf("registrar = %q (set=%v), want a registrar", v.Str(), ok)
	}
	// A registration date has no identity worth converging on, so it must not
	// have become a node of its own.
	if n := len(g.Nodes()); n != 2 {
		t.Errorf("graph has %d nodes, want the domain and the registrant only: %s", n, observetest.Dump(g))
	}
	observetest.WantStatus(t, g, observe.TypeDomain, "example.com", "whois", graph.StatusOK)
}

// TestWhoisNotFoundIsEmpty: "no match" is the registry authoritatively saying
// the name is free. That is the single most valuable answer a squatting scan
// gets and it must never be filed as an error.
func TestWhoisNotFoundIsEmpty(t *testing.T) {
	client := fakeWhois{raw: map[string]string{"nosuchdomain-xyz.com": unregisteredRaw}}
	g := observetest.Run(t, observe.TypeDomain, "nosuchdomain-xyz.com", newWhois(observe.Options{}, client))

	observetest.WantStatus(t, g, observe.TypeDomain, "nosuchdomain-xyz.com", "whois", graph.StatusEmpty)
	if n := len(g.Nodes()); n != 1 {
		t.Errorf("an unregistered domain admitted %d extra nodes", n-1)
	}
}

// TestWhoisTransportErrorIsFailed: an unreachable whois server says nothing
// about the name.
func TestWhoisTransportErrorIsFailed(t *testing.T) {
	client := fakeWhois{err: errors.New("dial tcp: connection refused")}
	g := observetest.Run(t, observe.TypeDomain, "example.com", newWhois(observe.Options{}, client))
	observetest.WantStatus(t, g, observe.TypeDomain, "example.com", "whois", graph.StatusFailed)
}

// TestWhoisRateLimitIsFailedNotEmpty: a throttled query is the failure mode
// most likely to be mistaken for a free name during a wide scan.
func TestWhoisRateLimitIsFailedNotEmpty(t *testing.T) {
	client := fakeWhois{raw: map[string]string{
		"example.com": "WHOIS LIMIT EXCEEDED - SEE WWW.PIR.ORG/WHOIS FOR DETAILS\n",
	}}
	g := observetest.Run(t, observe.TypeDomain, "example.com", newWhois(observe.Options{}, client))
	observetest.WantStatus(t, g, observe.TypeDomain, "example.com", "whois", graph.StatusFailed)
}

// TestContactNameNormalizes keeps two spellings of one organization from
// becoming two registrant nodes.
func TestContactNameNormalizes(t *testing.T) {
	client := fakeWhois{raw: map[string]string{
		"a.com": "Domain Name: a.com\nRegistrant Organization: ACME  Inc\nCreation Date: 2020-01-01T00:00:00Z\n",
	}}
	g := observetest.Run(t, observe.TypeDomain, "a.com", newWhois(observe.Options{}, client))
	if !observetest.HasNode(g, observe.TypeRegistrant, "acme inc") {
		t.Errorf("registrant not normalized; graph has %s", observetest.Dump(g))
	}
}

// Most registrars answer a redacted contact by filling the field with a
// placeholder rather than omitting it. Taking that at face value gave every
// redacted domain the same registrant node, and the campaign analyzer -- which
// clusters on REGISTERED_BY -- reported unrelated domains as one operation.
func TestContactNameRejectsRedactionPlaceholders(t *testing.T) {
	for _, in := range []string{
		"REDACTED FOR PRIVACY",
		"Redacted for privacy",
		"redacted",
		"Data Protected",
		"GDPR Masked",
		"Statutory Masking Enabled",
		"Whois Privacy Service",
		"Perfect Privacy, LLC",
		"Domains By Proxy, LLC",
		"Not Disclosed",
		"Anonymize, Inc.",
		"N/A",
		"-",
		"none",
	} {
		if got := contactName(&parser.Contact{Organization: in}); got != "" {
			t.Errorf("contactName(%q) = %q, want no registrant", in, got)
		}
	}
}

// The fix must not swallow real registrants, which are the entire point of the
// registrant node.
func TestContactNameKeepsRealNames(t *testing.T) {
	for _, in := range []string{
		"ACME Inc",
		"Google LLC",
		"Cloudflare, Inc.",
		"Wikimedia Foundation, Inc.",
	} {
		if got := contactName(&parser.Contact{Organization: in}); got != in {
			t.Errorf("contactName(%q) = %q, want it kept", in, got)
		}
	}
}

// The organization is preferred, but a redacted organization must fall through
// to the personal name rather than stopping at the placeholder.
func TestContactNameFallsThroughARedactedOrganization(t *testing.T) {
	c := &parser.Contact{Organization: "REDACTED FOR PRIVACY", Name: "ACME Inc"}
	if got := contactName(c); got != "ACME Inc" {
		t.Errorf("contactName = %q, want the personal name", got)
	}
}
