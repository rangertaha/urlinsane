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
package bn

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/db"
	"github.com/rangertaha/urlinsane/internal/plugins/collectors"
	// log "github.com/sirupsen/logrus"
)

type Plugin struct {
	collectors.Plugin
}

func (p *Plugin) Exec(domain *db.Domain) (vaiant *db.Domain, err error) {
	for _, ip := range domain.IPs {
		p.Banner(ip)
	}

	return domain, err
}

func (p *Plugin) Banner(addr *db.Address) (err error) {

	// address := fmt.Sprintf("%s:%d", addr.Addr, 80)
	// conn, err := net.DialTimeout("http", address, 10*time.Second)
	// if err != nil {
	// 	log.Error("Error:", err.Error())
	// 	// conn.Close()
	// 	return
	// }

	// // Send the request to the server
	// // fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: %s\r\n\r\n", host)

	// // Read the response from the server
	// buffer := make([]byte, 1024)
	// n, err := conn.Read(buffer)
	// if err != nil {
	// 	log.Error("Error:", err.Error())
	// 	return
	// }
	// var service db.Service
	// svc := db.Service{
	// 	Name:   "HTTP",
	// 	Banner: string(buffer[:n]),
	// }
	// db.DB.FirstOrInit(service, svc)
	// // var ip db.Address
	// // db.DB.FirstOrInit(&ip, db.Address{Addr: addr.Addr})
	// addr.Ports = append(addr.Ports, &db.Port{Number: 80, Service: &service})

	// defer conn.Close()
	return err
}

func (p *Plugin) Close() {}

// Register the plugin
func init() {
	var CODE = "bn"
	collectors.Add(CODE, func() internal.Collector {
		return &Plugin{
			Plugin: collectors.Plugin{
				Num:       5,
				Code:      CODE,
				Title:     "Banner Grabber",
				Summary:   "Capturing HTTP/SMTP banners",
				DependsOn: []string{"ip"},
			},
		}
	})
}
