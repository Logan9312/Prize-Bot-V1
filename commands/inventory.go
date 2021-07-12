package commands

import (
	"github.com/bwmarrin/discordgo"
)

func Inventory(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userID := i.ApplicationCommandData().Options[0].UserValue(s).ID
	username := i.ApplicationCommandData().Options[0].UserValue(s).Username
	Discriminator := i.ApplicationCommandData().Options[0].UserValue(s).Discriminator

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         "",
			Components:      []discordgo.MessageComponent{
			},
			Embeds:          []*discordgo.MessageEmbed{
				{
					Title: "**__" + username + "__**" + "#" + Discriminator,
					Description: "Inventory For: <@" + userID + ">",
					Timestamp: "",
					Color: 0x8073ff,
					Fields:    []*discordgo.MessageEmbedField{},
				},
			},
		},
	})
}