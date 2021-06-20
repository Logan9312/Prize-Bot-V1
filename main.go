package main

import (
	"fmt"

	"example.com/m/commands"
	"example.com/m/keys"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var Prefix = "!"
var GuildID = "835209409109557289"
var AppID = keys.AppID

const Token string = keys.Token

var slashCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand, 
	&commands.AuctionCommand,
	&commands.InventoryCommand,
	&commands.RolesCommand,
}

func main() {
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

	for _, v := range slashCommands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, GuildID, v)
			fmt.Println(err)
			if err != nil {
				fmt.Println(err)
			}
		fmt.Println("Command Finished")
	}

commands.HelpBuilder(slashCommands)

	fmt.Println("Bot is running")

	<-make(chan struct{})

}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.ApplicationCommandData().Name {
	case "help", "bid":
		commands.Help(s, i)
	case "auction":
		commands.Auction(s, i)
	case "inventory":
		commands.Inventory(s, i)
	case "roles":
		commands.Roles(s,i)
	}

	switch i.MessageComponentData().CustomID {
	case "Help":
		commands.HelpButton(s, i)
	case "auction1":
		commands.AuctionButtons(s, i)
	case "roles1":
		commands.RolesButton(s, i)
	}
}
