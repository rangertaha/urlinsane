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

	LocationID *uint
	Location   *Location `json:"location,omitempty"`
	IPs        []*IP     `gorm:"many2many:ipaddrs;"  json:"ips,omitempty"`
}

type IP struct {
	gorm.Model
	Address string  `json:"address,omitempty"`
	Type    string  `json:"type,omitempty"`
	Ports   []*Port `gorm:"many2many:ipports;"   json:"ports,omitempty"`
}

type Port struct {
	gorm.Model
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
