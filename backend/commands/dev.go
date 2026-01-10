package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/config"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gitlab.com/logan9312/discord-auction-bot/logger"
	"gorm.io/gorm/clause"
)

var DevCommand = discordgo.ApplicationCommand{
	Name:        "dev",
	Description: "Developer only commands!",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "version",
			Description: "Change the version of the bot",
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "resend_claim",
			Description: "Resend a claim message by message ID (format: message_id,channel_id)",
		},
	},
}

var StatsCommand = discordgo.ApplicationCommand{
	Name:        "stats",
	Description: "Developer only commands!",
}

func Dev(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := h.ParseSlashCommand(i)

	// Check for nil Member (DM context)
	if i.Member == nil || i.Member.User == nil {
		return fmt.Errorf("this command cannot be used in DMs")
	}

	if !config.IsDevUser(i.Member.User.ID) {
		return fmt.Errorf("User must be a developer to run this command")
	}

	result := database.DB.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Model(database.DevSetup{}).Create(map[string]interface{}{
		"bot_id": s.State.User.ID,
	})
	if result.Error != nil {
		logger.Sugar.Warnw("dev command database error", "error", result.Error)
		return result.Error
	}

	if options["version"] != nil {
		err := s.UpdateGameStatus(0, "Bot Version "+options["version"].(string))
		if err != nil {
			return err
		}
	}

	if options["resend_claim"] != nil {
		return devResendClaim(s, i, options["resend_claim"].(string))
	}

	result = database.DB.Model(database.DevSetup{
		BotID: s.State.User.ID,
	}).Updates(options)
	if result.Error != nil {
		logger.Sugar.Warnw("dev command database error", "error", result.Error)
		return result.Error
	}

	fields := []*discordgo.MessageEmbedField{}

	for key, value := range options {
		if value != nil && value != "" {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   key,
				Value:  fmt.Sprint(value),
				Inline: true,
			})
		}
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:  "Success!",
		Fields: fields,
	})

	return nil
}

func Stats(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	members := 0
	largestServerCount := 0
	largestServerName := ""

	for _, v := range s.State.Guilds {
		members += v.MemberCount
		if v.MemberCount > largestServerCount {
			largestServerCount = v.MemberCount
			largestServerName = v.Name
		}
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title: "Statistics",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Members",
				Value:  fmt.Sprint(members),
				Inline: false,
			},
			{
				Name:   "Servers",
				Value:  fmt.Sprint(len(s.State.Guilds)),
				Inline: false,
			},
			{
				Name:   largestServerName,
				Value:  fmt.Sprint(largestServerCount),
				Inline: false,
			},
		},
	})

	return nil
}

func devResendClaim(s *discordgo.Session, i *discordgo.InteractionCreate, input string) error {
	// Parse input: message_id,channel_id
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid format. Use: message_id,channel_id")
	}

	messageID := strings.TrimSpace(parts[0])
	channelID := strings.TrimSpace(parts[1])

	// Look up the claim by message_id
	claimMap := map[string]interface{}{}
	result := database.DB.Model(database.Claim{}).First(&claimMap, messageID)
	if result.Error != nil {
		return fmt.Errorf("claim not found with message_id: %s", messageID)
	}

	// Prepare the claim data for resending
	claimMap["log_channel"] = channelID
	claimMap["old_id"] = claimMap["message_id"]

	if claimMap["type"] == nil || claimMap["type"] == "" {
		claimMap["type"] = "Auction"
	}

	h.SuccessResponse(s, i, h.PresetResponse{
		Title:       "Resending Claim",
		Description: fmt.Sprintf("Resending claim for **%s**...", claimMap["item"]),
	})

	// Resend the claim
	_, err := ClaimOutput(s, claimMap, claimMap["type"].(string))
	if err != nil {
		logger.Sugar.Errorw("failed to resend claim", "message_id", messageID, "error", err)
		h.FollowUpErrorResponse(s, i, fmt.Sprintf("Failed to resend claim: %s", err.Error()))
		return nil
	}

	h.FollowUpSuccessResponse(s, i, h.PresetResponse{
		Title:       "Claim Resent",
		Description: fmt.Sprintf("Claim for **%s** has been resent to <#%s>", claimMap["item"], channelID),
	})
	return nil
}
