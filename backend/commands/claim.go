package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gitlab.com/logan9312/discord-auction-bot/logger"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ClaimCommand = discordgo.ApplicationCommand{
	Name:        "claim",
	Description: "Manage your claims!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create a claimable prize.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionMentionable,
					Name:        "target",
					Description: "The user who will receive the prize.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The prize you are giving",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "log_channel",
					Description: "The output channel for the prizes.",
					ChannelTypes: []discordgo.ChannelType{
						0,
						5,
					},
					Required: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "Set a custom item description",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "inventory",
			Description: "Displays a user's unclaimed prizes",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user who you want to display.",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "refresh",
			Description: "Resends all claim messages. Will not delete old ones",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The output channel for the claims..",
					ChannelTypes: []discordgo.ChannelType{
						0,
						5,
					},
					Required: true,
				},
			},
		},
	},
	DefaultMemberPermissions: h.Ptr(int64(discordgo.PermissionManageServer)),
}

func Claim(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	logger.Sugar.Debugw("claim subcommand invoked", "subcommand", i.ApplicationCommandData().Options[0].Name)
	switch i.ApplicationCommandData().Options[0].Name {
	case "create":
		return ClaimCreate(s, i)
	case "inventory":
		return ClaimInventory(s, i)
	case "refresh":
		return claimRefresh(s, i)
	}

	return fmt.Errorf("Unknown Claim command, please contact support")
}

func ClaimSetupClearButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return GenericSetupClear(s, i, &database.ClaimSetup{GuildID: i.GuildID}, SetupClearConfig{
		SetupType: "Claim",
		SetupCmd:  "/claim setup",
	})
}

func ClaimCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	claimMap := h.ParseSubCommand(i)

	claimMap["host"] = i.Member.User.ID
	claimMap["guild_id"] = i.GuildID

	if i.ApplicationCommandData().Resolved.Roles[claimMap["target"].(string)] != nil {
		if !CheckPremiumGuild(i.GuildID) {
			h.PremiumError(s, i, "Giving out a claim to an entire role is restricted to premium users. Please select an user instead or upgrade to premium.")
			return nil
		}

		claimMap["interaction"] = i

		h.SaveChunkData(i.ID, claimMap)

		err := s.RequestGuildMembers(i.GuildID, "", 0, "claim_create:"+i.ID, false)
		if err != nil {
			return err
		}

		err = h.ExperimentalResponse(s, i, h.PresetResponse{
			Title:       "Claims are being created!",
			Description: "Check out <#" + claimMap["log_channel"].(string) + "> to see the claims. This might take a while.",
		})
		if err != nil {
			logger.Sugar.Warnw("claim operation error", "error", err)
		}
	}

	if i.ApplicationCommandData().Resolved.Users[claimMap["target"].(string)] != nil {
		claimMap["winner"] = claimMap["target"].(string)
		err := ClaimOutput(s, claimMap, "Custom Claim")
		if err != nil {
			return err
		}
		err = h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "Claim Successfully Created!",
			Description: "Check out <#" + claimMap["channel_id"].(string) + "> to see the claim",
		})
		if err != nil {
			logger.Sugar.Warnw("claim operation error", "error", err)
		}
	}

	return nil
}

func ClaimCreateRole(s *discordgo.Session, g *discordgo.GuildMembersChunk) error {

	details := strings.Split(g.Nonce, ":")

	claimMap := h.ReadChunkData(details[1])

	for _, v := range g.Members {
		if g.GuildID != claimMap["target"].(string) && !HasRole(v, claimMap["target"].(string)) {
			continue
		}
		if v.User.Bot {
			continue
		}

		claimMap["winner"] = v.User.ID
		err := ClaimOutput(s, claimMap, "Custom Claim")
		if err != nil {
			h.FollowUpErrorResponse(s, claimMap["interaction"].(*discordgo.InteractionCreate), fmt.Sprintf("There was an issue creating a claim for <@%s>. Error Message: %s", v.User.ID, err))
		}
	}

	h.FollowUpSuccessResponse(s, claimMap["interaction"].(*discordgo.InteractionCreate), h.PresetResponse{
		Title:       "__**Claim Create Role**__",
		Description: fmt.Sprintf("Claims are currently being created for all users in <@&%s>", claimMap["role"]),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Progress**",
				Value:  fmt.Sprintf("`%d`/`%d` chunks completed", g.ChunkIndex+1, g.ChunkCount),
				Inline: false,
			},
		},
	})
	return nil
}

// ClaimOutput Seems like using a map here overcomplicates it. Possibly need to go back to fix if I run into issues.
// If tx is nil, uses the global database connection. Otherwise, uses the provided transaction.
func ClaimOutput(s *discordgo.Session, claimMap map[string]interface{}, eventType string) error {
	return ClaimOutputWithTx(s, claimMap, eventType, nil)
}

// ClaimOutputWithTx allows passing a transaction for atomic operations
func ClaimOutputWithTx(s *discordgo.Session, claimMap map[string]interface{}, eventType string, tx *gorm.DB) error {

	p := message.NewPrinter(language.English)
	mentionUser := ""
	finalBid := ""
	primaryKey := ""
	claimMap["cost"] = claimMap["bid"]
	claimSetup := map[string]interface{}{}

	result := database.DB.Model(database.ClaimSetup{}).First(&claimSetup, claimMap["guild_id"].(string))
	if result.Error != nil {
		logger.Sugar.Warnw("claim setup database error", "error", result.Error)
	}

	if claimMap["formatted_price"] != nil {
		finalBid = p.Sprint(claimMap["formatted_price"].(string))
	}
	if claimMap["buyout_message"] != nil {
		finalBid = p.Sprint(claimMap["buyout_message"].(string))
	}
	if claimMap["target_message"] != nil {
		finalBid = p.Sprint(claimMap["target_message"].(string))
	}
	if claimMap["cost"] != nil {
		finalBid = p.Sprint(claimMap["cost"].(float64))
	}
	if claimMap["winner"] == nil {
		claimMap["winner"] = "No Winner Detected"
	}

	buttons := []discordgo.MessageComponent{
		/*discordgo.Button{
			Label: "Support",
			Style: discordgo.LinkButton,
			Emoji: discordgo.ComponentEmoji{
				Name:     "logo",
				ID:       "889025400120950804",
				Animated: false,
			},
			URL: "https://discord.gg/RxP2z5NGtj",
		},*/
	}

	if claimMap["bid_history"] != nil {
		buttons = append(buttons, discordgo.Button{
			Label: "Bids",
			Style: 2,
			Emoji: &discordgo.ComponentEmoji{
				Name: "📜",
			},
			CustomID: "bid_history",
		})
	}

	user, err := s.User(fmt.Sprint(claimMap["winner"]))
	if err != nil {
		user = &discordgo.User{}
		user.Username = claimMap["winner"].(string)
	} else {
		mentionUser = user.Mention()
		user.Username = fmt.Sprintf("%s (%s#%s)", user.Mention(), user.Username, user.Discriminator)
		if claimSetup["disable_claiming"] != true {
			buttons = append(buttons, discordgo.Button{
				Label: "Claim!",
				Style: 3,
				Emoji: &discordgo.ComponentEmoji{
					Name: "cryopod",
					ID:   "889307390690885692",
				},
				CustomID: "claim_prize:" + claimMap["winner"].(string),
			})
		}
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: buttons,
		},
	}

	if len(buttons) == 0 {
		components = []discordgo.MessageComponent{}
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "**Winner**",
			Value:  user.Username,
			Inline: true,
		},
		{
			Name:   "**Host**",
			Value:  fmt.Sprintf("This Auction was hosted by: <@!%s>", claimMap["host"]),
			Inline: true,
		},
	}

	if claimMap["cost"] != nil && claimMap["cost"] != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**Final Bid**",
			Value:  finalBid,
			Inline: true,
		})
	}

	if claimMap["description"] != nil && claimMap["description"] != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**Description**",
			Value:  claimMap["description"].(string),
			Inline: true,
		})
	}

	if claimMap["image_url"] == nil {
		claimMap["image_url"] = ""
	}

	if claimMap["log_channel"] == nil {
		return fmt.Errorf("No logging channel set.")
	}

	message, err := h.SuccessMessage(s, claimMap["log_channel"].(string), h.PresetResponse{
		Content: mentionUser,
		Title:   fmt.Sprintf("%s Prize: __**%s**__", eventType, claimMap["item"]),
		Fields:  fields,
		/*Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: claimMap["image_url"].(string),
		},*/
		Components: components,
	})

	if err != nil {
		// Check for common Discord API errors and provide user-friendly messages
		errStr := err.Error()
		if strings.Contains(errStr, "Unknown Channel") || strings.Contains(errStr, "10003") {
			return fmt.Errorf("claim channel no longer exists. Please update your claim settings at https://prizebot.info/dashboard")
		}
		if strings.Contains(errStr, "Missing Access") || strings.Contains(errStr, "50001") {
			return fmt.Errorf("bot doesn't have access to the claim channel. Please check channel permissions")
		}
		if strings.Contains(errStr, "Missing Permissions") || strings.Contains(errStr, "50013") {
			return fmt.Errorf("bot doesn't have permission to send messages in the claim channel")
		}
		return fmt.Errorf("failed to post claim message: %w", err)
	}

	if claimMap["old_id"] != nil {
		primaryKey = claimMap["old_id"].(string)
	} else {
		primaryKey = message.ID
	}

	claimMap["channel_id"] = claimMap["log_channel"].(string)
	claimMap["type"] = eventType
	claimMap["message_id"] = message.ID

	// Use transaction if provided, otherwise use global DB
	db := database.DB
	if tx != nil {
		db = tx
	}

	result = db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.Claim{}).Create(map[string]interface{}{
		"message_id": primaryKey,
	})
	if result.Error != nil {
		return result.Error
	}

	result = db.Model(database.Claim{
		MessageID: primaryKey,
	}).Select([]string{"message_id", "channel_id", "guild_id", "item", "type", "winner", "cost", "host", "bid_history", "note", "image_url", "description"}).Updates(claimMap)
	if result.Error != nil {
		logger.Sugar.Errorw("critical error updating claim", "message_id", primaryKey, "error", result.Error)
		return fmt.Errorf("failed to save claim data. Please contact support")
	}

	return err
}

func ClaimTicket(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	return nil
}

func ClaimPrizeButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	idSlice := strings.Split(i.MessageComponentData().CustomID, ":")
	if len(idSlice) == 2 {
		if idSlice[1] != i.Member.User.ID {
			return fmt.Errorf("UserID does not match. You can only claim your own prizes")
		}
	} else {
		err := h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "Claim Prize",
			Description: "This button was created before the Claim Prize update, and has no ID saved so the bot cannot easily verify if you are the winner. Contact the auction host to claim.",
		})
		return err
	}

	claimMap := map[string]interface{}{}
	claimSetup := map[string]interface{}{}

	result := database.DB.Model(database.Claim{}).First(&claimMap, i.Message.ID)
	if result.Error != nil {
		return result.Error
	}
	database.DB.Model(database.ClaimSetup{}).First(&claimSetup, i.GuildID)

	if claimSetup["disable_claiming"] == true {
		return fmt.Errorf("Claiming has been disabled. A server administrator must use `/claim setup disable_claiming:` to re-enable.")
	}
	if claimSetup["category"] == nil {
		claimSetup["category"] = ""
	}
	if claimMap["item"] == nil {
		return fmt.Errorf("No item to claim saved in database. Contact support server for help.")
	}
	if claimMap["host"] == nil {
		claimMap["host"] = "0"
	}

	// Check if claim has already been completed or cancelled
	if claimMap["status"] != nil {
		status := claimMap["status"].(string)
		if status == database.ClaimStatusClaimed {
			return fmt.Errorf("This prize has already been claimed.")
		}
		if status == database.ClaimStatusCancelled {
			return fmt.Errorf("This claim was cancelled. Please contact the host if you believe this is an error.")
		}
	}

	// Check for existing open ticket
	if claimMap["ticket_id"] != nil && claimMap["ticket_id"].(string) != "" {
		currentChannel, err := s.Channel(claimMap["ticket_id"].(string))
		if err == nil {
			return fmt.Errorf("You already have a ticket open for this prize. Please go to <#%s> to claim.", currentChannel.ID)
		}
		// If channel doesn't exist but ticket_id is set, clear it (orphaned reference)
		database.DB.Model(&database.Claim{MessageID: claimMap["message_id"].(string)}).Update("ticket_id", "")
	}

	category, err := s.Channel(claimSetup["category"].(string))
	if err != nil {
		logger.Sugar.Warnw("error fetching category", "error", err)
		category = &discordgo.Channel{}
	}

	perms := []*discordgo.PermissionOverwrite{
		{
			ID:    i.Member.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
		},
		{
			ID:    s.State.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
		},
		{
			ID:    claimMap["host"].(string),
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
		},
		{
			ID:   i.GuildID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionViewChannel,
		},
	}

	for _, catPerm := range category.PermissionOverwrites {
		duplicatePerm := false
		for _, newPerm := range perms {
			if catPerm.ID == newPerm.ID {
				duplicatePerm = true
			}
		}
		if !duplicatePerm {
			perms = append(perms, catPerm)
		}
	}

	//Add permissions for the opener to see the channels. Plus add support role.
	channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name:                 "🎁│" + i.Member.User.Username + i.Member.User.Discriminator,
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             claimSetup["category"].(string),
		PermissionOverwrites: perms,
	})
	if err != nil {
		return fmt.Errorf("Error creating claim channel: %w", err)
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "__**ITEM**__",
			Value: claimMap["item"].(string),
		},
	}
	if claimMap["description"] != nil && claimMap["description"] != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "__**DESCRIPTION**__",
			Value: claimMap["description"].(string),
		})
	}

	if claimMap["cost"] != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "__**COST**__",
			Value: fmt.Sprint(claimMap["cost"].(float64)),
		})
	}

	if claimSetup["instructions"] != nil && claimSetup["instructions"] != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "__**INSTRUCTIONS**__",
			Value: claimSetup["instructions"].(string),
		})
	}
	if claimMap["note"] != nil && claimMap["note"] != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "__**NOTE**__",
			Value: claimMap["note"].(string),
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  "**To Close this Ticket:**",
		Value: "• **Complete**: Click when you have successfully claimed your prize.\n• **Cancel** (Admin only): Close the ticket without claiming.",
	})

	/*thumbnail := &discordgo.MessageEmbedThumbnail{}

	if len(i.Message.Embeds) > 0 && i.Message.Embeds[0].Thumbnail != nil {
		thumbnail = i.Message.Embeds[0].Thumbnail
	}*/

	_, err = h.SuccessMessage(s, channel.ID, h.PresetResponse{
		Content:     fmt.Sprintf("<@%s> %s", claimMap["host"].(string), i.Member.Mention()),
		Title:       "Ticket!",
		Description: "This is where you claim your prize.",
		Fields:      fields,
		//Thumbnail:   thumbnail,
		Image:       &discordgo.MessageEmbedImage{},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Complete",
						Style:    3,
						CustomID: fmt.Sprintf("claim_complete:%s:%s", i.ChannelID, i.Message.ID),
					},
					discordgo.Button{
						Label:    "Cancel",
						Style:    4,
						CustomID: fmt.Sprintf("claim_cancel:%s", i.Message.ID),
					},
				},
			},
		},
		Embeds: []*discordgo.MessageEmbed{},
		Files:  []*discordgo.File{},
	})
	if err != nil {
		logger.Sugar.Warnw("claim operation error", "error", err)
	}

	database.DB.Model(database.Claim{
		MessageID: claimMap["message_id"].(string),
	}).Update("ticket_id", channel.ID)

	logger.Sugar.Debugw("claim channel created", "channel_id", channel.ID)

	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Claim Prize",
		Description: "Please visit the ticket channel to claim your prize. <#" + channel.ID + ">",
	})
	if err != nil {
		logger.Sugar.Warnw("claim operation error", "error", err)
	}
	return nil
}

func CompleteButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	issues := ""
	//thumbnailImage := &discordgo.MessageEmbedThumbnail{}

	customID := strings.Split(i.MessageComponentData().CustomID, ":")

	if len(customID) < 3 {
		return fmt.Errorf("Failed parsing button ID. No messageID found. Please contact __AFTM Prize Manager support__ and I can look into the issue.")
	}

	claimMap := map[string]interface{}{}
	claimFound := true

	result := database.DB.Model(database.Claim{}).First(claimMap, customID[2])
	if result.Error != nil {
		logger.Sugar.Warnw("claim not found, will close ticket gracefully", "message_id", customID[2], "error", result.Error)
		claimFound = false
		// Set minimal data for orphaned ticket handling
		claimMap["winner"] = i.Member.User.ID
	}

	claimSetup := map[string]interface{}{}

	result = database.DB.Model(database.ClaimSetup{}).First(claimSetup, i.GuildID)
	if result.Error != nil {
		logger.Sugar.Errorw("failed to find claim setup", "guild_id", i.GuildID, "error", result.Error)
		return fmt.Errorf("failed to find claim setup. Please use `/settings claiming` at least once")
	}

	if claimMap["image_url"] != nil {
		//thumbnailImage.URL = claimMap["image_url"].(string)
	}

	if claimSetup["log_channel"] == nil {
		return fmt.Errorf("No Log Channel set for claim tickets. Please have an administrator run `/claim setup log_channel:` to set it.")
	}

	message, err := s.ChannelMessage(customID[1], customID[2])
	if err != nil {
		logger.Sugar.Warnw("original claim message not found", "channel_id", customID[1], "message_id", customID[2], "error", err)
		issues += "Original message not found (may have been deleted). "
	}

	if message != nil && message.Embeds != nil && len(message.Embeds) > 0 {
		//0x14f7b2 0x50c878
		message.Embeds[0].Color = 0x00c940

		components := []discordgo.MessageComponent{}
		embeds := message.Embeds
		_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content:         &message.Content,
			Components:      &components,
			Embeds:          &embeds,
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			ID:              customID[2],
			Channel:         customID[1],
		})
		if err != nil {
			return fmt.Errorf("There was an issue editing the old log embed: %w", err)
		}
	} else if message != nil {
		components := []discordgo.MessageComponent{}
		_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content:    &message.Content,
			Components: &components,
			ID:         customID[2],
			Channel:    customID[1],
		})
		if err != nil {
			logger.Sugar.Warnw("failed to remove original message buttons", "channel_id", customID[1], "message_id", customID[2], "error", err)
			issues += "Original message buttons could not be removed. "
		}
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "**Opened By:**",
			Value:  fmt.Sprintf("<@%s>", claimMap["winner"]),
			Inline: true,
		},
		{
			Name:   "**Closed By:**",
			Value:  fmt.Sprintf("<@%s>", i.Member.User.ID),
			Inline: true,
		},
	}

	if claimMap["item"] != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**Item Claimed:**",
			Value:  claimMap["item"].(string),
			Inline: true,
		})
	}

	if claimMap["bid"] != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**Price**",
			Value:  PriceFormat(claimMap["bid"].(float64), i.GuildID, claimMap["currency"]),
			Inline: false,
		})
	}

	if claimMap["host"] != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**Sold By:**",
			Value:  fmt.Sprintf("<@%s>", claimMap["host"]),
			Inline: true,
		})
	}

	if issues != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**Issues**",
			Value:  issues,
			Inline: true,
		})
	}

	_, err = h.SuccessMessage(s, claimSetup["log_channel"].(string), h.PresetResponse{
		Title:       "**Ticket Log**",
		Description: "This prize has been claimed successfully, soon I should have a transcript of the ticket saved as well.",
		Fields:      fields,
		//Thumbnail:   thumbnailImage,
		Components: []discordgo.MessageComponent{
			/*discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Re-Open",
						Style: discordgo.SuccessButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "🔄",
						},
						Disabled: true,
						CustomID: fmt.Sprintf("reopen_ticket:%s", claimMap["winner"]),
					},
				},
			},*/
		},
	})

	if err != nil {
		return err
	}

	// Respond to user before deleting channel
	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Prize Claimed Successfully",
		Description: "The ticket will now be closed.",
	})
	if err != nil {
		logger.Sugar.Warnw("claim operation error", "error", err)
	}

	// Delete channel FIRST (before updating DB status)
	_, err = s.ChannelDelete(i.ChannelID)
	if err != nil {
		logger.Sugar.Errorw("ticket could not be closed", "channel_id", i.ChannelID, "error", err)
		// Don't update claim status if channel deletion failed
		return fmt.Errorf("Ticket could not be closed: %w. Claim status was NOT updated.", err)
	}

	// Only update status AFTER successful channel deletion
	if claimFound {
		result = database.DB.Model(&database.Claim{MessageID: customID[2]}).Updates(map[string]interface{}{
			"status":    database.ClaimStatusClaimed,
			"ticket_id": "", // Clear ticket_id since channel is deleted
		})
		if result.Error != nil {
			// Log error but don't fail - channel is already deleted, claim is functionally complete
			logger.Sugar.Errorw("failed to update claim status after channel deletion",
				"message_id", customID[2], "error", result.Error)
		}
	}

	return nil
}

func CancelButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	customID := strings.Split(i.MessageComponentData().CustomID, ":")

	if len(customID) < 2 {
		return fmt.Errorf("Failed parsing button ID. No messageID found. Please contact __AFTM Prize Manager support__ and I can look into the issue.")
	}

	claimMap := map[string]interface{}{}
	claimFound := true

	result := database.DB.Model(database.Claim{}).First(claimMap, customID[1])
	if result.Error != nil {
		logger.Sugar.Warnw("claim not found, will attempt graceful close", "message_id", customID[1], "error", result.Error)
		claimFound = false
		// Set minimal data for orphaned ticket handling
		claimMap["winner"] = i.Member.User.ID
	}

	// Permission check: admins can always close, hosts can close their own tickets
	// For orphaned tickets (claimFound=false), only admins can close
	isAdmin := i.Member.Permissions&(1<<3) == 8
	isHost := claimFound && claimMap["host"] != nil && i.Member.User.ID == claimMap["host"]
	if !isAdmin && !isHost {
		if claimFound {
			return fmt.Errorf("User must have administrator permissions or be the host (<@%s>) to run this command", fmt.Sprint(claimMap["host"]))
		}
		return fmt.Errorf("User must have administrator permissions to close orphaned tickets")
	}

	claimSetup := map[string]interface{}{}

	result = database.DB.Model(database.ClaimSetup{}).First(claimSetup, i.GuildID)
	if result.Error != nil {
		logger.Sugar.Errorw("failed to fetch claim setup", "guild_id", i.GuildID, "error", result.Error)
		return fmt.Errorf("failed to fetch claim settings. Please have an admin run `/settings claiming` at least once")
	}

	if claimSetup["log_channel"] == nil {
		return fmt.Errorf("No Log Channel set for claim tickets. Please have an administrator run `/settings claiming log_channel:` to set it.")
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "**Opened By:**",
			Value:  fmt.Sprintf("<@%s>", claimMap["winner"]),
			Inline: true,
		},
		{
			Name:   "**Closed By:**",
			Value:  fmt.Sprintf("<@%s>", i.Member.User.ID),
			Inline: true,
		},
	}

	if claimMap["item"] != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**Item:**",
			Value:  claimMap["item"].(string),
			Inline: true,
		})
	}

	if claimMap["host"] != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "**Host:**",
			Value:  fmt.Sprintf("<@%s>", claimMap["host"]),
			Inline: true,
		})
	}

	_, err := h.SuccessMessage(s, claimSetup["log_channel"].(string), h.PresetResponse{
		Title:       "**Ticket Log**",
		Description: "Prize claiming was cancelled!",
		Fields:      fields,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Re-Open",
						Style: discordgo.SuccessButton,
						Emoji: &discordgo.ComponentEmoji{
							Name: "🔄",
						},
						Disabled: true,
						CustomID: fmt.Sprintf("reopen_ticket:%s", claimMap["winner"]),
					},
				},
			},
		},
	})
	if err != nil {
		logger.Sugar.Warnw("claim operation error", "error", err)
	}

	// Delete channel FIRST
	_, err = s.ChannelDelete(i.ChannelID)
	if err != nil {
		logger.Sugar.Errorw("ticket could not be closed", "channel_id", i.ChannelID, "error", err)
		return fmt.Errorf("Ticket could not be closed: %w. Claim status was NOT updated.", err)
	}

	// Only update status AFTER successful channel deletion
	if claimFound {
		result = database.DB.Model(&database.Claim{MessageID: customID[1]}).Updates(map[string]interface{}{
			"status":    database.ClaimStatusCancelled,
			"ticket_id": "", // Clear ticket_id since channel is deleted
		})
		if result.Error != nil {
			logger.Sugar.Errorw("failed to update claim status after cancel",
				"message_id", customID[1], "error", result.Error)
		}
	}

	return nil
}

func claimRefresh(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := h.ParseSubCommand(i)

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Claim Refresh",
		Description: "Starting to refresh your missing claims.",
	})

	claimMap := []map[string]interface{}{}
	// Only refresh pending claims (not claimed or cancelled)
	// Handle legacy data: status IS NULL OR status = '' OR status = 'pending'
	result := database.DB.Model([]database.Claim{}).Where(map[string]interface{}{
		"guild_id": i.GuildID,
	}).Where("status = ? OR status IS NULL OR status = ''", database.ClaimStatusPending).Find(&claimMap)
	if result.Error != nil {
		logger.Sugar.Errorw("failed to fetch claims for refresh", "guild_id", i.GuildID, "error", result.Error)
		h.FollowUpErrorResponse(s, i, "Failed to fetch claims. Please try again or contact support.")
		return nil
	}

	var restored int

	for _, v := range claimMap {
		v["log_channel"] = options["channel"]
		if v["type"] == nil {
			v["type"] = "unknown"
		}

		if v["channel_id"] != nil {
			_, err := s.ChannelMessage(v["channel_id"].(string), v["message_id"].(string))
			v["old_id"] = v["message_id"]
			if err != nil {
				logger.Sugar.Debugw("processing claim message", "message_id", v["message_id"])
				err = ClaimOutput(s, v, v["type"].(string))
				restored++
				if err != nil {
					logger.Sugar.Warnw("claim restore failed", "message_id", v["message_id"], "error", err)
					h.FollowUpErrorResponse(s, i, "Failed to restore a claim. Please try again or contact support.")
				}
			}
		}
	}

	h.FollowUpSuccessResponse(s, i, h.PresetResponse{
		Title:       "Claim Refresh Complete",
		Description: fmt.Sprintf("%d claim(s) have been restored to <#%s>", restored, options["channel"]),
	})
	return nil
}

func ClaimInventory(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	options := h.ParseSubCommand(i)
	claimSlice := []map[string]interface{}{}
	fields := []*discordgo.MessageEmbedField{}

	options["guild_id"] = i.GuildID

	// Only show pending claims (not claimed or cancelled)
	// Handle legacy data: status IS NULL OR status = '' OR status = 'pending'
	result := database.DB.Model(database.Claim{}).Where(map[string]interface{}{
		"winner":   options["user"],
		"guild_id": options["guild_id"],
	}).Where("status = ? OR status IS NULL OR status = ''", database.ClaimStatusPending).Limit(25).Find(&claimSlice)
	if result.Error != nil {
		return result.Error
	}

	for _, v := range claimSlice {
		if v["description"] == nil {
			v["description"] = "No description."
		}
		if v["item"] != nil {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   v["item"].(string),
				Value:  v["description"].(string),
				Inline: false,
			})
		}
	}
	userAvatar := ""
	user, err := s.User(options["user"].(string))
	if err != nil {
		logger.Sugar.Warnw("claim operation error", "error", err)
	} else {
		userAvatar = user.AvatarURL("")
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:  "**Inventory**",
		Fields: fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: userAvatar,
		},
		Image: &discordgo.MessageEmbedImage{},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "▶",
						Style:    0,
						Disabled: false,
						Emoji: &discordgo.ComponentEmoji{
							Name:     "",
							ID:       "",
							Animated: false,
						},
						URL:      "",
						CustomID: "inventory_right:",
					},
				},
			},
		},
		Embeds: []*discordgo.MessageEmbed{},
		Files:  []*discordgo.File{},
	})
	return nil
}

func TicketEnd() {

}

func ReopenTicket(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	return nil
}
