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
package txt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
	log "github.com/sirupsen/logrus"
)

const (
	CODE        = "json"
	DESCRIPTION = "Deeply nested JSON structured output"
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

func (n *Plugin) Read(domain internal.Domain) {
	n.domains = append(n.domains, domain)

	if !n.config.Progress() {
		fmt.Println(domain.Json())
	}
}

func (n *Plugin) Write() {
	if n.config.Progress() {
		for _, domain := range n.domains {
			fmt.Println(domain.Json())
		}
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

func (n *Plugin) Summary(report map[string]string) {}

func (n *Plugin) Save(fname string) {
	// Open the file for writing
	file, err := os.Create(fname)
	if err != nil {
		log.Error("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffered writer for efficiency
	writer := bufio.NewWriter(file)

	// Stream data to the file
	for _, domain := range n.domains {
		_, err := writer.WriteString(domain.Json())
		if err != nil {
			log.Error("Error writing to file:", err)
			return
		}
	}

	// Flush the buffered data to the file
	writer.Flush()

	log.Info("Data streamed to file successfully")
}

// Register the plugin
func init() {
	outputs.Add(CODE, func() internal.Output {
		return &Plugin{}
	})
}
