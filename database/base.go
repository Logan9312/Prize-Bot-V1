package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

func DatabaseConnect(host, password string) {
	fmt.Println("Connecting to Database...")

	dbuser := "auctionbot"
	port := "3306"
	dbname := "auction"

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, dbuser, dbname, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Database Connected!")

	user := User{Name: "Logan"}

	db.AutoMigrate(&User{})
	db.Create(user)
	test := User{}

	db.First(&test, 1)
	fmt.Println(test)
}