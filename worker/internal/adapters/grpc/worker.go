package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/paulja/go-work/proto/worker/v1"
	"github.com/paulja/go-work/shared"
	"github.com/paulja/go-work/worker/config"
	"github.com/paulja/go-work/worker/internal/app"
	"github.com/paulja/go-work/worker/internal/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var _ worker.WorkerServiceServer = (*WorkerServer)(nil)

type WorkerServer struct {
	worker.UnimplementedWorkerServiceServer

	work          *app.Worker
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
	workerTLS, err := tls.WorkerServerTLSConfig(config.GetServerName())
	if err != nil {
		return fmt.Errorf("failed to server TLS: %s", err)
	}
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(workerTLS)),
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

	if w.work != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Already working")
	}

	work := new(app.Worker)
	go func() {
		w.heartbeatFunc(HeartbeatStatusBusy)
		err := work.Start(req.Id, req.Payload)

		ss := NewSchedulerClient()
		ss.Connect()
		ss.TaskComplete(req.Id, errorString(err))
		ss.Close()

		w.heartbeatFunc(HeartbeatStatusIdle)
		w.work = nil
	}()
	w.work = work

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

	w.heartbeatFunc(HeartbeatStatusIdle)
	if w.work != nil {
		w.work.Stop(req.Id)
		w.work = nil
	}

	return &worker.StopWorkResponse{
		Success: true,
	}, nil
}

func errorString(e error) *string {
	if e != nil {
		s := e.Error()
		return &s
	}
	return nil
}
