package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ClaimCommand = discordgo.ApplicationCommand{
	Name:        "claim",
	Description: "Put an item up for auction!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "setup",
			Description: "Setup the claiming system for your prizes.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "category",
					Description: "The category to claim prizes in.",
					ChannelTypes: []discordgo.ChannelType{
						4,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "log_channel",
					Description: "The output channel for completed tickets.",
					ChannelTypes: []discordgo.ChannelType{
						0,
						5,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "instructions",
					Description: "Leave instructions for whoever opens the ticket.",
				},
			},
			Autocomplete: false,
			Choices:      []*discordgo.ApplicationCommandOptionChoice{},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Name:        "create",
			Description: "Create a claimable prize.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "user",
					Description: "Create a claimable prize for someone.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "winner",
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
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "description",
							Description: "Set a custom item description",
							Required:    false,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "role",
					Description: "Create a claimable prize for a role.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionRole,
							Name:        "role",
							Description: "The role of users who will receive the prize.",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "item",
							Description: "The prize you are giving",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "description",
							Description: "Set a custom item description",
							Required:    false,
						},
					},
				},
			},
		},
	},
}

func Claim(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[0].Name {
	case "setup":
		ClaimSetup(s, i)
	case "create":
		ClaimCreate(s, i)
	}
}

func ClaimSetup(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&(1<<3) != 8 {
		h.ErrorResponse(s, i, "User must have administrator permissions to run this command")
		return
	}

	options := h.ParseSubCommand(i)

	claimSetup := map[string]interface{}{
		"guild_id": i.GuildID,
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.ClaimSetup{}).Create(&claimSetup)
	if result.Error != nil {
		fmt.Println(result.Error)
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	result = database.DB.Model(database.ClaimSetup{
		GuildID: i.GuildID,
	}).Updates(options)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	//Now check what options are set
	setOptions := map[string]interface{}{}

	database.DB.Model(database.ClaimSetup{}).First(&setOptions, i.GuildID)

	responseFields := []*discordgo.MessageEmbedField{}

	for _, v := range ClaimCommand.Options[0].Options {

		//Need to see if it's still set after it says not set because of empty string
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

	menuOptions := []discordgo.SelectMenuOption{}

	for _, v := range ClaimCommand.Options[0].Options {
		menuOptions = append(menuOptions, discordgo.SelectMenuOption{
			Label:       strings.Title(strings.ReplaceAll(v.Name, "_", " ")),
			Value:       v.Name,
			Description: v.Description,
		})
	}

	err := h.SuccessResponse(s, i, h.PresetResponse{
		Title:  "Claim Setup",
		Fields: responseFields,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "clear_claim_setup",
						Placeholder: "Clear Setup Options",
						MinValues:   1,
						MaxValues:   len(ClaimCommand.Options[0].Options),
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

func ClaimSetupClearButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.MessageComponentData().Values
	clearedMap := map[string]interface{}{}

	info := database.ClaimSetup{
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

func ClaimCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&(1<<3) != 8 {
		h.ErrorResponse(s, i, "User must have administrator permissions to run this command")
		return
	}

	claimSetup := map[string]any{}
	auctionSetup := map[string]any{}

	claimMap := h.ParseSubSubCommand(i)

	result := database.DB.Model(database.AuctionSetup{}).First(&auctionSetup, i.GuildID)
	if result.Error != nil {
		h.ErrorResponse(s, i, fmt.Sprintf("Error fetching setup, try running `/auction setup` to fix. Error: %s", result.Error.Error()))
		fmt.Println(result.Error)
		return
	}
	result = database.DB.Model(database.ClaimSetup{}).First(&claimSetup, i.GuildID)
	if result.Error != nil {
		h.ErrorResponse(s, i, fmt.Sprintf("Error fetching setup, try running `/claim setup` to fix. Error: %s", result.Error.Error()))
		fmt.Println(result.Error)
		return
	}

	claimMap["log_channel"] = auctionSetup["log_channel"]
	claimMap["host"] = i.Member.User.ID

	switch i.ApplicationCommandData().Options[0].Options[0].Name {
	case "role":
		err := h.ExperimentalResponse(s, i, h.PresetResponse{
			Title:       "Claims are being created!",
			Description: "Check out <#" + claimMap["log_channel"].(string) + "> to see the claims. The bot will respond here when complete, or if there is an error.",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "**Free Premium Feature!**",
					Value: "Eventually this is gonna be a premium feature, since it's just for convenience and it took more effort to create. For now though, its free for all users!"},
			},
		})
		if err != nil {
			fmt.Println(err)
		}
		afterID := ""
		for {

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})

			memberSlice, err := s.GuildMembers(i.GuildID, afterID, 1000)

			fmt.Println("Printing members: ")

			for _, v := range memberSlice {
				fmt.Println(v.User.Username)
				for _, role := range v.Roles {
					if role == claimMap["role"] {
						claimMap["winner"] = v.User.ID
						err = ClaimOutput(s, claimMap, "Custom Claim")
						if err != nil {
							h.FollowUpErrorResponse(s, i, err.Error())
							fmt.Println(err)
							return
						}
					}
				}
			}

			if len(memberSlice) < 1000 {
				break
			}
			afterID = memberSlice[len(memberSlice)-1].User.ID
		}
		_, err = h.FollowUpSuccessResponse(s, i, h.PresetResponse{
			Title:       "Claims Successfully Created!",
			Description: "All claims should now be created in: <#" + claimMap["channel_id"].(string) + ">",
		})
		if err != nil {
			fmt.Println(err)
		}

	case "user":
		err := ClaimOutput(s, claimMap, "Custom Claim")
		if err != nil {
			h.ErrorResponse(s, i, err.Error())
			fmt.Println(err)
			return
		}
		err = h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "Claim Successfully Created!",
			Description: "Check out <#" + claimMap["channel_id"].(string) + "> to see the claim",
		})
		if err != nil {
			fmt.Println(err)
		}
	}

}

//Seems like using a map here overcomplicates it. Possibly need to go back to fix if I run into issues.
func ClaimOutput(s *discordgo.Session, claimMap map[string]interface{}, claimType string) error {

	mentionUser := ""
	finalBid := ""
	claimMap["cost"] = claimMap["bid"]

	/*result := database.DB.Model(database.ClaimSetup{}).First(&claimSetup, claimMap["guild_id"])
	if result.Error != nil {
		return fmt.Errorf("Error Fetching Claim Setup database. Try running `/claim setup` to make sure everything is set properly. Error message: " + result.Error.Error())
	}*/

	if claimMap["formatted_price"] != nil {
		finalBid = claimMap["formatted_price"].(string)
	}
	if claimMap["buyout_message"] != nil {
		finalBid = claimMap["buyout_message"].(string)
	}
	if claimMap["target_message"] != nil {
		finalBid = claimMap["target_message"].(string)
	}
	if claimMap["winner"] == nil {
		claimMap["winner"] = "No Winner Detected"
	}

	buttons := []discordgo.MessageComponent{
		discordgo.Button{
			Label: "Support",
			Style: discordgo.LinkButton,
			Emoji: discordgo.ComponentEmoji{
				Name:     "logo",
				ID:       "889025400120950804",
				Animated: false,
			},
			URL: "https://discord.gg/RxP2z5NGtj",
		},
	}

	if claimMap["bid_history"] != nil {
		buttons = append(buttons, discordgo.Button{
			Label: "Bids",
			Style: 2,
			Emoji: discordgo.ComponentEmoji{
				Name: "📜",
			},
			CustomID: "bid_history",
		})
	}

	user, err := Session.User(fmt.Sprint(claimMap["winner"]))
	if err != nil {
		user = &discordgo.User{}
		user.Username = claimMap["winner"].(string)
	} else {
		mentionUser = user.Mention()
		user.Username = fmt.Sprintf("%s (%s#%s)", user.Mention(), user.Username, user.Discriminator)
		buttons = append(buttons, discordgo.Button{
			Label: "Claim!",
			Style: 3,
			Emoji: discordgo.ComponentEmoji{
				Name: "cryopod",
				ID:   "889307390690885692",
			},
			CustomID: "claim_prize:" + claimMap["winner"].(string),
		})
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

	message, err := h.SuccessMessage(Session, claimMap["log_channel"].(string), h.PresetResponse{
		Content: mentionUser,
		//Add in the type of prize
		Title:       fmt.Sprintf("%s Prize: __**%s**__", claimType, claimMap["item"]),
		Description: "**Please Note:** Any old claims without this message may not work correctly.",
		Fields:      fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: claimMap["image_url"].(string),
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/9wo7diC.png",
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: buttons,
			},
		},
	})

	if err != nil {
		return err
	}

	claimMap["message_id"] = message.ID
	claimMap["channel_id"] = message.ChannelID

	result := database.DB.Model(database.Claim{}).Select([]string{"message_id", "channel_id", "item", "type", "winner", "cost", "host", "bid_history", "note", "image_url"}).Create(claimMap)
	if result.Error != nil {
		return result.Error
	}

	return err
}

func ClaimTicket(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func ClaimPrizeButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	idSlice := strings.Split(i.MessageComponentData().CustomID, ":")
	if len(idSlice) == 2 {
		if idSlice[1] != i.Member.User.ID {
			h.ErrorResponse(s, i, "UserID does not match. You can only claim your own prizes")
			return
		}
	} else {
		err := h.SuccessResponse(s, i, h.PresetResponse{
			Title:       "Claim Prize",
			Description: "This button was created before the Claim Prize update, and has no ID saved so the bot cannot easily verify if you are the winner. Contact the auction host to claim.",
		})
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	claimMap := map[string]interface{}{}
	claimSetup := map[string]interface{}{}

	result := database.DB.Model(database.Claim{}).First(&claimMap, i.Message.ID)
	if result.Error == nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}
	database.DB.Model(database.ClaimSetup{}).First(&claimSetup, i.GuildID)

	if claimSetup["category"] == nil {
		h.ErrorResponse(s, i, "The server admins have not set a claiming channel. Please use the command `/claim setup category:` to set it.")
		return
	}
	if claimMap["item"] == nil {
		h.ErrorResponse(s, i, "No item to claim saved in database. Contact support server for help.")
		return
	}
	if claimMap["host"] == nil {
		claimMap["host"] = "0"
	}
	if claimMap["ticket_id"] != nil {
		currentChannel, err := s.Channel(claimMap["ticket_id"].(string))
		if err == nil {
			h.ErrorResponse(s, i, "You already have a ticket open for this prize. Please go to <#"+currentChannel.ID+"> to claim.")
			return
		}
	}

	//Add permissions for the opener to see the channels. Plus add support role.
	channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name: "🎁│" + i.Member.User.Username + i.Member.User.Discriminator,
		Type: 0,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    i.Member.User.ID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Deny:  0,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
			{
				ID:    claimMap["host"].(string),
				Type:  discordgo.PermissionOverwriteTypeMember,
				Deny:  0,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages,
			},
			{
				ID:    i.GuildID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Deny:  discordgo.PermissionViewChannel,
				Allow: 0,
			},
		},
		ParentID: claimSetup["category"].(string),
	})
	if err != nil {
		fmt.Println(err)
		h.ErrorResponse(s, i, err.Error())
		return
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "__**ITEM**__",
			Value: claimMap["item"].(string),
		},
	}
	if claimSetup["description"] != nil && claimSetup["description"] != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "__**DESCRIPTION**__",
			Value: claimSetup["description"].(string),
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

	thumbnail := &discordgo.MessageEmbedThumbnail{}

	if i.Message.Embeds[0].Thumbnail != nil {
		thumbnail = i.Message.Embeds[0].Thumbnail
	}

	h.SuccessMessage(s, channel.ID, h.PresetResponse{
		Content:     i.Member.Mention(),
		Title:       "Ticket!",
		Description: "This is where you claim your prize.",
		Fields:      fields,
		Thumbnail:   thumbnail,
		Image:       &discordgo.MessageEmbedImage{},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Complete",
						Style:    3,
						Emoji:    discordgo.ComponentEmoji{},
						CustomID: fmt.Sprintf("claim_complete:%s:%s", i.ChannelID, i.Message.ID),
					},
					discordgo.Button{
						Label:    "Cancel",
						Style:    4,
						Emoji:    discordgo.ComponentEmoji{},
						CustomID: fmt.Sprintf("claim_cancel:%s", i.Message.ID),
					},
				},
			},
		},
		Embeds: []*discordgo.MessageEmbed{},
		Files:  []*discordgo.File{},
	})

	database.DB.Model(database.Claim{
		MessageID: claimMap["message_id"].(string),
	}).Update("ticket_id", channel.ID)

	fmt.Println("channel ID 2", channel.ID)

	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Claim Prize",
		Description: "Please visit the ticket channel to claim your prize. <#" + channel.ID + ">",
	})
	if err != nil {
		fmt.Println(err)
	}
}

func CompleteButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	issues := ""
	thumbnail := ""

	customID := strings.Split(i.MessageComponentData().CustomID, ":")

	if len(customID) < 3 {
		h.ErrorResponse(s, i, "Failed parsing button ID. No messageID found. Please contact __AFTM Prize Manager support__ and I can look into the issue.")
		return
	}

	claimMap := map[string]interface{}{}

	result := database.DB.Model(database.Claim{}).First(claimMap, customID[2])
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	claimSetup := map[string]interface{}{}

	result = database.DB.Model(database.ClaimSetup{}).First(claimSetup, i.GuildID)
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	if claimMap["image_url"] != nil {
		thumbnail = claimMap["image_url"].(string)
	}

	if claimSetup["log_channel"] == nil {
		h.ErrorResponse(s, i, "No Log Channel set for claim tickets. Please have an administrator run `/claim setup log_channel:` to set it.")
		return
	}

	message, err := s.ChannelMessage(customID[1], customID[2])
	if err != nil {
		issues += "Original message not found: " + err.Error()
	}

	if message.Embeds != nil {
		//0x14f7b2 0x50c878
		message.Embeds[0].Color = 0x00c940

		_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content:         &message.Content,
			Components:      []discordgo.MessageComponent{},
			Embeds:          message.Embeds,
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			ID:              customID[2],
			Channel:         customID[1],
		})
		if err != nil {
			h.ErrorResponse(s, i, "There was an issue editing the old log embed: "+err.Error())
			fmt.Println("There was an issue editing the old log embed: " + err.Error())
			return
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
			Value:  PriceFormat(map[string]interface{}{}, claimMap["bid"].(float64)),
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
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: thumbnail,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/9wo7diC.png",
		},
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
		h.ErrorResponse(s, i, err.Error())
		return
	}

	result = database.DB.Delete(database.Claim{}, customID[2])
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	err = h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Prize Claimed Successfully",
		Description: "The ticket will now be closed, please reopen if you have any issues.",
	})

	_, err = s.ChannelDelete(i.ChannelID)
	if err != nil {
		h.ErrorMessage(s, i.ChannelID, "Ticket could not be closed: "+err.Error())
		fmt.Println("Ticket could not be closed: " + err.Error())
	}

}

func CancelButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	customID := strings.Split(i.MessageComponentData().CustomID, ":")

	if len(customID) < 2 {
		h.ErrorResponse(s, i, "Failed parsing button ID. No messageID found. Please contact __AFTM Prize Manager support__ and I can look into the issue.")
		return
	}

	claimMap := map[string]interface{}{}

	result := database.DB.Model(database.Claim{}).First(claimMap, customID[1])
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != claimMap["host"] {
		h.ErrorResponse(s, i, fmt.Sprintf("User must have administrator permissions or be the host (%s) to run this command", fmt.Sprint(claimMap["host"])))
		return
	}

	claimSetup := map[string]interface{}{}

	result = database.DB.Model(database.ClaimSetup{}).First(claimSetup, i.GuildID)
	if result.Error != nil {
		h.ErrorResponse(s, i, result.Error.Error())
		return
	}

	if claimSetup["log_channel"] == nil {
		h.ErrorResponse(s, i, "No Log Channel set for claim tickets. Please have an administrator run `/claim setup log_channel:` to set it.")
		return
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
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/9wo7diC.png",
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
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
			},
		},
	})

	_, err = s.ChannelDelete(i.ChannelID)
	if err != nil {
		h.ErrorResponse(s, i, err.Error())
		return
	}

}

func TicketEnd() {

}

func ReopenTicket(s *discordgo.Session, i *discordgo.InteractionCreate) {

}
