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
	defer fmt.Println("Bot has finished attempting to connect to the database!")

	dbuser := "auctionbot"
	port := "3306"
	dbname := "auction"

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, dbuser, dbname, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("Database connected")

	test := User{}

	logan := User{Name: "Logan"}

	err = db.AutoMigrate(User{})
	if err != nil {
		fmt.Println(err)
		return
	}

	db.Create(&logan)

	db.First(&test)

	_, err = fmt.Println(test.Name)
	if err != nil {
		fmt.Println("Error printing name " + err.Error())
		return
	}
}