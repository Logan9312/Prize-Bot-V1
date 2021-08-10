package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var SpawnExactDinoCommand = discordgo.ApplicationCommand{
	Name:        "spawn-exact-dino",
	Description: "A command for spawning a dino with setup stat points. (Local)",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "dinoblueprintpath",
			Description: "The blueprint path of the dino. Can be found in the ark wiki",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "basestats",
			Description: "Enter the stats found with the green gun. Set to 0 if none",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "addedstats",
			Description: "Points added after being tamed. Set to 0 if none",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "saddleblueprintpath",
			Description: "The blueprint path of the saddle. Can be found in the ark wiki",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "saddlequality",
			Description: "must be an integer",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "dinoname",
			Description: "The name of the dino. Can be anything",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "cloned",
			Description: "True if cloned, false if not",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "neutered",
			Description: "True if neutered, false if not",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "tamed-on",
			Description: "The location where the dino was tamed",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "uploaded-from",
			Description: "The location where the dino was uploaded",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "imprinter-name",
			Description: "If empty there will be no imprint",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "imprinter-id",
			Description: "The player ID of the user who imprinted",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "imprint-quality",
			Description: "The quality of the imptint. Multiply quality by 100 to get percentage",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "colors",
			Description: "the colors of the dino by region",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "dino-id",
			Description: "The ID of the dino",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "exp",
			Description: "experience gained by the dino",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "spawndistance",
			Description: "Distance the dino will spawn at",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "yoffset",
			Description: "The Y offset of the dino",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "zoffset",
			Description: "The Z offset of the dino",
			Required:    false,
		},
	},
}

func SpawnExactDino(s *discordgo.Session, i *discordgo.InteractionCreate) {

	saddleBP := ""
	saddleQuality := int64(0)
	dinoName := ""
	tamedOn := ""
	uploadedOn := ""
	imprinterName := ""
	imprinterID := int64(0)
	imprinterQuality := int64(0)
	colors := "0,0,0,0,0,0"
	dinoID := int64(0)
	exp := int64(0)
	spawnDistance := int64(0)
	yOffset := int64(20)
	zOffset := int64(20)
	blueprint := ""
	basestats := ""
	addedstats := ""
	isCloned := false
	isNeutered := false

	for _, v := range i.ApplicationCommandData().Options {
		switch v.Name {
		case "dinoblueprintpath":
			blueprint = v.StringValue()
		case "saddleblueprintpath":
			saddleBP = v.StringValue()
		case "basestats":
			basestats = v.StringValue()
		case "addedstats":
			addedstats = v.StringValue()
		case "saddlequality":
			saddleQuality = v.IntValue()
		case "dinoname":
			dinoName = v.StringValue()
		case "cloned":
			isCloned = v.BoolValue()
		case "neutered":
			isNeutered = v.BoolValue()
		case "tamed-on":
			tamedOn = v.StringValue()
		case "uploaded-from":
			uploadedOn = v.StringValue()
		case "imprinter-name":
			imprinterName = v.StringValue()
		case "imprinter-id":
			imprinterID = v.IntValue()
		case "imprinter-quality":
			imprinterQuality = v.IntValue()
		case "colors":
			colors = v.StringValue()
		case "dino-id":
			dinoID = v.IntValue()
		case "exp":
			exp = v.IntValue()
		case "spawndistance":
			spawnDistance = v.IntValue()
		case "yoffset":
			yOffset = v.IntValue()
		case "zoffset":
			zOffset = v.IntValue()
		}
	}

	cloned := 0
	if isCloned {
		cloned = 1
	}

	fmt.Println(isNeutered)

	neutered := 0
	if isNeutered {
		neutered = 1
	}

	splitBase := strings.Split(basestats, " ")
	newbase := []string{"0", "0", "0", "0", "0", "0", "0"}
	baseLevel := 0

	for num, v := range splitBase {
		points, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		newbase[num] = v
		baseLevel = baseLevel + points
	}
	baseLevel = baseLevel + 1
	basestats = strings.Join(newbase, ",")

	splitAdded := strings.Split(addedstats, " ")
	newadded := []string{"0", "0", "0", "0", "0", "0", "0"}
	extraLevels := 0

	for num, v := range splitAdded {
		points, err := strconv.Atoi(v)

		if err != nil {
			fmt.Println(err)
			return
		}

		newadded[num] = v
		extraLevels = extraLevels + points
	}
	addedstats = strings.Join(newadded, ",")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("cheat SpawnExactDino \"%s\" \"%s\" %d %d %d \"%s\" \"%s\" \"%s\" %d %d \"%s\" \"%s\" \"%s\" %d %d \"%s\" %d %d %d %d %d ", blueprint, saddleBP, saddleQuality, baseLevel, extraLevels, basestats, addedstats, dinoName, cloned, neutered, tamedOn, uploadedOn, imprinterName, imprinterID, imprinterQuality, colors, dinoID, exp, spawnDistance, yOffset, zOffset),
		},
	})
}
