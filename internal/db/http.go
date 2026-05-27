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

type Page struct {
	Domain      *Domain `json:"domain,omitempty"`
	Uri         string  `json:"uri,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Body        string  `json:"body,omitempty"`

	// Media relations
	Images []*Image `json:"images,omitempty"`
	Pages  []*Page  `json:"pages,omitempty"`
	Files  []*File  `json:"files,omitempty"`
	Har    string   `json:"har,omitempty"`
}

type Image struct {
	Url    string            `json:"url,omitempty"`
	Hashes map[string]string `json:"hashes,omitempty"`
}

type File struct {
	Url    string            `json:"uri,omitempty"`
	Hashes map[string]string `json:"hashes,omitempty"`
}
