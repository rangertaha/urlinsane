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

// type Service struct {
// 	Port   string `json:"port,omitempty"`
// 	Banner string    `json:"banner,omitempty"`
// }

// type Services map[string]string

// func (d *Services) First(ports ...string) (value string) {
// 	for _, record := range *d {
// 		for _, port := range ports {
// 			if strings.EqualFold(record.Port, port) {
// 				return record.Banner
// 			}
// 		}
// 	}
// 	return
// }

// func (d *Services) Array(ports ...string) (values []string) {
// 	for _, record := range *d {
// 		for _, port := range ports {
// 			if strings.EqualFold(record.Port, port) {
// 				values = append(values, record.Banner)
// 			}
// 		}
// 	}
// 	return
// }

// func (d *Services) String(rtypes ...string) (values string) {
// 	for _, record := range *d {
// 		for _, rtype := range rtypes {
// 			if strings.EqualFold(record.Port, rtype) {
// 				values = fmt.Sprintf("%s\n%s", values, record.Banner)
// 			}
// 		}
// 	}
// 	return
// }

// func (d *Services) Add(port string, banner string) {
// 	deduped := make(map[string]bool)
// 	for _, record := range *d {
// 		deduped[record.Port] = true
// 	}
// 	for _, value := range banners {
// 		if _, ok := deduped[value]; !ok {
// 			*d = append(*d, Service{Port: port, Banner: value})
// 		}
// 	}
// }

// func (d *Services) Json() json.RawMessage {
// 	records, err := json.Marshal(d)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	return json.RawMessage(records)
// }

// func (d *Services) Has(rtype string) bool {
// 	return len(*d) > 0
// }
