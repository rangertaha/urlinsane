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
		// &Scan{},
		&Contact{},
		&Domain{},
		&WhoisRecord{},
		&DnsRecord{},

		// Networking
		&Address{},
		&Service{},
		&Port{},
		// &Device{},


		// Geography
		&Place{},
		&Location{},

		// Files
		// &Page{},
		// &File{},
		// &Image{},
	)
}

type Scan struct {
	gorm.Model
	DomainID uint
	Domain   *Domain
	Results  []*Domain
}
