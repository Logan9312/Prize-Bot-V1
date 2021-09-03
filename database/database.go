package database

import (
	"fmt"

	"gorm.io/gorm"
)

type Auction struct {
	gorm.Model
	Bid       string
	MessageID string
	ChannelID string
}

type GuildInfo struct {
	GuildID         string `gorm:"primaryKey"`
	AuctionCategory string
}

var db *gorm.DB

func DatabaseConnect(password, host, env string) *gorm.DB {
	fmt.Println("Connecting to Database...")
	defer fmt.Println("Bot has finished attempting to connect to the database!")

	if env == "prod" {
		db = ProdDB(password, host)
	} else if env == "local" {
		db = LocalDB()
	}

	err := db.AutoMigrate(GuildInfo{}, Auction{})
	if err != nil {
		fmt.Println(err)
	}

	return db

}

func GuildSetup(guildinfo GuildInfo) {

}
