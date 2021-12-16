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
				/*{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "log_channel",
					Description: "Sets the channel where auctions will send outputs when they end",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						0,
					},
				},*/
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
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "image",
					Description: "Must be a link",
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

	database.DB.Model(&info).Updates(options)

	//Now check what options are set
	setOptions := map[string]interface{}{}

	database.DB.Model(database.GiveawaySetup{}).First(&setOptions, i.GuildID)

	responseFields := []*discordgo.MessageEmbedField{}

	for n, v := range GiveawayCommand.Options {
		if v.Name == "setup" {
			for _, v := range GiveawayCommand.Options[n].Options {
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

			for _, v := range GiveawayCommand.Options[n].Options {
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
								MaxValues:   len(GiveawayCommand.Options[n].Options),
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
	}
}

func GiveawayCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	GiveawaySetup := database.GiveawaySetup{
		GuildID: i.GuildID,
	}
	database.DB.First(&GiveawaySetup, i.GuildID)

	options := h.ParseSubCommand(i)

	item := options["item"].(string)
	winners := options["winners"].(float64)
	channel := i.ChannelID
	duration := options["duration"].(string)
	canHost := false

	if GiveawaySetup.HostRole != "" {
		for _, v := range i.Member.Roles {
			if v == GiveawaySetup.HostRole {
				canHost = true
			}
		}
		if i.Member.Permissions&(1<<3) == 8 {
			canHost = true
		}
		if !canHost {
			h.ErrorResponse(s, i, "User must be administrator or have the role <@&"+GiveawaySetup.HostRole+"> to host giveaways.")
			return
		}
	}

	if winners <= 0 {
		h.ErrorResponse(s, i, "Must have 1 or more winners. Winners entered: "+fmt.Sprint(winners))
		fmt.Println("Must have 1 or more winners. Winners entered: ", winners)
		return
	}

	if options["channel"] != nil {
		channel = options["channel"].(string)
	}

	var imageURL string
	if options["image"] != nil {
		imageURL = options["image"].(string)
	}

	description := fmt.Sprintf("**%s** Winners!\n", fmt.Sprint(winners))
	if options["description"] != nil {
		description += options["description"].(string)
	}
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		fmt.Println(err)
	}

	endTimeDuration, err := h.ParseTime(duration)
	if err != nil {
		fmt.Println(err)
		h.ErrorResponse(s, i, err.Error())
		return
	}

	endTime := time.Now().Add(endTimeDuration)

	content := ""
	if GiveawaySetup.AlertRole != "" {
		content = fmt.Sprintf("<@&%s>", GiveawaySetup.AlertRole)
	}

	message, err := h.PresetMessageSend(s, channel, h.PresetResponse{
		Content:     content,
		Title:       "__**" + item + "**__",
		Description: fmt.Sprintf("<@%s> has hosted a giveaway! To enter, press the button below!", i.Member.User.ID),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "__**Details:**__",
				Value:  description,
				Inline: false,
			},
			{
				Name:   "**End Time**",
				Value:  fmt.Sprintf("<t:%d:R>", endTime.Unix()),
				Inline: false,
			},
		},
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
						Label: "Enter!",
						Style: discordgo.PrimaryButton,
						Emoji: discordgo.ComponentEmoji{
							Name:     "meeting",
							ID:       "759178932892729344",
							Animated: true,
						},
						CustomID: "enter_giveaway",
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	database.DB.Create(&database.Giveaway{
		MessageID:   message.ID,
		ChannelID:   channel,
		GuildID:     i.GuildID,
		Item:        item,
		EndTime:     endTime,
		Description: description,
		Host:        i.Member.User.ID,
		Winners:     winners,
	})

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Giveaway Started!**",
		Description: fmt.Sprintf("Giveaway has started in <#%s>", channel),
	})

	time.Sleep(endTimeDuration)
	GiveawayEnd(s, message.ID)
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

func GiveawayEnter(s *discordgo.Session, i *discordgo.InteractionCreate) {

	giveawayInfo := database.Giveaway{
		MessageID: i.Message.ID,
	}
	database.DB.First(&giveawayInfo, i.Message.ID)

	for _, v := range strings.Split(giveawayInfo.Entries, " ") {
		if v == i.Member.User.ID {
			h.ErrorResponse(s, i, "You have already entered this giveaway!")
			fmt.Println("User has already entered giveaway.")
			return
		}
	}
	if giveawayInfo.Entries == "" {
		giveawayInfo.Entries = i.Member.User.ID
	} else {
		giveawayInfo.Entries += fmt.Sprintf(" %s", i.Member.User.ID)
	}

	database.DB.Model(&giveawayInfo).Updates(giveawayInfo)

	fmt.Println(strings.Split(giveawayInfo.Entries, " "))

	if len(i.Message.Embeds[0].Fields) == 3 {
		i.Message.Embeds[0].Fields[2].Value = fmt.Sprint(len(strings.Split(giveawayInfo.Entries, " ")))
	} else {
		i.Message.Embeds[0].Fields = append(i.Message.Embeds[0].Fields, &discordgo.MessageEmbedField{
			Name:   "**Number of Entries**",
			Value:  fmt.Sprint(len(strings.Split(giveawayInfo.Entries, " "))),
			Inline: false,
		})
	}

	messageEdit := discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)

	messageEdit.Embeds = i.Message.Embeds
	messageEdit.Components = i.Message.Components

	_, err := s.ChannelMessageEditComplex(messageEdit)

	if err != nil {
		fmt.Println(err.Error())
		h.ErrorResponse(s, i, err.Error())
		return
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "**Successful Entry!**",
		Description: "You have successfully been entered into the giveaway.",
	})
}

func GiveawayEnd(s *discordgo.Session, messageID string) {

	winnerTags, winnerList, giveawayInfo, err := GiveawayRoll(s, messageID)

	if err != nil {
		fmt.Println(err)
		return
	}

	m, err := s.ChannelMessage(giveawayInfo.ChannelID, messageID)
	if err != nil {
		fmt.Println(err)
		return
	}

	updateWinner := false

	for num, v := range m.Embeds[0].Fields {
		if v.Name == "**Giveaway Has Ended!**" {
			updateWinner = true
			m.Embeds[0].Fields[num] = &discordgo.MessageEmbedField{
				Name:   "**Giveaway Has Ended!**",
				Value:  "**Winners:**\n" + winnerList,
				Inline: false,
			}
		}
	}

	if !updateWinner {
		m.Embeds[0].Fields = append(m.Embeds[0].Fields, &discordgo.MessageEmbedField{
			Name:   "**Giveaway Has Ended!**",
			Value:  "**Winners:**\n" + winnerList,
			Inline: false,
		})
	}

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Content: &m.Content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Reroll",
						Style:    1,
						CustomID: "reroll_giveaway",
					},
				},
			},
		},
		Embeds:  m.Embeds,
		ID:      messageID,
		Channel: m.ChannelID,
	})
	if err != nil {
		fmt.Println(err)
	}

	message, err := h.PresetMessageSend(s, giveawayInfo.ChannelID, h.PresetResponse{
		Content: winnerTags,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Claim!",
						Style: 3,
						Emoji: discordgo.ComponentEmoji{
							Name: "cryopod",
							ID:   "889307390690885692",
						},
						CustomID: "claim_giveaway",
					},
				},
			},
		},
		Title:       "Giveaway Completed!",
		Description: "You have 24 hours to reroll the winners if you would like.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Giveaway Host**",
				Value:  fmt.Sprintf("This Giveaway was hosted by: <@!%s>", giveawayInfo.Host),
				Inline: true,
			},
			{
				Name:   "**Item Won**",
				Value:  giveawayInfo.Item,
				Inline: false,
			},
			{
				Name:   "**Winners:**",
				Value:  winnerList,
				Inline: false,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: giveawayInfo.ImageURL,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	giveawayInfo.Finished = true
	if giveawayInfo.WinnerOutput != "" {
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			fmt.Println(err)
		}
	}
	giveawayInfo.WinnerOutput = message.ID
	database.DB.Model(&giveawayInfo).Updates(giveawayInfo)

	time.Sleep(24 * time.Hour)
	database.DB.Delete(database.Giveaway{}, messageID)
}

func GiveawayRoll(s *discordgo.Session, messageID string) (string, string, database.Giveaway, error) {

	var winnerList string
	var winnerTags string
	giveawayInfo := database.Giveaway{
		MessageID: messageID,
	}
	database.DB.First(&giveawayInfo, messageID)

	for n := float64(0); n < giveawayInfo.Winners; {

		entryString := strings.Split(giveawayInfo.Entries, " ")

		if len(entryString) == 0 || giveawayInfo.Entries == "" {
			winnerList += "• Not enough entrants!"
			fmt.Println("Not enough Entrants")
			break
		}

		result, err := random.Int(random.Reader, big.NewInt(int64(len(entryString))))
		if err != nil {
			fmt.Println("Random Error: ", err)
			return "", "", database.Giveaway{}, err
		}

		winnerID := entryString[result.Int64()]

		user, err := s.User(winnerID)
		if err != nil {
			fmt.Println("User Error:", err)
			return "", "", database.Giveaway{}, err
		}

		winner := fmt.Sprintf("<@%s> (%s#%s)", user.ID, user.Username, user.Discriminator)

		winnerList += fmt.Sprintf("• %s\n", winner)
		winnerTags += fmt.Sprintf("<@%s>, ", winnerID)

		giveawayInfo.Entries = strings.Trim(strings.ReplaceAll(" "+giveawayInfo.Entries, " "+winnerID, ""), " ")

		n++
	}

	return winnerTags, winnerList, giveawayInfo, nil
}

func ClaimGiveawayButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var GiveawaySetup database.GiveawaySetup
	database.DB.First(&GiveawaySetup, i.GuildID)

	if GiveawaySetup.Claiming == "" {
		GiveawaySetup.Claiming = "The discord owner has not set a claiming message. Common ways to claim include: Opening a ticket or contacting the auction host. \nTo customize this message, use the command: `/giveaway setup claiming:`."
	}

	err := h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Claim Prize",
		Description: GiveawaySetup.Claiming,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func RerollGiveawayButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	giveawayInfo := database.Giveaway{
		MessageID: i.Message.ChannelID,
	}

	database.DB.First(&giveawayInfo)

	if i.Member.Permissions&(1<<3) != 8 && i.Member.User.ID != giveawayInfo.Host {
		h.ErrorResponse(s, i, "User must be host or have administrator permissions to run this command")
		return
	}

	GiveawayEnd(s, i.Message.ID)
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
