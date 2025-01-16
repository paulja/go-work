package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/paulja/go-work/proto/scheduler/v1"
	"github.com/paulja/go-work/scheduler/config"
	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/paulja/go-work/scheduler/internal/ports"
	"github.com/paulja/go-work/shared"
	"github.com/paulja/go-work/shared/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var _ scheduler.SchedulerServiceServer = (*ScheduleServer)(nil)

type ScheduleServer struct {
	scheduler.UnimplementedSchedulerServiceServer

	logger    *slog.Logger
	scheduler ports.Scheduler
	conn      net.Listener
}

func NewScheduleServer(logger *slog.Logger, scheduler ports.Scheduler) *ScheduleServer {
	return &ScheduleServer{
		logger:    logger,
		scheduler: scheduler,
	}
}

func (s *ScheduleServer) Start() error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetRPCPort()))
	if err != nil {
		return fmt.Errorf("failed to listen on port: %s", err)
	}
	s.conn = listen
	schedulerTLS, err := tls.SchedulerTLSConfig(config.GetServerName())
	if err != nil {
		return fmt.Errorf("failed to server TLS: %s", err)
	}
	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(schedulerTLS)),
		grpc.UnaryInterceptor(
			grpc.UnaryServerInterceptor(shared.CreateLogInterceptor(*s.logger)),
		),
	)
	if config.GetEnvironment() == "development" {
		reflection.Register(grpcServer)
	}
	scheduler.RegisterSchedulerServiceServer(grpcServer, s)
	go func() {
		err = grpcServer.Serve(listen)
	}()
	if err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
	return nil
}

func (s *ScheduleServer) Stop() error {
	return s.conn.Close()
}

func (s *ScheduleServer) ScheduleTask(
	ctx context.Context,
	req *scheduler.ScheduleTaskRequest,
) (
	*scheduler.ScheduleTaskResponse,
	error,
) {
	if req.Task == nil {
		return nil, status.Errorf(codes.InvalidArgument, domain.ErrTaskRequired.Error())
	}
	if req.Task.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, domain.ErrTaskIdRequired.Error())
	}
	if req.Task.Payload == "" {
		return nil, status.Errorf(codes.InvalidArgument, domain.ErrTaskPayloadRequired.Error())
	}

	task := domain.NewTask(req.Task.Id, req.Task.Payload)
	task.Priority = CastInPriority(req.Task.Priority)
	if err := s.scheduler.Schedule(task); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &scheduler.ScheduleTaskResponse{}, nil
}

func (s *ScheduleServer) CancelTask(
	ctx context.Context,
	req *scheduler.CancelTaskRequest,
) (
	*scheduler.CancelTaskResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, domain.ErrTaskIdRequired.Error())
	}

	if err := s.scheduler.Unschedule(req.Id); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &scheduler.CancelTaskResponse{}, nil
}

func (s *ScheduleServer) GetTasks(
	ctx context.Context,
	req *scheduler.GetTasksRequest,
) (
	*scheduler.GetTasksResponse,
	error,
) {
	tasks := make([]*scheduler.Task, 0, 8)
	for _, t := range s.scheduler.List() {
		tasks = append(tasks, CastTask(t))
	}

	return &scheduler.GetTasksResponse{
		Tasks: tasks,
	}, nil
}

func (s *ScheduleServer) TaskComplete(
	ctx context.Context,
	req *scheduler.TaskCompleteRequest,
) (
	*scheduler.TaskCompleteResponse,
	error,
) {
	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, domain.ErrTaskIdRequired.Error())
	}
	var taskErr error
	if req.Error != nil {
		taskErr = fmt.Errorf(*req.Error)
	}
	if err := s.scheduler.Completed(req.Id, taskErr); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &scheduler.TaskCompleteResponse{}, nil
}

func CastInPriority(in *scheduler.TaskPriority) domain.TaskPriority {
	if in == nil {
		return domain.TaskPriorityLow
	}

	switch *in {
	case scheduler.TaskPriority_TASK_PRIORITY_HIGH:
		return domain.TaskPriorityHigh
	case scheduler.TaskPriority_TASK_PRIORITY_MEDIUM:
		return domain.TaskPriorityMedium
	case scheduler.TaskPriority_TASK_PRIORITY_LOW:
		return domain.TaskPriorityLow
	default:
		return domain.TaskPriorityLow
	}
}

func CastOutPriority(in domain.TaskPriority) *scheduler.TaskPriority {
	switch in {
	case domain.TaskPriorityHigh:
		return scheduler.TaskPriority_TASK_PRIORITY_HIGH.Enum()
	case domain.TaskPriorityMedium:
		return scheduler.TaskPriority_TASK_PRIORITY_MEDIUM.Enum()
	case domain.TaskPriorityLow:
		return scheduler.TaskPriority_TASK_PRIORITY_LOW.Enum()
	default:
		return scheduler.TaskPriority_TASK_PRIORITY_LOW.Enum()
	}
}

func CastOutStatus(in domain.TaskStatus) *scheduler.TaskStatus {
	switch in {
	case domain.TaskStatusCancelled:
		return scheduler.TaskStatus_TASK_STATUS_CANCELLED.Enum()
	case domain.TaskStatusCompleted:
		return scheduler.TaskStatus_TASK_STATUS_COMPLETED.Enum()
	case domain.TaskStatusError:
		return scheduler.TaskStatus_TASK_STATUS_ERROR.Enum()
	case domain.TaskStatusPending:
		return scheduler.TaskStatus_TASK_STATUS_PENDING.Enum()
	case domain.TaskStatusRunning:
		return scheduler.TaskStatus_TASK_STATUS_RUNNING.Enum()
	default:
		return scheduler.TaskStatus_TASK_STATUS_UNSPECIFIED.Enum()
	}
}

func CastTask(in *domain.Task) *scheduler.Task {
	return &scheduler.Task{
		Id:       in.Id,
		Payload:  in.Payload,
		Priority: CastOutPriority(in.Priority),
		Status:   CastOutStatus(in.Status),
	}
}
