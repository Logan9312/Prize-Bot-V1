package commands

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gorm.io/gorm/clause"
)

var ProfileCommand = discordgo.ApplicationCommand{
	Name:        "profile",
	Description: "Displays a user's profile.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "Chose who's profile to display",
			Required:    true,
		},
	},
}

func Profile(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := h.ParseSlashCommand(i)

	var balance float64
	userMap := map[string]any{}

	member, err := s.GuildMember(i.GuildID, options["user"].(string))
	if err != nil {
		return err
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.UserProfile{}).Create(map[string]any{
		"guild_id": i.GuildID,
		"user_id":  member.User.ID,
		"balance":  0,
	})
	if result.Error != nil {
		return result.Error
	}

	result = database.DB.Model(database.UserProfile{}).First(userMap, map[string]any{
		"guild_id": i.GuildID,
		"user_id":  member.User.ID,
	})
	if result.Error != nil {
		return result.Error
	}

	balance = userMap["balance"].(float64)

	return h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**__" + member.User.Username + "__**" + "#" + member.User.Discriminator,
		Description: "Profile for: <@" + member.User.ID + ">",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Currency**",
				Value:  PriceFormat(balance, i.GuildID, nil),
				Inline: true,
			},
		},
	})
}
