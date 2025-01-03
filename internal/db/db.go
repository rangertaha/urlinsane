package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Meta *gorm.DB

func Config(filepath string) {
	var err error

	if DB, err = gorm.Open(sqlite.Open(filepath), &gorm.Config{CreateBatchSize: 1000}); err != nil {
		fmt.Println(err)
	}

	// Migrate the schema
	DB.AutoMigrate(
		// Domains
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
