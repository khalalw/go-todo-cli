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

	// Initialize the todo list
	todoList := &todo.Todos{}
	// Load todos with error handling
	if err := todoList.Load(fileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *add != "":
		todoList.Add(*add)
		fmt.Printf("Added task: %s, length = %d\n", *add, len(*todoList))
	case *complete >= 0:
		if err := todoList.Complete(*complete - 1); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Marked task %d as complete\n", *complete)
	case *del >= 0:
		if err := todoList.Delete(*del - 1); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("Deleted task %d\n", *del)
	case *list:
		todo.Print(todoList)
	default:
		fmt.Fprintln(os.Stderr, "Invalid command")
		flag.Usage()
		os.Exit(1)
	}

	// Save changes after adding
	if err := todoList.Save(fileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
