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

var defaultFileToWrite = "todos.json"

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
	Search    []string `arg:"--search" help:"Search for tasks containing the given keyword"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	_, err := parseArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseArgs(args Args) (todo.Todos, error) {
	todoList := &todo.Todos{}
	err := handleFileLoading(todoList, defaultFileToWrite)
	if err != nil {
		return *todoList, err
	}

	err = executeCommand(args, todoList)
	if err != nil {
		return *todoList, err
	}

	return *todoList, todoList.Save(defaultFileToWrite)
}

func handleFileLoading(todoList *todo.Todos, filename string) error {
	if err := todoList.Load(filename); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File not found, creating a new todos.json file.")
			return todoList.Save(filename)
		}
		return fmt.Errorf("error loading go-todo-cli file: %w", err)
	}
	return nil
}

func executeCommand(args Args, todoList *todo.Todos) error {
	switch {
	case len(args.Add) > 0:
		return handleAddCommand(args, todoList)
	case args.Complete > 0:
		commands.CompleteCommand([]string{fmt.Sprint(args.Complete)}, todoList)
	case args.Delete > 0:
		commands.DeleteCommand([]string{fmt.Sprint(args.Delete)}, todoList)
	case args.List:
		commands.ListCommand(todoList)
	case args.Clear:
		commands.ClearTasksCommand(todoList)
	case args.Edit > 0:
		commands.EditCommand(args.Edit, todoList)
	case len(args.AddTag) == 2:
		commands.AddTagCommand(args.AddTag, todoList)
	case len(args.RemoveTag) == 2:
		commands.RemoveTagCommand(args.RemoveTag, todoList)
	case args.FilterTag != "":
		commands.FilterByTagCommand([]string{args.FilterTag}, todoList)
	case len(args.Search) > 0:
		commands.SearchCommand(args.Search, todoList)
	default:
		return fmt.Errorf("invalid command. Use --help for usage information")
	}
	return nil
}

func handleAddCommand(args Args, todoList *todo.Todos) error {
	task := strings.Join(args.Add, " ")
	dueDate, err := parseDueDate(args.DueDate)
	if err != nil {
		return err
	}
	priority, err := parsePriority(args.Priority)
	if err != nil {
		return err
	}
	tags := parseTags(args.Tags)
	commands.AddCommand([]string{task}, dueDate, priority, todoList, tags)
	return nil
}

func parseDueDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %s. Use YYYY-MM-DD", dateStr)
	}
	return &parsedDate, nil
}

func parsePriority(priorityStr string) (todo.Priority, error) {
	if priorityStr == "" {
		return todo.Low, nil
	}
	priority, err := todo.ParsePriority(priorityStr)
	if err != nil {
		return todo.Low, fmt.Errorf("invalid priority: %s. Use low, medium, or high", priorityStr)
	}
	return priority, nil
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
