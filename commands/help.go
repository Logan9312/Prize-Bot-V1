package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"

)

func Help (s *discordgo.Session, m *discordgo.MessageCreate, details []string) {


	helpmessage := &discordgo.MessageSend{
		Content:         "Test",
		Embed:           &discordgo.MessageEmbed{
			Title:       "Discord Bot Help",
			Description: "This is the help function",
		},
	}
	
	_, err := s.ChannelMessageSendComplex(m.ChannelID, helpmessage)
	fmt.Println(err)
}