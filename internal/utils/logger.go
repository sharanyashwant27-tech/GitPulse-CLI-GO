package utils

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Init configures the global structured logger.
func Init(level string, development bool) error {
	var err error
	once.Do(func() {
		var cfg zap.Config
		if development {
			cfg = zap.NewDevelopmentConfig()
		} else {
			cfg = zap.NewProductionConfig()
			cfg.Encoding = "json"
		}

		switch level {
		case "debug":
			cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		case "warn":
			cfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		case "error":
			cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		default:
			cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		}

		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		log, err = cfg.Build(zap.AddCallerSkip(1))
		if err != nil {
			return
		}
	})
	return err
}

// L returns the global logger, initializing a no-op logger if needed.
func L() *zap.Logger {
	if log == nil {
		_ = Init("info", false)
		if log == nil {
			log = zap.NewNop()
		}
	}
	return log
}

// Sync flushes buffered log entries.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// Fatal logs at fatal level and exits.
func Fatal(msg string, fields ...zap.Field) {
	L().Fatal(msg, fields...)
	os.Exit(1)
}
