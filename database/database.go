package database

import (
	"fmt"

	"time"

	"gorm.io/gorm"
)

type Auction struct {
	ChannelID string `gorm:"primaryKey"`
	Bid       float64
	MessageID string
	EndTime   time.Time
	Winner    string
	GuildID   string
	Item	  string
	
}

type GuildInfo struct {
	GuildID         string `gorm:"primaryKey"`
	AuctionCategory string
	Currency        string
	LogChannel      string
}

var DB *gorm.DB

func DatabaseConnect(password, host, env string) {
	fmt.Println("Connecting to Database...")
	defer fmt.Println("Bot has finished attempting to connect to the database!")

	if env == "prod" {
		DB = ProdDB(password, host)
	} else if env == "local" {
		DB = LocalDB()
	}

	err := DB.AutoMigrate(GuildInfo{}, Auction{})
	if err != nil {
		fmt.Println(err)
	}

}
