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
	"errors"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
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
	g := run(t, TypeDomain, "example.com", newWhois(Options{}, client))

	if !hasNode(g, TypeRegistrant, "example holdings inc") {
		t.Fatalf("no registrant node; graph has %s", dump(g))
	}
	if !hasEdge(g, TypeDomain, "example.com", RelRegisteredBy, TypeRegistrant, "example holdings inc") {
		t.Error("no REGISTERED_BY edge")
	}
	if v, ok := prop(t, g, TypeDomain, "example.com", FieldCreated); !ok || v.Time().Year() != 1995 {
		t.Errorf("created = %v (set=%v), want 1995", v.Time(), ok)
	}
	if v, ok := prop(t, g, TypeDomain, "example.com", FieldRegistrar); !ok || v.Str() == "" {
		t.Errorf("registrar = %q (set=%v), want a registrar", v.Str(), ok)
	}
	// A registration date has no identity worth converging on, so it must not
	// have become a node of its own.
	if n := len(g.Nodes()); n != 2 {
		t.Errorf("graph has %d nodes, want the domain and the registrant only: %s", n, dump(g))
	}
	wantStatus(t, g, TypeDomain, "example.com", "whois", graph.StatusOK)
}

// TestWhoisNotFoundIsEmpty: "no match" is the registry authoritatively saying
// the name is free. That is the single most valuable answer a squatting scan
// gets and it must never be filed as an error.
func TestWhoisNotFoundIsEmpty(t *testing.T) {
	client := fakeWhois{raw: map[string]string{"nosuchdomain-xyz.com": unregisteredRaw}}
	g := run(t, TypeDomain, "nosuchdomain-xyz.com", newWhois(Options{}, client))

	wantStatus(t, g, TypeDomain, "nosuchdomain-xyz.com", "whois", graph.StatusEmpty)
	if n := len(g.Nodes()); n != 1 {
		t.Errorf("an unregistered domain admitted %d extra nodes", n-1)
	}
}

// TestWhoisTransportErrorIsFailed: an unreachable whois server says nothing
// about the name.
func TestWhoisTransportErrorIsFailed(t *testing.T) {
	client := fakeWhois{err: errors.New("dial tcp: connection refused")}
	g := run(t, TypeDomain, "example.com", newWhois(Options{}, client))
	wantStatus(t, g, TypeDomain, "example.com", "whois", graph.StatusFailed)
}

// TestWhoisRateLimitIsFailedNotEmpty: a throttled query is the failure mode
// most likely to be mistaken for a free name during a wide scan.
func TestWhoisRateLimitIsFailedNotEmpty(t *testing.T) {
	client := fakeWhois{raw: map[string]string{
		"example.com": "WHOIS LIMIT EXCEEDED - SEE WWW.PIR.ORG/WHOIS FOR DETAILS\n",
	}}
	g := run(t, TypeDomain, "example.com", newWhois(Options{}, client))
	wantStatus(t, g, TypeDomain, "example.com", "whois", graph.StatusFailed)
}

// TestContactNameNormalizes keeps two spellings of one organization from
// becoming two registrant nodes.
func TestContactNameNormalizes(t *testing.T) {
	client := fakeWhois{raw: map[string]string{
		"a.com": "Domain Name: a.com\nRegistrant Organization: ACME  Inc\nCreation Date: 2020-01-01T00:00:00Z\n",
	}}
	g := run(t, TypeDomain, "a.com", newWhois(Options{}, client))
	if !hasNode(g, TypeRegistrant, "acme inc") {
		t.Errorf("registrant not normalized; graph has %s", dump(g))
	}
}
