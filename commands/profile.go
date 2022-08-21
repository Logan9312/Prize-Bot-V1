package commands

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
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

	member := i.Member
	var balance float64
	userMap := map[string]any{}

	if options["user"] != nil {
		member = i.ApplicationCommandData().Resolved.Members[options["user"].(string)]
	}

	database.DB.Model(database.UserProfile{}).First(userMap, map[string]any{
		"guild_id": i.GuildID,
		"user_id":  member.User.ID,
	})

	if userMap["balance"] != nil {
		balance = userMap["balance"].(float64)
	}

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
