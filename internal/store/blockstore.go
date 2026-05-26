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

package store

import (
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
)

// Blockstore is a content-addressed store of raw IPLD blocks.
type Blockstore interface {
	Has(c cid.Cid) (bool, error)
	Get(c cid.Cid) ([]byte, error)
	Put(c cid.Cid, block []byte) error
}

// fsStore is a filesystem blockstore. Blocks live under <root>/<shard>/<cid>.cbor
// where the shard is the last two characters of the CID string, keeping any one
// directory from accumulating tens of thousands of entries.
type fsStore struct {
	root string
}

func newFSStore(dir string) (*fsStore, error) {
	root := filepath.Join(dir, "blocks")
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	return &fsStore{root: root}, nil
}

func (s *fsStore) shardDir(c cid.Cid) string {
	name := c.String()
	shard := "_"
	if len(name) >= 2 {
		shard = name[len(name)-2:]
	}
	return filepath.Join(s.root, shard)
}

func (s *fsStore) path(c cid.Cid) string {
	return filepath.Join(s.shardDir(c), c.String()+".cbor")
}

func (s *fsStore) Has(c cid.Cid) (bool, error) {
	_, err := os.Stat(s.path(c))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *fsStore) Get(c cid.Cid) ([]byte, error) {
	return os.ReadFile(s.path(c))
}

// Put writes the block. Because the CID is derived from the content, Put is
// idempotent: storing identical content rewrites the same path.
func (s *fsStore) Put(c cid.Cid, block []byte) error {
	if ok, _ := s.Has(c); ok {
		return nil
	}
	if err := os.MkdirAll(s.shardDir(c), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.path(c), block, 0o644)
}
