package domain

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrTaskRequired        = errors.New("task required")
	ErrTaskIdRequired      = errors.New("task id required")
	ErrTaskPayloadRequired = errors.New("task payload required")
	ErrTaskNotFound        = errors.New("task not found")
)

type TaskStatus int

const (
	TaskStatusUnspecified TaskStatus = iota
	TaskStatusPending
	TaskStatusRunning
	TaskStatusCompleted
	TaskStatusCancelled
	TaskStatusError
)

type TaskPriority int

const (
	TaskPriorityUnspecified TaskPriority = iota
	TaskPriorityHigh
	TaskPriorityMedium
	TaskPriorityLow
)

type Task struct {
	Id       string
	Payload  string
	Priority TaskPriority

	Status TaskStatus
	Worker string
	Error  error

	Next *Task
}

func NewTask(id, payload string) *Task {
	return &Task{
		Id:      id,
		Payload: payload,
		Status:  TaskStatusUnspecified,
	}
}

type Tasks struct {
	sync.Mutex

	head *Task
	c    int
}

func (s *Tasks) Add(t *Task) {
	s.Lock()
	defer s.Unlock()

	t.Next = s.head
	s.head = t
	s.c += 1
}

func (s *Tasks) Stream() <-chan *Task {
	stream := make(chan *Task)
	node := s.head
	go func() {
		s.Lock()
		defer func() {
			s.Unlock()
			close(stream)
		}()

		for {
			if node == nil {
				return
			}
			stream <- node
			node = node.Next
		}
	}()
	return stream
}

func (s *Tasks) Id(id string) *Task {
	for t := range s.Stream() {
		if strings.Compare(t.Id, id) == 0 {
			return t
		}
	}
	return nil
}

func (s *Tasks) Status(status TaskStatus) []*Task {
	tasks := make([]*Task, 0, 8)
	for t := range s.Stream() {
		if t.Status == status {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

func (s *Tasks) ComparePriority(a, b *Task) int {
	if a.Priority == b.Priority {
		return 0
	}
	if a.Priority < b.Priority {
		return -1
	} else {
		return 1
	}
}

func (s *Tasks) Count() int {
	return s.c
}
