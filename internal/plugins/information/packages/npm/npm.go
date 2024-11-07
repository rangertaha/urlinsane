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
package npm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information/packages"
)

const (
	ORDER       = 1
	CODE        = "npm"
	DESCRIPTION = "Node packages repository"
)

type Plugin struct {
	conf internal.Config
}

func (p *Plugin) Id() string {
	return CODE
}

func (p *Plugin) Order() int {
	return ORDER
}

func (p *Plugin) Init(conf internal.Config) {
	p.conf = conf
}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (p *Plugin) Headers() []string {
	return []string{"NPM"}
}

func (p *Plugin) Exec(in internal.Typo) (out internal.Typo) {
	if p.Exists(in.Variant().Name()) {
		in.Variant().Add("NPM", "YES")
	}

	return in
}

func (p *Plugin) Exists(name string) (out bool) {
	name = strings.TrimSpace(name)
	url := fmt.Sprintf("https://www.npmjs.com/package/%s", name)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true
	}
	// fmt.Println(string(body))

	return false
}

// Register the plugin
func init() {
	packages.Add(CODE, func() internal.Information {
		return &Plugin{}
	})
}
