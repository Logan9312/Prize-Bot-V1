package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DatabaseConnect (host, password string){
	user := "auctionbot"
	port := "3306"
	dbname := "bot-prod-auction"

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbname, password)
	db, err := gorm.Open(postgres.Open(dsn), )
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(db)

}