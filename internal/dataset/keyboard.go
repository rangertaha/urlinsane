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
package dataset

const (
	ISO105 = `
	
	`
	ANSI104 = `
	
	`
	ABNT = `
	
	`
	OADG109A = `
	
	`
)

type Keyboard struct {
	ID          uint
	Code        string `gorm:"unique"`
	Name        string
	Description string
	Arrangement string      // ISO105, ANSI104, ABNT, OADG109A
	Languages   []*Language `gorm:"many2many:langboards;"`
	Layout      []Key       `gorm:"serializer:json"`
}

type Key struct {
	SC    string `json:"sc,omitempty"`
	VK    string `json:"vk,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
