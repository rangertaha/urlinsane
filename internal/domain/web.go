// Copyright (C) 2024 Rangertaha
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
package domain

type URL struct {
	Address  string   `json:"ip,omitempty"`
}

// IPv4       []IP   `json:"ipv4,omitempty"`
// IPv6       []IP   `json:"ipv6,omitempty"`
// Banner     Banner `json:"response,omitempty"`
// Screenshot string `json:"screenshot,omitempty"`
// Html       string `json:"html,omitempty"`
// Ssdeep     string `json:"ssdeep,omitempty"`
