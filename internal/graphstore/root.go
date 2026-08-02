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

package graphstore

import (
	"fmt"

	"github.com/ipfs/go-cid"
)

// Root is a scan root: the single block that names one stored graph.
//
// It carries no timestamp and no run id. That is the whole point — a scan root
// exists so "what changed since last week" is a CID comparison, and a clock
// reading inside it would make every re-scan differ (docs/DESIGN.md §1.2). If a
// caller needs to know when a scan ran, that belongs beside the root in an
// index, not inside it.
type Root struct {
	Version  int
	SeedType string
	SeedKey  string // canonical
	Nodes    []cid.Cid
	Edges    []cid.Cid
	Side     cid.Cid
}

// Scan is a loaded scan root with its blocks decoded. Nodes and edges keep the
// root's order, which is graph.Nodes()' (type, key) and graph.Edges()'
// (from, rel, to) — the canonical report order of §11.
type Scan struct {
	CID   cid.Cid
	Root  *Root
	Nodes []*NodeBlock
	Edges []*EdgeBlock
	Side  *Side
}

// Node finds a loaded node block by canonical identity.
func (s *Scan) Node(typeName, key string) (*NodeBlock, bool) {
	for _, n := range s.Nodes {
		if n.Type == typeName && n.Key == key {
			return n, true
		}
	}
	return nil, false
}

// encodeRoot writes the root block:
//
//	[version, seedType, seedKey, [nodeLinks], [edgeLinks], sideLink]
//
// Node and edge CIDs are inline link lists rather than a linked index block:
// the diff wants them and nothing else, so an extra hop would cost a block read
// on the one path that matters most.
func encodeRoot(r *Root) ([]byte, cid.Cid, error) {
	return encodeList(6, func(e *enc) {
		e.i64(int64(r.Version))
		e.str(r.SeedType)
		e.str(r.SeedKey)
		e.sub(len(r.Nodes), func(e *enc) {
			for _, c := range r.Nodes {
				e.link(c)
			}
		})
		e.sub(len(r.Edges), func(e *enc) {
			for _, c := range r.Edges {
				e.link(c)
			}
		})
		e.link(r.Side)
	})
}

// decodeRoot parses a root block.
func decodeRoot(block []byte) (*Root, error) {
	n, err := decodeBlock(block)
	if err != nil {
		return nil, err
	}
	d := newDec(n).expect("scan root", 6)
	r := &Root{
		Version:  int(d.at(0).i64()),
		SeedType: d.at(1).str(),
		SeedKey:  d.at(2).str(),
	}
	if d.err() == nil && r.Version != FormatVersion {
		return nil, fmt.Errorf("graphstore: scan root format version %d, this build writes %d", r.Version, FormatVersion)
	}
	d.at(3).each(func(v *dec) { r.Nodes = append(r.Nodes, v.link()) })
	d.at(4).each(func(v *dec) { r.Edges = append(r.Edges, v.link()) })
	r.Side = d.at(5).link()
	if err := d.err(); err != nil {
		return nil, err
	}
	return r, nil
}
