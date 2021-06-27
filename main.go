package main

import (
	"fmt"

	"example.com/m/commands"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var Prefix = "!"
var GuildID = "835209409109557289"
var AppID = "829527477268774953"

const Token string = "ODI5NTI3NDc3MjY4Nzc0OTUz.YG5bqg.5qESTPXLoiooMNTr3jUv_BXZWcY"

var slashCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand, 
	&commands.AuctionCommand,
	&commands.InventoryCommand,
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
			if err != nil {
				fmt.Println(err)
			}
		fmt.Println("Command Finished")
	}

	commands.HelpBuilder(slashCommands)

	fmt.Println("Bot is running")

	if err != nil {
		fmt.Println()
		return
	}

	<-make(chan struct{})

}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	fmt.Println("Selecting Function")

	switch i.ApplicationCommandData().Name {
	case "help":
		commands.Help(s, i)
	case "auction":
		commands.Auction(s, i, AppID)
	case "inventory":
		commands.Inventory(s, i)
	case "bidtest":
		commands.BidTest(s, i, AppID)
	}

	switch i.MessageComponentData().CustomID {
	case "Help":
		commands.HelpButton(s, i)
	case "auction1":
		commands.AuctionButtons(s, i)
	}
}
