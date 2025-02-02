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
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Config(filepath string) {
	var err error

	if DB, err = gorm.Open(sqlite.Open(filepath), &gorm.Config{CreateBatchSize: 10000}); err != nil {
		fmt.Println(err)
	}

	// Migrate the schema
	DB.AutoMigrate(
		&Keyboard{},

		// Language
		&Language{},
		&Word{},
		&Char{},
		&Sym{},

		// Domain
		&Prefix{},
		&Suffix{},
		&Domain{},

		// NLP
		// Topic
		// 
	)
}
