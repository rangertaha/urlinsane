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
package dns

// // mxLookupFunc
// func mxLookupFunc(tr Result) (results []Result) {
// 	records, _ := net.LookupMX(tr.Variant.String())
// 	tr.Variant.Meta.DNS.MX = dnsLib.NewMX(records...)
// 	for _, record := range records {
// 		record := strings.TrimSuffix(record.Host, ".")
// 		if !strings.Contains(tr.Data["MX"], record) {
// 			tr.Data["MX"] = strings.TrimSpace(tr.Data["MX"] + "\n" + record)
// 		}
// 	}
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

// // nsLookupFunc
// func nsLookupFunc(tr Result) (results []Result) {
// 	records, _ := net.LookupNS(tr.Variant.String())
// 	tr.Variant.Meta.DNS.NS = dnsLib.NewNS(records...)
// 	for _, record := range records {
// 		record := strings.TrimSuffix(record.Host, ".")
// 		if !strings.Contains(tr.Data["NS"], record) {
// 			tr.Data["NS"] = strings.TrimSpace(tr.Data["NS"] + "\n" + record)
// 		}
// 	}
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

// // cnameLookupFunc
// func cnameLookupFunc(tr Result) (results []Result) {
// 	records, _ := net.LookupCNAME(tr.Variant.String())
// 	// tr.Variant.Meta.DNS.CName = records
// 	for _, record := range records {
// 		tr.Data["CNAME"] = strings.TrimSuffix(string(record), ".")
// 		tr.Variant.Meta.DNS.CName = append(tr.Variant.Meta.DNS.CName, string(record))
// 	}
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

// // ipLookupFunc
// func ipLookupFunc(tr Result) (results []Result) {
// 	results = append(results, checkIP(tr))
// 	return
// }

// // txtLookupFunc
// func txtLookupFunc(tr Result) (results []Result) {
// 	records, _ := net.LookupTXT(tr.Variant.String())
// 	tr.Variant.Meta.DNS.TXT = records
// 	for _, record := range records {
// 		tr.Data["TXT"] = strings.TrimSpace(tr.Data["TXT"] + "\n" + record)
// 		tr.Variant.Meta.DNS.TXT = append(tr.Variant.Meta.DNS.TXT, record)
// 	}
// 	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
// 	return
// }

import (
	"github.com/rangertaha/urlinsane/internal"
	"github.com/rangertaha/urlinsane/internal/plugins/information"
)

const (
	CODE        = "dns"
	NAME        = "DNS Records"
	DESCRIPTION = "Retrieve DNS records"
)

type None struct {
	types []string
}

func (n *None) Id() string {
	return CODE
}

func (n *None) Name() string {
	return NAME
}

func (n *None) Description() string {
	return DESCRIPTION
}

func (n *None) Headers() []string {
	return []string{"DNS"}
}

func (n *None) Exec(in internal.Typo) (out internal.Typo) {
	in.Variant().Add("DNS", []string{"one", "two"})
	return in
}

// Register the plugin
func init() {
	information.Add(CODE, func() internal.Information {
		return &None{}
	})
}
