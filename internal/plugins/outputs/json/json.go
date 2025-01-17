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
package json

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
	log "github.com/sirupsen/logrus"
)

type Plugin struct {
	outputs.Plugin
}

func (p *Plugin) Read(domain *db.Domain) {
	p.Domains = append(p.Domains, domain)

	if !p.Config.Progress() {
		fmt.Println(domain.Json())
	}
}

func (p *Plugin) Write() {
	if p.Config.Progress() {
		for _, domain := range p.Domains {
			fmt.Println(domain.Json())
		}
	}
}

// func (p *Plugin) Filter(header string) bool {
// 	header = strings.TrimSpace(header)
// 	header = strings.ToLower(header)
// 	for _, filter := range p.Config.Filters() {
// 		filter = strings.TrimSpace(filter)
// 		filter = strings.ToLower(filter)
// 		if filter == header {
// 			return true
// 		}
// 	}
// 	return false
// }

func (p *Plugin) Summary(report map[string]string) {}

func (p *Plugin) Save(fname string) {
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
	for _, domain := range p.Domains {
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
	var CODE = "json"
	outputs.Add(CODE, func() internal.Output {
		return &Plugin{
			Plugin: outputs.Plugin{
				ID:      CODE,
				Summary: "nested JSON structured output",
			},
		}
	})
}
