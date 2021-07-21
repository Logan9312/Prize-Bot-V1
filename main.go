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
	{
		Name:        "help",
		Description: "Basic bot functionality",
	},
	{
		Name:        "inventory",
		Description: "Displays a user's inventory.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "username",
				Description: "Chose who's inventory to display",
				Required:    true,
				Choices:     []*discordgo.ApplicationCommandOptionChoice{},
				Options:     []*discordgo.ApplicationCommandOption{},
			},
		},
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
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "startingbid",
				Description: "Starting Bid Amount",
				Required:    true,
			},
		},
	},
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

	err = dg.UpdateGameStatus(0, "Aftermath Ark")

	if err != nil {
		fmt.Println(err)
	}

	for _, v := range slashCommands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, GuildID, v)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Command Finished")
	}

	commands.HelpBuilder(slashCommands)

	fmt.Println("Bot is running!")

	<-make(chan struct{})

}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type == 2 {
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
	}
	if i.Type == 3 {
		switch i.MessageComponentData().CustomID {
		case "Help":
			commands.HelpButton(s, i)
		case "auction1":
			commands.AuctionButtons(s, i)
		}
	}
}