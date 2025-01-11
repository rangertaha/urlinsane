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
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
)

type Plugin struct {
	outputs.Plugin

	config  internal.Config
	domains []*db.Domain
	elapsed time.Duration
	started time.Time
	offline int64
	online  int64
	total   int64
}

func (p *Plugin) Init(conf internal.Config) {
	p.started = time.Now()
	p.config = conf
}

func (p *Plugin) Read(in *db.Domain) {
	p.domains = append(p.domains, in)

	p.total++
	if in.Live() {
		p.online++
	} else {
		p.offline++
	}

	if p.config.Registered() {
		if !in.Live() {
			return
		}
	}

	if p.config.Unregistered() {
		if in.Live() {
			return
		}
	}

	if !p.config.Progress() {
		fmt.Println(p.Row(in))
	}

}

func (p *Plugin) Write() {
	if p.config.Progress() {
		fmt.Println(strings.Join(p.Rows(p.domains...), "\n"))
	}
}

func (p *Plugin) Summary() {
	p.elapsed = time.Since(p.started)
	summary := map[string]string{
		"  TIME:":  p.elapsed.String(),
		"  TOTAL:": fmt.Sprintf("%d", p.total),
	}
	if len(p.config.Collectors()) > 0 {
		summary[text.FgGreen.Sprintf("%s", "  LIVE:")] = fmt.Sprintf("%d", p.online)
		summary[text.FgRed.Sprintf("%s", "  OFFLINE")] = fmt.Sprintf("%d", p.offline)
	}

	fmt.Println("")
	for k, v := range summary {
		fmt.Printf("%s %s   ", k, v)
	}
	fmt.Println("")
}

func (p *Plugin) Save(fname string) {
	output := strings.Join(p.Rows(p.domains...), "\n")
	if err := os.WriteFile(fname, []byte(output), 0644); err != nil {
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
