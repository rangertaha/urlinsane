// Copyright (C) 2024 Rangertaha
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
package text

import (
	"log"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
)

type Store struct {
	config internal.Config
	db     *badger.DB
}

func (s *Store) Init(conf internal.Config) {
	s.config = conf
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	var err error
	s.db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer s.db.Close()
}

// Read(key string) (error, interface{})
// Write(key string, value interface{}) error
func (n *Store) Read(key string) (err error, i interface{}) {
	return
}

func (n *Store) Write(key string, value interface{}) (err error) {
	return
}

// Register the plugin
func init() {
	outputs.Add("bdb", func() internal.Output {
		return &Store{}
	})
}
