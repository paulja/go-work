package grpc

import (
	"context"

	"github.com/paulja/go-work/cli/config"
	"github.com/paulja/go-work/proto/scheduler/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SchedulerClient struct {
	conn   *grpc.ClientConn
	client scheduler.SchedulerServiceClient
}

func (s *SchedulerClient) Connect() error {
	conn, err := grpc.NewClient(
		config.GetSchedulerAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	s.conn = conn
	s.client = scheduler.NewSchedulerServiceClient(conn)
	return nil
}

func (s *SchedulerClient) GetTasks() ([]*scheduler.Task, error) {
	resp, err := s.client.GetTasks(context.Background(), &scheduler.GetTasksRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Tasks, nil
}

func (s *SchedulerClient) AddTask(id, payload string, priority scheduler.TaskPriority) error {
	_, err := s.client.ScheduleTask(context.Background(), &scheduler.ScheduleTaskRequest{
		Task: &scheduler.Task{
			Id:       id,
			Payload:  payload,
			Priority: &priority,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *SchedulerClient) RemoveTask(id string) error {
	_, err := s.client.CancelTask(context.Background(), &scheduler.CancelTaskRequest{
		Id: id,
	})
	if err != nil {
		return err
	}
	return nil
}

func ConvPriority(v scheduler.TaskPriority) string {
	return v.String()
}

func ConvStatus(v scheduler.TaskStatus) string {
	return v.String()
}
