package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/stripe/stripe-go/v72"
	"gitlab.com/logan9312/discord-auction-bot/connect"
	"gitlab.com/logan9312/discord-auction-bot/database"
	"gitlab.com/logan9312/discord-auction-bot/routers"
	"gitlab.com/logan9312/discord-auction-bot/website"
)

// Environment struct
type Environment struct {
	Environment  string `env:"ENVIRONMENT,required"`
	DiscordToken string `env:"DISCORD_TOKEN,required"`
	Migrate      bool   `env:"MIGRATE"`
	Host         string `env:"DB_HOST"`
	Password     string `env:"DB_PASSWORD"`
	StripeToken  string `env:"STRIPE_TOKEN"`
}

func main() {

	environment := Environment{}

	if err := env.Parse(&environment); err != nil {
		fmt.Println(err)
		log.Fatal("FAILED TO LOAD ENVIRONMENT VARIABLES")
	}

	stripe.Key = environment.StripeToken

	//Connects database
	database.DatabaseConnect(environment.Password, environment.Host, environment.Environment)

	//Connects main bot
	connect.BotConnect(environment.DiscordToken, environment.Environment, "main")

	fmt.Println("Bot is running!")
	go website.SetupWebsite()

	routers.BotStatus()
}
