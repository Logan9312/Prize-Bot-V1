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
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "host_role",
					Description: "Set a role that can host auctions. Choosing @everyone will reset it to default.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "claiming",
					Description: "Set the message that will appear when someone tries to claim an auction prize",
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
					Name:        "currency",
					Description: "A one time currency to use for this auction.",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "increment_max",
					Description: "The max amount someone can bid at once",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "increment_min",
					Description: "The minimum amount someone can bid at once",
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
			},
		},
	},
}
var BidCommand = discordgo.ApplicationCommand{
	Name:        "bid",
	Description: "Bid on an Auction",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        10,
			Name:        "amount",
			Description: "Place your bid here",
			Required:    true,
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
	ErrorResponse(s, i, "Help command has not been setup yet")
}

func AuctionSetup(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&(1<<3) != 8 {
		ErrorResponse(s, i, "User must have administrator permissions to run this command")
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
		info.AuctionCategory = category

		ch, err := s.Channel(category)

		if err != nil {
			ErrorResponse(s, i, err.Error())
			return
		}

		if ch.Type != 4 {
			content = content + "• ERROR: Auction Category must be a category, not a channel.\n"
		} else {

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
			ErrorResponse(s, i, err.Error())
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

	if options["host_role"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}

		info.AuctionHostRole = options["host_role"].(string)

		if i.GuildID == options["host_role"].(string) {
			info.AuctionHostRole = ""
		}

		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"auction_host_role": info.AuctionHostRole}),
		}).Create(&info)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		content = content + "• Host Role Message has been successfully set.\n"
	}

	if options["claiming"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		info.Claiming = options["claiming"].(string)
		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"claiming": info.Claiming}),
		}).Create(&info)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		content = content + "• Claiming Message has been successfully set.\n"
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
			ErrorResponse(s, i, err.Error())
			return
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
	if info.Claiming == "" {
		info.Claiming = "Not Set"
	}
	if info.AuctionHostRole == "" {
		info.AuctionHostRole = "Not Set"
	} else {
		info.AuctionHostRole = "<@&" + info.AuctionHostRole + ">"
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: componentValue,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Auction Setup",
					Description: content,
					Color:       0x8073ff,
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
							Value: info.Currency + "\n*Note: The bot can only use emojis from discord servers that it is in.*",
						},
						{
							Name:  "**Alert Role**",
							Value: info.AuctionRole,
						},
						{
							Name:  "**Host Role**",
							Value: info.AuctionHostRole,
						},
						{
							Name:  "**Claiming Message**",
							Value: info.Claiming,
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
	image := ""
	var maxBid float64
	var minBid float64
	info := database.GuildInfo{
		GuildID: i.GuildID,
	}
	incCurrency := "+"

	currentTime := time.Now()
	inputDuration := options["duration"].(string)
	if strings.HasSuffix(strings.ToLower(inputDuration), "d") {
		inputDuration = strings.TrimSuffix(inputDuration, "d")
		float, err := strconv.ParseFloat(inputDuration, 64)
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}
		inputDuration = fmt.Sprint(float*24) + "h"
	}
	duration, err := time.ParseDuration(inputDuration)
	if err != nil {
		fmt.Println(err)
		ErrorResponse(s, i, err.Error())
		return
	}
	endTime := currentTime.Add(duration)

	if len(item) > 100 {
		return
	}

	database.DB.First(&info, i.GuildID)

	canHost := false

	if info.AuctionHostRole != "" {
		for _, v := range i.Member.Roles {
			if v == info.AuctionHostRole {
				canHost = true
			}
		}
		if i.Member.Permissions&(1<<3) == 8 {
			canHost = true
		}
		if !canHost {
			ErrorResponse(s, i, "User must be administrator or have the role <@&"+info.AuctionHostRole+"> to host auctions.")
			return
		}
	}

	currency := info.Currency

	channelInfo := discordgo.GuildChannelCreateData{
		Name:     "💸│" + item,
		Type:     0,
		ParentID: info.AuctionCategory,
	}

	channel, err := s.GuildChannelCreateComplex(i.Interaction.GuildID, channelInfo)

	if err != nil {
		fmt.Println(err)
		ErrorResponse(s, i, err.Error())
		return
	}

	description := fmt.Sprintf("%s has hosted an auction! To bid, use the command `/auction bid` in the channel below.", i.Member.Mention())
	details := fmt.Sprintf("**Auction End Time:\n%s**", fmt.Sprintf("<t:%d>", endTime.Unix()))
	savedDescription := ""

	if options["description"] != nil {
		details += "\n**Description:**\n" + options["description"].(string)
		savedDescription = options["description"].(string)
	}
	if options["image"] != nil {
		image = options["image"].(string)
	}
	if options["currency"] != nil {
		currency = options["currency"].(string)
	}

	if currency != "" {
		incCurrency = currency
	}

	if options["increment_min"] != nil {
		minBid = options["increment_min"].(float64)
		details += "\n**Minimum Bid Increment:**\n" + incCurrency + " " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", minBid), "0"), ".")
	}
	if options["increment_max"] != nil {
		maxBid = options["increment_max"].(float64)
		details += "\n**Max Bid Increment:**\n" + incCurrency + " " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", maxBid), "0"), ".")
	}

	details += "\n\u200b"

	guild, err := s.Guild(i.GuildID)

	if err != nil {
		fmt.Println(err)
	}

	message, err := s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: info.AuctionRole,
		Embed: &discordgo.MessageEmbed{
			Title:       "Item: __**" + item + "**__",
			Description: description,
			Color:       0x8073ff,
			Image: &discordgo.MessageEmbedImage{
				URL: image,
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: guild.IconURL(),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "__**Details:**__",
					Value:  details,
					Inline: true,
				},
				{
					Name:   "__**Starting Bid:**__",
					Value:  fmt.Sprintf("%s %s\n\u200b", currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", initialBid), "0"), ".")),
					Inline: true,
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "End Auction",
						Style: 4,
						Emoji: discordgo.ComponentEmoji{
							Name: "🛑",
						},
						CustomID: "endauction",
					},
					discordgo.Button{
						Label:    "Clear Chat",
						Style:    3,
						CustomID: "clearauction",
						Emoji: discordgo.ComponentEmoji{
							Name: "restart",
							ID:   "835685528917114891",
						},
						Disabled: false,
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		ErrorResponse(s, i, err.Error())
		return
	}

	database.DB.Create(&database.Auction{
		ChannelID:   message.ChannelID,
		Bid:         initialBid,
		MessageID:   message.ID,
		EndTime:     endTime,
		Winner:      "No bidders",
		GuildID:     i.GuildID,
		Item:        item,
		Host:        i.Member.User.ID,
		Currency:    currency,
		MaxBid:      maxBid,
		MinBid:      minBid,
		Description: savedDescription,
	})

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "",
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "**Auction Started**",
					Description: "Auction has successfully been started in <#" + message.ChannelID + ">.",
					Color:       0x8073ff,
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
	currency := guildInfo.Currency
	var Content string

	if info.Currency != "" {
		currency = info.Currency
	}

	if info.EndTime.Before(time.Now()) {
		ErrorResponse(s, i, "Cannot Bid, Auction has ended")
		return
	}

	if bidAmount > info.Bid {
		if bidAmount-info.Bid < info.MinBid {
			ErrorResponse(s, i, "Bid must be higher than the previous bid by: "+info.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", info.MinBid), "0"), "."))
			return
		}

		if bidAmount-info.Bid > info.MaxBid && info.MaxBid != 0 {
			ErrorResponse(s, i, "Bid must be no more than "+info.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", info.MaxBid), "0"), ".")+" Higher than the previous bid.")
			return
		}

		info.Bid = bidAmount
		info.Winner = i.Member.User.ID
		Winner := info.Winner

		database.DB.Model(&info).Updates(info)

		updateAuction, err := s.ChannelMessage(info.ChannelID, info.MessageID)
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}

		bidHistory := ""
		username := i.Member.Nick
		if username == "" {
			username = i.Member.User.Username
		}

		if len(updateAuction.Embeds[0].Fields) == 4 {
			bidHistory = updateAuction.Embeds[0].Fields[3].Value + "\n-> " + username + ": " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), ".")
		} else if len(updateAuction.Embeds) == 2 {
			bidHistory = updateAuction.Embeds[1].Description + "\n-> " + username + ": " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), ".")
		} else {
			bidHistory = "-> " + username + ": " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), ".")
		}

		if len(strings.ReplaceAll(bidHistory, " ", "")) >= 4096 {
			bidHistory = "BidHistory was too long and has been reset to prevent a crash.\n-> " + username + ": " + fmt.Sprint(bidAmount)
		}

		updateAuction.Embeds[0].Fields = []*discordgo.MessageEmbedField{
			updateAuction.Embeds[0].Fields[0],
			{
				Name:   "__**Current Highest Bid:**__",
				Value:  fmt.Sprintf("%s %s\n\u200b", currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), ".")),
				Inline: true,
			},
			{
				Name:   "__**Current Winner**__",
				Value:  fmt.Sprintf("<@%s>", Winner) + "\n\u200b",
				Inline: true,
			},
		}

		if len(updateAuction.Embeds) != 2 {
			updateAuction.Embeds = append(updateAuction.Embeds, &discordgo.MessageEmbed{
				Title:       "**Bid History**",
				Description: bidHistory,
				Color:       0x8073ff,
				Image: &discordgo.MessageEmbedImage{
					URL: "https://i.imgur.com/9wo7diC.png",
				},
			})
		} else {
			updateAuction.Embeds[1].Description = bidHistory
		}

		_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Components: updateAuction.Components,
			Embeds:     updateAuction.Embeds,
			ID:         info.MessageID,
			Channel:    info.ChannelID,
		})
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}
		Content = "Bid has successfully been placed"
	} else {
		ErrorResponse(s, i, "You must bid higher than: "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", info.Bid), "0"), "."))
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: Content,
					Color: 0x8073ff,
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
	username := ""
	auctionInfo.ChannelID = ChannelID

	result := database.DB.First(&auctionInfo, ChannelID)
	if result.Error != nil {
		fmt.Println("Error checking for auction: " + result.Error.Error())

		_, err := Session.ChannelDelete(ChannelID)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	result = database.DB.First(&guildInfo, GuildID)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	message := discordgo.NewMessageEdit(auctionInfo.ChannelID, auctionInfo.MessageID)
	messageEmbeds, err := Session.ChannelMessage(auctionInfo.ChannelID, auctionInfo.MessageID)
	if err != nil {
		fmt.Println(err)
	}

	message.Embeds = append(messageEmbeds.Embeds, &discordgo.MessageEmbed{
		Title:       "Auction has ended!",
		Description: "Thank you for participating!",
		Color:       0x32CD32,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://c.tenor.com/MvFFZxXwdpwAAAAC/sold-ray.gif",
		},
	})
	message.Components = 		[]discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "Support Server",
					Style: discordgo.LinkButton,
					Emoji: discordgo.ComponentEmoji{
						Name:     "logo",
						ID:       "889025400120950804",
						Animated: false,
					},
					URL: "https://discord.gg/RxP2z5NGtj",
				},
			},
		},
	}
	Session.ChannelMessageEditComplex(message)

	description := fmt.Sprintf("**Item:** %s", auctionInfo.Item)
	if auctionInfo.Description != "" {
		description += fmt.Sprintf("\n**Description:** %s", auctionInfo.Description)
	}
	if auctionInfo.Winner == "" {
		auctionInfo.Winner = "No winner detected. Please contact support to report this bug"
	} else if auctionInfo.Winner != "No bidders" {
		user, err := Session.User(strings.Trim(auctionInfo.Winner, "<@!>"))
		if err != nil {
			fmt.Println(err)
		}
		username = fmt.Sprintf("(%s#%s)", user.Username, user.Discriminator)
		auctionInfo.Winner = "<@" + auctionInfo.Winner + ">"
	}

	messageSend := discordgo.MessageSend{
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Claim prize!",
						Style: 3,
						Emoji: discordgo.ComponentEmoji{
							Name: "cryopod",
							ID:   "889307390690885692",
						},
						CustomID: "claim_prize",
					},
				},
			},
		},
		Embed: &discordgo.MessageEmbed{
			Title:       "Auction Completed!",
			Description: description,
			Color:       0x8073ff,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "**Winner**",
					Value:  auctionInfo.Winner + " " + username,
					Inline: true,
				},
				{
					Name:   "**Payment Due**",
					Value:  fmt.Sprintf("%s %s\n\u200b", guildInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Bid), "0"), ".")),
					Inline: true,
				},
				{
					Name:   "**Auction Host**",
					Value:  fmt.Sprintf("This Auction was hosted by: <@!%s>", auctionInfo.Host),
					Inline: false,
				},
			},
		},
	}

	if guildInfo.LogChannel == "" {
		fmt.Println("Log channel has not been set for guild: " + GuildID)
		ErrorMessage(Session, ChannelID, "Auction cannot end because log channel has not been set. Please setup an auction log using `/auction setup`")
		return
	}

	_, err = Session.ChannelMessageSendComplex(guildInfo.LogChannel, &messageSend)
	if err != nil {
		fmt.Println(err)
		ErrorMessage(Session, ChannelID, err.Error())
		return
	}

	time.Sleep(30 * time.Second)

	_, err = Session.ChannelDelete(ChannelID)
	if err != nil {
		fmt.Println(err)
	}

	result = database.DB.Delete(&auctionInfo, ChannelID)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

}

func AuctionEndButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	guildInfo := database.GuildInfo{}
	info := database.Auction{}

	database.DB.First(&info, i.ChannelID)
	database.DB.First(&guildInfo, i.GuildID)

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != info.Host {
		ErrorResponse(s, i, "You must have an administrator role to end the auction early!")
		return
	}

	if guildInfo.LogChannel == "" {
		fmt.Println("Log channel has not been set for guild: " + i.GuildID)
		ErrorResponse(s, i, "Auction cannot end because log channel has not been set. Please setup an auction log using `/auction setup`")
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "",
			Components: []discordgo.MessageComponent{},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "End Auction",
					Description: "Auction Ending...",
					Timestamp:   "",
					Color:       0x8073ff,
				},
			},
			Flags: 64,
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	AuctionEnd(i.ChannelID, i.GuildID)
}

func ClaimPrizeButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var auctionInfo database.Auction
	var guildInfo database.GuildInfo
	database.DB.First(&auctionInfo, i.ChannelID)
	database.DB.First(&guildInfo, i.GuildID)

	if guildInfo.Claiming == "" {
		guildInfo.Claiming = "The discord owner has not set a claiming message. Common ways to claim include: Opening a ticket or contacting the auction host. \nTo customize this message, use the command: `/auction setup claiming:`."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Claim Prize",
					Description: guildInfo.Claiming,
					Color:       0x8073ff,
				},
			},
			Flags: 64,
		},
	})
}

func ClearAuctionButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var auctionInfo database.Auction

	result := database.DB.First(&auctionInfo, i.ChannelID)
	if result.Error != nil {
		fmt.Println("Error checking for auction. Auction has likely ended")
		auctionInfo.Host = ""
	}

	if i.Member.Permissions&(1<<3) != 8 && i.User.ID != auctionInfo.Host {
		ErrorResponse(s, i, "User must be host or have administrator permissions to run this command")
		return
	}

	for {
		messageIDs := make([]string, 0)
		messages, err := s.ChannelMessages(i.ChannelID, 100, "", i.Message.ID, "")
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}

		if len(messages) == 0 {
			break
		}

		for _, v := range messages {
			messageIDs = append(messageIDs, v.ID)
		}

		err = s.ChannelMessagesBulkDelete(i.ChannelID, messageIDs)
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Success!",
					Description: "Chat has been cleared",
					Color:       0x8073ff,
				},
			},
			Flags: 64,
		},
	})
}
