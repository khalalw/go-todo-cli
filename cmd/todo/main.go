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
	Add       []string `arg:"-a,--add" help:"Add a task to the TODO list"`
	DueDate   string   `arg:"-d,--due" help:"Set a due date for the task (format: YYYY-MM-DD)"`
	Priority  string   `arg:"-p,--priority" help:"Set a priority for the task (low, medium, high)"`
	Tags      string   `arg:"-t,--tags" help:"Comma-separated tags for the task"`
	Complete  int      `arg:"-c,--complete" help:"Mark a task as complete"`
	Delete    int      `arg:"-r,--delete" help:"Delete a task"`
	List      bool     `arg:"-l,--list" help:"List all tasks"`
	Clear     bool     `arg:"-x,--clear-tasks" help:"Clear all tasks"`
	Edit      int      `arg:"-e,--edit" help:"Edit a task"`
	AddTag    []string `arg:"--add-tag" help:"Add a tag to a task: <task_number> <tag>"`
	RemoveTag []string `arg:"--remove-tag" help:"Remove a tag from a task: <task_number> <tag>"`
	FilterTag string   `arg:"--filter-tag" help:"Filter tasks by tag"`
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
		var priority todo.Priority
		var err error
		if args.Priority != "" {
			priority, err = todo.ParsePriority(args.Priority)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid priority: %s. Use low, medium, or high.\n", args.Priority)
				os.Exit(1)
			}
		} else {
			priority = todo.Low // Default priority
		}
		tags := parseTags(args.Tags)
		commands.AddCommand([]string{task}, dueDate, priority, todoList, tags)
	} else if args.Complete > 0 {
		commands.CompleteCommand([]string{fmt.Sprint(args.Complete)}, todoList)
	} else if args.Delete > 0 {
		commands.DeleteCommand([]string{fmt.Sprint(args.Delete)}, todoList)
	} else if args.List {
		commands.ListCommand(todoList)
	} else if args.Clear {
		commands.ClearTasksCommand(todoList)
	} else if args.Edit > 0 {
		commands.EditCommand(args.Edit, todoList)
	} else if len(args.AddTag) == 2 {
		commands.AddTagCommand(args.AddTag, todoList)
	} else if len(args.RemoveTag) == 2 {
		commands.RemoveTagCommand(args.RemoveTag, todoList)
	} else if args.FilterTag != "" {
		commands.FilterByTagCommand([]string{args.FilterTag}, todoList)
	} else {
		fmt.Fprintln(os.Stderr, "Invalid command. Use --help for usage information.")
	}
}

func parseTags(tagString string) []string {
	if tagString == "" {
		return nil
	}
	tags := strings.Split(tagString, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return tags
}
