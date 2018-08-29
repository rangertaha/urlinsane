// Copyright © 2018 rangertaha rangertaha@gmail.com
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

package urlinsane

import (
	"net"
	"strings"
)

// The registry for extra functions
var FREGISTRY = make(map[string][]Extra)

var mxLookup = Extra{
	Code:        "MX",
	Name:        "MX Lookup",
	Description: "Checking for DNS's MX records",
	Exec:        mxLookupFunc,
	Headers:     []string{"MX"},
}

var txtLookup = Extra{
	Code:        "TXT",
	Name:        "TXT Lookup",
	Description: "Checking for DNS's TXT records",
	Exec:        txtLookupFunc,
	Headers:     []string{"TXT"},
}

var ipLookup = Extra{
	Code:        "IP",
	Name:        "IP Lookup",
	Description: "Checking for IP address",
	Exec:        ipLookupFunc,
	Headers:     []string{"IPv4", "IPv6"},
}

var nsLookup = Extra{
	Code:        "NS",
	Name:        "NS Lookup",
	Description: "Checks DNS NS records",
	Exec:        nsLookupFunc,
	Headers:     []string{"NS"},
}

var cnameLookup = Extra{
	Code:        "CNAME",
	Name:        "CNAME Lookup",
	Description: "Checks DNS CNAME records",
	Exec:        cnameLookupFunc,
	Headers:     []string{"CNAME"},
}

var geoIPLookup = Extra{
	Code:        "GEO",
	Name:        "GeoIP Lookup",
	Description: "Looks up geopgraphic information via IP address",
	Exec:        geoIPLookupFunc,
	Headers:     []string{"GEO"},
}

var idnaLookup = Extra{
	Code:        "IDNA",
	Name:        "IDNA Domain",
	Description: "Show international domain name",
	Exec:        idnaFunc,
	Headers:     []string{"IDNA"},
}

//var allFuncGroup = ExtraGroup{
//	Code:        "ALL",
//	Name:        "All Functions",
//	Description: "Apply all post typosquating functions",
//	Funcs:       []Extra{
//		mxLookup,
//		ipLookup,
//		idnaLookup,
//		txtLookup,
//		nsLookup,
//		cnameLookup,
//		//geoIPLookup,
//	},
//}

func init() {
	FRegister("IDNA", idnaLookup)
	FRegister("MX", mxLookup)
	FRegister("IP", ipLookup)
	FRegister("TXT", txtLookup)
	FRegister("NS", nsLookup)
	FRegister("CNAME", cnameLookup)

	//FRegister("GEO", geoIPLookup)

	FRegister("ALL",
		mxLookup,
		ipLookup,
		idnaLookup,
		txtLookup,
		nsLookup,
		cnameLookup,

		//geoIPLookup,
	)
}

// mxLookupFunc
func mxLookupFunc(tr TypoResult) (results []TypoResult) {
	records, _ := net.LookupMX(tr.Domain.String())
	for _, record := range records {
		tr.Data["MX"] = strings.TrimSuffix(record.Host, ".")
	}
	results = append(results, TypoResult{tr.Domain, tr.Typo, tr.Live, tr.Data})
	return
}

// nsLookupFunc
func nsLookupFunc(tr TypoResult) (results []TypoResult) {
	records, _ := net.LookupMX(tr.Domain.String())
	//fmt.Println(records)
	for _, record := range records {
		tr.Data["NS"] = strings.TrimSuffix(record.Host, ".")
	}
	results = append(results, TypoResult{tr.Domain, tr.Typo, tr.Live, tr.Data})
	return
}

// cnameLookupFunc
func cnameLookupFunc(tr TypoResult) (results []TypoResult) {
	records, _ := net.LookupCNAME(tr.Domain.String())
	//fmt.Println(records)
	for _, record := range records {
		tr.Data["CNAME"] = strings.TrimSuffix(string(record), ".")
	}
	results = append(results, TypoResult{tr.Domain, tr.Typo, tr.Live, tr.Data})
	return
}

// ipLookupFunc
func ipLookupFunc(tr TypoResult) (results []TypoResult) {
	records, _ := net.LookupIP(tr.Domain.String())
	for _, record := range records {
		IPv4 := record.To4().String()
		if IPv4 != "" {
			tr.Data["IPv4"] = IPv4
			tr.Live = true
		}
		IPv6 := record.To16().String()
		if IPv6 != "" {
			tr.Data["IPv6"] = IPv6
			tr.Live = true
		}
	}
	results = append(results, TypoResult{tr.Domain, tr.Typo, tr.Live, tr.Data})
	return
}

// txtLookupFunc
func txtLookupFunc(tr TypoResult) (results []TypoResult) {
	records, _ := net.LookupTXT(tr.Domain.String())
	for _, record := range records {
		tr.Data["TXT"] = record
	}
	results = append(results, TypoResult{tr.Domain, tr.Typo, tr.Live, tr.Data})
	return
}

// geoIPLookupFunc
func geoIPLookupFunc(tr TypoResult) (results []TypoResult) {
	tr.Data["GEO"] = "GEO GEO"

	results = append(results, TypoResult{tr.Domain, tr.Typo, tr.Live, tr.Data})
	return
}

// idnaFunc
func idnaFunc(tr TypoResult) (results []TypoResult) {
	tr.Data["IDNA"] = tr.Domain.Idna()
	results = append(results, TypoResult{tr.Domain, tr.Typo, tr.Live, tr.Data})
	return
}

// FRegister
func FRegister(name string, efunc ...Extra) {
	_, registered := FREGISTRY[strings.ToUpper(name)]
	if !registered {
		FREGISTRY[strings.ToUpper(name)] = efunc
	}
}

// FRetrieve
func FRetrieve(strs ...string) (results []Extra) {
	for _, f := range strs {
		value, ok := FREGISTRY[strings.ToUpper(f)]
		if ok {
			results = append(results, value...)
		} else {
			panic("Invalid option! ")
		}
	}
	if len(strs) == 0 {
		return FRetrieve("all")
	}
	return
}
