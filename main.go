package main

import (
	"fmt"

	"example.com/m/commands"
	"example.com/m/cfunctions"
	"example.com/m/keys"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var Prefix = "!"
var GuildID = "835209409109557289"
var AppID = keys.AppID
const Token string = keys.Token

var Scommands = []*discordgo.ApplicationCommand{
	cfunctions.Help, 
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

	for _, v := range Scommands {
        acc, err := dg.ApplicationCommandCreate(dg.State.User.ID, GuildID, v)
        if err != nil {
            fmt.Println(err)
		err = dg.ApplicationCommandDelete(AppID, "", acc.ID)
		if err != nil {
            fmt.Println(err)
        	}
		}
	}



	<-make(chan struct{})

}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.ApplicationCommandData().Name {
	case "help":
		commands.Help(s, i)
	case "auction":
		commands.Auction(s, i)
	}

	switch i.MessageComponentData().CustomID {
	case "Help":
		commands.HelpButton(s, i)
	case "auction1", "auction2", "auction3":
		commands.AuctionButton(s, i)
	}
}
