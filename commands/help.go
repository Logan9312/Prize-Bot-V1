package commands

import "github.com/bwmarrin/discordgo"

var HelpCommand = discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Shows the list of commands.",
}

func Help () {
	
}