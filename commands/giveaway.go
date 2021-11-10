package commands

import (
	random "crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
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
					Description: "Set a role to get pinged whenever an auction starts. Choosing @everyone will reset it to default.",
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
					Description: "Set a role that can host auctions. Choosing @everyone will reset it to default.",
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

	options := ParseSubCommand(i)
	content := ""

	if i.Member.Permissions&(1<<3) != 8 {
		ErrorResponse(s, i, "User must have administrator permissions to run this command")
		return
	}

	if options["alert_role"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		info.GiveawayRole = options["alert_role"].(string)

		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"giveaway_role": info.GiveawayRole}),
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
		info.GiveawayHostRole = options["host_role"].(string)

		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"giveaway_host_role": info.GiveawayHostRole}),
		}).Create(&info)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		content = content + "• Host Role has been successfully set.\n"
	}

	if options["claiming"] != nil {
		info := database.GuildInfo{
			GuildID: i.GuildID,
		}
		info.GiveawayClaiming = options["claiming"].(string)
		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "guild_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"giveaway_claiming": info.GiveawayClaiming}),
		}).Create(&info)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
		}
		content = content + "• Claiming Message has been successfully set.\n"
	}

	guildInfo := database.GuildInfo{
		GuildID: i.GuildID,
	}

	database.DB.First(&guildInfo, i.GuildID)

	if guildInfo.GiveawayHostRole == "" {
		guildInfo.GiveawayHostRole = "Not Set"
	} else {
		guildInfo.GiveawayHostRole = fmt.Sprintf("<@&%s>", guildInfo.GiveawayHostRole)
	}
	if guildInfo.GiveawayClaiming == "" {
		guildInfo.GiveawayClaiming = "Not Set"
	}
	if guildInfo.GiveawayRole == "" {
		guildInfo.GiveawayRole = "Not Set"
	} else {
		guildInfo.GiveawayRole = fmt.Sprintf("<@&%s>", guildInfo.GiveawayRole)
	}

	err := SuccessResponse(s, i, PresetResponse{
		Title:       "Giveaway Setup",
		Description: content,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Host Role**",
				Value: guildInfo.GiveawayHostRole,
			},
			{
				Name:  "**Alert Role**",
				Value: guildInfo.GiveawayRole,
			},
			{
				Name:  "**Claiming Message**",
				Value: guildInfo.GiveawayClaiming,
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}

}

func GiveawayCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	guildInfo := database.GuildInfo{
		GuildID: i.GuildID,
	}
	database.DB.First(&guildInfo, i.GuildID)

	options := ParseSubCommand(i)

	item := options["item"].(string)
	winners := options["winners"].(float64)
	channel := i.ChannelID
	duration := options["duration"].(string)
	canHost := false

	if guildInfo.GiveawayHostRole != "" {
		for _, v := range i.Member.Roles {
			if v == guildInfo.GiveawayHostRole {
				canHost = true
			}
		}
		if i.Member.Permissions&(1<<3) == 8 {
			canHost = true
		}
		if !canHost {
			ErrorResponse(s, i, "User must be administrator or have the role <@&"+guildInfo.GiveawayHostRole+"> to host giveaways.")
			return
		}
	}

	if winners <= 0 {
		ErrorResponse(s, i, "Must have 1 or more winners. Winners entered: "+fmt.Sprint(winners))
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

	description := fmt.Sprintf("**%s** Winners!", fmt.Sprint(winners))
	if options["description"] != nil {
		description = options["description"].(string)
	}
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		fmt.Println(err)
	}

	endTimeDuration, err := ParseTime(duration)
	if err != nil {
		fmt.Println(err)
		ErrorResponse(s, i, err.Error())
		return
	}

	endTime := time.Now().Add(endTimeDuration)

	content := ""
	if guildInfo.GiveawayRole != "" {
		content = fmt.Sprintf("<@&%s>", guildInfo.GiveawayRole)
	}

	message, err := PresetMessageSend(s, channel, PresetResponse{
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

	SuccessResponse(s, i, PresetResponse{
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
		options := ParseSubCommand(i)

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
			ErrorResponse(s, i, "You have already entered this giveaway!")
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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: i.Message.Embeds,
			Flags:  0,
		},
	})
}

func GiveawayEnd(s *discordgo.Session, messageID string) {

	var winnerList string
	giveawayInfo := database.Giveaway{
		MessageID: messageID,
	}
	database.DB.First(&giveawayInfo, messageID)

	for n := float64(0); n < giveawayInfo.Winners; {

		entryString := strings.Split(giveawayInfo.Entries, " ")

		if len(entryString) == 0 || giveawayInfo.Entries == "" {
			winnerList += fmt.Sprintf("• Not enough entrants!")
			fmt.Println("Not enough Entrants")
			break
		}

		result, err := random.Int(random.Reader, big.NewInt(int64(len(entryString))))
		if err != nil {
			fmt.Println("Random Error: ", err)
			return
		}

		winnerID := entryString[result.Int64()]

		user, err := s.User(winnerID)
		if err != nil {
			fmt.Println("User Error:", err)
			return
		}

		winner := fmt.Sprintf("<@%s> (%s#%s)", user.ID, user.Username, user.Discriminator)

		winnerList += fmt.Sprintf("• %s\n", winner)

		giveawayInfo.Entries = strings.Trim(strings.ReplaceAll(" "+giveawayInfo.Entries, " "+winnerID, ""), " ")

		n++
	}

	_, err := PresetMessageSend(s, giveawayInfo.ChannelID, PresetResponse{
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
						CustomID: "claim_giveaway",
					},
				},
			},
		},
		Title:       "Giveaway Completed!",
		Description: giveawayInfo.Description,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Giveaway Host**",
				Value:  fmt.Sprintf("This Giveaway was hosted by: <@!%s>", giveawayInfo.Host),
				Inline: true,
			},
			{
				Name:   "**Item Won**",
				Value:  giveawayInfo.Item,
				Inline: true,
			},
			{
				Name:   "**Winners**",
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
	}
}

func ClaimGiveawayButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var guildInfo database.GuildInfo
	database.DB.First(&guildInfo, i.GuildID)

	if guildInfo.GiveawayClaiming == "" {
		guildInfo.GiveawayClaiming = "The discord owner has not set a claiming message. Common ways to claim include: Opening a ticket or contacting the auction host. \nTo customize this message, use the command: `/giveaway setup claiming:`."
	}

	err := SuccessResponse(s, i, PresetResponse{
		Title:       "Claim Prize",
		Description: guildInfo.GiveawayClaiming,
	})
	if err != nil {
		fmt.Println(err)
	}
}