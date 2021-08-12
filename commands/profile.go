package commands

import (
	"github.com/bwmarrin/discordgo"
)

var ProfileCommand = discordgo.ApplicationCommand{
	Name:        "profile",
	Description: "Displays a user's profile.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "username",
			Description: "Chose who's profile to display",
			Required:    true,
			Choices:     []*discordgo.ApplicationCommandOptionChoice{},
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	},
}

func Profile(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
	username := i.ApplicationCommandData().Options[0].UserValue(s).Username
	Discriminator := i.ApplicationCommandData().Options[0].UserValue(s).Discriminator

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         "",
			Components:      []discordgo.MessageComponent{},
			Embeds:          []*discordgo.MessageEmbed{
				{
					Title: "**__" + username + "__**" + "#" + Discriminator,
					Description: "Inventory For: <@" + userID + ">",
					Color: 0x8073ff,
					Fields:    []*discordgo.MessageEmbedField{},
				},
			},
		},
	})
}