package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gorm.io/gorm/clause"
)

var DevCommand = discordgo.ApplicationCommand{
	Name:        "dev",
	Description: "Developer only commands!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "version",
			Description: "Change the version of the bot",
		},
	},
}

func Dev(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := h.ParseSlashCommand(i)
	guildMap := map[string]interface{}{}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.AuctionSetup{}).Create(&database.DevSetup{
		Bot: "main",
	})
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	result = database.DB.Model(database.DevSetup{}).First(guildMap)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	if options["version"] != nil {
		err := s.UpdateGameStatus(0, "Bot Version "+options["version"].(string))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	result = database.DB.Model(database.DevSetup{
		Bot:     "main",
	}).Updates(options)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

}
