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
package badger

import (
	"errors"
	"fmt"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

type KV struct {
	db  *badger.DB
	ttl time.Duration
}

func NewBadgerDb(pathToDb string, ttl time.Duration) (*KV, error) {
	opts := badger.DefaultOptions(pathToDb)

	opts.Logger = nil
	badgerInstance, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("opening kv: %w", err)
	}

	return &KV{db: badgerInstance, ttl: ttl}, nil
}

func (k *KV) Close() error {
	return k.db.Close()
}

// nolint:wrapcheck
func (k *KV) Exists(key string) (bool, error) {
	var exists bool
	err := k.db.View(
		func(tx *badger.Txn) error {
			if val, err := tx.Get([]byte(key)); err != nil {
				return err
			} else if val != nil {
				exists = true
			}
			return nil
		})
	if errors.Is(err, badger.ErrKeyNotFound) {
		err = nil
	}
	return exists, err
}

func (k *KV) Get(key string) (string, error) {
	var value string
	return value, k.db.View(
		func(tx *badger.Txn) error {
			item, err := tx.Get([]byte(key))
			if err != nil {
				return fmt.Errorf("getting value: %w", err)
			}
			valCopy, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("copying value: %w", err)
			}
			value = string(valCopy)
			return nil
		})
}

//	  err := db.Update(func(txn *badger.Txn) error {
//		e := badger.NewEntry([]byte("answer"), []byte("42")).WithTTL(time.Hour)
//		err := txn.SetEntry(e)
//		return err
//	  })
func (k *KV) Set(key, value string) error {
	return k.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(value)).WithTTL(k.ttl)
		err := txn.SetEntry(e)
		return err
	})
}

// func (k *KV) Set(key, value string) error {
// 	return k.db.Update(
// 		func(txn *badger.Txn) error {
// 			return txn.Set([]byte(key), []byte(value))
// 		})
// }

func (k *KV) Delete(key string) error {
	return k.db.Update(
		func(txn *badger.Txn) error {
			return txn.Delete([]byte(key))
		})
}
