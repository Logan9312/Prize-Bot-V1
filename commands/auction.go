package commands

import (
	"fmt"
	"sort"
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
					Description: "Sets the category to create auctions in.",
					ChannelTypes: []discordgo.ChannelType{
						4,
					},
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
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						0,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "alert_role",
					Description: "Set a role to get pinged whenever an auction starts. Choosing @everyone will reset it to default.",
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "snipe_extension",
					Description: "Set 0 to disable. Duration an auction by when a bid is placed within the snipe range. (Example: 5m)",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "snipe_range",
					Description: "Set 0 to disable. The remaining time needed to activate Anti-Snipe (Example: 24h, or 1d)",
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
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "duration",
					Description:  "Time that auction will run for. (Example: 24h, or 1d)",
					Required:     true,
					Autocomplete: true,
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
					Type:        10,
					Name:        "increment_max",
					Description: "The max amount someone can bid at once",
					Required:    false,
				},
				{
					Type:        10,
					Name:        "increment_min",
					Description: "The minimum amount someone can bid at once",
					Required:    false,
				},
				{
					Type:        10,
					Name:        "buyout",
					Description: "Set a price that someone can immediately win the auction for if they bid it or higher.",
					Required:    false,
				},
				{
					Type:        10,
					Name:        "target_price",
					Description: "If this hidden price is not reached, no winner will be chosen",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "category",
					Description: "Sets the category to create auctions in.",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						4,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "image",
					Description: "Must be a link",
					Required:    false,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "schedule",
					Description:  "Set how long until the auction starts. (Example: 24h, or 1d)",
					Required:     false,
					Autocomplete: true,
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
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "queue",
			Description: "Display the current auction Queue",
			Required:    false,
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
		AuctionPlanner(s, i)
	case "bid":
		AuctionBid(s, i)
	case "queue":
		AuctionQueue(s, i)
	}
}

func AuctionAutoComplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch i.ApplicationCommandData().Options[0].Name {
	case "create":
		options := ParseSubCommand(i)

		if options["duration"] != nil {
			choices = TimeSuggestions(options["duration"].(string))
		} else if options["schedule"] != nil {
			choices = TimeSuggestions(options["schedule"].(string))
		} else {
			choices = []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "",
					Value: "",
				},
			}
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func TimeSuggestions(input string) []*discordgo.ApplicationCommandOptionChoice {
	_, err := strconv.ParseFloat(input, 64)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	choices := []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  input + "d",
			Value: input + "d",
		},
		{
			Name:  input + "h",
			Value: input + "h",
		},
		{
			Name:  input + "m",
			Value: input + "m",
		},
		{
			Name:  input + "s",
			Value: input + "s",
		},
	}
	return choices
}

func AuctionHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ErrorResponse(s, i, "Help command has not been setup yet")
}

func AuctionSetup(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&(1<<3) != 8 {
		ErrorResponse(s, i, "User must have administrator permissions to run this command")
		return
	}

	var err error

	options := ParseSubCommand(i)
	content := ""
	category := &discordgo.Channel{}

	if options["category"] != nil {

		info := database.GuildInfo{
			GuildID: i.GuildID,
		}

		category := options["category"].(string)
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

		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"log_channel": info.LogChannel}),
		}).Create(&info)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		content = content + "• Log Channel has been successfully set.\n"
	}

	if options["alert_role"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		info.AuctionRole = fmt.Sprintf("<@&%s>", options["alert_role"].(string))

		if i.GuildID == options["alert_role"] {
			info.AuctionRole = ""
		}

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

	if options["snipe_extension"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		info.SnipeExtension, err = ParseTime(options["snipe_extension"].(string))
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}
		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"snipe_extension": info.SnipeExtension}),
		}).Create(&info)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		content = content + "• Snipe Extension has been successfully set.\n"
	}

	if options["snipe_range"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		info.SnipeRange, err = ParseTime(options["snipe_range"].(string))
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}
		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"snipe_range": info.SnipeRange}),
		}).Create(&info)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		content = content + "• Snipe Range has been successfully set.\n"
	}

	info := database.GuildInfo{
		GuildID: i.GuildID,
	}
	database.DB.First(&info, i.GuildID)

	antiSnipeDescription := fmt.Sprintf("If a bid is placed within %s of the auction ending, it will be extended by %s.", info.SnipeRange.String(), info.SnipeExtension.String())

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
	if info.SnipeExtension == 0 || info.SnipeRange == 0 {
		antiSnipeDescription = "Anti Snipe Disabled. To enable, set both snipe_extension and snipe_range"
	}

	err = SuccessResponse(s, i, PresetResponse{
		Title:       "Auction Setup",
		Description: content,
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
			{
				Name:  "**Anti Snipe**",
				Value: antiSnipeDescription,
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionCreate(s *discordgo.Session, auctionInfo database.AuctionQueue) {

	guildInfo := database.GuildInfo{
		GuildID: auctionInfo.GuildID,
	}
	database.DB.First(&guildInfo, auctionInfo.GuildID)
	currency := auctionInfo.Currency
	incCurrency := "+"
	if currency != "" {
		incCurrency = currency
	}

	details := fmt.Sprintf("**Auction End Time:\n%s**", fmt.Sprintf("<t:%d:R>", auctionInfo.EndTime.Unix()))

	if auctionInfo.Description != "" {
		details += "\n**Description:**\n" + auctionInfo.Description
	}

	if auctionInfo.MinBid != 0 {
		details += "\n**Minimum Bid Increment:**\n" + incCurrency + " " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.MinBid), "0"), ".")
	}

	if auctionInfo.MaxBid != 0 {
		details += "\n**Max Bid Increment:**\n" + incCurrency + " " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.MaxBid), "0"), ".")
	}

	if auctionInfo.TargetPrice != 0 {
		details += "\n**Target Price:**\nThe host has set a hidden target price for this auction."
	}

	if guildInfo.AntiSnipe {
		details += "\n**Anti-Snipe:**\nAnti-Snipe has been enabled for this auction. It can be disabled at any time with `/Auction Setup`."
	}

	if auctionInfo.Buyout != 0 {
		details += "\n**Buyout Price:**\n" + incCurrency + " " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Buyout), "0"), ".")
	}

	details += "\n\u200b"

	guild, err := s.Guild(auctionInfo.GuildID)

	if err != nil {
		fmt.Println(err)
	}

	channelInfo := discordgo.GuildChannelCreateData{
		Name:     "💸│" + auctionInfo.Item,
		Type:     0,
		ParentID: guildInfo.AuctionCategory,
	}

	channel, err := s.GuildChannelCreateComplex(auctionInfo.GuildID, channelInfo)

	if err != nil {
		fmt.Println(err)
		return
	}

	host, err := s.GuildMember(auctionInfo.GuildID, auctionInfo.Host)
	if err != nil {
		fmt.Println(err)
	}

	message, err := PresetMessageSend(s, channel.ID, PresetResponse{
		Content:     guildInfo.AuctionRole,
		Title:       "Auction Item: __**" + auctionInfo.Item + "**__",
		Description: fmt.Sprintf("%s has hosted an auction! To bid, use the command `/auction bid` in the channel below.", host.Mention()),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "__**Details:**__",
				Value:  details,
				Inline: true,
			},
			{
				Name:   "__**Starting Bid:**__",
				Value:  fmt.Sprintf("%s %s\n\u200b", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Bid), "0"), ".")),
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: guild.IconURL(),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: auctionInfo.ImageURL,
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
		ErrorMessage(s, channel.ID, "Error starting auction: " + err.Error())
		return
	}

	database.DB.Create(&database.Auction{
		ChannelID:   message.ChannelID,
		Bid:         auctionInfo.Bid,
		MessageID:   message.ID,
		EndTime:     auctionInfo.EndTime,
		Winner:      auctionInfo.Host,
		GuildID:     auctionInfo.GuildID,
		Item:        auctionInfo.Item,
		Host:        auctionInfo.Host,
		Currency:    auctionInfo.Currency,
		MinBid:      auctionInfo.MinBid,
		MaxBid:      auctionInfo.MaxBid,
		Description: auctionInfo.Description,
		ImageURL:    auctionInfo.ImageURL,
		TargetPrice: auctionInfo.TargetPrice,
		Buyout:      auctionInfo.Buyout,
	})

	database.DB.Delete(auctionInfo, auctionInfo.ID)

	time.Sleep(time.Until(auctionInfo.EndTime))
	AuctionEnd(channel.ID, auctionInfo.GuildID)
}

func AuctionPlanner(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var startTime time.Time
	var endTime time.Time

	options := ParseSubCommand(i)

	endTimeDuration, err := ParseTime(strings.ToLower(options["duration"].(string)))
	if err != nil {
		fmt.Println(err)
		ErrorResponse(s, i, err.Error())
		return
	}

	item := options["item"].(string)
	var image string
	var maxBid float64
	var minBid float64
	var buyout float64
	var targetPrice float64
	var description string
	info := database.GuildInfo{
		GuildID: i.GuildID,
	}

	if len(item) > 100 {
		ErrorResponse(s, i, "Title cannot be over 100 characters long")
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

	if options["category"] != nil {
		info.AuctionCategory = options["category"].(string)
	}

	if options["description"] != nil {
		description = options["description"].(string)
	}
	if options["image"] != nil {
		image = options["image"].(string)
	}
	if options["currency"] != nil {
		currency = options["currency"].(string)
	}
	if options["increment_min"] != nil {
		minBid = options["increment_min"].(float64)
	}
	if options["increment_max"] != nil {
		maxBid = options["increment_max"].(float64)
	}
	if options["buyout"] != nil {
		buyout = options["buyout"].(float64)
	}
	if options["target_price"] != nil {
		targetPrice = options["target_price"].(float64)
	}

	if options["schedule"] != nil {

		if i.Member.Permissions&(1<<3) != 8 {
			ErrorResponse(s, i, "User have administrator permissions to schedule auctions in advance. If you are not an admin you can still run auctions normally without the `schedule` parameter.")
			return
		}

		var AuctionQueue []database.AuctionQueue
		var num int

		database.DB.Find(&AuctionQueue)

		for _, v := range AuctionQueue {
			if v.GuildID == i.GuildID {
				num += 1
			}
		}

		if num >= 25 {
			ErrorResponse(s, i, "You can only schedule 25 auctions in advance.")
			return
		}

		startTimeDuration, err := ParseTime(strings.ToLower(options["schedule"].(string)))
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}

		endTime = time.Now().Add(endTimeDuration).Add(startTimeDuration)
		startTime = time.Now().Add(startTimeDuration)
	} else {
		endTime = time.Now().Add(endTimeDuration)
	}

	auctionData := database.AuctionQueue{
		Bid:         options["startingbid"].(float64),
		StartTime:   startTime,
		EndTime:     endTime,
		GuildID:     i.GuildID,
		Item:        item,
		Host:        i.Member.User.ID,
		Currency:    currency,
		MinBid:      minBid,
		MaxBid:      maxBid,
		Description: description,
		ImageURL:    image,
		Category:    info.AuctionCategory,
		Buyout:      buyout,
		TargetPrice: targetPrice,
	}

	if options["schedule"] != nil {

		database.DB.Create(&auctionData)
		err = PremiumResponse(s, i, PresetResponse{
			Title: "Auction has been Scheduled!",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "**Auction Start Time:**",
					Value:  fmt.Sprintf("<t:%d:R>", startTime.Unix()),
					Inline: false,
				},
			},
		})
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Until(startTime))
	} else {
		err = SuccessResponse(s, i, PresetResponse{
			Title:       "**Auction Starting**",
			Description: "Auction has successfully been started!",
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	AuctionCreate(s, auctionData)
}

func AuctionBid(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var options map[string]interface{}
	var responseFields []*discordgo.MessageEmbedField

	if i.ApplicationCommandData().Options[0].Type == discordgo.ApplicationCommandOptionSubCommand {
		options = ParseSubCommand(i)
	} else {
		options = ParseSlashCommand(i)
	}

	bidAmount := options["amount"].(float64)
	var auctionInfo database.Auction
	var guildInfo database.GuildInfo
	auctionInfo.ChannelID = i.ChannelID
	database.DB.First(&auctionInfo, i.ChannelID)
	database.DB.First(&guildInfo, i.GuildID)
	currency := guildInfo.Currency
	var Content string
	var antiSnipeFlag bool

	if auctionInfo.Currency != "" {
		currency = auctionInfo.Currency
	}

	if auctionInfo.EndTime.Before(time.Now()) {
		ErrorResponse(s, i, "Cannot Bid, Auction has ended")
		return
	}

	if i.Member.User.ID == auctionInfo.Winner {
		ErrorResponse(s, i, "Cannot out bid yourself")
		return
	}

	if time.Until(auctionInfo.EndTime) < guildInfo.SnipeRange && guildInfo.SnipeExtension != 0 {
		auctionInfo.EndTime = auctionInfo.EndTime.Add(guildInfo.SnipeExtension)
		responseFields = []*discordgo.MessageEmbedField{
			{
				Name:   "**Anti-Snipe Activated!**",
				Value:  fmt.Sprintf("New End Time: <t:%d>", auctionInfo.EndTime.Unix()),
				Inline: false,
			},
		}
		antiSnipeFlag = true
	}

	if bidAmount >= auctionInfo.Buyout && auctionInfo.Buyout != 0 {

		auctionInfo.Bid = bidAmount
		auctionInfo.Winner = i.Member.User.ID

		database.DB.Model(&auctionInfo).Updates(auctionInfo)

		AuctionEnd(i.ChannelID, i.GuildID)

		SuccessResponse(s, i, PresetResponse{
			Title:       "Success!",
			Description: "Auction has successfully been bought out!",
		})
		return
	} else if bidAmount > auctionInfo.Bid {
		if bidAmount-auctionInfo.Bid < auctionInfo.MinBid {
			ErrorResponse(s, i, "Bid must be higher than the previous bid by: "+auctionInfo.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.MinBid), "0"), "."))
			return
		}

		if bidAmount-auctionInfo.Bid > auctionInfo.MaxBid && auctionInfo.MaxBid != 0 {
			ErrorResponse(s, i, "Bid must be no more than "+auctionInfo.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.MaxBid), "0"), ".")+" Higher than the previous bid.")
			return
		}

		auctionInfo.Bid = bidAmount
		auctionInfo.Winner = i.Member.User.ID
		Winner := auctionInfo.Winner

		database.DB.Model(&auctionInfo).Updates(auctionInfo)

		updateAuction, err := s.ChannelMessage(auctionInfo.ChannelID, auctionInfo.MessageID)
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

		if antiSnipeFlag {
			updateAuction.Embeds[0].Fields = append(updateAuction.Embeds[0].Fields, &discordgo.MessageEmbedField{
				Name:   "**Anti-Snipe Activated!**",
				Value:  fmt.Sprintf("New End Time: <t:%d>", auctionInfo.EndTime.Unix()),
				Inline: false,
			})
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
			ID:         auctionInfo.MessageID,
			Channel:    auctionInfo.ChannelID,
		})
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}
		Content = "Bid has successfully been placed"
	} else {
		ErrorResponse(s, i, "You must bid higher than: "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Bid), "0"), "."))
	}

	err := SuccessResponse(s, i, PresetResponse{
		Title:  Content,
		Fields: responseFields,
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
	imageURL := "https://i.imgur.com/9wo7diC.png"

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

	if auctionInfo.EndTime.After(time.Now()) {
		time.Sleep(time.Until(auctionInfo.EndTime))
		AuctionEnd(ChannelID, GuildID)
		return
	}

	message := discordgo.NewMessageEdit(auctionInfo.ChannelID, auctionInfo.MessageID)
	messageEmbeds, err := Session.ChannelMessage(auctionInfo.ChannelID, auctionInfo.MessageID)
	if err != nil {
		fmt.Println(err)
		result = database.DB.Delete(&auctionInfo, ChannelID)
		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		return
	}

	if message != nil {
		message.Embeds = append(messageEmbeds.Embeds, &discordgo.MessageEmbed{
			Title:       "Auction has ended!",
			Description: "Thank you for participating!",
			Color:       0x32CD32,
			Image: &discordgo.MessageEmbedImage{
				URL: "https://c.tenor.com/MvFFZxXwdpwAAAAC/sold-ray.gif",
			},
		})
		message.Components = []discordgo.MessageComponent{
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
	}

	if auctionInfo.ImageURL != "" {
		imageURL = auctionInfo.ImageURL
	}

	description := fmt.Sprintf("**Item:** %s", auctionInfo.Item)
	if auctionInfo.Description != "" {
		description += fmt.Sprintf("\n**Description:** %s", auctionInfo.Description)
	}
	if auctionInfo.Winner == "" {
		auctionInfo.Winner = "No Bidders"
	} else if auctionInfo.Winner != "No bidders" {
		user, err := Session.User(strings.Trim(auctionInfo.Winner, "<@!>"))
		if err != nil {
			fmt.Println(err)
		}
		username = fmt.Sprintf("(%s#%s)", user.Username, user.Discriminator)
		auctionInfo.Winner = "<@" + auctionInfo.Winner + ">"
	}

	fmt.Println("Target Price: ", auctionInfo.TargetPrice)

	auctionWinner := fmt.Sprintf("The host had set a target price of %s that has not been reached.", fmt.Sprintf("%s %s", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.TargetPrice), "0"), ".")))
	if auctionInfo.TargetPrice <= auctionInfo.Bid {
		auctionWinner = fmt.Sprintf("%s %s", auctionInfo.Winner, username)
	}

	finalBid := fmt.Sprintf("%s %s\n\u200b", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Bid), "0"), "."))

	if auctionInfo.Bid >= auctionInfo.Buyout && auctionInfo.Buyout != 0 {
		finalBid = fmt.Sprintf("%s %s\n\u200b", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Buyout), "0"), ".")) + " BUYOUT!"
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
					Value:  auctionWinner,
					Inline: true,
				},
				{
					Name:   "**Final Bid**",
					Value:  finalBid,
					Inline: true,
				},
				{
					Name:   "**Auction Host**",
					Value:  fmt.Sprintf("This Auction was hosted by: <@!%s>", auctionInfo.Host),
					Inline: false,
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: imageURL,
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

func AuctionQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var AuctionQueueInfo []database.AuctionQueue
	var AuctionQueue []database.AuctionQueue
	var fields []*discordgo.MessageEmbedField
	var selectOptions []discordgo.SelectMenuOption
	var marker string

	database.DB.Find(&AuctionQueueInfo)

	for _, v := range AuctionQueueInfo {
		if v.GuildID == i.GuildID {
			AuctionQueue = append(AuctionQueue, v)
		}
	}

	sort.Slice(AuctionQueue, func(i, j int) bool { return AuctionQueue[i].StartTime.Unix() < AuctionQueue[j].StartTime.Unix() })

	for _, v := range AuctionQueue {
		if len(fields) == 0 {
			marker = "⭐ "
		} else {
			marker = ""
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("**%s%s**", marker, v.Item),
			Value: fmt.Sprintf("**Start time:** <t:%d:R>\n**End Time:** <t:%d>\n**Starting Price:** %s %s\n\u200b", v.StartTime.Unix(), v.EndTime.Unix(), v.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", v.Bid), "0"), ".")),
		})
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label:       v.Item,
			Value:       fmt.Sprint(v.ID),
			Description: fmt.Sprintf("Start time: %s", v.StartTime.String()),
		})
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "delete_auction_queue",
					Placeholder: "💣 Remove auction from queue",
					MinValues:   0,
					MaxValues:   len(AuctionQueue),
					Options:     selectOptions,
				},
			},
		},
	}

	if len(AuctionQueue) == 0 {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**No Scheduled Auctions**",
			Value:  "Use the `schedule` parameter when creating auctions to plan them in advance!",
			Inline: false,
		})
		components = []discordgo.MessageComponent{}
	}

	err := SuccessResponse(s, i, PresetResponse{
		Title:       "**Auction Queue**",
		Description: "Displays upcoming auctions!",
		Fields:      fields,
		Components:  components,
	})

	if err != nil {
		fmt.Println(err)
	}
}

func AuctionEndButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	guildInfo := database.GuildInfo{}
	auctionInfo := database.Auction{}

	database.DB.First(&auctionInfo, i.ChannelID)
	database.DB.First(&guildInfo, i.GuildID)

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != auctionInfo.Host {
		ErrorResponse(s, i, "You must have an administrator role to end the auction early!")
		return
	}

	if guildInfo.LogChannel == "" {
		fmt.Println("Log channel has not been set for guild: " + i.GuildID)
		ErrorResponse(s, i, "Auction cannot end because log channel has not been set. Please setup an auction log using `/auction setup`")
		return
	}

	err := SuccessResponse(s, i, PresetResponse{
		Title:       "End Auction",
		Description: "Auction Ending...",
	})

	if err != nil {
		fmt.Println(err)
	}

	auctionInfo.EndTime = time.Now()

	database.DB.Model(&auctionInfo).Updates(auctionInfo)

	AuctionEnd(i.ChannelID, i.GuildID)
}

func ClaimPrizeButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var guildInfo database.GuildInfo
	database.DB.First(&guildInfo, i.GuildID)

	if guildInfo.Claiming == "" {
		guildInfo.Claiming = "The discord owner has not set a claiming message. Common ways to claim include: Opening a ticket or contacting the auction host. \nTo customize this message, use the command: `/auction setup claiming:`."
	}

	err := SuccessResponse(s, i, PresetResponse{
		Title:       "Claim Prize",
		Description: guildInfo.Claiming,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func ClearAuctionButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var auctionInfo database.Auction

	result := database.DB.First(&auctionInfo, i.ChannelID)
	if result.Error != nil {
		fmt.Println("Error checking for auction. Auction has likely ended")
		auctionInfo.Host = ""
	}

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != auctionInfo.Host {
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

		for _, v := range messages {
			if !v.Pinned {
				messageIDs = append(messageIDs, v.ID)
			}
		}

		if len(messageIDs) == 0 {
			break
		}

		err = s.ChannelMessagesBulkDelete(i.ChannelID, messageIDs)
		if err != nil {
			fmt.Println(err)
			ErrorResponse(s, i, err.Error())
			return
		}
	}

	err := SuccessResponse(s, i, PresetResponse{
		Title:       "Success!",
		Description: "Chat has been cleared",
	})
	if err != nil {
		fmt.Println(err)
	}
}

func DeleteAuctionQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&(1<<3) != 8 {
		ErrorResponse(s, i, "User have administrator permissions to run this command")
		return
	}

	IDs := i.MessageComponentData().Values

	for _, v := range IDs {
		database.DB.Delete(database.AuctionQueue{}, v)
	}

	AuctionQueue(s, i)
}

func ParseTime(inputDuration string) (time.Duration, error) {
	if strings.HasSuffix(strings.ToLower(inputDuration), "d") {
		inputDuration = strings.TrimSuffix(strings.ToLower(inputDuration), "d")
		float, err := strconv.ParseFloat(inputDuration, 64)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		inputDuration = fmt.Sprint(float*24) + "h"
	}

	return time.ParseDuration(inputDuration)

}
