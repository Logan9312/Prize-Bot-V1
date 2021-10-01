module gitlab.com/logan9312/discord-auction-bot

go 1.17

require (
	github.com/bwmarrin/discordgo v0.23.3-0.20210821175000-0fad116c6c2a
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/gorilla/mux v1.8.0
	gorm.io/driver/postgres v1.1.1
	gorm.io/driver/sqlite v1.1.5
	gorm.io/gorm v1.21.15
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.1.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.8.1 // indirect
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.8 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/sys v0.0.0-20210910150752-751e447fb3d0 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/bwmarrin/discordgo => github.com/xIceArcher/discordgo v0.23.3-0.20210918115336-a3ed5a759b50
