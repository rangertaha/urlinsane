package resolver

import (
	"github.com/miekg/dns"
)

type resolverDNS struct {
	serverBalancer balancer
	parser         messageParser

	retryCodes map[int]struct{}
	retryCount int
}

type balancer interface {
	Next() server
	Close()
}

type messageParser interface {
	Parse(msg *dns.Msg) []Record
}

func newResolverDNS(retryCount int, retryCodes []int, serverBalancer balancer, parser messageParser) *resolverDNS {
	retryMap := make(map[int]struct{}, len(retryCodes))
	for _, code := range retryCodes {
		retryMap[code] = struct{}{}
	}

	return &resolverDNS{
		serverBalancer: serverBalancer,
		parser:         parser,
		retryCodes:     retryMap,
		retryCount:     retryCount,
	}
}

func (s *resolverDNS) Resolve(query string, rrtype RRtype) []Record {
	var err error
	var msg *dns.Msg

	// Send the request to a server, retrying on error
	var i int
	for i = 0; i < s.retryCount; i++ {
		server := s.serverBalancer.Next()
		msg, _, err = server.Query(query, rrtype)

		if err != nil {
			continue
		}

		if _, ok := s.retryCodes[msg.Rcode]; ok {
			continue
		}

		break
	}

	// Max retry count reached
	if i == s.retryCount {
		return nil
	}

	// Parse records only on success or NXDOMAIN
	if msg.Rcode != dns.RcodeSuccess && msg.Rcode != dns.RcodeNameError {
		return nil
	}

	return s.parser.Parse(msg)
}

func (s *resolverDNS) Close() {
	s.serverBalancer.Close()
}
