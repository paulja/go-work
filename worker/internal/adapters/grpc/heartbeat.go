package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/paulja/go-work/proto/cluster/v1"
	"github.com/paulja/go-work/worker/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HeartbeatStatus int

const (
	HeartbeatStatusUnknown HeartbeatStatus = iota
	HeartbeatStatusIdle
	HeartbeatStatusBusy
	HeartbeatStatusFailed
)

type Adapter struct {
	ctx    context.Context
	conn   *grpc.ClientConn
	client cluster.LeaderServiceClient
	stop   chan interface{}

	id     string
	addr   string
	status cluster.HeartbeatStatus
}

func NewHeartbeat() *Adapter {
	return &Adapter{
		ctx:  context.Background(),
		id:   config.GetName(),
		addr: config.GetLocalAddr(),
		stop: make(chan interface{}),
	}
}

func (a *Adapter) Start() error {
	conn, err := grpc.NewClient(
		config.GetLeaderAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to leader: %s", err)
	}
	a.conn = conn
	a.client = cluster.NewLeaderServiceClient(conn)
	_, err = a.client.Join(a.ctx, &cluster.JoinRequest{
		Id:      a.id,
		Address: a.addr,
	})
	if err != nil {
		return fmt.Errorf("failed to join with the leader: %s", err)
	}
	go a.heartbeatHandler()

	return nil
}

func (a *Adapter) Stop() error {
	close(a.stop)
	_, err := a.client.Leave(a.ctx, &cluster.LeaveRequest{
		Id: a.id,
	})
	if err != nil {
		return fmt.Errorf("failed to leave leader: %s", err)
	}
	return a.conn.Close()
}

func (a *Adapter) ApplyStatus(s HeartbeatStatus) {
	switch s {
	case HeartbeatStatusIdle:
		a.status = cluster.HeartbeatStatus_IDLE
	case HeartbeatStatusBusy:
		a.status = cluster.HeartbeatStatus_BUSY
	case HeartbeatStatusFailed:
		a.status = cluster.HeartbeatStatus_FAILED
	default:
		a.status = cluster.HeartbeatStatus_UNSPECIFIED
	}
}

func (a *Adapter) heartbeatHandler() {
	timeout := config.GetHeartbeatTimeout() * time.Second
	for {
		select {
		case <-time.Tick(timeout):
			a.client.Heartbeat(a.ctx, &cluster.HeartbeatRequest{
				Id:     a.id,
				Status: a.status,
			})
		case <-a.stop:
			return // stop the timer
		}
	}
}
