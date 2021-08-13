package commands

import (
	"github.com/bwmarrin/discordgo"
)

var QueueCommand = discordgo.ApplicationCommand{
	Name:        "auction",
	Description: "Put an item up for auction!",
	Options:     []*discordgo.ApplicationCommandOption{},
}

func Queue(s *discordgo.Session, i *discordgo.InteractionCreate) {

}
