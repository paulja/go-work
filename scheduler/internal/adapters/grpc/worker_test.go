package grpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/paulja/go-work/proto/worker/v1"
	grpcint "github.com/paulja/go-work/scheduler/internal/adapters/grpc"
	"github.com/paulja/go-work/worker/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWorker(t *testing.T) {
	s := new(WorkerServerMock)
	assert.NoError(t, s.Start())
	defer s.Stop()

	t.Run("can connect to worker client", func(t *testing.T) {
		s.Reset()

		w := grpcint.NewWorkerClient()
		assert.NoError(t, w.Connect(":40041"), "should be able to connect")
		assert.NoError(t, w.Close(), "should be able to close connection")
	})
	t.Run("can start work", func(t *testing.T) {
		s.Reset()

		w := grpcint.NewWorkerClient()
		assert.NoError(t, w.Connect(":40041"), "should be able to connect")
		assert.NoError(t, w.StartWork("id", "test"), "should be able start work")
		assert.Equal(t, 1, s.StartWorkCount, "unexpected value")
		assert.Error(t, w.StartWork("", ""), "expected an error")
		assert.Error(t, w.StartWork("1", ""), "expected an error")
		assert.NoError(t, w.Close(), "should be able to close connection")
	})
	t.Run("can stop work", func(t *testing.T) {
		s.Reset()

		w := grpcint.NewWorkerClient()
		assert.NoError(t, w.Connect(":40041"), "should be able to connect")
		assert.NoError(t, w.StopWork("id"), "should be able start work")
		assert.Equal(t, 1, s.StopWorkCount, "unexpected value")
		assert.Error(t, w.StopWork(""), "expected an error")
		assert.NoError(t, w.Close(), "should be able to close connection")
	})
}

/// -- MOCKS ---

var _ worker.WorkerServiceServer = (*WorkerServerMock)(nil)

type WorkerServerMock struct {
	worker.UnimplementedWorkerServiceServer

	conn net.Listener

	StartWorkCount int
	StopWorkCount  int
}

func (w *WorkerServerMock) Start() error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetWorkerPort()))
	if err != nil {
		return fmt.Errorf("failed to listen on port: %s", err)
	}
	w.conn = listen
	grpcServer := grpc.NewServer()
	worker.RegisterWorkerServiceServer(grpcServer, w)
	go func() {
		err = grpcServer.Serve(listen)
	}()
	if err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

func (w *WorkerServerMock) Stop() error {
	return w.conn.Close()
}

func (w *WorkerServerMock) Reset() {
	w.StartWorkCount = 0
	w.StopWorkCount = 0
}

func (w *WorkerServerMock) StartWork(
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

	w.StartWorkCount += 1

	return &worker.StartWorkResponse{
		Success: true,
	}, nil
}

func (w *WorkerServerMock) StopWork(
	ctx context.Context,
	req *worker.StopWorkRequest,
) (
	*worker.StopWorkResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Id required")
	}

	w.StopWorkCount += 1

	return &worker.StopWorkResponse{
		Success: true,
	}, nil
}
