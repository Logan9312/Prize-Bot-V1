package connect

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/database"
)

type slashCommands struct {
	local, prod []*discordgo.ApplicationCommand
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
	},
	prod: []*discordgo.ApplicationCommand{
		&commands.HelpCommand,
		&commands.AuctionCommand,
		&commands.BidCommand,
		&commands.GiveawayCommand,
		&commands.PrivacyCommand,
		&commands.DevCommand,
	},
}

func BotConnect(token, environment, botName string) {

	fmt.Println(botName + " Starting Up...")
	var s *discordgo.Session
	var err error

	s, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("discordgo.New error:" + err.Error())
	}

	defer s.ChannelMessageSend("915768633752449054", "Bot has finished restarting")

	commands.Session = s

	err = s.Open()

	if err != nil {
		fmt.Println("s.Open error: " + err.Error())
		return
	}

	//Builds local commands
	if environment == "local" {
		for _, v := range s.State.Guilds {
			_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, botCommands.local)
			fmt.Println("Commands added to guild: " + v.Name)
			if err != nil {
				fmt.Println("Bulk Overwrite Error:", err)
			}
		}
		database.DB.Model(database.AuctionSetup{
			GuildID: "915767892467920967",
		}).Create(map[string]interface{}{
			"GuildID":     "915767892467920967",
			"category":    "915768615742103625",
			"log_channel": "937519464671490118",
		})
	}

	//Builds prod commands
	if environment == "prod" {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", botCommands.prod)
		if err != nil {
			fmt.Println("Bulk Overwrite Error:", err)
		}
	}

	s.AddHandler(CommandHandler)
	s.AddHandler(MessageHandler)

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
}

func AuctionEndTimer(v map[string]interface{}, s *discordgo.Session) {
	fmt.Println("Auction Timer Re-Started: ", v["item"], "GuildID: ", v["guild_id"], "ImageURL", v["image_url"], "Host", v["host"], "End Time", v["end_time"].(time.Time).String())
	if v["end_time"].(time.Time).Before(time.Now()) {
		commands.AuctionEnd(v)
	} else {
		time.Sleep(time.Until(v["end_time"].(time.Time)))
		commands.AuctionEnd(v)
	}
}

func AuctionStartTimer(v map[string]interface{}, s *discordgo.Session) {
	fmt.Println("Auction Re-Queued: ", v["item"], "GuildID: ", v["guild_id"], "ImageURL", v["image_url"], "Host", v["host"], "Start Time", v["start_time"].(time.Time).String())
	if v["start_time"].(time.Time).Before(time.Now()) {
		commands.AuctionCreate(s, v)
	} else {
		time.Sleep(time.Until(v["start_time"].(time.Time)))
		commands.AuctionCreate(s, v)
	}
}

func GiveawayEndTimer(v map[string]interface{}, s *discordgo.Session) {
	fmt.Println("Giveaway Timer Re-Started: ", v["item"], "GuildID: ", v["guild_id"], "ImageURL", v["image_url"], "Host", v["host"], "End Time", v["end_time"].(time.Time).String())
	if v["end_time"].(time.Time).Before(time.Now()) {
		if v["finished"].(bool) {
			time.Sleep(time.Until(v["end_time"].(time.Time).Add(24 * time.Hour)))
			database.DB.Delete(database.Giveaway{}, v["message_id"].(string))
		} else {
			commands.GiveawayEnd(commands.Session, v["message_id"].(string))
		}
	} else {
		time.Sleep(time.Until(v["end_time"].(time.Time)))
		commands.GiveawayEnd(commands.Session, v["message_id"].(string))
	}
}
