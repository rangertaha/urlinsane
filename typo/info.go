// Copyright (C) 2024  Rangertaha <rangertaha@gmail.com>
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
package typo

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	// "github.com/davecgh/go-spew/spew"
	"github.com/glaslos/ssdeep"
	dnsLib "github.com/rangertaha/urlinsane/hackingo/net/dns"
	geoLib "github.com/rangertaha/urlinsane/hackingo/net/geoip"
	httpLib "github.com/rangertaha/urlinsane/hackingo/net/http"
	nlpLib "github.com/rangertaha/urlinsane/hackingo/nlp"
)

// Extras is the registry for extra functions
var Extras = NewRegistry()

var levenshteinDistance = Module{
	Code:        "LD",
	Name:        "Levenshtein Distance",
	Description: "The Levenshtein distance between strings",
	Exe:         levenshteinDistanceFunc,
	Fields:      []string{"LD"},
}

var mxLookup = Module{
	Code:        "MX",
	Name:        "MX Lookup",
	Description: "Checking for DNS's MX records",
	Exe:         mxLookupFunc,
	Fields:      []string{"MX"},
}

var txtLookup = Module{
	Code:        "TXT",
	Name:        "TXT Lookup",
	Description: "Checking for DNS's TXT records",
	Exe:         txtLookupFunc,
	Fields:      []string{"TXT"},
}

var ipLookup = Module{
	Code:        "IP",
	Name:        "IP Lookup",
	Description: "Checking for IP address",
	Exe:         ipLookupFunc,
	Fields:      []string{"IPv4", "IPv6"},
}

var nsLookup = Module{
	Code:        "NS",
	Name:        "NS Lookup",
	Description: "Checks DNS NS records",
	Exe:         nsLookupFunc,
	Fields:      []string{"NS"},
}

var cnameLookup = Module{
	Code:        "CNAME",
	Name:        "CNAME Lookup",
	Description: "Checks DNS CNAME records",
	Exe:         cnameLookupFunc,
	Fields:      []string{"CNAME"},
}

var geoIPLookup = Module{
	Code:        "GEO",
	Name:        "GeoIP Lookup",
	Description: "Show country location of ip address",
	Exe:         geoIPLookupFunc,
	Fields:      []string{"IPv4", "IPv6", "GEO"},
}

var idnaLookup = Module{
	Code:        "IDNA",
	Name:        "IDNA Domain",
	Description: "Show international domain name",
	Exe:         idnaFunc,
	Fields:      []string{"IDNA"},
}

var ssdeepLookup = Module{
	Code:        "SIM",
	Name:        "Domain Similarity",
	Description: "Show domain content similarity",
	Exe:         ssdeepFunc,
	Fields:      []string{"IPv4", "IPv6", "SIM"},
}

var whoisLookup = Module{
	Code:        "WHOIS",
	Name:        "Show whois info",
	Description: "Query whois for additional information",
	Exe:         whoisLookupFunc,
	Fields:      []string{"WHOIS?"},
}

var httpLookup = Module{
	Code:        "HTTP",
	Name:        "Get Request",
	Description: "Get http related information",
	Exe:         httpLookupFunc,
	Fields:      []string{"IPv4", "IPv6", "SIZE", "Redirect"},
}

func init() {
	Extras.Set("LD", levenshteinDistance)
	Extras.Set("IDNA", idnaLookup)
	Extras.Set("IP", ipLookup)
	Extras.Set("HTTP", httpLookup)
	Extras.Set("MX", mxLookup)
	Extras.Set("TXT", txtLookup)
	Extras.Set("NS", nsLookup)
	Extras.Set("CNAME", cnameLookup)
	Extras.Set("SIM", ssdeepLookup)

	//FRegister("WHOIS", whoisLookup)
	Extras.Set("GEO", geoIPLookup)

	Extras.Set("ALL",
		levenshteinDistance,
		idnaLookup,
		ipLookup,
		httpLookup,
		mxLookup,
		txtLookup,
		nsLookup,
		cnameLookup,
		ssdeepLookup,

		//whoisLookup,
		geoIPLookup,
	)
}

// httpLookupFunc
func httpLookupFunc(tr Result) (results []Result) {
	if tr := checkIP(tr); tr.Variant.Live {
		httpReq, gerr := http.Get("http://" + tr.Variant.String())
		if gerr == nil {
			tr.Variant.Meta.HTTP = httpLib.NewResponse(httpReq)
			// spew.Dump(original)

			str := httpReq.Request.URL.String()
			subdomain := domainutil.Subdomain(str)
			domain := domainutil.DomainPrefix(str)
			suffix := domainutil.DomainSuffix(str)
			if domain == "" {
				domain = str
			}
			dm := Domain{subdomain, domain, suffix, tr.Variant.Meta, true}
			if tr.Variant.String() != dm.String() {
				tr.Data["Redirect"] = dm.String()
				tr.Variant.Meta.Redirect = dm.String()
			}
		}
	}
	results = append(results, tr)
	return
}

// levenshteinDistanceFunc
func levenshteinDistanceFunc(tr Result) (results []Result) {
	domain := tr.Original.String()
	variant := tr.Variant.String()
	tr.Data["LD"] = strconv.Itoa(nlpLib.Levenshtein(domain, variant))
	tr.Variant.Meta.Levenshtein = nlpLib.Levenshtein(domain, variant)
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// mxLookupFunc
func mxLookupFunc(tr Result) (results []Result) {
	records, _ := net.LookupMX(tr.Variant.String())
	tr.Variant.Meta.DNS.MX = dnsLib.NewMX(records...)
	for _, record := range records {
		record := strings.TrimSuffix(record.Host, ".")
		if !strings.Contains(tr.Data["MX"], record) {
			tr.Data["MX"] = strings.TrimSpace(tr.Data["MX"] + "\n" + record)
		}
	}
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// nsLookupFunc
func nsLookupFunc(tr Result) (results []Result) {
	records, _ := net.LookupNS(tr.Variant.String())
	tr.Variant.Meta.DNS.NS = dnsLib.NewNS(records...)
	for _, record := range records {
		record := strings.TrimSuffix(record.Host, ".")
		if !strings.Contains(tr.Data["NS"], record) {
			tr.Data["NS"] = strings.TrimSpace(tr.Data["NS"] + "\n" + record)
		}
	}
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// cnameLookupFunc
func cnameLookupFunc(tr Result) (results []Result) {
	records, _ := net.LookupCNAME(tr.Variant.String())
	// tr.Variant.Meta.DNS.CName = records
	for _, record := range records {
		tr.Data["CNAME"] = strings.TrimSuffix(string(record), ".")
		tr.Variant.Meta.DNS.CName = append(tr.Variant.Meta.DNS.CName, string(record))
	}
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// ipLookupFunc
func ipLookupFunc(tr Result) (results []Result) {
	results = append(results, checkIP(tr))
	return
}

// txtLookupFunc
func txtLookupFunc(tr Result) (results []Result) {
	records, _ := net.LookupTXT(tr.Variant.String())
	tr.Variant.Meta.DNS.TXT = records
	for _, record := range records {
		tr.Data["TXT"] = strings.TrimSpace(tr.Data["TXT"] + "\n" + record)
		tr.Variant.Meta.DNS.TXT = append(tr.Variant.Meta.DNS.TXT, record)
	}
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// geoIPLookupFunc
func geoIPLookupFunc(tr Result) (results []Result) {
	tr = checkIP(tr)
	if tr.Variant.Live {
		_, ok := tr.Data["IPv4"]
		if ok {
			for _, ip4 := range tr.Variant.Meta.DNS.IPv4 {
				if ip4 != "" {
					record, _ := geoLib.GeoIP(ip4)
					// if err != nil {
					// 	fmt.Print(err)
					// }
					tr.Data["GEO"] = fmt.Sprint(record.Country.Names["en"])
					tr.Variant.Meta.Geo = *record
				}
			}
		}
	}

	// If you are using strings that may be invalid, check that ip is not nil
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// idnaFunc
func idnaFunc(tr Result) (results []Result) {
	tr.Data["IDNA"] = tr.Variant.Idna()
	tr.Variant.Meta.IDNA = tr.Variant.Idna()
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

func ssdeepFunc(tr Result) (results []Result) {
	tr = checkIP(tr)
	if tr.Original.Live {
		var h1, h2 string
		h1, _ = ssdeep.FuzzyBytes([]byte(tr.Original.Meta.HTTP.Body))
		tr.Original.Meta.SSDeep = h1

		if tr.Variant.Live {
			h2, _ = ssdeep.FuzzyBytes([]byte(tr.Variant.Meta.HTTP.Body))
			tr.Variant.Meta.SSDeep = h2
		}

		if compare, err := ssdeep.Distance(h1, h2); err == nil {
			tr.Data["SIM"] = fmt.Sprintf("%d%s", compare, "%")
			tr.Variant.Meta.Similarity = compare
		}
	}
	results = append(results, tr)
	return
}

func whoisLookupFunc(tr Result) (results []Result) {
	return
}

func checkIP(tr Result) Result {
	if tr.Variant.Meta.DNS.ipCheck == false {
		records, _ := net.LookupIP(tr.Variant.String())
		// if err != nil {
		// 	fmt.Println(err)
		// }
		for _, record := range uniqIP(records) {
			dotlen := strings.Count(record, ".")
			if dotlen == 3 {
				if !strings.Contains(tr.Data["IPv4"], record) {
					tr.Data["IPv4"] = strings.TrimSpace(tr.Data["IPv4"] + "\n" + record)
					tr.Variant.Meta.DNS.IPv4 = append(tr.Variant.Meta.DNS.IPv4, record)
				}
				tr.Variant.Live = true
			}
			clen := strings.Count(record, ":")
			if clen == 5 {
				if !strings.Contains(tr.Data["IPv6"], record) {
					tr.Data["IPv6"] = strings.TrimSpace(tr.Data["IPv6"] + "\n" + record)
					tr.Variant.Meta.DNS.IPv6 = append(tr.Variant.Meta.DNS.IPv6, record)
				}
				tr.Variant.Live = true
			}
		}
		tr.Variant.Meta.DNS.ipCheck = true
	}

	return Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data}
}

func uniqIP(list []net.IP) (ulist []string) {
	uinq := map[string]bool{}
	for _, l := range list {
		uinq[l.String()] = true
	}
	for k := range uinq {
		ulist = append(ulist, k)
	}
	return
}
