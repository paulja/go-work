package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/paulja/go-work/worker/config"
	"github.com/paulja/go-work/worker/internal/adapters/grpc"
)

func main() {
	logger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)

	hb := grpc.NewHeartbeat()
	if err := hb.Start(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("worker running",
		"WORKER_NAME", config.GetName(),
		"Local Addr", config.GetLocalAddr(),
	)
	hb.ApplyStatus(grpc.HeartbeatStatusIdle)

	notifyStream := make(chan os.Signal, 1)
	signal.Notify(notifyStream, syscall.SIGINT, syscall.SIGTERM)
	<-notifyStream

	if err := hb.Stop(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
