package loglint

import (
	"log/slog"

	"go.uber.org/zap"
)

type Logger struct{}

func (l *Logger) Info(msg string)  {}
func (l *Logger) Error(msg string) {}
func (l *Logger) Warn(msg string)  {}
func (l *Logger) Debug(msg string) {}

func testSlog() {
	slog.Info("starting server")

	slog.Info("Starting server") // want "log message should start with lowercase letter"

	slog.Info("запуск сервера") // want "log message should contain only English text"

	slog.Info("server started!!!") // want "log message should not contain special characters or emoji"

	slog.Info("user password") // want "log message should not contain sensitive data keywords"

	slog.Info("api_key=" + "123") // want "log message should not contain special characters or emoji" "log message should not contain sensitive data keywords"
}

func testZap(logger *zap.Logger) {
	logger.Info("Starting zap logger") // want "log message should start with lowercase letter"

	logger.Error("user token") // want "log message should not contain sensitive data keywords"
}

func testCustomLogger(logger *Logger) {
	logger.Info("Starting logger") // want "log message should start with lowercase letter"
}