package main

import (
	"fmt"

	"example.com/m/commands"
	"github.com/bwmarrin/discordgo"
)

const token string = "ODI5NTI3NDc3MjY4Nzc0OTUz.YG5bqg.5qESTPXLoiooMNTr3jUv_BXZWcY"

var BotID string
var Prefix = "!"
var GuildID = "835209409109557289"
var AppID = "829527477268774953"

var Scommands = []*discordgo.ApplicationCommand{
	{ID: "deletehelp",
		Name:        "help",
		Description: "Basic bot functionality",
		Version:     "",
		Options:     []*discordgo.ApplicationCommandOption{},
	},
	{Name: "auction",
		Description: "Put an item up for auction!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item",
				Description: "Choose an Item to put up for auction",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "bid",
				Description: "Starting Bid Amount",
				Required:    true,
			},
		},
	},
}

func main() {
	dg, err := discordgo.New("Bot " + token)

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

	fmt.Println("Bot is running!")

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
