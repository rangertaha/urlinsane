// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package dns

import (
	"context"
	"net"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/observetest"
)

// TestAddressesEmitIPNodesAndEdges is the base case: an A/AAAA answer becomes
// ip nodes joined by RESOLVES_TO, not props on the domain. Addresses are nodes
// because two variants sharing one is the finding.
func TestAddressesEmitIPNodesAndEdges(t *testing.T) {
	res := &observetest.FakeResolver{IPs: map[string][]net.IP{
		"example.com": {net.ParseIP("93.184.216.34"), net.ParseIP("2606:2800:220:1:248:1893:25c8:1946")},
	}}
	g := observetest.Run(t, observe.TypeDomain, "example.com", newAddresses(observe.Options{}, res))

	for _, want := range []string{"93.184.216.34", "2606:2800:220:1:248:1893:25c8:1946"} {
		if !observetest.HasNode(g, observe.TypeIP, want) {
			t.Fatalf("no ip node %s; graph has %s", want, observetest.Dump(g))
		}
		if !observetest.HasEdge(g, observe.TypeDomain, "example.com", observe.RelResolvesTo, observe.TypeIP, want) {
			t.Errorf("no RESOLVES_TO edge to %s", want)
		}
	}
	if v, ok := observetest.Prop(t, g, observe.TypeIP, "93.184.216.34", observe.FieldIPVersion); !ok || v.Num() != 4 {
		t.Errorf("ipv4 version observetest.Prop = %v (set=%v), want 4", v.Num(), ok)
	}
	if v, ok := observetest.Prop(t, g, observe.TypeIP, "2606:2800:220:1:248:1893:25c8:1946", observe.FieldIPVersion); !ok || v.Num() != 6 {
		t.Errorf("ipv6 version observetest.Prop = %v (set=%v), want 6", v.Num(), ok)
	}
	observetest.WantStatus(t, g, observe.TypeDomain, "example.com", "dns-a", graph.StatusOK)
}

// TestNXDOMAINIsEmptyNotFailed is the distinction the whole tool exists to
// collect. An unregistered name is a successful determination of absence: it
// must be Empty, and it must admit nothing.
func TestNXDOMAINIsEmptyNotFailed(t *testing.T) {
	res := &observetest.FakeResolver{Err: observetest.Nxdomain()}
	g := observetest.Run(t, observe.TypeDomain, "nx-example.com", newAddresses(observe.Options{}, res))

	if n := len(g.Nodes()); n != 1 {
		t.Fatalf("graph has %d nodes, want only the seed: %s", n, observetest.Dump(g))
	}
	if n := len(g.Edges()); n != 0 {
		t.Fatalf("graph has %d edges, want none", n)
	}
	observetest.WantStatus(t, g, observe.TypeDomain, "nx-example.com", "dns-a", graph.StatusEmpty)

	// And the operator's own judgement, not just what the scheduler stored.
	d, out := newAddresses(observe.Options{}, res).(addresses).Exec(context.Background(), observetest.StubView{NodeType: observe.TypeDomain, NodeKey: "nx-example.com"})
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
	res := &observetest.FakeResolver{Err: observetest.Servfail()}
	g := observetest.Run(t, observe.TypeDomain, "broken.com", newAddresses(observe.Options{}, res))
	observetest.WantStatus(t, g, observe.TypeDomain, "broken.com", "dns-a", graph.StatusFailed)
}

// TestTimeoutIsTimeout covers both shapes a deadline arrives in: the resolver's
// own timeout flag and a context deadline.
func TestTimeoutIsTimeout(t *testing.T) {
	g := observetest.Run(t, observe.TypeDomain, "slow.com", newAddresses(observe.Options{}, &observetest.FakeResolver{Err: observetest.DNSTimeout()}))
	observetest.WantStatus(t, g, observe.TypeDomain, "slow.com", "dns-a", graph.StatusTimeout)

	res := &observetest.FakeResolver{Err: context.DeadlineExceeded}
	_, out := newAddresses(observe.Options{}, res).(addresses).Exec(context.Background(), observetest.StubView{NodeType: observe.TypeDomain, NodeKey: "slow.com"})
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
	g := observetest.Run(t, observe.TypeDomain, "nodata.com", newMailHosts(observe.Options{}, &observetest.FakeResolver{}))
	observetest.WantStatus(t, g, observe.TypeDomain, "nodata.com", "dns-mx", graph.StatusEmpty)
}

// TestPartialAnswerIsOK: records plus an error is still something learned.
func TestPartialAnswerIsOK(t *testing.T) {
	res := &observetest.FakeResolver{
		IPs: map[string][]net.IP{"partial.com": {net.ParseIP("10.0.0.1")}},
		Err: observetest.Servfail(),
	}
	g := observetest.Run(t, observe.TypeDomain, "partial.com", newAddresses(observe.Options{}, res))
	observetest.WantStatus(t, g, observe.TypeDomain, "partial.com", "dns-a", graph.StatusOK)
}

// TestNameserversAreDomainNodes guards DESIGN §2: a nameserver is a domain
// reached by an NS edge, not a type of its own. A separate type would split
// ns1.example.com into two nodes and destroy convergence.
func TestNameserversAreDomainNodes(t *testing.T) {
	res := &observetest.FakeResolver{NS: map[string][]*net.NS{
		"example.com": {{Host: "ns1.example.com."}, {Host: "ns2.example.com."}},
	}}
	g := observetest.Run(t, observe.TypeDomain, "example.com", newNameservers(observe.Options{}, res))

	for _, want := range []string{"ns1.example.com", "ns2.example.com"} {
		n, ok := g.Node(observetest.NodeID(t, g, observe.TypeDomain, want))
		if !ok {
			t.Fatalf("nameserver %s is not a domain node", want)
		}
		if n.Type.Name() != observe.TypeDomain {
			t.Errorf("nameserver %s has type %s, want domain", want, n.Type.Name())
		}
		if !observetest.HasEdge(g, observe.TypeDomain, "example.com", observe.RelNS, observe.TypeDomain, want) {
			t.Errorf("no NS edge to %s", want)
		}
	}
	observetest.WantStatus(t, g, observe.TypeDomain, "example.com", "dns-ns", graph.StatusOK)
}

// TestMailHostsCarryPreference: the preference describes the relation, not the
// host, so it belongs on the edge.
func TestMailHostsCarryPreference(t *testing.T) {
	res := &observetest.FakeResolver{MX: map[string][]*net.MX{
		"example.com": {{Host: "mail.example.com.", Pref: 10}},
	}}
	g := observetest.Run(t, observe.TypeDomain, "example.com", newMailHosts(observe.Options{}, res))

	if !observetest.HasEdge(g, observe.TypeDomain, "example.com", observe.RelMX, observe.TypeDomain, "mail.example.com") {
		t.Fatalf("no MX edge; graph has %s", observetest.Dump(g))
	}
	if v, ok := observetest.EdgeProp(t, g, observe.RelMX, "mail.example.com", observe.FieldPreference); !ok || v.Num() != 10 {
		t.Errorf("MX preference = %d (set=%v), want 10", v.Num(), ok)
	}
}

// TestTextRecordsAreJoinedProps: TXT has no entity behind it, so it is a observetest.Prop.
// Props hold one value, so the record set is joined rather than truncated.
func TestTextRecordsAreJoinedProps(t *testing.T) {
	res := &observetest.FakeResolver{TXT: map[string][]string{
		"example.com": {"v=spf1 -all", "  ", "google-site-verification=x"},
	}}
	g := observetest.Run(t, observe.TypeDomain, "example.com", newText(observe.Options{}, res))

	v, ok := observetest.Prop(t, g, observe.TypeDomain, "example.com", observe.FieldTXT)
	if !ok {
		t.Fatal("txt observetest.Prop not set")
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
	self := &observetest.FakeResolver{CNAME: map[string]string{"example.com": "example.com."}}
	g := observetest.Run(t, observe.TypeDomain, "example.com", newCanonicalName(observe.Options{}, self))
	if _, ok := observetest.Prop(t, g, observe.TypeDomain, "example.com", observe.FieldCNAME); ok {
		t.Error("self-referential CNAME was recorded")
	}
	observetest.WantStatus(t, g, observe.TypeDomain, "example.com", "dns-cname", graph.StatusEmpty)

	alias := &observetest.FakeResolver{CNAME: map[string]string{"example.com": "cdn.provider.net."}}
	g = observetest.Run(t, observe.TypeDomain, "example.com", newCanonicalName(observe.Options{}, alias))
	if v, ok := observetest.Prop(t, g, observe.TypeDomain, "example.com", observe.FieldCNAME); !ok || v.Str() != "cdn.provider.net" {
		t.Errorf("cname = %q (set=%v), want cdn.provider.net", v.Str(), ok)
	}
}

// TestReverseBindsToIPAndEmitsDomain covers the type-flow cycle
// domain → ip → domain. ptr binds On: ip, so it runs against an address
// whatever produced it — here a bare ip seed with no domain in sight, which the
// old DependsOn("ip") collector could not have handled.
// TestDNSOperatorsDeclareNoConditions: any Where here would be a producer
// dependency in disguise. DESIGN §4.1 sketches Where HasProp(punycode), but §2
// canonicalizes a domain key to punycode at admission, so the condition can
// only couple these operators to whoever sets that prop.
func TestDNSOperatorsDeclareNoConditions(t *testing.T) {
	res := &observetest.FakeResolver{}
	for _, op := range New(observe.Options{}, res) {
		if w := op.Trigger().Where; len(w) != 0 {
			t.Errorf("%s declares %d conditions, want none", op.Id(), len(w))
		}
	}
}
