package database

import (
	"fmt"

	"time"

	"gorm.io/gorm"
)

type DevSetup struct {
	BotID   string `gorm:"primaryKey"`
	Version string
}

type WhiteLabels struct {
	BotID    string `gorm:"primaryKey;autoIncrement:false"`
	UserID   string `gorm:"primaryKey;autoIncrement:false"`
	BotToken string
}

type AuctionSetup struct {
	GuildID         string `gorm:"primaryKey"`
	Category        string
	AlertRole       string
	Currency        string
	LogChannel      string
	HostRole        string
	SnipeExtension  time.Duration
	SnipeRange      time.Duration
	CurrencySide    string
	IntegerOnly     bool
	ChannelOverride string
	ChannelLock     bool
	ChannelPrefix   string
	UseCurrency     bool
}

type Auction struct {
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

type AuctionQueue struct {
	ID              int `gorm:"primaryKey"`
	ChannelID       string
	Bid             float64
	StartTime       time.Time
	EndTime         time.Time
	GuildID         string
	Item            string
	Host            string
	Currency        string
	IncrementMin    float64
	IncrementMax    float64
	Description     string
	ImageURL        string
	Category        string
	TargetPrice     float64
	Buyout          float64
	CurrencySide    string
	IntegerOnly     bool
	SnipeExtension  time.Duration
	SnipeRange      time.Duration
	AlertRole       string
	Note            string
	ChannelOverride string
	ChannelLock     bool
	UseCurrency     bool
}

// ClaimSetup FromMake sure to remove LogChannel and ClaimMessage from auction log
type ClaimSetup struct {
	GuildID         string `gorm:"primaryKey"`
	Category        string
	StaffRole       string
	Instructions    string
	LogChannel      string
	Expiration      string
	DisableClaiming bool
	ChannelPrefix   string
}

type Claim struct {
	MessageID   string `gorm:"primaryKey"`
	ChannelID   string
	GuildID     string
	Item        string
	Type        string
	Winner      string
	Cost        float64
	Host        string
	BidHistory  string
	Note        string
	ImageURL    string
	TicketID    string
	Description string
	UseCurrency bool
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
	Finished    bool
}

type GiveawaySetup struct {
	GuildID    string `gorm:"primaryKey"`
	HostRole   string
	AlertRole  string
	LogChannel string
}

type ShopSetup struct {
	GuildID    string `gorm:"primaryKey"`
	HostRole   string
	AlertRole  string
	LogChannel string
}

type CurrencySetup struct {
	GuildID  string `gorm:"primaryKey"`
	Currency string
	Side     string
}

type UserProfile struct {
	UserID  string `gorm:"primaryKey;autoIncrement:false"`
	GuildID string `gorm:"primaryKey;autoIncrement:false"`
	Balance float64
}

type Quest struct {
	MessageID string `gorm:"primaryKey;autoIncrement:false"`
}

type Errors struct {
	ErrorID string `gorm:"primaryKey"`
	UserID  string `gorm:"primaryKey;autoIncrement:false"`
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

	err := DB.AutoMigrate(AuctionSetup{}, Auction{}, AuctionQueue{}, GiveawaySetup{}, Giveaway{}, ClaimSetup{}, CurrencySetup{}, Claim{}, DevSetup{}, UserProfile{}, ShopSetup{}, WhiteLabels{})
	if err != nil {
		fmt.Println(err)
	}

}
