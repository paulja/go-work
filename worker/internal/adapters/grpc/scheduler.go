package grpc

import (
	"context"

	"github.com/paulja/go-work/proto/scheduler/v1"
	"github.com/paulja/go-work/worker/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SchedulerClient struct {
	ctx    context.Context
	conn   *grpc.ClientConn
	client scheduler.SchedulerServiceClient
}

func NewSchedulerClient() *SchedulerClient {
	return &SchedulerClient{
		ctx: context.Background(),
	}
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

func (s *SchedulerClient) Close() error {
	return s.conn.Close()
}

func (s *SchedulerClient) TaskComplete(id string, errorMessage *string) error {
	_, err := s.client.TaskComplete(s.ctx, &scheduler.TaskCompleteRequest{
		Id:    id,
		Error: errorMessage,
	})
	if err != nil {
		return err
	}
	return nil
}
