package database

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Auction struct {
	gorm.Model
	Bid string
	Message discordgo.Message
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

	err = db.AutoMigrate(Auction{})
	if err != nil {
		fmt.Println(err)
		return
	}
}