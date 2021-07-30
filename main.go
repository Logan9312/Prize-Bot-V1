package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/m/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
)

var BotID string
var Prefix = "!"
var GuildID = "835209409109557289"
var AppID = "829527477268774953"

const Token string = "ODI5NTI3NDc3MjY4Nzc0OTUz.YG5bqg.5qESTPXLoiooMNTr3jUv_BXZWcY"

var slashCommands = []*discordgo.ApplicationCommand{
	&commands.DynamicVotingCommand,
	&commands.HelpCommand,
	&commands.ProfileCommand,
	&commands.AuctionCommand,
	&commands.SelectCommand,
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

	//Builds local commands
	_, err = dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, GuildID, slashCommands)
	if err != nil {
		fmt.Println(err)
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
			commands.Profile(s, i)
		case "bidtest":
			commands.BidTest(s, i, AppID)
		case "select-test":
			commands.Select(s, i)
		}
		switch i.ApplicationCommandData().Options[0].Name {
		case "create":
			commands.DynamicCreate(s, i)
		case "add":
			commands.DynamicEdit(s, i)
		}
	}
	if i.Type == 3 {
		switch i.MessageComponentData().CustomID {
		case "Help":
			commands.HelpButton(s, i)
		}
	}
}

type StatusOutput struct {
	Message string `json:"message"`
}
  
func HandleRequests() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/auction-bot/status", GetStatus).Methods("GET")
}
  
// GetStatus responds with the availability status of this service
func GetStatus(w http.ResponseWriter, r *http.Request) {
	status := StatusOutput{
	Message: "Bot is available",
	}
  
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(200)
json.NewEncoder(w).Encode(status)

}