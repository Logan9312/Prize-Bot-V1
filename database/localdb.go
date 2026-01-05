package database

import (
	"gitlab.com/logan9312/discord-auction-bot/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func LocalDB() *gorm.DB {
	log := logger.Database("local_connect")
	log.Info("connecting to local SQLite database")

	db, err := gorm.Open(sqlite.Open("/tmp/test.db"), &gorm.Config{
		Logger: logger.NewGormLogger(),
	})
	if err != nil {
		log.Errorw("failed to connect to SQLite database", "error", err)
		return nil
	}

	log.Info("SQLite database connection successful")
	return db
}
