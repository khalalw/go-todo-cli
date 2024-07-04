package main

import (
	"flag"
	"fmt"
	"go-todo-cli/todo"
	"os"
)

const fileToWrite = "todos.json"

func main() {
	// Parse command-line flags
	add, complete, del, list := parseFlags()

	// Initialize an empty todo list
	todoList := &todo.Todos{}

	// Load existing todos from file or create a new file if it doesn't exist
	handleFileLoading(todoList, fileToWrite)

	// Execute the appropriate command based on the flags provided
	executeCommand(add, complete, del, list, todoList)

	// Save the updated todo list to the file
	if err := todoList.Save(fileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// parseFlags parses command-line flags and returns pointers to the flag values
func parseFlags() (add *string, complete *int, del *int, list *bool) {
	add = flag.String("add", "", "Task to add to the TODO list")
	complete = flag.Int("complete", -1, "Task number to mark as complete")
	del = flag.Int("delete", -1, "Task number to delete")
	list = flag.Bool("list", false, "List all tasks")
	flag.Parse()
	return
}

// handleFileLoading loads the todo list from a file or creates a new file if it doesn't exist
func handleFileLoading(todoList *todo.Todos, filename string) {
	if err := todoList.Load(filename); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File not found, creating a new todos.json file.")
			if err := todoList.Save(filename); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

// executeCommand executes the command based on the parsed flags
func executeCommand(add *string, complete *int, del *int, list *bool, todoList *todo.Todos) {
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
}
