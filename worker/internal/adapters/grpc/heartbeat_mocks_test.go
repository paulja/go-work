package grpc_test

// LeaderServerService mock implementation go support testing in the Worker module.

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/paulja/go-work/proto/cluster/v1"
	"github.com/paulja/go-work/worker/config"
	grpcint "github.com/paulja/go-work/worker/internal/adapters/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ cluster.LeaderServiceServer = (*LeaderMock)(nil)

type LeaderMock struct {
	cluster.UnimplementedLeaderServiceServer

	t    *testing.T
	conn net.Listener

	JoinCallCount      int
	LeaveCallCount     int
	HeartbeatCallCount int
	Status             grpcint.HeartbeatStatus
}

func NewLeaderMock(t *testing.T) *LeaderMock {
	t.Helper()
	return &LeaderMock{
		t: t,
	}
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
	l.JoinCallCount = 0
	l.LeaveCallCount = 0
	l.HeartbeatCallCount = 0
	l.Status = grpcint.HeartbeatStatusUnknown
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
	l.JoinCallCount += 1
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
	l.LeaveCallCount += 1
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
	case cluster.HeartbeatStatus_UNSPECIFIED:
		return nil, status.Errorf(codes.InvalidArgument, "UNSPECIFIED is an invalid status")
	case cluster.HeartbeatStatus_IDLE:
		s = grpcint.HeartbeatStatusIdle
	case cluster.HeartbeatStatus_BUSY:
		s = grpcint.HeartbeatStatusBusy
	case cluster.HeartbeatStatus_FAILED:
		s = grpcint.HeartbeatStatusFailed
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid status")
	}

	l.HeartbeatCallCount += 1
	l.Status = s

	return &cluster.HeartbeatResponse{}, nil
}
