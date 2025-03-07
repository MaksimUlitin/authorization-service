package logger

import (
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func Info(msg string, args ...interface{}) {
	Logger.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	Logger.Error(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	Logger.Warn(msg, args...)
}
