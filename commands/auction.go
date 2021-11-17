package commands

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
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
						5,
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
					//Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "snipe_range",
					Description: "Set 0 to disable. The remaining time needed to activate Anti-Snipe (Example: 24h, or 1d)",
					//Autocomplete: true,
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
					//Autocomplete: true,
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
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "schedule",
					Description: "Set how long until the auction starts. (Example: 24h, or 1d)",
					Required:    false,
					//Autocomplete: true,
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
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "setup_clear",
			Description: "Clear one or more of the settings in /auction setup.",
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
	case "setup_clear":
		AuctionSetupClear(s, i)
	}
}

func AuctionAutoComplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch i.ApplicationCommandData().Options[0].Name {
	case "create":
		options := h.ParseSubCommand(i)

		if options["duration"] != nil {
			choices = TimeSuggestions(options["duration"].(string))
		} else if options["schedule"] != nil {
			choices = TimeSuggestions(options["schedule"].(string))
		} else if options["snipe_range"] != nil {
			choices = TimeSuggestions(options["snipe_range"].(string))
		} else if options["snipe_extension"] != nil {
			choices = TimeSuggestions(options["snipe_extension"].(string))
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

	input = strings.Trim(input, "dhms")

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
	h.ErrorResponse(s, i, "Help command has not been setup yet")
}

func AuctionSetup(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&(1<<3) != 8 {
		h.ErrorResponse(s, i, "User must have administrator permissions to run this command")
		return
	}

	var err error

	options := h.ParseSubCommand(i)
	content := ""

	info := database.AuctionSetup{
		GuildID: i.GuildID,
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&info)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	for key := range options {
		content += fmt.Sprintf("• %s has been successfully set.\n", strings.Title(strings.ReplaceAll(key, "_", " ")))

		switch key {
		case "snipe_extension", "snipe_range":
			options[key], err = h.ParseTime(options[key].(string))
		}
		if err != nil {
			fmt.Println(err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

	}

	database.DB.Model(&info).Updates(options)

	//Now check what options are set
	setOptions := map[string]interface{}{}

	database.DB.Model(database.AuctionSetup{}).First(&setOptions, i.GuildID)

	antiSnipeDescription := fmt.Sprintf("If a bid is placed within %s of the auction ending, it will be extended by %s.", info.SnipeRange.String(), info.SnipeExtension.String())

	responseFields := []*discordgo.MessageEmbedField{}

	for _, v := range AuctionCommand.Options[1].Options {
		if !strings.Contains(v.Name, "snipe") {
			switch {
			case setOptions[v.Name] == "":
				setOptions[v.Name] = "Not Set"
			case strings.Contains(v.Name, "role"):
				setOptions[v.Name] = fmt.Sprintf("<@&%s>", setOptions[v.Name])
			case strings.Contains(v.Name, "channel"):
				setOptions[v.Name] = fmt.Sprintf("<#%s>", setOptions[v.Name])
			case strings.Contains(v.Name, "category"):
				category, err := s.Channel(info.Category)
				if err != nil {
					fmt.Println("Category Error:", err)
					setOptions[v.Name] = err.Error()
				} else {
					setOptions[v.Name] = category.Name
				}
			}
			responseFields = append(responseFields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("**%s**", strings.Title(strings.ReplaceAll(v.Name, "_", " "))),
				Value: setOptions[v.Name].(string),
			})
		}
	}

	if setOptions["snipe_range"].(time.Duration) == 0 || setOptions["snipe_extension"].(time.Duration) == 0 {
		antiSnipeDescription = "Anti Snipe Disabled. To enable, set both snipe_extension and snipe_range"
	}

	responseFields = append(responseFields, &discordgo.MessageEmbedField{
		Name:  "**Anti Snipe**",
		Value: antiSnipeDescription,
	})

	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Auction Setup",
		Description: content,
		Fields:      responseFields,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionSetupClear(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.Permissions&(1<<3) != 8 {
		h.ErrorResponse(s, i, "User must have administrator permissions to run this command")
		return
	}

	menuOptions := []discordgo.SelectMenuOption{}

	for _, v := range AuctionCommand.Options[1].Options {
		menuOptions = append(menuOptions, discordgo.SelectMenuOption{
			Label:       strings.Title(strings.ReplaceAll(v.Name, "_", " ")),
			Value:       v.Name,
			Description: v.Description,
		})
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Clear Auction Setup**",
		Description: "Please select which setup options you would like to clear",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "clear_auction_setup",
						Placeholder: "Select options to clear!",
						MinValues:   1,
						MaxValues:   len(AuctionCommand.Options[1].Options),
						Options:     menuOptions,
					},
				},
			},
		},
	})
}

func AuctionCreate(s *discordgo.Session, auctionInfo database.AuctionQueue) {

	AuctionSetup := database.AuctionSetup{
		GuildID: auctionInfo.GuildID,
	}
	database.DB.First(&AuctionSetup, auctionInfo.GuildID)
	currency := auctionInfo.Currency
	incCurrency := "+"
	if currency != "" {
		incCurrency = currency
	}

	guild, err := s.Guild(auctionInfo.GuildID)

	if err != nil {
		fmt.Println(err)
	}

	channel, err := s.GuildChannelCreateComplex(auctionInfo.GuildID, discordgo.GuildChannelCreateData{
		Name:     "💸│" + auctionInfo.Item,
		Type:     0,
		ParentID: auctionInfo.Category,
	})

	if err != nil {
		fmt.Println("Error, Bot trying to start a no category auction: " + err.Error())
		channel, err = s.GuildChannelCreateComplex(auctionInfo.GuildID, discordgo.GuildChannelCreateData{
			Name:     "💸│" + auctionInfo.Item,
			Type:     0,
			ParentID: "",
		})

		if err != nil {
			fmt.Println("Bot failed starting a no category auction: ", err.Error())
			return
		}
	}

	var message *discordgo.Message
	auctionfields := []*discordgo.MessageEmbedField{
		{
			Name:   "__**Hosted By**__",
			Value:  fmt.Sprintf("<@%s>", auctionInfo.Host),
			Inline: false,
		},
		{
			Name:   "__**End Time**__",
			Value:  fmt.Sprintf("<t:%d:R>", auctionInfo.EndTime.Unix()),
			Inline: false,
		},
	}

	if auctionInfo.Description != "" {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Item Description**__",
			Value:  auctionInfo.Description,
			Inline: false,
		})
	}

	if auctionInfo.MinBid != 0 || auctionInfo.MaxBid != 0 {
		bidRanges := ""
		if auctionInfo.MinBid != 0 {
			bidRanges += fmt.Sprintf("• Minimum %s %s above previous bid\n", incCurrency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.MinBid), "0"), "."))
		}
		if auctionInfo.MaxBid != 0 {
			bidRanges += fmt.Sprintf("• Maximum %s %s above previous bid", incCurrency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.MaxBid), "0"), "."))
		}
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Bid Range**__",
			Value:  bidRanges,
			Inline: false,
		})
	}

	if auctionInfo.TargetPrice != 0 {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Target Price**__",
			Value:  "The host has set a hidden target price for this auction.",
			Inline: false,
		})
	}

	if AuctionSetup.SnipeExtension != 0 && AuctionSetup.SnipeRange != 0 {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Anti-Snipe**__",
			Value:  fmt.Sprintf("Anti-Snipe has been enabled for this auction. If a bid is placed within the last %s, the auction will be extended by %s.", AuctionSetup.SnipeRange.String(), AuctionSetup.SnipeExtension.String()),
			Inline: false,
		})
	}

	if auctionInfo.Buyout != 0 {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Buyout Price**__",
			Value:  fmt.Sprintf("%s %s", incCurrency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Buyout), "0"), ".")),
			Inline: false,
		})
	}

	auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
		Name:   "__**Current Highest Bid:**__",
		Value:  fmt.Sprintf("%s %s", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Bid), "0"), ".")),
		Inline: true,
	}, &discordgo.MessageEmbedField{
		Name:   "__**Current Winner**__",
		Value:  fmt.Sprintf("<@%s>", auctionInfo.Host),
		Inline: true,
	}, &discordgo.MessageEmbedField{
		Name:   "__**How to Bid**__",
		Value:  "Use the command `/auction bid` below.\n• Ex: `/auction bid 550`",
		Inline: false,
	})

	message, err = h.PresetMessageSend(s, channel.ID, h.PresetResponse{
		Content: fmt.Sprintf("<@&%s>", AuctionSetup.AlertRole),
		Title:   fmt.Sprintf("Auction Item: __**%s**__", auctionInfo.Item),
		Fields:  auctionfields,
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
		if strings.Contains(err.Error(), "http") {
			message, err = h.PresetMessageSend(s, channel.ID, h.PresetResponse{
				Content: fmt.Sprintf("<@&%s>", AuctionSetup.AlertRole),
				Title:   "__**" + auctionInfo.Item + "**__",
				Fields:  auctionfields,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: guild.IconURL(),
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
		}
		fmt.Println(err.Error())
		h.ErrorMessage(s, channel.ID, err.Error())
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

	options := h.ParseSubCommand(i)

	endTimeDuration, err := h.ParseTime(strings.ToLower(options["duration"].(string)))
	if err != nil {
		fmt.Println(err)
		h.ErrorResponse(s, i, err.Error())
		return
	}

	item := options["item"].(string)
	var image string
	var maxBid float64
	var minBid float64
	var buyout float64
	var targetPrice float64
	var description string
	info := database.AuctionSetup{
		GuildID: i.GuildID,
	}

	if len(item) > 100 {
		h.ErrorResponse(s, i, "Title cannot be over 100 characters long")
		return
	}

	database.DB.First(&info, i.GuildID)

	canHost := false

	if info.HostRole != "" {
		for _, v := range i.Member.Roles {
			if v == info.HostRole {
				canHost = true
			}
		}
		if i.Member.Permissions&(1<<3) == 8 {
			canHost = true
		}
		if !canHost {
			h.ErrorResponse(s, i, "User must be administrator or have the role <@&"+info.HostRole+"> to host auctions.")
			return
		}
	}

	currency := info.Currency

	if options["category"] != nil {
		info.Category = options["category"].(string)
	}
	if options["description"] != nil {
		description = options["description"].(string)
	}
	if options["image"] != nil {
		image = options["image"].(string)
		if !strings.Contains(image, "http") {
			h.ErrorResponse(s, i, "Image must be an http or https link.")
			fmt.Println("Image must be an http or https link. Entered image: " + image)
			return
		}
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
			h.ErrorResponse(s, i, "User have administrator permissions to schedule auctions in advance. If you are not an admin you can still run auctions normally without the `schedule` parameter.")
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
			h.ErrorResponse(s, i, "You can only schedule 25 auctions in advance.")
			return
		}

		startTimeDuration, err := h.ParseTime(strings.ToLower(options["schedule"].(string)))
		if err != nil {
			fmt.Println(err)
			h.ErrorResponse(s, i, err.Error())
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
		Category:    info.Category,
		Buyout:      buyout,
		TargetPrice: targetPrice,
	}

	if options["schedule"] != nil {

		database.DB.Create(&auctionData)
		err = h.PremiumResponse(s, i, h.PresetResponse{
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
		err = h.SuccessResponse(s, i, h.PresetResponse{
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
		options = h.ParseSubCommand(i)
	} else {
		options = h.ParseSlashCommand(i)
	}

	bidAmount := options["amount"].(float64)
	var auctionInfo database.Auction
	var AuctionSetup database.AuctionSetup
	auctionInfo.ChannelID = i.ChannelID
	database.DB.First(&auctionInfo, i.ChannelID)
	database.DB.First(&AuctionSetup, i.GuildID)
	currency := AuctionSetup.Currency
	var Content string
	var antiSnipeFlag bool

	if auctionInfo.Currency != "" {
		currency = auctionInfo.Currency
	}

	if auctionInfo.EndTime.Before(time.Now()) {
		h.ErrorResponse(s, i, "Cannot Bid, Auction has ended")
		return
	}

	if i.Member.User.ID == auctionInfo.Winner && i.Member.User.ID != "280812467775471627" && auctionInfo.MaxBid != 0 {
		h.ErrorResponse(s, i, "Cannot out bid yourself on a capped bid auction!")
		return
	}

	if time.Until(auctionInfo.EndTime) < AuctionSetup.SnipeRange && AuctionSetup.SnipeExtension != 0 {
		auctionInfo.EndTime = auctionInfo.EndTime.Add(AuctionSetup.SnipeExtension)
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

		err := h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "Success!",
			Description: "Auction has successfully been bought out!",
		})

		AuctionEnd(i.ChannelID, i.GuildID)

		if err != nil {
			fmt.Println(err.Error())
		}
		return
	} else if bidAmount > auctionInfo.Bid {
		if bidAmount-auctionInfo.Bid < auctionInfo.MinBid {
			h.ErrorResponse(s, i, "Bid must be higher than the previous bid by: "+auctionInfo.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.MinBid), "0"), "."))
			return
		}

		if bidAmount-auctionInfo.Bid > auctionInfo.MaxBid && auctionInfo.MaxBid != 0 {
			h.ErrorResponse(s, i, "Bid must be no more than "+auctionInfo.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.MaxBid), "0"), ".")+" Higher than the previous bid.")
			return
		}

		auctionInfo.Bid = bidAmount
		auctionInfo.Winner = i.Member.User.ID
		Winner := auctionInfo.Winner

		database.DB.Model(&auctionInfo).Updates(auctionInfo)

		updateAuction, err := s.ChannelMessage(auctionInfo.ChannelID, auctionInfo.MessageID)
		if err != nil {
			fmt.Println(err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

		bidHistory := ""
		username := i.Member.Nick
		if username == "" {
			username = i.Member.User.Username
		}

		if len(updateAuction.Embeds) == 2 {
			bidHistory = updateAuction.Embeds[1].Description + "\n-> " + username + ": " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), ".")
		} else {
			bidHistory = "-> " + username + ": " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), ".")
		}

		if len(strings.ReplaceAll(bidHistory, " ", "")) >= 4096 {
			bidHistory = "BidHistory was too long and has been reset to prevent a crash.\n-> " + username + ": " + fmt.Sprint(bidAmount)
		}

		for n, v := range updateAuction.Embeds[0].Fields {
			switch v.Name {
			case "__**Current Highest Bid:**__":
				updateAuction.Embeds[0].Fields[n].Value = fmt.Sprintf("%s %s\n\u200b", currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), "."))
			case "__**Starting Bid:**__":
				updateAuction.Embeds[0].Fields[n].Value = fmt.Sprintf("%s %s\n\u200b", currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), "."))
				updateAuction.Embeds[0].Fields[n].Name = "__**Current Highest Bid:**__"
			case "__**Current Winner**__":
				updateAuction.Embeds[0].Fields[n].Value = fmt.Sprintf("<@%s>", Winner)
			case "__**End Time**__":
				if antiSnipeFlag {
					updateAuction.Embeds[0].Fields[n].Value = fmt.Sprintf("New End Time: <t:%d:R>", auctionInfo.EndTime.Unix())
				}
			}
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
			h.ErrorResponse(s, i, err.Error())
			return
		}
		Content = "Bid has successfully been placed"
	} else {
		h.ErrorResponse(s, i, "You must bid higher than: "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Bid), "0"), "."))
		return
	}

	err := h.SuccessResponse(s, i, h.PresetResponse{
		Title:  Content,
		Fields: responseFields,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionEnd(ChannelID, GuildID string) {
	var auctionInfo database.Auction
	var AuctionSetup database.AuctionSetup
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

	result = database.DB.First(&AuctionSetup, GuildID)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	if auctionInfo.Bid < auctionInfo.Buyout || auctionInfo.Buyout == 0 {
		if auctionInfo.EndTime.After(time.Now()) {
			time.Sleep(time.Until(auctionInfo.EndTime))
			AuctionEnd(ChannelID, GuildID)
			return
		}
	}

	message := discordgo.NewMessageEdit(ChannelID, auctionInfo.MessageID)
	messageEmbeds, err := Session.ChannelMessage(auctionInfo.ChannelID, auctionInfo.MessageID)
	if err != nil {
		fmt.Println(err)
		result = database.DB.Delete(&auctionInfo, ChannelID)
		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		return
	}

	if AuctionSetup.LogChannel == "" {
		fmt.Println("Log channel has not been set for guild: " + GuildID)
		_, err := h.ErrorMessage(Session, ChannelID, "Auction cannot end because log channel has not been set. Please setup an auction log using `/auction setup`. You might need to end the auction manually after setting the channel.")
		if err != nil {
			fmt.Println(err)
			return
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

	auctionWinner := fmt.Sprintf("The host had set a target price of %s that has not been reached.", fmt.Sprintf("%s %s", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.TargetPrice), "0"), ".")))
	if auctionInfo.TargetPrice <= auctionInfo.Bid {
		auctionWinner = fmt.Sprintf("%s %s", auctionInfo.Winner, username)
	}

	finalBid := fmt.Sprintf("%s %s\n\u200b", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Bid), "0"), "."))

	if auctionInfo.Bid >= auctionInfo.Buyout && auctionInfo.Buyout != 0 {
		finalBid = fmt.Sprintf("%s %s\n\u200b", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Buyout), "0"), ".")) + " BUYOUT!"
	}

	messageSend := discordgo.MessageSend{
		Content: auctionInfo.Winner,
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

	_, err = Session.ChannelMessageSendComplex(AuctionSetup.LogChannel, &messageSend)
	if err != nil {
		fmt.Println(err)
		messageSend.Embed.Image.URL = ""
		_, err = Session.ChannelMessageSendComplex(AuctionSetup.LogChannel, &messageSend)
		if err != nil {
			fmt.Println(err)
			h.ErrorMessage(Session, ChannelID, err.Error())
			return
		}
		h.ErrorMessage(Session, ChannelID, err.Error())
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

	err := h.SuccessResponse(s, i, h.PresetResponse{
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

	AuctionSetup := database.AuctionSetup{}
	auctionInfo := database.Auction{}

	database.DB.First(&auctionInfo, i.ChannelID)
	database.DB.First(&AuctionSetup, i.GuildID)

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != auctionInfo.Host {
		h.ErrorResponse(s, i, "You must have an administrator role to end the auction early!")
		return
	}

	if AuctionSetup.LogChannel == "" {
		fmt.Println("Log channel has not been set for guild: " + i.GuildID)
		h.ErrorResponse(s, i, "Auction cannot end because log channel has not been set. Please setup an auction log using `/auction setup`")
		return
	}

	err := h.SuccessResponse(s, i, h.PresetResponse{
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
	var AuctionSetup database.AuctionSetup
	database.DB.First(&AuctionSetup, i.GuildID)

	if AuctionSetup.Claiming == "" {
		AuctionSetup.Claiming = "The discord owner has not set a claiming message. Common ways to claim include: Opening a ticket or contacting the auction host. \nTo customize this message, use the command: `/auction setup claiming:`."
	}

	err := h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Claim Prize",
		Description: AuctionSetup.Claiming,
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
		h.ErrorResponse(s, i, "User must be host or have administrator permissions to run this command")
		fmt.Println("User must be host or have administrator permissions to run this command")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})

	for {
		messageIDs := make([]string, 0)
		messages, err := s.ChannelMessages(i.ChannelID, 100, "", i.Message.ID, "")
		if err != nil {
			fmt.Println(err)
			h.DeferredErrorResponse(s, i, err.Error())
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
			h.DeferredErrorResponse(s, i, err.Error())
			return
		}
	}

	_, err := h.DeferredResponse(s, i, h.PresetResponse{
		Title:       "Success!",
		Description: "Clearing Chat",
	})

	if err != nil {
		fmt.Println(err)
	}
}

func DeleteAuctionQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&(1<<3) != 8 {
		h.ErrorResponse(s, i, "User have administrator permissions to run this command")
		return
	}

	IDs := i.MessageComponentData().Values

	for _, v := range IDs {
		database.DB.Delete(database.AuctionQueue{}, v)
	}

	AuctionQueue(s, i)
}

func AuctionSetupClearButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.MessageComponentData().Values

	clearedValues := map[string]interface{}{}

	clearedValues["guild_id"] = i.GuildID
	info := database.AuctionSetup{
		GuildID: i.GuildID,
	}

	clearedSettings := "No Settings Cleared!"
	if len(options) > 0 {
		clearedSettings = ""
	}

	for _, v := range options {
		switch v {
		case "snipe_extension", "snipe_range":
			clearedValues[v] = 0
		default:
			clearedValues[v] = ""
		}
		clearedSettings += fmt.Sprintf("• %s\n", strings.Title(strings.ReplaceAll(v, "_", " ")))
	}

	database.DB.Model(&info).Updates(clearedValues)

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Cleared Auction Settings**",
		Description: "You have successfully cleared the following settings. Run `/auction setup` to see your changes.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Cleared Settings**",
				Value: clearedSettings,
			},
		},
	})
}
