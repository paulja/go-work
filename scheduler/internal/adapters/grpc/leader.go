package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/paulja/go-work/proto/cluster/v1"
	"github.com/paulja/go-work/scheduler/config"
	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/paulja/go-work/scheduler/internal/ports"
	"github.com/paulja/go-work/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var _ cluster.LeaderServiceServer = (*LeaderServer)(nil)

type LeaderServer struct {
	cluster.UnimplementedLeaderServiceServer

	logger *slog.Logger
	conn   net.Listener
	store  ports.MembershipPort
}

func NewLeaderServer(logger *slog.Logger, store ports.MembershipPort) *LeaderServer {
	return &LeaderServer{
		logger: logger,
		store:  store,
	}
}

func (l *LeaderServer) Start() error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetLeaderPort()))
	if err != nil {
		return fmt.Errorf("failed to listen on port: %s", err)
	}
	l.conn = listen
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc.UnaryServerInterceptor(shared.CreateLogInterceptor(*l.logger)),
		),
	)
	env := config.GetEnvironment()
	if env == "development" {
		reflection.Register(grpcServer)
	}
	cluster.RegisterLeaderServiceServer(grpcServer, l)
	go func() {
		err = grpcServer.Serve(listen)
	}()
	if err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

func (l *LeaderServer) Stop() error {
	return l.conn.Close()
}

func (l *LeaderServer) Join(
	ctx context.Context,
	req *cluster.JoinRequest,
) (
	*cluster.JoinResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, domain.ErrIdRequired.Error())
	}
	if req.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, domain.ErrAddressRequired.Error())
	}

	err := l.store.AddMember(domain.NewMember(req.Id, req.Address))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add member: %s", err)
	}
	return &cluster.JoinResponse{}, nil
}

func (l *LeaderServer) Leave(
	ctx context.Context,
	req *cluster.LeaveRequest,
) (
	*cluster.LeaveResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, domain.ErrIdRequired.Error())
	}

	err := l.store.RemoveMember(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove member: %s", err)
	}
	return &cluster.LeaveResponse{}, nil
}

func (l *LeaderServer) Members(
	context.Context,
	*cluster.MembersRequest,
) (
	*cluster.MembersResponse,
	error,
) {
	members, err := l.store.ListMembers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list members: %s", err)
	}
	resp := &cluster.MembersResponse{
		Members: make([]*cluster.Member, 0, len(members)),
	}
	for _, m := range members {
		resp.Members = append(resp.Members, &cluster.Member{
			Id:      m.Id,
			Address: m.Address,
			Status:  m.StatusString(),
		})
	}
	return resp, nil
}

func (l *LeaderServer) Heartbeat(stream cluster.LeaderService_HeartbeatServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return status.Error(codes.Unknown, err.Error())
		}
		if req.Id == "" {
			return status.Error(codes.InvalidArgument, domain.ErrIdRequired.Error())
		}

		var s domain.HeartbeatStatus
		switch req.Status {
		case cluster.HeartbeatStatus_HEARTBEAT_STATUS_UNSPECIFIED:
			return status.Errorf(codes.InvalidArgument, "UNSPECIFIED is an invalid status")
		case cluster.HeartbeatStatus_HEARTBEAT_STATUS_IDLE:
			s = domain.HeartbeatStatusIdle
		case cluster.HeartbeatStatus_HEARTBEAT_STATUS_BUSY:
			s = domain.HeartbeatStatusBusy
		case cluster.HeartbeatStatus_HEARTBEAT_STATUS_FAILED:
			s = domain.HeartbeatStatusFailed
		default:
			return status.Errorf(codes.InvalidArgument, "invalid status")
		}

		err = l.store.UpdateHeartbeatStatus(req.Id, s)
		if err != nil {
			switch err.Error() {
			case domain.ErrMemberNotFound.Error():
				return status.Errorf(codes.NotFound, "member not found")
			default:
				return status.Errorf(codes.Internal, "failed to update heartbeat status: %s", err)
			}
		}
		stream.Send(&cluster.HeartbeatResponse{})
	}
}
