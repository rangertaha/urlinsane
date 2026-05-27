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

// Source is a place to look up an entity: a username platform, a package
// registry, or an email provider. Template is the human display URL (%s replaced
// by the entity value; for email providers, just the provider domain). Check is
// the URL used to test existence (often an API endpoint that returns a clean
// 404 for missing entries); when empty, Template is used.
type Source struct {
	ID       uint
	Type     string `gorm:"index"` // username | package | email
	Code     string
	Template string
	CheckURL string // existence-check URL (API); falls back to Template
}

// Package is a known package name in a code registry — a corpus of common
// supply-chain typosquatting targets.
type Package struct {
	ID       uint
	Name     string `gorm:"index"`
	Registry string
	Rank     int64
}
