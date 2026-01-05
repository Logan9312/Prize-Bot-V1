package connect

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	c "gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/config"
	"gitlab.com/logan9312/discord-auction-bot/database"
	"gitlab.com/logan9312/discord-auction-bot/logger"
)

type slashCommands struct {
	Local, Prod, Dev []*discordgo.ApplicationCommand
}

var BotCommands = slashCommands{
	Local: []*discordgo.ApplicationCommand{
		&HelpCommand,
		&c.ShopCommand,
		&c.QuestCommand,
		//&QuestContextMenu,
	},
	Prod: []*discordgo.ApplicationCommand{
		&c.ProfileCommand,
		&c.AuctionCommand,
		&c.BidCommand,
		&c.GiveawayCommand,
		&c.PrivacyCommand,
		&c.ClaimCommand,
		&c.PremiumCommand,
		&c.SettingsCommand,
		&c.CurrencyCommand,
		&WhitelabelCommand,
	},
	Dev: []*discordgo.ApplicationCommand{
		&c.DevCommand,
		&c.StatsCommand,
	},
}

func BotConnect(token, environment string) (*discordgo.Session, error) {
	BotCommands.Local = append(BotCommands.Local, BotCommands.Prod...)
	BotCommands.Local = append(BotCommands.Local, BotCommands.Dev...)

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return s, err
	}

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages

	// Create a channel to wait for the READY event
	ready := make(chan bool)
	readyHandler := func(s *discordgo.Session, r *discordgo.Ready) {
		select {
		case ready <- true:
		default:
		}
	}
	removeHandler := s.AddHandlerOnce(readyHandler)
	defer removeHandler()

	RegisterHandlers(s)

	err = s.Open()
	if err != nil {
		return s, err
	}

	log := logger.Bot(s.State.User.ID, s.State.User.Username)
	log.Info("bot starting up")

	// Wait for READY event with timeout
	select {
	case <-ready:
		log.Info("bot received READY event")
	case <-time.After(30 * time.Second):
		return s, fmt.Errorf("timed out waiting for READY event after 30 seconds")
	}

	// Give Discord more time to send all GUILD_CREATE events
	// This prevents rate limiting by not trying to access guilds before they're fully loaded
	time.Sleep(2 * time.Second)

	// Builds local commands
	if environment == "local" {
		s.LogLevel = discordgo.LogInformational
		guildCount := len(s.State.Guilds)
		log.Infow("registering commands to guilds", "guild_count", guildCount)

		for i, v := range s.State.Guilds {
			// Add a delay every 5 guilds to avoid rate limiting
			if i > 0 && i%5 == 0 {
				log.Infow("pausing to avoid rate limits", "processed", i, "total", guildCount)
				time.Sleep(2 * time.Second)
			}

			_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, BotCommands.Local)
			if err != nil {
				log.Warnw("failed to add commands to guild",
					"guild_name", v.Name,
					"guild_id", v.ID,
					"error", err,
				)
				continue
			}
			log.Debugw("commands added to guild",
				"guild_name", v.Name,
				"progress", i+1,
				"total", guildCount,
			)
		}
	}

	// Builds prod commands
	if environment == "prod" {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", BotCommands.Prod)
		if err != nil {
			return s, err
		}

		// Builds dev commands
		_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, config.C.DevCommandsGuild, BotCommands.Dev)
		if err != nil {
			return s, err
		}
	}

	log.Info("bot startup complete")
	return s, nil
}

func Timers(s *discordgo.Session) {
	log := logger.Bot(s.State.User.ID, s.State.User.Username)
	log.Info("initializing startup timers")

	Auctions := []map[string]interface{}{}
	AuctionQueue := []map[string]interface{}{}
	Giveaways := []map[string]interface{}{}

	// Load active auctions with error checking
	result := database.DB.Model([]database.Auction{}).Find(&Auctions)
	if result.Error != nil {
		log.Errorw("failed to load active auctions", "error", result.Error)
	} else {
		log.Infow("loaded active auctions", "count", len(Auctions))
		for _, v := range Auctions {
			go AuctionEndHandler(v, s)
		}
	}

	// Load queued auctions with error checking
	result = database.DB.Model([]database.AuctionQueue{}).Find(&AuctionQueue)
	if result.Error != nil {
		log.Errorw("failed to load queued auctions", "error", result.Error)
	} else {
		log.Infow("loaded queued auctions", "count", len(AuctionQueue))
		for _, v := range AuctionQueue {
			go AuctionStartHandler(v, s)
		}
	}

	// Load active giveaways with error checking
	result = database.DB.Model([]database.Giveaway{}).Find(&Giveaways)
	if result.Error != nil {
		log.Errorw("failed to load active giveaways", "error", result.Error)
	} else {
		log.Infow("loaded active giveaways", "count", len(Giveaways))
		for _, v := range Giveaways {
			go GiveawayEndHandler(v, s)
		}
	}

	log.Info("startup timers initialized")
}

func AuctionEndHandler(v map[string]interface{}, s *discordgo.Session) {
	defer func() {
		if r := recover(); r != nil {
			logger.Sugar.Errorw("panic in AuctionEndHandler",
				"panic", r,
			)
		}
	}()

	// Safe type assertions for required fields
	channelID, ok := v["channel_id"].(string)
	if !ok {
		logger.Sugar.Error("invalid or missing channel_id in auction data")
		return
	}

	guildID, ok := v["guild_id"].(string)
	if !ok {
		logger.Sugar.Error("invalid or missing guild_id in auction data")
		return
	}

	item := ""
	if v["item"] != nil {
		if itemStr, ok := v["item"].(string); ok {
			item = itemStr
		}
	}

	log := logger.Auction(channelID, guildID, item)
	log.Infow("auction timer restarted",
		"host", v["host"],
		"end_time", v["end_time"],
	)

	// Check if the auction has already ended or is still active
	endTime, ok := v["end_time"].(time.Time)
	if !ok {
		log.Error("invalid end_time for auction")
		return
	}

	// If auction hasn't ended yet, wait until end time before making any Discord API calls
	if endTime.After(time.Now()) {
		timeUntilEnd := time.Until(endTime)
		log.Infow("auction scheduled to end", "time_until_end", timeUntilEnd.String())
		time.Sleep(timeUntilEnd)
	}

	// Now the auction has ended, proceed with ending logic
	c.AuctionEnd(s, channelID, guildID)
}

func AuctionStartHandler(v map[string]interface{}, s *discordgo.Session) {
	defer func() {
		if r := recover(); r != nil {
			logger.Sugar.Errorw("panic in AuctionStartHandler",
				"panic", r,
			)
		}
	}()

	// Safe type assertion for guild_id
	guildID, ok := v["guild_id"].(string)
	if !ok {
		logger.Sugar.Error("invalid or missing guild_id in auction queue data")
		return
	}

	item := ""
	if v["item"] != nil {
		if itemStr, ok := v["item"].(string); ok {
			item = itemStr
		}
	}

	log := logger.Timer("auction_queue", guildID)
	log.Infow("auction re-queued",
		"item", item,
		"host", v["host"],
		"start_time", v["start_time"],
	)

	// Safe type assertion for start_time
	startTime, ok := v["start_time"].(time.Time)
	if !ok {
		log.Error("invalid or missing start_time in auction queue data")
		return
	}

	if startTime.Before(time.Now()) {
		c.AuctionStart(s, v)
	} else {
		timeUntilStart := time.Until(startTime)
		log.Infow("auction scheduled to start", "time_until_start", timeUntilStart.String())
		time.Sleep(timeUntilStart)
		c.AuctionStart(s, v)
	}
}

func GiveawayEndHandler(v map[string]interface{}, s *discordgo.Session) {
	defer func() {
		if r := recover(); r != nil {
			logger.Sugar.Errorw("panic in GiveawayEndHandler",
				"panic", r,
			)
		}
	}()

	// Safe type assertions for required fields
	messageID, ok := v["message_id"].(string)
	if !ok {
		logger.Sugar.Error("invalid or missing message_id in giveaway data")
		return
	}

	guildID, ok := v["guild_id"].(string)
	if !ok {
		logger.Sugar.Error("invalid or missing guild_id in giveaway data")
		return
	}

	item := ""
	if v["item"] != nil {
		if itemStr, ok := v["item"].(string); ok {
			item = itemStr
		}
	}

	log := logger.Giveaway(messageID, guildID, item)
	log.Infow("giveaway timer restarted",
		"host", v["host"],
		"end_time", v["end_time"],
	)

	endTime, ok := v["end_time"].(time.Time)
	if !ok {
		log.Error("invalid end_time for giveaway")
		return
	}

	if endTime.Before(time.Now()) {
		if v["finished"] == true {
			cleanupTime := time.Until(endTime.Add(24 * time.Hour))
			log.Infow("giveaway finished, scheduling cleanup", "cleanup_in", cleanupTime.String())
			time.Sleep(cleanupTime)
			result := database.DB.Delete(database.Giveaway{}, messageID)
			if result.Error != nil {
				log.Errorw("failed to delete finished giveaway", "error", result.Error)
			}
		} else {
			c.GiveawayEnd(s, messageID)
		}
	} else {
		timeUntilEnd := time.Until(endTime)
		log.Infow("giveaway scheduled to end", "time_until_end", timeUntilEnd.String())
		time.Sleep(timeUntilEnd)
		c.GiveawayEnd(s, messageID)
	}
}
