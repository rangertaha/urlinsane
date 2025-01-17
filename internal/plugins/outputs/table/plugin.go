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
package table

// import (
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/jedib0t/go-pretty/text"
// 	"github.com/jedib0t/go-pretty/v6/table"
// 	"github.com/rangertaha/urlinsane/internal"
// 	"github.com/rangertaha/urlinsane/internal/db"
// 	"github.com/rangertaha/urlinsane/internal/pkg"
// 	"github.com/rangertaha/urlinsane/internal/plugins/outputs"
// 	log "github.com/sirupsen/logrus"
// 	"golang.org/x/term"
// )

// type Plugin struct {
// 	outputs.Plugin
// 	table  table.Writer
// 	output string
// }

// func (p *Plugin) Init(conf internal.Config) {
// 	p.Plugin.Init(conf)

// 	p.table = table.NewWriter()

// 	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
// 		p.table.SetAllowedRowLength(width - 4)
// 	}

// 	p.table.SetOutputMirror(os.Stdout)
// 	p.table.AppendHeader(p.Header())
// 	p.table.AppendFooter(p.Header())
// 	p.Conf()
// }

// func (p *Plugin) Read(in *db.Domain) {
// 	p.table.AppendRow(p.Row(in))
// }

// func (p *Plugin) Header() (row table.Row) {
// 	row = append(row, "LD")
// 	row = append(row, "TYPE")
// 	row = append(row, "TYPO")

// 	for _, collector := range p.Config.Collectors() {
// 		if collector.Id() == "ip" {

// 		}
// 	}

// 	return
// }

// func (p *Plugin) Row(domain *db.Domain) (row table.Row) {
// 	p.Domains = append(p.Domains, domain)
// 	row = append(row, domain.Levenshtein)
// 	if p.Config.Verbose() {
// 		row = append(row, domain.Algorithm.Name)
// 	} else {
// 		row = append(row, strings.ToUpper(domain.Algorithm.Code))
// 	}
// 	row = append(row, domain.Name)

// 	// for _, info := range p.Config.Collectors() {
// 	// 	// for _, header := range info.Headers() {
// 	// 	// 	meta := domain.Meta()
// 	// 	// 	if col, ok := meta[header]; ok {
// 	// 	// 		row = append(row, col)
// 	// 	// 	} else {
// 	// 	// 		row = append(row, "")
// 	// 	// 	}
// 	// 	// }
// 	// }
// 	return
// }

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

// func (p *Plugin) Conf() (row table.Row) {
// 	p.table.SetStyle(pkg.StyleDefault)

// 	// nameTransformer := text.Transformer(func(val interface{}) string {
// 	// 	if val.(string) == "MD" {
// 	// 		return text.Colors{text.BgBlack, text.FgGreen}.Sprint(val)
// 	// 	}
// 	// 	return fmt.Sprint(val)
// 	// })

// 	// n.table.SetRowPainter()

// 	p.table.SetColumnConfigs(ColumnConfig)
// 	return
// }
// func (p *Plugin) Progress(typo <-chan internal.Domain) <-chan internal.Domain {
// 	return typo
// }

// func (p *Plugin) Write() {
// 	p.output = p.table.Render()
// }

// func (p *Plugin) Save(fname string) {
// 	results := []byte(p.output)
// 	if err := os.WriteFile(fname, results, 0644); err != nil {
// 		log.Errorf("Error: %s", err)
// 	}
// }

// // Register the plugin
// func init() {
// 	var CODE = "table"
// 	outputs.Add(CODE, func() internal.Output {
// 		return &Plugin{
// 			Plugin: outputs.Plugin{
// 				ID:      CODE,
// 				Summary: "pretty table output format with color",
// 			},
// 		}
// 	})
// }
