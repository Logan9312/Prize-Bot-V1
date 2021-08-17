package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/commands"
)

var mainID = "829527477268774953"
var grungyID = "864930428639772692"

type slashCommands struct{
	local, prod, grungerson, auction []*discordgo.ApplicationCommand
}

func BotConnect(token, environment, botName string) {

	var c = slashCommands{
		local: []*discordgo.ApplicationCommand{
			&commands.HelpCommand,
			&commands.ProfileCommand,
			&commands.AuctionCommand,
		},
		prod: []*discordgo.ApplicationCommand{
			&commands.HelpCommand,
		},
		auction: []*discordgo.ApplicationCommand{},
		grungerson: []*discordgo.ApplicationCommand{
			&commands.ReviewCommand,
			&commands.ReviewEditCommand,
		},
	}
		


	var status string

	fmt.Println(botName + " Starting Up...")

	s, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = s.Open()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("User ID: " + s.State.User.ID)
	switch s.State.User.ID {
	case mainID:
		status = "Aftermath Ark"

	case grungyID:
		status = "suggon"
		c.prod = append(c.prod, c.grungerson...)
	}

	CommandBuilder(s, environment, c.local, c.prod)

	s.AddHandler(commands.CommandHandler)

	err = s.UpdateGameStatus(0, status)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer fmt.Println(botName + " Startup Complete!")
}
