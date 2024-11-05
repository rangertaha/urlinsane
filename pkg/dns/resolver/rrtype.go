package resolver

import "github.com/miekg/dns"

// RRtype resource record type.
type RRtype uint16

const (
	// TypeA A record
	TypeA RRtype = RRtype(dns.TypeA)

	// TypeAAAA AAAA record
	TypeAAAA RRtype = RRtype(dns.TypeAAAA)

	// TypeCNAME CNAME record
	TypeCNAME RRtype = RRtype(dns.TypeCNAME)

	// TypeTXT TXT record
	TypeTXT RRtype = RRtype(dns.TypeTXT)

	// TypeMX MX record
	TypeMX RRtype = RRtype(dns.TypeMX)

	// TypeNS NS record
	TypeNS RRtype = RRtype(dns.TypeNS)
)

var rrtypeToString = map[RRtype]string{
	TypeA:     "A",
	TypeAAAA:  "AAAA",
	TypeCNAME: "CNAME",
	TypeTXT:   "TXT",
	TypeMX:    "MX",
	TypeNS:    "NS",
}

func (t RRtype) String() string {
	return rrtypeToString[t]
}
