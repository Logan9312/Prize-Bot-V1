package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"gitlab.com/logan9312/discord-auction-bot/commands"
	//"gitlab.com/logan9312/discord-auction-bot/database"
	"gitlab.com/logan9312/discord-auction-bot/routers"
	"github.com/bwmarrin/discordgo"
)

var prodCommands = []*discordgo.ApplicationCommand{
	&commands.HelpCommand,

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

	s, err := discordgo.New("Bot " + environment.DiscordToken)
	if err != nil {
		fmt.Println(err)
	}

	//Connects main bot
	//BotConnect(environment.DiscordToken, environment.Environment, "Main Bot")

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
	

	//Connects Sir Grungerson
	// BotConnect(environment.Grungerson, environment.Environment, "Sir Grungerson")

	//Connects database
	//database.DatabaseConnect(environment.Host, environment.Password)

	fmt.Println("Bot is running! To stop, use: docker kill $(docker ps -q)")

	routers.BotStatus()
}

