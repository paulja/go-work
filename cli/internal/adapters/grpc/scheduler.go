package grpc

import (
	"context"
	"time"

	"github.com/paulja/go-work/cli/config"
	"github.com/paulja/go-work/proto/scheduler/v1"
	"github.com/paulja/go-work/shared/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type SchedulerClient struct {
	conn   *grpc.ClientConn
	client scheduler.SchedulerServiceClient
}

func (s *SchedulerClient) Connect() error {
	cliTLS, err := tls.CliTLSConfig(config.GetServerName())
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(
		config.GetSchedulerAddr(),
		grpc.WithTransportCredentials(credentials.NewTLS(cliTLS)),
	)
	if err != nil {
		return err
	}
	s.conn = conn
	s.client = scheduler.NewSchedulerServiceClient(conn)
	return nil
}

func (c *SchedulerClient) Close() error {
	return c.conn.Close()
}

func (s *SchedulerClient) GetTasks() ([]*scheduler.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.client.GetTasks(ctx, &scheduler.GetTasksRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Tasks, nil
}

func (s *SchedulerClient) AddTask(id, payload string, priority scheduler.TaskPriority) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.client.ScheduleTask(ctx, &scheduler.ScheduleTaskRequest{
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.client.CancelTask(ctx, &scheduler.CancelTaskRequest{
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
