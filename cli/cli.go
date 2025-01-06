package cli

import (
	"fmt"
	"os"
	"text/template"

	"github.com/paulja/go-work/cli/internal/adapters/grpc"
	"github.com/paulja/go-work/proto/scheduler/v1"
)

const (
	memberTmpl = "{{ range . }}{{.Id}}\t{{.Address}}\t{{.Status}}\n{{end}}"
	tasksTmpl  = "{{ range . }}{{.Id}}\t{{.Payload}}\t{{priority .Priority}}\t{{status .Status}}\n{{end}}"
)

func MembersCommand() error {
	cc := new(grpc.ClusterClient)
	if err := cc.Connect(); err != nil {
		return err
	}
	out, err := cc.GetMembers()
	if err != nil {
		return err
	}
	if len(out) == 0 {
		fmt.Println("no members found")
	} else {
		if err := printTemplate(memberTmpl, out); err != nil {
			return err
		}
	}
	return nil
}

func TasksCommand() error {
	sc := new(grpc.SchedulerClient)
	if err := sc.Connect(); err != nil {
		return err
	}
	out, err := sc.GetTasks()
	if err != nil {
		return err
	}
	if len(out) == 0 {
		fmt.Println("no tasks found")
	} else {
		if err := printTemplate(tasksTmpl, out); err != nil {
			return err
		}
	}
	return nil
}

func AddCommand(id, payload string) error {
	sc := new(grpc.SchedulerClient)
	if err := sc.Connect(); err != nil {
		return err
	}
	err := sc.AddTask(id, payload, scheduler.TaskPriority_TASK_PRIORITY_LOW)
	if err != nil {
		return err
	}
	fmt.Printf("success: %s\n", id)
	return nil
}

func RemoveCommand(id string) error {
	sc := new(grpc.SchedulerClient)
	if err := sc.Connect(); err != nil {
		return err
	}
	err := sc.RemoveTask(id)
	if err != nil {
		return err
	}
	fmt.Printf("success: %s\n", id)
	return nil
}

func printTemplate(tmpl string, data any) error {
	funcMap := template.FuncMap{
		"priority": grpc.ConvPriority,
		"status":   grpc.ConvStatus,
	}
	t, err := template.New("tmpl").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}
	err = t.Execute(os.Stdout, data)
	if err != nil {
		return err
	}
	return nil
}
