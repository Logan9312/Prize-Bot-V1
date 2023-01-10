package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

)

func LocalDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("/tmp/test.db"), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err)
	}

	return db
}
