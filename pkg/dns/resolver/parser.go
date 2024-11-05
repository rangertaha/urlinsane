package resolver

import (
	"strings"

	"github.com/miekg/dns"
)

type msgParser struct{}

func (s *msgParser) Parse(msg *dns.Msg) []Record {
	records := []Record{}

	for _, answer := range msg.Answer {
		// Invalid question, discard answer
		if len(msg.Question) == 0 {
			continue
		}

		switch t := answer.(type) {
		case *dns.A:
			rec := Record{
				Question: strings.TrimSuffix(msg.Question[0].Name, "."),
				Type:     TypeA,
				Answer:   t.A.String(),
			}
			records = append(records, rec)
		case *dns.AAAA:
			rec := Record{
				Question: strings.TrimSuffix(msg.Question[0].Name, "."),
				Type:     TypeAAAA,
				Answer:   t.AAAA.String(),
			}
			records = append(records, rec)
		case *dns.CNAME:
			rec := Record{
				Question: strings.TrimSuffix(msg.Question[0].Name, "."),
				Type:     TypeCNAME,
				Answer:   strings.TrimSuffix(t.Target, "."),
			}
			records = append(records, rec)
		case *dns.TXT:
			for _, txt := range t.Txt {
				rec := Record{
					Question: strings.TrimSuffix(msg.Question[0].Name, "."),
					Type:     TypeTXT,
					Answer:   txt,
				}
				records = append(records, rec)
			}
		case *dns.MX:
			rec := Record{
				Question: strings.TrimSuffix(msg.Question[0].Name, "."),
				Type:     TypeMX,
				Answer:   strings.TrimSuffix(t.Mx, "."),
			}
			records = append(records, rec)
		case *dns.NS:
			rec := Record{
				Question: strings.TrimSuffix(msg.Question[0].Name, "."),
				Type:     TypeNS,
				Answer:   strings.TrimSuffix(t.Ns, "."),
			}
			records = append(records, rec)
		}
	}

	return records
}
