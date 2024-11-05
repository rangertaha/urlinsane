package resolver

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type spyBalancer struct {
	calls  int32
	closes int32
	server server
}

func (s *spyBalancer) Next() server {
	s.calls++

	return s.server
}

func (s *spyBalancer) Close() {
	s.closes++
}

type fakeServer struct {
	errors int
	msg    *dns.Msg
}

func (s *fakeServer) Query(query string, ttype RRtype) (*dns.Msg, time.Duration, error) {
	if s.errors > 0 {
		s.errors--
		return nil, time.Duration(0), errors.New("error")
	}

	return s.msg, time.Duration(0), nil
}

func (s *fakeServer) Close() {
}

func createStubResolver(retryCount int, retryCodes []int, errorCount int, msg *dns.Msg) (*resolverDNS, *spyBalancer) {
	stubServer := &fakeServer{
		errors: errorCount,
		msg:    msg,
	}
	spyBalancer := &spyBalancer{server: stubServer}
	resolver := newResolverDNS(retryCount, retryCodes, spyBalancer, &msgParser{})

	return resolver, spyBalancer
}

var stubARecord = Record{"www.example.com", TypeA, "127.0.0.1"}
var stubCNAMERecord = Record{"www.example.com", TypeCNAME, "cname.example.com"}

func createStubMessage(code int, rrtype RRtype) *dns.Msg {
	stubMsg := &dns.Msg{}
	stubMsg.Question = []dns.Question{{Name: "www.example.com"}}
	stubMsg.Rcode = code

	switch rrtype {
	case TypeA:
		stubA := &dns.A{}
		stubA.Hdr.Rrtype = uint16(stubARecord.Type)
		stubA.A = net.ParseIP(stubARecord.Answer)
		stubMsg.Answer = []dns.RR{stubA}
	case TypeCNAME:
		stubCNAME := &dns.CNAME{}
		stubCNAME.Hdr.Rrtype = uint16(stubCNAMERecord.Type)
		stubCNAME.Target = stubCNAMERecord.Answer
		stubMsg.Answer = []dns.RR{stubCNAME}
	}

	return stubMsg
}

func TestResolve_Simple(t *testing.T) {
	msg := createStubMessage(dns.RcodeSuccess, TypeA)
	resolver, spyBalancer := createStubResolver(3, nil, 0, msg)

	got := resolver.Resolve("www.example.com", TypeA)

	assert.EqualValues(t, []Record{stubARecord}, got)
	assert.EqualValues(t, 1, spyBalancer.calls)
}

func TestResolve_Retry(t *testing.T) {
	msg := createStubMessage(dns.RcodeSuccess, TypeA)
	resolver, spyBalancer := createStubResolver(3, nil, 2, msg)

	got := resolver.Resolve("www.example.com", TypeA)

	assert.EqualValues(t, []Record{stubARecord}, got)
	assert.EqualValues(t, 3, spyBalancer.calls)
}

func TestResolve_MaxRetry(t *testing.T) {
	msg := createStubMessage(dns.RcodeSuccess, TypeA)
	resolver, spyBalancer := createStubResolver(1, nil, 1, msg)

	got := resolver.Resolve("www.example.com", TypeA)

	assert.EqualValues(t, []Record(nil), got)
	assert.EqualValues(t, 1, spyBalancer.calls)
}

func TestResolve_NXDOMAINWithCNAME(t *testing.T) {
	msg := createStubMessage(dns.RcodeNameError, TypeCNAME)
	resolver, spyBalancer := createStubResolver(1, nil, 0, msg)

	got := resolver.Resolve("www.example.com", TypeCNAME)

	assert.EqualValues(t, []Record{stubCNAMERecord}, got)
	assert.EqualValues(t, 1, spyBalancer.calls)
}

func TestResolve_SERVFAILNoRetry(t *testing.T) {
	msg := createStubMessage(dns.RcodeServerFailure, TypeA)
	resolver, spyBalancer := createStubResolver(3, nil, 0, msg)

	got := resolver.Resolve("www.example.com", TypeCNAME)

	assert.EqualValues(t, []Record(nil), got)
	assert.EqualValues(t, 1, spyBalancer.calls)
}

func TestResolve_SERVFAILWithRetry(t *testing.T) {
	msg := createStubMessage(dns.RcodeServerFailure, TypeA)
	resolver, spyBalancer := createStubResolver(3, []int{dns.RcodeServerFailure}, 0, msg)

	got := resolver.Resolve("www.example.com", TypeCNAME)

	assert.EqualValues(t, []Record(nil), got)
	assert.EqualValues(t, 3, spyBalancer.calls)
}

func TestClose(t *testing.T) {
	msg := createStubMessage(dns.RcodeServerFailure, TypeA)
	resolver, spyBalancer := createStubResolver(3, nil, 0, msg)

	resolver.Close()
	assert.EqualValues(t, 1, spyBalancer.closes)
}
