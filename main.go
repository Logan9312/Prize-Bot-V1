package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/gorilla/mux"
	"gitlab.com/logan9312/discord-auction-bot/database"
	"gitlab.com/logan9312/discord-auction-bot/connect"
)

// Environment struct
type Environment struct {
	Environment  string `env:"ENVIRONMENT,required"`
	DiscordToken string `env:"DISCORD_TOKEN,required"`
	Migrate      bool   `env:"MIGRATE"`
	Host 		 string `env:"DB_HOST"`
	Password	 string `env:"DB_PASSWORD"`
	Grungerson 	string 	`env:"GRUNGERSON"`
}

func main() {

	environment := Environment{}

	if err := env.Parse(&environment); err != nil {
		log.Fatal("FAILED TO LOAD ENVIRONMENT VARIABLES")
	}

	//Connects main bot
	connect.BotConnect(environment.DiscordToken, environment.Environment)

	//Connects Sir Grungerson
	connect.BotConnect(environment.Grungerson, environment.Environment)

	//Connects database
	database.DatabaseConnect(environment.Host, environment.Password)

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
