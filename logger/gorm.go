package logger

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger is a custom GORM logger that uses Zap
type GormLogger struct {
	ZapLogger                 *zap.SugaredLogger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// NewGormLogger creates a new GORM logger using Zap
func NewGormLogger() *GormLogger {
	return &GormLogger{
		ZapLogger:                 Sugar.With("component", "gorm"),
		LogLevel:                  gormlogger.Warn,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

// LogMode sets the log level and returns a new logger
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info logs info level messages
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Infof(msg, args...)
	}
}

// Warn logs warn level messages
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.ZapLogger.Warnf(msg, args...)
	}
}

// Error logs error level messages
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.ZapLogger.Errorf(msg, args...)
	}
}

// Trace logs SQL queries with timing information
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error:
		if !(errors.Is(err, gorm.ErrRecordNotFound) && l.IgnoreRecordNotFoundError) {
			l.ZapLogger.Errorw("database error",
				"error", err.Error(),
				"elapsed_ms", float64(elapsed.Nanoseconds())/1e6,
				"rows", rows,
				"sql", sql,
			)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		l.ZapLogger.Warnw("slow query detected",
			"elapsed_ms", float64(elapsed.Nanoseconds())/1e6,
			"threshold_ms", float64(l.SlowThreshold.Nanoseconds())/1e6,
			"rows", rows,
			"sql", sql,
		)
	case l.LogLevel >= gormlogger.Info:
		l.ZapLogger.Debugw("query executed",
			"elapsed_ms", float64(elapsed.Nanoseconds())/1e6,
			"rows", rows,
			"sql", sql,
		)
	}
}
