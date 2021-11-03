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
			Options:     []*discordgo.ApplicationCommandOption{},
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
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The channel where the giveaway will take place",
					Required:    true,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
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
	if i.Member.Permissions&(1<<3) != 8 {
		ErrorResponse(s, i, "User must have administrator permissions to run this command")
		return
	}

	guildInfo := database.GuildInfo{
		GuildID: i.GuildID,
	}

	result := database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "guild_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"giveaway_role": guildInfo.GiveawayRole}),
	}).Create(&guildInfo)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
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
	channel := options["channel"].(string)
	duration := options["duration"].(string)

	var imageURL string
	if options["image"] != nil {
		imageURL = options["image"].(string)
	}

	description := "No Description"
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

	message, err := PresetMessageSend(s, channel, PresetResponse{
		Content:     fmt.Sprintf("<@&%s>", guildInfo.GiveawayRole),
		Title:       "Giveaway Item: __**" + item + "**__",
		Description: fmt.Sprintf("<@%s> has hosted a giveaway! To enter, press the button below!", i.Member.User.ID),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "__**Details:**__",
				Value:  description,
				Inline: false,
			},
			{
				Name:   "**End Time**",
				Value:  fmt.Sprintf("<t:%s:R>", endTime),
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

	giveawayInfo.Entries += i.Member.User.ID

	SuccessResponse(s, i, PresetResponse{
		Content:     "",
		Title:       "Success!",
		Description: "You have joined the giveaway",
		Fields:      []*discordgo.MessageEmbedField{},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{},
		Image:       &discordgo.MessageEmbedImage{},
		Components:  []discordgo.MessageComponent{},
	})
}

func GiveawayEnd(s *discordgo.Session, messageID string) {

	giveawayInfo := database.Giveaway{
		MessageID: messageID,
	}
	database.DB.First(&giveawayInfo, messageID)

	entryString := strings.Split(giveawayInfo.Entries, " ")

	result, err := random.Int(random.Reader, big.NewInt(int64(len(entryString))))
	if err != nil {
		fmt.Println(err)
		return
	}

	winnerID := entryString[result.Int64()]

	user, err := s.User(winnerID)
	if err != nil {
		fmt.Println(err)
		return
	}

	winner := fmt.Sprintf("<@%s> (%s#%s)", user.ID, user.Username, user.Discriminator)

	PresetMessageSend(s, giveawayInfo.ChannelID, PresetResponse{
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
		Title:       "Giveaway Completed!",
		Description: giveawayInfo.Description,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Winner**",
				Value:  winner,
				Inline: true,
			},
			{
				Name:   "**Giveaway Host**",
				Value:  fmt.Sprintf("This Auction was hosted by: <@!%s>", giveawayInfo.Host),
				Inline: false,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: giveawayInfo.ImageURL,
		},
	})
}
