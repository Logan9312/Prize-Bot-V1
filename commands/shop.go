package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
)

var ShopCommand = discordgo.ApplicationCommand{
	Name:        "shop",
	Description: "Buy, Sell and Trade items with others!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "setup",
			Description: "Setup shop settings on your server",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "sell",
			Description: "Select an item to sell",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The item you wish to sell",
					Required:    false,
					Choices:     []*discordgo.ApplicationCommandOptionChoice{},
					Options:     []*discordgo.ApplicationCommandOption{},
				},
			},
		},
	},
}

func Shop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[0].Name {
	case "sell":
		ShopCreate(s, i)
	}
}

func ShopCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	shopMap := h.ParseSubCommand(i)

	shopMap["guild_id"] = i.GuildID
	shopMap["host"] = i.Member.User.ID

	response, err := AuctionFormat(s, shopMap, "Shop")
	if err != nil {
		fmt.Println(err)
	}
	h.SuccessResponse(s, i, response)
}

func AddItem(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := h.SuccessResponse(s, i, h.PresetResponse{
		Content:     "",
		Title:       "",
		Description: "",
		Fields:      []*discordgo.MessageEmbedField{},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{},
		Image:       &discordgo.MessageEmbedImage{},
		Components:  []discordgo.MessageComponent{},
		Embeds:      []*discordgo.MessageEmbed{},
		Files:       []*discordgo.File{},
	})

	if err != nil {
		fmt.Println(err)
	}
}
