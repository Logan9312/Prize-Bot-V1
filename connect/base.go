package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/commands"
)

var prodCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
	&commands.SpawnExactDinoCommand,
}

var localCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
	&commands.ProfileCommand,
	&commands.AuctionCommand,
}

var guilds []*discordgo.Guild

var slashCommands []*discordgo.ApplicationCommand

func BotConnect(token, environment string) {

	s, err := discordgo.New("Bot " + token)
	
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if environment == "local" {
		slashCommands = localCommands
		s.AddHandler(commands.CommandHandlerLocal)

		for _, v := range slashCommands {
			v.Description = "EXPERIMENTAL: " + v.Description
		}
		guilds = s.State.Guilds
	}

	if environment == "prod" {
		slashCommands = prodCommands
		s.AddHandler(commands.CommandHandlerProd)
		_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", prodCommands)
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, v := range guilds {
		_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, slashCommands)
		fmt.Println("Commands added to guild: " + v.Name)
		if err != nil {
			fmt.Println(err)
		}
	}

	commands.HelpBuilder(slashCommands)

	err = s.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotStatus(s)
}