package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(env string) *Logger {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		Logger: logger,
	}
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	l.Logger.Fatal(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.Logger.Sync()
}
