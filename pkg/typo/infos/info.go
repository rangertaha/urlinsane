// The MIT License (MIT)
//
// Copyright © 2018 Rangertaha <rangertaha@gmail.com>
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

package infos

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/glaslos/ssdeep"

	"github.com/cybersectech-org/urlinsane/pkg/typo"
)

type (
	Meta struct {
		Geo GeoCountry `json:"geo,omitempty"`
	}
)

// FREGISTRY is the registry for extra functions
var FREGISTRY = make(map[string][]typo.Module)

var levenshteinDistance = typo.Module{
	Code:        "LD",
	Name:        "Levenshtein Distance",
	Description: "The Levenshtein distance is a string metric for measuring the difference between two domains",
	Exe:         levenshteinDistanceFunc,
	Fields:      []string{"LD"},
}

var mxLookup = typo.Module{
	Code:        "MX",
	Name:        "MX Lookup",
	Description: "Checking for DNS's MX records",
	Exe:         mxLookupFunc,
	Fields:      []string{"MX"},
}

var txtLookup = typo.Module{
	Code:        "TXT",
	Name:        "TXT Lookup",
	Description: "Checking for DNS's TXT records",
	Exe:         txtLookupFunc,
	Fields:      []string{"TXT"},
}

var ipLookup = typo.Module{
	Code:        "IP",
	Name:        "IP Lookup",
	Description: "Checking for IP address",
	Exe:         ipLookupFunc,
	Fields:      []string{"IPv4", "IPv6"},
}

var nsLookup = typo.Module{
	Code:        "NS",
	Name:        "NS Lookup",
	Description: "Checks DNS NS records",
	Exe:         nsLookupFunc,
	Fields:      []string{"NS"},
}

var cnameLookup = typo.Module{
	Code:        "CNAME",
	Name:        "CNAME Lookup",
	Description: "Checks DNS CNAME records",
	Exe:         cnameLookupFunc,
	Fields:      []string{"CNAME"},
}

var idnaLookup = typo.Module{
	Code:        "IDNA",
	Name:        "IDNA Domain",
	Description: "Show international domain name",
	Exe:         idnaFunc,
	Fields:      []string{"IDNA"},
}

var ssdeepLookup = typo.Module{
	Code:        "SIM",
	Name:        "Domain Similarity",
	Description: "Show domain content similarity",
	Exe:         ssdeepFunc,
	Fields:      []string{"IPv4", "IPv6", "SIM"},
}

var redirectLookup = typo.Module{
	Code:        "301",
	Name:        "Redirected Domain",
	Description: "Show domains redirects",
	Exe:         redirectLookupFunc,
	Fields:      []string{"IPv4", "IPv6", "Redirect"},
}

var whoisLookup = typo.Module{
	Code:        "WHOIS",
	Name:        "Show whois info",
	Description: "Query whois for additional information",
	Exe:         whoisLookupFunc,
	Fields:      []string{"WHOIS?"},
}

func init() {
	FRegister("LD", levenshteinDistance)
	FRegister("IDNA", idnaLookup)
	FRegister("MX", mxLookup)
	FRegister("IP", ipLookup)
	FRegister("TXT", txtLookup)
	FRegister("NS", nsLookup)
	FRegister("CNAME", cnameLookup)
	FRegister("SIM", ssdeepLookup)
	FRegister("301", redirectLookup)

	//FRegister("WHOIS", whoisLookup)
	FRegister("GEO", geoIPLookup)

	FRegister("ALL",
		levenshteinDistance,
		mxLookup,
		ipLookup,
		idnaLookup,
		txtLookup,
		nsLookup,
		cnameLookup,
		ssdeepLookup,
		// liveFilter,
		redirectLookup,

		//whoisLookup,
		geoIPLookup,
	)
}

// levenshteinDistanceFunc
func levenshteinDistanceFunc(tr typo.Result) (results []typo.Result) {
	domain := tr.Original.String()
	variant := tr.Variant.String()
	tr.Data["LD"] = strconv.Itoa(Levenshtein(domain, variant))
	tr.Meta["levenshtein"] = Levenshtein(domain, variant)
	results = append(results, typo.Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Live: tr.Live, Data: tr.Data, Meta: tr.Meta})
	return
}

// mxLookupFunc
func mxLookupFunc(tr typo.Result) (results []typo.Result) {
	records, _ := net.LookupMX(tr.Variant.String())
	tr.Meta["MX"] = records
	for _, record := range records {
		record := strings.TrimSuffix(record.Host, ".")
		if !strings.Contains(tr.Data["MX"], record) {
			tr.Data["MX"] = strings.TrimSpace(tr.Data["MX"] + "\n" + record)
		}
	}
	results = append(results, typo.Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Live: tr.Live, Data: tr.Data, Meta: tr.Meta})
	return
}

// nsLookupFunc
func nsLookupFunc(tr typo.Result) (results []typo.Result) {
	records, _ := net.LookupNS(tr.Variant.String())
	tr.Meta["NS"] = records
	for _, record := range records {
		record := strings.TrimSuffix(record.Host, ".")
		if !strings.Contains(tr.Data["NS"], record) {
			tr.Data["NS"] = strings.TrimSpace(tr.Data["NS"] + "\n" + record)
		}
	}
	results = append(results, typo.Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Live: tr.Live, Data: tr.Data, Meta: tr.Meta})
	return
}

// cnameLookupFunc
func cnameLookupFunc(tr typo.Result) (results []typo.Result) {
	records, _ := net.LookupCNAME(tr.Variant.String())
	tr.Meta["CNAME"] = records
	for _, record := range records {
		tr.Data["CNAME"] = strings.TrimSuffix(string(record), ".")
	}
	results = append(results, typo.Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Live: tr.Live, Data: tr.Data, Meta: tr.Meta})
	return
}

// ipLookupFunc
func ipLookupFunc(tr typo.Result) (results []typo.Result) {
	results = append(results, checkIP(tr))
	return
}

// txtLookupFunc
func txtLookupFunc(tr typo.Result) (results []typo.Result) {
	records, _ := net.LookupTXT(tr.Variant.String())
	tr.Meta["TXT"] = records
	for _, record := range records {
		tr.Data["TXT"] = record
	}
	results = append(results, typo.Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Live: tr.Live, Data: tr.Data, Meta: tr.Meta})
	return
}

// idnaFunc
func idnaFunc(tr typo.Result) (results []typo.Result) {

	tr.Data["IDNA"] = tr.Variant.Idna()
	tr.Meta["IDNA"] = tr.Variant.Idna()
	results = append(results, typo.Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Live: tr.Live, Data: tr.Data, Meta: tr.Meta})
	return
}

func ssdeepFunc(tr typo.Result) (results []typo.Result) {
	tr = checkIP(tr)
	if tr.Live {
		var h1, h2 string
		{
			original, gerr := http.Get("http://" + tr.Original.String())
			tr.Meta["original"] = "" //original
			if gerr == nil {
				if o, err := ioutil.ReadAll(original.Body); err == nil {
					h1, _ = ssdeep.FuzzyBytes(o)
					tr.Meta["original-ssdeep"] = "" //h1
				}
			}
		}
		{
			variation, gerr := http.Get("http://" + tr.Variant.String())
			if gerr == nil {
				if v, err := ioutil.ReadAll(variation.Body); err == nil {
					h2, _ = ssdeep.FuzzyBytes(v)
					tr.Meta["variant-ssdeep"] = "" // h2
				}
			}
		}
		if h1 != "" && h2 != "" {
			if compare, err := ssdeep.Distance(h1, h2); err == nil {
				//fmt.Println(compare, h2, err)
				tr.Data["SIM"] = fmt.Sprintf("%d%s", compare, "%")
				tr.Meta["SIM"] = fmt.Sprintf("%d%s", compare, "%")
			}
		}
	}
	results = append(results, tr)
	return
}

// redirectLookupFunc
func redirectLookupFunc(tr typo.Result) (results []typo.Result) {
	tr = checkIP(tr)
	if tr.Live {
		variation, err := http.Get("http://" + tr.Variant.String())
		if err == nil {
			tr.Meta["Variant"] = "" //variation
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
				tr.Meta["Redirect"] = dm
			}
		}
	}
	results = append(results, tr)
	return
}

func whoisLookupFunc(tr typo.Result) (results []typo.Result) {
	return
}

func checkIP(tr typo.Result) typo.Result {
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
		tr.Meta["IP"] = records
	}

	return typo.Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Live: tr.Live, Data: tr.Data, Meta: tr.Meta}
}

// FRegister ...
func FRegister(name string, efunc ...typo.Module) {
	_, registered := FREGISTRY[strings.ToUpper(name)]
	if !registered {
		FREGISTRY[strings.ToUpper(name)] = efunc
	}
}

// FRetrieve ...
func FRetrieve(strs ...string) (results []typo.Module) {
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
