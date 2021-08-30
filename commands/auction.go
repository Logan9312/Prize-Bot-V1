package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var AuctionCommand = discordgo.ApplicationCommand{
	Name:        "auction",
	Description: "Put an item up for auction!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "item",
			Description: "Choose an Item to put up for auction",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "startingbid",
			Description: "Starting Bid Amount",
			Required:    true,
		},
	},
}

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate, AppID string) {

	options := ParseSlashCommand(i)
	item := options["item"].(string)
	initialBid := options["startingbid"].(float64)
	bidder := i.Member.User.Mention()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: 4,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Bid",
							Style:    1,
							Emoji:    discordgo.ComponentEmoji{},
							CustomID: "startbid",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Item: " + item,
					Description: "Current Highest Bid: " + fmt.Sprint(initialBid) + " 🍓",
					Color:       0x00bfff,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "**Current Winner:**",
							Value:  bidder,
							Inline: false,
						},
					},
				},
			},
		},
	})

}

func AuctionButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: i.Member.User.Mention(),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "increment",
							Placeholder: "Select Increment",
							MinValues:   1,
							MaxValues:   1,
							Options: []discordgo.SelectMenuOption{
								{
									Label:       "1",
									Value:       "1",
									Description: "Change the bid by 1",
									Default:     true,
								},
								{
									Label:       "10",
									Value:       "10",
									Description: "Change the bid by 10",
									Default:     false,
								},
								{
									Label:       "100",
									Value:       "100",
									Description: "Change the bid by 100",
									Default:     false,
								},
								{
									Label:       "1000",
									Value:       "1000",
									Description: "Change the bid by 1000",
									Default:     false,
								},
								{
									Label:       "10 000",
									Value:       "10000",
									Description: "Change the bid by 10 000",
									Default:     false,
								},
							},
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Raise",
							Style:    1,
							Disabled: false,
							Emoji:    discordgo.ComponentEmoji{},
							CustomID: "raisebid",
						},
						discordgo.Button{
							Label:    "Lower",
							Style:    1,
							Disabled: false,
							Emoji:    discordgo.ComponentEmoji{},
							CustomID: "lowerbid",
						},
						discordgo.Button{
							Label:    "Bid",
							Style:    2,
							Disabled: false,
							Emoji:    discordgo.ComponentEmoji{},
							CustomID: "placebid",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Select Bid Amount: ",
							Value:  "0",
							Inline: false,
						},
					},
				},
			},
			Flags: 64,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func Bid(s *discordgo.Session, i *discordgo.InteractionCreate) {

}
