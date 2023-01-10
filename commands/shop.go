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

func Shop(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Options[0].Name {
	case "sell":
		return ShopCreate(s, i)
	}
	return fmt.Errorf("Unknown Shop command, please contact support")
}

func ShopCreate(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	shopMap := h.ParseSubCommand(i)

	shopMap["guild_id"] = i.GuildID
	shopMap["host"] = i.Member.User.ID

	response, err := EventFormat(s, shopMap, EventTypeShop, i.GuildID)
	if err != nil {
		fmt.Println(err)
	}
	h.SuccessResponse(s, i, response)
	return nil
}

func AddItem(s *discordgo.Session, i *discordgo.InteractionCreate) error {

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
	return nil
}
