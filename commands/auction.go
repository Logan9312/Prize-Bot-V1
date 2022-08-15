package commands

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
)

type EventType string

const (
	EventTypeAuction  EventType = "Auction"
	EventTypeShop     EventType = "Shop"
	EventTypeGiveaway EventType = "Giveaway"
)

var AuctionCommand = discordgo.ApplicationCommand{
	Name:        "auction",
	Description: "Put an item up for auction!",
	Options: []*discordgo.ApplicationCommandOption{
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
				/*{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "use_currency",
					Description: "The winner will pay with their currency balance.",
				},*/
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
					Type:        11,
					Name:        "image",
					Description: "Attach an image to your auction",
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
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "image",
					Description: "Edit the image",
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

func Auction(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Options[0].Name {
	case "create":
		return AuctionCreate(s, i)
	case "queue":
		return AuctionQueue(s, i)
	case "edit":
		return AuctionEdit(s, i)
	}
	return fmt.Errorf("Unknown Auction command, please contact support")
}

func AuctionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := h.ParseSubCommand(i)

	if options["schedule"] != nil {
		if !CheckPremiumGuild(i.GuildID) {
			return h.PremiumError(s, i)
		}
	}

	if options["image"] != nil {
		options["image_url"] = i.ApplicationCommandData().Resolved.Attachments[options["image"].(string)].URL
		delete(options, "image")
	}

	channelID, err := AuctionHandler(s, options, i.Member, i.GuildID)
	if err != nil {
		return err
	}

	if channelID != "" {
		return h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "**Auction Starting**",
			Description: fmt.Sprintf("Auction has successfully been started in <#%s>!", channelID),
		})
	} else {
		exampleMessage, err := AuctionFormat(s, options, EventTypeAuction)
		if err != nil {
			return err
		}

		return h.SuccessResponse(s, i, h.PresetResponse{
			Title: "Auction has been Scheduled!",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "**Auction Start Time:**",
					Value:  fmt.Sprintf("<t:%d:R>", options["start_time"].(time.Time).Unix()),
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
	}
}

func AuctionHandler(s *discordgo.Session, auctionMap map[string]any, member *discordgo.Member, guildID string) (channelID string, err error) {
	auctionSetup := map[string]interface{}{}
	currencyMap := map[string]interface{}{}

	if !AuctionHostCheck(auctionMap, member) {
		return "", fmt.Errorf("User must be administrator or have the role <@&" + auctionSetup["host_role"].(string) + "> to host auctions.")
	}

	result := database.DB.Model(&database.AuctionSetup{}).First(&auctionSetup, guildID)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	result = database.DB.Model(database.CurrencySetup{}).First(&currencyMap, guildID)
	if result.Error != nil {
		fmt.Println("Error getting currency setup: " + result.Error.Error())
	}

	auctionMap["guild_id"] = guildID
	auctionMap["host"] = member.User.ID

	if auctionMap["currency"] == nil {
		auctionMap["currency"] = currencyMap["currency"]
	}

	for _, key := range []string{"category", "snipe_extension", "snipe_range", "currency_side", "integer_only", "alert_role", "channel_override"} {
		if auctionMap[key] == nil {
			auctionMap[key] = auctionSetup[key]
		}
	}

	auctionMap["duration"], err = h.ParseTime(strings.ToLower(auctionMap["duration"].(string)))
	if err != nil {
		return "", fmt.Errorf("Error parsing duration input: %w", err)
	}

	if auctionMap["schedule"] != nil {
		err = AuctionSchedule(s, auctionMap)
		if err != nil {
			return "", fmt.Errorf("Error scheduling auction: %w", err)
		}
	} else {
		auctionMap["end_time"] = time.Now().Add(auctionMap["duration"].(time.Duration))
		delete(auctionMap, "duration")

		channelID, err := AuctionStart(s, auctionMap)
		if err != nil {
			return channelID, fmt.Errorf("Error starting auction: %w", err)
		}

		return channelID, nil

	}

	return "", nil
}

func AuctionSchedule(s *discordgo.Session, auctionMap map[string]any) error {

	var AuctionQueue []database.AuctionQueue

	database.DB.Where(map[string]interface{}{"guild_id": auctionMap["guild_id"].(string)}).Find(&AuctionQueue)
	if len(AuctionQueue) >= 25 {
		return fmt.Errorf("You can only schedule 25 auctions in advance.")
	}

	startTimeDuration, err := h.ParseTime(strings.ToLower(auctionMap["schedule"].(string)))
	if err != nil {
		return err
	}

	auctionMap["end_time"] = time.Now().Add(auctionMap["duration"].(time.Duration)).Add(startTimeDuration)
	auctionMap["start_time"] = time.Now().Add(startTimeDuration)
	delete(auctionMap, "schedule")
	delete(auctionMap, "duration")

	result := database.DB.Model(database.AuctionQueue{}).Create(&auctionMap)
	if result.Error != nil {
		return fmt.Errorf("Error saving queued auction to database: %w", result.Error)
	}

	go AuctionStartTimer(s, auctionMap)

	return nil
}

func AuctionStartTimer(s *discordgo.Session, auctionMap map[string]any) {
	time.Sleep(time.Until(auctionMap["start_time"].(time.Time)))

	_, err := AuctionStart(s, auctionMap)
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionEndTimer(s *discordgo.Session, auctionMap map[string]any) {
	time.Sleep(time.Until(auctionMap["end_time"].(time.Time)))
	err := AuctionEnd(s, auctionMap["channel_id"].(string), auctionMap["guild_id"].(string))
	if err != nil {
		fmt.Println(err)
	}
}

func AuctionHostCheck(auctionSetup map[string]any, member *discordgo.Member) bool {
	if auctionSetup["host_role"] == nil {
		return true
	}
	for _, v := range member.Roles {
		if v == auctionSetup["host_role"].(string) {
			return true
		}
	}

	return member.Permissions&discordgo.PermissionAdministrator == 8
}

func AuctionStart(s *discordgo.Session, auctionMap map[string]interface{}) (string, error) {
	if auctionMap["id"] != nil {
		result := database.DB.Delete(database.AuctionQueue{}, auctionMap["id"])
		if result.Error != nil {
			fmt.Println(result.Error)
		}
	}

	auctionMessage, err := AuctionFormat(s, auctionMap, EventTypeAuction)
	if err != nil {
		return "", err
	}

	delete(auctionMap, "start_time")
	delete(auctionMap, "id")
	delete(auctionMap, "alert_role")
	delete(auctionMap, "snipe_extension")
	delete(auctionMap, "snipe_range")

	channel, err := MakeAuctionChannel(s, auctionMap)
	if err != nil {
		return channel.ID, err
	}
	auctionMap["channel_id"] = channel.ID

	message, err := h.SuccessMessage(s, auctionMap["channel_id"].(string), auctionMessage)
	if err != nil {
		return channel.ID, err
	}
	auctionMap["message_id"] = message.ID

	delete(auctionMap, "category")
	result := database.DB.Model(database.Auction{}).Create(auctionMap)
	if result.Error != nil {
		return channel.ID, result.Error
	}

	go AuctionEndTimer(s, auctionMap)

	return channel.ID, nil
}

func MakeAuctionChannel(s *discordgo.Session, auctionMap map[string]any) (channel *discordgo.Channel, err error) {

	var category string

	if auctionMap["category"] != nil {
		category = auctionMap["category"].(string)
	}

	if auctionMap["channel_override"] == nil {
		channel, err = s.GuildChannelCreateComplex(auctionMap["guild_id"].(string), discordgo.GuildChannelCreateData{
			Name:     "💸│" + auctionMap["item"].(string),
			Type:     discordgo.ChannelTypeGuildText,
			ParentID: category,
		})
		if err != nil {
			err = fmt.Errorf("Error creating auction channel: %w", err)
		}

	} else {
		channel, err = s.Channel(auctionMap["channel_override"].(string))
		if err != nil {
			err = fmt.Errorf("Error finding auction channel: %w", err)
		}
	}

	return
}

func AuctionBid(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := h.ParseSlashCommand(i)
	auctionMap := map[string]any{}

	result := database.DB.Model(database.Auction{}).First(&auctionMap, i.ChannelID)
	if result.Error != nil {
		return fmt.Errorf("Error fetching auction data from the database. Error Message: " + result.Error.Error())
	}

	err := AuctionBidPlace(s, options["amount"].(float64), i.Member, i.ChannelID, i.GuildID)
	if err != nil {
		return err
	}

	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title: "Bid has been successfully placed!",
	})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func AuctionBidPlace(s *discordgo.Session, amount float64, member *discordgo.Member, channelID, guildID string) error {

	auctionSetup := map[string]interface{}{}
	auctionMap := map[string]interface{}{}
	p := message.NewPrinter(language.English)

	result := database.DB.Model(database.Auction{}).First(&auctionMap, channelID)
	if result.Error != nil {
		return fmt.Errorf("Error fetching auction data from the database. Error Message: " + result.Error.Error())
	}
	result = database.DB.Model(database.AuctionSetup{}).First(&auctionSetup, guildID)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	if auctionSetup["snipe_range"] != nil && auctionSetup["snipe_extension"] != nil {
		if time.Until(auctionMap["end_time"].(time.Time)) < auctionSetup["snipe_range"].(time.Duration) {
			auctionMap["end_time"] = auctionMap["end_time"].(time.Time).Add(auctionSetup["snipe_extension"].(time.Duration))
			h.SuccessMessage(s, channelID, h.PresetResponse{
				Title:       "**Anti-Snipe Activated!**",
				Description: fmt.Sprintf("New End Time: <t:%d>", auctionMap["end_time"].(time.Time).Unix()),
			})
		}
	}

	//Checking if the auction has ended.
	if auctionMap["end_time"].(time.Time).Before(time.Now()) {
		return fmt.Errorf("cannot Bid, Auction has ended")
	}

	//Checking if the auction is capped and the current winner is bidding.
	if member.User.ID == auctionMap["winner"] && auctionMap["increment_max"] != nil {
		return fmt.Errorf("cannot out bid yourself on a capped bid auction")
	}

	//Checking if integer only bidding is enabled.
	if auctionMap["integer_only"] != nil && auctionMap["integer_only"].(bool) && amount != math.Floor(amount) {
		return fmt.Errorf("Your bid must be an integer for this auction! For example: " + fmt.Sprint(math.Floor(amount)) + " instead of " + PriceFormat(amount, guildID, auctionMap["currency"]))
	}

	//Checking if bid is higher than minimum increment.
	if auctionMap["increment_min"] != nil && amount-auctionMap["bid"].(float64) < auctionMap["increment_min"].(float64) {
		return fmt.Errorf("Bid must be higher than the previous bid by: %s\n\u200b", PriceFormat(auctionMap["increment_min"].(float64), guildID, auctionMap["currency"]))
	}

	//Checking if bid is lower than maximum increment.
	if auctionMap["increment_max"] != nil && amount-auctionMap["bid"].(float64) > auctionMap["increment_max"].(float64) {
		return fmt.Errorf("Bid must be no more than %s higher than the previous bid. \n\u200b", PriceFormat(auctionMap["increment_max"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]))
	}

	if amount < auctionMap["bid"].(float64) {
		return fmt.Errorf("You must bid higher than: " + PriceFormat(auctionMap["bid"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]))
	}

	if amount == auctionMap["bid"].(float64) && auctionMap["winner"] != nil {
		return fmt.Errorf("You must bid higher than: " + PriceFormat(auctionMap["bid"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]))
	}

	if auctionMap["bid_history"] == nil {
		auctionMap["bid_history"] = ""
	}

	auctionMap["bid"] = amount
	auctionMap["winner"] = member.User.ID
	auctionMap["bid_history"] = auctionMap["bid_history"].(string) + "\n-> " + member.User.Username + ": " + strings.TrimRight(strings.TrimRight(p.Sprintf("%f", amount), "0"), ".")

	result = database.DB.Model(database.Auction{
		ChannelID: channelID,
	}).Updates(auctionMap)
	if result.Error != nil {
		return result.Error
	}

	auctionMap["snipe_extension"] = auctionSetup["snipe_extension"]
	auctionMap["snipe_range"] = auctionSetup["snipe_range"]

	m, err := AuctionFormat(s, auctionMap, EventTypeAuction)
	if err != nil {
		return err
	}

	_, err = h.SuccessMessageEdit(s, channelID, auctionMap["message_id"].(string), m)
	if err != nil {
		return err
	}

	if auctionMap["buyout"] != nil && amount >= auctionMap["buyout"].(float64) {

		auctionMap["end_time"] = time.Now()

		err = AuctionEnd(s, channelID, guildID)
		if err != nil {
			return err
		}
	}

	return nil
}

// AuctionEdit handles the AuctionEdit event
func AuctionEdit(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := h.ParseSubCommand(i)

	if options["image"] != nil {
		options["image_url"] = i.ApplicationCommandData().Resolved.Attachments[options["image"].(string)].URL
		delete(options, "image")
	}

	if options["queue_number"] != nil {
		err := AuctionQueueUpdate(options, i.GuildID)
		if err != nil {
			return err
		}
	} else {
		err := AuctionUpdate(s, options, i.Member, i.ChannelID, i.GuildID)
		if err != nil {
			return err
		}
	}

	return h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Success",
		Description: "Auction has successfully been edited",
	})
}

// AuctionUpdate updates a currently running auction
func AuctionUpdate(s *discordgo.Session, options map[string]any, member *discordgo.Member, channelID, guildID string) error {

	auctionMap := map[string]any{}
	p := message.NewPrinter(language.English)

	if member.Permissions&(1<<3) != 8 && member.User.ID != auctionMap["host"] {
		return fmt.Errorf("User must have be host or have administrator permissions to run this command")
	}

	result := database.DB.Model(database.Auction{}).First(&auctionMap, channelID)
	if result.Error != nil {
		return result.Error
	}

	editedOptions := ""

	for key, value := range options {
		if key == "bid" {
			editedOptions += fmt.Sprintf("\n\u3000- %s set to: %s", key, strings.TrimRight(strings.TrimRight(p.Sprintf("%f", options["bid"].(float64)), "0"), "."))
		} else {
			editedOptions += fmt.Sprintf("\n\u3000- %s set to: %s", key, p.Sprint(value))
		}
	}

	if options["extend"] != nil {
		extraDuration, err := h.ParseTime(strings.ToLower(options["extend"].(string)))
		if err != nil {
			return err
		}
		options["end_time"] = auctionMap["end_time"].(time.Time).Add(extraDuration)
		delete(options, "extend")
	}

	if auctionMap["bid_history"] == nil {
		auctionMap["bid_history"] = ""
	}

	options["bid_history"] = auctionMap["bid_history"].(string) + "\n-> Auction edited by " + member.User.Username + ":" + editedOptions

	fmt.Println(options["bid_history"])

	result = database.DB.Model(database.Auction{
		ChannelID: channelID,
	}).Updates(options)
	if result.Error != nil {
		return result.Error
	}

	if options["bid"] != nil && options["winner"] != nil {
		member, err := s.GuildMember(guildID, options["winner"].(string))
		if err != nil {
			return err
		}
		username := member.User.Username
		options["bid_history"] = options["bid_history"].(string) + "\n-> " + username + ": " + strings.TrimRight(strings.TrimRight(p.Sprintf("%f", options["bid"].(float64)), "0"), ".")
	}

	if options["item"] != nil {
		channel, err := s.Channel(channelID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = s.ChannelEditComplex(channelID, &discordgo.ChannelEdit{
			Name:     "💸│" + options["item"].(string),
			Position: channel.Position,
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	result = database.DB.Model(database.Auction{}).First(&auctionMap, channelID)
	if result.Error != nil {
		return result.Error
	}

	message, err := AuctionFormat(s, auctionMap, EventTypeAuction)
	if err != nil {
		return err
	}

	_, err = h.SuccessMessageEdit(s, channelID, auctionMap["message_id"].(string), message)
	if err != nil {
		return err
	}

	return nil
}

// Fix this
func AuctionQueueUpdate(options map[string]any, guildID string) error {

	//Need to fix this since the database covers all guilds.
	guildQueue := []database.AuctionQueue{}

	result := database.DB.Where(map[string]interface{}{"guild_id": guildID}).Find(&guildQueue)
	if result.Error != nil {
		return result.Error
	}

	queueNumber := guildQueue[int(options["queue_number"].(float64))-1].ID

	delete(options, "queue_number")

	result = database.DB.Model(database.AuctionQueue{
		ID: queueNumber,
	}).Updates(options)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func AuctionFormat(s *discordgo.Session, auctionMap map[string]interface{}, eventType EventType) (h.PresetResponse, error) {

	content := ""
	imageURL := ""
	embeds := []*discordgo.MessageEmbed{}

	if auctionMap["item"] != nil && len(auctionMap["item"].(string)) > 100 {
		return h.PresetResponse{}, fmt.Errorf("title cannot be over 100 characters long")
	}

	if auctionMap["image_url"] != nil {
		imageURL = auctionMap["image_url"].(string)
	}

	description := fmt.Sprintf("**Host:** <@%s>.\n", auctionMap["host"])

	if auctionMap["description"] != nil {
		description += fmt.Sprintf("**Description:** %s\n", auctionMap["description"])
	}

	if auctionMap["winners"] != nil {
		description += fmt.Sprintf("**Winners:** %d\n", int(auctionMap["winners"].(float64)))
	}

	if auctionMap["increment_min"] != nil {
		description += fmt.Sprintf("**Minimum Bid:** + %s above previous.\n", PriceFormat(auctionMap["increment_min"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]))
	}

	if auctionMap["increment_max"] != nil {
		description += fmt.Sprintf("**Maximum Bid:** + %s above previous.\n", PriceFormat(auctionMap["increment_max"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]))
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
		description += fmt.Sprintf("**Buyout Price:** %s.\n", PriceFormat(auctionMap["buyout"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]))
	}

	auctionfields := []*discordgo.MessageEmbedField{
		{
			Name:  "__**Details:**__",
			Value: description,
		},
	}

	if auctionMap["end_time"] != nil {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**End Time**__",
			Value:  fmt.Sprintf("<t:%d:R>", auctionMap["end_time"].(time.Time).Unix()),
			Inline: true,
		})
	}

	if auctionMap["bid"] != nil {
		if auctionMap["winner"] != nil {
			auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
				Name:   "__**Current Highest Bid:**__",
				Value:  PriceFormat(auctionMap["bid"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]),
				Inline: true,
			})
		} else {
			auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
				Name:   "__**Starting Bid:**__",
				Value:  PriceFormat(auctionMap["bid"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]),
				Inline: true,
			})
		}
	}

	if auctionMap["winner"] != nil {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Current Winner**__",
			Value:  fmt.Sprintf("<@%s>", auctionMap["winner"]),
			Inline: true,
		})
	}

	if eventType == EventTypeAuction {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**How to Bid**__",
			Value:  "Use the command `/bid` below.\n• Ex: `/bid 550`.\n**Alternate Method:** reply to this auction or @ the bot with `bid <amount>`",
			Inline: false,
		})
	}

	if eventType == EventTypeAuction {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**How to Enter**__",
			Value:  "To enter, select the 🎁 reaction below! Removing your reaction will remove your entry.",
			Inline: false,
		})
	}

	guild, err := s.Guild(auctionMap["guild_id"].(string))
	if err != nil {
		fmt.Println("Error fetching guild: ", err)
		return h.PresetResponse{}, err
	}

	if auctionMap["alert_role"] != nil {
		content = fmt.Sprintf("<@&%s>", strings.Trim(auctionMap["alert_role"].(string), " "))
	}

	components := []discordgo.MessageComponent{}

	if eventType == EventTypeAuction {
		components = []discordgo.MessageComponent{
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
		}
	}

	if eventType == EventTypeShop {
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Purchase",
						Style:    discordgo.SuccessButton,
						CustomID: "additem",
					},
				},
			},
		}
	}

	if auctionMap["bid_history"] != nil {
		if len(auctionMap["bid_history"].(string)) > 4095 {
			auctionMap["bid_history"] = auctionMap["bid_history"].(string)[len(auctionMap["bid_history"].(string))-4095:]
		}
		embeds = []*discordgo.MessageEmbed{{
			Title:       "**Bid History**",
			Description: auctionMap["bid_history"].(string),
			Color:       0x8073ff,
			Image: &discordgo.MessageEmbedImage{
				URL: "https://i.imgur.com/9wo7diC.png",
			},
		}}
	}

	return h.PresetResponse{
		Content:    content,
		Title:      fmt.Sprintf("%s Item: __**%s**__", eventType, auctionMap["item"]),
		Fields:     auctionfields,
		Thumbnail:  &discordgo.MessageEmbedThumbnail{URL: guild.IconURL()},
		Image:      &discordgo.MessageEmbedImage{URL: imageURL},
		Components: components,
		Embeds:     embeds,
		Files:      []*discordgo.File{},
	}, nil
}

//Buttons

func AuctionBidHistory(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	claimMap := map[string]interface{}{}

	result := database.DB.Model(database.Claim{}).First(claimMap, i.Message.ID)

	if result.Error != nil {
		return result.Error
	}

	if claimMap["bid_history"] == nil {
		return fmt.Errorf("No bid history found for this auction.")
	}

	bidHistory := claimMap["bid_history"].(string)

	if len(bidHistory) > 4095 {
		bidHistory = bidHistory[len(bidHistory)-4095:]
	}

	return h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Bid History**",
		Description: bidHistory,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/9wo7diC.png",
		},
	})

}

//Extra Responses

func AuctionAutoComplete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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
		return err
	}

	return nil
}

func TimeSuggestions(input string) []*discordgo.ApplicationCommandOptionChoice {

	re := regexp.MustCompile(`[0-9]+\.?[0-9]*`)
	input = re.FindString(input)

	if input == "" {
		return []*discordgo.ApplicationCommandOptionChoice{}
	}

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

func AuctionEnd(s *discordgo.Session, channelID, guildID string) error {

	auctionMap := map[string]any{}
	auctionSetup := map[string]any{}

	result := database.DB.Model(database.AuctionSetup{}).First(&auctionSetup, guildID)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	result = database.DB.Model(database.Auction{}).First(&auctionMap, channelID)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	//Pause auction ending until end time if the auction is not over yet.
	if auctionMap["end_time"].(time.Time).After(time.Now()) {
		time.Sleep(time.Until(auctionMap["end_time"].(time.Time)))
		err := AuctionEnd(s, channelID, guildID)
		return err
	}

	message := discordgo.NewMessageEdit(channelID, auctionMap["message_id"].(string))
	messageEmbeds, err := s.ChannelMessage(channelID, auctionMap["message_id"].(string))
	if err != nil {
		result = database.DB.Delete(database.Auction{
			ChannelID: auctionMap["channel_id"].(string),
		})
		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		return err
	}

	if auctionSetup["log_channel"] == nil {
		fmt.Println("Log channel has not been set for guild: " + guildID)
		_, err := h.ErrorMessage(s, channelID, "Auction cannot end because log channel has not been set. Please setup an auction log using `/settings auction`. You might need to end the auction manually after setting the channel.")
		if err != nil {
			return err
		}
		return err
	}

	if auctionMap["target_price"] != nil && auctionMap["target_price"].(float64) > auctionMap["bid"].(float64) {
		auctionMap["target_message"] = fmt.Sprintf("The host had set a target price of %s that has not been reached.", PriceFormat(auctionMap["target_price"].(float64), guildID, auctionMap["currency"]))
		delete(auctionMap, "winner")
	}

	if auctionMap["buyout"] != nil && auctionMap["buyout"].(float64) != 0 {
		if auctionMap["bid"].(float64) >= auctionMap["buyout"].(float64) {
			auctionMap["buyout_message"] = fmt.Sprintf("%s\n\u200bBUYOUT!", PriceFormat(auctionMap["buyout"].(float64), auctionMap["guild_id"].(string), auctionMap["currency"]))
		}
	}

	auctionMap["log_channel"] = auctionSetup["log_channel"]

	err = ClaimOutput(s, auctionMap, "Auction")
	if err != nil {
		return fmt.Errorf("Claim Output Error: " + err.Error())
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
		_, err = s.ChannelMessageEditComplex(message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if auctionMap["channel_override"] == nil {
		time.Sleep(30 * time.Second)
		_, err = s.ChannelDelete(channelID)
		if err != nil {
			fmt.Println(err)
		}
	}

	result = database.DB.Delete(database.Auction{}, channelID)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

	return nil
}

func AuctionQueue(s *discordgo.Session, i *discordgo.InteractionCreate) error {

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
			Name:  fmt.Sprintf("**%s. %s**", fmt.Sprint(len(fields)+1), v.Item),
			Value: fmt.Sprintf("**Start time:** <t:%d:R>\n**End Time:** <t:%d>\n**Starting Price:** %s\n\u200b", v.StartTime.Unix(), v.EndTime.Unix(), PriceFormat(v.Bid, v.GuildID, v.Currency)),
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
					MinValues:   h.Ptr(0),
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

	return nil
}

func AuctionEndButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {

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
		return nil
	}
	database.DB.First(&auctionSetup, i.GuildID)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	if auctionMap["log_channel"] == "" {
		return fmt.Errorf("Auction cannot end because log channel has not been set. Please setup an auction log using `/settings auction`")
	}

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != auctionMap["host"] {
		return fmt.Errorf("User must have administrator permissions to run this command")
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

	err = AuctionEnd(s, i.ChannelID, i.GuildID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAuctionChannel(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	_, err := s.ChannelDelete(i.ChannelID)
	if err != nil {
		fmt.Println(err)
		h.ErrorResponse(s, i, err.Error())
	}
	return nil
}

func ClearAuctionButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	var auctionInfo database.Auction

	result := database.DB.First(&auctionInfo, i.ChannelID)
	if result.Error != nil {
		fmt.Println("Error checking for auction. Auction has likely ended")
		auctionInfo.Host = ""
	}

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != auctionInfo.Host {
		return fmt.Errorf("User must be host or have administrator permissions to run this command")
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{},
	})

	for {
		messageIDs := make([]string, 0)
		messages, err := s.ChannelMessages(i.ChannelID, 100, "", i.Message.ID, "")
		if err != nil {
			h.DeferredErrorResponse(s, i, err.Error())
			return nil
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
			h.DeferredErrorResponse(s, i, err.Error())
			return nil
		}
	}

	_, err := h.DeferredResponse(s, i, h.PresetResponse{
		Title:       "Success!",
		Description: "Clearing Chat",
	})

	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func DeleteAuctionQueue(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	if i.Member.Permissions&(1<<3) != 8 {
		return fmt.Errorf("User must have administrator permissions to run this command")
	}

	IDs := i.MessageComponentData().Values

	for _, v := range IDs {
		database.DB.Delete(database.AuctionQueue{}, v)
	}

	AuctionQueue(s, i)
	return nil
}

func AuctionSetupClearButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {

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
		Description: "You have successfully cleared the following settings. Run `/settings auction` to see your changes.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Cleared Settings**",
				Value: clearedSettings,
			},
		},
	})
	return nil
}
