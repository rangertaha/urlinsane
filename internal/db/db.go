package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// func init() {
// 	var err error

// 	if DB, err = gorm.Open(sqlite.Open("urlinsane.db"), &gorm.Config{}); err != nil {
// 		fmt.Println(err)
// 	}

// 	// Migrate the schema
// 	DB.AutoMigrate(
// 		// Languages
// 		&Keyboard{},
// 		&Language{},
// 		&Word{},
// 		&Char{},

// 		// Domains
// 		&Prefix{},
// 		&Suffix{},
// 		&Contact{},
// 		&Domain{},
// 		&Whois{},
// 		&Dns{},

// 		// Networking
// 		&Server{},
// 		&Service{},

// 		// Geography
// 		&Place{},
// 		&Location{},

// 		// Files
// 		&Page{},
// 		&File{},
// 		&Image{},
// 	)
// }


func Config(filepath string) {
	var err error

	if DB, err = gorm.Open(sqlite.Open(filepath), &gorm.Config{}); err != nil {
		fmt.Println(err)
	}

	// Migrate the schema
	DB.AutoMigrate(
		// Languages
		&Keyboard{},
		&Language{},
		&Word{},
		&Char{},

		// Domains
		&Prefix{},
		&Suffix{},
		&Contact{},
		&Domain{},
		&Whois{},
		&Dns{},

		// Networking
		&Server{},
		&Service{},

		// Geography
		&Place{},
		&Location{},

		// Files
		&Page{},
		&File{},
		&Image{},
	)
}
