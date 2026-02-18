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

import (
	"gorm.io/gorm"
)

type Vector struct {
	gorm.Model
}

type Topic struct {
	gorm.Model
	Name      string
	TermFreqs []*Word `gorm:"many2many:termfreqs;"`
}

type TermFreq struct {
	TopicID uint `gorm:"primaryKey"`
	WordID  uint `gorm:"primaryKey"`
	Freq    int64
}

type NGram struct {
	gorm.Model
	Tokens []uint `gorm:"serializer:json"`
}
