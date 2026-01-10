package database

import (
	"time"

	"gitlab.com/logan9312/discord-auction-bot/logger"
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
	SnipeLimit      time.Duration
	SnipeCap        time.Duration
	CurrencySide    string
	IntegerOnly     bool
	ChannelOverride string
	ChannelLock     bool
	ChannelPrefix   string
	UseCurrency     bool
}

type Auction struct {
	ChannelID           string `gorm:"primaryKey"`
	Bid                 float64
	MessageID           string `gorm:"index"`
	StartTime           time.Time
	EndTime             time.Time `gorm:"index"`
	Winner              string
	GuildID             string `gorm:"index"`
	Item                string
	Host                string
	Currency            string
	IncrementMin        float64
	IncrementMax        float64
	Description         string
	ImageURL            string
	TargetPrice         float64
	Buyout              float64
	CurrencySide        string
	IntegerOnly         bool
	BidHistory          string
	Note                string
	ChannelOverride     string
	ChannelLock         bool
	UseCurrency         bool
	TotalSnipeExtension time.Duration
}

type AuctionQueue struct {
	ID              int `gorm:"primaryKey"`
	ChannelID       string
	Bid             float64
	StartTime       time.Time
	EndTime         time.Time
	GuildID         string `gorm:"index"`
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
	SnipeLimit      time.Duration
	SnipeCap        time.Duration
	AlertRole       string
	Note            string
	ChannelOverride string
	ChannelLock     bool
	UseCurrency     bool
	ChannelPrefix   string
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
	GuildID     string `gorm:"index"`
	Item        string
	Type        string
	Winner      string `gorm:"index"`
	Cost        float64
	Host        string
	BidHistory  string
	Note        string
	ImageURL    string
	TicketID    string
	Description string
	UseCurrency bool
	Status      string `gorm:"index;default:pending"` // "pending", "claimed", "cancelled"
}

// Claim status constants
const (
	ClaimStatusPending   = "pending"
	ClaimStatusClaimed   = "claimed"
	ClaimStatusCancelled = "cancelled"
)

type Giveaway struct {
	MessageID   string `gorm:"primaryKey"`
	ChannelID   string
	GuildID     string `gorm:"index"`
	Item        string
	EndTime     time.Time `gorm:"index"`
	Description string
	Host        string
	Winners     int // Changed from float64 - winner count should be an integer
	ImageURL    string
	Finished    bool `gorm:"index"`
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

// Subscription caches Stripe subscription data locally
// This follows the recommendation to avoid querying Stripe API on every premium check
type Subscription struct {
	ID                 string    `gorm:"primaryKey"`              // Stripe subscription ID
	CustomerID         string    `gorm:"index"`                   // Stripe customer ID
	DiscordUserID      string    `gorm:"index"`                   // Discord user ID from metadata
	GuildID            string    `gorm:"index"`                   // Linked guild ID from metadata (optional)
	Status             string    `gorm:"index"`                   // active, canceled, past_due, etc.
	PriceID            string                                     // Stripe price ID
	CurrentPeriodStart time.Time                                  // Subscription period start
	CurrentPeriodEnd   time.Time                                  // Subscription period end
	CancelAtPeriodEnd  bool                                       // Whether subscription cancels at period end
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`          // Last sync time
}

var DB *gorm.DB

func DatabaseConnect(password, host, env string) {
	log := logger.Database("connect")
	log.Infow("connecting to database", "environment", env)

	switch env {
	case "prod", "dev":
		DB = ProdDB()
	case "local":
		DB = LocalDB()
	}

	if DB == nil {
		log.Fatal("failed to initialize database connection")
		return
	}

	log.Info("database connection established")

	err := DB.AutoMigrate(AuctionSetup{}, Auction{}, AuctionQueue{}, GiveawaySetup{}, Giveaway{}, ClaimSetup{}, CurrencySetup{}, Claim{}, DevSetup{}, UserProfile{}, ShopSetup{}, WhiteLabels{}, Subscription{})
	if err != nil {
		log.Errorw("database migration failed", "error", err)
	} else {
		log.Info("database migration completed")
	}
}
