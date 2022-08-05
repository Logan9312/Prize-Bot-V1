module gitlab.com/logan9312/discord-auction-bot

go 1.18

require (
	github.com/bwmarrin/discordgo v0.25.1-0.20220804185119-c0803d021f34
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/stripe/stripe-go/v72 v72.122.0
	golang.org/x/text v0.3.7
	gorm.io/driver/postgres v1.3.8
	gorm.io/driver/sqlite v1.3.6
	gorm.io/gorm v1.23.8
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.12.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	github.com/jackc/pgx/v4 v4.16.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.0.0-20220805013720-a33c5aa5df48 // indirect
	golang.org/x/sys v0.0.0-20220804214406-8e32c043e418 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/bwmarrin/discordgo => github.com/FedorLap2006/discordgo v0.22.1-0.20220210004113-86c1f8e9af7a
