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
package bn

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	log "github.com/sirupsen/logrus"
)

type Banner struct {
	Port  string `json:"port,omitempty"`
	String string `json:"string,omitempty"`
}

type Banners []Banner

type Plugin struct {
	// resolver resolver.Client
	conf internal.Config
}

func (i *Plugin) Init(c internal.Config) {
	i.conf = c
}

func (i *Plugin) Exec(acc internal.Accumulator) (err error) {
	ports := []string{"80", "21", "587"}
	banners := Banners{}

	for _, port := range ports {
		b := Banner{
			Port:   port,
			String: i.Banner("tcp", acc.Domain().String(), port),
		}
		banners = append(banners, b)
	}
	if len(banners) > 0 {
		acc.SetMeta("BANNER", banners.String("80"))
		acc.SetJson("BANNER", banners.Json())
		acc.Domain().Live(true)
		acc.Save("banner.txt", []byte(banners.String("80")))
	}

	return acc.Next()
}

func (i *Plugin) Close() {}

func (i *Plugin) Banner(proto, host, port string) (bn string) {
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.DialTimeout(proto, address, 10*time.Second)
	if err != nil {
		log.Error("Error:", err.Error())
		return
	}
	defer conn.Close()

	// Send the request to the server
	fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: %s\r\n\r\n", host)

	// Read the response from the server
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Error("Error:", err.Error())
		return
	}

	return string(buffer[:n])
}

func (b *Banners) Json() json.RawMessage {
	records, err := json.Marshal(b)
	if err != nil {
		log.Error(err)
	}
	return json.RawMessage(records)
}

func (b *Banners) String(p string) (values string) {
	for _, banner := range *b {
		if banner.Port == p {
			return banner.String
		}
	}
	return
}
