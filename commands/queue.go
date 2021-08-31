package commands

import (
	"github.com/bwmarrin/discordgo"
)

var QueueCommand = discordgo.ApplicationCommand{
	Name:        "queue",
	Description: "Queue commands",
	Options:     []*discordgo.ApplicationCommandOption{},
}

func Queue(s *discordgo.Session, i *discordgo.InteractionCreate) {

}
