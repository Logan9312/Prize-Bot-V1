package commands 

import (
	"reflect"
	"strings"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type InfoType struct {}

func CommandHandlerLocal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var info InfoType

	input := make([]reflect.Value, 2)
	input[0] = reflect.ValueOf(s)
	input[1] = reflect.ValueOf(i)
	fmt.Println(input)

	if i.Type == 2 {
	name := strings.Title(i.ApplicationCommandData().Name)
	reflect.ValueOf(info).MethodByName(name).Call(input)
	} else if i.Type == 3 {
		name := strings.Title(i.MessageComponentData().CustomID) + "Button"
		reflect.ValueOf(info).MethodByName(name).Call(input)
	}
}

func CommandHandlerProd(s *discordgo.Session, i *discordgo.InteractionCreate) {
var info InfoType
	if i.Type == 2 {
		switch i.ApplicationCommandData().Name {
		case "help":
			info.Help(s, i)
		case "auction":
			info.Auction(s, i, s.State.User.ID)
		case "inventory":
			info.Profile(s, i)
		case "bid":
			info.Bid(s, i)
		default:
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds:          []*discordgo.MessageEmbed{
						{
							Title:       "Command Selection Error",
							Description: "Command response has not been set properly, please contact Logan to fix",
						},
					},
					Flags:           64,
				},
			})
		}
		/*switch i.ApplicationCommandData().Options[0].Name {
		case "create":
			commands.DynamicCreate(s, i)
		case "add":
			commands.DynamicEdit(s, i)
		}*/
	}
	if i.Type == 3 {
		switch i.MessageComponentData().CustomID {
			//auction bid alert button
		}
	}
}
