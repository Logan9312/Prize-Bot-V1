package connect

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/database"
)

type slashCommands struct {
	local, prod, dev []*discordgo.ApplicationCommand
}

var botCommands = slashCommands{
	local: []*discordgo.ApplicationCommand{
		&commands.HelpCommand,
		&commands.ProfileCommand,
		&commands.AuctionCommand,
		&commands.BidCommand,
		&commands.GiveawayCommand,
		&commands.PrivacyCommand,
		&commands.DevCommand,
		&commands.ClaimCommand,
		&commands.ShopCommand,
		&commands.PremiumCommand,
		&commands.SettingsCommand,
		&commands.QuestCommand,
		//&commands.QuestContextMenu,
		&commands.CurrencyCommand,
	},
	prod: []*discordgo.ApplicationCommand{
		&commands.ProfileCommand,
		&commands.AuctionCommand,
		&commands.BidCommand,
		&commands.GiveawayCommand,
		&commands.PrivacyCommand,
		&commands.ClaimCommand,
		&commands.PremiumCommand,
		&commands.SettingsCommand,
		&commands.CurrencyCommand,
	},
	dev: []*discordgo.ApplicationCommand{
		&commands.DevCommand,
		&commands.StatsCommand,
	},
}

func BotConnect(token, environment, botName string) {

	fmt.Println(botName + " Starting Up...")

	var err error

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("discordgo.New error:" + err.Error())
	}

	s.AddHandler(ReadyHandler)
	s.AddHandler(InteractionHandler)
	s.AddHandler(MessageHandler)
	s.AddHandler(GuildMemberChunkHandler)
	s.AddHandler(GuildCreateHandler)

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages
	err = s.Open()

	if err != nil {
		fmt.Println("s.Open error: " + err.Error())
		return
	}

	//Builds local commands
	if environment == "local" {
		s.LogLevel = discordgo.LogInformational
		for _, v := range s.State.Guilds {
			_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, botCommands.local)
			fmt.Println("Commands added to guild: " + v.Name)
			if err != nil {
				log.Fatal("Bulk Overwrite Error:", err)
			}
		}
	}

	//Builds prod commands
	if environment == "prod" {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", botCommands.prod)
		if err != nil {
			fmt.Println("Bulk Overwrite Error:", err)
		}

		//Builds dev commands
		_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "915767892467920967", botCommands.dev)
		if err != nil {
			fmt.Println("Bulk Overwrite Error:", err)
		}
	}

	go commands.SetRoles(s)

	Timers(s)

	devData := database.DevSetup{
		Bot: botName,
	}
	result := database.DB.First(&devData)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	err = s.UpdateGameStatus(0, "Bot Version "+devData.Version)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(botName + " bot startup complete!")
}

func Timers(s *discordgo.Session) {

	Auctions := []map[string]interface{}{}
	AuctionQueue := []map[string]interface{}{}
	Giveaways := []map[string]interface{}{}
	Claims := []map[string]interface{}{}

	fmt.Println("Beginning Startup Timers")

	database.DB.Model([]database.Auction{}).Find(&Auctions)
	for _, v := range Auctions {
		go AuctionEndTimer(v, s)
	}

	database.DB.Model([]database.AuctionQueue{}).Find(&AuctionQueue)
	for _, v := range AuctionQueue {
		go AuctionStartTimer(v, s)
	}

	database.DB.Model([]database.Giveaway{}).Find(&Giveaways)
	for _, v := range Giveaways {
		go GiveawayEndTimer(v, s)
	}

	database.DB.Model([]database.Claim{}).Find(&Claims)
	for _, v := range Claims {
		if v["item"] == nil {
			database.DB.Delete(database.Claim{
				MessageID: v["message_id"].(string),
			})
		}
	}
}

func AuctionEndTimer(v map[string]interface{}, s *discordgo.Session) {
	fmt.Println("Auction Timer Re-Started: ", v["item"], "GuildID: ", v["guild_id"], "ImageURL", v["image_url"], "Host", v["host"], "End Time", v["end_time"].(time.Time).String())
	commands.AuctionEnd(s, v["channel_id"].(string), v["guild_id"].(string))
}

func AuctionStartTimer(v map[string]interface{}, s *discordgo.Session) {
	fmt.Println("Auction Re-Queued: ", v["item"], "GuildID: ", v["guild_id"], "ImageURL", v["image_url"], "Host", v["host"], "Start Time", v["start_time"].(time.Time).String())
	if v["start_time"].(time.Time).Before(time.Now()) {
		commands.AuctionStart(s, v)
	} else {
		time.Sleep(time.Until(v["start_time"].(time.Time)))
		commands.AuctionStart(s, v)
	}
}

func GiveawayEndTimer(v map[string]interface{}, s *discordgo.Session) {
	fmt.Println("Giveaway Timer Re-Started: ", v["item"], "GuildID: ", v["guild_id"], "ImageURL", v["image_url"], "Host", v["host"], "End Time", v["end_time"].(time.Time).String())
	if v["end_time"].(time.Time).Before(time.Now()) {
		if v["finished"] == true {
			time.Sleep(time.Until(v["end_time"].(time.Time).Add(24 * time.Hour)))
			database.DB.Delete(database.Giveaway{}, v["message_id"].(string))
		} else {
			commands.GiveawayEnd(s, v["message_id"].(string))
		}
	} else {
		time.Sleep(time.Until(v["end_time"].(time.Time)))
		commands.GiveawayEnd(s, v["message_id"].(string))
	}
}
