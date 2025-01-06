package main

import (
	"fmt"
	"os"

	"github.com/paulja/go-work/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR not enough arguments")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "members":
		fmt.Println("go-work", "members")
		fmt.Println()
		if err := cli.MembersCommand(); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	case "tasks":
		fmt.Println("go-work", "tasks")
		fmt.Println()
		if err := cli.TasksCommand(); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	case "add":
		fmt.Println("go-work", "add task")
		if len(os.Args) < 4 {
			fmt.Println("ERROR not enough arguments")
			os.Exit(2)
		}
		if err := cli.AddCommand(os.Args[2], os.Args[3]); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	case "remove":
		fmt.Println("go-work", "remove task")
		if len(os.Args) < 3 {
			fmt.Println("ERROR not enough arguments")
			os.Exit(2)
		}
		if err := cli.RemoveCommand(os.Args[2]); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	default:
		fmt.Println("ERROR command not recognised")
		os.Exit(2)
	}
}
