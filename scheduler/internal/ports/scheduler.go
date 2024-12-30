package ports

import "github.com/paulja/go-work/scheduler/internal/domain"

type Scheduler interface {
	Schedule(t *domain.Task) error
	Unschedule(id string) error
	List() []*domain.Task

	Completed(id string, err error) error
}
