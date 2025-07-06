package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProdDB() *gorm.DB {
	databaseURL := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	return db
}