package slashcommands

import (
	"github.com/bwmarrin/discordgo"
)

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Item:" + "Grapes",
					Description: "",
					Timestamp:   "",
					Color:       0,
				},
			},
		},
	})
}

func AuctionButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	content := ""
	switch i.MessageComponentData().CustomID {
	case "1":
		content += "Bid Successful"
	case "2":
		content += " no"

	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			TTS:     false,
			Content: content,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "Click here to Bid",
							Style:    1,
							Disabled: false,
							Emoji:    discordgo.ButtonEmoji{},
							CustomID: "1",
						},
						&discordgo.Button{
							Label:    "Bid Successful",
							Style:    3,
							Disabled: true,
							Emoji:    discordgo.ButtonEmoji{},
							Link:     "",
							CustomID: "2",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Item: Test Item",
					Description: "Bid Successful!",
				},
			},
		},
	})

}
