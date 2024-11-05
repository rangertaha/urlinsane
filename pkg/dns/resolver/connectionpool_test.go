package resolver

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerPoolNew(t *testing.T) {
	tests := []struct {
		name           string
		haveInitCount  int
		haveMaxCount   int
		haveIPAddrPort string
		wantCount      int
		wantErr        bool
	}{
		{name: "Single", haveInitCount: 1, haveMaxCount: 1, haveIPAddrPort: "8.8.8.8:53", wantCount: 1, wantErr: false},
		{name: "Couple", haveInitCount: 2, haveMaxCount: 2, haveIPAddrPort: "8.8.8.8:53", wantCount: 2, wantErr: false},
		{name: "Invalid", haveInitCount: 1, haveMaxCount: 1, haveIPAddrPort: "invalid", wantCount: 0, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool, gotErr := newConnectionPool(test.haveInitCount, test.haveMaxCount, test.haveIPAddrPort)

			if gotErr == nil {
				gotCount := pool.Count()
				assert.Equal(t, test.wantCount, gotCount)
			}

			assert.Equal(t, test.wantErr, gotErr != nil)
		})
	}
}

func TestServerPoolGet(t *testing.T) {
	tests := []struct {
		name           string
		haveInitCount  int
		haveMaxCount   int
		haveIPAddrPort string
		wantCount      int
		wantErr        bool
	}{
		{name: "On Demand", haveInitCount: 0, haveMaxCount: 1, haveIPAddrPort: "8.8.8.8:53", wantCount: 1, wantErr: false},
		{name: "Max Conn", haveInitCount: 1, haveMaxCount: 1, haveIPAddrPort: "8.8.8.8:53", wantCount: 1, wantErr: false},
		{name: "No Conn", haveInitCount: 0, haveMaxCount: 0, haveIPAddrPort: "8.8.8.8:53", wantCount: 0, wantErr: true},
		{name: "Failed", haveInitCount: 0, haveMaxCount: 1, haveIPAddrPort: "invalid", wantCount: 0, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pool, err := newConnectionPool(test.haveInitCount, test.haveMaxCount, test.haveIPAddrPort)

			if err != nil {
				t.Fatal("pool creation failed")
			}

			_, _, gotErr := pool.Get()

			assert.Equal(t, test.wantCount, pool.Count())
			assert.Equal(t, test.wantErr, gotErr != nil)
		})
	}
}

func TestServerPoolReturn(t *testing.T) {
	pool, err := newConnectionPool(0, 1, "8.8.8.8:53")

	if err != nil {
		t.Fatal("pool creation failed")
	}

	// First request
	_, connA, err := pool.Get()
	assert.Nil(t, err)

	// Blocking request
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_, connB, err := pool.Get()
		assert.Nil(t, err)
		assert.Equal(t, connA, connB)

		wg.Done()
	}()

	// Unblock by returning connection
	time.Sleep(100 * time.Millisecond)
	pool.Return(connA)

	wg.Wait()
}

func TestServerPoolGetFail(t *testing.T) {
	pool, err := newConnectionPool(0, 2, "8.8.8.8:53")

	if err != nil {
		t.Fatal("pool creation failed")
	}

	// First request
	_, connA, err := pool.Get()
	assert.Nil(t, err)

	// Ensure the second request fails to create a connection
	pool.IPAddrPort = "invalid"

	// Start the second request
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_, connB, _ := pool.Get()
		assert.Equal(t, connA, connB)
		wg.Done()
	}()

	// Return connA to unblock goroutine
	time.Sleep(100 * time.Millisecond)
	pool.Return(connA)

	wg.Wait()
}

func TestServerPoolClose(t *testing.T) {
	pool, err := newConnectionPool(2, 2, "8.8.8.8:53")

	if err != nil {
		t.Fatal("pool creation failed")
	}

	pool.Close()

	_, open := <-pool.channel

	assert.False(t, open)
}
