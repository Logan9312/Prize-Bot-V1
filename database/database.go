package database

import (
	"fmt"

	"time"

	"gorm.io/gorm"
)

type AuctionSetup struct {
	GuildID        string `gorm:"primaryKey"`
	Category       string
	AlertRole      string
	Currency       string
	LogChannel     string
	Claiming       string
	HostRole       string
	SnipeExtension time.Duration
	SnipeRange     time.Duration
	CurrencySide   string
	IntegerOnly    bool
}

type Auction struct {
	ChannelID      string `gorm:"primaryKey"`
	Bid            float64
	MessageID      string
	EndTime        time.Time
	Winner         string
	GuildID        string
	Item           string
	Host           string
	Currency       string
	IncrementMin   float64
	IncrementMax   float64
	Description    string
	ImageURL       string
	TargetPrice    float64
	Buyout         float64
	SnipeExtension time.Duration
	SnipeRange     time.Duration
	CurrencySide   string
	IntegerOnly    bool
}

type AuctionQueue struct {
	ID           int `gorm:"primaryKey"`
	Bid          float64
	StartTime    time.Time
	EndTime      time.Time
	GuildID      string
	Item         string
	Host         string
	Currency     string
	IncrementMin float64
	IncrementMax float64
	Description  string
	ImageURL     string
	Category     string
	TargetPrice  float64
	Buyout       float64
	CurrencySide string
	IntegerOnly  bool
}

type AuctionClaim struct {
}

type Giveaway struct {
	MessageID    string `gorm:"primaryKey"`
	ChannelID    string
	GuildID      string
	Item         string
	EndTime      time.Time
	Description  string
	Host         string
	Winners      float64
	ImageURL     string
	Entries      string
	Finished     bool
	WinnerOutput string
}

type GiveawaySetup struct {
	GuildID   string `gorm:"primaryKey"`
	HostRole  string
	AlertRole string
	Claiming  string
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

	err := DB.AutoMigrate(AuctionSetup{}, Auction{}, AuctionQueue{}, GiveawaySetup{}, Giveaway{})
	if err != nil {
		fmt.Println(err)
	}

}
