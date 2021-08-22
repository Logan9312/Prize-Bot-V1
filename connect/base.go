package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/commands"
)

type slashCommands struct{
	local, prod []*discordgo.ApplicationCommand
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

	status = "Aftermath Ark"

	CommandBuilder(s, environment, c.local, c.prod)

	s.AddHandler(commands.CommandHandler)
	s.AddHandler(commands.MessageHandler)

	err = s.UpdateGameStatus(0, status)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(botName + " Startup Complete!")
}
