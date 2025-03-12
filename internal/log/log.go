package log

import (
	"log/slog"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
)

// Aliases
var (
	Error = zap.Error
)

type Logger struct {
	*zap.Logger
}

func New() *Logger {
	return &Logger{
		zap.NewExample(),
	}
}

func (l *Logger) Zap() *zap.Logger {
	return l.Logger
}

func (l *Logger) Slog() *slog.Logger {
	level := slog.LevelInfo
	for s, z := range slogzap.LogLevels {
		if z == l.Level() {
			level = s
			break
		}
	}

	return slog.New(slogzap.Option{
		Level: level, Logger: l.Zap(),
	}.NewZapHandler())
}
