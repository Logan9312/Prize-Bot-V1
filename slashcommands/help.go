package slashcommands

import (
	"github.com/bwmarrin/discordgo"
)

func Help (s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components:      []discordgo.MessageComponent{},
			Embeds:          []*discordgo.MessageEmbed{
				{
					URL:         "",
					Type:        "",
					Title:       "Discord Bot Help",
					Description: "This bot currently only supports Auctions.",
					Timestamp: "",
					Color:     0,
					Footer:    &discordgo.MessageEmbedFooter{},
					Image:     &discordgo.MessageEmbedImage{},
					Thumbnail: &discordgo.MessageEmbedThumbnail{},
					Video:     &discordgo.MessageEmbedVideo{},
					Provider:  &discordgo.MessageEmbedProvider{},
					Author:    &discordgo.MessageEmbedAuthor{},
					Fields:    []*discordgo.MessageEmbedField{
						{
							Name: "`/Help`:",
							Value: "A command that displays bot functionality.",
						},
						{
							Name: "`/Auction`:",
							Value: "Allows you to put an item up for auction.",
						},
					},
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			Flags:           0,
		},
	})
}