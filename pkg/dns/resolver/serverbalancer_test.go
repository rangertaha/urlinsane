package resolver

import (
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type stubServer struct {
	ipAddrPort string
	closes     int
}

func (s *stubServer) Query(query string, rrtype RRtype) (*dns.Msg, time.Duration, error) {
	return &dns.Msg{}, time.Duration(0), nil
}

func (s *stubServer) Close() {
	s.closes++
}

func TestServerBalancerNext(t *testing.T) {
	stubServerA := &stubServer{ipAddrPort: "8.8.8.8:53"}
	stubServerB := &stubServer{ipAddrPort: "8.8.4.4:53"}

	testTable := []struct {
		name        string
		haveServers []server
		haveCount   int
		want        server
	}{
		{name: "No Items", haveServers: []server{}, haveCount: 1, want: nil},
		{name: "Nil Resolver", haveServers: nil, haveCount: 1, want: nil},
		{name: "Single", haveServers: []server{stubServerA}, haveCount: 1, want: stubServerA},
		{name: "Second", haveServers: []server{stubServerA, stubServerB}, haveCount: 2, want: stubServerB},
		{name: "Wrap Around", haveServers: []server{stubServerA, stubServerB}, haveCount: 3, want: stubServerA},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			list := newServerBalancer(test.haveServers)

			var got server
			for i := 0; i < test.haveCount; i++ {
				got = list.Next()
			}

			assert.Equal(t, test.want, got)
		})
	}
}

func TestServerBalancerClose(t *testing.T) {
	tests := []struct {
		name string
		have []string
		want int
	}{
		{name: "Empty", have: []string{}, want: 0},
		{name: "Single", have: []string{"8.8.8.8:53"}, want: 1},
		{name: "Double", have: []string{"8.8.8.8:53", "8.8.4.4:53"}, want: 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serverList := []server{}

			for _, ipAddrPort := range test.have {
				serverList = append(serverList, &stubServer{ipAddrPort: ipAddrPort})
			}

			list := newServerBalancer(serverList)
			list.Close()

			got := 0
			for _, server := range serverList {
				got += server.(*stubServer).closes
			}

			assert.Equal(t, test.want, got)
		})
	}
}
