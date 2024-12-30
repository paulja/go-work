package grpc_test

import "log/slog"

func createLogger() *slog.Logger {
	logger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)
	return logger
}
