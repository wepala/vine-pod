package database

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	applogger "github.com/wepala/vine-pod/pkg/logger"
)

// GormLogger adapts our Zap logger to work with GORM
type GormLogger struct {
	logger               applogger.Logger
	logLevel             logger.LogLevel
	ignoreRecordNotFound bool
	slowThreshold        time.Duration
}

// NewGormLogger creates a new GORM logger adapter
func NewGormLogger(log applogger.Logger) logger.Interface {
	return &GormLogger{
		logger:               log,
		logLevel:             logger.Info,
		ignoreRecordNotFound: true,
		slowThreshold:        200 * time.Millisecond,
	}
}

// LogMode sets the log level
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info logs info level messages
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		l.logger.Info("GORM: "+msg, zap.Any("data", data))
	}
}

// Warn logs warning level messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		l.logger.Warn("GORM: "+msg, zap.Any("data", data))
	}
}

// Error logs error level messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		l.logger.Error("GORM: "+msg, zap.Any("data", data))
	}
}

// Trace logs SQL statements
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rowsAffected := fc()

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.String("sql", sql),
		zap.Int64("rows_affected", rowsAffected),
	}

	switch {
	case err != nil && l.logLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFound):
		l.logger.Error("GORM SQL Error", append(fields, zap.Error(err))...)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		l.logger.Warn("GORM Slow SQL", append(fields, zap.Duration("slow_threshold", l.slowThreshold))...)
	case l.logLevel == logger.Info:
		l.logger.Debug("GORM SQL", fields...)
	}
}
