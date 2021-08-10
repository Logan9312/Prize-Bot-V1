package grungerson

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var BotID string
var GuildID = "863122466485174323"
var AppID = "864930428639772692"

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

var slashCommands = []*discordgo.ApplicationCommand{
	{
		ApplicationID: AppID,
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

func GrungersonSetup (Token string) {
	dg, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := dg.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	dg.AddHandler(InteractionHandler)

	err = dg.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, GuildID, slashCommands)
	fmt.Println("Command Finished")

	fmt.Println("Bot is running")

	if err != nil {
		fmt.Println(err)
		return
	}

}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type == 2 {
		fmt.Println("Command Selecting")
		switch i.ApplicationCommandData().Name {
		case "review":
			Review(s, i)
		case "reviewedit":
			ReviewEdit(s, i)
		}
	}

}
