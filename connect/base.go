package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/commands"
)

var mainID = "829527477268774953"
var grungyID = "864930428639772692"

var prodCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
}

var localCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
	&commands.ProfileCommand,
	&commands.AuctionCommand,

}

var status string

var MainCommands = []*discordgo.ApplicationCommand{}

var GrungyCommands = []*discordgo.ApplicationCommand{
	&commands.ReviewCommand,
	&commands.ReviewEditCommand,
}

func BotConnect(token, environment, botName string) {

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

	switch s.State.User.ID {
	case mainID:
		status = "Aftermath Ark"
		prodCommands = append(prodCommands, MainCommands...)
	case grungyID:
		status = "Grungy's Gaming Garage"
		prodCommands = append(prodCommands, GrungyCommands...)
	}

	CommandBuilder(s, environment, localCommands, prodCommands)

	err = s.UpdateGameStatus(0, status)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer fmt.Println(botName + " Startup Complete!")
}
