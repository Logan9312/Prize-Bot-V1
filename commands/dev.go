package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/config"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gitlab.com/logan9312/discord-auction-bot/logger"
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

var StatsCommand = discordgo.ApplicationCommand{
	Name:        "stats",
	Description: "Developer only commands!",
}

func Dev(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := h.ParseSlashCommand(i)

	// Check for nil Member (DM context)
	if i.Member == nil || i.Member.User == nil {
		return fmt.Errorf("this command cannot be used in DMs")
	}

	if !config.IsDevUser(i.Member.User.ID) {
		return fmt.Errorf("User must be a developer to run this command")
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.DevSetup{}).Create(map[string]interface{}{
		"bot_id": s.State.User.ID,
	})
	if result.Error != nil {
		logger.Sugar.Warnw("dev command database error", "error", result.Error)
		return result.Error
	}

	if options["version"] != nil {
		err := s.UpdateGameStatus(0, "Bot Version "+options["version"].(string))
		if err != nil {
			return err
		}
	}

	result = database.DB.Model(database.DevSetup{
		BotID: s.State.User.ID,
	}).Updates(options)
	if result.Error != nil {
		logger.Sugar.Warnw("dev command database error", "error", result.Error)
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

func Stats(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	members := 0
	largestServerCount := 0
	largestServerName := ""

	for _, v := range s.State.Guilds {
		members += v.MemberCount
		if v.MemberCount > largestServerCount {
			largestServerCount = v.MemberCount
			largestServerName = v.Name
		}
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title: "Statistics",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Members",
				Value:  fmt.Sprint(members),
				Inline: false,
			},
			{
				Name:   "Servers",
				Value:  fmt.Sprint(len(s.State.Guilds)),
				Inline: false,
			},
			{
				Name:   largestServerName,
				Value:  fmt.Sprint(largestServerCount),
				Inline: false,
			},
		},
	})

	return nil
}
