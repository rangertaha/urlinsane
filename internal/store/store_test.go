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

package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	gdb, err := gorm.Open(sqlite.Open(filepath.Join(dir, "index.db")), &gorm.Config{})
	if err != nil {
		t.Fatalf("open index db: %v", err)
	}
	s, err := Open(dir, gdb)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	return s
}

func sampleDomain() *db.Domain {
	return &db.Domain{
		Type: entity.Domain,
		Name: "exmaple.com",
		Rank: 7,
		Dns: []*db.Dns{
			{Type: "A", Value: "93.184.216.34", Ttl: "300"},
			{Type: "NS", Value: "a.iana-servers.net"},
		},
		IPs: []*db.Address{
			{Addr: "93.184.216.34", Type: "IPv4", Location: &db.Location{Code: "US", Name: "United States"}},
		},
		Whois: []db.Whois{
			{Registrant: &db.Contact{Name: "ACME", Country: "US"}},
		},
	}
}

func TestConvertRoundTrip(t *testing.T) {
	d := sampleDomain()
	got := ToDomain(ToEntity(d))
	if got.EntityType() != entity.Domain || got.Name != d.Name || got.Rank != d.Rank {
		t.Fatalf("scalar mismatch: %+v", got)
	}
	if len(got.Dns) != 2 || len(got.IPs) != 1 || len(got.Whois) != 1 {
		t.Fatalf("slice length mismatch: dns=%d ips=%d whois=%d", len(got.Dns), len(got.IPs), len(got.Whois))
	}
	if got.IPs[0].Location == nil || got.IPs[0].Location.Code != "US" {
		t.Fatalf("location lost: %+v", got.IPs[0])
	}
	if got.Whois[0].Registrant == nil || got.Whois[0].Registrant.Name != "ACME" {
		t.Fatalf("whois contact lost: %+v", got.Whois[0])
	}
}

func TestPutGetRoundTrip(t *testing.T) {
	s := testStore(t)
	if _, err := s.Put(sampleDomain()); err != nil {
		t.Fatalf("put: %v", err)
	}
	got, ok, err := s.Get(entity.Domain, "exmaple.com")
	if err != nil || !ok {
		t.Fatalf("get: ok=%v err=%v", ok, err)
	}
	if got.Name != "exmaple.com" || len(got.Dns) != 2 {
		t.Fatalf("loaded domain mismatch: %+v", got)
	}

	if _, ok, _ := s.Get(entity.Domain, "absent.com"); ok {
		t.Fatal("expected absent lookup to report not found")
	}
}

func TestGetFresh(t *testing.T) {
	s := testStore(t)
	if _, err := s.Put(sampleDomain()); err != nil { // exmaple.com
		t.Fatal(err)
	}
	// Fresh within a generous ttl.
	if _, ok, err := s.GetFresh(entity.Domain, "exmaple.com", time.Hour); err != nil || !ok {
		t.Fatalf("expected fresh cache hit: ok=%v err=%v", ok, err)
	}
	// ttl <= 0 disables caching.
	if _, ok, _ := s.GetFresh(entity.Domain, "exmaple.com", 0); ok {
		t.Fatal("ttl=0 should disable caching")
	}
	// Stale: a ttl shorter than the time already elapsed since Put.
	time.Sleep(5 * time.Millisecond)
	if _, ok, _ := s.GetFresh(entity.Domain, "exmaple.com", time.Millisecond); ok {
		t.Fatal("expected stale miss")
	}
}

func TestPutNormalizationDeterministic(t *testing.T) {
	s := testStore(t)
	// Same logical domain, DNS records appended in a different order.
	d1 := &db.Domain{Type: entity.Domain, Name: "x.com", Dns: []*db.Dns{
		{Type: "A", Value: "1.1.1.1"}, {Type: "NS", Value: "n2"}, {Type: "NS", Value: "n1"},
	}}
	d2 := &db.Domain{Type: entity.Domain, Name: "x.com", Dns: []*db.Dns{
		{Type: "NS", Value: "n1"}, {Type: "A", Value: "1.1.1.1"}, {Type: "NS", Value: "n2"},
	}}
	c1, err := s.Put(d1)
	if err != nil {
		t.Fatal(err)
	}
	c2, err := s.Put(d2)
	if err != nil {
		t.Fatal(err)
	}
	if !c1.Equals(c2) {
		t.Fatalf("normalization failed: %s != %s", c1, c2)
	}
}

func TestPutIdempotent(t *testing.T) {
	s := testStore(t)
	c1, _ := s.Put(sampleDomain())
	c2, _ := s.Put(sampleDomain())
	if !c1.Equals(c2) {
		t.Fatalf("idempotent put gave different CIDs: %s != %s", c1, c2)
	}
	if ok, err := s.bs.Has(c1); err != nil || !ok {
		t.Fatalf("block missing after put: ok=%v err=%v", ok, err)
	}
}

func TestScanDiff(t *testing.T) {
	s := testStore(t)
	// Scan 1: a (with A record) and b.
	a1 := &db.Domain{Type: entity.Domain, Name: "a.com", Dns: []*db.Dns{{Type: "A", Value: "1.1.1.1"}}}
	b := &db.Domain{Type: entity.Domain, Name: "b.com"}
	if _, err := s.PutScan("q", []*db.Domain{a1, b}); err != nil {
		t.Fatal(err)
	}
	// Scan 2: a changed (different IP), b unchanged, c added.
	a2 := &db.Domain{Type: entity.Domain, Name: "a.com", Dns: []*db.Dns{{Type: "A", Value: "2.2.2.2"}}}
	b2 := &db.Domain{Type: entity.Domain, Name: "b.com"}
	c := &db.Domain{Type: entity.Domain, Name: "c.com"}
	if _, err := s.PutScan("q", []*db.Domain{a2, b2, c}); err != nil {
		t.Fatal(err)
	}

	diff, err := s.Diff("q")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(diff.Changed, "a.com") {
		t.Fatalf("expected a.com changed, got %+v", diff)
	}
	if !contains(diff.Same, "b.com") {
		t.Fatalf("expected b.com same, got %+v", diff)
	}
	if !contains(diff.Added, "c.com") {
		t.Fatalf("expected c.com added, got %+v", diff)
	}
}

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
