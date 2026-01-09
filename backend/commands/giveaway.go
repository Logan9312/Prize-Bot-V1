package commands

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gitlab.com/logan9312/discord-auction-bot/logger"
	"gorm.io/gorm"
)

var GiveawayCommand = discordgo.ApplicationCommand{
	Name:        "giveaway",
	Description: "Start a giveaway!",
	Options: []*discordgo.ApplicationCommandOption{
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
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "duration",
					Description:  "The amount of time that the giveaway will run for. (Example: 1d, 12h, etc.)",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel where the giveaway will take place",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
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
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "delete",
			Description: "Delete a giveaway",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message_id",
					Description: "The message ID of the giveaway to delete",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel containing the giveaway",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
						5,
					},
				},
			},
		},
	},
}

func Giveaway(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Options[0].Name {
	case "create":
		return GiveawayCreate(s, i)
	case "delete":
		return GiveawayDelete(s, i)
	}
	return fmt.Errorf("Unknown Giveaway command, please contact support")
}

func GiveawayCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	GiveawaySetup := map[string]interface{}{}

	result := database.DB.Model(database.GiveawaySetup{}).First(&GiveawaySetup, i.GuildID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Sugar.Errorw("failed to fetch giveaway setup", "guild_id", i.GuildID, "error", result.Error)
		return fmt.Errorf("failed to load giveaway settings. Please try again or contact support")
	}

	// Check if user has permission to host giveaways
	if !CheckHostPermission(i.Member, GiveawaySetup["host_role"], i.GuildID) {
		return errors.New(GetHostRoleErrorMessage(GiveawaySetup["host_role"], "host giveaways"))
	}

	giveawayMap := h.ParseSubCommand(i)

	giveawayMap["guild_id"] = i.GuildID
	giveawayMap["host"] = i.Member.User.ID
	giveawayMap["channel_id"] = giveawayMap["channel"]
	delete(giveawayMap, "channel")

	if giveawayMap["channel_id"] == nil {
		giveawayMap["channel_id"] = i.ChannelID
	}

	endTimeDuration, err := h.ParseTime(strings.ToLower(giveawayMap["duration"].(string)))
	if err != nil {
		return err
	}
	if giveawayMap["image"] != nil {
		giveawayMap["image_url"] = i.ApplicationCommandData().Resolved.Attachments[giveawayMap["image"].(string)].URL
	}

	delete(giveawayMap, "duration")
	delete(giveawayMap, "image")

	if getWinnersCount(giveawayMap["winners"]) < 1 {
		return fmt.Errorf("Must have 1 or more winners.")
	}

	giveawayMap["end_time"] = time.Now().Add(endTimeDuration)
	giveawayMap["alert_role"] = GiveawaySetup["alert_role"]

	formattedMessage, err := EventFormat(s, giveawayMap, EventTypeGiveaway, i.GuildID)
	if err != nil {
		return err
	}

	message, err := h.PresetMessageSend(s, giveawayMap["channel_id"].(string), formattedMessage)

	if err != nil {
		return err
	}

	delete(giveawayMap, "alert_role")
	giveawayMap["message_id"] = message.ID

	err = s.MessageReactionAdd(giveawayMap["channel_id"].(string), message.ID, "🎁")

	if err != nil {
		return err
	}

	result = database.DB.Model(database.Giveaway{}).Create(giveawayMap)
	if result.Error != nil {
		logger.Sugar.Errorw("failed to save giveaway to database", "guild_id", giveawayMap["guild_id"], "item", giveawayMap["item"], "error", result.Error)
		return fmt.Errorf("giveaway could not be saved. Please contact support - your giveaway may not function properly")
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Giveaway Started!**",
		Description: fmt.Sprintf("Giveaway has started in <#%s>", giveawayMap["channel_id"]),
	})

	time.Sleep(endTimeDuration)
	err = GiveawayEnd(s, message.ID)
	if err != nil {
		logger.Sugar.Errorw("giveaway end failed", "message_id", message.ID, "error", err)
		h.ErrorMessage(s, giveawayMap["channel_id"].(string), "An error occurred ending the giveaway. Please contact support.")
	}
	return nil
}

func GiveawayDelete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	giveawaySetup := map[string]interface{}{}
	result := database.DB.Model(database.GiveawaySetup{}).First(&giveawaySetup, i.GuildID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Sugar.Errorw("failed to fetch giveaway setup for delete", "guild_id", i.GuildID, "error", result.Error)
		return fmt.Errorf("failed to load giveaway settings. Please try again or contact support")
	}

	// Check if user has permission to delete giveaways
	if !CheckHostPermission(i.Member, giveawaySetup["host_role"], i.GuildID) {
		return errors.New(GetHostRoleErrorMessage(giveawaySetup["host_role"], "delete giveaways"))
	}
	
	// Parse command options
	options := h.ParseSubCommand(i)
	messageID := options["message_id"].(string)
	channelID := options["channel"]
	
	if channelID == nil {
		channelID = i.ChannelID
	} else {
		channelID = channelID.(string)
	}
	
	// Find giveaway in database
	giveawayInfo := database.Giveaway{
		MessageID: messageID,
	}
	
	result = database.DB.First(&giveawayInfo)
	if result.Error != nil {
		return fmt.Errorf("Giveaway not found. Please check the message ID and try again.")
	}
	
	// Delete the giveaway message if possible
	err := s.ChannelMessageDelete(channelID.(string), messageID)
	if err != nil {
		logger.Sugar.Warnw("error deleting giveaway message", "error", err)
		// Continue even if message delete fails - we still want to remove from DB
	}
	
	// Delete from database
	result = database.DB.Delete(&giveawayInfo)
	if result.Error != nil {
		logger.Sugar.Errorw("failed to delete giveaway from database", "message_id", messageID, "error", result.Error)
		return fmt.Errorf("failed to delete giveaway. Please try again or contact support")
	}
	
	// Send success response
	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Giveaway Deleted!**",
		Description: "The giveaway has been successfully deleted.",
	})
	
	return nil
}

func GiveawayAutoComplete(s *discordgo.Session, i *discordgo.InteractionCreate) error {
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
			return err
		}
	}
	return nil
}

func GiveawayEnd(s *discordgo.Session, messageID string) error {

	giveawayMap := map[string]interface{}{}
	giveawaySetup := map[string]interface{}{}

	result := database.DB.Model(database.Giveaway{}).First(&giveawayMap, messageID)
	if result.Error != nil {
		return result.Error
	}

	result = database.DB.Model(database.GiveawaySetup{}).First(&giveawaySetup, giveawayMap["guild_id"].(string))
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Sugar.Warnw("error fetching giveaway setups", "error", result.Error)
	}

	fm, err := EventFormat(s, giveawayMap, EventTypeGiveaway, giveawayMap["guild_id"].(string))
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
		// Fetch reactions in batches of 100 (Discord's max per request)
		// Loop continues until all reactions are collected via pagination
		users, err := s.MessageReactions(giveawayMap["channel_id"].(string), messageID, "🎁", 100, "", afterID)
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

	winnerList, err := GiveawayRoll(s, entrants, giveawayMap)
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

	winnersCount := getWinnersCount(giveawayMap["winners"])
	if len(entrants) < winnersCount {
		formattedWinnerList += fmt.Sprintf("• Only %d users entered out of a maximum of %d winners.\n", len(entrants), winnersCount)
	}

	fm.Fields = append(fm.Fields, &discordgo.MessageEmbedField{
		Name:   "**Giveaway Complete! Winners:**",
		Value:  formattedWinnerList,
		Inline: false,
	})

	fm.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Reroll",
					Style:    1,
					Disabled: false,
					Emoji: &discordgo.ComponentEmoji{
						Name:     "🔄",
						ID:       "",
						Animated: false,
					},
					CustomID: "reroll_giveaway",
				},
			},
		},
	}

	_, err = h.SuccessMessageEdit(s, giveawayMap["channel_id"].(string), messageID, fm)
	if err != nil {
		return err
	}

	giveawayMap["log_channel"] = giveawaySetup["log_channel"]

	// Process claims for each winner (not in transaction since we want partial success)
	for _, v := range winnerList {
		// Create a copy of the map for each winner to avoid race conditions
		winnerMap := make(map[string]interface{})
		for k, val := range giveawayMap {
			winnerMap[k] = val
		}
		winnerMap["winner"] = v
		claimErr := ClaimOutput(s, winnerMap, "Giveaway")
		if claimErr != nil {
			logger.Sugar.Errorw("failed to create claim for giveaway winner", "winner", v, "message_id", messageID, "error", claimErr)
			h.ErrorMessage(s, giveawayMap["channel_id"].(string), fmt.Sprintf("An error occurred creating a claim for winner <@%s>. Please contact support.", v))
		}
	}

	// Update giveaway as finished
	result = database.DB.Model(database.Giveaway{
		MessageID: messageID,
	}).Update("finished", true)
	if result.Error != nil {
		logger.Sugar.Errorw("error saving giveaway finished status", "message_id", messageID, "error", result.Error)
		return result.Error
	}

	return nil
}

func GiveawayRoll(s *discordgo.Session, entries []string, giveawayMap map[string]interface{}) ([]string, error) {

	winnerList := []string{}

	if len(entries) == 0 {
		return winnerList, fmt.Errorf("No entries found.")
	}
	logger.Sugar.Debug("rolling giveaway")

	// Create a new random source with a time-based seed
	// This is the modern approach as of Go 1.20+ (rand.Seed is deprecated)
	rSource := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rSource)

	winnersCount := getWinnersCount(giveawayMap["winners"])
	for n := 0; n < winnersCount; {

		if len(entries) == 0 {
			break
		}

		printlist := ""
		for i, v := range entries {
			username := ""
			u, err := s.User(v)
			if err == nil {
				username = u.Username
			}
			printlist += fmt.Sprintf("Entry %d: <@%s> (%s)", i, v, username)
		}

		logger.Sugar.Debugw("giveaway entries", "entries", printlist)

		// Use the local random generator instance instead of the global one
		index := rng.Intn(len(entries))
		logger.Sugar.Debugw("giveaway winner selected", "index", index)
		winnerID := entries[index]

		if len(entries) > 1 {
			entries = append(entries[:index], entries[index+1:]...)
		} else {
			entries = []string{}
		}

		winnerList = append(winnerList, winnerID)

		n++
	}

	return winnerList, nil
}

func RerollGiveawayButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	giveawayInfo := database.Giveaway{
		MessageID: i.Message.ID,
	}

	result := database.DB.First(&giveawayInfo)
	if result.Error != nil {
		return result.Error
	}

	if i.Member.Permissions&(1<<3) == 0 && i.Member.User.ID != giveawayInfo.Host {
		return fmt.Errorf("User must be host or have administrator permissions to run this command")
	}

	err := GiveawayEnd(s, i.Message.ID)
	if err != nil {
		return err
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Reroll Successful!*",
		Description: "New winners have been selected.",
	})
	return nil
}

func GiveawaySetupClearButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return GenericSetupClear(s, i, &database.GiveawaySetup{GuildID: i.GuildID}, SetupClearConfig{
		SetupType: "Giveaway",
		SetupCmd:  "/giveaway setup",
	})
}
