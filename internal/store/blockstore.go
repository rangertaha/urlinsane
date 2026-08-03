// Copyright 2024 Rangertaha. All rights reserved.
// SPDX-License-Identifier: GPL-3.0-or-later

package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ipfs/go-cid"
)

// Blockstore is a content-addressed store of raw dag-cbor blocks. Because the
// key is derived from the content, Put is idempotent and a Get can be verified
// against the key it was fetched with.
type Blockstore interface {
	Has(c cid.Cid) (bool, error)
	Get(c cid.Cid) ([]byte, error)
	Put(c cid.Cid, block []byte) error
}

// FSBlockstore is a filesystem blockstore. Blocks live at
// <root>/<shard>/<cid>.cbor, sharded on the last two characters of the CID so
// that a wide scan's 10^5 nodes do not land in one directory.
type FSBlockstore struct {
	root string
}

// NewFSBlockstore opens (creating if needed) a blockstore under dir/blocks.
func NewFSBlockstore(dir string) (*FSBlockstore, error) {
	root := filepath.Join(dir, "blocks")
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &FSBlockstore{root: root}, nil
}

func (s *FSBlockstore) shardDir(c cid.Cid) string {
	name := c.String()
	shard := "_"
	if len(name) >= 2 {
		shard = name[len(name)-2:]
	}
	return filepath.Join(s.root, shard)
}

func (s *FSBlockstore) path(c cid.Cid) string {
	return filepath.Join(s.shardDir(c), c.String()+".cbor")
}

// Has reports whether the block is present.
func (s *FSBlockstore) Has(c cid.Cid) (bool, error) {
	_, err := os.Stat(s.path(c))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Get reads a block.
func (s *FSBlockstore) Get(c cid.Cid) ([]byte, error) {
	b, err := os.ReadFile(s.path(c))
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("store: block %s not found", c)
	}
	return b, err
}

// Put writes a block, skipping the write when it is already stored. Identical
// content always lands at the same path, so a re-scan that changes nothing
// writes nothing.
//
// The write is atomic: a temporary in the shard directory, then a rename. That
// is what makes the skip above safe to trust. os.WriteFile creates and
// truncates before it writes, so a crash or a full disk in the middle left a
// short file at the block's path — and because the path exists, Has says the
// block is stored and every later Put skips it. In a content-addressed store
// that is the one corruption that cannot be noticed: the bytes at path(c) no
// longer hash to c, nothing re-checks them, and nothing will ever overwrite
// them. A rename cannot half-happen, so a block is either absent or whole.
func (s *FSBlockstore) Put(c cid.Cid, block []byte) error {
	if ok, _ := s.Has(c); ok {
		return nil
	}
	if err := os.MkdirAll(s.shardDir(c), 0o755); err != nil {
		return err
	}
	return writeAtomic(s.path(c), block, 0o644)
}

// MemBlockstore is an in-memory blockstore, for tests and for a scan that is
// diffed against an earlier one but never itself persisted.
type MemBlockstore struct {
	mu     sync.RWMutex
	blocks map[string][]byte
}

// NewMemBlockstore returns an empty in-memory blockstore.
func NewMemBlockstore() *MemBlockstore {
	return &MemBlockstore{blocks: map[string][]byte{}}
}

// Has reports whether the block is present.
func (m *MemBlockstore) Has(c cid.Cid) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.blocks[c.KeyString()]
	return ok, nil
}

// Get reads a block.
func (m *MemBlockstore) Get(c cid.Cid) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.blocks[c.KeyString()]
	if !ok {
		return nil, fmt.Errorf("store: block %s not found", c)
	}
	return b, nil
}

// Put stores a block.
func (m *MemBlockstore) Put(c cid.Cid, block []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.blocks[c.KeyString()] = append([]byte(nil), block...)
	return nil
}
