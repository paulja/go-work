package main

import (
	"log/slog"
	"os"

	"github.com/paulja/go-work/scheduler/config"
	"github.com/paulja/go-work/scheduler/internal/adapters/grpc"
	"github.com/paulja/go-work/scheduler/internal/adapters/membership"
)

func main() {
	logger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	s := membership.NewAdapter()
	l := grpc.NewLeaderServer(logger, s)
	if err := l.Start(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("scheduler running",
		"LEADER_PORT", config.GetLeaderPort(),
		"RPC_PORT", config.GetRPCPort(),
	)

	stop := make(chan interface{})
	<-stop
	os.Exit(0)
}
