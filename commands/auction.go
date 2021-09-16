package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	"gorm.io/gorm/clause"
)

var Session *discordgo.Session

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
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "category",
					Description: "Sets the category to create auctions in. Name must be an exact match",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "currency",
					Description: "Sets the auction currency",
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "log_channel",
					Description: "Sets the channel where auctions will send outputs when they end",
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "alert_role",
					Description: "Set a role to get pinged whenever an auction starts",
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "duration",
					Description: "Time that auction will run for. (Example: 24h, or 1d)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "Set a custom item description",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "image",
					Description: "Must be a link",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "bid",
			Description: "Bid on an Auction",
			Required:    false,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        10,
					Name:        "amount",
					Description: "Place your bid here",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "secret_bidder",
					Description: "Turn this on to protect your identity for the next bid.",
					Required:    false,
				},
			},
		},
	},
}

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[0].Name {
	case "help":
		AuctionHelp(s, i)
	case "setup":
		AuctionSetup(s, i)
	case "create":
		AuctionCreate(s, i)
	case "bid":
		AuctionBid(s, i)
	}
}

func AuctionHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "/Auction Help has not been setup yet",
					Color: 0x00bfff,
				},
			},
			Flags: 64,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionSetup(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&(1<<3) != 8 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "User must be administrator to change auction settings",
						Color: 0x00bfff,
					},
				},
				Flags: 64,
			},
		})

		if err != nil {
			fmt.Println(err)
		}
		return
	}

	options := ParseSubCommand(i)
	content := ""
	componentValue := []discordgo.MessageComponent{}
	category := &discordgo.Channel{}

	if options["category"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		category := options["category"].(string)
		ch, err := s.Channel(category)
		if err != nil {
			fmt.Println(err)
			return
		}
		if ch.Type != 4 {
			content = content + "• ERROR: Auction Category must be a category, not a channel.\n"
		} else {

			info.AuctionCategory = category
			result := database.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "guild_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"auction_category": info.AuctionCategory}),
			}).Create(&info)

			if result.Error != nil {
				fmt.Println(result.Error.Error())
			}
			content = content + "• Category has been successfully set.\n"
		}
	}

	if options["currency"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		info.Currency = options["currency"].(string)
		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"currency": info.Currency}),
		}).Create(&info)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		content = content + "• Currency has been successfully set.\n"
	}

	if options["log_channel"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		info.LogChannel = options["log_channel"].(string)

		ch, err := s.Channel(info.LogChannel)
		if err != nil {
			fmt.Println(err)
			return
		}

		if ch.Type != 0 {
			content = content + "• ERROR: Auction Log must be a text channel\n"
		} else {
			result := database.DB.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "guild_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{"log_channel": info.LogChannel}),
			}).Create(&info)

			if result.Error != nil {
				fmt.Println(result.Error.Error())
			}
			content = content + "• Log Channel has been successfully set.\n"
		}
	}

	if options["alert_role"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		info.AuctionRole = fmt.Sprintf("<@&%s>", options["alert_role"].(string))
		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"auction_role": info.AuctionRole}),
		}).Create(&info)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		content = content + "• Alert Role has been successfully set.\n"
	}

	info := database.GuildInfo{
		GuildID: i.GuildID,
	}
	database.DB.First(&info, i.GuildID)

	if info.AuctionCategory == "" {
		category.Name = "Not Set"
	} else {
		var err error
		category, err = s.Channel(info.AuctionCategory)
		if err != nil {
			fmt.Println(err)
		}
	}
	if info.Currency == "" {
		info.Currency = "Not Set"
	}
	if info.AuctionRole == "" {
		info.AuctionRole = "Not Set"
	}
	if info.LogChannel == "" {
		info.LogChannel = "Not Set"
	} else {
		info.LogChannel = fmt.Sprintf("<#%s>", info.LogChannel)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: componentValue,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Auction Setup",
					Description: content,
					Color:       0x00bfff,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "**Category**",
							Value: category.Name,
						},
						{
							Name:  "**Log Channel**",
							Value: info.LogChannel,
						},
						{
							Name:  "**Currency**",
							Value: info.Currency,
						},
						{
							Name:  "**Alert Role**",
							Value: info.AuctionRole,
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

func AuctionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := ParseSubCommand(i)
	item := options["item"].(string)
	initialBid := options["startingbid"].(float64)
	description := ""
	image := ""
	info := database.GuildInfo{
		GuildID: i.GuildID,
	}

	currentTime := time.Now()
	inputDuration := options["duration"].(string)
	if strings.HasSuffix(strings.ToLower(inputDuration), "d") {
		inputDuration = strings.TrimSuffix(inputDuration, "d")
		float, err := strconv.ParseFloat(inputDuration, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		inputDuration = fmt.Sprint(float*24) + "h"
	}

	duration, err := time.ParseDuration(inputDuration)
	if err != nil {
		fmt.Println(err)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:    "",
				Components: []discordgo.MessageComponent{},
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "**Time Duration Error**",
						Description: "Error message: " + err.Error(),
						Color:       0x00bfff,
					},
				},
				Flags: 64,
			},
		})

		if err != nil {
			fmt.Println(err)
		}
		return
	}
	endTime := currentTime.Add(duration)

	if len(item) > 100 {
		return
	}

	database.DB.First(&info, i.GuildID)

	channelInfo := discordgo.GuildChannelCreateData{
		Name:     "💸│" + item,
		Type:     0,
		ParentID: info.AuctionCategory,
	}

	channel, err := s.GuildChannelCreateComplex(i.Interaction.GuildID, channelInfo)

	if err != nil {
		fmt.Println(err)
	}

	if options["description"] != nil {
		description = options["description"].(string)
	}
	if options["image"] != nil {
		image = options["image"].(string)
	}

	message, err := s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: info.AuctionRole,
		Embed: &discordgo.MessageEmbed{
			Title:       "Item: " + item,
			Description: description,
			Color:       0x00bfff,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "To bid in this auction use /auction bid in the channel below!",
			},
			Image: &discordgo.MessageEmbedImage{
				URL: image,
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "**Auction End Time:**",
					Value:  fmt.Sprintf("<t:%d>", endTime.Unix()),
					Inline: true,
				},
				{
					Name:   "**Starting Bid:**",
					Value:  fmt.Sprint(initialBid),
					Inline: true,
				},
				{
					Name:   "**Auction Host:**",
					Value:  i.Member.Mention(),
					Inline: true,
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "End Auction Early",
						Style: 4,
						Emoji: discordgo.ComponentEmoji{
							Name: "🛑",
						},
						CustomID: "endauction",
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	database.DB.Create(&database.Auction{
		ChannelID: message.ChannelID,
		Bid:       initialBid,
		MessageID: message.ID,
		EndTime:   endTime,
		Winner:    "No bidders",
		GuildID:   i.GuildID,
		Item:      item,
		Host:      i.Member.Mention(),
	})

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "",
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "**Auction Started**",
					Description: "Auction has successfully been started, I might have some bugs to work out so please contact me if there is a failure.",
					Timestamp:   "",
					Color:       0x00bfff,
				},
			},
			Flags: 64,
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(duration)
	AuctionEnd(channel.ID, i.GuildID)
}

func AuctionBid(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := ParseSubCommand(i)
	bidAmount := options["amount"].(float64)
	var info database.Auction
	var guildInfo database.GuildInfo
	info.ChannelID = i.ChannelID
	database.DB.First(&info, i.ChannelID)
	database.DB.First(&guildInfo, i.GuildID)
	host := info.Host
	currency := guildInfo.Currency
	var Content string

	if info.Host == "" {
		host = "No Host Recorded"
	}

	if info.EndTime.Before(time.Now()) {
		Content = "Cannot Bid, auction has ended"
	} else if bidAmount > info.Bid {
		info.Bid = bidAmount
		info.Winner = i.Member.Mention()
		Winner := info.Winner

		database.DB.Model(&info).Updates(info)

		if options["secret_bidder"] != nil {
			Winner = "🤫"
		}

		updateAuction, err := s.ChannelMessage(info.ChannelID, info.MessageID)
		if err != nil {
			fmt.Println(err)
		}

		updateAuction.Embeds[0].Fields = []*discordgo.MessageEmbedField{
			{
				Name:   "**Auction End Time:**",
				Value:  fmt.Sprintf("<t:%d>", info.EndTime.Unix()),
				Inline: true,
			},
			{
				Name:   "**Current Highest Bid:**",
				Value:  fmt.Sprintf("%s %s", currency, fmt.Sprint(info.Bid)),
				Inline: true,
			},
			{
				Name:   "**Current Winner**",
				Value:  fmt.Sprint(Winner),
				Inline: true,
			},
			{
				Name:   "**Auction Host:**",
				Value:  host,
				Inline: false,
			},
		}

		_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Components: updateAuction.Components,
			Embed:      updateAuction.Embeds[0],
			ID:         info.MessageID,
			Channel:    info.ChannelID,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		Content = "Bid has successfully been placed"
	} else {
		Content = "You must bid higher than: " + fmt.Sprint(info.Bid)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: Content,
					Color: 0x00bfff,
				},
			},
			Flags: 64,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionEnd(ChannelID, GuildID string) {
	var auctionInfo database.Auction
	var guildInfo database.GuildInfo
	auctionInfo.ChannelID = ChannelID
	database.DB.First(&auctionInfo, ChannelID)
	database.DB.First(&guildInfo, GuildID)

	messageSend := discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			Title:       "Auction Completed!",
			Description: auctionInfo.Item,
			Timestamp:   "",
			Color:       0x00bfff,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "**Winner**",
					Value:  auctionInfo.Winner,
					Inline: false,
				},
				{
					Name:   "**Payment Due**",
					Value:  guildInfo.Currency + " " + fmt.Sprint(auctionInfo.Bid),
					Inline: false,
				},
			},
		},
	}

	if guildInfo.LogChannel == "" {
		Session.ChannelMessageSend(ChannelID, "ERROR: Auction cannot end because log channel has not been set. Please setup an auction log using `/auction setup`")
		return
	}

	Session.ChannelMessageSend(ChannelID, "Auction has ended, channel will automatically delete in 1 hour")



	_, err := Session.ChannelMessageSendComplex(guildInfo.LogChannel, &messageSend)
	if err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Hour)
	
	auctionInfo.EndTime = time.Date(1, time.January, 1, 1, 0, 0, 0, time.UTC)

	database.DB.Model(&auctionInfo).Updates(auctionInfo)

	_, err = Session.ChannelDelete(ChannelID)
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionEndButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	content := ""

	if i.Member.Permissions&(1<<3) != 8 {
		content = "You must have an administrator role to end the auction!"
	} else {
		content = "Auction Ending..."
		defer AuctionEnd(i.ChannelID, i.GuildID)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "",
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "End Auction",
					Description: content,
					Timestamp:   "",
					Color:       0x00bfff,
				},
			},
			Flags: 64,
		},
	})

	if err != nil {
		fmt.Println(err)
	}
}
