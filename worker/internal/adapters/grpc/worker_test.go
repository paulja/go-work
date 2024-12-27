package grpc_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/paulja/go-work/proto/worker/v1"
	"github.com/paulja/go-work/worker/internal/adapters/grpc"
	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	t.Run("can start and stop worker", func(t *testing.T) {
		logger := createLogger()
		w := grpc.NewWorkerServer(logger, updateHeartbeatStatus)
		assert.NoError(t, w.Start(), "should be able to start worker")
		assert.NoError(t, w.Stop(), "should be able to stop worker")
	})
	t.Run("can start work", func(t *testing.T) {
		logger := createLogger()
		w := grpc.NewWorkerServer(logger, updateHeartbeatStatus)
		assert.NoError(t, w.Start(), "should be able to start worker")
		resp, err := w.StartWork(context.Background(), &worker.StartWorkRequest{
			Id:      "WorkItem1",
			Payload: "do some work",
		})
		assert.NoError(t, err, "should be able to start work item")
		assert.True(t, resp.Success, "start work should be a success")
		assert.NoError(t, w.Stop(), "should be able to stop worker")
	})
	t.Run("cannot start work with invalid request", func(t *testing.T) {
		logger := createLogger()
		w := grpc.NewWorkerServer(logger, updateHeartbeatStatus)

		// no Id
		assert.NoError(t, w.Start(), "should be able to start worker")
		_, err := w.StartWork(context.Background(), &worker.StartWorkRequest{})
		assert.Error(t, err, "should get error with missing work id")
		assert.NoError(t, w.Stop(), "should be able to stop worker")

		// no payload
		assert.NoError(t, w.Start(), "should be able to start worker")
		_, err = w.StartWork(context.Background(), &worker.StartWorkRequest{
			Id: "WorkItem1",
		})
		assert.Error(t, err, "should get error with missing work payload")
		assert.NoError(t, w.Stop(), "should be able to stop worker")
	})
	t.Run("cannot start work if work is already started", func(t *testing.T) {
		logger := createLogger()
		w := grpc.NewWorkerServer(logger, updateHeartbeatStatus)
		assert.NoError(t, w.Start(), "should be able to start worker")
		resp, err := w.StartWork(context.Background(), &worker.StartWorkRequest{
			Id:      "WorkItem1",
			Payload: "do some work",
		})
		assert.NoError(t, err, "should be able to start work item")
		assert.True(t, resp.Success, "start work should be a success")
		_, err = w.StartWork(context.Background(), &worker.StartWorkRequest{
			Id:      "WorkItem1",
			Payload: "do some work",
		})
		assert.Error(t, err, "should not be able to start work if already working")
		assert.NoError(t, w.Stop(), "should be able to stop worker")
	})
	t.Run("doing work changes heartbeat status", func(t *testing.T) {
		logger := createLogger()
		w := grpc.NewWorkerServer(logger, updateHeartbeatStatus)
		assert.NoError(t, w.Start(), "should be able to start worker")
		resetHeartbeatStatus()
		assert.Equal(t, grpc.HeartbeatStatusIdle, heartbeatStatus, "heartbeat status unexpected")
		resp, err := w.StartWork(context.Background(), &worker.StartWorkRequest{
			Id:      "WorkItem1",
			Payload: "do some work",
		})
		assert.NoError(t, err, "should be able to start work item")
		assert.True(t, resp.Success, "start work should be a success")
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, grpc.HeartbeatStatusBusy, heartbeatStatus, "heartbeat status unexpected")
		assert.NoError(t, w.Stop(), "should be able to stop worker")
	})
	t.Run("cannot stop work with invalid request", func(t *testing.T) {
		logger := createLogger()
		w := grpc.NewWorkerServer(logger, updateHeartbeatStatus)

		assert.NoError(t, w.Start(), "should be able to start worker")
		_, err := w.StopWork(context.Background(), &worker.StopWorkRequest{})
		assert.Error(t, err, "should get error with missing work id")
		assert.NoError(t, w.Stop(), "should be able to stop worker")
	})
	t.Run("can stop work", func(t *testing.T) {
		logger := createLogger()
		w := grpc.NewWorkerServer(logger, updateHeartbeatStatus)
		assert.NoError(t, w.Start(), "should be able to start worker")
		resp1, err := w.StartWork(context.Background(), &worker.StartWorkRequest{
			Id:      "WorkItem1",
			Payload: "do some work",
		})
		assert.NoError(t, err, "should be able to start work item")
		assert.True(t, resp1.Success, "start work should be a success")
		resp2, err := w.StopWork(context.Background(), &worker.StopWorkRequest{
			Id: "WorkItem1",
		})
		assert.NoError(t, err, "should be able to start work item")
		assert.True(t, resp2.Success, "start stop should be a success")
		assert.NoError(t, w.Stop(), "should be able to stop worker")
	})
}

func createLogger() *slog.Logger {
	logger := slog.Default()
	slog.SetLogLoggerLevel(slog.LevelDebug)
	return logger
}

var heartbeatStatus grpc.HeartbeatStatus

func updateHeartbeatStatus(s grpc.HeartbeatStatus) {
	heartbeatStatus = s
}

func resetHeartbeatStatus() {
	heartbeatStatus = grpc.HeartbeatStatusIdle
}
