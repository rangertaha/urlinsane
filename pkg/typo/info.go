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

package typo

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/bobesa/go-domain-util/domainutil"
	// "github.com/davecgh/go-spew/spew"
	"github.com/glaslos/ssdeep"
	"github.com/oschwald/geoip2-golang"

	dnsLib "github.com/cybint/hackingo/net/dns"
	geoLib "github.com/cybint/hackingo/net/geoip"
	httpLib "github.com/cybint/hackingo/net/http"

	"github.com/cybersectech-org/urlinsane/pkg/datasets"
)

type (
	// Continent struct {
	// 	Code      string            `json:"code,omitempty"`
	// 	GeoNameID uint              `json:"geo_name,omitempty"`
	// 	Names     map[string]string `json:"names,omitempty"`
	// }
	// Country struct {
	// 	GeoNameID         uint              `json:"code,omitempty"`
	// 	IsInEuropeanUnion bool              `json:"european,omitempty"`
	// 	IsoCode           string            `json:"iso_code,omitempty"`
	// 	Names             map[string]string `json:"names,omitempty"`
	// }
	// RegisteredCountry struct {
	// 	GeoNameID         uint              `json:"geo_name,omitempty"`
	// 	IsInEuropeanUnion bool              `json:"european,omitempty"`
	// 	IsoCode           string            `json:"iso_code,omitempty"`
	// 	Names             map[string]string `json:"names,omitempty"`
	// }
	// RepresentedCountry struct {
	// 	GeoNameID         uint              `json:"geo_name,omitempty"`
	// 	IsInEuropeanUnion bool              `json:"european,omitempty"`
	// 	IsoCode           string            `json:"iso_code,omitempty"`
	// 	Names             map[string]string `json:"names,omitempty"`
	// 	Type              string            `json:"type,omitempty"`
	// }
	// Traits struct {
	// 	IsAnonymousProxy    bool `json:"is_anonymous_proxy,omitempty"`
	// 	IsSatelliteProvider bool `json:"is_satellite_provider,omitempty"`
	// }
	// GeoCountry struct {
	// 	Continent          Continent          `json:"continent,omitempty"`
	// 	Country            Country            `json:"country,omitempty"`
	// 	RegisteredCountry  RegisteredCountry  `json:"registered_country,omitempty"`
	// 	RepresentedCountry RepresentedCountry `json:"represented_country,omitempty"`
	// 	Traits             Traits             `json:"traits,omitempty"`
	// }

	// ConnectionState records basic TLS details about the connection.
	// ConnectionState struct {
	// 	Version                     uint16                `json:"version,omitempty"`                      // TLS version used by the connection (e.g. VersionTLS12); added in Go 1.3
	// 	Complete                    bool                  `json:"complete,omitempty"`                     // TLS handshake is complete
	// 	DidResume                   bool                  `json:"did_resume,omitempty"`                   // connection resumes a previous TLS connection; added in Go 1.1
	// 	CipherSuite                 uint16                `json:"cipher_suite,omitempty"`                 // cipher suite in use (TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, ...)
	// 	NegotiatedProtocol          string                `json:"negotiated_protocol,omitempty"`          // negotiated next protocol (not guaranteed to be from Config.NextProtos)
	// 	NegotiatedProtocolIsMutual  bool                  `json:"negotiated_protocolIsMutual,omitempty"`  // negotiated protocol was advertised by server (client side only)
	// 	ServerName                  string                `json:"server_name,omitempty"`                  // server name requested by client, if any (server side only)
	// 	PeerCertificates            []*x509.Certificate   `json:"-"`                                      // certificate chain presented by remote peer
	// 	VerifiedChains              [][]*x509.Certificate `json:"-"`                                      // verified chains built from PeerCertificates
	// 	SignedCertificateTimestamps [][]byte              `json:"signed_certificateTimestamps,omitempty"` // SCTs from the peer, if any; added in Go 1.5
	// 	OCSPResponse                []byte                `json:"ocsp_response,omitempty"`                // stapled OCSP response from peer, if any; added in Go 1.5
	// }

	// // Response represents the response from an HTTP request.
	// Response struct {
	// 	Status     string `json:"status,omitempty"`      // e.g. "200 OK"
	// 	StatusCode int    `json:"status_code,omitempty"` // e.g. 200
	// 	Proto      string `json:"proto,omitempty"`       // e.g. "HTTP/1.0"
	// 	ProtoMajor int    `json:"proto_major,omitempty"`
	// 	ProtoMinor int    `json:"proto_minor,omitempty"`

	// 	// A Header represents the key-value pairs in an HTTP header.
	// 	Header map[string][]string `json:"header,omitempty"`

	// 	// Body represents the response body.
	// 	Body string `json:"body,omitempty"`

	// 	// Length records the length of the associated content.
	// 	Length int64 `json:"length,omitempty"`

	// 	SSDeep string `json:"ssdeep,omitempty"`

	// 	// Contains transfer encodings from outer-most to inner-most. Value is
	// 	// nil, means that "identity" encoding is used.
	// 	Encoding []string `json:"encoding,omitempty"`

	// 	// Uncompressed reports whether the response was sent compressed but
	// 	// was decompressed by the http package.
	// 	Uncompressed bool `json:"uncompressed,omitempty"`

	// 	// After the body is read, Trailer will contain any trailer values
	// 	// sent by the server.
	// 	Trailer map[string][]string `json:"trailer,omitempty"`

	// 	// TLS contains information about the TLS connection on which the
	// 	// response was received. It is nil for unencrypted responses.
	// 	// The pointer is shared between responses and should not be
	// 	// modified.
	// 	TLS ConnectionState `json:"tls,omitempty"`
	// }

	DNS struct {
		IPv4  []string    `json:"ipv4,omitempty"`
		IPv6  []string    `json:"ip46,omitempty"`
		NS    []dnsLib.NS `json:"ns,omitempty"`
		MX    []dnsLib.MX `json:"mx,omitempty"`
		CName []string    `json:"cname,omitempty"`
		TXT   []string    `json:"txt,omitempty"`
	}

	Meta struct {
		Levenshtein int              `json:"Levenshtein,omitempty"`
		IDNA        string           `json:"idna,omitempty"`
		IP          []string         `json:"ip,omitempty"`
		Redirect    string           `json:"redirect,omitempty"`
		HTTP        httpLib.Response `json:"http,omitempty"`
		Geo         geoLib.Country   `json:"geo,omitempty"`
		DNS         DNS              `json:"dns,omitempty"`
		SSDeep      string           `json:"ssdeep,omitempty"`
		Similarity  int              `json:"similarity,omitempty"`
		// Whois    Whois      `json:"whois,omitempty"`
	}
)

// Extras is the registry for extra functions
var Extras = NewRegistry()

var levenshteinDistance = Module{
	Code:        "LD",
	Name:        "Levenshtein Distance",
	Description: "The Levenshtein distance is a string metric for measuring the difference between two domains",
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
	tr = checkIP(tr)
	if tr.Original.Live {
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
			if tr.Original.String() != dm.String() {
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
	tr.Data["LD"] = strconv.Itoa(Levenshtein(domain, variant))
	tr.Variant.Meta.Levenshtein = Levenshtein(domain, variant)
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// mxLookupFunc
func mxLookupFunc(tr Result) (results []Result) {
	tr = checkIP(tr)
	if tr.Original.Live {
		records, _ := net.LookupMX(tr.Variant.String())
		tr.Variant.Meta.DNS.MX = dnsLib.NewMX(records...)
		for _, record := range records {
			record := strings.TrimSuffix(record.Host, ".")
			if !strings.Contains(tr.Data["MX"], record) {
				tr.Data["MX"] = strings.TrimSpace(tr.Data["MX"] + "\n" + record)
			}
		}
	}
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// nsLookupFunc
func nsLookupFunc(tr Result) (results []Result) {
	tr = checkIP(tr)
	if tr.Original.Live {
		records, _ := net.LookupNS(tr.Variant.String())
		tr.Variant.Meta.DNS.NS = dnsLib.NewNS(records...)
		for _, record := range records {
			record := strings.TrimSuffix(record.Host, ".")
			if !strings.Contains(tr.Data["NS"], record) {
				tr.Data["NS"] = strings.TrimSpace(tr.Data["NS"] + "\n" + record)
			}
		}
	}
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// cnameLookupFunc
func cnameLookupFunc(tr Result) (results []Result) {
	tr = checkIP(tr)
	if tr.Original.Live {
		records, _ := net.LookupCNAME(tr.Variant.String())
		// tr.Variant.Meta.DNS.CName = records
		for _, record := range records {
			tr.Data["CNAME"] = strings.TrimSuffix(string(record), ".")
		}
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
	}
	results = append(results, Result{Original: tr.Original, Variant: tr.Variant, Typo: tr.Typo, Data: tr.Data})
	return
}

// geoIPLookupFunc
func geoIPLookupFunc(tr Result) (results []Result) {
	tr = checkIP(tr)
	if tr.Variant.Live {
		geolite2CityMmdb, err := datasets.Asset("GeoLite2-Country.mmdb")
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
					tr.Variant.Meta.Geo = *geoLib.NewCountry(record)
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
		{
			if tr.Original.Live {
				original, gerr := http.Get("http://" + tr.Original.String())
				tr.Original.Meta.HTTP = httpLib.NewResponse(original)
				if gerr == nil {
					if o, err := ioutil.ReadAll(original.Body); err == nil {
						h1, _ = ssdeep.FuzzyBytes(o)
						tr.Original.Meta.SSDeep = h1
					}
				}
			}
		}

		{
			variation, gerr := http.Get("http://" + tr.Variant.String())
			if gerr == nil {
				if v, err := ioutil.ReadAll(variation.Body); err == nil {
					h2, _ = ssdeep.FuzzyBytes(v)
					tr.Variant.Meta.SSDeep = h2
				}
			}
		}
		if h1 != "" && h2 != "" {
			if compare, err := ssdeep.Distance(h1, h2); err == nil {
				//fmt.Println(compare, h2, err)
				tr.Data["SIM"] = fmt.Sprintf("%d%s", compare, "%")
				tr.Variant.Meta.Similarity = compare
			}
		}
	}
	results = append(results, tr)
	return
}

// // redirectLookupFunc
// func redirectLookupFunc(tr Result) (results []Result) {
// 	tr = checkIP(tr)
// 	if tr.Live {
// 		variation, err := http.Get("http://" + tr.Variant.String())
// 		if err == nil {
// 			tr.Meta["Variant"] = "" //variation
// 			str := variation.Request.URL.String()
// 			subdomain := domainutil.Subdomain(str)
// 			domain := domainutil.DomainPrefix(str)
// 			suffix := domainutil.DomainSuffix(str)
// 			if domain == "" {
// 				domain = str
// 			}
// 			dm := Domain{subdomain, domain, suffix}
// 			if tr.Original.String() != dm.String() {
// 				tr.Data["Redirect"] = dm.String()
// 				tr.Meta["Redirect"] = dm
// 			}
// 		}
// 	}
// 	results = append(results, tr)
// 	return
// }

func whoisLookupFunc(tr Result) (results []Result) {
	return
}

func checkIP(tr Result) Result {
	if tr.Variant.Live == false {
		records, err := net.LookupIP(tr.Variant.String())
		if err != nil {
			fmt.Println(err)
		}
		for _, record := range uniqIP(records) {
			dotlen := strings.Count(record, ".")
			if dotlen == 3 {
				if !strings.Contains(tr.Data["IPv4"], record) {
					tr.Data["IPv4"] = strings.TrimSpace(tr.Data["IPv4"] + "\n" + record)
				}
				tr.Variant.Live = true
			}
			clen := strings.Count(record, ":")
			if clen == 5 {
				if !strings.Contains(tr.Data["IPv6"], record) {
					tr.Data["IPv6"] = strings.TrimSpace(tr.Data["IPv6"] + "\n" + record)
				}
				tr.Variant.Live = true
			}

		}
		tr.Variant.Meta.IPv4 = records
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
