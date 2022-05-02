package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const timeFormat = "2006-01-02 15:04:05"

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func NewZapLogger(level string) (*ZapLogger, error) {
	zLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("unknown log level %s: %w", level, err)
	}

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		formatted := t.Format(timeFormat)
		encoder.AppendString(formatted)
	}
	config.Level = zLevel
	config.Encoding = "console"
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.OutputPaths = []string{"stderr"}
	config.ErrorOutputPaths = []string{"stderr"}

	zl, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("build ZapLogger: %w", err)
	}

	return &ZapLogger{zl}, nil
}

type ZapLogger struct {
	logger *zap.Logger
}

func (l *ZapLogger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *ZapLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *ZapLogger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *ZapLogger) Error(msg string) {
	l.logger.Error(msg)
}
