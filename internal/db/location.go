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

type Location struct {
	ID         uint
	ParentID   uint
	Code       string     `gorm:"unique" json:"code,omitempty"`
	Name       string     `json:"name,omitempty"`
	TimeZone   string     `json:"timezone,omitempty"`
	Latitude   float64    `json:"lat,omitempty"`
	Longitude  float64    `json:"lon,omitempty"`
	PostalCode string     `json:"postcode,omitempty"`
}
