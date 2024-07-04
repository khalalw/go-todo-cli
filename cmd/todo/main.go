package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"go-todo-cli/internal/commands"
	"go-todo-cli/internal/todo"
	"os"
	"strings"
	"time"
)

const defaultFileToWrite = "todos.json"

// Args defines the command-line arguments structure
type Args struct {
	Add      []string `arg:"-a,--add" help:"Add a task to the TODO list"`
	DueDate  string   `arg:"-d,--due" help:"Set a due date for the task (format: YYYY-MM-DD)"`
	Priority string   `arg:"-p,--priority" help:"Set a priority for the task (low, medium, high)"`
	Complete int      `arg:"-c,--complete" help:"Mark a task as complete"`
	Delete   int      `arg:"-r,--delete" help:"Delete a task"`
	List     bool     `arg:"-l,--list" help:"List all tasks"`
	Clear    bool     `arg:"-x,--clear-tasks" help:"Clear all tasks"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	todoList := &todo.Todos{}
	handleFileLoading(todoList, defaultFileToWrite)

	executeCommand(args, todoList)

	if err := todoList.Save(defaultFileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving go-todo-cli list:", err)
		os.Exit(1)
	}
}

// handleFileLoading loads the TODO list from a file, creating a new file if it doesn't exist
func handleFileLoading(todoList *todo.Todos, filename string) {
	if err := todoList.Load(filename); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File not found, creating a new todos.json file.")
			if err := todoList.Save(filename); err != nil {
				fmt.Fprintln(os.Stderr, "Error creating new go-todo-cli file:", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Error loading go-todo-cli file:", err)
			os.Exit(1)
		}
	}
}

// executeCommand executes the appropriate command based on the provided arguments
func executeCommand(args Args, todoList *todo.Todos) {
	if len(args.Add) > 0 {
		task := strings.Join(args.Add, " ")
		var dueDate *time.Time
		if args.DueDate != "" {
			parsedDate, err := time.Parse("2006-01-02", args.DueDate)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid date format: %s. Use YYYY-MM-DD.\n", args.DueDate)
				os.Exit(1)
			}
			dueDate = &parsedDate
		}
		priority, err := todo.ParsePriority(args.Priority)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid priority: %s. Use low, medium, or high.\n", args.Priority)
			os.Exit(1)
		}
		commands.AddCommand([]string{task}, dueDate, priority, todoList)
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
