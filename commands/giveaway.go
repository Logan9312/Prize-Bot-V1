package commands

import (
	random "crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var GiveawayCommand = discordgo.ApplicationCommand{
	Name:        "giveaway",
	Description: "Start a giveaway!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "setup",
			Description: "Setup giveaways",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "alert_role",
					Description: "Set a role to get pinged whenever an auction starts.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "log_channel",
					Description: "Sets the channel where giveaway will send outputs when they end",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						0,
						5,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "host_role",
					Description: "Set a role that can host auctions.",
				},
				/*{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "auto_enter",
					Description: "Anyone with this role will be automatically entered.",
				},*/
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Start a giveaway!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The item you want to give away",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "winners",
					Description: "The number of winners.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "duration",
					Description: "The amount of time that the giveaway will run for. (Example: 1d, 12h, etc.)",
					Required:    true,
					//Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel where the giveaway will take place",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						0,
						5,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "Description of the giveaway item.",
					Required:    false,
				},
				{
					Type:        11,
					Name:        "image",
					Description: "Attach an image to your giveaway",
					Required:    false,
				},
			},
		},
	},
}

func Giveaway(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[0].Name {
	case "setup":
		GiveawaySetup(s, i)
	case "create":
		GiveawayCreate(s, i)
	}
}

func GiveawaySetup(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var err error

	options := h.ParseSubCommand(i)
	content := ""

	info := database.GiveawaySetup{
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
	}

	result = database.DB.Model(&info).Updates(options)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	//Now check what options are set
	setOptions := map[string]interface{}{}

	database.DB.Model(database.GiveawaySetup{}).First(&setOptions, i.GuildID)

	responseFields := []*discordgo.MessageEmbedField{}

	for _, v := range GiveawayCommand.Options[0].Options {
		switch {
		case setOptions[v.Name] == "", setOptions[v.Name] == 0, setOptions[v.Name] == nil:
			setOptions[v.Name] = "Not Set"
		case strings.Contains(v.Name, "role"):
			setOptions[v.Name] = fmt.Sprintf("<@&%s>", setOptions[v.Name])
		case strings.Contains(v.Name, "channel"):
			setOptions[v.Name] = fmt.Sprintf("<#%s>", setOptions[v.Name])
		}
		responseFields = append(responseFields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("**%s**", strings.Title(strings.ReplaceAll(v.Name, "_", " "))),
			Value: setOptions[v.Name].(string),
		})
	}

	menuOptions := []discordgo.SelectMenuOption{}

	for _, v := range GiveawayCommand.Options[0].Options {
		menuOptions = append(menuOptions, discordgo.SelectMenuOption{
			Label:       strings.Title(strings.ReplaceAll(v.Name, "_", " ")),
			Value:       v.Name,
			Description: v.Description,
		})
	}

	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Giveaway Setup",
		Description: content,
		Fields:      responseFields,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "clear_giveaway_setup",
						Placeholder: "Clear Setup Options",
						MinValues:   1,
						MaxValues:   len(GiveawayCommand.Options[0].Options),
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

func GiveawayCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	canHost := false

	GiveawaySetup := map[string]interface{}{}

	result := database.DB.Model(database.GiveawaySetup{}).First(&GiveawaySetup, i.GuildID)
	if result.Error != nil {
		fmt.Println("Error fetching Giveaway Setup DB")
	}

	if GiveawaySetup["host_role"] != "" && GiveawaySetup["host_role"] != nil {
		for _, v := range i.Member.Roles {
			if v == GiveawaySetup["host_role"].(string) {
				canHost = true
			}
		}
		if i.Member.Permissions&(1<<3) == 8 {
			canHost = true
		}
		if !canHost {
			h.ErrorResponse(s, i, "User must be administrator or have the role <@&"+GiveawaySetup["host_role"].(string)+"> to host giveaways.")
			return
		}
	}

	giveawayMap := h.ParseSubCommand(i)

	giveawayMap["guild_id"] = i.GuildID
	giveawayMap["host"] = i.Member.User.ID

	if giveawayMap["channel_id"] == nil {
		giveawayMap["channel_id"] = i.ChannelID
	}

	endTimeDuration, err := h.ParseTime(strings.ToLower(giveawayMap["duration"].(string)))
	if err != nil {
		h.ErrorResponse(s, i, err.Error())
		return
	}
	if giveawayMap["image"] != nil {
		giveawayMap["image_url"] = i.ApplicationCommandData().Resolved.Attachments[giveawayMap["image"].(string)].URL
	}

	delete(giveawayMap, "duration")
	delete(giveawayMap, "image")

	if giveawayMap["winners"].(float64) < 1 {
		h.ErrorResponse(s, i, "Must have 1 or more winners.")
		return
	}

	giveawayMap["end_time"] = time.Now().Add(endTimeDuration)

	formattedMessage, err := AuctionFormat(s, giveawayMap, "Giveaway")
	if err != nil {
		h.ErrorResponse(s, i, err.Error())
		return
	}

	message, err := h.PresetMessageSend(s, giveawayMap["channel_id"].(string), formattedMessage)

	if err != nil {
		h.ErrorResponse(s, i, err.Error())
		return
	}

	giveawayMap["message_id"] = message.ID

	err = s.MessageReactionAdd(giveawayMap["channel_id"].(string), message.ID, "🎁")

	if err != nil {
		h.ErrorResponse(s, i, err.Error())
		return
	}

	delete(giveawayMap, "channel")

	result = database.DB.Model(database.Giveaway{}).Create(giveawayMap)
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Giveaway Started!**",
		Description: fmt.Sprintf("Giveaway has started in <#%s>", giveawayMap["channel_id"]),
	})

	time.Sleep(endTimeDuration)
	err = GiveawayEnd(s, message.ID)
	if err != nil {
		h.ErrorMessage(s, giveawayMap["channel_id"].(string), err.Error())
	}
}

func GiveawayAutoComplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var choices []*discordgo.ApplicationCommandOptionChoice
	switch i.ApplicationCommandData().Options[0].Name {
	case "create":
		options := h.ParseSubCommand(i)

		if options["duration"] != nil {
			choices = TimeSuggestions(options["duration"].(string))
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

func GiveawayEnd(s *discordgo.Session, messageID string) error {

	giveawayMap := map[string]interface{}{}
	giveawaySetup := map[string]interface{}{}

	result := database.DB.Model(database.Giveaway{}).First(&giveawayMap, messageID)
	if result.Error != nil {
		return result.Error
	}

	result = database.DB.Model(database.GiveawaySetup{}).First(&giveawaySetup, giveawayMap["guild_id"].(string))
	if result.Error != nil {
		fmt.Println("Error fetching giveaway setups", result.Error)
	}

	fm, err := AuctionFormat(s, giveawayMap, "Giveaway")
	if err != nil {
		return err
	}

	if giveawaySetup["log_channel"] == nil {
		giveawaySetup["log_channel"] = giveawayMap["channel_id"]
		fm.Fields = append(fm.Fields, &discordgo.MessageEmbedField{
			Name:   "⚙ **Please Note:**",
			Value:  "No logging channel was set, so the giveaway end output defaulted to the current channel. To fix this for future giveaways please run `/claim setup log_channel`",
			Inline: false,
		})
	}

	entrants := []string{}
	afterID := ""
	for {

		users, err := s.MessageReactions(giveawayMap["channel_id"].(string), messageID, "🎁", 0, "", afterID)
		if err != nil {
			return err
		}

		for _, v := range users {
			if v.ID != s.State.User.ID {
				entrants = append(entrants, v.ID)
			}
		}

		if len(users) < 100 {
			break
		}
		afterID = users[len(users)-1].ID
	}

	winnerList, err := GiveawayRoll(entrants, giveawayMap)
	if err != nil {
		return err
	}

	formattedWinnerList := ""

	for _, v := range winnerList {
		user, err := s.User(v)
		if err != nil {
			return err
		}
		formattedWinnerList += fmt.Sprintf("• %s (%s#%s)\n", user.Mention(), user.Username, user.Discriminator)
	}

	if len(entrants) < int(giveawayMap["winners"].(float64)) {
		formattedWinnerList += fmt.Sprintf("• Only %d users entered out of a maximum of %d winners.\n", len(entrants), int(giveawayMap["winners"].(float64)))
	}

	fm.Fields = append(fm.Fields, &discordgo.MessageEmbedField{
		Name:   "**Giveaway Complete! Winners:**",
		Value:  formattedWinnerList,
		Inline: false,
	})

	_, err = h.SuccessMessageEdit(s, giveawayMap["channel_id"].(string), messageID, fm)
	if err != nil {
		return err
	}

	giveawayMap["log_channel"] = giveawaySetup["log_channel"]

	for _, v := range winnerList {
		giveawayMap["winner"] = v
		err = ClaimOutput(s, giveawayMap, "Giveaway")
		if err != nil {
			h.ErrorMessage(s, giveawayMap["channel_id"].(string), err.Error())
		}
	}

	result = database.DB.Model(database.Giveaway{
		MessageID: messageID,
	}).Update("finished", true)
	if result.Error != nil {
		fmt.Println("Error saving giveaway finished status.", result.Error)
	}

	time.Sleep(24 * time.Hour)
	database.DB.Delete(database.Giveaway{}, messageID)

	return nil
}

func GiveawayRoll(entries []string, giveawayMap map[string]interface{}) ([]string, error) {

	winnerList := []string{}

	if len(entries) == 0 {
		return winnerList, fmt.Errorf("No entries found.")
	}

	for n := float64(0); n < giveawayMap["winners"].(float64); {

		if len(entries) == 0 {
			break
		}

		bigInt, err := random.Int(random.Reader, big.NewInt(int64(len(entries))))
		if err != nil {
			fmt.Println("Random Error: ", err)
			return winnerList, err
		}

		winnerID := entries[bigInt.Int64()]

		entries[bigInt.Int64()] = entries[len(entries)-1]

		if len(entries) >= 2 {
			entries = entries[0 : len(entries)-2]
		} else {
			entries = []string{}
		}

		winnerList = append(winnerList, winnerID)

		n++
	}

	return winnerList, nil
}

func RerollGiveawayButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	giveawayInfo := database.Giveaway{
		MessageID: i.Message.ChannelID,
	}

	result := database.DB.First(&giveawayInfo)
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != giveawayInfo.Host {
		h.ErrorResponse(s, i, "User must be host or have administrator permissions to run this command")
		return
	}

	err := GiveawayEnd(s, i.Message.ID)
	if err != nil {
		h.ErrorResponse(s, i, err.Error())
		return
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Reroll Successful!*",
		Description: "New winners have been selected.",
	})
}

func GiveawaySetupClearButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.MessageComponentData().Values

	clearedValues := map[string]interface{}{}

	clearedValues["guild_id"] = i.GuildID
	info := database.GiveawaySetup{
		GuildID: i.GuildID,
	}

	clearedSettings := "No Settings Cleared!"
	if len(options) > 0 {
		clearedSettings = ""
	}

	//Might need editing
	for _, v := range options {
		clearedValues[v] = gorm.Expr("NULL")
		clearedSettings += fmt.Sprintf("• %s\n", strings.Title(strings.ReplaceAll(v, "_", " ")))
	}

	database.DB.Model(&info).Updates(clearedValues)

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Cleared Giveaway Settings**",
		Description: "You have successfully cleared the following settings. Run `/giveaway setup` to see your changes.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Cleared Settings**",
				Value: clearedSettings,
			},
		},
	})
}
