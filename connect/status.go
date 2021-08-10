package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var status string 
func BotStatus(s *discordgo.Session) {

switch s.State.User.ID {
case "864930428639772692":
	status = "Grungy's Gaming Garage"
case "829527477268774953":
	status = "Aftermath Ark"
}

err := s.UpdateGameStatus(0, status)
if err != nil {
	fmt.Println(err)
	return
}

}