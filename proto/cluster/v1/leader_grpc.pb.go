// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: cluster/v1/leader.proto

package cluster

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	LeaderService_Join_FullMethodName      = "/cluster.v1.LeaderService/Join"
	LeaderService_Leave_FullMethodName     = "/cluster.v1.LeaderService/Leave"
	LeaderService_Members_FullMethodName   = "/cluster.v1.LeaderService/Members"
	LeaderService_Heartbeat_FullMethodName = "/cluster.v1.LeaderService/Heartbeat"
)

// LeaderServiceClient is the client API for LeaderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LeaderServiceClient interface {
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error)
	Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveResponse, error)
	Members(ctx context.Context, in *MembersRequest, opts ...grpc.CallOption) (*MembersResponse, error)
	Heartbeat(ctx context.Context, opts ...grpc.CallOption) (LeaderService_HeartbeatClient, error)
}

type leaderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLeaderServiceClient(cc grpc.ClientConnInterface) LeaderServiceClient {
	return &leaderServiceClient{cc}
}

func (c *leaderServiceClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(JoinResponse)
	err := c.cc.Invoke(ctx, LeaderService_Join_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderServiceClient) Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LeaveResponse)
	err := c.cc.Invoke(ctx, LeaderService_Leave_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderServiceClient) Members(ctx context.Context, in *MembersRequest, opts ...grpc.CallOption) (*MembersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MembersResponse)
	err := c.cc.Invoke(ctx, LeaderService_Members_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderServiceClient) Heartbeat(ctx context.Context, opts ...grpc.CallOption) (LeaderService_HeartbeatClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &LeaderService_ServiceDesc.Streams[0], LeaderService_Heartbeat_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &leaderServiceHeartbeatClient{ClientStream: stream}
	return x, nil
}

type LeaderService_HeartbeatClient interface {
	Send(*HeartbeatRequest) error
	Recv() (*HeartbeatResponse, error)
	grpc.ClientStream
}

type leaderServiceHeartbeatClient struct {
	grpc.ClientStream
}

func (x *leaderServiceHeartbeatClient) Send(m *HeartbeatRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *leaderServiceHeartbeatClient) Recv() (*HeartbeatResponse, error) {
	m := new(HeartbeatResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LeaderServiceServer is the server API for LeaderService service.
// All implementations must embed UnimplementedLeaderServiceServer
// for forward compatibility
type LeaderServiceServer interface {
	Join(context.Context, *JoinRequest) (*JoinResponse, error)
	Leave(context.Context, *LeaveRequest) (*LeaveResponse, error)
	Members(context.Context, *MembersRequest) (*MembersResponse, error)
	Heartbeat(LeaderService_HeartbeatServer) error
	mustEmbedUnimplementedLeaderServiceServer()
}

// UnimplementedLeaderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLeaderServiceServer struct {
}

func (UnimplementedLeaderServiceServer) Join(context.Context, *JoinRequest) (*JoinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedLeaderServiceServer) Leave(context.Context, *LeaveRequest) (*LeaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedLeaderServiceServer) Members(context.Context, *MembersRequest) (*MembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Members not implemented")
}
func (UnimplementedLeaderServiceServer) Heartbeat(LeaderService_HeartbeatServer) error {
	return status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedLeaderServiceServer) mustEmbedUnimplementedLeaderServiceServer() {}

// UnsafeLeaderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LeaderServiceServer will
// result in compilation errors.
type UnsafeLeaderServiceServer interface {
	mustEmbedUnimplementedLeaderServiceServer()
}

func RegisterLeaderServiceServer(s grpc.ServiceRegistrar, srv LeaderServiceServer) {
	s.RegisterService(&LeaderService_ServiceDesc, srv)
}

func _LeaderService_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderServiceServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LeaderService_Join_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderServiceServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderService_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderServiceServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LeaderService_Leave_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderServiceServer).Leave(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderService_Members_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderServiceServer).Members(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LeaderService_Members_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderServiceServer).Members(ctx, req.(*MembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderService_Heartbeat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LeaderServiceServer).Heartbeat(&leaderServiceHeartbeatServer{ServerStream: stream})
}

type LeaderService_HeartbeatServer interface {
	Send(*HeartbeatResponse) error
	Recv() (*HeartbeatRequest, error)
	grpc.ServerStream
}

type leaderServiceHeartbeatServer struct {
	grpc.ServerStream
}

func (x *leaderServiceHeartbeatServer) Send(m *HeartbeatResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *leaderServiceHeartbeatServer) Recv() (*HeartbeatRequest, error) {
	m := new(HeartbeatRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LeaderService_ServiceDesc is the grpc.ServiceDesc for LeaderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LeaderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cluster.v1.LeaderService",
	HandlerType: (*LeaderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _LeaderService_Join_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _LeaderService_Leave_Handler,
		},
		{
			MethodName: "Members",
			Handler:    _LeaderService_Members_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Heartbeat",
			Handler:       _LeaderService_Heartbeat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "cluster/v1/leader.proto",
}
