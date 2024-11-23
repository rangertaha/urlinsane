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
package none

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/analyzers"
)

const (
	CODE        = "a"
	ORDER       = 1
	DESCRIPTION = "Default "
)

type Plugin struct{}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Order() int {
	return ORDER
}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (n *Plugin) Init(conf internal.Config) {

}

func (n *Plugin) Headers() []string {
	return []string{"NONE"}
}

func (n *Plugin) Exec(original, variant internal.Domain, acc internal.Accumulator) (err error) {
	return
}

// Register the plugin
func init() {
	analyzers.Add(CODE, func() internal.Analyzer {
		return &Plugin{}
	})
}
