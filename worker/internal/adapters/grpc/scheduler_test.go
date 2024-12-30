package grpc_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/paulja/go-work/proto/scheduler/v1"
	"github.com/paulja/go-work/worker/config"
	"github.com/paulja/go-work/worker/internal/adapters/grpc"
	"github.com/stretchr/testify/assert"
	ggrpc "google.golang.org/grpc"
)

func TestScheduler(t *testing.T) {
	t.Run("can create scheduler client", func(t *testing.T) {
		s := grpc.NewSchedulerClient()
		assert.NotNil(t, s, "should be able to create object")
	})
	t.Run("can connect and close client", func(t *testing.T) {
		s := grpc.NewSchedulerClient()
		assert.NoError(t, s.Connect(), "should be able to connect client")
		assert.NoError(t, s.Close(), "should be able to close client")
	})
	t.Run("can mark task complete", func(t *testing.T) {
		ss := new(ScheduleServerMock)
		ss.Start()
		defer ss.Stop()

		sc := grpc.NewSchedulerClient()
		assert.NoError(t, sc.Connect(), "should be able to connect client")
		defer sc.Close()
		assert.NoError(t, sc.TaskComplete("1", nil), "should be able to mark complete")
		assert.Equal(t, 1, ss.TaskCompleteCount)
	})
}

// --- Mocks ---\

var _ scheduler.SchedulerServiceServer = (*ScheduleServerMock)(nil)

type ScheduleServerMock struct {
	scheduler.UnimplementedSchedulerServiceServer
	conn net.Listener

	TaskCompleteCount int
}

func (s *ScheduleServerMock) Start() error {
	listen, err := net.Listen("tcp", config.GetSchedulerAddr())
	if err != nil {
		return fmt.Errorf("failed to listen on port: %s", err)
	}
	s.conn = listen
	grpcServer := ggrpc.NewServer()
	scheduler.RegisterSchedulerServiceServer(grpcServer, s)
	go func() {
		err = grpcServer.Serve(listen)
	}()
	if err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

func (s *ScheduleServerMock) Stop() error {
	return s.conn.Close()
}

func (s *ScheduleServerMock) Reset() {
	s.TaskCompleteCount = 0
}

func (s *ScheduleServerMock) TaskComplete(
	context.Context,
	*scheduler.TaskCompleteRequest,
) (
	*scheduler.TaskCompleteResponse,
	error,
) {
	s.TaskCompleteCount += 1
	return &scheduler.TaskCompleteResponse{}, nil
}
