// The MIT License (MIT)
//
// # Copyright © 2019 CYBINT
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package dns

import "net"

type (
	NS struct {
		net.NS
		Host string
	}
	MX struct {
		net.MX
		Host string `json:"host,omitempty"`
		Pref uint16 `json:"pref,omitempty"`
	}
)

// NewMX ...
func NewMX(mx ...*net.MX) (nmx []MX) {
	for _, m := range mx {
		nmx = append(nmx, MX{Host: m.Host, Pref: m.Pref})
	}
	return
}

// NewNS ...
func NewNS(ns ...*net.NS) (nns []NS) {
	for _, n := range ns {
		nns = append(nns, NS{Host: n.Host})
	}
	return
}

// func checkIP(tr Result) Result {
// 	if tr.Variant.Meta.DNS.ipCheck == false {
// 		records, _ := net.LookupIP(tr.Variant.String())
// 		// if err != nil {
// 		// 	fmt.Println(err)
// 		// }
// 		for _, record := range uniqIP(records) {
// 			dotlen := strings.Count(record, ".")
// 			if dotlen == 3 {
// 				if !strings.Contains(tr.Data["IPv4"], record) {
// 					tr.Data["IPv4"] = strings.TrimSpace(tr.Data["IPv4"] + "\n" + record)
// 					tr.Variant.Meta.DNS.IPv4 = append(tr.Variant.Meta.DNS.IPv4, record)
// 				}
// 				tr.Variant.Live = true
// 			}
// 			clen := strings.Count(record, ":")
// 			if clen == 5 {
// 				if !strings.Contains(tr.Data["IPv6"], record) {
// 					tr.Data["IPv6"] = strings.TrimSpace(tr.Data["IPv6"] + "\n" + record)
// 					tr.Variant.Meta.DNS.IPv6 = append(tr.Variant.Meta.DNS.IPv6, record)
// 				}
// 				tr.Variant.Live = true
// 			}
// 		}
// 		tr.Variant.Meta.DNS.ipCheck = true
// 	}

// 	return Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data}
// }

// func uniqIP(list []net.IP) (ulist []string) {
// 	uinq := map[string]bool{}
// 	for _, l := range list {
// 		uinq[l.String()] = true
// 	}
// 	for k := range uinq {
// 		ulist = append(ulist, k)
// 	}
// 	return
// }
