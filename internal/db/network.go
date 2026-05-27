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

// Device is a network host with one or more addresses.
type Device struct {
	Name     string     `json:"name,omitempty"`
	Addreses []*Address `json:"ips,omitempty"`
}

// Address is an IP address (IPv4/IPv6) with optional open ports and geolocation.
type Address struct {
	Addr     string    `json:"address"`
	Type     string    `json:"type"`
	Ports    []*Port   `json:"ports,omitempty"`
	Location *Location `json:"location,omitempty"`
}

// Port is a network port and the service detected on it.
type Port struct {
	Proto   string `json:"proto,omitempty"`
	Number  int    `json:"num,omitempty"`
	State   string `json:"state,omitempty"`
	Service string `json:"service,omitempty"`
}
