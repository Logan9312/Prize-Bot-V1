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
	Item      string
	Host      string
	Currency  string
	MinBid    float64
	MaxBid    float64
}

type GuildInfo struct {
	GuildID         string `gorm:"primaryKey"`
	AuctionCategory string
	AuctionRole     string
	Currency        string
	LogChannel      string
	Claiming        string
	AuctionHostRole string
}

type UserProfile struct {
	UserID string `gorm:"primaryKey"`
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
