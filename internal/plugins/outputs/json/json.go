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
package txt

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/models"
	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
)

const (
	CODE        = "json"
	DESCRIPTION = "Deeply nested JSON structured output"
)

type Record struct {
	Original string        `json:"domain"`
	Variant  models.Domain `json:"variant"`
	Distance int           `json:"distance"`
}

type Json struct {
	config internal.Config
	output string
	typos  []internal.Typo
}

func (p *Record) Json() string {
	// Marshal the struct into JSON
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(jsonData)
}

func (n *Json) Id() string {
	return CODE
}

func (n *Json) Description() string {
	return DESCRIPTION
}

func (n *Json) Init(conf internal.Config) {
	n.config = conf
}

func (n *Json) Write(in internal.Typo) {
	if n.config.Progress() {
		n.typos = append(n.typos, in)
	} else {
		n.Stream(in)
	}
}

func (n *Json) Stream(in internal.Typo) {
	orig, vari := in.Get()

	record := Record{
		Original: orig.Fqdn(),
		Variant:  vari,
		Distance: in.Dist(),
	}
	fmt.Println(record.Json())
}

func (n *Json) Filter(header string) bool {
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

func (n *Json) Summary(report map[string]int64) {
	//
}

func (n *Json) Save() {
	if n.config.Progress() {
		for _, typo := range n.typos {
			n.Stream(typo)
		}
	}

	if n.config.File() != "" {
		results := []byte(n.output)
		if err := os.WriteFile(n.config.File(), results, 0644); err != nil {
			fmt.Printf("Error: %s", err)
		}
	}
}

// Register the plugin
func init() {
	outputs.Add(CODE, func() internal.Output {
		return &Json{}
	})
}
