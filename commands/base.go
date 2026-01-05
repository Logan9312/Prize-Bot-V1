package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/logan9312/discord-auction-bot/database"
	h "gitlab.com/logan9312/discord-auction-bot/helpers"
	"gitlab.com/logan9312/discord-auction-bot/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// getWinnersCount safely extracts the winners count from an interface{}
// It handles both int (from database) and float64 (from JSON parsing)
func getWinnersCount(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}

const (
	EventTypeAuction  = "Auction"
	EventTypeShop     = "Shop"
	EventTypeGiveaway = "Giveaway"
)

func AddCommand() {

}

func EventFormat(s *discordgo.Session, data map[string]interface{}, eventType string, guildID string) (h.PresetResponse, error) {

	content := ""
	embeds := []*discordgo.MessageEmbed{}

	if data["item"] != nil && len(data["item"].(string)) > 100 {
		return h.PresetResponse{}, fmt.Errorf("title cannot be over 100 characters long")
	}

	description := fmt.Sprintf("**Host:** <@%s>.\n", data["host"])

	if data["description"] != nil {
		description += fmt.Sprintf("**Description:** %s\n", data["description"])
	}

	if data["winners"] != nil {
		description += fmt.Sprintf("**Winners:** %d\n", getWinnersCount(data["winners"]))
	}

	if data["increment_min"] != nil {
		description += fmt.Sprintf("**Minimum Bid:** + %s above previous.\n", PriceFormat(data["increment_min"].(float64), guildID, data["currency"]))
	}

	if data["increment_max"] != nil {
		description += fmt.Sprintf("**Maximum Bid:** + %s above previous.\n", PriceFormat(data["increment_max"].(float64), guildID, data["currency"]))
	}

	if data["target_price"] != nil {
		description += "**Target Price:** The host has set a hidden target price for this auction.\n"
	}

	if data["integer_only"] != nil {
		description += fmt.Sprintf("**Integer Only:** %t.\n", data["integer_only"].(bool))
	}

	if data["snipe_extension"] != nil && data["snipe_range"] != nil {
		antiSnipeDesc := fmt.Sprintf("**Anti Snipe:** If a bid is placed within the last %s, the auction will be extended by %s", data["snipe_range"], data["snipe_extension"].(time.Duration).String())
		if data["snipe_limit"] != nil {
			antiSnipeDesc += fmt.Sprintf(", up to a maximum total extension of %s", data["snipe_limit"].(time.Duration).String())
		}
		description += antiSnipeDesc + ".\n"
	}

	if data["buyout"] != nil {
		description += fmt.Sprintf("**Buyout Price:** %s.\n", PriceFormat(data["buyout"].(float64), guildID, data["currency"]))
	}

	auctionfields := []*discordgo.MessageEmbedField{
		{
			Name:  "__**Details:**__",
			Value: description,
		},
	}

	if data["end_time"] != nil {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**End Time**__",
			Value:  fmt.Sprintf("<t:%d:R>", data["end_time"].(time.Time).Unix()),
			Inline: true,
		})
	}

	if data["bid"] != nil {
		if data["winner"] != nil {
			auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
				Name:   "__**Current Highest Bid:**__",
				Value:  PriceFormat(data["bid"].(float64), guildID, data["currency"]),
				Inline: true,
			})
		} else {
			auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
				Name:   "__**Starting Bid:**__",
				Value:  PriceFormat(data["bid"].(float64), guildID, data["currency"]),
				Inline: true,
			})
		}
	}

	if data["winner"] != nil {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**Current Winner**__",
			Value:  fmt.Sprintf("<@%s>", data["winner"]),
			Inline: true,
		})
	}

	if eventType == EventTypeAuction {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:  "__**How to Bid**__",
			Value: "Use the /bid command or type `/bid {value}` in chat\nEx: `/bid 550`\n",
		})
	}

	if eventType == EventTypeGiveaway {
		auctionfields = append(auctionfields, &discordgo.MessageEmbedField{
			Name:   "__**How to Enter**__",
			Value:  "To enter, select the 🎁 reaction below! Removing your reaction will remove your entry.",
			Inline: false,
		})
	}

	guild, err := s.State.Guild(guildID)
	if err != nil {
		logger.Sugar.Warnw("error fetching guild", "error", err)
		return h.PresetResponse{}, err
	}

	if data["alert_role"] != nil {
		content = fmt.Sprintf("<@&%s>", strings.Trim(data["alert_role"].(string), " "))
	}
	if data["alert_role"] == guildID {
		content = "@everyone"
	}

	components := []discordgo.MessageComponent{}

	if eventType == EventTypeAuction {
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "End Auction",
						Style: 4,
						Emoji: &discordgo.ComponentEmoji{
							Name: "🛑",
						},
						CustomID: "endauction",
					},
					discordgo.Button{
						Label:    "Clear Chat",
						Style:    3,
						CustomID: "clearauction",
						Emoji: &discordgo.ComponentEmoji{
							Name: "restart",
							ID:   "835685528917114891",
						},
						Disabled: false,
					},
				},
			},
		}
	}

	if eventType == EventTypeShop {
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Purchase",
						Style:    discordgo.SuccessButton,
						CustomID: "additem",
					},
				},
			},
		}
	}

	if data["bid_history"] != nil {
		if len(data["bid_history"].(string)) > 4095 {
			data["bid_history"] = data["bid_history"].(string)[len(data["bid_history"].(string))-4095:]
		}
		embeds = []*discordgo.MessageEmbed{{
			Title:       "**Bid History**",
			Description: data["bid_history"].(string),
			Color:       0x8073ff,
		}}
	}
	EventMessage := h.PresetResponse{
		Content:    content,
		Title:      fmt.Sprintf("%s Item: __**%s**__", eventType, data["item"]),
		Fields:     auctionfields,
		Thumbnail:  &discordgo.MessageEmbedThumbnail{URL: guild.IconURL("")},
		Components: components,
		Embeds:     embeds,
		Files:      []*discordgo.File{},
	}

	if data["image_url"] != nil {
		EventMessage.Image = &discordgo.MessageEmbedImage{
			URL: data["image_url"].(string),
		}
	}

	return EventMessage, nil
}

func HasRole(m *discordgo.Member, roleID string) bool {
	if m == nil {
		return false
	}
	for _, role := range m.Roles {
		if role == roleID {
			return true
		}
	}
	return false
}

// CheckHostPermission checks if a member has permission to host events (auctions, giveaways, etc.)
// It checks:
// 1. If no host_role is set, returns true (anyone can host)
// 2. If host_role equals guildID, it represents @everyone, so anyone can host
// 3. If the member has the host_role
// 4. If the member has administrator permissions
func CheckHostPermission(member *discordgo.Member, hostRole interface{}, guildID string) bool {
	// If member is nil, deny permission
	if member == nil {
		return false
	}

	// If no host role is configured, anyone can host
	if hostRole == nil || hostRole == "" {
		return true
	}

	hostRoleStr, ok := hostRole.(string)
	if !ok || hostRoleStr == "" {
		return true
	}

	// If host_role equals guildID, it represents @everyone role
	if hostRoleStr == guildID {
		return true
	}

	// Check if member has the host role
	if HasRole(member, hostRoleStr) {
		return true
	}

	// Check if member has administrator permissions (1<<3 == 8)
	if member.Permissions&(1<<3) != 0 {
		return true
	}

	return false
}

// GetHostRoleErrorMessage returns a formatted error message for when a user lacks host permission
func GetHostRoleErrorMessage(hostRole interface{}, actionType string) string {
	hostRoleStr, ok := hostRole.(string)
	if !ok || hostRoleStr == "" {
		return fmt.Sprintf("User must have permission to %s", actionType)
	}
	return fmt.Sprintf("User must be administrator or have the role <@&%s> to %s", hostRoleStr, actionType)
}

// SetupClearConfig holds the configuration for a setup clear operation
type SetupClearConfig struct {
	SetupType   string // e.g., "Auction", "Giveaway", "Claim", "Currency"
	SetupCmd    string // e.g., "/settings auction", "/giveaway setup"
}

// GenericSetupClear handles clearing settings for any setup type
// model should be a pointer to the database struct (e.g., &database.AuctionSetup{GuildID: guildID})
func GenericSetupClear(s *discordgo.Session, i *discordgo.InteractionCreate, model interface{}, config SetupClearConfig) error {
	options := i.MessageComponentData().Values
	clearedMap := map[string]interface{}{}

	clearedSettings := "No Settings Cleared!"
	if len(options) > 0 {
		clearedSettings = ""
	}

	for _, v := range options {
		clearedSettings += fmt.Sprintf("• %s\n", cases.Title(language.English).String(strings.ReplaceAll(v, "_", " ")))
		clearedMap[v] = gorm.Expr("NULL")
	}

	result := database.DB.Model(model).Updates(clearedMap)
	if result.Error != nil {
		logger.Sugar.Errorw("failed to clear settings",
			"setup_type", config.SetupType,
			"guild_id", i.GuildID,
			"error", result.Error,
		)
		return fmt.Errorf("failed to clear settings. Please try again or contact support")
	}

	return h.SuccessResponse(s, i, h.PresetResponse{
		Title:       fmt.Sprintf("**Cleared %s Settings**", config.SetupType),
		Description: fmt.Sprintf("You have successfully cleared the following settings. Run `%s` to see your changes.", config.SetupCmd),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Cleared Settings**",
				Value: clearedSettings,
			},
		},
	})
}
