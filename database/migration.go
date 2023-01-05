package database

import "time"

type AuctionOld struct {
	ChannelID       string `gorm:"primaryKey"`
	Bid             float64
	MessageID       string
	EndTime         time.Time
	Winner          string
	GuildID         string
	Item            string
	Host            string
	Currency        string
	IncrementMin    float64
	IncrementMax    float64
	Description     string
	ImageURL        string
	TargetPrice     float64
	Buyout          float64
	CurrencySide    string
	IntegerOnly     bool
	BidHistory      string
	Note            string
	ChannelOverride string
	ChannelLock     bool
	UseCurrency     bool
}