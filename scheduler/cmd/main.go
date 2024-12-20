package main

import (
	"github.com/paulja/go-work/scheduler/internal/adapters/grpc"
	"github.com/paulja/go-work/scheduler/internal/adapters/membership"
)

func main() {
	s := membership.NewAdapter()
	l := grpc.NewLeaderServer(s)
	l.Start()
}
