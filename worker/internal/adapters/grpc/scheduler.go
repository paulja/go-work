package grpc

import (
	"context"
	"time"

	"github.com/paulja/go-work/proto/scheduler/v1"
	"github.com/paulja/go-work/shared/tls"
	"github.com/paulja/go-work/worker/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	workerTLS, err := tls.WorkerTLSConfig(config.GetServerName())
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(
		config.GetSchedulerAddr(),
		grpc.WithTransportCredentials(credentials.NewTLS(workerTLS)),
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
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()

	_, err := s.client.TaskComplete(ctx, &scheduler.TaskCompleteRequest{
		Id:    id,
		Error: errorMessage,
	})
	if err != nil {
		return err
	}
	return nil
}
