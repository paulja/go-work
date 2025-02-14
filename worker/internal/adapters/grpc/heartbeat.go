package grpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/paulja/go-work/proto/cluster/v1"
	"github.com/paulja/go-work/worker/config"
	"github.com/paulja/go-work/worker/internal/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type HeartbeatStatus int

const (
	HeartbeatStatusUnknown HeartbeatStatus = iota
	HeartbeatStatusIdle
	HeartbeatStatusBusy
	HeartbeatStatusFailed
)

type HeartbeatAdapter struct {
	sync.Mutex

	ctx    context.Context
	conn   *grpc.ClientConn
	client cluster.LeaderServiceClient
	stop   chan interface{}

	id     string
	addr   string
	status cluster.HeartbeatStatus
}

func NewHeartbeat() *HeartbeatAdapter {
	return &HeartbeatAdapter{
		ctx:  context.Background(),
		id:   config.GetName(),
		addr: config.GetLocalAddr(),
		stop: make(chan interface{}),
	}
}

func (a *HeartbeatAdapter) Start() error {
	workerTLS, err := tls.WorkerTLSConfig(config.GetServerName())
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(
		config.GetLeaderAddr(),
		grpc.WithTransportCredentials(credentials.NewTLS(workerTLS)),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to leader: %s", err)
	}
	a.conn = conn
	a.client = cluster.NewLeaderServiceClient(conn)
	_, err = a.client.Join(a.ctx, &cluster.JoinRequest{
		Id:      a.id,
		Address: config.GetAddr(),
	})
	if err != nil {
		return fmt.Errorf("failed to join with the leader: %s", err)
	}
	go a.heartbeatHandler()

	return nil
}

func (a *HeartbeatAdapter) Stop() error {
	close(a.stop)
	_, err := a.client.Leave(a.ctx, &cluster.LeaveRequest{
		Id: a.id,
	})
	if err != nil {
		return fmt.Errorf("failed to leave leader: %s", err)
	}
	return a.conn.Close()
}

func (a *HeartbeatAdapter) ApplyStatus(s HeartbeatStatus) {
	a.Lock()
	defer a.Unlock()

	switch s {
	case HeartbeatStatusIdle:
		a.status = cluster.HeartbeatStatus_HEARTBEAT_STATUS_IDLE
	case HeartbeatStatusBusy:
		a.status = cluster.HeartbeatStatus_HEARTBEAT_STATUS_BUSY
	case HeartbeatStatusFailed:
		a.status = cluster.HeartbeatStatus_HEARTBEAT_STATUS_FAILED
	default:
		a.status = cluster.HeartbeatStatus_HEARTBEAT_STATUS_UNSPECIFIED
	}
}

func (a *HeartbeatAdapter) heartbeatHandler() {
	timeout := config.GetHeartbeatTimeout() * time.Second
	for {
		select {
		case <-time.Tick(timeout):
			ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
			a.Lock()
			a.client.Heartbeat(ctx, &cluster.HeartbeatRequest{
				Id:     a.id,
				Status: a.status,
			})
			a.Unlock()
			cancel()
			// TODO: handle failures and rejoin/send as needed
		case <-a.stop:
			return // stop the timer
		}
	}
}
