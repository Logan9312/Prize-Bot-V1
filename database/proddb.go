package database

import (
	"os"

	"gitlab.com/logan9312/discord-auction-bot/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ProdDB() *gorm.DB {
	log := logger.Database("prod_connect")
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Error("DATABASE_URL environment variable is not set")
		return nil
	}

	// Only show first 20 chars for security
	if len(databaseURL) > 20 {
		log.Infow("connecting to postgres database", "url_prefix", databaseURL[:20]+"...")
	} else {
		log.Info("connecting to postgres database")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.NewGormLogger(),
	})
	if err != nil {
		log.Errorw("failed to connect to postgres database", "error", err)
		return nil
	}

	log.Info("postgres database connection successful")
	return db
}
