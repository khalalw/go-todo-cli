package commands

import (
	"fmt"
	"go-todo-cli/pkg/todo"
	"os"
	"strconv"
	"strings"
)

const fileToWrite = "todos.json"

// exitFunc is a variable that holds the function to exit the program. By default, it is set to os.Exit.
var exitFunc = os.Exit

func ExecuteCommand(cmd string, args []string, todoList *todo.Todos) {
	switch cmd {
	case "add":
		addCommand(args, todoList)
	case "complete":
		completeCommand(args, todoList)
	case "delete":
		deleteCommand(args, todoList)
	case "list":
		listCommand(todoList)
	case "clear-tasks":
		clearTasksCommand(todoList)
	case "exit":
		exitCommand(todoList)
	case "help":
		printHelp()
	default:
		fmt.Println("Unknown command. Type 'help' for a list of commands.")
	}
}

func addCommand(args []string, todoList *todo.Todos) {
	if len(args) < 1 {
		fmt.Println("Usage: add <task>")
		return
	}
	todoList.Add(strings.Join(args, " "))
	fmt.Println("Task added.")
	saveTodoList(todoList)
}

func completeCommand(args []string, todoList *todo.Todos) {
	if len(args) != 1 {
		fmt.Println("Usage: complete <task_number>")
		return
	}
	index := parseIndex(args[0])
	if index >= 0 && index < len(*todoList) {
		if err := todoList.Complete(index); err == nil {
			fmt.Println("Task marked as complete.")
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Invalid task number.")
	}
	saveTodoList(todoList)
}

func deleteCommand(args []string, todoList *todo.Todos) {
	if len(args) != 1 {
		fmt.Println("Usage: delete <task_number>")
		return
	}
	index := parseIndex(args[0])
	if index >= 0 && index < len(*todoList) {
		if err := todoList.Delete(index); err == nil {
			fmt.Println("Task deleted.")
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Invalid task number.")
	}
	saveTodoList(todoList)
}

func listCommand(todoList *todo.Todos) {
	todo.Print(todoList)
}

func clearTasksCommand(todoList *todo.Todos) {
	*todoList = todo.Todos{}
	fmt.Println("All tasks cleared.")
	saveTodoList(todoList)
}

func exitCommand(todoList *todo.Todos) {
	fmt.Println("Exiting the TODO CLI.")
	saveTodoList(todoList)
	exitFunc(0)
}

func saveTodoList(todoList *todo.Todos) {
	if err := todoList.Save(fileToWrite); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving todo list:", err)
	}
}

func parseIndex(input string) int {
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || index <= 0 {
		return -1 // Return -1 for invalid or non-positive inputs
	}
	return index - 1 // Convert to zero-based index
}

func printHelp() {
	fmt.Println("Commands:")
	fmt.Println("  add <task>           - Add a new task")
	fmt.Println("  complete <task_number> - Mark a task as complete")
	fmt.Println("  delete <task_number> - Delete a task")
	fmt.Println("  list                 - List all tasks")
	fmt.Println("  clear-tasks          - Clear all tasks")
	fmt.Println("  help                 - Show this help message")
	fmt.Println("  exit                 - Exit the TODO CLI")
}
