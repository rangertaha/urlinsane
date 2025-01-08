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

type Device struct {
	gorm.Model
	Name string `json:"name,omitempty"`

	Addreses []*Address `gorm:"many2many:devaddrs;"  json:"ips,omitempty"`
}

type Address struct {
	gorm.Model
	Addr       string    `gorm:"unique"                json:"address"`
	Type       string    `                             json:"type"`
	Ports      []*Port   `gorm:"many2many:addrports;"  json:"ports,omitempty"`
	Domians    []*Domain `gorm:"many2many:domaddrs;"   json:"domains,omitempty"`
	LocationID *uint
	Location   *Location `                             json:"location,omitempty"`
}

type Port struct {
	gorm.Model
	Proto     string   `json:"proto,omitempty"`
	Number    int      `json:"num,omitempty"`
	State     string   `json:"state,omitempty"`
	Service   *Service `json:"service,omitempty"`
	ServiceID *uint
}

type Service struct {
	gorm.Model
	Name   string `json:"name,omitempty"`
	Banner string `json:"banner,omitempty"`
}
