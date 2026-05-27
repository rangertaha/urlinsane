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
package list

import (
	"fmt"
	"os"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
)

type Plugin struct {
	outputs.Plugin
}

func (p *Plugin) Init(conf internal.Config) {
	p.Started = time.Now()
	p.Config = conf
}

func (p *Plugin) Read(in *db.Domain) {
	p.Domains = append(p.Domains, in)

	p.Total++
	if in.Live() {
		p.Online++
	} else {
		p.Offline++
	}
}

// Write renders all collected variants as one aligned table once the scan
// completes (filtering by --registered/--unregistered is applied in Table).
func (p *Plugin) Write() {
	fmt.Println(p.Table())
}

func (p *Plugin) Save(fname string) {
	if err := os.WriteFile(fname, []byte(p.Table()+"\n"), 0644); err != nil {
		fmt.Printf("Saving file: %s \n", err)
	}
}

// Register the plugin
func init() {
	var CODE = "list"
	outputs.Add(CODE, func() internal.Output {
		return &Plugin{
			Plugin: outputs.Plugin{
				ID:      CODE,
				Summary: "outputs one record per line",
			},
		}
	})
}
