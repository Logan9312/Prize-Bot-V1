package commands

import (
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestAuction(t *testing.T) {

	guildID := "915767892467920967"

	s, err := discordgo.New("Bot " + "ODcyNjE1MzMxMTY2ODM0Nzc4.YQscWw.tEeyle-QBEVLyWW0UX1WCqEDWQw")
	if err != nil {
		fmt.Println("discordgo.New error:" + err.Error())
	}

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages
	err = s.Open()
	if err != nil {
		fmt.Println("s.Open error: " + err.Error())
		return
	}

	member, err := s.GuildMember(guildID, "280812467775471627")
	if err != nil {
		t.Error(err)
	}

	auctionMaps := []map[string]any{
		{
			"item":          "Automated Test Auction",
			"bid":           500,
			"duration":      "30s",
			"description":   "Test Description",
			"currency":      "$$",
			"integer_only":  true,
			"increment_max": -1000,
			"increment_min": -500,
			"buyout":        0,
			"target_price":  0,
			"category":      "915768615742103625",
			"schedule":      "30s",
		},
		{
			"item":          "Automated Test Auction",
			"bid":           -500,
			"duration":      "-1m",
			"description":   "",
			"currency":      "123",
			"integer_only":  false,
			"increment_max": 0,
			"increment_min": 0,
			"buyout":        -400,
			"target_price":  -600,
			"category":      "",
			"schedule":      "0s",
		},
		{
			"item":          "Automated Test Auction",
			"bid":           0,
			"duration":      "30s",
			"description":   nil,
			"currency":      nil,
			"integer_only":  nil,
			"increment_max": nil,
			"increment_min": nil,
			"buyout":        nil,
			"target_price":  nil,
			"category":      nil,
			"schedule":      nil,
		},
	}

	for _, auctionMap := range auctionMaps {
		_, err = AuctionHandler(s, auctionMap, member, guildID)
		if err != nil {
			t.Error(err)
		}
	}
}
