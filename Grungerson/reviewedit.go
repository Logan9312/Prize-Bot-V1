package grungerson

import (
	"fmt"
	"strings"

	"math"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var editContent string

func ReviewEdit(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var OutOfRange bool = false

	var name string = ""
	var genre string = ""
	var platform string = ""
	var price string = ""
	var rating string = ""
	var reason string = ""
	var messageID string = ""
	var thumbnail string = ""
	var FloatPrice float64
	var floatRating float64
	var percentage string = ""
	var err error
	var channelID string

	for _, v := range i.ApplicationCommandData().Options {
		switch v.Name {
		case "message-id":
			messageID = v.StringValue()
		case "game-name":
			name = v.StringValue()
		case "genre":
			genre = v.StringValue()
		case "platform":
			platform = v.StringValue()
		case "price":
			price = strings.ReplaceAll(v.StringValue(), " ", "")
		case "rating":
			rating = strings.ReplaceAll(v.StringValue(), " ", "")
		case "reason":
			reason = v.StringValue()
		case "image":
			thumbnail = v.StringValue()
		}
	}

	if i.GuildID == "835209409109557289" {
		channelID = "837391772736159754"
	}

	if i.GuildID == "863122466485174323" {
		channelID = "863201827033710632"
	}

	message, err := s.ChannelMessage(channelID, messageID)

	if err != nil {
		fmt.Println("Message Error: " + err.Error())
	}

	if strings.ToLower(price) == "free" {
		price = "0"
	}

	price = strings.Trim(price, "$")

	embed := message.Embeds[0]

	footerID := strings.Split(embed.Footer.Text, " | ")

	if strings.TrimPrefix(footerID[len(footerID)-1], "ID: ") != i.Member.User.ID {
		OutOfRange = true
		editContent = "Cannot edit another user's message."
	}

	if name != "" {
		embed.Title = fmt.Sprintf("Review for: %s", name)
	}

	if thumbnail != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: thumbnail,
		}
	}

	if rating != "" {
		floatRating, err := strconv.ParseFloat(rating, 64)
		if err != nil {
			fmt.Println(err)
			editContent = "There has been an error with the rating causing your review to fail. **Make sure you only use numbers for the rating**. Error Message: " + err.Error()
			OutOfRange = true
		}

		percentage := fmt.Sprintf("%d%%", int(floatRating*10))
		embed.Fields[0].Value = fmt.Sprintf("**%s/10** (%s)", rating, percentage)
	}
	if reason != "" {
		embed.Fields[1].Value = reason
	}
	if genre != "" {
		embed.Fields[2].Value = genre
	}
	if platform != "" {
		embed.Fields[3].Value = platform
	}

	if price != "" {
		FloatPrice, err = strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println(err)
			editContent = "There has been an error with the price causing your review to fail. **Make sure you only use numbers for the price**. Error Message: " + err.Error()
			OutOfRange = true
		}
		if price == "0" {
			price = "Free"
		} else {
			price = "$" + price + " USD"
		}
		embed.Fields[2].Value = price
	}

	switch {
	case math.Signbit(floatRating):
		editContent = fmt.Sprintf("If *Crypt of the Serpent King* can get 73%% I'm sure this game can get better than %s. **Please try again with a number between 0 and 10**", percentage)
		OutOfRange = true
	case math.Signbit(FloatPrice):
		editContent = "<@" + i.Member.User.ID + "> I know you think the game is trash but no one would pay *you* to play it. **Please try again with a positive cash value**"
		OutOfRange = true
	case floatRating == 0:
		fmt.Println(floatRating)
	case floatRating <= 1:
		embed.Color = 0xFF0000
	case floatRating <= 2:
		embed.Color = 0xFF4000
	case floatRating <= 3:
		embed.Color = 0xFF6200
	case floatRating <= 4:
		embed.Color = 0xFF7400
	case floatRating <= 5:
		embed.Color = 0xFFB100
	case floatRating <= 6:
		embed.Color = 0xFFAB00
	case floatRating <= 7:
		embed.Color = 0xFFCF00
	case floatRating <= 8:
		embed.Color = 0xFFFF00
	case floatRating <= 9:
		embed.Color = 0xCDFF00
	case floatRating < 10:
		embed.Color = 0x9CFF00
	case floatRating == 10:
		embed.Color = 0x00FF00
	case floatRating > 10:
		editContent = "<@" + i.Member.User.ID + "> Dude, it wasn't that good. **Please try again with a rating between 0 and 10**."
		fmt.Println(floatRating)
		OutOfRange = true
	}

	if FloatPrice > 100 {
		s.ChannelMessageSend(i.ChannelID, "Stop flexing your over $100 games! <@"+i.Member.User.ID+">")
	}

	if !OutOfRange {
		fmt.Println("Channel ID: " + channelID)
		fmt.Println("Message ID: " + messageID)
		_, err := s.ChannelMessageEditEmbed(channelID, messageID, embed)

		if err != nil {
			fmt.Println(err)
			fmt.Println("embed could not be sent")
		}

		editContent = "<@" + i.Member.User.ID + "> your review has been successfully submitted!"

	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: editContent,
			Flags:   64,
		},
	})

}
