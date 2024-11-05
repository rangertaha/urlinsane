package resolver

import (
	"net"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func stubRR() map[string]dns.RR {
	records := make(map[string]dns.RR)

	A := &dns.A{}
	A.Hdr.Rrtype = dns.TypeA
	A.A = net.ParseIP("127.0.0.1")
	records["A"] = A

	AAAA := &dns.AAAA{}
	AAAA.Hdr.Rrtype = dns.TypeAAAA
	AAAA.AAAA = net.ParseIP("::1")
	records["AAAA"] = AAAA

	CNAME := &dns.CNAME{}
	CNAME.Hdr.Rrtype = dns.TypeCNAME
	CNAME.Target = "foo.bar."
	records["CNAME"] = CNAME

	TXT := &dns.TXT{}
	TXT.Hdr.Rrtype = dns.TypeCNAME
	TXT.Txt = []string{"foo", "bar"}
	records["TXT"] = TXT

	MX := &dns.MX{}
	MX.Hdr.Rrtype = dns.TypeMX
	MX.Mx = "mx.foo.bar."
	records["MX"] = MX

	NS := &dns.NS{}
	NS.Hdr.Rrtype = dns.TypeNS
	NS.Ns = "ns.foo.bar."
	records["NS"] = NS

	return records
}

func stubRecords() map[string]Record {
	records := make(map[string]Record)

	A := Record{}
	A.Type = TypeA
	A.Answer = "127.0.0.1"
	records["A"] = A

	AAAA := Record{}
	AAAA.Type = TypeAAAA
	AAAA.Answer = "::1"
	records["AAAA"] = AAAA

	CNAME := Record{}
	CNAME.Type = TypeCNAME
	CNAME.Answer = "foo.bar"
	records["CNAME"] = CNAME

	TXT1 := Record{}
	TXT1.Type = TypeTXT
	TXT1.Answer = "foo"
	records["TXT1"] = TXT1

	TXT2 := Record{}
	TXT2.Type = TypeTXT
	TXT2.Answer = "bar"
	records["TXT2"] = TXT2

	MX := Record{}
	MX.Type = TypeMX
	MX.Answer = "mx.foo.bar"
	records["MX"] = MX

	NS := Record{}
	NS.Type = TypeNS
	NS.Answer = "ns.foo.bar"
	records["NS"] = NS

	return records
}

func TestParse(t *testing.T) {
	RR := stubRR()
	records := stubRecords()

	testTable := []struct {
		name       string
		haveQuery  string
		haveAnswer *dns.Msg
		want       []Record
	}{
		{
			name:      "Two Records",
			haveQuery: "foo.bar",
			haveAnswer: &dns.Msg{
				Question: []dns.Question{
					{Name: "foo.bar."},
				},
				Answer: []dns.RR{
					RR["A"],
					RR["CNAME"],
				},
			},
			want: []Record{
				records["A"],
				records["CNAME"],
			},
		},
		{
			name:      "AAAA",
			haveQuery: "foo.bar",
			haveAnswer: &dns.Msg{
				Question: []dns.Question{
					{Name: "foo.bar."},
				},
				Answer: []dns.RR{
					RR["AAAA"],
				},
			},
			want: []Record{
				records["AAAA"],
			},
		},
		{
			name:      "TXT",
			haveQuery: "foo.bar",
			haveAnswer: &dns.Msg{
				Question: []dns.Question{
					{Name: "foo.bar."},
				},
				Answer: []dns.RR{
					RR["TXT"]}},
			want: []Record{
				records["TXT1"],
				records["TXT2"],
			},
		},
		{
			name:      "MX",
			haveQuery: "foo.bar",
			haveAnswer: &dns.Msg{
				Question: []dns.Question{
					{Name: "foo.bar."},
				},
				Answer: []dns.RR{
					RR["MX"],
				},
			},
			want: []Record{
				records["MX"],
			},
		},
		{
			name:      "NS",
			haveQuery: "foo.bar",
			haveAnswer: &dns.Msg{
				Question: []dns.Question{
					{Name: "foo.bar."},
				},
				Answer: []dns.RR{
					RR["NS"],
				},
			},
			want: []Record{
				records["NS"],
			},
		},
		{
			name:      "Empty Answer",
			haveQuery: "foo.bar",
			haveAnswer: &dns.Msg{
				Question: []dns.Question{},
				Answer:   []dns.RR{},
			},
			want: []Record{},
		},
		{
			name:      "Empty Question",
			haveQuery: "foo.bar",
			haveAnswer: &dns.Msg{
				Question: []dns.Question{},
				Answer: []dns.RR{
					RR["A"],
				},
			},
			want: []Record{},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			for i := range test.want {
				test.want[i].Question = test.haveQuery
			}

			msgParser := msgParser{}
			records := msgParser.Parse(test.haveAnswer)

			assert.EqualValues(t, test.want, records, test.name)
		})
	}
}
