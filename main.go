package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"gitlab.com/logan9312/discord-auction-bot/connect"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	//"gitlab.com/logan9312/discord-auction-bot/database"
	"gitlab.com/logan9312/discord-auction-bot/routers"
	"github.com/bwmarrin/discordgo"
)

var prodCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
	&commands.SpawnExactDinoCommand,
}

var localCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,
	&commands.ProfileCommand,
	&commands.AuctionCommand,
}



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
	BotConnect(environment.DiscordToken, environment.Environment, "Main Bot")

	//Connects Sir Grungerson
	BotConnect(environment.Grungerson, environment.Environment, "Sir Grungerson")

	//Connects database
	//database.DatabaseConnect(environment.Host, environment.Password)

	fmt.Println("Bot is running! To stop, use: docker kill $(docker ps -q)")

	routers.BotStatus()
}

func BotConnect(token, environment, botName string) {

	fmt.Println(botName + " Starting Up...")

	s, err := discordgo.New("Bot " + token)

	prodCommands = append(prodCommands, commands.ReviewCommands[0], commands.ReviewCommands[1])
	
	if err != nil {
		fmt.Println(err)
		return
	}

	if environment == "local" {
		s.AddHandler(commands.CommandHandlerLocal)

		for _, v := range localCommands {
			v.Description = "EXPERIMENTAL: " + v.Description
		}

		for _, v := range s.State.Guilds {
			_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, v.ID, localCommands)
			fmt.Println("Commands added to guild: " + v.Name)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		commands.HelpBuilder(localCommands)
	}

	if environment == "prod" {

		s.AddHandler(commands.CommandHandlerProd)
		_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", prodCommands)
		if err != nil {
			fmt.Println(err)
			return
		}
		commands.HelpBuilder(prodCommands)
	}

	err = s.Open()

	if err != nil {
		fmt.Println(err)
		return
	}

	connect.BotStatus(s)

	defer fmt.Println(botName + " Startup Complete!")
}