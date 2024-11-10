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
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/databases"
)

const (
	CODE = "badger"
)

type Plugin struct {
	config internal.Config
}

func (n *Plugin) Init(conf internal.Config) {
	n.config = conf
}

func (n *Plugin) Read(key string) (i interface{}, err error) {
	return
}

func (n *Plugin) Write(key string, i interface{}) (err error) {
	return
}

// Register the plugin
func init() {
	databases.Add(CODE, func() internal.Database {
		return &Plugin{}
	})
}
