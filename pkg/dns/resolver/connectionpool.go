package resolver

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"
)

type connectionPool struct {
	sync.Mutex

	client     *dns.Client
	IPAddrPort string

	channel chan *dns.Conn

	maxCount  int32
	initCount int32
	count     int32
}

func newConnectionPool(initCount int, maxCount int, IPAddrPort string) (*connectionPool, error) {
	pool := &connectionPool{
		initCount:  int32(initCount),
		maxCount:   int32(maxCount),
		IPAddrPort: IPAddrPort,
	}

	pool.client = new(dns.Client)
	pool.client.Dialer = &net.Dialer{Timeout: 10 * time.Second}

	pool.channel = make(chan *dns.Conn, maxCount)

	for i := 0; i < initCount; i++ {
		conn, err := pool.createConn()

		if err != nil {
			return nil, err
		}

		pool.channel <- conn
	}

	return pool, nil
}

func (s *connectionPool) Count() int {
	return int(s.count)
}

func (s *connectionPool) Get() (*dns.Client, *dns.Conn, error) {
	for {
		// Can't create new connections, return an existing one
		if atomic.LoadInt32(&s.count) == s.maxCount {
			if s.count == 0 {
				return nil, nil, errors.New("connection pool empty and unable to create new connections")
			}

			return s.client, <-s.channel, nil
		}

		// Select an available connection or try to create a new one
		select {
		case conn := <-s.channel:
			return s.client, conn, nil

		default:
			s.Lock()
			if s.count < s.maxCount {
				conn, err := s.createConn()

				// Return the connection that was created
				if err == nil {
					s.Unlock()
					return s.client, conn, nil
				}

				// Unable to create a new connection and no connection in the pool, return error
				if s.count == 0 {
					s.Unlock()
					return nil, nil, err
				}

				// Unable to create a new connection, stop trying
				s.maxCount = s.count
			}
			s.Unlock()
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (s *connectionPool) Return(conn *dns.Conn) {
	s.channel <- conn
}

func (s *connectionPool) Close() {
	close(s.channel)

	for conn := range s.channel {
		conn.Close()
	}

	s.count = 0
}

func (s *connectionPool) createConn() (*dns.Conn, error) {
	c, err := s.client.Dial(s.IPAddrPort)

	if err != nil {
		return nil, err
	}

	atomic.AddInt32(&s.count, 1)

	return c, nil
}
