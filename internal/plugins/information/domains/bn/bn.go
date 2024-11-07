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
	"net"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information/domains"
)

const (
	ORDER       = 8
	CODE        = "bn"
	DESCRIPTION = "Banner grabbing "
)

type Plugin struct {
	conf internal.Config
}

func (p *Plugin) Id() string {
	return CODE
}

func (p *Plugin) Order() int {
	return ORDER
}

func (p *Plugin) Init(conf internal.Config) {
	p.conf = conf
}

func (n *Plugin) Description() string {
	return DESCRIPTION
}

func (p *Plugin) Headers() []string {
	return []string{"BANNER"}
}

func (p *Plugin) Exec(in internal.Typo) (out internal.Typo) {
	// if v := in.Variant(); v.Live() {
	banner := p.Banner(in.Variant().Name())
	in.Variant().Add("BANNER", banner)
	// }
	return in
}

func (p *Plugin) Banner(domain string) (out string) {

	// host := os.Args[1]
	// port := os.Args[2]

	// Connect to the target host and port
	conn, err := net.DialTimeout("tcp", domain+":80", 5*time.Second)
	if err != nil {
		// fmt.Println("Error:", err.Error())
		return
	}
	defer conn.Close()

	// Send the request to the server
	// fmt.Fprintf(conn, "GET / HTTP/1.1\r\nHost: %s\r\n\r\n", domain)

	// Read the response from the server
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)

	// Print the response
	response := string(buffer[:n])
	// fmt.Println(response)
	return response
}

// Register the plugin
func init() {
	domains.Add(CODE, func() internal.Information {
		return &Plugin{}
	})
}
