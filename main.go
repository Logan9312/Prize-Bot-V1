package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"example.com/m/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
)

var BotID string
var Prefix = "!"
var GuildID = "835209409109557289"

// Environment struct
type Environment struct {
	Environment  string `env:"ENVIRONMENT,required"`
	DiscordToken string `env:"DISCORD_TOKEN,required"`
	Migrate      bool   `env:"MIGRATE"`
}

var slashCommands = []*discordgo.ApplicationCommand{
	&commands.DynamicVotingCommand,
	&commands.HelpCommand,
	&commands.ProfileCommand,
	&commands.AuctionCommand,
	&commands.SelectCommand,
}

func main() {
	environment := Environment{}

	if err := env.Parse(&environment); err != nil {
		log.Fatal("FAILED TO LOAD ENVIRONMENT VARIABLES")
	}

	s, err := discordgo.New("Bot " + environment.DiscordToken)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s.AddHandler(InteractionHandler)

	err = s.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = s.UpdateGameStatus(0, "Aftermath Ark")
	if err != nil {
		fmt.Println(err)
	}

	//Builds local commands
	if environment.Environment == "local" {
		for range s.State.Guilds {
			_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildID, slashCommands)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if environment.Environment == "prod" {
		_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildID, slashCommands)
		if err != nil {
			fmt.Println(err)
		}
	}

	commands.HelpBuilder(slashCommands)

	fmt.Println("Bot is running!")

	r := mux.NewRouter().StrictSlash(true)
	HandleRequests(r)
	log.Fatal(http.ListenAndServe(":8080", r))

}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type == 2 {
		switch i.ApplicationCommandData().Name {
		case "help":
			commands.Help(s, i)
		case "auction":
			commands.Auction(s, i, s.State.User.ID)
		case "inventory":
			commands.Profile(s, i)
		case "bidtest":
			commands.BidTest(s, i, s.State.User.ID)
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

func HandleRequests(r *mux.Router) {
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
