package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProdDB() *gorm.DB {
	databaseURL := os.Getenv("DATABASE_URL")
	
	if databaseURL == "" {
		fmt.Println("ERROR: DATABASE_URL environment variable is not set!")
		return nil
	}

	fmt.Printf("Connecting to database with URL: %s\n", databaseURL[:20] + "...") // Only show first 20 chars for security

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return nil
	}

	return db
}