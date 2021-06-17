package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)
var HelpCommand =	discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Basic bot functionality",
}

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Close Help menu",
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
					Description: "With slash commands you can use any of the following features. In the future I will be adding Inventory tracking and other related commands.",
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

func HelpButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)

	if err != nil {
		fmt.Println(err)
	}
}
