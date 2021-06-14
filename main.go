package main

import (
	"fmt"
	"strings"

	"example.com/m/commands"
	"github.com/bwmarrin/discordgo"
)

const token string = "ODI5NTI3NDc3MjY4Nzc0OTUz.YG5bqg.5qESTPXLoiooMNTr3jUv_BXZWcY"

var BotID string
var Prefix = "!"

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

	dg.AddHandler(MessageHandler)
	dg.AddHandler(InteractionHandler)

	err = dg.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	<-make(chan struct{})

}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	CommandSlice := strings.Split(m.Content, " ")

	if m.Author.ID == BotID {
		return
	}

	Builder(s, m, CommandSlice)
}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commands.AuctionButton(s, i)
}

func Builder(s *discordgo.Session, m *discordgo.MessageCreate, CommandSlice []string) {

	if strings.HasPrefix(CommandSlice[0], Prefix) {
		fmt.Println("Command is being Built")
		phrase := strings.TrimPrefix(CommandSlice[0], Prefix)
		switch phrase {
		case "auction":
			commands.Auction(s, m, CommandSlice[1:])
		case "help":
			commands.Help(s, m, CommandSlice[1:]) 
		}
	}
}
