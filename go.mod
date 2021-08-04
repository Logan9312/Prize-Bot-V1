module gitlab.com/logan9312/discord-auction-bot

go 1.16

require (
	github.com/bwmarrin/discordgo v0.23.3-0.20210725173931-083bf5c1d918
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/stretchr/testify v1.7.0 // indirect
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
)

replace github.com/bwmarrin/discordgo => github.com/FedorLap2006/discordgo v0.22.1-0.20210730172215-d44abc7f915e
