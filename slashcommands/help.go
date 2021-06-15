package slashcommands

import (
	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Close Help Menu",
							Style:    3,
							Disabled: false,
							Emoji: discordgo.ButtonEmoji{
								Name:     "",
								ID:       "",
								Animated: false,
							},
							CustomID: "Help",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					URL:         "",
					Type:        "",
					Title:       "Discord Bot Help",
					Description: "This bot currently only supports Auctions.",
					Timestamp:   "",
					Color:       0,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "**/Help**:",
							Value: "A command that displays bot functionality.",
						},
						{
							Name:  "**/Auction**:",
							Value: "Allows you to put an item up for auction.",
						},
					},
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		},
	})

}

func HelpButton (s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.WebhookDelete(i.Message.WebhookID)
}