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
package db

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	DomainID    uint
	Domain      *Domain `json:"domain,omitempty"`
	Uri         string  `json:"uri,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Body        string  `json:"body,omitempty"`

	// Ssdeep string `json:"ssdeep,omitempty"`

	Images []*Image `gorm:"many2many:images;" json:"images,omitempty"`
	Pages  []*Page  `gorm:"many2many:pages;" json:"pages,omitempty"`
	Files  []*File  `gorm:"many2many:files;" json:"files,omitempty"`

	// Language Analysis
	// Languages
	// Keywords
	// Topics
	// Vector
	// SSDeep
}

type Image struct {
	gorm.Model
	Url    string            `json:"uri,omitempty"`
	Hashes map[string]string `gorm:"serializer:json"      json:"hashes,omitempty"`
}

type File struct {
	gorm.Model
	Url    string            `json:"uri,omitempty"`
	Hashes map[string]string `gorm:"serializer:json"      json:"hashes,omitempty"`
}

// // Banner Details
// Status        string `json:"status,omitempty"`
// Protocol      string `json:"protocol,omitempty"`
// Headers       string `json:"headers,omitempty"`
// ContentLength int64  `json:"length,omitempty"`
// Proto         string `json:"proto,omitempty"`
// StatusCode    int    `json:"code,omitempty"`
// TLS           string `json:"tls,omitempty"`
// Cookies       string `json:"cookies,omitempty"`
// Trailer       string `json:"trailer,omitempty"`
