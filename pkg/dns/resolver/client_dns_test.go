package resolver

import (
	"fmt"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type stubResolver struct {
	records []Record
	index   int32
	sleep   time.Duration

	closes int
}

func (s *stubResolver) Resolve(query string, rrtype RRtype) []Record {
	nextIndex := int(atomic.AddInt32(&s.index, 1) - 1)
	record := s.records[nextIndex%len(s.records)]

	time.Sleep(s.sleep)

	return []Record{record}
}

func (s *stubResolver) Close() {
	s.closes++
}

func TestClientDNSResolve(t *testing.T) {
	testTable := []struct {
		name           string
		haveConcurrent int
		haveDomains    []string
		haveRRtype     RRtype
		want           []Record
	}{
		{
			name:           "Empty",
			haveConcurrent: 5,
			haveDomains:    []string{},
			haveRRtype:     TypeA,
			want:           []Record{},
		},
		{
			name:           "Simple",
			haveConcurrent: 5,
			haveDomains:    []string{"foo.bar"},
			haveRRtype:     TypeA,
			want: []Record{
				{
					Question: "foo.bar",
					Type:     TypeA,
					Answer:   "127.0.0.1",
				},
			},
		},
		{
			name:           "Concurrency",
			haveConcurrent: 2,
			haveDomains:    []string{"foo.bar", "abc.xyz"},
			haveRRtype:     TypeA,
			want: []Record{
				{
					Question: "foo.bar",
					Type:     TypeA,
					Answer:   "127.0.0.1",
				},
				{
					Question: "abc.xyz",
					Type:     TypeA,
					Answer:   "127.0.1.1",
				},
			},
		},
		{
			name:           "Max Concurrency",
			haveConcurrent: 1,
			haveDomains:    []string{"foo.bar", "abc.xyz", "wine.bar"},
			haveRRtype:     TypeA,
			want: []Record{
				{
					Question: "foo.bar",
					Type:     TypeA,
					Answer:   "127.0.0.1",
				},
				{
					Question: "abc.xyz",
					Type:     TypeA,
					Answer:   "127.0.1.1",
				},
				{
					Question: "wine.bar",
					Type:     TypeA,
					Answer:   "127.1.1.1",
				},
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			resolver := &stubResolver{sleep: time.Duration(10 * time.Millisecond), records: test.want}

			client := newClientDNS(resolver, test.haveConcurrent)

			got := client.Resolve(test.haveDomains, test.haveRRtype)

			sort.SliceStable(test.want, func(i, j int) bool {
				return test.want[i].Question < test.want[j].Question
			})

			sort.SliceStable(got, func(i, j int) bool {
				return got[i].Question < got[j].Question
			})

			assert.Equal(t, test.want, got)
		})
	}
}

func TestClientDNSResolveLarge(t *testing.T) {
	const iterations int = 32768

	resolver := &stubResolver{sleep: time.Duration(0), records: []Record{{Question: "foo.bar", Type: TypeA, Answer: "127.0.0.1"}}}

	client := newClientDNS(resolver, 10)

	list := make([]string, iterations)
	for i := range list {
		list[i] = fmt.Sprintf("query-%d", i)
	}

	got := client.Resolve(list, TypeA)

	assert.Equal(t, iterations, len(got))
}

func TestClientDNSQueryCount(t *testing.T) {
	tests := []struct {
		name        string
		haveDomains []string
		haveThreads int
		want        int
	}{
		{name: "No Query", haveDomains: []string{}, haveThreads: 1, want: 0},
		{name: "Single Query", haveDomains: []string{"test.com"}, haveThreads: 1, want: 1},
		{name: "Multiple Queries Single Thread", haveDomains: []string{"test.com", "foo.com", "bar.com"}, haveThreads: 1, want: 3},
		{name: "Multiple Queries Multi Thread", haveDomains: []string{"test.com", "foo.com", "bar.com"}, haveThreads: 10, want: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolver := &stubResolver{sleep: time.Duration(0), records: []Record{{Question: "foo.bar", Type: TypeA, Answer: "127.0.0.1"}}}

			client := newClientDNS(resolver, 10)
			client.Resolve(test.haveDomains, TypeA)

			got := client.QueryCount()

			assert.Equal(t, test.want, got)
		})
	}
}

func TestClientDNSClose(t *testing.T) {
	resolver := &stubResolver{}

	client := newClientDNS(resolver, 10)
	client.Close()

	assert.Equal(t, 1, resolver.closes)
}
