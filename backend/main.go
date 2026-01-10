package main

import (
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v82"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/config"
	"gitlab.com/logan9312/discord-auction-bot/connect"
	"gitlab.com/logan9312/discord-auction-bot/crypto"
	"gitlab.com/logan9312/discord-auction-bot/database"
	"gitlab.com/logan9312/discord-auction-bot/logger"
	"gitlab.com/logan9312/discord-auction-bot/routers"
)

// Environment struct
type Environment struct {
	Environment  string `env:"ENVIRONMENT,required"`
	DiscordToken string `env:"DISCORD_TOKEN,required"`
	Migrate      bool   `env:"MIGRATE"`
	Host         string `env:"DB_HOST"`
	Password     string `env:"DB_PASSWORD"`
	StripeToken  string `env:"STRIPE_TOKEN"`
	LogLevel     string `env:"LOG_LEVEL" envDefault:"info"`
}

func main() {
	environment := Environment{}

	// Load .env file first (before logger init so we can get ENVIRONMENT)
	err := godotenv.Load(".env")
	if err != nil {
		// Can't use logger yet, use stderr
		os.Stderr.WriteString("Warning: Error loading .env file: " + err.Error() + "\n")
	}

	if err := env.Parse(&environment); err != nil {
		os.Stderr.WriteString("FATAL: Failed to load environment variables: " + err.Error() + "\n")
		os.Exit(1)
	}

	// Initialize config from environment variables
	config.Init()

	// Initialize encryption for whitelabel tokens
	if err := crypto.Init(); err != nil {
		os.Stderr.WriteString("Warning: Encryption not initialized - whitelabel tokens will not be encrypted: " + err.Error() + "\n")
	}

	// Initialize logger as early as possible
	if err := logger.Init(environment.Environment, environment.LogLevel); err != nil {
		os.Stderr.WriteString("FATAL: Failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Sugar.Infow("starting bot",
		"environment", environment.Environment,
		"log_level", environment.LogLevel,
	)

	stripe.Key = environment.StripeToken

	// Connects database
	database.DatabaseConnect(environment.Password, environment.Host, environment.Environment)

	// Migrate existing whitelabel tokens to encrypted format
	if err := database.MigrateWhitelabelTokens(); err != nil {
		logger.Sugar.Errorw("failed to migrate whitelabel tokens", "error", err)
	}

	// Connects main bot
	mainSession, err := connect.BotConnect(environment.DiscordToken, environment.Environment)
	if err != nil {
		logger.Sugar.Fatalw("failed to connect main bot", "error", err)
	}

	devData := database.DevSetup{
		BotID: mainSession.State.User.ID,
	}
	result := database.DB.First(&devData)
	if result.Error != nil {
		logger.Sugar.Warnw("failed to fetch dev setup", "error", result.Error)
	}

	err = mainSession.UpdateGameStatus(0, "Bot Version "+devData.Version)
	if err != nil {
		logger.Sugar.Warnw("failed to set game status", "error", err)
	}

	whitelabels := []database.WhiteLabels{}
	result = database.DB.Find(&whitelabels)
	if result.Error != nil {
		logger.Sugar.Errorw("failed to fetch whitelabels", "error", result.Error)
	}

	connect.Timers(mainSession)

	// Start whitelabel bots with delay to avoid rate limiting
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Sugar.Errorw("panic in whitelabel bot goroutine",
					"panic", r,
				)
			}
		}()

		for i, wl := range whitelabels {
			// Add 2 second delay between each whitelabel bot connection
			if i > 0 {
				time.Sleep(2 * time.Second)
			}

			// Decrypt token if encrypted
			botToken := wl.BotToken
			if wl.Encrypted && crypto.IsInitialized() {
				decrypted, err := crypto.Decrypt(wl.BotToken)
				if err != nil {
					logger.Sugar.Errorw("failed to decrypt whitelabel token",
						"bot_id", wl.BotID,
						"error", err,
					)
					continue
				}
				botToken = decrypted
			}

			s, err := connect.BotConnect(botToken, environment.Environment)
			if err != nil {
				logger.Sugar.Errorw("failed to connect whitelabel bot",
					"error", err,
					"index", i,
				)
				continue
			}
			logger.Bot(s.State.User.ID, s.State.User.Username).Info("whitelabel bot connected")

			if s.State.User.ID == config.C.WhitelabelBotID {
				err = s.UpdateGameStatus(0, "Bot Version "+devData.Version)
				if err != nil {
					logger.Sugar.Warnw("failed to set whitelabel status", "error", err)
				}
			}
			connect.Timers(s)
		}
	}()

	go commands.SetRoles(mainSession)

	logger.Sugar.Infow("bot is running",
		"bot_id", mainSession.State.User.ID,
		"bot_name", mainSession.State.User.Username,
	)

	routers.StartServer(mainSession)
}
