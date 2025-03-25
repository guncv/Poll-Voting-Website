package util

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerInstance *Logger
	once           sync.Once
)

// LoggerInterface defines the methods for custom logger
type LoggerInterface interface {
	ErrorWithID(ctx context.Context, args ...interface{})
	DebugWithID(ctx context.Context, args ...interface{})
	InfoWithID(ctx context.Context, args ...interface{})
}

// Logger wraps zap.SugaredLogger and implements LoggerInterface
type Logger struct {
	*zap.SugaredLogger
}

// Initialize sets up the logger based on the application environment
func Initialize(appEnv string) *Logger {
	once.Do(func() {
		var err error
		var config zap.Config

		if appEnv == "dev" {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder // Use human-readable time
			config.EncoderConfig.TimeKey = ""                            // Include time in logs
			config.EncoderConfig.CallerKey = "caller"                    // Include caller information
			config.EncoderConfig.MessageKey = "msg"                      // Message field key
			config.EncoderConfig.LevelKey = "level"                      // Level field key
			config.EncoderConfig.ConsoleSeparator = " | "                // Separator for better readability
		} else {
			config = zap.NewProductionConfig()
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		}

		baseLogger, err := config.Build(zap.AddCaller()) // Base logger includes caller
		if err != nil {
			panic("failed to initialize zap logger: " + err.Error())
		}

		// Convert zap.Logger to SugaredLogger for easier logging
		loggerInstance = &Logger{baseLogger.Sugar()}
	})

	return loggerInstance
}

// Sync flushes any buffered log entries
func Sync() {
	if loggerInstance != nil {
		_ = loggerInstance.Sync()
	}
}

// GetLogger returns the singleton logger instance
func GetLogger() *Logger {
	if loggerInstance == nil {
		panic("logger is not initialized. Call Initialize() first.")
	}
	return loggerInstance
}

// ErrorWithID logs an error with custom context information, adjusting caller skip
func (l *Logger) ErrorWithID(ctx context.Context, args ...interface{}) {
	// Create a new logger instance with caller skip set to 1 to point to the handler
	loggerWithSkip := l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()
	fields := l.buildLogFields(ctx)
	loggerWithSkip.With(fields...).Error(args...)
}

// DebugWithID logs a debug message with custom context information, adjusting caller skip
func (l *Logger) DebugWithID(ctx context.Context, args ...interface{}) {
	// Create a new logger instance with caller skip set to 1 to point to the handler
	loggerWithSkip := l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()
	fields := l.buildLogFields(ctx)
	loggerWithSkip.With(fields...).Debug(args...)
}

// InfoWithID logs an info message with custom context information, adjusting caller skip
func (l *Logger) InfoWithID(ctx context.Context, args ...interface{}) {
	// Create a new logger instance with caller skip set to 1 to point to the handler
	loggerWithSkip := l.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar()
	fields := l.buildLogFields(ctx)
	loggerWithSkip.With(fields...).Info(args...)
}

// buildLogFields extracts and builds log fields from context and log type
func (l *Logger) buildLogFields(ctx context.Context) []interface{} {
	// Example implementation: extract information from context keys
	userID := ctx.Value("userID")
	if userID == nil {
		userID = "unknown"
	}

	return []interface{}{
		"userID", userID,
		"timestamp", time.Now().Format(time.RFC3339),
	}
}
