package domain_test

import (
	"slices"
	"testing"

	"github.com/paulja/go-work/scheduler/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestTasks(t *testing.T) {
	t.Run("can create tasks", func(t *testing.T) {
		tasks := new(domain.Tasks)
		assert.NotNil(t, tasks, "should be able to create tasks")
	})
	t.Run("can append tasks", func(t *testing.T) {
		tasks := new(domain.Tasks)
		assert.NotNil(t, tasks, "should be able to create tasks")

		tasks.Add(domain.NewTask("1", "testing"))
		tasks.Add(domain.NewTask("2", "testing"))
	})
	t.Run("can count tasks", func(t *testing.T) {
		tasks := new(domain.Tasks)
		assert.NotNil(t, tasks, "should be able to create tasks")

		tasks.Add(domain.NewTask("1", "testing"))
		tasks.Add(domain.NewTask("2", "testing"))
		tasks.Add(domain.NewTask("3", "testing"))
		tasks.Add(domain.NewTask("4", "testing"))

		assert.Equal(t, 4, tasks.Count(), "unexpected number of items")
	})
	t.Run("can stream tasks", func(t *testing.T) {
		tasks := new(domain.Tasks)
		assert.NotNil(t, tasks, "should be able to create tasks")

		tasks.Add(domain.NewTask("1", "testing"))
		tasks.Add(domain.NewTask("2", "testing"))
		tasks.Add(domain.NewTask("3", "testing"))
		tasks.Add(domain.NewTask("4", "testing"))

		count := 0
		for t := range tasks.Stream() {
			if t != nil {
				count += 1
			}
		}
		assert.Equal(t, 4, count, "unexpected number of items")
	})
	t.Run("can find tasks", func(t *testing.T) {
		tasks := new(domain.Tasks)
		assert.NotNil(t, tasks, "should be able to create tasks")

		tasks.Add(domain.NewTask("1", "testing1"))
		tasks.Add(domain.NewTask("2", "testing2"))
		tasks.Add(domain.NewTask("3", "testing3"))
		tasks.Add(domain.NewTask("4", "testing4"))

		task := tasks.Id("3")
		assert.NotNil(t, task, "should not be nil")
		assert.Equal(t, "testing3", task.Payload, "unexpected number of items")
	})
	t.Run("can group tasks", func(t *testing.T) {
		tasks := new(domain.Tasks)
		assert.NotNil(t, tasks, "should be able to create tasks")

		tasks.Add(&domain.Task{Id: "1", Payload: "testing1", Status: domain.TaskStatusPending})
		tasks.Add(&domain.Task{Id: "2", Payload: "testing2", Status: domain.TaskStatusRunning})
		tasks.Add(&domain.Task{Id: "3", Payload: "testing3", Status: domain.TaskStatusRunning})
		tasks.Add(&domain.Task{Id: "4", Payload: "testing4", Status: domain.TaskStatusPending})

		group := tasks.Status(domain.TaskStatusPending)
		assert.NotNil(t, group, "should not be nil")
		assert.Equal(t, 2, len(group), "unexpected number of items")
		assert.Equal(t, "testing4", group[0].Payload, "unexpected result")
		assert.Equal(t, "testing1", group[1].Payload, "unexpected result")
	})
	t.Run("can sort tasks", func(t *testing.T) {
		tasks := new(domain.Tasks)
		assert.NotNil(t, tasks, "should be able to create tasks")

		tasks.Add(&domain.Task{Id: "1", Payload: "testing1", Priority: domain.TaskPriorityLow})
		tasks.Add(&domain.Task{Id: "2", Payload: "testing2", Priority: domain.TaskPriorityHigh})
		tasks.Add(&domain.Task{Id: "3", Payload: "testing3", Priority: domain.TaskPriorityMedium})
		tasks.Add(&domain.Task{Id: "4", Payload: "testing4", Priority: domain.TaskPriorityHigh})

		list := make([]*domain.Task, 0, tasks.Count())
		for t := range tasks.Stream() {
			list = append(list, t)
		}
		slices.SortFunc(list, tasks.ComparePriority)

		assert.Equal(t, 4, len(list), "unexpected number of items")
		assert.Equal(t, "4", list[0].Id, "unexpected result")
		assert.Equal(t, "2", list[1].Id, "unexpected result")
		assert.Equal(t, "3", list[2].Id, "unexpected result")
		assert.Equal(t, "1", list[3].Id, "unexpected result")
	})
}
