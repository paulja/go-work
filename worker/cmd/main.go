package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/paulja/go-work/worker/internal/adapters/grpc"
)

func main() {
	hb := grpc.NewHeartbeat()
	hb.Start()

	notifyStream := make(chan os.Signal, 1)
	signal.Notify(notifyStream, syscall.SIGINT, syscall.SIGTERM)
	<-notifyStream
}
