// Copyright 2024 Rangertaha. All Rights Reserved.
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
package pkg

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strings"

// 	log "github.com/sirupsen/logrus"
// )

// type DnsRecord struct {
// 	Type  string `json:"type,omitempty"`
// 	Value string `json:"value,omitempty"`
// 	TTL   int    `json:"ttl,omitempty"`
// }

// type DnsRecords []DnsRecord

// func (d *DnsRecords) First(rtypes ...string) (value string) {
// 	for _, record := range *d {
// 		for _, rtype := range rtypes {
// 			if strings.EqualFold(record.Type, rtype) {
// 				return record.Value
// 			}
// 		}
// 	}
// 	return
// }

// func (d *DnsRecords) Array(rtypes ...string) (values []string) {
// 	for _, record := range *d {
// 		for _, rtype := range rtypes {
// 			if strings.EqualFold(record.Type, rtype) {
// 				values = append(values, record.Value)
// 			}
// 		}
// 	}
// 	return
// }

// func (d *DnsRecords) String(rtypes ...string) (values string) {
// 	for _, record := range *d {
// 		for _, rtype := range rtypes {
// 			if strings.EqualFold(record.Type, rtype) {
// 				values = fmt.Sprintf("%s %s", values, record.Value)
// 			}
// 		}
// 	}
// 	return
// }

// func (d *DnsRecords) Add(rtype string, ttl int, values ...string) {
// 	deduped := make(map[string]bool)
// 	for _, record := range *d {
// 		deduped[record.Value] = true
// 	}
// 	for _, value := range values {
// 		value = strings.TrimSpace(value)
// 		if _, ok := deduped[value]; !ok {
// 			*d = append(*d, DnsRecord{Type: rtype, TTL: ttl, Value: value})
// 		}
// 	}
// }

// func (d *DnsRecords) Json() json.RawMessage {
// 	records, err := json.Marshal(d)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	return json.RawMessage(records)
// }

// func (d *DnsRecords) Has(rtype string) bool {
// 	return len(*d) > 0
// }
