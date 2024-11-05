package resolver

import (
	"errors"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type spyLimiter struct {
	calls int
}

func (s *spyLimiter) Take() time.Time {
	s.calls++

	return time.Time{}
}

type fakePool struct {
	connectionPool
}

func (s *fakePool) Get() (*dns.Client, *dns.Conn, error) {
	return nil, nil, errors.New("error")
}

func TestRateLimitedServerNewList(t *testing.T) {
	tests := []struct {
		name string
		have []string
		want int
	}{
		{name: "Valid Resolver", have: []string{"8.8.8.8:53"}, want: 1},
		{name: "Empty Resolver", have: []string{}, want: 0},
		{name: "Multiple Resolvers", have: []string{"8.8.8.8:53", "invalid", "127.0.0.1"}, want: 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotList := newRateLimitedServerList(test.have, 10)

			assert.Equal(t, test.want, len(gotList))
		})
	}
}

func TestRateLimitedServerQuery(t *testing.T) {
	limiter := &spyLimiter{}

	server, err := newRateLimitedServer("8.8.8.8:53", 10)

	if err != nil {
		t.Fatal("unable to connect to resolver")
	}

	server.limiter = limiter

	_, _, gotErr := server.Query("www.google.com", TypeA)

	assert.Equal(t, 1, limiter.calls)
	assert.Nil(t, gotErr)
}

func TestRateLimitedServerQueryPoolErr(t *testing.T) {
	server, err := newRateLimitedServer("8.8.8.8:53", 10)

	if err != nil {
		t.Fatal("unable to connect to resolver")
	}

	// Make sure the next Query receives an error from the connection pool
	server.pool = &fakePool{}

	_, _, gotErr := server.Query("www.google.com", TypeA)

	assert.NotNil(t, gotErr)
}

func TestRateLimitedServerClose(t *testing.T) {
	server, err := newRateLimitedServer("8.8.8.8:53", 10)

	if err != nil {
		t.Fatal("unable to connect to resolver")
	}

	server.Close()
}
