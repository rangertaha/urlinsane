package resolver

import (
	"sync/atomic"
)

type serverBalancer struct {
	servers  []server
	curIndex int64
	count    int
}

func newServerBalancer(servers []server) *serverBalancer {
	count := len(servers)

	list := serverBalancer{
		servers:  make([]server, count),
		curIndex: -1,
		count:    count,
	}

	copy(list.servers, servers)

	return &list
}

func (s *serverBalancer) Next() server {
	if len(s.servers) == 0 {
		return nil
	}

	nextIndex := int(atomic.AddInt64(&s.curIndex, 1))
	server := s.servers[nextIndex%s.count]

	return server
}

func (s *serverBalancer) Close() {
	for _, server := range s.servers {
		server.Close()
	}
}
