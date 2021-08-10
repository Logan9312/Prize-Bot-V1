package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/database"
	"gitlab.com/logan9312/discord-auction-bot/grungerson"
)

// Environment struct
type Environment struct {
	Environment  string `env:"ENVIRONMENT,required"`
	DiscordToken string `env:"DISCORD_TOKEN,required"`
	Migrate      bool   `env:"MIGRATE"`
	Host 		 string `env:"DB_HOST"`
	Password	 string `env:"DB_PASSWORD"`
}

var prodCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
	&commands.SpawnExactDinoCommand,
}

var localCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
	&commands.ProfileCommand,
	&commands.AuctionCommand,
}

func main() {

	environment := Environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal("FAILED TO LOAD ENVIRONMENT VARIABLES")
	}

	s, err := discordgo.New("Bot " + environment.DiscordToken)
	grungerson.GrungersonSetup()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if environment.Environment == "local" {
		s.AddHandler(commands.CommandHandlerLocal)
	}

	if environment.Environment == "prod" {
		s.AddHandler(commands.CommandHandlerProd)
	}

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

		for _, v := range localCommands {
			v.Description = "EXPERIMENTAL: " + v.Description
		}

		for _, v := range s.State.Guilds {
			_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, localCommands)
			fmt.Println("Commands added to guild: " + v.Name)
			if err != nil {
				fmt.Println(err)
			}
		}
		commands.HelpBuilder(localCommands)
	}

	//Builds prod commands
	if environment.Environment == "prod" {

		_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", prodCommands)
		if err != nil {
			fmt.Println(err)
		}
		commands.HelpBuilder(prodCommands)
	}

	database.DatabaseConnect (environment.Host, environment.Password)

	fmt.Println("Bot is running! To stop, use: docker kill $(docker ps -q)")

	r := mux.NewRouter().StrictSlash(true)
	HandleRequests(r)
	log.Fatal(http.ListenAndServe(":8080", r))
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
