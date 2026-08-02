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
	"net"
	"testing"

	"github.com/rangertaha/urlinsane/internal/graph"
)

// fakeGeo answers from a table, so no MaxMind database is needed.
type fakeGeo struct {
	places map[string]*Location
	err    error
}

func (f fakeGeo) Locate(addr string) (*Location, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.places[addr], nil
}

// recordingGeo notes every address it was asked about, which is how a test
// proves geo never ran against a domain.
type recordingGeo struct {
	places map[string]*Location
	seen   *[]string
}

func (r recordingGeo) Locate(addr string) (*Location, error) {
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
		places: map[string]*Location{"93.184.216.34": {
			City: "Norwell", Country: "United States", CountryCode: "US",
			Continent: "North America", TimeZone: "America/New_York",
			PostalCode: "02061", Latitude: 42.1596, Longitude: -70.8217,
		}},
		seen: &seen,
	}
	res := &fakeResolver{ips: map[string][]net.IP{
		"example.com": {net.ParseIP("93.184.216.34")},
	}}

	g := run(t, TypeDomain, "example.com",
		newAddresses(Options{}, res), newGeo(Options{}, locator))

	if len(seen) != 1 || seen[0] != "93.184.216.34" {
		t.Fatalf("geo was asked about %v, want exactly the address", seen)
	}
	wantStatus(t, g, TypeIP, "93.184.216.34", "geo", graph.StatusOK)

	// It never matched the domain: the trigger binds to a type, and domain is
	// not it.
	if _, ok := g.Status(nodeID(t, g, TypeDomain, "example.com"), "geo"); ok {
		t.Error("geo ran against a domain node")
	}
}

// TestGeoEmitsPropsNotNodes guards DESIGN §2: a location has no identity worth
// converging on, so it is props on the address. A location node would be a new
// node per address with nothing ever joining on it.
func TestGeoEmitsPropsNotNodes(t *testing.T) {
	locator := fakeGeo{places: map[string]*Location{"93.184.216.34": {
		City: "Norwell", Country: "United States", CountryCode: "US",
		Continent: "North America", TimeZone: "America/New_York",
		PostalCode: "02061", Latitude: 42.1596, Longitude: -70.8217,
	}}}
	g := run(t, TypeIP, "93.184.216.34", newGeo(Options{}, locator))

	if n := len(g.Nodes()); n != 1 {
		t.Fatalf("geo admitted %d nodes beyond the address; a location is a prop: %s", n-1, dump(g))
	}
	if n := len(g.Edges()); n != 0 {
		t.Fatalf("geo emitted %d edges, want none", n)
	}
	for field, want := range map[string]string{
		FieldCity:        "Norwell",
		FieldCountry:     "United States",
		FieldCountryCode: "US",
		FieldContinent:   "North America",
		FieldTimeZone:    "America/New_York",
		FieldPostalCode:  "02061",
	} {
		if v, ok := prop(t, g, TypeIP, "93.184.216.34", field); !ok || v.Str() != want {
			t.Errorf("%s = %q (set=%v), want %q", field, v.Str(), ok, want)
		}
	}
	if v, ok := prop(t, g, TypeIP, "93.184.216.34", FieldLatitude); !ok || v.Real() != 42.1596 {
		t.Errorf("latitude = %v (set=%v), want 42.1596", v.Real(), ok)
	}

	// The operator declares props and no node type at all, which is why Effects
	// covers props: plan compilation would otherwise not see it produce
	// anything.
	eff := newGeo(Options{}, locator).Emits()
	if len(eff.Nodes) != 0 || len(eff.Rels) != 0 || len(eff.Props) == 0 {
		t.Errorf("geo effects = %+v, want props only", eff)
	}
}

// TestGeoUnknownAddressIsEmpty: the database answered and has nothing. That is
// absence, and it must not be confused with a database that failed to open.
func TestGeoUnknownAddressIsEmpty(t *testing.T) {
	g := run(t, TypeIP, "10.0.0.1", newGeo(Options{}, fakeGeo{}))
	wantStatus(t, g, TypeIP, "10.0.0.1", "geo", graph.StatusEmpty)
}

// TestGeoLookupErrorIsFailed keeps a broken database from reading as a world
// with no locations in it.
func TestGeoLookupErrorIsFailed(t *testing.T) {
	g := run(t, TypeIP, "10.0.0.1", newGeo(Options{}, fakeGeo{err: errors.New("corrupt database")}))
	wantStatus(t, g, TypeIP, "10.0.0.1", "geo", graph.StatusFailed)
}

// TestGeoIgnoresNullIsland: MaxMind uses 0,0 for "unknown", and recording it
// would place half the internet in the Gulf of Guinea.
func TestGeoIgnoresNullIsland(t *testing.T) {
	locator := fakeGeo{places: map[string]*Location{
		"10.0.0.1": {Country: "Anonymous Proxy", Latitude: 0, Longitude: 0},
	}}
	g := run(t, TypeIP, "10.0.0.1", newGeo(Options{}, locator))
	if _, ok := prop(t, g, TypeIP, "10.0.0.1", FieldLatitude); ok {
		t.Error("recorded 0,0 as a location")
	}
	if v, ok := prop(t, g, TypeIP, "10.0.0.1", FieldCountry); !ok || v.Str() != "Anonymous Proxy" {
		t.Errorf("country = %q (set=%v), want the country to survive", v.Str(), ok)
	}
}
