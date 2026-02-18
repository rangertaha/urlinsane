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
package web

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type Header map[string][]string

type Response struct {
	URL        string   `json:"url,omitempty"`
	Headers    Header   `json:"headers,omitempty"`
	StatusCode int      `json:"status,omitempty"`
	SSDeep     string   `json:"ssdeep,omitempty"`
	Size       int      `json:"size,omitempty"`
	File       string   `json:"file,omitempty"`
	Keywords   []string `json:"keywords,omitempty"`
	HTML       HTML     `json:"html,omitempty"`
}

type HTML struct {
	Title    string     `json:"title,omitempty"`
	Keywords string     `json:"keywords,omitempty"`
	Meta     []Metatags `json:"meta,omitempty"`
	Texts    []string   `json:"text,omitempty"`
	Links    []string   `json:"links,omitempty"`
	Images   []string   `json:"images,omitempty"`
}

type Metatags struct {
	Property string `json:"property,omitempty"`
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
}

func (s *Response) Json() []byte {
	res, err := json.Marshal(s)
	if err != nil {
		log.Error(err)
	}
	return res
}

func (s *Response) String() string {

	return ""
}
