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
package cn

import (
	"net"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/pkg"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
	log "github.com/sirupsen/logrus"
)

const (
	ORDER       = 3
	CODE        = "cn"
	DESCRIPTION = "DNS CNAME records"
)

type Plugin struct {
	// resolver resolver.Client
	conf internal.Config
	log  *log.Entry
}

func (n *Plugin) Id() string {
	return CODE
}

func (n *Plugin) Order() int {
	return ORDER
}

func (i *Plugin) Init(c internal.Config) {
	i.log = log.WithFields(log.Fields{"plugin": CODE, "method": "Exec"})
	i.conf = c
}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (n *Plugin) Headers() []string {
	return []string{"CNAME"}
}

func (i *Plugin) Exec(acc internal.Accumulator) (err error) {
	l := i.log.WithFields(log.Fields{"domain": acc.Domain().String()})
	// if acc.Domain().Cached() {
	// 	// l.Debug("Returning cache domain: ", acc.Domain().String())
	// 	return acc.Next()
	// }

	dns := make(pkg.DnsRecords, 0)
	if err := acc.Unmarshal("DNS", &dns); err != nil {
		l.Error("Unmarshal DNS: ", err)
	}

	cname, err := net.LookupCNAME(acc.Domain().String())
	if err != nil {
		l.Error(err)
	}
	if cname != "" {
		dns.Add("CNAME", 0, cname)

		acc.SetMeta("CNAME", cname)
		acc.SetJson("DNS", dns.Json())
		acc.Domain().Live(true)
	}

	return acc.Next()
}

func (i *Plugin) Close() {}

// Register the plugin
func init() {
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{}
	})
}
