// // Copyright (C) 2024 Rangertaha
// //
// // This program is free software: you can redistribute it and/or modify
// // it under the terms of the GNU General Public License as published by
// // the Free Software Foundation, either version 3 of the License, or
// // (at your option) any later version.
// //
// // This program is distributed in the hope that it will be useful,
// // but WITHOUT ANY WARRANTY; without even the implied warranty of
// // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// // GNU General Public License for more details.
// //
// // You should have received a copy of the GNU General Public License
// // along with this program.  If not, see <http://www.gnu.org/licenses/>.
package target

import (
	"encoding/json"
	"strings"

	"github.com/rangertaha/urlinsane/internal/pkg/domain"
)

// Target ...
type Target struct {
	name  string
	meta  map[string]interface{}
	ready bool
	live  bool
}

func New(name string) *Target {
	name = strings.TrimSpace(name)
	return &Target{
		name:  name,
		live:  false,
		ready: false,
		meta:  make(map[string]interface{}),
	}
}

func (d *Target) Meta() map[string]interface{} {
	return d.meta
}

func (d *Target) Domain() (string, string, string) {
	dm := domain.Parse(d.name)
	return dm.Subdomain, dm.Prefix, dm.Suffix
}

func (d *Target) Add(key string, value interface{}) {
	d.meta[key] = value
}
func (d *Target) Get(key string) (value interface{}) {
	if value, ok := d.meta[key]; ok {
		return value
	}
	return nil
}

func (d *Target) Name() string {
	return d.name
}

func (d *Target) Live(v ...bool) bool {
	if len(v) > 0 {
		d.live = v[0]
	}
	return d.live
}
func (d *Target) Ready(v ...bool) bool {
	if len(v) > 0 {
		d.ready = v[0]
	}
	return d.ready
}

func (d *Target) Json() (j []byte, e error) {
	d.meta["name"] = d.name
	return json.Marshal(d.meta)
}
