package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"gitlab.com/logan9312/discord-auction-bot/connect"
	"gitlab.com/logan9312/discord-auction-bot/database"
	"gitlab.com/logan9312/discord-auction-bot/routers"
)

// Environment struct
type Environment struct {
	Environment  string `env:"ENVIRONMENT,required"`
	DiscordToken string `env:"DISCORD_TOKEN,required"`
	Migrate      bool   `env:"MIGRATE"`
	Host         string `env:"DB_HOST"`
	Password     string `env:"DB_PASSWORD"`
}

func main() {

	environment := Environment{}

	if err := env.Parse(&environment); err != nil {
		fmt.Println(err)
		log.Fatal("FAILED TO LOAD ENVIRONMENT VARIABLES")
	}

	//Connects database
	database.DatabaseConnect(environment.Password, environment.Host, environment.Environment)

	//Connects main bot
	connect.BotConnect(environment.DiscordToken, environment.Environment, "main")

	fmt.Println("Bot is running!")

	routers.BotStatus()
}
