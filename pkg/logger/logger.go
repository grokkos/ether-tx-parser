package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var log *zap.Logger

// InitLogger initializes the logger with proper configuration
func InitLogger(env string) (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return log, nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		log, _ = InitLogger("development")
	}
	return log
}
