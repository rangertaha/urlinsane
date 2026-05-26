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
	"bytes"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/multiformats/go-multihash"
)

// cidPrefix defines content-addressing: CIDv1, dag-cbor codec, sha2-256.
var cidPrefix = cid.Prefix{
	Version:  1,
	Codec:    cid.DagCBOR,
	MhType:   multihash.SHA2_256,
	MhLength: -1,
}

// cidOf returns the deterministic CID for a dag-cbor block.
func cidOf(block []byte) (cid.Cid, error) {
	return cidPrefix.Sum(block)
}

// encodeEntity encodes an Entity to a canonical dag-cbor block.
func encodeEntity(e *Entity) ([]byte, error) {
	node := bindnode.Wrap(e, entityProto.Type())
	var buf bytes.Buffer
	if err := dagcbor.Encode(node.Representation(), &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// decodeEntity decodes a dag-cbor block back into an Entity.
func decodeEntity(block []byte) (*Entity, error) {
	builder := entityProto.Representation().NewBuilder()
	if err := dagcbor.Decode(builder, bytes.NewReader(block)); err != nil {
		return nil, err
	}
	return bindnode.Unwrap(builder.Build()).(*Entity), nil
}

// encodeScan encodes a Scan node to a dag-cbor block.
func encodeScan(s *Scan) ([]byte, error) {
	node := bindnode.Wrap(s, scanProto.Type())
	var buf bytes.Buffer
	if err := dagcbor.Encode(node.Representation(), &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// encodeEntityJSON encodes an Entity as dag-json (nested, for output plugins).
func encodeEntityJSON(e *Entity) ([]byte, error) {
	node := bindnode.Wrap(e, entityProto.Type())
	var buf bytes.Buffer
	if err := dagjson.Encode(node.Representation(), &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
