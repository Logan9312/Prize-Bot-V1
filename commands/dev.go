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

func Dev(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	devUsers := map[string]interface{}{
		"280812467775471627": "Logan",
	}

	options := h.ParseSlashCommand(i)
	guildMap := map[string]interface{}{}

	if devUsers[i.Member.User.ID] == nil {
		return fmt.Errorf("User must be a developer to run this command")
	}

	devCreate := map[string]interface{}{
		"bot": "main",
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.DevSetup{}).Create(&devCreate)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return result.Error
	}

	result = database.DB.Model(database.DevSetup{}).First(guildMap)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return result.Error
	}

	if options["version"] != nil {
		err := s.UpdateGameStatus(0, "Bot Version "+options["version"].(string))
		if err != nil {
			return err
		}
	}

	result = database.DB.Model(database.DevSetup{
		Bot: "main",
	}).Updates(options)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return result.Error
	}

	fields := []*discordgo.MessageEmbedField{}

	for key, value := range options {
		if value != nil && value != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   key,
				Value:  fmt.Sprint(value),
				Inline: true,
			})
		}
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:  "Success!",
		Fields: fields,
	})

	return nil
}
