package resolver

import (
	"sync/atomic"
	"time"
)

type question struct {
	question string
	rrtype   RRtype
	channel  chan []Record
}

type clientDNS struct {
	resolver       resolver
	maxConcurrency int

	queryChan  chan question
	queryCount int32
}

type resolver interface {
	Resolve(query string, rrtype RRtype) []Record
	Close()
}

func newClientDNS(resolver resolver, maxConcurrency int) *clientDNS {
	client := clientDNS{
		resolver:       resolver,
		maxConcurrency: maxConcurrency,
	}

	client.startThreads()

	return &client
}

func (s *clientDNS) startThreads() {
	s.queryChan = make(chan question, s.maxConcurrency)

	for i := 0; i < s.maxConcurrency; i++ {
		go func(queryChan chan question) {
			for {
				query, open := <-queryChan

				if !open {
					return
				}

				results := s.resolver.Resolve(query.question, query.rrtype)
				query.channel <- results
			}
		}(s.queryChan)
	}
}

func (s *clientDNS) Resolve(queries []string, rrtype RRtype) []Record {
	records := []Record{}

	queryCount := len(queries)

	if queryCount == 0 {
		return records
	}

	// Start result reader
	var received int32
	resultChan := make(chan []Record, s.maxConcurrency)
	go func(resultChan chan []Record) {
		for {
			response, open := <-resultChan

			if !open {
				return
			}

			records = append(records, response...)
			atomic.AddInt32(&received, 1)
		}
	}(resultChan)

	// Send work to goroutines
	for _, query := range queries {
		question := question{
			question: query,
			rrtype:   rrtype,
			channel:  resultChan,
		}

		s.queryChan <- question
	}

	// Wait for results
	for int(atomic.LoadInt32(&received)) < queryCount {
		time.Sleep(1 * time.Millisecond)
	}

	// Work done
	close(resultChan)
	atomic.AddInt32(&s.queryCount, int32(queryCount))

	return records
}

func (s *clientDNS) QueryCount() int {
	return int(atomic.LoadInt32(&s.queryCount))
}

func (s *clientDNS) Close() {
	close(s.queryChan)
	s.resolver.Close()
}
