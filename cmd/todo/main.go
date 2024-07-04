package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"go-todo-cli/internal/commands"
	"go-todo-cli/internal/todo"
	"os"
)

const defaultFileToWrite = "todos.json"

// Args defines the command-line arguments structure
type Args struct {
	Add      string `arg:"-a,--add,help:Add a task to the TODO list"`
	Complete int    `arg:"-c,--complete,help:Mark a task as complete"`
	Delete   int    `arg:"-d,--delete,help:Delete a task"`
	List     bool   `arg:"-l,--list,help:List all tasks"`
	Clear    bool   `arg:"-x,--clear-tasks,help:Clear all tasks"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	todoList := &todo.Todos{}
	handleFileLoading(todoList, defaultFileToWrite)

	executeCommand(args, todoList)

	if err := todoList.Save(defaultFileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving todo list:", err)
		os.Exit(1)
	}
}

// handleFileLoading loads the TODO list from a file, creating a new file if it doesn't exist
func handleFileLoading(todoList *todo.Todos, filename string) {
	if err := todoList.Load(filename); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File not found, creating a new todos.json file.")
			if err := todoList.Save(filename); err != nil {
				fmt.Fprintln(os.Stderr, "Error creating new todo file:", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Error loading todo file:", err)
			os.Exit(1)
		}
	}
}

// executeCommand executes the appropriate command based on the provided arguments
func executeCommand(args Args, todoList *todo.Todos) {
	if args.Add != "" {
		commands.AddCommand([]string{args.Add}, todoList)
	} else if args.Complete > 0 {
		commands.CompleteCommand([]string{fmt.Sprint(args.Complete)}, todoList)
	} else if args.Delete > 0 {
		commands.DeleteCommand([]string{fmt.Sprint(args.Delete)}, todoList)
	} else if args.List {
		commands.ListCommand(todoList)
	} else if args.Clear {
		commands.ClearTasksCommand(todoList)
	} else {
		fmt.Fprintln(os.Stderr, "Invalid command. Use --help for usage information.")
	}
}
