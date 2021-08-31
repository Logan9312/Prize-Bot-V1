package database

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type Auction struct {
	gorm.Model
	Bid     string
	Message discordgo.Message
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
