package logger

import (
	"errors"
	"io"
	"log/slog"
	"strings"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/interfaces"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type SLogger struct {
	logger *slog.Logger
	Level  Level
}

func New(level Level, writer io.Writer) appinterfaces.Logger {
	return &SLogger{
		Level: level,
		logger: slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: slog.Level((level - 1) * 4),
		})),
	}
}

func (l *SLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *SLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *SLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *SLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

var (
	LevelMap = map[string]Level{
		"debug":   Debug,
		"info":    Info,
		"warning": Warn,
		"error":   Error,
	}
	ErrUnknownLevel = errors.New("unknown level")
)

func GetLevelOrPanic(level string) Level {
	l, found := LevelMap[strings.ToLower(level)]
	if !found {
		panic(ErrUnknownLevel)
	}
	return l
}
