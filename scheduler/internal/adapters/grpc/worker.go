package grpc

import (
	"context"
	"errors"

	"github.com/paulja/go-work/proto/worker/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
	resp, err := w.client.StartWork(w.ctx, &worker.StartWorkRequest{
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
	resp, err := w.client.StopWork(w.ctx, &worker.StopWorkRequest{
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
