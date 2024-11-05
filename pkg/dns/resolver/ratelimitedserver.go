package resolver

import (
	"strings"
	"time"

	"github.com/miekg/dns"
	"go.uber.org/ratelimit"
)

type rateLimitedServer struct {
	limiter    ratelimit.Limiter
	ipAddrPort string

	pool pool
}

type pool interface {
	Get() (*dns.Client, *dns.Conn, error)
	Return(conn *dns.Conn)
	Close()
}

type server interface {
	Query(query string, rrtype RRtype) (*dns.Msg, time.Duration, error)
	Close()
}

func newRateLimitedServerList(ipAddrPort []string, queriesPerSecond int) []server {
	list := []server{}

	for i := range ipAddrPort {
		server, err := newRateLimitedServer(ipAddrPort[i], queriesPerSecond)

		if err != nil {
			continue
		}

		list = append(list, server)
	}

	return list
}

func newRateLimitedServer(IPAddrPort string, queriesPerSecond int) (*rateLimitedServer, error) {
	if !strings.Contains(IPAddrPort, ":") {
		IPAddrPort += ":53"
	}

	server := rateLimitedServer{}
	server.limiter = ratelimit.New(queriesPerSecond, ratelimit.WithoutSlack)
	server.ipAddrPort = IPAddrPort

	pool, err := newConnectionPool(1, queriesPerSecond, IPAddrPort)

	if err != nil {
		return nil, err
	}

	server.pool = pool

	return &server, nil
}

func (s *rateLimitedServer) Query(query string, rrtype RRtype) (*dns.Msg, time.Duration, error) {
	msg := new(dns.Msg)
	msg.Id = s.newID()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{Name: query + ".", Qtype: uint16(rrtype), Qclass: dns.ClassINET}

	s.limiter.Take()

	client, conn, err := s.pool.Get()

	if err != nil {
		return nil, 0, err
	}

	msg, dur, err := client.ExchangeWithConn(msg, conn)

	s.pool.Return(conn)

	return msg, dur, err
}

func (s *rateLimitedServer) Close() {
	s.pool.Close()
}

func (s *rateLimitedServer) newID() uint16 {
	return uint16(randID.Int31())
}

var randID *saferand

func init() {
	randID = newSafeRand(true)
}
