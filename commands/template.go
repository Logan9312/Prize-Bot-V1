package commands

import (
	"github.com/bwmarrin/discordgo"
)

var CommandNameCommand = discordgo.ApplicationCommand{
	Name:        "auction",
	Description: "Put an item up for auction!",
	Options: []*discordgo.ApplicationCommandOption{},
}

func CommandName(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

