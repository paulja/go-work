package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net"
	"time"

	"github.com/paulja/go-work/proto/worker/v1"
	"github.com/paulja/go-work/shared"
	"github.com/paulja/go-work/worker/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var _ worker.WorkerServiceServer = (*WorkerServer)(nil)

type WorkerServer struct {
	worker.UnimplementedWorkerServiceServer

	cancel        chan interface{}
	logger        *slog.Logger
	heartbeatFunc func(HeartbeatStatus)
	conn          net.Listener
}

func NewWorkerServer(logger *slog.Logger, hb func(HeartbeatStatus)) *WorkerServer {
	return &WorkerServer{
		logger:        logger,
		heartbeatFunc: hb,
	}
}

func (w *WorkerServer) Start() error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetWorkerPort()))
	if err != nil {
		return fmt.Errorf("failed to listen on port: %s", err)
	}
	w.conn = listen
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc.UnaryServerInterceptor(shared.CreateLogInterceptor(*w.logger)),
		),
	)
	env := config.GetEnvironment()
	if env == "development" {
		reflection.Register(grpcServer)
	}
	worker.RegisterWorkerServiceServer(grpcServer, w)
	go func() {
		err = grpcServer.Serve(listen)
	}()
	if err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

func (w *WorkerServer) Stop() error {
	return w.conn.Close()
}

func (w *WorkerServer) StartWork(
	ctx context.Context,
	req *worker.StartWorkRequest,
) (
	*worker.StartWorkResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Id required")
	}

	if req.Payload == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Payload required")
	}

	if w.cancel != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Already working")
	}

	w.heartbeatFunc(HeartbeatStatusBusy)

	w.cancel = make(chan interface{})
	defer close(w.cancel)
	min, max := uint(30), uint(60)
	num := rand.UintN(max-min) + min

	w.logger.Debug("working ", "duration", fmt.Sprintf("%ds", num))
	select {
	case <-time.After(time.Duration(num) * time.Second):
		w.logger.Debug("finished")
		w.cancel = nil
	case <-w.cancel:
		w.logger.Debug("cancelled")
		w.cancel = nil
	}

	w.heartbeatFunc(HeartbeatStatusIdle)

	return &worker.StartWorkResponse{
		Success: true,
	}, nil
}

func (w *WorkerServer) StopWork(
	ctx context.Context,
	req *worker.StopWorkRequest,
) (
	*worker.StopWorkResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Id required")
	}

	if w.cancel != nil {
		close(w.cancel)
	}

	return &worker.StopWorkResponse{
		Success: true,
	}, nil
}
