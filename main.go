package main

import (
	"flag"
	"fmt"
	"go-todo-cli/todo"
	"os"
)

const fileToWrite = "todos.json"

func main() {
	// -add, -complete, -delete, -list
	add := flag.String("add", "", "Task to add to the TODO list")
	complete := flag.Int("complete", -1, "Task number to mark as complete")
	del := flag.Int("delete", -1, "Task number to delete")
	list := flag.Bool("list", false, "List all tasks")
	flag.Parse()

	// initialize todolist
	todos := &todo.Todos{}
	// load todos with error handling
	if err := todos.Load(fileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *add != "":
		todos.Add(*add)
		fmt.Printf("Added task %s\n, length = %d", *add, len(*todos))
	case *complete >= 0:
		if err := todos.Complete(*complete - 1); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Marked task %d as complete", *complete)
	case *del >= 0:
		if err := todos.Delete(*del - 1); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Deleted task %d\n", *del)
	case *list:
		todo.Print(todos)
	default:
		fmt.Fprintln(os.Stderr, "Invalid command")
		flag.Usage()
		os.Exit(1)
	}

	// save changes after adding
	if err := todos.Save(fileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
