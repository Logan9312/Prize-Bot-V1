package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func LocalDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("Test DB"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	return db
}