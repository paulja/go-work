package app_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/paulja/go-work/proto/worker/v1"
	"github.com/paulja/go-work/scheduler/internal/adapters/membership"
	"github.com/paulja/go-work/scheduler/internal/app"
	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/paulja/go-work/worker/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestScheduler(t *testing.T) {
	t.Run("can start and stop scheduler", func(t *testing.T) {
		logger := createLogger()
		store := membership.NewAdapter()
		s := app.NewTaskScheduler(logger, store)
		assert.NotNil(t, s, "should be able to create a scheduler")
		assert.NoError(t, s.Start(), "should be able to start scheduler")
		assert.NoError(t, s.Stop(), "should be able to stop scheduler")
	})
	t.Run("can schedule tasks", func(t *testing.T) {
		logger := createLogger()
		store := membership.NewAdapter()
		s := app.NewTaskScheduler(logger, store)
		assert.NotNil(t, s, "should be able to create a scheduler")

		assert.ErrorIs(t, domain.ErrTaskRequired, s.Schedule(nil), "unexpect error")

		task := domain.NewTask("1", "testing")
		assert.Equal(t, domain.TaskStatusUnspecified, task.Status, "unexpected state")
		assert.Equal(t, 0, len(s.List()), "unexpected list length")
		s.Schedule(task)
		assert.Equal(t, domain.TaskStatusPending, task.Status, "unexpected state")
		assert.Equal(t, 1, len(s.List()), "unexpected list length")
	})
	t.Run("can unschule tasks", func(t *testing.T) {
		logger := createLogger()
		store := membership.NewAdapter()
		s := app.NewTaskScheduler(logger, store)
		assert.NotNil(t, s, "should be able to create a scheduler")

		assert.ErrorIs(t, domain.ErrTaskNotFound, s.Unschedule("?"), "unexpect error")

		task := domain.NewTask("1", "testing")
		s.Schedule(task)
		assert.Equal(t, domain.TaskStatusPending, task.Status, "unexpected state")
		assert.NoError(t, s.Unschedule("1"), "should be able to unschedule tasks")
		assert.Equal(t, domain.TaskStatusCancelled, task.Status, "unexpected state")

		// TODO test when a worker has been engaged that the work gets cancelled
	})
	t.Run("can complete task", func(t *testing.T) {
		logger := createLogger()
		store := membership.NewAdapter()
		s := app.NewTaskScheduler(logger, store)
		assert.NotNil(t, s, "should be able to create a scheduler")

		assert.ErrorIs(t, domain.ErrTaskNotFound, s.Completed("?", nil), "unexpect error")

		task := domain.NewTask("1", "testing")
		s.Schedule(task)
		s.Completed("1", nil)
		assert.Equal(t, domain.TaskStatusCompleted, task.Status, "unexpected state")

		s.Completed("1", fmt.Errorf("test error"))
		assert.Equal(t, domain.TaskStatusError, task.Status, "unexpected state")
		assert.NotNil(t, task.Error, "expected an error")
	})
	t.Run("can apply", func(t *testing.T) {
		w := new(WorkerServerMock)
		w.Start()
		defer w.Stop()

		logger := createLogger()
		store := membership.NewAdapter()
		m := domain.NewMember("1", ":40041")
		m.SetHeartbeatStatus(domain.HeartbeatStatusIdle)
		m.SetMembershipStatus(domain.MembershipStatusAlive)
		store.AddMember(m)

		s := app.NewTaskScheduler(logger, store)
		assert.NotNil(t, s, "should be able to create a scheduler")

		task := domain.NewTask("1", "testing")
		s.Schedule(task)

		os.Setenv("POLL_INTERVAL", "1")
		assert.NoError(t, s.Start(), "should be able to start scheduler")
		time.Sleep(1200 * time.Millisecond)
		assert.NoError(t, s.Stop(), "should be able to stop scheduler")
		os.Unsetenv("POLL_INTERVAL")

		assert.Equal(t, domain.TaskStatusRunning, task.Status, "unexpected state")

		// TODO test the correct worker name is set for the task
	})
}

func createLogger() *slog.Logger {
	logger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)
	return logger
}

// -- Mocks --

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

func (w *WorkerServerMock) StartWork(
	context.Context,
	*worker.StartWorkRequest,
) (
	*worker.StartWorkResponse,
	error,
) {
	return &worker.StartWorkResponse{
		Success: true,
	}, nil
}

func (w *WorkerServerMock) StopWork(
	context.Context,
	*worker.StopWorkRequest,
) (
	*worker.StopWorkResponse,
	error,
) {
	return &worker.StopWorkResponse{
		Success: true,
	}, nil
}
