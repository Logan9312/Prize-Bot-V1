package connect

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	c "gitlab.com/logan9312/discord-auction-bot/commands"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gitlab.com/logan9312/discord-auction-bot/logger"
	"gorm.io/gorm/clause"
)

var WhitelabelCommand = discordgo.ApplicationCommand{
	Name:        "whitelabel",
	Description: "Configure your custom bot",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add_token",
			Description: "Add or swap your custom bot token.",
		},
	},
}

func Whitelabel(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	if !c.CheckPremiumUser(i.Member.User.ID) {
		return h.PremiumError(s, i, "Whitelabelling is restricted to premium users only")
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "add_token":
		return WhitelabelToken(s, i)
	}
	return fmt.Errorf("Unknown Currency command, please contact support")
}

func WhitelabelToken(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "token",
							Label:       "Add your token here",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "MTAwNzQyMjQ3MzYyNjczMDUxNg.GYviWS.r9elKUmUCCYx7wT3PnVDS9GoP8DbLYiULsNsLo",
							Required:    true,
						},
					},
				},
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			CustomID:        "whitelabel_token",
			Title:           "Custom Bot Token",
		},
	})
}

func WhitelabelTokenModal(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	token := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	newSession, err := BotConnect(token, "prod")
	if err != nil {
		newSession.Close()
		return fmt.Errorf("Error connecting bot: %w", err)
	}

	whitelabelData := map[string]any{
		"bot_id":    newSession.State.User.ID,
		"user_id":   i.Member.User.ID,
		"bot_token": token,
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.WhiteLabels{}).Create(whitelabelData)
	if result.Error != nil {
		logger.Sugar.Errorw("failed to create whitelabel database entry", "user_id", i.Member.User.ID, "error", result.Error)
		return fmt.Errorf("failed to save bot configuration. Please try again or contact support")
	}

	result = database.DB.Model(database.WhiteLabels{}).Where(map[string]any{
		"user_id": i.Member.User.ID,
	}).Updates(whitelabelData)
	if result.Error != nil {
		logger.Sugar.Errorw("failed to update whitelabel data", "user_id", i.Member.User.ID, "error", result.Error)
		return fmt.Errorf("failed to save bot data. Please try again or contact support")
	}

	return h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Custom Bot Token Added",
		Description: fmt.Sprintf("<@%s> should now be functional.", newSession.State.User.ID),
	})
}
