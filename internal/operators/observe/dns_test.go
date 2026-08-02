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
	"net"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// TestAddressesEmitIPNodesAndEdges is the base case: an A/AAAA answer becomes
// ip nodes joined by RESOLVES_TO, not props on the domain. Addresses are nodes
// because two variants sharing one is the finding.
func TestAddressesEmitIPNodesAndEdges(t *testing.T) {
	res := &fakeResolver{ips: map[string][]net.IP{
		"example.com": {net.ParseIP("93.184.216.34"), net.ParseIP("2606:2800:220:1:248:1893:25c8:1946")},
	}}
	g := run(t, TypeDomain, "example.com", newAddresses(Options{}, res))

	for _, want := range []string{"93.184.216.34", "2606:2800:220:1:248:1893:25c8:1946"} {
		if !hasNode(g, TypeIP, want) {
			t.Fatalf("no ip node %s; graph has %s", want, dump(g))
		}
		if !hasEdge(g, TypeDomain, "example.com", RelResolvesTo, TypeIP, want) {
			t.Errorf("no RESOLVES_TO edge to %s", want)
		}
	}
	if v, ok := prop(t, g, TypeIP, "93.184.216.34", FieldIPVersion); !ok || v.Num() != 4 {
		t.Errorf("ipv4 version prop = %v (set=%v), want 4", v.Num(), ok)
	}
	if v, ok := prop(t, g, TypeIP, "2606:2800:220:1:248:1893:25c8:1946", FieldIPVersion); !ok || v.Num() != 6 {
		t.Errorf("ipv6 version prop = %v (set=%v), want 6", v.Num(), ok)
	}
	wantStatus(t, g, TypeDomain, "example.com", "dns-a", graph.StatusOK)
}

// TestNXDOMAINIsEmptyNotFailed is the distinction the whole tool exists to
// collect. An unregistered name is a successful determination of absence: it
// must be Empty, and it must admit nothing.
func TestNXDOMAINIsEmptyNotFailed(t *testing.T) {
	res := &fakeResolver{err: nxdomain()}
	g := run(t, TypeDomain, "nx-example.com", newAddresses(Options{}, res))

	if n := len(g.Nodes()); n != 1 {
		t.Fatalf("graph has %d nodes, want only the seed: %s", n, dump(g))
	}
	if n := len(g.Edges()); n != 0 {
		t.Fatalf("graph has %d edges, want none", n)
	}
	wantStatus(t, g, TypeDomain, "nx-example.com", "dns-a", graph.StatusEmpty)

	// And the operator's own judgement, not just what the scheduler stored.
	d, out := newAddresses(Options{}, res).(addresses).Exec(context.Background(), stubView{TypeDomain, "nx-example.com"})
	if out.Status != graph.StatusEmpty {
		t.Errorf("outcome = %s, want empty", out.Status)
	}
	if len(d.Nodes) != 0 {
		t.Errorf("delta carried %d nodes on NXDOMAIN, want none", len(d.Nodes))
	}
}

// TestServfailIsFailed: a broken lookup proves nothing about the name, so it
// must not read as absence.
func TestServfailIsFailed(t *testing.T) {
	res := &fakeResolver{err: servfail()}
	g := run(t, TypeDomain, "broken.com", newAddresses(Options{}, res))
	wantStatus(t, g, TypeDomain, "broken.com", "dns-a", graph.StatusFailed)
}

// TestTimeoutIsTimeout covers both shapes a deadline arrives in: the resolver's
// own timeout flag and a context deadline.
func TestTimeoutIsTimeout(t *testing.T) {
	g := run(t, TypeDomain, "slow.com", newAddresses(Options{}, &fakeResolver{err: dnsTimeout()}))
	wantStatus(t, g, TypeDomain, "slow.com", "dns-a", graph.StatusTimeout)

	res := &fakeResolver{err: context.DeadlineExceeded}
	_, out := newAddresses(Options{}, res).(addresses).Exec(context.Background(), stubView{TypeDomain, "slow.com"})
	if out.Status != graph.StatusTimeout {
		t.Errorf("context deadline outcome = %s, want timeout", out.Status)
	}
	if out.Err == nil {
		t.Error("timeout outcome carried no error; the cause is part of the finding")
	}
}

// TestNoDataIsEmpty: the name resolves but holds no record of this type. That
// is an authoritative absence, the same as NXDOMAIN.
func TestNoDataIsEmpty(t *testing.T) {
	g := run(t, TypeDomain, "nodata.com", newMailHosts(Options{}, &fakeResolver{}))
	wantStatus(t, g, TypeDomain, "nodata.com", "dns-mx", graph.StatusEmpty)
}

// TestPartialAnswerIsOK: records plus an error is still something learned.
func TestPartialAnswerIsOK(t *testing.T) {
	res := &fakeResolver{
		ips: map[string][]net.IP{"partial.com": {net.ParseIP("10.0.0.1")}},
		err: servfail(),
	}
	g := run(t, TypeDomain, "partial.com", newAddresses(Options{}, res))
	wantStatus(t, g, TypeDomain, "partial.com", "dns-a", graph.StatusOK)
}

// TestNameserversAreDomainNodes guards DESIGN §2: a nameserver is a domain
// reached by an NS edge, not a type of its own. A separate type would split
// ns1.example.com into two nodes and destroy convergence.
func TestNameserversAreDomainNodes(t *testing.T) {
	res := &fakeResolver{ns: map[string][]*net.NS{
		"example.com": {{Host: "ns1.example.com."}, {Host: "ns2.example.com."}},
	}}
	g := run(t, TypeDomain, "example.com", newNameservers(Options{}, res))

	for _, want := range []string{"ns1.example.com", "ns2.example.com"} {
		n, ok := g.Node(nodeID(t, g, TypeDomain, want))
		if !ok {
			t.Fatalf("nameserver %s is not a domain node", want)
		}
		if n.Type.Name() != TypeDomain {
			t.Errorf("nameserver %s has type %s, want domain", want, n.Type.Name())
		}
		if !hasEdge(g, TypeDomain, "example.com", RelNS, TypeDomain, want) {
			t.Errorf("no NS edge to %s", want)
		}
	}
	wantStatus(t, g, TypeDomain, "example.com", "dns-ns", graph.StatusOK)
}

// TestMailHostsCarryPreference: the preference describes the relation, not the
// host, so it belongs on the edge.
func TestMailHostsCarryPreference(t *testing.T) {
	res := &fakeResolver{mx: map[string][]*net.MX{
		"example.com": {{Host: "mail.example.com.", Pref: 10}},
	}}
	g := run(t, TypeDomain, "example.com", newMailHosts(Options{}, res))

	if !hasEdge(g, TypeDomain, "example.com", RelMX, TypeDomain, "mail.example.com") {
		t.Fatalf("no MX edge; graph has %s", dump(g))
	}
	if v, ok := edgeProp(t, g, RelMX, "mail.example.com", FieldPreference); !ok || v.Num() != 10 {
		t.Errorf("MX preference = %d (set=%v), want 10", v.Num(), ok)
	}
}

// TestTextRecordsAreJoinedProps: TXT has no entity behind it, so it is a prop.
// Props hold one value, so the record set is joined rather than truncated.
func TestTextRecordsAreJoinedProps(t *testing.T) {
	res := &fakeResolver{txt: map[string][]string{
		"example.com": {"v=spf1 -all", "  ", "google-site-verification=x"},
	}}
	g := run(t, TypeDomain, "example.com", newText(Options{}, res))

	v, ok := prop(t, g, TypeDomain, "example.com", FieldTXT)
	if !ok {
		t.Fatal("txt prop not set")
	}
	if v.Str() != "v=spf1 -all\ngoogle-site-verification=x" {
		t.Errorf("txt = %q, want both records joined and the blank dropped", v.Str())
	}
	if n := len(g.Nodes()); n != 1 {
		t.Errorf("txt records created %d nodes, want none beyond the seed", n-1)
	}
}

// TestCNAMEIgnoresSelfReference: resolvers answer with the queried name when
// there is no alias, and recording that would make every domain look aliased to
// itself.
func TestCNAMEIgnoresSelfReference(t *testing.T) {
	self := &fakeResolver{cname: map[string]string{"example.com": "example.com."}}
	g := run(t, TypeDomain, "example.com", newCanonicalName(Options{}, self))
	if _, ok := prop(t, g, TypeDomain, "example.com", FieldCNAME); ok {
		t.Error("self-referential CNAME was recorded")
	}
	wantStatus(t, g, TypeDomain, "example.com", "dns-cname", graph.StatusEmpty)

	alias := &fakeResolver{cname: map[string]string{"example.com": "cdn.provider.net."}}
	g = run(t, TypeDomain, "example.com", newCanonicalName(Options{}, alias))
	if v, ok := prop(t, g, TypeDomain, "example.com", FieldCNAME); !ok || v.Str() != "cdn.provider.net" {
		t.Errorf("cname = %q (set=%v), want cdn.provider.net", v.Str(), ok)
	}
}

// TestReverseBindsToIPAndEmitsDomain covers the type-flow cycle
// domain → ip → domain. ptr binds On: ip, so it runs against an address
// whatever produced it — here a bare ip seed with no domain in sight, which the
// old DependsOn("ip") collector could not have handled.
func TestReverseBindsToIPAndEmitsDomain(t *testing.T) {
	res := &fakeResolver{addr: map[string][]string{
		"93.184.216.34": {"example.com.", "www.example.com."},
	}}
	g := run(t, TypeIP, "93.184.216.34", newReverse(Options{}, res))

	for _, want := range []string{"example.com", "www.example.com"} {
		if !hasNode(g, TypeDomain, want) {
			t.Fatalf("ptr did not emit domain %s; graph has %s", want, dump(g))
		}
		if !hasEdge(g, TypeIP, "93.184.216.34", RelPTRTo, TypeDomain, want) {
			t.Errorf("no PTR_TO edge to %s", want)
		}
	}
	wantStatus(t, g, TypeIP, "93.184.216.34", "ptr", graph.StatusOK)
}

// TestReverseDomainsAreOutsideTheSeedClosure: PTR domains are Nameable, so
// without the closure rule every reverse name on every variant's address would
// become a new variant root — the explosion DESIGN §8 exists to prevent.
func TestReverseDomainsAreOutsideTheSeedClosure(t *testing.T) {
	res := &fakeResolver{
		ips:  map[string][]net.IP{"example.com": {net.ParseIP("93.184.216.34")}},
		addr: map[string][]string{"93.184.216.34": {"shared-hosting.example.net."}},
	}
	g := run(t, TypeDomain, "example.com",
		newAddresses(Options{}, res), newReverse(Options{}, res))

	id := nodeID(t, g, TypeDomain, "shared-hosting.example.net")
	if g.InClosure(id) {
		t.Error("a PTR-derived domain entered the seed closure; it could then root variant generation")
	}
	if got := g.Depth(id); got != 2 {
		t.Errorf("depth of the PTR domain = %d, want 2 observation hops", got)
	}
}

// TestDNSOperatorsDeclareNoConditions: any Where here would be a producer
// dependency in disguise. DESIGN §4.1 sketches Where HasProp(punycode), but §2
// canonicalizes a domain key to punycode at admission, so the condition can
// only couple these operators to whoever sets that prop.
func TestDNSOperatorsDeclareNoConditions(t *testing.T) {
	res := &fakeResolver{}
	for _, op := range []graph.Operator{
		newAddresses(Options{}, res), newNameservers(Options{}, res),
		newMailHosts(Options{}, res), newText(Options{}, res),
		newCanonicalName(Options{}, res), newReverse(Options{}, res),
	} {
		if w := op.Trigger().Where; len(w) != 0 {
			t.Errorf("%s declares %d conditions, want none", op.Id(), len(w))
		}
	}
}
