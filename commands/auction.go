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
		/*{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "help",
			Description: "auction info",
		},*/
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
					Description: "Set a role to get pinged whenever an auction starts.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "host_role",
					Description: "Set a role that can host auctions.",
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
				{
					Type:        10,
					Name:        "quick_bid",
					Description: "Set the value of the quickbid reaction.",
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
					Name:        "bid",
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
					Name:        "image_url",
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
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "edit",
			Description: "Edit any auction details",
			Options: []*discordgo.ApplicationCommandOption{
				/*{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "queue_number",
					Description: "The number of the auction if you are editing one in queue.",
					Required:    false,
					//Autocomplete: true,
				},*/
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "Change the auction item",
					Required:    false,
					//Autocomplete: true,
				},
				{
					Type:        10,
					Name:        "bid",
					Description: "Change the original bid, or edit the current bid",
					Required:    false,
					//Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "extend",
					Description: "Extend the length of the auction, use a negative value to reduce the time. (Example: 24h, or 1d)",
					Required:    false,
					//Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "Change the description",
					Required:    false,
					//Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "currency",
					Description: "Change the currency",
					Required:    false,
					//Autocomplete: true,
				},
				{
					Type:        10,
					Name:        "increment_max",
					Description: "The max amount someone can bid at once",
					Required:    false,
					//Autocomplete: true,
				},
				{
					Type:        10,
					Name:        "increment_min",
					Description: "The minimum amount someone can bid at once",
					Required:    false,
					//Autocomplete: true,
				},
				{
					Type:        10,
					Name:        "buyout",
					Description: "Edit the buyout price",
					Required:    false,
					//Autocomplete: true,
				},
				{
					Type:        10,
					Name:        "target_price",
					Description: "Edit the hidden target price",
					Required:    false,
					//Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "image_url",
					Description: "Edit the image link",
					Required:    false,
					//Autocomplete: true,
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
		AuctionPlanner(s, i)
	case "bid":
		AuctionBid(s, i)
	case "queue":
		AuctionQueue(s, i)
	case "edit":
		AuctionEdit(s, i)
	}
}

func AuctionFormat(s *discordgo.Session, auctionInfo database.Auction) discordgo.MessageEmbed {

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

	if auctionInfo.IncrementMin != 0 || auctionInfo.IncrementMax != 0 {
		bidRanges := ""
		if auctionInfo.IncrementMin != 0 {
			bidRanges += fmt.Sprintf("• Minimum %s %s above previous bid\n", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.IncrementMin), "0"), "."))
		}
		if auctionInfo.IncrementMax != 0 {
			bidRanges += fmt.Sprintf("• Maximum %s %s above previous bid", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.IncrementMax), "0"), "."))
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

	if auctionInfo.SnipeExtension != 0 && auctionInfo.SnipeRange != 0 {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Anti-Snipe**__",
			Value:  fmt.Sprintf("Anti-Snipe has been enabled for this auction. If a bid is placed within the last %s, the auction will be extended by %s.", auctionInfo.SnipeRange.String(), auctionInfo.SnipeExtension.String()),
			Inline: false,
		})
	}

	if auctionInfo.Buyout != 0 {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Buyout Price**__",
			Value:  fmt.Sprintf("%s %s", auctionInfo.Currency, strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.Buyout), "0"), ".")),
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
		Value:  "Use the command `/auction bid` below.\n• Ex: `/auction bid 550`.\n**Alternate Method:** reply to this auction or @ the bot with `bid <amount>`",
		Inline: false,
	})

	auctionEmbed := discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Auction Item: __**%s**__", auctionInfo.Item),
		Description: "",
		Fields:      auctionfields,
	}

	return auctionEmbed
}

func AuctionAutoComplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var choices []*discordgo.ApplicationCommandOptionChoice
	var focusedData discordgo.ApplicationCommandInteractionDataOption
	var choiceName string

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
	case "edit":
		currentValues := map[string]interface{}{}

		for _, v := range i.ApplicationCommandData().Options[0].Options {
			if v.Focused {
				focusedData = *v
			}
		}

		database.DB.Model(database.Auction{}).First(&currentValues, i.ChannelID)

		name := focusedData.Name

		if currentValues[name] != nil || name == "extend" {
			switch name {
			case "extend":
				choiceName = fmt.Sprintf("Time remaining: %s", time.Until(currentValues["end_time"].(time.Time)))
			default:
				choiceName = fmt.Sprintf("Current Value: %v", currentValues[name])
			}
		} else {
			choiceName = "Not Configured"
		}

		choices = []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  choiceName,
				Value: focusedData.Value,
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
		fmt.Println("Response Error:", err)
		return
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

	var err error

	options := h.ParseSubCommand(i)
	content := ""

	if i.Member.Permissions&(1<<3) != 8 {
		h.ErrorResponse(s, i, "User must have administrator permissions to run this command")
		return
	}

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
			case setOptions[v.Name] == "", setOptions[v.Name] == 0, setOptions[v.Name] == nil:
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
			case v.Name == "quick_bid":
				setOptions[v.Name] = fmt.Sprintf("%s %s", setOptions["currency"], strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", setOptions[v.Name]), "0"), "."))
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

	menuOptions := []discordgo.SelectMenuOption{}

	for _, v := range AuctionCommand.Options[1].Options {
		menuOptions = append(menuOptions, discordgo.SelectMenuOption{
			Label:       strings.Title(strings.ReplaceAll(v.Name, "_", " ")),
			Value:       v.Name,
			Description: v.Description,
		})
	}

	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Auction Setup",
		Description: content,
		Fields:      responseFields,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "clear_auction_setup",
						Placeholder: "Clear Setup Options",
						MinValues:   1,
						MaxValues:   len(AuctionCommand.Options[1].Options),
						Options:     menuOptions,
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionCreate(s *discordgo.Session, auctionInfo database.AuctionQueue) {

	AuctionSetup := database.AuctionSetup{
		GuildID: auctionInfo.GuildID,
	}
	database.DB.First(&AuctionSetup, auctionInfo.GuildID)

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
	auctionfields := AuctionFormat(s, database.Auction{
		Bid:          auctionInfo.Bid,
		EndTime:      auctionInfo.EndTime,
		Winner:       auctionInfo.Host,
		GuildID:      auctionInfo.GuildID,
		Item:         auctionInfo.Item,
		Host:         auctionInfo.Host,
		Currency:     auctionInfo.Currency,
		IncrementMin: auctionInfo.IncrementMin,
		IncrementMax: auctionInfo.IncrementMax,
		Description:  auctionInfo.Description,
		ImageURL:     auctionInfo.ImageURL,
		TargetPrice:  auctionInfo.TargetPrice,
		Buyout:       auctionInfo.Buyout,
	}).Fields

	alertRole := ""
	if AuctionSetup.AlertRole != "" {
		alertRole = fmt.Sprintf("<@&%s>", AuctionSetup.AlertRole)
	}

	message, err = h.PresetMessageSend(s, channel.ID, h.PresetResponse{
		Content: alertRole,
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
		ChannelID:    message.ChannelID,
		Bid:          auctionInfo.Bid,
		MessageID:    message.ID,
		EndTime:      auctionInfo.EndTime,
		Winner:       auctionInfo.Host,
		GuildID:      auctionInfo.GuildID,
		Item:         auctionInfo.Item,
		Host:         auctionInfo.Host,
		Currency:     auctionInfo.Currency,
		IncrementMin: auctionInfo.IncrementMin,
		IncrementMax: auctionInfo.IncrementMax,
		Description:  auctionInfo.Description,
		ImageURL:     auctionInfo.ImageURL,
		TargetPrice:  auctionInfo.TargetPrice,
		Buyout:       auctionInfo.Buyout,
	})

	database.DB.Delete(auctionInfo, auctionInfo.ID)

	time.Sleep(time.Until(auctionInfo.EndTime))
	AuctionEnd(channel.ID, auctionInfo.GuildID)
}

func AuctionEdit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := h.ParseSubCommand(i)
	auctionInfo := database.Auction{
		ChannelID: i.ChannelID,
	}
	database.DB.First(&auctionInfo)

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != auctionInfo.Host{
		h.ErrorResponse(s, i, "User must have be host or have administrator permissions to run this command")
		return
	}

	for key, value := range options {
		switch key {
		case "extend":
			extraDuration, err := h.ParseTime(strings.ToLower(value.(string)))
			if err != nil {
				fmt.Println(err)
				h.ErrorResponse(s, i, err.Error())
				return
			}
			options["end_time"] = auctionInfo.EndTime.Add(extraDuration)
			delete(options, "extend")
		}
	}

	if options["queue_number"] != nil {

		queueNumber, err := strconv.Atoi(options["queue_number"].(string))
		if err != nil {
			fmt.Println(err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

		delete(options, "queue_number")

		result := database.DB.Model(database.AuctionQueue{
			ID: queueNumber,
		}).Updates(options)

		if result.Error != nil {
			fmt.Println(result.Error)
			h.ErrorResponse(s, i, result.Error.Error())
			return
		}

	} else {

		result := database.DB.Model(database.Auction{
			ChannelID: i.ChannelID,
		}).Updates(options)

		if result.Error != nil {
			fmt.Println(result.Error)
			h.ErrorResponse(s, i, result.Error.Error())
			return
		}

		AuctionSetup := database.AuctionSetup{
			GuildID: auctionInfo.GuildID,
		}

		database.DB.First(&AuctionSetup, auctionInfo.GuildID)

		database.DB.Model(database.Auction{}).First(&auctionInfo, i.ChannelID)

		message, err := s.ChannelMessage(i.ChannelID, auctionInfo.MessageID)
		if err != nil {
			fmt.Println(err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

		if options["item"] != nil {
			message.Embeds[0].Title = fmt.Sprintf("Auction Item: __**%s**__", auctionInfo.Item)
		}

		message.Embeds[0].Fields = AuctionFormat(s, auctionInfo).Fields

		_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content:    &message.Content,
			Components: message.Components,
			Embeds:     message.Embeds,
			ID:         message.ID,
			Channel:    message.ChannelID,
		})
		if err != nil {
			fmt.Println(err)
			h.ErrorResponse(s, i, err.Error())
			return
		}
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Content:     "",
		Title:       "Success",
		Description: "Auction has successfully been edited",
	})
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

		var AuctionQueue []database.AuctionQueue

		database.DB.Where(map[string]interface{}{"guild_id": info.GuildID}).Find(&AuctionQueue)

		if len(AuctionQueue) >= 25 {
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
		Bid:          options["bid"].(float64),
		StartTime:    startTime,
		EndTime:      endTime,
		GuildID:      i.GuildID,
		Item:         item,
		Host:         i.Member.User.ID,
		Currency:     currency,
		IncrementMin: minBid,
		IncrementMax: maxBid,
		Description:  description,
		ImageURL:     image,
		Category:     info.Category,
		Buyout:       buyout,
		TargetPrice:  targetPrice,
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

func AuctionBidAlternate(s *discordgo.Session, bidData database.Auction) {

	var Content string
	var antiSnipeFlag bool
	var responseFields []*discordgo.MessageEmbedField
	var auctionSetup database.AuctionSetup
	var auction database.Auction
	database.DB.First(&auction, bidData.ChannelID)
	database.DB.First(&auctionSetup, bidData.GuildID)
	currency := auctionSetup.Currency

	fmt.Println("bidding")

	if auction.Currency != "" {
		currency = auction.Currency
	}

	if time.Until(auction.EndTime) < auction.SnipeRange && auction.SnipeExtension != 0 {
		auction.EndTime = auction.EndTime.Add(auction.SnipeExtension)
		responseFields = []*discordgo.MessageEmbedField{
			{
				Name:   "**Anti-Snipe Activated!**",
				Value:  fmt.Sprintf("New End Time: <t:%d>", auction.EndTime.Unix()),
				Inline: false,
			},
		}
		antiSnipeFlag = true
	}

	switch {
	case auction.EndTime.Before(time.Now()):
		h.ErrorMessage(s, bidData.ChannelID, "Cannot Bid, Auction has ended")
		return
	case bidData.Winner == auction.Winner && bidData.Winner != "280812467775471627" && auction.IncrementMax != 0:
		h.ErrorMessage(s, bidData.ChannelID, "Cannot out bid yourself on a capped bid auction!")
		return
	case bidData.Bid >= auction.Buyout && auction.Buyout != 0:
		auction.Bid = bidData.Bid
		auction.Winner = bidData.Winner

		database.DB.Model(&auction).Updates(auction)

		_, err := h.SuccessMessage(s, bidData.ChannelID, h.PresetResponse{
			Title:       "Success!",
			Description: "Auction has successfully been bought out!",
		})
		if err != nil {
			fmt.Println(err.Error())
		}

		AuctionEnd(auction.ChannelID, auction.GuildID)

		return
	case bidData.Bid > auction.Bid:
		switch {
		case bidData.Bid-auction.Bid < auction.IncrementMin:
			h.ErrorMessage(s, bidData.ChannelID, "Bid must be higher than the previous bid by: "+auction.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auction.IncrementMin), "0"), "."))
			return
		case bidData.Bid-auction.Bid > auction.IncrementMax && auction.IncrementMax != 0:
			h.ErrorMessage(s, bidData.ChannelID, "Bid must be no more than "+auction.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auction.IncrementMax), "0"), ".")+" Higher than the previous bid.")
			return
		}

		auction.Bid = bidData.Bid
		auction.Winner = bidData.Winner

		database.DB.Model(&auction).Updates(auction)

		updateAuction, err := s.ChannelMessage(auction.ChannelID, auction.MessageID)
		if err != nil {
			fmt.Println(err)
			h.ErrorMessage(s, bidData.ChannelID, err.Error())
			return
		}

		user, err := s.GuildMember(bidData.GuildID, bidData.Winner)
		if err != nil {
			fmt.Println(err)
			h.ErrorMessage(s, bidData.ChannelID, err.Error())
			return
		}

		bidHistory := ""
		username := user.Nick
		if username == "" {
			username = user.User.Username
		}
		bidAmount := bidData.Bid

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
				updateAuction.Embeds[0].Fields[n].Value = fmt.Sprintf("<@%s>", auction.Winner)
			case "__**End Time**__":
				if antiSnipeFlag {
					updateAuction.Embeds[0].Fields[n].Value = fmt.Sprintf("New End Time: <t:%d:R>", auction.EndTime.Unix())
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
			ID:         auction.MessageID,
			Channel:    auction.ChannelID,
		})
		if err != nil {
			fmt.Println(err)
			h.ErrorMessage(s, bidData.ChannelID, err.Error())
			return
		}
		Content = "Bid has successfully been placed"

	default:
		h.ErrorMessage(s, bidData.ChannelID, "You must bid higher than: "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auction.Bid), "0"), "."))
		return
	}

	_, err := h.SuccessMessage(s, bidData.ChannelID, h.PresetResponse{
		Title:  Content,
		Fields: responseFields,
	})
	if err != nil {
		fmt.Println(err)
	}
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

	if i.Member.User.ID == auctionInfo.Winner && i.Member.User.ID != "280812467775471627" && auctionInfo.IncrementMax != 0 {
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
		if bidAmount-auctionInfo.Bid < auctionInfo.IncrementMin {
			h.ErrorResponse(s, i, "Bid must be higher than the previous bid by: "+auctionInfo.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.IncrementMin), "0"), "."))
			return
		}

		if bidAmount-auctionInfo.Bid > auctionInfo.IncrementMax && auctionInfo.IncrementMax != 0 {
			h.ErrorResponse(s, i, "Bid must be no more than "+auctionInfo.Currency+" "+strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionInfo.IncrementMax), "0"), ".")+" Higher than the previous bid.")
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

	database.DB.Find(&AuctionQueueInfo)

	for _, v := range AuctionQueueInfo {
		if v.GuildID == i.GuildID {
			AuctionQueue = append(AuctionQueue, v)
		}
	}

	sort.Slice(AuctionQueue, func(i, j int) bool { return AuctionQueue[i].StartTime.Unix() < AuctionQueue[j].StartTime.Unix() })

	for _, v := range AuctionQueue {

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("**%s%s**", fmt.Sprint(len(fields)+1), v.Item),
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

	if AuctionSetup.LogChannel == "" {
		fmt.Println("Log channel has not been set for guild: " + i.GuildID)
		h.ErrorResponse(s, i, "Auction cannot end because log channel has not been set. Please setup an auction log using `/auction setup`")
		return
	}

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != auctionInfo.Host {
		h.ErrorResponse(s, i, "User must have administrator permissions to run this command")
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
