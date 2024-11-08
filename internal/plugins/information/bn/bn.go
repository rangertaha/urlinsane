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
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/models"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
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
	_, vari := in.Get()
	// fmt.Println(vari.FQDN)
	if vari.Live {

		_, err := p.Response("tcp", vari.IP(), 80)
		if err != nil {
			fmt.Println("Error dialing:", err)
			return in
		}
		// fmt.Println(vari.FQDN, vari.IP())
		// in.SetMeta("BANNER", res.Status)
		// vari.Response = res
		// in.Set(orig, vari)
	}

	return in
}

func (p *Plugin) Response(proto, host string, port int) (out models.Banner, err error) {
	domain := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout(proto, domain, 1*time.Second)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		return
	}
	fmt.Println(resp)

	return
	//	return models.Banner{
	//		Status:        resp.Status,
	//		Protocol:      resp.Proto,
	//		ContentLength: resp.ContentLength,
	//		StatusCode:    resp.StatusCode,
	//		// Headers:       resp.Header.,
	//		// Cookies:       resp.Cookies(),
	//		// Trailer:       resp.Trailer,
	//		// TLS:           resp.TLS,
	//	}, nil
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &Plugin{}
	})
}
