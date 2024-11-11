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
package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
	log "github.com/sirupsen/logrus"
)

const (
	CODE        = "list"
	DESCRIPTION = "outputs one record per line"
)

type Plugin struct {
	config  internal.Config
	domains []internal.Domain
}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (n *Plugin) Init(conf internal.Config) {
	n.config = conf
}

func (n *Plugin) Read(in internal.Domain) {
	n.domains = append(n.domains, in)
	if !n.config.Progress() {
		fmt.Println(n.Row(in))
	}
}

func (n *Plugin) Write() {
	if n.config.Progress() {
		fmt.Println(strings.Join(n.Rows(n.domains...), "\n"))
	}
}

func (n *Plugin) Filter(header string) bool {
	header = strings.TrimSpace(header)
	header = strings.ToLower(header)
	for _, filter := range n.config.Filters() {
		filter = strings.TrimSpace(filter)
		filter = strings.ToLower(filter)
		if filter == header {
			return true
		}
	}
	return false
}

func (n *Plugin) Summary(report map[string]string) {
	fmt.Println("")
	for k, v := range report {
		fmt.Printf("%s %s   ", k, v)
	}
	fmt.Println("")
}

func (n *Plugin) Save(fname string) {
	output := strings.Join(n.Rows(n.domains...), "\n")
	if err := os.WriteFile(fname, []byte(output), 0644); err != nil {
		log.Errorf("Saving file: %s", err)
	}
}

// Register the plugin
func init() {
	outputs.Add(CODE, func() internal.Output {
		return &Plugin{}
	})
}
