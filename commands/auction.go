package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	"gorm.io/gorm"
)

var AuctionCommand = discordgo.ApplicationCommand{
	Name:        "auction",
	Description: "Put an item up for auction!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "help",
			Description: "auction info",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "setup",
			Description: "Setup auctions on your server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "Sets the category to create auctions in. Name must be an exact match",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create an Auction",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The item you wish to auction off",
					Required:    true,
				},
				{
					Type:        10,
					Name:        "startingbid",
					Description: "The starting price to bid on",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "bid",
			Description: "Bid on an Auction",
		},
	},
}

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	switch i.ApplicationCommandData().Options[0].Name {
	case "setup":
		AuctionSetup(s, i, db)
	case "create":
		AuctionCreate(s, i, db)
	}
}

func AuctionSetup(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	options := ParseSubCommand(i)
	category := options["category"].(string)
	catIDs := make([]string, 0)
	catMatch := false
	status := "FAILED"
	content := "There is no category in your discord server that matches the name: " + category
	catMenu := make([]discordgo.SelectMenuOption, 0)
	componentValue := []discordgo.MessageComponent{}

	channels, err := s.GuildChannels(i.GuildID)

	if err != nil {
		fmt.Println(err)
	}

	for n, v := range channels {
		if v.Type == 4 {
			if strings.EqualFold(category, v.Name) {
				catMatch = true
				catIDs = append(catIDs, v.ID)
				menuOption := discordgo.SelectMenuOption{
					Label:       fmt.Sprintf("%s (%s)", v.Name, v.ID),
					Value:       v.ID,
					Description: "Discord Channel number (includes all channels/categories): " + fmt.Sprint(n+1),
					Emoji:       discordgo.ComponentEmoji{},
					Default:     false,
				}
				catMenu = append(catMenu, menuOption)
			}
		}
	}

	if catMatch {
		status = "SUCCESS"
		content = "Successfully set the output to the category: `" + category + "`"
		if len(catIDs) > 1 {
			status = "PENDING INPUT"
			content = "You have multiple categories that match the name: " + category + "**. Please select the correct one below."
			catOptions := catMenu
			componentValue = []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "categorymenu",
							Placeholder: "Select Category",
							MinValues:   1,
							MaxValues:   1,
							Options:     catOptions,
						},
					},
				},
			}
		} else {
			info := database.GuildInfo{
				GuildID:         i.GuildID,
				AuctionCategory: catIDs[0],
			}
			db.Create(&info)
		}
	}

	fmt.Println("Catmatch: ", catMatch)
	fmt.Println("catIDs: ", catIDs)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: componentValue,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       fmt.Sprintf("Auction Category Setup: __%s__", status),
					Description: content,
				},
			},
			Flags: 64,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionCreate(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {

	options := ParseSubCommand(i)
	item := options["item"].(string)
	initialBid := options["startingbid"].(float64)
	bidder := i.Member.User.Mention()
	info := database.GuildInfo{}

	db.First(&info, i.Interaction.GuildID)

	channelInfo := discordgo.GuildChannelCreateData{
		Name:     "Auction Test",
		Type:     0,
		ParentID: info.AuctionCategory,
	}

	channel, err := s.GuildChannelCreateComplex(i.Interaction.GuildID, channelInfo)

	if err != nil {
		fmt.Println(err)
	}

	_, err = s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content:         "",
		Embed:           &discordgo.MessageEmbed{
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
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "",
						Style: 2,
						Emoji: discordgo.ComponentEmoji{
							Name: "📩",
						},
						CustomID: "startbid",
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         "",
			Components:      []discordgo.MessageComponent{},
			Embeds:          []*discordgo.MessageEmbed{
				{
					Title:       "**Auction Started**",
					Description: "Auction has successfully been started, I might have some bugs to work out so please contact me if there is a failure.",
					Timestamp:   "",
					Color:       0,
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			Flags:           64,
		},
	})

	if err != nil {
		fmt.Println(err)
	}
}

func CategorySelect(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {

	categoryID := i.MessageComponentData().Values[0]
	fmt.Println("Category ID Stored: " + categoryID)

	info := database.GuildInfo{
		GuildID:         i.GuildID,
		AuctionCategory: categoryID,
	}

	db.Create(&info)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "ID Saved Successfully!",
					Description: "You may now create an auction",
				},
			},
			Components: []discordgo.MessageComponent{},
			Flags:      64,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
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
