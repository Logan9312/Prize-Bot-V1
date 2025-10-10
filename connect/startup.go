package connect

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	c "gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/database"
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
		return s, fmt.Errorf("Discordgo.New Error: %w", err)
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
		return s, fmt.Errorf("Failed to open a websocket connection with discord. Likely due to an invalid token. %w", err)
	}

	fmt.Println(s.State.User.Username, " Starting Up...")

	// Wait for READY event with timeout
	select {
	case <-ready:
		fmt.Println("Bot received READY event")
	case <-time.After(30 * time.Second):
		return s, fmt.Errorf("timeout waiting for READY event")
	}

	// Give Discord more time to send all GUILD_CREATE events
	// This prevents rate limiting by not trying to access guilds before they're fully loaded
	time.Sleep(2 * time.Second)

	//Builds local commands
	if environment == "local" {
		s.LogLevel = discordgo.LogInformational
		guildCount := len(s.State.Guilds)
		fmt.Printf("Registering commands to %d guilds...\n", guildCount)

		for i, v := range s.State.Guilds {
			// Add a delay every 5 guilds to avoid rate limiting
			if i > 0 && i%5 == 0 {
				fmt.Printf("Processed %d/%d guilds, pausing to avoid rate limits...\n", i, guildCount)
				time.Sleep(2 * time.Second)
			}

			_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, BotCommands.Local)
			if err != nil {
				fmt.Printf("Warning: Failed to add commands to guild %s: %v\n", v.Name, err)
				// Don't return error, just log and continue with other guilds
				continue
			}
			fmt.Printf("Commands added to guild: %s (%d/%d)\n", v.Name, i+1, guildCount)
		}
	}

	//Builds prod commands
	if environment == "prod" {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", BotCommands.Prod)
		if err != nil {
			return s, fmt.Errorf("Bulk Overwrite Prod Command Error: %w", err)
		}

		//Builds dev commands

		_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "915767892467920967", BotCommands.Dev)
		if err != nil {
			return s, fmt.Errorf("Bulk Overwrite Dev Command Error: %w", err)
		}

	}

	fmt.Println(s.State.User.Username + " bot startup complete!")

	return s, nil
}

func Timers(s *discordgo.Session) {

	Auctions := []map[string]interface{}{}
	AuctionQueue := []map[string]interface{}{}
	Giveaways := []map[string]interface{}{}
	fmt.Println("Beginning Startup Timers")

	database.DB.Model([]database.Auction{}).Find(&Auctions)
	fmt.Printf("Found %d active auctions to process\n", len(Auctions))
	for _, v := range Auctions {
		go AuctionEndHandler(v, s)
	}

	//TODO Fix this with whitelabels
	database.DB.Model([]database.AuctionQueue{}).Find(&AuctionQueue)
	fmt.Printf("Found %d queued auctions to process\n", len(AuctionQueue))
	for _, v := range AuctionQueue {
		go AuctionStartHandler(v, s)
	}

	database.DB.Model([]database.Giveaway{}).Find(&Giveaways)
	fmt.Printf("Found %d active giveaways to process\n", len(Giveaways))
	for _, v := range Giveaways {
		go GiveawayEndHandler(v, s)
	}
}

func AuctionEndHandler(v map[string]interface{}, s *discordgo.Session) {
	fmt.Println("Auction Timer Re-Started: ", v["item"], "GuildID: ", v["guild_id"], "ImageURL", v["image_url"], "Host", v["host"], "End Time", v["end_time"].(time.Time).String())

	// Check if the auction has already ended or is still active
	endTime, ok := v["end_time"].(time.Time)
	if !ok {
		fmt.Println("Error: Invalid end_time for auction", v["channel_id"])
		return
	}

	// If auction hasn't ended yet, wait until end time before making any Discord API calls
	if endTime.After(time.Now()) {
		timeUntilEnd := time.Until(endTime)
		fmt.Printf("Auction '%s' will end in %s\n", v["item"], timeUntilEnd)
		time.Sleep(timeUntilEnd)
	}

	// Now the auction has ended, proceed with ending logic
	c.AuctionEnd(s, v["channel_id"].(string), v["guild_id"].(string))
}

func AuctionStartHandler(v map[string]interface{}, s *discordgo.Session) {
	fmt.Println("Auction Re-Queued: ", v["item"], "GuildID: ", v["guild_id"], "ImageURL", v["image_url"], "Host", v["host"], "Start Time", v["start_time"].(time.Time).String())
	if v["start_time"].(time.Time).Before(time.Now()) {
		c.AuctionStart(s, v)
	} else {
		time.Sleep(time.Until(v["start_time"].(time.Time)))
		c.AuctionStart(s, v)
	}
}

func GiveawayEndHandler(v map[string]interface{}, s *discordgo.Session) {
	fmt.Println("Giveaway Timer Re-Started: ", v["item"], "GuildID: ", v["guild_id"], "ImageURL", v["image_url"], "Host", v["host"], "End Time", v["end_time"].(time.Time).String())

	endTime, ok := v["end_time"].(time.Time)
	if !ok {
		fmt.Println("Error: Invalid end_time for giveaway", v["message_id"])
		return
	}

	if endTime.Before(time.Now()) {
		if v["finished"] == true {
			time.Sleep(time.Until(endTime.Add(24 * time.Hour)))
			database.DB.Delete(database.Giveaway{}, v["message_id"].(string))
		} else {
			c.GiveawayEnd(s, v["message_id"].(string))
		}
	} else {
		time.Sleep(time.Until(endTime))
		c.GiveawayEnd(s, v["message_id"].(string))
	}
}
