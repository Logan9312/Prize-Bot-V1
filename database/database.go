package database

import (
	"fmt"

	"gorm.io/gorm"
)

type Auction struct {
	gorm.Model
	Bid     string
	MessageID string
	ChannelID string
}

func DatabaseConnect(password, host, env string) {
	fmt.Println("Connecting to Database...")
	defer fmt.Println("Bot has finished attempting to connect to the database!")

	var db *gorm.DB
	if env == "prod" {
		db = ProdDB(password, host)
	} else if env == "local" {
		db = LocalDB()
	}
	
	fmt.Println("Database connected")

	err := db.AutoMigrate(Auction{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
