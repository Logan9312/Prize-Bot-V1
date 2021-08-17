package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/commands"
)

func CommandBuilder(s *discordgo.Session, environment string, localCommands, prodCommands []*discordgo.ApplicationCommand) {

	//Builds local commands
	if environment == "local" {
		for _, v := range s.State.Guilds {
			_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, localCommands)
			fmt.Println("Commands added to guild: " + v.Name)
			if err != nil {
				fmt.Println(err)
			}
		}
		commands.HelpBuilder(localCommands)
	}

	//Builds prod commands
	if environment == "prod" {

		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", prodCommands)
		if err != nil {
			fmt.Println(err)
		}
		commands.HelpBuilder(prodCommands)
	}
}
