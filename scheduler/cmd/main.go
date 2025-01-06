package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/paulja/go-work/scheduler/config"
	"github.com/paulja/go-work/scheduler/internal/adapters/grpc"
	"github.com/paulja/go-work/scheduler/internal/adapters/membership"
	"github.com/paulja/go-work/scheduler/internal/app"
)

func main() {
	logger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	m := membership.NewAdapter()
	l := grpc.NewLeaderServer(logger, m)
	if err := l.Start(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	t := app.NewTaskScheduler(logger, m)
	if err := t.Start(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	d := grpc.NewScheduleServer(logger, t)
	if err := d.Start(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("scheduler running",
		"LEADER_PORT", config.GetLeaderPort(),
		"RPC_PORT", config.GetRPCPort(),
	)

	notifyStream := make(chan os.Signal, 1)
	signal.Notify(notifyStream, syscall.SIGINT, syscall.SIGTERM)
	<-notifyStream

	if err := l.Stop(); err != nil {
		slog.Error("problem stopping leader", "err", err.Error())
		os.Exit(1)
	}
	if err := t.Stop(); err != nil {
		slog.Error("problem stopping task scheduler", "err", err.Error())
		os.Exit(1)
	}
	if err := d.Stop(); err != nil {
		slog.Error("problem stopping scheduler", "err", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
