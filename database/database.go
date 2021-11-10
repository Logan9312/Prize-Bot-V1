package database

import (
	"fmt"

	"time"

	"gorm.io/gorm"
)

type Auction struct {
	ChannelID   string `gorm:"primaryKey"`
	Bid         float64
	MessageID   string
	EndTime     time.Time
	Winner      string
	GuildID     string
	Item        string
	Host        string
	Currency    string
	MinBid      float64
	MaxBid      float64
	Description string
	ImageURL    string
	TargetPrice float64
	Buyout      float64
}

type Giveaway struct {
	MessageID   string `gorm:"primaryKey"`
	ChannelID   string
	GuildID     string
	Item        string
	EndTime     time.Time
	Description string
	Host        string
	Winners     float64
	ImageURL    string
	Entries     string
	Finished    bool
}

type AuctionQueue struct {
	ID          uint `gorm:"primaryKey"`
	Bid         float64
	StartTime   time.Time
	EndTime     time.Time
	GuildID     string
	Item        string
	Host        string
	Currency    string
	MinBid      float64
	MaxBid      float64
	Description string
	ImageURL    string
	Category    string
	TargetPrice float64
	Buyout      float64
}

type GuildInfo struct {
	GuildID          string `gorm:"primaryKey"`
	AuctionCategory  string
	AuctionRole      string
	Currency         string
	LogChannel       string
	Claiming         string
	AuctionHostRole  string
	GiveawayHostRole string
	AntiSnipe        bool
	SnipeExtension   time.Duration
	SnipeRange       time.Duration
	GiveawayRole     string
	GiveawayClaiming string
}

type UserProfile struct {
	UserID  string `gorm:"primaryKey;autoIncrement:false"`
	GuildID string `gorm:"primaryKey;autoIncrement:false"`
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

	err := DB.AutoMigrate(GuildInfo{}, Auction{}, AuctionQueue{}, Giveaway{})
	if err != nil {
		fmt.Println(err)
	}
}
