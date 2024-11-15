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
package badger

import (
	"path/filepath"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/databases"
	log "github.com/sirupsen/logrus"
)

const (
	CODE = "badger"
)

type Plugin struct {
	config internal.Config
	db     *KV
	log    *log.Entry
}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Init(conf internal.Config) {
	n.log = log.WithFields(log.Fields{"plugin": CODE})
	n.config = conf
	var err error

	dbpath := filepath.Join(conf.Dir(), "db")
	ilog := n.log.WithFields(log.Fields{"path": dbpath, "ttl": conf.TTL()})
	if n.db, err = NewBadgerDb(dbpath, conf.TTL()); err != nil {
		ilog.Error(err)
	}

}

func (n *Plugin) Read(keys ...string) (value string, err error) {
	key := strings.Join(keys, ":")
	if value, err = n.db.Get(key); err != nil {
		n.log.WithError(err).Errorf("Read(%s) -> len(%d)", key, len(value))
		return value, err
	}
	n.log.Debugf("Read(%s) -> len(%d)", key, len(value))
	return
}

func (n *Plugin) Write(value string, keys ...string) (err error) {
	key := strings.Join(keys, ":")
	if err := n.db.Set(key, value); err != nil {
		n.log.WithError(err).Errorf("Write(%s) <- len(%d)", key, len(value))
		return err
	}
	n.log.Debugf("Write(%s) <- len(%d)", key, len(value))
	return
}

func (n *Plugin) Close() {
	n.db.Close()
	n.log.Debug("Connection closed")
}

// Register the plugin
func init() {
	databases.Add(CODE, func() internal.Database {
		return &Plugin{}
	})
}
