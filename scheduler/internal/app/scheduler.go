package app

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/paulja/go-work/scheduler/config"
	"github.com/paulja/go-work/scheduler/internal/adapters/grpc"
	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/paulja/go-work/scheduler/internal/ports"
)

var (
	ErrNoWorkers = errors.New("no workers available")
)

type TaskScheduler struct {
	sync.Mutex
	logger *slog.Logger

	store  ports.MembershipPort
	tasks  *domain.Tasks
	cancel chan interface{}
}

var _ ports.Scheduler = (*TaskScheduler)(nil)

func NewTaskScheduler(logger *slog.Logger, store ports.MembershipPort) *TaskScheduler {
	return &TaskScheduler{
		logger: logger,
		store:  store,
		tasks:  new(domain.Tasks),
		cancel: make(chan interface{}),
	}
}

func (s *TaskScheduler) Start() error {
	go func() {
		for {
			select {
			case <-time.After(config.GetPollInterval() * time.Second):
				err := s.scheduleWork()
				if err != nil {
					s.logger.Error(err.Error())
				}
			case <-s.cancel:
				return
			}
		}
	}()
	return nil
}

func (s *TaskScheduler) Stop() error {
	close(s.cancel)
	return nil
}

func (s *TaskScheduler) Schedule(t *domain.Task) error {
	if t == nil {
		return domain.ErrTaskRequired
	}
	t.Status = domain.TaskStatusPending
	s.tasks.Add(t)
	return nil
}

func (s *TaskScheduler) Unschedule(id string) error {
	t := s.tasks.Id(id)
	if t == nil {
		return domain.ErrTaskNotFound
	}

	switch t.Status {
	case domain.TaskStatusPending:
		t.Status = domain.TaskStatusCancelled
	case domain.TaskStatusRunning:
		t.Status = domain.TaskStatusCancelled
		if len(t.Worker) > 0 {
			if err := s.unscheduleWork(t.Worker, id); err != nil {
				return err
			}
		}
	default:
		break
	}
	return nil
}

func (s *TaskScheduler) Completed(id string, err error) error {
	t := s.tasks.Id(id)
	if t == nil {
		return domain.ErrTaskNotFound
	}
	if err != nil {
		t.Status = domain.TaskStatusError
		t.Error = err
	} else {
		t.Status = domain.TaskStatusCompleted
	}
	return nil
}

func (s *TaskScheduler) List() []*domain.Task {
	ctx := context.Background()

	tasks := make([]*domain.Task, 0, s.tasks.Count())
	for t := range s.tasks.Stream(ctx) {
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *TaskScheduler) unscheduleWork(name, id string) error {
	worker := s.findMember(name)
	if worker != nil {
		return workerStopWork(worker.Address, id)
	}
	return nil
}

func (s *TaskScheduler) scheduleWork() error {
	s.Lock()
	defer s.Unlock()

	tasks := s.tasks.Status(domain.TaskStatusPending)
	if len(tasks) == 0 {
		return nil
	}
	slices.SortFunc(tasks, s.tasks.ComparePriority)
	for _, t := range tasks {
		worker, ok := s.firstAvailbleMember()
		if !ok {
			return ErrNoWorkers
		}
		err := workerStartWork(worker.Address, t)
		if err != nil {
			return err
		}
		t.Status = domain.TaskStatusRunning
		t.Worker = worker.Id
	}
	return nil
}

func workerStartWork(address string, t *domain.Task) error {
	wc := grpc.NewWorkerClient()
	if err := wc.Connect(address); err != nil {
		return err
	}
	defer wc.Close()
	if err := wc.StartWork(t.Id, t.Payload); err != nil {
		return err
	}
	return nil
}

func workerStopWork(address, id string) error {
	wc := grpc.NewWorkerClient()
	if err := wc.Connect(address); err != nil {
		return err
	}
	defer wc.Close()
	if err := wc.StopWork(id); err != nil {
		return err
	}
	return nil
}

func (s *TaskScheduler) firstAvailbleMember() (*domain.Member, bool) {
	list, _ := s.store.ListMembers()
	for _, v := range list {
		if v.MembershipStatus() == domain.MembershipStatusAlive &&
			v.HeartbeatStatus() == domain.HeartbeatStatusIdle {
			return v, true
		}
	}
	return nil, false
}

func (s *TaskScheduler) findMember(id string) *domain.Member {
	list, _ := s.store.ListMembers()
	for _, v := range list {
		if strings.Compare(id, v.Id) == 0 {
			return v
		}
	}
	return nil
}
