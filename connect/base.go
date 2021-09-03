package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gorm.io/gorm"
)

type slashCommands struct {
	local, prod []*discordgo.ApplicationCommand
}

func BotConnect(token, environment, botName string, db *gorm.DB) {

	var c = slashCommands{
		local: []*discordgo.ApplicationCommand{
			&commands.HelpCommand,
			&commands.ProfileCommand,
			&commands.AuctionCommand,
			&commands.SpawnExactDinoCommand,
			&commands.QueueCommand,
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

	//Builds local commands
	if environment == "local" {
		for _, v := range s.State.Guilds {
			_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, c.local)
			fmt.Println("Commands added to guild: " + v.Name)
			if err != nil {
				fmt.Println(err)
			}
		}
		commands.HelpBuilder(c.local)
	}

	//Builds prod commands
	if environment == "prod" {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", c.prod)
		if err != nil {
			fmt.Println(err)
		}
		commands.HelpBuilder(c.prod)
	}

	s.AddHandler(commands.CommandHandler)

	err = s.UpdateGameStatus(0, status)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(botName + " Startup Complete!")
}
