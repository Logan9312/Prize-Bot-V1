package main

import (
	"fmt"
	"strings"

	"example.com/m/commands"
	"example.com/m/slashcommands"
	"github.com/bwmarrin/discordgo"
)

const token string = "ODI5NTI3NDc3MjY4Nzc0OTUz.YG5bqg.5qESTPXLoiooMNTr3jUv_BXZWcY"

var BotID string
var Prefix = "!"

var Scommands = []*discordgo.ApplicationCommand{
	{
		ID:          BotID,
		Name:        "help",
		Description: "A Helpful Command",
	},
	{
		Name:        "auction",
		Description: "Put an item up for auction!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "item",
				Description: "Choose an Item to put up for auction",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "apples",
						Value: "1",
					},
				},
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

	dg.AddHandler(MessageHandler)
	dg.AddHandler(InteractionHandler)

	err = dg.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range Scommands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, "840768614142967809", v)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Command Applied")
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

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//commands.AuctionButton(s, i)

	switch i.ApplicationCommandData().Name {
	case "help":
		slashcommands.Help(s, i)
	case "Auction":
		slashcommands.Auction(s, i)
	}

	switch i.MessageComponentData().CustomID {
	case "Help":
		slashcommands.HelpButton(s, i)
	case "Auction":
		slashcommands.AuctionButton(s, i)
	}
}
