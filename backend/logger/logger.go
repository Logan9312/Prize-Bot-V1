package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Sugar is the global sugared logger for convenient logging
	Sugar *zap.SugaredLogger
	// Log is the underlying zap logger for structured logging
	Log *zap.Logger
)

// Init initializes the logger based on environment and log level
func Init(environment, logLevel string) error {
	var config zap.Config

	if environment == "prod" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set log level from environment variable
	level := parseLogLevel(logLevel)
	config.Level = zap.NewAtomicLevelAt(level)

	// Add caller information to logs
	var err error
	Log, err = config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	Sugar = Log.Sugar()
	return nil
}

// parseLogLevel converts a string log level to zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		// Default to info in prod, debug in dev
		if os.Getenv("ENVIRONMENT") == "prod" {
			return zapcore.InfoLevel
		}
		return zapcore.DebugLevel
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// Command returns a logger with command context fields
func Command(commandName, guildID, userID, username string) *zap.SugaredLogger {
	return Sugar.With(
		"command", commandName,
		"guild_id", guildID,
		"user_id", userID,
		"username", username,
	)
}

// Button returns a logger with button interaction context fields
func Button(buttonID, guildID, userID, username string) *zap.SugaredLogger {
	return Sugar.With(
		"button", buttonID,
		"guild_id", guildID,
		"user_id", userID,
		"username", username,
	)
}

// Auction returns a logger with auction context fields
func Auction(channelID, guildID, item string) *zap.SugaredLogger {
	return Sugar.With(
		"event_type", "auction",
		"channel_id", channelID,
		"guild_id", guildID,
		"item", item,
	)
}

// Giveaway returns a logger with giveaway context fields
func Giveaway(messageID, guildID, item string) *zap.SugaredLogger {
	return Sugar.With(
		"event_type", "giveaway",
		"message_id", messageID,
		"guild_id", guildID,
		"item", item,
	)
}

// Database returns a logger with database operation context
func Database(operation string) *zap.SugaredLogger {
	return Sugar.With(
		"component", "database",
		"operation", operation,
	)
}

// Bot returns a logger with bot instance context
func Bot(botID, botName string) *zap.SugaredLogger {
	return Sugar.With(
		"bot_id", botID,
		"bot_name", botName,
	)
}

// Timer returns a logger for timer/scheduler operations
func Timer(eventType, guildID string) *zap.SugaredLogger {
	return Sugar.With(
		"component", "timer",
		"event_type", eventType,
		"guild_id", guildID,
	)
}
