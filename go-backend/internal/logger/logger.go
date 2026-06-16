package logger

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

// Init initializes the global Zap logger.
// Call this once at application startup.
// env should be "production" or "development".
func Init(env string) {
	once.Do(func() {
		var err error

		switch env {
		case "production":
			log, err = newProductionLogger()
		default:
			log, err = newDevelopmentLogger()
		}

		if err != nil {
			// Fallback to a no-op logger if initialization fails.
			log = zap.NewNop()
		}
	})
}

// Get returns the global logger instance.
// If Init has not been called, it initializes a development logger.
func Get() *zap.Logger {
	if log == nil {
		Init("development")
	}
	return log
}

// Sugar returns the global sugared logger for printf-style logging.
func Sugar() *zap.SugaredLogger {
	return Get().Sugar()
}

// Sync flushes any buffered log entries. Call this before application exit.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

func newDevelopmentLogger() (*zap.Logger, error) {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

func newProductionLogger() (*zap.Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}
