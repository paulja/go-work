package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/paulja/go-work/proto/worker/v1"
	"github.com/paulja/go-work/scheduler/config"
	"github.com/paulja/go-work/scheduler/internal/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	ErrFailedToStartWork = errors.New("failed to start work")
	ErrFailedToStopWork  = errors.New("failed to stop work")
)

type WorkerClient struct {
	ctx    context.Context
	conn   *grpc.ClientConn
	client worker.WorkerServiceClient
}

func NewWorkerClient() *WorkerClient {
	return &WorkerClient{
		ctx: context.Background(),
	}
}

func (w *WorkerClient) Connect(addr string) error {
	workerTLS, err := tls.WorkerTLSConfig(config.GetServerName())
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(credentials.NewTLS(workerTLS)),
	)
	if err != nil {
		return err
	}
	w.conn = conn
	w.client = worker.NewWorkerServiceClient(conn)
	return nil
}

func (w *WorkerClient) Close() error {
	return w.conn.Close()
}

func (w *WorkerClient) StartWork(id, payload string) error {
	ctx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
	defer cancel()

	resp, err := w.client.StartWork(ctx, &worker.StartWorkRequest{
		Id:      id,
		Payload: payload,
	})
	if err != nil {
		return errors.Join(ErrFailedToStartWork, err)
	}
	if !resp.Success {
		return ErrFailedToStartWork
	}
	return err
}

func (w *WorkerClient) StopWork(id string) error {
	ctx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
	defer cancel()

	resp, err := w.client.StopWork(ctx, &worker.StopWorkRequest{
		Id: id,
	})
	if err != nil {
		return errors.Join(ErrFailedToStopWork, err)
	}
	if !resp.Success {
		return ErrFailedToStopWork
	}
	return nil
}
