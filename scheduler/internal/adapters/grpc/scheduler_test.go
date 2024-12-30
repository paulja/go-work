package grpc_test

import (
	"context"
	"testing"

	"github.com/paulja/go-work/proto/scheduler/v1"
	"github.com/paulja/go-work/scheduler/internal/adapters/grpc"
	"github.com/paulja/go-work/scheduler/internal/adapters/membership"
	"github.com/paulja/go-work/scheduler/internal/app"
	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	t.Run("can start and stop service", func(t *testing.T) {
		s, _ := createScheduleServer()
		assert.NoError(t, s.Start(), "should be able to start scheduler")
		assert.NoError(t, s.Stop(), "should be able to stop scheduler")
	})
	t.Run("can schedule task", func(t *testing.T) {
		s, a := createScheduleServer()
		assert.NoError(t, s.Start(), "should be able to start scheduler")

		_, err := s.ScheduleTask(context.Background(), &scheduler.ScheduleTaskRequest{})
		assert.Error(t, err, "unexpected error")
		_, err = s.ScheduleTask(context.Background(), &scheduler.ScheduleTaskRequest{
			Task: &scheduler.Task{},
		})
		assert.Error(t, err, "unexpected error")
		_, err = s.ScheduleTask(context.Background(), &scheduler.ScheduleTaskRequest{
			Task: &scheduler.Task{Id: "2"},
		})
		assert.Error(t, err, "unexpected error")

		_, err = s.ScheduleTask(context.Background(), &scheduler.ScheduleTaskRequest{
			Task: &scheduler.Task{
				Id:       "1",
				Payload:  "testing",
				Priority: scheduler.TaskPriority_TASK_PRIORITY_MEDIUM.Enum(),
			},
		})
		assert.NoError(t, err, "unexpected error")
		assert.NoError(t, s.Stop(), "should be able to stop scheduler")
		assert.Equal(t, 1, len(a.List()), "item should be queued")
	})
	t.Run("can cancel task", func(t *testing.T) {
		s, a := createScheduleServer()
		assert.NoError(t, s.Start(), "should be able to start scheduler")
		s.ScheduleTask(context.Background(), &scheduler.ScheduleTaskRequest{
			Task: &scheduler.Task{
				Id:       "1",
				Payload:  "testing",
				Priority: scheduler.TaskPriority_TASK_PRIORITY_MEDIUM.Enum(),
			},
		})

		_, err := s.CancelTask(context.Background(), &scheduler.CancelTaskRequest{})
		assert.Error(t, err, "unexpected error")

		_, err = s.CancelTask(context.Background(), &scheduler.CancelTaskRequest{
			Id: "1",
		})
		assert.NoError(t, err, "unexpected error")
		task := a.List()[0]
		assert.Equal(t, domain.TaskStatusCancelled, task.Status, "unexpected status")
		assert.NoError(t, s.Stop(), "should be able to stop scheduler")
	})
	t.Run("can get tasks", func(t *testing.T) {
		s, _ := createScheduleServer()
		assert.NoError(t, s.Start(), "should be able to start scheduler")
		s.ScheduleTask(context.Background(), &scheduler.ScheduleTaskRequest{
			Task: &scheduler.Task{Id: "1", Payload: "testing",
				Priority: scheduler.TaskPriority_TASK_PRIORITY_HIGH.Enum()},
		})
		s.ScheduleTask(context.Background(), &scheduler.ScheduleTaskRequest{
			Task: &scheduler.Task{Id: "2", Payload: "testing",
				Priority: scheduler.TaskPriority_TASK_PRIORITY_LOW.Enum()},
		})
		s.ScheduleTask(context.Background(), &scheduler.ScheduleTaskRequest{
			Task: &scheduler.Task{Id: "3", Payload: "testing"},
		})
		resp, err := s.GetTasks(context.Background(), &scheduler.GetTasksRequest{})
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, 3, len(resp.Tasks), "item should be queued")
		assert.NoError(t, s.Stop(), "should be able to stop scheduler")
	})

	// TODO add tests for TaskComplete
}

func createScheduleServer() (*grpc.ScheduleServer, *app.TaskScheduler) {
	logger := createLogger()
	store := membership.NewAdapter()
	scheduler := app.NewTaskScheduler(logger, store)
	return grpc.NewScheduleServer(logger, scheduler), scheduler
}
