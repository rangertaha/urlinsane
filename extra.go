// Copyright © 2018 CyberSecTech Inc
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

	"github.com/glaslos/ssdeep"
	"net/http"

	"fmt"
	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/oschwald/geoip2-golang"
	"io/ioutil"
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
	Description: "Show country location of ip address",
	Exec:        geoIPLookupFunc,
	Headers:     []string{"IPv4", "IPv6", "GEO"},
}

var idnaLookup = Extra{
	Code:        "IDNA",
	Name:        "IDNA Domain",
	Description: "Show international domain name",
	Exec:        idnaFunc,
	Headers:     []string{"IDNA"},
}

var ssdeepLookup = Extra{
	Code:        "SIM",
	Name:        "Domain Similarity",
	Description: "Show domain content similarity",
	Exec:        ssdeepFunc,
	Headers:     []string{"IPv4", "IPv6", "SIM"},
}

var liveFilter = Extra{
	Code:        "LIVE",
	Name:        "Online domians",
	Description: "Show domains with ip addresses only",
	Exec:        liveFilterFunc,
	Headers:     []string{"IPv4", "IPv6"},
}

var redirectLookup = Extra{
	Code:        "301",
	Name:        "Redirected Domain",
	Description: "Show domains redirects",
	Exec:        redirectLookupFunc,
	Headers:     []string{"IPv4", "IPv6", "Redirect"},
}

var whoisLookup = Extra{
	Code:        "WHOIS",
	Name:        "Show whois info",
	Description: "Query whois for additional information",
	Exec:        whoisLookupFunc,
	Headers:     []string{"WHOIS?"},
}

func init() {
	FRegister("IDNA", idnaLookup)
	FRegister("MX", mxLookup)
	FRegister("IP", ipLookup)
	FRegister("TXT", txtLookup)
	FRegister("NS", nsLookup)
	FRegister("CNAME", cnameLookup)
	FRegister("SIM", ssdeepLookup)
	FRegister("LIVE", liveFilter)
	FRegister("301", redirectLookup)

	//FRegister("WHOIS", whoisLookup)
	FRegister("GEO", geoIPLookup)

	FRegister("ALL",
		mxLookup,
		ipLookup,
		idnaLookup,
		txtLookup,
		nsLookup,
		cnameLookup,
		ssdeepLookup,
		liveFilter,
		redirectLookup,

		//whoisLookup,
		geoIPLookup,
	)
}

// mxLookupFunc
func mxLookupFunc(tr TypoResult) (results []TypoResult) {
	records, _ := net.LookupMX(tr.Variant.String())
	for _, record := range records {
		record := strings.TrimSuffix(record.Host, ".")
		if !strings.Contains(tr.Data["MX"], record) {
			tr.Data["MX"] = strings.TrimSpace(tr.Data["MX"] + "\n" + record)
		}
	}
	results = append(results, TypoResult{tr.Original, tr.Variant, tr.Typo, tr.Live, tr.Data})
	return
}

// nsLookupFunc
func nsLookupFunc(tr TypoResult) (results []TypoResult) {
	records, _ := net.LookupNS(tr.Variant.String())
	for _, record := range records {
		record := strings.TrimSuffix(record.Host, ".")
		if !strings.Contains(tr.Data["NS"], record) {
			tr.Data["NS"] = strings.TrimSpace(tr.Data["NS"] + "\n" + record)
		}
	}
	results = append(results, TypoResult{tr.Original, tr.Variant, tr.Typo, tr.Live, tr.Data})
	return
}

// cnameLookupFunc
func cnameLookupFunc(tr TypoResult) (results []TypoResult) {
	records, _ := net.LookupCNAME(tr.Variant.String())
	for _, record := range records {
		tr.Data["CNAME"] = strings.TrimSuffix(string(record), ".")
	}
	results = append(results, TypoResult{tr.Original, tr.Variant, tr.Typo, tr.Live, tr.Data})
	return
}

// ipLookupFunc
func ipLookupFunc(tr TypoResult) (results []TypoResult) {
	results = append(results, checkIP(tr))
	return
}

// txtLookupFunc
func txtLookupFunc(tr TypoResult) (results []TypoResult) {
	records, _ := net.LookupTXT(tr.Variant.String())
	for _, record := range records {
		tr.Data["TXT"] = record
	}
	results = append(results, TypoResult{tr.Original, tr.Variant, tr.Typo, tr.Live, tr.Data})
	return
}

// geoIPLookupFunc
func geoIPLookupFunc(tr TypoResult) (results []TypoResult) {
	tr = checkIP(tr)
	if tr.Live {
		geolite2CityMmdb, err := Asset("GeoLite2-Country.mmdb")
		if err != nil {
			// Asset was not found.
		}

		db, err := geoip2.FromBytes(geolite2CityMmdb)
		if err != nil {
			fmt.Print(err)
		}
		defer db.Close()

		ipv4s, ok := tr.Data["IPv4"]
		if ok {
			ips := strings.Split(ipv4s, "\n")
			for _, ip4 := range ips {
				ip := net.ParseIP(ip4)
				if ip != nil {
					record, err := db.Country(ip)
					if err != nil {
						fmt.Print(err)
					}
					tr.Data["GEO"] = fmt.Sprint(record.Country.Names["en"])

				}

			}
		}
	}

	// If you are using strings that may be invalid, check that ip is not nil

	results = append(results, TypoResult{tr.Original, tr.Variant, tr.Typo, tr.Live, tr.Data})
	return
}

// idnaFunc
func idnaFunc(tr TypoResult) (results []TypoResult) {

	tr.Data["IDNA"] = tr.Variant.Idna()
	results = append(results, TypoResult{tr.Original, tr.Variant, tr.Typo, tr.Live, tr.Data})
	return
}

func ssdeepFunc(tr TypoResult) (results []TypoResult) {
	tr = checkIP(tr)
	if tr.Live {
		var h1, h2 string
		{
			original, gerr := http.Get("http://" + tr.Original.String())
			if gerr == nil {
				if o, err := ioutil.ReadAll(original.Body); err == nil {
					h1, _ = ssdeep.FuzzyBytes(o)
					//fmt.Println(h1, err)
				}
			}
		}
		{
			variation, gerr := http.Get("http://" + tr.Variant.String())
			if gerr == nil {
				if v, err := ioutil.ReadAll(variation.Body); err == nil {
					h2, _ = ssdeep.FuzzyBytes(v)
					//fmt.Println(h2, err)
				}
			}
		}
		if h1 != "" && h2 != "" {
			if compare, err := ssdeep.Distance(h1, h2); err == nil {
				//fmt.Println(compare, h2, err)
				tr.Data["SIM"] = fmt.Sprintf("%d%s", compare, "%")
			}
		}
	}
	results = append(results, tr)
	return
}

// liveFilterFunc
func liveFilterFunc(tr TypoResult) (results []TypoResult) {
	tr = checkIP(tr)
	if tr.Live {
		results = append(results, TypoResult{tr.Original, tr.Variant, tr.Typo, tr.Live, tr.Data})
	}
	return
}

// redirectLookupFunc
func redirectLookupFunc(tr TypoResult) (results []TypoResult) {
	tr = checkIP(tr)
	if tr.Live {
		variation, err := http.Get("http://" + tr.Variant.String())
		if err == nil {
			str := variation.Request.URL.String()
			subdomain := domainutil.Subdomain(str)
			domain := domainutil.DomainPrefix(str)
			suffix := domainutil.DomainSuffix(str)
			if domain == "" {
				domain = str
			}
			dm := Domain{subdomain, domain, suffix}
			if tr.Original.String() != dm.String() {
				tr.Data["Redirect"] = dm.String()
			}
		}
	}
	results = append(results, tr)
	return
}

func whoisLookupFunc(tr TypoResult) (results []TypoResult) {
	return
}

func checkIP(tr TypoResult) TypoResult {
	ip4, _ := tr.Data["IPv4"]
	ip6, _ := tr.Data["IPv6"]
	if ip4 == "" || ip6 == "" {
		records, _ := net.LookupIP(tr.Variant.String())
		for _, record := range records {
			dotlen := strings.Count(record.String(), ".")
			if dotlen == 3 {
				if !strings.Contains(tr.Data["IPv4"], record.String()) {
					tr.Data["IPv4"] = strings.TrimSpace(tr.Data["IPv4"] + "\n" + record.String())
				}
				tr.Live = true
			}
			clen := strings.Count(record.String(), ":")
			if clen == 5 {
				if !strings.Contains(tr.Data["IPv6"], record.String()) {
					tr.Data["IPv6"] = strings.TrimSpace(tr.Data["IPv6"] + "\n" + record.String())
				}
			}
			tr.Live = true
		}
	}

	return TypoResult{tr.Original, tr.Variant, tr.Typo, tr.Live, tr.Data}
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
		}
	}
	if len(strs) == 0 {
		return FRetrieve("all")
	}
	return
}
