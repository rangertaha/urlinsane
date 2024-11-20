package db

import (
	"fmt"

	// "gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)


var DB *gorm.DB

func init() {
	var err error

	if DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{}); err != nil {
		fmt.Println(err)
	}

}



// Domain
// 	Name
// 	IDN
// 	DNS []DnsRecord
// 	IPs []IP
// 	Subdomains []Domain
	
// DnsRecord

// WhoisRecord


// IP 
// 	Address
// 	Location
// 	Ports []Port

// Port
// 	Number
// 	State
// 	Service


// Location

