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
package engine

import (
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

type Database struct {
	opts badger.Options
	db   *badger.DB
	ttl  int
}

func (d *Database) Init() {
	var err error
	d.db, err = badger.Open(d.opts)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Database) Read() {

}

func (d *Database) Write() {

}

func (d *Database) Close() {
	d.db.Close()
}

// func main() {
// 	// Open the Badger database located in the /tmp/badger directory.
// 	// It will be created if it doesn't exist.

// 	// Your code here…
// }
