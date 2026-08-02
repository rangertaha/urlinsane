// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package geo

import (
	"errors"
	"net"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
	"github.com/rangertaha/urlinsane/internal/plugins/observe"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/dns"
	"github.com/rangertaha/urlinsane/internal/plugins/observe/observetest"
)

// fakeGeo answers from a table, so no MaxMind database is needed.
type fakeGeo struct {
	places map[string]*observe.Location
	err    error
}

func (f fakeGeo) Locate(addr string) (*observe.Location, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.places[addr], nil
}

// recordingGeo notes every address it was asked about, which is how a test
// proves geo never ran against a domain.
type recordingGeo struct {
	places map[string]*observe.Location
	seen   *[]string
}

func (r recordingGeo) Locate(addr string) (*observe.Location, error) {
	*r.seen = append(*r.seen, addr)
	return r.places[addr], nil
}

// TestGeoBindsToIPNotItsProducer is the point of DESIGN §4.1. geo runs on the
// address the dns operator produced without declaring anything about the dns
// operator; it is dispatched by the pattern "there is an ip", in the round
// after the address appears.
func TestGeoBindsToIPNotItsProducer(t *testing.T) {
	var seen []string
	locator := recordingGeo{
		places: map[string]*observe.Location{"93.184.216.34": {
			City: "Norwell", Country: "United States", CountryCode: "US",
			Continent: "North America", TimeZone: "America/New_York",
			PostalCode: "02061", Latitude: 42.1596, Longitude: -70.8217,
		}},
		seen: &seen,
	}
	res := &observetest.FakeResolver{IPs: map[string][]net.IP{
		"example.com": {net.ParseIP("93.184.216.34")},
	}}

	g := observetest.Run(t, observe.TypeDomain, "example.com",
		dns.New(observe.Options{}, res)[0], newGeo(observe.Options{}, locator))

	if len(seen) != 1 || seen[0] != "93.184.216.34" {
		t.Fatalf("geo was asked about %v, want exactly the address", seen)
	}
	observetest.WantStatus(t, g, observe.TypeIP, "93.184.216.34", "geo", graph.StatusOK)

	// It never matched the domain: the trigger binds to a type, and domain is
	// not it.
	if _, ok := g.Status(observetest.NodeID(t, g, observe.TypeDomain, "example.com"), "geo"); ok {
		t.Error("geo ran against a domain node")
	}
}

// TestGeoEmitsPropsNotNodes guards DESIGN §2: a location has no identity worth
// converging on, so it is props on the address. A location node would be a new
// node per address with nothing ever joining on it.
func TestGeoEmitsPropsNotNodes(t *testing.T) {
	locator := fakeGeo{places: map[string]*observe.Location{"93.184.216.34": {
		City: "Norwell", Country: "United States", CountryCode: "US",
		Continent: "North America", TimeZone: "America/New_York",
		PostalCode: "02061", Latitude: 42.1596, Longitude: -70.8217,
	}}}
	g := observetest.Run(t, observe.TypeIP, "93.184.216.34", newGeo(observe.Options{}, locator))

	if n := len(g.Nodes()); n != 1 {
		t.Fatalf("geo admitted %d nodes beyond the address; a location is a observetest.Prop: %s", n-1, observetest.Dump(g))
	}
	if n := len(g.Edges()); n != 0 {
		t.Fatalf("geo emitted %d edges, want none", n)
	}
	for field, want := range map[string]string{
		observe.FieldCity:        "Norwell",
		observe.FieldCountry:     "United States",
		observe.FieldCountryCode: "US",
		observe.FieldContinent:   "North America",
		observe.FieldTimeZone:    "America/New_York",
		observe.FieldPostalCode:  "02061",
	} {
		if v, ok := observetest.Prop(t, g, observe.TypeIP, "93.184.216.34", field); !ok || v.Str() != want {
			t.Errorf("%s = %q (set=%v), want %q", field, v.Str(), ok, want)
		}
	}
	if v, ok := observetest.Prop(t, g, observe.TypeIP, "93.184.216.34", observe.FieldLatitude); !ok || v.Real() != 42.1596 {
		t.Errorf("latitude = %v (set=%v), want 42.1596", v.Real(), ok)
	}

	// The operator declares props and no node type at all, which is why Effects
	// covers props: plan compilation would otherwise not see it produce
	// anything.
	eff := newGeo(observe.Options{}, locator).Emits()
	if len(eff.Nodes) != 0 || len(eff.Rels) != 0 || len(eff.Props) == 0 {
		t.Errorf("geo effects = %+v, want props only", eff)
	}
}

// TestGeoUnknownAddressIsEmpty: the database answered and has nothing. That is
// absence, and it must not be confused with a database that failed to open.
func TestGeoUnknownAddressIsEmpty(t *testing.T) {
	g := observetest.Run(t, observe.TypeIP, "10.0.0.1", newGeo(observe.Options{}, fakeGeo{}))
	observetest.WantStatus(t, g, observe.TypeIP, "10.0.0.1", "geo", graph.StatusEmpty)
}

// TestGeoLookupErrorIsFailed keeps a broken database from reading as a world
// with no locations in it.
func TestGeoLookupErrorIsFailed(t *testing.T) {
	g := observetest.Run(t, observe.TypeIP, "10.0.0.1", newGeo(observe.Options{}, fakeGeo{err: errors.New("corrupt database")}))
	observetest.WantStatus(t, g, observe.TypeIP, "10.0.0.1", "geo", graph.StatusFailed)
}

// TestGeoIgnoresNullIsland: MaxMind uses 0,0 for "unknown", and recording it
// would place half the internet in the Gulf of Guinea.
func TestGeoIgnoresNullIsland(t *testing.T) {
	locator := fakeGeo{places: map[string]*observe.Location{
		"10.0.0.1": {Country: "Anonymous Proxy", Latitude: 0, Longitude: 0},
	}}
	g := observetest.Run(t, observe.TypeIP, "10.0.0.1", newGeo(observe.Options{}, locator))
	if _, ok := observetest.Prop(t, g, observe.TypeIP, "10.0.0.1", observe.FieldLatitude); ok {
		t.Error("recorded 0,0 as a location")
	}
	if v, ok := observetest.Prop(t, g, observe.TypeIP, "10.0.0.1", observe.FieldCountry); !ok || v.Str() != "Anonymous Proxy" {
		t.Errorf("country = %q (set=%v), want the country to survive", v.Str(), ok)
	}
}
