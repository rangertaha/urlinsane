// Copyright 2024 Rangertaha. All Rights Reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package graph

import (
	"crypto/sha256"
	"encoding/binary"
	"sync"
)

// CacheKey identifies one operator result.
type CacheKey [32]byte

// Cache stores operator results. The key covers everything the operator reads,
// so a hit is only served when nothing it depends on has changed.
//
// A hit still occupies its round: the delta is applied at the same point a live
// call's would be, so a warm run and a cold run produce the same rounds, the
// same barriers and the same graph. A cache that short-circuited the round
// structure would make plan pinning depend on cache state.
type Cache struct {
	mu      sync.Mutex
	entries map[CacheKey]cacheEntry

	// modelCIDs and resourceConfig are mixed into every key. A retrained plugin
	// model or a changed --nameservers must invalidate; the operator's own
	// Version() only covers its code.
	modelCIDs      map[string][]string
	resourceConfig map[string]string
}

type cacheEntry struct {
	delta   Delta
	outcome Outcome
}

func NewCache() *Cache {
	return &Cache{
		entries:        map[CacheKey]cacheEntry{},
		modelCIDs:      map[string][]string{},
		resourceConfig: map[string]string{},
	}
}

// SetModels records the model CIDs an operator reads.
func (c *Cache) SetModels(operator string, cids []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.modelCIDs[operator] = append([]string(nil), cids...)
}

// SetResourceConfig records the config digest for a resource class — the
// nameservers for dns, the registry URL and token for npm, the proxy and UA
// policy for http.
func (c *Cache) SetResourceConfig(resource, digest string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.resourceConfig[resource] = digest
}

// Key builds the cache key for one dispatch.
func (c *Cache) Key(op Operator, id NodeID, readDigest [32]byte) CacheKey {
	c.mu.Lock()
	models := c.modelCIDs[op.Id()]
	rcfg := c.resourceConfig[op.Resource()]
	c.mu.Unlock()

	h := sha256.New()
	writeField(h, op.Id())
	var v [8]byte
	binary.BigEndian.PutUint64(v[:], uint64(op.Version()))
	_, _ = h.Write(v[:])
	for _, m := range models {
		writeField(h, m)
	}
	_, _ = h.Write(id[:])
	_, _ = h.Write(readDigest[:])
	writeField(h, rcfg)

	var k CacheKey
	copy(k[:], h.Sum(nil))
	return k
}

// Get returns a cached result.
func (c *Cache) Get(k CacheKey) (Delta, Outcome, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[k]
	return e.delta, e.outcome, ok
}

// Put stores a result. Transient failures are never cached — the scheduler
// filters those before calling.
func (c *Cache) Put(k CacheKey, d Delta, o Outcome) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[k] = cacheEntry{delta: d, outcome: o}
}
