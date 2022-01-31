package commands

import (
	"fmt"
	"math"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var Session *discordgo.Session

var AuctionCommand = discordgo.ApplicationCommand{
	Name:        "auction",
	Description: "Put an item up for auction!",
	Options: []*discordgo.ApplicationCommandOption{
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
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "currency_side",
					Description: "Left/Right currency",
					//Autocomplete: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Left",
							Value: "left",
						},
						{
							Name:  "Right",
							Value: "right",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "integer_only",
					Description: "Only allow integer bids (no decimals).",
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
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "currency_side",
					Description: "Left/Right currency",
					//Autocomplete: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Left",
							Value: "left",
						},
						{
							Name:  "Right",
							Value: "right",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "integer_only",
					Description: "Only allow integer bids (no decimals).",
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
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "host",
					Description: "Changes the host",
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
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "winner",
					Description: "Set the current winner",
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
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "currency_side",
					Description: "Left/Right currency",
					//Autocomplete: true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Left",
							Value: "left",
						},
						{
							Name:  "Right",
							Value: "right",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "integer_only",
					Description: "Only allow integer bids (no decimals).",
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
		/*{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "help",
			Description: "auction info",
		},*/
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

	info := map[string]interface{}{
		"guild_id": i.GuildID,
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.AuctionSetup{}).Create(&info)

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

	result = database.DB.Model(database.AuctionSetup{
		GuildID: i.GuildID,
	}).Updates(options)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	//Now check what options are set
	setOptions := map[string]interface{}{}

	database.DB.Model(database.AuctionSetup{}).First(&setOptions, i.GuildID)

	antiSnipeDescription := "Anti Snipe Disabled. To enable, set both snipe_extension and snipe_range"

	responseFields := []*discordgo.MessageEmbedField{}

	for _, v := range AuctionCommand.Options[0].Options {
		if !strings.Contains(v.Name, "snipe") {
			switch {
			case setOptions[v.Name] == nil:
				setOptions[v.Name] = "Not Set"
			case strings.Contains(v.Name, "role"):
				setOptions[v.Name] = fmt.Sprintf("<@&%s>", setOptions[v.Name])
			case strings.Contains(v.Name, "channel"):
				setOptions[v.Name] = fmt.Sprintf("<#%s>", setOptions[v.Name])
			case strings.Contains(v.Name, "category"):
				category, err := s.Channel(setOptions[v.Name].(string))
				if err != nil {
					fmt.Println("Category Error:", err)
					setOptions[v.Name] = "Error Displaying Category: " + err.Error()
				} else {
					setOptions[v.Name] = category.Name
				}
			}
			responseFields = append(responseFields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("**%s**", strings.Title(strings.ReplaceAll(v.Name, "_", " "))),
				Value: fmt.Sprint(setOptions[v.Name]),
			})
		}
	}

	if setOptions["snipe_range"] != nil && setOptions["snipe_extension"] != nil {
		antiSnipeDescription = fmt.Sprintf("If a bid is placed within %s of the auction ending, it will be extended by %s.", setOptions["snipe_range"].(time.Duration).String(), setOptions["snipe_extension"].(time.Duration).String())
	}

	responseFields = append(responseFields, &discordgo.MessageEmbedField{
		Name:  "**Anti Snipe**",
		Value: antiSnipeDescription,
	})

	menuOptions := []discordgo.SelectMenuOption{}

	for _, v := range AuctionCommand.Options[0].Options {
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
						MaxValues:   len(AuctionCommand.Options[0].Options),
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

func AuctionFormat(s *discordgo.Session, auctionMap map[string]interface{}) (h.PresetResponse, error) {

	content := ""
	imageURL := ""

	if auctionMap["item"] != nil && len(auctionMap["item"].(string)) > 100 {
		return h.PresetResponse{}, fmt.Errorf("Title cannot be over 100 characters long")
	}

	if auctionMap["image_url"] != nil {
		imageURL = auctionMap["image_url"].(string)
	}

	description := fmt.Sprintf("**Host:** <@%s>.\n", auctionMap["host"])

	if auctionMap["description"] != nil {
		description += fmt.Sprintf("**Description:** %s\n", auctionMap["description"])
	}
	if auctionMap["increment_min"] != nil {
		description += fmt.Sprintf("**Minimum Bid:** + %s above previous.\n", PriceFormat(auctionMap, auctionMap["increment_min"].(float64)))
	}

	if auctionMap["increment_max"] != nil {
		description += fmt.Sprintf("**Maximum Bid:** + %s above previous.\n", PriceFormat(auctionMap, auctionMap["increment_max"].(float64)))
	}

	if auctionMap["target_price"] != nil {
		description += "**Target Price:** The host has set a hidden target price for this auction.\n"
	}

	if auctionMap["integer_only"] != nil {
		description += fmt.Sprintf("**Integer Only:** %t.\n", auctionMap["integer_only"].(bool))
	}

	if auctionMap["snipe_extension"] != nil && auctionMap["snipe_range"] != nil {
		description += fmt.Sprintf("**Anti Snipe:** If a bid is placed within the last %s, the auction will be extended by %s.\n", auctionMap["snipe_range"], auctionMap["snipe_extension"].(time.Duration).String())
	}

	if auctionMap["buyout"] != nil {
		description += fmt.Sprintf("**Buyout Price:** %s.\n", PriceFormat(auctionMap, auctionMap["buyout"].(float64)))
	}

	auctionfields := []*discordgo.MessageEmbedField{
		{
			Name:   "__**Auction Details:**__",
			Value:  description,
			Inline: true,
		},
		{
			Name:   "__**End Time**__",
			Value:  fmt.Sprintf("<t:%d:R>", auctionMap["end_time"].(time.Time).Unix()),
			Inline: true,
		},
	}

	if auctionMap["winner"] != nil {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Current Highest Bid:**__",
			Value:  PriceFormat(auctionMap, auctionMap["bid"].(float64)),
			Inline: true,
		}, &discordgo.MessageEmbedField{
			Name:   "__**Current Winner**__",
			Value:  fmt.Sprintf("<@%s>", auctionMap["winner"]),
			Inline: true,
		})
	} else {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Starting Bid:**__",
			Value:  PriceFormat(auctionMap, auctionMap["bid"].(float64)),
			Inline: true,
		})
	}

	auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
		Name:   "__**How to Bid**__",
		Value:  "Use the command `/auction bid` below.\n• Ex: `/auction bid 550`.\n**Alternate Method:** reply to this auction or @ the bot with `bid <amount>`",
		Inline: false,
	})

	guild, err := s.Guild(auctionMap["guild_id"].(string))
	if err != nil {
		fmt.Println("Error fetching guild: ", err)
		return h.PresetResponse{}, err
	}

	if auctionMap["alert_role"] != nil {
		content = fmt.Sprintf("<@&%s>", strings.Trim(auctionMap["alert_role"].(string), " "))
	}

	return h.PresetResponse{
		Content: content,
		Title:   fmt.Sprintf("Item: __**%s**__", auctionMap["item"]),
		Fields:  auctionfields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: guild.IconURL(),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
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
	}, nil
}

func AuctionPlanner(s *discordgo.Session, i *discordgo.InteractionCreate) {

	auctionMap := h.ParseSubCommand(i)
	auctionSetup := map[string]interface{}{}

	database.DB.Model(&database.AuctionSetup{}).First(&auctionSetup, i.GuildID)
	auctionMap["guild_id"] = i.GuildID
	auctionMap["host"] = i.Member.User.ID

	for _, key := range []string{"category", "currency", "snipe_extension", "snipe_range", "currency_side", "integer_only", "alert_role"} {
		if auctionMap[key] == nil {
			auctionMap[key] = auctionSetup[key]
		}
	}

	endTimeDuration, err := h.ParseTime(strings.ToLower(auctionMap["duration"].(string)))
	if err != nil {
		fmt.Println(err)
		h.ErrorResponse(s, i, err.Error())
		return
	}

	delete(auctionMap, "duration")

	canHost := false

	if auctionSetup["host_role"] != nil {
		for _, v := range i.Member.Roles {
			if v == auctionSetup["host_role"].(string) {
				canHost = true
			}
		}
		if i.Member.Permissions&(1<<3) == 8 {
			canHost = true
		}
		if !canHost {
			h.ErrorResponse(s, i, "User must be administrator or have the role <@&"+auctionSetup["host_role"].(string)+"> to host auctions.")
			return
		}
	}

	if auctionMap["schedule"] != nil {

		var AuctionQueue []database.AuctionQueue

		database.DB.Where(map[string]interface{}{"guild_id": i.GuildID}).Find(&AuctionQueue)

		if len(AuctionQueue) >= 25 {
			h.ErrorResponse(s, i, "You can only schedule 25 auctions in advance.")
			return
		}

		startTimeDuration, err := h.ParseTime(strings.ToLower(auctionMap["schedule"].(string)))
		if err != nil {
			fmt.Println(err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

		auctionMap["end_time"] = time.Now().Add(endTimeDuration).Add(startTimeDuration)
		auctionMap["start_time"] = time.Now().Add(startTimeDuration)
		delete(auctionMap, "schedule")

		result := database.DB.Model(database.AuctionQueue{}).Create(&auctionMap)
		if result.Error != nil {
			fmt.Println(result.Error)
			h.ErrorResponse(s, i, result.Error.Error())
			return
		}
		exampleMessage, err := AuctionFormat(s, auctionMap)
		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(file, line, err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

		err = h.PremiumResponse(s, i, h.PresetResponse{
			Title: "Auction has been Scheduled!",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "**Auction Start Time:**",
					Value:  fmt.Sprintf("<t:%d:R>", auctionMap["start_time"].(time.Time).Unix()),
					Inline: false,
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "[__**PREVIEW:**__] " + exampleMessage.Title,
					Description: exampleMessage.Description,
					Color:       0x8073ff,
					Image:       exampleMessage.Image,
					Thumbnail:   exampleMessage.Thumbnail,
					Fields:      exampleMessage.Fields,
				},
			},
		})
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Until(auctionMap["start_time"].(time.Time)))

		err = AuctionCreate(s, auctionMap)
		if err != nil {
			fmt.Println(err)
			h.ErrorMessage(s, i.ChannelID, err.Error())
			return
		}

	} else {

		auctionMap["end_time"] = time.Now().Add(endTimeDuration)

		err = AuctionCreate(s, auctionMap)
		if err != nil {
			fmt.Println(err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

		err = h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "**Auction Starting**",
			Description: "Auction has successfully been started!",
		})

		if err != nil {
			fmt.Println(err)
		}

	}

	time.Sleep(time.Until(auctionMap["end_time"].(time.Time)))
	AuctionEnd(auctionMap)
}

func AuctionCreate(s *discordgo.Session, auctionMap map[string]interface{}) error {

	result := database.DB.Delete(database.AuctionQueue{}, auctionMap["id"])
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	delete(auctionMap, "start_time")
	delete(auctionMap, "id")

	auctionMessage, err := AuctionFormat(s, auctionMap)
	if err != nil {
		return err
	}

	delete(auctionMap, "alert_role")
	delete(auctionMap, "snipe_extension")
	delete(auctionMap, "snipe_range")

	if auctionMap["category"] == nil {
		auctionMap["category"] = ""
	}

	channel, err := s.GuildChannelCreateComplex(auctionMap["guild_id"].(string), discordgo.GuildChannelCreateData{
		Name:     "💸│" + auctionMap["item"].(string),
		Type:     0,
		ParentID: auctionMap["category"].(string),
	})
	if err != nil {
		return err
	}

	message, err := h.SuccessMessage(s, channel.ID, auctionMessage)
	if err != nil {
		return err
	}

	auctionMap["channel_id"] = message.ChannelID
	auctionMap["message_id"] = message.ID
	delete(auctionMap, "category")

	result = database.DB.Model(database.Auction{}).Create(auctionMap)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func AuctionEdit(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := h.ParseSubCommand(i)

	auctionMap := map[string]interface{}{}
	auctionSetup := map[string]interface{}{}
	result := database.DB.Model(database.Auction{}).First(&auctionMap, i.ChannelID)

	if result.Error != nil {
		fmt.Println(result.Error)
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != auctionMap["host"] {
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
			options["end_time"] = auctionMap["end_time"].(time.Time).Add(extraDuration)
			delete(options, "extend")
		case "queue_number":
		default:
			auctionMap[key] = value
		}

	}

	if options["queue_number"] != nil {

		//Need to fix this since the databaser covers all guilds.
		guildQueue := []database.AuctionQueue{}

		result := database.DB.Where(map[string]interface{}{"guild_id": i.GuildID}).Find(&guildQueue)
		if result.Error != nil {
			fmt.Println(result.Error)
			h.ErrorResponse(s, i, result.Error.Error())
			return
		}

		queueNumber := guildQueue[int(options["queue_number"].(float64))-1].ID

		delete(options, "queue_number")

		result = database.DB.Model(database.AuctionQueue{
			ID: queueNumber,
		}).Updates(options)

		if result.Error != nil {
			fmt.Println(result.Error)
			h.ErrorResponse(s, i, result.Error.Error())
			return
		}

	} else {

		result = database.DB.Model(database.AuctionSetup{}).First(auctionSetup, i.GuildID)

		if result.Error != nil {
			fmt.Println(result.Error)
		}

		result := database.DB.Model(database.Auction{
			ChannelID: i.ChannelID,
		}).Updates(options)

		if result.Error != nil {
			fmt.Println(result.Error)
			h.ErrorResponse(s, i, result.Error.Error())
			return
		}

		if auctionMap["message_id"] == nil {
			h.ErrorResponse(s, i, "No auction found in this channel")
			return
		}
		if auctionMap["bid_history"] == nil {
			auctionMap["bid_history"] = ""
		}
		if auctionMap["bid"] != nil && auctionMap["winner"] != nil {
			member, err := s.GuildMember(i.GuildID, auctionMap["winner"].(string))
			if err != nil {
				fmt.Println(result.Error)
			}
			username := member.Nick
			if username == "" {
				username = member.User.Username
			}
			auctionMap["bid_history"] = auctionMap["bid_history"].(string) + "\n-> " + username + ": " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", auctionMap["bid"]), "0"), ".")
		}
		auctionMap["snipe_extension"] = auctionSetup["snipe_extension"]
		auctionMap["snipe_range"] = auctionSetup["snipe_range"]

		message, err := s.ChannelMessage(i.ChannelID, auctionMap["message_id"].(string))
		if err != nil {
			fmt.Println("Error fetching message", err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

		formattedMessage, err := AuctionFormat(s, auctionMap)
		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(file, line, err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

		message.Embeds[0] = &discordgo.MessageEmbed{
			Title:       formattedMessage.Title,
			Description: formattedMessage.Description,
			Color:       message.Embeds[0].Color,
			Image:       formattedMessage.Image,
			Thumbnail:   formattedMessage.Thumbnail,
			Fields:      formattedMessage.Fields,
		}
		message.Content = formattedMessage.Content

		_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content:    &message.Content,
			Components: message.Components,
			Embeds:     message.Embeds,
			ID:         message.ID,
			Channel:    i.ChannelID,
		})

		if err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(file, line, err)
			h.ErrorResponse(s, i, err.Error())
			return
		}

		if options["item"] != nil {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				fmt.Println(err)
			}
			_, err = s.ChannelEditComplex(i.ChannelID, &discordgo.ChannelEdit{
				Name:     "💸│" + options["item"].(string),
				Position: channel.Position,
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Content:     "",
		Title:       "Success",
		Description: "Auction has successfully been edited",
	})
}

func AuctionBid(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := map[string]interface{}{}

	if i.ApplicationCommandData().Name == "bid" {
		options = h.ParseSlashCommand(i)
	} else {
		options = h.ParseSubCommand(i)
	}

	message, err := AuctionBidFormat(s, database.Auction{
		ChannelID: i.ChannelID,
		Bid:       options["amount"].(float64),
		Winner:    i.Member.User.ID,
		GuildID:   i.GuildID,
	})

	if err != nil {
		fmt.Println(err)
		h.ErrorResponse(s, i, err.Error())
		return
	}

	err = h.SuccessResponse(s, i, message)
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionBidFormat(s *discordgo.Session, bidData database.Auction) (h.PresetResponse, error) {

	response := h.PresetResponse{}
	var Content string
	var responseFields []*discordgo.MessageEmbedField
	auctionSetup := map[string]interface{}{}
	auctionMap := map[string]interface{}{}

	result := database.DB.Model(database.Auction{}).First(&auctionMap, bidData.ChannelID)
	if result.Error != nil {
		return response, result.Error
	}
	result = database.DB.Model(database.AuctionSetup{}).First(&auctionSetup, bidData.GuildID)
	if result.Error != nil {
		return response, result.Error
	}

	auctionMap["snipe_extension"] = auctionSetup["snipe_extension"]
	auctionMap["snipe_range"] = auctionSetup["snipe_range"]

	if auctionMap["snipe_range"] != nil && auctionMap["snipe_extension"] != nil {
		if time.Until(auctionMap["end_time"].(time.Time)) < auctionMap["snipe_range"].(time.Duration) && auctionMap["snipe_extension"] != 0 {
			auctionMap["end_time"] = auctionMap["end_time"].(time.Time).Add(auctionMap["snipe_extension"].(time.Duration))
			responseFields = []*discordgo.MessageEmbedField{
				{
					Name:   "**Anti-Snipe Activated!**",
					Value:  fmt.Sprintf("New End Time: <t:%d>", auctionMap["end_time"].(time.Time).Unix()),
					Inline: false,
				},
			}
		}
	}

	delete(auctionMap, "snipe_range")
	delete(auctionMap, "snipe_extension")

	switch {
	case auctionMap["end_time"].(time.Time).Before(time.Now()):
		return response, fmt.Errorf("Cannot Bid, Auction has ended")
	case bidData.Winner == auctionMap["winner"] && auctionMap["increment_max"] != nil:
		return response, fmt.Errorf("Cannot out bid yourself on a capped bid auction!")
	case auctionMap["integer_only"] != nil && auctionMap["integer_only"].(bool) && bidData.Bid != math.Floor(bidData.Bid):
		return response, fmt.Errorf("Your bid must be an integer for this auction! For example: " + fmt.Sprint(math.Floor(bidData.Bid)) + " instead of " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidData.Bid), "0"), "."))
	case auctionMap["buyout"] != nil && bidData.Bid >= auctionMap["buyout"].(float64) && auctionMap["buyout"] != 0:
		auctionMap["bid"] = bidData.Bid
		auctionMap["winner"] = bidData.Winner

		user, err := s.GuildMember(bidData.GuildID, bidData.Winner)
		if err != nil {
			return response, err
		}
		username := user.Nick
		if username == "" {
			username = user.User.Username
		}

		bidAmount := bidData.Bid

		if auctionMap["bid_history"] == nil {
			auctionMap["bid_history"] = ""
		}

		auctionMap["bid_history"] = auctionMap["bid_history"].(string) + "\n-> " + username + ": " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), ".")

		database.DB.Model(database.Auction{
			ChannelID: bidData.ChannelID,
		}).Updates(auctionMap)

		response = h.PresetResponse{
			Title:       "Success!",
			Description: "Auction has successfully been bought out!",
		}

		AuctionEnd(auctionMap)

		return response, nil
	case bidData.Bid >= auctionMap["bid"].(float64):
		switch {
		case auctionMap["increment_min"] != nil && bidData.Bid-auctionMap["bid"].(float64) < auctionMap["increment_min"].(float64):
			return response, fmt.Errorf("Bid must be higher than the previous bid by: " + fmt.Sprintf("%s\n\u200b", PriceFormat(auctionMap, auctionMap["increment_min"].(float64))))
		case auctionMap["increment_max"] != nil && bidData.Bid-auctionMap["bid"].(float64) > auctionMap["increment_max"].(float64):
			return response, fmt.Errorf("Bid must be no more than " + fmt.Sprintf("%s\n\u200b", PriceFormat(auctionMap, auctionMap["increment_max"].(float64))) + " Higher than the previous bid.")
		case auctionMap["winner"] != nil && bidData.Bid <= auctionMap["bid"].(float64):
			return response, fmt.Errorf("You must bid higher than: " + PriceFormat(auctionMap, auctionMap["bid"].(float64)))
		}

		auctionMap["bid"] = bidData.Bid
		auctionMap["winner"] = bidData.Winner

		user, err := s.GuildMember(bidData.GuildID, bidData.Winner)
		if err != nil {
			return response, err
		}
		username := user.Nick
		if username == "" {
			username = user.User.Username
		}

		bidAmount := bidData.Bid

		if auctionMap["bid_history"] == nil {
			auctionMap["bid_history"] = ""
		}

		auctionMap["bid_history"] = auctionMap["bid_history"].(string) + "\n-> " + username + ": " + strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", bidAmount), "0"), ".")

		database.DB.Model(database.Auction{
			ChannelID: bidData.ChannelID,
		}).Updates(auctionMap)

		updateAuction, err := s.ChannelMessage(auctionMap["channel_id"].(string), auctionMap["message_id"].(string))
		if err != nil {
			return response, err
		}

		bidHistory := auctionMap["bid_history"].(string)
		if len(bidHistory) > 4095 {
			bidHistory = bidHistory[len(bidHistory)-4095:]
		}

		m, err := AuctionFormat(s, auctionMap)
		if err != nil {
			return response, err
		}
		updateAuction.Embeds[0].Fields = m.Fields
		updateAuction.Embeds[0].Description = m.Description

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
			ID:         auctionMap["message_id"].(string),
			Channel:    auctionMap["channel_id"].(string),
		})
		if err != nil {
			return response, err
		}
		Content = "Bid has successfully been placed"
	default:
		return response, fmt.Errorf("You must bid higher than: " + PriceFormat(auctionMap, auctionMap["bid"].(float64)))
	}

	response = h.PresetResponse{
		Title:  Content,
		Fields: responseFields,
	}

	return response, nil
}

func AuctionBidHistory(s *discordgo.Session, i *discordgo.InteractionCreate) {

	claimMap := map[string]interface{}{}

	result := database.DB.Model(database.Claim{}).First(claimMap, i.Message.ID)

	if result.Error != nil {
		fmt.Println(result.Error)
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	if claimMap["bid_history"] == nil {
		h.ErrorResponse(s, i, "No bid history found for this auction.")
		return
	}

	bidHistory := claimMap["bid_history"].(string)

	if len(bidHistory) > 4095 {
		bidHistory = bidHistory[len(bidHistory)-4095:]
	}

	err := h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Bid History**",
		Description: bidHistory,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/9wo7diC.png",
		},
	})
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(file, line, err)
	}
}

func AuctionEnd(auctionMap map[string]interface{}) error {

	AuctionSetup := database.AuctionSetup{
		GuildID: auctionMap["guild_id"].(string),
	}

	result := database.DB.First(&AuctionSetup)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	result = database.DB.Model(database.Auction{}).First(&auctionMap, auctionMap["channel_id"])
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	if auctionMap["end_time"].(time.Time).After(time.Now()) {
		time.Sleep(time.Until(auctionMap["end_time"].(time.Time)))
		AuctionEnd(auctionMap)
		return nil
	}

	if auctionMap["buyout"] != nil && auctionMap["buyout"].(float64) != 0 {
		if auctionMap["bid"].(float64) < auctionMap["buyout"].(float64) && auctionMap["end_time"].(time.Time).After(time.Now()) {
			time.Sleep(time.Until(auctionMap["end_time"].(time.Time)))
			AuctionEnd(auctionMap)
			return nil
		}
	}

	message := discordgo.NewMessageEdit(auctionMap["channel_id"].(string), auctionMap["message_id"].(string))
	messageEmbeds, err := Session.ChannelMessage(auctionMap["channel_id"].(string), auctionMap["message_id"].(string))
	if err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(file, line, err)
		result = database.DB.Delete(database.Auction{
			ChannelID: auctionMap["channel_id"].(string),
		})
		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		return nil
	}

	if AuctionSetup.LogChannel == "" {
		fmt.Println("Log channel has not been set for guild: " + auctionMap["guild_id"].(string))
		_, err := h.ErrorMessage(Session, auctionMap["channel_id"].(string), "Auction cannot end because log channel has not been set. Please setup an auction log using `/auction setup`. You might need to end the auction manually after setting the channel.")
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return nil
	}

	if auctionMap["target_price"] != nil && auctionMap["target_price"].(float64) > auctionMap["bid"].(float64) {
		auctionMap["target_message"] = fmt.Sprintf("The host had set a target price of %s that has not been reached.", PriceFormat(auctionMap, auctionMap["target_price"].(float64)))
		auctionMap["winner"] = nil
	}

	auctionMap["formatted_price"] = fmt.Sprintf("%s\n\u200b", PriceFormat(auctionMap, auctionMap["bid"].(float64)))

	if auctionMap["buyout"] != nil && auctionMap["buyout"].(float64) != 0 {
		if auctionMap["bid"].(float64) >= auctionMap["buyout"].(float64) {
			auctionMap["buyout_message"] = fmt.Sprintf("%s\n\u200b", PriceFormat(auctionMap, auctionMap["buyout"].(float64))) + " BUYOUT!"
		}
	}

	auctionMap["log_channel"] = AuctionSetup.LogChannel
	deleteChannel := auctionMap["channel_id"].(string)

	err = ClaimOutput(Session, auctionMap, "Auction")
	if err != nil {
		return err
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

	time.Sleep(30 * time.Second)

	_, err = Session.ChannelDelete(deleteChannel)
	if err != nil {
		fmt.Println(err)
	}

	result = database.DB.Delete(database.Auction{}, deleteChannel)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	return nil
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
			Name: fmt.Sprintf("**%s. %s**", fmt.Sprint(len(fields)+1), v.Item),
			Value: fmt.Sprintf("**Start time:** <t:%d:R>\n**End Time:** <t:%d>\n**Starting Price:** %s\n\u200b", v.StartTime.Unix(), v.EndTime.Unix(), PriceFormat(map[string]interface{}{
				"currency":      v.Currency,
				"currency_side": v.CurrencySide,
			}, v.Bid)),
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

	auctionSetup := database.AuctionSetup{}
	auctionMap := map[string]interface{}{}

	result := database.DB.Model(database.Auction{}).First(&auctionMap, i.ChannelID)
	if result.Error != nil {
		fmt.Println(result.Error)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label: "Delete Auction Channel",
								Style: 4,
								Emoji: discordgo.ComponentEmoji{
									Name: "🛑",
								},
								CustomID: "delete_auction_channel",
							},
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
				},
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Error",
						Description: "There was an error. Please contact support if you need help.\n**If you would like to close this auction channel, press the button below.**",
						Color:       0xff0000,
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:  "**Error Message:**",
								Value: result.Error.Error(),
							},
						},
					},
				},
				Flags: 64,
			},
		})
		return
	}
	database.DB.First(&auctionSetup, i.GuildID)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	if auctionMap["log_channel"] == "" {
		fmt.Println("Log channel has not been set for guild: " + i.GuildID)
		h.ErrorResponse(s, i, "Auction cannot end because log channel has not been set. Please setup an auction log using `/auction setup`")
		return
	}

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != auctionMap["host"] {
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

	auctionMap["end_time"] = time.Now()

	result = database.DB.Model(database.Auction{
		ChannelID: i.ChannelID,
	}).Updates(auctionMap)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	err = AuctionEnd(auctionMap)
	if err != nil {
		h.ErrorMessage(s, i.ChannelID, err.Error())
		return
	}
}

func DeleteAuctionChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := s.ChannelDelete(i.ChannelID)
	if err != nil {
		fmt.Println(err)
		h.ErrorResponse(s, i, err.Error())
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
		h.ErrorResponse(s, i, "User must have administrator permissions to run this command")
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
	clearedMap := map[string]interface{}{}

	info := database.AuctionSetup{
		GuildID: i.GuildID,
	}

	clearedSettings := "No Settings Cleared!"
	if len(options) > 0 {
		clearedSettings = ""
	}

	for _, v := range options {
		clearedSettings += fmt.Sprintf("• %s\n", strings.Title(strings.ReplaceAll(v, "_", " ")))
		clearedMap[v] = gorm.Expr("NULL")
	}

	database.DB.Model(&info).Updates(clearedMap)

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

func PriceFormat(auctionMap map[string]interface{}, price float64) string {

	if auctionMap["currency"] == nil {
		if auctionMap["currency_side"] == "right" {
			return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", price), "0"), ".")
		} else {
			return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", price), "0"), ".")
		}
	} else {
		if auctionMap["currency_side"] == "right" {
			return fmt.Sprintf("%s %s", strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", price), "0"), "."), auctionMap["currency"])
		} else {
			return fmt.Sprintf("%s %s", auctionMap["currency"], strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", price), "0"), "."))
		}
	}
}
