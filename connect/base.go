package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/commands"
)

var prodCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
}

var localCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
	&commands.ProfileCommand,
	&commands.AuctionCommand,
}

func BotConnect(token, environment, botName string) {

	fmt.Println(botName + " Starting Up...")

	s, err := discordgo.New("Bot " + token)

	prodCommands = append(prodCommands, commands.ReviewCommands[0], commands.ReviewCommands[1])

	if err != nil {
		fmt.Println(err)
		return
	}

	//Builds local commands
	if environment == "local" {

		s.AddHandler(commands.CommandHandlerLocal)

		for _, v := range localCommands {
			v.Description = "EXPERIMENTAL: " + v.Description
		}

		for _, v := range s.State.Guilds {
			_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, localCommands)
			fmt.Println("Commands added to guild: " + v.Name)
			if err != nil {
				fmt.Println(err)
			}
		}
		commands.HelpBuilder(localCommands)
	}

	//Builds prod commands
	if environment == "prod" {

		s.AddHandler(commands.CommandHandlerProd)

		_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", prodCommands)
		if err != nil {
			fmt.Println(err)
		}
		
		commands.HelpBuilder(prodCommands)
	}

	err = s.Open()

	if err != nil {
		fmt.Println(err)
		return
	}

	BotStatus(s)

	defer fmt.Println(botName + " Startup Complete!")
}
