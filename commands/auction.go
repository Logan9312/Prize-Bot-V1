package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Auction(s *discordgo.Session, m *discordgo.MessageCreate, details []string) {

	AuctionPrize := strings.Join(details, " ")
	content := "Press the Button to Bid"

	embed := &discordgo.MessageEmbed{
		Title:       "Item: " + AuctionPrize,
		Description: "React below to Bid!",
		Color:       0x00BFFF,
	}

	ar := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			&discordgo.Button{
				Label:    content,
				Style:    1,
				Disabled: false,
				Emoji:    discordgo.ButtonEmoji{},
				CustomID: "1",
			},
			&discordgo.Button{
				Label:    "Click to Confirm",
				Style:    3,
				Disabled: true,
				Emoji:    discordgo.ButtonEmoji{},
				Link:     "",
				CustomID: "2",
			},
		},
	}

	messageSend := &discordgo.MessageSend{
		Content:    "Item: " + AuctionPrize,
		Embed:      embed,
		Components: []discordgo.MessageComponent{ar},
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, messageSend)
	fmt.Println(err)

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
