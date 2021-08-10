package commands

import (
	"fmt"
	"strings"

	"math"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var reviewOptions = []*discordgo.ApplicationCommandOption{
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "game-name",
		Description: "The name of the game you are reviewing",
		Required:    true,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "genre",
		Description: "The type of game",
		Required:    true,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Action",
				Value: "Action",
			},
			{
				Name:  "Action-Adventure",
				Value: "Action-Adventure",
			},
			{
				Name:  "Adventure",
				Value: "Adventure",
			},
			{
				Name:  "Role Playing",
				Value: "Role Playing",
			},
			{
				Name:  "Simulation",
				Value: "Sim",
			},
			{
				Name:  "Strategy",
				Value: "Strategy",
			},
			{
				Name:  "Sports",
				Value: "Sports",
			},
			{
				Name:  "MMO",
				Value: "MMO",
			},
			{
				Name:  "Sandbox/Open World",
				Value: "Sandbox",
			},
			{
				Name:  "Other",
				Value: "Other",
			},
		},
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "platform",
		Description: "The platform the game is played on",
		Required:    true,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "PC and Console",
				Value: "PC and Console",
			},
			{
				Name:  "Console",
				Value: "Console",
			},
			{
				Name:  "PC",
				Value: "PC",
			},
			{
				Name:  "Xbox Exclusive",
				Value: "Xbox Exclusive",
			},
			{
				Name:  "Playstation Exclusive",
				Value: "Playstation Exclusive",
			},
			{
				Name:  "Switch",
				Value: "Switch",
			},
			{
				Name:  "Mobile",
				Value: "Mobile",
			},
			{
				Name:  "KFConsole",
				Value: "KFConsole",
			},
			{
				Name:  "Other",
				Value: "Other",
			},
		},
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "price",
		Description: "Enter the price in USD, I'm working on supporting multiple currencies.",
		Required:    true,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "rating",
		Description: "A rating from 0-10. Numbers outside the range will return an error",
		Required:    true,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "reason",
		Description: "Explain your rating",
		Required:    true,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "image",
		Description: "Post a link to the image of the game",
		Required:    false,
	},
}

var editOptions = []*discordgo.ApplicationCommandOption{
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "message-id",
		Description: "The ID of the message you want to edit",
		Required:    true,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "game-name",
		Description: "The name of the game you are reviewing",
		Required:    false,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "genre",
		Description: "The type of game",
		Required:    false,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Action",
				Value: "Action",
			},
			{
				Name:  "Action-Adventure",
				Value: "Action-Adventure",
			},
			{
				Name:  "Adventure",
				Value: "Adventure",
			},
			{
				Name:  "Role Playing",
				Value: "Role Playing",
			},
			{
				Name:  "Simulation",
				Value: "Sim",
			},
			{
				Name:  "Strategy",
				Value: "Strategy",
			},
			{
				Name:  "Sports",
				Value: "Sports",
			},
			{
				Name:  "MMO",
				Value: "MMO",
			},
			{
				Name:  "Sandbox/Open World",
				Value: "Sandbox",
			},
			{
				Name:  "Other",
				Value: "Other",
			},
		},
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "platform",
		Description: "The platform the game is played on",
		Required:    false,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "PC and Console",
				Value: "PC and Console",
			},
			{
				Name:  "Console",
				Value: "Console",
			},
			{
				Name:  "PC",
				Value: "PC",
			},
			{
				Name:  "Xbox Exclusive",
				Value: "Xbox Exclusive",
			},
			{
				Name:  "Playstation Exclusive",
				Value: "Playstation Exclusive",
			},
			{
				Name:  "Switch",
				Value: "Switch",
			},
			{
				Name:  "Mobile",
				Value: "Mobile",
			},
			{
				Name:  "KFConsole",
				Value: "KFConsole",
			},
			{
				Name:  "Other",
				Value: "Other",
			},
		},
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "price",
		Description: "Enter the price in USD, I'm working on supporting multiple currencies.",
		Required:    false,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "rating",
		Description: "A rating from 0-10. Numbers outside the range will return an error",
		Required:    false,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "reason",
		Description: "Why did you choose that rating?",
		Required:    false,
	},
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "image",
		Description: "Post a link to the image of the game",
		Required:    false,
	},
}

var reviewCommands = []*discordgo.ApplicationCommand{
	{
		Name:          "review",
		Description:   "Leave your review here!",
		Options:       reviewOptions,
	},
	{
		Name:        "reviewedit",
		Description: "This command is WIP and does not function properly now",
		Options:     editOptions,
	},
}

func Review(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var channelID string
	var Content string

	var OutOfRange bool = false

	name := i.ApplicationCommandData().Options[0].StringValue()
	genre := i.ApplicationCommandData().Options[1].StringValue()
	platform := i.ApplicationCommandData().Options[2].StringValue()
	price := strings.ReplaceAll(i.ApplicationCommandData().Options[3].StringValue(), " ", "")
	rating := strings.ReplaceAll(i.ApplicationCommandData().Options[4].StringValue(), " ", "")
	reason := i.ApplicationCommandData().Options[5].StringValue()

	if strings.ToLower(price) == "free" {
		price = "0"
	}

	price = strings.Trim(price, "$")

	floatPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println(err)
		Content = "There has been an error with the price causing your review to fail. **Make sure you only use numbers for the price**. Error Message: " + err.Error()
		OutOfRange = true
	}

	floatRating, err := strconv.ParseFloat(rating, 64)
	if err != nil {
		fmt.Println(err)
		Content = "There has been an error with the rating causing your review to fail. **Make sure you only use numbers for the rating**. Error Message: " + err.Error()
		OutOfRange = true
	}

	if price == "0" {
		price = "Free"
	} else {
		price = "$" + price + " USD"
	}

	percentage := fmt.Sprintf("%d%%", int(floatRating*10))

	embed := &discordgo.MessageEmbed{
		Type:  "rich",
		Title: fmt.Sprintf("**__Review for__: %s**", strings.Title(strings.ToLower(name))),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Review by: " + i.Member.User.Username + " | ID: " + i.Member.User.ID,
			IconURL: i.Member.User.AvatarURL(""),
		},
		Description: "Use the /reviewedit command to change any info. Voting coming soon.",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: i.Member.User.AvatarURL(""),
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Sir Grungerson Reviews",
			IconURL: "https://media.discordapp.net/attachments/863122466485174325/864938572372115456/image0.png?width=481&height=473",
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://media.discordapp.net/attachments/859512244956823612/861860829511024640/INVISIBLE_BACKGROUND.png?width=1025&height=9",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "**Rating:**",
				Value:  fmt.Sprintf("**%s/10** (%s)", rating, percentage),
				Inline: false,
			},
			{
				Name:   "**Review:**",
				Value:  reason,
				Inline: false,
			},
			{
				Name:   "**Genre:**",
				Value:  genre,
				Inline: true,
			},
			{
				Name:   "**Platform:**",
				Value:  platform,
				Inline: true,
			},
			{
				Name:   "**Game Price:**",
				Value:  price,
				Inline: true,
			},
		},
	}

	if len(i.ApplicationCommandData().Options) == 7 {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: i.ApplicationCommandData().Options[6].StringValue(),
		}
	}

	if i.GuildID == "835209409109557289" {
		channelID = "837391772736159754"
	}

	if i.GuildID == "863122466485174323" {
		channelID = "863201827033710632"
	}

	switch {
	case math.Signbit(floatRating):
		Content = fmt.Sprintf("If *Crypt of the Serpent King* can get 73%% I'm sure this game can get better than %s. **Please try again with a number between 0 and 10**", percentage)
		OutOfRange = true
	case math.Signbit(floatPrice):
		Content = "<@" + i.Member.User.ID + "> I know you think the game is trash but no one would pay *you* to play it. **Please try again with a positive cash value**"
		OutOfRange = true
	case floatRating == 0:
		fmt.Println(floatRating)
	case floatRating <= 1:
		embed.Color = 0xFF0000
	case floatRating <= 2:
		embed.Color = 0xFF4000
	case floatRating <= 3:
		embed.Color = 0xFF6200
	case floatRating <= 4:
		embed.Color = 0xFF7400
	case floatRating <= 5:
		embed.Color = 0xFFB100
	case floatRating <= 6:
		embed.Color = 0xFFAB00
	case floatRating <= 7:
		embed.Color = 0xFFCF00
	case floatRating <= 8:
		embed.Color = 0xFFFF00
	case floatRating <= 9:
		embed.Color = 0xCDFF00
	case floatRating < 10:
		embed.Color = 0x9CFF00
	case floatRating == 10:
		embed.Color = 0x00FF00
	case floatRating > 10:
		Content = "<@" + i.Member.User.ID + "> Dude, it wasn't that good. **Please try again with a rating between 0 and 10**."
		fmt.Println(floatRating)
		OutOfRange = true
	}

	if floatPrice > 100 {
		s.ChannelMessageSend(i.ChannelID, "Stop flexing your over $100 games! <@"+i.Member.User.ID+">")
	}

	if !OutOfRange {
		_, err = s.ChannelMessageSendEmbed(channelID, embed)

		if err != nil {
			fmt.Println(err)
			fmt.Println("embed could not be sent")
		}

		Content = "<@" + i.Member.User.ID + "> your review has been successfully submitted!"

	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: Content,
			Flags:   64,
		},
	})

}
