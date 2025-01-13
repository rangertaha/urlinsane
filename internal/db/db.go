package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Config(filepath string) {
	var err error

	if DB, err = gorm.Open(sqlite.Open(filepath), &gorm.Config{CreateBatchSize: 1000}); err != nil {
		fmt.Println(err)
	}

	// Migrate the schema
	DB.AutoMigrate(
		// Domains
		&Scan{},
		&Contact{},
		&Domain{},
		&Whois{},
		&Dns{},

		// Networking
		&Address{},
		&Port{},

		// Geography
		&Location{},

		// Files
		// &Page{},
		// &File{},
		// &Image{},
	)
}

type Scan struct {
	gorm.Model
	Query   string
	Results []*Domain `gorm:"many2many:results;"`
}
