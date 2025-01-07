package grpc_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/paulja/go-work/proto/cluster/v1"
	"github.com/paulja/go-work/worker/config"
	grpcint "github.com/paulja/go-work/worker/internal/adapters/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHeartbeat(t *testing.T) {
	leader := NewLeaderMock(t)
	assert.NoError(t, leader.Start())

	t.Run("can start and stop heartbeat", func(t *testing.T) {
		leader.Reset()

		hb := grpcint.NewHeartbeat()
		assert.NoError(t, hb.Start(), "should be able to start heartbeat")
		assert.NoError(t, hb.Stop(), "should be able to stop heartbeat")
		assert.Equal(t, 1, leader.JoinCallCount(), "unexpected join call count")
		assert.Equal(t, 1, leader.LeaveCallCount(), "unexpected leave call count")
	})
	t.Run("can apply status", func(t *testing.T) {
		hb := grpcint.NewHeartbeat()
		assert.NoError(t, hb.Start(), "should be able to start heartbeat")
		hb.ApplyStatus(grpcint.HeartbeatStatusBusy)
		assert.NoError(t, hb.Stop(), "should be able to stop heartbeat")
	})
	t.Run("cannot apply invalid status", func(t *testing.T) {
		leader.Reset()

		os.Setenv("HEARTBEAT_TIMEOUT", "1")

		hb := grpcint.NewHeartbeat()
		assert.NoError(t, hb.Start(), "should be able to start heartbeat")

		hb.ApplyStatus(9)
		time.Sleep(1100 * time.Millisecond)
		assert.Equal(t, 0, leader.HeartbeatCallCount(), "unexpected heartbeat call count")
		assert.Equal(t, grpcint.HeartbeatStatusUnknown, leader.Status(), "unexpected status")
	})
	t.Run("heartbeat handler sends correct status", func(t *testing.T) {
		leader.Reset()

		os.Setenv("HEARTBEAT_TIMEOUT", "1")

		hb := grpcint.NewHeartbeat()
		assert.NoError(t, hb.Start(), "should be able to start heartbeat")

		hb.ApplyStatus(grpcint.HeartbeatStatusIdle)
		time.Sleep(1100 * time.Millisecond)
		assert.Equal(t, 1, leader.HeartbeatCallCount(), "unexpected heartbeat call count")
		assert.Equal(t, grpcint.HeartbeatStatusIdle, leader.Status(), "unexpected status")

		hb.ApplyStatus(grpcint.HeartbeatStatusBusy)
		time.Sleep(1100 * time.Millisecond)
		assert.Equal(t, 2, leader.HeartbeatCallCount(), "unexpected heartbeat call count")
		assert.Equal(t, grpcint.HeartbeatStatusBusy, leader.Status(), "unexpected heartbeat status")

		hb.ApplyStatus(grpcint.HeartbeatStatusFailed)
		time.Sleep(1100 * time.Millisecond)
		assert.Equal(t, 3, leader.HeartbeatCallCount(), "unexpected heartbeat call count")
		assert.Equal(t, grpcint.HeartbeatStatusFailed, leader.Status(), "unexpected status")

		assert.NoError(t, hb.Stop(), "should be able to stop heartbeat")
	})

	assert.NoError(t, leader.Stop())
}

/// -- MOCKS ---

var _ cluster.LeaderServiceServer = (*LeaderMock)(nil)

type LeaderMock struct {
	cluster.UnimplementedLeaderServiceServer

	sync.Mutex

	t    *testing.T
	conn net.Listener

	joinCallCount      int
	leaveCallCount     int
	heartbeatCallCount int
	status             grpcint.HeartbeatStatus
}

func NewLeaderMock(t *testing.T) *LeaderMock {
	t.Helper()
	return &LeaderMock{
		t: t,
	}
}

func (l *LeaderMock) JoinCallCount() int {
	l.Lock()
	defer l.Unlock()

	return l.joinCallCount
}

func (l *LeaderMock) LeaveCallCount() int {
	l.Lock()
	defer l.Unlock()

	return l.leaveCallCount
}

func (l *LeaderMock) HeartbeatCallCount() int {
	l.Lock()
	defer l.Unlock()

	return l.heartbeatCallCount
}

func (l *LeaderMock) Status() grpcint.HeartbeatStatus {
	l.Lock()
	defer l.Unlock()

	return l.status
}

func (l *LeaderMock) Start() error {
	listen, err := net.Listen("tcp", config.GetLeaderAddr())
	if err != nil {
		return fmt.Errorf("failed to listen: %s", config.GetLeaderAddr())
	}
	l.conn = listen
	grpcServer := grpc.NewServer()
	cluster.RegisterLeaderServiceServer(grpcServer, l)
	go func() {
		err = grpcServer.Serve(listen)
	}()
	if err != nil {
		l.t.Fatalf("failed to serve: %s", err)
	}
	return nil
}

func (l *LeaderMock) Stop() error {
	l.conn.Close()
	return nil
}

func (l *LeaderMock) Reset() {
	l.Lock()
	defer l.Unlock()

	l.joinCallCount = 0
	l.leaveCallCount = 0
	l.heartbeatCallCount = 0
	l.status = grpcint.HeartbeatStatusUnknown
}

func (l *LeaderMock) Join(
	ctx context.Context,
	req *cluster.JoinRequest,
) (
	*cluster.JoinResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id required")
	}
	if req.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, "address requied")
	}

	l.Lock()
	l.joinCallCount += 1
	l.Unlock()
	return &cluster.JoinResponse{}, nil
}

func (l *LeaderMock) Leave(
	ctx context.Context,
	req *cluster.LeaveRequest,
) (
	*cluster.LeaveResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id required")
	}
	l.Lock()
	l.leaveCallCount += 1
	l.Unlock()
	return &cluster.LeaveResponse{}, nil
}

func (l *LeaderMock) Heartbeat(
	ctx context.Context,
	req *cluster.HeartbeatRequest,
) (
	*cluster.HeartbeatResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id required")
	}

	var s grpcint.HeartbeatStatus
	switch req.Status {
	case cluster.HeartbeatStatus_HEARTBEAT_STATUS_UNSPECIFIED:
		return nil, status.Errorf(codes.InvalidArgument, "UNSPECIFIED is an invalid status")
	case cluster.HeartbeatStatus_HEARTBEAT_STATUS_IDLE:
		s = grpcint.HeartbeatStatusIdle
	case cluster.HeartbeatStatus_HEARTBEAT_STATUS_BUSY:
		s = grpcint.HeartbeatStatusBusy
	case cluster.HeartbeatStatus_HEARTBEAT_STATUS_FAILED:
		s = grpcint.HeartbeatStatusFailed
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid status")
	}

	l.Lock()
	l.heartbeatCallCount += 1
	l.status = s
	l.Unlock()

	return &cluster.HeartbeatResponse{}, nil
}
