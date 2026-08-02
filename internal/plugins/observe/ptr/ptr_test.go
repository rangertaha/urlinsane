// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later
package ptr

import (
	"net"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/dns"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/observetest"
)

// TestReverseBindsToIPAndEmitsDomain covers the type-flow cycle
// domain → ip → domain. ptr binds On: ip, so it runs against an address
// whatever produced it — here a bare ip seed with no domain in sight, which the
// old DependsOn("ip") collector could not have handled.
func TestReverseBindsToIPAndEmitsDomain(t *testing.T) {
	res := &observetest.FakeResolver{Addr: map[string][]string{
		"93.184.216.34": {"example.com.", "www.example.com."},
	}}
	g := observetest.Run(t, observe.TypeIP, "93.184.216.34", newReverse(observe.Options{}, res))

	for _, want := range []string{"example.com", "www.example.com"} {
		if !observetest.HasNode(g, observe.TypeDomain, want) {
			t.Fatalf("ptr did not emit domain %s; graph has %s", want, observetest.Dump(g))
		}
		if !observetest.HasEdge(g, observe.TypeIP, "93.184.216.34", observe.RelPTRTo, observe.TypeDomain, want) {
			t.Errorf("no PTR_TO edge to %s", want)
		}
	}
	observetest.WantStatus(t, g, observe.TypeIP, "93.184.216.34", "ptr", graph.StatusOK)
}

// TestReverseDomainsAreOutsideTheSeedClosure: PTR domains are Nameable, so
// without the closure rule every reverse name on every variant's address would
// become a new variant root — the explosion DESIGN §8 exists to prevent.
func TestReverseDomainsAreOutsideTheSeedClosure(t *testing.T) {
	res := &observetest.FakeResolver{
		IPs:  map[string][]net.IP{"example.com": {net.ParseIP("93.184.216.34")}},
		Addr: map[string][]string{"93.184.216.34": {"shared-hosting.example.net."}},
	}
	g := observetest.Run(t, observe.TypeDomain, "example.com",
		dns.New(observe.Options{}, res)[0], newReverse(observe.Options{}, res))

	id := observetest.NodeID(t, g, observe.TypeDomain, "shared-hosting.example.net")
	if g.InClosure(id) {
		t.Error("a PTR-derived domain entered the seed closure; it could then root variant generation")
	}
	if got := g.Depth(id); got != 2 {
		t.Errorf("depth of the PTR domain = %d, want 2 observation hops", got)
	}
}

// TestPTRDeclaresNoConditions: a Where here would be a producer dependency in
// disguise — the very coupling binding On: ip exists to remove.
func TestPTRDeclaresNoConditions(t *testing.T) {
	op := New(observe.Options{}, &observetest.FakeResolver{})
	if w := op.Trigger().Where; len(w) != 0 {
		t.Errorf("ptr declares %d conditions, want none", len(w))
	}
}
