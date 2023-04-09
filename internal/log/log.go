package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newCore(debug bool) zapcore.Core {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "@timestamp", // Logstash standard timestamp field.
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // Logstash expects ISO 8601 times.
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	if debug {
		return zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), zapcore.DebugLevel)
	}
	return zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), zapcore.InfoLevel)
}

func new(debug bool) *zap.Logger {
	logger := zap.New(newCore(debug))
	logger = logger.WithOptions(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	logger.Sync()
	return logger
}

// New is exported
func New(loglevel string) *zap.Logger {
	if strings.ToLower(loglevel) == "debug" {
		return new(true)
	}
	return new(false)
}
