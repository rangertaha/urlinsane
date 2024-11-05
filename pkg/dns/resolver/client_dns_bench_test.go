package resolver

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkClientDNS(b *testing.B) {
	for _, queries := range []int{1, 16, 256, 4096, 65536} {
		for _, routines := range []int{1, 4, 16, 64, 256, 1024} {
			runner(b, queries, routines)
		}
	}
}

func runner(b *testing.B, queries int, routines int) bool {
	return b.Run(fmt.Sprintf("queries:%d-routines:%d", queries, routines), func(b *testing.B) {
		resolver := &stubResolver{sleep: time.Duration(100 * time.Nanosecond), records: []Record{{}}}

		client := newClientDNS(resolver, routines)
		client.Resolve(make([]string, queries), TypeA)
	})
}
